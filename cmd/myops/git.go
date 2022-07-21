package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

const clonePath string = "cloneTemp"

const shortHashSize int = 7

func remoteShorthash(repoUrl string, branch string) string {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoUrl},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		// fmt.Println("Couldn't get remote refs")
		// panic(err)
		return ""
	}

	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name().Short() == branch {
			return ref.Hash().String()[:shortHashSize]
		}
	}

	return ""
}

func clone(repoUrl string, branch string) *bytes.Buffer {
	fmt.Println("cloning", repoUrl+"#"+branch)
	git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:           repoUrl,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})

	defer os.RemoveAll(clonePath)

	ret := &bytes.Buffer{}

	gzw := gzip.NewWriter(ret)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	filepath.Walk(clonePath, func(file string, fi fs.FileInfo, err error) error {
		if err != nil {
			return (err)
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return (err)
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, clonePath, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return (err)
		}

		f, err := os.Open(file)
		if err != nil {
			return (err)
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})

	return ret
}

package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

const shortHashSize int = 7

func remoteShorthash(repoUrl string, branch string) string {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoUrl},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		fmt.Println("Couldn't get remote refs for", repoUrl, ", tagging with 'no_hash'")
		return "no_hash"
	}

	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name().Short() == branch {
			return ref.Hash().String()[:shortHashSize]
		}
	}

	return ""
}

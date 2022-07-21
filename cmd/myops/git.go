package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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

// func clone(repoUrl string, branch string) {
// 	git.PlainClone(clonePath, false, &git.CloneOptions{
// 		URL:           repoUrl,
// 		ReferenceName: plumbing.NewBranchReferenceName(branch),
// 		SingleBranch:  true,
// 	})
// }

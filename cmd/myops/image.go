package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	cli "github.com/docker/docker/client"
)

func buildImage(ctx context.Context, client *cli.Client, config Config, projectTag string) {
	resp, err := client.ImageBuild(ctx, nil, types.ImageBuildOptions{
		RemoteContext: config.RepoUrl + "/archive/refs/heads/" + config.Branch + ".tar.gz",
		Tags:          []string{projectTag},
		Dockerfile:    config.Dockerfile,
	})
	if err != nil {
		fmt.Println("Couldn't build image for", projectTag, err)
		fmt.Println("Resp:", resp)
		return
	}

	defer resp.Body.Close()

	var bytes []byte
	resp.Body.Read(bytes)

	fmt.Println("Response for building image for", projectTag, ":", string(bytes))
}

package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	cli "github.com/docker/docker/client"
)

func buildImage(ctx context.Context, client *cli.Client, config Config, projectTag string) {
	resp, err := client.ImageBuild(ctx, nil, types.ImageBuildOptions{
		RemoteContext: config.RepoUrl + "#" + config.Branch,
		Tags:          []string{projectTag},
		Dockerfile:    config.Dockerfile,
		Remove:        true,
	})
	if err != nil {
		fmt.Println("Couldn't build image for", projectTag, err)
		fmt.Println("Resp:", resp)
		return
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println("Building", projectTag, scanner.Text())
	}
}

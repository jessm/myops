package main

import (
	"context"
	"fmt"

	cli "github.com/docker/docker/client"
)

func update() {
	fmt.Println("--- MyOps Update ---")

	ctx := context.Background()
	client, err := cli.NewClientWithOpts(cli.FromEnv, cli.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	configs := parseConfig()

	cleanupContainers(ctx, client, configs)
	cleanupImages(ctx, client, configs)
	cleanupVolumes(ctx, client, configs)

	for projectName, config := range configs {
		shortHash := remoteShorthash(config.RepoUrl, config.Branch)
		projectTag := projectName + ":" + shortHash

		buildImage(ctx, client, config, projectName)
		err := runContainer(ctx, client, config, projectName, projectTag)
		if err != nil {
			panic(err)
		}
	}

	renderCaddyfile(configs)

	runCaddy()
}

func main() {
	update()
}

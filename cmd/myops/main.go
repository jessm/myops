package main

import (
	"context"
	"fmt"
	"os"

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
		// shortHash := remoteShorthash(config.RepoUrl, config.Branch)
		// projectTag := projectName + ":" + shortHash

		fmt.Println("cloning", projectName)
		tarBuffer := clone(config.RepoUrl, config.Branch)
		file, err := os.Create(clonePath)
		if err != nil {
			panic(err)
		}
		file.Write(tarBuffer.Bytes())

		// buildImage(ctx, client, config, projectName)
		// err := runContainer(ctx, client, config, projectName, projectTag)
		// if err != nil {
		// 	panic(err)
		// }
	}

	renderCaddyfile(configs)

	runCaddy()
}

func main() {
	update()
}

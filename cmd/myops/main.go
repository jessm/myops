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

	fmt.Println("Project short hashes:")
	for projectName, config := range configs {
		shortHash := remoteShorthash(config.RepoUrl, config.Branch)
		fmt.Printf("  - %s: %s\n", projectName, shortHash)
	}

	renderCaddyfile(configs)

	runCaddy()
}

func main() {
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	panic(err)
	// }

	// imagename := "nginxdemos/hello"
	// containername := "demo"
	// portopening := "80"
	// inputEnv := []string{}
	// err = runContainer(cli, imagename, containername, portopening, inputEnv)
	// if err != nil {
	// 	log.Println(err)
	// }

	update()
}

package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	cli "github.com/docker/docker/client"
)

func update() {
	fmt.Println("--- MyOps Update ---")

	ctx := context.Background()
	client, err := cli.NewClientWithOpts(cli.FromEnv, cli.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	configs := getConfigs()
	oldConfigs := getOldConfigs()

	cleanupContainers(ctx, client, configs)
	cleanupImages(ctx, client, configs)
	cleanupVolumes(ctx, client, configs)

	projectsUpdated := []string{}

	projectIPs := map[string]string{}

	for projectName, config := range configs {
		shortHash := remoteShorthash(config.RepoUrl, config.Branch)
		projectTag := projectName + ":" + shortHash

		oldConfig, oldConfigExists := oldConfigs[projectName]
		// If no old config exists, we count that as a config change
		configChanged := !oldConfigExists || checkConfigChanged(config, oldConfig, projectName)
		shortHashChanged := checkShortHashChanged(images, projectName, shortHash)

		// If config or shorthash changed, we want to rerun the container, so stop the old one
		if configChanged || shortHashChanged {
			removeContainerByProject(ctx, client, projectName)
		}

		// If shorthash changed, also remove the image and rebuild it
		if shortHashChanged || !oldConfigExists {
			removeImageByProject(ctx, client, projectName)
			buildImage(ctx, client, config, projectTag)
		}

		// Finally if either changed, rerun the container
		if configChanged || shortHashChanged {
			_, err := runContainer(ctx, client, config, config.HostPort, projectTag, projectName)
			if err != nil {
				panic(err)
			}

			projectsUpdated = append(projectsUpdated, projectName)
		}
	}

	for projectName := range configs {
		ip := containerByProject(ctx, client, projectName).NetworkSettings.Networks["bridge"].IPAddress
		projectIPs[projectName] = ip
	}

	renderCaddyfile(configs, projectIPs)
	fmt.Println("Caddyfile created:")
	printCaddyfile()
	fmt.Println("End of Caddyfile")

	runCaddy(ctx, client)

	fmt.Println("Finished updating, writing to old config file")
	writeConfigToOldConfig()

	fmt.Println("Projects updated:")
	for _, name := range projectsUpdated {
		fmt.Println("  - " + name)
	}
}

func main() {
	update()
}

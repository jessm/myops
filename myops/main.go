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
		if !oldConfigExists || newVersionRequired(images, config, oldConfig, projectName, shortHash) {
			removeContainerByProject(ctx, client, projectName)
			removeImageByProject(ctx, client, projectName)

			buildImage(ctx, client, config, projectTag)

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

	runCaddy()

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

package main

import (
	"context"
	"fmt"
	"strconv"

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

	portNum := 8001
	portMap := map[string]string{}

	projectsUpdated := []string{}

	for projectName, config := range configs {
		shortHash := remoteShorthash(config.RepoUrl, config.Branch)
		projectTag := projectName + ":" + shortHash

		oldConfig, oldConfigExists := oldConfigs[projectName]
		if !oldConfigExists || newVersionRequired(images, config, oldConfig, projectName, shortHash) {
			removeContainerByProject(ctx, client, projectName)
			removeImageByProject(ctx, client, projectName)

			buildImage(ctx, client, config, projectTag)

			err := runContainer(ctx, client, config, strconv.Itoa(portNum), projectTag, projectName)
			if err != nil {
				panic(err)
			}

			projectsUpdated = append(projectsUpdated, projectName)
		}

		portMap[projectName] = strconv.Itoa(portNum)

		portNum++
	}

	renderCaddyfile(configs, portMap)
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	cli "github.com/docker/docker/client"
)

func removeImage(ctx context.Context, client *cli.Client, id string) {
	_, err := client.ImageRemove(ctx, id, types.ImageRemoveOptions{})
	if err != nil {
		fmt.Println("  Couldn't remove image", id)
		fmt.Println("  Error:", err)
	}
}

// Just do best effort image cleanup to save some space
func cleanupImages(ctx context.Context, client *cli.Client, configs Configs) {
	images, err := client.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: "myops"}),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Images for cleaning:")
	for _, i := range images {
		if len(i.RepoTags) > 1 {
			fmt.Println("Found image with multiple tags:", i.RepoTags, ", skipping")
			continue
		}
		if len(i.RepoTags) < 1 {
			fmt.Println("Found image with less than 1 tag:", i.ID, ", skipping")
			fmt.Println(i.RepoDigests, i.RepoTags)
			continue
		}
		for _, tag := range i.RepoTags {
			// Skip the caddy image and myops image
			if tag == CaddyImage || strings.Split(tag, ":")[0] == "myops" {
				break
			}
			// Skip if it's in the project list, it'll get cleaned later if it's outdated
			if _, ok := configs[strings.Split(tag, ":")[0]]; ok {
				break
			}
			fmt.Println("  - " + tag + " - " + i.ID)
			removeImage(ctx, client, i.ID)
		}
	}
}

func removeContainer(ctx context.Context, client *cli.Client, id string) {
	err := client.ContainerStop(ctx, id, nil)
	if err != nil {
		fmt.Println("  Couldn't stop container", id)
	}
	err = client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	if err != nil {
		fmt.Println("  Couldn't remove container", id)
		fmt.Println("  Error:", err)
	}
}

func cleanupContainers(ctx context.Context, client *cli.Client, configs Configs) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("error getting hostname:", err)
	}
	fmt.Println("Hostname:", hostname)

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: "myops"}),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Containers for cleaning:")
	for _, c := range containers {
		if len(c.Names) > 1 {
			fmt.Println("Found container with multiple names:", c.Names, ", skipping")
			continue
		}
		// Skip the myops container
		if hostname == c.ID[:len(hostname)] {
			continue
		}
		// Skip if it's in the project list
		if _, ok := configs[strings.Split(c.Image, ":")[0]]; ok {
			continue
		}
		for _, name := range c.Names {
			// Skip caddy
			if name == "/"+CaddyContainer {
				continue
			}
			fmt.Println("  - " + name + ": " + c.ID)
			removeContainer(ctx, client, c.ID)
		}
	}
}

func removeVolume(ctx context.Context, client *cli.Client, name string) {
	err := client.VolumeRemove(ctx, name, false)
	if err != nil {
		fmt.Println("  Couldn't remove volume", name)
		fmt.Println("  Error:", err)
	}
}

func cleanupVolumes(ctx context.Context, client *cli.Client, configs Configs) {
	resp, err := client.VolumeList(ctx, filters.NewArgs(filters.KeyValuePair{Key: "label", Value: "myops"}))
	if err != nil {
		panic(err)
	}

	volumes := resp.Volumes

	fmt.Println("Volumes for cleaning:")
	for _, v := range volumes {
		// Skip caddy volumes
		if v.Name == CaddyDataVolumeName || v.Name == CaddyConfigVolumeName {
			continue
		}
		// Skip if it's in the project list
		if _, ok := configs[v.Name]; ok {
			continue
		}
		// Note: bind mount volumes aren't shown, so we're safe for those
		fmt.Println("  - " + v.Name)
		removeVolume(ctx, client, v.Name)
	}
}

func removeContainerByProject(ctx context.Context, client *cli.Client, projectName string) {
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		if strings.Split(c.Image, ":")[0] == projectName {
			removeContainer(ctx, client, c.ID)
		}
	}
}

func removeImageByProject(ctx context.Context, client *cli.Client, projectName string) {
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, i := range images {
		for _, tag := range i.RepoTags {
			if strings.Split(tag, ":")[0] == projectName {
				removeImage(ctx, client, i.ID)
				return
			}
		}
	}
}

func checkShortHashChanged(images []types.ImageSummary, projectName string, shortHash string) bool {
	// Check if repo updated
	for _, i := range images {
		for _, tag := range i.RepoTags {
			tagSplit := strings.Split(tag, ":")
			if len(tagSplit) != 2 {
				continue
			}

			if tagSplit[0] == projectName && tagSplit[1] != shortHash {
				return true
			}
		}
	}

	return false
}

func checkConfigChanged(config Config, oldConfig Config, projectName string) bool {
	// Check if config changed
	jsonConfig, _ := json.Marshal(config)
	oldJsonConfig, _ := json.Marshal(oldConfig)
	return string(jsonConfig) != string(oldJsonConfig)
}

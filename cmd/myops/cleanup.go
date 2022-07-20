package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	cli "github.com/docker/docker/client"
)

func removeImage(id string) {

}

func cleanupImages(ctx context.Context, client *cli.Client, configs Configs) {
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Images for cleaning:")
	for _, i := range images {
		for _, repo := range i.RepoDigests {
			fmt.Println("  - " + repo + " - " + i.ID)
		}
	}
}

func removeContainer(id string) {

}

func cleanupContainers(ctx context.Context, client *cli.Client, configs Configs) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("error getting hostname:", err)
	}
	fmt.Println("Hostname:", hostname)

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
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
		for _, name := range c.Names {
			// Skip caddy
			if name == CaddyContainer {
				continue
			}
			fmt.Println("  - " + name + ": " + c.ID)
		}
	}
}

func removeVolume(id string) {

}

func cleanupVolumes(ctx context.Context, client *cli.Client, configs Configs) {
	resp, err := client.VolumeList(ctx, filters.Args{})
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
		// Note: bind mount volumes aren't shown, so we're safe for those
		fmt.Println("  - " + v.Name)
	}
}

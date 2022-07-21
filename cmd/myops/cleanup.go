package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	cli "github.com/docker/docker/client"
)

func removeImage(ctx context.Context, client *cli.Client, id string) {
	resp, err := client.ImageRemove(ctx, id, types.ImageRemoveOptions{})
	if err != nil {
		fmt.Println("  Couldn't remove image", id)
		fmt.Println("  Error:", err)
		fmt.Println("  Response:", resp)
	}
}

func cleanupImages(ctx context.Context, client *cli.Client, configs Configs) {
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Images for cleaning:")
	for _, i := range images {
		if len(i.RepoDigests) > 1 {
			fmt.Println("Found container with multiple repos:", i.RepoDigests, ", skipping")
			continue
		}
		for _, repo := range i.RepoDigests {
			// Skip the caddy image
			if strings.Split(repo, "@")[0] == strings.Split(CaddyImage, ":")[0] {
				continue
			}
			fmt.Println("  - " + repo + " - " + i.ID)
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
		removeVolume(ctx, client, v.Name)
	}
}

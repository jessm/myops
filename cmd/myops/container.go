package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	cli "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func runContainer(ctx context.Context, client *cli.Client, config Config, imageName string, containerName string) error {
	containerPort, err := nat.NewPort("tcp", config.Port)
	if err != nil {
		fmt.Println("Unable to create port")
		return err
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: config.Port,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	exposedPorts := map[nat.Port]struct{}{
		containerPort: {},
	}

	envList := make([]string, 0, len(config.EnvVars))
	for key, value := range config.EnvVars {
		envList = append(envList, key+"="+value)
	}

	containerConfig := &container.Config{
		Image:        imageName,
		Env:          envList,
		ExposedPorts: exposedPorts,
	}

	container, err := client.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})

	return nil
}

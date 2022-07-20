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

func runContainer(imagename string, containername string, port string, envVars map[string]string) error {
	ctx := context.Background()
	client, err := cli.NewClientWithOpts(cli.FromEnv, cli.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containerPort, err := nat.NewPort("tcp", port)
	if err != nil {
		fmt.Println("Unable to create port")
		return err
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
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

	envList := make([]string, 0, len(envVars))
	for key, value := range envVars {
		envList = append(envList, key+"="+value)
	}

	containerConfig := &container.Config{
		Image:        imagename,
		Env:          envList,
		ExposedPorts: exposedPorts,
	}

	container, err := client.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		containername,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})

	return nil
}

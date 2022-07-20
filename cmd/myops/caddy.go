package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	cli "github.com/docker/docker/client"
)

type templateConfigs struct {
	C Configs
}

const (
	Caddyfile      string = "/var/myops/Caddyfile"
	CaddyfileMount string = "/etc/caddy/Caddyfile"
)

const (
	CaddyDataVolumeName    string = "CaddyData"
	CaddyConfigVolumeName  string = "CaddyConfig"
	CaddyDataVolumeMount   string = "/data"
	CaddyConfigVolumeMount string = "/config"
)

const CaddyImage string = "caddy:2.5.2"

const caddyTemplate string = `{{ range .C }}
{{ .DomainMatcher }} {
	reverse_proxy localhost:{{ .Port }}
}
{{ end }}`

func printCaddyfile() {
	bytes, err := ioutil.ReadFile(Caddyfile)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func renderCaddyfile(configs Configs) {
	t, err := template.New("Caddyfile").Parse(caddyTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(Caddyfile)
	if err != nil {
		panic(err)
	}

	err = t.Execute(file, templateConfigs{configs})
	if err != nil {
		panic(err)
	}
}

func runCaddy() {
	ctx := context.Background()
	client, err := cli.NewClientWithOpts(cli.FromEnv, cli.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// If the volumes don't exist yet, create them
	volumes, err := client.VolumeList(ctx, filters.Args{})
	if err != nil {
		panic(err)
	}

	dataVolFound := false
	configVolFound := false
	for _, vol := range volumes.Volumes {
		if vol.Name == CaddyDataVolumeName {
			dataVolFound = true
		}
		if vol.Name == CaddyConfigVolumeName {
			configVolFound = true
		}
	}

	if dataVolFound != configVolFound {
		panic("Found only one of config or data volume for caddy")
	}

	if !dataVolFound || !configVolFound {
		_, err := client.VolumeCreate(ctx, volume.VolumeCreateBody{
			Name: CaddyDataVolumeName,
		})
		if err != nil {
			panic(err)
		}

		_, err = client.VolumeCreate(ctx, volume.VolumeCreateBody{
			Name: CaddyConfigVolumeName,
		})
		if err != nil {
			panic(err)
		}
	}

	// Set up and run the caddy container
	containerConfig := &container.Config{
		Image: CaddyImage,
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "host",
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Binds: []string{
			Caddyfile + ":" + CaddyfileMount,
			CaddyDataVolumeName + ":" + CaddyDataVolumeMount,
			CaddyConfigVolumeName + ":" + CaddyConfigVolumeMount,
		},
	}

	container, err := client.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		"caddy",
	)

	if err != nil {
		panic(err)
	}

	client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
}

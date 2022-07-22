package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	cli "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type templateConfigs struct {
	C   Configs
	IPs map[string]string
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
const CaddyContainer string = "caddy"

const caddyTemplate string = `{{ range $name, $config := .C }}
{{ $config.DomainMatcher }} {
	reverse_proxy {{ index $.IPs $name }}:{{ $config.Port }}
}
{{ end }}`

func printCaddyfile() {
	bytes, err := ioutil.ReadFile(Caddyfile)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func renderCaddyfile(configs Configs, projectIPs map[string]string) {
	t, err := template.New("Caddyfile").Parse(caddyTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(Caddyfile)
	if err != nil {
		panic(err)
	}

	err = t.Execute(file, templateConfigs{
		C:   configs,
		IPs: projectIPs,
	})
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

	// Check if we're already running
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+CaddyContainer {
				if c.State == "running" {
					// Reload the configuration
					resp, err := client.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
						WorkingDir: "/etc/caddy",
						Cmd:        []string{"caddy", "reload"},
					})
					if err != nil {
						panic(err)
					}

					if err := client.ContainerExecStart(ctx, resp.ID, types.ExecStartCheck{}); err != nil {
						panic(err)
					}

					return
				} else {
					// If state isn't running, just remove the container
					removeContainer(ctx, client, c.ID)
				}
			}
		}
	}

	// If we haven't pulled the image yet, pull it
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	caddyImageFound := false
	for _, i := range images {
		for _, repo := range i.RepoDigests {
			if strings.Split(repo, "@")[0] == CaddyContainer {
				caddyImageFound = true
				break
			}
		}
	}

	if !caddyImageFound {
		out, err := client.ImagePull(ctx, CaddyImage, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		defer out.Close()
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

	httpPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		fmt.Println("can't get port 80 for caddy")
		panic(err)
	}

	httpsPort, err := nat.NewPort("tcp", "443")
	if err != nil {
		fmt.Println("can't get port 443 for caddy")
		panic(err)
	}

	portMap := nat.PortMap{
		httpPort:  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "80"}},
		httpsPort: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "443"}},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portMap,
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
		CaddyContainer,
	)

	if err != nil {
		panic(err)
	}

	client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
}

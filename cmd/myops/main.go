package main

import (
	"fmt"
)

func update() {
	configs := parseConfig()

	for projectName, config := range configs {
		shortHash := remoteShorthash(config.RepoUrl, config.Branch)
		fmt.Printf("%s: %s\n", projectName, shortHash)
	}

	renderCaddyfile(configs)
	fmt.Println("Caddyfile contents:")
	printCaddyfile()

	runCaddy()
}

func main() {
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	panic(err)
	// }

	// imagename := "nginxdemos/hello"
	// containername := "demo"
	// portopening := "80"
	// inputEnv := []string{}
	// err = runContainer(cli, imagename, containername, portopening, inputEnv)
	// if err != nil {
	// 	log.Println(err)
	// }

	update()
}

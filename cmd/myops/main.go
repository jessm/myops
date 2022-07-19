package main

import (
	"errors"
	"fmt"
	"os"
)

func update() {
	_, err := os.Stat(configFile)
	if errors.Is(err, os.ErrNotExist) {
		writeSampleConfig()
	}

	configs := parseConfig()

	fmt.Println(configs)
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

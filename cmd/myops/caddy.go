package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

const Caddyfile string = "/Caddyfile"

type caddyEntry struct {
	matcher string
	port    string
}

const caddyTemplate string = `
{{ range entries }}
{{ .matcher }} {
	reverse_proxy localhost:{{ .port }}
}
{{ end }}
`

func printCaddyfile() {
	bytes, err := ioutil.ReadFile(Caddyfile)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func renderCaddyfile(configs Configs) {
	entries := []caddyEntry{}
	for _, c := range configs {
		entries = append(entries, caddyEntry{
			matcher: c.DomainMatcher,
			port:    c.Port,
		})
	}

	t, err := template.New("Caddyfile").Parse(caddyTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(Caddyfile)
	if err != nil {
		panic(err)
	}

	err = t.Execute(file, entries)
	if err != nil {
		panic(err)
	}
}

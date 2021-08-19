package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Services []struct {
		Name    string `yaml:"name"`
		Host    string `yaml:"host"`
		Sockets []struct {
			Port int    `yaml:"port"`
			App  string `yaml:"app"`
		} `yaml:"sockets"`
	} `yaml:"services"`

	Apps []struct {
		Name   string `yaml:"name"`
		Routes []struct {
			Path     string `yaml:"path"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"routes"`
	} `yaml:"apps"`
}

type FastCfg map[int]FastPort
type FastPort map[string]FastHost
type FastHost []FastRoute
type FastRoute struct {
	Path     string
	Endpoint string
}

func (config Config) toFastCfg() (fastCfg FastCfg) {
	fastCfg = FastCfg{}

	for _, service := range config.Services {

		for _, socket := range service.Sockets {

			if fastCfg[socket.Port] == nil {
				fastCfg[socket.Port] = FastPort{}
			}
			if fastCfg[socket.Port][service.Host] == nil {
				fastCfg[socket.Port][service.Host] = FastHost{}
			}

			// Passende App finden
			for _, app := range config.Apps {

				if app.Name == socket.App {
					sort.Slice(app.Routes, func(i, j int) bool {
						if app.Routes[i].Path == "/" {
							return false
						}
						if app.Routes[j].Path == "/" {
							return true
						}

						return strings.Count(app.Routes[i].Path, "/") < strings.Count(app.Routes[j].Path, "/")
					})

					// Routen in FastCfg einfügen
					for _, route := range app.Routes {

						if len(route.Path) > 1 {
							route.Path = strings.TrimSuffix(route.Path, "/")
						}
						route.Endpoint = strings.TrimSuffix(route.Endpoint, "/")

						fastCfg[socket.Port][service.Host] = append(fastCfg[socket.Port][service.Host], FastRoute{route.Path, route.Endpoint})
					}
				}
			}
		}
	}

	return
}

func loadConfig() (config Config) {
	if len(os.Args[1:]) == 0 {
		fmt.Println("Keine Konfiguration angegeben")
		os.Exit(1)
	}

	filePath, _ := filepath.Abs(os.Args[1:][0])
	yamlFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Konfiguration konnte nicht geöffnet werden")
		os.Exit(1)
	}

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		fmt.Println("Konfiguration ist fehlerhaft")
		os.Exit(1)
	}

	return
}

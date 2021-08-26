package main

type Config struct {
	Services []struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Host        string   `yaml:"host"`
		Hosts       []string `yaml:"hosts"`
		Port        int      `yaml:"port"`
		App         string   `yaml:"app"`
		DevPort     int      `yaml:"dev_port"`
	} `yaml:"services"`

	Apps []struct {
		Name   string `yaml:"name"`
		Routes []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Path        string `yaml:"path"`
			Endpoint    string `yaml:"endpoint"`
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

type DevCfg map[int]DevPort
type DevPort []*DevService

type DevServerCfg map[int]DevServerPort
type DevServerPort []*DevService

type DevService struct {
	Name        string
	Description string
	Host        string
	Port        int
	DevPort     int
	Routes      []DevRoute
}
type DevRoute struct {
	Name            string
	Description     string
	Path            string
	Endpoint        string
	RecentEndpoints []string
}

type RequestBody struct {
	Route    int
	Endpoint string
}

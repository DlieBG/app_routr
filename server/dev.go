package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (devPort DevPort) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, devService := range devPort {

		if devService.Host == strings.Split(req.Host, ":")[0] {

			for _, route := range devService.Routes {

				if strings.HasPrefix(req.URL.Path, route.Path) {
					req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Path)
					if !strings.HasPrefix(req.URL.Path, "/") {
						req.URL.Path = "/" + req.URL.Path
					}

					url, _ := url.Parse(route.Endpoint)

					proxy := httputil.NewSingleHostReverseProxy(url)
					proxy.ServeHTTP(res, req)

					return
				}
			}

			fmt.Println("404 - Keine Route gedunden")
			res.WriteHeader(404)
			res.Write([]byte("404 - Keine Route gefunden"))
			return
		}
	}

	fmt.Println("404 - Keinen Host gefunden")
	res.WriteHeader(404)
	res.Write([]byte("404 - Keinen Host gefunden"))
}

func (devServerPort DevServerPort) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, devService := range devServerPort {
		if devService.Host == strings.Split(req.Host, ":")[0] {
			if strings.HasPrefix(req.URL.Path, "/api") {
				devService.ServeAPI(res, req)
			} else {
				if _, err := os.Stat("./frontend/dist/frontend"); !os.IsNotExist(err) {
					http.FileServer(http.Dir("./frontend/dist/frontend")).ServeHTTP(res, req)
				}
				if _, err := os.Stat("../frontend/dist/frontend"); !os.IsNotExist(err) {
					http.FileServer(http.Dir("../frontend/dist/frontend")).ServeHTTP(res, req)
				}
			}
		}
	}
}

func (devService DevService) ServeAPI(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		serviceJson, _ := json.Marshal(devService)

		res.Header().Add("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Write(serviceJson)
	}

	if req.Method == "POST" {
		var body RequestBody
		json.NewDecoder(req.Body).Decode(&body)

		if !includes(devService.Routes[body.Route].RecentEndpoints, body.Endpoint) {
			devService.Routes[body.Route].RecentEndpoints = append(devService.Routes[body.Route].RecentEndpoints, body.Endpoint)
		}

		devService.Routes[body.Route].Endpoint = body.Endpoint

		serviceJson, _ := json.Marshal(devService.Routes[body.Route])

		res.Header().Add("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Write(serviceJson)
	}

	if req.Method == "OPTIONS" {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "*")
		res.Header().Set("Allow", "OPTIONS, GET, POST")
		res.Write([]byte(""))
	}
}

func includes(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func scheduleDev(config Config, waiter *sync.WaitGroup) {
	devCfg, devServerCfg := config.toDevCfgAndDevServerCfg()

	fmt.Println("Listening Dev:")

	for port, devPort := range devCfg {
		fmt.Println(port, devPort)
		waiter.Add(1)
		go func(port int, devPort DevPort) {
			http.ListenAndServe(":"+strconv.FormatInt(int64(port), 10), devPort)
		}(port, devPort)
	}

	for port, devServerPort := range devServerCfg {
		fmt.Println(port, devServerPort)
		waiter.Add(1)
		go func(port int, devServerPort DevServerPort) {
			http.ListenAndServe(":"+strconv.FormatInt(int64(port), 10), devServerPort)
		}(port, devServerPort)
	}

	fmt.Println()
}

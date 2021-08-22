package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func (fastPort FastPort) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for hostname, routes := range fastPort {

		if hostname == strings.Split(req.Host, ":")[0] {

			for _, route := range routes {

				if strings.HasPrefix(req.URL.Path, route.Path) {
					req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Path)
					if !strings.HasPrefix(req.URL.Path, "/") {
						req.URL.Path = "/" + req.URL.Path
					}

					url, _ := url.Parse(route.Endpoint)

					proxy := httputil.NewSingleHostReverseProxy(url)
					proxy.ServeHTTP(res, req)

					// fmt.Println(req.Host + req.URL.Path)

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

func (devService DevService) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

				// fmt.Println(req.Host + req.URL.Path)

				return
			}
		}

		fmt.Println("404 - Keine Route gedunden")
		res.WriteHeader(404)
		res.Write([]byte("404 - Keine Route gefunden"))
		return
	}

	fmt.Println("404 - Keinen Host gefunden")
	res.WriteHeader(404)
	res.Write([]byte("404 - Keinen Host gefunden"))
}

func main() {
	fmt.Println("app-routr")
	fmt.Println()

	waiter := sync.WaitGroup{}
	config := loadConfig()
	fastCfg := config.toFastCfg()
	devCfg := config.toDevCfg()

	fmt.Println("Listening:")

	for portNumber, port := range fastCfg {
		fmt.Println(portNumber, port)
		waiter.Add(1)
		go func(portNumber int, port FastPort) {
			http.ListenAndServe(":"+strconv.FormatInt(int64(portNumber), 10), port)
		}(portNumber, port)
	}

	fmt.Println()
	fmt.Println("Listening Dev:")

	for _, devService := range devCfg {
		fmt.Println(devService)
		waiter.Add(1)
		go func(devService DevService) {
			registerDev(&devService)
			http.ListenAndServe(devService.Host+":"+strconv.FormatInt(int64(devService.Port), 10), &devService)
		}(devService)
	}

	fmt.Println()

	waiter.Wait()
}

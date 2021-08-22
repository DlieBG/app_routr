package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type DevServer struct {
	DevService *DevService
}

type RequestBody struct {
	Route    int
	Endpoint string
}

func (devServer DevServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api") {
		devServer.ServeAPI(res, req)
	} else {
		http.FileServer(http.Dir("./frontend/dist/frontend")).ServeHTTP(res, req)
	}
}

func (devServer DevServer) ServeAPI(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		serviceJson, _ := json.Marshal(devServer.DevService)

		res.Header().Add("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Write(serviceJson)
	}

	if req.Method == "POST" {
		var body RequestBody
		json.NewDecoder(req.Body).Decode(&body)

		if !includes(devServer.DevService.Routes[body.Route].RecentEndpoints, body.Endpoint) {
			devServer.DevService.Routes[body.Route].RecentEndpoints = append(devServer.DevService.Routes[body.Route].RecentEndpoints, body.Endpoint)
		}

		devServer.DevService.Routes[body.Route].Endpoint = body.Endpoint

		serviceJson, _ := json.Marshal(devServer.DevService.Routes[body.Route])

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

func registerDev(devService *DevService) {
	go func(devServer DevServer) {
		http.ListenAndServe(devServer.DevService.Host+":"+strconv.FormatInt(int64(devServer.DevService.DevPort), 10), devServer)
	}(DevServer{devService})
}

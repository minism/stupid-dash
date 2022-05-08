package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client
var httpClient *http.Client

type ContainerData struct {
	Name string
	Url  string
}

type TemplateContext struct {
	Title      string
	Containers []ContainerData
}

func init() {
	// Init docker client.
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Init http client.
	httpClient = &http.Client{
		Timeout: time.Second * 5,
	}
}

func getTemplateData(hostname string, title string) TemplateContext {
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	containerDatas := make([]ContainerData, 0)
	for _, container := range containers {
		for _, port := range container.Ports {
			// Check if port is HTTP.
			if isPortHttp(container, port) {
				containerDatas = append(containerDatas, ContainerData{
					Name: container.Names[0][1:],
					Url:  fmt.Sprintf("%v:%v", hostname, port.PublicPort),
				})
				break
			}
		}
	}

	return TemplateContext{
		Title:      title,
		Containers: containerDatas,
	}
}

func isPortHttp(container types.Container, port types.Port) bool {
	// Try each container IP.
	for _, network := range container.NetworkSettings.Networks {
		if tryHttpConnection(network.IPAddress, port.PrivatePort) {
			return true
		}
	}
	return false
}

func tryHttpConnection(ip string, port uint16) bool {
	url := fmt.Sprintf("http://%v:%v/", ip, port)
	log.Printf("Trying container at %v", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-StupidDash", "1")

	_, err = httpClient.Do(req)
	return err == nil
}

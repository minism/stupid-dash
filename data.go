package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

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
			if isPortHttp(port.PublicPort) {
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

func isPortHttp(port uint16) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%v", port), time.Duration(10*time.Second))
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false
	}

	match, _ := regexp.MatchString("^HTTP", response)
	return match
}

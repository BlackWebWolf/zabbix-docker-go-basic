package main

//import "time"
import (
	. "zabbix_sender/zabbix_lib"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"fmt"
	"context"
	// "encoding/json"
	// "time"
)
// type zbdata struct {
//     data   int      `json:"data"`
// }
//
// type container struct {
// 	name string `json:{#CONTAINER}`
// }
const (
	defaultHost  = `localhost`
	defaultPort  = 10051
)

func main() {
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All:true})
	if err != nil {
		panic(err)
	}
			var metrics []*Metric
	// test := `{"data":[{"{#CONTAINER}":"`+ container.Image+`"}]}`
	test1 := `{"data":[`
	for i, container := range containers {
		fmt.Printf("%s %s\n", container.State, container.Image)

		var app string
		if i < len(containers)-1 {
			app = `{"{#CONTAINER}":"` + container.Image+`"},`
		} else {
			app = `{"{#CONTAINER}":"` + container.Image+`"}`
		}
		test1 = test1 + app



		metrics = append(metrics, NewMetric("ZBagent", "docker.["+container.Image+"]", container.State))

		//// Create instance of Packet class

	}
	test1 = test1 + `]}`
			metrics = append(metrics, NewMetric("ZBagent", "docker.discovery", test1))

			packet := NewPacket(metrics)
			//
			//// Send packet to zabbix
			z := NewSender(defaultHost, defaultPort)
			z.Send(packet)
}

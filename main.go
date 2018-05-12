package main

//import "time"
import (
	. "zabbix_sender/zabbix_lib"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"fmt"
	"context"
	"flag"
	// "reflect"
	// "encoding/json"
	// "time"
)
// const (
// 	defaultHost  = `localhost`
// 	defaultPort  = 10051
// )

func main() {


	defaultHost := flag.String("zabbix","localhost","zabbix server e.g. '127.0.0.1'")
targetHost := flag.String("host","ZBagent","zabbix target host e.g. 'myhost'")
defaultPort := flag.Int("port",10051,"zabbix server port e.g. 10051")
flag.Parse()

// make sure required fields 'zabbix' and 'host' are populated
if *defaultHost=="" || *targetHost=="" {
	flag.PrintDefaults()
	return
}
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All:true})
	if err != nil {
		panic(err)
	}
			var metrics []*Metric
	test1 := `{"data":[`
	for i, container := range containers {

		fmt.Printf("%s %s\n", container.State, container.Names[0])

		var app string
		if i < len(containers)-1 {
			app = `{"{#CONTAINER}":"` + container.Names[0]+`"},`
		} else {
			app = `{"{#CONTAINER}":"` + container.Names[0]+`"}`
		}
		test1 = test1 + app



		metrics = append(metrics, NewMetric(*targetHost, "docker.["+container.Names[0]+"]", container.State))

		//// Create instance of Packet class

	}
	test1 = test1 + `]}`
			metrics = append(metrics, NewMetric(*targetHost, "docker.discovery", test1))

			packet := NewPacket(metrics)
			//
			//// Send packet to zabbix
			z := NewSender(*defaultHost, *defaultPort)
			z.Send(packet)
}

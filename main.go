package main

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"flag"
	. "zabbix_sender/zabbix_lib"
	 "encoding/json"
)
type containerData struct {
		Name string `json:"{#CONTAINER}"`
}
func main() {

	defaultHost := flag.String("zabbix", "localhost", "zabbix server e.g. '127.0.0.1'")
	targetHost := flag.String("host", "Agent", "zabbix target host e.g. 'myhost'")
	defaultPort := flag.Int("port", 10051, "zabbix server port e.g. 10051")
	flag.Parse()

	// make sure required fields 'zabbix' and 'host' are populated
	if *defaultHost == "" || *targetHost == "" {
		flag.PrintDefaults()
		return
	}
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	var metrics []*Metric
	dataToJsonStruct := []*containerData{}
	for _, container := range containers {
		cont := &containerData{container.Names[0]}
		dataToJsonStruct = append(dataToJsonStruct, cont)
		metrics = append(metrics, NewMetric(*targetHost, "docker.["+container.Names[0]+"]", container.State))
	}

	dataJson := make(map[string][]*containerData)
	dataJson["data"] = dataToJsonStruct
	bufdataJson, _ := json.Marshal(dataJson)

	metrics = append(metrics, NewMetric(*targetHost, "docker.discovery", string(bufdataJson)))

	packet := NewPacket(metrics)
	//// Send packet to zabbix
	z := NewSender(*defaultHost, *defaultPort)
	z.Send(packet)
}

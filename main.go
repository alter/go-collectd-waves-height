package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"collectd.org/api"
	"collectd.org/exec"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

type nodeStatus struct {
	blockchainHeight int32
	stateHeight      int32
	updatedTimestamp int64
	updatedDate      string
}

func main() {
	nodestatus := new(nodeStatus)
	getJson("http://127.0.0.1:6869/node/status", nodestatus)
	vl := api.ValueList{
		Identifier: api.Identifier{
			Host:   exec.Hostname(),
			Plugin: "node.height",
			Type:   "gauge",
		},
		Time:     time.Now(),
		Interval: exec.Interval(),
		Values:   []api.Value{api.Gauge(nodestatus.blockchainHeight)},
	}
	exec.Putval.Write(context.Background(), &vl)
}

package main

import (
	"collectd.org/api"
	"collectd.org/exec"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type nodeStatus struct {
	BlockchainHeight int    `json:"blockchainHeight"`
	StateHeight      int    `json:"stateHeight"`
	UpdatedTimestamp int    `json:"updatedTimestamp"`
	UpdatedDate      string `json:"updatedDate"`
}

func main() {
	url := "http://127.0.0.1:6869/node/status"
	nodeClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "spacecount-tutorial")
	res, getErr := nodeClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	nodestatus := nodeStatus{}
	jsonErr := json.Unmarshal(body, &nodestatus)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	vl := api.ValueList{
		Identifier: api.Identifier{
			Host:   exec.Hostname(),
			Plugin: "node-height",
			Type:   "gauge",
		},
		Time:     time.Now(),
		Interval: exec.Interval(),
		Values:   []api.Value{api.Gauge(nodestatus.BlockchainHeight)},
	}
	exec.Putval.Write(context.Background(), &vl)
}

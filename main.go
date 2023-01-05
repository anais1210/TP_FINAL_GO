package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	KrakenAPI = "https://api.kraken.com/"
)

type Status struct {
	Result struct {
		Status string `json:"status"`
	}`json:"result"`
}

func main() {
	resp, err := http.Get(KrakenAPI + "/0/public/SystemStatus")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	var status Status
	errors := json.Unmarshal(body, &status)
	if errors != nil {
		fmt.Printf("Error: %s", errors)
		return
	}


	fmt.Printf("The current status of Kraken's API is %s", status.Result.Status)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	KrakenAPI = "https://api.kraken.com/"
)

type Status struct {
	Result struct {
		Status string `json:"status"`
	}`json:"result"`
}

type Pair struct {
	Altname string `json:"altname"`
	Wsname string `json:"wsname"`
	Base string `json:"base"`
	Quote string `json:"quote"`
	CostDecimals string `json:"cost_decimals"`
	PairDecimals string `json:"pair_decimals"`
}

type Pairs struct{
		PairList map[string]Pair `json:"result"`
}


func getPair(){
	resp, err := http.Get(KrakenAPI + "/0/public/AssetPairs")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	//Write in file
	
	error := os.WriteFile("Archive/pairsKraken.csv", []byte(i.Altname), 0644)
	if error != nil {
		panic(error)
	}

	var m Pairs

	errors := json.Unmarshal(body, &m)
	if errors != nil {
		fmt.Printf("Error: %s", errors)
		return
	}
	var listOfPairs []string
	for _, i := range m.PairList {
		listOfPairs = append(listOfPairs, i.Altname)

	}

}
func writeData(){

}

func getStatus(){
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
func createFolder(){
	err := os.Mkdir("Archive", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	file, err := os.Create("Archive/pairsKakren.csv")
	// err = os.WriteFile("Arhive/pairsKraken.csv", []byte("Hello, Gophers!"), 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}
func main() {
	// getStatus()
	getPair()
	// createFolder()
}

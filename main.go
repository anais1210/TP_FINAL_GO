package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
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
	// CostDecimal int `json:"cost_decimals"`
	// PairDecimal int `json:"pair_decimals"`
}

type Pairs struct{
	PairList map[string]Pair `json:"result"`
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "toto"
	password = "mysecretpassword"
	dbname   = "mydatabase"
	schema   = "public"
)

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
	var m Pairs

	errors := json.Unmarshal(body, &m)
	if errors != nil {
		fmt.Printf("Error: %s", errors)
		return
	}

	// err := os.Mkdir("Archive", 0750)
	// if err != nil && !os.IsExist(err) {
	// 	log.Fatal(err)
	// }
	file, err := os.Create("Archive/pairsKraken.csv")
	// err = os.WriteFile("Arhive/pairsKraken.csv", []byte("Hello, Gophers!"), 0660)
	defer file.Close()
    if err != nil {
        log.Fatalln("failed to open file", err)
    }
    w := csv.NewWriter(file)
    defer w.Flush()

	header := []string{"altname", "wsname", "base", "quote"}

	if err := w.Write(header); err != nil {
        log.Fatalln("error writing csv:", err)
    }
	for _, i := range m.PairList {
		var csvRow []string
		csvRow = append(csvRow, i.Altname, i.Wsname,i.Base, i.Quote)
        if err := w.Write(csvRow); err != nil {
            log.Fatalln("error writing csv:", err)
        }
		// records := []string {
		// 	{"altname", "wsname", "base", "quote"},
		// 	{i.Altname, i.Wsname,i.Base, i.Quote},
		// }
		
		// w.WriteAll(records)
	}
}
func getStatus(){
	resp, err := http.Get(KrakenAPI + "/0/public/SystemStatus")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

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

func connectDatabase(datas []string){
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs ( id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, base character varying, quote character varying, cost real, pair real ,PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
	_, errors := db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
	}

	now := time.Now().Unix()
	fmt.Println(datas)
	for _, data := range datas {
		// Insertion des données dans la base de données
		sqlStatement := fmt.Sprintf("INSERT INTO pairs (id, altname) VALUES (%d, '%s')", now, data)
		_, err := db.Exec(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		// Eviter les collisions de clé primaire
		now++
		// fmt.Println(data[i])
	}

}
func main() {
	// getStatus()
	getPair()
	// createFolder()
}

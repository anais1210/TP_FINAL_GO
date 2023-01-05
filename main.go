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

	// ---------------------------------------------------------------- WRITE IN CSV FILE ----------------------------------------------------------------

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
	// go func(){
		for _, i := range m.PairList {
			var csvRow []string
			csvRow = append(csvRow, i.Altname, i.Wsname,i.Base, i.Quote)
			if err := w.Write(csvRow); err != nil {
				log.Fatalln("error writing csv:", err)
			}
		}
	// }()

	// time.Sleep(time.Second)
// ---------------------------------------------------------------- INSERT IN DATABASE ----------------------------------------------------------------

	go func(){
		connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		fmt.Println("im here1")

		sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs (id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, base character varying, quoteK character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
		_, errors := db.Exec(sqlStat)
		if errors != nil {
			fmt.Println(errors)
		}

		file, err := os.Open("Archive/pairsKraken.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Lecture du fichier CSV
		// reader := csv.NewReader(file)
		// dataPairs, err := reader.ReadAll()
		// if err != nil {
		// 	panic(err)
		// }

		// now := time.Now().Unix()
		// for _, data := range dataPairs {
			// Insertion des données dans la base de données
			// altname := data[0]
			// wsname := data[1]
			// base := data[2]
			// quote := data[3]

			// sqlStatement := fmt.Sprintf("INSERT INTO pairs (id, altname, wsname, base, quote) VALUES (%d, '%s', '%s', '%s', '%s')", now, altname,wsname,base,quote)
			// _, err := db.Exec(sqlStatement)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// Eviter les collisions de clé primaire
			// now++
			// fmt.Println(data[i])
		// }
	}()
	// time.Sleep(time.Second)

}

func downloadFile(){
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("Archive/pairsKraken.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		dataPairs, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", `attachment; filename="pairsKraken.csv"`)

		writer := csv.NewWriter(os.Stdout)
		writer.WriteAll(dataPairs)
		
	})
	http.ListenAndServe(":8080", nil)
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

// func connectDatabase(datas []string){
// 	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
// 	db, err := sql.Open("postgres", connectionString)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs ( id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, base character varying, quote character varying, cost real, pair real ,PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
// 	_, errors := db.Exec(sqlStat)
// 	if errors != nil {
// 		fmt.Println(errors)
// 	}

// 	now := time.Now().Unix()
// 	fmt.Println(datas)
// 	for _, data := range datas {
// 		// Insertion des données dans la base de données
// 		sqlStatement := fmt.Sprintf("INSERT INTO pairs (id, altname) VALUES (%d, '%s')", now, data)
// 		_, err := db.Exec(sqlStatement)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// Eviter les collisions de clé primaire
// 		now++
// 		// fmt.Println(data[i])
// 	}

// }
func main() {
	// getStatus()
	getPair()
	downloadFile()
	// createFolder()
}

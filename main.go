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
		Timestamp string `json:"timestamp"`
	}`json:"result"`
}

type Pair struct {
	Altname string `json:"altname"`
	Wsname string `json:"wsname"`
	Base string `json:"base"`
	Quote string `json:"quote"`
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

	fmt.Printf("Status of Kraken's API : %s at %s", status.Result.Status, status.Result.Timestamp)
}

func getPair()([]string){
	resp, err := http.Get(KrakenAPI + "/0/public/AssetPairs")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var m Pairs

	errors := json.Unmarshal(body, &m)
	if errors != nil {
		panic(errors)
	}
	createFolder("Archives")
	file, err := os.Create("Archives/pairsKraken.csv")
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	header := []string{"altname", "wsname", "base", "quote"}
	w := csv.NewWriter(file)
	defer w.Flush()
	if err := w.Write(header); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	
	var assetDatas []string
	for _, i := range m.PairList {
		assetDatas = nil
		assetDatas = append(assetDatas, i.Altname, i.Wsname,i.Base, i.Quote)
		if err := w.Write(assetDatas); err != nil {
			log.Fatalln("error writing csv:", err)
		}
	}
	return assetDatas
}
func createFolder(folderName string){
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			panic(err)
		}
	}
}
// ---------------------------------------------------------------- INSERT IN DATABASE ----------------------------------------------------------------
func databaseConnection(dataPairs []string){
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs (id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, base character varying, quoteK character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
	_, errors := db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
	}

	now := time.Now().Unix()
	for _, data := range dataPairs {
		// Insertion des données dans la base de données
		altname := data[0]
		wsname := data[1]
		base := data[2]
		quote := data[3]

		sqlStatement := fmt.Sprintf("INSERT INTO pairs (id, altname, wsname, base, quoteK) VALUES (%d, '%s', '%s', '%s', '%s')", now, altname,wsname,base,quote)
		_, err := db.Exec(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		now++
	}
}


// }
// func databaseView(){
// 	http.HandleFunc("/database", func(w http.ResponseWriter, r *http.Request) {
		
// 		w.Header().Set("Content-Type", "text/plain")
// 		w.Header().Set("Content-Disposition", `attachment; filename="pairsKraken.csv"`)
		
// 	})
// 	http.ListenAndServe(":8080", nil)
// }

// func downloadFile(){
// 	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
// 		file, err := os.Open("Archive/pairsKraken.csv")
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer file.Close()

// 		reader := csv.NewReader(file)
// 		dataPairs, err := reader.ReadAll()
// 		if err != nil {
// 			panic(err)
// 		}

// 		w.Header().Set("Content-Type", "text/plain")
// 		w.Header().Set("Content-Disposition", `attachment; filename="pairsKraken.csv"`)

// 		writer := csv.NewWriter(os.Stdout)
// 		writer.WriteAll(dataPairs)
		
// 	})
// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	getStatus()
	dataPairs := getPair()
	databaseConnection(dataPairs)
	// downloadFile()
}



//nom du fichier change en fonction de la date
//load env package
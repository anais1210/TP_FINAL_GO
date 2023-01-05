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

type AssetPrice struct {
	ActualPrice  []string `json:"c"`
	Volume       []string `json:"v"`
	High         []string `json:"h"`
	Low          []string `json:"l"`
	OpeningPrice string   `json:"o"`
}

type AssetPrices struct {
	FieldMap map[string]AssetPrice `json:"result"`
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
	connection := connectDB()
	createTable(connection)

	now := time.Now().Unix()
	var altname []string 
	for _, i := range m.PairList {
		assetDatas = nil
		altname = append(altname, i.Altname)

		assetDatas = append(assetDatas, i.Altname, i.Wsname,i.Base, i.Quote)
		if err := w.Write(assetDatas); err != nil {
		log.Fatalln("error writing csv:", err)
		}
		// sqlStatement := fmt.Sprintf("INSERT INTO pairs (id, altname, wsname, base, quoteK) VALUES (%d, '%s', '%s', '%s', '%s')", now, i.Altname, i.Wsname, i.Base,i.Quote)
		// _, err := connection.Exec(sqlStatement)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		now++
			
	}
	return altname
}
func getAssetPrice(altnames []string){
	for _, i := range altnames{
		resp, err := http.Get(KrakenAPI + "/0/public/Ticker?pair="+i)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var assetPrice AssetPrices
		errors := json.Unmarshal(body, &assetPrice)
		if errors != nil {
			// Handle the error
			fmt.Printf("Error: %s", errors)
			panic(errors)
		}
		fmt.Println(assetPrice)

	}
}
func createFolder(folderName string){
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			panic(err)
		}
	}
}
func connectDB()(db *sql.DB){
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
		// defer db.Close()
	return db
}

func createTable(db *sql.DB){
	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs (id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, base character varying, quoteK character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
	_, errors := db.Exec(sqlStat)
	if errors != nil {
	fmt.Println(errors)
	}
}
func selectDataFromDB(db *sql.DB)(){
	rows, err := db.Query("SELECT altname, wsname, base, quotek FROM pairs")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var v1, v2, v3, v4 string
	for rows.Next() {
		err = rows.Scan(&v1, &v2, &v3, &v4)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%s, %s, %s, %s", v1, v2, v3, v4)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
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
	// getStatus()
	altname :=getPair()
	db := connectDB()
	selectDataFromDB(db)
	getAssetPrice(altname)
	// databaseConnection(dataPairs)
	// downloadFile()
}



//nom du fichier change en fonction de la date
//load env package
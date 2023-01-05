package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "toto"
	password = "mysecretpassword"
	dbname   = "mydatabase"
	schema   = "public"
)
func connectDB()(db *sql.DB){
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
		defer db.Close()
	return db
}

func createTable(db *sql.DB){
	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs (id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, base character varying, quoteK character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to toto;"
	_, errors := db.Exec(sqlStat)
	if errors != nil {
	fmt.Println(errors)
	}
}

func insertData(dataPairs []string){
	now := time.Now().Unix()
	for _, data := range dataPairs {
		// Insertion des données dans la base de données
		altname := data[0]
		wsname := data[1]
		base := data[2]
		quote := data[3]

		sqlStatement := fmt.Sprintf("INSERT INTO pairs (id, altname, wsname, base, quote) VALUES (%d, '%s', '%s', '%s', '%s')", now, altname,wsname,base,quote)
		_, err := db.Exec(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		// Eviter les collisions de clé primaire
		now++
	}
}
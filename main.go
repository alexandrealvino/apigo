package main

import (
	"apigo/entities"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	//"database/sql"
	//"bytes"
	//"github.com/wcharczuk/go-chart" //exposes "chart"
	"text/template"

	"apigo/database"
)
//
var db *sql.DB // instanciating
var tmpl = template.Must(template.ParseGlob("form/*")) // template
//
func init() {
	database.CalculateQuota()
	//var x []float64
	//for _, value := range results  {
	//	x = append(x, value.Value)
	//}
	fmt.Println("Successfuly connected to MySQL database!")
}
//
func getWallet(w http.ResponseWriter, r *http.Request) { // get wallet
	var results []entities.Ticker
	results = database.GetWallet()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
//
func getTickerById(w http.ResponseWriter, r *http.Request) { // get ticker by id
	params := mux.Vars(r)
	id,_ := strconv.Atoi(params["id"])
	var result entities.Ticker
	result,_ = database.GetTickerById(id)
	if (result != entities.Ticker{}) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		println("share not found!")
	}
}
//
func getTickerBySymbol(w http.ResponseWriter, r *http.Request) { // get ticker by symbol
	params := mux.Vars(r)
	symbol := params["symbol"]
	var result entities.Ticker
	result, _ = database.GetTickerBySymbol(symbol)
	if (result != entities.Ticker{}) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		println("share not found!")
	}
}
//
func insertTicker(w http.ResponseWriter, r *http.Request) { // insert ticker
	r.Header.Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var ticker entities.Ticker  // Unmarshal
	err = json.Unmarshal(b, &ticker)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	shareInWallet, _ := database.TickerInWallet(ticker.Symbol)
	if (shareInWallet == entities.Ticker{})  {
		database.InsertTicker(ticker)
	}	else {
		ticker.ID = shareInWallet.ID
		database.UpdateTicker(ticker)
	}
}
//
func updateTicker(w http.ResponseWriter, r *http.Request) {  // update ticker
	r.Header.Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var ticker entities.Ticker  // Unmarshal
	err = json.Unmarshal(b, &ticker)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	shareInWallet, _ := database.TickerInWallet(ticker.Symbol)
	if (shareInWallet == entities.Ticker{})  {
		println("share not in found!")
	}	else {
		ticker.ID = shareInWallet.ID
		database.UpdateTicker(ticker)
	}
}
//
func deleteTicker(w http.ResponseWriter, r *http.Request) {  // delete ticker
	params := mux.Vars(r)
	id,_ := strconv.Atoi(params["id"])
	idNotInWallet,_ := database.GetTickerById(id)
	if (idNotInWallet == entities.Ticker{}) {
		println("id not found!")
	} else {
		database.DeleteTicker(id)
	}
}
//
func calculateQuota(w http.ResponseWriter, r*http.Request)  { // calculates all quotas
	database.CalculateQuota()
}
//
func main() {
	router := mux.NewRouter()  // init router
	log.Println("Server started on: http://localhost:8000")
	// router handlers
	router.HandleFunc("/api/wallet", getWallet).Methods("GET")
	router.HandleFunc("/api/wallet/{id}", getTickerById).Methods("GET")
	router.HandleFunc("/api/wallet/{symbol}", getTickerBySymbol).Methods("GET")
	router.HandleFunc("/api/wallet", insertTicker).Methods("POST")
	router.HandleFunc("/api/wallet", updateTicker).Methods("PUT")
	router.HandleFunc("/api/wallet/{id}", deleteTicker).Methods("DELETE")
	router.HandleFunc("/api/wallet/quota/", calculateQuota).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router)) // if error return fatal log
}
//

package main

import (
	"apigo/entities"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	//"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
	"io/ioutil"
	"strconv"

	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	//"database/sql"
	//"bytes"
	//"github.com/wcharczuk/go-chart" //exposes "chart"

	"apigo/database"

	//"github.com/piquette/finance-go/quote"
)
//
var db *sql.DB // instanciating
//
func init() {
	err := database.CalculateQuota()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly connected to MySQL database!")
	//var w http.ResponseWriter
	//var r *http.Request
	//getPrices(w,r)
}
//
func getWallet(w http.ResponseWriter, r *http.Request) { // get wallet
	var results []entities.Ticker
	results, err := database.GetWallet()
	if err != nil {
		panic(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	encoder.Encode(results)
	//json.NewEncoder(w).Encode(results)
}
//
func getTickerById(w http.ResponseWriter, r *http.Request) { // get ticker by id
	params := mux.Vars(r)
	id,_ := strconv.Atoi(params["id"])
	var result entities.Ticker
	result, err := database.GetTickerById(id)
	if err != nil {
		panic(err.Error())
	}
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
	result, err := database.GetTickerBySymbol(symbol)
	if err != nil {
		panic(err.Error())
	}
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
		panic(err.Error())
	}
	var ticker entities.Ticker  // Unmarshal
	err = json.Unmarshal(b, &ticker)
	if err != nil {
		panic(err.Error())
	}
	shareInWallet, _ := database.TickerInWallet(ticker.Symbol)
	if (shareInWallet == entities.Ticker{})  {
		err := database.InsertTicker(ticker)
		if err != nil {
			panic(err.Error())
		}
	}	else {
		ticker.ID = shareInWallet.ID
		err := database.UpdateTicker(ticker)
		if err != nil {
			panic(err.Error())
		}
	}
}
//
func updateTicker(w http.ResponseWriter, r *http.Request) {  // update ticker
	r.Header.Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err.Error())
	}
	var ticker entities.Ticker  // Unmarshal
	err = json.Unmarshal(b, &ticker)
	if err != nil {
		panic(err.Error())
	}
	shareInWallet, _ := database.TickerInWallet(ticker.Symbol)
	if (shareInWallet == entities.Ticker{})  {
		println("share not in found!")
	}	else {
		ticker.ID = shareInWallet.ID
		err := database.UpdateTicker(ticker)
		if err != nil {
			panic(err.Error())
		}
	}
}
//
func deleteTicker(w http.ResponseWriter, r *http.Request) {  // delete ticker
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err.Error())
	}
	idNotInWallet, err := database.GetTickerById(id)
	if err != nil {
		panic(err.Error())
	}
	if (idNotInWallet == entities.Ticker{}) {
		println("id not found!")
	} else {
		err := database.DeleteTicker(id)
		if err != nil {
			panic(err.Error())
		}
	}
}
//
func calculateQuota(w http.ResponseWriter, r*http.Request)  { // calculates all quotas
	err := database.CalculateQuota()
	if err != nil {
		panic(err.Error())
	}
}
//
func getPrices(w http.ResponseWriter, r*http.Request)  { // get prices and update closes
	var results []entities.Ticker
	results, err := database.GetWallet()
	if err != nil {
		panic(err.Error())
	}
	var tickersList []string
	var ticker string
	for _, element := range results {
		ticker = element.Symbol + ".SA"
		tickersList = append(tickersList, ticker)
	}
	//fmt.Println(tickersList)
	//var q finance.Quote
	var closePrices []float64
	var previousClosePrice, lastChangePercent float64
	for _, ticker := range tickersList {
		q, err := quote.Get(ticker)
		if err != nil {
			// Uh-oh.
			panic(err)
		}
		//fmt.Println(q)
		ticker = strings.Replace(ticker, ".SA", "", 1)
		//fmt.Println(ticker)
		var tic entities.Ticker
		previousClosePrice = q.RegularMarketPrice
		lastChangePercent = q.RegularMarketChangePercent
		closePrices = append(closePrices, q.RegularMarketPreviousClose)
		tic.Symbol = ticker
		tic.PreviousClose = previousClosePrice
		tic.LastChangePercent = lastChangePercent
		if tic.Symbol != "" {
			err = database.UpdatePrices(tic)
			if err != nil {
				panic(err.Error())
			}
		} else {
			return
		}
	}
	//k,_ := quote.Get("WEGE3.SA")
	//var close float64
	//close := q.RegularMarketPreviousClose
	// Success!
	//fmt.Println(q,closePrices)
	//w.Header().Add("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(k)
}
//
func getAvgPrice(w http.ResponseWriter, r*http.Request)  { // calculates avg price of ticker
	params := mux.Vars(r)
	symbol := params["symbol"]
	result, err := database.GetAvgPrice(symbol)
	if err != nil {
		panic(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
//
func calculateAvgPrice(w http.ResponseWriter, r*http.Request)  { // calculates all quotas
	err := database.CalculateAvgPrice()
	if err != nil {
		panic(err.Error())
	}
}
//
func calculateChangeFromAvgPrice(w http.ResponseWriter, r*http.Request)  { // calculates all changes from avg
	err := database.CalculateChangeFromAvgPrice()
	if err != nil {
		panic(err.Error())
	}
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
	// get stock price
	router.HandleFunc("/api/wallet/quota/price", getPrices).Methods("GET")
	router.HandleFunc("/api/wallet/quota/avgprice", calculateAvgPrice).Methods("GET")
	router.HandleFunc("/api/wallet/quota/changeAll", calculateChangeFromAvgPrice).Methods("GET")
	router.HandleFunc("/api/wallet/quota/{symbol}", getAvgPrice).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router)) // if error return fatal log
}
//

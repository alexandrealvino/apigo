package main

import (
	"apigo/adapter"
	"fmt"

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
//var db *sql.DB // instanciating
//
func init() {
	err := database.CalculateQuotas()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly connected to MySQL database!")
	//var w http.ResponseWriter
	//var r *http.Request
	//getPrices(w,r)
}
//
func main() {
	router := mux.NewRouter()  // init router
	log.Println("Server started on: http://localhost:8000")

	// router handlers
	router.HandleFunc("/api", adapter.GetIndexT).Methods("GET")
	router.HandleFunc("/api/addticker", adapter.AddPageT).Methods("GET")
	router.HandleFunc("/api/addticker", adapter.AddStockBuy).Methods("POST")
	router.HandleFunc("/api/addticker", adapter.AddStockSell).Methods("DELETE")

	router.HandleFunc("/api/wallet", adapter.GetWallet).Methods("GET") // get wallet and return data in json
	router.HandleFunc("/api/walletup/", adapter.GetWalletTRefreshingAllValues).Methods("GET") // get wallet and return data with updated quotas in json
	router.HandleFunc("/api/wallet/{id}", adapter.GetTickerById).Methods("GET") // returns ticker by id
	router.HandleFunc("/api/wallet/{symbol}", adapter.GetTickerBySymbol).Methods("GET") // return ticker by symbol
	router.HandleFunc("/api/wallet", adapter.InsertTicker).Methods("POST") // insert ticker in database
	router.HandleFunc("/api/wallet", adapter.UpdateTicker).Methods("PUT") // update ticker in the database
	router.HandleFunc("/api/wallet/{id}", adapter.DeleteTicker).Methods("DELETE") // delete ticker from database
	router.HandleFunc("/api/wallet/quota/", adapter.CalculateQuotas).Methods("GET") // calculates all quotas
	// get stock price
	router.HandleFunc("/api/wallet/quota/price", adapter.GetPrices).Methods("GET") // fetch stocks price from !yahoo finance
	router.HandleFunc("/api/wallet/quota/avgprice", adapter.CalculateAllAvgPrice).Methods("GET") // calculates all avg prices
	router.HandleFunc("/api/wallet/quota/changeAll", adapter.CalculateAllChangeFromAvgPrice).Methods("GET") // calculate earnings from all tickers in wallet
	router.HandleFunc("/api/wallet/quota/{symbol}", adapter.GetAvgPrice).Methods("GET") // calculates and returns the avg price of the ticker
	router.HandleFunc("/api/wallet/prices/up", adapter.UpdatePricesTable).Methods("GET") // // get prices from !yahoo finance, calculates results and updates Prices table

	log.Fatal(http.ListenAndServe(":8000", router)) // if error return fatal log
}
//

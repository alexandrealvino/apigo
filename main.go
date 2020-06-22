package main

import (
	"apigo/entities"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

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

func getIndex(w http.ResponseWriter, r *http.Request) { // render "wallet" template
	var results []entities.Ticker
	results, err := database.GetWallet()
	if err != nil {	panic(err.Error())} else {}
	tmpl, _ := template.ParseFiles("static/wallet.html")
	err = tmpl.Execute(w, results)
	if err != nil {
		panic(err.Error())
	}
} // render "wallet" template
//
func addPageT(w http.ResponseWriter, r *http.Request) { // execute "addticker" template
	Title := "New buy"
	tmpl, _ := template.ParseFiles("static/addticker.html")
	err := tmpl.Execute(w, Title)
	if err != nil {
		panic(err.Error())
	}
} // execute "addticker" template
//
func addStockBuy(w http.ResponseWriter, r *http.Request) { // add stock buy to the "buys" table in database
	err := r.ParseForm()
	if err != nil {
		panic(err.Error()) // Handle error here via logging and then return
	}
	var buy entities.StockBuy
	buy.Symbol = r.PostFormValue("ticker")
	buy.Quantity, _ = strconv.Atoi(r.PostFormValue("quantity"))
	buy.Value, _ = strconv.ParseFloat(r.PostFormValue("price"), 64)
	_ = database.InsertBuy(buy)
} // add stock buy to the "buys" table in database
//
func getWallet(w http.ResponseWriter, r *http.Request) { // get wallet and returns json format
	var results []entities.Ticker
	results, err := database.GetWallet()
	if err != nil {	panic(err.Error())} else {}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	_ = encoder.Encode(results)
	//json.NewEncoder(w).Encode(results)
} // get wallet in json format
//
func getWalletTRefreshingAllValues(w http.ResponseWriter, r *http.Request) { // get wallet refreshing all values and execute wallet template
	var results []entities.Ticker
	getPrices(w,r)
	calculateAllChangeFromAvgPrice(w,r)
	results, err := database.GetWalletRefreshingAllValues()
	if err != nil {	panic(err.Error())} else {}
	tmpl, _ := template.ParseFiles("static/wallet.html")
	err = tmpl.Execute(w, results)
	if err != nil {
		panic(err.Error())
	}
} // get wallet refreshing all values and execute wallet template
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
		_ = json.NewEncoder(w).Encode(result)
	} else {
		println("share not found!")
	}
} // get ticker by id
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
		_ = json.NewEncoder(w).Encode(result)
	} else {
		println("share not found!")
	}
} // get ticker by symbol
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
} // insert ticker
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
} // update ticker
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
} // delete ticker
//
func calculateQuotas(w http.ResponseWriter, r*http.Request)  { // calculates all quotas
	err := database.CalculateQuotas()
	if err != nil {
		panic(err.Error())
	}
} // calculates all quotas
//
// from !yahoo finance API
func getPrices(w http.ResponseWriter, r*http.Request)  { // get prices from !yahoo finance and update closes
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
	var closePrices []float64
	var previousClosePrice, lastChangePercent float64
	for _, ticker := range tickersList {
		q, err := quote.Get(ticker)
		if err != nil {
			// Uh-oh.
			panic(err)
		}
		ticker = strings.Replace(ticker, ".SA", "", 1)
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
} // get prices from !yahoo finance and update closes
//
func getAvgPrice(w http.ResponseWriter, r*http.Request)  { // calculates and get the avg price of the ticker in wallet
	params := mux.Vars(r)
	symbol := params["symbol"]
	result, err := database.GetAvgPrice(symbol)
	if err != nil {
		panic(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		panic(err.Error())
	}
} // calculates and get the avg price of the ticker in wallet
//
func calculateAllAvgPrice(w http.ResponseWriter, r*http.Request)  { // calculates avg price of ticker in wallet
	err := database.CalculateAllAvgPrice()
	if err != nil {
		panic(err.Error())
	}
} // calculates avg price of ticker in wallet
//
func calculateAllChangeFromAvgPrice(w http.ResponseWriter, r*http.Request)  { // calculates the return of the ticker from the avg price
	err := database.CalculateAllChangeFromAvgPrice()
	if err != nil {
		panic(err.Error())
	}
} // calculates the return of the ticker from the avg price
//
func updatePricesTable(w http.ResponseWriter, r*http.Request)  { // get prices from !yahoo finance, calculates results and updates Prices table
	var results []string
	results, err := database.GetTickersList()
	if err != nil {
		panic(err.Error())
	}
	var tickersList []string
	var ticker string
	for _, element := range results {
		ticker = element + ".SA"
		tickersList = append(tickersList, ticker)
	}
	var lastPrice, preMarketPrice, weekResult, monthResult, yearResult float64
	for _, ticker := range tickersList {
		q, err := quote.Get(ticker)
		if err != nil {
			panic(err)
		}
		ticker = strings.Replace(ticker, ".SA", "", 1)
		lastPrice = q.RegularMarketPrice
		preMarketPrice = q.RegularMarketPreviousClose
		weekDay := time.Now().Weekday()
		monthDay := time.Now().Day()
		month := time.Now().Month()
		hour := time.Now().Hour()
		year := time.Now().Year()
		resultsList , lastUpdate, _ := database.GetResults(ticker)
		weekResult =  resultsList[0]
		monthResult =  resultsList[1]
		yearResult =  resultsList[2]
		date := strconv.Itoa(monthDay) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(year)
		if weekDay == 1 {
			weekResult = lastPrice
		}
		if weekDay != 0 && weekDay !=6  && date != lastUpdate && hour >=18 && (q.MarketState == "CLOSED" || q.MarketState == "POSTPOST") {
			weekResult = weekResult + 100*(lastPrice - q.RegularMarketPreviousClose)/q.RegularMarketPreviousClose
			if month == 1 && monthDay == 1 {
				yearResult = 0.00
			}
			if monthDay == 1 && month !=1 {
				monthResult = 0.00
			}
			if monthDay != 1 {
				monthResult = monthResult + 100*(lastPrice - q.RegularMarketPreviousClose)/q.RegularMarketPreviousClose
				yearResult = yearResult + 100*(lastPrice - q.RegularMarketPreviousClose)/q.RegularMarketPreviousClose
			}
		} // check if market is closed and if is not a weekend day
		err = database.UpdateTablePrices(lastPrice, preMarketPrice, weekResult, monthResult, yearResult, date, ticker) // calls database func
		if err != nil {panic(err)}
		}
	q, _ := quote.Get("WEGE3.SA")
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	_ = encoder.Encode(q)
	fmt.Println("Successfuly updated table prices!")
} // get prices from !yahoo finance, calculates results and updates Prices table
//
func main() {
	router := mux.NewRouter()  // init router
	log.Println("Server started on: http://localhost:8000")

	// router handlers
	router.HandleFunc("/api", getIndex).Methods("GET")
	router.HandleFunc("/api/addticker", addPageT).Methods("GET")
	router.HandleFunc("/api/addticker", addStockBuy).Methods("POST")

	router.HandleFunc("/api/wallet", getWallet).Methods("GET") // get wallet and return data in json
	router.HandleFunc("/api/walletup/", getWalletTRefreshingAllValues).Methods("GET") // get wallet and return data with updated quotas in json
	router.HandleFunc("/api/wallet/{id}", getTickerById).Methods("GET") // returns ticker by id
	router.HandleFunc("/api/wallet/{symbol}", getTickerBySymbol).Methods("GET") // return ticker by symbol
	router.HandleFunc("/api/wallet", insertTicker).Methods("POST") // insert ticker in database
	router.HandleFunc("/api/wallet", updateTicker).Methods("PUT") // update ticker in the database
	router.HandleFunc("/api/wallet/{id}", deleteTicker).Methods("DELETE") // delete ticker from database
	router.HandleFunc("/api/wallet/quota/", calculateQuotas).Methods("GET") // calculates all quotas
	// get stock price
	router.HandleFunc("/api/wallet/quota/price", getPrices).Methods("GET") // fetch stocks price from !yahoo finance
	router.HandleFunc("/api/wallet/quota/avgprice", calculateAllAvgPrice).Methods("GET") // calculates all avg prices
	router.HandleFunc("/api/wallet/quota/changeAll", calculateAllChangeFromAvgPrice).Methods("GET") // calculate earnings from all tickers in wallet
	router.HandleFunc("/api/wallet/quota/{symbol}", getAvgPrice).Methods("GET") // calculates and returns the avg price of the ticker
	router.HandleFunc("/api/wallet/prices/up", updatePricesTable).Methods("GET") // // get prices from !yahoo finance, calculates results and updates Prices table
	log.Fatal(http.ListenAndServe(":8000", router)) // if error return fatal log
}
//

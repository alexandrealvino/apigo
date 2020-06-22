package database

import (
	"apigo/config"
	"apigo/entities"
	"fmt"
	"math"
	"strings"
	"time"
)
//
func init() {
	fmt.Println("Successfuly connected to MySQL database!")
}
//
func GetWallet() ([]entities.Ticker, error){  // get wallet
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT id, symbol, value , quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice FROM tickers ORDER BY value DESC")
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	tickersList := []entities.Ticker{}
	for results.Next() {
		err = results.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota, &tic.AvgPrice, &tic.PreviousClose, &tic.LastChangePercent, &tic.ChangeFromAvgPrice)
		if err != nil {
			panic(err.Error())
		}
		tickersList = append(tickersList, tic)
	}
	return tickersList, err
}
//
func GetWalletRefreshingAllValues() ([]entities.Ticker, error){  // get wallet refreshing all values
	db := config.DbConn()
	defer db.Close()
	err := CalculateQuotas()
	if err != nil {
		panic(err.Error())
	}
	err = CalculateAllAvgPrice()
	if err != nil {
		panic(err.Error())
	}
	results, err := db.Query("SELECT id, symbol, value , quota, avgPrice, previousClose, lastChangePercent, changeFromAvgPrice FROM tickers ORDER BY value DESC")
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	tickersOfWallet := []entities.Ticker{}
	for results.Next() {
		err = results.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota, &tic.AvgPrice, &tic.PreviousClose, &tic.LastChangePercent, &tic.ChangeFromAvgPrice)
		if err != nil {
			panic(err.Error())
		}
		tickersOfWallet = append(tickersOfWallet, tic)
	}
	fmt.Println("All values refreshed!")
	return tickersOfWallet, err
}
//
func GetTickerById(id int) (entities.Ticker, error){  // get ticker by id
	db := config.DbConn()
	defer db.Close()
	err := CalculateQuotas()
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query("SELECT id, symbol, value , quota FROM tickers WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	for result.Next() {
		err = result.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota)
		if err != nil {
			panic(err.Error())
		}
	}
	return tic, err
}
//
func GetTickerBySymbol(symbol string) (entities.Ticker, error){  // get ticker by symbol
	db := config.DbConn()
	defer db.Close()
	err := CalculateQuotas()
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query("SELECT id, symbol, value , quota FROM tickers WHERE symbol = ?", symbol)
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	for result.Next() {
		err = result.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota)
		if err != nil {
			panic(err.Error())
		}
	}
	return tic, err
}
//
func InsertTicker(ticker entities.Ticker) (error) {  // insert ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("INSERT INTO tickers (symbol, value, quota) VALUES (?, ?, ?);", ticker.Symbol, ticker.Value, ticker.Quota)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly inserted share!")
	return err
}
//
func UpdateTicker(ticker entities.Ticker) (error) {  // update ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("UPDATE tickers  SET value = ?, quota = ? WHERE id = ?;", ticker.Value, ticker.Quota, ticker.ID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly updated share!")
	return err
}
//
func DeleteTicker(id int) (error) { // delete ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("DELETE FROM tickers WHERE id = ?;", id)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly deleted share!")
	return err
}
//
func CalculateQuotas() (error)  {  // calculates all quotas
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT id, symbol, value, quota FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	var totalInvested float64
	tic := entities.Ticker{}
	for results.Next() {
		err = results.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota)
		if err != nil {
			panic(err.Error())
		}
		totalInvested = totalInvested + tic.Value
	}
	results2, err := db.Query("SELECT id, symbol, value, quota FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	for results2.Next() {
		err = results2.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota)
		if err != nil {
			panic(err.Error())
		}
		tic.Quota = tic.Value/totalInvested*100
		_, err = db.Query("UPDATE tickers  SET quota = ? WHERE id = ?;", tic.Quota, tic.ID)
		if err != nil {
			panic(err.Error())
		}
	}
	return err
}
//
func TickerInWallet(symbol string) (entities.Ticker, error){  // checks if ticker is already in wallet
	db := config.DbConn()
	defer db.Close()
	result, err := db.Query("SELECT id , value FROM tickers WHERE symbol = ?", symbol)
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	for result.Next() {
		err = result.Scan(&tic.ID, &tic.Value)
		if err != nil {
			panic(err.Error())
		}
	}
	return tic, err
}
//
func GetAvgPrice(symbol string) (float64, error){  // calculates and returns the avg price of the ticker
	db := config.DbConn()
	defer db.Close()
	result, err := db.Query("SELECT price, quantity FROM buys WHERE symbol = ?", symbol)
	if err != nil {
		panic(err.Error())
	}
	var price, priceSum, priceAvg, quantity, quantitySum float64
	for result.Next() {
		err = result.Scan(&price, &quantity)
		if err != nil {
			panic(err.Error())
		}
		priceSum += price
		quantitySum += quantity
	}
	priceAvg = priceSum/quantitySum
	return priceAvg, err
}
//
func CalculateAllAvgPrice() (error) {  // calculates all avg prices and update their values in the database
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT symbol FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	var ticker string
	var tickersList []string
	var priceAvgList []float64
	for results.Next() {
		err = results.Scan(&ticker)
		if err != nil {
			panic(err.Error())
		}
		tickersList = append(tickersList, ticker)
	}
	for _, ticker := range tickersList{
		priceAvg, err := GetAvgPrice(ticker)
		if err != nil {
			panic(err.Error())
		}
		priceAvgList = append(priceAvgList, priceAvg)
		_, err = db.Query("UPDATE tickers  SET avgPrice = ? WHERE symbol = ?;", priceAvg, ticker)
	}
	return err
}
//
func UpdatePrices(ticker entities.Ticker) (error) {  // update previousClose and lastChangePercent by ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("UPDATE tickers SET previousClose = ?, lastChangePercent = ? WHERE symbol = ?;", ticker.PreviousClose, ticker.LastChangePercent, ticker.Symbol)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Println("Successfuly updated previous close and last change percent!")
	return err
}
//
func CalculateAllChangeFromAvgPrice() (error) {  // calculates all changes from avg prices
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT symbol, avgPrice, previousClose FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	ticker := entities.Ticker{}
	//var priceAvg, priceClose,changePrice float64
	for results.Next() {
		err = results.Scan(&ticker.Symbol, &ticker.AvgPrice, &ticker.PreviousClose)
		if err != nil {
			panic(err.Error())
		}
		ticker.ChangeFromAvgPrice = (ticker.PreviousClose - ticker.AvgPrice)/ticker.AvgPrice*100
		_, err = db.Query("UPDATE tickers  SET changeFromAvgPrice = ? WHERE symbol = ?;", ticker.ChangeFromAvgPrice, ticker.Symbol)
		if err != nil {
			panic(err.Error())
		}
	}
	err = CalculateEarnings()
	if err != nil {
		panic(err.Error())
	}
	return err
}
//
func InsertBuy(buy entities.StockBuy) (error) {  // insert stock buy to "buys" table and insert or update ticker in "tickers" table
	db := config.DbConn()
	defer db.Close()
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	date = strings.Replace(date, "-", "/", 2)
	year := date[0:4]
	month := date[5:7]
	day := date[8:10]
	date = day + "/" + month + "/" + year
	_, err := db.Query("INSERT INTO buys (symbol, price, quantity, date) VALUES (?, ?, ?, ?);", buy.Symbol, buy.Value, buy.Quantity , date)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly bought stock!")
	// handling exists cases
	tickerInWallet ,_ := TickerInWallet(buy.Symbol)
	var ticker entities.Ticker
	if (tickerInWallet == entities.Ticker{}) {
		ticker.Value = buy.Value
		ticker.Symbol = buy.Symbol
		ticker.AvgPrice = buy.Value/float64(buy.Quantity)
		_, err := db.Query("INSERT INTO tickers (symbol, value, avgPrice) VALUES (?, ?, ?);", ticker.Symbol, ticker.Value, ticker.AvgPrice)
		if err != nil {
			panic(err.Error())
		}
	} else {
		ticker.ID = tickerInWallet.ID
		ticker.Value = tickerInWallet.Value + buy.Value
		ticker.Symbol = buy.Symbol
		_, err := db.Query("UPDATE tickers  SET value = ? WHERE id = ?;", ticker.Value, ticker.ID)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Successfuly updated share!")
	}
	//
	return err
}
//
func GetResults(symbol string) ([]float64, string, error) { // get results from Prices table
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT weekResult, monthResult, yearResult, lastUpdate FROM prices WHERE symbol = ?;", symbol)
	if err != nil {
		panic(err.Error())
	}
	var weekResult, monthResult, yearResult float64
	var lastUpdate string
	var resultsList []float64
	for results.Next(){
		err = results.Scan(&weekResult, &monthResult, &yearResult, &lastUpdate)
		if err != nil {
			panic(err.Error())
		}
		resultsList = append(resultsList, weekResult, monthResult, yearResult)
	}
	return resultsList,lastUpdate, err
}
//
func UpdateTablePrices(lastPrice, preMarketPrice, weekResult, monthResult, yearResult float64, lastUpdate , symbol string) (error) {  // update table Prices
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("UPDATE prices  SET lastPrice = ?, lastClosePrice = ?, weekResult = ?, monthResult = ?, yearResult = ?, lastUpdate = ? WHERE symbol = ?;", lastPrice, preMarketPrice, weekResult, monthResult, yearResult, lastUpdate, symbol)
	if err != nil {
		panic(err.Error())
	}
	return err
} // update table Prices
//

// FUNCTIONS FOR INTERNAL USE
func GetTickersList() ([]string, error) {
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT symbol FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	var ticker string
	var tickersList []string
	for results.Next(){
		err = results.Scan(&ticker)
		if err != nil {
			panic(err.Error())
		}
		tickersList = append(tickersList, ticker)
	}
	return tickersList, err
} // returns list of tickers in wallet
//
func CalculateEarnings() (error)  {  // calculates earnings of wallet and prints in console
	db := config.DbConn()
	defer db.Close()
	results, err := db.Query("SELECT id, symbol, value, changeFromAvgPrice FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	var earnings, valueSum float64
	tic := entities.Ticker{}
	for results.Next() {
		err = results.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.ChangeFromAvgPrice)
		if err != nil {
			panic(err.Error())
		}
		valueSum = valueSum + tic.Value
		earnings = earnings + tic.Value*tic.ChangeFromAvgPrice
	}
	earnings = math.Round(earnings/valueSum*100)/100
	absoluteResult := math.Round(earnings*valueSum/100*100)/100
	fmt.Println("Wallet earnings: ",earnings , " ", absoluteResult)
	return err
} // calculates earnings of wallet and prints in console
//

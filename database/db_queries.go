package database

import (
	"apigo/config"
	"apigo/entities"
	"fmt"
)
//
func init() {
	fmt.Println("Successfuly connected to MySQL database!")
}
//
func GetWallet() ([]entities.Ticker, error){  // get wallet
	db := config.DbConn()
	defer db.Close()
	err := CalculateQuota()
	if err != nil {
		panic(err.Error())
	}
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
func GetTickerById(id int) (entities.Ticker, error){  // get ticker by id
	db := config.DbConn()
	defer db.Close()
	err := CalculateQuota()
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
	err := CalculateQuota()
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
func CalculateQuota() (error)  {  // calculates quota
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
	result, err := db.Query("SELECT id FROM tickers WHERE symbol = ?", symbol)
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	for result.Next() {
		err = result.Scan(&tic.ID)
		if err != nil {
			panic(err.Error())
		}
	}
	return tic, err
}
//
func GetAvgPrice(symbol string) (float64, error){  // get the avg price of the ticker
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
func CalculateAvgPrice() (error) {  // calculates all avg prices
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
func UpdatePrices(ticker entities.Ticker) (error) {  // update previous close by ticker
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
func CalculateChangeFromAvgPrice() (error) {  // calculates all changes from avg prices
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
		fmt.Println(ticker.ChangeFromAvgPrice)
		_, err = db.Query("UPDATE tickers  SET changeFromAvgPrice = ? WHERE symbol = ?;", ticker.ChangeFromAvgPrice, ticker.Symbol)
		if err != nil {
			panic(err.Error())
		}
	}
	return err
}
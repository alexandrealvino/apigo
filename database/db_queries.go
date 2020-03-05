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
func GetWallet() []entities.Ticker{  // get wallet
	db := config.DbConn()
	defer db.Close()
	CalculateQuota()
	results, err := db.Query("SELECT id, symbol, value , quota FROM tickers ORDER BY value DESC")
	if err != nil {
		panic(err.Error())
	}
	tic := entities.Ticker{}
	tickersList := []entities.Ticker{}
	for results.Next() {
		err = results.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota)
		if err != nil {
			panic(err.Error())
		}
		tickersList = append(tickersList, tic)
	}
	return tickersList
}
//
func GetTickerById(id int) (entities.Ticker, error){  // get ticker
	db := config.DbConn()
	defer db.Close()
	CalculateQuota()
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
func GetTickerBySymbol(symbol string) (entities.Ticker, error){  // get ticker
	db := config.DbConn()
	defer db.Close()
	CalculateQuota()
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
func InsertTicker(ticker entities.Ticker) {  // insert ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("INSERT INTO tickers (symbol, value, quota) VALUES (?, ?, ?);", ticker.Symbol, ticker.Value, ticker.Quota)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly inserted share!")
}
//
func UpdateTicker(ticker entities.Ticker) {  // update ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("UPDATE tickers  SET value = ?, quota = ? WHERE id = ?;", ticker.Value, ticker.Quota, ticker.ID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly updated share!")
}
//
func DeleteTicker(id int) { // delete ticker
	db := config.DbConn()
	defer db.Close()
	_, err := db.Query("DELETE FROM tickers WHERE id = ?;", id)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfuly deleted share!")
}
//
func CalculateQuota()  {  // calculates quota
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
		totalInvested = totalInvested + tic.Value
	}
	results2, err := db.Query("SELECT id, symbol, value, quota FROM tickers")
	if err != nil {
		panic(err.Error())
	}
	for results2.Next() {
		err = results2.Scan(&tic.ID, &tic.Symbol, &tic.Value, &tic.Quota)
		tic.Quota = tic.Value/totalInvested*100
		_, err = db.Query("UPDATE tickers  SET quota = ? WHERE id = ?;", tic.Quota, tic.ID)
	}
}
//
func TickerInWallet(symbol string) (entities.Ticker, error){  // checks if ticker is in wallet
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

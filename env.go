package main

import "os"

var dbName = os.Getenv("MYSQL_USER")
var dbUser = os.Getenv("MYSQL_USERNAME")
var dbPass = os.Getenv("MYSQL_PASSWORD")
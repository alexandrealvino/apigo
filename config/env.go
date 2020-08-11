package config

import "os"

var dbDriver = os.Getenv("MYSQL_DRIVER")
var dbName = os.Getenv("MYSQL_DATABASE")
var dbUser = os.Getenv("MYSQL_USER")
var dbPass = os.Getenv("MYSQL_PASSWORD")

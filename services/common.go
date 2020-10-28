package services

import (
	"io/ioutil"
	"strings"

	"github.com/jinzhu/gorm"
)

var username string
var password string
var dbHost string
var dbPort string
var dbName string

var Db *gorm.DB

func readProperties() {
	content, _ := ioutil.ReadFile("config/db.config")

	lines := strings.Split(string(content), "\n")

	if len(lines) >= 6 {
		username = lines[1]
		password = lines[2]
		dbHost = lines[3]
		dbPort = lines[4]
		dbName = lines[5]
	}

}

func OpenDatabase() {
	//open a db connection
	readProperties()
	var err error
	Db, err = gorm.Open("postgres", "postgres://"+username+":"+password+"@"+dbHost+":"+dbPort+"/"+dbName+"?sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
}

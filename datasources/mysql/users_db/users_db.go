package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	mysqlUsersUsername = "mysql_users_username"
	mysqlUsersPassword = "mysql_users_password"
	mysqlUsersHost     = "mysql_users_host"
	mysqlUsersSchema   = "mysql_users_schema"
)

var (
	Client *sql.DB

	username = os.Getenv(mysqlUsersUsername)
	password = os.Getenv(mysqlUsersPassword)
	host     = os.Getenv(mysqlUsersHost)
	schema   = os.Getenv(mysqlUsersSchema)
)

//if import users_db package called init
func init() {
	//user - password - host - database name
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, password, host, schema)

	log.Println(fmt.Sprintf("about to connect to %s", dataSourceName)) //TODO : delete because log is for test

	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err) //try to connect
	}
	if err := Client.Ping(); err != nil {
		panic(err) //try to connect
	}
	log.Println("database successfully configure")
}

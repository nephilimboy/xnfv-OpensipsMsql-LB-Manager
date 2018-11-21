//package opensips_mysql_lb_adder
package main

import (
	_ "github.com/go-sql-driver/mysql"
)
import (
	"fmt"
	"log"
	"database/sql"
)

const dbURL = "root:rootroot@tcp(127.0.0.1:3306)/opensips"

/*
 * Tag... - a very simple struct
 */
type Tag struct {
	id       int    `json:"id"`
	username string `json:"username"`
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb

}

func write() {
	db, err := sql.Open("mysql", dbURL)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
}

func read() {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	results, err := db.Query("SELECT id, username FROM subscriber")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.id, &tag.username)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		log.Printf(tag.username)
	}
	defer db.Close()
}

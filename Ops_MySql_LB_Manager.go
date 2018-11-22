//package opensips_mysql_lb_adder
package main

import (
	_ "github.com/go-sql-driver/mysql"
)
import (
	"database/sql"
	"net/http"
	"encoding/json"
	"sync"
	"fmt"
	"time"
)

const dbURL = "root:rootroot@tcp(127.0.0.1:3306)/opensips"

type LoadBalancer struct {
	Group_id    string `json:"group_id"`
	Dst_uri     string `json:"dst_uri"`
	Resources   string `json:"resources"`
	Probe_mode  string `json:"probe_mode"`
	Description string `json:"description"`
}

func deleteLB(lb LoadBalancer, wg *sync.WaitGroup) string{
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	query := "DELETE FROM opensips.load_balancer WHERE dst_uri='" + lb.Dst_uri + "'"
	delete, err := db.Query(query)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer delete.Close()
	wg.Done() // Need to signal to waitgroup that this goroutine is done
	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05") + " --- " + "Load Balancer: " + lb.Dst_uri + " deleted from DataBase")
	return "delete LB: " + lb.Dst_uri
}

func addNewLB(lb LoadBalancer, wg *sync.WaitGroup) string {

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	query := "INSERT INTO opensips.load_balancer (group_id, dst_uri, resources, probe_mode, description) VALUES " +
		"('" + lb.Group_id + "', '" + lb.Dst_uri + "', '" + lb.Resources + "', '" + lb.Probe_mode + "', '" + lb.Description + "'); "
	insert, err := db.Query(query)
	if err != nil {
		fmt.Printf("%s", err.Error())
		//panic(err.Error())
	}
	defer insert.Close()
	wg.Done() // Need to signal to waitgroup that this goroutine is done
	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05") + " --- " + "Load Balancer: " + lb.Dst_uri + " added to DataBase")
	return "Add LB: " + lb.Dst_uri
}

// ================  Request Handlers ========================//

func createNewLB_Handler(w http.ResponseWriter, r *http.Request) {
	lb := LoadBalancer{}
	err := json.NewDecoder(r.Body).Decode(&lb)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	if len(lb.Group_id) == 0 || len(lb.Dst_uri) == 0 || len(lb.Resources) == 0 || len(lb.Probe_mode) == 0{
		w.WriteHeader(http.StatusBadRequest)
	} else {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		am := addNewLB(lb, wg)
		wg.Wait()
		fmt.Fprintf(w, am)
	}
}

func deleteLB_Handler(w http.ResponseWriter, r *http.Request) {
	lb := LoadBalancer{}
	err := json.NewDecoder(r.Body).Decode(&lb)
	if err != nil {
		panic(err)
	}
	if len(lb.Dst_uri) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		am := deleteLB(lb, wg)
		wg.Wait()
		fmt.Fprintf(w, am)
	}
}

func read() {
	//db, err := sql.Open("mysql", dbURL)
	//if err != nil {
	//	panic(err.Error())
	//}
	//results, err := db.Query("SELECT id, username FROM subscriber")
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//}
	//for results.Next() {
	//	var tag Tag
	//	// for each row, scan the result into our tag composite object
	//	err = results.Scan(&tag.id, &tag.username)
	//	if err != nil {
	//		panic(err.Error()) // proper error handling instead of panic in your app
	//	}
	//	log.Printf(tag.username)
	//}
	//defer db.Close()
}


func main() {

	fmt.Println(" ")
	fmt.Println("****  OpenSips Mysql LoadBalancer Manager  ****")
	fmt.Println("****  By AH.GHORAB          ****")
	fmt.Println("****  Version 0.1           ****")
	fmt.Println("****  Fall 2018             ****")
	fmt.Println("------------------------------------")
	fmt.Println("[*] App Running at localhost:8000")
	fmt.Println("[*] Valid rest URLs")
	fmt.Println("[#] - /CreateNewLB")
	fmt.Println("[#] - /deleteLB")
	fmt.Println(" ")
	fmt.Println("------------ App Logs ------------")
	fmt.Println(" ")

	http.HandleFunc("/CreateNewLB", createNewLB_Handler)
	http.HandleFunc("/deleteLB", deleteLB_Handler)
	http.ListenAndServe(":8000", nil)
}
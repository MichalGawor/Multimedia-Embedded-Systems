package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "pi"
	password = "qwerty"
	dbname   = "test"
)

var db *sql.DB

//Environment data to be filled into tables
type Environment1 struct {
	ID          int `json:"id"`
	DayTime     int `json:"day_time"`
	NightTime   int `json:"night_time"`
	Moisture    int `json:"moisture"`
	Temperature int `json:"temperature"`
}

type Environment2 struct {
	ID          int `json:"id"`
	Moisture    int `json:"moisture"`
	Water		int `json:"water_level"`
	Temperature int `json:"temperature"`
}


func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)

	err = db.Ping()

	if err != nil {
		fmt.Println("PANIC")
		panic(err)
	}
	initialization()

	http.HandleFunc("/insertExpected", insertExpected)
	http.HandleFunc("/getExpected", getExpected)
	http.HandleFunc("/insertReal", insertReal)
	http.HandleFunc("/getReal", getReal)
	http.ListenAndServe(":8091", nil)
}

func insertExpected(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Environment1
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	insert(msg)
	json.NewEncoder(w).Encode(msg)
}

func insert(environment Environment1) {
	sqlStatement := `INSERT INTO expected_table(moisture, day_time, night_time, temperature) VALUES($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, environment.Moisture, environment.DayTime, environment.NightTime, environment.Temperature)
	if err != nil {
		panic(err)
	}
	fmt.Println("INSERT")
}

func getExpected(w http.ResponseWriter, r *http.Request) {
	env := getLastExpected()
	fmt.Println(env)
	json.NewEncoder(w).Encode(env)
}

func getLastExpected() Environment1 {
	sqlStatement := `SELECT * FROM expected_table ORDER by id DESC LIMIT 1;`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	var environment Environment1

	for rows.Next() {
		err = rows.Scan(&environment.ID, &environment.Moisture, &environment.DayTime, &environment.NightTime, &environment.Temperature)
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(environment)

	return environment
}

func insertReal(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Environment2
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	insertRealTable(msg)
	json.NewEncoder(w).Encode(msg)
}

func insertRealTable(environment Environment2) {
	sqlStatement := `INSERT INTO real_table(moisture, water, temperature) VALUES($1, $2, $3)`
	_, err := db.Exec(sqlStatement, environment.Moisture, environment.Water, environment.Temperature)
	if err != nil {
		panic(err)
	}
}

func getReal(w http.ResponseWriter, r *http.Request) {
	env := getLastReal()
	json.NewEncoder(w).Encode(env)
}

func getLastReal() Environment2 {
	sqlStatement := `SELECT * FROM real_table ORDER by id DESC LIMIT 1;`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	var environment Environment2

	for rows.Next() {
		err = rows.Scan(&environment.ID, &environment.Moisture, &environment.Water, &environment.Temperature)
		fmt.Println(environment)
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(environment)

	return environment
}

func initialization() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)

	result, err := db.Exec("CREATE TABLE IF NOT EXISTS real_table (id SERIAL PRIMARY KEY, moisture INT, water_level INT, temperature INT);")
	if err != nil {
		log.Fatal(err)
	}

	result, err = db.Exec("CREATE TABLE IF NOT EXISTS expected_table (id SERIAL PRIMARY KEY, day_time INT, night_time INT, moisture INT, temperature INT);")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result, err, "INITIALIZATION ENDED")
}

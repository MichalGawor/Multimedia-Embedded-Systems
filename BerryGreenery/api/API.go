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
type Environment struct {
	ID          int `json:"id"`
	Moisture    int `json:"moisture"`
	DayTime     float32 `json:"day_time"`
	NightTime   float32 `json:"night_time"`
	Temperature int `json:"temperature"`
	Water		int `json:"water_level"`
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
	var msg Environment
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	insert(msg)
	json.NewEncoder(w).Encode(msg)
}

func insert(environment Environment) {
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

func getLastExpected() Environment {
	sqlStatement := `SELECT * FROM expected_table ORDER by id DESC LIMIT 1;`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	var environment Environment

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
	var msg Environment
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	insertRealTable(msg)
	json.NewEncoder(w).Encode(msg)
}

func insertRealTable(environment Environment) {
	sqlStatement := `INSERT INTO real_table(moisture, water, temperature) VALUES($1, $2, $3)`
	_, err := db.Exec(sqlStatement, environment.Moisture, environment.Water, environment.Temperature)
	if err != nil {
		panic(err)
	}
}

func getReal(w http.ResponseWriter, r *http.Request) {
	env := getLastExpected()
	fmt.Println(env)
	json.NewEncoder(w).Encode(env)
}

func getLastReal() Environment {
	sqlStatement := `SELECT * FROM real_table ORDER by id DESC LIMIT 1;`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	var environment Environment

	for rows.Next() {
		err = rows.Scan(&environment.ID, &environment.Moisture, &environment.Water, &environment.Temperature)
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

	result, err := db.Exec("CREATE TABLE IF NOT EXISTS real_table (id SERIAL PRIMARY KEY, moisture INT, water INT, temperature INT);")
	if err != nil {
		log.Fatal(err)
	}

	result, err = db.Exec("CREATE TABLE IF NOT EXISTS expected_table (id SERIAL PRIMARY KEY, day_time REAL, night_time REAL, moisture INT, temperature INT);")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result, err, "INITIALIZATION ENDED")
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db_file    string
	cnx        *sql.DB
	opened     bool
	cache_file string
	Types      struct {
		Cultures           []string
		Cultures_speciales []string
		Locaux             []string
	}
}

func (db *Database) setup(filename string) {
	db.db_file = filename
	db.opened = false
	db.cache_file = "./data/db_distincts.json"
}

func (db *Database) open() {
	if !db.opened {
		db.cnx, _ = sql.Open("sqlite3", db.db_file)
		db.opened = true
	}
}

func (db *Database) close() {
	db.cnx.Close()
	db.opened = false
}

//
// Create database file and table if needed
//
func (db *Database) create_if_no_exists() {
	_, err := os.Stat(db.db_file)
	if os.IsNotExist(err) {
		file, err := os.Create(db.db_file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		logInfo("%s created", db.db_file)

		db.open()

		logInfo("%s", "Reading sql file")
		statements, err := os.ReadFile("create_table.sql")
		if err != nil {
			os.Remove(db.db_file)
			log.Fatal(err)
		}

		logInfo("%s", "Creating table")
		statement, err := db.cnx.Prepare(string(statements))
		if err != nil {
			log.Fatal(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal(err)
		}
		logInfo("%s", "Table created")
	} else {
		logInfo("%s", "Database already exists")
	}
}

//
// CSV import into database, the csv filename contains the year
//
func (db *Database) import_csv_data(csv_file string) {
	year := strings.Split(path.Base(csv_file), "-")[0]
	logInfo("Removing old data for year %s from database if needed", year)
	statement, _ := db.cnx.Prepare(fmt.Sprintf("DELETE FROM dvf WHERE substr(id_mutation, 1, 4) = '%s'", year))
	_, err := statement.Exec()
	checkError(err)

	logInfo("Importing csv data from %s", csv_file)
	import_command := fmt.Sprintf(".import --csv --skip 1 %s dvf", csv_file)
	cmd := exec.Command("sqlite3", db.db_file, "-cmd", import_command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logError("cmd failed with %s", err.Error())
	}
}

//
// Clearing the distinct values cache is required if database have been updated
//
func (db *Database) clear_distinct_cache() {
	_, err := os.Stat(db.cache_file)
	if os.IsExist(err) {
		err = os.Remove(db.cache_file)
		if err != nil {
			logError("Could not delete db cache file - %s", err.Error())
		}
	}
}

//
// The distinct values cache file avoid scanning the database for distinct natures on every app start
// This function create the cache file as a json one
//
func (db *Database) init_distinct_values() {
	_, err := os.Stat(db.cache_file)
	if os.IsNotExist(err) {
		logInfo("%s", "Checking for distinct natures - Please wait")
		db.Types.Cultures = make([]string, 0)
		db.Types.Locaux = make([]string, 0)
		fields := []string{"nature_culture", "type_local", "nature_culture_speciale"}

		for i := 0; i < len(fields); i++ {
			// rows, err := db.cnx.Query("SELECT DISTINCT ? FROM dvf", fields[i]) // --> Problem !!!
			rows, err := db.cnx.Query(fmt.Sprintf("SELECT DISTINCT %s FROM dvf", fields[i]))
			if err != nil {
				logError("%s", err.Error())
				return
			}
			defer rows.Close()
			for rows.Next() {
				var s string
				err := rows.Scan(&s)
				if err != nil {
					log.Fatal(err)
				}
				println(s)
				if len(s) > 0 {
					switch fields[i] {
					case "nature_culture":
						db.Types.Cultures = append(db.Types.Cultures, s)
					case "nature_culture_speciale":
						db.Types.Cultures_speciales = append(db.Types.Cultures_speciales, s)
					case "type_local":
						db.Types.Locaux = append(db.Types.Locaux, s)
					}
				}
			}
		}
		b, _ := json.Marshal(db.Types)
		err = os.WriteFile(db.cache_file, b, 0644)
		checkError(err)
	} else {
		logInfo("%s", "Loading distinct natures from cache")
		file, _ := os.ReadFile(db.cache_file)
		_ = json.Unmarshal([]byte(file), &db.Types)
	}
}

func (db *Database) has_data_for_year(year int) bool {
	id_mutation := ""
	err := db.cnx.QueryRow("SELECT id_mutation FROM dvf WHERE substr(id_mutation, 1, 4) = ?", strconv.Itoa(year)).Scan(&id_mutation)
	return err != sql.ErrNoRows
}

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	db *Database
)

func leveledLog(level int, format string, v ...any) {
	level_str := map[int]string{0: "[INFO ]", 5: "[WARN ]", 10: "[ERROR]", 20: "[HTTP ]"}
	format = level_str[level] + " " + format
	log.Printf(format, v...)
}

func logInfo(format string, v ...any) {
	leveledLog(0, format, v...)
}

func logWarn(format string, v ...any) {
	leveledLog(5, format, v...)
}

func logError(format string, v ...any) {
	leveledLog(10, format, v...)
}

func checkError(err error) {
	if err != nil {
		logError("%s", err.Error())
	}
}

//
// File unzipping
//
func ungz_file(src string) error {
	logInfo("Unziping %s", src)

	gzipfile, err := os.Open(src)
	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		return err
	}
	defer reader.Close()

	newfilename := strings.TrimSuffix(src, ".gz")
	writer, err := os.Create(newfilename)
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}

//
// Kindly retry passed function attemps time
//
func retry_fn(attempts int, f func() error) error {
	var err error
	cont := true

	for cont {
		err = f()
		if err == nil {
			return nil
		}
		attempts--
		if attempts == 0 {
			break
		}
		time.Sleep(60 * time.Second)
	}
	return err
}

//
// Get info about a remote file
//
func get_remote_file_head(url string, attempts int) (*http.Response, error) {
	var resp *http.Response
	var err error

	err = retry_fn(attempts, func() error {
		logInfo("Checking remote file %s", url)
		resp, err = http.Head(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP status was %d", resp.StatusCode)
		} else {
			return nil
		}
	})
	return resp, err
}

//
// Download data file for a specific year
//
func download_file(url string, file_path string, attempts int) error {
	client := http.Client{}

	err := retry_fn(attempts, func() error {
		logInfo("Downloading %s", url)

		resp, err := client.Get(url) // Download data file
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP status was %d", resp.StatusCode)
		}

		file, err := os.Create(file_path)
		if err != nil {
			return err
		}
		defer file.Close()
		defer resp.Body.Close()

		logInfo("Writing file %s to disk", file_path)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		} else {
			return nil
		}
	})
	return err
}

//
// Check if a new or update remote file is available
// If so, or if local database update is needed then import data
//
func check_update_year(year int, attempts int) {
	var resp *http.Response
	var err error

	download_path := "./data/csv"
	base_url := "https://files.data.gouv.fr/geo-dvf/latest/csv"
	var url = fmt.Sprintf("%s/%d/full.csv.gz", base_url, year)

	resp, err = get_remote_file_head(url, attempts)
	if err != nil {
		logError("%s", err.Error())
		return
	}

	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	downloadSize := int64(size)

	var file_path = fmt.Sprintf("%s/%d-full.csv.gz", download_path, year)
	stats, err := os.Stat(file_path)
	if os.IsNotExist(err) || stats.Size() != downloadSize {
		err = download_file(url, file_path, attempts)
		if err != nil {
			logError("%s", err.Error())
			return
		} else {
			err = ungz_file(file_path)
			if err != nil {
				logError("%s", err.Error())
				return
			}
			db.import_csv_data(strings.TrimSuffix(file_path, ".gz"))
			db.clear_distinct_cache()
		}
	} else {
		if !db.has_data_for_year(year) {
			db.import_csv_data(strings.TrimSuffix(file_path, ".gz"))
			db.clear_distinct_cache()
		} else {
			logInfo("Found data for year %d", year)
		}
	}
}

//
// Call check update for every year from 2017
//
func update_local_data_files() {
	current_year := time.Now().Year()

	logInfo("%s", "Checking local files")
	for year := 2017; year < current_year; year++ {
		check_update_year(year, 3)
	}
	check_update_year(current_year, 1)
}

//
// DB Setup
//
func setup_db(dbfile string) *Database {
	db = new(Database)
	db.setup(dbfile)
	db.create_if_no_exists()
	db.open()
	db.init_distinct_values()
	return db
}

func main() {
	err := os.MkdirAll("./data/csv", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	logfile, err := os.Create("app.log")
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(w)
	defer logfile.Close()

	db = setup_db("./data/database.db")
	defer db.close()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		update_local_data_files()
		wg.Done()
	}()

	go func() {
		server_start()
		wg.Done()
	}()

	wg.Wait()
}

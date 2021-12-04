package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host    = "localhost"
	port    = 5432
	user    = "postgres"
	pasword = "password-here"
	dbname  = "phone_normalizer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, pasword)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)

	err = resetDB(db, dbname)
	checkErr(err)

	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	checkErr(err)

	defer db.Close()

	checkErr(db.Ping())
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func normalizeRegex(phone string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phone, "")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

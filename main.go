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
	pasword = ""
	dbname  = "phone_normalizer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pasword, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	defer db.Close()

	checkErr(createPhoneNumbersTable(db))

	phones := [...]string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, phone := range phones {
		id, err := insertPhone(db, phone)
		checkErr(err)
		fmt.Printf("id=%d\n", id)

	}

}

func insertPhone(db *sql.DB, phone string) (int, error) {
	var id int
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`

	_, err := db.Exec(statement)
	return err

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

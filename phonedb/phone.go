package phonedb

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var ErrNoRecord = errors.New("record not found")

type DB struct {
	db *sql.DB
}

// Phone represents the phone_numbers table in the database
type Phone struct {
	ID     int
	Number string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	phones := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, p := range phones {
		if _, err := insertPhone(db.db, p); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		phones = append(phones, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return phones, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	statement := `SELECT id, value FROM phone_numbers WHERE value=$1`
	row := db.db.QueryRow(statement, number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &p, nil
}
func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}

func (db *DB) DeletePhone(id int) error {
	_, err := db.db.Exec("DELETE FROM phone_numbers WHERE id=$1", id)
	return err
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func createDB(db *sql.DB, name string) error {
	statement := "CREATE DATABASE " + name
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	statement := "DROP DATABASE IF EXISTS " + name
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return createDB(db, name)
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

func insertPhone(db *sql.DB, phone string) (int, error) {
	var id int
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	statement := `SELECT value FROM phone_numbers WHERE id=$1`
	err := db.QueryRow(statement, id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

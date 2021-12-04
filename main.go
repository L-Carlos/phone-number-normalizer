package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/L-Carlos/phone-number-normalizer/phonedb"
	_ "github.com/lib/pq"
)

const (
	host       = "localhost"
	port       = 5432
	user       = "postgres"
	pasword    = ""
	dbname     = "phone_normalizer"
	driverName = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s  sslmode=disable",
		host, port, user, pasword)

	must(phonedb.Reset(driverName, psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate(driverName, psqlInfo))

	db, err := phonedb.Open(driverName, psqlInfo)
	must(err)
	defer db.Close()

	must(db.Seed())

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			if errors.Is(err, phonedb.ErrNoRecord) {
				p.Number = number
				must(db.UpdatePhone(&p))
			} else if existing != nil {
				must(db.DeletePhone(p.ID))
			} else {
				must(err)
			}
		} else {
			fmt.Println("No changes required")
		}
	}

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

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

package models

import (
	"database/sql"
	"fmt"
	"github.com/zjbztianya/gophercises/phone/conf"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func Init() error {
	var err error
	db, err = sql.Open(conf.Conf.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?)",
		conf.Conf.User,
		conf.Conf.Password,
		conf.Conf.Host,
		conf.Conf.DbName))
	if err != nil {
		return err
	}

	return createTable()
}

func Close() {
	db.Close()
}

type PhoneNumber struct {
	ID     int
	Number string
}

func createTable() error {
	statement := `
    CREATE TABLE IF NOT EXISTS phone_numbers (
      id int(10) unsigned NOT NULL AUTO_INCREMENT,
      value varchar(255) DEFAULT '',
      PRIMARY KEY (id))`
	_, err := db.Exec(statement)
	return err
}

func Seed() error {
	numbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range numbers {
		if err := AddPhoneNumber(number); err != nil {
			return err
		}
	}
	return nil
}

func AddPhoneNumber(number string) error {
	statement := `INSERT INTO phone_numbers(value) VALUES(?)`
	_, err := db.Exec(statement, number)
	return err
}

func GetPhoneNumbers() ([]PhoneNumber, error) {
	statement := `SELECT id,value FROM phone_numbers`
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []PhoneNumber
	for rows.Next() {
		var p PhoneNumber
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func UpdatePhoneNumber(p PhoneNumber) error {
	statement := `UPDATE phone_numbers SET value =? WHERE id=?`
	_, err := db.Exec(statement, p.Number, p.ID)
	return err
}
func DeletePhoneNumber(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id = ?`
	_, err := db.Exec(statement, id)
	return err
}

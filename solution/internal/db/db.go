package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Country struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Alpha2 string `json:"alpha2"`
	Alpha3 string `json:"alpha3"`
	Region string `json:"region"`
}

func scanCountry(row *sql.Rows) (Country, error) {
	c := Country{}
	err := row.Scan(&c.Id, &c.Name, &c.Alpha2, &c.Alpha3, &c.Region)
	return c, err
}

type DB struct {
	db    *sql.DB
	pgUrl string
}

func NewDB(pgUrl string) *DB {
	return &DB{pgUrl: pgUrl}
}

func (d *DB) Connect() error {
	pg, err := sql.Open("postgres", d.pgUrl)
	if err != nil {
		return err
	}
	for {
		err = pg.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}

	d.db = pg
	return nil

}

func (d *DB) GetCountries() ([]Country, error) {
	return getAll[Country](d.db, `SELECT * FROM countries`, scanCountry)
}

func (d *DB) GetCountriesOfRegion(region string) ([]Country, error) {
	return getAll[Country](d.db, fmt.Sprintf(`SELECT * FROM countries WHERE region='%v'`, region), scanCountry)
}

func getAll[T any](db *sql.DB, query string, scan func(*sql.Rows) (T, error)) ([]T, error) {
	res := []T{}
	rows, err := db.Query(query)
	if err != nil {
		return res, err
	}

	defer rows.Close()
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			continue
		}
		res = append(res, item)
	}

	return res, nil
}

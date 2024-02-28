package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Country struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Alpha2 string `json:"alpha2"`
	Alpha3 string `json:"alpha3"`
	Region string `json:"region"`
}

func scanCountry(row pgx.Rows) (Country, error) {
	c := Country{}
	err := row.Scan(&c.Id, &c.Name, &c.Alpha2, &c.Alpha3, &c.Region)
	return c, err
}

type DB struct {
	pool  *pgxpool.Pool
	pgUrl string
}

func NewDB(pgUrl string) *DB {
	return &DB{pgUrl: pgUrl}
}

func (d *DB) Connect() error {
	pool, err := pgxpool.New(context.Background(), d.pgUrl)
	if err != nil {
		return err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return err
	}
	d.pool = pool
	return nil

}

func (d *DB) GetCountries() ([]Country, error) {
	return getAll[Country](d.pool, `SELECT * FROM countries`, scanCountry)
}

func (d *DB) GetCountriesOfRegion(region string) ([]Country, error) {
	return getAll[Country](d.pool, fmt.Sprintf(`SELECT * FROM countries WHERE region='%v'`, region), scanCountry)
}

func getAll[T any](pool *pgxpool.Pool, query string, scan func(pgx.Rows) (T, error)) ([]T, error) {
	res := []T{}
	rows, err := pool.Query(context.Background(), query)
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

package service

import (
	"solution/internal/contract"

	"github.com/jackc/pgx/v5"
)

type Scanner[T any] interface {
	Scan(rows pgx.Row) (T, error)
}

type TableCreater interface {
	InitTable() string
}

type CountryDriver struct {
}

func (c CountryDriver) Scan(row pgx.Row) (contract.Country, error) {
	var country contract.Country
	err := row.Scan(&country.Id, &country.Name, &country.Alpha2, &country.Alpha3, &country.Region)
	return country, err
}

func (c CountryDriver) InitTable() string {
	return `CREATE TABLE IF NOT exists countries (
		id SERIAL PRIMARY KEY,
		name TEXT,
		alpha2 TEXT,
		alpha3 TEXT,
		region TEXT
	  );`
}

type SliceScanner[T any] struct {
	Length int
}

func (s SliceScanner[T]) Scan(row pgx.Row) ([]T, error) {
	sl := make([]T, s.Length)
	for i := 0; i < s.Length; i++ {
		err := row.Scan(&sl[i])
		if err != nil {
			return sl, err
		}
	}
	return sl, nil
}

func NewSliceScanner[T any](length int) SliceScanner[T] {
	return SliceScanner[T]{Length: length}
}

type StringScanner struct {
}

func (s StringScanner) Scan(row pgx.Row) (string, error) {
	var res string
	err := row.Scan(&res)
	return res, err
}

package service

import (
	"fmt"
	"solution/internal/contract"

	"github.com/jackc/pgx/v5"
)

type Scanner[T any] interface {
	Scan(rows pgx.Row) (T, error)
}

type TableCreater interface {
	InitTable() string
}

type Adder[T any] interface {
	Add(T) string
}

type CountryDriver struct {
}

func (c CountryDriver) Scan(row pgx.Row) (contract.Country, error) {
	var country contract.Country
	err := row.Scan(&country.Id, &country.Name, &country.Alpha2, &country.Alpha3, &country.Region)
	return country, err
}

func (c CountryDriver) InitTable() string {
	return `CREATE TABLE IF NOT EXISTS countries (
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

type UserDriver struct{}

func (u UserDriver) Scan(row pgx.Row) (contract.User, error) {
	user := contract.User{}
	err := row.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.CountryCode, &user.IsPublic, &user.Phone, &user.Image)
	return user, err
}

func (u UserDriver) InitTable() string {
	return `CREATE TABLE IF NOT EXISTS users (
		id TEXT NOT NULL PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		countryCode CHAR(2) NOT NULL,
		isPublic BOOL NOT NULL,
		phone TEXT UNIQUE,
		image TEXT
	);`
}

func (u UserDriver) Add(user contract.User) string {
	q := fmt.Sprintf(`INSERT INTO users VALUES ('%v', '%v', '%v', '%v', '%v', %v, '%v', '%v')`,
		user.Id, user.Login, user.Email, user.Password, user.CountryCode,
		user.IsPublic, user.Phone, user.Image,
	)
	return q
}

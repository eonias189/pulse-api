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

type Updater[T any] interface {
	Update(T, T) string
}
type Deleter[T any] interface {
	Delete(T) string
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
	err := row.Scan(&user.Login, &user.Email,
		&user.Password, &user.CountryCode, &user.IsPublic,
		&user.Phone, &user.Image, &user.PasswordChanged)
	return user, err
}

func (u UserDriver) InitTable() string {
	return `CREATE TABLE IF NOT EXISTS users (
		login TEXT NOT NULL UNIQUE PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		countryCode CHAR(2) NOT NULL,
		isPublic BOOL NOT NULL,
		phone TEXT,
		image TEXT,
		passwordChanged INTEGER
	);`
}

func (u UserDriver) Add(user contract.User) string {
	q := fmt.Sprintf(`INSERT INTO users VALUES ('%v', '%v', '%v', '%v', %v, '%v', '%v', %v)`,
		user.Login, user.Email, user.Password, user.CountryCode,
		user.IsPublic, user.Phone, user.Image, user.PasswordChanged,
	)
	return q
}

func (u UserDriver) Update(old, newUser contract.User) string {
	return fmt.Sprintf(`UPDATE users SET password='%v', countryCode='%v', isPublic=%v, phone='%v', image='%v', passwordChanged=%v WHERE login='%v'`,
		newUser.Password, newUser.CountryCode, newUser.IsPublic, newUser.Phone, newUser.Image, newUser.PasswordChanged, old.Login)
}

type RelationDriver struct {
}

func (r RelationDriver) InitTable() string {
	return `CREATE TABLE IF NOT EXISTS relations (
		id SERIAL PRIMARY KEY,
		senderLogin TEXT NOT NULL,
		accepterLogin TEXT NOT NULL,
		createTime INTEGER NOT NULL
	)`
}

func (r RelationDriver) Scan(row pgx.Row) (contract.Relation, error) {
	relation := contract.Relation{}
	err := row.Scan(&relation.Id, &relation.SenderLogin, &relation.AccepterLogin, &relation.CreateTime)
	return relation, err
}

func (r RelationDriver) Add(relation contract.Relation) string {
	return fmt.Sprintf(`INSERT INTO relations (senderLogin, accepterLogin, createTime) VALUES ('%v', '%v', %v)`,
		relation.SenderLogin, relation.AccepterLogin, relation.CreateTime)
}

func (r RelationDriver) Delete(relation contract.Relation) string {
	return fmt.Sprintf(`DELETE FROM relations WHERE id=%v`, relation.Id)
}

type AccepterRelation struct {
	contract.Relation
	contract.User
}

type AccepterRelationDriver struct{}

func (a AccepterRelationDriver) Scan(row pgx.Row) (AccepterRelation, error) {
	ar := AccepterRelation{}
	err := row.Scan(&ar.Id, &ar.SenderLogin, &ar.AccepterLogin, &ar.CreateTime,
		&ar.Login, &ar.Email, &ar.Password, &ar.CountryCode, &ar.IsPublic,
		&ar.Phone, &ar.Image, &ar.PasswordChanged,
	)
	return ar, err
}

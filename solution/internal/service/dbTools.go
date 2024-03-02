package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func QueryAll[T any, D Scanner[T]](driver D, pool *pgxpool.Pool, query string) ([]T, error) {
	res := []T{}
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return res, err
	}

	defer rows.Close()
	for rows.Next() {
		item, err := driver.Scan(rows)
		if err != nil {
			continue
		}
		res = append(res, item)
	}

	return res, nil
}

func QuerySingle[T any, D Scanner[T]](driver D, pool *pgxpool.Pool, query string) (T, error) {
	row := pool.QueryRow(context.Background(), query)
	return driver.Scan(row)
}

func InitTable[D TableCreater](driver D, pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), driver.InitTable())
	return err
}

func Insert[T any, A Adder[T]](adder A, pool *pgxpool.Pool, item T) error {
	_, err := pool.Exec(context.Background(), adder.Add(item))
	return err
}

func Update[T any, U Updater[T]](updater U, pool *pgxpool.Pool, old T, newItem T) error {
	_, err := pool.Exec(context.Background(), updater.Update(old, newItem))
	return err
}

package service

import (
	"context"
	"fmt"
	"solution/internal/contract"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool  *pgxpool.Pool
	pgUrl string
}

func New(pgUrl string) (*Service, error) {
	s := &Service{pgUrl: pgUrl}
	err := s.connect()
	if err != nil {
		return s, err
	}
	err = s.initTables()
	return s, err
}

func (s *Service) connect() error {
	pool, err := pgxpool.New(context.Background(), s.pgUrl)
	if err != nil {
		return err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return err
	}
	s.pool = pool
	return nil

}

func (s *Service) initTables() error {
	err := InitTable(CountryDriver{}, s.pool)
	return err
}

func (s *Service) GetCountries() ([]contract.Country, error) {
	return QueryAll(CountryDriver{}, s.pool, `SELECT * FROM countries`)
}

func (s *Service) GetCountriesOfRegion(region string) ([]contract.Country, error) {
	return QueryAll(CountryDriver{}, s.pool, fmt.Sprintf(`SELECT * FROM countries WHERE region='%v'`, region))
}

func (s *Service) GetRegions() ([]string, error) {
	return QueryAll(StringScanner{}, s.pool, `SELECT DISTINCT (region) FROM countries`)
}

func (s *Service) GetCountryByAlpha2(alpha2 string) (contract.Country, error) {
	query := fmt.Sprintf(`SELECT * FROM countries WHERE alpha2='%v'`, alpha2)
	return QuerySingle(CountryDriver{}, s.pool, query)
}

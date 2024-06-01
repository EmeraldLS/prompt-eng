package respository

import "database/sql"

type respository struct {
	DB *sql.DB
}

func NewRepository() (*respository, error) {
	connString := ""
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return &respository{
		DB: conn,
	}, nil
}

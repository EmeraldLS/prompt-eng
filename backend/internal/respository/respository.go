package respository

import "database/sql"

type respository struct {
	db *sql.DB
}

func NewRepository() (*respository, error) {
	connString := ""
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return &respository{
		db: conn,
	}, nil
}

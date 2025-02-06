package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"parecer-gen/pkg/date"
	"parecer-gen/pkg/parecer"
)

type SQLClient struct {
	DBClient *sql.DB
}

func NewSQLClient(connStr string) SQLClient {
	dbClient, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %s", err)
	}

	err = dbClient.Ping()
	if err != nil {
		log.Fatalf("Cannot ping the database: %s", err)
	}
	return SQLClient{DBClient: dbClient}
}

func (cli SQLClient) SaveParecer(data *parecer.Data) error {
	dateTime := date.StringToTime(data.Date)
	data.ID = fmt.Sprintf("parecer-%s", dateTime.Format("2-01-2006"))
	if _, err := cli.DBClient.ExecContext(context.Background(), `INSERT INTO pareceres (id, "user", creci, date, content) VALUES ($1, $2, $3, $4, $5)`, data.ID, data.User, data.Creci, data.Date, data.Content); err != nil {
		return fmt.Errorf("error inserting parecer into database: %w", err)
	}

	return nil
}

func (cli SQLClient) GetAllParecer() ([]parecer.Data, error) {
	rows, err := cli.DBClient.QueryContext(context.Background(), `SELECT id, "user", creci, date, content FROM pareceres`)
	if err != nil {
		return nil, fmt.Errorf("error getting pareceres from database: %w", err)
	}
	defer rows.Close()

	var res []parecer.Data

	for rows.Next() {
		var id, user, creci, date, content string
		if err := rows.Scan(&id, &user, &creci, &date, &content); err != nil {
			return nil, fmt.Errorf("error scanning pareceres from database: %w", err)
		}
		res = append(res, parecer.Data{
			ID:      id,
			User:    user,
			Creci:   creci,
			Date:    date,
			Content: content,
		})
	}

	return res, nil
}

func (cli SQLClient) GetParecer(id string) (*parecer.Data, error) {
	rows, err := cli.DBClient.QueryContext(context.Background(), `SELECT "user", creci, date, content FROM pareceres WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user, creci, date, content string

	for rows.Next() {
		if err := rows.Scan(&user, &creci, &date, &content); err != nil {
			return nil, err
		}
	}

	return &parecer.Data{
		User:    user,
		Creci:   creci,
		Date:    date,
		Content: content,
	}, nil
}

package parecer

import (
	"fmt"
	"time"
)

type Data struct {
	ID      string    `json:"id"`
	User    string    `json:"user"`
	Creci   string    `json:"creci"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

func NewData(user, creci, content string, date time.Time) (*Data, error) {

	if user == "" || creci == "" || content == "" {
		return nil, fmt.Errorf("missing data to generate parecer")
	}

	if date.IsZero() {
		date = time.Now()
	}

	return &Data{
		User:    user,
		Creci:   creci,
		Date:    date,
		Content: content,
	}, nil
}

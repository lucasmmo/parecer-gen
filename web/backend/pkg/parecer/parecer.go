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

func NewData(user, creci, content string) (*Data, error) {
	if user == "" || creci == "" || content == "" {
		return nil, fmt.Errorf("missing data to generate parecer")
	}

	return &Data{
		User:    user,
		Creci:   creci,
		Date:    time.Now(),
		Content: content,
	}, nil
}

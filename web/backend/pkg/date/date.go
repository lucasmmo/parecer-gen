package date

import (
	"fmt"
	"time"
)

var BRMonthMap = map[time.Month]string{
	time.January:   "Janeiro",
	time.February:  "Fevereiro",
	time.March:     "Março",
	time.April:     "Abril",
	time.May:       "Maio",
	time.June:      "Junho",
	time.July:      "Julho",
	time.August:    "Agosto",
	time.September: "Setembro",
	time.October:   "Outubro",
	time.November:  "Novembro",
	time.December:  "Dezembro",
}

func TimeToBRString(dateTime time.Time) string {
	// Definir o fuso horário UTC-3
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		location = time.FixedZone("UTC-3", -3*60*60)
	}
	dateTime = dateTime.In(location)

	year, month, day := dateTime.Date()
	monthStr := BRMonthMap[month]
	return fmt.Sprintf("%d de %s de %d", day, monthStr, year)
}

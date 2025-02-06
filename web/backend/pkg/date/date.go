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

func StringToTime(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

func TimeToBRString(dateTime time.Time) string {
	year, month, day := dateTime.Date()
	monthStr := BRMonthMap[month]
	return fmt.Sprintf("%d de %s de %d", day, monthStr, year)
}

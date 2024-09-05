package scraper

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/arewabolu/csvmanager"
)

func init() {
	if _, err := os.Stat(database()); os.IsNotExist(err) {
		err := os.Mkdir(database(), os.ModePerm)
		if err != nil {
			if err != os.ErrExist {
				panic(err)
			}
		}
	}
}

func database() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	return home + string(filepath.Separator) + "chromedp" + string(filepath.Separator)
}

func writeHeader(header []string, file string) error {
	nwFile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer nwFile.Close()
	writer := csvmanager.WriteFrame{
		File:    nwFile,
		Headers: header,
		Row:     true,
	}

	err = writer.WriteCSV()
	if err != nil {
		return err
	}
	return nil
}

func WriteBody(data [][]string, file string) error {
	nwFile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer nwFile.Close()
	writer2 := csvmanager.WriteFrame{
		File:   nwFile,
		Row:    false,
		Arrays: data,
	}
	err = writer2.WriteCSV()
	if err != nil {
		return err
	}
	return nil
}

func removeCountry(team string) string {
	if strings.Contains(team, "(") {
		split := strings.Split(team, "(")
		return split[0]
	}
	return team
}

func checkSeason(season string) (string, bool) {
	if strings.Contains(season, "-") {
		splitYear := strings.Split(season, "-")
		return splitYear[0], true
	} else {
		return season, false
	}
}

func getLastDate(file string, ch chan<- int) {
	frame, _ := csvmanager.ReadCsv(file, true)
	dates := frame.Col("date").String()
	lastDate := dates[len(dates)-1]
	var sep []string
	switch {
	case strings.Contains(lastDate, "."):
		sep = strings.Split(lastDate, ".")

	case strings.Contains(lastDate, "/"):
		sep = strings.Split(lastDate, "/")
	}
	day, err := strconv.Atoi(sep[0])
	if err != nil {
		panic("Invalid day format")
	}
	month, err := strconv.Atoi(sep[1])
	if err != nil {
		panic("Invalid month format")
	}

	ch <- day
	ch <- month
}

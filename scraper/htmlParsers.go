package scraper

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	gohaskell "github.com/arewabolu/GoHaskell"
	"github.com/gocolly/colly"
)

func generator(html, year string, summerStart bool) []match {
	buf := bytes.NewBufferString(html)
	query, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		panic(err)
	}
	matches := make([]match, 0)
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		panic(err)
	}
	if summerStart {
		year = fmt.Sprint(yearInt + 1)
	}
	query.Find("div.sportName").Each(func(i int, s *goquery.Selection) {
		s.Find("div.event__match").Each(func(i int, s *goquery.Selection) {
			tm := s.Find("div.event__time").Text()
			splitTM := strings.Split(tm, ".")
			day := splitTM[0]
			monthStr := splitTM[1]
			monthInt, _ := strconv.Atoi(monthStr)
			if len(matches) > 0 {
				splitDates := strings.Split(matches[len(matches)-1].Date, ".")
				oldMonthInt, _ := strconv.Atoi(splitDates[1])

				if oldMonthInt >= 1 && oldMonthInt <= 3 {
					if monthInt >= 10 && monthInt <= 12 {
						year = fmt.Sprint(yearInt)
					}
				}
			}

			nwTM := day + "." + splitTM[1] + "." + year
			match := match{
				Date:     nwTM,
				Hometeam: removeCountry(s.Find("div.event__participant--home").Text()),
				Awayteam: removeCountry(s.Find("div.event__participant--away").Text()),
				HomeGoal: s.Find("div.event__score--home").Text(),
				AwayGoal: s.Find("div.event__score--away").Text(), //
			}

			matches = append(matches, match)
		})
	})
	matches2 := gohaskell.Reverse(matches)

	return matches2
}

func Generator2(html, year string, summerStart bool) []halfMatch {
	buf := bytes.NewBufferString(html)
	query, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		panic(err)
	}
	matches := make([]halfMatch, 0)
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		panic(err)
	}
	if summerStart {
		year = fmt.Sprint(yearInt + 1)
	}
	query.Find("div.sportName").Each(func(i int, s *goquery.Selection) {
		s.Find("div.event__match").Each(func(i int, s *goquery.Selection) {
			tm := s.Find("div.event__time").Text()
			splitTM := strings.Split(tm, ".")
			day := splitTM[0]
			monthStr := splitTM[1]
			monthInt, _ := strconv.Atoi(monthStr)
			if len(matches) > 0 {
				splitDates := strings.Split(matches[len(matches)-1].Date, ".")
				oldMonthInt, _ := strconv.Atoi(splitDates[1])

				if oldMonthInt >= 1 && oldMonthInt <= 3 {
					if monthInt >= 10 && monthInt <= 12 {
						year = fmt.Sprint(yearInt)
					}
				}
			}

			nwTM := day + "." + splitTM[1] + "." + year

			fullMatch := halfMatch{
				match: match{
					Date:     nwTM,
					Hometeam: removeCountry(s.Find("div.event__participant--home").Text()),
					Awayteam: removeCountry(s.Find("div.event__participant--away").Text()),
					HomeGoal: s.Find("div.event__score--home").Text(),
					AwayGoal: s.Find("div.event__score--away").Text(), //
				},
				firstHalfHomeGoal: removeBrackets(s.Find("div.event__part--home").Text()),
				firstHalfAwayGoal: removeBrackets(s.Find("div.event__part--away").Text()),
			}
			fullMatch.setSecondhalfScore()

			//set second half scores
			matches = append(matches, fullMatch)
		})
	})
	matches2 := gohaskell.Reverse(matches)

	return matches2
}

func updater(league, year string) []match {
	//url := "https://www.soccer24.com/england/league-one-2020-2021/"
	chanel := make(chan int, 2)
	go getLastDate(league, chanel)

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	matchces := make([]match, 0)

	c := colly.NewCollector()
	c.WithTransport(t)
	oldDay := <-chanel
	oldMonth := <-chanel

	c.OnHTML("div.sportName", func(e *colly.HTMLElement) {
		query := e.DOM

		query.Each(func(i int, s *goquery.Selection) {
			s.Find("div.event__match").Each(func(i int, s *goquery.Selection) {
				tm := s.Find("div.event__time").Text()
				splitTM := strings.Split(tm, ".")
				dayInt, _ := strconv.Atoi(splitTM[0])
				monthInt, _ := strconv.Atoi(splitTM[1])

				switch {
				case dayInt > oldDay && monthInt >= oldMonth:
					nwTM := splitTM[0] + "." + splitTM[1] + "." + year
					match := match{
						Date:     nwTM,
						Hometeam: s.Find("div.event__participant--home").Text(),
						Awayteam: s.Find("div.event__participant--away").Text(),
						HomeGoal: s.Find("div.event__score--home").Text(),
						AwayGoal: s.Find("div.event__score--away").Text(), //
					}
					matchces = append(matchces, match)
				case dayInt != oldDay && monthInt > oldMonth:
					nwTM := splitTM[0] + "." + splitTM[1] + "." + year
					match := match{
						Date:     nwTM,
						Hometeam: s.Find("div.event__participant--home").Text(),
						Awayteam: s.Find("div.event__participant--away").Text(),
						HomeGoal: s.Find("div.event__score--home").Text(),
						AwayGoal: s.Find("div.event__score--away").Text(), //
					}
					matchces = append(matchces, match)
				}
			})
		})
		//
		//
		//
		//
		//	}

		//fmt.Println(s.Children().Html())

	})
	c.SetRequestTimeout(30 * time.Second)
	matchces2 := make([]match, len(matchces))

	// Define the OnResponse callback to wait for the full response
	c.Visit("file://./v3.html")
	for i := len(matchces) - 1; i >= 0; i-- {
		matchces2 = append(matchces2, matchces[i])
	}

	return matchces2
}

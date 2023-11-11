package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	gohaskell "github.com/arewabolu/GoHaskell"
	"github.com/arewabolu/csvmanager"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

type AppConfig struct {
	Cfg *Config
	Log *slog.Logger
}

type Config struct {
	sport   string
	country string
	league  string
	season  string
	save    string
}

func (a *AppConfig) genUrl() string {
	return fmt.Sprintf("https://www.flashscore.com/%s/%s/%s-%s/results/", a.Cfg.sport, a.Cfg.country, a.Cfg.league, a.Cfg.season)
}
func (c *Config) genUrl() string {
	return fmt.Sprintf("/%s-%s-%s.csv", c.country, c.league, c.season)
}

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	wd, err := os.Getwd()
	if err != nil {
		logger.Error("unable to get current working directory")
	}

	var cfg Config
	flag.StringVar(&cfg.country, "country", "", "set country for the league")
	flag.StringVar(&cfg.league, "league", "", "choose league to get match results")
	flag.StringVar(&cfg.sport, "sport", "", "available sports are:football,basketball,hockey,...")
	flag.StringVar(&cfg.season, "season", "", "set season to search for match results.\n Multi-year seasons should be of the form `start-end`\n e.g `2012-2022`")
	flag.StringVar(&cfg.save, "save", wd, "saves file as csv to specified directory, default value is the present working directory")
	flag.String("help", "", "show this help dialog")

	flag.Parse()

	appConfig := &AppConfig{
		Cfg: &cfg,
		Log: logger,
	}
	var matches []Match
	html := VisitSite(appConfig)
	if strings.Contains(cfg.season, "-") {
		splitYear := strings.Split(cfg.season, "-")
		matches = generator(html, splitYear[0], true)
	} else {
		matches = generator(html, cfg.season, false)
	}
	writeHeader([]string{"date", "homeTeam", "awayTeam", "homeScore", "awayScore"}, fmt.Sprintf("%s/%s", wd, cfg.genUrl()))
	writer(matches, fmt.Sprintf("%s/%s", wd, cfg.genUrl()))
}

func VisitSite(appConfig *AppConfig) string {
	showMoreAction :=
		`
	 (async function() {
		 while (true) {
   			try {
        		await new Promise((resolve) => setTimeout(resolve, 1500));
       			const element = document.querySelector(
          		"a.event__more.event__more--static"
        	);
        	element.scrollIntoView();
        	element.click();
    		} catch (error) {
      			break;
    		}
  		}
	})();
		`
	newCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	newCtx, cancel = chromedp.NewContext(newCtx)
	defer cancel()

	var html string
	var isElementPresent bool
	err := chromedp.Run(newCtx,
		// navigate to a page,
		chromedp.Navigate(appConfig.genUrl()),
		// wait for footer element i.e, page is loaded
		chromedp.WaitVisible(`body > footer`),
	)
	if err != nil {
		appConfig.Log.DebugContext(newCtx, fmt.Sprintf("%s could not load", appConfig.genUrl()))
	}

	//evaluate javascript scroll and click
	err = chromedp.Run(
		newCtx,
		chromedp.Evaluate(showMoreAction, nil),
	)
	if err != nil {
		appConfig.Log.Error(err.Error())
	}

	//wait for action to complete since async is not supported
	time.Sleep(7 * time.Second)

	err = chromedp.Run(newCtx, chromedp.Evaluate(`!!document.querySelector("a.event__more.event__more--static")`, &isElementPresent))
	if err != nil {
		appConfig.Log.Error(err.Error())
	}

	err = chromedp.Run(newCtx, chromedp.InnerHTML(`.leagues--static.event--leagues.results`, &html, chromedp.AtLeast(1)))
	if err != nil {
		log.Fatal(err)
	}

	return html
}

// Getter visits the Url,
// gathers and edits html content.
// Getter then opens the file,
// writes the html content to the file,
// then closes the file

type Match struct {
	Date     string
	Hometeam string
	Awayteam string
	HomeGoal string
	AwayGoal string
}

func confirmandQuit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeHeader(header []string, file string) {
	nwFile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	confirmandQuit(err)
	defer nwFile.Close()
	writer := csvmanager.WriteFrame{
		File:    nwFile,
		Headers: header,
		Row:     true,
	}

	writer.WriteCSV()
}

func WriteBody(data []string, file string) error {
	nwFile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	confirmandQuit(err)
	defer nwFile.Close()
	writer2 := csvmanager.WriteFrame{
		File:   nwFile,
		Row:    true,
		Arrays: [][]string{data},
	}
	err = writer2.WriteCSV()
	return err
}

func getLastDate(file string, ch chan<- int) {
	frame, _ := csvmanager.ReadCsv(file, 0755, true)
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

// main.container__liveTableWrapper>div.container__livetable
func updater(league, year string) []Match {
	//url := "https://www.soccer24.com/england/league-one-2020-2021/"
	chanel := make(chan int, 2)
	go getLastDate(league, chanel)

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	matchces := make([]Match, 0)

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
					match := Match{
						Date:     nwTM,
						Hometeam: s.Find("div.event__participant--home").Text(),
						Awayteam: s.Find("div.event__participant--away").Text(),
						HomeGoal: s.Find("div.event__score--home").Text(),
						AwayGoal: s.Find("div.event__score--away").Text(), //
					}
					matchces = append(matchces, match)
				case dayInt != oldDay && monthInt > oldMonth:
					nwTM := splitTM[0] + "." + splitTM[1] + "." + year
					match := Match{
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
	matchces2 := make([]Match, len(matchces))

	// Define the OnResponse callback to wait for the full response
	c.Visit("file://./v3.html")
	for i := len(matchces) - 1; i >= 0; i-- {
		matchces2 = append(matchces2, matchces[i])
	}

	//fmt.Println(matchces[0])
	//fmt.Println(matchces2[len(matchces2)-1])
	//csvmanager.PrependRow("./EngLeague2.csv", 0755, true, []string{match.Date, strings.ToUpper(strings.TrimSpace(match.Hometeam)), strings.ToUpper(strings.TrimSpace(match.Awayteam)), strings.ToUpper(strings.TrimSpace(match.HomeGoal)), strings.ToUpper(strings.TrimSpace(match.AwayGoal))})

	return matchces2
}

func generator(html, year string, summerStart bool) []Match {
	buf := bytes.NewBufferString(html)
	query, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		panic(err)
	}
	matches := make([]Match, 0)
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
			match := Match{
				Date:     nwTM,
				Hometeam: s.Find("div.event__participant--home").Text(),
				Awayteam: s.Find("div.event__participant--away").Text(),
				HomeGoal: s.Find("div.event__score--home").Text(),
				AwayGoal: s.Find("div.event__score--away").Text(), //
			}

			matches = append(matches, match)
		})
	})
	matches2 := gohaskell.Reverse(matches)

	return matches2
}

func writer(matches []Match, league string) error {
	for _, match := range matches {
		err := WriteBody([]string{match.Date, strings.ToUpper(strings.TrimSpace(match.Hometeam)), strings.ToUpper(strings.TrimSpace(match.Awayteam)), strings.ToUpper(strings.TrimSpace(match.HomeGoal)), strings.ToUpper(strings.TrimSpace(match.AwayGoal))}, league)
		if err != nil {
			return err
		}
	}

	return nil
}

//chromedp.WaitSelected()
/*	//
	var count int
	for {
		// Check if the element is still present on the refreshed page
		var isElementPresent bool
		if err := chromedp.Run(ctx, chromedp.Evaluate(`!!document.querySelector("a.event__more.event__more--static")`, &isElementPresent)); err != nil {
			log.Fatal(err)
		}
		if !isElementPresent {
			break
		}
		count++

		// Perform a click on the element
		//if err := chromedp.Run(ctx,
		//	chromedp.ScrollIntoView(showMore, chromedp.ByQuery),
		//	chromedp.Click(showMore, chromedp.NodeNotVisible)); err != nil {
		//	log.Fatal(err)
		//}
		//
		//// Use WaitAction to wait for a specific action to complete
		//if err := chromedp.Run(ctx, chromedp.WaitNotPresent(showMore)); err != nil {
		//	log.Fatal(err)
		//}

		// Wait for the page to refresh
		//if err := chromedp.Run(ctx, chromedp.WaitVisible(showMore, chromedp.ByQuery)); err != nil {
		//	log.Fatal(err)
		//}

		// If the element is not present, break out of the loop

		// You can add a delay between clicks if needed

	}
	fmt.Println(count)
	err = chromedp.InnerHTML(`.leagues--static.event--leagues.results`, &html, chromedp.AtLeast(1)).Do(ctx)
	if err != nil {
		panic(err)
	}
*/

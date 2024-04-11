package scraper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/arewabolu/Flashscore-Scraper/config"
	gohaskell "github.com/arewabolu/GoHaskell"
	"github.com/chromedp/chromedp"
)

const (
	PossessionIndex       int = 1
	ShotsIndex            int = 2
	ShotsOTIndex          int = 3
	ShotsOfffIndex        int = 4
	CornersIndex          int = 7
	SavesIndex            int = 10
	DangerousAttacksIndex int = 14
)

// use fixtures to setup url to either visit the results or fixtures page
func VisitSite(appConfig *config.Config, fixtures bool) (string, error) {
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
	newCtx, cancel := context.WithTimeout(context.Background(), time.Duration(appConfig.TimeOut)*time.Second)
	defer cancel()

	newCtx, cancel = chromedp.NewContext(newCtx)
	defer cancel()

	var html string
	var isElementPresent bool
	url := appConfig.GenUrl(fixtures)

	err := chromedp.Run(newCtx,
		// navigate to a page,
		chromedp.Navigate(url),
		// wait for footer element i.e, page is loaded
		chromedp.WaitVisible(`body>footer`),
	)
	if err != nil {
		return "", errors.Join(err, fmt.Errorf("%s could not load", url))
	}
	//evaluate javascript scroll and click
	err = chromedp.Run(
		newCtx,
		chromedp.Evaluate(showMoreAction, nil),
	)
	if err != nil {
		return "", errors.Join(err, errors.New("unable to evaluate show more action"))
	}

	//wait for action to complete since async is not supported
	time.Sleep(10 * time.Second)
	err = chromedp.Run(newCtx, chromedp.Evaluate(`document.querySelector("a.event__more.event__more--static") !== null`, &isElementPresent))
	if err != nil {
		return "", err
	}

	err = chromedp.Run(newCtx, chromedp.InnerHTML(`div.leagues--static`, &html, chromedp.AtLeast(1)))
	if err != nil {
		return "", err
	}

	return html, nil
}

func GetBasicMatchInfo(appConfig *config.Config, fixtures bool) error {
	html, err := VisitSite(appConfig, fixtures)
	if err != nil {
		return err
	}
	var matches []match
	if strings.Contains(appConfig.Season, "-") {
		splitYear := strings.Split(appConfig.Season, "-")
		matches = Generator(html, splitYear[0], true)
	} else {
		matches = Generator(html, appConfig.Season, false)
	}

	header := [5]string{"date", "homeTeam", "awayTeam", "homeScore", "awayScore"}
	err = writeHeader(header[:], fmt.Sprintf("%s/%s", appConfig.Path, appConfig.GenFilePath()))
	if err != nil {
		return err
	}
	reverseMatches := gohaskell.Reverse(matches)
	err = WriteBody(stringifyMatch(reverseMatches), fmt.Sprintf("%s/%s", appConfig.Path, appConfig.GenFilePath()))
	if err != nil {
		return err
	}
	return nil
}

// use fixtures to setup url to either visit the results or fixtures page
func GetHalfMatchInfo(appConfig *config.Config, fixtures bool) error {
	html, err := VisitSite(appConfig, fixtures)
	if err != nil {
		return err
	}
	year, val := checkSeason(appConfig.Season)
	matches := Generator2(html, year, val)
	header := [9]string{"date", "homeTeam", "awayTeam", "homeScore 1st half", "awayScore 1st half", "homeScore 2nd half", "awayScore 2nd half", "homeScore", "awayScore"}
	err = writeHeader(header[:], fmt.Sprintf("%s/%s", appConfig.Path, appConfig.GenFilePath()))
	if err != nil {
		return err
	}

	reverseMatches := gohaskell.Reverse(matches)
	err = WriteBody(stringifyMatch2(reverseMatches), fmt.Sprintf("%s/%s", appConfig.Path, appConfig.GenFilePath()))
	if err != nil {
		return err
	}
	return nil
}

/*
func GetMatchIds(url string) []string {
	var divIDs []string
	getDivIds := `(() => {
		let divs = document.querySelectorAll('div');
		return Array.from(divs).map(div => div.id);
	})()`

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
	newCtx, cancel := context.WithTimeout(context.Background(), 28*time.Second)
	defer cancel()

	newCtx, cancel = chromedp.NewContext(newCtx)
	defer cancel()

	//var html string
	var isElementPresent bool
	err := chromedp.Run(newCtx,
		// navigate to a page,
		chromedp.Navigate(url),
		// wait for footer element i.e, page is loaded
		chromedp.WaitVisible(`body > footer`),
	)
	if err != nil {
		panic(fmt.Sprintf("%s could not load %s", url, err.Error()))
	}

	//evaluate javascript scroll and click
	err = chromedp.Run(
		newCtx,
		chromedp.Evaluate(showMoreAction, nil),
	)
	if err != nil {
		panic(err)
	}

	//wait for action to complete since async is not supported
	time.Sleep(7 * time.Second)

	err = chromedp.Run(newCtx, chromedp.Evaluate(`!!document.querySelector("a.event__more.event__more--static")`, &isElementPresent))
	if err != nil {
		log.Panic(err.Error())
	}

	err = chromedp.Run(newCtx, chromedp.Evaluate(getDivIds, &divIDs))
	if err != nil {
		log.Fatal(err)
	}

	divIds := make([]string, 0)
	for _, id := range divIDs {
		if strings.HasPrefix(id, "g_1_") {
			divIds = append(divIds, strings.TrimPrefix(id, "g_1_"))
		}
	}
	return divIds
}

func Flashscore(matchId string) string {
	return fmt.Sprintf("https://www.flashscore.com/match/%s/#/match-summary", matchId)
}

//https://www.flashscore.com/match/ns2qTfdf/#/match-summary

func FlashscoreStat(matchId string) string {
	return fmt.Sprintf("https://www.flashscore.com/match/%s/#/match-summary/match-statistics/0", matchId)
}

func GetMatchData(matchId string, urlFunc func(string) string) string {
	url := urlFunc(matchId)
	idCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	idCtx, cancel = chromedp.NewContext(idCtx)
	defer cancel()
	var html string
	err := chromedp.Run(idCtx, chromedp.Navigate(url), chromedp.InnerHTML(`.container__detailInner`, &html, chromedp.AtLeast(1)))
	if err != nil {
		panic("error getting match data with id: " + matchId)
	}
	return html
}
*/

func generateDOM(html string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func GetBasicMatchData(dom *goquery.Document) MatchInfo {
	matchInfo := &MatchInfo{}
	matchInfo.Date = strings.SplitAfter(dom.Find("div.duelParticipant__startTime").Text(), " ")[0]
	matchInfo.Hometeam = removeCountry(dom.Find("div.duelParticipant__home").Find("a.participant__participantName").Text())
	matchInfo.Awayteam = removeCountry(dom.Find("div.duelParticipant__away").Find("a.participant__participantName").Text())

	dom.Find("div.duelParticipant__score").Find("div.detailScore__wrapper").Each(
		func(i int, s *goquery.Selection) {
			scores := s.Find("span:not(span.detailScore__divider)").Text()
			matchInfo.HomeGoal = strings.Split(scores, "")[0]
			matchInfo.AwayGoal = strings.Split(scores, "")[1]
		},
	)
	return *matchInfo
}

func parseDom(matchInfo *MatchInfo, dom *goquery.Document) MatchInfo {
	// used to pass collect exttra data for match
	homeDatas := make([]string, 0, 15)
	awayDatas := make([]string, 0, 15)
	dom.Find("div.section").Each(
		func(i int, s *goquery.Selection) {
			s.Find("div._category_1gfjz_16").Each(
				func(i int, s *goquery.Selection) {
					homeDatas = append(homeDatas, s.ChildrenFiltered("div._homeValue_v26p1_10").Text())
					awayDatas = append(awayDatas, s.ChildrenFiltered("div._awayValue_v26p1_14").Text())
				},
			)
		},
	)
	if len(homeDatas) == 0 && len(awayDatas) == 0 {
		return *matchInfo
	}

	matchInfo.HomeMatchData = makeMatchData(homeDatas)
	matchInfo.AwayMatchData = makeMatchData(awayDatas)
	return *matchInfo
}

func makeMatchData(datas []string) MatchData {

	//	if len(datas) ==  {
	return MatchData{
		Possession:       datas[PossessionIndex],
		Saves:            datas[SavesIndex],
		Shots:            datas[ShotsIndex],
		ShotsOT:          datas[ShotsOTIndex],
		ShotsOfff:        datas[ShotsOfffIndex],
		Corners:          datas[CornersIndex],
		DangerousAttacks: datas[DangerousAttacksIndex],
	}
}

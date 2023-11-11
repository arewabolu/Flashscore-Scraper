package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arewabolu/Flashscore-Scraper/config"
	"github.com/chromedp/chromedp"
)

func VisitSite(appConfig *config.AppConfig) string {
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
	newCtx, cancel := context.WithTimeout(context.Background(), time.Duration(appConfig.Cfg.TimeOut)*time.Second)
	defer cancel()

	newCtx, cancel = chromedp.NewContext(newCtx)
	defer cancel()

	var html string
	var isElementPresent bool
	err := chromedp.Run(newCtx,
		// navigate to a page,
		chromedp.Navigate(appConfig.GenUrl()),
		// wait for footer element i.e, page is loaded
		chromedp.WaitVisible(`body > footer`),
	)
	if err != nil {
		appConfig.Log.DebugContext(newCtx, fmt.Sprintf("%s could not load", appConfig.GenUrl()))
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

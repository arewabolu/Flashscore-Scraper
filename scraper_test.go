package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/arewabolu/Flashscore-Scraper/config"
	"github.com/arewabolu/Flashscore-Scraper/scraper"
	gohaskell "github.com/arewabolu/GoHaskell"
	"github.com/gocolly/colly"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
)

var leaguefile = "/home/arthemis/bettor/database/portugal_u23.csv"

func TestHeader(t *testing.T) {
	//writeHeader([]string{"date", "homeTeam", "awayTeam", "homeScore", "awayScore"}, leaguefile)
}

func TestUpdater(t *testing.T) {
	//matches := updater(leaguefile, "2023")
	////t.Error(matches[0], matches[len(matches)-1].Date)
	//writer(matches, leaguefile)
}

func TestGptScraper(t *testing.T) {
	c := colly.NewCollector()
	c.OnXML("/html/body/div[4]/div[1]/div/div/main/div[4]/div[2]/div[1]/div[1]/div", func(e *colly.XMLElement) {
		// Extract and print the content of the selected element
		elementContent := e.Text
		elementContent = trimWhitespace(elementContent)
		t.Error(elementContent)
	})

	c.Visit("https://www.soccer24.com/netherlands/eredivisie-2012-2013/results/")
	c.Wait()

}

// Helper function to trim whitespace from a string
func trimWhitespace(s string) string {
	return strings.TrimSpace(s)
}

func TestGetDate(t *testing.T) {
	chanel := make(chan int, 2)
	//go getLastDate(leaguefile, chanel)
	t.Error(<-chanel)
	t.Error(<-chanel)

}

func TestGetter(t *testing.T) {
	resp, err := http.Get("https://www.soccer24.com/belgium/jupiler-pro-league-2016-2017/results/#")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read and process chunks
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // End of response
			}
			fmt.Println("Error reading chunk:", err)
			break
		}
		// Process the chunk (e.g., print it)
		t.Error(string(buffer[:n]))
	}
}

/*c.OnHTML("a[href='your-show-more-options-link']", func(e *colly.HTMLElement) {
	// Click the link to load more options
	link := e.Attr("href")
	c.Visit(e.Request.AbsoluteURL(link))
})
c.OnHTML("your-selector-for-additional-content", func(e *colly.HTMLElement) {
	// Extract data from the fully loaded page
})*/
//Find("div.event").Find("leagues--static").Find("div.sportName").Size()
//		t.Error(sz)

func TestTemp(t *testing.T) {
	tempFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer tempFile.Close()
	t.Error(tempFile.Name())
}

func TestViewMore(t *testing.T) {
	trans := &http.Transport{}
	trans.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))

	c := colly.NewCollector()
	c.WithTransport(trans)
	c.OnHTML("div.container__fsbody", func(h *colly.HTMLElement) {
		fl, _ := os.Create("test.html")
		fl.Write(h.Response.Body)
		size := h.DOM.Find("div.event__more").Size()
		t.Error(size)
	})

	c.Visit("https://www.soccer24.com/belgium/jupiler-pro-league-2016-2017/results/")
}

func TestView(t *testing.T) {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Chrome())
	bow.AddRequestHeader("Connection", "keep-alive")
	bow.SetAttribute(browser.SendReferer, true)

	err := bow.Open("https://www.flashscore.com/football/andorra/andorra-cup-2021-2022/results/")
	if err != nil {
		panic(err)
	}
	size := bow.Dom().Text()
	fl, _ := os.Create("test2.html")
	fl.Write([]byte(size))

	found := bow.Dom().Find("div#live-table").Size()
	t.Error(found)

}

func TestPiped(t *testing.T) {
	buf := make([]byte, 1000)
	file, _ := os.Open("")
	file.Read(buf)

}

func TestHeadless(t *testing.T) {
	//html := scraper.GetMatchData("fmcnUxqt")
	//dom := scraper.GenerateDOM(html)
	//data := scraper.ParseDom(dom)
	//t.Error(data)
}

func Test2HalvesData(t *testing.T) {
	cfg := config.NewConfig()
	cfg.SetCountry("england")
	cfg.SetSport("football")
	cfg.SetLeague("premier-league")
	cfg.SetSeason("2022-2023")
	cfg.TimeOut = 100
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	appCfg := config.NewAppConfig(*cfg)
	appCfg.Log = logger
	html := scraper.VisitSite(appCfg)
	splitYear := strings.Split(cfg.Season, "-")
	matches := scraper.Generator2(html, splitYear[0], true)
	reverse := gohaskell.Reverse(matches)
	t.Error(reverse[10].String())
	//dom := scraper.GenerateDOM(html)
	//data := scraper.ParseDom(dom)
	//t.Error(data)
}

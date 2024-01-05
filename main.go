package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/arewabolu/Flashscore-Scraper/config"
	"github.com/arewabolu/Flashscore-Scraper/scraper"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	wd, err := os.Getwd()
	if err != nil {
		logger.Error("unable to get current working directory")
	}

	var detailed bool

	cfg := config.NewConfig()
	flag.StringVar(&cfg.Country, "c", "", "set country for the league")
	flag.StringVar(&cfg.League, "league", "", "choose league to get match results")
	flag.StringVar(&cfg.Season, "season", "", "set season to search for match results.\n Multi-year seasons should be of the form `start-end`\n e.g `2012-2022`")
	flag.StringVar(&cfg.Path, "path", wd, "saves file as csv to specified directory, default value is the present directory")
	flag.StringVar(&cfg.FileName, "f", cfg.GenFilePath(), "name of the file to be saved,defaults to country-league-season.csv. I.e. the country-league-season must be specified")
	flag.UintVar(&cfg.TimeOut, "t", 30, "timeout (in seconds) for request. default 30")
	flag.BoolVar(&detailed, "d", false, "")
	flag.String("h", "", "show this help dialog")

	flag.Parse()

	appConfig := &config.AppConfig{
		Cfg: cfg,
		Log: logger,
	}

	switch {
	case detailed:
		scraper.GetMatchesWithExtraData(appConfig)
	default:
		scraper.GetBasicMatchInfo(appConfig)
	}

}

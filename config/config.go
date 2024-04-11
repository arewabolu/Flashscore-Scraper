package config

import (
	"fmt"
)

type Config struct {
	FileName string
	Sport    string
	Country  string
	League   string
	Season   string
	Path     string
	TimeOut  uint
}

func (c *Config) GenFilePath() string {
	if c.FileName != "" {
		return c.FileName
	}
	return fmt.Sprintf("%s-%s-%s.csv", c.Country, c.League, c.Season)
}

// use fixtures to setup url to either visit the results or fixtures page
func (c *Config) GenUrl(fixtures bool) string {
	if fixtures {
		return fmt.Sprintf("https://www.flashscore.com/%s/%s/%s/fixtures/", c.Sport, c.Country, c.League)
	}
	return fmt.Sprintf("https://www.flashscore.com/%s/%s/%s-%s/results/", c.Sport, c.Country, c.League, c.Season)
}

func NewConfig() *Config {
	return &Config{Sport: "football"}
}

func (c *Config) SetSport(sport string) {
	c.Sport = sport
}

func (c *Config) SetCountry(country string) {
	c.Country = country
}

func (c *Config) SetLeague(league string) {
	c.League = league
}

func (c *Config) SetSeason(season string) {
	c.Season = season
}

package config

import (
	"fmt"
	"log/slog"
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

type AppConfig struct {
	Cfg *Config
	Log *slog.Logger
}

func (a *AppConfig) GenUrl() string {
	return fmt.Sprintf("https://www.flashscore.com/%s/%s/%s-%s/results/", a.Cfg.Sport, a.Cfg.Country, a.Cfg.League, a.Cfg.Season)
}

func NewAppConfig(cfg Config) *AppConfig {
	return &AppConfig{Cfg: &cfg}
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

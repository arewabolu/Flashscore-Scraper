package scraper

import "strings"

type match struct {
	Date     string
	Hometeam string
	Awayteam string
	HomeGoal string
	AwayGoal string
}

func (m *match) String() []string {
	return []string{m.Date, strings.ToUpper(m.Hometeam), strings.ToUpper(m.Awayteam), m.HomeGoal, m.AwayGoal}
}

func stringifyMatch(matches []match) [][]string {
	strMatches := make([][]string, len(matches))
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		if match.HomeGoal == "-" || match.AwayGoal == "-" {
			continue
		}
		strMatches[i] = match.String()
	}
	return strMatches
}

package scraper

import "strings"

type MatchData struct {
	Possession       string
	Shots            string
	ShotsOT          string
	ShotsOfff        string
	Corners          string
	Saves            string
	DangerousAttacks string
}

type MatchInfo struct {
	Date          string
	Hometeam      string
	Awayteam      string
	HomeGoal      string
	AwayGoal      string
	HomeMatchData MatchData
	AwayMatchData MatchData
}

func (m *MatchInfo) String() []string {
	return []string{m.Date,
		removeCountry(strings.ToUpper(m.Hometeam)), removeCountry(strings.ToUpper(m.Awayteam)),
		m.HomeGoal, m.AwayGoal,
		m.HomeMatchData.Possession, m.HomeMatchData.Shots,
		m.HomeMatchData.ShotsOT, m.HomeMatchData.ShotsOfff,
		m.HomeMatchData.Corners, m.HomeMatchData.Saves, m.HomeMatchData.DangerousAttacks,
		m.AwayMatchData.Possession, m.AwayMatchData.Shots,
		m.AwayMatchData.ShotsOT, m.AwayMatchData.ShotsOfff,
		m.AwayMatchData.Corners, m.AwayMatchData.Saves, m.AwayMatchData.DangerousAttacks}
}

func stringifyMatchInfo(matches []MatchInfo) [][]string {
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

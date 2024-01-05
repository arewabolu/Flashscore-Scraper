package scraper

import (
	"strconv"
	"strings"
)

type halfMatch struct {
	match
	firstHalfHomeGoal  string
	secondHalfHomeGoal string
	firstHalfAwayGoal  string
	secondHalfAwayGoal string
}

func (m halfMatch) String() []string {
	return []string{m.Date, strings.ToUpper(m.Hometeam), strings.ToUpper(m.Awayteam), m.firstHalfHomeGoal, m.firstHalfAwayGoal, m.secondHalfHomeGoal, m.secondHalfAwayGoal, m.HomeGoal, m.AwayGoal}
}

func removeBrackets(team string) string {
	if strings.Contains(team, "(") {
		split := strings.Split(team, "")
		return split[1]
	}
	return team
}

func stringifyMatch2(matches []halfMatch) [][]string {
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

func (m *halfMatch) setSecondhalfScore() {
	fullGameHome, err := strconv.Atoi(m.HomeGoal)
	if err != nil {
		panic(err)
	}
	fullGameAway, err := strconv.Atoi(m.AwayGoal)
	if err != nil {
		panic(err)
	}
	halfGameHome, err := strconv.Atoi(m.firstHalfHomeGoal)
	if err != nil {
		panic(err)
	}
	halfGameAway, err := strconv.Atoi(m.firstHalfAwayGoal)
	if err != nil {
		panic(err)
	}
	m.secondHalfHomeGoal = strconv.Itoa(fullGameHome - halfGameHome)
	m.secondHalfAwayGoal = strconv.Itoa(fullGameAway - halfGameAway)
}

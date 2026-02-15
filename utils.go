package mango

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type leaderboardItem struct {
	Rank     int
	Username string
	Mana     int
	Traders  int
}

// LeadUser represents a member of one of the Manifold leaderboards, and their associated user information
type LeadUser struct {
	Rank     int
	Username string
	Mana     int
	Traders  int
	User     User
}

// UsernameType represents a special category of users on Manifold
type UsernameType string

const (
	Bot   UsernameType = "BOT"
	Check UsernameType = "CHECK"
	Core  UsernameType = "CORE"
)

// LeaderType represents a Manifold leaderboard which an item can be part of
type LeaderType int

const (
	Trader LeaderType = iota
	Creator
	// Referrer
)

// LeaderPeriod represents a time period that a Manifold leaderboard can represent
type LeaderPeriod int

const (
	Daily LeaderPeriod = iota
	Weekly
	Monthly
	All
)

// KellyBet returns the percentage of your bankroll that you ought to bet, for a given computed probability and payout.
//
// Payout is in decimal odds.
//
// See [the Kelly Criterion Wikipedia page] for more details.
//
// [the Kelly Criterion Wikipedia page]: https://en.wikipedia.org/wiki/Kelly_criterion
func KellyBet(prob, payout float64) float64 {
	// Ensure prob is within valid range (0 to 1)
	if prob <= 0 || prob >= 1 {
		return 0.0
	}

	// Compute the Kelly bet using the formula: f* = p - ((1-p) / b)
	k := prob - ((1 - prob) / (payout / 2))

	// Round the Kelly bet to 2 decimal places
	k = math.Round(k*100) / 100

	return k
}

func scrapeConstants(url string) (string, error) {
	var textContent strings.Builder

	c := colly.NewCollector()

	c.OnHTML("table.highlight", func(e *colly.HTMLElement) {
		e.ForEach("td.blob-code", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				textContent.WriteString(text)
				textContent.WriteString("\n")
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return textContent.String(), nil
}

func scrapeLeaderboards(url string) ([]leaderboardItem, error) {
	leaders := make([]leaderboardItem, 160)

	c := colly.NewCollector()

	i := 0

	c.OnHTML("main div", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, el *colly.HTMLElement) {
			el.ForEachWithBreak("tbody tr", func(_ int, ele *colly.HTMLElement) bool {
				if i > 159 {
					return false
				}
				tds := ele.ChildTexts("td")
				rank, _ := strconv.Atoi(tds[0])
				uname := ele.ChildAttr("td a", "href")
				mana, traders := 0, 0
				if tds[2][0] == 225 {
					mana, _ = strconv.Atoi(strings.Replace(tds[2][3:], ",", "", -1))
				} else {
					traders, _ = strconv.Atoi(strings.Replace(tds[2], ",", "", -1))
				}
				l := leaderboardItem{
					Rank:     rank,
					Username: uname[1:],
					Mana:     mana,
					Traders:  traders,
				}

				leaders[i] = l

				i++

				return i < 160
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return leaders, nil
}

func getUsernames(t UsernameType, text string) []string {
	// Find correct USERNAMES list
	startPattern := fmt.Sprintf(`export const %v_USERNAMES = [`, t)
	endPattern := `]`
	startIndex := strings.Index(text, startPattern)
	endIndex := strings.Index(text[startIndex:], endPattern)
	usernamesText := text[startIndex+len(startPattern) : startIndex+endIndex]

	// Extract strings from USERNAMES list
	pattern := `'([^']+)'`
	r := regexp.MustCompile(pattern)
	matches := r.FindAllStringSubmatch(usernamesText, -1)

	var usernames []string

	for _, match := range matches {
		usernames = append(usernames, match[1])
	}

	return usernames
}

// GetUsersOfType returns a slice of [User] and an error. It takes a UsernameType, which can be one of the following values:
//   - Bot - representing users with the `Bot` tag
//   - Core - representing Manifold employees
//   - Check - representing Manifold users with the `Trustworthy. ish.` label.
func (mc *Client) GetUsersOfType(t UsernameType) (*[]User, error) {
	text, _ := scrapeConstants(manifoldConstantsUrl)

	m := getUsernames(t, text)

	var us []User

	for _, b := range m {
		u, _ := mc.GetUserByUsername(b)
		us = append(us, *u)
	}

	return &us, nil
}

func (mc *Client) getLeaderboard() []leaderboardItem {
	leaders, _ := scrapeLeaderboards(manifoldLeaderboards)

	return leaders
}

// GetLeaders returns a slice of [LeadUser]. It takes a LeaderType, which can have one of the following values:
//   - Trader - an item on the top traders leaderboard
//   - Creator - an item on the top creators leaderboard
//
// And a LeaderPeriod, which can have one of the following values:
//   - Daily
//   - Weekly
//   - Monthly
//   - All
func (mc *Client) GetLeaders(t LeaderType, p LeaderPeriod) *[]LeadUser {
	leaders := mc.getLeaderboard()
	var leadUsers []LeadUser

	start := int(40*p) + int(20*t)
	stop := start + 20

	for i := start; i < stop; i++ {
		u, _ := mc.GetUserByUsername(leaders[i].Username)
		fmt.Println(u)
		leadUsers = append(leadUsers, LeadUser{
			Rank:     leaders[i].Rank,
			Username: leaders[i].Username,
			Mana:     leaders[i].Mana,
			Traders:  leaders[i].Traders,
			User:     *u,
		})
	}

	return &leadUsers
}

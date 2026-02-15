package mango

// LeagueEntry represents a single entry in league standings.
type LeagueEntry struct {
	UserId              string             `json:"userId"`
	ManaEarned          float64            `json:"manaEarned"`
	ManaEarnedBreakdown map[string]float64 `json:"manaEarnedBreakdown,omitempty"`
	Rank                int                `json:"rank,omitempty"`
	RankSnapshot        int                `json:"rankSnapshot,omitempty"`
	Season              int                `json:"season"`
	Cohort              string             `json:"cohort"`
	Division            int                `json:"division"`
	CreatedTime         int64              `json:"createdTime"`
}

// GetLeaguesRequest represents the parameters for fetching league standings.
type GetLeaguesRequest struct {
	UserId string `json:"userId,omitempty"`
	Season int    `json:"season,omitempty"`
	Cohort string `json:"cohort,omitempty"`
}

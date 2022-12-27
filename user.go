package mango

type ProfitCached struct {
	Weekly  float64 `json:"weekly"`
	Daily   float64 `json:"daily"`
	AllTime float64 `json:"allTime"`
	Monthly float64 `json:"monthly"`
}

type User struct {
	Id            string  `json:"id"`
	CreatedTime   int64   `json:"createdTime"`
	Name          string  `json:"name"`
	Username      string  `json:"username"`
	Url           string  `json:"url"`
	AvatarUrl     string  `json:"avatarUrl"`
	Balance       float64 `json:"balance"`
	TotalDeposits float64 `json:"totalDeposits"`
	ProfitCached  ProfitCached `json:"profitCached"`
	Bio           string `json:"bio,omitempty"`
	Website       string `json:"website,omitempty"`
	TwitterHandle string `json:"twitterHandle,omitempty"`
	DiscordHandle string `json:"discordHandle,omitempty"`
}
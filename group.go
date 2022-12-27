package mango

type group struct {
	AboutPostId            string `json:"aboutPostId,omitempty"`
	MostRecentActivityTime int64  `json:"mostRecentActivityTime"`
	AnyoneCanJoin          bool   `json:"anyoneCanJoin"`
	TotalContracts         int    `json:"totalContracts"`
	Name                   string `json:"name"`
	PinnedItems            []struct {
		ItemId string `json:"itemId"`
		Type   string `json:"type"`
	} `json:"pinnedItems,omitempty"`
	TotalMembers      int    `json:"totalMembers"`
	CreatedTime       int64  `json:"createdTime"`
	Slug              string `json:"slug"`
	CachedLeaderboard struct {
		TopTraders []struct {
			UserId string  `json:"userId"`
			Score  float64 `json:"score"`
		} `json:"topTraders"`
		TopCreators []struct {
			Score  int    `json:"score"`
			UserId string `json:"userId"`
		} `json:"topCreators"`
	} `json:"cachedLeaderboard"`
	About                       string   `json:"about"`
	MostRecentContractAddedTime int64    `json:"mostRecentContractAddedTime,omitempty"`
	CreatorId                   string   `json:"creatorId"`
	Id                          string   `json:"id"`
	PostIds                     []string `json:"postIds,omitempty"`
	BannerUrl                   string   `json:"bannerUrl,omitempty"`
	MostRecentChatActivityTime  int64    `json:"mostRecentChatActivityTime,omitempty"`
	ChatDisabled                bool     `json:"chatDisabled,omitempty"`
}
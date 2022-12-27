package mango

type PinnedItem struct {
	ItemId string `json:"itemId"`
	Type   string `json:"type"`
}

type Leader struct {
	UserId string  `json:"userId"`
	Score  float64 `json:"score"`
}

type CachedLeaderboard struct {
	TopTraders  []Leader `json:"topTraders"`
	TopCreators []Leader `json:"topCreators"`
}

type Group struct {
	AboutPostId                 string            `json:"aboutPostId,omitempty"`
	MostRecentActivityTime      int64             `json:"mostRecentActivityTime"`
	AnyoneCanJoin               bool              `json:"anyoneCanJoin"`
	TotalContracts              int64             `json:"totalContracts"`
	Name                        string            `json:"name"`
	PinnedItems                 []PinnedItem      `json:"pinnedItems,omitempty"`
	TotalMembers                int64             `json:"totalMembers"`
	CreatedTime                 int64             `json:"createdTime"`
	Slug                        string            `json:"slug"`
	CachedLeaderboard           CachedLeaderboard `json:"cachedLeaderboard"`
	About                       string            `json:"about"`
	MostRecentContractAddedTime int64             `json:"mostRecentContractAddedTime,omitempty"`
	CreatorId                   string            `json:"creatorId"`
	Id                          string            `json:"id"`
	PostIds                     []string          `json:"postIds,omitempty"`
	BannerUrl                   string            `json:"bannerUrl,omitempty"`
	MostRecentChatActivityTime  int64             `json:"mostRecentChatActivityTime,omitempty"`
	ChatDisabled                bool              `json:"chatDisabled,omitempty"`
}

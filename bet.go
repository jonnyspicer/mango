package mango

type Fees struct {
	LiquidityFee int64 `json:"liquidityFee"`
	PlatformFee  int64 `json:"platformFee"`
	CreatorFee   int64 `json:"creatorFee"`
}

type Fill struct {
	MatchedBetId *string `json:"matchedBetId"`
	Amount       float64 `json:"amount"`
	Shares       float64 `json:"shares"`
	Timestamp    int64   `json:"timestamp"`
	IsSale       bool    `json:"isSale,omitempty"`
}

type BetRequest struct {
	Amount     float64     `json:"amount"`
	ContractID string  `json:"contractId"`
	Outcome    string  `json:"outcome"`
	LimitProb  *float64 `json:"limitProb,omitempty"`
}

type Bet struct {
	Outcome       string  `json:"outcome"`
	Fees          Fees    `json:"fees"`
	IsAnte        bool    `json:"isAnte"`
	IsCancelled   bool    `json:"isCancelled,omitempty"`
	UserId        string  `json:"userId"`
	ProbBefore    float64 `json:"probBefore"`
	LoanAmount    float64 `json:"loanAmount"`
	ContractId    string  `json:"contractId"`
	UserUsername  string  `json:"userUsername"`
	CreatedTime   int64   `json:"createdTime"`
	UserAvatarUrl string  `json:"userAvatarUrl"`
	Id            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Fills         []Fill  `json:"fills,omitempty"`
	Shares        float64 `json:"shares"`
	IsRedemption  bool    `json:"isRedemption"`
	IsFilled      bool    `json:"isFilled,omitempty"`
	UserName      string  `json:"userName"`
	OrderAmount   float64 `json:"orderAmount,omitempty"`
	IsChallenge   bool    `json:"isChallenge"`
	ProbAfter     float64 `json:"probAfter"`
}

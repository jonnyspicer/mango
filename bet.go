package mango

// Fees represents the fees paid on a [Bet].
type Fees struct {
	LiquidityFee float64 `json:"liquidityFee"`
	PlatformFee  float64 `json:"platformFee"`
	CreatorFee   float64 `json:"creatorFee"`
}

// Fill represents the portion of a limit order that has been filled
type Fill struct {
	MatchedBetId *string `json:"matchedBetId"`
	Amount       float64 `json:"amount"`
	Shares       float64 `json:"shares"`
	Timestamp    int64   `json:"timestamp"`
	IsSale       bool    `json:"isSale,omitempty"`
}

// PostBetRequest represents the parameters required to post a new [Bet] via the API
type PostBetRequest struct {
	Amount     float64  `json:"amount"`
	ContractId string   `json:"contractId"`
	Outcome    string   `json:"outcome"`
	LimitProb  *float64 `json:"limitProb,omitempty"`
	AnswerId   string   `json:"answerId,omitempty"`
}

// GetBetsRequest represents the optional parameters that can be supplied to
// get bets via the API
type GetBetsRequest struct {
	UserId       string `json:"userId,omitempty"`
	Username     string `json:"username,omitempty"`
	ContractId   string `json:"contractID,omitempty"`
	ContractSlug string `json:"contractSlug,omitempty"`
	Before       string `json:"before,omitempty"`
	Limit        int64  `json:"limit,omitempty"`
	Kinds        string `json:"kinds,omitempty"`
}

// PostMultiBetRequest represents the parameters for placing multiple YES bets
// on a multiple choice market.
type PostMultiBetRequest struct {
	ContractId string   `json:"contractId"`
	AnswerIds  []string `json:"answerIds"`
	Amount     float64  `json:"amount"`
	LimitProb  *float64 `json:"limitProb,omitempty"`
	ExpiresAt  *int64   `json:"expiresAt,omitempty"`
}

// Bet represents a Bet object in the Manifold backend.
//
// See [the Manifold API docs for GET /v0/bets] for more details
//
// [the Manifold API docs for GET /v0/bets]: https://docs.manifold.markets/api#get-v0bets
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
	BetId         string  `json:"betId,omitempty"`
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

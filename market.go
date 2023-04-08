package mango

type outcomeType string

const (
	Binary         outcomeType = "BINARY"
	FreeResponse   outcomeType = "FREE_RESPONSE"
	MultipleChoice outcomeType = "MULTIPLE_CHOICE"
	Numeric        outcomeType = "NUMERIC"
	PseudoNumeric  outcomeType = "PSEUDO-NUMERIC"
)

type Pool struct {
	No       float64 `json:"NO,omitempty"`
	Yes      float64 `json:"YES,omitempty"`
	Option0  float64 `json:"0,omitempty"`
	Option1  float64 `json:"1,omitempty"`
	Option2  float64 `json:"2,omitempty"`
	Option3  float64 `json:"3,omitempty"`
	Option4  float64 `json:"4,omitempty"`
	Option5  float64 `json:"5,omitempty"`
	Option6  float64 `json:"6,omitempty"`
	Option7  float64 `json:"7,omitempty"`
	Option8  float64 `json:"8,omitempty"`
	Option9  float64 `json:"9,omitempty"`
	Option10 float64 `json:"10,omitempty"`
	Option11 float64 `json:"11,omitempty"`
	Option12 float64 `json:"12,omitempty"`
	Option13 float64 `json:"13,omitempty"`
	Option14 float64 `json:"14,omitempty"`
	Option15 float64 `json:"15,omitempty"`
	Option16 float64 `json:"16,omitempty"`
	Option17 float64 `json:"17,omitempty"`
	Option18 float64 `json:"18,omitempty"`
	Option19 float64 `json:"19,omitempty"`
}

type Answer struct {
	Id             string        `json:"id"`
	Username       string        `json:"username"`
	Name           string        `json:"name"`
	UserId string `json:"userId"`
	CreatedTime           int64         `json:"createdTime"`
	AvatarUrl      string        `json:"avatarUrl"`
	Number int64 `json:"number"`
	ContractId string `json:"contractId"`
	Text string `json:"text"`
	Probability float64 `json:"probability"`
}


type MarketRequest struct {
	OutcomeType string   `json:"outcomeType"`
	Question    string   `json:"question"`
	Description string   `json:"description,omitempty"`
	CloseTime   int64    `json:"closeTime,omitempty"`
	Visibility  string   `json:"visibility,omitempty"`
	GroupId     string   `json:"groupId,omitempty"`
	InitialProb int      `json:"initialProb,omitempty"`
	Min         int      `json:"min,omitempty"`
	Max         int      `json:"max,omitempty"`
	IsLogScale  bool     `json:"isLogScale,omitempty"`
	InitialVal  int      `json:"initialValue,omitempty"`
	Answers     []string `json:"answers,omitempty"`
}

type MarketResponse struct {
	Id string `json:"id"`
}

type LiquidityAmount struct {
	Amount int64 `json:"amount"`
}

type CloseTimestamp struct {
	CloseTime int64 `json:"closeTime,omitempty"`
}

type GroupId struct {
	GroupId string `json:"groupId,omitempty"`
}


type ResolveMarketRequest struct {
	Outcome        string               `json:"outcome"`
	Resolutions    []Resolution         `json:"resolutions,omitempty"`
	ProbabilityInt int                  `json:"probabilityInt,omitempty"`
	Value          float64              `json:"value,omitempty"`
}

type Resolution struct {
	Answer int `json:"answer"`
	Pct    int `json:"pct"`
}

type SellSharesRequest struct {
	Outcome string `json:"outcome,omitempty"`
	Shares  int    `json:"shares,omitempty"`
}



type Market interface{}

type LiteMarket struct {
	Id                    string        `json:"id"`
	CreatorId             string        `json:"creatorId"`
	CreatorUsername       string        `json:"creatorUsername"`
	CreatorName           string        `json:"creatorName"`
	CreatedTime           int64         `json:"createdTime"`
	CreatorAvatarUrl      string        `json:"creatorAvatarUrl"`
	CloseTime             int64         `json:"closeTime"`
	Question              string        `json:"question"`
	Tags                  []interface{} `json:"tags"`
	Url                   string        `json:"url"`
	Pool                  Pool          `json:"pool,omitempty"`
	Probability           float64       `json:"probability,omitempty"`
	P                     float64       `json:"p,omitempty"`
	TotalLiquidity        float64       `json:"totalLiquidity,omitempty"`
	OutcomeType           outcomeType   `json:"outcomeType"`
	Mechanism             string        `json:"mechanism"`
	Volume                float64       `json:"volume"`
	Volume24Hours         float64       `json:"volume24Hours"`
	IsResolved            bool          `json:"isResolved"`
	LastUpdatedTime       int64         `json:"lastUpdatedTime,omitempty"`
	Min                   float64       `json:"min,omitempty"`
	Max                   float64       `json:"max,omitempty"`
	IsLogScale            bool          `json:"isLogScale,omitempty"`
	Resolution            string        `json:"resolution,omitempty"`
	ResolutionTime        int64         `json:"resolutionTime,omitempty"`
	ResolutionProbability float64       `json:"resolutionProbability,omitempty"`
}

type FullMarket struct {
	Id                    string   `json:"id"`
	CreatorId             string   `json:"creatorId"`
	CreatorUsername       string   `json:"creatorUsername"`
	CreatorName           string   `json:"creatorName"`
	CreatedTime           int64    `json:"createdTime"`
	CreatorAvatarUrl      string   `json:"creatorAvatarUrl"`
	CloseTime             int64    `json:"closeTime"`
	Question              string   `json:"question"`
	Answers []Answer `json:"answers,omitempty"`
	Tags                  []string `json:"tags"`
	Url                   string   `json:"url"`
	Pool                  Pool     `json:"pool"`
	Probability           float64  `json:"probability"`
	P                     float64  `json:"p"`
	TotalLiquidity        float64  `json:"totalLiquidity"`
	OutcomeType           string   `json:"outcomeType"`
	Mechanism             string   `json:"mechanism"`
	Volume                float64  `json:"volume"`
	Volume24Hours         float64  `json:"volume24Hours"`
	IsResolved            bool     `json:"isResolved"`
	Resolution            string   `json:"resolution"`
	ResolutionTime        int64    `json:"resolutionTime"`
	ResolutionProbability float64  `json:"resolutionProbability"`
	LastUpdatedTime       int64    `json:"lastUpdatedTime"`
	// temporarily ignore `Description` field, as it's html-in-json, see https://tiptap.dev/guide/output#option-1-json
	// Description           string  `json:"description"`
	TextDescription string `json:"textDescription"`
}

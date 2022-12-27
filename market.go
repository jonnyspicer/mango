package mango

type outcomeType string

const (
	Binary outcomeType = "BINARY"
	FreeResponse outcomeType = "FREE_REPONSE"
	MultipleChoice outcomeType = "MULTIPLE_CHOICE"
	Numeric outcomeType = "NUMERIC"
	PseudoNumeric outcomeType = "PSEUDO-NUMERIC"
)

type Market interface {}

type FullMarket struct {}

type LiteMarket struct {
	Id               string        `json:"id"`
	CreatorId        string        `json:"creatorId"`
	CreatorUsername  string        `json:"creatorUsername"`
	CreatorName      string        `json:"creatorName"`
	CreatedTime      int64         `json:"createdTime"`
	CreatorAvatarUrl string        `json:"creatorAvatarUrl"`
	CloseTime        int64         `json:"closeTime"`
	Question         string        `json:"question"`
	Tags             []interface{} `json:"tags"`
	Url              string        `json:"url"`
	Pool             struct { // TODO: fix this
		NO      float64 `json:"NO,omitempty"`
		YES     float64 `json:"YES,omitempty"`
		Field3  float64 `json:"0,omitempty"`
		Field4  float64 `json:"1,omitempty"`
		Field5  float64 `json:"2,omitempty"`
		Field6  float64 `json:"3,omitempty"`
		Field7  float64 `json:"4,omitempty"`
		Field8  float64 `json:"5,omitempty"`
		Field9  float64 `json:"6,omitempty"`
		Field10 float64 `json:"7,omitempty"`
		Field11 float64 `json:"8,omitempty"`
		Field12 float64 `json:"9,omitempty"`
		Field13 float64 `json:"10,omitempty"`
		Field14 float64 `json:"11,omitempty"`
		Field15 int     `json:"12,omitempty"`
		Field16 int     `json:"13,omitempty"`
		Field17 int     `json:"14,omitempty"`
		Field18 int     `json:"15,omitempty"`
		Field19 int     `json:"16,omitempty"`
		Field20 int     `json:"17,omitempty"`
		Field21 int     `json:"18,omitempty"`
		Field22 int     `json:"19,omitempty"`
	} `json:"pool"`
	Probability           float64 `json:"probability,omitempty"`
	P                     float64 `json:"p,omitempty"`
	TotalLiquidity        int64     `json:"totalLiquidity,omitempty"`
	OutcomeType           outcomeType  `json:"outcomeType"`
	Mechanism             string  `json:"mechanism"`
	Volume                float64 `json:"volume"`
	Volume24Hours         float64 `json:"volume24Hours"`
	IsResolved            bool    `json:"isResolved"`
	LastUpdatedTime       int64   `json:"lastUpdatedTime,omitempty"`
	Min                   float64     `json:"min,omitempty"`
	Max                   float64 `json:"max,omitempty"`
	IsLogScale            bool    `json:"isLogScale,omitempty"`
	Resolution            string  `json:"resolution,omitempty"`
	ResolutionTime        int64   `json:"resolutionTime,omitempty"`
	ResolutionProbability float64 `json:"resolutionProbability,omitempty"`
}
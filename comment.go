package mango

type Comment struct {
	UserUsername             string  `json:"userUsername"`
	ContractSlug             string  `json:"contractSlug"`
	UserName                 string  `json:"userName"`
	CommentType              string  `json:"commentType"`
	Id                       string  `json:"id"`
	Text                     string  `json:"text"`
	CommenterPositionProb    float64 `json:"commenterPositionProb,omitempty"`
	ReplyToCommentId         string  `json:"replyToCommentId,omitempty"`
	ContractQuestion         string  `json:"contractQuestion"`
	CreatedTime              int64   `json:"createdTime"`
	UserAvatarUrl            string  `json:"userAvatarUrl"`
	ContractId               string  `json:"contractId"`
	UserId                   string  `json:"userId"`
	CommenterPositionShares  int     `json:"commenterPositionShares,omitempty"`
	CommenterPositionOutcome string  `json:"commenterPositionOutcome,omitempty"`
	BetAmount                float64 `json:"betAmount,omitempty"`
	BetId                    string  `json:"betId,omitempty"`
	BetOutcome               string  `json:"betOutcome,omitempty"`
}

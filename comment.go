package mango

// Comment represents a Comment object in the Manifold backend.
//
// This type isn't documented by Manifold and its structure was inferred from API calls.
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
	CommenterPositionShares  float64 `json:"commenterPositionShares,omitempty"`
	CommenterPositionOutcome string  `json:"commenterPositionOutcome,omitempty"`
	BetAmount                float64 `json:"betAmount,omitempty"`
	BetId                    string  `json:"betId,omitempty"`
	BetOutcome               string  `json:"betOutcome,omitempty"`
}

// PostCommentRequest represents the parameters required to post a
// new [Comment] via the API
type PostCommentRequest struct {
	ContractId string `json:"contractId"`
	Content    string `json:"content,omitempty"`
	Html       string `json:"html,omitempty"`
	Markdown   string `json:"markdown,omitempty"`
}

// GetCommentsRequest represents the optional parameters that can be supplied to
// get comments via the API
type GetCommentsRequest struct {
	ContractId   string `json:"contractId,omitempty"`
	ContractSlug string `json:"contractSlug,omitempty"`
}

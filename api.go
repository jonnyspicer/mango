// Package mango provides a client that can be used to make calls
// to the Manifold Markets API.
//
// Currently, it provides wrapper functions for every documented
// API call that Manifold offers. It also offers data types representing
// each data structure that can be returned by the API.
//
// See [the Manifold API docs] for more details.
//
// [the Manifold API docs]: https://docs.manifold.markets/api
package mango

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// TODO: POST doc comments need a note that they won't work without an API key

// GetAuthenticatedUser returns the [User] associated with the current API key.
// If the key is invalid or not present, returns nil and an error.
//
// See [the Manifold API docs for GET /v0/me] for more details.
//
// [the Manifold API docs for GET /v0/me]: https://docs.manifold.markets/api#get-v0me
func (mc *Client) GetAuthenticatedUser() (*User, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL(mc.url, GetMe, "", ""), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

// GetBets returns a slice of [Bet] and an error. It takes a [GetBetsRequest] which has the following
// optional parameters:
//   - [GetBetsRequest.UserId]
//   - [GetBetsRequest.Username]
//   - [GetBetsRequest.ContractId]
//   - [GetBetsRequest.ContractSlug]
//   - [GetBetsRequest.Before] - the ID of the bet before which the list will start.
//   - [GetBetsRequest.Limit] - the maximum and default limit is 1000.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/bets] for more details.
//
// [the Manifold API docs for GET /v0/bets]: https://docs.manifold.markets/api#get-v0bets
func (mc *Client) GetBets(gbr GetBetsRequest) (*[]Bet, error) {
	if gbr.Limit == 0 {
		gbr.Limit = DefaultLimit
	}
	resp, err := http.Get(requestURL(
		mc.url,
		GetBets, "", "",
		"userId", gbr.UserId,
		"username", gbr.Username,
		"contractId", gbr.ContractId,
		"contractSlug", gbr.ContractSlug,
		"before", gbr.Before,
		"limit", strconv.FormatInt(gbr.Limit, 10),
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []Bet{})
}

// GetComments returns a slice of [Comment] and an error. It takes a [GetCommentsRequest] which has the following
// optional parameters:
//   - [GetCommentsRequest.ContractId]
//   - [GetCommentsRequest.ContractSlug]
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/comments] for more details.
//
// [the Manifold API docs for GET /v0/comments]: https://docs.manifold.markets/api#get-v0comments
func (mc *Client) GetComments(gcr GetCommentsRequest) (*[]Comment, error) {
	if gcr.ContractId == "" && gcr.ContractSlug == "" {
		return nil, fmt.Errorf("either contractID or contractSlug must be specified")
	}

	resp, err := http.Get(requestURL(mc.url, GetComments, "", "",
		"contractId", gcr.ContractId,
		"contractSlug", gcr.ContractSlug,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []Comment{})
}

// GetGroupById returns a [Group] by its unique id.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/group/by-id/id] for more details.
//
// [the Manifold API docs for GET /v0/group/by-id/id]: https://docs.manifold.markets/api#get-v0groupby-idid
func (mc *Client) GetGroupById(id string) (*Group, error) {
	resp, err := http.Get(requestURL(mc.url, GetGroupByID, id, ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, Group{})
}

// GetGroupBySlug returns a [Group] by its slug.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/group/slug] for more details.
//
// [the Manifold API docs for GET /v0/group/slug]: https://docs.manifold.markets/api#get-v0groupslug
func (mc *Client) GetGroupBySlug(slug string) (*Group, error) {
	resp, err := http.Get(requestURL(mc.url, GetGroupBySlug, slug, ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, Group{})
}

// GetGroups returns a slice of [Group]. Optionally a userId can be passed:
// in this case, only groups available to the given user will be returned.
// Results are unordered.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/groups] for more details.
//
// [the Manifold API docs for GET /v0/groups]: https://docs.manifold.markets/api#get-v0groups
func (mc *Client) GetGroups(userId *string) (*[]Group, error) {
	uid := ""
	if userId != nil {
		uid = *userId
	}
	resp, err := http.Get(requestURL(
		mc.url, GetGroups, "", "",
		"availableToUserId", uid,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []Group{})
}

// GetMarketByID returns a [FullMarket] by its unique id.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/market/marketId] for more details.
//
// [the Manifold API docs for GET /v0/market/marketId]: https://docs.manifold.markets/api#get-v0marketmarketid
func (mc *Client) GetMarketByID(id string) (*FullMarket, error) {
	resp, err := http.Get(requestURL(mc.url, GetMarketByID, id, ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, FullMarket{})
}

// GetMarketBySlug returns a [FullMarket] by its unique slug.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/slug/marketSlug] for more details.
//
// [the Manifold API docs for GET /v0/slug/marketSlug]: https://docs.manifold.markets/api#get-v0slugmarketslug
func (mc *Client) GetMarketBySlug(slug string) (*FullMarket, error) {
	resp, err := http.Get(requestURL(mc.url, GetMarketBySlug, slug, ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, FullMarket{})
}

// GetMarkets returns a slice of [LiteMarket] and an error. It takes a [GetMarketsRequest] which has the following
// optional parameters:
//   - [GetMarketsRequest.Before] - the ID of the market before which the list will start.
//   - [GetMarketsRequest.Limit] - the maximum limit is 1000 and the default is 500.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/markets] for more details.
//
// [the Manifold API docs for GET /v0/markets]: https://docs.manifold.markets/api#get-v0markets
func (mc *Client) GetMarkets(gmr GetMarketsRequest) (*[]LiteMarket, error) {
	resp, err := http.Get(requestURL(
		mc.url, GetMarkets,
		"",
		"",
		"limit", strconv.FormatInt(gmr.Limit, 10), "before", gmr.Before,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []LiteMarket{})
}

// GetMarketsForGroup returns a slice of [LiteMarket] and an error. It takes a group ID to retrieve the markets for.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/group/by-id/id/markets] for more details
//
// [the Manifold API docs for GET /v0/group/by-id/id/markets]: https://docs.manifold.markets/api#get-v0groupby-ididmarkets
func (mc *Client) GetMarketsForGroup(id string) (*[]LiteMarket, error) {
	resp, err := http.Get(requestURL(mc.url, GetGroupByID, id, MarketsSuffix))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []LiteMarket{})
}

// GetMarketPositions returns positions for a given [Market] ID. It takes a [GetMarketPositionsRequest]
// which has the following parameters:
//   - [GetMarketPositionsRequest.MarketId] - Required.
//   - [GetMarketPositionsRequest.Order] - Optional. The field to order results by. Default: "profit". Options: "shares" or "profit". // TODO: make this an enum
//   - [GetMarketPositionsRequest.Top] - Optional. The number of top positions (ordered by order) to return. Default: null.
//   - [GetMarketPositionsRequest.Bottom] - Optional. The number of bottom positions (ordered by order) to return. Default: null.
//   - [GetMarketPositionsRequest.UserId] - Optional. The user ID to query by. Default: null. If provided, only the position for this user will be returned.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/market/marketId/positions] for more details.
//
// [the Manifold API docs for GET /v0/market/marketId/positions]: https://docs.manifold.markets/api#get-v0marketmarketidpositions
func (mc *Client) GetMarketPositions(gmpr GetMarketPositionsRequest) (*[]ContractMetric, error) {
	if gmpr.MarketId == "" {
		return nil, fmt.Errorf("no market ID provided")
	}

	var t, b string
	if gmpr.Top != 0 {
		t = fmt.Sprintf("%d", gmpr.Top)
	} else {
		t = "null"
	}

	if gmpr.Bottom != 0 {
		b = fmt.Sprintf("%d", gmpr.Bottom)
	} else {
		b = "null"
	}

	resp, err := http.Get(requestURL(mc.url, GetMarketByID, gmpr.MarketId, PositionsSuffix,
		"order", gmpr.Order,
		"top", t,
		"bottom", b,
		"userId", gmpr.UserId,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []ContractMetric{})
}

// SearchMarkets returns a slice of [FullMarket] that are the results of searching for provided terms.
// It can take any number of strings. Each string should be an individual search term.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/search-markets] for more details.
//
// [the Manifold API docs for GET /v0/search-markets]: https://docs.manifold.markets/api#get-v0search-markets
func (mc *Client) SearchMarkets(terms ...string) (*[]FullMarket, error) {
	ts := ""

	for i, t := range terms {
		ts += strings.TrimSpace(t)
		if i+1 < len(terms) {
			ts += " "
		}
	}

	resp, err := http.Get(requestURL(mc.url, GetSearchMarkets, "", "",
		"terms", ts,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []FullMarket{})
}

// GetUserByID returns a [User] by user id.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/user/by-id/id] for more details.
//
// [the Manifold API docs for GET /v0/user/by-id/id]: https://docs.manifold.markets/api#get-v0userby-idid
func (mc *Client) GetUserByID(id string) (*User, error) {
	resp, err := http.Get(requestURL(mc.url, GetUserByID, id, ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

// GetUserByUsername returns a [User] by username
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/user/username] for more details.
//
// [the Manifold API docs for GET /v0/user/username]: https://docs.manifold.markets/api#get-v0userusername
func (mc *Client) GetUserByUsername(un string) (*User, error) {
	resp, err := http.Get(requestURL(mc.url, GetUserByUsername, un, ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

// GetUsers returns a slice of [User] and an error. It takes a [GetUsersRequest] which has the following parameters:
//   - [GetUsersRequest.Before] - Optional. The ID of the user before which the list will start.
//   - [GetUsersRequest.Limit] - Optional. The default and maximum limit is 1000.
//
// If there is an error making the request, then nil and an error
// will be returned.
//
// See [the Manifold API docs for GET /v0/markets] for more details.
//
// [the Manifold API docs for GET /v0/markets]: https://docs.manifold.markets/api#get-v0markets
func (mc *Client) GetUsers(gur GetUsersRequest) (*[]User, error) {
	resp, err := http.Get(requestURL(
		mc.url, GetUsers,
		"",
		"",
		"limit", strconv.FormatInt(gur.Limit, 10), "before", gur.Before,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []User{})
}

// PostBet makes a new bet on a market. It takes a [PostBetRequest] which has the following parameters:
//   - [PostBetRequest.Amount] - Required.
//   - [PostBetRequest.ContractId] - Required.
//   - [PostBetRequest.Outcome] - Required.
//   - [PostBetRequest.LimitProb] - Optional. A number between 0.001 and 0.999 inclusive representing the limit probability for your bet
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/bet] for more details.
//
// [the Manifold API docs for POST /v0/bet]: https://docs.manifold.markets/api#post-v0bet
func (mc *Client) PostBet(br PostBetRequest) error {
	jsonBody, err := json.Marshal(br)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostBet,
		"",
		""), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bet placement failed with status %d", resp.StatusCode)
	}

	return nil
}

// CancelBet cancels an existing limit order for the given betId.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/bet/cancel/id] for more details.
//
// [the Manifold API docs for POST /v0/bet/cancel/id]: https://docs.manifold.markets/api#post-v0betcancelid
func (mc *Client) CancelBet(betId string) error {
	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostCancellation,
		betId,
		""), nil)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bet cancellation failed with status %d", resp.StatusCode)
	}

	return nil
}

// CreateMarket creates a new market. It takes a [PostMarketRequest] which has the following parameters:
//   - [PostMarketRequest.OutcomeType] - Required. One of [OutcomeType]
//   - [PostMarketRequest.Question] - Required.
//   - [PostMarketRequest.Description] - Optional. Non-rich text description.
//   - [PostMarketRequest.DescriptionHtml] - Optional.
//   - [PostMarketRequest.DescriptionMarkdown] - Optional.
//   - [PostMarketRequest.CloseTime] - Optional. In epoch time. Default 7 days from time of creation.
//   - [PostMarketRequest.Visibility] - Optional. One of "public" or "unlisted" TODO: make this an enum
//   - [PostMarketRequest.GroupId] - Optional. A group to show the market under.
//   - [PostMarketRequest.InitialProb] - Required for binary markets. Must be between 1 and 99.
//   - [PostMarketRequest.Min] - Required for numeric markets. The minimum value that the market may resolve to.
//   - [PostMarketRequest.Max] - Required for numeric markets. The maximum value that the market may resolve to.
//   - [PostMarketRequest.IsLogScale] - Required for numeric markets. If true, your numeric market will increase exponentially from min to max.
//   - [PostMarketRequest.InitialVal] - Required for numeric markets. An initial value for the market, between min and max, exclusive.
//   - [PostMarketRequest.Answers] - Required for multiple choice markets. An array of strings, each of which will be a valid answer for the market.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/market] for more details.
//
// [the Manifold API docs for POST /v0/market]: https://docs.manifold.markets/api#post-v0market
func (mc *Client) CreateMarket(mr PostMarketRequest) (*string, error) {
	// TODO: add input validation
	jsonBody, err := json.Marshal(mr)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostMarket,
		"",
		""), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("market creation failed with status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	marketResp := struct {
		id string `json:"id"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&marketResp); err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return &marketResp.id, nil
}

// AddLiquidity adds a given amount of liquidity to a given market.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/market/marketId/add-liquidity] for more details.
//
// [the Manifold API docs for POST /v0/market/marketId/add-liquidity]: https://docs.manifold.markets/api#post-v0marketmarketidadd-liquidity
func (mc *Client) AddLiquidity(marketId string, amount int64) error {
	amt := struct {
		Amount int64 `json:"amount"`
	}{amount}

	jsonBody, err := json.Marshal(amt)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostMarket,
		marketId,
		LiquiditySuffix), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("liquidity addition failed with status %d", resp.StatusCode)
	}

	return nil
}

// CloseMarket updates the closing time of a given market to be the given epoch timestamp,
// or closes it immediately if no new time is provided.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/market/marketId/close] for more details.
//
// [the Manifold API docs for POST /v0/market/marketId/close]: https://docs.manifold.markets/api#post-v0marketmarketidclose
func (mc *Client) CloseMarket(marketId string, ct *int64) error {
	c := struct {
		closeTime int64 `json:"closeTime,omitempty"`
	}{*ct}

	jsonBody, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostMarket,
		marketId,
		ClosureSuffix), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("market closure failed with status %d", resp.StatusCode)
	}

	return nil
}

// AddMarketToGroup adds a given market to a given group.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/market/marketId/group] for more details.
//
// [the Manifold API docs for POST /v0/market/marketId/group]: https://docs.manifold.markets/api#post-v0marketmarketidgroup
func (mc *Client) AddMarketToGroup(marketId, gi string) error {
	g := struct {
		groupId string `json:"groupId,omitempty"`
	}{gi}

	jsonBody, err := json.Marshal(g)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostMarket,
		marketId,
		GroupSuffix), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("adding market to group failed with status %d", resp.StatusCode)
	}

	return nil
}

// TODO: add more information about input params to this comment

// ResolveMarket creates a new market. It takes a [ResolveMarketRequest] which has the following parameters:
//   - [ResolveMarketRequest.Outcome] - Required. One of "YES", "NO", "MKT", "CANCEL" or a number depending on market type. TODO: make this an enum
//   - [ResolveMarketRequest.Resolutions] - Required for free response or multiple choice markets.
//   - [ResolveMarketRequest.ProbabilityInt] - Required if value is present.
//   - [ResolveMarketRequest.Value] - Optional, only relevant to numeric markets.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/market/marketId/resolve] for more details.
//
// [the Manifold API docs for POST /v0/market/marketId/resolve]: https://docs.manifold.markets/api#post-v0marketmarketidresolve
func (mc *Client) ResolveMarket(marketId string, rmr ResolveMarketRequest) error {
	jsonBody, err := json.Marshal(rmr)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostMarket,
		marketId,
		ResolutionSuffix), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("market resolution failed with status %d", resp.StatusCode)
	}

	return nil
}

// SellShares creates a new market. It takes a [SellSharesRequest] which has the following parameters:
//   - [SellSharesRequest.Outcome] - Optional. One of "YES" or "NO". If omitted and only one kind of shares are held, sells those. TODO: make this an enum
//   - [SellSharesRequest.Shares] - Optional. If omitted, all shares held will be sold.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/market/marketId/sell] for more details.
//
// [the Manifold API docs for POST /v0/market/marketId/sell]: https://docs.manifold.markets/api#post-v0marketmarketidsell
func (mc *Client) SellShares(marketId string, ssr SellSharesRequest) error {
	jsonBody, err := json.Marshal(ssr)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostMarket,
		marketId,
		SellSuffix), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("selling shares failed with status %d", resp.StatusCode)
	}

	return nil
}

// PostComment makes a new bet on a market. It takes a [PostCommentRequest] which has the following parameters:
//   - [PostCommentRequest.ContractId] - Required.
//   - [PostCommentRequest.Content] - Optional. A plaintext string.
//   - [PostCommentRequest.Html] - Optional.
//   - [PostCommentRequest.Markdown] - Optional.
//
// If there is an error making the request, then an error will be returned.
//
// See [the Manifold API docs for POST /v0/comment] for more details.
//
// [the Manifold API docs for POST /v0/comment]: https://docs.manifold.markets/api#post-v0comment
func (mc *Client) PostComment(marketId string, cr PostCommentRequest) error {
	jsonBody, err := json.Marshal(cr)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, PostComment,
		marketId,
		""), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.postRequest(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("posting comment failed with status %d", resp.StatusCode)
	}

	return nil
}

// parseResponse takes an HTTP response and a type and attempts to unmarshal
// the body from the response into the given type.
func parseResponse[S any](r *http.Response, s S) (*S, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code found: %v, message: %v", r.StatusCode, string(body))
	}

	if err = json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &s, nil
}

func (mc *Client) postRequest(req *http.Request) (*http.Response, error) {
	if mc.key == "" {
		return nil, fmt.Errorf("no API key found")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", mc.key))

	return mc.postRequest(req)
}

package mango

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jonnyspicer/mango/endpoint"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// TODO: move all this to another file
type ManiClient struct {
	client http.Client
	key    string
	url    string
}

var lock = &sync.Mutex{}
var mcInstance *ManiClient

func maniClientInstance(client *http.Client, url *string) *ManiClient {
	if mcInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if mcInstance == nil {
			if client == nil {
				client = &http.Client{
					Timeout: 10,
				}
			}

			if url == nil {
				u := endpoint.Base
				url = &u
			}

			mcInstance = &ManiClient{
				client: *client,
				key:    apiKey(),
				url:    *url,
			}
		}
	}
	return mcInstance
}

func defaultManiClient() *ManiClient {
	return maniClientInstance(nil, nil)
}

func (mc *ManiClient) destroy() {
	if mcInstance != nil {
		lock.Lock()
		defer lock.Unlock()
		mcInstance = nil
	}
}

func apiKey() string {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Errorf("fatal error config file: %w", err)
	}

	return viper.GetString("MANIFOLD_API_KEY")
}

func (mc *ManiClient) GetAuthenticatedUser() User {
	req, err := http.NewRequest(http.MethodGet, endpoint.RequestURL(mc.url, endpoint.GetMe, "", ""), nil)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

func (mc *ManiClient) GetBets(userID, username, contractID, contractSlug, before string, limit int) []Bet {
	if limit == 0 {
		limit = endpoint.DefaultLimit
	}
	resp, err := http.Get(endpoint.RequestURL(
		mc.url,
		endpoint.GetBets, "", "",
		"userId", userID,
		"username", username,
		"contractId", contractID,
		"contractSlug", contractSlug,
		"before", before,
		"limit", strconv.Itoa(limit),
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []Bet{})
}

func (mc *ManiClient) GetComments(contractID, contractSlug string) []Comment {
	if contractID == "" && contractSlug == "" {
		log.Println("either contractID or contractSlug must be specified")
		return nil
	}

	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetComments, "", "",
		"contractId", contractID,
		"contractSlug", contractSlug,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []Comment{})
}

func (mc *ManiClient) GetGroupByID(id string) Group {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetGroupByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, Group{})
}

func (mc *ManiClient) GetGroupBySlug(slug string) Group {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetGroupBySlug, slug, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, Group{})
}

func (mc *ManiClient) GetGroups(userID string) []Group {
	resp, err := http.Get(endpoint.RequestURL(
		mc.url, endpoint.GetGroups, "", "",
		"availableToUserId", userID,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []Group{})
}

func (mc *ManiClient) GetMarketByID(id string) FullMarket {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetMarketByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, FullMarket{})
}

func (mc *ManiClient) GetMarketBySlug(slug string) FullMarket {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetMarketBySlug, slug, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, FullMarket{})
}

func (mc *ManiClient) GetMarkets(before string, limit int) []LiteMarket {
	resp, err := http.Get(endpoint.RequestURL(
		mc.url, endpoint.GetMarkets,
		"",
		"",
		"limit", strconv.Itoa(limit), "before", before,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []LiteMarket{})
}

func (mc *ManiClient) GetMarketsForGroup(id string) []LiteMarket {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetGroupByID, id, endpoint.MarketsSuffix))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []LiteMarket{})
}

func (mc *ManiClient) GetMarketPositions(marketId string, order string, top *int, bottom *int, userId string) []ContractMetric {
	var t, b string
	if top != nil {
		t = fmt.Sprintf("%d", *top)
	} else {
		t = "null"
	}

	if bottom != nil {
		b = fmt.Sprintf("%d", *bottom)
	} else {
		b = "null"
	}

	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetMarketByID, marketId, endpoint.PositionsSuffix,
		"order", order,
		"top", t,
		"bottom", b,
		"userId", userId,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []ContractMetric{})
}

func (mc *ManiClient) SearchMarkets(terms string) []FullMarket {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetSearchMarkets, "", "",
		"terms", terms,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []FullMarket{})
}

func (mc *ManiClient) GetUserByID(id string) User {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetUserByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

func (mc *ManiClient) GetUserByUsername(un string) User {
	resp, err := http.Get(endpoint.RequestURL(mc.url, endpoint.GetUserByUsername, un, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

func (mc *ManiClient) GetUsers(before string, limit int) []User {
	resp, err := http.Get(endpoint.RequestURL(
		mc.url, endpoint.GetUsers,
		"",
		"",
		"limit", strconv.Itoa(limit), "before", before,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []User{})
}

func (mc *ManiClient) PostBet(br BetRequest) {
	jsonBody, err := json.Marshal(br)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostBet,
		"",
		""), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("bet placement failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) CancelBet(betId string) {
	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostCancellation,
		betId,
		""), nil)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("bet cancellation failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) CreateMarket(mr MarketRequest) string {
	jsonBody, err := json.Marshal(mr)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostMarket,
		"",
		""), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("market creation failed with status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var marketResp MarketResponse
	if err = json.NewDecoder(resp.Body).Decode(&marketResp); err != nil {
		log.Printf("error reading response body: %v", err)
		return ""
	}

	return marketResp.Id
}

func (mc *ManiClient) AddLiquidity(marketId string, amount LiquidityAmount) {
	jsonBody, err := json.Marshal(amount)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostMarket,
		marketId,
		endpoint.LiquiditySuffix), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("liquidity addition failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) CloseMarket(marketId string, ct CloseTimestamp) {
	jsonBody, err := json.Marshal(ct)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostMarket,
		marketId,
		endpoint.ClosureSuffix), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("market closure failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) AddMarketToGroup(marketId string, gi GroupId) {
	jsonBody, err := json.Marshal(gi)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostMarket,
		marketId,
		endpoint.GroupSuffix), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("adding market to group failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) ResolveMarket(marketId string, rmr ResolveMarketRequest) {
	jsonBody, err := json.Marshal(rmr)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostMarket,
		marketId,
		endpoint.ResolutionSuffix), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("market resolution failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) SellShares(marketId string, ssr SellSharesRequest) {
	jsonBody, err := json.Marshal(ssr)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostMarket,
		marketId,
		endpoint.SellSuffix), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("selling shares failed with status %d", resp.StatusCode)
	}
}

func (mc *ManiClient) PostComment(marketId string, cr CommentRequest) {
	jsonBody, err := json.Marshal(cr)
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, endpoint.RequestURL(
		mc.url, endpoint.PostComment,
		marketId,
		""), bodyReader)
	if err != nil {
		log.Printf("error creating http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", apiKey()))

	resp, err := mc.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("posting comment failed with status %d", resp.StatusCode)
	}
}

func parseResponse[S any](r *http.Response, s S) S {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	if r.StatusCode != http.StatusOK {
		fmt.Errorf("non-200 status code found: %v, message: %v", r.StatusCode, string(body))
		return s
	}

	if err = json.Unmarshal(body, &s); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return s
}

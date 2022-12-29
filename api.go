package mango

import (
	"encoding/json"
	"github.com/jonnyspicer/mango/endpoint"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetBets(userID, username, contractID, contractSlug, before string, limit int) []Bet {
	if limit == 0 { limit = endpoint.DefaultLimit}
	resp, err := http.Get(endpoint.RequestURL(
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

func GetComments(contractID, contractSlug string) []Comment {
	if contractID == "" && contractSlug == "" {
		log.Println("either contractID or contractSlug must be specified")
		return nil
	}

	resp, err := http.Get(endpoint.RequestURL(endpoint.GetComments, "", "",
		"contractId", contractID,
		"contractSlug", contractSlug,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []Comment{})
}

func GetGroupByID(id string) Group {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetGroupByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, Group{})
}

func GetGroupBySlug(slug string) Group {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetGroupBySlug, slug, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, Group{})
}

func GetGroups(userID string) []Group {
	resp, err := http.Get(endpoint.RequestURL(
		endpoint.GetGroups, "", "",
		"availableToUserId", userID,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []Group{})
}

func GetMarketByID(id string) FullMarket {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetMarketByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, FullMarket{})
}

func GetMarketBySlug(slug string) FullMarket {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetMarketBySlug, slug, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, FullMarket{})
}

func GetMarkets(before string, limit int) []LiteMarket {
	resp, err := http.Get(endpoint.RequestURL(
		endpoint.GetMarkets,
		"",
		"",
		"limit", strconv.Itoa(limit), "before", before,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []LiteMarket{})
}

func GetMarketsForGroup(id string) []LiteMarket {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetGroupByID, id, endpoint.MarketsSuffix))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []LiteMarket{})
}

func GetUserByID(id string) User {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetUserByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

func GetUserByUsername(un string) User {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetUserByUsername, un, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}

func GetUsers(before string, limit int) []User {
	resp, err := http.Get(endpoint.RequestURL(
		endpoint.GetUsers,
		"",
		"",
		"limit", strconv.Itoa(limit), "before", before,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	return parseResponse(resp, []User{})
}

func parseResponse[S any](r *http.Response, s S) S {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading r body: %v", err)
	}

	if err = json.Unmarshal(body, &s); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return s
}
package mango

import (
	"encoding/json"
	"github.com/jonnyspicer/mango/endpoint"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetUsers(limit int, before string) []User {
	resp, err := http.Get(endpoint.RequestURL(
		endpoint.GetUsers,
		"",
		"",
		"limit", strconv.Itoa(limit), "before", before,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var us []User

	if err = json.Unmarshal(body, &us); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return us
}

func GetUserByUsername(un string) User {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetUserByUsername, un, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var u User

	if err = json.Unmarshal(body, &u); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return u
}

func GetUserByID(id string) User {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetUserByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var u User

	if err = json.Unmarshal(body, &u); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return u
}

func GetBets(userID, username, contractID, contractSlug, before string, limit int) []Bet {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var bs []Bet

	if err = json.Unmarshal(body, &bs); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return bs
}

func GetGroups(userID string) []Group {
	resp, err := http.Get(endpoint.RequestURL(
		endpoint.GetGroups, "", "",
		"availableToUserId", userID,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var gs []Group

	if err = json.Unmarshal(body, &gs); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return gs
}

func GetGroupBySlug(slug string) Group {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetGroupBySlug, slug, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var g Group

	if err = json.Unmarshal(body, &g); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return g
}

func GetGroupByID(id string) Group {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetGroupByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var g Group

	if err = json.Unmarshal(body, &g); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return g
}

func GetMarketsForGroup(id string) []Market {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetGroupByID, id, endpoint.MarketsSuffix))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var ms []Market

	if err = json.Unmarshal(body, &ms); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return ms
}

func GetMarkets(limit int, before string) []LiteMarket {
	resp, err := http.Get(endpoint.RequestURL(
		endpoint.GetMarkets,
		"",
		"",
		"limit", strconv.Itoa(limit), "before", before,
	))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var ms []LiteMarket

	if err = json.Unmarshal(body, &ms); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return ms
}

func GetMarketBySlug(slug string) FullMarket {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetMarketBySlug, slug, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var fm FullMarket

	if err = json.Unmarshal(body, &fm); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return fm
}

func GetMarketByID(id string) FullMarket {
	resp, err := http.Get(endpoint.RequestURL(endpoint.GetMarketByID, id, ""))
	if err != nil {
		log.Printf("error making http request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var fm FullMarket

	if err = json.Unmarshal(body, &fm); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return fm
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	var cs []Comment

	if err = json.Unmarshal(body, &cs); err != nil {
		log.Printf("error unmarshalling JSON: %v", err)
	}

	return cs
}
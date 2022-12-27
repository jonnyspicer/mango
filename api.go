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
		"limit", strconv.Itoa(limit),  "before", before,
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


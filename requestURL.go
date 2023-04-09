package mango

import (
	"log"
	"net/url"
)

// requestURL returns a fully-formed URL that HTTP requests can be sent to.
// It includes the base domain, path, and any query parameters supplied.
func requestURL(base, path, value, suffix string, params ...string) string {
	// if query parameters are supplied, they must be in key:value pairs
	if len(params)%2 != 0 {
		log.Println("number of params passed to requestURL() must be divisible by 2")
		return ""
	}

	query, err := url.Parse(base + "/" + Version)
	if err != nil {
		log.Fatalf("error parsing base URL: %v", err)
	}

	query.Path += path

	if value != "" {
		query.Path += value
	} else {
		ps := url.Values{}

		for i := 0; i < len(params); i += 2 {
			if params[i+1] != "" {
				ps.Add(params[i], params[i+1])
			}
		}

		query.RawQuery = ps.Encode()
	}

	// the suffix represents part of the path that is supplied after an identifier
	// eg for the path: /v0/market/[marketId]/add-liquidity
	// the string "add-liquidity" would be the suffix
	if suffix != "" {
		query.Path += suffix
	}

	return query.String()
}

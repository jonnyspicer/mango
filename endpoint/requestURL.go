package endpoint

import (
	"log"
	"net/url"
)

func RequestURL(path, value, suffix string, params ...string) string {
	if len(params) % 2 != 0 {
		log.Println("number of params passed to RequestURL() must be divisible by 2")
		return ""
	}

	query, err := url.Parse(Base + Version)
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

	if suffix != "" {
		query.Path += suffix
	}

	return query.String()
}
package mango

import (
	"fmt"
	"io"
	"net/http"
)

func GetUsers() []User {
	resp, err := http.Get("https://manifold.markets/api/v0/users")
	if err != nil {
		fmt.Println("oh no")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("OH NO")
	}

	fmt.Println(string(body))

	return nil
}


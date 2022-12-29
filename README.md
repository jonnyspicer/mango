# Mango

Mango is a Go library for interacting with [Manifold.](https://manifold.markets)

Currently it provides wrappers for all GET endpoints in the [Manifold API](https://docs.manifold.markets/api) that don't require authentication.

## Installation

`go get github.com/jonnyspicer/mango`

## Usage

Mango offers custom structs representing different data structures used by Manifold, as well as methods to call the Manifold API and retrieve those objects.
All parameters that are optional must be passed in with values - in practice this means that string parameters you don't want to be set must be passed as `""` and int parameters as `0`.

The available methods are:

```go
package mango

func GetBets(userID, username, contractID, contractSlug, before string, limit int) []Bet {}

func GetComments(contractID, contractSlug string) []Comment {}

func GetGroupBySlug(slug string) Group {}

func GetGroups(userID string) []Group {}

func GetMarketByID(id string) FullMarket {}

func GetMarketBySlug(slug string) FullMarket {}

func GetMarkets(before string, limit int) []LiteMarket {}

func GetMarketsForGroup(id string) []LiteMarket {}

func GetUserByID(id string) User {}

func GetUserByUsername(un string) User {}

func GetUsers(before string, limit int) []User {}
```
# Mango 🥭

Mango is a Go library for interacting with [Manifold Markets.](https://manifold.markets) It provides wrapper functions 
for every documented API call that Manifold offers. It also offers data types representing each data structure 
that can be returned by the API.

See the [Manifold API docs](https://docs.manifold.markets/api) for more details.

## Installation

`go get github.com/jonnyspicer/mango`

## Usage

Mango offers custom structs representing different data structures used by Manifold, as well as methods to call the Manifold API and retrieve those objects.


In order for some functions to work correctly, you will need to have a `MANIFOLD_API_KEY` set in the
`.env` file in the root of your project. Your key can be found on the edit profile screen in the Manifold UI.

### Basics

Full documentation, including all available functions and types, is available on [pkg.go.dev](https://pkg.go.dev/github.com/jonnyspicer/mango#section-documentation).

Get information about the currently authenticated user:

```go
package main

import (
	"github.com/jonnyspicer/mango"
	"fmt"
)

// Initialize the default mango client
// This will try to read the value of `MANIFOLD_API_KEY` from your .env file
// and will use the base URL `https://manifold.markets` for all requests.
mc := mango.DefaultClientInstance()

user, err := mc.GetAuthenticatedUser()
if err != nil {
  log.Errorf("error getting authenticated user: %v", err)
}

fmt.Printf("authenticated user: %+v", *user)
```

```shell
$ authenticated user: {Id:xN67Q0mAhddL0X9wVYP2YfOrYH42 CreatedTime:1653515196337 Name:Jonny Spicer Username:jonny Url: AvatarUrl:https://lh3.googleuser
content.com/a-/AOh14GikSB2nbgbE_S2n-QUj9ydaNOX1w3QHIQrkvSsQHA=s96-c BannerUrl: Balance:10701.116370604414 TotalDeposits:12267.829283182942 ProfitCach
ed:{Weekly:107.5333580138431 Daily:19.18702587612779 AllTime:3409.646711972091 Monthly:307.27444250182816} Bio: Website:https://jonnyspicer.com Twitt
erHandle:https://twitter.com/jjspicer DiscordHandle:}
```

Create a new market:

```go
mc := mango.DefaultClientInstance()

pmr := mango.PostMarketRequest{
    OutcomeType: mango.Binary,
    Question:    "How much wood would a woodchuck chuck if a woodchuck could chuck wood?",
    Description: "Will resolve based on some completely arbitrary criteria",
    InitialProb: 50,
    CloseTime:   1704067199000, // Sunday, December 31, 2023 11:59:59 PM
}

marketId, err := mc.CreateMarket(pmr)
if err != nil {
    fmt.Printf("error creating market: %v", err)
}

fmt.Printf("created market id: %v", *marketId)
```

```shell
$ created market id: 1LZpVeeTGAjkF4IgPAMk
```

Bet on a market:

```go
mc := mango.DefaultClientInstance()

pbr := mango.PostBetRequest{
    Amount:     10,
    ContractId: "1LZpVeeTGAjkF4IgPAMk",
    Outcome:    "YES",
}

bet, err := mc.PostBet(pbr)
if err != nil {
    fmt.Printf("error posting bet: %v", err)
}
fmt.Printf("placed bet %s, shares: %f", bet.Id, bet.Shares)
```
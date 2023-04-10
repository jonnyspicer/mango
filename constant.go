package mango

// Various constants associated with the Manifold REST API.
// See [the Manifold API docs] for more details
//
// [the Manifold API docs]: https://docs.manifold.markets/api

const base string = "https://manifold.markets/api"
const version string = "v0/"

const defaultLimit = 1000

const getBets string = "bets/"
const getComments string = "comments/"
const getGroupBySlug string = "group/"
const getGroupByID string = "group/by-id/"
const getGroups string = "groups/"
const getMarketBySlug string = "slug/"
const getMarketByID string = "market/"
const getMarkets string = "markets/"
const getSearchMarkets string = "search-markets/"
const getMe string = "me/"
const getUserByUsername string = "user/"
const getUserByID string = "user/by-id/"
const getUsers string = "users/"

const postBet string = "bet/"
const postCancellation string = "bet/cancel/"
const postMarket string = "market/"
const postComment string = "comment/"

const marketsSuffix = "/markets/"
const positionsSuffix = "/positions/"
const liquiditySuffix = "/liquidity/"
const closureSuffix = "/close/"
const groupSuffix = "/group/"
const resolutionSuffix = "/resolve/"
const sellSuffix = "/sell/"

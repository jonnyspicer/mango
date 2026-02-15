# Mango API Library Update Design

**Date:** 2026-02-15
**Goal:** Full modernization of the Mango Go library to match the current Manifold Markets API

## Approach: Parallel Domain Teams

Five agents working in parallel on isolated domains, then integrating and testing together.

## Agent Responsibilities

### Agent 1: Core (shared infrastructure)
- Update `go.mod` from Go 1.19 to 1.21+
- Rename `postRequest` to `doRequest` and add `getRequest` helper that sends auth headers on GET requests
- Add response body to POST error messages
- Fix `Pool` struct: replace hardcoded `Option0`..`Option19` with `map[string]float64`
- Fix exported field bugs in `CloseMarket` (`closeTime` -> `CloseTime`) and `AddMarketToGroup` (`groupId` -> `GroupId`)
- Update dependencies

### Agent 2: Users & Portfolio
**Files:** `user.go`, new `portfolio.go`
- Add `DisplayUser` struct (lite user info)
- Add `GetUserLite(username)` and `GetUserByIDLite(id)`
- Add `LivePortfolioMetrics` struct
- Add `GetUserPortfolio(userId)` and `GetUserPortfolioHistory(userId, period)`
- Add `GetUserContractMetricsWithContracts()`
- Audit `User` struct against current API responses

### Agent 3: Markets & Search
**Files:** `market.go`, `api.go`
- Add `GetMarketProb(marketId)` and `GetMarketProbs(ids)`
- Add `PostAnswer(marketId, text)`
- Create `SearchMarketsRequest` struct with full filter/sort params, update `SearchMarkets`
- Add `AddBounty(marketId, amount)` and `AwardBounty(marketId, amount, commentId)`
- Audit `LiteMarket`/`FullMarket` structs against current API responses

### Agent 4: Bets & Transactions
**Files:** `bet.go`, new `transaction.go`
- Add `PostMultiBet()` with `PostMultiBetRequest`
- Update `GetBetsRequest` with `Kinds` field
- Add `Txn` struct
- Add `GetTransactions()` with `GetTransactionsRequest`
- Add `SendManagram()` with `SendManagramRequest`
- Add `GetLeagues()` with `GetLeaguesRequest`
- Audit `Bet` struct against current API responses

### Agent 5: Test Coordinator
**Files:** `api_test.go`, new `integration_test.go`
- Write unit tests (httptest mocks) for all new endpoints
- Write integration tests behind `//go:build integration` tag
- Integration tests use minimal mana (1 mana bets, single shared test market)
- Verify all existing tests pass
- Move existing real-API tests behind build tag

## Data Structure Changes

### New Types
- `DisplayUser` - lightweight user for /lite endpoints
- `LivePortfolioMetrics` - portfolio data
- `PortfolioMetrics` - portfolio history entry
- `Txn` - transaction record
- `PostMultiBetRequest`, `SendManagramRequest`, `GetTransactionsRequest`, `GetLeaguesRequest`, `SearchMarketsRequest`, etc.

### Updated Types
- `Pool` - `map[string]float64` instead of hardcoded options
- `GetBetsRequest` - add `Kinds` field
- `LiteMarket`/`FullMarket` - audit for new fields
- `User` - audit for new fields

### Bug Fixes
- `CloseMarket`: unexported `closeTime` field never marshals to JSON
- `AddMarketToGroup`: unexported `groupId` field never marshals to JSON

## Testing Strategy

### Unit Tests
- Every function gets an httptest mock test
- POST tests verify request body, HTTP method, and auth headers

### Integration Tests
- Build tag: `//go:build integration`
- Read-only tests: fetch known users, markets, bets, comments, etc.
- Write tests: create one shared test market, bet/comment/resolve on it
- Mana conservation: 1 mana bets, minimal market creation
- Sequential execution for write tests

## Out of Scope
- WebSocket support
- `POST /unresolve` (internal API)
- Major client architecture changes (singleton pattern stays)
- Removing `colly` scraping functions (evaluate separately)

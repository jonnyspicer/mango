# Mango API Library Update Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Update the Mango Go library to cover all current Manifold Markets API endpoints, fix existing bugs, and add comprehensive unit + integration tests.

**Architecture:** Parallel domain teams working on isolated files. Core infrastructure changes land first, then domain agents work in parallel on users, markets, bets/transactions. Integration tests land last.

**Tech Stack:** Go 1.21+, net/http, httptest, Manifold Markets API v0

**API Docs:** https://docs.manifold.markets/api

**Phases:**
- Phase 1: Core infrastructure (blocking - must complete first)
- Phase 2: Domain endpoints (parallel - 3 agents)
- Phase 3: Integration tests (after Phase 2)

---

## Phase 1: Core Infrastructure

### Task 1: Fix Pool struct

The `Pool` type hardcodes `Option0`..`Option19` which breaks for markets with >20 answers and is brittle. Replace with a dynamic map.

**Files:**
- Modify: `market.go` (lines 21-44)

**Step 1: Replace Pool struct with map type alias**

In `market.go`, replace the `Pool` struct with:

```go
// Pool represents the potential outcomes for a market.
// Keys are outcome names like "YES", "NO", or answer indices like "0", "1", etc.
type Pool map[string]float64
```

Delete the entire existing `Pool struct { ... }` definition (lines 21-44).

**Step 2: Fix LiteMarket and FullMarket Pool field tags**

The `Pool` field in both `LiteMarket` and `FullMarket` already uses `json:"pool"` tags which work fine with a map type. No changes needed here.

**Step 3: Run existing tests to check for breakage**

Run: `go build ./...`
Expected: PASS (Pool is only used as a field, not accessed by named fields in tests)

Run: `go test ./... -run TestSearchMarkets`
Expected: PASS (the test creates a Pool{} which is now an empty map)

**Step 4: Commit**

```bash
git add market.go
git commit -m "refactor: replace hardcoded Pool struct with map[string]float64"
```

---

### Task 2: Fix unexported field bugs in CloseMarket and AddMarketToGroup

Both functions use anonymous structs with unexported (lowercase) fields, so `json.Marshal` produces `{}`.

**Files:**
- Modify: `api.go` (lines 560-562 for CloseMarket, lines 599-601 for AddMarketToGroup)

**Step 1: Write failing test for CloseMarket serialization**

In `api_test.go`, add:

```go
func TestCloseMarketSendsBody(t *testing.T) {
	var receivedBody []byte

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	ct := int64(1704067199000)
	err := mc.CloseMarket("123", &ct)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(receivedBody) == "{}" || string(receivedBody) == "null" {
		t.Errorf("expected non-empty body with closeTime, got: %s", receivedBody)
	}
}
```

Add `"io"` to the import block in `api_test.go` if not already present.

**Step 2: Run test to verify it fails**

Run: `go test ./... -run TestCloseMarketSendsBody -v`
Expected: FAIL - body is `{}`

**Step 3: Fix CloseMarket - export the field**

In `api.go`, change the anonymous struct in `CloseMarket` (around line 560-562) from:

```go
c := struct {
    closeTime int64 `json:"closeTime,omitempty"`
}{*ct}
```

to:

```go
c := struct {
    CloseTime int64 `json:"closeTime,omitempty"`
}{*ct}
```

**Step 4: Run test to verify it passes**

Run: `go test ./... -run TestCloseMarketSendsBody -v`
Expected: PASS

**Step 5: Write failing test for AddMarketToGroup serialization**

In `api_test.go`, add:

```go
func TestAddMarketToGroupSendsBody(t *testing.T) {
	var receivedBody []byte

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	err := mc.AddMarketToGroup("123marketid", "456groupid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(receivedBody) == "{}" || string(receivedBody) == "null" {
		t.Errorf("expected non-empty body with groupId, got: %s", receivedBody)
	}
}
```

**Step 6: Run test to verify it fails**

Run: `go test ./... -run TestAddMarketToGroupSendsBody -v`
Expected: FAIL

**Step 7: Fix AddMarketToGroup - export the field**

In `api.go`, change the anonymous struct in `AddMarketToGroup` (around line 599-601) from:

```go
g := struct {
    groupId string `json:"groupId,omitempty"`
}{gi}
```

to:

```go
g := struct {
    GroupId string `json:"groupId,omitempty"`
}{gi}
```

**Step 8: Run tests to verify both pass**

Run: `go test ./... -run "TestCloseMarketSendsBody|TestAddMarketToGroupSendsBody" -v`
Expected: PASS

**Step 9: Commit**

```bash
git add api.go api_test.go
git commit -m "fix: export anonymous struct fields so json.Marshal includes them"
```

---

### Task 3: Add authenticated GET request helper and rename postRequest

Currently only POST requests send the `Authorization` header. GET endpoints that need auth (like `/v0/me`, `/v0/get-user-portfolio`) won't work properly. Also, bare `http.Get()` calls bypass the client's configured timeout.

**Files:**
- Modify: `api.go`

**Step 1: Rename postRequest to doRequest**

In `api.go`, rename the method (around line 766):

```go
func (mc *Client) doRequest(req *http.Request) (*http.Response, error) {
	if mc.key == "" {
		return nil, fmt.Errorf("no API key found")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %v", mc.key))

	return mc.client.Do(req)
}
```

Update all call sites in `api.go` that reference `mc.postRequest(` to `mc.doRequest(`.
These are in: `GetAuthenticatedUser`, `PostBet`, `CancelBet`, `CreateMarket`, `AddLiquidity`, `CloseMarket`, `AddMarketToGroup`, `ResolveMarket`, `SellShares`, `PostComment`.

**Step 2: Add getRequest helper**

Add this method to `api.go` (after `doRequest`):

```go
// getRequest makes an authenticated GET request to the given URL.
// Unlike http.Get(), this sends the Authorization header and uses the client's timeout.
func (mc *Client) getRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Key %v", mc.key))

	return mc.client.Do(req)
}
```

**Step 3: Update GET endpoints to use getRequest**

Replace all `http.Get(requestURL(...))` calls with `mc.getRequest(requestURL(...))`.

Functions to update: `GetBets`, `GetComments`, `GetGroupById`, `GetGroupBySlug`, `GetGroups`, `GetMarketByID`, `GetMarketBySlug`, `GetMarkets`, `GetMarketsForGroup`, `GetMarketPositions`, `SearchMarkets`, `GetUserByID`, `GetUserByUsername`, `GetUsers`.

For each, change from:
```go
resp, err := http.Get(requestURL(...))
```
to:
```go
resp, err := mc.getRequest(requestURL(...))
```

Also simplify `GetAuthenticatedUser` - it currently manually builds a request. Replace with:
```go
func (mc *Client) GetAuthenticatedUser() (*User, error) {
	resp, err := mc.getRequest(requestURL(mc.url, getMe, "", ""))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, User{})
}
```

**Step 4: Run all existing tests**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 5: Commit**

```bash
git add api.go
git commit -m "refactor: add authenticated GET helper, use client for all requests"
```

---

### Task 4: Improve POST error messages to include response body

**Files:**
- Modify: `api.go`

**Step 1: Add helper for reading error body**

Add to `api.go`:

```go
// readErrorBody reads up to 512 bytes from a response body for error reporting.
func readErrorBody(resp *http.Response) string {
	if resp.Body == nil {
		return ""
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil || len(body) == 0 {
		return ""
	}
	return string(body)
}
```

Add `"io"` to the import block if not already present.

**Step 2: Update POST error handling**

For each POST function that checks `resp.StatusCode != http.StatusOK`, change the error from:
```go
return fmt.Errorf("bet placement failed with status %d", resp.StatusCode)
```
to:
```go
return fmt.Errorf("bet placement failed with status %d: %s", resp.StatusCode, readErrorBody(resp))
```

Functions to update: `PostBet`, `CancelBet`, `CreateMarket`, `AddLiquidity`, `CloseMarket`, `AddMarketToGroup`, `ResolveMarket`, `SellShares`, `PostComment`.

**Step 3: Run all tests**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 4: Commit**

```bash
git add api.go
git commit -m "feat: include response body in POST error messages"
```

---

### Task 5: Update go.mod and dependencies

**Files:**
- Modify: `go.mod`

**Step 1: Update Go version**

```bash
go mod edit -go=1.21
```

**Step 2: Update dependencies**

```bash
go get -u ./...
go mod tidy
```

**Step 3: Run all tests**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: update Go to 1.21 and refresh dependencies"
```

---

## Phase 2: Domain Endpoints (Parallel)

> Phase 2 agents can start after Phase 1 completes. All three domain agents work in parallel.

### Task 6: Add new endpoint constants

**Files:**
- Modify: `request.go`

**Step 1: Add new route constants**

Add these constants to `request.go`:

```go
const getUserPortfolio string = "get-user-portfolio/"
const getUserPortfolioHistory string = "get-user-portfolio-history/"
const getMarketProb string = "market/"
const getMarketProbs string = "market-probs/"
const getLeagues string = "leagues/"
const getTxns string = "txns/"
const getUserContractMetrics string = "get-user-contract-metrics-with-contracts/"

const postMultiBet string = "multi-bet/"
const postManagram string = "managram/"

const answerSuffix string = "/answer/"
const bountySuffix string = "/add-bounty/"
const awardBountySuffix string = "/award-bounty/"
const liteSuffix string = "/lite/"
const probSuffix string = "/prob/"
```

**Step 2: Commit**

```bash
git add request.go
git commit -m "feat: add route constants for new API endpoints"
```

---

### Task 7: Users & Portfolio - DisplayUser and lite endpoints

**Files:**
- Modify: `user.go`
- Modify: `api.go`

**Step 1: Add DisplayUser struct to user.go**

```go
// DisplayUser represents a lightweight user object with only display information.
//
// Returned by the /lite user endpoints.
type DisplayUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
}
```

**Step 2: Write failing test for GetUserLite**

In `api_test.go`:

```go
func TestGetUserLite(t *testing.T) {
	expected := DisplayUser{
		Id:        "abc123",
		Name:      "Test User",
		Username:  "testuser",
		AvatarUrl: "https://example.com/avatar.png",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetUserLite("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Id != expected.Id || result.Username != expected.Username {
		t.Errorf("got %+v, want %+v", *result, expected)
	}
}
```

**Step 3: Run test to verify it fails**

Run: `go test ./... -run TestGetUserLite -v`
Expected: FAIL - method does not exist

**Step 4: Implement GetUserLite and GetUserByIDLite**

Add to `api.go`:

```go
// GetUserLite returns a [DisplayUser] by username with only basic display fields.
func (mc *Client) GetUserLite(username string) (*DisplayUser, error) {
	resp, err := mc.getRequest(requestURL(mc.url, getUserByUsername, username, liteSuffix))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, DisplayUser{})
}

// GetUserByIDLite returns a [DisplayUser] by user ID with only basic display fields.
func (mc *Client) GetUserByIDLite(id string) (*DisplayUser, error) {
	resp, err := mc.getRequest(requestURL(mc.url, getUserByID, id, liteSuffix))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, DisplayUser{})
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./... -run TestGetUserLite -v`
Expected: PASS

**Step 6: Commit**

```bash
git add user.go api.go api_test.go
git commit -m "feat: add DisplayUser type and lite user endpoints"
```

---

### Task 8: Users & Portfolio - Portfolio endpoints

**Files:**
- Create: `portfolio.go`
- Modify: `api.go`

**Step 1: Create portfolio.go with types**

```go
package mango

// LivePortfolioMetrics represents a user's current portfolio state.
type LivePortfolioMetrics struct {
	InvestmentValue     float64 `json:"investmentValue"`
	CashInvestmentValue float64 `json:"cashInvestmentValue"`
	Balance             float64 `json:"balance"`
	CashBalance         float64 `json:"cashBalance"`
	SpiceBalance        float64 `json:"spiceBalance"`
	TotalDeposits       float64 `json:"totalDeposits"`
	TotalCashDeposits   float64 `json:"totalCashDeposits"`
	LoanTotal           float64 `json:"loanTotal"`
	Timestamp           int64   `json:"timestamp"`
	Profit              float64 `json:"profit,omitempty"`
	UserId              string  `json:"userId"`
	DailyProfit         float64 `json:"dailyProfit"`
}

// PortfolioMetrics represents a snapshot of portfolio metrics at a point in time.
type PortfolioMetrics struct {
	InvestmentValue float64 `json:"investmentValue"`
	Balance         float64 `json:"balance"`
	TotalDeposits   float64 `json:"totalDeposits"`
	LoanTotal       float64 `json:"loanTotal"`
	Timestamp       int64   `json:"timestamp"`
	Profit          float64 `json:"profit,omitempty"`
}

// PortfolioPeriod represents a time period for portfolio history queries.
type PortfolioPeriod string

const (
	PeriodDaily   PortfolioPeriod = "daily"
	PeriodWeekly  PortfolioPeriod = "weekly"
	PeriodMonthly PortfolioPeriod = "monthly"
	PeriodAllTime PortfolioPeriod = "allTime"
)
```

**Step 2: Write failing test for GetUserPortfolio**

In `api_test.go`:

```go
func TestGetUserPortfolio(t *testing.T) {
	expected := LivePortfolioMetrics{
		InvestmentValue: 5000.0,
		Balance:         1000.0,
		UserId:          "user123",
		DailyProfit:     50.0,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("userId") != "user123" {
			t.Errorf("expected userId query param 'user123', got '%s'", r.URL.Query().Get("userId"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetUserPortfolio("user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserId != expected.UserId || result.Balance != expected.Balance {
		t.Errorf("got %+v, want %+v", *result, expected)
	}
}
```

**Step 3: Run test to verify it fails**

Run: `go test ./... -run TestGetUserPortfolio -v`
Expected: FAIL

**Step 4: Implement GetUserPortfolio and GetUserPortfolioHistory**

Add to `api.go`:

```go
// GetUserPortfolio returns a user's current [LivePortfolioMetrics].
func (mc *Client) GetUserPortfolio(userId string) (*LivePortfolioMetrics, error) {
	resp, err := mc.getRequest(requestURL(mc.url, getUserPortfolio, "", "",
		"userId", userId,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, LivePortfolioMetrics{})
}

// GetUserPortfolioHistory returns a slice of [PortfolioMetrics] for a user over a given period.
func (mc *Client) GetUserPortfolioHistory(userId string, period PortfolioPeriod) (*[]PortfolioMetrics, error) {
	resp, err := mc.getRequest(requestURL(mc.url, getUserPortfolioHistory, "", "",
		"userId", userId,
		"period", string(period),
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []PortfolioMetrics{})
}
```

**Step 5: Run tests**

Run: `go test ./... -run TestGetUserPortfolio -v`
Expected: PASS

**Step 6: Commit**

```bash
git add portfolio.go api.go api_test.go
git commit -m "feat: add portfolio types and endpoints"
```

---

### Task 9: Users & Portfolio - GetUserContractMetricsWithContracts

**Files:**
- Modify: `contract.go`
- Modify: `api.go`

**Step 1: Add request/response types to contract.go**

```go
// GetUserContractMetricsRequest represents the parameters for fetching user contract metrics.
type GetUserContractMetricsRequest struct {
	UserId    string `json:"userId"`
	Limit     int64  `json:"limit"`
	Offset    int64  `json:"offset,omitempty"`
	Order     string `json:"order,omitempty"`
	PerAnswer bool   `json:"perAnswer,omitempty"`
}

// UserContractMetricsResponse contains contract metrics grouped by contract ID
// alongside the full contract objects.
type UserContractMetricsResponse struct {
	MetricsByContract map[string][]ContractMetric `json:"metricsByContract"`
	Contracts         []FullMarket                `json:"contracts"`
}
```

**Step 2: Write failing test**

In `api_test.go`:

```go
func TestGetUserContractMetricsWithContracts(t *testing.T) {
	expected := UserContractMetricsResponse{
		MetricsByContract: map[string][]ContractMetric{
			"contract1": {{ContractId: "contract1", Profit: 100.0}},
		},
		Contracts: []FullMarket{{Id: "contract1", Question: "Test?"}},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetUserContractMetricsWithContracts(GetUserContractMetricsRequest{
		UserId: "user1",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Contracts) != 1 {
		t.Errorf("expected 1 contract, got %d", len(result.Contracts))
	}
}
```

**Step 3: Implement the endpoint**

Add to `api.go`:

```go
// GetUserContractMetricsWithContracts returns a user's contract metrics alongside the contracts.
func (mc *Client) GetUserContractMetricsWithContracts(req GetUserContractMetricsRequest) (*UserContractMetricsResponse, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("userId is required")
	}
	if req.Limit == 0 {
		req.Limit = defaultLimit
	}

	var offset, perAnswer string
	if req.Offset > 0 {
		offset = strconv.FormatInt(req.Offset, 10)
	}
	if req.PerAnswer {
		perAnswer = "true"
	}

	resp, err := mc.getRequest(requestURL(mc.url, getUserContractMetrics, "", "",
		"userId", req.UserId,
		"limit", strconv.FormatInt(req.Limit, 10),
		"offset", offset,
		"order", req.Order,
		"perAnswer", perAnswer,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, UserContractMetricsResponse{})
}
```

**Step 4: Run tests and commit**

Run: `go test ./... -run TestGetUserContractMetricsWithContracts -v`
Expected: PASS

```bash
git add contract.go api.go api_test.go
git commit -m "feat: add GetUserContractMetricsWithContracts endpoint"
```

---

### Task 10: Markets - GetMarketProb and GetMarketProbs

**Files:**
- Modify: `market.go`
- Modify: `api.go`

**Step 1: Add response types to market.go**

```go
// MarketProb represents the probability/probabilities for a market.
type MarketProb struct {
	Prob        float64            `json:"prob,omitempty"`
	AnswerProbs map[string]float64 `json:"answerProbs,omitempty"`
}
```

**Step 2: Write failing test for GetMarketProb**

In `api_test.go`:

```go
func TestGetMarketProb(t *testing.T) {
	expected := MarketProb{Prob: 0.75}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetMarketProb("market123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Prob != 0.75 {
		t.Errorf("expected prob 0.75, got %f", result.Prob)
	}
}
```

**Step 3: Implement GetMarketProb and GetMarketProbs**

Add to `api.go`:

```go
// GetMarketProb returns the current probability for a market.
// For binary markets, the Prob field is set.
// For non-binary markets, the AnswerProbs map is set.
func (mc *Client) GetMarketProb(marketId string) (*MarketProb, error) {
	resp, err := mc.getRequest(requestURL(mc.url, getMarketByID, marketId, probSuffix))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, MarketProb{})
}

// GetMarketProbs returns probabilities for multiple markets at once.
// Takes a slice of market IDs (max 100).
func (mc *Client) GetMarketProbs(ids []string) (*map[string]MarketProb, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one market ID is required")
	}
	if len(ids) > 100 {
		return nil, fmt.Errorf("maximum 100 market IDs allowed, got %d", len(ids))
	}

	params := make([]string, 0, len(ids)*2)
	for _, id := range ids {
		params = append(params, "ids", id)
	}

	resp, err := mc.getRequest(requestURL(mc.url, getMarketProbs, "", "", params...))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, map[string]MarketProb{})
}
```

Note: `GetMarketProbs` passes multiple `ids` query params. The `requestURL` helper uses `url.Values.Add()` which supports duplicate keys, so this works correctly.

**Step 4: Run tests and commit**

Run: `go test ./... -run TestGetMarketProb -v`
Expected: PASS

```bash
git add market.go api.go api_test.go
git commit -m "feat: add GetMarketProb and GetMarketProbs endpoints"
```

---

### Task 11: Markets - PostAnswer

**Files:**
- Modify: `api.go`

**Step 1: Write failing test**

In `api_test.go`:

```go
func TestPostAnswer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var parsed map[string]string
		json.Unmarshal(body, &parsed)
		if parsed["text"] != "New answer" {
			t.Errorf("expected text 'New answer', got '%s'", parsed["text"])
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Answer{Id: "ans123", Text: "New answer"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.PostAnswer("market123", "New answer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Id != "ans123" {
		t.Errorf("expected answer ID 'ans123', got '%s'", result.Id)
	}
}
```

**Step 2: Implement PostAnswer**

Add to `api.go`:

```go
// PostAnswer adds a new answer to a multiple choice market.
func (mc *Client) PostAnswer(marketId, text string) (*Answer, error) {
	body := struct {
		Text string `json:"text"`
	}{text}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, postMarket, marketId, answerSuffix), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, Answer{})
}
```

**Step 3: Run tests and commit**

Run: `go test ./... -run TestPostAnswer -v`
Expected: PASS

```bash
git add api.go api_test.go
git commit -m "feat: add PostAnswer endpoint for multiple choice markets"
```

---

### Task 12: Markets - Update SearchMarkets with full params

**Files:**
- Modify: `market.go`
- Modify: `api.go`

**Step 1: Add SearchMarketsRequest to market.go**

```go
// SearchMarketsRequest represents the parameters for searching markets.
type SearchMarketsRequest struct {
	Term         string `json:"term,omitempty"`
	Sort         string `json:"sort,omitempty"`
	Filter       string `json:"filter,omitempty"`
	ContractType string `json:"contractType,omitempty"`
	TopicSlug    string `json:"topicSlug,omitempty"`
	CreatorId    string `json:"creatorId,omitempty"`
	Limit        int64  `json:"limit,omitempty"`
	Offset       int64  `json:"offset,omitempty"`
}
```

**Step 2: Update SearchMarkets in api.go**

Replace the existing `SearchMarkets` method with:

```go
// SearchMarkets returns a slice of [FullMarket] matching the search criteria.
// It takes a [SearchMarketsRequest] with the following optional parameters:
//   - Term - search text
//   - Sort - one of "score", "newest", "liquidity", etc.
//   - Filter - one of "all", "open", "closed", "resolved"
//   - ContractType - one of "ALL", "BINARY", "MULTIPLE_CHOICE", etc.
//   - TopicSlug - filter by topic
//   - CreatorId - filter by market creator
//   - Limit - max results (default 100)
//   - Offset - pagination offset
func (mc *Client) SearchMarkets(req SearchMarketsRequest) (*[]FullMarket, error) {
	var limit, offset string
	if req.Limit > 0 {
		limit = strconv.FormatInt(req.Limit, 10)
	}
	if req.Offset > 0 {
		offset = strconv.FormatInt(req.Offset, 10)
	}

	resp, err := mc.getRequest(requestURL(mc.url, getSearchMarkets, "", "",
		"term", req.Term,
		"sort", req.Sort,
		"filter", req.Filter,
		"contractType", req.ContractType,
		"topicSlug", req.TopicSlug,
		"creatorId", req.CreatorId,
		"limit", limit,
		"offset", offset,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []FullMarket{})
}
```

**Important:** This is a breaking change - the old signature was `SearchMarkets(terms ...string)`. The test `TestSearchMarkets` in `api_test.go` must be updated to use the new signature:

```go
result, err := mc.SearchMarkets(SearchMarketsRequest{Term: "apple banana celery damson"})
```

**Step 3: Update existing test and run**

Run: `go test ./... -run TestSearchMarkets -v`
Expected: PASS

**Step 4: Commit**

```bash
git add market.go api.go api_test.go
git commit -m "feat: update SearchMarkets with full filter/sort params (breaking change)"
```

---

### Task 13: Markets - Bounty endpoints

**Files:**
- Modify: `api.go`

**Step 1: Write failing tests**

In `api_test.go`:

```go
func TestAddBounty(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "txn123"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	err := mc.AddBounty("market123", 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAwardBounty(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "txn456"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	err := mc.AwardBounty("market123", 50, "comment456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

**Step 2: Implement AddBounty and AwardBounty**

Add to `api.go`:

```go
// AddBounty adds mana to a bounty question.
func (mc *Client) AddBounty(marketId string, amount int64) error {
	body := struct {
		Amount int64 `json:"amount"`
	}{amount}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, postMarket, marketId, bountySuffix), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.doRequest(req)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("adding bounty failed with status %d: %s", resp.StatusCode, readErrorBody(resp))
	}

	return nil
}

// AwardBounty distributes a bounty reward to a comment.
func (mc *Client) AwardBounty(marketId string, amount int64, commentId string) error {
	body := struct {
		Amount    int64  `json:"amount"`
		CommentId string `json:"commentId"`
	}{amount, commentId}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, postMarket, marketId, awardBountySuffix), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.doRequest(req)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("awarding bounty failed with status %d: %s", resp.StatusCode, readErrorBody(resp))
	}

	return nil
}
```

**Step 3: Run tests and commit**

Run: `go test ./... -run "TestAddBounty|TestAwardBounty" -v`
Expected: PASS

```bash
git add api.go api_test.go
git commit -m "feat: add AddBounty and AwardBounty endpoints"
```

---

### Task 14: Bets - PostMultiBet

**Files:**
- Modify: `bet.go`
- Modify: `api.go`

**Step 1: Add PostMultiBetRequest to bet.go**

```go
// PostMultiBetRequest represents the parameters for placing multiple YES bets
// on a multiple choice market.
type PostMultiBetRequest struct {
	ContractId string   `json:"contractId"`
	AnswerIds  []string `json:"answerIds"`
	Amount     float64  `json:"amount"`
	LimitProb  *float64 `json:"limitProb,omitempty"`
	ExpiresAt  *int64   `json:"expiresAt,omitempty"`
}
```

**Step 2: Write failing test**

In `api_test.go`:

```go
func TestPostMultiBet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"betId": "bet1", "betGroupId": "bg1"},
			{"betId": "bet2", "betGroupId": "bg1"},
		})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	err := mc.PostMultiBet(PostMultiBetRequest{
		ContractId: "contract123",
		AnswerIds:  []string{"ans1", "ans2"},
		Amount:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

**Step 3: Implement PostMultiBet**

Add to `api.go`:

```go
// PostMultiBet places multiple YES bets on answers in a multiple choice market.
func (mc *Client) PostMultiBet(req PostMultiBetRequest) error {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, postMultiBet, "", ""), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.doRequest(httpReq)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("multi-bet placement failed with status %d: %s", resp.StatusCode, readErrorBody(resp))
	}

	return nil
}
```

**Step 4: Run tests and commit**

Run: `go test ./... -run TestPostMultiBet -v`
Expected: PASS

```bash
git add bet.go api.go api_test.go
git commit -m "feat: add PostMultiBet endpoint for multiple choice markets"
```

---

### Task 15: Bets - Update GetBetsRequest with Kinds field

**Files:**
- Modify: `bet.go`
- Modify: `api.go`

**Step 1: Add Kinds field to GetBetsRequest**

In `bet.go`, add the field:

```go
type GetBetsRequest struct {
	UserId       string `json:"userId,omitempty"`
	Username     string `json:"username,omitempty"`
	ContractId   string `json:"contractID,omitempty"`
	ContractSlug string `json:"contractSlug,omitempty"`
	Before       string `json:"before,omitempty"`
	Limit        int64  `json:"limit,omitempty"`
	Kinds        string `json:"kinds,omitempty"`
}
```

**Step 2: Update GetBets in api.go to pass Kinds**

In the `GetBets` method, add `"kinds", gbr.Kinds` to the `requestURL` params.

**Step 3: Run existing TestGetBets to verify no breakage**

Run: `go test ./... -run TestGetBets -v`
Expected: PASS (Kinds defaults to empty string, which is omitted)

**Step 4: Commit**

```bash
git add bet.go api.go
git commit -m "feat: add Kinds filter to GetBetsRequest"
```

---

### Task 16: Transactions - Txn type and GetTransactions

**Files:**
- Create: `transaction.go`
- Modify: `api.go`

**Step 1: Create transaction.go**

```go
package mango

// Txn represents a transaction in the Manifold system.
type Txn struct {
	Id          string  `json:"id"`
	CreatedTime int64   `json:"createdTime"`
	FromId      string  `json:"fromId"`
	FromType    string  `json:"fromType"`
	ToId        string  `json:"toId"`
	ToType      string  `json:"toType"`
	Amount      float64 `json:"amount"`
	Token       string  `json:"token"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
}

// GetTransactionsRequest represents the parameters for fetching transactions.
type GetTransactionsRequest struct {
	Token    string `json:"token,omitempty"`
	Offset   int64  `json:"offset,omitempty"`
	Limit    int64  `json:"limit,omitempty"`
	Before   int64  `json:"before,omitempty"`
	After    int64  `json:"after,omitempty"`
	ToId     string `json:"toId,omitempty"`
	FromId   string `json:"fromId,omitempty"`
	Category string `json:"category,omitempty"`
}

// SendManagramRequest represents the parameters for sending mana to users.
type SendManagramRequest struct {
	ToIds   []string `json:"toIds"`
	Amount  float64  `json:"amount"`
	Message string   `json:"message,omitempty"`
}
```

**Step 2: Write failing test for GetTransactions**

In `api_test.go`:

```go
func TestGetTransactions(t *testing.T) {
	expected := []Txn{
		{Id: "txn1", Amount: 100, Category: "MANA_PAYMENT"},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetTransactions(GetTransactionsRequest{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*result) != 1 || (*result)[0].Id != "txn1" {
		t.Errorf("unexpected result: %+v", result)
	}
}
```

**Step 3: Implement GetTransactions**

Add to `api.go`:

```go
// GetTransactions returns a slice of [Txn] matching the given filters.
func (mc *Client) GetTransactions(req GetTransactionsRequest) (*[]Txn, error) {
	if req.Limit == 0 {
		req.Limit = 100
	}

	var before, after, offset string
	if req.Before > 0 {
		before = strconv.FormatInt(req.Before, 10)
	}
	if req.After > 0 {
		after = strconv.FormatInt(req.After, 10)
	}
	if req.Offset > 0 {
		offset = strconv.FormatInt(req.Offset, 10)
	}

	resp, err := mc.getRequest(requestURL(mc.url, getTxns, "", "",
		"token", req.Token,
		"offset", offset,
		"limit", strconv.FormatInt(req.Limit, 10),
		"before", before,
		"after", after,
		"toId", req.ToId,
		"fromId", req.FromId,
		"category", req.Category,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []Txn{})
}
```

**Step 4: Run tests and commit**

Run: `go test ./... -run TestGetTransactions -v`
Expected: PASS

```bash
git add transaction.go api.go api_test.go
git commit -m "feat: add Txn type and GetTransactions endpoint"
```

---

### Task 17: Transactions - SendManagram

**Files:**
- Modify: `api.go`

**Step 1: Write failing test**

In `api_test.go`:

```go
func TestSendManagram(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var parsed SendManagramRequest
		json.Unmarshal(body, &parsed)
		if parsed.Amount != 10 || len(parsed.ToIds) != 1 {
			t.Errorf("unexpected request body: %+v", parsed)
		}
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	err := mc.SendManagram(SendManagramRequest{
		ToIds:   []string{"user123"},
		Amount:  10,
		Message: "Test managram",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

**Step 2: Implement SendManagram**

Add to `api.go`:

```go
// SendManagram sends mana to one or more users. Minimum amount is 10.
func (mc *Client) SendManagram(req SendManagramRequest) error {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, requestURL(
		mc.url, postManagram, "", ""), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating http request: %v", err)
	}

	resp, err := mc.doRequest(httpReq)
	if err != nil {
		return fmt.Errorf("error making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sending managram failed with status %d: %s", resp.StatusCode, readErrorBody(resp))
	}

	return nil
}
```

**Step 3: Run tests and commit**

Run: `go test ./... -run TestSendManagram -v`
Expected: PASS

```bash
git add api.go api_test.go
git commit -m "feat: add SendManagram endpoint"
```

---

### Task 18: Transactions - GetLeagues

**Files:**
- Create: `league.go`
- Modify: `api.go`

**Step 1: Create league.go with types**

```go
package mango

// LeagueEntry represents a single entry in league standings.
type LeagueEntry struct {
	UserId          string  `json:"userId"`
	ManaEarned      float64 `json:"manaEarned"`
	ManaEarnedBreakdown map[string]float64 `json:"manaEarnedBreakdown,omitempty"`
	Rank            int     `json:"rank,omitempty"`
	RankSnapshot    int     `json:"rankSnapshot,omitempty"`
	Season          int     `json:"season"`
	Cohort          string  `json:"cohort"`
	Division        int     `json:"division"`
	CreatedTime     int64   `json:"createdTime"`
}

// GetLeaguesRequest represents the parameters for fetching league standings.
type GetLeaguesRequest struct {
	UserId string `json:"userId,omitempty"`
	Season int    `json:"season,omitempty"`
	Cohort string `json:"cohort,omitempty"`
}
```

**Step 2: Write failing test**

In `api_test.go`:

```go
func TestGetLeagues(t *testing.T) {
	expected := []LeagueEntry{
		{UserId: "user1", Season: 1, Division: 3, ManaEarned: 500},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetLeagues(GetLeaguesRequest{UserId: "user1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*result) != 1 || (*result)[0].UserId != "user1" {
		t.Errorf("unexpected result: %+v", result)
	}
}
```

**Step 3: Implement GetLeagues**

Add to `api.go`:

```go
// GetLeagues returns league standings data.
func (mc *Client) GetLeagues(req GetLeaguesRequest) (*[]LeagueEntry, error) {
	var season string
	if req.Season > 0 {
		season = strconv.Itoa(req.Season)
	}

	resp, err := mc.getRequest(requestURL(mc.url, getLeagues, "", "",
		"userId", req.UserId,
		"season", season,
		"cohort", req.Cohort,
	))
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}

	return parseResponse(resp, []LeagueEntry{})
}
```

**Step 4: Run tests and commit**

Run: `go test ./... -run TestGetLeagues -v`
Expected: PASS

```bash
git add league.go api.go api_test.go
git commit -m "feat: add LeagueEntry type and GetLeagues endpoint"
```

---

### Task 19: Run full unit test suite

After all domain tasks are complete, verify everything compiles and passes together.

**Step 1: Build**

Run: `go build ./...`
Expected: No errors

**Step 2: Run all tests**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 3: Commit any fixes if needed**

---

## Phase 3: Integration Tests

### Task 20: Set up integration test infrastructure

**Files:**
- Create: `integration_test.go`

**Step 1: Create the integration test file with build tag and helpers**

Create `integration_test.go`:

```go
//go:build integration

package mango

import (
	"os"
	"testing"
)

// testClient returns a mango client configured for integration testing.
// Requires MANIFOLD_API_KEY env var to be set.
func testClient(t *testing.T) *Client {
	t.Helper()
	key := os.Getenv("MANIFOLD_API_KEY")
	if key == "" {
		t.Skip("MANIFOLD_API_KEY not set, skipping integration test")
	}
	mc := ClientInstance(nil, nil, &key)
	t.Cleanup(func() { mc.Destroy() })
	return mc
}

// Known stable IDs for integration tests.
// These reference long-lived markets/users on Manifold.
const (
	testUsername = "Austin"
	testUserId   = "igi2zGXsfxYPgB0DJTXVJVmwCOr2"
)
```

**Step 2: Commit**

```bash
git add integration_test.go
git commit -m "feat: add integration test infrastructure with build tag"
```

---

### Task 21: Integration tests - Read-only user endpoints

**Files:**
- Modify: `integration_test.go`

**Step 1: Add read-only user tests**

Append to `integration_test.go`:

```go
func TestIntegrationGetUserByUsername(t *testing.T) {
	mc := testClient(t)
	user, err := mc.GetUserByUsername(testUsername)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if user.Username != testUsername {
		t.Errorf("expected username %s, got %s", testUsername, user.Username)
	}
}

func TestIntegrationGetUserByID(t *testing.T) {
	mc := testClient(t)
	user, err := mc.GetUserByID(testUserId)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if user.Id != testUserId {
		t.Errorf("expected id %s, got %s", testUserId, user.Id)
	}
}

func TestIntegrationGetUserLite(t *testing.T) {
	mc := testClient(t)
	user, err := mc.GetUserLite(testUsername)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if user.Username != testUsername {
		t.Errorf("expected username %s, got %s", testUsername, user.Username)
	}
}

func TestIntegrationGetAuthenticatedUser(t *testing.T) {
	mc := testClient(t)
	user, err := mc.GetAuthenticatedUser()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if user.Id == "" {
		t.Error("expected non-empty user ID")
	}
}

func TestIntegrationGetUserPortfolio(t *testing.T) {
	mc := testClient(t)
	user, err := mc.GetAuthenticatedUser()
	if err != nil {
		t.Fatalf("error getting auth user: %v", err)
	}
	portfolio, err := mc.GetUserPortfolio(user.Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if portfolio.UserId != user.Id {
		t.Errorf("expected userId %s, got %s", user.Id, portfolio.UserId)
	}
}

func TestIntegrationGetUsers(t *testing.T) {
	mc := testClient(t)
	users, err := mc.GetUsers(GetUsersRequest{Limit: 5})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(*users) != 5 {
		t.Errorf("expected 5 users, got %d", len(*users))
	}
}
```

**Step 2: Run integration tests**

Run: `MANIFOLD_API_KEY=<key> go test ./... -tags=integration -run TestIntegration -v`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add integration_test.go
git commit -m "feat: add integration tests for user endpoints"
```

---

### Task 22: Integration tests - Read-only market endpoints

**Files:**
- Modify: `integration_test.go`

**Step 1: Add market integration tests**

Append to `integration_test.go`:

```go
func TestIntegrationGetMarkets(t *testing.T) {
	mc := testClient(t)
	markets, err := mc.GetMarkets(GetMarketsRequest{Limit: 5})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(*markets) != 5 {
		t.Errorf("expected 5 markets, got %d", len(*markets))
	}
}

func TestIntegrationSearchMarkets(t *testing.T) {
	mc := testClient(t)
	markets, err := mc.SearchMarkets(SearchMarketsRequest{
		Term:  "bitcoin",
		Limit: 3,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(*markets) == 0 {
		t.Error("expected at least 1 search result for 'bitcoin'")
	}
}

func TestIntegrationGetMarketBySlug(t *testing.T) {
	mc := testClient(t)
	// First get a market to find a valid slug
	markets, err := mc.GetMarkets(GetMarketsRequest{Limit: 1})
	if err != nil {
		t.Fatalf("error getting markets: %v", err)
	}
	if len(*markets) == 0 {
		t.Fatal("no markets returned")
	}
	// Use the slug from a LiteMarket - extract from URL
	market, err := mc.GetMarketByID((*markets)[0].Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if market.Id == "" {
		t.Error("expected non-empty market ID")
	}
}

func TestIntegrationGetMarketProb(t *testing.T) {
	mc := testClient(t)
	markets, err := mc.GetMarkets(GetMarketsRequest{Limit: 1})
	if err != nil {
		t.Fatalf("error getting markets: %v", err)
	}
	if len(*markets) == 0 {
		t.Fatal("no markets returned")
	}
	prob, err := mc.GetMarketProb((*markets)[0].Id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if prob.Prob == 0 && len(prob.AnswerProbs) == 0 {
		t.Error("expected either prob or answerProbs to be set")
	}
}

func TestIntegrationGetBets(t *testing.T) {
	mc := testClient(t)
	bets, err := mc.GetBets(GetBetsRequest{Limit: 5})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(*bets) == 0 {
		t.Error("expected at least 1 bet")
	}
}

func TestIntegrationGetComments(t *testing.T) {
	mc := testClient(t)
	// Use a well-known market with comments
	markets, err := mc.SearchMarkets(SearchMarketsRequest{Term: "manifold", Limit: 1})
	if err != nil || len(*markets) == 0 {
		t.Skip("could not find a market to test comments on")
	}
	comments, err := mc.GetComments(GetCommentsRequest{ContractId: (*markets)[0].Id})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// Just verify we got a valid response (may be empty)
	_ = comments
}

func TestIntegrationGetTransactions(t *testing.T) {
	mc := testClient(t)
	txns, err := mc.GetTransactions(GetTransactionsRequest{Limit: 5})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	_ = txns // just verify no error
}

func TestIntegrationGetLeagues(t *testing.T) {
	mc := testClient(t)
	leagues, err := mc.GetLeagues(GetLeaguesRequest{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	_ = leagues // just verify no error
}
```

**Step 2: Run integration tests**

Run: `MANIFOLD_API_KEY=<key> go test ./... -tags=integration -run TestIntegration -v`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add integration_test.go
git commit -m "feat: add integration tests for market, bet, and transaction endpoints"
```

---

### Task 23: Integration tests - Write endpoints (mana-conscious)

**Files:**
- Modify: `integration_test.go`

**Step 1: Add write tests that share a single test market**

These tests run sequentially and share a market to minimize mana usage. Only one market is created for the entire write test suite.

Append to `integration_test.go`:

```go
// TestIntegrationWriteEndpoints tests all write endpoints using a single shared market.
// Run with: go test -tags=integration -run TestIntegrationWriteEndpoints -v
// WARNING: This test spends mana. Use a test account.
func TestIntegrationWriteEndpoints(t *testing.T) {
	mc := testClient(t)

	// Create a test market (costs mana)
	t.Run("CreateMarket", func(t *testing.T) {
		pmr := PostMarketRequest{
			OutcomeType: Binary,
			Question:    "Mango integration test market - please ignore",
			Description: "Automated test market created by mango library integration tests",
			InitialProb: 50,
		}

		marketId, err := mc.CreateMarket(pmr)
		if err != nil {
			t.Fatalf("error creating market: %v", err)
		}
		if *marketId == "" {
			t.Fatal("expected non-empty market ID")
		}

		// Store marketId for subsequent subtests
		t.Setenv("TEST_MARKET_ID", *marketId)
	})

	marketId := os.Getenv("TEST_MARKET_ID")
	if marketId == "" {
		t.Fatal("no test market ID available")
	}

	t.Run("PostBet", func(t *testing.T) {
		err := mc.PostBet(PostBetRequest{
			Amount:     1, // minimum mana
			ContractId: marketId,
			Outcome:    "YES",
		})
		if err != nil {
			t.Fatalf("error posting bet: %v", err)
		}
	})

	t.Run("PostComment", func(t *testing.T) {
		err := mc.PostComment(marketId, PostCommentRequest{
			ContractId: marketId,
			Markdown:   "Automated test comment from mango integration tests",
		})
		if err != nil {
			t.Fatalf("error posting comment: %v", err)
		}
	})

	t.Run("SellShares", func(t *testing.T) {
		err := mc.SellShares(marketId, SellSharesRequest{
			Outcome: "YES",
		})
		if err != nil {
			t.Fatalf("error selling shares: %v", err)
		}
	})

	t.Run("CloseMarket", func(t *testing.T) {
		err := mc.CloseMarket(marketId, nil)
		if err != nil {
			t.Fatalf("error closing market: %v", err)
		}
	})

	t.Run("ResolveMarket", func(t *testing.T) {
		err := mc.ResolveMarket(marketId, ResolveMarketRequest{
			Outcome: "CANCEL",
		})
		if err != nil {
			t.Fatalf("error resolving market: %v", err)
		}
	})
}
```

**Step 2: Run write integration tests**

Run: `MANIFOLD_API_KEY=<key> go test ./... -tags=integration -run TestIntegrationWriteEndpoints -v`
Expected: ALL PASS (each subtest runs sequentially)

**Step 3: Commit**

```bash
git add integration_test.go
git commit -m "feat: add mana-conscious write integration tests"
```

---

### Task 24: Move existing real-API tests behind build tag

Some existing tests in `api_test.go` hit the real API (e.g., `TestGetBets`, `TestGetComments`, `TestGetGroups`, etc.) without mocks. These should be moved to `integration_test.go` so that `go test ./...` always works without credentials.

**Files:**
- Modify: `api_test.go`
- Modify: `integration_test.go`

**Step 1: Identify tests that use DefaultClientInstance() and hit real API**

These tests call `DefaultClientInstance()` and make real HTTP requests:
- `TestGetBets`
- `TestGetComments`
- `TestGetGroupByID`
- `TestGetGroupBySlug`
- `TestGetGroups`
- `TestGetMarketByID`
- `TestGetMarketBySlug`
- `TestGetMarkets`
- `TestGetMarketsForGroup`
- `TestGetUserByID`
- `TestGetUserByUsername`
- `TestGetUsers`

**Step 2: Move these tests to integration_test.go**

Move each test function from `api_test.go` to `integration_test.go`, renaming them with `TestIntegrationLegacy` prefix and updating them to use `testClient(t)` instead of `DefaultClientInstance()`.

**Step 3: Verify unit tests pass without credentials**

Run: `go test ./... -v`
Expected: ALL PASS (only httptest-based tests remain in api_test.go)

**Step 4: Verify integration tests pass with credentials**

Run: `MANIFOLD_API_KEY=<key> go test ./... -tags=integration -v`
Expected: ALL PASS

**Step 5: Commit**

```bash
git add api_test.go integration_test.go
git commit -m "refactor: move real-API tests behind integration build tag"
```

---

### Task 25: Final verification

**Step 1: Run all unit tests**

Run: `go test ./... -v`
Expected: ALL PASS, no credentials needed

**Step 2: Run all integration tests**

Run: `MANIFOLD_API_KEY=<key> go test ./... -tags=integration -v`
Expected: ALL PASS

**Step 3: Verify build**

Run: `go build ./...`
Expected: Clean build, no errors

**Step 4: Run go vet**

Run: `go vet ./...`
Expected: No issues

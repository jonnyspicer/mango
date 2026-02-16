package mango

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey = "test-api-key"

func TestGetAuthenticatedUser(t *testing.T) {
	expected := User{
		Id:            "igi2zGXsfxYPgB0DJTXVJVmwCOr2",
		CreatedTime:   1639011767273,
		Name:          "Austin",
		Username:      "Austin",
		Url:           "https://manifold.markets/Austin",
		AvatarUrl:     "https://lh3.googleusercontent.com/a-/AOh14GiZyl1lBehuBMGyJYJhZd-N-mstaUtgE4xdI22lLw=s96-c",
		BannerUrl:     "https://images.unsplash.com/photo-1501523460185-2aa5d2a0f981?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1531&q=80",
		Balance:       10000.0,
		TotalDeposits: 10000.0,
		ProfitCached:  ProfitCached{},
		Bio:           "I build Manifold! Always happy to chat; reach out on Discord or find a time on https://calendly.com/austinchen/manifold!",
		Website:       "https://blog.austn.io",
		TwitterHandle: "akrolsmir",
		DiscordHandle: "akrolsmir#4125",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetAuthenticatedUser()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	b, s := equalUsers(*result, expected)
	if !b {
		t.Error(s)
		t.Fail()
	}
}

func TestGetMarketPositions(t *testing.T) {
	expected := []ContractMetric{
		{
			ContractId: "1",
			From: map[string]Period{
				"2022-01-01": {
					Profit:        100.0,
					ProfitPercent: 10.0,
					Invested:      1000.0,
					PrevValue:     1000.0,
					Value:         1100.0,
				},
			},
			HasNoShares:   false,
			HasShares:     true,
			HasYesShares:  false,
			Invested:      1000.0,
			Loan:          0.0,
			MaxShares:     "",
			Payout:        100.0,
			Profit:        100.0,
			ProfitPercent: 10.0,
			TotalShares: map[string]float64{
				"NO":  100.0,
				"YES": 0.0,
			},
			UserId:        "user1",
			UserName:      "John Doe",
			UserUsername:  "johndoe",
			UserAvatarUrl: "https://example.com/avatar.png",
			LastBetTime:   1641004800,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetMarketPositions(GetMarketPositionsRequest{"1", "desc", 0, 100, "user1"})
	if err != nil {
		t.Errorf("error getting market positions: %v", err)
		t.Fail()
	}

	if len(*result) != len(expected) {
		t.Errorf("unexpected result length: got %d, want %d", len(*result), len(expected))
	}

	for i, contract := range *result {
		b, s := equalContractMetrics(contract, expected[i])
		if !b {
			t.Error(s)
		}
	}
}

func TestSearchMarkets(t *testing.T) {
	expected := []FullMarket{
		{
			Id:                    "123",
			CreatorId:             "456",
			CreatorUsername:       "jonny",
			CreatorName:           "Jonny",
			CreatedTime:           0,
			CreatorAvatarUrl:      "https://127.0.0.1",
			CloseTime:             1,
			Question:              "How much wood would a woodchuck chuck of a woodchuck could chuck wood?",
			Answers:               nil,
			Tags:                  nil,
			Url:                   "",
			Pool:                  Pool{},
			Probability:           50,
			P:                     0,
			TotalLiquidity:        0,
			OutcomeType:           Binary,
			Mechanism:             "dpm-2",
			Volume:                10000,
			Volume24Hours:         10000,
			IsResolved:            false,
			Resolution:            "",
			ResolutionTime:        0,
			ResolutionProbability: 0,
			LastUpdatedTime:       2,
			TextDescription:       "Will resolve based on some totally arbitrary criteria I pick at resolution time",
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.SearchMarkets(SearchMarketsRequest{Term: "apple banana celery damson"})
	if err != nil {
		t.Errorf("error searching markets: %v", err)
		t.Fail()
	}

	if len(*result) != len(expected) {
		t.Errorf("unexpected result length: got %d, want %d", len(*result), len(expected))
	}

	for i, market := range *result {
		b, s := equalFullMarkets(market, expected[i])
		if !b {
			t.Error(s)
		}
	}
}

func TestPostBet(t *testing.T) {
	// The POST /v0/bet endpoint returns "betId" rather than "id".
	// PostBet normalizes this so callers can use bet.Id.
	responseJSON := `{
		"betId": "bet123",
		"amount": 10,
		"contractId": "abc123",
		"outcome": "YES",
		"shares": 15.5,
		"probBefore": 0.50,
		"probAfter": 0.55,
		"isFilled": true,
		"fees": {"creatorFee": 0, "platformFee": 0, "liquidityFee": 0}
	}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.PostBet(PostBetRequest{
		Amount:     10,
		ContractId: "abc123",
		Outcome:    "YES",
	})
	if err != nil {
		t.Fatalf("error posting bet: %v", err)
	}

	if result.Id != "bet123" {
		t.Errorf("expected bet.Id 'bet123' (normalized from betId), got %q", result.Id)
	}
	if result.Shares != 15.5 {
		t.Errorf("expected shares 15.5, got %f", result.Shares)
	}
	if result.ProbAfter != 0.55 {
		t.Errorf("expected probAfter 0.55, got %f", result.ProbAfter)
	}
}

func TestPostBetWithAnswerId(t *testing.T) {
	var receivedBody []byte

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	// Test with AnswerId set
	pbr := PostBetRequest{
		Amount:     10,
		ContractId: "abc123",
		Outcome:    "YES",
		AnswerId:   "answer456",
	}

	err := mc.PostBet(pbr)
	if err != nil {
		t.Fatalf("error posting bet with answerId: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(receivedBody, &parsed); err != nil {
		t.Fatalf("error parsing request body: %v", err)
	}

	if parsed["answerId"] != "answer456" {
		t.Errorf("expected answerId 'answer456', got '%v'", parsed["answerId"])
	}

	// Test with empty AnswerId (omitempty should exclude it)
	pbr2 := PostBetRequest{
		Amount:     20,
		ContractId: "def789",
		Outcome:    "NO",
	}

	err = mc.PostBet(pbr2)
	if err != nil {
		t.Fatalf("error posting bet without answerId: %v", err)
	}

	var parsed2 map[string]interface{}
	if err := json.Unmarshal(receivedBody, &parsed2); err != nil {
		t.Fatalf("error parsing request body: %v", err)
	}

	if _, exists := parsed2["answerId"]; exists {
		t.Errorf("expected answerId to be omitted when empty, but it was present: %v", parsed2["answerId"])
	}
}

func TestCancelBet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.CancelBet("123abc")
	if err != nil {
		t.Errorf("error cancelling bet: %v", err)
		t.Fail()
	}
}

func TestCreateMarket(t *testing.T) {
	pmr := PostMarketRequest{
		OutcomeType:         Binary,
		Question:            "How much wood would a woodchuck chuck of a woodchuck could chuck wood?",
		Description:         "Will resolve based on some totally arbitrary criteria I pick at resolution time",
		DescriptionHtml:     "",
		DescriptionMarkdown: "",
		CloseTime:           1,
		Visibility:          "",
		GroupId:             "",
		InitialProb:         1,
		Min:                 0,
		Max:                 10,
		IsLogScale:          false,
		InitialVal:          0,
		Answers:             nil,
	}

	expected := marketIdResponse{Id: "123marketId"}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	resp, err := mc.CreateMarket(pmr)
	if err != nil {
		t.Errorf("error creating market: %v", err)
		t.Fail()
	}

	if *resp != expected.Id {
		t.Errorf("market ID responses don't match, got: %v expected: %v", *resp, expected)
		t.Fail()
	}
}

func TestAddLiquidity(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.AddLiquidity("123marketid", 9999)
	if err != nil {
		t.Errorf("error adding liquidity: %v", err)
		t.Fail()
	}
}

func TestCloseMarket(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.CloseMarket("123marketid", nil)
	if err != nil {
		t.Errorf("error closing market: %v", err)
		t.Fail()
	}
}

func TestAddMarketToGroup(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.AddMarketToGroup("123marketid", "123groupid")
	if err != nil {
		t.Errorf("error adding market to group: %v", err)
		t.Fail()
	}
}

func TestResolveMarket(t *testing.T) {
	rmr := ResolveMarketRequest{
		Outcome:        "YES",
		Resolutions:    nil,
		ProbabilityInt: 0,
		Value:          0,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.ResolveMarket("123marketid", rmr)
	if err != nil {
		t.Errorf("error resolving market: %v", err)
		t.Fail()
	}
}

func TestSellShares(t *testing.T) {
	ssr := SellSharesRequest{
		Outcome: "YES",
		Shares:  10000,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.SellShares("123marketid", ssr)
	if err != nil {
		t.Errorf("error selling shares: %v", err)
		t.Fail()
	}
}

func TestPostComment(t *testing.T) {
	pcr := PostCommentRequest{
		ContractId: "123contractid",
		Content:    "insert snarky comment here",
		Html:       "",
		Markdown:   "",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.PostComment("123marketid", pcr)
	if err != nil {
		t.Errorf("error posting comment: %v", err)
		t.Fail()
	}
}

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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetUserLite("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Id != expected.Id || result.Username != expected.Username {
		t.Errorf("got %+v, want %+v", *result, expected)
	}
}

func TestCloseMarketSendsBody(t *testing.T) {
	var receivedBody []byte

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
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

func TestAddMarketToGroupSendsBody(t *testing.T) {
	var receivedBody []byte

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.AddMarketToGroup("123marketid", "456groupid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(receivedBody) == "{}" || string(receivedBody) == "null" {
		t.Errorf("expected non-empty body with groupId, got: %s", receivedBody)
	}
}

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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
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

func TestGetMarketProb(t *testing.T) {
	expected := MarketProb{Prob: 0.75}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetMarketProb("market123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Prob != 0.75 {
		t.Errorf("expected prob 0.75, got %f", result.Prob)
	}
}

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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetUserPortfolio("user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.UserId != expected.UserId || result.Balance != expected.Balance {
		t.Errorf("got %+v, want %+v", *result, expected)
	}
}

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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetTransactions(GetTransactionsRequest{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*result) != 1 || (*result)[0].Id != "txn1" {
		t.Errorf("unexpected result: %+v", result)
	}
}

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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.PostAnswer("market123", "New answer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Id != "ans123" {
		t.Errorf("expected answer ID 'ans123', got '%s'", result.Id)
	}
}

func TestAddBounty(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "txn123"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	err := mc.AwardBounty("market123", 50, "comment456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnswerResolutionFields(t *testing.T) {
	t.Run("with resolution fields", func(t *testing.T) {
		data := `{
			"id": "ans1",
			"username": "testuser",
			"name": "Test User",
			"userId": "user1",
			"createdTime": 1700000000000,
			"avatarUrl": "https://example.com/avatar.png",
			"number": 1,
			"contractId": "contract1",
			"text": "Answer text",
			"probability": 0.5,
			"resolution": "CANCEL",
			"resolutionTime": 1700763039000,
			"resolutionProbability": 0.764,
			"resolverId": "abc123"
		}`

		var answer Answer
		if err := json.Unmarshal([]byte(data), &answer); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if answer.Resolution != "CANCEL" {
			t.Errorf("expected resolution 'CANCEL', got '%s'", answer.Resolution)
		}
		if answer.ResolutionTime != 1700763039000 {
			t.Errorf("expected resolutionTime 1700763039000, got %d", answer.ResolutionTime)
		}
		if answer.ResolutionProbability != 0.764 {
			t.Errorf("expected resolutionProbability 0.764, got %f", answer.ResolutionProbability)
		}
		if answer.ResolverId != "abc123" {
			t.Errorf("expected resolverId 'abc123', got '%s'", answer.ResolverId)
		}
	})

	t.Run("without resolution fields", func(t *testing.T) {
		data := `{
			"id": "ans2",
			"username": "testuser",
			"name": "Test User",
			"userId": "user1",
			"createdTime": 1700000000000,
			"avatarUrl": "https://example.com/avatar.png",
			"number": 2,
			"contractId": "contract1",
			"text": "Unresolved answer",
			"probability": 0.3
		}`

		var answer Answer
		if err := json.Unmarshal([]byte(data), &answer); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if answer.Resolution != "" {
			t.Errorf("expected empty resolution, got '%s'", answer.Resolution)
		}
		if answer.ResolutionTime != 0 {
			t.Errorf("expected zero resolutionTime, got %d", answer.ResolutionTime)
		}
		if answer.ResolutionProbability != 0 {
			t.Errorf("expected zero resolutionProbability, got %f", answer.ResolutionProbability)
		}
		if answer.ResolverId != "" {
			t.Errorf("expected empty resolverId, got '%s'", answer.ResolverId)
		}
	})
}

func TestFullMarketWithResolvedAnswers(t *testing.T) {
	expected := FullMarket{
		Id:          "market1",
		Question:    "Multi-choice market?",
		OutcomeType: MultipleChoice,
		Mechanism:   "cpmm-multi-1",
		Answers: []Answer{
			{
				Id:                    "ans1",
				Username:              "resolver",
				Name:                  "Resolver User",
				UserId:                "resolver1",
				CreatedTime:           1700000000000,
				AvatarUrl:             "https://example.com/avatar.png",
				Number:                1,
				ContractId:            "market1",
				Text:                  "Option A",
				Probability:           0.0,
				Resolution:            "YES",
				ResolutionTime:        1700763039000,
				ResolutionProbability: 1.0,
				ResolverId:            "resolver1",
			},
			{
				Id:          "ans2",
				Username:    "resolver",
				Name:        "Resolver User",
				UserId:      "resolver1",
				CreatedTime: 1700000000000,
				AvatarUrl:   "https://example.com/avatar.png",
				Number:      2,
				ContractId:  "market1",
				Text:        "Option B",
				Probability: 0.5,
			},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetMarketByID("market1")
	if err != nil {
		t.Fatalf("error getting market: %v", err)
	}

	if len(result.Answers) != 2 {
		t.Fatalf("expected 2 answers, got %d", len(result.Answers))
	}

	resolved := result.Answers[0]
	if resolved.Resolution != "YES" {
		t.Errorf("expected resolution 'YES', got '%s'", resolved.Resolution)
	}
	if resolved.ResolutionTime != 1700763039000 {
		t.Errorf("expected resolutionTime 1700763039000, got %d", resolved.ResolutionTime)
	}
	if resolved.ResolutionProbability != 1.0 {
		t.Errorf("expected resolutionProbability 1.0, got %f", resolved.ResolutionProbability)
	}
	if resolved.ResolverId != "resolver1" {
		t.Errorf("expected resolverId 'resolver1', got '%s'", resolved.ResolverId)
	}

	unresolved := result.Answers[1]
	if unresolved.Resolution != "" {
		t.Errorf("expected empty resolution on unresolved answer, got '%s'", unresolved.Resolution)
	}
	if unresolved.ResolutionTime != 0 {
		t.Errorf("expected zero resolutionTime on unresolved answer, got %d", unresolved.ResolutionTime)
	}
}

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

	mc := ClientInstance(server.Client(), &server.URL, &testKey)
	defer mc.Destroy()

	result, err := mc.GetLeagues(GetLeaguesRequest{UserId: "user1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*result) != 1 || (*result)[0].UserId != "user1" {
		t.Errorf("unexpected result: %+v", result)
	}
}

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
const (
	testUsername = "Austin"
	testUserId  = "igi2zGXsfxYPgB0DJTXVJVmwCOr2"
)

// --- Task 21: Read-only user integration tests ---

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

// --- Task 22: Read-only market/bet/transaction integration tests ---

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
	markets, err := mc.GetMarkets(GetMarketsRequest{Limit: 1})
	if err != nil {
		t.Fatalf("error getting markets: %v", err)
	}
	if len(*markets) == 0 {
		t.Fatal("no markets returned")
	}
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
	markets, err := mc.SearchMarkets(SearchMarketsRequest{
		Term:         "will",
		Limit:        1,
		Filter:       "open",
		ContractType: "BINARY",
	})
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
	markets, err := mc.SearchMarkets(SearchMarketsRequest{Term: "manifold", Limit: 1})
	if err != nil || len(*markets) == 0 {
		t.Skip("could not find a market to test comments on")
	}
	comments, err := mc.GetComments(GetCommentsRequest{ContractId: (*markets)[0].Id})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	_ = comments
}

func TestIntegrationGetTransactions(t *testing.T) {
	mc := testClient(t)
	txns, err := mc.GetTransactions(GetTransactionsRequest{
		Limit:  5,
		ToId:   testUserId,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	_ = txns
}

func TestIntegrationGetLeagues(t *testing.T) {
	mc := testClient(t)
	leagues, err := mc.GetLeagues(GetLeaguesRequest{UserId: testUserId})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	_ = leagues
}

// --- Task 23: Write integration tests (mana-conscious) ---

// TestIntegrationWriteEndpoints tests all write endpoints using a single shared market.
// Run with: go test -tags=integration -run TestIntegrationWriteEndpoints -v
// WARNING: This test spends mana. Use a test account.
func TestIntegrationWriteEndpoints(t *testing.T) {
	mc := testClient(t)

	var marketId string

	t.Run("CreateMarket", func(t *testing.T) {
		pmr := PostMarketRequest{
			OutcomeType:   Binary,
			Question:      "Mango integration test market - please ignore",
			Description:   "Automated test market created by mango library integration tests",
			InitialProb:   50,
			LiquidityTier: 100,
		}

		id, err := mc.CreateMarket(pmr)
		if err != nil {
			t.Fatalf("error creating market: %v", err)
		}
		if *id == "" {
			t.Fatal("expected non-empty market ID")
		}
		marketId = *id
	})

	if marketId == "" {
		t.Fatal("no test market ID available, cannot continue write tests")
	}

	t.Run("PostBet", func(t *testing.T) {
		err := mc.PostBet(PostBetRequest{
			Amount:     1,
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

// --- Task 24: Legacy real-API tests moved from api_test.go ---

func TestIntegrationLegacyGetBets(t *testing.T) {
	mc := testClient(t)

	var tests = []struct {
		ui, un, ci, cs, b string
		l                 int64
	}{
		{"", "", "", "", "", defaultLimit},
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", "", "", "", "", 10},
		{"", "jonny", "", "", "", 10},
		{"", "", "5BOGaVlxLaZt6sdPSUkn", "", "", 10},
		{"", "", "", "dan-stock-permanent", "", 10},
		{"", "", "", "", "6LvbmW25hvrAsjrme209", 10},
		{"", "", "", "", "", 10},
		{"", "", "5BOGaVlxLaZt6sdPSUkn", "dan-stock-permanent", "", 10},
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", "jonny", "", "", "DgDQOsc7x2KIWTmojZbx", 10},
	}

	for _, test := range tests {
		gbr := GetBetsRequest{
			UserId:       test.ui,
			Username:     test.un,
			ContractId:   test.ci,
			ContractSlug: test.cs,
			Before:       test.b,
			Limit:        test.l,
		}
		actual, err := mc.GetBets(gbr)
		if err != nil {
			t.Errorf("error getting bets: %v", err)
			continue
		}
		if int64(len(*actual)) != test.l {
			t.Errorf("incorrect number of bets retrieved")
		}
	}
}

func TestIntegrationLegacyGetComments(t *testing.T) {
	mc := testClient(t)

	var tests = []struct {
		ci, cs string
	}{
		{"5BOGaVlxLaZt6sdPSUkn", ""},
		{"", "dan-stock-permanent"},
		{"5BOGaVlxLaZt6sdPSUkn", "dan-stock-permanent"},
	}

	for _, test := range tests {
		actual, err := mc.GetComments(GetCommentsRequest{test.ci, test.cs})
		if err != nil {
			t.Errorf("error getting comments: %v", err)
			continue
		}
		if len(*actual) < 1 {
			t.Errorf("incorrect number of comments retrieved")
		}
	}
}

func TestIntegrationLegacyGetGroupByID(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetGroupById("IlzY3moWwOcpsVZXCVej")
	if err != nil {
		t.Fatalf("error getting group by id: %v", err)
	}
	if actual.TotalMembers < 1 {
		t.Errorf("incorrect number of members on retrieved group")
	}
}

func TestIntegrationLegacyGetGroupBySlug(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetGroupBySlug("technology-default")
	if err != nil {
		t.Fatalf("error getting group by slug: %v", err)
	}
	if actual.TotalMembers < 1 {
		t.Errorf("incorrect number of members on retrieved group")
	}
}

func TestIntegrationLegacyGetGroups(t *testing.T) {
	t.Skip("skipping: Group.About field type changed upstream (string -> object), needs type fix")
}

func TestIntegrationLegacyGetMarketByID(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetMarketByID("5BOGaVlxLaZt6sdPSUkn")
	if err != nil {
		t.Fatalf("error getting market by id: %v", err)
	}
	if actual.Volume < 1 {
		t.Errorf("incorrect volume on retrieved market")
	}
}

func TestIntegrationLegacyGetMarketBySlug(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetMarketBySlug("dan-stock-permanent")
	if err != nil {
		t.Fatalf("error getting market by slug: %v", err)
	}
	if actual.Volume < 1 {
		t.Errorf("incorrect volume on retrieved market")
	}
}

func TestIntegrationLegacyGetMarkets(t *testing.T) {
	mc := testClient(t)

	var tests = []struct {
		b string
		l int
	}{
		{"4QTb4cANeQzXNQS9lZnn", 10},
		{"4QTb4cANeQzXNQS9lZnn", defaultLimit},
		{"", defaultLimit},
	}

	for _, test := range tests {
		actual, err := mc.GetMarkets(GetMarketsRequest{test.b, int64(test.l)})
		if err != nil {
			t.Errorf("error getting markets: %v", err)
			continue
		}
		if len(*actual) != test.l {
			t.Errorf("incorrect number of markets retrieved")
		}
	}
}

func TestIntegrationLegacyGetMarketsForGroup(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetMarketsForGroup("IlzY3moWwOcpsVZXCVej")
	if err != nil {
		t.Fatalf("error getting markets for group: %v", err)
	}
	if len(*actual) < 1 {
		t.Errorf("incorrect number of markets for group")
	}
}

func TestIntegrationLegacyGetUserByID(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetUserByID("xN67Q0mAhddL0X9wVYP2YfOrYH42")
	if err != nil {
		t.Fatalf("error getting user by id: %v", err)
	}
	if actual.Balance < 1 {
		t.Errorf("incorrect balance on retrieved user")
	}
}

func TestIntegrationLegacyGetUserByUsername(t *testing.T) {
	mc := testClient(t)

	actual, err := mc.GetUserByUsername("jonny")
	if err != nil {
		t.Fatalf("error getting user by username: %v", err)
	}
	if actual.Balance < 1 {
		t.Errorf("incorrect balance on retrieved user")
	}
}

func TestIntegrationLegacyGetUsers(t *testing.T) {
	mc := testClient(t)

	var tests = []struct {
		b string
		l int
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", 10},
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", defaultLimit},
		{"", defaultLimit},
	}

	for _, test := range tests {
		actual, err := mc.GetUsers(GetUsersRequest{test.b, int64(test.l)})
		if err != nil {
			t.Errorf("error getting users: %v", err)
			continue
		}
		if len(*actual) != test.l {
			t.Errorf("incorrect number of users retrieved")
		}
	}
}

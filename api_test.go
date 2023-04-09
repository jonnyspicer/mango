package mango

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: extract non-deterministic tests out as e2e tests, replace them with deterministic unit tests
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

	mc := ClientInstance(server.Client(), &server.URL, nil)
	defer mc.Destroy()

	result, err := mc.GetAuthenticatedUser()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	b, s := equalUsers(*result, expected)
	if !b {
		t.Errorf(s)
		t.Fail()
	}
}

func TestGetBets(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		ui, un, ci, cs, b string
		l                 int64
	}{
		{"", "", "", "", "", DefaultLimit},
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
			test.ui, test.un, test.ci, test.cs, test.b, test.l
		}
		actual, err := mc.GetBets(gbr)
		if err != nil {
			t.Errorf("error getting bets: %v", err)
		}
		if int64(len(*actual)) != test.l {
			t.Errorf("incorrect number of bets retrieved")
			t.Fail()
		}
	}
}

func TestGetComments(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		ci, cs string
	}{
		{"5BOGaVlxLaZt6sdPSUkn", ""},
		{"", "dan-stock-permanent"},
		{"5BOGaVlxLaZt6sdPSUkn", "dan-stock-permanent"},
	}

	for _, test := range tests {
		actual := mc.GetComments(test.ci, test.cs)
		if len(actual) < 1 {
			t.Errorf("incorrect number of comments retrieved")
			t.Fail()
		}
	}
}

func TestGetGroupByID(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		s string
	}{
		{"IlzY3moWwOcpsVZXCVej"},
	}

	for _, test := range tests {
		actual := mc.GetGroupById(test.s)
		if actual.TotalMembers < 1 {
			t.Errorf("incorrect number of members on retrieved group")
			t.Fail()
		}
	}
}

func TestGetGroupBySlug(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		s string
	}{
		{"technology-default"},
	}

	for _, test := range tests {
		actual := mc.GetGroupBySlug(test.s)
		if actual.TotalMembers < 1 {
			t.Errorf("incorrect number of members on retrieved group")
			t.Fail()
		}
	}
}

func TestGetGroups(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		ui string
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42"},
		{""},
	}

	for _, test := range tests {
		actual := mc.GetGroups(test.ui)
		if len(actual) < 1 {
			t.Errorf("incorrect number of groups retrieved")
			t.Fail()
		}
	}
}

func TestGetMarketByID(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		mi string
	}{
		{"5BOGaVlxLaZt6sdPSUkn"},
	}

	for _, test := range tests {
		actual := mc.GetMarketByID(test.mi)
		if actual.Volume < 1 {
			t.Errorf("incorrect volume on retrieved market")
			t.Fail()
		}
	}
}

func TestGetMarketBySlug(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		ms string
	}{
		{"dan-stock-permanent"},
	}

	for _, test := range tests {
		actual := mc.GetMarketBySlug(test.ms)
		if actual.Volume < 1 {
			t.Errorf("incorrect volume on retrieved market")
			t.Fail()
		}
	}
}

func TestGetMarkets(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		b string
		l int
	}{
		{"4QTb4cANeQzXNQS9lZnn", 10},
		{"4QTb4cANeQzXNQS9lZnn", DefaultLimit},
		{"", DefaultLimit},
	}

	for _, test := range tests {
		actual := mc.GetMarkets(test.b, test.l)
		if len(actual) != test.l {
			t.Errorf("incorrect number of markets retrieved")
			t.Fail()
		}
	}
}

func TestGetMarketsForGroup(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		gi string
	}{
		{"IlzY3moWwOcpsVZXCVej"},
	}

	for _, test := range tests {
		actual := mc.GetMarketsForGroup(test.gi)
		if len(actual) < 1 {
			t.Errorf("incorrect number of markets for group")
			t.Fail()
		}
	}
}

func TestGetUserByID(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		ui string
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42"},
	}

	for _, test := range tests {
		actual := mc.GetUserByID(test.ui)
		if actual.Balance < 1 {
			t.Errorf("incorrect balance on retrieved user")
			t.Fail()
		}
	}
}

func TestGetUserByUsername(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		un string
	}{
		{"jonny"},
	}

	for _, test := range tests {
		actual := mc.GetUserByUsername(test.un)
		if actual.Balance < 1 {
			t.Errorf("incorrect balance on retrieved user")
			t.Fail()
		}
	}
}

func TestGetUsers(t *testing.T) {
	mc := DefaultClientInstance()
	defer mc.Destroy()

	var tests = []struct {
		b string
		l int
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", 10},
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", DefaultLimit},
		{"", DefaultLimit},
	}

	for _, test := range tests {
		actual := mc.GetUsers(test.b, test.l)
		if len(actual) != test.l {
			t.Errorf("incorrect number of users retrieved")
			t.Fail()
		}
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

	mc := ClientInstance(server.Client(), &server.URL)
	defer mc.Destroy()

	result := mc.GetMarketPositions("1", "desc", nil, nil, "user1")

	if len(result) != len(expected) {
		t.Errorf("unexpected result length: got %d, want %d", len(result), len(expected))
	}

	for i, contract := range result {
		b, s := equalContractMetrics(contract, expected[i])
		if !b {
			t.Errorf(s)
		}
	}
}

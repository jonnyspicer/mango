package mango

import (
	"encoding/json"
	"github.com/jonnyspicer/mango/endpoint"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetBets(t *testing.T) {
	mc := defaultManiClient()
	defer mc.destroy()

	var tests = []struct {
		ui, un, ci, cs, b string
		l                 int
	}{
		{"", "", "", "", "", endpoint.DefaultLimit},
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
		actual := mc.GetBets(test.ui, test.un, test.ci, test.cs, test.b, test.l)
		if len(actual) != test.l {
			t.Errorf("incorrect number of bets retrieved")
			t.Fail()
		}
	}
}

func TestGetComments(t *testing.T) {
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

	var tests = []struct {
		s string
	}{
		{"IlzY3moWwOcpsVZXCVej"},
	}

	for _, test := range tests {
		actual := mc.GetGroupByID(test.s)
		if actual.TotalMembers < 1 {
			t.Errorf("incorrect number of members on retrieved group")
			t.Fail()
		}
	}
}

func TestGetGroupBySlug(t *testing.T) {
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

	var tests = []struct {
		b string
		l int
	}{
		{"4QTb4cANeQzXNQS9lZnn", 10},
		{"4QTb4cANeQzXNQS9lZnn", endpoint.DefaultLimit},
		{"", endpoint.DefaultLimit},
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
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

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
	mc := defaultManiClient()
	defer mc.destroy()

	var tests = []struct {
		b string
		l int
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", 10},
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", endpoint.DefaultLimit},
		{"", endpoint.DefaultLimit},
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

	mc := maniClientInstance(server.Client(), &server.URL)
	defer mc.destroy()

	result := mc.GetMarketPositions("1", "desc", nil, nil, "user1")

	if len(result) != len(expected) {
		t.Errorf("unexpected result length: got %d, want %d", len(result), len(expected))
	}

	for i, contract := range result {
		b, s := contract.Equals(expected[i])
		if !b {
			t.Errorf(s)
		}
	}
}

type mockTransport struct {
	response string
	status   int
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := http.Response{
		StatusCode: t.status,
		Body:       io.NopCloser(strings.NewReader(t.response)),
	}

	return &resp, nil
}

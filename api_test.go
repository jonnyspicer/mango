package mango

import (
	"github.com/jonnyspicer/mango/endpoint"
	"testing"
)

func TestGetBets(t *testing.T) {
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
		actual := GetBets(test.ui, test.un, test.ci, test.cs, test.b, test.l)
		if len(actual) != test.l {
			t.Errorf("incorrect number of bets retrieved")
			t.Fail()
		}
	}
}

func TestGetComments(t *testing.T) {
	var tests = []struct {
		ci, cs string
	}{
		{"5BOGaVlxLaZt6sdPSUkn", ""},
		{"", "dan-stock-permanent"},
		{"5BOGaVlxLaZt6sdPSUkn", "dan-stock-permanent"},
	}

	for _, test := range tests {
		actual := GetComments(test.ci, test.cs)
		if len(actual) < 1 {
			t.Errorf("incorrect number of comments retrieved")
			t.Fail()
		}
	}
}

func TestGetGroupByID(t *testing.T) {
	var tests = []struct {
		s string
	}{
		{"IlzY3moWwOcpsVZXCVej"},
	}

	for _, test := range tests {
		actual := GetGroupByID(test.s)
		if actual.TotalMembers < 1 {
			t.Errorf("incorrect number of members on retrieved group")
			t.Fail()
		}
	}
}

func TestGetGroupBySlug(t *testing.T) {
	var tests = []struct {
		s string
	}{
		{"technology-default"},
	}

	for _, test := range tests {
		actual := GetGroupBySlug(test.s)
		if actual.TotalMembers < 1 {
			t.Errorf("incorrect number of members on retrieved group")
			t.Fail()
		}
	}
}

func TestGetGroups(t *testing.T) {
	var tests = []struct {
		ui string
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42"},
		{""},
	}

	for _, test := range tests {
		actual := GetGroups(test.ui)
		if len(actual) < 1 {
			t.Errorf("incorrect number of groups retrieved")
			t.Fail()
		}
	}
}

func TestGetMarketByID(t *testing.T) {
	var tests = []struct {
		mi string
	}{
		{"5BOGaVlxLaZt6sdPSUkn"},
	}

	for _, test := range tests {
		actual := GetMarketByID(test.mi)
		if actual.Volume < 1 {
			t.Errorf("incorrect volume on retrieved market")
			t.Fail()
		}
	}
}

func TestGetMarketBySlug(t *testing.T) {
	var tests = []struct {
		ms string
	}{
		{"dan-stock-permanent"},
	}

	for _, test := range tests {
		actual := GetMarketBySlug(test.ms)
		if actual.Volume < 1 {
			t.Errorf("incorrect volume on retrieved market")
			t.Fail()
		}
	}
}

func TestGetMarkets(t *testing.T) {
	var tests = []struct {
		b string
		l int
	}{
		{"4QTb4cANeQzXNQS9lZnn", 10},
		{"4QTb4cANeQzXNQS9lZnn", endpoint.DefaultLimit},
		{"", endpoint.DefaultLimit},
	}

	for _, test := range tests {
		actual := GetMarkets(test.b, test.l)
		if len(actual) != test.l {
			t.Errorf("incorrect number of markets retrieved")
			t.Fail()
		}
	}
}

func TestGetMarketsForGroup(t *testing.T) {
	var tests = []struct {
		gi string
	}{
		{"IlzY3moWwOcpsVZXCVej"},
	}

	for _, test := range tests {
		actual := GetMarketsForGroup(test.gi)
		if len(actual) < 1 {
			t.Errorf("incorrect number of markets for group")
			t.Fail()
		}
	}
}

func TestGetUserByID(t *testing.T) {
	var tests = []struct {
		ui string
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42"},
	}

	for _, test := range tests {
		actual := GetUserByID(test.ui)
		if actual.Balance < 1 {
			t.Errorf("incorrect balance on retrieved user")
			t.Fail()
		}
	}
}

func TestGetUserByUsername(t *testing.T) {
	var tests = []struct {
		un string
	}{
		{"jonny"},
	}

	for _, test := range tests {
		actual := GetUserByUsername(test.un)
		if actual.Balance < 1 {
			t.Errorf("incorrect balance on retrieved user")
			t.Fail()
		}
	}
}

func TestGetUsers(t *testing.T) {
	var tests = []struct {
		b string
		l int
	}{
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", 10},
		{"xN67Q0mAhddL0X9wVYP2YfOrYH42", endpoint.DefaultLimit},
		{"", endpoint.DefaultLimit},
	}

	for _, test := range tests {
		actual := GetUsers(test.b, test.l)
		if len(actual) != test.l {
			t.Errorf("incorrect number of users retrieved")
			t.Fail()
		}
	}
}

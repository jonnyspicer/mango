package mango

import "testing"

func TestKellyBet(t *testing.T) {
	tests := []struct {
		name     string
		prob     float64
		payout   float64
		expected float64
	}{
		{"standard case", 0.6, 2.0, 0.2},
		{"higher prob and payout", 0.7, 3.0, 0.55},
		{"low payout margin", 0.55, 1.5, -0.35},
		{"break even", 0.5, 2.0, 0.0},
		{"zero net odds", 0.5, 1.0, 0.0},
		{"payout less than 1", 0.5, 0.5, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KellyBet(tt.prob, tt.payout)
			if result != tt.expected {
				t.Errorf("KellyBet(%v, %v) = %v, want %v", tt.prob, tt.payout, result, tt.expected)
			}
		})
	}
}

func TestGetUsernames(t *testing.T) {
	text := `export const BOT_USERNAMES = [
		'pos',
		'v',
		'acc',
		'jerk',
	]
	export const CORE_USERNAMES = [
		'Austin',
		'JamesGrugett',
		'SG',
	]
	export const CHECK_USERNAMES = [
		'EliezerYudkowsky',
		'memestiny',
		'ScottAlexander',
	]`

	expectedBotUsernames := []string{"pos", "v", "acc", "jerk"}
	expectedCoreUsernames := []string{"Austin", "JamesGrugett", "SG"}
	expectedCheckUsernames := []string{"EliezerYudkowsky", "memestiny", "ScottAlexander"}

	botUsernames := getUsernames(Bot, text)
	coreUsernames := getUsernames(Core, text)
	checkUsernames := getUsernames(Check, text)

	if !stringSliceEqual(botUsernames, expectedBotUsernames) {
		t.Errorf("Expected BOT_USERNAMES %v, but got %v", expectedBotUsernames, botUsernames)
	}

	if !stringSliceEqual(coreUsernames, expectedCoreUsernames) {
		t.Errorf("Expected CORE_USERNAMES %v, but got %v", expectedCoreUsernames, coreUsernames)
	}

	if !stringSliceEqual(checkUsernames, expectedCheckUsernames) {
		t.Errorf("Expected CHECK_USERNAMES %v, but got %v", expectedCheckUsernames, checkUsernames)
	}
}

// Helper function to compare string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

package mango

import "testing"

func TestKellyBet(t *testing.T) {
	prob := 0.6
	payout := 2.0
	expected := 0.2

	result := KellyBet(prob, payout)

	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
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

package util

import (
	"strings"
	"testing"
)

func TestRandomInt(t *testing.T) {
	min, max := int64(10), int64(20)
	n := RandomInt(min, max)
	if n < min || n > max {
		t.Errorf("RandomInt(%d, %d) = %d; want between %d and %d", min, max, n, min, max)
	}
}

func TestRandomString(t *testing.T) {
	n := 8
	s := RandomString(n)
	if len(s) != n {
		t.Errorf("RandomString(%d) = %q; length = %d; want %d", n, s, len(s), n)
	}
}

func TestRandomOwner(t *testing.T) {
	owner := RandomOwner()
	if len(owner) != 10 {
		t.Errorf("RandomOwner() = %q; length = %d; want 10", owner, len(owner))
	}
}

func TestRandomMoney(t *testing.T) {
	money := RandomMoney()
	if money < 0 || money > 1000 {
		t.Errorf("RandomMoney() = %d; want between 0 and 1000", money)
	}
}

func TestRandomCurrency(t *testing.T) {
	validCurrencies := map[string]bool{
		SGD: true,
		IDR: true,
		USD: true,
	}
	curr := RandomCurrency()
	if !validCurrencies[curr] {
		t.Errorf("RandomCurrency() = %q; want one of SGD, IDR, USD", curr)
	}
}

func TestRandomEmail(t *testing.T) {
	email := RandomEmail()
	if !strings.Contains(email, "@") {
		t.Errorf("RandomEmail() = %q; missing '@'", email)
	}
	if !strings.HasSuffix(email, "@umarhadi.dev") {
		t.Errorf("RandomEmail() = %q; expected suffix '@umarhadi.dev'", email)
	}
}

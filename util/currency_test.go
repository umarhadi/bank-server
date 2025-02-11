package util

import "testing"

func TestIsSupportedCurrency(t *testing.T) {
	cases := []struct {
		name      string
		currency  string
		supported bool
	}{
		{"Supported SGD", SGD, true},
		{"Supported IDR", IDR, true},
		{"Supported USD", USD, true},
		{"Unsupported EUR", "EUR", false},
		{"Unsupported JPY", "JPY", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsSupportedCurrency(tc.currency)
			if got != tc.supported {
				t.Errorf("IsSupportedCurrency(%q) = %v; want %v", tc.currency, got, tc.supported)
			}
		})
	}
}

package val

import (
	"testing"
)

func TestValidateString(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		min, max  int
		expectErr bool
	}{
		{"empty string", "", 3, 10, true},
		{"too short", "ab", 3, 5, true},
		{"too long", "abcdefghijk", 3, 10, true},
		{"valid", "hello", 3, 10, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateString(tc.value, tc.min, tc.max)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"valid username", "user_123", false},
		{"invalid due to capitals", "User_123", true},
		{"invalid due to symbols", "user-123", true},
		{"too short", "ab", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUsername(tc.value)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateFullName(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"valid full name", "John Doe", false},
		{"invalid with numbers", "John Doe2", true},
		{"invalid symbol", "John@Doe", true},
		{"too short", "Jo", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFullName(tc.value)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"valid password", "secret123", false},
		{"too short", "123", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.value)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"valid email", "user@example.com", false},
		{"missing @", "userexample.com", true},
		{"too short", "a@b", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.value)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateEmailId(t *testing.T) {
	cases := []struct {
		name      string
		value     int64
		expectErr bool
	}{
		{"valid id", 10, false},
		{"zero id", 0, true},
		{"negative id", -5, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmailId(tc.value)
			if tc.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateSecretCode(t *testing.T) {
	// valid secret code length between 32 and 128 characters.
	validCode := "abcdefghijklmnopqrstuvwxyz123456" // length 32
	err := ValidateSecretCode(validCode)
	if err != nil {
		t.Errorf("expected valid secret code, got error: %v", err)
	}

	// too short
	shortCode := "short"
	if err := ValidateSecretCode(shortCode); err == nil {
		t.Errorf("expected error for short secret code, got nil")
	}

	// too long
	longCode := make([]byte, 129)
	for i := range longCode {
		longCode[i] = 'a'
	}
	if err := ValidateSecretCode(string(longCode)); err == nil {
		t.Errorf("expected error for long secret code, got nil")
	}
}

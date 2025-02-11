package api

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestValidCurrency(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("currency", validCurrency)

	testCases := []struct {
		name      string
		currency  interface{}
		expectErr bool
	}{
		{
			name:      "valid currency",
			currency:  "USD",
			expectErr: false,
		},
		{
			name:      "invalid currency",
			currency:  "XYZ",
			expectErr: true,
		},
		{
			name:      "empty currency",
			currency:  "",
			expectErr: true,
		},
		{
			name:      "non-string currency",
			currency:  123,
			expectErr: true,
		},
	}

	type TestStruct struct {
		Currency interface{} `json:"currency" validate:"currency"`
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test := TestStruct{
				Currency: tc.currency,
			}

			err := validate.Struct(test)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

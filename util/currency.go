package util

const (
	SGD = "SGD"
	IDR = "IDR"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case SGD, IDR, USD:
		return true
	}
	return false
}

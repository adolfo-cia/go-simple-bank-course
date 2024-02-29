package utils

const (
	USD = "USD"
	EUR = "EUR"
	ARS = "ARS"
)

func IsSupportedCurrency(curreny string) bool {
	switch curreny {
	case USD, EUR, ARS:
		return true
	}
	return false
}

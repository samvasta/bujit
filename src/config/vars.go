package config

const version string = "0.0.1"

func Version() string {
	return version
}

var currencySymbol string = "$"

func CurrencySymbol() string {
	return currencySymbol
}
func SetCurrencySymbol(symbol string) {
	currencySymbol = symbol
}

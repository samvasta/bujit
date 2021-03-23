package models

import (
	"fmt"
	"math"

	"github.com/dustin/go-humanize"
	"samvasta.com/bujit/session"
	"samvasta.com/bujit/util"
)

type Money int64

func MakeMoney(value float64) Money {
	dollars, cents := math.Modf(value)

	return Money(int(dollars*100) + int(cents*100))
}

func (m Money) Cents() int64 {
	return int64(m) % 100
}

func (m Money) Dollars() int64 {
	return int64(m) / 100
}

func (m Money) IsNegative() bool {
	return int64(m) < 0
}

func (m Money) Value() int64 {
	return int64(m)
}

func (m Money) String(s *session.Session) string {
	if m.IsNegative() {
		return fmt.Sprintf("(%s%s.%d)%s", s.CurrencyPrefix, humanize.Comma(util.AbsI64(m.Dollars())), m.Cents(), s.CurrencySuffixWithSpace())
	} else {
		return fmt.Sprintf("%s%s.%d%s", s.CurrencyPrefix, humanize.Comma(m.Dollars()), m.Cents(), s.CurrencySuffixWithSpace())
	}
}

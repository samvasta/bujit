package models

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"samvasta.com/bujit/config"
	"samvasta.com/bujit/util"
)

type Money int64

func (m *Money) Cents() int64 {
	return int64(*m) % 100
}

func (m *Money) Dollars() int64 {
	return int64(*m) / 100
}

func (m *Money) IsNegative() bool {
	return int64(*m) < 0
}

func (m *Money) Value() int64 {
	return int64(*m)
}

func (m *Money) String() string {
	if m.IsNegative() {
		return fmt.Sprintf("(%s%s.%d)", config.CurrencySymbol(), humanize.Comma(util.AbsI64(m.Dollars())), m.Cents())
	} else {
		return fmt.Sprintf("%s%s.%d", config.CurrencySymbol(), humanize.Comma(m.Dollars()), m.Cents())
	}
}

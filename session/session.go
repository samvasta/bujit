package session

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Session struct {
	CurrencyPrefix string
	CurrencySuffix string
	Db             *gorm.DB
}

func (s *Session) CurrencySuffixWithSpace() string {
	if len(s.CurrencySuffix) > 0 {
		return " " + s.CurrencySuffix
	}
	return ""
}

func SQLiteSession(pathToDb string, initDb func(*gorm.DB)) Session {
	db, err := gorm.Open(sqlite.Open(pathToDb), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	initDb(db)

	return Session{
		CurrencyPrefix: "",
		CurrencySuffix: "USD",
		Db:             db}
}

func InMemorySession(initDb func(*gorm.DB)) Session {
	return SQLiteSession("file::memory:", initDb)
}

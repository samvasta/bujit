package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"samvasta.com/bujit/session"
)

type Category struct {
	ID              uint `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	Description     string
	SuperCategoryID *uint
	SubCategories   []Category       `gorm:"foreignkey:SuperCategoryID"`
	Accounts        []Account        `gorm:"foreignkey:ID"`
	Session         *session.Session `gorm:"-"` // Ignored by ORM
}

func (cat Category) CurrentBalance() Money {
	var total int64 = 0
	for _, account := range cat.Accounts {
		total += account.Balance().Value()
	}
	for _, subCat := range cat.SubCategories {
		total += subCat.CurrentBalance().Value()
	}
	return Money(total)
}

func (cat Category) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = cat.ID
	details["name"] = cat.Name
	details["description"] = cat.Description
	details["createdAt"] = cat.CreatedAt
	details["updatedAt"] = cat.UpdatedAt
	details["currentBalance"] = cat.CurrentBalance()

	subCategoryIds := []uint{}
	for _, subCat := range cat.SubCategories {
		subCategoryIds = append(subCategoryIds, subCat.ID)
	}
	details["subCategories"] = subCategoryIds

	accountIds := []uint{}
	for _, account := range cat.Accounts {
		accountIds = append(accountIds, account.ID)
	}
	details["accounts"] = accountIds

	return json.Marshal(details)
}

type AccountState struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	Balance     Money
	PrevStateID *uint
	PrevState   *AccountState
	Session     *session.Session `gorm:"-"` // Ignored by ORM
}

func (as AccountState) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = as.ID

	details["timestamp"] = as.CreatedAt
	details["balance"] = as.Balance.String(as.Session)

	return json.Marshal(details)
}

type Account struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	Name           string
	Description    string
	IsActive       bool
	CurrentStateID *uint
	CurrentState   AccountState     `gorm:"foreignkey:CurrentStateID"`
	Session        *session.Session `gorm:"-"` // Ignored by ORM
}

func (account *Account) Balance() Money {
	return account.CurrentState.Balance
}

func (account Account) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = account.ID
	details["name"] = account.Name

	if account.IsActive {
		details["status"] = "open"
	} else {
		details["status"] = "closed"
	}

	details["description"] = account.Description
	details["createdAt"] = account.CreatedAt
	details["updatedAt"] = account.CurrentState.CreatedAt
	details["currentBalance"] = account.Balance().String(account.Session)

	return json.Marshal(details)
}

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	Change        Money
	SourceID      *uint
	Source        *Account `gorm:"foreignkey:SourceID"`
	DestinationID *uint
	Destination   *Account `gorm:"foreignkey:DestinationID"`
	Memo          string
	Session       *session.Session `gorm:"-"` // Ignored by ORM
}

func (tran *Transaction) SourceExists() bool {
	return tran.Source != nil
}
func (tran *Transaction) DestinationExists() bool {
	return tran.Destination != nil
}

func (tran Transaction) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = tran.ID
	details["timestamp"] = tran.CreatedAt
	details["amount"] = tran.Change.String(tran.Session)

	if tran.SourceExists() {
		details["fromAccount"] = tran.Source.Name
	}
	if tran.DestinationExists() {
		details["toAccount"] = tran.Destination.Name
	}

	details["memo"] = tran.Memo

	return json.Marshal(details)
}

func MigrateSchema(db *gorm.DB) {
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&AccountState{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})
}

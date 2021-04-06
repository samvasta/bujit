package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"samvasta.com/bujit/session"
)

type Category struct {
	ID                 uint  `gorm:"primaryKey"`
	CreatedAt          int64 `gorm:"autoCreateTime"`
	UpdatedAt          int64 `gorm:"autoUpdateTime"`
	Name               string
	FullyQualifiedName string `gorm:"unique"` // the compound name of this category and the parents' names. Ex. "grandparent/parent/this"
	Description        string
	SuperCategoryID    *uint
	SubCategories      []Category       `gorm:"foreignkey:SuperCategoryID"`
	Accounts           []Account        `gorm:"foreignkey:ID"`
	Session            *session.Session `gorm:"-"` // Ignored by ORM
}

func (this Category) GetSession() *session.Session {
	return this.Session
}

func MakeCategory(name, description string, superCategory *Category, accounts ...Account) Category {
	category := Category{Name: name, Description: description, Accounts: accounts}
	category.SetParent(superCategory)
	return category
}

func (cat *Category) SetParent(parent *Category) {
	if parent != nil {
		cat.SuperCategoryID = &parent.ID
		cat.FullyQualifiedName = fmt.Sprintf("%s/%s", parent.FullyQualifiedName, cat.Name)
	} else {
		cat.SuperCategoryID = nil
		cat.FullyQualifiedName = cat.Name
	}

	for _, subCategory := range cat.SubCategories {
		subCategory.SetParent(cat)
	}
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
	details["createdAt"] = time.Unix(cat.CreatedAt, 0).UTC()
	details["updatedAt"] = time.Unix(cat.UpdatedAt, 0).UTC()
	details["currentBalance"] = cat.CurrentBalance().String(cat.Session)

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
	ID          uint  `gorm:"primaryKey"`
	CreatedAt   int64 `gorm:"autoCreateTime"`
	Balance     Money
	PrevStateID *uint
	PrevState   *AccountState
	IsClosed    bool
	Session     *session.Session `gorm:"-"` // Ignored by ORM
}

func (this AccountState) GetSession() *session.Session {
	return this.Session
}

func (as AccountState) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = as.ID

	details["timestamp"] = time.Unix(as.CreatedAt, 0).UTC()
	details["balance"] = as.Balance.String(as.Session)

	return json.Marshal(details)
}

type Account struct {
	ID             uint   `gorm:"primaryKey"`
	CreatedAt      int64  `gorm:"autoCreateTime"`
	Name           string `gorm:"unique"`
	Description    string
	IsActive       bool
	CurrentStateID *uint
	CurrentState   AccountState `gorm:"foreignkey:CurrentStateID"`
	CategoryID     *uint
	Category       Category
	Session        *session.Session `gorm:"-"` // Ignored by ORM
}

func (this Account) GetSession() *session.Session {
	return this.Session
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
	details["createdAt"] = time.Unix(account.CreatedAt, 0).UTC()
	details["updatedAt"] = time.Unix(account.CurrentState.CreatedAt, 0).UTC()
	details["currentBalance"] = account.Balance().String(account.Session)

	return json.Marshal(details)
}

type Transaction struct {
	ID            uint  `gorm:"primaryKey"`
	CreatedAt     int64 `gorm:"autoCreateTime"`
	Change        Money
	SourceID      *uint
	Source        *Account `gorm:"foreignkey:SourceID"`
	DestinationID *uint
	Destination   *Account `gorm:"foreignkey:DestinationID"`
	Memo          string
	Session       *session.Session `gorm:"-"` // Ignored by ORM
}

func (this Transaction) GetSession() *session.Session {
	return this.Session
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
	details["timestamp"] = time.Unix(tran.CreatedAt, 0).UTC()
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

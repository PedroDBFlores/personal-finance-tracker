package database

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"gorm.io/gorm"
)

type DbAccountType domain.AccountType
type DbCurrency domain.Currency
type DbMoney domain.Money
type DbMovementType domain.MovementType

type Account struct {
	gorm.Model
	Name      string          `gorm:"not null;type:varchar(128)"`
	Type      DbAccountType   `gorm:"not null;type:varchar(64)"`
	Currency  domain.Currency `gorm:"not null;type:varchar(3)"`
	Movements []Movement      `gorm:"joinForeignKey:account_id"`
}

type Movement struct {
	gorm.Model
	AccountId    uint   `gorm:"not null"`
	Description  string `gorm:"not null;type:varchar(255)"`
	Value        domain.Money
	MovementType domain.MovementType
	Date         time.Time
}

type Investment struct {
	gorm.Model
	AccountId    uint                  `gorm:"not null;index"`
	Name         string                `gorm:"not null;type:varchar(128)"`
	Type         domain.InvestmentType `gorm:"not null;type:varchar(64)"`
	CurrentValue domain.Money          `gorm:"not null"`
	Currency     domain.Currency       `gorm:"not null;type:varchar(3)"`
	PurchaseDate time.Time             `gorm:"not null"`
	MaturityDate *time.Time
}

type InvestmentTransaction struct {
	gorm.Model
	InvestmentId uint                            `gorm:"not null;index"`
	Type         domain.InvestmentTransactionType `gorm:"not null;type:varchar(64)"`
	Amount       domain.Money                     `gorm:"not null"`
	Quantity     domain.Money                     `gorm:"not null"`
	Price        domain.Money                     `gorm:"not null"`
	Date         time.Time                        `gorm:"not null"`
	Description  string                           `gorm:"type:varchar(255)"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *DbAccountType) Scan(value interface{}) error {
	bytes, ok := value.(string)
	if !ok {
		return errors.New("oops")
	}
	*j = DbAccountType(bytes)
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j DbAccountType) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

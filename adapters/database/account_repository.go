package database

import (
	"encoding/json"
	"fmt"
	"time"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/models"
)

type Account struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

const prefix = "account"

func accountKey(id uuid.UUID) []byte {
	return fmt.Appendf(nil, "%s:%s", prefix, id)
}

type AccountRepository struct {
	db *badger.DB
}

func NewAccountRepository(db *badger.DB) AccountRepository {
	return AccountRepository{db}
}

func (a *AccountRepository) Create(account models.Account) (uuid.UUID, error) {
	id := uuid.NewV7()
	key := accountKey(id)
	account.Id = id
	b, err := json.Marshal(account)
	if err != nil {
		return uuid.Nil(), err
	}
	if err := a.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	}); err != nil {
		return uuid.Nil(), err
	}

	return id, nil
}

func (a *AccountRepository) Update(account models.Account) error {
	key := accountKey(account.Id)
	ok, err := IsKeyPresent(a.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Account", fmt.Sprintf("account:%s", account.Id))
	}
	return a.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err != nil {
			return err
		}
		b, err := json.Marshal(account)
		if err != nil {
			return err
		}
		return txn.Set(key, b)
	})
}

func (a *AccountRepository) Get(id uuid.UUID) (*models.Account, error) {
	key := accountKey(id)
	out, err := GetByKey[Account](a.db, key)
	if err != nil {
		return nil, err
	}
	modelAccount := toModelAccount(out)
	return &modelAccount, nil
}

func (a *AccountRepository) List() ([]models.Account, error) {
	accounts, err := ListAllByKey[Account](a.db, []byte(prefix))
	if err != nil {
		return []models.Account{}, err
	}
	var modelAccounts []models.Account
	for _, acc := range accounts {
		modelAccounts = append(modelAccounts, toModelAccount(acc))
	}

	return modelAccounts, nil
}

func (a *AccountRepository) Delete(id uuid.UUID) error {
	return a.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(accountKey(id))
	})
}

func toModelAccount(account Account) models.Account {
	return models.Account{
		Id:        account.Id,
		Name:      account.Name,
		Type:      models.AccountType(account.Type),
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

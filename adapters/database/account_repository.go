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

type account struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

const accountPrefix = "account"

func accountKey(id uuid.UUID) []byte {
	return fmt.Appendf(nil, "%s:%s", accountPrefix, id)
}

type AccountRepository struct {
	db *badger.DB
}

func NewAccountRepository(db *badger.DB) AccountRepository {
	return AccountRepository{db}
}

func (r *AccountRepository) Create(account models.Account) (uuid.UUID, error) {
	id := uuid.NewV7()
	key := accountKey(id)
	account.Id = id
	b, err := json.Marshal(toRepoAccount(account))
	if err != nil {
		return uuid.Nil(), err
	}
	if err := r.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	}); err != nil {
		return uuid.Nil(), err
	}

	return id, nil
}

func (r *AccountRepository) Update(account models.Account) error {
	key := accountKey(account.Id)
	ok, err := IsKeyPresent(r.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Account", string(key))
	}
	return r.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err != nil {
			return err
		}
		b, err := json.Marshal(toRepoAccount(account))
		if err != nil {
			return err
		}
		return txn.Set(key, b)
	})
}

func (r *AccountRepository) Get(id uuid.UUID) (*models.Account, error) {
	key := accountKey(id)
	out, err := GetByKey[account](r.db, key)
	if err != nil {
		return nil, err
	}
	modelAccount := toModelAccount(out)
	return &modelAccount, nil
}

func (r *AccountRepository) List() ([]models.Account, error) {
	accounts, err := ListAllByKey[account](r.db, []byte(accountPrefix))
	if err != nil {
		return []models.Account{}, err
	}
	var modelAccounts []models.Account
	for _, acc := range accounts {
		modelAccounts = append(modelAccounts, toModelAccount(acc))
	}

	return modelAccounts, nil
}

func (r *AccountRepository) Delete(id uuid.UUID) error {
	key := accountKey(id)
	ok, err := IsKeyPresent(r.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Account", string(key))
	}
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func toModelAccount(account account) models.Account {
	return models.Account{
		Id:        account.Id,
		Name:      account.Name,
		Type:      models.AccountType(account.Type),
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func toRepoAccount(modelAccount models.Account) account {
	return account{
		Id:        modelAccount.Id,
		Name:      modelAccount.Name,
		Type:      string(modelAccount.Type),
		Currency:  modelAccount.Currency,
		CreatedAt: modelAccount.CreatedAt,
	}
}

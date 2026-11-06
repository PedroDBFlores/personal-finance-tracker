package database

import (
	"encoding/json"
	"fmt"
	"time"
	"uuid"

	"github.com/pedrodbflores/personal-finance-tracker/domain"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
)

type account struct {
	Id        domain.Id          `json:"id"`
	Name      string             `json:"name"`
	Type      domain.AccountType `json:"type"`
	Currency  domain.Currency    `json:"currency"`
	CreatedAt time.Time          `json:"createdAt"`
}

const accountPrefix = "account"

func accountKey(id domain.Id) []byte {
	return fmt.Appendf(nil, "%s:%s", accountPrefix, id)
}

type AccountRepository struct {
	db *badger.DB
}

func NewAccountRepository(db *badger.DB) AccountRepository {
	return AccountRepository{db}
}

func (r *AccountRepository) Create(account domain.Account) (domain.Id, error) {
	id := domain.Id(uuid.NewV7())
	key := accountKey(id)
	account.Id = id
	b, err := json.Marshal(domainToRepo(account))
	if err != nil {
		return domain.Id(uuid.Nil()), err
	}
	if err := r.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	}); err != nil {
		return domain.Id(uuid.Nil()), err
	}

	return id, nil
}

func (r *AccountRepository) Update(account domain.Account) error {
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
		b, err := json.Marshal(domainToRepo(account))
		if err != nil {
			return err
		}
		return txn.Set(key, b)
	})
}

func (r *AccountRepository) Get(id domain.Id) (*domain.Account, error) {
	key := accountKey(id)
	out, err := GetByKey[account](r.db, key)
	if err != nil {
		return nil, err
	}
	modelAccount := accountToDomain(out)
	return &modelAccount, nil
}

func (r *AccountRepository) List() ([]domain.Account, error) {
	accounts, err := ListAllByKey[account](r.db, []byte(accountPrefix))
	if err != nil {
		return []domain.Account{}, err
	}
	var modelAccounts []domain.Account
	for _, acc := range accounts {
		modelAccounts = append(modelAccounts, accountToDomain(acc))
	}

	return modelAccounts, nil
}

func (r *AccountRepository) Delete(id domain.Id) error {
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

func accountToDomain(account account) domain.Account {
	return domain.Account{
		Id:        account.Id,
		Name:      account.Name,
		Type:      account.Type,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func domainToRepo(modelAccount domain.Account) account {
	return account{
		Id:        modelAccount.Id,
		Name:      modelAccount.Name,
		Type:      modelAccount.Type,
		Currency:  modelAccount.Currency,
		CreatedAt: modelAccount.CreatedAt,
	}
}

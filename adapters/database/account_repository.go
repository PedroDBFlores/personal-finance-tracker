package database

import (
	"encoding/json"
	"fmt"
	"time"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/models"
)

type Account struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

func accountKey(id uuid.UUID) []byte {
	return fmt.Appendf(nil, "account#%s", id)
}

type AccountRepository struct {
	db *badger.DB
}

func NewAccountRepository(db *badger.DB) AccountRepository {
	return AccountRepository{db}
}

func (a *AccountRepository) Create(account models.Account) (uuid.UUID, error) {
	id := uuid.NewV7()
	key := fmt.Sprintf("account#%s", id)
	account.Id = id
	b, _ := json.Marshal(account)
	_ = a.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), b)
	})

	var out Account
	_ = a.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &out)
		})
	})

	return id, nil
}

func toRepoAccount(modelAccount models.Account) Account {
	return Account{
		Id:        modelAccount.Id,
		Name:      modelAccount.Name,
		Type:      string(modelAccount.Type),
		Currency:  modelAccount.Currency,
		CreatedAt: modelAccount.CreatedAt,
	}
}

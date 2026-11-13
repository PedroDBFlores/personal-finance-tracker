package database

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"gorm.io/gorm"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return AccountRepository{db}
}

func (r *AccountRepository) Create(ctx context.Context, account domain.Account) (*domain.Id, error) {
	acc := domainAccountToDBAccount(account)
	if err := gorm.G[Account](r.db).Create(ctx, &acc); err != nil {
		return nil, err
	}
	id := domain.Id(acc.ID)
	return &id, nil
}

func (r *AccountRepository) Update(ctx context.Context, account domain.Account) error {
	if rowsAffected, err := gorm.G[Account](r.db).Updates(ctx, domainAccountToDBAccount(account)); rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Account", account.Id.String())
	} else if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Get(ctx context.Context, id domain.Id) (*domain.Account, error) {
	if acc, err := gorm.G[Account](r.db).Where("id = ?", id).First(ctx); err != nil {
		return nil, app_errors.NewResourceNotFoundError("Account", id.String())
	} else {
		account := dbAccountToDomainAccount(acc)
		return &account, nil
	}
}

func (r *AccountRepository) List(ctx context.Context) ([]domain.Account, error) {
	accounts, err := gorm.G[Account](r.db).Find(ctx)
	if err != nil {
		return []domain.Account{}, err
	}
	var modelAccounts []domain.Account
	for _, acc := range accounts {
		modelAccounts = append(modelAccounts, dbAccountToDomainAccount(acc))
	}

	return modelAccounts, nil
}

func (r *AccountRepository) Delete(ctx context.Context, id domain.Id) error {
	if rowsAffected, err := gorm.G[Account](r.db).Where("id = ?", id).Delete(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Account", id.String())
	}
	return nil
}

func dbAccountToDomainAccount(account Account) domain.Account {
	return domain.Account{
		Id:        domain.Id(account.ID),
		Name:      account.Name,
		Type:      domain.AccountType(account.Type),
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func domainAccountToDBAccount(modelAccount domain.Account) Account {
	return Account{
		ID:        uint(modelAccount.Id),
		Name:      modelAccount.Name,
		Type:      DbAccountType(modelAccount.Type),
		Currency:  modelAccount.Currency,
		CreatedAt: modelAccount.CreatedAt,
	}
}

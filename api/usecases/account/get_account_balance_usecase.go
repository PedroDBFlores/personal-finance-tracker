package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type GetAccountBalanceUseCase struct {
	accountGetter    usecases.AccountGetter
	movementLister   usecases.MovementLister
	investmentLister usecases.InvestmentLister
	txLister         usecases.InvestmentTransactionLister
	logger           *slog.Logger
}

func NewGetAccountBalanceUseCase(accountGetter usecases.AccountGetter, movementLister usecases.MovementLister, investmentLister usecases.InvestmentLister, txLister usecases.InvestmentTransactionLister, logger *slog.Logger) GetAccountBalanceUseCase {
	return GetAccountBalanceUseCase{accountGetter, movementLister, investmentLister, txLister, logger}
}

func (u *GetAccountBalanceUseCase) Execute(ctx context.Context, accountId domain.Id) (*dtos.AccountBalanceDto, error) {
	acc, err := u.accountGetter.Get(ctx, accountId)
	if err != nil {
		return nil, err
	}

	u.logger.DebugContext(ctx, "listing movements for balance", slog.Any("accountId", accountId))
	movements, err := u.movementLister.List(ctx, accountId)
	if err != nil {
		return nil, err
	}
	movementsBalance := domain.Money(0)
	for _, m := range movements {
		if m.Type == domain.MovementTypeCredit {
			movementsBalance = movementsBalance.Add(m.Value)
		} else if m.Type == domain.MovementTypeDebit {
			movementsBalance = movementsBalance.Subtract(m.Value)
		}
	}

	u.logger.DebugContext(ctx, "listing investments for balance", slog.Any("accountId", accountId))
	investments, err := u.investmentLister.List(ctx, accountId)
	if err != nil {
		return nil, err
	}
	investmentsBalance := domain.Money(0)
	for i := range investments {
		inv := &investments[i]
		transactions, err := u.txLister.List(ctx, accountId, inv.Id)
		if err != nil {
			return nil, err
		}
		inv.UpdateCurrentValue(transactions)
		investmentsBalance = investmentsBalance.Add(inv.CurrentValue)
	}

	dto := dtos.NewAccountBalanceDto(acc.Id, acc.Currency, movementsBalance, investmentsBalance)
	return &dto, nil
}

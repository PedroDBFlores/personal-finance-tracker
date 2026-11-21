package dtos

import (
	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

// AccountBalanceDto is the current balance of an account, split into the
// movements part and the investments part, plus the total.
type AccountBalanceDto struct {
	// AccountId identifies the account the balance belongs to.
	AccountId domain.Id
	// Currency is the account's currency.
	Currency domain.Currency
	// MovementsBalance is the sum of credits minus debits (minor units).
	MovementsBalance domain.Money
	// InvestmentsBalance is the sum of the current values of the account's investments (minor units).
	InvestmentsBalance domain.Money
	// TotalBalance is MovementsBalance + InvestmentsBalance (minor units).
	TotalBalance domain.Money
}

func NewAccountBalanceDto(accountId domain.Id, currency domain.Currency, movementsBalance, investmentsBalance domain.Money) AccountBalanceDto {
	return AccountBalanceDto{
		AccountId:          accountId,
		Currency:           currency,
		MovementsBalance:   movementsBalance,
		InvestmentsBalance: investmentsBalance,
		TotalBalance:       movementsBalance.Add(investmentsBalance),
	}
}

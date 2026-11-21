package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type ListMovementMonthsUseCase struct {
	accountGetter usecases.AccountGetter
	monthLister   usecases.MovementMonthLister
	logger        *slog.Logger
}

func NewListMovementMonthsUseCase(accountGetter usecases.AccountGetter, monthLister usecases.MovementMonthLister, logger *slog.Logger) ListMovementMonthsUseCase {
	return ListMovementMonthsUseCase{accountGetter, monthLister, logger}
}

func (u *ListMovementMonthsUseCase) Execute(ctx context.Context, accountId domain.Id) ([]dtos.MovementMonthDto, error) {
	if _, err := u.accountGetter.Get(ctx, accountId); err != nil {
		return nil, err
	}
	u.logger.DebugContext(ctx, "listing movement months", slog.Any("accountId", accountId))
	months, err := u.monthLister.ListMonths(ctx, accountId)
	if err != nil {
		return nil, err
	}
	result := make([]dtos.MovementMonthDto, 0, len(months))
	for _, ym := range months {
		result = append(result, dtos.NewMovementMonthDto(ym))
	}
	return result, nil
}

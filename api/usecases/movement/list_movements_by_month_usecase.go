package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type ListMovementsByMonthUseCase struct {
	accountGetter usecases.AccountGetter
	byMonthLister usecases.MovementByMonthLister
	logger        *slog.Logger
}

func NewListMovementsByMonthUseCase(accountGetter usecases.AccountGetter, byMonthLister usecases.MovementByMonthLister, logger *slog.Logger) ListMovementsByMonthUseCase {
	return ListMovementsByMonthUseCase{accountGetter, byMonthLister, logger}
}

func (u *ListMovementsByMonthUseCase) Execute(ctx context.Context, accountId domain.Id, year, month int) ([]dtos.MovementDto, error) {
	if _, err := u.accountGetter.Get(ctx, accountId); err != nil {
		return nil, err
	}
	u.logger.DebugContext(ctx, "listing movements by month", slog.Any("accountId", accountId), slog.Int("year", year), slog.Int("month", month))
	movements, err := u.byMonthLister.ListByMonth(ctx, accountId, year, month)
	if err != nil {
		return nil, err
	}
	result := make([]dtos.MovementDto, 0, len(movements))
	for _, m := range movements {
		result = append(result, dtos.NewMovementDto(m))
	}
	return result, nil
}

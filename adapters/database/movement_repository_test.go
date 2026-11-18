package database_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/stretchr/testify/require"
)

func TestMovementRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	baseAccount := domain.Account{
		Name:      "Test Account",
		Type:      domain.AccountTypeBank,
		Currency:  domain.ParseCurrency("USD"),
		CreatedAt: time.Now(),
	}

	t.Run("it creates a movement in the database", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewMovementRepository(db)
		movement := domain.NewMovement(
			*accountId,
			"Test Movement",
			domain.Money(1000),
			domain.MovementTypeDebit,
			time.Now(),
		)
		id, err := repo.Create(ctx, *accountId, movement)
		require.NoError(t, err)
		require.Greater(t, *id, uint(0))

		mov, err := repo.Get(ctx, *accountId, *id)
		require.NoError(t, err)
		require.NotNil(t, mov)
		require.Equal(t, *accountId, mov.AccountId)
		require.Equal(t, domain.MovementTypeDebit, mov.Type)
		require.Equal(t, domain.Money(1000), mov.Value)
		require.WithinDuration(t, movement.Date, mov.Date, 0)
	})

	t.Run("get", func(t *testing.T) {
		t.Run("gets a movement by id", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			accRepo := database.NewAccountRepository(db)
			accountId, err := accRepo.Create(ctx, baseAccount)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)
			movement := domain.NewMovement(
				*accountId,
				"Test Movement",
				domain.Money(1000),
				domain.MovementTypeDebit,
				time.Now(),
			)
			id, err := repo.Create(ctx, *accountId, movement)
			require.NoError(t, err)

			mov, err := repo.Get(ctx, *accountId, *id)

			require.NoError(t, err)
			require.NotNil(t, mov)
			require.Equal(t, movement.Description, mov.Description)
			require.Equal(t, movement.Value, mov.Value)
			require.Equal(t, movement.Type, mov.Type)
			require.WithinDuration(t, movement.Date, mov.Date, 0)
		})

		t.Run("returns ResourceNotFoundError if the movement doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)

			id := domain.NewId()
			mov, err := repo.Get(ctx, domain.NewId(), id)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", id.String()))
			require.Nil(t, mov)
		})

		t.Run("returns ResourceNotFoundError if the movement belongs to another account", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			accRepo := database.NewAccountRepository(db)
			accountId, err := accRepo.Create(ctx, baseAccount)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)
			movement := domain.NewMovement(*accountId, "Test Movement", domain.Money(1000), domain.MovementTypeDebit, time.Now())
			id, err := repo.Create(ctx, *accountId, movement)
			require.NoError(t, err)

			mov, err := repo.Get(ctx, domain.NewId(), *id)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", id.String()))
			require.Nil(t, mov)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("updates the given fields of a movement", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			accRepo := database.NewAccountRepository(db)
			accountId, err := accRepo.Create(ctx, baseAccount)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)
			movement := domain.NewMovement(
				*accountId,
				"Test Movement",
				domain.Money(1000),
				domain.MovementTypeDebit,
				time.Now(),
			)
			id, err := repo.Create(ctx, *accountId, movement)
			require.NoError(t, err)

			newDescription := "Updated Movement"
			newValue := domain.Money(2000)
			err = repo.Update(ctx, *accountId, *id, domain.MovementPatch{
				Description: &newDescription,
				Value:       &newValue,
			})
			require.NoError(t, err)

			updatedMov, err := repo.Get(ctx, *accountId, *id)
			require.NoError(t, err)
			require.NotNil(t, updatedMov)
			require.Equal(t, "Updated Movement", updatedMov.Description)
			require.Equal(t, domain.Money(2000), updatedMov.Value)
			// untouched fields keep their values
			require.Equal(t, domain.MovementTypeDebit, updatedMov.Type)
			require.WithinDuration(t, movement.Date, updatedMov.Date, 0)
		})

		t.Run("returns ResourceNotFoundError if updating a non-existing movement", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)

			accountId := domain.NewId()
			movementId := domain.NewId()
			newValue := domain.Money(2000)

			err = repo.Update(ctx, accountId, movementId, domain.MovementPatch{Value: &newValue})
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", movementId.String()))
		})
	})

	t.Run("list the movements for a given account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewMovementRepository(db)
		firstMovement := domain.NewMovement(
			*accountId,
			"Test Movement",
			domain.Money(1000),
			domain.MovementTypeDebit,
			time.Now(),
		)
		firstId, err := repo.Create(ctx, *accountId, firstMovement)
		require.NoError(t, err)
		secondMovement := domain.NewMovement(
			*accountId,
			"Test Movement 2",
			domain.Money(10000),
			domain.MovementTypeCredit,
			time.Now(),
		)
		secondId, err := repo.Create(ctx, *accountId, secondMovement)
		require.NoError(t, err)

		accountMovements, err := repo.List(ctx, *accountId)
		require.NoError(t, err)
		require.Len(t, accountMovements, 2)
		require.Equal(t, *firstId, accountMovements[0].Id)
		require.Equal(t, *secondId, accountMovements[1].Id)

		otherMovements, err := repo.List(ctx, domain.NewId())
		require.NoError(t, err)
		require.Empty(t, otherMovements)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("deletes a movement", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			accRepo := database.NewAccountRepository(db)
			accountId, err := accRepo.Create(ctx, baseAccount)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)
			movement := domain.NewMovement(
				*accountId,
				"Test Movement",
				domain.Money(1000),
				domain.MovementTypeDebit,
				time.Now(),
			)
			id, err := repo.Create(ctx, *accountId, movement)
			require.NoError(t, err)

			err = repo.Delete(ctx, *accountId, *id)
			require.NoError(t, err)

			mov, err := repo.Get(ctx, *accountId, *id)
			require.Error(t, err)
			require.Nil(t, mov)
		})

		t.Run("returns ResourceNotFoundError if the movement doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)

			id := domain.NewId()
			err = repo.Delete(ctx, domain.NewId(), id)

			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", id.String()))
		})

		t.Run("returns ResourceNotFoundError if the movement belongs to another account", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			accRepo := database.NewAccountRepository(db)
			accountId, err := accRepo.Create(ctx, baseAccount)
			require.NoError(t, err)
			repo := database.NewMovementRepository(db)
			movement := domain.NewMovement(*accountId, "Test Movement", domain.Money(1000), domain.MovementTypeDebit, time.Now())
			id, err := repo.Create(ctx, *accountId, movement)
			require.NoError(t, err)

			err = repo.Delete(ctx, domain.NewId(), *id)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", id.String()))

			// the movement still exists for the right account
			mov, err := repo.Get(ctx, *accountId, *id)
			require.NoError(t, err)
			require.NotNil(t, mov)
		})
	})
}

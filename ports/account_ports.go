package ports

import (
	"uuid"

	"github.com/pedrodbflores/personal-finance-tracker/models"
)

//go:generate go tool counterfeiter -generate

//counterfeiter:generate . AccountCreator
type AccountCreator interface {
	Create(account models.Account) (uuid.UUID, error)
}

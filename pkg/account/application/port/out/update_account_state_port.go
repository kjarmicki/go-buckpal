package account_application_port_out

import (
	"context"

	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
)

type UpdateAccountStatePort interface {
	UpdateActivities(ctx context.Context, account *account_domain.Account) error
}

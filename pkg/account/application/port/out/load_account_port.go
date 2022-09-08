package account_application_port_out

import (
	"context"
	"time"

	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
)

type LoadAccountPort interface {
	LoadAccount(ctx context.Context, accountId account_domain.AccountId, baselineDate time.Time) (*account_domain.Account, error)
}

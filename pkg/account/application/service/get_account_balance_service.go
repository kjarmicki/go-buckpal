package account_application_service

import (
	"context"
	"time"

	account_application_port_out "github.com/kjarmicki/go-buckpal/pkg/account/application/port/out"
	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
)

type GetAccountBalanceService struct {
	loadAccountPort account_application_port_out.LoadAccountPort
}

func NewGetAccountBalanceService(loadAccountPort account_application_port_out.LoadAccountPort) *GetAccountBalanceService {
	return &GetAccountBalanceService{
		loadAccountPort: loadAccountPort,
	}
}

func (gabs *GetAccountBalanceService) GetAccountBalance(ctx context.Context, accountId account_domain.AccountId) (account_domain.Money, error) {
	account, err := gabs.loadAccountPort.LoadAccount(ctx, accountId, time.Now())
	if err != nil {
		return account_domain.Money{}, err
	}
	return account.CalculateBalance(), nil
}

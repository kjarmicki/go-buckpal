package account_adapter_out_persistence

import (
	"context"
	"time"

	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
)

type AccountPersistenceAdapter struct {
	accountRepository  account_domain.AccountRepository
	activityRepository account_domain.ActivityWindowRepository
}

func NewAccountPersistenceAdapter(accountRepository account_domain.AccountRepository, activityRepository account_domain.ActivityWindowRepository) *AccountPersistenceAdapter {
	return &AccountPersistenceAdapter{
		accountRepository:  accountRepository,
		activityRepository: activityRepository,
	}
}

func (apa *AccountPersistenceAdapter) LoadAccount(ctx context.Context, accountId account_domain.AccountId, baselineDate time.Time) (*account_domain.Account, error) {
	return apa.accountRepository.FindById(ctx, accountId, baselineDate)
}

func (apa *AccountPersistenceAdapter) UpdateActivities(ctx context.Context, account *account_domain.Account) error {
	for _, activity := range account.GetActivities() {
		if activity.Id == 0 {
			err := apa.activityRepository.SaveActivity(ctx, activity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

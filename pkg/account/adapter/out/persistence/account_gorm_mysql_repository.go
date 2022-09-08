package account_adapter_out_persistence

import (
	"context"
	"time"

	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
	"gorm.io/gorm"
)

type AccountGormMySqlRepository struct {
	table                    *gorm.DB
	activityWindowRepository account_domain.ActivityWindowRepository
}

type AccountGormMySqlRepositoryEntity struct {
	Id int
}

func NewAccountGormMySqlRepository(db *gorm.DB, activityWindowRepository account_domain.ActivityWindowRepository) *AccountGormMySqlRepository {
	return &AccountGormMySqlRepository{
		table:                    db.Table("accounts"),
		activityWindowRepository: activityWindowRepository,
	}
}

func (agr *AccountGormMySqlRepository) FindById(ctx context.Context, accountId account_domain.AccountId, baselineDate time.Time) (*account_domain.Account, error) {
	var account AccountGormMySqlRepositoryEntity
	if err := agr.table.First(&account, accountId).Error; err != nil {
		return nil, err
	}

	activities, err := agr.activityWindowRepository.FindByOwnerSince(ctx, accountId, baselineDate)
	if err != nil {
		return nil, err
	}

	withdrawalBalance, err := agr.activityWindowRepository.GetWithdrawalBalanceUntil(ctx, accountId, baselineDate)
	if err != nil {
		return nil, err
	}

	depositBalance, err := agr.activityWindowRepository.GetDepositBalanceUntil(ctx, accountId, baselineDate)
	if err != nil {
		return nil, err
	}

	return account_domain.NewAccount(accountId, depositBalance.Minus(withdrawalBalance), activities), nil
}

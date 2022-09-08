package account_adapter_out_persistence

import (
	"context"
	"time"

	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
	"gorm.io/gorm"
)

type ActivityWindowGormMySqlRepository struct {
	db *gorm.DB
}

type ActivityWindowGormMySqlRepositoryData struct {
	Id              int
	Timestamp       time.Time
	OwnerAccountId  int
	SourceAccountId int
	TargetAccountId int
	Amount          int64
}

func NewActivityWindowGormRepository(db *gorm.DB) *ActivityWindowGormMySqlRepository {
	return &ActivityWindowGormMySqlRepository{
		db: db,
	}
}

func (agr *ActivityWindowGormMySqlRepository) SaveActivity(ctx context.Context, activity account_domain.Activity) error {
	activityData := MapActivityToGormMySqlRepositoryData(activity)
	if err := agr.db.
		Table("activities").
		Create(activityData).
		Error; err != nil {
		return err
	}
	return nil
}

func (agr *ActivityWindowGormMySqlRepository) FindByOwnerSince(ctx context.Context, ownerAccountId account_domain.AccountId, since time.Time) (*account_domain.ActivityWindow, error) {
	var entities []ActivityWindowGormMySqlRepositoryData
	if err := agr.db.
		Table("activities").
		Where("owner_account_id = ? AND timestamp >= ?", ownerAccountId, since).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}
	return MapGormMySqlRepositoryDataToActivityWindow(entities), nil
}

func (agr *ActivityWindowGormMySqlRepository) GetDepositBalanceUntil(ctx context.Context, accountId account_domain.AccountId, until time.Time) (account_domain.Money, error) {
	var total int64
	if err := agr.db.
		Table("activities").
		Select("COALESCE(SUM(amount), 0)").
		Where("target_account_id = ? AND owner_account_id = ? AND timestamp < ?", accountId, accountId, until).
		Row().
		Scan(&total); err != nil {
		return account_domain.NewMoney(0), err
	}
	return account_domain.NewMoney(total), nil
}

func (agr *ActivityWindowGormMySqlRepository) GetWithdrawalBalanceUntil(ctx context.Context, accountId account_domain.AccountId, until time.Time) (account_domain.Money, error) {
	var total int64
	if err := agr.db.
		Table("activities").
		Select("COALESCE(SUM(amount), 0)").
		Where("source_account_id = ? AND owner_account_id = ? AND timestamp < ?", accountId, accountId, until).
		Row().
		Scan(&total); err != nil {
		return account_domain.NewMoney(0), err
	}
	return account_domain.NewMoney(total), nil
}

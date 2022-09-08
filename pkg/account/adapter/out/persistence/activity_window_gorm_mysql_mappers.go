package account_adapter_out_persistence

import account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"

func MapActivityToGormMySqlRepositoryData(activity account_domain.Activity) ActivityWindowGormMySqlRepositoryData {
	amount := activity.Money.GetAmount()
	return ActivityWindowGormMySqlRepositoryData{
		Id:              int(activity.Id),
		Timestamp:       activity.Timestamp,
		OwnerAccountId:  int(activity.OwnerAccountId),
		SourceAccountId: int(activity.SourceAccountId),
		TargetAccountId: int(activity.TargetAccountId),
		Amount:          amount,
	}
}

func MapGormMySqlRepositoryDataToActivity(data ActivityWindowGormMySqlRepositoryData) account_domain.Activity {
	return account_domain.Activity{
		Id:              account_domain.ActivityId(data.Id),
		Timestamp:       data.Timestamp,
		OwnerAccountId:  account_domain.AccountId(data.OwnerAccountId),
		SourceAccountId: account_domain.AccountId(data.SourceAccountId),
		TargetAccountId: account_domain.AccountId(data.TargetAccountId),
		Money:           account_domain.NewMoney(data.Amount),
	}
}

func MapGormMySqlRepositoryDataToActivityWindow(data []ActivityWindowGormMySqlRepositoryData) *account_domain.ActivityWindow {
	activities := make([]account_domain.Activity, 0, len(data))
	for _, one := range data {
		activities = append(activities, MapGormMySqlRepositoryDataToActivity(one))
	}
	return account_domain.NewActivityWindow(activities)
}

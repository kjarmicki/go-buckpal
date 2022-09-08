package account_domain

import (
	"context"
	"time"
)

type AccountId int

type Account struct {
	id              AccountId
	baselineBalance Money
	activityWindow  *ActivityWindow
}

type AccountRepository interface {
	FindById(ctx context.Context, id AccountId, baselineDate time.Time) (*Account, error)
}

func NewAccount(id AccountId, baselineBalance Money, activityWindow *ActivityWindow) *Account {
	return &Account{
		id:              id,
		baselineBalance: baselineBalance,
		activityWindow:  activityWindow,
	}
}

func (a *Account) GetId() AccountId {
	return a.id
}

func (a *Account) GetActivities() []Activity {
	return a.activityWindow.activities
}

func (a *Account) CalculateBalance() Money {
	return MoneyAdd(
		a.baselineBalance,
		a.activityWindow.calculateBalance(a.id),
	)
}

func (a *Account) Withdraw(money Money, targetAccountId AccountId) bool {
	if !a.mayWithdraw(money) {
		return false
	}
	withdrawal := Activity{
		OwnerAccountId:  a.id,
		SourceAccountId: a.id,
		TargetAccountId: targetAccountId,
		Timestamp:       time.Now(),
		Money:           money,
	}
	a.activityWindow.AddActivity(withdrawal)
	return true
}

func (a *Account) Deposit(money Money, sourceAccountId AccountId) bool {
	deposit := Activity{
		OwnerAccountId:  a.id,
		SourceAccountId: sourceAccountId,
		TargetAccountId: a.id,
		Timestamp:       time.Now(),
		Money:           money,
	}
	a.activityWindow.AddActivity(deposit)
	return true
}

func (a *Account) mayWithdraw(money Money) bool {
	return a.CalculateBalance().Minus(money).IsPositiveOrZero()
}

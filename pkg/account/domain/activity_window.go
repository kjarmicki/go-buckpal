package account_domain

import (
	"context"
	"time"
)

type ActivityWindow struct {
	activities []Activity
}

type ActivityWindowRepository interface {
	SaveActivity(ctx context.Context, activity Activity) error
	FindByOwnerSince(ctx context.Context, ownerAccountId AccountId, since time.Time) (*ActivityWindow, error)
	GetDepositBalanceUntil(ctx context.Context, accountId AccountId, until time.Time) (Money, error)
	GetWithdrawalBalanceUntil(ctx context.Context, accountId AccountId, until time.Time) (Money, error)
}

func NewActivityWindow(activities []Activity) *ActivityWindow {
	return &ActivityWindow{
		activities: activities,
	}
}

func (aw *ActivityWindow) GetActivities() []Activity {
	return aw.activities
}

func (aw *ActivityWindow) AddActivity(activity Activity) {
	aw.activities = append(aw.activities, activity)
}

// the timestamp of the first activity within this window
func (aw *ActivityWindow) getStartTimestamp() (time.Time, error) {
	if len(aw.activities) == 0 {
		return time.Time{}, NewIllegalStateError("no activities within a window")
	}
	min := aw.activities[0]
	for _, activity := range aw.activities {
		if activity.Timestamp.Before(min.Timestamp) {
			min = activity
		}
	}
	return min.Timestamp, nil
}

// the timestamp of the last activity within this window
func (aw *ActivityWindow) getEndTimestamp() (time.Time, error) {
	if len(aw.activities) == 0 {
		return time.Time{}, NewIllegalStateError("no activities within a window")
	}
	max := aw.activities[0]
	for _, activity := range aw.activities {
		if activity.Timestamp.After(max.Timestamp) {
			max = activity
		}
	}
	return max.Timestamp, nil
}

// calculate the balance by summing up the values of all activities within this window
func (aw *ActivityWindow) calculateBalance(accountId AccountId) Money {
	depositBalance := NewMoney(0)
	for _, activity := range aw.activities {
		if activity.TargetAccountId == accountId {
			depositBalance = depositBalance.Plus(activity.Money)
		}
	}

	withdrawalBalance := NewMoney(0)
	for _, activity := range aw.activities {
		if activity.SourceAccountId == accountId {
			withdrawalBalance = withdrawalBalance.Minus(activity.Money)
		}
	}

	return depositBalance.Plus(withdrawalBalance)
}

package account_domain

import (
	"time"
)

type ActivityId int

type Activity struct {
	Id              ActivityId
	OwnerAccountId  AccountId
	SourceAccountId AccountId
	TargetAccountId AccountId
	Timestamp       time.Time
	Money           Money
}

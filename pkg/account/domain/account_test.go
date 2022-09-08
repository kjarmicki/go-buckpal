package account_domain_test

import (
	"testing"

	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	t.Run("withdrawal success", func(t *testing.T) {
		accountId := account_domain.AccountId(1)
		account := account_domain.NewAccount(
			accountId,
			account_domain.NewMoney(555),
			account_domain.NewActivityWindow([]account_domain.Activity{
				{
					TargetAccountId: accountId,
					Money:           account_domain.NewMoney(999),
				},
				{
					TargetAccountId: accountId,
					Money:           account_domain.NewMoney(1),
				},
			}),
		)

		success := account.Withdraw(account_domain.NewMoney(555), account_domain.AccountId(99))

		assert.True(t, success)
		assert.Equal(t, 3, len(account.GetActivities()))
		assert.Equal(t, int64(1000), account.CalculateBalance().GetAmount())
	})
}

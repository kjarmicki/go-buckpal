package account_application_service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	account_application_port_in "github.com/kjarmicki/go-buckpal/pkg/account/application/port/in"
	mock_account_application_port_out "github.com/kjarmicki/go-buckpal/pkg/account/application/port/out/mock"
	account_application_service "github.com/kjarmicki/go-buckpal/pkg/account/application/service"
	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
	"github.com/stretchr/testify/assert"
)

func TestSendMoneyService(t *testing.T) {
	t.Run("transaction success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		loadAccountPort := mock_account_application_port_out.NewMockLoadAccountPort(ctrl)
		accountLock := mock_account_application_port_out.NewMockAccountLock(ctrl)
		updateAccountStatePort := mock_account_application_port_out.NewMockUpdateAccountStatePort(ctrl)
		sendMoneyService := account_application_service.NewSendMoneyService(loadAccountPort, accountLock, updateAccountStatePort)

		sourceAccount := givenSourceAccount()
		targetAccount := givenTargetAccount()
		money := account_domain.NewMoney(500)
		command, err := account_application_port_in.NewSendMoneyCommand(
			sourceAccount.GetId(),
			targetAccount.GetId(),
			money,
		)
		if err != nil {
			t.Fatalf("error while creating command: %v", err)
		}

		loadAccountPort.EXPECT().LoadAccount(context.Background(), sourceAccount.GetId(), gomock.Any()).
			Return(sourceAccount, nil).
			Times(1)
		loadAccountPort.EXPECT().LoadAccount(context.Background(), targetAccount.GetId(), gomock.Any()).
			Return(targetAccount, nil).
			Times(1)

		accountLock.EXPECT().ReleaseAccount(sourceAccount.GetId()).
			After(accountLock.EXPECT().LockAccount(sourceAccount.GetId()).Times(1))
		accountLock.EXPECT().ReleaseAccount(targetAccount.GetId()).
			After(accountLock.EXPECT().LockAccount(targetAccount.GetId()).Times(1))

		updateAccountStatePort.EXPECT().UpdateActivities(context.Background(), sourceAccount).Times(1)
		updateAccountStatePort.EXPECT().UpdateActivities(context.Background(), targetAccount).Times(1)

		success, err := sendMoneyService.SendMoney(context.Background(), command)
		if err != nil {
			t.Fatalf("error while sending money: %v", err)
		}

		assert.True(t, success)
		assert.Equal(t, sourceAccount.CalculateBalance().GetAmount(), account_domain.NewMoney(500).GetAmount())
		assert.Equal(t, targetAccount.CalculateBalance().GetAmount(), account_domain.NewMoney(500).GetAmount())
	})
}

func givenSourceAccount() *account_domain.Account {
	return givenAnAccountWithId(account_domain.AccountId(41), account_domain.NewMoney(1000))
}

func givenTargetAccount() *account_domain.Account {
	return givenAnAccountWithId(account_domain.AccountId(42), account_domain.NewMoney(0))
}

func givenAnAccountWithId(accountId account_domain.AccountId, baselineBalance account_domain.Money) *account_domain.Account {
	return account_domain.NewAccount(accountId, baselineBalance, account_domain.NewActivityWindow(make([]account_domain.Activity, 0)))
}

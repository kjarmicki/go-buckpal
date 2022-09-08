package account_application_service

import (
	"context"
	"time"

	account_application_port_in "github.com/kjarmicki/go-buckpal/pkg/account/application/port/in"
	account_application_port_out "github.com/kjarmicki/go-buckpal/pkg/account/application/port/out"
)

type SendMoneyService struct {
	loadAccountPort        account_application_port_out.LoadAccountPort
	accountLock            account_application_port_out.AccountLock
	updateAccountStatePort account_application_port_out.UpdateAccountStatePort
}

func NewSendMoneyService(
	loadAccountPort account_application_port_out.LoadAccountPort,
	accountLock account_application_port_out.AccountLock,
	updateAccountStatePort account_application_port_out.UpdateAccountStatePort,
) *SendMoneyService {
	return &SendMoneyService{
		loadAccountPort:        loadAccountPort,
		accountLock:            accountLock,
		updateAccountStatePort: updateAccountStatePort,
	}
}

func (sms *SendMoneyService) SendMoney(ctx context.Context, command *account_application_port_in.SendMoneyCommand) (bool, error) {
	baselineDate := time.Now().AddDate(0, 0, -10)

	sourceAccount, err := sms.loadAccountPort.LoadAccount(ctx, command.GetSourceAccountId(), baselineDate)
	if err != nil {
		return false, err
	}

	targetAccount, err := sms.loadAccountPort.LoadAccount(ctx, command.GetTargetAccountId(), baselineDate)
	if err != nil {
		return false, err
	}

	err = sms.accountLock.LockAccount(sourceAccount.GetId())
	if err != nil {
		return false, err
	}
	defer func() {
		sms.accountLock.ReleaseAccount(sourceAccount.GetId())
	}()
	if !sourceAccount.Withdraw(command.GetMoney(), targetAccount.GetId()) {
		return false, nil
	}

	err = sms.accountLock.LockAccount(targetAccount.GetId())
	if err != nil {
		return false, err
	}
	defer func() {
		sms.accountLock.ReleaseAccount(targetAccount.GetId())
	}()
	if !targetAccount.Deposit(command.GetMoney(), sourceAccount.GetId()) {
		return false, nil
	}

	err = sms.updateAccountStatePort.UpdateActivities(ctx, sourceAccount)
	if err != nil {
		return false, err
	}

	err = sms.updateAccountStatePort.UpdateActivities(ctx, targetAccount)
	if err != nil {
		return false, err
	}

	return true, nil
}

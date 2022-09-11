package account_application_port_in

import "context"

type SendMoneyUseCase interface {
	SendMoney(ctx context.Context, command *SendMoneyCommand) (bool, error)
}

package account_application_port_in

type SendMoneyUseCase interface {
	SendMoney(command *SendMoneyCommand) (bool, error)
}

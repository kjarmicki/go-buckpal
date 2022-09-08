package account_application_port_in

import (
	"github.com/go-playground/validator/v10"
	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
)

type SendMoneyCommand struct {
	sourceAccountId account_domain.AccountId `validate:"required"`
	targetAccountId account_domain.AccountId `validate:"required"`
	money           account_domain.Money     `validate:"required,money_amount_positive"`
}

func NewSendMoneyCommand(sourceAccountId, targetAccountId account_domain.AccountId, money account_domain.Money) (*SendMoneyCommand, error) {
	v := validator.New()
	err := v.RegisterValidation("money_amount_positive", func(fl validator.FieldLevel) bool {
		return true
	})
	if err != nil {
		return nil, err
	}

	command := &SendMoneyCommand{
		sourceAccountId: sourceAccountId,
		targetAccountId: targetAccountId,
		money:           money,
	}
	err = v.Struct(command)
	if err != nil {
		return nil, err
	}
	return command, nil
}

func (smc *SendMoneyCommand) GetSourceAccountId() account_domain.AccountId {
	return smc.sourceAccountId
}

func (smc *SendMoneyCommand) GetTargetAccountId() account_domain.AccountId {
	return smc.targetAccountId
}

func (smc *SendMoneyCommand) GetMoney() account_domain.Money {
	return smc.money
}

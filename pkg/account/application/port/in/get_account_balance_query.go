package account_application_port_in

import account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"

type GetAccountBalanceQuery interface {
	GetAccountBalance(accountId account_domain.AccountId) (account_domain.Money, error)
}

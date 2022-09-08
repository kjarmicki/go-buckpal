package account_application_port_out

import account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"

type AccountLock interface {
	LockAccount(accountId account_domain.AccountId) error
	ReleaseAccount(accountId account_domain.AccountId)
}

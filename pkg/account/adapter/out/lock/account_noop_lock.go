package account_adapter_out_lock

import account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"

type AccountNoopLockAdapter struct{}

func (l *AccountNoopLockAdapter) LockAccount(accountId account_domain.AccountId) error {
	return nil
}

func (l *AccountNoopLockAdapter) ReleaseAccount(accountId account_domain.AccountId) {

}

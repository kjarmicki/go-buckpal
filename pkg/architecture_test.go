package pkg_test

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestDomainDependencies(t *testing.T) {
	t.Run("domain doesn't depend on application", func(t *testing.T) {
		archtest.Package(t, "github.com/kjarmicki/go-buckpal/pkg/account/domain").
			ShouldNotDependOn("github.com/kjarmicki/go-buckpal/pkg/account/application/...")
	})

	t.Run("domain doesn't depend on adapters", func(t *testing.T) {
		archtest.Package(t, "github.com/kjarmicki/go-buckpal/pkg/account/domain").
			ShouldNotDependOn("github.com/kjarmicki/go-buckpal/pkg/account/adapter/...")
	})
}

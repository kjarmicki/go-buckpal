package account_adapter_in_web_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	account_adapter_in_web "github.com/kjarmicki/go-buckpal/pkg/account/adapter/in/web"
	account_application_port_in "github.com/kjarmicki/go-buckpal/pkg/account/application/port/in"
	mock_account_application_port_in "github.com/kjarmicki/go-buckpal/pkg/account/application/port/in/mock"
	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router
var server *httptest.Server
var client *http.Client

func TestMain(m *testing.M) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}))
	client = server.Client()
	defer server.Close()
	os.Exit(m.Run())
}

func TestSendMoneyController(t *testing.T) {
	t.Run("test send money", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sendMoneyUseCase := mock_account_application_port_in.NewMockSendMoneyUseCase(ctrl)
		sourceAccountNumber := 41
		targetAccountNumber := 42
		amountOfMoney := int64(500)
		command, _ := account_application_port_in.NewSendMoneyCommand(
			account_domain.AccountId(sourceAccountNumber),
			account_domain.AccountId(targetAccountNumber),
			account_domain.NewMoney(amountOfMoney),
		)
		controller := account_adapter_in_web.NewSendMoneyController(sendMoneyUseCase)
		router = mux.NewRouter()
		controller.AttachToRouter(router)
		sendMoneyUseCase.EXPECT().SendMoney(command).Return(true, nil).Times(1)

		req, _ := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/accounts/send/%d/%d/%d", server.URL, sourceAccountNumber, targetAccountNumber, amountOfMoney),
			nil,
		)

		res, err := client.Do(req)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

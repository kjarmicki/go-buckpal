package account_adapter_in_web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	account_application_port_in "github.com/kjarmicki/go-buckpal/pkg/account/application/port/in"
	account_domain "github.com/kjarmicki/go-buckpal/pkg/account/domain"
)

type SendMoneyController struct {
	sendMoneyUseCase account_application_port_in.SendMoneyUseCase
}

func NewSendMoneyController(sendMoneyUseCase account_application_port_in.SendMoneyUseCase) *SendMoneyController {
	return &SendMoneyController{
		sendMoneyUseCase: sendMoneyUseCase,
	}
}

func (smc *SendMoneyController) AttachToRouter(router *mux.Router) {
	router.
		HandleFunc("/accounts/send/{sourceAccountId}/{targetAccountId}/{amount}", smc.SendMoney).
		Methods("POST")
}

func (smc *SendMoneyController) SendMoney(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceAccountId, err := strconv.Atoi(vars["sourceAccountId"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	targetAccountId, err := strconv.Atoi(vars["targetAccountId"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	amount, err := strconv.Atoi(vars["amount"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	sendMoneyCommand, err := account_application_port_in.NewSendMoneyCommand(
		account_domain.AccountId(sourceAccountId), account_domain.AccountId(targetAccountId), account_domain.NewMoney(int64(amount)),
	)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := smc.sendMoneyUseCase.SendMoney(r.Context(), sendMoneyCommand)
	if err != nil || !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

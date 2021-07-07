package usecase

import (
	"encoding/json"
	"os"
	"time"

	"github.com/welitonhs/fccodebank/domain"
	"github.com/welitonhs/fccodebank/dto"
	"github.com/welitonhs/fccodebank/infrastructure/kafka"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
	KafkaProducer         kafka.KafkaProducer
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transaction dto.Transaction) (domain.Transaction, error) {
	creditCard := u.hydrateCreditCard(transaction)
	ccBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}

	creditCard.ID = ccBalanceAndLimit.ID
	creditCard.Limit = ccBalanceAndLimit.Limit
	creditCard.Balance = ccBalanceAndLimit.Balance

	t := u.NewTransaction(transaction, ccBalanceAndLimit)
	t.ProcessAndValidate(creditCard)

	err = u.TransactionRepository.SaveTransaction(*t, *creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	transaction.ID = t.ID
	transaction.CreatedAt = t.CreatedAt

	transactionJson, err := json.Marshal(transaction)

	if err != nil {
		return *t, err
	}

	err = u.KafkaProducer.Publish(string(transactionJson), os.Getenv("KafkaTransactionsTopic"))

	if err != nil {
		return *t, err
	}

	return *t, nil
}

func (u UseCaseTransaction) hydrateCreditCard(transactionDTO dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDTO.Name
	creditCard.Number = transactionDTO.Number
	creditCard.ExpirationMonth = transactionDTO.ExpirationMonth
	creditCard.ExpirationYear = transactionDTO.ExpirationYear
	creditCard.CVV = transactionDTO.CVV
	return creditCard
}

func (u UseCaseTransaction) NewTransaction(transactionDTO dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	transaction := domain.NewTransaction()
	transaction.CreditCard = cc.ID
	transaction.Amount = transactionDTO.Amount
	transaction.Store = transactionDTO.Store
	transaction.Description = transactionDTO.Description
	transaction.CreatedAt = time.Now()
	return transaction
}

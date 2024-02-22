package domain

import (
	"context"
	"gorm.io/gorm"
	errorsDom "jugaldb.com/byob_task/src/internal/domain/errors"
)

type Transaction struct {
	db *gorm.DB
}

type TxnContext struct {
	gormTx *gorm.DB
}

func (tc *TxnContext) GetGormTxFromTx() (*gorm.DB, error) {
	if tc.gormTx == nil {
		return nil, errorsDom.InvalidDbTransaction()
	}
	return tc.gormTx, nil
}

// Atomic Used for creating atomic db transaction.
// IN: function containing all transactions to perform
// OUT: error if there is any otherwise nil
func (t *Transaction) Atomic(ctx context.Context, fn func(tc *TxnContext) error) error {
	return t.db.Transaction(func(gormTx *gorm.DB) error {
		transactionContext := &TxnContext{gormTx: gormTx}
		err := fn(transactionContext)
		// Handles context timeout
		ctx_err := ctx.Err()
		if ctx_err != nil {
			return errorsDom.ContextTimeoutError(ctx_err)
		}
		return err
	})
}

// NewTxnContext Usually used for creating new transaction.
// IN: *gorm.DB instance
// OUT: Transaction Context
func NewTxnContext(initialisedDB *gorm.DB) TxnContext {
	return TxnContext{gormTx: initialisedDB}
}

// NewTransactionHelper Return Transaction instance which can be used for Atomic transaction
func NewTransactionHelper(db *gorm.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

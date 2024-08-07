package transaction

import (
	"context"

	"github.com/andredubov/golibs/pkg/client/database"
	"github.com/andredubov/golibs/pkg/client/database/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db database.Transactor
}

// NewTransactionManager creates a new transaction manager that satisfies the db.txManager interface
func NewTransactionManager(db database.Transactor) database.TxManager {
	return &manager{
		db,
	}
}

// transaction is the main function that executes a user-specified handler in a transaction
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn database.Handler) (err error) {
	// if this is a nested transaction, skip the initiation of a new transaction and execute the handler
	tx, ok := ctx.Value(postgres.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// starting a new transaction
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// putting the transaction in context.
	ctx = postgres.MakeContextTx(ctx, tx)

	// Настраиваем функцию отсрочки для отката или коммита транзакции.
	defer func() {
		// восстанавливаемся после паники
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// откатываем транзакцию, если произошла ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		// если ошибок не было, коммитим транзакцию
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	// Выполните код внутри транзакции.
	// Если функция терпит неудачу, возвращаем ошибку, и функция отсрочки выполняет откат
	// или в противном случае транзакция коммитится.
	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f database.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	return m.transaction(ctx, txOpts, f)
}

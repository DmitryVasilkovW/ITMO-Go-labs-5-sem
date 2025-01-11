//go:build !solution

package ledger

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ledger struct {
	db *sql.DB
}

func (l ledger) Close() error {
	if l.db != nil {
		return l.db.Close()
	}
	return nil
}

func (l ledger) CreateAccount(ctx context.Context, id ID) error {
	_, err := l.db.ExecContext(
		ctx,
		"INSERT INTO accounts(id) VALUES($1)",
		id,
	)
	return err
}

func (l ledger) Deposit(ctx context.Context, id ID, amount Money) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	result, err := l.updateAccountBalance(ctx, id, amount)
	if err != nil {
		return err
	}

	if err := checkRowsAffected(result); err != nil {
		return err
	}

	return nil
}

func validateAmount(amount Money) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	return nil
}

func (l ledger) updateAccountBalance(ctx context.Context, id ID, amount Money) (sql.Result, error) {
	return l.db.ExecContext(
		ctx,
		`
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2;
		`,
		amount,
		id,
	)
}

func checkRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("account does not exist")
	}

	return nil
}

func (l ledger) GetBalance(ctx context.Context, id ID) (Money, error) {
	var money Money

	err := l.db.QueryRowContext(
		ctx,
		"SELECT balance FROM accounts WHERE id = $1",
		id,
	).Scan(&money)

	if err != nil {
		return 0, err
	}

	return money, nil
}

func (l ledger) Transfer(ctx context.Context, from ID, to ID, amount Money) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := l.checkAccountsExist(ctx, tx, from, to); err != nil {
		return err
	}

	if err := l.executeTransfer(ctx, tx, from, to, amount); err != nil {
		return err
	}

	return tx.Commit()
}

func (l ledger) checkAccountsExist(ctx context.Context, tx *sql.Tx, from ID, to ID) error {
	rows, err := tx.QueryContext(
		ctx,
		"SELECT balance FROM accounts WHERE id IN ($1, $2) ORDER BY id FOR UPDATE",
		from,
		to,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	cnt := 0
	for rows.Next() {
		cnt++
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	if cnt < 2 {
		return errors.New("some accounts do not exist")
	}

	return nil
}

func (l ledger) executeTransfer(ctx context.Context, tx *sql.Tx, from ID, to ID, amount Money) error {
	_, err := tx.ExecContext(
		ctx,
		`
		CALL transfer_money($1, $2, $3);
		`,
		from,
		to,
		amount,
	)
	if err != nil {
		if err.Error() == `ERROR: new row for relation "accounts" violates check constraint "accounts_balance_check" (SQLSTATE 23514)` {
			return ErrNoMoney
		}
		return err
	}
	return nil
}

func (l ledger) Withdraw(ctx context.Context, id ID, amount Money) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	result, err := l.decreaseBalance(ctx, id, amount)
	if err != nil {
		return handleWithdrawError(err)
	}

	if err := _checkRowsAffected(result); err != nil {
		return err
	}

	return nil
}

func (l ledger) decreaseBalance(ctx context.Context, id ID, amount Money) (sql.Result, error) {
	return l.db.ExecContext(
		ctx,
		`
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2;
		`,
		amount,
		id,
	)
}

func handleWithdrawError(err error) error {
	if err.Error() == `ERROR: new row for relation "accounts" violates check constraint "accounts_balance_check" (SQLSTATE 23514)` {
		return ErrNoMoney
	}
	return err
}

func _checkRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("account does not exist")
	}

	return nil
}

func New(ctx context.Context, dsn string) (Ledger, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	l := ledger{
		db: db,
	}

	_, err = db.ExecContext(
		ctx,
		`
		CREATE TABLE IF NOT EXISTS accounts(
			id 				TEXT PRIMARY KEY,
			balance 		BIGINT NOT NULL DEFAULT 0 CHECK(balance >= 0)
		);

		CREATE OR REPLACE PROCEDURE transfer_money(
			id_from TEXT,
			id_to TEXT, 
			amount BIGINT
		)
		LANGUAGE plpgsql    
		AS $$
		BEGIN

		UPDATE accounts 
		SET
			balance = balance - amount 
		WHERE id = id_from;

		UPDATE accounts 
		SET
			balance = balance + amount 
		WHERE id = id_to;

		END;
		$$
		
		`,
	)

	if err != nil {
		defer db.Close()
		return nil, err
	}

	return l, nil
}

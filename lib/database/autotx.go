package database

import (
	"database/sql"
	"errors"
	"log"
	"runtime/debug"

	"github.com/go-openapi/runtime/middleware"
	"github.com/leporo/sqlf"
)

type AutoTx struct {
	tx    *sql.Tx
	query string
	rows  *sql.Rows
	res   sql.Result
	err   error
}

func (tx *AutoTx) Query(query string, args ...any) *AutoTx {
	if tx.HasQueryError() {
		return tx
	}

	if tx.rows != nil {
		tx.rows.Close()
	}

	tx.query = query
	tx.rows, tx.err = tx.tx.Query(query, args...)

	return tx
}

func (tx *AutoTx) QueryStmt(stmt *sqlf.Stmt) *AutoTx {
	tx.Query(stmt.String(), stmt.Args()...)
	stmt.Close()
	return tx
}

func (tx *AutoTx) Scan(dest ...any) bool {
	if tx.err != nil {
		return false
	}

	if !tx.rows.Next() {
		tx.err = tx.rows.Err()
		tx.rows = nil
		if tx.err == nil {
			tx.err = sql.ErrNoRows
		}

		return false
	}

	tx.err = tx.rows.Scan(dest...)

	return true
}

func (tx *AutoTx) ScanBool() bool {
	var result sql.NullBool
	tx.Scan(&result)
	return result.Bool
}

func (tx *AutoTx) ScanBools() []bool {
	var result []bool
	var value bool
	for tx.Scan(&value) {
		result = append(result, value)
	}

	return result
}

func (tx *AutoTx) ScanInt64() int64 {
	var result sql.NullInt64
	tx.Scan(&result)
	return result.Int64
}

func (tx *AutoTx) ScanInt64s() []int64 {
	var result []int64
	var value int64
	for tx.Scan(&value) {
		result = append(result, value)
	}

	return result
}

func (tx *AutoTx) ScanFloat64() float64 {
	var result sql.NullFloat64
	tx.Scan(&result)
	return result.Float64
}

func (tx *AutoTx) ScanString() string {
	var result sql.NullString
	tx.Scan(&result)
	return result.String
}

func (tx *AutoTx) ScanStrings() []string {
	var result []string
	var value string
	for tx.Scan(&value) {
		result = append(result, value)
	}

	return result
}

func (tx *AutoTx) QueryBool(query string, args ...any) bool {
	return tx.Query(query, args...).ScanBool()
}

func (tx *AutoTx) QueryBools(query string, args ...any) []bool {
	return tx.Query(query, args...).ScanBools()
}

func (tx *AutoTx) QueryInt64(query string, args ...any) int64 {
	return tx.Query(query, args...).ScanInt64()
}

func (tx *AutoTx) QueryInt64s(query string, args ...any) []int64 {
	return tx.Query(query, args...).ScanInt64s()
}

func (tx *AutoTx) QueryFloat64(query string, args ...any) float64 {
	return tx.Query(query, args...).ScanFloat64()
}

func (tx *AutoTx) QueryString(query string, args ...any) string {
	return tx.Query(query, args...).ScanString()
}

func (tx *AutoTx) QueryStrings(query string, args ...any) []string {
	return tx.Query(query, args...).ScanStrings()
}

func (tx *AutoTx) Close() {
	if tx.rows != nil {
		tx.rows.Close()
		tx.rows = nil
	}
}

func (tx *AutoTx) Error() error {
	return tx.err
}

func (tx *AutoTx) HasQueryError() bool {
	return tx.err != nil && !errors.Is(tx.err, sql.ErrNoRows)
}

func (tx *AutoTx) LastQuery() string {
	return tx.query
}

func (tx *AutoTx) Exec(query string, args ...any) {
	if tx.HasQueryError() {
		return
	}

	if tx.rows != nil {
		tx.rows.Close()
		tx.rows = nil
	}

	tx.query = query
	tx.res, tx.err = tx.tx.Exec(query, args...)
}

func (tx *AutoTx) ExecStmt(stmt *sqlf.Stmt) *AutoTx {
	tx.Exec(stmt.String(), stmt.Args()...)
	stmt.Close()
	return tx
}

func (tx *AutoTx) RowsAffected() int64 {
	var cnt int64
	if tx.err == nil {
		cnt, tx.err = tx.res.RowsAffected()
	}

	return cnt
}

func NewAutoTx(db *sql.DB) *AutoTx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	return &AutoTx{tx: tx}
}

func (tx *AutoTx) Finish() {
	p := recover()
	tx.Close()

	var err error

	if p != nil {
		err = tx.tx.Rollback()
		log.Println(p, " (recovered by AutoTx)")
		log.Println(tx.LastQuery())
		debug.PrintStack()
	} else if !tx.HasQueryError() {
		err = tx.tx.Commit()
	} else {
		log.Println(tx.Error())
		log.Println(tx.LastQuery())
		err = tx.tx.Rollback()
	}

	if err != nil {
		log.Println(err)
	}
}

// Transact wraps func in an SQL transaction.
// Responder will be just passed through.
func Transact(db *sql.DB, txFunc func(*AutoTx) middleware.Responder) middleware.Responder {
	atx := NewAutoTx(db)
	defer atx.Finish()

	return txFunc(atx)
}

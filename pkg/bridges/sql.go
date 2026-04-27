package bridges

import (
	"database/sql"

	"github.com/Pure-Company/purekernels/pkg/monoid"
	"github.com/Pure-Company/purekernels/pkg/result"
)

// ScanRow scans a single row into a value
func ScanRow[T any](
	rows *sql.Rows,
	scanner func(*sql.Rows) (T, error),
) result.Result[monoid.Option[T]] {
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return result.Err[monoid.Option[T]](err)
		}
		return result.Ok(monoid.None[T]())
	}

	val, err := scanner(rows)
	if err != nil {
		return result.Err[monoid.Option[T]](err)
	}

	return result.Ok(monoid.Some(val))
}

// ScanAll scans all rows into a slice
func ScanAll[T any](
	rows *sql.Rows,
	scanner func(*sql.Rows) (T, error),
) result.Result[[]T] {
	items := []T{}

	for rows.Next() {
		val, err := scanner(rows)
		if err != nil {
			return result.Err[[]T](err)
		}
		items = append(items, val)
	}

	if err := rows.Err(); err != nil {
		return result.Err[[]T](err)
	}

	return result.Ok(items)
}

// QueryOne executes a query expecting one result
func QueryOne[T any](
	db *sql.DB,
	query string,
	scanner func(*sql.Rows) (T, error),
	args ...any,
) result.Result[monoid.Option[T]] {
	rows, err := db.Query(query, args...)
	if err != nil {
		return result.Err[monoid.Option[T]](err)
	}
	defer rows.Close()

	return ScanRow(rows, scanner)
}

// QueryAll executes a query expecting multiple results
func QueryAll[T any](
	db *sql.DB,
	query string,
	scanner func(*sql.Rows) (T, error),
	args ...any,
) result.Result[[]T] {
	rows, err := db.Query(query, args...)
	if err != nil {
		return result.Err[[]T](err)
	}
	defer rows.Close()

	return ScanAll(rows, scanner)
}

// Transaction wraps a function in a transaction
func Transaction[T any](
	db *sql.DB,
	f func(*sql.Tx) result.Result[T],
) result.Result[T] {
	tx, err := db.Begin()
	if err != nil {
		return result.Err[T](err)
	}

	res := f(tx)
	if res.IsErr() {
		tx.Rollback()
		return res
	}

	if err := tx.Commit(); err != nil {
		return result.Err[T](err)
	}

	return res
}

// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package storage

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
)

type history2TransactionConvertToValues func(row history2.Transaction) []interface{}

func history2TransactionBatchInsert(repo *pgdb.DB, rows []history2.Transaction, tableName string, columns []string, converter history2TransactionConvertToValues) error {
	if len(rows) == 0 {
		return nil
	}

	sql := sq.Insert(tableName).Columns(columns...)

	paramsInQueue := 0
	for _, row := range rows {
		paramsInQueue += len(columns)
		if paramsInQueue > maxPostgresParams {
			err := repo.Exec(sql)
			if err != nil {
				return errors.Wrap(err, "failed to perform batch insert", logan.F{"rows_len": len(rows)})
			}

			sql = sq.Insert(tableName).Columns(columns...)
			paramsInQueue = len(columns)
		}

		sql = sql.Values(converter(row)...)
	}

	if paramsInQueue == 0 {
		return nil
	}

	err := repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to perform batch insert", logan.F{"rows_len": len(rows)})
	}

	return nil
}

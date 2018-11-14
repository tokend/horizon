package ingest2

import (
	"gitlab.com/tokend/horizon/db2"
	"database/sql"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	history "gitlab.com/tokend/horizon/db2/history2"
)

// storage - entity used to perform write operations into the persistent storage by ingest
type storage struct {
	// DB is the sql repo to be used for writing any rows into the horizon
	// database.
	db *db2.Repo
	// allStms - all prepared statements
	allStms []*sql.Stmt
	// account - prepared statement for the account insert
	account *sql.Stmt
}

// newStorage - creates new storage. db should not have any transactions started.
func newStorage(db *db2.Repo) (*storage, error) {
	err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin tx for ingest storage")
	}

	st := &storage{
		db: db,
	}

	err = st.init()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init storage for ingest")
	}

	return st, nil
}

func (s *storage) init() error {
	err := s.initAccount()
	return err
}

func (s *storage) initAccount() error {
	var err error
	s.account, err = s.db.DB.Prepare("INSERT INTO accounts(id, address, account_type) VALUES(?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement for accounts")
	}
	s.allStms = append(s.allStms, s.account)
	return nil
}

func (s *storage) InsertAccounts(accounts []history.Account) error {
	for _, account := range accounts {
		_, err := s.account.Exec(account.ID, account.Address, account.AccountType)
		if err != nil {
			return errors.Wrap(err, "failed to insert account")
		}
	}

	return nil
}

func (s *storage) closeStatements() error {
	for i := range s.allStms {
		err := s.allStms[i].Close()
		if err != nil {
			return errors.Wrap(err, "failed to close statements", logan.Field("i of statement", i))
		}
	}

	return nil
}

func (s *storage) Commit() error {
	err := s.closeStatements()
	if err != nil {
		return errors.Wrap(err, "failed to close statements")
	}

	err = s.db.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit db tx")
	}

	return nil
}

func (s *storage) Rollback() error {
	err := s.closeStatements()
	if err != nil {
		return errors.Wrap(err, "failed to close statements")
	}

	err = s.db.Rollback()
	if err != nil {
		return errors.Wrap(err, "failed to rollback db tx")
	}

	return nil
}
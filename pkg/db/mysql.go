package db

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/okiww/billing-loan-system/configs"
	log "github.com/sirupsen/logrus"
)

type DBMySQL struct {
	cfg    *configs.DBConfig
	DB     *sqlx.DB
	ExecTx TxExecutor
}

type DB struct {
	DB     *sqlx.DB
	ExecTx TxExecutor
}

type DBInterface interface {
	Connect() (*DBMySQL, error)
	CloseDB() error
}

func InitDB(cfg *configs.DBConfig) DBInterface {
	return &DBMySQL{
		cfg,
		nil,
		ExecTx,
	}
}

func (d *DBMySQL) Connect() (*DBMySQL, error) {
	dbMySQL, err := sqlx.Connect("mysql", d.cfg.Source)
	if err != nil {
		return nil, err
	}

	dbMySQL.SetMaxOpenConns(200)
	dbMySQL.SetMaxIdleConns(10)

	log.WithFields(log.Fields{
		"dsn":  d.cfg.Source,
		"name": d.cfg.DBName,
	}).Info("Success connect to db")

	return &DBMySQL{
		DB:     dbMySQL,
		ExecTx: ExecTx,
	}, nil
}

// TxExecutor accesses ExecTx from outer package
type TxExecutor func(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error

// ExecTx runs fn inside tx which should already have begun.
func ExecTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (d *DBMySQL) CloseDB() error {
	if d.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	err := d.DB.Close()
	if err != nil {
		log.WithError(err).Error("Failed to close the database connection")
		return err
	}

	log.Info("Database connection closed successfully")
	return nil
}

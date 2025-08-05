package main

import (
	"database/sql"
	"fmt"
)

func TransferMoney(db *sql.DB, from, to int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Print("Rollback failed")
				return
			}
		}
	}()

	if _, err = tx.Exec("UPDATE...", from, -amount); err != nil {
		return err
	}

	if _, err = tx.Exec("UPDATE...", to, amount); err != nil {
		return err
	}

	return tx.Commit()
}

package tests

import (
	"database/sql"
	"fmt"
)

func Up_20240615204715(txn *sql.Tx) {
	fmt.Println("Hello from migration 20240615204715 Up!")
}

func Down_20240615204715(txn *sql.Tx) {
	fmt.Println("Hello from migration 20240615204715 Down!")
}

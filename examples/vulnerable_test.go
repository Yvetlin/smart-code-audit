package examples

import (
	"crypto/md5"
	"crypto/tls"
	"database/sql"
	"fmt"
	"testing"
)

func TestSQLInjection(t *testing.T) {
	user := "admin' OR '1'='1"
	query := "SELECT * FROM users WHERE name = '" + user + "'"

	var db *sql.DB
	_, _ = db.Query(query)

	t.Log("SQL injection test")
}

func TestHardcodedCrypto(t *testing.T) {
	sum := md5.Sum([]byte("password"))
	fmt.Printf("%x\n", sum)
}

func TestInsecureTLS(t *testing.T) {
	_ = &tls.Config{
		InsecureSkipVerify: true,
	}
}

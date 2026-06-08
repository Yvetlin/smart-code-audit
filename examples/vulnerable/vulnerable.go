package vulnerable

import (
	"crypto/md5"
	"crypto/tls"
	"database/sql"
	"fmt"
)

// ❌ SQL Injection (G201)
func SQLInjection(db *sql.DB, user string) {
	query := "SELECT * FROM users WHERE name = '" + user + "'"
	_, _ = db.Query(query)
}

// ❌ Weak crypto (G401)
func WeakCrypto() {
	sum := md5.Sum([]byte("password"))
	fmt.Printf("%x\n", sum)
}

// ❌ Insecure TLS (G402)
func InsecureTLS() {
	_ = &tls.Config{
		InsecureSkipVerify: true,
	}
}

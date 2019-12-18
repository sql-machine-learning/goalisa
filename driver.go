package goalisa

import (
	"database/sql"
	"database/sql/driver"
)

// register driver
func init() {
	sql.Register("alisa", &Driver{})
}

// Driver implements database/sql/driver.Driver
type Driver struct{}

// Open returns a connection
func (dr Driver) Open(dsn string) (driver.Conn, error) {
	return nil, nil
}

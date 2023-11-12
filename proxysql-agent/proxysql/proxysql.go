package proxysql

import (
	"database/sql"
	// "fmt"
	//"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

type ProxySQL struct {
	dsn    string
	conn   *sql.DB
	logger zerolog.Logger
}

func (p *ProxySQL) Conn() *sql.DB {
	return p.conn
}

func (p *ProxySQL) Ping() error {
	return p.conn.Ping()
}

func (p *ProxySQL) Close() {
	p.conn.Close()
}

func (p *ProxySQL) PersistChanges() error {
	// SAVE ALL THINGS
	// _, err := p.conn.Exec("save mysql servers to disk")

	// if err != nil {
	// 	return err
	// }
	// _, err = p.conn.Exec("load mysql servers to runtime")
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (p *ProxySQL) GetBackends() (map[string]int, error) {
	entries := make(map[string]int)

	rows, err := p.conn.Query("SELECT hostgroup_id, hostname, port FROM runtime_mysql_servers ORDER BY hostgroup_id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var hostgroup int
		var hostname string
		var port int

		err := rows.Scan(&hostgroup, &hostname, &port)
		if err != nil {
			return nil, err
		}

		entries[hostname] = hostgroup
		if rows.Err() != nil && rows.Err() != sql.ErrNoRows {
			return nil, rows.Err()
		}
	}

	return entries, nil
}

// CheckConnection periodically queries the ProxySQL runtime_mysql_servers table to check the status of backend servers.
// It retrieves the hostgroup ID and hostname of each backend server and logs the information using the provided logger.
//
// The function runs in an infinite loop, periodically querying the table. The sleep interval between queries is set to 10 seconds,
// but you can configure this interval as needed for your specific use case.
//
// Parameters:
// - p: A pointer to the ProxySQL instance with an active database connection.
//
// Notes:
// - The query results are logged using the logger associated with the ProxySQL instance.
func (p *ProxySQL) CheckConnection() {
	for {
		query := "SELECT hostgroup_id, hostname FROM runtime_mysql_servers"
		rows, err := p.conn.Query(query)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		// Process the query results (you can replace this with your logic)
		for rows.Next() {
			var col1, col2 string
			if err := rows.Scan(&col1, &col2); err != nil {
				panic(err)
			}

			p.logger.Debug().Str("hg", col1).Str("hostname", col2).Msg("Backends")
		}

		// FIXME: make this configurable
		time.Sleep(10 * time.Second)
	}
}

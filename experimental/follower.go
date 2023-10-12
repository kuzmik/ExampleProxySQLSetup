package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:6032)/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	missingCommand := `SELECT COUNT(hostname) FROM stats_proxysql_servers_metrics WHERE last_check_ms > 30000 AND hostname != 'proxysql-leader' AND Uptime_s > 0`
	allCommand := `SELECT COUNT(hostname) FROM stats_proxysql_servers_metrics`

	missingCount, err := executeMySQLQuery(db, missingCommand)
	if err != nil {
		log.Fatal(err)
	}

	allCount, err := executeMySQLQuery(db, allCommand)
	if err != nil {
		log.Fatal(err)
	}

	if missingCount > 0 {
		fmt.Printf("%d/%d proxysql leaders haven't been seen in over 30s, resetting leader state\n", missingCount, allCount)

		commands := []string{
			"DELETE FROM proxysql_servers",
			"LOAD PROXYSQL SERVERS FROM CONFIG",
			"LOAD PROXYSQL SERVERS TO RUNTIME",
		}
		combinedCommands := strings.Join(commands, "; ")

		_, err := db.Exec(combinedCommands)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func executeMySQLQuery(db *sql.DB, query string) (int, error) {
	row := db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

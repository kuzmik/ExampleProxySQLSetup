package proxysql

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
)

func New(dsn string) (*ProxySQL, error) {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	// FIXME: dont log full dsn, it has passwords in it
	logger.Info().Str("Host", dsn).Msg("Connected to ProxySQL admin")

	return &ProxySQL{dsn, conn, logger}, nil
}



package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/kuzmik/proxysql-cluster-agent/proxysql"
)

var (
	logger zerolog.Logger
	psql   *proxysql.ProxySQL

	coreModeFlag      bool
	satelliteModeFlag bool
	userFlag          string
	passwordFlag      string
	addressFlag       string
	pauseFlag         int
)

func main() {
	setupLogger()

	userEnv := os.Getenv("PROXYSQL_USER")
	passwordEnv := os.Getenv("PROXYSQL_PASSWORD")
	addressEnv := os.Getenv("PROXYSQL_ADDRESS")

	// If environment variables are not set, use command line arguments
	flag.StringVar(&userFlag, "user", userEnv, "ProxySQL username")
	flag.StringVar(&passwordFlag, "password", passwordEnv, "ProxySQL password")
	flag.StringVar(&addressFlag, "address", addressEnv, "ProxySQL address")

	flag.IntVar(&pauseFlag, "pause", 0, "Seconds to pause before attempting to start")

	flag.BoolVar(&coreModeFlag, "core", false, "Run the functions required for core pods")
	flag.BoolVar(&satelliteModeFlag, "satellite", false, "Run the functions required for satellite pods")

	flag.Parse()

	// If command line arguments are not set, use default values
	if userFlag == "" {
		userFlag = "radmin"
	}
	if passwordFlag == "" {
		passwordFlag = "radmin"
	}
	if addressFlag == "" {
		addressFlag = "127.0.0.1:6032"
	}

	logger.Debug().
		Str("username", userFlag).
		Str("password", passwordFlag).
		Str("address", addressFlag).
		Msg("ProxySQL configuration")

	if pauseFlag > 0 {
		logger.Info().Int("seconds", pauseFlag).Msg("Pausing before boot")
		time.Sleep(time.Duration(pauseFlag) * time.Second)
	}

	setupProxySQL()

	if coreModeFlag {
		go psql.Core()
	} else if satelliteModeFlag {
		go psql.Satellite()
	}

	setupUnixSocket()
}

func setupLogger() {
	logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
}

func setupProxySQL() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", userFlag, passwordFlag, addressFlag)
	psql, err = proxysql.New(dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to connect to ProxySQL")
	}
}

func setupUnixSocket() {
	socketPath := "/tmp/proxysql_cnc.sock"
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Error().Err(err).Str("socket", socketPath).Msg("Unable to create unix socket")
		return
	}
	defer listener.Close()

	logger.Info().Str("path", socketPath).Msg("Unix socket listening for commands")

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error().Err(err).Str("socket", socketPath).Msg("Unable to read from unix socket")
			continue
		}

		go handleSocketCommand(conn)
	}
}

func handleSocketCommand(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Error().Err(err).Msg("Error reading from socket")
		return
	}

	command := strings.TrimRight(string(buffer[:n]), "\n")
	logger.Debug().Str("command", command).Msg("Got command from socket")

	var msg string

	switch command {
	case "ping":
		go psql.Ping()
		msg = "PONG"

	case "get_backends":
		backends, err := psql.GetBackends()
		if err != nil {
			logger.Error().Err(err).Msg("Error in running Backends")
		}

		for host, id := range backends {
			logger.Info().Str("hostname", host).Int("hg", id).Send()
			msg = fmt.Sprintf("hg:%d, host:%s", id, host)
			conn.Write([]byte(msg + "\n"))
		}

	case "resync":
		go psql.SatelliteResync()
		msg = "Running resync"

	case "get_core_pods":
		// go clustering.GetPods("core")
		msg = "Got core pods" //FIXME

	case "get_satellite_pods":
		// go clustering.GetPods("satellite")
		msg = "Got satellite pods" //FIXME

	default:
		logger.Warn().Str("command", command).Msg("Unprocessable command received")
		msg = fmt.Sprintf("Unprocessable command: %s", command)
	}
	conn.Write([]byte(msg + "\n"))
}

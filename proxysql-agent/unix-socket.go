package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// FIXME: make configurable
var socketPath = "/tmp/proxysql_cnc.sock"

func setupUnixSocket() {
	handleSignals()

	//delete socket if it still exists for some reason
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Error().Err(err).Str("socket", socketPath).Msg("Unable to create unix socket")
		return
	}

	defer listener.Close()
	defer removeUnixSocket()

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

// read from the socket and process any commands that are valid. can also write back
// to the socket if required.
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

// function to trap signals (ctrl+c) and cleanup the socket
func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		removeUnixSocket()
		os.Exit(1)
	}()
}

func removeUnixSocket() {
	logger.Debug().Str("socket", socketPath).Msg("Removing unix socket")
	if err := os.Remove(socketPath); err != nil {
		logger.Err(err).Str("socket", socketPath).Msg("Error removing unix socket")
	}
}

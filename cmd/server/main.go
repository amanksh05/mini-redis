package main

import (
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)

	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found; using environment variables")
	}
	protocol := os.Getenv("PROTOCOL")
	if protocol == "" {
		logger.Error("Protocol string not set.")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		logger.Error("Port string not set.")
		return
	}

	listener, err := net.Listen(protocol, port)
	if err != nil {
		logger.Error("Failed to initialize TCP server", "error", err)
		return
	}
	defer listener.Close()

	logger.Info("Server listening", "address", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Failed to accept connection", "error", err)
			continue
		}

		logger.Info("Client connected", "remote_addr", conn.RemoteAddr())

		go handleConnection(logger, conn)
	}
}

func handleConnection(logger *slog.Logger, conn net.Conn) {
	defer conn.Close()

	connLogger := logger.With(
		"remote_addr", conn.RemoteAddr().String(),
	)

	connLogger.Info("Handling client")

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			message := string(buffer[:n])
			connLogger.Info("Message received", "message", message)

			_, err = conn.Write(buffer[:n])
			if err != nil {
				connLogger.Error("Failed to write back to the client", "error", err)
				return
			}

		}
		if err != nil {
			if err == io.EOF {
				connLogger.Info("Client closed the connection!")
			} else {
				connLogger.Error("Read failed", "error", err)
			}
			return
		}

	}
}

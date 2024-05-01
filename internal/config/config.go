package config

import (
	"errors"
	"flag"
	"os"
	"time"
)

const (
	MaxRoutine                  = 20
	DefaultServerAddress        = "localhost:8080"
	DefaultAccuralSystemAddress = "htpp://localhost:8081"
	TokenExpiration             = 24 * time.Hour
	JwtKeySecret                = "very-very-secret-key"
)

type contextKey string

const UserContextKey = contextKey("login")

var (
	ServerAddr           string
	DBDSN                string
	AccuralSystemAddress string
)

var (
	ErrNoServerAddress        = errors.New("server address is empty")
	ErrNoDatabaseDestination  = errors.New("database destination is empty")
	ErrNoAccuralSystemAddress = errors.New("no address for accural system")
	ErrTokenInvalid           = errors.New("token is not valid")
	ErrLoginExists            = errors.New("user is already exists")
	ErrWrongPassword          = errors.New("wrong password")
	ErrGone                   = errors.New("user was deleted")
	ErrNoFunds                = errors.New("insufficient funds")
	ErrNotFound               = errors.New("order not found")
)

func ConfigInit() error {
	flag.StringVar(&ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
	flag.StringVar(&DBDSN, "d", "", "base file path to save URLs")
	flag.StringVar(&AccuralSystemAddress, "r", "", "address for accural system")

	flag.Parse()

	if serverAddr := os.Getenv("RUN_ADDRESS"); serverAddr != "" {
		ServerAddr = serverAddr
	}
	if dbdsn := os.Getenv("DATABASE_URI"); dbdsn != "" {
		DBDSN = dbdsn
	}
	if accural := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accural != "" {
		AccuralSystemAddress = accural
	}

	// DBDSN = "postgres://postgres:7513@localhost:5432/postgres"
	// ServerAddr = DefaultServerAddress
	// AccuralSystemAddress = DefaultAccuralSystemAddress

	return checkConfig()
}

func checkConfig() error {
	switch {
	case ServerAddr == "":
		return ErrNoServerAddress
	case DBDSN == "":
		return ErrNoDatabaseDestination
	case AccuralSystemAddress == "":
		return ErrNoAccuralSystemAddress
	default:
		return nil
	}
}

package helper

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage/sql"
)

func sslModeToString(sslEnable bool) string {
	if sslEnable {
		return "enable"
	}

	return "disable"
}

// BuildDsn build dsn string from params.
func BuildDsn(host string, port int, user string, password string, dbName string, sslEnable bool) string {
	u := &url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, strconv.Itoa(port)),
		Path:     dbName,
		RawQuery: strings.Join([]string{"sslmode", sslModeToString(sslEnable)}, "="),
	}

	return u.String()
}

// CreateStorage creates storage by type.
func CreateStorage(dbType string, host string, port int, user string, password string, dbName string, sslMode bool) (storage.Storage, error) {
	switch strings.ToLower(dbType) {
	case "memory":
		return memorystorage.New(), nil
	case "postgres":
		dsn := BuildDsn(
			host,
			port,
			user,
			password,
			dbName,
			sslMode,
		)

		return sqlstorage.New(dsn), nil
	default:
		return nil, errors.New("unsupported storage type")
	}
}

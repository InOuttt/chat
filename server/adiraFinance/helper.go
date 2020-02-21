package adiraFinance

import (
	"fmt"
	"strings"

	"github.com/hako/branca"
)

// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

// DecodeConnection ...
func DecodeConnection(dsn, salt string) string {
	var (
		err       error
		hostname  string
		username  string
		password  string
		port      string
		database  string
		dsnFormat = "server=%s;user id=%s;password=%s;port=%s;database=%s"
		b         *branca.Branca
	)

	if len(salt) >= 32 {
		b = branca.NewBranca(salt[:32]) // This key must be exactly 32 bytes long.
	} else {
		b = branca.NewBranca("@ry0b4RoNpRoJ3ct4d!r@SmG4BaRoKAH") // This key must be exactly 32 bytes long.
	}

	hostname, err = b.DecodeToString(GetStringInBetween(dsn, "server=", ";"))
	if err != nil {
		hostname = GetStringInBetween(dsn, "server=", ";")
		Log.Ln("The hostname is not encrypt mode")
	}

	username, err = b.DecodeToString(GetStringInBetween(dsn, "user id=", ";"))
	if err != nil {
		username = GetStringInBetween(dsn, "user id=", ";")
		Log.Ln("The username is not encrypt mode")
	}

	password, err = b.DecodeToString(GetStringInBetween(dsn, "password=", ";"))
	if err != nil {
		password = GetStringInBetween(dsn, "password=", ";")
		Log.Ln("The password is not encrypt mode")
	}

	port, err = b.DecodeToString(GetStringInBetween(dsn, "port=", ";"))
	if err != nil {
		port = GetStringInBetween(dsn, "port=", ";")
		Log.Ln("The port is not encrypt mode")
	}

	dsn += ";"
	database, err = b.DecodeToString(GetStringInBetween(dsn, "database=", ";"))
	if err != nil {
		database = GetStringInBetween(dsn, "database=", ";")
		Log.Ln("The database is not encrypt mode")
	}

	return fmt.Sprintf(dsnFormat, hostname, username, password, port, database)
}

// DecodeDBName ...
func DecodeDBName(dbname, salt string) string {
	var (
		err      error
		database string
		b        *branca.Branca
	)

	if len(salt) >= 32 {
		b = branca.NewBranca(salt[:32]) // This key must be exactly 32 bytes long.
	} else {
		b = branca.NewBranca("@ry0b4RoNpRoJ3ct4d!r@SmG4BaRoKAH") // This key must be exactly 32 bytes long.
	}

	database, err = b.DecodeToString(dbname)
	if err != nil {
		return dbname
	}

	return database
}

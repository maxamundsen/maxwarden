package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// create
func maybeCreateSqliteDb() {
	if _, err := os.Stat("./passwords.db"); err != nil {
		fmt.Printf("Creating new sqlite database")
		err := os.WriteFile("./passwords.db", nil, 0755)
		if err != nil {
			fmt.Printf("Error creating Sqlite database.")
			os.Exit(1)
		}

		m, err := migrate.New(
			"file://./migrations",
			"sqlite3://passwords.db",
		)

		if err != nil {
			printStatus(false)
			fmt.Println(err.Error())
			os.Exit(1)
		}

		mErr := m.Up()
		if mErr != nil {
			printStatus(false)
			fmt.Println(mErr.Error())
			os.Exit(1)
		}

		printStatus(true)
	}
}

// handle the running and creation of migrations
func migrations(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: metagen migrate [up, down, goto {V}, create {migration name}]")
		os.Exit(1)
	}

	maybeCreateSqliteDb()

	m, err := migrate.New(
		"file://./migrations",
		"sqlite3://passwords.db",
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	migrateNum := 0

	if len(args) >= 3 && args[1] != "create" {
		var parseErr error
		migrateNum, parseErr = strconv.Atoi(args[2])
		if parseErr != nil {
			fmt.Println("Please provide a valid migration number.")
			os.Exit(1)
		}
	}

	m.PrefetchMigrations = migrate.DefaultPrefetchMigrations

	switch args[1] {
	case "up":
		err := m.Up()
		if err != nil {
			fmt.Println(err.Error())
		}
	case "down":
		err := m.Down()
		if err != nil {
			fmt.Println(err.Error())
		}
	case "goto":
		err := m.Migrate(uint(migrateNum))
		if err != nil {
			fmt.Println(err.Error())
		}
	case "create":
		if len(args) < 3 {
			fmt.Println("Please provide a name for the new migration.")
			os.Exit(1)
		}
		createCmd("./migrations", time.Now(), defaultTimeFormat, args[2], "sql", true, 7, true)
	}
}

const (
	defaultTimeFormat = "20060102150405"
	defaultTimezone   = "UTC"
)

var (
	errInvalidSequenceWidth     = errors.New("Digits must be positive")
	errIncompatibleSeqAndFormat = errors.New("The seq and format options are mutually exclusive")
	errInvalidTimeFormat        = errors.New("Time format may not be empty")
)

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	// os.Create() specifies 0666 as the FileMode, so we're doing the same
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}

func nextSeqVersion(matches []string, seqDigits int) (string, error) {
	if seqDigits <= 0 {
		return "", errInvalidSequenceWidth
	}

	nextSeq := uint64(1)

	if len(matches) > 0 {
		filename := matches[len(matches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("Malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("Next sequence number %s too large. At most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func timeVersion(startTime time.Time, format string) (version string, err error) {
	switch format {
	case "":
		err = errInvalidTimeFormat
	case "unix":
		version = strconv.FormatInt(startTime.Unix(), 10)
	case "unixNano":
		version = strconv.FormatInt(startTime.UnixNano(), 10)
	default:
		version = startTime.Format(format)
	}

	return
}

func createCmd(dir string, startTime time.Time, format string, name string, ext string, seq bool, seqDigits int, print bool) error {
	if seq && format != defaultTimeFormat {
		return errIncompatibleSeqAndFormat
	}

	var version string
	var err error

	dir = filepath.Clean(dir)
	ext = "." + strings.TrimPrefix(ext, ".")

	if seq {
		matches, err := filepath.Glob(filepath.Join(dir, "*"+ext))

		if err != nil {
			return err
		}

		version, err = nextSeqVersion(matches, seqDigits)

		if err != nil {
			return err
		}
	} else {
		version, err = timeVersion(startTime, format)

		if err != nil {
			return err
		}
	}

	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)

	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			return err
		}

		if print {
			absPath, _ := filepath.Abs(filename)
			log.Println(absPath)
		}
	}

	return nil
}

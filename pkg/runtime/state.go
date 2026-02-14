// Package runtime contains state references for a given runtime mode.
package runtime

import (
	"database/sql"
	"database/sql/driver"
	"gg/pkg/crypt/aesgcm"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"modernc.org/sqlite"
)

func init() {
	sqlite.MustRegisterScalarFunction("uuid", 0,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			return uuid.New().String(), nil
		},
	)

}

const (
	ApiVersion = "v0"
)

type KeyGetter func(uuid.UUID) []byte

// State contains references for a given runtime mode.
type State struct {
	ApiVersion           string
	DB                   *sql.DB
	EncryptionKeyVersion uuid.UUID
	GetKey               KeyGetter
	Logger               *slog.Logger
	Mode                 string
}

// NewForTest creates a State instance for test mode.
func NewForTest() *State {
	// Set DB.
	rootPath, ok := os.LookupEnv("GG_ROOT")
	if !ok {
		log.Fatal("GG_ROOT env var not defined")
	}

	root, err := os.OpenRoot(rootPath)
	if err != nil {
		log.Fatal(err)
	}
	defer root.Close() // nolint:errcheck

	schemaPath, ok := os.LookupEnv("GG_SCHEMA")
	if !ok {
		log.Fatal("GG_SCHEMA env var not defined")
	}

	schemaBs, err := root.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(schemaBs))
	if err != nil {
		log.Fatal(err)
	}

	// Set EncryptionKeyVersion.
	encryptionKeyVersion := uuid.New()

	// Set GetKey.

	// For test mode, add in two random keys in
	// addition to encryptionKeyVersion.
	keys := map[uuid.UUID][]byte{
		encryptionKeyVersion: aesgcm.NewKey(),
		uuid.New():           aesgcm.NewKey(),
		uuid.New():           aesgcm.NewKey(),
	}

	getKey := func(version uuid.UUID) []byte {
		if key, ok := keys[version]; ok {
			return key
		} else {
			return nil
		}
	}

	// Set Logger.
	logger := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{AddSource: true, Level: slog.LevelError},
	))

	return &State{
		ApiVersion:           ApiVersion,
		DB:                   db,
		EncryptionKeyVersion: encryptionKeyVersion,
		GetKey:               getKey,
		Logger:               logger,
		Mode:                 "test",
	}
}

// Close closes long-held references.
func (s *State) Close() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

// Package runtime contains state references for a given runtime mode.
package runtime

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"modernc.org/sqlite"
)

func init() {
	sqlite.MustRegisterScalarFunction("uuid", 0,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			v, err := uuid.NewRandom()
			if err != nil {
				panic(err)
			}
			return v.String(), nil
		},
	)

}

const (
	ApiVersion = "v0"
)

// State contains references for a given runtime mode.
type State struct {
	ApiVersion string
	DB         *sql.DB
	Logger     *slog.Logger
	Mode       string
}

// NewForTest creates a State instance for test mode.
func NewForTest() *State {
	rootPath, ok := os.LookupEnv("GG_ROOT")
	if !ok {
		log.Fatal("GG_ROOT env var not defined")
	}

	root, err := os.OpenRoot(rootPath)
	if err != nil {
		log.Fatal(err)
	}
	defer root.Close()

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

	logger := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{AddSource: true, Level: slog.LevelError},
	))

	return &State{
		ApiVersion: ApiVersion,
		DB:         db,
		Logger:     logger,
		Mode:       "test",
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

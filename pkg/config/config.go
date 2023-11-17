package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	FileName = "imperium.yaml"
	Prefix   = "imperium"
)

type Config interface {
	GetPort() int
	GetLogLevel() int8
}

type ImperiumConfig struct {
	Port      int              `mapstructure:"port"`
	LogLevel  int8             `mapstructure:"logLevel"`
	Datastore *DatastoreConfig `mapstructure:"datastore"`
	Check     *CheckConfig     `mapstructure:"check"`
}

func (config ImperiumConfig) GetPort() int {
	return config.Port
}

func (config ImperiumConfig) GetLogLevel() int8 {
	return config.LogLevel
}

func (config ImperiumConfig) GetDatastore() *DatastoreConfig {
	return config.Datastore
}

func (config ImperiumConfig) GetCheck() *CheckConfig {
	return config.Check
}

type DatastoreConfig struct {
	Postgres *PostgresConfig `mapstructure:"postgres"`
	SQLite   *SQLiteConfig   `mapstructure:"sqlite"`
}

type PostgresConfig struct {
	Username                 string        `mapstructure:"username"`
	Password                 string        `mapstructure:"password"`
	Hostname                 string        `mapstructure:"hostname"`
	Database                 string        `mapstructure:"database"`
	SSLMode                  string        `mapstructure:"sslmode"`
	MigrationSource          string        `mapstructure:"migrationSource"`
	MaxIdleConnections       int           `mapstructure:"maxIdleConnections"`
	ConnMaxIdleTime          time.Duration `mapstructure:"connMaxIdleTime"`
	MaxOpenConnections       int           `mapstructure:"maxOpenConnections"`
	ConnMaxLifetime          time.Duration `mapstructure:"connMaxLifetime"`
	ReaderHostname           string        `mapstructure:"readerHostname"`
	ReaderMaxIdleConnections int           `mapstructure:"readerMaxIdleConnections"`
	ReaderMaxOpenConnections int           `mapstructure:"readerMaxOpenConnections"`
}

type SQLiteConfig struct {
	Database           string        `mapstructure:"database"`
	InMemory           bool          `mapstructure:"inMemory"`
	MigrationSource    string        `mapstructure:"migrationSource"`
	MaxIdleConnections int           `mapstructure:"maxIdleConnections"`
	ConnMaxIdleTime    time.Duration `mapstructure:"connMaxIdleTime"`
	MaxOpenConnections int           `mapstructure:"maxOpenConnections"`
	ConnMaxLifetime    time.Duration `mapstructure:"connMaxLifetime"`
}

type CheckConfig struct {
	Concurrency    int           `mapstructure:"concurrency"`
	MaxConcurrency int           `mapstructure:"maxConcurrency"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

func NewConfig() ImperiumConfig {
	viper.SetConfigFile(FileName)

	viper.SetDefault("port", 8000)
	viper.SetDefault("logLevel", zerolog.DebugLevel)

	viper.SetDefault("datastore.postgres.connMaxIdleTime", 4*time.Hour)
	viper.SetDefault("datastore.postgres.connMaxLifetime", 6*time.Hour)

	viper.SetDefault("datastore.sqlite.connMaxIdleTime", 4*time.Hour)
	viper.SetDefault("datastore.sqlite.connMaxLifetime", 6*time.Hour)

	viper.SetDefault("check.concurrency", 4)
	viper.SetDefault("check.maxConcurrency", 1000)
	viper.SetDefault("check.timeout", 1*time.Minute)

	_, err := os.ReadFile(FileName)
	if err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Msg("init: error while reading imperium.yaml. Shutting down.")
		}
	} else {
		if os.IsNotExist(err) {
			log.Info().Msg("init: could not find imperium.yaml. Attempting to use environment variables.")
		} else {
			log.Fatal().Err(err).Msg("init: error while reading imperium.yaml. Shutting down.")
		}
	}

	var config ImperiumConfig
	for _, fieldName := range getFlattenedStructFields(reflect.TypeOf(config)) {
		envKey := strings.ToUpper(fmt.Sprintf("%s_%s", Prefix, strings.ReplaceAll(fieldName, ".", "-")))
		envVar := os.Getenv(envKey)
		if envVar != "" {
			viper.Set(fieldName, envVar)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("init: error while creating config. Shutting down.")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(config.GetLogLevel()))

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return config
}

func getFlattenedStructFields(t reflect.Type) []string {
	return getFlattenedStructFieldsHelper(t, []string{})
}

func getFlattenedStructFieldsHelper(t reflect.Type, prefixes []string) []string {
	unwrappedT := t
	if t.Kind() == reflect.Pointer {
		unwrappedT = t.Elem()
	}

	flattenedFields := make([]string, 0)
	for i := 0; i < unwrappedT.NumField(); i++ {
		field := unwrappedT.Field(i)
		fieldName := field.Tag.Get("mapstructure")
		switch field.Type.Kind() {
		case reflect.Struct, reflect.Pointer:
			flattenedFields = append(flattenedFields, getFlattenedStructFieldsHelper(field.Type, append(prefixes, fieldName))...)
		default:
			flattenedField := fieldName
			if len(prefixes) > 0 {
				flattenedField = fmt.Sprintf("%s.%s", strings.Join(prefixes, "."), fieldName)
			}
			flattenedFields = append(flattenedFields, flattenedField)
		}
	}
	return flattenedFields
}

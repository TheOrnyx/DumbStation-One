/*
 * Used for the logging in my emulator
 * Was easier to make like a package than try and switch between etc
 */
package log

import (
	"os"

	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
)

var plog zerolog.Logger


func init() {
	plog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}


// Info logs an informational message
func Info(msg string) {
	plog.Info().Msg(msg)
}

// Infof logs a formatted informational message
func Infof(format string, args ...interface{}) {
	plog.Info().Msgf(format, args...)
}

// Log logs a general log message (equivalent to Info)
func Log(msg string) {
	plog.Info().Msg(msg)
}

// Logf logs a formatted general log message
func Logf(format string, args ...interface{}) {
	plog.Info().Msgf(format, args...)
}

// Warn logs a warning message
func Warn(msg string) {
	plog.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	plog.Warn().Msgf(format, args...)
}

// Error logs an error message
func Error(msg string) {
	plog.Error().Msg(msg)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	plog.Error().Msgf(format, args...)
}

// Fatal logs a fatal error message and exits the application
func Fatal(msg string) {
	plog.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal error message and exits the application
func Fatalf(format string, args ...interface{}) {
	plog.Fatal().Msgf(format, args...)
}

// Panic logs a panic message and panics
func Panic(msg string) {
	plog.Panic().Msg(msg)
}

// Panicf logs a formatted panic message and panics
func Panicf(format string, args ...interface{}) {
	plog.Panic().Msgf(format, args...)
}

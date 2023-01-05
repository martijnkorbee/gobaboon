// Logger package abstracts the creation of service oriented loggers with log rotating build in. We are using zerologger and lumberjack.
// Refer to https://zerolog.io/ and https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2 for full documentation.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

type LoggerConfig struct {
	// Rootpath is the root path of the application.
	// We need it here to create the correct file path which is a full path to the log file.
	Rootpath string

	// Debug sets the mode for the loggers. With debug mode on logs contain extra info.
	Debug bool

	// Console sets if the logger should log to console
	// default is false
	Console bool

	// ToFile sets if the logger should log to file
	// default is false
	ToFile bool

	// Service sets the service name for the logger. Default service name is default.
	Service string

	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory. Should be the full path. Defaults to /logs/log.log.
	Filename string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool
}

// Starts creates and starts the logger and returns a pointer to the Logger.
func (l *LoggerConfig) Start() *Logger {
	if l.Filename == "" {
		l.Filename = "/logs/log.log"
	}
	if l.Service == "" {
		l.Service = "default"
	}

	var log zerolog.Logger
	var writers []io.Writer

	// log default settings
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log = log.Level(zerolog.InfoLevel).With().Str("service", l.Service).Timestamp().Logger()

	// log debug settings
	if l.Debug {
		log = log.Level(zerolog.DebugLevel).With().Caller().Logger()
	}

	// create writers
	if l.Console {
		cw := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		writers = append(writers, cw)
	}
	if l.ToFile {
		fw := lumberjack.Logger{
			Filename:   l.Rootpath + l.Filename,
			MaxSize:    l.MaxSize,
			MaxAge:     l.MaxAge,
			MaxBackups: l.MaxBackups,
			LocalTime:  l.LocalTime,
			Compress:   l.Compress,
		}
		writers = append(writers, &fw)
	}
	mw := zerolog.MultiLevelWriter(writers...)

	// set writers to log output
	log = log.Output(mw)

	return &Logger{
		&log,
	}

}

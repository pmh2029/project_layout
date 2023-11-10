package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

type CustomLogger struct{}

func (logger *CustomLogger) Format(
	entry *logrus.Entry,
) ([]byte, error) {
	var message string
	dataMap := make(map[string]interface{})
	for i, data := range entry.Data {
		if data == nil {
			continue
		}
		dataMap[i] = data
	}

	if len(dataMap) > 0 {
		message = entry.Message
		var keys []string
		for k := range dataMap {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			v := dataMap[k]
			message = color.Magenta.Sprintf("[%s:%v]", k, v) + " " + message
		}
	} else {
		message = entry.Message
	}

	// message = entry.Message + message

	var buff *bytes.Buffer
	if entry.Buffer != nil {
		buff = entry.Buffer
	} else {
		buff = &bytes.Buffer{}
	}

	var logLevel string
	switch entry.Level.String() {
	case "info":
		logLevel = color.Info.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	case "warn":
		logLevel = color.Warn.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	case "error":
		logLevel = color.Red.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	case "debug":
		logLevel = color.Debug.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	case "panic":
		logLevel = color.Error.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	case "fatal":
		logLevel = color.Error.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	default:
		logLevel = color.Info.Sprint("[" + strings.ToUpper(entry.Level.String()) + "]")
	}

	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	var newLog string

	if entry.HasCaller() {
		fileName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("%s %s %s%s %s\n",
			color.LightBlue.Sprint("["+timestamp+"]"),
			logLevel,
			color.LightYellow.Sprint("["+fileName+":"),
			color.LightYellow.Sprint(strconv.Itoa(entry.Caller.Line)+"]"),
			message,
		)
	} else {
		newLog = fmt.Sprintf("%s %s %s\n",
			color.Cyan.Sprint("["+timestamp+"]"),
			color.Info.Sprint("["+strings.ToUpper(entry.Level.String())+"]"),
			entry.Message,
		)
	}

	buff.WriteString(newLog)
	return buff.Bytes(), nil
}

func NewCustomLogger() *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&CustomLogger{})
	log.SetLevel(logrus.InfoLevel)

	return log
}

// Logger based on logrus, but compatible with gorm
type GormLogger struct {
	logger *logrus.Entry
}

func NewGormLogger(
	logger *logrus.Logger,
) GormLogger {
	return GormLogger{
		logger.WithField("service", "database"),
	}
}

// We ignore this setting, because the log level is already decided by logrus
func (logger GormLogger) LogMode(gormLog.LogLevel) gormLog.Interface {
	return logger
}

func (logger GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logger.logger.WithContext(ctx).Infof(msg, args...)
}

func (logger GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logger.logger.WithContext(ctx).Warnf(msg, args...)
}

func (logger GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logger.logger.WithContext(ctx).Errorf(msg, args...)
}

// We want the SQL logs with the info level, while it's defined as trace by gorm
func (logger GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	sql, rows := fc()
	duration := time.Since(begin)
	logEntry := logger.logger.
		WithContext(ctx).
		WithField("duration", duration.String()).
		WithField("rows", rows).
		WithField("query", sql)

	if err == nil {
		logEntry.Info("Performed SQL Query")
	} else {
		logEntry = logEntry.WithField("error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logEntry.Info("Performed SQL Query")
		} else {
			logEntry.Error("SQL Query failed")
		}
	}
}

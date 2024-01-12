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
	"gorm.io/gorm/utils"
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
			if k != "" {
				message = color.LightCyan.Sprintf("[%s:%v]", k, v) + " " + message
			} else {
				message = color.LightCyan.Sprintf("[%v]", v) + " " + message
			}
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
			color.LightBlue.Sprint("["+timestamp+"]"),
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
	Logger                    *logrus.Entry
	LogLevel                  gormLog.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	FileWithLineNumField      string
}

func NewGormLogger(
	opts GormLogger,
) *GormLogger {
	if opts.Logger == nil {
		opts.Logger = logrus.NewEntry(logrus.New())
	}

	if opts.LogLevel == 0 {
		opts.LogLevel = gormLog.Silent
	}
	return &opts
}

func (l *GormLogger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormLog.Info {
		l.Logger.WithContext(ctx).Infof(s, args...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormLog.Warn {
		l.Logger.WithContext(ctx).Warnf(s, args...)
	}
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormLog.Error {
		l.Logger.WithContext(ctx).Errorf(s, args...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLog.Silent {
		return
	}

	fields := logrus.Fields{}
	fields[l.FileWithLineNumField] = filepath.Base(utils.FileWithLineNum())

	sql, rows := fc()
	if rows == -1 {
		fields["rows"] = "-"
	} else {
		fields["rows"] = rows
	}

	elapsed := time.Since(begin)
	fields["durations"] = elapsed

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) && l.LogLevel >= gormLog.Error:
		l.Logger.WithContext(ctx).WithFields(fields).Errorf("[QUERY:%s] [ERROR:%v]", sql, err)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLog.Warn:
		l.Logger.WithContext(ctx).WithFields(fields).Warnf("[SLOW SQL >= %v], [QUERY: %s]", l.SlowThreshold, sql)
	case l.LogLevel == gormLog.Info:
		l.Logger.WithContext(ctx).WithFields(fields).Infof("[QUERY:%s]", sql)
	}
}

package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

var (
	Log *Zerolog
)

const (
	LogFile = "ttsBot.log"
)

type Logger interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Debugf(string, ...interface{})
	Debug(...interface{})
}

type Zerolog struct {
	Logger zerolog.Logger
}

func NewZerolog() *Zerolog {
	runLogFile, _ := os.OpenFile(
		LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)

	return &Zerolog{
		Logger: zerolog.New(multi).With().Timestamp().Logger(),
	}
}

func (z *Zerolog) Infof(format string, args ...interface{}) {
	z.Logger.Info().Msgf(format, args...)
}

func (z *Zerolog) Info(args ...interface{}) {
	z.Logger.Info().Msg(fmt.Sprint(args...))
}

func (z *Zerolog) Errorf(format string, args ...interface{}) {
	z.Logger.Error().Msgf(format, args...)
}

func (z *Zerolog) Error(args ...interface{}) {
	z.Logger.Error().Msg(fmt.Sprint(args...))
}

func (z *Zerolog) Warnf(format string, args ...interface{}) {
	z.Logger.Warn().Msgf(format, args...)
}

func (z *Zerolog) Warn(args ...interface{}) {
	z.Logger.Warn().Msg(fmt.Sprint(args...))
}

func (z *Zerolog) Debugf(format string, args ...interface{}) {
	z.Logger.Debug().Msgf(format, args...)
}

func (z *Zerolog) Debug(args ...interface{}) {
	z.Logger.Debug().Msg(fmt.Sprint(args...))
}

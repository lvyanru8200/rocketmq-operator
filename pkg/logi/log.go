package logi

import (
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func init() {
	SetLogger(Build(logConfig()))
}

func logConfig() uzap.Config {
	var uzapConfig uzap.Config
	level := logLevelFromEnv()
	switch level {
	case uzap.DebugLevel:
		uzapConfig = uzap.NewDevelopmentConfig()
		uzapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		uzapConfig = uzap.NewProductionConfig()
	}

	uzapConfig.Level = uzap.NewAtomicLevelAt(level)

	return uzapConfig
}

func SetLogger(l *uzap.Logger) {
	uzap.ReplaceGlobals(l)
}

func GetLogger() *uzap.Logger {
	return uzap.L()
}

func GetSugaredLogger() *uzap.SugaredLogger {
	return uzap.S()
}

func Build(uzapConfig uzap.Config) *uzap.Logger {
	log, err := uzapConfig.Build()
	if err != nil {
		panic(err)
	}
	return log
}

func logLevelFromEnv() zapcore.Level {
	levelStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return uzap.DebugLevel
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		panic(err)
	}
	return level
}

func IsDebug() bool {
	return uzap.L().Core().Enabled(uzap.DebugLevel)
}

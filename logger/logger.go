package logger

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type loggerCategory string

var (
	LoggerApp      loggerCategory = "app"
	LoggerRouter   loggerCategory = "router"
	LoggerUsecase  loggerCategory = "usecase"
	LoggerExternal loggerCategory = "external"
)

type LoggerData struct {
	Category       loggerCategory         `json:"category,omitempty"`
	Service        string                 `json:"service,omitempty"`
	Module         string                 `json:"module,omitempty"`
	Method         string                 `json:"method,omitempty"`
	Duration       int64                  `json:"duration,omitempty"`
	Version        string                 `json:"version,omitempty"`
	External       *ExternalLoggerData    `json:"external,omitempty"`
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
}

type ExternalLoggerData struct {
	Service        string                 `json:"service,omitempty"`
	Module         string                 `json:"module,omitempty"`
	Method         string                 `json:"method,omitempty"`
	Version        string                 `json:"version,omitempty"`
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
}

type Logger struct {
	logger *logrus.Logger
}

func Init() *Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Info(data *LoggerData, message string) {
	dataMap, err := buildLogData(data)
	if err != nil {
		panic(err)
	}

	l.logger.WithFields(dataMap).Info(message)
}

func (l *Logger) Error(data *LoggerData, e string) {
	dataMap, err := buildLogData(data)
	if err != nil {
		panic(err)
	}

	l.logger.WithFields(dataMap).Error(e)
}

func buildLogData(data *LoggerData) (map[string]interface{}, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var dataMap map[string]interface{}

	err = json.Unmarshal(dataStr, &dataMap)
	if err != nil {
		return nil, err
	}

	return dataMap, nil
}

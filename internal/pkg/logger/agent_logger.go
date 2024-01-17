package logger

import "go.uber.org/zap"

var ALog zap.SugaredLogger

func InitALogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()
	ALog = *logger.Sugar()
	return nil
}

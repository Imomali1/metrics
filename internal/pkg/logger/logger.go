package logger

import "go.uber.org/zap"

var Log zap.SugaredLogger

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()
	Log = *logger.Sugar()
	return nil
}

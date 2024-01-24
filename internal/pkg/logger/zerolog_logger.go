package logger

//type ILogger interface {
//	Trace() IEvent
//	Debug() IEvent
//	Info() IEvent
//	Warn() IEvent
//	Error() IEvent
//	Err(err error) IEvent
//	Fatal() IEvent
//	Panic() IEvent
//}
//
//type loggerImplementation struct {
//	zero *zerolog.Logger
//}
//
//func (l *loggerImplementation) Trace() IEvent {
//	return &event{zlEvent: l.zero.Trace()}
//}
//
//func (l *loggerImplementation) Debug() IEvent {
//	return &event{zlEvent: l.zero.Debug()}
//}
//
//func (l *loggerImplementation) Info() IEvent {
//	return &event{zlEvent: l.zero.Info()}
//}
//
//func (l *loggerImplementation) Warn() IEvent {
//	return &event{zlEvent: l.zero.Warn()}
//}
//
//func (l *loggerImplementation) Error() IEvent {
//	return &event{zlEvent: l.zero.Error()}
//}
//
//func (l *loggerImplementation) Err(err error) IEvent {
//	return &event{zlEvent: l.zero.Err(err)}
//}
//
//func (l *loggerImplementation) Fatal() IEvent {
//	return &event{zlEvent: l.zero.Fatal()}
//}
//
//func (l *loggerImplementation) Panic() IEvent {
//	return &event{zlEvent: l.zero.Panic()}
//}

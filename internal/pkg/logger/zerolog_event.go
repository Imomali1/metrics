package logger

//type IEvent interface {
//	Send()
//	Msg(msg string)
//	Str(key, val string) IEvent
//	Err(err error) IEvent
//	Int(key string, i int) IEvent
//	Int64(key string, i int64) IEvent
//	Float64(key string, f float64) IEvent
//	Time(key string, t time.Time) IEvent
//	Dur(key string, d time.Duration) IEvent
//	Any(key string, i interface{}) IEvent
//}
//
//type event struct {
//	zlEvent *zerolog.Event
//}
//
//func (e *event) Send() {
//	e.zlEvent.Send()
//}
//
//func (e *event) Msg(msg string) {
//	e.zlEvent.Msg(msg)
//}
//
//func (e *event) Str(key, val string) IEvent {
//	e.zlEvent.Str(key, val)
//	return e
//}
//
//func (e *event) Err(err error) IEvent {
//	e.zlEvent.Err(err)
//	return e
//}
//
//func (e *event) Int(key string, i int) IEvent {
//	e.zlEvent.Int(key, i)
//	return e
//}
//
//func (e *event) Int64(key string, i int64) IEvent {
//	e.zlEvent.Int64(key, i)
//	return e
//}
//
//func (e *event) Float64(key string, f float64) IEvent {
//	e.zlEvent.Float64(key, f)
//	return e
//}
//
//func (e *event) Time(key string, t time.Time) IEvent {
//	e.zlEvent.Time(key, t)
//	return e
//}
//
//func (e *event) Dur(key string, d time.Duration) IEvent {
//	e.zlEvent.Dur(key, d)
//	return e
//}
//
//func (e *event) Any(key string, i interface{}) IEvent {
//	e.zlEvent.Any(key, i)
//	return e
//}

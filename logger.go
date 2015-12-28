package protolog

import (
	"fmt"
	"io"
	"os"

	"go.pedge.io/protolog/pb"

	"github.com/golang/protobuf/proto"
)

type logger struct {
	pusher        Pusher
	enableID      bool
	idAllocator   IDAllocator
	timer         Timer
	errorHandler  ErrorHandler
	level         Level
	contexts      []proto.Message
	genericFields *protologpb.Fields
}

func newLogger(pusher Pusher, options ...LoggerOption) *logger {
	logger := &logger{
		pusher,
		false,
		DefaultIDAllocator,
		DefaultTimer,
		DefaultErrorHandler,
		DefaultLevel,
		make([]proto.Message, 0),
		&protologpb.Fields{
			Value: make(map[string]string, 0),
		},
	}
	for _, option := range options {
		option(logger)
	}
	return logger
}

func (l *logger) Flush() error {
	return l.pusher.Flush()
}

func (l *logger) AtLevel(level Level) Logger {
	return &logger{
		l.pusher,
		l.enableID,
		l.idAllocator,
		l.timer,
		l.errorHandler,
		level,
		l.contexts,
		l.genericFields,
	}
}

func (l *logger) WithContext(context proto.Message) Logger {
	return &logger{
		l.pusher,
		l.enableID,
		l.idAllocator,
		l.timer,
		l.errorHandler,
		l.level,
		append(l.contexts, context),
		l.genericFields,
	}
}

func (l *logger) Debug(event proto.Message) {
	l.print(LevelDebug, event)
}

func (l *logger) Info(event proto.Message) {
	l.print(LevelInfo, event)
}

func (l *logger) Warn(event proto.Message) {
	l.print(LevelWarn, event)
}

func (l *logger) Error(event proto.Message) {
	l.print(LevelError, event)
}

func (l *logger) Fatal(event proto.Message) {
	l.print(LevelFatal, event)
	os.Exit(1)
}

func (l *logger) Panic(event proto.Message) {
	l.print(LevelPanic, event)
	panic(fmt.Sprintf("%+v", event))
}

func (l *logger) Print(event proto.Message) {
	l.print(LevelNone, event)
}

func (l *logger) DebugWriter() io.Writer {
	return l.printWriter(LevelDebug)
}

func (l *logger) InfoWriter() io.Writer {
	return l.printWriter(LevelInfo)
}

func (l *logger) WarnWriter() io.Writer {
	return l.printWriter(LevelWarn)
}

func (l *logger) ErrorWriter() io.Writer {
	return l.printWriter(LevelError)
}

func (l *logger) Writer() io.Writer {
	return l.printWriter(LevelNone)
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	contextFields := make(map[string]string, len(fields))
	for key, value := range fields {
		contextFields[key] = fmt.Sprintf("%v", value)
	}
	for key, value := range l.genericFields.Value {
		contextFields[key] = value
	}
	return &logger{
		l.pusher,
		l.enableID,
		l.idAllocator,
		l.timer,
		l.errorHandler,
		l.level,
		l.contexts,
		&protologpb.Fields{
			Value: contextFields,
		},
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Debugln(args ...interface{}) {
	l.Debug(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Infoln(args ...interface{}) {
	l.Info(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Warnln(args ...interface{}) {
	l.Warn(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Errorln(args ...interface{}) {
	l.Error(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Fatalln(args ...interface{}) {
	l.Fatal(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.Panic(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Panicln(args ...interface{}) {
	l.Panic(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.Print(&protologpb.Event{Message: fmt.Sprintf(format, args...)})
}

func (l *logger) Println(args ...interface{}) {
	l.Print(&protologpb.Event{Message: fmt.Sprint(args...)})
}

func (l *logger) print(level Level, event proto.Message) {
	if err := l.printWithError(level, event); err != nil {
		l.errorHandler.Handle(err)
	}
}

func (l *logger) printWriter(level Level) io.Writer {
	// TODO(pedge): think more about this
	//if !l.isLoggedLevel(level) {
	//return ioutil.Discard
	//}
	return newLogWriter(l, level)
}

func (l *logger) printWithError(level Level, event proto.Message) error {
	if !l.isLoggedLevel(level) {
		return nil
	}
	// TODO(pedge): should copy this but has performance hit
	contexts := l.contexts
	if len(l.genericFields.Value) > 0 {
		contexts = append(contexts, l.genericFields)
	}
	goEntry := &GoEntry{
		Level:    level,
		Time:     l.timer.Now(),
		Contexts: contexts,
		Event:    event,
	}
	if l.enableID {
		goEntry.ID = l.idAllocator.Allocate()
	}
	return l.pusher.Push(goEntry)
}

func (l *logger) isLoggedLevel(level Level) bool {
	return level >= l.level || level == LevelNone
}

type logWriter struct {
	logger *logger
	level  Level
}

func newLogWriter(logger *logger, level Level) *logWriter {
	return &logWriter{logger, level}
}

func (w *logWriter) Write(p []byte) (int, error) {
	if err := w.logger.printWithError(w.level, &protologpb.WriterOutput{Value: p}); err != nil {
		return 0, err
	}
	return len(p), nil
}

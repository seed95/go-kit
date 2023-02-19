package log

import (
	"fmt"
	"github.com/seed95/go-kit/log/keyval"
	kitzap "github.com/seed95/go-kit/log/zap"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TODO change name this struct
type bundle struct {
	kitzap.Logger
	enable bool
}

var _ kitzap.Logger = (*bundle)(nil)

var (
	// log is a bundle instance to log globally in project
	log bundle

	// logPackageName is used to cache package name for bundle.
	// Cached at first use (initPackageName).
	logPackageName string

	// projectPath is used to cache root of project path.
	// Cached at first use (initPackageName).
	projectPath string

	// Positions in the call stack when tracing to report the calling method.
	//
	// Note: copy from https://github.com/sirupsen/logrus/blob/master/entry.go
	minimumCallerDepth int
)

const (
	// maximumCallerDepth used for runtime.Callers
	// runtime.Callers comment: `return program counters of function invocations on the calling goroutine's stack.`
	maximumCallerDepth int = 25

	// relativeLogPath is used in initPackageName to make projectPath
	//
	// Note: log package must be in located in the following folder relatively to the root of project,
	// Otherwise, the following value should be changed.
	relativeLogPath = "/pkg/log"
)

func init() {
	//log = bundle{}
	initPackageName()

	// Create default std bundle
	stdCore, err := kitzap.NewStandardCore(true, kitzap.ErrorLevel)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create std log instance, error: %v", err))
		return
	}
	log.Logger = kitzap.NewZapLoggerWithCores(stdCore)
	log.enable = false
}

// initPackageName Get the package name for log and root project file name,
// and set the minimum caller depth dynamically
func initPackageName() {
	// Get project path
	_, b, _, _ := runtime.Caller(0)
	projectPath = strings.Replace(filepath.Dir(b), relativeLogPath, "", -1) + "/" // e.g.: /home/sajad/go/src/ap/marketing-gateway/

	// Get package name for log
	pcs := make([]uintptr, maximumCallerDepth)
	_ = runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs)
	var i int
	for f, again := frames.Next(); again; f, again = frames.Next() {
		i++
		if strings.Contains(f.Function, "log.initPackageName") {
			logPackageName = getPackageName(f.Function)
			break
		}
	}

	// Set minimum caller depth
	minimumCallerDepth = i + 2 // +1 for ReqRes method in this package
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be a better way...
//
// Note: copy from https://github.com/sirupsen/logrus/blob/master/entry.go
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func SetEnable(enable bool) {
	log.enable = enable
}

// SetLogger assign new kitlogger.Logger to used for global log instance
func SetLogger(logger kitzap.Logger) {
	log.Logger = logger
}

func GetLogger() kitzap.Logger {
	return log
}

func Debug(message string, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	go log.Debug(message, keyVal...)
}

func Info(message string, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	go log.Info(message, keyVal...)
}

func Warn(message string, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	go log.Warn(message, keyVal...)
}

func Error(message string, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	go log.Error(message, keyVal...)
}

func Panic(message string, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	go log.Panic(message, keyVal...)
}

// ReqRes is used to log request response in all functions
// that used defer structure in domain or handler or service.
func ReqRes(startTime time.Time, err error, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	log.ReqRes(startTime, err, keyVal...)
}

func (l *bundle) ReqRes(startTime time.Time, err error, keyVal ...keyval.Pair) {
	if !log.enable {
		return
	}

	message := ""
	// Get caller function name for message
	caller := getCaller()
	if caller != nil {
		message = strings.Replace(caller.File, projectPath, "", 1)
		message = strings.Replace(message, ".go", "", 1) + getFuncName(caller.Function)
	}

	// Append duration to keyVal
	keyVal = append(keyVal, keyval.String("duration", time.Since(startTime).String()))

	if err != nil {
		keyVal = append(keyVal, keyval.Error(err))
		l.Error(message, keyVal...)
	} else {
		l.Info(message, keyVal...)
	}
}

func getFuncName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".func")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]

			lastPeriod = strings.LastIndex(f, ".")
			f = f[lastPeriod:]

		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-log calling function
//
// Note: copy from https://github.com/sirupsen/logrus/blob/master/entry.go
func getCaller() *runtime.Frame {
	// Restrict the lookups frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		if strings.Contains(f.Function, ".func") {
			pkg := getPackageName(f.Function)
			// If the caller isn't part of this package, we're done
			if pkg != logPackageName {
				return &f //nolint:scopelint
			}
		}

	}

	// if we got here, we failed to find the caller's context
	return nil
}

// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"strings"

	"fmt"
	"github.com/Sirupsen/logrus"
	//"github.com/astaxie/beego/logs"
	"os"
	"runtime"
	"sync"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var (
	goroLableMap = gorolable{gomap: make(map[int64]interface{})}
)

type gorolable struct {
	sync.Mutex
	gomap map[int64]interface{}
}

func RegisterLable(lable interface{}) {
	goroLableMap.Lock()
	defer goroLableMap.Unlock()
	goroLableMap.gomap[runtime.GoID()] = lable
}

func UnregisterLable() {
	goroLableMap.Lock()
	defer goroLableMap.Unlock()
	delete(goroLableMap.gomap, runtime.GoID())
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l string) {
	lvl, err := logrus.ParseLevel(l)
	if err != nil {
		return
	}
	logrus.SetLevel(lvl)
	//BeeLogger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	//BeeLogger.EnableFuncCallDepth(b)
	//BeeLogger.SetLogFuncCallDepth(3)
}

// logger references the used application logger.
//var BeeLogger *logs.BeeLogger

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	//err := BeeLogger.SetLogger(adaptername, config)
	f, err := os.OpenFile(adaptername, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	logrus.SetOutput(f)
	return nil
}

func Emergency(v ...interface{}) {
	entry().Fatalln(v)
	//BeeLogger.Emergency(generateFmtStr(len(v)), v...)
}

func Alert(v ...interface{}) {
	entry().Warnln(v)
	//BeeLogger.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	entry().Fatalln(v)
	//BeeLogger.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	entry().Errorln(v)
	//BeeLogger.Error(generateFmtStr(len(v)), v...)
}

func Errorf(format string, v ...interface{}) {
	entry().Errorf(format, v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	entry().Warning(v)
	//BeeLogger.Warning(generateFmtStr(len(v)), v...)
}

// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Warn(v ...interface{}) {
	entry().Warnln(v)
	//Warning(v...)
}

func Notice(v ...interface{}) {
	entry().Println(v)
	//BeeLogger.Notice(generateFmtStr(len(v)), v...)
}

// Info logs a message at info level.
func Informational(v ...interface{}) {
	entry().Infoln(v)
	//BeeLogger.Informational(generateFmtStr(len(v)), v...)
}

func Infof(format string, v ...interface{}) {
	entry().Infof(format, v...)
}

// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Info(v ...interface{}) {
	entry().Infoln(v)
	//Informational(v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	entry().Debugln(v)
	//BeeLogger.Debug(generateFmtStr(len(v)), v...)
}

func Debugf(format string, v ...interface{}) {
	entry().Debugf(format, v...)
}

// Trace logs a message at trace level.
// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Trace(v ...interface{}) {
	entry().Infoln(v)
	//BeeLogger.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	v, ok := goroLableMap.gomap[runtime.GoID()]
	if ok == true {
		return fmt.Sprintf("%-12s", fmt.Sprintf("[%-8v]", v)) + strings.Repeat("%v ", n)
	} else {
		return strings.Repeat("%v ", n)
	}

}

func entry() *logrus.Entry {
	v, ok := goroLableMap.gomap[runtime.GoID()]
	if ok == true {
		return logrus.WithField("ID", v)
	}
	return logrus.NewEntry(logrus.StandardLogger())

}

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00"})
	//logrus.
	//BeeLogger = logs.NewLogger(10000)
	//SetLogger("console", "")
}

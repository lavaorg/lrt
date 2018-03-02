/*
Copyright (C) 2017 Verizon. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mlog

import (
	"bytes"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	golog "log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	convey.Convey("Test Log Framework", t, func() {

		// working variables
		pid := strconv.Itoa(os.Getpid())
		testOne := "test one"
		testTwo := "test two"
		alarmStarts := "*|1|ALARM|0|" + pid + "|tgrp|tapp|mlog.test|0|mlog_test.go"
		errorStarts := "*|1|ERROR|0|" + pid + "|tgrp|tapp|mlog.test|0|mlog_test.go"
		debugStarts := "*|1|DEBUG|0|" + pid + "|tgrp|tapp|mlog.test|0|mlog_test.go"
		statStarts := "*|1|STAT|0|" + pid + "|tgrp|tapp|mlog.test|0|mlog_test.go"
		infoStarts := "*|1|INFO|0|" + pid + "|tgrp|tapp|mlog.test|0|mlog_test.go"
		eventStarts := "*|1|EVENT|0|" + pid + "|tgrp|tapp|mlog.test|0|mlog_test.go"

		oneEnds := "|" + testOne + "\n"
		twoEnds := "|" + testTwo + "\n"

		golangInfoStarts := "*|1|INFO|0|" + pid + "|tgrp|tapp|mlog.test|0|"
		golangStatStarts := "*|1|STAT|0|" + pid + "|tgrp|tapp|mlog.test|0|print.go"
		golangPanicInfoStarts := "*|1|ALARM|0|" + pid + "|tgrp|tapp|mlog.test|0|"
		golangOneEnds := " " + testOne + "\n"
		golangTwoEnds := " " + testTwo + "\n"

		logMessage := ""

		convey.Convey("Check Default Init", func() {
			// check default init
			convey.So(severity, convey.ShouldEqual, INFO)
		})

		convey.Convey("Check Constants", func() {
			convey.So(NOP, convey.ShouldEqual, 0)
			convey.So(ALARM, convey.ShouldEqual, 1)
			convey.So(ERROR, convey.ShouldEqual, 2)
			convey.So(STAT, convey.ShouldEqual, 3)
			convey.So(INFO, convey.ShouldEqual, 4)
			convey.So(EVENT, convey.ShouldEqual, 5)
			convey.So(DEBUG, convey.ShouldEqual, 6)
			convey.So(UNKNOWN, convey.ShouldEqual, 7)
			convey.So(len(severityToStream), convey.ShouldEqual, int(UNKNOWN)+1)
			convey.So(len(SeverityToString), convey.ShouldEqual, int(UNKNOWN)+1)
		})

		// prep env values
		os.Setenv("LRT_APP", "tapp")
		os.Setenv("LRT_GROUP", "tgrp")
		initialize()

		// create mock io.Writers (beside Stdout and Stderr)
		output1k := make([]byte, 1024)
		buffer1k := bytes.NewBuffer(output1k)
		severityToStream[0] = nil
		for i := 1; i < len(severityToStream); i++ {
			severityToStream[i] = buffer1k
		}
		buffer1k.Reset()

		convey.Convey("Check Severities with various enabled and disabled DEBUG", func() {

			// Test three conditions: default, disabled->enabled, enabled->disabled
			done := make(chan bool, 0)
			debugEnabled := false
			for i := 0; i < 3; i++ {
				// set log level
				buffer1k.Reset()
				EnableDebug(debugEnabled)

				// check reinitialization
				if debugEnabled {
					convey.So(severity, convey.ShouldEqual, DEBUG)
				} else {
					convey.So(severity, convey.ShouldEqual, EVENT)
				}
				debugEnabled = !debugEnabled

				// ALARM
				Alarm(testOne)
				logMessage = buffer1k.String()
				if severity >= ALARM {
					convey.So(logMessage, convey.ShouldStartWith, alarmStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// ALARM
				Alarm(testTwo)
				logMessage = buffer1k.String()
				if severity >= ALARM {
					convey.So(logMessage, convey.ShouldStartWith, alarmStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// ERROR
				Error(testOne)
				logMessage = buffer1k.String()
				if severity >= ERROR {
					convey.So(logMessage, convey.ShouldStartWith, errorStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// ERROR
				Error(testTwo)
				logMessage = buffer1k.String()
				if severity >= ERROR {
					convey.So(logMessage, convey.ShouldStartWith, errorStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// STAT
				Stat(testOne)
				logMessage = buffer1k.String()
				if severity >= STAT {
					convey.So(logMessage, convey.ShouldStartWith, statStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// STAT
				Stat(testTwo)
				logMessage = buffer1k.String()
				if severity >= STAT {
					time.Sleep(time.Millisecond)
					convey.So(logMessage, convey.ShouldStartWith, statStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// EVENT
				Event(testOne)
				logMessage = buffer1k.String()
				if severity >= EVENT {
					convey.So(logMessage, convey.ShouldStartWith, eventStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// EVENT
				Event(testTwo)
				logMessage = buffer1k.String()
				if severity >= EVENT {
					convey.So(logMessage, convey.ShouldStartWith, eventStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// INFO
				Info(testOne)
				logMessage = buffer1k.String()
				if severity >= INFO {
					convey.So(logMessage, convey.ShouldStartWith, infoStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// INFO
				Info(testTwo)
				logMessage = buffer1k.String()
				if severity >= INFO {
					convey.So(logMessage, convey.ShouldStartWith, infoStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// DEBUG
				Debug(testOne)
				logMessage = buffer1k.String()
				if severity >= DEBUG {
					convey.So(logMessage, convey.ShouldStartWith, debugStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// DEBUG
				Debug(testTwo)
				logMessage = buffer1k.String()
				if severity >= DEBUG {
					convey.So(logMessage, convey.ShouldStartWith, debugStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// STAT
				fmt.Fprintf(GetStatWriter(), testOne)
				logMessage = buffer1k.String()
				if severity >= STAT {
					convey.So(logMessage, convey.ShouldStartWith, golangStatStarts)
					convey.So(logMessage, convey.ShouldEndWith, oneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// STAT
				fmt.Fprintf(GetStatWriter(), testTwo)
				logMessage = buffer1k.String()
				if severity >= STAT {
					convey.So(logMessage, convey.ShouldStartWith, golangStatStarts)
					convey.So(logMessage, convey.ShouldEndWith, twoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// UNKNOWN - go log
				golog.Printf(testOne)
				logMessage = buffer1k.String()
				if severity >= INFO {
					convey.So(logMessage, convey.ShouldStartWith, golangInfoStarts)
					convey.So(logMessage, convey.ShouldEndWith, golangOneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// UNKNOWN - go log
				golog.Printf(testTwo)
				logMessage = buffer1k.String()
				if severity >= INFO {
					convey.So(logMessage, convey.ShouldStartWith, golangInfoStarts)
					convey.So(logMessage, convey.ShouldEndWith, golangTwoEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}

				// UNKNOWN - go panic
				go func() {
					defer func() {
						if r := recover(); r != nil {
							done <- true
						}
					}()
					golog.Panic(testOne)
				}()

				<-done
				logMessage = buffer1k.String()
				if severity >= ALARM {
					convey.So(logMessage, convey.ShouldStartWith, golangPanicInfoStarts)
					convey.So(logMessage, convey.ShouldEndWith, golangOneEnds)
					buffer1k.Reset()
				} else {
					convey.So(logMessage, convey.ShouldEqual, "")
				}
			}
		})
	})
}

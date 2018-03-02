/* Copyright(C) 2018 Larry Rau, All rights reserved.package mlog
/* See License Below */

package mlog

import (
	"fmt"
)

/* Provide a set of helper functions to consistently format common log output */

// Helper to dump ENV vars into a log message in a consistenet manner.
func DumpENV(envvar, envval string) string {
	return fmt.Sprintf("ENV: %v = [%v]", envvar, envval)
}

/// LICENSE
/*
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

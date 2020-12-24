/*
Copyright © 2020 Yevhen Lebid ealebed@gmail.com

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

package dockerhub

import (
	"github.com/fatih/color"
)

// BG returns bold green color
var BG = color.New(color.FgGreen, color.Bold).SprintFunc()

// BW returns bold white color
var BW = color.New(color.FgWhite, color.Bold).SprintFunc()

// BY returns bold yellow color
var BY = color.New(color.FgYellow, color.Bold).SprintFunc()

// BR returns bold red color
var BR = color.New(color.FgRed, color.Bold).SprintFunc()

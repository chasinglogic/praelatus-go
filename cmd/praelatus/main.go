// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found
// in the LICENSE file.

package main

import "github.com/praelatus/praelatus/commands"

var version = "master"
var commit = "HEAD"
var date = ""

func main() {
	commands.Version = version
	commands.Commit = commit
	commands.Execute()
}

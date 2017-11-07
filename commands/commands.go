// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/repo"
	"github.com/praelatus/praelatus/repo/mongo"
)

func loadRepo() repo.Repo {
	return mongo.New(config.DBURL())
}

func loadCache() repo.Cache {
	return mongo.NewCache(config.DBURL())
}

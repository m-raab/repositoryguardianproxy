/*
 * Copyright (c) 2019.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

var users = []User{
	{UserID: "admin", UserPassword: "admin123", UserName: "Admin User", Access: "all"},
	{UserID: "1234", UserPassword: "passwd1234", UserName: "Test User 1234", Access: "all"},
	{UserID: "1235", UserPassword: "passwd1235", UserName: "Test User 1235", Access: "all"},
	{UserID: "2367", UserPassword: "passwd2367", UserName: "Test User 2367", Access: "all"},
}

func (econf *ElasticConfig) insertUsers() {
	for _, user := range users {
		econf.insertUser(user)
	}
}

func (econf *ElasticConfig) insertUser(user User) {
	logger := log.WithFields(log.Fields{
		"service": "go-adduser-basic",
	})

	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()

	econf.ElasticClient.Index().
		Index(econf.IndexConf.userdb).
		BodyJson(user). // pass struct instance to BodyJson
		Do(ctx)         // Initiate API call with context object

	// Check for errors in each iteration
	logger.Info("Elasticsearch document indexed:", user)
}

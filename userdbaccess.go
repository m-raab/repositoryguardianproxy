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
	"encoding/json"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"time"
)

type User struct {
	UserID       string `json:"userid"`
	UserPassword string `json:"userpasswd"`
	UserName     string `json:"username"`
	Access       string `json:"access"`
}

func (econf *ElasticConfig) findUser(userid string) *User {
	logger := log.WithFields(log.Fields{
		"service": "go-user-access",
	})

	ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	var users []User

	var resultUser = &User{
		UserID:       "UNKNOWN",
		UserPassword: "",
		UserName:     "",
		Access:       "",
	}

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("userid", userid))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		logger.Fatal(err1, err2)
	}
	logger.Info(string(queryJs))

	searchService := econf.ElasticClient.Search().Index(econf.IndexConf.userdb).SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		var user User
		err = json.Unmarshal(hit.Source, &user)
		if err != nil {
			logger.Fatal(err)
		}
		users = append(users, user)
	}

	if err == nil {
		if len(users) > 1 {
			logger.Fatal("Only one user is allowed. Check user ", userid)
		}
		if len(users) < 1 {
			logger.Fatal("No user found. Check user ", userid)
		}
		resultUser.UserID = users[0].UserID
		resultUser.UserName = users[0].UserName
		resultUser.UserPassword = users[0].UserPassword
		resultUser.Access = users[0].Access
	}

	return resultUser
}

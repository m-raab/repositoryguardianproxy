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
	"time"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

//ElasticSetup should setup elastic search
func (econf *ElasticConfig) elasticSetup() {
	var ctx = context.Background()
	var connError error

	logger := log.WithFields(log.Fields{

		"service": "go-elastic",
	})

	errlogger := log.WithFields(log.Fields{
		"service": "go-elastic-info",
	})

	infologger := log.WithFields(log.Fields{
		"service": "go-elastic-error",
	})

	logger.Info("Connecting to elastic search on ", econf.Host)
	//create an elastic search client. connect to the running elastic search db
	econf.ElasticClient, connError = elastic.NewSimpleClient(
		elastic.SetURL(econf.Host),
		elastic.SetSniff(true),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(errlogger),
		elastic.SetInfoLog(infologger),
		elastic.SetBasicAuth(econf.User, econf.Password),
	)
	if connError != nil {
		panic(connError)
	}

	info, code, err := econf.ElasticClient.Ping(econf.Host).Do(ctx)

	if err != nil {
		logger.Fatal(err, info, code)
	}
}

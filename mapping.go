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

var requestmapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"remotehost":{
				"type":"keyword"
			},
			"userid":{
				"type":"keyword"
			},
			"artifactpath":{
				"type":"keyword"
			},
			"requested":{
				"type":"date"
			}
		}
	}
}`

var useridmapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"userid":{
				"type":"keyword"
			},
			"userpasswd":{
				"type":"keyword"
			},
			"username":{
				"type":"keyword"
			},
			"access":{
				"type":"text"
			}
		}
	}
}`

func (econf *ElasticConfig) checkIndeces() {
	checkIndex(econf, econf.IndexConf.userdb, useridmapping)
	checkIndex(econf, econf.IndexConf.requestdb, requestmapping)
}

func checkIndex(econf *ElasticConfig, indexname string, mapping string) {
	logger := log.WithFields(log.Fields{
		"service": "go-elastic-index", "index": indexname,
	})

	ctx, stop := context.WithTimeout(context.Background(), 30*time.Second)
	defer stop()

	// Use the IndexExists service to check if a specified index exists.
	exists, err := econf.ElasticClient.IndexExists(indexname).Do(ctx)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
	//if index does not exist, create a new one with the specified mapping
	if !exists {
		createIndex, err := econf.ElasticClient.CreateIndex(indexname).BodyString(mapping).Do(ctx)
		if err != nil {
			logger.Fatal(err)
			panic(err)
		}
		if !createIndex.Acknowledged {
			logger.Error(createIndex)
		} else {
			log.Info("successfully created index ", indexname)
		}
	} else {
		log.Info("Index ", indexname, " already exist")
	}
}

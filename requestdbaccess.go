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
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"time"
)

type ArtifactRequest struct {
	RemoteHost   string    `json:"remotehost"`
	UserID       string    `json:"userid"`
	ArtifactPath string    `json:"artifactpath"`
	Requested    time.Time `json:"requested"`
}

func (econf *ElasticConfig) insertRequest(remoteHost string, userid string, path string) {
	logger := log.WithFields(log.Fields{
		"service": "go-addrequest-access",
	})

	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()

	request := &ArtifactRequest{
		UserID: userid, ArtifactPath: path, RemoteHost: remoteHost, Requested: CreateTime()}

	econf.ElasticClient.Index().
		Index(econf.IndexConf.requestdb).
		BodyJson(request). // pass struct instance to BodyJson
		Do(ctx)            // Initiate API call with context object

	// Check for errors in each iteration
	logger.Info("Elasticsearch document indexed:", request)
}

func CreateTime() time.Time {
	return time.Now()
}

func (econf *ElasticConfig) countrequest(remoteHost string, userid string, path string) int64 {

	logger := log.WithFields(log.Fields{
		"service": "go-countrequest-access",
	})

	ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	q := elastic.NewMatchQuery("artifactpath", path)

	bq := elastic.NewBoolQuery().Must(q)
	if remoteHost != "" {
		bq.Must(elastic.NewMatchQuery("remotehost", remoteHost))
	}
	if userid != "" {
		bq.Must(elastic.NewMatchQuery("userid", userid))
	}

	result, err := econf.ElasticClient.Count().
		Index(econf.IndexConf.requestdb).
		Query(bq).
		Do(ctx)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Found Artifact: ", path, ", UserID: ", userid, ", Host: ", remoteHost)
	return result
}

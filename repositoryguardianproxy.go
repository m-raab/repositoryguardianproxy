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
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"os"
)

type ElasticConfig struct {
	Host     string
	User     string
	Password string

	ElasticClient *elastic.Client
	IndexConf     IndexConfig
}

type IndexConfig struct {
	userdb    string
	requestdb string
}

type ProxyConfig struct {
	Path string
	Host string
	Port string

	MaxCount int64
}

func main() {
	indexConf := &IndexConfig{
		userdb:    "proxyusers",
		requestdb: "artifactrequests",
	}

	econf := &ElasticConfig{
		Host:     "http://0.0.0.0:9200",
		User:     "elastic",
		Password: "changeme",

		IndexConf: *indexConf,
	}

	conf := &ProxyConfig{
		Host: "http://localhost:8081",
		Path: "/repository",
		Port: "8090",

		MaxCount: 10,
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	econf.elasticSetup()
	econf.checkIndeces()
	//econf.insertUsers()

	conf.runProxyServer(econf)
}

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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (conf *ProxyConfig) runProxyServer(econf *ElasticConfig) {

	logger := log.WithFields(log.Fields{
		"service": "go-reverse-proxy",
	})

	origin, err := url.Parse(conf.Host + conf.Path)

	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("Serve requests for ", origin)
	}

	path := "/"
	proxy := httputil.NewSingleHostReverseProxy(origin)

	http.HandleFunc(path, handler(proxy, econf, conf))
	err = http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), nil)
	if err != nil {
		logger.Error(err)
	}
}

func handler(p *httputil.ReverseProxy, econf *ElasticConfig, conf *ProxyConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		logger := log.WithFields(log.Fields{
			"service": "go-reverse-proxy-handler",
		})

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		var currentUser = &User{
			UserID:       "UNKNOWN",
			UserPassword: "",
			UserName:     "",
			Access:       "",
		}

		var validRequest = false

		if len(auth) == 2 {
			if auth[0] != "Basic" {
				http.Error(w, "authorization failed", http.StatusUnauthorized)
				return
			}
			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)

			if len(pair) == 2 {
				currentUser = econf.findUser(pair[0])
				validRequest = currentUser.UserID != "UNKNOWN" && currentUser.UserPassword == pair[1]
			}

			if !validRequest {
				logger.Warn("User ", pair[0], " is not valid!")
				http.Error(w, "Not authorized!", http.StatusForbidden)
			}
		} else {
			logger.Warn("Unauthorized request detected from ", r.URL.String(), ".")
		}

		if validRequest && currentUser.Access == "all" {
			r.Header.Set("Authorization", "Basic "+basicAuth("admin", "admin123"))
		}

		if r.Method != http.MethodHead {
			finalRemote := getIPFromRemote(r.RemoteAddr)

			if r.Method == http.MethodPut {
				http.Error(w, "Only read requests are allowed!", http.StatusForbidden)
				logger.Warn("Request method is not allowed for ", finalRemote,
					" from ", r.URL.String(), " (", currentUser.UserID, ").")
				return
			}

			econf.insertRequest(finalRemote, currentUser.UserID, r.URL.String())
			maxCount := conf.MaxCount

			if currentUser.UserID == "UNKNOWN" {
				maxCount = conf.MaxCount * 10
			}

			if econf.countrequest(finalRemote, currentUser.UserID, r.URL.String()) > maxCount {
				http.Error(w, "The number of allowed requests was reached!", http.StatusTooManyRequests)
				logger.Warn("Number of allowed requests was reached for ", finalRemote,
					" from ", r.URL.String(), " (", currentUser.UserID, ").")
				return
			}
		}

		p.ServeHTTP(w, r)
	}
}

func getIPFromRemote(remoteHost string) string {
	s := strings.Split(remoteHost, ":")
	if len(s) > 1 {
		return s[0]
	}
	return remoteHost
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goalisa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	//	"github.com/google/uuid"
)

const (
	httpTimeoutSec = 10
)

type popClient struct {
}

func (pop *popClient) request(params map[string]string, url, accessSecret string) (string, error) {
	sign := pop.signature(params, "POST", accessSecret)
	params["Signature"] = sign

	cli := &http.Client{
		Timeout: time.Second * httpTimeoutSec,
	}
	rsp, err := cli.PostForm(url, convertMap2URLValues(params))
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	body, _ := ioutil.ReadAll(rsp.Body)
	return string(body), nil
}

// Follow https://help.aliyun.com/document_detail/25492.html
func (pop *popClient) signature(params map[string]string, httpMethod, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	qry := ""
	for _, k := range keys {
		qry += fmt.Sprintf("&%s=%s", percentEncode(k), percentEncode(params[k]))
	}
	signSrc := percentEncode(httpMethod) + "&" + percentEncode("/") + "&" + percentEncode((qry[1:]))
	hm := hmac.New(sha1.New, []byte(secret+"&"))
	hm.Write([]byte(signSrc))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}

func percentEncode(s string) string {
	es := url.QueryEscape(s)
	es = strings.Replace(es, "+", "%20", -1)
	es = strings.Replace(es, "*", "%2A", -1)
	return strings.Replace(es, "%7E", "~", -1)
}

func convertMap2URLValues(src map[string]string) url.Values {
	uv := make(url.Values)
	for k, v := range src {
		uv.Add(k, v)
	}
	return uv
}

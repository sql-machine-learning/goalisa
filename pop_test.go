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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	a := assert.New(t)

	params := map[string]string{
		"name":     "由由",
		"age":      "3",
		"homepage": "http://little4.kg?true",
	}
	sign := signature(params, "POST", "test_secret_key")
	a.Equal("6kvgvUDEHtFdZKj8+HhtAS1ovHY=", sign)
}

func TestPercentEncoding(t *testing.T) {
	a := assert.New(t)
	s1 := percentEncode("由由")
	a.Equal("%E7%94%B1%E7%94%B1", s1)
	s2 := percentEncode("http://little4.kg?true")
	a.Equal("http%3A%2F%2Flittle4.kg%3Ftrue", s2)
}

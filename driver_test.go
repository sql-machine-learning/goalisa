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
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	a := assert.New(t)
	d := &Driver{}

	b64Str := base64.RawURLEncoding.EncodeToString([]byte(`{"value1": "param1"}`))
	conn, err := d.Open(fmt.Sprintf("pop_access_id:pop_access_key@pop_url?curr_project=proj&env=%s&with=%s", b64Str, b64Str))
	a.NoError(err)
	a.NotNil(conn)
}

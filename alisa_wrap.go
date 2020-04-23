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
	"fmt"
	"io"
	"os"
	"time"
)

const (
	waitInteveral    = time.Duration(2) * time.Second
	readResultsBatch = 20
)

func (ali *alisa) execOut(cmd string, w io.Writer) error {
	_, err := ali.run(cmd, false, w)
	return err
}

func (ali *alisa) exec(cmd string) error {
	return ali.execOut(cmd, os.Stdout)
}

func (ali *alisa) query(cmd string) (*alisaTaskResult, error) {
	return ali.run(cmd, true, os.Stdout)
}

func (ali *alisa) run(cmd string, resultExpected bool, w io.Writer) (*alisaTaskResult, error) {
	taskID, status, err := ali.createTask(cmd)
	if err != nil {
		return nil, err
	}
	if ali.Verbose {
		return ali.trackingTaskWithLog(taskID, status, resultExpected, w)
	}
	return ali.trackingTaskQuietly(taskID, status, resultExpected)
}

func (ali *alisa) trackingTaskWithLog(taskID string, status int, resultExpected bool, w io.Writer) (*alisaTaskResult, error) {
	var err error
	logOffset := 0
	for !ali.completed(status) {
		if status == alisaTaskWaiting || status == alisaTaskAllocate {
			io.WriteString(w, "waiting for resources")
		} else if status == alisaTaskRunning && logOffset >= 0 {
			if logOffset, err = ali.readLogs(taskID, logOffset, w); err != nil {
				return nil, err
			}
		}
		time.Sleep(waitInteveral)
		if status, err = ali.getStatus(taskID); err != nil {
			return nil, err
		}
	}

	if status == alisaTaskExpired {
		io.WriteString(w, "waiting for resources timeout")
	} else {
		if logOffset >= 0 {
			if logOffset, err = ali.readLogs(taskID, logOffset, w); err != nil {
				return nil, err
			}
		}
		if status == alisaTaskCompleted {
			if resultExpected {
				return ali.getResults(taskID, readResultsBatch)
			}
			return nil, nil
		}
	}
	return nil, fmt.Errorf("invalid task status=%d", status)
}

func (ali *alisa) trackingTaskQuietly(taskID string, status int, resultExpected bool) (*alisaTaskResult, error) {
	var err error
	for !ali.completed(status) {
		time.Sleep(waitInteveral)
		if status, err = ali.getStatus(taskID); err != nil {
			return nil, err
		}
	}

	if status == alisaTaskCompleted {
		if resultExpected {
			return ali.getResults(taskID, readResultsBatch)
		}
		return nil, nil
	}
	return nil, fmt.Errorf("invalid task status=%d", status)
}

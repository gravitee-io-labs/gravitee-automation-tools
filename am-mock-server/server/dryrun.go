// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const dryRunMessage = "Mock server is in dry-run-reject mode, all PUT request are rejected on purpose when ?dryRun=true"

var dryRunErrors = []DryRunError{{Severity: new(SeverityError), Message: new(dryRunMessage)}}

// DryRun skips PUT persistence when ?dryRun=true and returns a fixed DryRunError list.
func DryRun(reject bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isDryRun(r) && r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if reject {
					_ = json.NewEncoder(w).Encode(dryRunErrors)
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isDryRun(r *http.Request) bool {
	ok, err := strconv.ParseBool(r.URL.Query().Get("dryRun"))
	return err == nil && ok
}

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

package errors

// NoAuthProvided is returned when APIContext has neither bearer nor basic credentials.
const NoAuthProvided = ClientAuthError("no auth configured: provide basic auth or bearer auth credentials")

// ManyAuthProvided is returned when APIContext has both bearer and basic credentials.
const ManyAuthProvided = ClientAuthError("only one auth can be configured: provide basic auth or bearer auth credentials")

// ClientAuthError is a sentinel auth-config error. Compare with errors.Is.
type ClientAuthError string

func (e ClientAuthError) Error() string {
	return string(e)
}

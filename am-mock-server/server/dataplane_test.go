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
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/dataplane"
)

func TestListDataPlanes(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).DataPlanes.Put(DataPlane{Id: "test"})

	assertListEqual(t, dataPlanesURL(srv), []DataPlane{{Id: "test"}})
}

func TestGetDataPlane(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).DataPlanes.Put(DataPlane{Id: "test"})

	assertGetEqual(t, dataPlanesURL(srv), "test", DataPlane{Id: "test"})
}

func TestGetDataPlane404(t *testing.T) {
	_, srv := createAMServer(t)

	assertGet404(t, dataPlanesURL(srv), "test", "DataPlane")
}

func TestPutGetDataPlane(t *testing.T) {
	_, srv := createAMServer(t)
	body := DataPlane{Id: "test"}

	assertPutEqual(t, dataPlanesURL(srv), body)
	assertGetEqual(t, dataPlanesURL(srv), "test", body)
}

func TestDeleteDataPlane(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).DataPlanes.Put(DataPlane{Id: "test"})

	assertDeleteGone(t, dataPlanesURL(srv), "test", "DataPlane")
}

func TestListDataPlanesSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).DataPlanes.Put(DataPlane{Id: "test"})
	client := newAMClient(t, srv)

	res, err := client.DataPlanes.ListDataPlanesWithResponse(t.Context())
	assertSDKOK(t, res, err, []dataplane.DataPlane{{Id: "test"}})
}

func TestGetDataPlaneSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).DataPlanes.Put(DataPlane{Id: "test"})
	client := newAMClient(t, srv)

	res, err := client.DataPlanes.GetDataPlaneWithResponse(t.Context(), "test")
	assertSDKOK(t, res, err, dataplane.DataPlane{Id: "test"})
}

func TestGetDataPlane404SDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	res, err := client.DataPlanes.GetDataPlaneWithResponse(t.Context(), "test")
	assertSDK404(t, res, err)
}

func TestPutGetDataPlaneSDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)
	body := dataplane.DataPlane{Id: "test"}

	put, err := client.DataPlanes.UpsertDataPlaneWithResponse(t.Context(), body)
	assertSDKOK(t, put, err, body)

	get, err := client.DataPlanes.GetDataPlaneWithResponse(t.Context(), "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteDataPlaneSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).DataPlanes.Put(DataPlane{Id: "test"})
	client := newAMClient(t, srv)

	del, err := client.DataPlanes.DeleteDataPlaneWithResponse(t.Context(), "test")
	assertSDKNoContent(t, del, err)

	get, err := client.DataPlanes.GetDataPlaneWithResponse(t.Context(), "test")
	assertSDK404(t, get, err)
}

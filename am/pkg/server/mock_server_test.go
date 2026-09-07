package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDomains(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})

	res, err := http.Get(getDomainsUrl(srv))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	failOnNotOK(t, res)

	domains := decodeToSliceOf[Domain](t, res)
	assert.Len(t, domains, 1)
	assert.Equal(t, Domain{Key: "test", Name: "Test domain"}, domains[0])
}

func TestGetDomain(t *testing.T) {
	am, srv := createAMServer(t)

	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertGetEqual(t, getDomainsUrl(srv), "test", Domain{Key: "test", Name: "Test domain"})
}

func TestGetDomain404(t *testing.T) {
	am := NewMockAM(t.Context())

	srv := httptest.NewServer(New(am))
	t.Cleanup(srv.Close)

	assertGet404(t, getDomainsUrl(srv), "test")
}

func TestPutGetDomain(t *testing.T) {
	_, srv := createAMServer(t)

	resp := httpPut(t, getDomainsUrl(srv), Domain{Key: "test", Name: "Test domain"})
	defer resp.Body.Close()

	respDomain := decodeTo[Domain](t, resp)

	assert.Equal(t, Domain{Key: "test", Name: "Test domain"}, respDomain)

	assertGetEqual(t, getDomainsUrl(srv), "test", Domain{Key: "test", Name: "Test domain"})
}

func TestDeleteDomain(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})

	res := httpDelete(t, getDomainsUrl(srv), "test")
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	assertGet404(t, getDomainsUrl(srv), "test")

}

func getDomainsUrl(srv *httptest.Server) string {
	return srv.URL + BasePath + "/organizations/DEFAULT/environments/DEFAULT/domains"
}

func assertGetEqual[T any](t *testing.T, url, key string, expected T) {
	resp := httpGet(t, url, key)
	defer resp.Body.Close()

	failOnNotOK(t, resp)

	var given T
	if err := json.NewDecoder(resp.Body).Decode(&given); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, given)
}

func failOnNotOK(t *testing.T, resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d: %s", resp.StatusCode, body)
	}
}

func assertGet404(t *testing.T, url string, key string) {
	resp := httpGet(t, url, key)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	apiError := decodeTo[Error](t, resp)

	assert.Equal(t, Error{
		HttpStatus: new(int32(404)),
		Message:    new("Domain [test] not found"),
	}, apiError)
}

func createAMServer(t *testing.T) (*MockAM, *httptest.Server) {
	am := NewMockAM(t.Context())
	srv := httptest.NewServer(New(am))
	t.Cleanup(srv.Close)
	return am, srv
}

func httpGet(t *testing.T, url string, key string) *http.Response {
	resp, err := http.Get(url + "/" + key)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func httpPut(t *testing.T, url string, body any) *http.Response {
	request, err := http.NewRequest(http.MethodPut, url, encode(t, body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func httpDelete(t *testing.T, url, key string) *http.Response {
	request, err := http.NewRequest(http.MethodDelete, url+"/"+key, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func encode(t *testing.T, object any) *bytes.Buffer {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(object)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func decodeToSliceOf[T any](t *testing.T, resp *http.Response) []T {
	target := make([]T, 0)
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		t.Fatal(err)
	}
	return target
}

func decodeTo[T any](t *testing.T, resp *http.Response) T {
	target := new(T)
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		t.Fatal(err)
	}
	return *target
}

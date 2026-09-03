package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type listOK struct{ StrictUnimplemented }

func (listOK) AutomationListDomains(context.Context, AutomationListDomainsRequestObject) (AutomationListDomainsResponseObject, error) {
	return AutomationListDomains200JSONResponse{}, nil
}

func TestNew_ListDomains(t *testing.T) {
	srv := httptest.NewServer(New(listOK{}))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + BasePath + "/organizations/DEFAULT/environments/DEFAULT/domains")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("status %d: %s", res.StatusCode, body)
	}
	var domains []Domain
	if err := json.NewDecoder(res.Body).Decode(&domains); err != nil {
		t.Fatal(err)
	}
	if len(domains) != 0 {
		t.Fatalf("got %#v, want empty list", domains)
	}
}

func TestNew_Unimplemented(t *testing.T) {
	srv := httptest.NewServer(New(nil))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + BasePath + "/organizations/DEFAULT/environments/DEFAULT/domains")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotImplemented {
		t.Fatalf("status %d, want 501", res.StatusCode)
	}
}

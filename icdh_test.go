package icdh_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/icdh"
)

func newTestTLSServerWithPathValidator(respPayload string, wantURI string, t *testing.T) *httptest.Server {
	t.Helper()

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqURI := r.RequestURI
		verifyURIs(wantURI, gotReqURI, t)
		fmt.Fprint(w, respPayload)
	}))

	t.Cleanup(func() {
		ts.Close()
	})
	return ts
}

// verifyURIs is a test helper function that verifies if provided URIs are equal.
func verifyURIs(wanturi, goturi string, t *testing.T) {
	t.Helper()

	wantU, err := url.Parse(wanturi)
	if err != nil {
		t.Fatalf("error parsing URL %q, %v", wanturi, err)
	}
	gotU, err := url.Parse(goturi)
	if err != nil {
		t.Fatalf("error parsing URL %q, %v", wanturi, err)
	}

	if !cmp.Equal(wantU.Path, gotU.Path) {
		t.Fatal(cmp.Diff(wantU.Path, gotU.Path))
	}

	wantQuery, err := url.ParseQuery(wantU.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	gotQuery, err := url.ParseQuery(gotU.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantQuery, gotQuery) {
		t.Fatal(cmp.Diff(wantQuery, gotQuery))
	}
}

func TestStatsFor_CallsAPIWithValidPath(t *testing.T) {
	t.Parallel()

	var called bool
	wantPath := "/probe/foo.service.com"

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqPath := r.RequestURI
		verifyURIs(wantPath, gotReqPath, t)

		fmt.Fprint(w, validPayload)
		called = true
	}))
	defer ts.Close()

	client, err := icdh.NewClient(
		ts.URL,
		icdh.WithHTTPClient(ts.Client()),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetStats(context.Background(), "foo.service.com")
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("handler not called")
	}
}

func TestTSStatsFor_CallsAPIWithValidPath(t *testing.T) {
	t.Parallel()

	var called bool
	wantPath := "/probe/ts/foo"

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqPath := r.RequestURI
		verifyURIs(wantPath, gotReqPath, t)

		fmt.Fprint(w, validPayload)
		called = true
	}))
	defer ts.Close()

	client, err := icdh.NewClient(
		ts.URL,
		icdh.WithHTTPClient(ts.Client()),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetTSStats(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("handler not called")
	}
}

func TestCheck_RetrievesStatsForExistingHost(t *testing.T) {
	t.Parallel()

	wantPath := "/probe/foo.service.com"

	ts := newTestTLSServerWithPathValidator(validPayload, wantPath, t)

	client, err := icdh.NewClient(ts.URL, icdh.WithHTTPClient(ts.Client()))
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.GetStats(context.Background(), "foo.service.com")
	if err != nil {
		t.Fatal(err)
	}

	want := icdh.Stats{
		Total:     12,
		Up:        6,
		Unhealthy: 6,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCheck_RetrievesTSStatsForExistingName(t *testing.T) {
	t.Parallel()

	wantPath := "/probe/ts/foo"

	ts := newTestTLSServerWithPathValidator(validPayload, wantPath, t)

	client, err := icdh.NewClient(ts.URL, icdh.WithHTTPClient(ts.Client()))
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.GetTSStats(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}

	want := icdh.Stats{
		Total:     12,
		Up:        6,
		Unhealthy: 6,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestStatsFor_ErrorsOnNotExistingTSName(t *testing.T) {
	t.Parallel()

	wantPath := "/probe/ts/bogusname"

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqPath := r.RequestURI
		verifyURIs(wantPath, gotReqPath, t)

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, validPayload)
	}))
	defer ts.Close()

	client, err := icdh.NewClient(ts.URL, icdh.WithHTTPClient(ts.Client()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetTSStats(context.Background(), "bogusname")
	if err == nil {
		t.Fatal("want err, got nil")
	}
}

func TestStatsFor_ErrorsOnEmptyHost(t *testing.T) {
	t.Parallel()

	wantPath := "/probe/"

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqPath := r.RequestURI
		verifyURIs(wantPath, gotReqPath, t)

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, emptyPayload)
	}))
	defer ts.Close()

	client, err := icdh.NewClient(ts.URL, icdh.WithHTTPClient(ts.Client()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetStats(context.Background(), "")
	if err == nil {
		t.Fatal("want err, got nil")
	}
}

func TestStatsFor_ErrorsOnEmptyTSName(t *testing.T) {
	t.Parallel()

	wantPath := "/probe/ts/"

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqPath := r.RequestURI
		verifyURIs(wantPath, gotReqPath, t)

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, emptyPayload)
	}))
	defer ts.Close()

	client, err := icdh.NewClient(ts.URL, icdh.WithHTTPClient(ts.Client()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetTSStats(context.Background(), "")
	if err == nil {
		t.Fatal("want err, got nil")
	}
}

var (
	validPayload = `{"Total":12,"Up":6,"Unhealthy":6}`
	emptyPayload = ``
)

package txsub

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/distributed_lab/corer"
	"net/http"
	"testing"
)

func mustConnector(client *http.Client, url string) corer.Connector {
	c, err := corer.NewConnector(client, url)
	if err != nil {
		panic(err)
	}

	return c
}

func TestDefaultSubmitter(t *testing.T) {
	ctx := context.Background()

	Convey("submitter (The default Submitter implementation)", t, func() {

		Convey("submits to the configured stellar-core instance correctly", func() {
			server := NewStaticMockServer(`{
				"status": "PENDING",
				"error": null
				}`)
			defer server.Close()

			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, server.URL))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldBeNil)
			So(sr.Duration, ShouldBeGreaterThan, 0)
			So(server.LastRequest.URL.Query().Get("blob"), ShouldEqual, "hello")
		})

		Convey("succeeds when the stellar-core responds with DUPLICATE status", func() {
			server := NewStaticMockServer(`{
				"status": "DUPLICATE",
				"error": null
				}`)
			defer server.Close()

			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, server.URL))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldBeNil)
		})

		Convey("errors when the stellar-core url is not reachable", func() {
			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, "http://127.0.0.1:65535"))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldNotBeNil)
		})

		Convey("errors when the stellar-core returns an unparseable response", func() {
			server := NewStaticMockServer(`{`)
			defer server.Close()

			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, server.URL))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldNotBeNil)
		})

		Convey("errors when the stellar-core returns an exception response", func() {
			server := NewStaticMockServer(`{"exception": "Invalid XDR"}`)
			defer server.Close()

			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, server.URL))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldNotBeNil)
			So(sr.Err.Error(), ShouldContainSubstring, "Invalid XDR")
		})

		Convey("errors when the stellar-core returns an unrecognized status", func() {
			server := NewStaticMockServer(`{"status": "NOTREAL"}`)
			defer server.Close()

			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, server.URL))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldNotBeNil)
			So(sr.Err.Error(), ShouldContainSubstring, "NOTREAL")
		})

		Convey("errors when the stellar-core returns an error response", func() {
			server := NewStaticMockServer(`{"status": "ERROR", "error": "1234"}`)
			defer server.Close()

			s := NewDefaultSubmitter(mustConnector(http.DefaultClient, server.URL))
			sr := s.Submit(ctx, "hello")
			So(sr.Err, ShouldHaveSameTypeAs, &txSubError{})
			ferr := sr.Err.(Error)
			So(ferr.ResultXDR(), ShouldEqual, "1234")
		})
	})
}

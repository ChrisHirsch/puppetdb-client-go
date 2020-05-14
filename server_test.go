package puppetdb

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestSetHTTPTimeout(t *testing.T) {
	s := NewServer("http://localhost:8080")
	expect := time.Second * 15
	s.SetHTTPTimeout(expect)
	if s.HTTPTimeout != expect {
		t.Errorf("Expected HTTPTimeout set to 15 got %s", s.HTTPTimeout)
	}
}

func TestServer_SetHeader(t *testing.T) {
	type fields struct {
		BaseUrl       string
		HTTPTransport http.RoundTripper
		HTTPTimeout   time.Duration
		Headers       map[string]string
	}
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				BaseUrl:       tt.fields.BaseUrl,
				HTTPTransport: tt.fields.HTTPTransport,
				HTTPTimeout:   tt.fields.HTTPTimeout,
				Headers:       tt.fields.Headers,
			}
			s.SetHTTPTimeout(tt.args.t)
			// I don't think this is really doing anything??
			os.Setenv("PUPPET_TOKEN", "bad_suuuuuuupersecrettoke")
			s.SetHeader("X-Authentication", os.Getenv("PUPPET_TOKEN"))
			if s.Headers["PUPPET_TOKEN"] != "suuuuuuupersecrettoken" {
				panic("Header wasn't properly set")
			}
		})
	}
}

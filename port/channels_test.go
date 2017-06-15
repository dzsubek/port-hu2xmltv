package port

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	mockSrv *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup()  {
	mockSrv = http.NewServeMux()
	server = httptest.NewServer(mockSrv)
	client = NewClient(nil)
	client.baseURL = server.URL
	fmt.Println(server.URL)
}

func teardown() {
	server.Close()
}

func Test_Example(t *testing.T) {
	setup()
	defer teardown()

	mockSrv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Date": "2010-10-10T10:10:10+02:00"}`)
		if "GET" != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, "GET")
		}
	})

	var zoneinfo, _ = time.LoadLocation("CET")
	fmt.Println(zoneinfo)

	a := time.Date(2010, 10, 10, 10, 10, 10, 0, zoneinfo)
	list := client.GetChannelList()
	expected := PortInitResponse{Date: a}
	fmt.Println(expected.Date.Format("Mon Jan 2 15:04:05.000 -0700 MST 2006"))
	fmt.Println(list.Date.Format("Mon Jan 2 15:04:05.000 -0700 MST 2006"))

	if expected.Date.Truncate(24*time.Hour).Equal(list.Date.Truncate(24*time.Hour))  {
		t.Fatalf("unexpected response" + list.Date.Format("Mon Jan 2 15:04:05.000 -0700 MST 2006") + " " + expected.Date.Format("Mon Jan 2 15:04:05.000 -0700 MST 2006"))
	}
}
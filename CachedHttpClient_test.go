package CachedHttpClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func startTestServer() {

	if serverStarted {
		return
	}

	if !serverStartedTLS {
		http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {

			rand.Int()

			fmt.Fprintf(writer, "%d", rand.Int())

		})
	}
	serve := func() {
		log.Fatal(http.ListenAndServe(":8082", nil))
	}
	go serve()
	time.Sleep(10000) //wait for the server to accept connections
	serverStarted = true
}

var server = "http://localhost:8082"
var serverTLS = "https://localhost:8081"
var serverStartedTLS = false
var serverStarted = false

func startTestServerTLS() {
	DefaultCashedClient.Transport.(*CachedTransport).Fallback = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if serverStartedTLS {
		return
	}

	if !serverStarted {
		http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {

			rand.Int()

			fmt.Fprintf(writer, "%d", rand.Int())

		})
	}
	serve := func() {
		log.Fatal(http.ListenAndServeTLS(":8081", "testdata/cert.pem", "testdata/key.pem", nil))
	}
	go serve()

	time.Sleep(10000) //wait for the server to accept connections
	serverStartedTLS = true
}

func TestCachedTransport_RoundTrip(t *testing.T) {

	startTestServerTLS()

	request, err := http.NewRequest("GET", serverTLS, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	responseNotCached, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(DefaultCashedClient.Transport.(*CachedTransport).Cache.(*MapCache).cache) != 1 {
		t.Error("request was not save to cache")
		t.FailNow()
	}

	responseCached, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	notCached, err := ioutil.ReadAll(responseNotCached.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cached, err := ioutil.ReadAll(responseCached.Body)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(notCached) != string(cached) {
		t.Error(string(notCached), "!=", string(cached))

	}

}
func TestCachedTransport_RoundTrip_RealServer(t *testing.T) {

	request, err := http.NewRequest("GET", "https://httpbin.org/anything", nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	startTestServerTLS()

	client := DefaultCashedClient
	client.Transport.(*CachedTransport).Cache.(*MapCache).cache = map[string]*http.Response{}
	responseNotCached, err := client.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(client.Transport.(*CachedTransport).Cache.(*MapCache).cache) != 1 {
		t.Error("request was not save to cache")
		t.FailNow()
	}

	responseCached, err := client.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	notCached, err := ioutil.ReadAll(responseNotCached.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cached, err := ioutil.ReadAll(responseCached.Body)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(notCached) != string(cached) {
		t.Error(string(notCached), "!=", string(cached))

	}

}

func TestResponseToJSON(t *testing.T) {
	startTestServerTLS()

	request, err := http.NewRequest("GET", "https://localhost:8081", nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	response, err := DefaultCashedClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = responseToJSON(response)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

func TestCopyResponse(t *testing.T) {

	startTestServerTLS()

	request, err := http.NewRequest("GET", "https://localhost:8081", nil)

	response, err := DefaultCashedClient.Do(request)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	copyResponse, err := CopyResponse(response)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	orgResponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	copyResponseBody, err := ioutil.ReadAll(copyResponse.Body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if bytes.Compare(orgResponseBody, copyResponseBody) != 0 {
		t.Error("bodies not equal")
		t.FailNow()
	}

}

func BenchmarkCopyResponse(b *testing.B) {
	startTestServerTLS()
	request, err := http.NewRequest("GET", "https://localhost:8081", nil)

	response, err := DefaultCashedClient.Do(request)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {

		_, err := CopyResponse(response)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}

	}
}

func TestJSONResponse_ToResponse(t *testing.T) {

	startTestServerTLS()

	request, err := http.NewRequest("GET", "https://localhost:8081", nil)

	response, err := DefaultCashedClient.Do(request)

	jsonResponse, err := NewJsonResponse(response)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	recreatedResponse := jsonResponse.ToResponse()

	recreatedResponse.Request = response.Request
	recreatedResponse.Body = response.Body
	//removing TLS because the ekm can not be recreated and is therefor not equal
	tlsTmp := recreatedResponse.TLS
	recreatedResponse.TLS = response.TLS

	if !reflect.DeepEqual(response, recreatedResponse) {
		t.Error("response not equal")
		t.FailNow()
	}

	recreatedResponse.TLS = tlsTmp

	originalValue := reflect.ValueOf(*response.TLS)
	originalType := reflect.TypeOf(*response.TLS)
	recreatedValue := reflect.ValueOf(*recreatedResponse.TLS)

	//checking deep equal of all TLS attributes except ekm

	for i := 0; i < originalType.NumField(); i++ {
		name := originalType.Field(i).Name
		if name == "ekm" {
			continue
		}

		orgValue := originalValue.FieldByName(name)
		recValue := recreatedValue.FieldByName(name)
		orgInterface := orgValue.Interface()
		recInterface := recValue.Interface()
		if !reflect.DeepEqual(orgInterface, recInterface) {
			t.Error("tlsTmp not equal in field", name)
			t.Log(fmt.Printf("%#v\n", orgInterface.([]*x509.Certificate)[0]))
			t.Log(fmt.Printf("%#v\n", recInterface.([]*x509.Certificate)[0]))
			t.Fail()
		}

	}

}

func BenchmarkNewJSONResponse(b *testing.B) {
	startTestServerTLS()

	request, err := http.NewRequest("GET", "https://localhost:8081", nil)

	response, err := DefaultCashedClient.Do(request)

	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {

		_, err := NewJsonResponse(response)

		if err != nil {
			b.Error(err)
			b.FailNow()
		}

	}
}

func BenchmarkJSONResponse_ToResponse(b *testing.B) {
	benchmarks := []struct {
		name string
	}{
		{"Only ToResponse"},
		{"Always to JSON and ToResponse"},
	}

	startTestServerTLS()

	request, err := http.NewRequest("GET", "https://localhost:8081", nil)

	response, err := DefaultCashedClient.Do(request)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {

			switch bm.name {
			case "Only ToResponse":
				jsonResponse, err := NewJsonResponse(response)
				if err != nil {
					b.Error(err)
					b.FailNow()
				}
				for i := 0; i < b.N; i++ {
					_ = jsonResponse.ToResponse()

				}
			case "Always to JSON and ToResponse":
				for i := 0; i < b.N; i++ {

					jsonResponse, err := NewJsonResponse(response)
					if err != nil {
						b.Error(err)
						b.FailNow()
					}
					_ = jsonResponse.ToResponse()

				}
			}

		})
	}
}

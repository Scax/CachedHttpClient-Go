package CachedHttpClient

//Package CachedHttpClient provides struts and a interface to add a custom cache to a http.client
import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type Cacher interface {
	Get(req *http.Request) (*http.Response, error)
	Set(req *http.Request, res *http.Response) error
}

type CachedTransport struct {
	Cache                         Cacher
	fallback                      http.RoundTripper
	ContinueRoundTripWithSetError func(transport *CachedTransport, err error, request *http.Request, response *http.Response) bool
}

var DefaultCashedClient = &http.Client{
	Transport: DefaultCachedTransport,
}
var DefaultCachedTransport = &CachedTransport{
	Cache:                         NewMapCache(),
	fallback:                      http.DefaultTransport,
	ContinueRoundTripWithSetError: nil,
}

//RoundTrip checks if the cache has a response for the request and return it, if not save the response of the fallback
//RoundTripper to the cache. If the set function returns a error ContinueRoundTripWithSetError will be called if not nil
func (c *CachedTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	if res, err := c.Cache.Get(req); err == nil {
		res.Request = req
		return res, nil

	} else if !errors.Is(err, NotInCacheError) {
		return nil, err
	}
	response, err := c.fallback.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	err = c.Cache.Set(req, response)

	if err == nil {
		return response, nil

	}
	if c.ContinueRoundTripWithSetError == nil {
		return nil, err
	}
	if !c.ContinueRoundTripWithSetError(c, err, req, response) {
		return nil, err
	}

	return response, err

}

var NotInCacheError = errors.New("request not in the cache")

//DumpRequest dumps the request to bytes using httputil.DumpRequest if includeAllHeaders httputil.DumpRequestOut is used
func DumpRequest(req *http.Request, ignoreBody bool, dontIncludeAllHeaders bool) ([]byte, error) {

	var dump []byte
	var err error
	if !dontIncludeAllHeaders {
		dump, err = httputil.DumpRequestOut(req, !ignoreBody)

	} else {
		dump, err = httputil.DumpRequest(req, !ignoreBody)
	}
	if err != nil {
		return nil, err
	}

	return dump, nil
}

//CopyResponse creates a light copy of the response and the body and
//reads the body of the input response into a buffer and places a ReaderCloser of the buffers content in both responses
func CopyResponse(response *http.Response) (*http.Response, error) {

	cRes := *response

	if response.Body == http.NoBody {
		return &cRes, nil
	}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}

	response.Body = ioutil.NopCloser(&buf)
	cRes.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
	return &cRes, nil

}

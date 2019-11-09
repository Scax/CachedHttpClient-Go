package CachedHttpClient

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

//MapCache caches the response in a map string -> *http.Response
//
type MapCache struct {
	cache map[string]*http.Response
	MapCacheOptions
}

type MapCacheOptions struct {
	IgnoreRequestBody            bool
	DontIncludeAllRequestHeaders bool
}

func NewMapCache(options ...MapCacheOptions) *MapCache {

	mapCache := &MapCache{cache: map[string]*http.Response{}}

	if options != nil {
		mapCache.MapCacheOptions = options[0]
	}

	return mapCache
}

func (m *MapCache) Get(req *http.Request) (*http.Response, error) {

	dumpRequest, err := DumpRequest(req, !m.IgnoreRequestBody, m.DontIncludeAllRequestHeaders)
	if err != nil {
		return nil, err
	}

	res, ok := m.cache[string(dumpRequest)]
	if ok {
		cRep, err := CopyResponse(res)
		if err != nil {
			return nil, err
		}
		return cRep, nil
	}
	return nil, NotInCacheError

}

func (m *MapCache) Set(req *http.Request, res *http.Response) error {

	var buf bytes.Buffer
	if res.Body != http.NoBody {
		_, err := buf.ReadFrom(res.Body)
		if err != nil {
			return err
		}
		err = res.Body.Close()
		if err != nil {
			return err
		}
		res.Body = ioutil.NopCloser(&buf)
	}

	dumpRequest, err := DumpRequest(req, !m.IgnoreRequestBody, m.DontIncludeAllRequestHeaders)
	if err != nil {
		return err
	}
	m.cache[string(dumpRequest)] = res

	return nil
}

# CachedHTTPClient-Go



## How to use
Use the DefaultCachedClient
```gotemplate
request,err := http.NewRequest("GET","http://example.com",nil)
DefaultCashedClient.Do(request) //not cached
...
DefaultCachedClient.Do(request) //cached
...
```

Use the DefaultCachedTransport
```gotemplate

request, err := http.NewRequest("GET", "http://example.com", nil)
client := http.Client{
	Transport: DefaultCachedTransport,
}

client.Do(request) //not cached
client.Do(request) //cached
```

Create own CachedTransport
```gotemplate
request, err := http.NewRequest("GET", "http://example.com", nil)

cachedTransport := CachedTransport{
	Cache:    NewMapCache(),
	Fallback: http.DefaultTransport,
}

client := http.Client{
	Transport: &cachedTransport,
}

client.Do(request) //not cached
client.Do(request) //cached
```

Create own Cacher
```gotemplate
type MyCache struct {

}

func (m MyCache) Get(req *http.Request) (*http.Response, error) {
	panic("implement me")
}

func (m MyCache) Set(req *http.Request, res *http.Response) error {
	panic("implement me")
}

func someFunction() {

	request, err := http.NewRequest("GET", "http://example.com", nil)

	cachedTransport := CachedTransport{
		Cache:    NewMapCache(),
		Fallback: http.DefaultTransport,
	}

	client := http.Client{
		Transport: &cachedTransport,
	}

	client.Do(request) //not cached
	client.Do(request) //cached
}
```

## Caches

### Interface Cacher
```gotemplate
type Cacher interface {
	Get(req *http.Request) (*http.Response, error)
	Set(req *http.Request, res *http.Response) error
}
```

### MapCache

```gotemplate
type MapCache struct {
	cache map[string]*http.Response
	MapCacheOptions
}

type MapCacheOptions struct {
	IgnoreRequestBody            bool
	DontIncludeAllRequestHeaders bool
}
```

### FileCache
```gotemplate
type FileCache struct {
	*MapCache
	filePath string
	file     *os.File
}
```
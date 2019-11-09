package CachedHttpClient

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type FileCache struct {
	*MapCache
	filePath string
	file     *os.File
}

func (f *FileCache) Get(req *http.Request) (*http.Response, error) {

	return f.MapCache.Get(req)

}

type FileCacheEntry struct {
	Request  string
	Response *JsonResponse
}

func (f *FileCache) Set(req *http.Request, res *http.Response) error {

	dumpRequest, err := DumpRequest(req, !f.IgnoreRequestBody, f.DontIncludeAllRequestHeaders)

	if err != nil {
		return err
	}

	newJSONResponse, err := NewJsonResponse(res)
	if err != nil {
		return err
	}

	err = json.NewEncoder(f.file).Encode(FileCacheEntry{
		Request:  string(dumpRequest),
		Response: newJSONResponse,
	})

	err = f.MapCache.Set(req, res)

	if err != nil {
		return err
	}

	return err

}

func newFileCache(filePath string, file *os.File, cache *MapCache) *FileCache {

	return &FileCache{
		filePath: filePath,
		file:     file,
		MapCache: cache,
	}

}

//OpenFileCache loaded the cache from an existing cache file
func OpenFileCache(filePath string) (*FileCache, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	fileR, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	mapCache, err := loadMapCacheFromFile(fileR)
	if err != nil {
		return nil, err
	}
	err = fileR.Close()
	if err != nil {
		return nil, err
	}
	return newFileCache(filePath, file, mapCache), nil

}

func loadMapCacheFromFile(file *os.File) (*MapCache, error) {

	scanner := bufio.NewScanner(file)
	responses := map[string]*http.Response{}
	for scanner.Scan() {

		readBytes := scanner.Bytes()

		var entry FileCacheEntry
		err := json.Unmarshal(readBytes, &entry)
		if err != nil {
			return nil, err
		}
		responses[entry.Request] = entry.Response.ToResponse()

	}

	return &MapCache{
		cache: responses,
	}, nil

}

//OpenOrCreateFileCache open the existing cache file or creates a new
func OpenOrCreateFileCache(filePath string) (*FileCache, error) {

	_, err := os.Stat(filePath)
	if err == nil {
		return OpenFileCache(filePath)
	}
	if errors.Is(err, os.ErrNotExist) {
		return NewFileCache(filePath)
	}

	return nil, err
}

//NewFileCache create a new FileCache overriding the cache file
func NewFileCache(filePath string) (*FileCache, error) {
	create, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	err = create.Close()
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return newFileCache(filePath, file, NewMapCache()), nil

}

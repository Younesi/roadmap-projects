package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type CacheEntry struct {
	ResponseBody []byte
	Response     *http.Response
}

type CachingProxyServer struct {
	origin string
	data   sync.Map
}

func NewCachingProxyServer(origin string) *CachingProxyServer {
	return &CachingProxyServer{
		origin: origin,
		data:   sync.Map{},
	}
}

func (cp *CachingProxyServer) Start(w http.ResponseWriter, r *http.Request) {
	CACHE_KEY := r.Method + ":" + r.URL.String()

	if c, ok := cp.data.Load(CACHE_KEY); ok {
		cacheEntry := c.(CacheEntry)
		respondWithHeaders(w, cacheEntry.Response, cacheEntry.ResponseBody, "HIT", CACHE_KEY)
		return
	}

	resp, err := cp.get(r.URL.String())
	if err != nil {
		log.Println("error while requsting : ", err)
		http.Error(w, "Error Forwarding Request", http.StatusInternalServerError)
		return
	}

	cacheEntry, err := cp.cache(resp, CACHE_KEY)
	if err != nil {
		log.Fatal("error while caching : ", err)
		http.Error(w, "Error Forwarding Response", http.StatusInternalServerError)
		return
	}

	respondWithHeaders(w, cacheEntry.Response, cacheEntry.ResponseBody, "MISS", CACHE_KEY)
	return
}

func (cp *CachingProxyServer) ClearCache() error {
	cp.data = sync.Map{}

	return nil
}

func (cp *CachingProxyServer) get(url string) (*http.Response, error) {
	orginURL := cp.origin + url
	resp, err := http.Get(orginURL)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}

func (cp *CachingProxyServer) cache(response *http.Response, cacheKey string) (*CacheEntry, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	ce := CacheEntry{
		ResponseBody: body,
		Response:     response,
	}
	cp.data.Store(cacheKey, ce)

	return &ce, nil
}

func respondWithHeaders(w http.ResponseWriter, response *http.Response, body []byte, cacheHeader, key string) {
	fmt.Printf("\n Cache %s : %s \n", cacheHeader, key)
	w.Header().Set("X-cache", cacheHeader)
	w.WriteHeader(response.StatusCode)

	for k, v := range response.Header {
		w.Header()[k] = v
	}

	w.Write(body)
}

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/cors"
	"github.com/tidwall/gjson"
)

type Image struct {
	Index         int    `json:"index"`
	Url           string `json:"url"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`
	Title         string `json:"title"`
}

const (
	urlBase = "https://cn.bing.com/"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/todaybing", GetLatest7Days)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Use default cors options
	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":5033", handler); err != nil {
		log.Println(err)
	}
}

func GetLatest7Days(w http.ResponseWriter, _ *http.Request) {
	urls := make([]*Image, 7)
	var wg sync.WaitGroup
	for i := 0; i < 7; i++ {
		i := i
		url := fmt.Sprintf("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=%d&n=1&mkt=zh-CN", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			image, err := GetImageData(url)
			if err != nil {
				log.Println(err)
				return
			}
			image.Index = i
			urls[i] = image
		}()
	}
	wg.Wait()
	_ = json.NewEncoder(w).Encode(urls)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetImageData(url string) (*Image, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: tr,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var image Image
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	imageData := gjson.GetBytes(data, "images.0").String()
	if err := json.UnmarshalFromString(imageData, &image); err != nil {
		return nil, err
	}
	image.Url = urlBase + image.Url
	return &image, nil
}

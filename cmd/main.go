package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/cors"
)

type Image struct {
	Index         int    `json:"index"`
	Url           string `json:"url"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`
}

type Response struct {
	Image []*Image `json:"images"`
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
	ch := make(chan Response, 7)
	urls := make([]*Image, 7)
	for i := 0; i < 7; i++ {
		url := fmt.Sprintf("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=%d&n=1&mkt=zh-CN", i)
		go GetLatestDay(url, i, ch)
	}
	for i := 0; i < 7; i++ {
		data := <-ch
		urls[data.Image[0].Index] = data.Image[0]
	}
	close(ch)
	bytes, err := jsoniter.Marshal(urls)
	if err != nil {
		log.Printf("marshal urls err:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bytes)
}

func GetLatestDay(url string, index int, ch chan Response) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		ch <- Response{}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var response Response
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Println(err)
		ch <- Response{}
	}
	response.Image[0].Index = index
	response.Image[0].Url = urlBase + response.Image[0].Url
	ch <- response
}

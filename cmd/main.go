package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
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
	r := gin.Default()
	r.GET("/v1/todaybing", GetLatest7Days)
	if err := http.ListenAndServe(":5033", r); err != nil {
		log.Println(err)
	}
}

func GetLatest7Days(ctx *gin.Context) {
	ch := make(chan Response, 7)
	urls := make([]*Image, 7)
	for i := 0; i < 7; i++ {
		url := "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=" + strconv.Itoa(i) + "&n=1&mkt=zh-CN"
		go GetLatestDay(url, i, ch)
	}
	for i := 0; i < 7; i++ {
		data := <-ch
		urls[data.Image[0].Index] = data.Image[0]
	}
	close(ch)
	ctx.JSON(200, urls)
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
	defer resp.Body.Close()
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

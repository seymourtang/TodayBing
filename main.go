package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
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
type ImageList []*Image

func (i ImageList) Len() int {
	return len(i)
}

func (i ImageList) Swap(m, n int) {
	i[m], i[n] = i[n], i[m]
}
func (i ImageList) Less(m, n int) bool {
	return i[m].Index < i[n].Index
}

const (
	urlBase = "https://cn.bing.com/"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/v1/todaybing", GetLatest7Days)
	if err := http.ListenAndServe(":8088", r); err != nil {
		log.Println(err)
	}
}

func GetLatest7Days(ctx *gin.Context) {
	ch := make(chan Response, 7)
	urls := make([]*Image, 0, 7)
	for i := 0; i < 7; i++ {
		url := "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=" + strconv.Itoa(i) + "&n=1&mkt=zh-CN"
		go GetLatestDay(url, i, ch)
	}
	for i := 0; i < 7; i++ {
		data := <-ch
		urls = append(urls, data.Image[0])
	}
	close(ch)
	sort.Sort(ImageList(urls))
	ctx.JSON(200, urls)
}

func GetLatestDay(url string, index int, ch chan Response) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		ch <- Response{}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err)
		ch <- Response{}
	}
	response.Image[0].Index = index
	response.Image[0].Url = urlBase + response.Image[0].Url
	ch <- response
}

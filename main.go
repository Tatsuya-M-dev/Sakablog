package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var hinataURL = "https://www.hinatazaka46.com"
var cdnHinataURL = "https://cdn.hinatazaka46.com/"
var downloadDir = "download"

func main() {
	checkDir(downloadDir)
	getNewUpdateBlog()
}

//日向坂のブログ更新情報を取得する
func getNewUpdateBlog() {
	url := "https://www.hinatazaka46.com/s/official/diary/member"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	topics := doc.Find("ul.p-blog-top__list")
	topics.Find("li").Each(func(_ int, srg *goquery.Selection) {
		fmt.Print("更新者：")
		name := srg.Find("div.c-blog-top__name")
		member := strings.TrimSpace(name.Text())
		fmt.Println(member)

		dirName := downloadDir + "/" + member
		checkDir(dirName)

		fmt.Print("タイトル：")
		title := srg.Find("p.c-blog-top__title")
		fmt.Println(strings.TrimSpace(title.Text()))

		fmt.Print("更新日時：")
		updDate := srg.Find("time.c-blog-top__date")
		fmt.Println(strings.TrimSpace(updDate.Text()))

		fmt.Print("url：")
		urlElm := srg.Find("a")
		url, _ := urlElm.Attr("href")
		url = hinataURL + url
		fmt.Println(url)
		downloadImageFronPost(url, member)
		fmt.Println("上記記事の画像DL完了")

	})

}

func downloadImageFronPost(url, member string) {
	imageURLSlice := getImageURLFromPost(url)
	for _, value := range imageURLSlice {
		downloadImage(value, member)
	}
}

func getImageURLFromPost(url string) []string {
	var imageURLSlice []string
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	//ブログ内容吐き出し
	//fmt.Println(doc.Html())

	topics := doc.Find("div.l-maincontents--blog")
	topics.Find("img").Each(func(_ int, srg *goquery.Selection) {
		imgURL, _ := srg.Attr("src")
		fmt.Println(imgURL)
		if strings.HasPrefix(imgURL, cdnHinataURL) {
			imageURLSlice = append(imageURLSlice, imgURL)
		}

	})
	return imageURLSlice
}

func downloadImage(url, member string) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	dirName := downloadDir + "/" + member

	filename := lastString(strings.Split(url, "/"))
	filename = dirName + "/" + filename

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	file.Write(body)

}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func checkDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.Mkdir(dirName, 0777)
	}

}

package main

import (
	"bufio"
	"earthling_all_known/dao"
	"earthling_all_known/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gocolly/colly"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// lisences url
const javURL = "http://www.javxxx.com"

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	// 创建数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close()

	dao.DB.AutoMigrate(&models.License{})

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	detailCollector := c.Clone()

	c.OnHTML("div#content div.videos div.video > a", func(e *colly.HTMLElement) {
		detailCollector.Visit(javURL + "/cn" + e.Attr("href")[1:])
	})

	c.OnHTML("div.page_selector a.next", func(e *colly.HTMLElement) {
		log.Println("next page: ", e.Attr("href"))
		c.Visit(javURL + e.Attr("href"))
	})

	saveDir := "./license/"
	exist, err := PathExists(saveDir)
	if err != nil {
		log.Fatalln("check dir failed, err:", err.Error())
		return
	}
	if !exist {
		err := os.Mkdir(saveDir, os.ModePerm)
		if err != nil {
			log.Fatalln("create dir failed, err: ", err.Error())
			return
		}
	}

	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		log.Println("visiting", e.Request.URL.String())
		identifier := e.ChildText("div#video_info div#video_id td.text")
		releaseDate := e.ChildText("div#video_info div#video_date td.text")
		rating := e.ChildText("div#video_info div#video_review span.score")
		actor := e.ChildText("div#video_info div#video_cast span.star")
		imgurl := "https:" + e.ChildAttr("div#video_jacket img#video_jacket_img", "src")

		log.Println("get item:", identifier, releaseDate, rating, actor, imgurl)

		fileName := path.Base(imgurl)

		li := models.License{
			Identifier:  identifier,
			ReleaseDate: releaseDate,
			Rating:      rating,
			Actor:       actor,
			CoverPath:   saveDir + fileName,
		}

		err := models.CreateALicense(&li)
		if err != nil {
			log.Fatalln("create a license failed, err:", err.Error())
		}

		res, err := http.Get(imgurl)
		if err != nil {
			fmt.Println("A error occurred, err: ", err.Error())
		}
		reader := bufio.NewReaderSize(res.Body, 32*1024)

		file, err := os.Create(saveDir + fileName)
		if err != nil {
			panic(err)
		}
		writer := bufio.NewWriter(file)
		io.Copy(writer, reader)

	})

	c.Visit(javURL + "/cn/vl_bestrated.php")

	log.Println("Scraping finished...")
}

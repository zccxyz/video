package web

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
)

type Video struct {
	Size       string
	VideoName  string
	Belong     string
	Uploader   string
	UploadDate string
	ID         string
	VideoUrl   string
}

//Search word视频名称，p页数，c分类
func Search(word, p, c string) (list []*Video, total uint, err error) {
	total = 0
	reqUrl := Host + Methods["search"] + "?p=" + p

	if word != "" {
		reqUrl += "&word=" + word
	}

	//获取网页
	res, err := Request("GET", reqUrl, "")
	if err != nil {
		return list, total, err
	}
	body := res.Body
	defer body.Close()

	//抓取数据
	doc, err := goquery.NewDocumentFromReader(body)
	//总页数
	strNum := doc.Find("body > div.container > div > div.t_p").Text()
	reg := regexp.MustCompile(`共找到([0-9]+)页`)
	subMatch := reg.FindStringSubmatch(strNum)
	if len(subMatch) == 2 {
		num, err := strconv.Atoi(subMatch[1])
		if err == nil {
			total = uint(num)
		}
	}
	doc.Find(`.list ul li`).Each(func(i int, s *goquery.Selection) {
		size := s.Find(".file_size").Text()
		name := s.Find(".file_name").Text()
		user := s.Find(".file_user").Text()
		cate := s.Find(".file_category").Text()
		date := s.Find(".file_dates").Text()
		video := &Video{
			Size:       size,
			VideoName:  name,
			Belong:     cate,
			Uploader:   user,
			UploadDate: date,
		}
		id, ok := s.Attr("data")
		if !ok {
			id = ""
		} else {
			//id有值解析出播放地址
			realUrl, err := ParseUrl(FirstHost + id)
			if err == nil {
				video.VideoUrl = realUrl
			} else {
				log.Println(err)
			}
		}
		video.ID = id
		list = append(list, video)
	})

	//返回数据
	return list, total, nil
}

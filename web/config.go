package web

import (
	"errors"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	Ua        string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"
	Host      string = "http://g.shumafen.cn/"
	Cookie    string = "PHPSESSID=c1jrt2hfr1hbqu824418l173cq; user_key=b945871b50a7d6db2e6550a3bf2513971e364cf5; username=xyz404; userpass=zcc199651"
	FirstHost string = "https://g.shumafen.cn/api/file/84c4e8797fe8d1e8/"
)

var (
	Methods map[string]string = map[string]string{
		"search": "files/shared_files.php", //参数?word="一拳超人"，c="分类"
	}
	Cate map[string]string = map[string]string{
		"movies": "电影",
		"tv":     "电视剧",
		"comic":  "动漫理番",
		"other":  "其他",
	}
)

func ParseUrl(url string) (realUrl string, err error) {
	//获取url
	res, err := Request("GET", url, "")
	if err != nil {
		return "", err
	}
	body := res.Body
	defer body.Close()

	//获取到第一次的加密链接
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	compile := regexp.MustCompile(`var u="../..([a-zA-Z_/]+)";`)
	strSlice := compile.FindStringSubmatch(string(bytes))
	if len(strSlice) != 2 {
		return "", errors.New("数据解析失败1")
	}

	//请求返回的地址
	res, err = Request("GET", Host+"api/"+strSlice[1], url)
	if err != nil {
		return "", err
	}
	body = res.Body
	defer body.Close()

	bytes, err = ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	//匹配得到加密js
	reg := regexp.MustCompile(`<script>((?U)[.\s\S]*)</script>`)
	s := reg.FindStringSubmatch(string(bytes))
	if len(s) < 2 {
		return "", errors.New("数据解析失败2")
	}

	//解密js，获取真实播放地址
	reg = regexp.MustCompile(`.setAttribute\("src",(.+)\)`)
	urlSlice := reg.FindStringSubmatch(string(bytes))
	if len(urlSlice) < 2 {
		return "", errors.New("数据解析失败3")
	}

	vm := otto.New()
	_, err = vm.Run("function getUrl(){" + s[1] + "return " + urlSlice[1] + "}")
	if err != nil {
		return "", err
	}
	val, err := vm.Call("getUrl", nil)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

func Request(method, url, referer string) (rs *http.Response, err error) {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", Ua)
	req.Header.Set("Cookie", Cookie)
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Msg struct {
	Code int
	Msg  string
	Data interface{}
}

package service

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "strconv"
    "time"
    "../dto"
)

type Weibo struct{
    date *time.Time
    text string
    isRetweet bool
}

var dateExp = regexp.MustCompile(`\d+\-\d+\-\d+`)
var dateExp2 = regexp.MustCompile(`\d+\-\d+`)

func GetAllWeibo(userid string) {
    page := 200
    for true {
        fmt.Printf("page:%d\n", page)
        weiboArray, nextpage := qryOnePage(userid, page)
        fmt.Printf("next page:%d\n", nextpage)
        if weiboArray == nil {
            break
        }
        page++
        //fmt.Println(weiboArray)
    }
}

func qryOnePage(userid string, page int) (*[]Weibo, int) {
    //生成client 参数为默认
    client := &http.Client{}

    //生成要访问的url
    url := "https://m.weibo.cn/api/container/getIndex?" +
        "type=uid&value=" + userid +
        "&containerid=107603" + userid
    if page > 1 {
        url = url + "&page=" + strconv.Itoa(page)
    }

    //提交请求
    reqest, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }

    //处理返回结果
    response, _ := client.Do(reqest)
    if 200 != response.StatusCode {
        fmt.Println(response.StatusCode)
        panic("查询失败：" + string(response.StatusCode));
    }

    body, err := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))

    var weiboDto dto.WeiboListQryRespDto
    err = json.Unmarshal([]byte(string(body)), &weiboDto)
    if err != nil {
        fmt.Println("Can't decode json message", err)
        panic(err)
    } else {
        weiboArray := parseWeibo(&weiboDto)
        if weiboArray == nil {
            return nil, 0
        }
        return weiboArray, weiboDto.Data.CardlistInfo.Page
    }
}

func parseWeibo(qryDto *dto.WeiboListQryRespDto) *[]Weibo {
    size := len(*qryDto.Data.Cards)
    if size == 0 {
        return nil
    }
    weiboArray := make([]Weibo, size)
    for i,card := range *qryDto.Data.Cards {
        weiboDate := parseDate(card.Mblog.Created_at)
        weiboText := *filterText(&card.Mblog.Text)
        fmt.Println(weiboText)
        weiboIsRetweet := false
        if card.Mblog.Retweeted_status != nil {
            weiboIsRetweet = true
        }
        weiboArray[i] = Weibo{date : weiboDate, text : weiboText, isRetweet : weiboIsRetweet}
    }
    return &weiboArray
}

func parseDate(date string) *time.Time{
    result := dateExp.FindAllStringSubmatch(date, -1)
    if result != nil {
        weiboDate, _ := time.Parse("2006-01-02", date)
        return &weiboDate
    }

    result = dateExp2.FindAllStringSubmatch(date, -1)
    if result != nil {
        date = time.Now().Format("2006") + "-" + date
        weiboDate, _ := time.Parse("2006-01-02", date)
        return &weiboDate
    }

    return nil
}

func filterText(text *string) *string {
    pattenList := []string {
        `<span.+</span>`, `</a>`, `<a href.+>`}
    for _, patten := range pattenList {
        re, _ := regexp.Compile(patten)
        if re != nil {
            *text = re.ReplaceAllString(*text, "")
        }
    }
    return text
}
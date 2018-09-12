package service

import (
    "../dto"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "strconv"
    "time"
)

type Weibo struct{
    date *time.Time
    text string
    isRetweet bool
    imgUrl *[]string
    largeImgUrl *[]string
}

//获取所有微博
func GetAllWeibo(userid string, savePath string) {
    doc := CreateWord()
    page := 1
    for true {
        fmt.Printf("page:%d\n", page)
        weiboArray, nextpage := qryOnePage(userid, page)
        fmt.Printf("next page:%d\n", nextpage)
        if weiboArray == nil {
            break
        }
        for _, weibo := range *weiboArray{
            if !weibo.isRetweet {
                run := CreateParaRun(doc)
                run = AddText(doc, (*weibo.date).Format("2006-01-02"), run)
                run = AddText(doc,"：", run)
                run = AddText(doc, weibo.text, run)
                if weibo.imgUrl != nil {
                    for i, url := range *weibo.imgUrl {
                        if i % 3 == 0 {
                            run = AddBreak(doc, run)
                        }
                        run = AddImage(doc, url, run, savePath)
                    }
                }
                run = AddBreak(doc, run)
            }
        }
        page++
    }
    Save(doc, savePath + userid + ".docx")
}

//请求weibo获取一页数据
func qryOnePage(userid string, page int) (*[]Weibo, int) {
    //生成要访问的url
    url := "https://m.weibo.cn/api/container/getIndex?" +
        "type=uid&value=" + userid +
        "&containerid=107603" + userid
    if page > 1 {
        url = url + "&page=" + strconv.Itoa(page)
    }

    //处理返回结果
    response, _ := http.Get(url)
    defer response.Body.Close()

    if 200 != response.StatusCode {
        fmt.Println(response.StatusCode)
        panic("查询失败：" + string(response.StatusCode))
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

//解析微博返回的数据
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

        var weiboImgUrl *[]string = nil
        var weiboLargeImgUrl *[]string = nil

        if card.Mblog.Pics != nil {
            ImgUrl := make([]string, len(*card.Mblog.Pics))
            LargeImgUrl := make([]string, len(*card.Mblog.Pics))
            for j, pic := range *card.Mblog.Pics {
                ImgUrl[j] = pic.Url
                LargeImgUrl[j] = pic.Large.Url
            }
            weiboImgUrl = &ImgUrl
            weiboLargeImgUrl = &LargeImgUrl
        }
        weiboArray[i] = Weibo{date : weiboDate,
            text : weiboText,
            isRetweet : weiboIsRetweet,
            imgUrl : weiboImgUrl,
            largeImgUrl : weiboLargeImgUrl}
    }
    return &weiboArray
}

//日期解析
func parseDate(date string) *time.Time {
    var dateExp = regexp.MustCompile(`\d+\-\d+\-\d+`)
    var dateExp2 = regexp.MustCompile(`\d+\-\d+`)
    if dateExp.Match([]byte(date)) {
        weiboDate, _ := time.Parse("2006-01-02", date)
        return &weiboDate
    } else if dateExp2.Match([]byte(date)) {
        date = time.Now().Format("2006") + "-" + date
        weiboDate, _ := time.Parse("2006-01-02", date)
        return &weiboDate
    }

    return nil
}

//过滤表情等
func filterText(text *string) *string {
    pattenList := []string {
        `<span.+</span>`, `</a>`, `<a +href.+>`, `<a +data.+>`, `<br />`}
    for _, patten := range pattenList {
        re, _ := regexp.Compile(patten)
        if re != nil {
            *text = re.ReplaceAllString(*text, "")
        }
    }
    return text
}
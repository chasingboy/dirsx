package httpx

import (
    "time"
    "net"
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "slices"

    "sync"

    "io/ioutil"
    "net/http"
    "crypto/tls"
    "crypto/md5"
    "encoding/hex"
    
    UrlParse "net/url"

    "github.com/schollz/progressbar/v3"
    "github.com/cckuailong/simHtml/simHtml"
    "github.com/spf13/viper"
    "dirsx/common/logger"

    "dirsx/common"

    S "dirsx/common/format"
)


// init logger
var Logger = logger.Logger{}
var TitleReg = regexp.MustCompile("<title>(.+?)</title>")

type Httpx struct {
    Targets chan string
    Method string
    Headers string
    Data string
    Timeout int
    Errnum int

    Threads int

    MaxRespone int
    TitleLen int

    Results [] map[string]string
    Checks [] string
    Excodes [] string
    Smart bool

    // mutex sync.Mutex

}


func (this *Httpx) requester(url string) (map[string] string, bool) {

    defer func() {
            if err := recover(); err != nil {
                // fmt.Println("err:", err)
            }
    }()

    var result = make(map[string] string)

    client := &http.Client {
        Timeout: time.Duration(this.Timeout)*time.Second,
        
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        },

        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }
    
    request, err := http.NewRequest(this.Method, url, nil)
    
    if err != nil {
        return result, false
    }

    this.setHeaders(request)
    respone, err := client.Do(request)

    if err, ok := err.(net.Error); ok && err.Timeout() {
        this.Errnum += 1
        return result, false
    }

    defer respone.Body.Close()
    
    body, err := ioutil.ReadAll(respone.Body)
    
    if err != nil {
        result["body"] = ""
    }

    result["url"] = url
    result["body"] = strings.TrimSpace(string(body))
    result["code"] = strings.ReplaceAll(strconv.Itoa(respone.StatusCode),"206","200")
    result["location"] = this.locationUrl(result["code"], respone.Header)
    result["ctype"] = this.contentType(respone.Header)
    result["clen"] = this.contentLength(respone.Header, string(body))
    result["title"] = this.title(result["code"], result["ctype"], string(body), respone.Header)

    return result, true
}


func (this *Httpx) setHeaders(req *http.Request) {
    req.Header.Add("User-Agent", common.GetRandUserAgent())
    req.Header.Add("Range", fmt.Sprintf("bytes=0-%d",this.MaxRespone))
    req.Header.Add("Connection", "close")
}


func (this *Httpx) contentLength(headers http.Header, body string) string {
    var clen string

    if headers.Get("Content-Range") == "" {
        clen = headers.Get("Content-Length")
    } else {
        clen = strings.Split(headers.Get("Content-Range"),"/")[1]
    }

    if x, _  := strconv.Atoi(clen); x > (this.MaxRespone - 100) {
        clen = fmt.Sprintf("%dMB", x/1024/1024)
    }
    
    if clen == "" {
        clen = S.F("{0}", len(body))
    }

    return clen
}


func (this *Httpx) contentType (headers http.Header) string {
    return headers.Get("Content-Type")
}


func (this *Httpx) locationUrl (code string, headers http.Header) string {
    if strings.HasPrefix(code, "3") {
        return headers.Get("Location")
    }
    return ""
}


func (this *Httpx) bodymd5(body string) string {
    hash := md5.Sum([]byte(body))
    return hex.EncodeToString(hash[:])
}


func (this *Httpx) title(code string, ctype string, body string, headers http.Header) string {

    if strings.HasPrefix(code,"3") {
        return S.F("{0}-> {1}", code, headers.Get("Location"))
    }

    if strings.Contains(ctype, "json") || strings.Contains(ctype, "plain") || ctype == "" {
        title := common.ReplaceStrings(body, " ", "\n", "[", "]")
        
        if len(title) > this.TitleLen {
            return title[:this.TitleLen] + "..."
        }

        return title
    }

    mathers := TitleReg.FindStringSubmatch(strings.ReplaceAll(body, "\n", ""))
    
    if mathers == nil {
        return strings.Split(ctype,";")[0]
    }
    return mathers[1]
}

func (this *Httpx) is_html(body string) bool {
    return (strings.Contains(body, "<html>") || strings.Contains(body, "<body>") || strings.Contains(body, "<script>"))
}

func (this *Httpx) exclude_codes(code string) bool {
    return slices.Contains(this.Excodes, code)
}


func (this *Httpx) is_in_black_title(title string) bool {
    return slices.Contains(common.BLACK_TITLE, title)
}


func (this *Httpx) filter(url string, code string, clen string, title string, body string, ctype string, location string) bool {
    // exclude codes [400,404,500,...]
    if this.exclude_codes(code) == true {
        return false
    }
    
    // if enable automated filtering
    if this.Smart == false {
        return true
    }
    
    // 待完善 ... ...
    
    this.Checks = append(this.Checks, this.code_length(code, clen))
    
    return true
}


func (this *Httpx) generate_basic_checks(url string) bool {
    sign := false

    urls := []string {url, common.JoinUrlAndWord(url,"indexxxxx.zip","","","")}

    for _, url := range urls {
        result, flag := this.requester(url)
        
        if flag {
            Logger.INF().Str("Check url: " + url).Str(" ").State(result["code"],"WAR").Str(" ").State(result["clen"],"WAR").Str(" ").State(this.bodymd5(result["body"]),"WAR").Msgf("")
            
            x := []string {
                result["body"],
                this.bodymd5(result["body"]),
                this.code_length(result["code"], result["clen"]),
            }

            this.Checks = append(this.Checks, x...)
            sign = true
        }
    }
    return sign
}


func (this *Httpx) Reset() *Httpx {
    this.Results, this.Errnum = []map[string]string {}, 0
    this.Checks = []string {}
    return this
}


func (this *Httpx) send_targets(targets []string) *Httpx {
    for _, url := range targets {
        this.Targets <- url
    }

    return this
}


func (this *Httpx) close_targets() {
    time.Sleep(2 * time.Second)
    close(this.Targets)
}


func (this *Httpx) threader(wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
    defer wg.Done()

    for url := range this.Targets {
        result, flag := this.requester(url)
        
        if flag && this.filter(result["url"], result["code"], result["clen"], result["title"], result["body"], result["ctype"], result["location"]) == true {
            if result["code"] == "403" {
                this.send_targets([]string{ strings.TrimRight(result["url"],"/") })
                continue
            }

            result = map[string]string {"url": result["url"], "code": result["code"], "clen": result["clen"], "title": result["title"]}
            
            // this.mutex.Lock()
            this.Results = append(this.Results, result)
            // this.mutex.Unlock()
        }

        bar.Add(1)
    }
}


func (this *Httpx) Runner(url string, targets []string) []map[string]string {
    if this.generate_basic_checks(url) == false {
        Logger.ERR().Msgf("Target requests timeout! Skipping...\n")
        return this.Results
    }

    bar := progressbar.NewOptions(
        len(targets),
        progressbar.OptionEnableColorCodes(true),
        progressbar.OptionShowCount(),
        progressbar.OptionSetWidth(50),
        progressbar.OptionSetDescription(Logger.State("INF","INF").Msg("")),
        progressbar.OptionSetTheme(progressbar.Theme{
            Saucer:        "[green]=[reset]",
            SaucerHead:    "[green]>[reset]",
            SaucerPadding: " ",
            BarStart:      "[",
            BarEnd:        "]",
    }))

    var wg sync.WaitGroup

    for thread := 0; thread < this.Threads; thread ++ {
        wg.Add(1)
        go this.threader(&wg, bar)
    }

    this.send_targets(targets).close_targets()

    wg.Wait()

    fmt.Println("\n")
    return this.Results
}


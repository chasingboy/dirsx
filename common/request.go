package common


import (
	"os"
    "time"
    "strings"
    "math/rand"

    // "fmt"

    UrlParse "net/url"
    S "dirsx/common/format"
)


func GetRandUserAgent() string {
    rand.Seed(time.Now().UnixNano())
    randnum := rand.Intn(len(USER_AGENT))
    return USER_AGENT[randnum]
}


func IsUrlValid(url string) bool {
    if _, err := UrlParse.ParseRequestURI(url); err != nil {
        return false
    }
    
    prs, err := UrlParse.Parse(url)
    if err != nil {
        return false
    }

   return prs.Scheme == "http" || prs.Scheme == "https"
}


func FormatHttpRespone(url string, code string, clen string, title string) string {
    if strings.HasPrefix(code, "2"){
        return Logger.RET().State("200","SUC").Str(" "+url+" ").State(clen,"WAR").Str(" ").State(title,"RET").Msg("")
    }

    return Logger.RET().State(code,"WAR").Str(" "+url+" ").State(clen,"WAR").Str(" ").State(title,"RET").Msg("")
}


func HandleScanResults(url string, results []map[string] string, consoleText string) (string,string) {
    var fileText string = S.F("\n# [RET:{0}] {1}\n", len(results), url)
    
    for _, ret := range results {
        url, code, clen, title := ret["url"], ret["code"], ret["clen"], ret["title"]
        consoleText += FormatHttpRespone(url, code, clen, title) + "\n"
        fileText += S.F("[RET] [{0}] {1}  [{2}]  [{3}]\n", code, url, clen, title)
    }

    if len(results) == 0 {
        consoleText, fileText = Logger.WAR().Msg("Not Result Found...\n"), "# Not Result Found...\n"
    }

    return consoleText, fileText
}


func HandleHttpHeaders(cookie string, headers []string, headersFile string) (map[string] string){

	defer func() {
        if err := recover(); err != nil {
            Logger.WAR().Msgf("Headers argument error, please check... ")
            os.Exit(0)
        }
    }()

	var headersMap = make(map[string] string)
	var _headers [] string

	if headersFile != "" {
		_headers = ReadFile(headersFile)
	}

	if len(headers) != 0 {
		_headers = headers
	}

	for _, header := range _headers {
		x := strings.Split(header, ":")
		headersMap[x[0]] = strings.TrimSpace(x[1])
	}

	if cookie != "" {
		headersMap["Cookie"] = ReplaceStrings(cookie, "", "Cookie: ", "Cookie:", "cookie: ", "cookie:", "COOKIE: ", "COOKIE:")
	}

	return headersMap
}


func SettingProxy(url string) (*UrlParse.URL) {
	proxy, err := UrlParse.Parse(url)

	if err != nil {
		Logger.WAR().Msgf("Proxy setting error, please check...")
		os.Exit(0)
	}

	return proxy
}

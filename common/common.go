package common

import (
	"strings"
	"strconv"
    "regexp"
    "time"
    "io/ioutil"
    "os"
    // "fmt"
    "math/rand"
    "path/filepath"

	UrlParse "net/url"
	S "dirsx/common/format"

    "dirsx/common/logger"

	"github.com/bobesa/go-domain-util/domainutil"
)


var common_backup_words = strings.Split(strings.TrimSpace(BAK_FILES), "\n")

var UrlRemoveReg = regexp.MustCompile(`(\?|\#)(.*)`)

var Logger = logger.Logger {}


func ReadFile(filename string) [] string {
    text, err := ioutil.ReadFile(filename)
    
    if err != nil {
        Logger.WAR().Msgf(S.F("the target file is not exist! => {0}", filename))
        os.Exit(0)
    }
    
    x := strings.Split(strings.TrimSpace(string(text)), "\n")
    
    return x
}


func OutputResultsToFile(filename string, fileText string) {
    err := os.WriteFile(filename, []byte(fileText), 0777)

    if err != nil {
        Logger.ERR().Msgf(S.F("Output results to file failed! => {0}", filename))
    }
}


func GetRandUserAgent() string {
    rand.Seed(time.Now().UnixNano())
    randnum := rand.Intn(len(USER_AGENT))
    return USER_AGENT[randnum]
}


func JoinUrlAndWord(url string, word string, prefix string, suffix string, end string) string {
	return S.F("{0}/{1}{2}{3}{4}", strings.TrimRight(url,"/"), prefix, strings.TrimLeft(word,"/"), suffix, end)
}


func ReplaceStrings(text string, replaces ...string) string {
	for _, replace := range replaces {
		text = strings.ReplaceAll(text, replace, "")
	}
	return text
}


func StringToInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
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


func RemoveExtensions(wordlist []string, removeExt string) []string {
    if removeExt == "" {
        return wordlist
    }

    var wordstring = strings.Join(wordlist, "\n")

    exts := strings.Split(strings.Trim(removeExt, ","), ",")

    for _, ext := range exts {
        wordstring = strings.ReplaceAll(wordstring, S.F(".{0}\n", ext), "\n")
    }

    return strings.Split(wordstring, "\n")
}


func FormatWords(word string, upperTitle bool, upperAll bool) string {
    if upperTitle == true {
        word = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
    }

    if upperAll == true {
        word = strings.ToUpper(word)
    }

    return word
}



func RemoveUrlParams(url string) string {
    return UrlRemoveReg.ReplaceAllString(url, "")
}


func SplitUrlPath(urls []string, isSplit bool) []string {
    var _urls = []string {}
    
    if isSplit == false {
        return RemoveDuplicates(urls)
    }

    for _, url := range urls {
        if IsUrlValid(url) == false {
            continue
        }

        prs, _ := UrlParse.Parse(RemoveUrlParams(url))

        baseurl := S.F("{0}://{1}/", prs.Scheme, prs.Host)
        paths := strings.Trim(filepath.Dir(prs.Path), "/")

        _urls = append(_urls, []string {RemoveUrlParams(url), baseurl}...)

        if paths == "" || paths == "/" {
            continue
        }

        for _, path := range strings.Split(paths, "/") {
            baseurl = JoinUrlAndWord(baseurl, path, "", "", "/")
            _urls = append(_urls, baseurl)
        }
    }

    return RemoveDuplicates(_urls)
}


func RemoveDuplicates(input []string) []string {
    var (
    	encountered map[string] bool = map[string]bool {}
    	result []string = []string {}
    )

    for _, v := range input {
        if encountered[v] == false {
            encountered[v] = true
            result = append(result, v)
        }
    }

    return result
}


func GenerateBackupWords(url string, wordlist []string, isbak bool) []string {
	if isbak == false {
		return wordlist
	}

	exts := []string {".zip",".rar",".war",".bak",".7z",".tar",".gz",".tgz",".tar.gz",".bz2",".tar.bz2",".jar",
    ".zip_bak",".rar_bak",".war_bak",".bak",".7z_bak",".tar_bak",".gz_bak",".tgz_bak",".tar.gz_bak",".bz2_bak",".tar.bz2_bak",}

    var words [] string

    // x.www.baidu.com/a/b
    domain := domainutil.Domain(url)							// baidu.com -> baidu.com.zip
    suffix := domainutil.DomainSuffix(url)  					// com
    fld := strings.Replace(domain, "."+suffix, "", 1)			// baidu	 -> baidu.zip
    subdomain := domainutil.Subdomain(url)						// x.www	 -> x.www.zip

    prs, _ := UrlParse.Parse(url)
    paths := strings.Trim(prs.Path,"/")

    for _, ext := range exts {
    	words = append(words, domain+ext)
    	words = append(words, fld+ext)
    	words = append(words, subdomain+ext)
    	words = append(words, subdomain+"."+domain+ext)

    	for _, sub := range strings.Split(subdomain, ".") {
    		words = append(words, sub+ext)						// x.zip | www.zip
    		words = append(words, sub+"."+domain+ext)
    	}

    	for _, path := range strings.Split(paths, "/") {
    		words = append(words, path+ext)						// a.zip | b.zip
    	}
    }

    // delete null
    for _, word := range words {
    	if strings.HasPrefix(word, ".") == false {
    		wordlist = append(wordlist, word)
    	}
    }

    return append(wordlist, common_backup_words...)
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



func FormatDictsOptions (dictspath string) ([]string,string) {
    files, err := ioutil.ReadDir(dictspath)
    
    if err != nil {
        Logger.WAR().Msgf("The dicts directory is not exists")
        os.Exit(0)
    }

    var (
        idxs [] string
        idx int = 0
        options string
        maxlen int = 0
    )


    for _, file := range files {
        if len(file.Name()) > maxlen {
            maxlen = len(file.Name())
        }
    }


    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".txt") == false && file.IsDir() == false {
            continue
        }

        idxs = append(idxs, file.Name())
        
        index := strings.Repeat("0", 2 - len(strconv.Itoa(idx))) + strconv.Itoa(idx)
        fileName := file.Name() + strings.Repeat(" ", maxlen - len(file.Name()) + 4)
        options += Logger.Inf(S.F("[{0}] {1}", index, fileName)).Msg("")
        
        if idx % 2 == 1 {
            options += "\n"
        }
        
        idx = idx + 1
    }
    
    return idxs, options
}

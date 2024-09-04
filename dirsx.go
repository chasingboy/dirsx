package main

import (    
    "os"
    "fmt"
    "time"
    "path"
    "strings"
    "strconv"
    "path/filepath"

    "dirsx/common/logger"
    "dirsx/common/httpx"
    "github.com/jessevdk/go-flags"

    "dirsx/common"

    S "dirsx/common/format"
)

// init logger
var Logger = logger.Logger {}

var opts common.Options

var MAX_RESPONE int = 10*1024*1024

// node all results
var ALL_RESULTS strings.Builder

var banner string = fmt.Sprintf(`

    ██████╗ ██╗██████╗ ███████╗██╗  ██╗              
    ██╔══██╗██║██╔══██╗██╔════╝╚██╗██╔╝  
    ██║  ██║██║██████╔╝███████╗ ╚███╔╝ 
    ██║  ██║██║██╔══██╗╚════██║ ██╔██╗ 
    ██████╔╝██║██║  ██║███████║██╔╝ ██╗
    ╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
                                       %s
                        xboy@遥遥领先

`, "1.1.0")


var InfoFormat string = `
-----------------------------------------------------------
    Target Num   |    {0}
    Threads      |    {1}
    Wordlist     |    {2}
    Status Code  |    [{3}]
    Current Time |    {4}
-----------------------------------------------------------
`

func PrintScanInfo(tgnum int, wordlist string, threads int) {
    code := Logger.Suc("200").Str(",").War("301,302,307").Str(",400,404,501,502,503").Msg("")
    _, filename := path.Split(wordlist)

    fmt.Println(
        S.F(
            strings.TrimLeft(InfoFormat,"\n"), 
            tgnum, threads, filename, code, 
            time.Now().Format("2006-01-02 15:04:05"),
        ),
    )
}


func ListAndSelectDicts(rootpath string) ([]string, string) {
    dictspath := filepath.Join(filepath.Dir(rootpath), "dicts")
    
    Logger.War("[+] You have not appoint payloads, so you can select from the list: ").Msgf("")

    var (
        idxs [] string
        selnum string
        paths []string
        options string
    )

    idxs, options = common.FormatDictsOptions(dictspath)
    fmt.Println(options)

    for {
        fmt.Printf("[+] Select payloads with number: ")
        fmt.Scan(&selnum)

        selnums := strings.Split(selnum, ",")

        for _, sn := range selnums {
            idx, err := strconv.Atoi(sn)
            
            if err != nil {
                Logger.ERR().Msgf("Number error! No this selectioin! Please input again!")
                os.Exit(0)
            }
            
            if idx >= len(idxs) {
                Logger.ERR().Msgf("Number error! No this selectioin! Please input again!")
                os.Exit(0)
            }

            paths = append(paths, filepath.Join(dictspath, idxs[idx]))
        }

        return paths, selnum
    }
}


func LoadUrls() []string {
    var urls [] string

    if opts.Url == "" && opts.List == "" {
        Logger.WAR().Msgf("Please input the url or[url file]")
        os.Exit(0)
    }
    
    if opts.Url != "" {
        urls = []string {opts.Url}
    } else {
        urls = common.ReadFile(opts.List)
    }

    return common.SplitUrlPath(urls, opts.IsSplit)
}


func LoadPayloads(rootpath string) ([]string, string){
    var wordlist [] string
    var wordname string = "payload list-> "

    if opts.Wordlist != "" {
        wordlist = common.ReadFile(opts.Wordlist)
        wordname = filepath.Base(opts.Wordlist)
    } else {
        paths, selnums := ListAndSelectDicts(rootpath)

        for _, path := range paths {
            wlist := common.ReadFile(path)
            wordlist = append(wordlist, wlist...)
        }

        wordname = wordname + selnums
    }

    return common.RemoveDuplicates(wordlist), wordname
}



func GenerateTargets(url string, wordlist []string) []string {
    var targets, exts []string

    // adding backup wordlist and removing extensions
    wordlist = common.RemoveExtensions(
        common.GenerateBackupWords(url, wordlist, opts.Isbak),
        opts.RemoveExt,
    )

    if opts.Ext != "" {
        exts = strings.Split(strings.Trim(opts.Ext, ","), ",")
        opts.Suffix = opts.Suffix + "."
    } else {
        exts = [] string {""}
    }

    for _, word := range wordlist {
        for _, ext := range exts {
            // remove %ext% of dirsearch wordlist
            if strings.Contains(word, "%EXT%") {
                continue
            }

            word = common.FormatWords(word, opts.UpperTitle, opts.UpperAll)
            
            targets = append(targets, common.JoinUrlAndWord(url, word, opts.Prefix, opts.Suffix, ext))
        }
    }
    // unique targets
    return common.RemoveDuplicates(targets)

}


func DirScan(urls []string, wordlist []string) {

    for index, url := range urls {

        httpx := httpx.Httpx {
            Targets: make(chan string),
            Method: opts.Method,
            Timeout: opts.Timeout,
            MaxRespone: MAX_RESPONE,
            TitleLen: opts.TitleLen,
            Threads: opts.Threads,
            Excodes: strings.Split(opts.Excodes, ","),
            Smart: !opts.Unsmart,
        }

        // unique targets
        targets := GenerateTargets(url, wordlist)
        // fmt.Println(targets)
        
        Logger.INF().Msgf(S.F("Start Scanning Target {0} -> {1}", index+1, url))
        results := httpx.Reset().Runner(url, targets)

        consoleText, fileText := common.HandleScanResults(url, results, "")

        ALL_RESULTS.WriteString(fileText)

        fmt.Println(consoleText)
    }

    if opts.Output != "" {
        common.OutputResultsToFile(opts.Output, ALL_RESULTS.String())
    }
}


func ScanRunner(rootpath string) {
    urls := LoadUrls()

    wordlist, wordname := LoadPayloads(rootpath)

    PrintScanInfo(len(urls), wordname, opts.Threads)

    DirScan(urls, wordlist)

}


func main() {
    fmt.Printf(banner)

    var rootpath, err = os.Executable()
    
    if err != nil {
        Logger.WAR().Msgf("No permission to read dirsx directory, Please check ...")
        return
    }
    
    start := time.Now()

    if _, err := flags.Parse(&opts); err != nil {
        // fmt.Println(err)
        return
    }

    ScanRunner(rootpath)
    
    Logger.INF().Str(S.F("Scanning {0} targets completed [{1}]", 1, time.Since(start))).Msgf("")

    
}

package common

type Options struct {
    Url         string `short:"u" long:"url" description:"input url of target"`
    List        string `short:"l" long:"list" description:"input file containing list of target"`
    Wordlist    string `short:"w" long:"wordlist" description:"appoint wordlist for scanning directory"`
    TitleLen    int    `long:"title-len" description:"set title display length" default:"30"`
    Threads     int    `short:"t" long:"threads" description:"number of threads to use" default:"20"`
    Timeout     int    `long:"timeout" description:"timeout in seconds" default:"5"`
    
    Output      string `short:"o" long:"output" description:"file to write output results"`

    Prefix      string `long:"prefix" description:"add prefix of payloads"`
    Suffix      string `long:"suffix" description:"add suffix of payloads"`
    Ext         string `short:"e" long:"extension" description:"add extension eg: -e php,html"`

    RemoveExt   string `long:"remove-ext" description:"remove extension eg: --remove-ext php | admin.php -> admin"`

    UpperTitle  bool   `long:"upper-title" description:"capitalize the first letter eg: admin -> Admin"`
    UpperAll    bool   `long:"upper-all" description:"capitalize the all letter eg: admin -> ADMIN"`
     
    Isbak       bool   `long:"bak" description:"enable scanning backup file (default:false)"`
    IsSplit     bool   `long:"split" description:"enable spliting the url path, eg: /a/b -> /a/, /a/b (default: false)"`

    Method      string `short:"X" description:"method of http requests" default:"GET"`

    Excodes     string `short:"x" long:"excode" description:"specify the status codes that be filtered eg: 400,404" default:"400,404,406,416,501,502,503"`

    Cookie      string `long:"cookie" description:"set request cookies, eg: --cookie \"session=admin\""`
    Headers   []string `short:"H" long:"headers" description:"set request headers, string[] eg: -H \"Token: admin=true\" -H \"Cookie: login=true\""`
    
    HeadersFile string `long:"headers-file" description:"set request headers file, eg: --headers-file headers.txt"`

    Proxy       string `long:"proxy" description:"set request proxy, eg: --proxy http://127.0.0.1:8080"`

    Unsmart     bool   `long:"no-smart" description:"disable smart mode (automated filtering)"`
}

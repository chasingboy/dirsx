<h1 align="center">dirsx</h1>
<h3 align="center">dirsx 是一款能够自动化过滤扫描结果的目录扫描工具</h3>
<p align="center">
  <img src="https://img.shields.io/badge/Version-V1.6.2-green?style=flat">
  <img src="https://img.shields.io/github/stars/chasingboy/dirsx?style=flat&labelColor=rgb(41%2C52%2C52)&color=green">
  <img src="https://img.shields.io/github/issues/chasingboy/dirsx">
  <img src="https://img.shields.io/github/downloads/chasingboy/dirsx/total?style=flat&labelColor=rgb(41%2C52%2C52)&color=green">
  <img src="https://visitor-badge.laobi.icu/badge?page_id=chasingboy.dirsx&left_color=green&right_color=#66ccff">
</p>

<img width="1154" alt="image" src="https://github.com/user-attachments/assets/87879581-7278-4e3f-8e89-02487f429acd">

### 前言
> 当时正值华为发布遥遥领先, 加上“遥遥领先”只是开个玩笑, 大佬们见笑了

平时使用过 dirsearch｜dirmap 等一些目录扫描工具，针对如今的 WEB 多样化，对扫描结果的过滤总感觉与预期不符合。因此下定决心造个轮子，就这样有了 dirsx。当时是使用 python 写的，但是可移植性不是很好。所以使用 golang 进行重构，顺便学习一下 golang。

### 功能
> 大部分功能其他工具都有, 只是根据个人习惯更改
* 使用 html 相似度对结果进行过滤
* 对 301、302、403 状态进行二次判断
* 对 json 返回结果进行判断
* 字典第一个字母大写｜全部字母大写｜添加前后缀
* 返回页面 title, 如无 title 返回内容前面 30 个字符串 (默认｜设置)
* 自动过滤模式, 默认开启 (开启｜关闭)

### 基本使用
指定字典进行扫描
```bash
dirsx -u https://www.baidu.com -w words.txt
```
指定目录递归扫描, 目前暂无添加结果递归功能扫描，担心目录误报
```bash
dirsx -u https://www.baidu.com -w words.txt --split

# https://www.baidu.com/a/b/
# -> https://www.baidu.com/a/
# -> https://www.baidu.com/a/b/
```
指定备份文件进行扫描
```bash
dirsx -u https://www.baidu.com -w words.txt --bak
```

![image](https://github.com/user-attachments/assets/735dc7f5-f60a-43b3-8d9f-fdf695139aad)

指定添加后缀进行扫描
```bash
dirsx -u https://www.baidu.com -w words.txt --suffix h5

# https://www.baidu.com/admin
# -> https://www.baidu.com/adminh5
```

指定添加 cookie | headers
```bash
# --cookie
dirsx -u https://www.baidu.com -w words.txt --cookie "session=admin"

# --headers
dirsx -u https://www.baidu.com -w words.txt --headers "Authorization: bearer eyJ0eX..." --headers "X-Forwarded-For: 127.0.0.1"

# --headers-file
dirsx -u https://www.baidu.com -w words.txt --headers-file headers.txt
```

内置一些常用字典选择, 在没有指定字典时显示该列表
* 常见目录字典
* dirsearch 的自带字典
* 长度为 1-5 的字母组合
* ... ...
```
~ kali$ dirsx -u http://127.0.0.1/


    ██████╗ ██╗██████╗ ███████╗██╗  ██╗              
    ██╔══██╗██║██╔══██╗██╔════╝╚██╗██╔╝  
    ██║  ██║██║██████╔╝███████╗ ╚███╔╝ 
    ██║  ██║██║██╔══██╗╚════██║ ██╔██╗ 
    ██████╔╝██║██║  ██║███████║██╔╝ ██╗
    ╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
                                       1.1.0
                        xboy@遥遥领先

[+] You have not appoint payloads, so you can select from the list: 
[0] fuzzing-dirs-common.txt
[1] fuzzing-dirs-dirsearch.txt
[2] fuzzing-files-php.txt
[3] fuzzing-letter-len1.txt
[4] fuzzing-letter-len2.txt
[5] fuzzing-letter-len3.txt
[6] fuzzing-letter-len4.txt
... ...
[+] Select payloads with number: 1

```

### --httpx 模式
增加 httpx 模式, 可以在没有 httpx 工具的情况下用来探测 WEB 服务
```
dirsx -u https://www.baidu.com --httpx
```
<img width="1316" alt="image" src="https://github.com/chasingboy/dirsx/blob/main/assets/httpx.png">

### --ffuf 模式
增加 ffuf 模式, 用法与 ffuf 工具一样, 使用 FUZZ 指定 Fuzzing 位置
```
# Fuzzing dirs
dirsx --ffuf -u http://127.0.0.1/FUZZ -w words.txt
dirsx --ffuf -u http://127.0.0.1/FUZZ.php -w words.txt
dirsx --ffuf -u http://127.0.0.1/FUZZ/index.php -w words.txt

# Fuzzing headers
dirsx --ffuf -u http://127.0.0.1/ -H "x-forwarded-for: FUZZ"
```

### dirsx -h

```bash
~ kali$ dirsx -h


    ██████╗ ██╗██████╗ ███████╗██╗  ██╗              
    ██╔══██╗██║██╔══██╗██╔════╝╚██╗██╔╝  
    ██║  ██║██║██████╔╝███████╗ ╚███╔╝ 
    ██║  ██║██║██╔══██╗╚════██║ ██╔██╗ 
    ██████╔╝██║██║  ██║███████║██╔╝ ██╗
    ╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
                                       1.1.0
                        xboy@遥遥领先

Usage:
  dirsx [OPTIONS]

Application Options:
  -u, --url=          input url of target
  -l, --list=         input file containing list of target
  -w, --wordlist=     appoint wordlist for scanning directory
      --title-len=    set title display length (default: 30)
  -t, --threads=      number of threads to use (default: 20)
      --timeout=      timeout in seconds (default: 5)
  -o, --output=       file to write output results
      --prefix=       add prefix of payloads
      --suffix=       add suffix of payloads
  -e, --extension=    add extension eg: -e php,html
      --remove-ext=   remove extension eg: --remove-ext php | admin.php -> admin
      --upper-title   capitalize the first letter eg: admin -> Admin
      --upper-all     capitalize the all letter eg: admin -> ADMIN
      --bak           enable scanning backup file (default:false)
      --split         enable spliting the url path, eg: /a/b -> /a/, /a/b (default: false)
  -X=                 method of http requests (default: GET)
  -x, --excode=       specify the status codes that be filtered eg: 400,404 (default:
                      400,404,406,416,501,502,503)
      --cookie=       set request cookies, eg: --cookie "session=admin"
  -H, --headers=      set request headers, string[] eg: -H "Token: admin=true" -H "Cookie:
                      login=true"
      --headers-file= set request headers file, eg: --headers-file headers.txt
      --proxy=        set request proxy, eg: --proxy http://127.0.0.1:8080
      --no-smart      disable smart mode (automated filtering)

Help Options:
  -h, --help          Show this help message
```

### 字典添加
可在 dicts 目录下根据个人需求更新常用字典
```
dirsx $ tree
.
├── dicts
│   ├── fuzzing-dirs-common.txt
│   ├── fuzzing-dirs-dirsearch.txt
│   ├── fuzzing-files-php.txt
│   ├── fuzzing-letter-len1.txt
│   ├── fuzzing-letter-len2.txt
│   ├── fuzzing-letter-len3.txt
│   ├── fuzzing-letter-len4.txt
│   ├── fuzzing-months-1-12.txt
│   ├── fuzzing-numbers-0-9.txt
│   ├── fuzzing-payloads-aspx.txt
│   ├── fuzzing-payloads-bakfile.txt
│   ├── fuzzing-payloads-common.txt
│   ├── fuzzing-payloads-java.txt
│   ├── fuzzing-payloads-null.txt
│   ├── fuzzing-payloads-php.txt
│   ├── fuzzing-routers-common.txt
│   ├── fuzzing-words-len1-5.txt
│   └── fuzzing-years-1990-2024.txt
├── dirsx
```

### dirsx 安装
根据对应系统类型下载执行文件 https://github.com/chasingboy/dirsx/releases

> ⚠️注意: 源代码中删除了部分还需要完善的代码，所以请不要使用源代码编译


> window 10 终端颜色显示问题, 可以更换系统终端为 window terminal 解决此问题 `https://github.com/microsoft/terminal`

### 公众号
该公众号用于编写 Xtools 系列工具使用文档和工具更新通知

<p align="center"><img width="300" alt="image" src="https://github.com/chasingboy/appsx/blob/main/assets/xsec.png"></p>

### 特别感谢
chainreactors@ https://github.com/chainreactors/spray

maurosoria@ https://github.com/maurosoria/dirsearch

ffuf@ https://github.com/ffuf/ffuf

### 更新记录
[+] 2024-09-21 修复 --split bug、修复 302 filter bug

[+] 2024-10-07 增加 cookie, header 功能、 增加 proxy 功能

[+] 2024-09-27 修复 302 filter 错误

[+] 2024-10-11 修复 tls handshake failure, basic 页面对比错误

[+] 2024-11-03 修复 title 特殊字符导致格式问题、Redirect 二次判断问题

[+] 2024-11-15 增加 --httpx, --ffuf 模式

[+] 2024-11-18 修复发生异常时，-o 没有输出结果的问题 #4

[+] 2024-12-17 修复 map error: concurrent map read and map write

[+] 2024-12-17 修改扫描结果实时打印, 增加进度条设置是否显示

[+] 2024-12-19 增加 @tony 师傅整理字典 fuzzing-payloads-vulnerability.txt

[+] 2024-12-20 修复 --no-smart 模式 30X 跳转丢包问题

[+] 2025-04-20 重构部分代码，修复 ffuf 模式 bug，修改字典目录 wordlist

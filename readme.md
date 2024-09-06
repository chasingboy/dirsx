<h1 align="center">dirsx</h1>
<h3 align="center">dirsx 是一款能够自动化过滤扫描结果的目录扫描工具</h3>
<p align="center">
  <img src="https://img.shields.io/badge/Version-V1.1.0-green?style=flat">
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
* 字典第一个字母大写|全部字母大写|添加前后缀
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
  -u, --url=         input url of target
  -l, --list=        input file containing list of target
  -w, --wordlist=    appoint wordlist for scanning directory
      --title-len=   set title display length (default: 30)
  -t, --threads=     number of threads to use (default: 20)
      --timeout=     timeout in seconds (default: 5)
  -o, --output=      file to write output results
      --prefix=      add prefix of payloads
      --suffix=      add suffix of payloads
  -e, --extension=   add extension eg: -e php,html
      --remove-ext=  remove extension eg: --remove-ext php | admin.php -> admin
      --upper-title  capitalize the first letter eg: admin -> Admin
      --upper-all    capitalize the all letter eg: admin -> ADMIN
      --bak          enable scanning backup file (default:false)
      --split        enable spliting the url path, eg: /a/b -> /a/, /a/b (default: false)
  -X=                method of http requests (default: GET)
  -x, --excode=      specify the status codes that be filtered eg: 400,404 (default: 400,404,406,416,501,502,503)
      --no-smart     disable smart mode (automated filtering)

Help Options:
  -h, --help         Show this help message

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
⚠️注意事项: 源代码中删除了部分还需要完善的代码，所以请不要使用源代码编译

### 特别感谢
chainreactors@ https://github.com/chainreactors/spray

maurosoria@ https://github.com/maurosoria/dirsearch

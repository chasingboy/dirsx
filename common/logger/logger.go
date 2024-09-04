package logger

import (
    "fmt"
    "strings"
)

// Terminal show color
type Logger struct {
    builder strings.Builder
}


func (logger *Logger) INF() *Logger {
    // color -> blue 
    logger.builder.WriteString("[\033[34mINF\033[0m] ")
    
    return logger
}

func (logger *Logger) WAR() *Logger {
    // color -> yellow
    logger.builder.WriteString("[\033[35mWAR\033[0m] ")
    
    return logger
}

func (logger *Logger) SUC() *Logger {
    // color -> yellow
    logger.builder.WriteString("[\033[32mWAR\033[0m] ")
    
    return logger
}

func (logger *Logger) RET() *Logger {
    // color -> blue
    logger.builder.WriteString("[\033[34mRET\033[0m] ")
    
    return logger
}


func (logger *Logger) ERR() *Logger {
    // color -> blue
    logger.builder.WriteString("[\033[31mERR\033[0m] ")
    
    return logger
}


func (logger *Logger) Inf(msg string) *Logger {
    // color -> yellow
    logger.builder.WriteString("\033[34m")
    logger.builder.WriteString(msg)
    logger.builder.WriteString("\033[0m")

    return logger
}

func (logger *Logger) War(msg string) *Logger {
    // color -> yellow
    logger.builder.WriteString("\033[35m")
    logger.builder.WriteString(msg)
    logger.builder.WriteString("\033[0m")

    return logger
}

func (logger *Logger) Suc(msg string) *Logger {
    // color -> green
    logger.builder.WriteString("\033[32m")
    logger.builder.WriteString(msg)
    logger.builder.WriteString("\033[0m")

    return logger
}

func (logger *Logger) Ret(msg string) *Logger {
    // color -> yellow
    logger.builder.WriteString("\033[36m")
    logger.builder.WriteString(msg)
    logger.builder.WriteString("\033[0m")

    return logger
}

func (logger *Logger) Err(msg string) *Logger {
    // color -> red
    logger.builder.WriteString("\033[31m")
    logger.builder.WriteString(msg)
    logger.builder.WriteString("\033[0m")

    return logger
}

func (logger *Logger) State(msg string, level string) *Logger {
    // color -> 
    logger.builder.WriteString("[\033")
    switch level {
    case "INF":
        logger.builder.WriteString("[34m")
    case "WAR":
        logger.builder.WriteString("[35m")
    case "SUC":
        logger.builder.WriteString("[32m")
    case "RET":
        logger.builder.WriteString("[36m")
    case "ERR":
        logger.builder.WriteString("[31m")
    default:
        logger.builder.WriteString("[34m")
    }
    logger.builder.WriteString(msg)
    logger.builder.WriteString("\033[0m]")

    return logger
}

func (logger *Logger) Str(msg string) *Logger {
    logger.builder.WriteString(msg)
    return logger
}

func (logger *Logger) Msg(msg string) string {
    logger.builder.WriteString(msg)
    text := logger.builder.String()
    logger.builder.Reset()
    
    return text
}

func (logger *Logger) Msgf(msg string) string {
    logger.builder.WriteString(msg)
    text := logger.builder.String()
    logger.builder.Reset()
    fmt.Println(text)
    
    return text
}

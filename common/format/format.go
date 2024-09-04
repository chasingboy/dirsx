package format

import (
    "fmt"
    "strings"
)

// Format string
func F(x string, args ...interface{}) string {
    var value string
    
    for index, arg := range args {
        switch arg.(type) {
        case int:
            value = fmt.Sprintf("%d", arg)
        default:
            value = fmt.Sprintf("%s", arg)
        }
        
        // fmt.Println(index, value)
        x = strings.Replace(x, fmt.Sprintf("{%d}",index), value, -1)
    }
    
    return x
}

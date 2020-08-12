# ICSReader
Minimalist ICal/ICS Reader

## How to use it ?
```
$ go get github.com/robinjulien/icsreader
```

```go
package main

import "github.com/lepompiste/icsreader"

func main() {
    events, err := icsreader.GetCalendarFromURL("http://example.com/myCalendar.ics")

    if err == nil {
        println("Error")
    } else {
        var json string = events.ParseJSON()
        println(json)
    }
}
```
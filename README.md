# go-linenotify

Golang client for the Line Notify API.

- https://notify-bot.line.me/

## Installation

```
$ go get github.com/timwata/go-linenotify
```

## Example

```go
package main

import (
    "fmt"

    "github.com/timwata/go-linenotify"  
)

func main() {
    cli := linenotify.New("YOUR_TOKEN")
    err := cli.Post("Hello, World!", nil)
    if err != nil {
        fmt.Println(err)
    }
}
```

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).

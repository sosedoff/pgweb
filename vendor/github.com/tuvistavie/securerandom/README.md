# securerandom

Port of Ruby `securerandom` module for Golang.
The following functions are implemented.

```go
func RandomBytes(n int) ([]byte, error)
func Base64(n int, padded bool) (string, error)
func UrlSafeBase64(n int, padded bool) (string, error)
func Hex(n int) (string, error)
func Uuid() (string, error)
```

## Sample usage

```go
package main

import (
  "fmt"
  sr "github.com/tuvistavie/securerandom"
)

func main() {
  b, _ := sr.Base64(10, true)
  fmt.Println(b)
  b, _ = sr.Hex(10)
  fmt.Println(b)
  b, _ = sr.Uuid()
  fmt.Println(b)
}
```

For more information, check out the [documentation of the ruby module](http://ruby-doc.org/stdlib-2.1.0/libdoc/securerandom/rdoc/SecureRandom.html).

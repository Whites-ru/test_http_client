package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

func main() {

	urls := []string{"https://pkg.go.dev/std",
		"https://go.dev/",
		"https://go.dev/blog/",
		"https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go",
		"https://www.digitalocean.com/community/tutorials/how-to-use-vuls-as-a-vulnerability-scanner-on-ubuntu-18-04",
		"https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2"}

	n := 0
	urls_count := float64(len(urls))
	var words_count atomic.Int64

	done := make(chan bool)

	for i := 0; i < int(math.Ceil(urls_count/5.0)); i++ {

		var urls_s []string
		if n > int(urls_count) {
			urls_s = urls[n:]
		} else {
			if (n + 5) > int(urls_count) {
				urls_s = urls[n:]
			} else {
				urls_s = urls[n : n+5]
			}
		}

		for i := 0; i < len(urls_s); i++ {

			go func(url string) {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)

				} else {
					if resp.StatusCode != 200 {
						fmt.Printf("Response failed with status code: %d\n", resp.StatusCode)
					} else {
						if err != nil {
							fmt.Println(err)
						} else {
							defer resp.Body.Close()
							body, err := io.ReadAll(resp.Body)

							if err != nil {
								fmt.Println(err)
							} else {
								count := int64(strings.Count(string(body), "Go"))
								fmt.Println("Count for " + url + ": " + strconv.FormatInt(count, 10))
								words_count.Store(words_count.Load() + count)
							}
						}
					}
				}

				done <- true
			}(urls_s[i])
		}
		for range urls_s {
			<-done
		}

		n = n + 5
	}
	fmt.Println("Total: " + strconv.FormatInt(words_count.Load(), 10))
}

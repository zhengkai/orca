package tmptest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"project/zj"
)

// Test ...
func Test() {
	fetch()
}

func fetch() {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://10.0.84.49/e500", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	zj.J(`StatusCode`, resp.StatusCode)
	fmt.Printf("%s\n", bodyText)

}

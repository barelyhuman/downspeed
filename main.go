package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/barelyhuman/go/term"
	"github.com/dustin/go-humanize"
)

type speedTest struct {
	token     string
	tokenUrl  string
	apiUrl    string
	serverUrl string
	requests  int
}

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	term.HideCursor()
	defer term.ShowCursor()

	st := &speedTest{
		tokenUrl: "https://fast.com/app-ed402d.js",
		apiUrl:   "https://api.fast.com/netflix/speedtest",
		requests: 100,
	}
	err := st.getToken()
	bail(err)

	// 2kb
	byteSizeOne := 2 * 1024

	// ~ 25MB
	byteSizeTwo := 25 * 1024 * 1024

	st.fetchServers()

	doneChannel := make(chan bool)
	// warm the network up!
	go st.fetchByRange("0", strconv.Itoa(byteSizeOne), doneChannel, &[]float64{})
	<-doneChannel
	go st.fetchByRange("0", strconv.Itoa(byteSizeTwo), doneChannel, &[]float64{})
	<-doneChannel
	speeds := &[]float64{}
	for i := 0; i <= st.requests; i++ {
		go st.fetchByRange("0", strconv.Itoa(byteSizeTwo), doneChannel, speeds)
		<-doneChannel
		term.ClearAll()
		term.MoveTo(1, 1)
		fmt.Printf("%v/s", getAvg(*speeds))
	}
}

func getAvg(speeds []float64) string {
	total := 0.0
	for _, spd := range speeds {
		total += spd
	}
	return humanize.Bytes(uint64(total / float64(len(speeds))))
}

func (st *speedTest) getToken() (err error) {
	tokenRe := regexp.MustCompile("token:\"[^\"]+")
	res, err := http.Get(st.tokenUrl)
	bail(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	bail(err)
	matched := tokenRe.FindStringSubmatch(string(body))
	tokenParts := strings.Split(matched[0], "\"")
	st.token = tokenParts[1]
	return
}

func (st *speedTest) fetchServers() (err error) {
	urlRe := regexp.MustCompile("https[^\"]+")
	url := st.apiUrl + "?https=true&token=" + st.token
	res, err := http.Get(url)
	bail(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	bail(err)

	match := urlRe.FindStringSubmatch(string(body))
	st.serverUrl = match[0]
	return
}

func (st speedTest) fetchByRange(rangeStart string, rangeEnd string, done chan bool, speeds *[]float64) (speed float64) {
	var url strings.Builder
	url.WriteString(strings.Replace(st.serverUrl, "/speedtest", "/speedtest/range/"+rangeStart+"-"+rangeEnd, 1))
	speed = GetSpeed(url.String())
	done <- true
	*speeds = append(*speeds, speed)
	return
}

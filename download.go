package main

import (
	"fmt"
	"net/http"
	"time"
)

func GetSpeed(url string) float64 {
	client := http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("GET", url, nil)
	bail(err)
	req.Header.Set("Referer", "https://fast.com/")
	req.Header.Set("Origin", "https://fast.com/")
	start := time.Now().UnixMilli()
	res, err := client.Do(req)
	bail(err)
	if res.StatusCode != 200 {
		err = fmt.Errorf("failed to download file")
		bail(err)
	}
	end := time.Now().UnixMilli()
	diff := float64(end) - float64(start)
	sizeInBytes := res.ContentLength
	inSeconds := diff / 10
	speedInBytes := float64(sizeInBytes) / inSeconds

	return speedInBytes
}

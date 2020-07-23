package main

import (
	"log"
	"sync"
	"time"

	"github.com/fengyfei/moon/models"
	"github.com/fengyfei/moon/service/stock"

	_ "github.com/mattn/go-sqlite3"
)

const (
	workers   = 10
	intervals = 10000
	stop      = "stop"
)

var (
	dateCh = make(chan string, workers)
	lock   = sync.Mutex{}
	wg     = sync.WaitGroup{}
)

func main() {
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFunc()
	}

	now := time.Now()
	for i := 0; i < intervals; i++ {
		t := now.AddDate(0, 0, -i)
		date := t.Format("2006-01-02")

		dateCh <- date
	}

	for i := 0; i < workers; i++ {
		dateCh <- stop
	}

	wg.Wait()
}

func workerFunc() {
	for {
		date := <-dateCh

		if date == stop {
			wg.Done()
			return
		}

		resp, err := stock.GetNorthDailyReport(date)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println(resp)

		lock.Lock()
		if err = models.Record(date, resp); err != nil {
			log.Printf("[North] Record error(%#v) for %s", err, date)
		}
		lock.Unlock()

		time.Sleep(2 * time.Second)
	}
}

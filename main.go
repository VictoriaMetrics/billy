package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	startDateStr = flag.String("startdate", "2019-01-01", "date to start sweep YYYY-MM-DD")
	endDateStr   = flag.String("enddate", "2019-01-31", "date to end sweep YYYY-MM-DD")
	startKey     = flag.Int("startkey", 1, "first sensor ID")
	endKey       = flag.Int("endkey", 2, "last sensor ID")
	workers      = flag.Int("workers", runtime.GOMAXPROCS(-1), "the number of concurrent workers used for data ingestion")
	sink         = flag.String("sink", "http://localhost:8428/api/v1/import", "HTTP address for the data ingestion sink")
	compress     = flag.Bool("compress", false, "Whether to compress data before sending it to sink. This saves network bandwidth at the cost of higher CPU usage")
	digits       = flag.Int("digits", 5, "The number of decimal digits after the point in the generated temperature")
)

func main() {
	flag.Parse()

	startTimestamp := mustParseDate(*startDateStr, "startdate")
	endTimestamp := mustParseDate(*endDateStr, "enddate")
	if startTimestamp > endTimestamp {
		log.Fatalf("-startdate=%s cannot exceed -enddate=%s", *startDateStr, *endDateStr)
	}
	endTimestamp += 24 * 3600 * 1000
	rowsCount := int((endTimestamp - startTimestamp) / (60 * 1000))
	if *startKey > *endKey {
		log.Fatalf("-startkey=%d cannot exceed -endkey=%d", *startKey, *endKey)
	}

	workCh := make(chan work)
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(workCh)
		}()
	}
	go statsReporter()
	keysCount := *endKey - *startKey + 1
	startTime = time.Now()
	rowsTotal = rowsCount * keysCount
	for startTimestamp < endTimestamp {
		for key := *startKey; key <= *endKey; key++ {
			w := work{
				key:            key,
				startTimestamp: startTimestamp,
				rowsCount:      24 * 60,
			}
			workCh <- w
		}
		startTimestamp += 24 * 3600 * 1000
	}
	close(workCh)
	wg.Wait()
}

var rowsTotal int
var rowsGenerated uint64
var startTime time.Time

func statsReporter() {
	for {
		time.Sleep(10 * time.Second)
		d := time.Since(startTime).Seconds()
		n := atomic.LoadUint64(&rowsGenerated)
		log.Printf("generated %d out of %d rows at %.0f rows/sec", n, rowsTotal, float64(n)/d)
	}
}

type work struct {
	key            int
	startTimestamp int64
	rowsCount      int
}

func worker(workCh <-chan work) {
	pr, pw := io.Pipe()
	req, err := http.NewRequest("POST", *sink, pr)
	if err != nil {
		log.Fatalf("cannot create request to %q: %s", *sink, err)
	}
	w := io.Writer(pw)
	if *compress {
		zw, err := gzip.NewWriterLevel(pw, 1)
		if err != nil {
			log.Fatalf("unexpected error when creating gzip writer: %s", err)
		}
		w = zw
		req.Header.Set("Content-Encoding", "gzip")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("unexpected error when performing request to %q: %s", *sink, err)
		}
		if resp.StatusCode != http.StatusNoContent {
			log.Fatalf("unexpected response code from %q: %s", *sink, err)
		}
	}()
	bw := bufio.NewWriterSize(w, 16*1024)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for w := range workCh {
		writeSeries(bw, r, w.key, w.rowsCount, w.startTimestamp)
		atomic.AddUint64(&rowsGenerated, uint64(w.rowsCount))
	}
	_ = bw.Flush()
	if *compress {
		_ = w.(*gzip.Writer).Close()
	}
	_ = pw.Close()
	wg.Wait()
}

func writeSeries(bw *bufio.Writer, r *rand.Rand, sensorID, rowsCount int, startTimestamp int64) {
	min := 68 + r.ExpFloat64()/3.0
	e := math.Pow10(*digits)
	fmt.Fprintf(bw, `{"metric":{"__name__":"temperature","sensor_id":"%d"},"values":[`, sensorID)
	for i := 0; i < rowsCount-1; i++ {
		t := generateTemperature(r, min, e)
		fmt.Fprintf(bw, "%g,", t)
	}
	t := generateTemperature(r, min, e)
	fmt.Fprintf(bw, `%g],"timestamps":[`, t)
	for i := 0; i < rowsCount-1; i++ {
		fmt.Fprintf(bw, "%d,", startTimestamp+int64(i)*60*1000)
	}
	fmt.Fprintf(bw, "%d]}\n", startTimestamp+int64(rowsCount-1)*60*1000)
}

func generateTemperature(r *rand.Rand, min, e float64) float64 {
	t := r.ExpFloat64()/1.5 + min
	return math.Round(t*e) / e
}

func mustParseDate(dateStr, flagName string) int64 {
	startTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("cannot parse -%s=%q: %s", flagName, dateStr, err)
	}
	return startTime.UnixNano() / 1e6
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klauspost/compress/gzip"
)

var (
	startDateStr     = flag.String("startdate", "2019-01-01", "Date to start sweep YYYY-MM-DD")
	endDateStr       = flag.String("enddate", "2019-01-31", "Date to end sweep YYYY-MM-DD")
	startKey         = flag.Int("startkey", 1, "First sensor ID")
	endKey           = flag.Int("endkey", 2, "Last sensor ID")
	workers          = flag.Int("workers", runtime.GOMAXPROCS(-1), "The number of concurrent workers used for data ingestion")
	sink             = flag.String("sink", "http://localhost:8428/api/v1/import", "HTTP address for the data ingestion sink. It depends on the `-format`")
	compress         = flag.Bool("compress", false, "Whether to compress data before sending it to sink. This saves network bandwidth at the cost of higher CPU usage")
	digits           = flag.Int("digits", 2, "The number of decimal digits after the point in the generated temperature. The original benchmark from ScyllaDB uses 2 decimal digits after the point. See query results at https://www.scylladb.com/2019/12/12/how-scylla-scaled-to-one-billion-rows-a-second/")
	reportInterval   = flag.Duration("report-interval", 10*time.Second, "Stats reporting interval")
	format           = flag.String("format", "vmimport", "Data ingestion format. Supported values: vmimport, influx")
	blocksPerRequest = flag.Int("blocks-per-request", 0, "The maximum number of blocks per request. Unlimited if set to 0. "+
		"This can be used for ingesting data into InfluxDB, which doesn't support request body streaming")
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
	var workersWg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			worker(workCh)
		}()
	}
	statsReporterStopCh := make(chan struct{})
	var statsReporterWG sync.WaitGroup
	statsReporterWG.Add(1)
	go func() {
		defer statsReporterWG.Done()
		statsReporter(statsReporterStopCh)
	}()
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
	workersWg.Wait()

	close(statsReporterStopCh)
	statsReporterWG.Wait()
}

var rowsTotal int
var rowsGenerated uint64
var startTime time.Time

func statsReporter(stopCh <-chan struct{}) {
	prevTime := time.Now()
	nPrev := uint64(0)
	ticker := time.NewTicker(*reportInterval)
	mustStop := false
	for !mustStop {
		select {
		case <-ticker.C:
		case <-stopCh:
			mustStop = true
		}
		t := time.Now()
		dAll := t.Sub(startTime).Seconds()
		dLast := t.Sub(prevTime).Seconds()
		nAll := atomic.LoadUint64(&rowsGenerated)
		nLast := nAll - nPrev
		log.Printf("created %d out of %d rows in %.3f seconds at %.0f rows/sec; instant speed %.0f rows/sec",
			nAll, rowsTotal, dAll, float64(nAll)/dAll, float64(nLast)/dLast)
		prevTime = t
		nPrev = nAll
	}
}

type work struct {
	key            int
	startTimestamp int64
	rowsCount      int
}

func (w *work) do(bw *bufio.Writer, r *rand.Rand) {
	switch *format {
	case "vmimport":
		writeSeriesVMImport(bw, r, w.key, w.rowsCount, w.startTimestamp)
	case "influx":
		writeSeriesInflux(bw, r, w.key, w.rowsCount, w.startTimestamp)
	default:
		log.Fatalf("unexpected `-format=%q`. Supported values: vmimport, influx", *format)
	}
	atomic.AddUint64(&rowsGenerated, uint64(w.rowsCount))
}

func worker(workCh <-chan work) {
	for w := range workCh {
		workerSingleRequest(workCh, w)
	}
}

func workerSingleRequest(workCh <-chan work, wk work) {
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
			log.Printf("unexpected response code from %q: %d", *sink, resp.StatusCode)
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("cannot read response body: %s", err)
			}
			log.Fatalf("response body:\n%s", data)
		}
	}()
	bw := bufio.NewWriterSize(w, 16*1024)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	blocks := 0
	ok := true
	for ok {
		wk.do(bw, r)
		blocks++
		if *blocksPerRequest > 0 && blocks >= *blocksPerRequest {
			break
		}
		wk, ok = <-workCh
	}
	_ = bw.Flush()
	if *compress {
		_ = w.(*gzip.Writer).Close()
	}
	_ = pw.Close()
	wg.Wait()
}

func writeSeriesVMImport(bw *bufio.Writer, r *rand.Rand, sensorID, rowsCount int, startTimestamp int64) {
	min := 68 + r.ExpFloat64()/3.0
	e := math.Pow10(*digits)
	fmt.Fprintf(bw, `{"metric":{"__name__":"temperature","sensor_id":"%d"},"values":[`, sensorID)
	var buf []byte
	t := generateTemperature(r, min, e)
	for i := 0; i < rowsCount-1; i++ {
		buf = strconv.AppendFloat(buf[:0], t, 'f', *digits, 64)
		buf = append(buf, ',')
		bw.Write(buf)
		t = generateTemperature(r, min, e)
	}
	fmt.Fprintf(bw, `%.*f],"timestamps":[`, *digits, t)
	timestamp := startTimestamp
	for i := 0; i < rowsCount-1; i++ {
		buf = strconv.AppendInt(buf[:0], timestamp, 10)
		buf = append(buf, ',')
		bw.Write(buf)
		timestamp = startTimestamp + int64(i+1)*60*1000
	}
	fmt.Fprintf(bw, "%d]}\n", timestamp)
}

func writeSeriesInflux(bw *bufio.Writer, r *rand.Rand, sensorID, rowsCount int, startTimestamp int64) {
	min := 68 + r.ExpFloat64()/3.0
	e := math.Pow10(*digits)
	var buf []byte
	for i := 0; i < rowsCount; i++ {
		t := generateTemperature(r, min, e)
		timestamp := (startTimestamp + int64(i)*60*1000) * 1e6
		buf = append(buf[:0], "temperature,sensor_id="...)
		buf = strconv.AppendInt(buf, int64(sensorID), 10)
		buf = append(buf, " value="...)
		buf = strconv.AppendFloat(buf, t, 'f', *digits, 64)
		buf = append(buf, ' ')
		buf = strconv.AppendInt(buf, timestamp, 10)
		buf = append(buf, '\n')
		bw.Write(buf)
	}
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

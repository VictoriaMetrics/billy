package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkWriteSeries(b *testing.B) {
	const rowsCount = 24 * 60
	const loopsCount = 10
	b.ReportAllocs()
	b.SetBytes(rowsCount * loopsCount)
	b.RunParallel(func(pb *testing.PB) {
		startTimestamp := time.Now().UnixNano() / 1e6
		bw := bufio.NewWriter(ioutil.Discard)
		r := rand.New(rand.NewSource(startTimestamp))
		sensorID := int(startTimestamp) % 1e6
		for pb.Next() {
			for i := 0; i < loopsCount; i++ {
				writeSeries(bw, r, sensorID, rowsCount, startTimestamp)
			}
			if err := bw.Flush(); err != nil {
				panic(fmt.Errorf("unexpected error on bufio.Writer.Flush: %s", err))
			}
		}
	})
}

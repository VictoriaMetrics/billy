package main

import (
	"bufio"
	"bytes"
	"math/rand"
	"testing"
)

func TestWriteSeries(t *testing.T) {
	startTimestamp := int64(1234)
	var bb bytes.Buffer
	bw := bufio.NewWriter(&bb)
	r := rand.New(rand.NewSource(startTimestamp))
	sensorID := 789
	rowsCount := 3
	writeSeries(bw, r, sensorID, rowsCount, startTimestamp)
	if err := bw.Flush(); err != nil {
		t.Fatalf("unexpected error in bw.Flush: %s", err)
	}
	result := bb.String()
	resultExpected := `{"metric":{"__name__":"temperature","sensor_id":"789"},"values":[68.34,69.61,69.86],"timestamps":[1234,61234,121234]}`+"\n"
	if result != resultExpected {
		t.Fatalf("unexpected result;\ngot\n%s\nwant\n%s", result, resultExpected)
	}
}
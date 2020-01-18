package main

import (
	"bufio"
	"bytes"
	"math/rand"
	"testing"
)

func TestWriteSeriesVMImport(t *testing.T) {
	startTimestamp := int64(1234)
	var bb bytes.Buffer
	bw := bufio.NewWriter(&bb)
	r := rand.New(rand.NewSource(startTimestamp))
	sensorID := 789
	rowsCount := 3
	writeSeriesVMImport(bw, r, sensorID, rowsCount, startTimestamp)
	if err := bw.Flush(); err != nil {
		t.Fatalf("unexpected error in bw.Flush: %s", err)
	}
	result := bb.String()
	resultExpected := `{"metric":{"__name__":"temperature","sensor_id":"789"},"values":[68.34,69.61,69.86],"timestamps":[1234,61234,121234]}` + "\n"
	if result != resultExpected {
		t.Fatalf("unexpected result;\ngot\n%s\nwant\n%s", result, resultExpected)
	}
}

func TestWriteSeriesInflux(t *testing.T) {
	startTimestamp := int64(1234)
	var bb bytes.Buffer
	bw := bufio.NewWriter(&bb)
	r := rand.New(rand.NewSource(startTimestamp))
	sensorID := 789
	rowsCount := 3
	writeSeriesInflux(bw, r, sensorID, rowsCount, startTimestamp)
	if err := bw.Flush(); err != nil {
		t.Fatalf("unexpected error in bw.Flush: %s", err)
	}
	result := bb.String()
	resultExpected := `temperature value=68.34 1234000000
temperature value=69.61 61234000000
temperature value=69.86 121234000000
`
	if result != resultExpected {
		t.Fatalf("unexpected result;\ngot\n%s\nwant\n%s", result, resultExpected)
	}
}

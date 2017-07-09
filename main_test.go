package main

import (
	"fmt"
	"testing"
)

func TestGetConsignments(t *testing.T) {
	path := "/home/erlend/go/src/github.com/erlfas/SplitTransportJob/TransportJob-big-20170706.xml"
	consignments := getConsignments(path)
	numConsignments := len(consignments)
	if numConsignments != 17878 {
		t.Error("Expected 17878 consignments but found ", numConsignments)
	}
	header := getHeader(path)
	batchSize := 50
	batchedConsignments := groupConsignments(header, consignments, batchSize)
	numBatchedConsignments := len(batchedConsignments)
	if numBatchedConsignments != 358 {
		t.Error("Expected 358 batches but found ", numBatchedConsignments)
	}
	fmt.Println(batchedConsignments[0])
}

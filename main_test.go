package main

import (
	"os"
	"testing"
)

func TestReadLogFile(t *testing.T) {
	// Create a temporary file with some sample log entries
	tmpfile, err := os.CreateTemp("", "example.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString("127.0.0.1 - - [20/Apr/2022:12:34:56 +0000] \"GET / HTTP/1.1\" 200\n")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Create allVisitors and uniqueVisitors structures
	allVisitors := allVisitors{visitors: new(visitors), yearly: make(map[int]*yearlyVisitors)}
	uniqueVisitors := make(map[string]struct{})

	// Call the function being tested
	readLogFile(tmpfile.Name(), &allVisitors, uniqueVisitors)

	// Perform assertions based on the expected behavior
	expectedNew := 1
	expectedOld := 0

	if allVisitors.visitors.new != expectedNew || allVisitors.visitors.old != expectedOld {
		t.Errorf("Expected new: %d, old: %d, got new: %d, old: %d", expectedNew, expectedOld, allVisitors.visitors.new, allVisitors.visitors.old)
	}

	// Clean up - remove the temporary file
	os.Remove(tmpfile.Name())
}

func TestAllVisitors_GetNew(t *testing.T) {
	// Create an example allVisitors structure for testing
	allVisitors := allVisitors{visitors: &visitors{new: 10, old: 5}, yearly: make(map[int]*yearlyVisitors)}
	allVisitors.yearly[2022] = &yearlyVisitors{visitors: &visitors{new: 5, old: 2}}

	// Call the function being tested
	//err, result := allVisitors.GetNew("20/04/2022")
	err, result := allVisitors.GetNew("2022")
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions based on the expected behavior
	expectedResult := 5

	if result != expectedResult {
		t.Errorf("Expected result: %d, got result: %d", expectedResult, result)
	}
}

package store

import (
	"log"
	"os"
	"testing"
)

var testStore Store

const path = "../.."

func TestMain(m *testing.M) {
	var err error
	testStore, err = NewFileStore(path)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	os.Exit(m.Run())
}

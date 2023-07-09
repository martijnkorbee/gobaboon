package cache

import (
	"fmt"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
)

var (
	testRedisClient *miniredis.Miniredis
	testBadgerPath  string = "./testdata/badger"
)

func TestMain(m *testing.M) {
	// setup redis test client
	if mr, err := miniredis.Run(); err != nil {
		panic(err)
	} else {
		testRedisClient = mr
		defer mr.Close()
	}

	// defere remove ./testdata to cleanup badger models
	defer func() {
		if err := os.RemoveAll("./testdata"); err != nil {
			fmt.Println("failed to remove test models:", err)
		}
	}()

	m.Run()
}

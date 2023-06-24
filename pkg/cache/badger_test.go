package cache

import "testing"

func Test_badger_CreateBadgerCache(t *testing.T) {
	_, err := CreateBadgerCache(BadgerConfig{
		Prefix: "test-baboon",
		Path:   testBadgerPath,
	})
	if err != nil {
		t.Error("failed to connect to badger:", err.Error())
	}
}

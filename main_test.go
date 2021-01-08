// main_test.go

package main

import (
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize("test/data.pdf")

	code := m.Run()
	os.Exit(code)
}

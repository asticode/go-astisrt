package astisrt

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var code int
	defer func(code *int) {
		os.Exit(*code)
	}(&code)

	if err := Startup(); err != nil {
		log.Fatal(fmt.Errorf("main: starting up failed: %w", err))
	}

	defer func() {
		if err := CleanUp(); err != nil {
			log.Fatal(fmt.Errorf("main: cleaning up failed: %w", err))
		}
	}()

	code = m.Run()
}

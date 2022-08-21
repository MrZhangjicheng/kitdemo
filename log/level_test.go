package log

import (
	"log"
	"testing"
)

func TestLevel(t *testing.T) {
	log.Fatal(LevelDebug)
}

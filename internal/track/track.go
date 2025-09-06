package track

import (
	"fmt"
	"time"
)

func Measure(start time.Time, name string, count int) {
	if count > 1 {
		elapsed := time.Since(start)
		fmt.Printf("%s", elapsed/time.Duration(count))
	}
}

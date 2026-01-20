package runner

import (
	"fmt"
	"testing"
)

func BenchmarkRunLinterChecksInParallel(b *testing.B) {
	for cnt := 1; cnt <= len(Checks); cnt++ {
		n := cnt
		b.Run(fmt.Sprintf("Running with %d checks", n), func(b *testing.B) {
			for b.Loop() {
				RunLinterChecks("./../test", Checks[:n], 0, 0, true)
			}
		})
	}

}

func BenchmarkRunLinterChecksSequential(b *testing.B) {
	for cnt := 1; cnt <= len(Checks); cnt++ {
		n := cnt
		b.Run(fmt.Sprintf("Running with %d checks", n), func(b *testing.B) {
			for b.Loop() {
				RunLinterChecks("./../test", Checks[:n], 0, 0, false)
			}
		})
	}
}

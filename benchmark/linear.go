package benchmark

import (
	"fmt"
	"time"

	"github.com/MariusVanDerWijden/FuzzyVM/executor"
	"github.com/MariusVanDerWijden/FuzzyVM/generator"
	"github.com/holiman/goevmlab/evms"
)

// testGeneration generates N programs.
func testGeneration(N int) (time.Duration, error) {
	f, err := newFiller()
	if err != nil {
		return time.Nanosecond, err
	}
	start := time.Now()
	for i := 0; i < N; i++ {
		generator.GenerateProgram(f)
	}
	return time.Since(start), nil
}

// verify verifies a programs result N times.
func verify(N int) (time.Duration, error) {
	outDir, _, err := createTempDirs()
	if err != nil {
		return time.Nanosecond, err
	}
	name := "BenchTest"
	if err := generateTest(name, outDir); err != nil {
		return time.Nanosecond, err
	}
	name = fmt.Sprintf("%v/%v.json", outDir, name)
	out, err := executor.ExecuteTest(name)
	if err != nil {
		return time.Nanosecond, err
	}
	start := time.Now()
	for i := 0; i < N; i++ {
		if !executor.Verify(name, out) {
			return time.Nanosecond, fmt.Errorf("Verification failed: %v", name)
		}
	}
	return time.Since(start), nil
}

// execution executes a program N times.
func execution(N int) (time.Duration, error) {
	outDir, crashers, err := createTempDirs()
	if err != nil {
		return time.Nanosecond, err
	}
	name := "BenchTest"
	if err := generateTest(name, outDir); err != nil {
		return time.Nanosecond, err
	}
	name = fmt.Sprintf("%v.json", name)
	executor.PrintTrace = false
	start := time.Now()
	for i := 0; i < N; i++ {
		executor.ExecuteFullTest(outDir, crashers, name, false)
	}
	return time.Since(start), nil
}

// linear runs N tests in sequence.
func linear(N int) (time.Duration, error) {
	evms.Docker = false
	return execLinearMultiple(N, false)
}

// linearBatch runs a batch of N tests.
func linearBatch(N int) (time.Duration, error) {
	evms.Docker = false
	return execLinearMultiple(N, true)
}

// linearDocker runs N tests in sequence on a docker container.
func linearDocker(N int) (time.Duration, error) {
	evms.Docker = true
	return execLinearMultiple(N, false)
}

// linearBatchDocker runs a batch of N tests on a docker container.
func linearBatchDocker(N int) (time.Duration, error) {
	evms.Docker = true
	return execLinearMultiple(N, true)
}

func execLinearMultiple(N int, batch bool) (time.Duration, error) {
	outDir, crashers, err := createTempDirs()
	if err != nil {
		return time.Nanosecond, err
	}
	var names []string
	for i := 0; i < N; i++ {
		name := fmt.Sprintf("BenchTest-%v", i)
		if err := generateTest(name, outDir); err != nil {
			return time.Nanosecond, err
		}
		names = append(names, fmt.Sprintf("%v.json", name))
	}
	executor.PrintTrace = false
	start := time.Now()
	if batch {
		if err := executor.ExecuteFullBatch(outDir, crashers, names, false); err != nil {
			return time.Nanosecond, err
		}
	} else {
		for _, n := range names {
			if err := executor.ExecuteFullTest(outDir, crashers, n, false); err != nil {
				return time.Nanosecond, err
			}
		}
	}
	return time.Since(start), nil
}
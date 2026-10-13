// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/dirodriguezm/xmatch/service/docs"
	"io"
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
	args []string,
) error {
	slog.Info("Starting xmatch service")

	if len(args) < 2 {
		panic("run: Missing arguments")
	}
	fs := flag.NewFlagSet("xmatch", flag.ExitOnError)

	var profile bool
	fs.BoolVar(&profile, "profile", false, "Enable profiling")

	if err := fs.Parse(args[2:]); err != nil {
		return err
	}

	command := args[1]

	if profile {
		slog.Info("Profiling enabled")

		// CPU profiling
		cpuFile, err := os.Create("cpu.prof")
		if err != nil {
			slog.Error("could not create CPU profile: ", "error", err)
		}
		pprof.StartCPUProfile(cpuFile)
		defer func() {
			pprof.StopCPUProfile()
			cpuFile.Close()
		}()

		// Memory profiling
		defer func() {
			memFile, err := os.Create("mem.prof")
			if err != nil {
				slog.Error("could not create memory profile: ", "error", err)
			}
			defer memFile.Close()

			// Run a GC to get up-to-date memory statistics
			runtime.GC()

			if err := pprof.WriteHeapProfile(memFile); err != nil {
				slog.Error("could not write memory profile: ", "error", err)
			}
		}()
	}

	switch command {
	case "server":
		return StartHttpServer(ctx, getenv, stdout)
	case "indexer":
		err := StartCatalogIndexer(ctx, getenv, stdout)
		return err
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

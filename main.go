package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/yusaer/orm-benchmark/benchs"

	_ "github.com/lib/pq"
)

type ListOpts []string

func (opts *ListOpts) String() string {
	return fmt.Sprint(*opts)
}

func (opts *ListOpts) Set(value string) error {
	if value == "all" || strings.Index(" "+strings.Join(benchs.BrandNames, " ")+" ", " "+value+" ") != -1 {
	} else {
		return fmt.Errorf("wrong run name %s", value)
	}
	*opts = append(*opts, value)
	return nil
}

func (opts ListOpts) Shuffle() {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(opts); i++ {
		a := rd.Intn(len(opts))
		b := rd.Intn(len(opts))
		opts[a], opts[b] = opts[b], opts[a]
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Setenv("PGSSLMODE", "disable")

	var (
		orms                 ListOpts
		enableCPU, enableMem bool
	)

	flag.IntVar(&benchs.ORM_MAX_IDLE, "max_idle", 200, "max idle conns")
	flag.IntVar(&benchs.ORM_MAX_CONN, "max_conn", 200, "max open conns")
	flag.StringVar(&benchs.ORM_SOURCE, "source", "postgres://postgres:postgres@localhost:5432/test?sslmode=disable", "postgres dsn source")
	flag.IntVar(&benchs.ORM_MULTI, "multi", 1, "base query nums x multi")
	flag.Var(&orms, "orm", "orm name: all, "+strings.Join(benchs.BrandNames, ", "))
	flag.BoolVar(&enableCPU, "cpu", false, "enable cpu profile")
	flag.BoolVar(&enableMem, "mem", false, "enable mem profile")

	flag.Parse()

	if enableCPU {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if enableMem {
		f, err := os.Create("mem.pprof")
		if err != nil {
			panic(err)
		}
		defer func() {
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

	var all bool

	if len(orms) == 0 {
		all = true
	} else {
		for _, n := range orms {
			if n == "all" {
				all = true
			}
		}
	}

	if all {
		orms = ListOpts(benchs.BrandNames)
	}

	orms.Shuffle()

	for _, n := range orms {
		fmt.Println(n)
		benchs.RunBenchmark(n)
	}

	fmt.Println("\nReports: ")
	fmt.Print(benchs.MakeReport())

}

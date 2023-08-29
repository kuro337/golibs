package profiling

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"time"
)

type Profiler struct {
	profileOutputPath string
	memFile           string
	cpuFile           string
	traceFile         string

	mem      bool
	cpu      bool
	trace    bool
	helpFlag bool

	linger   bool
	traceOut *os.File
	cpuOut   *os.File
	memOut   *os.File
}

/*
NewProfiler initializes a new Profiler instance.

The provided profileOutputPath specifies the directory where profiling data will be saved.

The Profiler can be initalized at launch at the beginning of the application in main().

Example usage:

1. Enabling Memory, CPU , and Tracing Profiling

	p := profiling.NewProfiler("outputFolder").Memory().CPU().Tracing().Help().Start()
	defer p.Stop()

2. Enabling Only Memory and CPU Profiling

	p := profiling.NewProfiler("outputFolder").Memory().CPU().Start()
	defer p.Stop()

NOTE : The profiler will linger for 5s after the program ends to collect memory profile data.

Use NoLinger() to disable this behavior.

	profiling.NewProfiler("outputFolder").Memory().NoLinger().Start()

	func main() {
		profiler := profiling.NewProfiler("outputFolder").Memory().NoLinger().Start()
		defer profiler.Stop()

		...rest of app code

	}
*/
func NewProfiler(profileOutputPath string) *Profiler {
	// Create a folder for profiling data
	profileDir := profileOutputPath
	if err := os.MkdirAll(profileDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return &Profiler{
		profileOutputPath: profileOutputPath,
		linger:            true,
		helpFlag:          false,
		traceFile:         fmt.Sprintf("%s/trace.out", profileOutputPath),
		cpuFile:           fmt.Sprintf("%s/cpu.pprof", profileOutputPath),
		memFile:           fmt.Sprintf("%s/mem.pprof", profileOutputPath),
	}
}

/*
Memory enables memory profiling and saves the mem.pprof file to provided output path to NewProfiler("path").

Example usage:

	profiling.NewProfiler("profile").Memory().Start()

Run the following command to view the memory profile once App has finished running:

	go tool pprof profile/mem.pprof
*/
func (p *Profiler) Memory() *Profiler {
	p.mem = true
	return p
}

/*
CPU enables CPU profiling and saves the cpu.pprof file to provided output path to NewProfiler("path").

Example usage:

	profiling.NewProfiler("profile").CPU().Start()

Run the following command to view the memory profile once App has finished running:

	go tool pprof profile/cpu.pprof
*/
func (p *Profiler) CPU() *Profiler {
	p.cpu = true
	return p
}

/*
Trace enables tracing and saves the trace.out file to provided output path to NewProfiler("path").

-----------------

Example Usage in Application:

	profiling.NewProfiler("profile").Trace().Start()

-----------------

Run the following command to view the memory profile once App has finished running:

	go tool trace profile/trace.out

-----------------

Usage Note:

To add Tracing for Latency insights and additional information - use Tasks and Regions.

Example Usage in Application:

Adding Tasks to capture Latency, GC, and Syscall information for Functions or Goroutines

	// Define a Task for SomeFunc()
	ctx, task1 := trace.NewTask(context.Background(), "somepkg.SomeFunc")
	somepkg.SomeFunc()
	task1.End()

	// Define a Task for OtherFunc()
	_, task2 := trace.NewTask(context.Background(), "somepkg.OtherFunc")
	somepkg.OtherFunc()
	task2.End()

	// Then once app ends - run the following command to view the trace profile
	go tool trace profile/trace.out

	// Go to the Web UI and select the Tasks tab to view the Latency, GC, and Syscall information for SomeFunc() and OtherFunc()

-----------------

Refer to official docs for package runtime/trace for more information.

https://pkg.go.dev/runtime/trace
*/
func (p *Profiler) Tracing() *Profiler {
	p.trace = true
	return p
}

/*
NoLinger disables the default behavior of lingering for 5s after the program ends to collect memory profile data.
It is recommended to NOT use this flag so all memory profile data can be collected.

Example usage:

	profiling.NewProfiler("profile").Memory().NoLinger().Start()
	defer p.Stop()

	Use this flag if you do not want the application to wait 5s after execution ends.
*/
func (p *Profiler) NoLinger() *Profiler {
	p.linger = false
	return p
}

/*
Help prints a help message to the console with instructions on how to view the profiling data.
Recommended to use this flag to help users view the profiling data and access profiling data.

Example usage:

	profiling.NewProfiler("profile").Memory().CPU().Tracing().Help().Start()
	defer p.Stop()

	Use this flag to print a help message to the console with instructions on how to view the profiling data.
*/
func (p *Profiler) Help() *Profiler {
	p.helpFlag = true
	return p
}

/*
Start starts the profiler and returns the Profiler instance.

Example usage:

	profiling.NewProfiler("outputPath").Memory().CPU().Tracing().Help().Start()
	defer p.Stop()

	Note : It is recommended to use this method with defer to ensure the profiler is gracefully stopped after the program ends.
*/
func (p *Profiler) Start() (*Profiler, error) {
	if !p.mem && !p.cpu && !p.trace {
		fmt.Printf("\n<------\nProfiler has not been enabled for any metrics.\nEnable by using Memory(), CPU(), and Tracing().\n\n> p := profiling.NewProfiler(\"outputFolder\").Memory().CPU().Tracing().Help().NoLinger().Start()\n> defer p.Stop()\n------>\n\n")
		return p, errors.New("Profiler has not been enabled for any metrics. Enable by using Memory(), CPU(), and Tracing().\n> p := profiling.NewProfiler(\"outputFolder\").Memory().CPU().Tracing().Help().NoLinger().Start()\n> defer p.Stop()")
	}

	if p.helpFlag {
		p.printHelpMessage()
	}

	if p.trace {
		traceOut, err := os.Create(p.traceFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := trace.Start(traceOut); err != nil {
			log.Fatal(err)
		}
		p.traceOut = traceOut
	}

	if p.cpu {
		cpuOut, err := os.Create(p.cpuFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(cpuOut); err != nil {
			log.Fatal(err)
		}
		p.cpuOut = cpuOut
	}

	if p.mem {
		memOut, err := os.Create(p.memFile)
		if err != nil {
			log.Fatal(err)
		}
		p.memOut = memOut
	}

	fmt.Printf("Successfully started profiler.\n-------->\n\n")

	return p, nil
}

/*
Stop stops the profiler and saves the profiling data to the provided output path to NewProfiler("path").

Example usage:

	profiling.NewProfiler("outputPath").Memory().CPU().Tracing().Help().Start()
	defer p.Stop()

	Note : It is recommended to use this method with defer to ensure the profiler is gracefully stopped after the program ends.
*/
func (p *Profiler) Stop() {
	if !p.mem && !p.cpu && !p.trace {
		fmt.Println("Warning : Using Profiler with no Profiling. Error handling is recommended for profiling.NewProfiler(\"path\").Start()")
		return
	}
	if p.trace {
		trace.Stop()
		p.traceOut.Close()
	}

	if p.cpu {
		pprof.StopCPUProfile()
		p.cpuOut.Close()
	}

	if p.mem {
		if p.linger {
			fmt.Printf("App has finished executing.\nProfiler will linger for 5s to collect memory profile data.\nUse NoLinger() to disable this behavior - check GoDoc Info for Profiler.\n")
			time.Sleep(time.Second * 5)
		}
		pprof.WriteHeapProfile(p.memOut)
		p.printEndMessage()

		p.memOut.Close()
	}
}

func (p *Profiler) printHelpMessage() {
	if !p.mem && !p.cpu && !p.trace {
		return
	}
	activeProfiles := []string{}
	if p.mem {
		activeProfiles = append(activeProfiles, "Memory")
	}

	if p.cpu {
		activeProfiles = append(activeProfiles, "CPU")
	}

	if p.trace {
		activeProfiles = append(activeProfiles, "Tracing")
	}

	activeProfileStr := "No profiling flags have been set."
	if len(activeProfiles) > 0 {
		activeProfileStr = fmt.Sprintf("Profiling for %s.", joinWithCommasAndAnd(activeProfiles))
	}

	fmt.Printf(`
<-----
Profiler is enabled!
%s

After the program ends - the debugger will linger for 5s to collect memory profile data.
* To disable this - use NoLinger

> p := profiling.NewProfiler("outputFolder").Memory().CPU().Tracing().Help().NoLinger().Start()
> defer p.Stop()
------
# Viewing Function Usage
go tool pprof %s
(pprof) list utils.BulkInsert
(pprof) top5
------
# Viewing Memory Usage
go tool pprof %s
(pprof) list utils.BulkInsert
(pprof) top10
------
# Viewing Trace
go tool trace %s
------
# Web View
go tool pprof -http=:8080 %s
go tool pprof -http=:8080 %s
-----
`, activeProfileStr, p.cpuFile, p.memFile, p.traceFile, p.cpuFile, p.memFile)
}

func (p *Profiler) printEndMessage() {
	fmt.Println("\n<-----")
	fmt.Println("Profiler has finished collecting data!")

	if p.cpu {
		fmt.Printf("\n# Viewing CPU Profile\n")
		fmt.Printf("go tool pprof %s\n", p.cpuFile)
		fmt.Printf("(pprof) list <FunctionName>\n")
		fmt.Printf("(pprof) top5\n")
		fmt.Printf("For Web View: go tool pprof -http=:8080 %s\n", p.cpuFile)
	}

	if p.mem {
		fmt.Printf("\n# Viewing Memory Profile\n")
		fmt.Printf("go tool pprof %s\n", p.memFile)
		fmt.Printf("(pprof) list <FunctionName>\n")
		fmt.Printf("(pprof) top10\n")
		fmt.Printf("For Web View: go tool pprof -http=:8080 %s\n", p.memFile)
	}

	if p.trace {
		fmt.Printf("\n# Viewing Trace\n")
		fmt.Printf("go tool trace %s\n", p.traceFile)
	}

	fmt.Println("------>\n")
}

func joinWithCommasAndAnd(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return fmt.Sprintf("%s and %s", items[0], items[1])
	}
	return fmt.Sprintf("%s, and %s", strings.Join(items[:len(items)-1], ", "), items[len(items)-1])
}

/*

Usage

	profiler := NewProfiler().
		Memory("profile/mem.pprof").
		CPU("profile/cpu.pprof").
		Tracing("trace.out").
		Start()

	defer profiler.Stop()


*/

```
 ██████╗  ██████╗ ██╗     ██╗██████╗ ███████╗
██╔════╝ ██╔═══██╗██║     ██║██╔══██╗██╔════╝
██║  ███╗██║   ██║██║     ██║██████╔╝███████╗
██║   ██║██║   ██║██║     ██║██╔══██╗╚════██║
╚██████╔╝╚██████╔╝███████╗██║██████╔╝███████║
 ╚═════╝  ╚═════╝ ╚══════╝╚═╝╚═════╝ ╚══════╝
```

# utility libraries for go

- `github.com/Chinmay337/golibs/profiling`

  - Profiler to instrument and gather perf metrics from applications
  - Interface that uses `runtime/pprof` and `gotrace` to gather metrics
  - Will set up `PGO` - Profile Guided Optimization for Go Applications on subsequent builds.

```go
import "github.com/Chinmay337/golibs/profiling"

func main() {

	p := profiling.NewProfiler("outputFolder").
			 Tracing().Memory().CPU().Optimize().
			 Help().Start()

	defer p.Stop()

	... // rest of the app
}
```

- `github.com/Chinmay337/golibs/websockets`

  - Opinionated Websockets server implementation using `gorilla/websockets`
  - Provides default routes to echo and broadcast messages , keep track of connections , and ability to easily add more functionality/routes.

```go
import "github.com/Chinmay337/golibs/websockets"

func main() {

		wsServer := server.NewWsServer("8080").EnableAll().Start()

}
```

- `github.com/Chinmay337/golibs/utils`
  - Common utilities such as for copying files

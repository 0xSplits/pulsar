# pulsar

This repository contains the daemon responsible for our transaction indexing.
The relevant components here are described by the `server` and `worker`
packages. Pulsar's sole responsibility is to keep transaction information up to
date, according to all accounts managed in [Teams].

![Indexing Topology](.github/assets/Indexing-Topology.svg)

### Server

The server is a simple HTTP backend exposing the Pulsar's `/metrics` endopoint.

```
curl -s http://127.0.0.1:7777/metrics
```

```
# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
```

### Worker

The worker is a [custom task engine] executing asynchronous worker handlers
iteratively. New worker handlers can be added easily by implementing the handler
interface and registering the handler in the worker engine.

```golang
type Interface interface {
  // Cooler is manadatory to be implemented for worker handlers executed by the
  // *parallel.Worker engine, because those worker handlers do all run inside
  // their own isolated failure domains, which require individual cooler
  // durations to be provided. Cooler is irrelevant for the worker handlers
  // executed by the *sequence.Worker engine, because those handlers run inside
  // a single pipeline with a cooler duration congiured on the engine level.
  //
	// Cooler is the amount of time that any given handler specifies to wait
	// before being executed again. This is not an interval on a strict schedule.
	// This is simply the time to sleep after execution, before another cycle
	// repeats.
	Cooler() time.Duration

	// Ensure is the minimal worker handler interface that all users have to
	// implement for their own business logic, regardless of the underlying worker
	// engine.
	//
	// Ensure executes the handler specific business logic in order to complete
	// the given task, if possible. Any error returned will be emitted using the
	// underlying logger interface, unless the injected metrics registry is
	// configured to filter the received error.
	Ensure() error
}
```

### Usage

At its core, Pulsar is a simple [Cobra] command line tool, providing e.g. the
daemon command to start the long running `server` and `worker` processes.

```
pulsar -h
```

```
Golang based operator microservice.

Usage:
  pulsar [flags]
  pulsar [command]

Available Commands:
  daemon      Execute Pulsar's long running process for running the operator.
  deploy      Manually trigger a CloudFormation stack update.
  lint        Validate the release configuration under the given path.
  version     Print the version information for this command line tool.

Flags:
  -h, --help   help for pulsar

Use "pulsar [command] --help" for more information about a command.
```

### Development

As a convention, Pulsar's `.env` file should remain simple and generic. A
reasonable setting within that config file is e.g. `PULSAR_LOG_LEVEL`.

- `PULSAR_ENVIRONMENT`, the environment Pulsar is running in, one of `development` `testing` `staging` `production`.

```
pulsar daemon
```

```
{ "time":"2025-07-04 14:09:06", "level":"info", "message":"daemon is launching procs", "environment":"development", "caller":".../pkg/daemon/daemon.go:38" }
{ "time":"2025-07-04 14:09:06", "level":"info", "message":"server is accepting calls", "address":"127.0.0.1:7777",  "caller":".../pkg/server/server.go:95" }
{ "time":"2025-07-04 14:09:06", "level":"info", "message":"worker is executing tasks", "pipelines":"1",             "caller":".../pkg/worker/worker.go:110" }
```

### Releases

In order to update the Docker image, prepare all desired changes within the
`main` branch and create a Github release for the desired Pulsar version. The
release tag should be in [Semver Format]. Creating the Github release triggers
the responsible [Github Action] to build and push the Docker image to the
configured [Amazon ECR].

```
v0.1.11
```

The version command `pulsar version` and the version endpoint `/version` provide
build specific version information about the build and runtime environment. A
live demo can be seen at https://pulsar.testing.splits.org/version.

# Docker

Pulsar's build artifact is a statically compiled binary running in a
[distroless] image for maximum security and minimum size. If you do not have Go
installed and just want to run Pulsar locally in a Docker container, then use
the following commands.

```
docker build \
  --build-arg SHA="local-test-sha" \
  --build-arg TAG="local-test-tag" \
  -t pulsar:local .
```

```
docker run \
  -e PULSAR_ENVIRONMENT=development \
  -p 7777:7777 \
  pulsar:local \
  daemon
```

[Amazon ECR]: https://docs.aws.amazon.com/ecr
[Cobra]: https://github.com/spf13/cobra
[custom task engine]: https://github.com/0xSplits/workit
[distroless]: https://github.com/GoogleContainerTools/distroless
[Github Action]: .github/workflows/docker-release.yaml
[Semver Format]: https://semver.org
[Teams]: https://teams.splits.org

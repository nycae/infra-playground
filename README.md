# Infra Playground

Infra Playground (if you come up with a better name, please tell me), is an
application for testing tracing in distributed environments. This application
emulates a microservice based system that generates a static webpage with data
about random people.

## Examples

In order to run a complete simulation you can use the docker compose files at
build/compose/, the file `jaeger-kafka-es.yml` deploys the application, using
jaeger as the collector, a kafka intermediary and an ElasticSearch cluster as
the sink. Deploy the system, visit http://localhost:8080/v1/ (the first time
can take up to two seconds, but worry not, subsequent calls will be faster) and
http://localhost:16686 in order to see the jaeger's ui and mess with it as much
as you need to.

The example `otel-osearch-dprep.yml` uses otel-collector to forward traces to
a data-prepper pipeline that transforms and sinks the trace data to an
OpenSearch ecosystem. Start by visiting localhost:8080/v1/ to generate data and
log in to http://localhost:5601 to start browsing OpenSearch. If you're not
that much familiar with it, browse the burger menu > under the section named
*OpenSearch Plugins* > Click Observability > Left menu, click Trace analytics >
*Traces* is used to see the raw data of the traces and *Services* for the
dependency graph of those traces having into account the service they came from.

One last example can be found at the file `jaeger-osearch-kafka.yml`, but it
just combines the two previous examples in one. The open telemetry collector
will forward the traces to both, a jaeger collector and a data prepper, making
both previous scenarios possible for the same source of traces.


## Application insights

If you want to have a deeper understanding on how this app works the following
paragraphs are for you. Also, if you want to help extending the system and
fixing bugs, you're welcome!

First of all, the application communicates with 2 different public API's, so
it's necessary to run the containers on an environment in which public internet
is reachable. The used third party API's are:
- Random User Generator, [Randomuser.me](https://randomuser.me/)
- Random Age Gesser, [Agify.io](https://agify.io/)

You can see the dependency graph on any of the observability tools, but, for
the most part, the system is composed of 6 services:
- Frontend service (A Go http server that calls the remaining services in order
to build a static frontpage)
- Namer server (A gRPC server that calls the random user generator and returns
the provided data)
- Ager server (a gRPC server that, given a name returns a random age by calling
the random age API)
- City server (a gRPC server that, given a name returns a random city within a
hardcoded list of 2 cities).
- Heighter server (a gRPC server that, given a Name generates a random height
for that name, needs to call the limiter server in order to get the minimum and
maxmimum height).
- Limiter server (a gRPC server that returns the hardcoded values for minimum
and maximum height).

All the Go code is stored in the directories cmd/ and pkg/, some will be
generated at api/, but you can consider it a compilation artifact, and should
not be edited. This code can be regenerated running `make protos` and deleted
running `make clean`. Imporant! don't push it to the repository.

There's a directory in cmd/ for each one of the services, and that's all
there's to it. On the pkg/ directory you can find the implementations of the
gRPC interfaces for all the services as well as the http handler implementation
in the directories with the names of the services. There's also a tracing dir
with the different implementations of the otel library for all the use cases.
The dir pkg/utils/ is something that should be never do with Go, but I don't
know were else to put that code.

Within the build/docker/ directory you can find the dockerfiles for all the
services and under build/compose/ all the example arquitectures. The remaining
folders are:
- data-prepper/: stores the configurations of the data-prepper pipelines.
- otel-collector/: stores all the configurations of the opentelemetry-collector.
- static/: some static files for the frontend.

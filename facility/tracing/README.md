# Reference
  Specification

    http://opentracing.io/documentation/
    https://github.com/opentracing/specification/blob/master/specification.md
    https://github.com/opentracing/specification/blob/master/semantic_conventions.md

  Go implementation framework

    https://github.com/bg451/opentracing-example
    https://github.com/opentracing/opentracing-go
    https://github.com/grpc-ecosystem/grpc-opentracing/tree/master/go/otgrpc

  Tracing system

    http://zipkin.io/
    https://github.com/openzipkin/zipkin/
    https://github.com/openzipkin/zipkin-go-opentracing
    http://jaeger.readthedocs.io/en/latest/architecture/
    https://github.com/sourcegraph/appdash

# Environment
  localhost: OS X 10.10.3

# Installation
##  [opentracing-example]
an example for tracing

    $ docker run --rm -ti -p 8080:8080 -p 8700 bg451/opentracing-example

##  [opentracing-go]
package is a Go platform API for OpenTracing

    $ go get github.com/opentracing/opentracing-go

##  [otgrpc]

package enables distributed tracing in gRPC clients and servers via The OpenTracing Project

    $ go get github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc

# Glossary
##  Trace
*   [from jaeger & Opentracing spec.](http://jaeger.readthedocs.io/en/latest/architecture/)

    > A Trace is a data/execution path through the system, and can be thought of as a directed acyclic graph of spans.

*   [from zipkin](http://zipkin.io/pages/architecture.html)

    > Tracers live in your applications and record timing and metadata about operations that took place.

##  Span
*   [from jaeger](http://jaeger.readthedocs.io/en/latest/architecture/)

    > A Span represents a logical unit of work in the system that has an operation name, the start time of the operation, and the duration. Spans may be nested and ordered to model causal relationships. An RPC call is an example of a span.

*   [from zipkin](http://zipkin.io/pages/architecture.html)

    > For example, an instrumented web server records when it received a request and when it sent a response. The trace data collected is called a Span.

*   [from Opentracing spec.](https://github.com/opentracing/specification/blob/master/specification.md)

      > Each Span encapsulates the following state:

      >* An operation name

      >* A start timestamp

      >* A finish timestamp

      >* A set of zero or more key:value Span Tags. The keys must be strings. The values may be strings, bools, or numeric types.
        >> *OpenTracing project documents certain "standard tags" that have prescribed semantic meanings.*

      >* A set of zero or more Span Logs, each of which is itself a key:value map paired with a timestamp. The keys must be strings, though the values may be of any type. Not all OpenTracing implementations must support every value type.
      >* A SpanContext (see below)
      >* References to zero or more causally-related Spans (via the SpanContext of those related Spans)*

##  SpanContext
*   [from Opentracing spec.](https://github.com/opentracing/specification/blob/master/specification.md)

      > Each SpanContext encapsulates the following state:

      >* Any OpenTracing-implementation-dependent state (for example, trace and span ids) needed to refer to a distinct Span across a process boundary

      >* Baggage Items, which are just key:value pairs that cross process boundaries

##  References
*   [from Opentracing spec.](https://github.com/opentracing/specification/blob/master/specification.md)

      > The edges between Spans are called References. A Span may reference zero or more other SpanContexts that are causally related.

      *ChildOf references*

      > A Span may be the ChildOf a parent Span. In a ChildOf reference, the parent Span depends on the child Span in some capacity. All of the following would constitute ChildOf relationships:
        >* A Span representing the server side of an RPC may be the ChildOf a Span representing the client side of that RPC
        >* A Span representing a SQL insert may be the ChildOf a Span representing an ORM save method
        >* Many Spans doing concurrent (perhaps distributed) work may all individually be the ChildOf a single parent Span that merges the results for all children that return within a deadline

      *FollowsFrom references*

      > Some parent Spans do not depend in any way on the result of their child Spans. In these cases, we say merely that the child Span FollowsFrom the parent Span in a causal sense. There are many distinct FollowsFromreference sub-categories, and in future versions of OpenTracing they may be distinguished more formally.

##  Baggage Items
*   [from Opentracing spec.](https://github.com/opentracing/specification/blob/master/specification.md)

      > Baggage items are key:value string pairs that apply to the given Span, its SpanContext, and all Spans which directly or transitively reference the local Span. That is, baggage items propagate in-band along with the trace itself.

      > Baggage items enable powerful functionality given a full-stack OpenTracing integration (for example, arbitrary application data from a mobile app can make it, transparently, all the way into the depths of a storage system), and with it some powerful costs: use this feature with care.

      > Use this feature thoughtfully and with care. Every key and value is copied into every local and remote child of the associated Span, and that can add up to a lot of network and cpu overhead.

##  Inject(serialize) / extract(deserialize) a SpanContext into / from a carrier
*   [from Opentracing spec.](https://github.com/opentracing/specification/blob/master/specification.md)

      > Required parameters / return values

      >* A SpanContext instance
      >* A format descriptor (typically but not necessarily a string constant) which tells the Tracer implementation how to encode the SpanContext in the carrier parameter
      >* A carrier, whose type is dictated by the format. The Tracer implementation will encode the SpanContext in this carrier object according to the format.

##  formats for injection and extraction
*   [from Opentracing spec.](https://github.com/opentracing/specification/blob/master/specification.md)

      > All of the following formats must be supported by all Tracer implementations.

      >* Text Map: an arbitrary string-to-string map with an unrestricted character set for both keys and values
      >* HTTP Headers: a string-to-string map with keys and values that are suitable for use in HTTP headers (a la RFC 7230. In practice, since there is such "diversity" in the way that HTTP headers are treated in the wild, it is strongly recommended that Tracer implementations use a limited HTTP header key space and escape values conservatively.
      >* Binary: a (single) arbitrary binary blob representing a SpanContext

# pie [![GoDoc](https://godoc.org/github.com/natefinch/pie?status.svg)](https://godoc.org/github.com/natefinch/pie) [ ![Codeship Status for natefinch/pie](https://app.codeship.com/projects/ea82a1c0-4bae-0135-2de1-02fedcef81c5/status?branch=master)](https://app.codeship.com/projects/232834)

    import "github.com/natefinch/pie"

package pie provides a toolkit for creating plugins for Go applications.

![pie](https://cloud.githubusercontent.com/assets/3185864/7804562/bc35d256-0332-11e5-8562-fe00ec4d10b2.png)

**Why is it called pie?**

Because if you pronounce API like "a pie", then all this consuming and serving
of APIs becomes a lot more palatable.  Also, pies are the ultimate pluggable
interface - depending on what's inside, you can get dinner, dessert, a snack, or
even breakfast.  Plus, then I get to say that plugins in Go are as easy as...
well, you know.

If you have to explain it to your boss, just say it's an acronym for Plug In
Executables. <sub>(but it's not, really)</sub>

## About Pie

Plugins using this toolkit and the applications managing those plugins
communicate via RPC over the plugin application's Stdin and Stdout.

Functions in this package with the prefix `New` are intended to be used by the
plugin to set up its end of the communication.  Functions in this package
with the prefix `Start` are intended to be used by the main application to set
up its end of the communication and start a plugin executable.

<img src="https://cloud.githubusercontent.com/assets/3185864/7915136/8487d69e-0849-11e5-9dfa-13fc868f258f.png" />

This package provides two conceptually different types of plugins, based on
which side of the communication is the server and which is the client.
Plugins which provide an API server for the main application to call are
called Providers.  Plugins which consume an API provided by the main
application are called Consumers.

The default codec for RPC for this package is Go's gob encoding, however you
may provide your own codec, such as JSON-RPC provided by net/rpc/jsonrpc.

There is no requirement that plugins for applications using this toolkit be
written in Go. As long as the plugin application can consume or provide an
RPC API of the correct codec, it can interoperate with main applications
using this process.  For example, if your main application uses JSON-RPC,
many languages are capable of producing an executable that can provide a
JSON-RPC API for your application to use.

Included in this repo are some simple examples of a master process and a
plugin process, to see how the library can be used.  An example of the
standard plugin that provides an API the master process consumes is in the
examples/provider directory.  master\_provider expects plugin\_provider to be
in the same directory or in your $PATH.  You can just go install both of
them, and it'll work correctly.

In addition to a regular plugin that provides an API, this package can be
used for plugins that consume an API provided by the main process.  To see an
example of this, look in the examples/consumer folder.


## func NewConsumer
``` go
func NewConsumer() *rpc.Client
```
NewConsumer returns an rpc.Client that will consume an API from the host
process over this application's Stdin and Stdout using gob encoding.


## func NewConsumerCodec
``` go
func NewConsumerCodec(f func(io.ReadWriteCloser) rpc.ClientCodec) *rpc.Client
```
NewConsumerCodec returns an rpc.Client that will consume an API from the host
process over this application's Stdin and Stdout using the ClientCodec
returned by f.


## func StartProvider
``` go
func StartProvider(output io.Writer, path string, args ...string) (*rpc.Client, error)
```
StartProvider start a provider-style plugin application at the given path and
args, and returns an RPC client that communicates with the plugin using gob
encoding over the plugin's Stdin and Stdout.  The writer passed to output
will receive output from the plugin's stderr.  Closing the RPC client
returned from this function will shut down the plugin application.


## func StartProviderCodec
``` go
func StartProviderCodec(
    f func(io.ReadWriteCloser) rpc.ClientCodec,
    output io.Writer,
    path string,
    args ...string,
) (*rpc.Client, error)
```
StartProviderCodec starts a provider-style plugin application at the given
path and args, and returns an RPC client that communicates with the plugin
using the ClientCodec returned by f over the plugin's Stdin and Stdout. The
writer passed to output will receive output from the plugin's stderr.
Closing the RPC client returned from this function will shut down the plugin
application.


## type Server
``` go
type Server struct {
    // contains filtered or unexported fields
}
```
Server is a type that represents an RPC server that serves an API over
stdin/stdout.


### func NewProvider
``` go
func NewProvider() Server
```
NewProvider returns a Server that will serve RPC over this
application's Stdin and Stdout.  This method is intended to be run by the
plugin application.


### func StartConsumer
``` go
func StartConsumer(output io.Writer, path string, args ...string) (Server, error)
```
StartConsumer starts a consumer-style plugin application with the given path
and args, writing its stderr to output.  The plugin consumes an API this
application provides.  The function returns the Server for this host
application, which should be used to register APIs for the plugin to consume.


### func (Server) Close
``` go
func (s Server) Close() error
```
Close closes the connection with the client.  If the client is a plugin
process, the process will be stopped.  Further communication using this
Server will fail.


### func (Server) Register
``` go
func (s Server) Register(rcvr interface{}) error
```
Register publishes in the provider the set of methods of the receiver value
that satisfy the following conditions:


	- exported method
	- two arguments, both of exported type
	- the second argument is a pointer
	- one return value, of type error

It returns an error if the receiver is not an exported type or has no
suitable methods. It also logs the error using package log. The client
accesses each method using a string of the form "Type.Method", where Type is
the receiver's concrete type.


### func (Server) RegisterName
``` go
func (s Server) RegisterName(name string, rcvr interface{}) error
```
RegisterName is like Register but uses the provided name for the type
instead of the receiver's concrete type.


### func (Server) Serve
``` go
func (s Server) Serve()
```
Serve starts the Server's RPC server, serving via gob encoding.  This call
will block until the client hangs up.


### func (Server) ServeCodec
``` go
func (s Server) ServeCodec(f func(io.ReadWriteCloser) rpc.ServerCodec)
```
ServeCodec starts the Server's RPC server, serving via the encoding returned
by f. This call will block until the client hangs up.

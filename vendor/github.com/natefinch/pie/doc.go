// Package pie provides a toolkit for creating plugins for Go applications.
//
// Plugins using this toolkit and the applications managing those plugins
// communicate via RPC over the plugin application's Stdin and Stdout.
//
// Functions in this package with the prefix New are intended to be used by the
// plugin to set up its end of the communication.  Functions in this package
// with the prefix Start are intended to be used by the main application to set
// up its end of the communication and run a plugin executable.
//
// This package provides two conceptually different types of plugins, based on
// which side of the communication is the server and which is the client.
// Plugins which provide an API server for the main application to call are
// called Providers.  Plugins which consume an API provided by the main
// application are called Consumers.
//
// The default codec for RPC for this package is Go's gob encoding, however you
// may provide your own codec, such as JSON-RPC provided by net/rpc/jsonrpc.
//
// There is no requirement that plugins for applications using this toolkit be
// written in Go. As long as the plugin application can consume or provide an
// RPC API of the correct codec, it can interoperate with main applications
// using this process.  For example, if your main application uses JSON-RPC,
// many languages are capable of producing an executable that can provide a
// JSON-RPC API for your application to use.
//
// Included in this repo are some simple examples of a master process and a
// plugin process, to see how the library can be used.  An example of the
// standard plugin that provides an API the master process consumes is in the
// exmaples/provider directory.  master_provider expects plugin_provider to be
// in the same directory or in your $PATH.  You can just go install both of
// them, and it'll work correctly.

// In addition to a regular plugin that provides an API, this package can be
// used for plugins that consume an API provided by the main process.  To see an
// example of this, look in the examples/consumer folder.
package pie

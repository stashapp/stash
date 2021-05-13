package pie

import (
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

var errProcStopTimeout = errors.New("process killed after timeout waiting for process to stop")

// NewProvider returns a Server that will serve RPC over this
// application's Stdin and Stdout.  This method is intended to be run by the
// plugin application.
func NewProvider() Server {
	return Server{
		server: rpc.NewServer(),
		rwc:    rwCloser{os.Stdin, os.Stdout},
	}
}

// Server is a type that represents an RPC server that serves an API over
// stdin/stdout.
type Server struct {
	server *rpc.Server
	rwc    io.ReadWriteCloser
	codec  rpc.ServerCodec
}

// Close closes the connection with the client.  If the client is a plugin
// process, the process will be stopped.  Further communication using this
// Server will fail.
func (s Server) Close() error {
	if s.codec != nil {
		return s.codec.Close()
	}
	return s.rwc.Close()
}

// Serve starts the Server's RPC server, serving via gob encoding.  This call
// will block until the client hangs up.
func (s Server) Serve() {
	s.server.ServeConn(s.rwc)
}

// ServeCodec starts the Server's RPC server, serving via the encoding returned
// by f. This call will block until the client hangs up.
func (s Server) ServeCodec(f func(io.ReadWriteCloser) rpc.ServerCodec) {
	s.server.ServeCodec(f(s.rwc))
}

// Register publishes in the provider the set of methods of the receiver value
// that satisfy the following conditions:
//
//	- exported method
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
//
// It returns an error if the receiver is not an exported type or has no
// suitable methods. It also logs the error using package log. The client
// accesses each method using a string of the form "Type.Method", where Type is
// the receiver's concrete type.
func (s Server) Register(rcvr interface{}) error {
	return s.server.Register(rcvr)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (s Server) RegisterName(name string, rcvr interface{}) error {
	return s.server.RegisterName(name, rcvr)
}

// StartProvider start a provider-style plugin application at the given path and
// args, and returns an RPC client that communicates with the plugin using gob
// encoding over the plugin's Stdin and Stdout.  The writer passed to output
// will receive output from the plugin's stderr.  Closing the RPC client
// returned from this function will shut down the plugin application.
func StartProvider(output io.Writer, path string, args ...string) (*rpc.Client, error) {
	pipe, err := start(makeCommand(output, path, args))
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(pipe), nil
}

// StartProviderCodec starts a provider-style plugin application at the given
// path and args, and returns an RPC client that communicates with the plugin
// using the ClientCodec returned by f over the plugin's Stdin and Stdout. The
// writer passed to output will receive output from the plugin's stderr.
// Closing the RPC client returned from this function will shut down the plugin
// application.
func StartProviderCodec(
	f func(io.ReadWriteCloser) rpc.ClientCodec,
	output io.Writer,
	path string,
	args ...string,
) (*rpc.Client, error) {
	pipe, err := start(makeCommand(output, path, args))
	if err != nil {
		return nil, err
	}
	return rpc.NewClientWithCodec(f(pipe)), nil
}

// StartConsumer starts a consumer-style plugin application with the given path
// and args, writing its stderr to output.  The plugin consumes an API this
// application provides.  The function returns the Server for this host
// application, which should be used to register APIs for the plugin to consume.
func StartConsumer(output io.Writer, path string, args ...string) (Server, error) {
	pipe, err := start(makeCommand(output, path, args))
	if err != nil {
		return Server{}, err
	}
	return Server{
		server: rpc.NewServer(),
		rwc:    pipe,
	}, nil
}

// NewConsumer returns an rpc.Client that will consume an API from the host
// process over this application's Stdin and Stdout using gob encoding.
func NewConsumer() *rpc.Client {
	return rpc.NewClient(rwCloser{os.Stdin, os.Stdout})
}

// NewConsumerCodec returns an rpc.Client that will consume an API from the host
// process over this application's Stdin and Stdout using the ClientCodec
// returned by f.
func NewConsumerCodec(f func(io.ReadWriteCloser) rpc.ClientCodec) *rpc.Client {
	return rpc.NewClientWithCodec(f(rwCloser{os.Stdin, os.Stdout}))
}

// start runs the plugin and returns an ioPipe that can be used to control the
// plugin.
func start(cmd commander) (_ ioPipe, err error) {
	in, err := cmd.StdinPipe()
	if err != nil {
		return ioPipe{}, err
	}
	defer func() {
		if err != nil {
			in.Close()
		}
	}()
	out, err := cmd.StdoutPipe()
	if err != nil {
		return ioPipe{}, err
	}
	defer func() {
		if err != nil {
			out.Close()
		}
	}()

	proc, err := cmd.Start()
	if err != nil {
		return ioPipe{}, err
	}
	return ioPipe{out, in, proc}, nil
}

// makeCommand is a function that just creates an exec.Cmd and the process in
// it. It exists to facilitate testing.
var makeCommand = func(w io.Writer, path string, args []string) commander {
	cmd := exec.Command(path, args...)
	cmd.Stderr = w
	return execCmd{cmd}
}

type execCmd struct {
	*exec.Cmd
}

func (e execCmd) Start() (osProcess, error) {
	if err := e.Cmd.Start(); err != nil {
		return nil, err
	}
	return e.Cmd.Process, nil
}

// commander is an interface that is fulfilled by exec.Cmd and makes our testing
// a little easier.
type commander interface {
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	// Start is like exec.Cmd's start, except it also returns the os.Process if
	// start succeeds.
	Start() (osProcess, error)
}

// osProcess is an interface that is fullfilled by *os.Process and makes our
// testing a little easier.
type osProcess interface {
	Wait() (*os.ProcessState, error)
	Kill() error
	Signal(os.Signal) error
}

// ioPipe simply wraps a ReadCloser, WriteCloser, and a Process, and coordinates
// them so they all close together.
type ioPipe struct {
	io.ReadCloser
	io.WriteCloser
	proc osProcess
}

// Close closes the pipe's WriteCloser, ReadClosers, and process.
func (iop ioPipe) Close() error {
	err := iop.ReadCloser.Close()
	if writeErr := iop.WriteCloser.Close(); writeErr != nil {
		err = writeErr
	}
	if procErr := iop.closeProc(); procErr != nil {
		err = procErr
	}
	return err
}

// procTimeout is the timeout to wait for a process to stop after being
// signalled.  It is adjustable to keep tests fast.
var procTimeout = time.Second

// closeProc sends an interrupt signal to the pipe's process, and if it doesn't
// respond in one second, kills the process.
func (iop ioPipe) closeProc() error {
	result := make(chan error, 1)
	go func() { _, err := iop.proc.Wait(); result <- err }()
	if err := iop.proc.Signal(os.Interrupt); err != nil {
		return err
	}
	select {
	case err := <-result:
		return err
	case <-time.After(procTimeout):
		if err := iop.proc.Kill(); err != nil {
			return fmt.Errorf("error killing process after timeout: %s", err)
		}
		return errProcStopTimeout
	}
}

// rwCloser just merges a ReadCloser and a WriteCloser into a ReadWriteCloser.
type rwCloser struct {
	io.ReadCloser
	io.WriteCloser
}

// Close closes both the ReadCloser and the WriteCloser, returning the last
// error from either.
func (rw rwCloser) Close() error {
	err := rw.ReadCloser.Close()
	if err := rw.WriteCloser.Close(); err != nil {
		return err
	}
	return err
}

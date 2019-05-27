<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://godoc.org/github.com/gobuffalo/genny"><img src="https://godoc.org/github.com/gobuffalo/genny?status.svg" alt="GoDoc" /></a>
<a href="https://travis-ci.org/gobuffalo/genny"><img src="https://travis-ci.org/gobuffalo/genny.svg?branch=master" alt="Build Status" /></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/genny"><img src="https://goreportcard.com/badge/github.com/gobuffalo/genny" alt="Go Report Card" /></a>
</p>

# Genny

## What Is Genny?

Genny is a _framework_ for writing modular generators, it however, doesn't actually generate anything. It just makes it easier for you to. :)

## Core Concepts

### Generators

A [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) is used to build a blue print of what you want to generate.

A few of things that can be added to a `Generator` are:

* [`github.com/gobuffalo/genny#File`](https://godoc.org/github.com/gobuffalo/genny#File)
* [`os/exec#Cmd`](https://godoc.org/os/exec#Cmd)
* [`github.com/gobuffalo/packd#Box`](https://godoc.org/github.com/gobuffalo/packd#Box)
* [`net/http#Request`](https://godoc.org/net/http#Request)
* and more

A [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) does *not* actually generate anything; a [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) is needed to run the generator.

```go
g := genny.New()

// add a file
g.File(genny.NewFileS("index.html", "Hello\n"))

// execute a command
g.Command(exec.Command("go", "env"))

// run a function at run time
g.RunFn(func(r *genny.Runner) error {
  // look for the `genny` executable
  if _, err := r.LookPath("genny"); err != nil {
    // it wasn't found, so install it
    c := gogen.Get("github.com/gobuffalo/genny/genny")
    if err := r.Exec(c); err != nil {
      return err
    }
  }
  // call the `genny` executable with the `-h` flag.
  return r.Exec(exec.Command("genny", "-h"))
})
```

When a [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) is run each item that was added to it will be run in FIFO order. In the above example this means the following will happen:

1. Create a new file `r.Root/index.html`
1. Run the command `go env`
1. Run a function that installs `genny`

#### Runtime Checks

Genny has two different components; the "generator" (or blueprint) and the "runner" which executes the generator. Often it is necessary to only run certain code when the generator is run, not built. For example, checking the existing of an executable and installing it if missing.

In these situations you will want to use a [`github.com/gobuffalo/genny#RunFn`](https://godoc.org/github.com/gobuffalo/genny#RunFn) function.

In this example at runtime the `RunFn` will be called given the `*Runner` that is calling it. When called the function will ask the [`github.com/gobuffalo/genny#Runner.LookPath`](https://godoc.org/github.com/gobuffalo/genny#Runner.LookPath) function to ask the location of the `genny` executable.

In [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner) this will simply echo back the name of the executable that has been asked for, in this case `return "genny", nil`.

In [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner) this will call the [`os/exec#LookPath`](https://godoc.org/os/exec#LookPath) and return its results.

If the `genny` binary is not found, it will attempt to install it. Should that succeed the method returns the execution of a call to `genny -h`.

```go
g.RunFn(func(r *genny.Runner) error {
  // look for the `genny` executable
  if _, err := r.LookPath("genny"); err != nil {
    // it wasn't found, so install it
    c := gogen.Get("github.com/gobuffalo/genny/genny")
    if err := r.Exec(c); err != nil {
      return err
    }
  }
  // call the `genny` executable with the `-h` flag.
  return r.Exec(exec.Command("genny", "-h"))
})
```

The flexibility of the `*Fn` functions, combined with [`github.com/gobuffalo/genny#RunFn`](https://godoc.org/github.com/gobuffalo/genny#RunFn) make for a powerful testing combination.

### Runners

A [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) is used to run generators and control the environment in which those generators are run.

Genny ships with three implementations of `Runner` that cover _most_ situations. They can also provide good starting points for customized implementations of `Runner`.

* [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner)
* [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner)
* [`github.com/gobuffalo/genny/gentest#NewRunner`](https://godoc.org/github.com/gobuffalo/genny/gentest#NewRunner)

#### Adding Generators

To add a [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) to a [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) the [`github.com/gobuffalo/genny#Runner.With`](https://godoc.org/github.com/gobuffalo/genny#Runner.With) function can be used.

```go
run := genny.DryRunner(context.Background())

// add a generator from the `simple` package
g := simple.New()
run.With(g)

// add a generator from the `notsimple` package
g := notsimple.New()
run.With(g)
```

Each [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) is run in FIFO order in which it was added to the [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner).

It is common to have a function that builds a new [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) or returns an `error` if there was a problem.

```go
func New() (*genny.Generator, error) {
  g := simple.New()
  // do work which might error
  return g, nil
}
```

The [`github.com/gobuffalo/genny#Runner.WithNew`](https://godoc.org/github.com/gobuffalo/genny#Runner.WithNew) function was designed to make adding a [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator) with this return argument signature easier.

```go
if err := run.WithNew(New()); err != nil {
  log.Fatal(err)
}
```

#### Dry Running (**NON-DESTRUCTIVE**)

The idea of "dry" running means that no commands are executed, no files are written to disk, no HTTP requests are made, etc... Instead these steps are run "dry", which in the case of [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner) is the case.

```go
func main() {
  run := genny.DryRunner(context.Background())

  g := simple.New()
  run.With(g)

  if err := run.Run(); err != nil {
    log.Fatal(err)
  }
}
```

```plain
// output
DEBU[2018-12-06T15:13:47-05:00] Step: 4eac628c
DEBU[2018-12-06T15:13:47-05:00] Chdir: /go/src/github.com/gobuffalo/genny/internal/_examples/dry
DEBU[2018-12-06T15:13:47-05:00] File: /go/src/github.com/gobuffalo/genny/internal/_examples/dry/index.html
DEBU[2018-12-06T15:13:47-05:00] Exec: go env
DEBU[2018-12-06T15:13:47-05:00] LookPath: genny
DEBU[2018-12-06T15:13:47-05:00] Exec: genny -h
```

```bash
// file list
.
└── main.go

0 directories, 1 file
```

Using a "dry" runner can make testing easier when you don't have to worry about commands running, files being written, etc... It can also make it easy to provide a "dry-run" flag to your generators to let people see what will be generated when the generator is run for real.

#### Wet Running (**DESTRUCTIVE**)

While "dry" means to not execute commands or write files, "wet" running means the exact opposite; it will write files and execute commands.

Use the [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner) when "wet" running is the desired outcome.

```go
func main() {
  run := genny.WetRunner(context.Background())

  g := simple.New()
  run.With(g)

  if err := run.Run(); err != nil {
    log.Fatal(err)
  }
}
```

```plain
GOARCH="amd64"
GOBIN=""
// ...
A brief description of your application

Usage:
  genny [command]

Available Commands:
  help        Help about any command
  new         generates a new genny stub

Flags:
  -h, --help   help for genny

Use "genny [command] --help" for more information about a command.
```

```bash
// file list
.
├── index.html
└── main.go

0 directories, 2 files
```

```bash
$ cat index.html

Hello
```

#### Changing Runner Behavior

The change the way [`github.com/gobuffalo/genny#DryRunner`](https://godoc.org/github.com/gobuffalo/genny#DryRunner) or [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner) work, or to build your own [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) you need to implement the `*Fn` attributes on the [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner).

```go
type Runner struct {
  // ...
  ExecFn     func(*exec.Cmd) error                                     // function to use when executing files
  FileFn     func(File) (File, error)                                  // function to use when writing files
  ChdirFn    func(string, func() error) error                          // function to use when changing directories
  DeleteFn   func(string) error                                        // function used to delete files/folders
  RequestFn  func(*http.Request, *http.Client) (*http.Response, error) // function used to make http requests
  LookPathFn func(string) (string, error)                              // function used to make exec.LookPath lookups
  // ...
}
```

These `*Fn` functions represent the **FINAL** end-point for the that is trying to be run.

Here are two implementations of the [`github.com/gobuffalo/genny#Runner.FileFn`](https://godoc.org/github.com/gobuffalo/genny#Runner.FileFn) function.

The first will result in the file being printed to the screen. The second implementation writes the file to disk.

```go
run.FileFn = func(f packd.SimpleFile) (packd.SimpleFile, error) {
  io.Copy(os.Stdout, f)
  return f, nil
}

run.FileFn = func(f genny.File) (genny.File, error) {
  if d, ok := f.(genny.Dir); ok {
    if err := os.MkdirAll(d.Name(), d.Perm); err != nil {
      return f, err
    }
    return d, nil
  }

  name := f.Name()
  if !filepath.IsAbs(name) {
    name = filepath.Join(run.Root, name)
  }
  dir := filepath.Dir(name)
  if err := os.MkdirAll(dir, 0755); err != nil {
    return f, err
  }
  ff, err := os.Create(name)
  if err != nil {
    return f, err
  }
  defer ff.Close()
  if _, err := io.Copy(ff, f); err != nil {
    return f, err
  }
  return f, nil
}
```

### Files

Working with files, both creating new ones as well as, existing ones, is a core component of writing a generator. Genny understands this and offers several ways of working with files that is flexible and helps to make writing and testing your generators easier.

The [`github.com/gobuffalo/genny#File`](https://godoc.org/github.com/gobuffalo/genny#File) interface is the heart of working with files in Genny.

Genny ships with several convenience method for creating a [`github.com/gobuffalo/genny#File`](https://godoc.org/github.com/gobuffalo/genny#File).

* [`github.com/gobuffalo/genny#NewFile`](https://godoc.org/github.com/gobuffalo/genny#NewFile)
* [`github.com/gobuffalo/genny#NewFileS`](https://godoc.org/github.com/gobuffalo/genny#NewFileS)
* [`github.com/gobuffalo/genny#NewFileB`](https://godoc.org/github.com/gobuffalo/genny#NewFileB)
* [`github.com/gobuffalo/genny#NewDir`](https://godoc.org/github.com/gobuffalo/genny#NewDir)

#### Writing Files

To write a file you can add a [`github.com/gobuffalo/genny#File`](https://godoc.org/github.com/gobuffalo/genny#File) to your [`github.com/gobuffalo/genny#Generator.File`](https://godoc.org/github.com/gobuffalo/genny#Generator.File) and your file will then be handled by your `*Runner` when your generator is run.

```go
g.File(genny.NewFile("index.html", strings.NewReader("Hello\n")))
g.File(genny.NewFileS("strings/string.html", "Hello\n"))
g.File(genny.NewFileB("bytes/byte.html", []byte("Hello\n")))
```

In the case of [`github.com/gobuffalo/genny#WetRunner`](https://godoc.org/github.com/gobuffalo/genny#WetRunner) will attemp to create any directories your files require.

#### Reading Files

When writing generators you may need to read an existing file, perhaps to modify it, or perhaps read it's contents. This presents a problem in generators.

The first problem is that anytime we have to read files from disk, we make testing more difficult.

The bigger problems, however, present themselves more with "dry" runners (for example testing), than they do with "wet" runners.

If generator `A` creates a new file and generator `B` wants to modify that file in testing and "dry" runners this is a problem as the file may not present on disk for generator `B` to access.

To work around this issue Genny has the concept of a [`github.com/gobuffalo/genny#Disk`](https://godoc.org/github.com/gobuffalo/genny#Disk).

Now, instead of asking for the file directly from the file system, we can ask for it from the [`github.com/gobuffalo/genny#Runner.Disk`](https://godoc.org/github.com/gobuffalo/genny#Runner.Disk) instead.

```go
g.RunFn(func(r *genny.Runner) error {
  // try to find main.go either in the virtual "disk"
  // or the physical one
  f, err := r.Disk.Find("main.go")
  if err != nil {
    return err
  }
  // print the contents of the file
  fmt.Println(f.String())
  return nil
})
```

When asking for files from [`github.com/gobuffalo/genny#Runner.Disk`](https://godoc.org/github.com/gobuffalo/genny#Runner.Disk) it will first check its internal cache for the file, returning it if found. If the file is not in the cache, then it try to read it from disk at `filepath.Join(r.Root, name)`.

#### Transforming Files

There are times that you may need to transform either certain files, or all files. This could be as simple as replacing a variable in a template's name to match some user input, or something more complex, such as running any templates with a given extension through a certain template engine.

The [`github.com/gobuffalo/genny#Transformer`](https://godoc.org/github.com/gobuffalo/genny#Transformer) type can be used to implement these types of file transformations.

To create a new [`github.com/gobuffalo/genny#Transformer`](https://godoc.org/github.com/gobuffalo/genny#Transformer) you can use the [`github.com/gobuffalo/genny#NewTransformer`](https://godoc.org/github.com/gobuffalo/genny#NewTransformer) function.

The example below is taken from the [`github.com/gobuffalo/plushgen`](https://godoc.org/github.com/gobuffalo/plushgen) package.

```go
// Transformer will plushify any file that has a ".plush" extension
func Transformer(ctx *plush.Context) genny.Transformer {
  t := genny.NewTransformer(".plush", func(f genny.File) (genny.File, error) {
    s, err := plush.RenderR(f, ctx)
    if err != nil {
      return f, errors.Wrap(err, f.Name())
    }
    return genny.NewFileS(f.Name(), s), nil
  })
  t.StripExt = true
  return t
}
```

The [`github.com/gobuffalo/genny#Transformer`](https://godoc.org/github.com/gobuffalo/genny#Transformer) that is returned in the example will only be run on files that have a `.plush` extension in their name.

Should a file have a `.plush` extension, it will be sent to [`github.com/gobuffalo/plush`](https://godoc.org/github.com/gobuffalo/plush) to be rendered. The result of that rendering is returned as a new [`github.com/gobuffalo/genny#File`](https://godoc.org/github.com/gobuffalo/genny#File). Finally, the extension `.plush` will be stripped from the file name.

```go
g := genny.New()

// add a file
g.File(genny.NewFileS("index.html.plush", "Hello <%= name %>\n"))

// add the plush transformer
ctx := plush.NewContext()
ctx.Set("name", "World")
g.Transformer(plushgen.Transformer(ctx))
```

```plain
// output
DEBU[2018-12-07T10:35:56-05:00] Step: 09c9663e
DEBU[2018-12-07T10:35:56-05:00] Chdir: /go/src/github.com/gobuffalo/genny/internal/_examples/dry
DEBU[2018-12-07T10:35:56-05:00] File: /go/src/github.com/gobuffalo/genny/internal/_examples/dry/index.html
Hello World
```

### Testing

Testing a generator can be difficult because creating, deleting, and modifying files can be painful to handle during testing. The same can be said of running functions and HTTP requests.

The `*Fn` attributes on [`github.com/gobuffalo/genny#Runner`](https://godoc.org/github.com/gobuffalo/genny#Runner) make it simplier to mock out different test cases.

Most of the time the out of the box defaults are "good enough" for testing. The [`github.com/gobuffalo/genny/gentest`](https://godoc.org/github.com/gobuffalo/genny/gentest) package offers several helpers to simplify testing further.

In this example we test the "happy" path of a [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator).

```go
func Test_Happy(t *testing.T) {
  r := require.New(t)

  run := gentest.NewRunner()
  run.Disk.Add(genny.NewFileS("main.go", "my main.go file"))

  g := New()
  run.With(g)

  r.NoError(run.Run())
  res := run.Results()

  cmds := []string{"go env", "genny -h"}
  r.NoError(gentest.CompareCommands(cmds, res.Commands))

  files := []string{"index.html", "main.go"}
  r.NoError(gentest.CompareFiles(files, res.Files))
}
```

Notice how in the above example we had to add `main.go` to the [`github.com/gobuffalo/genny#Runner.Disk`](https://godoc.org/github.com/gobuffalo/genny#Runner.Disk). That is because the file doesn't exist in our testing directory.

In the following example we test what happens when the `genny` executable can not be found when running the [`github.com/gobuffalo/genny#Generator`](https://godoc.org/github.com/gobuffalo/genny#Generator).

We can simulate this experience by using the [`github.com/gobuffalo/genny#Runner.LookPathFn`](https://godoc.org/github.com/gobuffalo/genny#Runner.LookPathFn) to return an error if it is asked about that particular executable.

```go
func Test_Missing_Genny(t *testing.T) {
  r := require.New(t)

  run := gentest.NewRunner()
  run.Disk.Add(genny.NewFileS("main.go", "my main.go file"))

  g := New()
  run.With(g)

  // pretend we can't find genny
  run.LookPathFn = func(s string) (string, error) {
    if s == "genny" {
      return "", errors.New("can't find genny")
    }
    return s, nil
  }

  r.NoError(run.Run())
  res := run.Results()

  cmds := []string{"go env", "go get github.com/gobuffalo/genny/genny", "genny -h"}
  r.NoError(gentest.CompareCommands(cmds, res.Commands))

  files := []string{"index.html", "main.go"}
  r.NoError(gentest.CompareFiles(files, res.Files))
}
```

## The `genny` Executable

Genny ships with an executable that helps to generate new generators.

### Installation

```bash
$ go get -u github.com/gobuffalo/genny/genny
```

### Usage

```bash
$ genny -h

tools for working with genny

Usage:
  genny [command]

Available Commands:
  help        Help about any command
  new         generates a new genny stub

Flags:
  -h, --help   help for genny

Use "genny [command] --help" for more information about a command.
```

### Generating a New Generator

```bash
$ genny new coke -h

DEBU[2018-12-07T11:07:01-05:00] Step: a1d8eb2f
DEBU[2018-12-07T11:07:01-05:00] Chdir: /go/src/github.com/gobuffalo
DEBU[2018-12-07T11:07:01-05:00] File: /go/src/github.com/gobuffalo/coke/coke.go
DEBU[2018-12-07T11:07:01-05:00] File: /go/src/github.com/gobuffalo/coke/coke_test.go
DEBU[2018-12-07T11:07:01-05:00] File: /go/src/github.com/gobuffalo/coke/options.go
DEBU[2018-12-07T11:07:01-05:00] File: /go/src/github.com/gobuffalo/coke/options_test.go
DEBU[2018-12-07T11:07:01-05:00] File: /go/src/github.com/gobuffalo/coke/templates/example.txt
```


mockery
=======
[![Release](https://github.com/vektra/mockery/actions/workflows/release.yml/badge.svg)](https://github.com/vektra/mockery/actions/workflows/release.yml) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/vektra/mockery/v2?tab=overview) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vektra/mockery) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/vektra/mockery) [![Go Report Card](https://goreportcard.com/badge/github.com/vektra/mockery)](https://goreportcard.com/report/github.com/vektra/mockery) [![codecov](https://codecov.io/gh/vektra/mockery/branch/master/graph/badge.svg)](https://codecov.io/gh/vektra/mockery)




mockery provides the ability to easily generate mocks for golang interfaces using the [stretchr/testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock?tab=doc) package. It removes
the boilerplate coding required to use mocks.

Table of Contents
-----------------

- [Installation](#installation)
  * [Github Release](#github-release)
  * [Docker](#docker)
  * [Homebrew](#homebrew)
  * [go install](#go-install)
- [Examples](#examples)
    + [Simplest case](#simplest-case)
    + [Next level case](#next-level-case)
- [Return Value Provider Functions](#return-value-provider-functions)
    + [Requirements](#requirements)
    + [Notes](#notes)
- [Extended Flag Descriptions](#extended-flag-descriptions)
- [Mocking interfaces in `main`](#mocking-interfaces-in-main)
- [Configuration](#configuration)
  * [Example](#example)
- [Semantic Versioning](#semantic-versioning)
- [Stargazers](#stargazers)


Installation
------------

### Github Release

Visit the [releases page](https://github.com/vektra/mockery/releases) to download one of the pre-built binaries for your platform. 

### Docker

Use the [Docker image](https://hub.docker.com/r/vektra/mockery)

    docker pull vektra/mockery

### Homebrew

Install through [brew](https://brew.sh/)

    brew install mockery
    brew upgrade mockery

### go install

Alternatively, you can use the go install method:

    go install github.com/vektra/mockery/v2@latest

Examples
--------

![](https://raw.githubusercontent.com/vektra/mockery/master/docs/Peek%202020-06-28%2000-08.gif)

#### Simplest case

Given this is in `string.go`

```go
package test

type Stringer interface {
	String() string
}
```

Run: `mockery --name=Stringer` and the following will be output to `mocks/Stringer.go`:

```go
package mocks

import "github.com/stretchr/testify/mock"

type Stringer struct {
	mock.Mock
}

func (m *Stringer) String() string {
	ret := m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
```

#### Function type case

Given this is in `send.go`

```go
package test

type SendFunc func(data string) (int, error)
```

Run: `mockery --name=SendFunc` and the following will be output to `mocks/SendFunc.go`:

```go
package mocks

import "github.com/stretchr/testify/mock"

type SendFunc struct {
	mock.Mock
}

func (_m *SendFunc) Execute(data string) (int, error) {
	ret := _m.Called(data)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
```

#### Next level case

See [github.com/jaytaylor/mockery-example](https://github.com/jaytaylor/mockery-example)
for the fully runnable version of the outline below.

```go
package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jaytaylor/mockery-example/mocks"
	"github.com/stretchr/testify/mock"
)

func main() {
	mockS3 := &mocks.S3API{}

	mockResultFn := func(input *s3.ListObjectsInput) *s3.ListObjectsOutput {
		output := &s3.ListObjectsOutput{}
		output.SetCommonPrefixes([]*s3.CommonPrefix{
			&s3.CommonPrefix{
				Prefix: aws.String("2017-01-01"),
			},
		})
		return output
	}

	// NB: .Return(...) must return the same signature as the method being mocked.
	//     In this case it's (*s3.ListObjectsOutput, error).
	mockS3.On("ListObjects", mock.MatchedBy(func(input *s3.ListObjectsInput) bool {
		return input.Delimiter != nil && *input.Delimiter == "/" && input.Prefix == nil
	})).Return(mockResultFn, nil)

	listingInput := &s3.ListObjectsInput{
		Bucket:    aws.String("foo"),
		Delimiter: aws.String("/"),
	}
	listingOutput, err := mockS3.ListObjects(listingInput)
	if err != nil {
		panic(err)
	}

	for _, x := range listingOutput.CommonPrefixes {
		fmt.Printf("common prefix: %+v\n", *x)
	}
}
```


Return Value Provider Functions
--------------------------------

If your tests need access to the arguments to calculate the return values,
set the return value to a function that takes the method's arguments as its own
arguments and returns the return value. For example, given this interface:

```go
package test

type Proxy interface {
  passthrough(ctx context.Context, s string) string
}
```

The argument can be passed through as the return value:

```go
import . "github.com/stretchr/testify/mock"

Mock.On("passthrough", mock.AnythingOfType("context.Context"), mock.AnythingOfType("string")).Return(func(ctx context.Context, s string) string {
    return s
})
```

#### Requirements

`Return` must be passed the same argument count and types as expected by the interface. Then, for each of the return values of the mocked function, `Return` needs a function which takes the same arguments as the mocked function, and returns one of the return values. For example, if the return argument signature of `passthrough` in the above example was instead `(string, error)` in the interface, `Return` would also need a second function argument to define the error value:

```go
type Proxy interface {
  passthrough(ctx context.Context, s string) (string, error)
}
```

```go
Mock.On("passthrough", mock.AnythingOfType("context.Context"), mock.AnythingOfType("string")).Return(
	func(ctx context.Context, s string) string {
		return s
	},
	func(ctx context.Context, s string) error {
		return nil
	})
```

Note that the following is incorrect (you can't return all the return values with one function):
```go
Mock.On("passthrough", mock.AnythingOfType("context.Context"), mock.AnythingOfType("string")).Return(
	func(ctx context.Context, s string) (string, error) {
		return s, nil
	})
```

If any return argument is missing, `github.com/stretchr/testify/mock.Arguments.Get` will emit a panic.

For example, `panic: assert: arguments: Cannot call Get(0) because there are 0 argument(s). [recovered]` indicates that `Return` was not provided any arguments but (at least one) was expected based on the interface. `Get(1)` would indicate that the `Return` call is missing a second argument, and so on.

#### Notes

This approach should be used judiciously, as return values should generally
not depend on arguments in mocks; however, this approach can be helpful for
situations like passthroughs or other test-only calculations.


Extended Flag Descriptions
--------------------------

The following descriptions provide additional elaboration on a few common parameters.

| flag name  | description  |
|---|---|
| `--name`  | The `--name` option takes either the name or matching regular expression of interface to generate mock(s) for. |
| `--all`  |  It's common for a big package to have a lot of interfaces, so mockery provides `--all`. This option will tell mockery to scan all files under the directory named by `--dir` ("." by default) and generates mocks for any interfaces it finds. This option implies `--recursive=true`. |
| `--recursive`  |  Use the `--recursive` option to search subdirectories for the interface(s). This option is only compatible with `--name`. The `--all` option implies `--recursive=true`. |
| `--output` | mockery always generates files with the package `mocks` to keep things clean and simple. You can control which mocks directory is used by using `--output`, which defaults to `./mocks`. |
| `--inpackage` and `--keeptree` | For some complex repositories, there could be multiple interfaces with the same name but in different packages. In that case, `--inpackage` allows generating the mocked interfaces directly in the package that it mocks. In the case you don't want to generate the mocks into the package but want to keep a similar structure, use the option `--keeptree`. |
| `--filename` | Use the `--filename` and `--structname` to override the default generated file and struct name. These options are only compatible with non-regular expressions in `--name`, where only one mock is generated. |
| `--case` | mockery generates files using the casing of the original interface name.  This can be modified by specifying `--case underscore` to format the generated file name using underscore casing. |
| `--print` | Use `mockery --print` to have the resulting code printed out instead of written to disk. |
| `--exported` | Use `mockery --exported` to generate public mocks for private interfaces. |

Mocking interfaces in `main`
----------------------------

When your interfaces are in the main package you should supply the `--inpackage` flag.
This will generate mocks in the same package as the target code avoiding import issues.

Configuration
--------------

mockery uses [spf13/viper](https://github.com/spf13/viper) under the hood for its configuration parsing. It is bound to three different configuration sources, in order of decreasing precedence:

1. Command line
2. Environment variables
3. Configuration file

### Example

	$ export MOCKERY_STRUCTNAME=config_from_env
	$ echo $MOCKERY_STRUCTNAME
	config_from_env
	$ grep structname .mockery.yaml
	structname: config_from_file
	$ ./mockery showconfig --structname config_from_cli | grep structname
	Using config file: /home/ltclipp/git/vektra/mockery/.mockery.yaml
	structname: config_from_cli
	$ ./mockery showconfig  | grep structname
	Using config file: /home/ltclipp/git/vektra/mockery/.mockery.yaml
	structname: config_from_env
	$ unset MOCKERY_STRUCTNAME
	$ ./mockery showconfig  | grep structname
	Using config file: /home/ltclipp/git/vektra/mockery/.mockery.yaml
	structname: config_from_file

By default it searches the current working directory for a file named `.mockery.[extension]` where [extension] is any of the [recognized extensions](https://pkg.go.dev/github.com/spf13/viper@v1.7.0?tab=doc#pkg-variables).

Semantic Versioning
-------------------

The versioning in this project applies only to the behavior of the mockery binary itself. This project explicitly does not promise a stable internal API, but rather a stable executable. The versioning applies to the following:

1. CLI arguments.
2. Parsing of Golang code. New features in the Golang language will be supported in a backwards-compatible manner, except during major version bumps.
3. Behavior of mock objects. Mock objects can be considered to be part of the public API.
4. Behavior of mockery given a set of arguments.

What the version does _not_ track:
1. The interfaces, objects, methods etc. in the vektra/mockery package.
2. Compatibility of `go get`-ing mockery with new or old versions of Golang.

Development Efforts
-------------------

> v2 is in a soft change freeze due to the complexity of the software and the fact that functionality addition generally requires messing with logic that has been thoroughly tested, but is sensitive to change.

### v1

v1 is the original version of the software, and is no longer supported.

### v2

`mockery` is currently in v2, which iterates on v1 and includes mostly cosmetic and configuration improvements. 

### v3

[v3](https://github.com/vektra/mockery/projects/3) will include a ground-up overhaul of the entire codebase and will completely change how mockery works internally and externally. The highlights of the project are:
- Moving towards a package-based model instead of a file-based model. `mockery` currently iterates over every file in a project and calls `package.Load` on each one, which is time consuming. Moving towards a model where the entire package is loaded at once will dramtically reduce runtime, and will simplify logic. Additionally, supporting only a single mode of operation (package mode) will greatly increase the intuitiveness of the software.
- Configuration-driven generation. `v3` will be entirely driven by configuration, meaning:
  * You specify the packages you want mocked, instead of relying on it auto-discovering your package. Auto-discovery in theory sounds great, but in practice it leads to a great amount of complexity for very little benefit.
  * Package- or interface-specific overrides can be given that change mock generation settings on a granular level. This will allow your mocks to be generated in a heterogenous manner, and will be made explicit by yaml configuration.
 - Proper error reporting. Errors across the board will be done in accordance with modern Golang practices
 - Variables in generated mocks will be given meaningful names. 
 


Stargazers
----------

[![Stargazers over time](https://starchart.cc/vektra/mockery.svg)](https://starchart.cc/vektra/mockery)

symwalk
=======

``` go
    import "github.com/facebookgo/symwalk"
```

Package symwalk provides an implementation of symbolic link aware filepath walk.

Walk calls [filepath.Walk](http://golang.org/pkg/path/filepath/#Walk) by providing it with a special WalkFn that wraps the given WalkFn.
This function calls the given WalkFn for regular files.
However, when it encounters a symbolic link, it resolves the link fully using
[filepath.EvalSymlinks](http://golang.org/pkg/path/filepath/#EvalSymlinks) and recursively calls symwalk.Walk on the resolved path.
This ensures that unlike filepath.Walk, traversal does not stop at symbolic links.

Using it can be as simple as:

``` go
    Walk(
      "/home/me/src",
      func(path string, info os.FileInfo, err error) error {
        fmt.Println(path)
        return nil
      },
    )
```

**CAVEAT**: Note that symwalk.Walk does not terminate if there are any non-terminating loops in the file structure.

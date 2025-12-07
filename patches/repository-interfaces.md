# Repository Interface Modifications

After merging from upstream, ensure these Faceter interfaces are embedded in the Reader interfaces.

## pkg/models/repository_scene.go

Add `SceneFaceter` to `SceneReader`:

```go
// SceneReader provides all methods to read scenes.
type SceneReader interface {
	SceneFinder
	SceneQueryer
	SceneCounter
	SceneFaceter  // <-- ADD THIS LINE

	URLLoader
	// ... rest of interface
}
```

## pkg/models/repository_performer.go

Add `PerformerFaceter` to `PerformerReader`:

```go
// PerformerReader provides all methods to read performers.
type PerformerReader interface {
	PerformerFinder
	PerformerQueryer
	PerformerAutoTagQueryer
	PerformerCounter
	PerformerFaceter  // <-- ADD THIS LINE

	AliasLoader
	// ... rest of interface
}
```

## pkg/models/repository_gallery.go

Add `GalleryFaceter` to `GalleryReader`:

```go
// GalleryReader provides all methods to read galleries.
type GalleryReader interface {
	GalleryFinder
	GalleryQueryer
	GalleryCounter
	GalleryFaceter  // <-- ADD THIS LINE

	URLLoader
	// ... rest of interface
}
```

## pkg/models/repository_group.go

Add `GroupFaceter` to `GroupReader`:

```go
// GroupReader provides all methods to read groups.
type GroupReader interface {
	GroupFinder
	GroupQueryer
	GroupCounter
	GroupFaceter  // <-- ADD THIS LINE
	URLLoader
	// ... rest of interface
}
```

## pkg/models/repository_studio.go

Add `StudioFaceter` to `StudioReader`:

```go
// StudioReader provides all methods to read studios.
type StudioReader interface {
	StudioFinder
	StudioQueryer
	StudioAutoTagQueryer
	StudioCounter
	StudioFaceter  // <-- ADD THIS LINE

	AliasLoader
	// ... rest of interface
}
```

## pkg/models/repository_tag.go

Add `TagFaceter` to `TagReader`:

```go
// TagReader provides all methods to read tags.
type TagReader interface {
	TagFinder
	TagQueryer
	TagAutoTagQueryer
	TagCounter
	TagFaceter  // <-- ADD THIS LINE

	AliasLoader
	// ... rest of interface
}
```

## Verification

After adding these, run:

```bash
go build ./...
```

If the interfaces are correctly embedded, the build will succeed. If not, the compiler will report missing method implementations.


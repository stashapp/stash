// Package models - Fork Extension Interfaces Documentation
//
// This file documents the Faceter interfaces used by the fork's facets system.
// The actual interface definitions are in the repository_*.go files.
//
// IMPORTANT: After merging from upstream, ensure these interfaces are
// embedded in their respective Reader interfaces.
//
// See also:
//   - /patches/repository-interfaces.md for exact lines to add
//   - ui/v2.5/src/extensions/docs/BACKEND-API.md for full documentation
//
// =============================================================================
// FACET INTERFACE LOCATIONS
// =============================================================================
//
// repository_scene.go:
//   - SceneFaceter interface (line ~40)
//   - Embedded in SceneReader (line ~99)
//
// repository_performer.go:
//   - PerformerFaceter interface (line ~30)
//   - Embedded in PerformerReader (line ~82)
//
// repository_gallery.go:
//   - GalleryFaceter interface (line ~33)
//   - Embedded in GalleryReader (line ~70)
//
// repository_group.go:
//   - GroupFaceter interface (line ~28)
//   - Embedded in GroupReader (line ~72)
//
// repository_studio.go:
//   - StudioFaceter interface (line ~29)
//   - Embedded in StudioReader (line ~81)
//
// repository_tag.go:
//   - TagFaceter interface (line ~38)
//   - Embedded in TagReader (line ~94)
//
// =============================================================================
// MERGE INSTRUCTIONS
// =============================================================================
//
// After merging from upstream, if repository_*.go files are overwritten:
//
// 1. Add the Faceter interface definition (copy from this fork's version)
// 2. Add the Faceter interface to the Reader interface
//
// Example for SceneReader:
//
//   type SceneReader interface {
//       SceneFinder
//       SceneQueryer
//       SceneCounter
//       SceneFaceter    // <-- Add this line
//       // ... rest of interface
//   }
//
// See /patches/repository-interfaces.md for complete examples.

package models

// This file intentionally contains no code declarations.
// All Faceter interfaces are defined in their respective repository_*.go files.

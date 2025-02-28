// Package autotag provides the autotagging functionality for the application.
//
// The autotag functionality sets media metadata based on the media's path.
// The functions in this package are in the form of {ObjectType}{TagTypes},
// where the ObjectType is the single object instance to run on, and TagTypes
// are the related types.
// For example, PerformerScenes finds and tags scenes with a provided performer,
// whereas ScenePerformers tags a single scene with any Performers that match.
package autotag

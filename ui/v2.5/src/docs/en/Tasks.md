# Tasks

This page allows you to direct the stash server to perform a variety of tasks.

# Scanning

The scan function walks through the stash directories you have configured for new and moved files. 

Stash currently identifies files by performing a quick file hash. This means that if the file is renamed for moved elsewhere within your configured stash directories, then the scan will detect this and update its database accordingly.

Stash currently ignores duplicate files. If two files contain identical content, only the first one it comes across is used.

The scan task accepts the following options:

| Option | Description |
|--------|-------------|
| Generate previews | Generates video previews which play when hovering over a scene. |
| Generate animated image previews | Generates animated webp previews. Only required if the Preview Type is set to Animated Image. Requires Generate previews to be enabled. |
| Generate sprites | Generates sprites for the scene scrubber. |
| Generate perceptual hashes | Generates perceptual hashes for scene deduplication and identification. |
| Generate thumbnails for images | Generates thumbnails for image files. | 
| Don't include file extension in title | By default, scenes, images and galleries have their title created using the file basename. When the flag is enabled, the file extension is stripped when setting the title. |
| Set name, date, details from embedded file metadata. | Parse the video file metadata (where supported) and set the scene attributes accordingly. It has previously been noted that this information is frequently incorrect, so only use this option where you are certain that the metadata is correct in the files. |

# Auto Tagging
See the [Auto Tagging](/help/AutoTagging.md) page.

# Scene Filename Parser
See the [Scene Filename Parser](/help/SceneFilenameParser.md) page.

# Generated Content

The scanning function automatically generates a screenshot of each scene. The generated content provides the following:
* Video or image previews that are played when mousing over the scene card
* Perceptual hashes - helps match against StashDB, and feeds the duplicate finder
* Sprites (scene stills for parts of each scene) that are shown in the scene scrubber
* Marker video previews that are shown in the markers page
* Transcoded versions of scenes. See below
* Image thumbnails of galleries

The generate task accepts the following options:

| Option | Description |
|--------|-------------|
| Previews | Generates video previews which play when hovering over a scene. |
| Animated image previews | Generates animated webp previews. Only required if the Preview Type is set to Animated Image. Requires Generate previews to be enabled. |
| Scene Scrubber Sprites | Generates sprites for the scene scrubber. |
| Markers Previews | Generates 20 second videos which begin at the marker timecode. |
| Marker Animated Image Previews | Generates animated webp previews for markers. Only required if the Preview Type is set to Animated Image. Requires Markers to be enabled. |
| Marker Screenshots | Generates static JPG images for markers. Only required if Preview Type is set to Static Image. Requires Marker Previews to be enabled. | 
| Transcodes | MP4 conversions of unsupported video formats. Allows direct streaming instead of live transcoding. |
| Perceptual hashes | Generates perceptual hashes for scene deduplication and identification. |
| Overwrite existing generated files | By default, where a generated file exists, it is not regenerated. When this flag is enabled, then the generated files are regenerated. |

## Transcodes

Web browsers support a limited number of video and audio codecs and containers. Stash will directly stream video files where the browser supports the codecs and container. Originally, stash did not support viewing scene videos where the browser did not support the codecs/container, and generating transcodes was a way of viewing these files.

Stash has since implemented live transcoding, so transcodes are essentially unnecessary now. Further, transcodes use up a significant amount of disk space and are not guaranteed to be lossless.

## Image gallery thumbnails

These are generated when the gallery is first viewed, so generating them beforehand is not necessary.

# Cleaning

This task will walk through your configured media directories and remove any scene from the database that can no longer be found. It will also remove generated files for scenes that subsequently no longer exist.

Care should be taken with this task, especially where the configured media directories may be inaccessible due to network issues.

# Exporting and Importing

The import and export tasks read and write JSON files to the configured metadata directory. Import from file will merge your database with a file.

> **⚠️ Note:** The full import task wipes the current database completely before importing.

See the [JSON Specification](/help/JSONSpec.md) page for details on the exported JSON format.

---

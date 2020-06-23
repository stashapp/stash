# Tasks

This page allows you to direct the stash server to perform a variety of tasks.

> **⚠️ Note:** It is currently only possible to run one task at a time. No queuing is currently implemented.

# Scanning

The scan function walks through the stash directories you have configured for new and moved files. 

Stash currently identifies files by performing a full MD5 hash on them. This means that if the file is renamed for moved elsewhere within your configured stash directories, then the scan will detect this and update its database accordingly.

Stash currently ignores duplicate files. If a file is detected with the same hash as a file already in the database (and that file still exists on the filesystem), then the duplicate file is ignored.

The "Set name, data, details from metadata" option will parse the files metadata (where supported) and set the scene attributes accordingly. It has previously been noted that this information is frequently incorrect, so only use this option where you are certain that the metadata is correct in the files.

# Auto Tagging
See the [Auto Tagging](/help/AutoTagging.md) page.

# Scene Filename Parser
See the [Scene Filename Parser](/help/SceneFilenameParser.md) page.

# Generated Content

The scanning function automatically generates a screenshot of each scene. The generated content provides the following:
* video or image previews that are played when mousing over the scene card
* sprites (scene stills for parts of each scene) that are shown in the scene scrubber 
* marker video previews that are shown in the markers page
* transcoded versions of scenes. See below
* image thumbnails of galleries

## Transcodes

Web browsers support a limited number of video and audio codecs and containers. Stash will directly stream video files where the browser supports the codecs and container. Originally, stash did not support viewing scene videos where the browser did not support the codecs/container, and generating transcodes was a way of viewing these files.

Stash has since implemented live transcoding, so transcodes are essentially unnecessary now. Further, transcodes use up a significant amount of disk space and are not guaranteed to be lossless.

## Image gallery thumbnails

These are generated when the gallery is first viewed, so generating them beforehand is not necessary.

# Cleaning

This task will walk through your configured media directories and remove any scene from the database that can no longer be found. It will also remove generated files for scenes that subsequently no longer exist.

Care should be taken with this task, especially where the configured media directories may be inaccessible due to network issues.

# Exporting and Importing

The import and export tasks read and write JSON files to the configured metadata directory. 

> **⚠️ Note:** The import task wipes the current database completely before importing.

See the [JSON Specification](/help/JSONSpec.md) page for details on the exported JSON format.

---

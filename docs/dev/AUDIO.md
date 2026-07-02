# Audio Datatype

The `Audio` datatype is similar to `Scene` but stores audio-only media (i.e. Audiobooks, music, ASMR, etc).

## Scope

- This ticket adds backend support for Audio Only, future tickets can add the UI elements
    - Database design
    - Graphql Support
    - Scanner Support
        - No transcodes right now, but will keep the infrastructure to more easily support adding transcodes in the future

- Audio metadata:
    - Title
    - Date
    - Studio
    - Performers
    - Tags
    - Details
    - Urls
    - Rating
    - Organized
    - O History
    - Play History
    - Groups
    - ~~Captions~~ [REMOVED PER PR COMMENT - can be added back later]
- Audio File metadata:
    - duration
    - audio codec
    - bitrate
    - sample rate

### Open Questions

- Should Audio's have `cover` photo?
    - ANSWER: yes
- Should Legacy/Deprecate features be copied over?
    - ANSWER: do not copy over deprecated features
- What should be done for `sortByOCounter`/`sortByPlayCount`?
    - These assume SCENES
    - I see 3 options
        - ignore
        - add `audios` into the calculation
        - split into `sortBySceneOCounter` and `sortByAudioOCounter`

- Flaky Unit Tests
    - pkg/sqlite/studio_test.go --> `TestStudioQuerySortOCounter`
        - Assumes an id, but doesn't enforce this in the test setup
        - This test will likely keep breaking due to this
        - might be better to for a very large o_count on a specific studio to ensure this does not keep failing

## Future Tickets

- UI
    - Audio using `video.js` (ref: https://videojs.org/blog/video-js-4-9-now-can-join-the-party)
    - Audio Waveform (ref: https://github.com/collab-project/videojs-wavesurfer)
    - New AudioPlayer.tsx (copy `ui/v2.5/src/components/ScenePlayer/ScenePlayer.tsx`)

## General TODO

- [x] Setup Database
- [x] Scanner to scan Audio Files and create Audios
    - [x] FFProbe for Audio Files
- [x] Graphql to return Audios (queries)
- [x] Graphql to update Audios (mutations)
- [x] Update test files


## Notes

- Phashes cannot be used on audio files; A future ticket might introduce Chromaprint (AcoustID)
- Gallery could be added to Audio, but I am removing to reduce PR complexity
- StashIDs was removed, audio is unlikely to be added immediately to stashbox
- Audio's could have interactive components, but removed to reduce PR complexity

## Last Steps
- [ ] Delete this file upon completion of the feature


## Manual Tests

### Setup

1. Copy `.mp3` files into `.local-data`
2. `make server-clean`
3. `make server-start` OR run go debugger (VSCode F5)
4. Create new instance with library at `./.local-data/`
5. go to <http://127.0.0.1:9999/playground>
    - Perform manual tests here

### Check Query

This is a manual test with all fields. The test ensures that the Querying is setup correctly.

Later you can reuse this to ensure that mutations correctly updated the database.

```graphql
query {
  findAudios(filter:{sort:"title" direction:DESC}){
    count
    audios {
			id title code details urls date rating100 organized o_counter created_at updated_at last_played_at resume_time play_duration play_count play_history o_history custom_fields

      files{
          id path basename mod_time size format duration audio_codec sample_rate bit_rate created_at updated_at
          parent_folder{id}
          zip_file{id}
          fingerprints{type value}
      }
      captions{language_code caption_type}
      paths{caption stream}
      studio{id}
      groups{group{id} audio_index}
      tags{id}
      performers{id}
      audioStreams{url mime_type label}
    }
  }
  # findScenes(filter:{sort:"title" direction:DESC}){
  #   count
  #   scenes {
  #     id sceneStreams{url mime_type label}
  #     files{id path fingerprints{type value}}
  #   }
  # }
}
```

### Check Mutations

```graphql
mutation audio_mut {
  audioAddO(id:1){count history}
  audioUpdate(input:{id:1 title:"testing 1"}){id title o_history}
  audiosUpdate(input:[{id:1 details:"details 1"}]){id title details}
}
```

### Check Streams

Currently only direct streams are implemented. Use the following to get the Stream URL.

1. Execute this GraphQL
2. Paste the `Direct stream` url into the browser, ensure that the audio plays

```graphql
query {
  findAudios(filter:{sort:"title" direction:DESC}){
    count
    audios {id audioStreams{url mime_type label}
    }
  }
}
```

### HTML Confirmation

```html
<audio controls>
  <source src="http://127.0.0.1:9999/audio/1/stream" type="audio/mp3">
Your browser does not support the audio element.
</audio>
```

You can also listen to audio using VIDEO tag

```html
<video controls>
  <source src="http://127.0.0.1:9999/audio/1/stream" type="audio/mp3">
Your browser does not support the video element.
</video>
```
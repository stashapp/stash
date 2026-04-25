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
- Audio File metadata:
    - duration
    - audio codec
    - bitrate
    - sample rate

### Open Questions

- Should Audio's have `cover` photo?
- Should Legacy/Deprecate features be copied over?
    - Since Audio's is NEW, it doesn't have to support deprecated features/naming/etc
    - I suggest removing them if easy to do, and for the more complicated ones to defer to a separate ticket

## Future Tickets

- UI
    - Audio using `video.js` (ref: https://videojs.org/blog/video-js-4-9-now-can-join-the-party)
    - Audio Waveform (ref: https://github.com/collab-project/videojs-wavesurfer)
    - New AudioPlayer.tsx (copy `ui/v2.5/src/components/ScenePlayer/ScenePlayer.tsx`)

## General TODO

- [x] Setup Database
- [ ] Scanner to scan Audio Files and create Audios
    - [ ] FFProbe for Audio Files
- [ ] Graphql to return Audios (queries)
- [ ] Graphql to update Audios (mutations)


## Notes

- Phashes cannot be used on audio files; A future ticket might introduce Chromaprint (AcoustID)
- Gallery could be added to Audio, but I am removing to reduce PR complexity
- StashIDs was removed, audio is unlikely to be added immediately to stashbox
- Audio's could have interactive components, but removed to reduce PR complexity

## Last Steps
- [ ] Delete this file upon completion of the feature
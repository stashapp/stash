# Audio Datatype

The `Audio` datatype is similar to `Scene` but stores audio-only media (i.e. Audiobooks, music, ASMR, etc).

## Scope

- This ticket adds backend support for Audio Only, future tickets can add the UI elements
- Audio metadata:
    - Title
    - Artists (string? like director)
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
    - Studio Code
    - NICE TO HAVES
        - Groups
- Audio File metadata:
    - duration
    - audio codec
    - FUTURE (to be considered at a later date)
        - channels (mono, stereo, 5.1, 7.1)
        - bitrate
        - sample rate


## TODO List

- [ ] `pkg/sqlite/migrations/86_audio.up.sql`
    - Create a migration for the Audio type, very similar to Scene
- [ ] Duplicate much of `pkg/scene/*` into `pkg/audio/*`
    - Exclude: markers, screenshot, preview, transcode, sprite
- [ ] Graphql
    - [ ] Copy/modify `graphql/schema/types/scene.graphql` to `graphql/schema/types/audio.graphql`

### Last Steps
- [ ] Delete this file upon completion of the feature
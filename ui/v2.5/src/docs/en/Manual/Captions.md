# Captions

Stash supports captioning with SRT and VTT files.

These files need to be named as follows:

## Scene

- {scene_name}.{language_code}.ext
- {scene_name}.ext

Where `{language_code}` is defined by the [ISO-6399-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (2 letters) standard and `ext` is the file extension. Captions files without a language code will be labeled as Unknown in the video player but will work fine.

Scenes with captions can be filtered with the `captions` criterion.

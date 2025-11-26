# Captions

Stash supports captioning with SRT and VTT files.

Captions will only be detected if they are located in the same folder as the corresponding scene file.

Ensure the caption files follow these naming conventions:

## Scene

- {scene_file_name}.{language_code}.ext
- {scene_file_name}.ext

Where `{language_code}` is defined by the [ISO-6399-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (2 letters) standard and `ext` is the file extension. Captions files without a language code will be labeled as Unknown in the video player but will work fine.

Scenes with captions can be filtered with the `captions` criterion.

**Note:** If the caption file was added after the scene was initially added during scan, you will need to run a Selective Scan task for it to show up.

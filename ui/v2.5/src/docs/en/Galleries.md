# Galleries

Stash offers support for image galleries.
Here are some remarks on using them:

- **Galleries are zip-folders with images (e.g. jpeg or png) in them.**
- Stash searches for zip galleries in the same paths it searches for videos.
- In order for a gallery to be associated with a scene, the zip file and the video file must be in the same folder.
- For best results, images in zip file should be stored without compression (copy, store or no compression options depending on the software you use. Eg on linux: `zip -0 -r gallery.zip foldertozip/`). This impacts **heavily** on the zip read performance.
- Stash uses the golang native (pure go) image decoders (more suitable for cross compilation). With huge images, decoding and converting to thumbnails can be slow and  in some cases cause visual errors or delays when loading the gallery page.
- Stash adds a gallery to its related scene during the scanning process if they have matching names. For example, gallery `/my/stash/collection/media_filename.zip` will be auto assigned to `/my/stash/collection/media_filename.mp4` (where **mp4** can any supported video extension).
- If an filename of an image in the gallery zip file ends with `cover.jpg`, it will be treated like a cover and presented first in the gallery view page and as a gallery cover in the gallery list view. If more than one images match the name the first one found in natural sort order is selected.
- Gallery thumbnails are cached. The first time you go to a gallery page the thumbnails of the images are created and stored to the disk cache for later use. If you want to populate the cache beforehand, this can be done using the Generate task.
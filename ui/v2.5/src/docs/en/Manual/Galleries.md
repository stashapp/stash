# Galleries

**Note:** images are now included during the scan process and are loaded independently of galleries. It is _no longer necessary_ to have images in zip files to be scanned into your library.

Galleries are automatically created from zip files found during scanning that contain images. It is also possible to automatically create galleries from folders containing images, by selecting the "Create galleries from folders containing images" checkbox in the Configuration page. It is also possible to manually create galleries.

For best results, images in zip file should be stored without compression (copy, store or no compression options depending on the software you use. Eg on linux: `zip -0 -r gallery.zip foldertozip/`). This impacts **heavily** on the zip read performance.

If an filename of an image in the gallery zip file ends with `cover.jpg`, it will be treated like a cover and presented first in the gallery view page and as a gallery cover in the gallery list view. If more than one images match the name the first one found in natural sort order is selected.

Images can be added to a gallery by navigating to the gallery's page, selecting the "Add" tab, querying for and selecting the images to add, then selecting "Add to Gallery" from the `...` menu button. Likewise, images may be removed from a gallery by selecting the "Images" tab, selecting the images to remove and selecting "Remove from Gallery" from the `...` menu button.


# Images and Galleries

Images are the parts which make up galleries, but you can also have them be scanned independently. To declare an image part of a gallery, there are four ways:

1. Group them in a folder together and activate the **Create galleries from folders containing images** option in the library section of your settings. The gallery will get the name of the folder.
2. Group them in a folder together and create a file in the folder called .forcegallery. The gallery will get the name of the folder.
3. Group them into a zip archive together. The gallery will get the name of the archive.
4. You can simply create a gallery in stash itself by clicking on **New** in the Galleries tab. 

You can add images to every gallery manually in the gallery detail page. Deleting can be done by selecting the according images in the same view and clicking on the minus next to the edit button.

For best results, images in zip file should be stored without compression (copy, store or no compression options depending on the software you use. Eg on linux: `zip -0 -r gallery.zip foldertozip/`). This impacts **heavily** on the zip read performance.

> **:warning: Note:** AVIF files in ZIP archives are currently unsupported.

If a filename of an image in the gallery zip file ends with `cover.jpg`, it will be treated like a cover and presented first in the gallery view page and as a gallery cover in the gallery list view. If more than one images match the name the first one found in natural sort order is selected.

You can also manually select any image from a gallery as its cover. On the gallery details page, select the desired cover image, and then select **Set as Cover** in the ⋯ menu.

## Image clips/gifs

Images can also be clips/gifs. These are meant to be short video loops. Right now they are not possible in zipfiles. To declare video files to be images, there are two ways:

1. Deactivate video scanning for all libraries that contain clips/gifs, but keep image scanning active. Set the **Scan Video Extensions as Image Clip** option in the library section of your settings. 
2. Make sure none of the file endings used by your clips/gifs are present in the **Video Extensions** and add them to the **Image Extensions** in the library section of your settings.

A clip/gif will be a stillframe in the wall and grid view by default. To view the loop, you can go into the Lightbox Carousel (e.g. by clicking on an image in the wall view) or the image detail page.

If you want the loop to be used as a preview on the wall and grid view, you will have to generate them. 
You can do this as you scan for the new clip file by activating **Generate previews for image clips** on the scan settings, or do it after by going to the **Generated Content** section in the task section of your settings, activating **Image Clip Previews** and clicking generate. This takes a while, as the files are transcoded.


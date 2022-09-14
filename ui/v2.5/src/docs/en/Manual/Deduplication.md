# Dupe Checker

[The dupe checker](/sceneDuplicateChecker) searches your collection for scenes that are perceptually similar. This means that the files don't need to be identical, and will be identified even with different bitrates, resolutions, and intros/outros.

To achieve this stash needs to generate what's called a phash, or perceptual hash. Similar to sprite generation stash will generate a set of 25 images from fixed points in the scene. These images will be stitched together, and then hashed using the phash algorithm. The phash can then be used to find scenes that are the same or similar to others in the database. Phash generation can be run during scan, or as a separate task. Note that generation can take a while due to the work involved with extracting screenshots.

The dupe checker can be run with four different levels of accuracy. `Exact` looks for scenes that have exactly the same phash. This is a fast and accurate operation that should not yield any false positives except in very rare cases. The other accuracy levels look for duplicate files within a set distance of each other. This means the scenes don't have exactly the same phash, but are very similar. `High` and `Medium` should still yield very good results with few or no false positives. `Low` is likely to produce some false positives, but might still be useful for finding dupes.

Note that to generate a phash stash requires an uncorrupted file. If any errors are encountered during sprite generation the phash will not be generated. This is to prevent false positives.

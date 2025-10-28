# Scene Tagger

Stash can be integrated with stash-box which acts as a centralized metadata database. This is in the early stages of development but can be used for fingerprint/keyword lookups and automated tagging of performers and scenes. The batch tagging interface can be accessed from the [scene view](/scenes?disp=3). For more information join our [Discord](https://discord.gg/2TsNFKt).

## Searching 

The fingerprint search matches your current selection of files against the remote stash-box instance. Any scenes with a matching fingerprint will be returned, although there is currently no validation of fingerprints so it&rsquo;s recommended to double-check the validity before saving.

If no fingerprint match is found it&rsquo;s possible to search by keywords. The search works by matching the query against a scene&rsquo;s _title_, _release date_, _studio name_, and _performer names_. By default the tagger uses metadata set on the file, or parses the filename, this can be changed in the config.

An important thing to note is that it only returns a match *if all query terms are a match*. As an example, if a scene is titled `"A Trip to the Mall"` with the performer `"Jane Doe"`, a search for `"Trip to the Mall 1080p"` will *not* match, however `"trip mall doe"` would. Usually a few pieces of info is enough, for instance performer name + release date or studio name. To avoid common non-related keywords you can add them to the blacklist in the tagger config. Any items in the blacklist are stripped out of the query.

## Saving
When a scene is matched stash will try to match the studio and performers against your local studios and performers. If you have previously matched them, they will automatically be selected. If not you either have to select the correct performer/studio from the dropdown, choose create to create a new entity, or skip to ignore it.

Once a scene is saved the scene and the matched studio/performers will have the `stash_id` saved which will then be used for future tagging.

By default male performers are not shown, this can be enabled in the tagger config. Likewise scene tags are by default not saved. They can be set to either merge with existing tags on the scene, or overwrite them. It is not recommended to set tags currently since they are hard to deduplicate and can litter your data.

## Submitting fingerprints
After a scene is saved you will prompted to submit the fingerprint back to the stash-box instance. This is optional, but can be helpful for other users who have an identical copy who will then be able to match via the fingerprint search. No other information than the `stash_id` and file fingerprint is submitted.

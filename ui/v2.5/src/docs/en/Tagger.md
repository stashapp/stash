# Scene Tagger

The search works by matching the query against a scene&rsquo;s  title_, release date_, _studio name_, and _performer names_. An important thing to note is that it only returns a match *if all query terms are a match*.

As an example, if a scene is titled `"A Trip to the Mall"`, a search for `"Trip to the Mall 1080p"` will *not* match, however `"trip mall"` would. Usually a few pieces of info is enough, for instance performer name + release date or studio name.

# Browsing

## Querying and Filtering

### Keyword searching

The text field allows you to search using keywords. Keyword searching matches on different fields depending on the object type:

| Type | Fields searched |
|------|-----------------|
| Scene | Title, Details, Path, OSHash, Checksum, Marker titles |
| Image | Title, Path, Checksum |
| Group | Title |
| Marker | Title, Scene title |
| Gallery | Title, Path, Checksum |
| Performer | Name, Aliases |
| Studio | Name, Aliases |
| Tag | Name, Aliases |

Keyword matching uses the following rules:

* all words are required in the matching field. For example, `foo bar` matches scenes with both `foo` and `bar` in the title.
* the `or` keyword or symbol (`|`) is used to match either fields. For example, `foo or bar` (or `foo | bar`) matches scenes with `foo` or `bar` in the title. Or sets can be combined. For example, `foo or bar or baz xyz or zyx` matches scenes with one of `foo`, `bar` and `baz`, *and* `xyz` or `zyx`.
* the not symbol (`-`) is used to exclude terms. For example, `foo -bar` matches scenes with `foo` and excludes those with `bar`. The not symbol cannot be combined with an or operand. That is, `-foo or bar` will be interpreted to match `-foo` or `bar`. On the other hand, `foo or bar -baz` will match `foo` or `bar` and exclude `baz`.
* surrounding a phrase in quotes (`"`) matches on that exact phrase. For example, `"foo bar"` matches scenes with `foo bar` in the title. Quotes may also be used to escape the keywords and symbols. For example, `foo "-bar"` will match scenes with `foo` and `-bar`.
* quoted phrases may be used with the or and not operators. For example, `"foo bar" or baz -"xyz zyx"` will match scenes with `foo bar` *or* `baz`, and exclude those with `xyz zyx`.
* `or` keywords or symbols at the start or end of a line will be treated literally. That is, `or foo` will match scenes with `or` and `foo`.
* all keyword matching is case-insensitive

### Filters

Filters can be accessed by clicking the filter button on the right side of the query text field. 

Note that only one filter criterion per criterion type may be assigned.

#### Regex modifiers

Some filters have regex modifier as an option. Regex modifiers are case-sensitive by default.

### Sorting and page size

The current sorting field is shown next to the query text field, indicating the current sort field and order. The page size dropdown allows selecting from a standard set of objects per page, and allows setting a custom page size.

### Saved filters

Saved filters can be accessed with the bookmark button on the left of the query text field. The current filter can be saved by entering a filter name and clicking on the save button. Existing saved filters may be overwritten with the current filter by clicking on the save button next to the filter name. Saved filters may also be deleted by pressing the delete button next to the filter name.

Saved filters are sorted alphabetically by title with capitalized titles sorted first.

### Default filter

The default filter for the top-level pages may be set to the current filter by clicking the `Set as default` button in the saved filter menu.

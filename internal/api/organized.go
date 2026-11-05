package api

// skipOrganizedInBulk reports whether a bulk edit must leave this item alone.
// Bulk edits never change organized items, except to unset organized or to
// change rating or favorite, which are the user's opinion rather than item data.
func skipOrganizedInBulk(organized bool, translator changesetTranslator, newOrganized *bool) bool {
	if !organized {
		return false
	}

	if translator.hasField("organized") && newOrganized != nil && !*newOrganized {
		return false
	}

	for _, f := range translator.getFields() {
		switch f {
		case "ids", "organized", "rating100", "favorite":
		default:
			return true
		}
	}

	return false
}

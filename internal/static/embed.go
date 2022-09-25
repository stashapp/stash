package static

import "embed"

//go:embed performer
var Performer embed.FS

//go:embed performer_male
var PerformerMale embed.FS

// scene.png modified from https://github.com/FortAwesome/Font-Awesome/blob/6.x/svgs/regular/circle-play.svg
// Font Awesome Free 6.1.1 by @fontawesome - https://fontawesome.com
// License CC BY 4.0 (https://creativecommons.org/licenses/by/4.0/)
// Converted to png and resized to 1280x720.

//go:embed scene
var Scene embed.FS

package models

import "time"

type Tag struct {
	ID            int             `db:"id" json:"id"`
	Name          string          `db:"name" json:"name"` // TODO make schema not null
	IgnoreAutoTag bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
	CreatedAt     SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type TagPartial struct {
	ID            int              `db:"id" json:"id"`
	Name          *string          `db:"name" json:"name"` // TODO make schema not null
	IgnoreAutoTag *bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
	CreatedAt     *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type TagPath struct {
	Tag
	Path string `db:"path" json:"path"`
}

func NewTag(name string) *Tag {
	currentTime := time.Now()
	return &Tag{
		Name:      name,
		CreatedAt: SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: SQLiteTimestamp{Timestamp: currentTime},
	}
}

type Tags []*Tag

func (t *Tags) Append(o interface{}) {
	*t = append(*t, o.(*Tag))
}

func (t *Tags) New() interface{} {
	return &Tag{}
}

type TagPaths []*TagPath

func (t *TagPaths) Append(o interface{}) {
	*t = append(*t, o.(*TagPath))
}

func (t *TagPaths) New() interface{} {
	return &TagPath{}
}

// Original Tag image from: https://fontawesome.com/icons/tag?style=solid
// Modified to change color and rotate
// Licensed under CC Attribution 4.0: https://fontawesome.com/license
var DefaultTagImage = []byte(`<svg
   xmlns:dc="http://purl.org/dc/elements/1.1/"
   xmlns:cc="http://creativecommons.org/ns#"
   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
   xmlns:svg="http://www.w3.org/2000/svg"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   width="200"
   height="200"
   id="svg2"
   version="1.1"
   inkscape:version="0.48.4 r9939"
   sodipodi:docname="tag.svg">
  <defs
     id="defs4" />
  <sodipodi:namedview
     id="base"
     pagecolor="#000000"
     bordercolor="#666666"
     borderopacity="1.0"
     inkscape:pageopacity="1"
     inkscape:pageshadow="2"
     inkscape:zoom="1"
     inkscape:cx="181.77771"
     inkscape:cy="279.72376"
     inkscape:document-units="px"
     inkscape:current-layer="layer1"
     showgrid="false"
     fit-margin-top="0"
     fit-margin-left="0"
     fit-margin-right="0"
     fit-margin-bottom="0"
     inkscape:window-width="1920"
     inkscape:window-height="1017"
     inkscape:window-x="-8"
     inkscape:window-y="-8"
     inkscape:window-maximized="1" />
  <metadata
     id="metadata7">
    <rdf:RDF>
      <cc:Work
         rdf:about="">
        <dc:format>image/svg+xml</dc:format>
        <dc:type
           rdf:resource="http://purl.org/dc/dcmitype/StillImage" />
        <dc:title></dc:title>
      </cc:Work>
    </rdf:RDF>
  </metadata>
  <g
     inkscape:label="Layer 1"
     inkscape:groupmode="layer"
     id="layer1"
     transform="translate(-157.84358,-524.69522)">
    <path
       id="path2987"
       d="m 229.94314,669.26549 -36.08466,-36.08466 c -4.68653,-4.68653 -4.68653,-12.28468 0,-16.97121 l 36.08466,-36.08467 a 12.000453,12.000453 0 0 1 8.4856,-3.5148 l 74.91443,0 c 6.62761,0 12.00041,5.3728 12.00041,12.00041 l 0,72.16933 c 0,6.62761 -5.3728,12.00041 -12.00041,12.00041 l -74.91443,0 a 12.000453,12.000453 0 0 1 -8.4856,-3.51481 z m -13.45639,-53.05587 c -4.68653,4.68653 -4.68653,12.28468 0,16.97121 4.68652,4.68652 12.28467,4.68652 16.9712,0 4.68653,-4.68653 4.68653,-12.28468 0,-16.97121 -4.68653,-4.68652 -12.28468,-4.68652 -16.9712,0 z"
       inkscape:connector-curvature="0"
       style="fill:#ffffff;fill-opacity:1" />
  </g>
</svg>`)

// var DefaultTagImage = []byte(`<svg
// xmlns:dc="http://purl.org/dc/elements/1.1/"
// xmlns:cc="http://creativecommons.org/ns#"
// xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
// xmlns:svg="http://www.w3.org/2000/svg"
// xmlns="http://www.w3.org/2000/svg"
// xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
// xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
// width="600"
// height="600"
// id="svg2"
// version="1.1"
// inkscape:version="0.48.4 r9939"
// sodipodi:docname="New document 1">
// <defs
//   id="defs4" />
// <sodipodi:namedview
//   id="base"
//   pagecolor="#000000"
//   bordercolor="#666666"
//   borderopacity="1.0"
//   inkscape:pageopacity="1"
//   inkscape:pageshadow="2"
//   inkscape:zoom="0.82173542"
//   inkscape:cx="181.77771"
//   inkscape:cy="159.72376"
//   inkscape:document-units="px"
//   inkscape:current-layer="layer1"
//   showgrid="false"
//   fit-margin-top="0"
//   fit-margin-left="0"
//   fit-margin-right="0"
//   fit-margin-bottom="0"
//   inkscape:window-width="1920"
//   inkscape:window-height="1017"
//   inkscape:window-x="-8"
//   inkscape:window-y="-8"
//   inkscape:window-maximized="1" />
// <metadata
//   id="metadata7">
//  <rdf:RDF>
//    <cc:Work
// 	  rdf:about="">
// 	 <dc:format>image/svg+xml</dc:format>
// 	 <dc:type
// 		rdf:resource="http://purl.org/dc/dcmitype/StillImage" />
// 	 <dc:title></dc:title>
//    </cc:Work>
//  </rdf:RDF>
// </metadata>
// <g
//   inkscape:label="Layer 1"
//   inkscape:groupmode="layer"
//   id="layer1"
//   transform="translate(-157.84358,-124.69522)">
//  <path
// 	id="path2987"
// 	d="M 346.24605,602.96957 201.91282,458.63635 c -18.7454,-18.7454 -18.7454,-49.13685 0,-67.88225 L 346.24605,246.42087 a 48,48 0 0 1 33.94111,-14.05869 l 299.64641,0 c 26.50943,0 47.99982,21.49039 47.99982,47.99982 l 0,288.66645 c 0,26.50943 -21.49039,47.99982 -47.99983,47.99982 l -299.64639,0 a 48,48 0 0 1 -33.94112,-14.0587 z M 292.42249,390.7541 c -18.7454,18.7454 -18.7454,49.13685 0,67.88225 18.7454,18.7454 49.13685,18.7454 67.88225,0 18.7454,-18.7454 18.7454,-49.13685 0,-67.88225 -18.7454,-18.7454 -49.13685,-18.7454 -67.88225,0 z"
// 	inkscape:connector-curvature="0"
// 	style="fill:#ffffff;fill-opacity:1" />
// </g>
// </svg>`)

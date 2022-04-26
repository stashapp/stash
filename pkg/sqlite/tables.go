package sqlite

import (
	"github.com/doug-martin/goqu/v9"

	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

var dialect = goqu.Dialect("sqlite3")

var (
	galleriesImagesJoinTable  = goqu.T(galleriesImagesTable)
	imagesTagsJoinTable       = goqu.T(imagesTagsTable)
	performersImagesJoinTable = goqu.T(performersImagesTable)

	galleriesTagsJoinTable       = goqu.T(galleriesTagsTable)
	performersGalleriesJoinTable = goqu.T(performersGalleriesTable)
	galleriesScenesJoinTable     = goqu.T(galleriesScenesTable)

	scenesTagsJoinTable       = goqu.T(scenesTagsTable)
	scenesPerformersJoinTable = goqu.T(performersScenesTable)
	scenesStashIDsJoinTable   = goqu.T("scene_stash_ids")
	scenesMoviesJoinTable     = goqu.T(moviesScenesTable)

	filesFingerprintsJoinTable = goqu.T(filesFingerprintsTable)
)

var (
	imageTableMgr = &table{
		table:    goqu.T(imageTable),
		idColumn: goqu.T(imageTable).Col(idColumn),
	}

	imageGalleriesTableMgr = &joinTable{
		table: table{
			table:    galleriesImagesJoinTable,
			idColumn: galleriesImagesJoinTable.Col(imageIDColumn),
		},
		fkColumn: galleriesImagesJoinTable.Col(galleryIDColumn),
	}

	imagesTagsTableMgr = &joinTable{
		table: table{
			table:    imagesTagsJoinTable,
			idColumn: imagesTagsJoinTable.Col(imageIDColumn),
		},
		fkColumn: imagesTagsJoinTable.Col(tagIDColumn),
	}

	imagesPerformersTableMgr = &joinTable{
		table: table{
			table:    performersImagesJoinTable,
			idColumn: performersImagesJoinTable.Col(imageIDColumn),
		},
		fkColumn: performersImagesJoinTable.Col(performerIDColumn),
	}
)

var (
	galleryTableMgr = &table{
		table:    goqu.T(galleryTable),
		idColumn: goqu.T(galleryTable).Col(idColumn),
	}

	galleriesTagsTableMgr = &joinTable{
		table: table{
			table:    galleriesTagsJoinTable,
			idColumn: galleriesTagsJoinTable.Col(galleryIDColumn),
		},
		fkColumn: galleriesTagsJoinTable.Col(tagIDColumn),
	}

	galleriesPerformersTableMgr = &joinTable{
		table: table{
			table:    performersGalleriesJoinTable,
			idColumn: performersGalleriesJoinTable.Col(galleryIDColumn),
		},
		fkColumn: performersGalleriesJoinTable.Col(performerIDColumn),
	}

	galleriesScenesTableMgr = &joinTable{
		table: table{
			table:    galleriesScenesJoinTable,
			idColumn: galleriesScenesJoinTable.Col(galleryIDColumn),
		},
		fkColumn: galleriesScenesJoinTable.Col(sceneIDColumn),
	}
)

var (
	sceneTableMgr = &table{
		table:    goqu.T(sceneTable),
		idColumn: goqu.T(sceneTable).Col(idColumn),
	}

	scenesTagsTableMgr = &joinTable{
		table: table{
			table:    scenesTagsJoinTable,
			idColumn: scenesTagsJoinTable.Col(sceneIDColumn),
		},
		fkColumn: scenesTagsJoinTable.Col(tagIDColumn),
	}

	scenesPerformersTableMgr = &joinTable{
		table: table{
			table:    scenesPerformersJoinTable,
			idColumn: scenesPerformersJoinTable.Col(sceneIDColumn),
		},
		fkColumn: scenesPerformersJoinTable.Col(performerIDColumn),
	}

	scenesGalleriesTableMgr = galleriesScenesTableMgr.invert()

	scenesStashIDsTableMgr = &stashIDTable{
		table: table{
			table:    scenesStashIDsJoinTable,
			idColumn: scenesStashIDsJoinTable.Col(sceneIDColumn),
		},
	}

	scenesMoviesTableMgr = &scenesMoviesTable{
		table: table{
			table:    scenesMoviesJoinTable,
			idColumn: scenesMoviesJoinTable.Col(sceneIDColumn),
		},
	}
)

var (
	fileTableMgr = &table{
		table:    goqu.T(fileTable),
		idColumn: goqu.T(fileTable).Col(idColumn),
	}

	videoFileTableMgr = &table{
		table:    goqu.T(videoFileTable),
		idColumn: goqu.T(videoFileTable).Col(idColumn),
	}

	imageFileTableMgr = &table{
		table:    goqu.T(imageFileTable),
		idColumn: goqu.T(imageFileTable).Col(idColumn),
	}

	folderTableMgr = &table{
		table:    goqu.T(folderTable),
		idColumn: goqu.T(folderTable).Col(idColumn),
	}

	fingerprintTableMgr = &table{
		table:    goqu.T(fingerprintTable),
		idColumn: goqu.T(fingerprintTable).Col(idColumn),
	}

	filesFingerprintsTableMgr = &joinTable{
		table: table{
			table:    filesFingerprintsJoinTable,
			idColumn: filesFingerprintsJoinTable.Col(fileIDColumn),
		},
		fkColumn: filesFingerprintsJoinTable.Col(fingerprintIDColumn),
	}
)

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
	imagesFilesJoinTable      = goqu.T(imagesFilesTable)

	galleriesFilesJoinTable      = goqu.T(galleriesFilesTable)
	galleriesTagsJoinTable       = goqu.T(galleriesTagsTable)
	performersGalleriesJoinTable = goqu.T(performersGalleriesTable)
	galleriesScenesJoinTable     = goqu.T(galleriesScenesTable)

	scenesFilesJoinTable      = goqu.T(scenesFilesTable)
	scenesTagsJoinTable       = goqu.T(scenesTagsTable)
	scenesPerformersJoinTable = goqu.T(performersScenesTable)
	scenesStashIDsJoinTable   = goqu.T("scene_stash_ids")
	scenesMoviesJoinTable     = goqu.T(moviesScenesTable)

	performersAliasesJoinTable  = goqu.T(performersAliasesTable)
	performersTagsJoinTable     = goqu.T(performersTagsTable)
	performersStashIDsJoinTable = goqu.T("performer_stash_ids")
)

var (
	imageTableMgr = &table{
		table:    goqu.T(imageTable),
		idColumn: goqu.T(imageTable).Col(idColumn),
	}

	imagesFilesTableMgr = &relatedFilesTable{
		table: table{
			table:    imagesFilesJoinTable,
			idColumn: imagesFilesJoinTable.Col(imageIDColumn),
		},
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

	galleriesFilesTableMgr = &relatedFilesTable{
		table: table{
			table:    galleriesFilesJoinTable,
			idColumn: galleriesFilesJoinTable.Col(galleryIDColumn),
		},
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

	sceneMarkerTableMgr = &table{
		table:    goqu.T(sceneMarkerTable),
		idColumn: goqu.T(sceneMarkerTable).Col(idColumn),
	}

	scenesFilesTableMgr = &relatedFilesTable{
		table: table{
			table:    scenesFilesJoinTable,
			idColumn: scenesFilesJoinTable.Col(sceneIDColumn),
		},
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
		idColumn: goqu.T(videoFileTable).Col(fileIDColumn),
	}

	imageFileTableMgr = &table{
		table:    goqu.T(imageFileTable),
		idColumn: goqu.T(imageFileTable).Col(fileIDColumn),
	}

	folderTableMgr = &table{
		table:    goqu.T(folderTable),
		idColumn: goqu.T(folderTable).Col(idColumn),
	}

	fingerprintTableMgr = &table{
		table:    goqu.T(fingerprintTable),
		idColumn: goqu.T(fingerprintTable).Col(idColumn),
	}
)

var (
	performerTableMgr = &table{
		table:    goqu.T(performerTable),
		idColumn: goqu.T(performerTable).Col(idColumn),
	}

	performersAliasesTableMgr = &stringTable{
		table: table{
			table:    performersAliasesJoinTable,
			idColumn: performersAliasesJoinTable.Col(performerIDColumn),
		},
		stringColumn: performersAliasesJoinTable.Col(performerAliasColumn),
	}

	performersTagsTableMgr = &joinTable{
		table: table{
			table:    performersTagsJoinTable,
			idColumn: performersTagsJoinTable.Col(performerIDColumn),
		},
		fkColumn: performersTagsJoinTable.Col(tagIDColumn),
	}

	performersStashIDsTableMgr = &stashIDTable{
		table: table{
			table:    performersStashIDsJoinTable,
			idColumn: performersStashIDsJoinTable.Col(performerIDColumn),
		},
	}
)

var (
	studioTableMgr = &table{
		table:    goqu.T(studioTable),
		idColumn: goqu.T(studioTable).Col(idColumn),
	}
)

var (
	tagTableMgr = &table{
		table:    goqu.T(tagTable),
		idColumn: goqu.T(tagTable).Col(idColumn),
	}
)

var (
	movieTableMgr = &table{
		table:    goqu.T(movieTable),
		idColumn: goqu.T(movieTable).Col(idColumn),
	}
)

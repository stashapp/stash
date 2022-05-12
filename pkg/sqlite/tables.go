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
	imagesQueryTable          = goqu.T("images_query")
	galleriesQueryTable       = goqu.T("galleries_query")
	scenesQueryTable          = goqu.T("scenes_query")

	galleriesFilesJoinTable      = goqu.T(galleriesFilesTable)
	galleriesTagsJoinTable       = goqu.T(galleriesTagsTable)
	performersGalleriesJoinTable = goqu.T(performersGalleriesTable)
	galleriesScenesJoinTable     = goqu.T(galleriesScenesTable)

	scenesFilesJoinTable      = goqu.T(scenesFilesTable)
	scenesTagsJoinTable       = goqu.T(scenesTagsTable)
	scenesPerformersJoinTable = goqu.T(performersScenesTable)
	scenesStashIDsJoinTable   = goqu.T("scene_stash_ids")
	scenesMoviesJoinTable     = goqu.T(moviesScenesTable)
)

var (
	imageTableMgr = &table{
		table:    goqu.T(imageTable),
		idColumn: goqu.T(imageTable).Col(idColumn),
	}

	imageQueryTableMgr = &table{
		table:    imagesQueryTable,
		idColumn: imagesQueryTable.Col(idColumn),
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

	galleryQueryTableMgr = &table{
		table:    galleriesQueryTable,
		idColumn: galleriesQueryTable.Col(idColumn),
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

	sceneQueryTableMgr = &table{
		table:    scenesQueryTable,
		idColumn: scenesQueryTable.Col(idColumn),
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

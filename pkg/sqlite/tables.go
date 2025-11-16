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
	imagesURLsJoinTable       = goqu.T(imagesURLsTable)

	galleriesFilesJoinTable      = goqu.T(galleriesFilesTable)
	galleriesTagsJoinTable       = goqu.T(galleriesTagsTable)
	performersGalleriesJoinTable = goqu.T(performersGalleriesTable)
	galleriesScenesJoinTable     = goqu.T(galleriesScenesTable)
	galleriesURLsJoinTable       = goqu.T(galleriesURLsTable)

	scenesFilesJoinTable      = goqu.T(scenesFilesTable)
	scenesTagsJoinTable       = goqu.T(scenesTagsTable)
	scenesPerformersJoinTable = goqu.T(performersScenesTable)
	scenesStashIDsJoinTable   = goqu.T("scene_stash_ids")
	scenesGroupsJoinTable     = goqu.T(groupsScenesTable)
	scenesURLsJoinTable       = goqu.T(scenesURLsTable)

	sceneMarkersTagsJoinTable = goqu.T(sceneMarkersTagsTable)

	performersAliasesJoinTable  = goqu.T(performersAliasesTable)
	performersURLsJoinTable     = goqu.T(performerURLsTable)
	performersTagsJoinTable     = goqu.T(performersTagsTable)
	performersStashIDsJoinTable = goqu.T("performer_stash_ids")
	performersCustomFieldsTable = goqu.T("performer_custom_fields")

	studiosAliasesJoinTable  = goqu.T(studioAliasesTable)
	studiosURLsJoinTable     = goqu.T(studioURLsTable)
	studiosTagsJoinTable     = goqu.T(studiosTagsTable)
	studiosStashIDsJoinTable = goqu.T("studio_stash_ids")

	groupsURLsJoinTable     = goqu.T(groupURLsTable)
	groupsTagsJoinTable     = goqu.T(groupsTagsTable)
	groupRelationsJoinTable = goqu.T(groupRelationsTable)

	tagsAliasesJoinTable  = goqu.T(tagAliasesTable)
	tagRelationsJoinTable = goqu.T(tagRelationsTable)
	tagsStashIDsJoinTable = goqu.T("tag_stash_ids")
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

	imageGalleriesTableMgr = &imageGalleriesTable{
		joinTable: joinTable{
			table: table{
				table:    galleriesImagesJoinTable,
				idColumn: galleriesImagesJoinTable.Col(imageIDColumn),
			},
			fkColumn: galleriesImagesJoinTable.Col(galleryIDColumn),
		},
	}

	imagesTagsTableMgr = &joinTable{
		table: table{
			table:    imagesTagsJoinTable,
			idColumn: imagesTagsJoinTable.Col(imageIDColumn),
		},
		fkColumn:     imagesTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
	}

	imagesPerformersTableMgr = &joinTable{
		table: table{
			table:    performersImagesJoinTable,
			idColumn: performersImagesJoinTable.Col(imageIDColumn),
		},
		fkColumn: performersImagesJoinTable.Col(performerIDColumn),
	}

	imagesURLsTableMgr = &orderedValueTable[string]{
		table: table{
			table:    imagesURLsJoinTable,
			idColumn: imagesURLsJoinTable.Col(imageIDColumn),
		},
		valueColumn: imagesURLsJoinTable.Col(imageURLColumn),
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
		fkColumn:     galleriesTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
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

	galleriesChaptersTableMgr = &table{
		table:    goqu.T(galleriesChaptersTable),
		idColumn: goqu.T(galleriesChaptersTable).Col(idColumn),
	}

	galleriesURLsTableMgr = &orderedValueTable[string]{
		table: table{
			table:    galleriesURLsJoinTable,
			idColumn: galleriesURLsJoinTable.Col(galleryIDColumn),
		},
		valueColumn: galleriesURLsJoinTable.Col(galleriesURLColumn),
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

	sceneMarkersTagsTableMgr = &joinTable{
		table: table{
			table:    sceneMarkersTagsJoinTable,
			idColumn: sceneMarkersTagsJoinTable.Col(sceneMarkerIDColumn),
		},
		fkColumn:     sceneMarkersTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
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
		fkColumn:     scenesTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
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

	scenesGroupsTableMgr = &scenesGroupsTable{
		table: table{
			table:    scenesGroupsJoinTable,
			idColumn: scenesGroupsJoinTable.Col(sceneIDColumn),
		},
	}

	scenesURLsTableMgr = &orderedValueTable[string]{
		table: table{
			table:    scenesURLsJoinTable,
			idColumn: scenesURLsJoinTable.Col(sceneIDColumn),
		},
		valueColumn: scenesURLsJoinTable.Col(sceneURLColumn),
	}

	scenesViewTableMgr = &viewHistoryTable{
		table: table{
			table:    goqu.T(scenesViewDatesTable),
			idColumn: goqu.T(scenesViewDatesTable).Col(sceneIDColumn),
		},
		dateColumn: goqu.T(scenesViewDatesTable).Col(sceneViewDateColumn),
	}

	scenesOTableMgr = &viewHistoryTable{
		table: table{
			table:    goqu.T(scenesODatesTable),
			idColumn: goqu.T(scenesODatesTable).Col(sceneIDColumn),
		},
		dateColumn: goqu.T(scenesODatesTable).Col(sceneODateColumn),
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

	performersURLsTableMgr = &orderedValueTable[string]{
		table: table{
			table:    performersURLsJoinTable,
			idColumn: performersURLsJoinTable.Col(performerIDColumn),
		},
		valueColumn: performersURLsJoinTable.Col(performerURLColumn),
	}

	performersTagsTableMgr = &joinTable{
		table: table{
			table:    performersTagsJoinTable,
			idColumn: performersTagsJoinTable.Col(performerIDColumn),
		},
		fkColumn:     performersTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
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

	studiosAliasesTableMgr = &stringTable{
		table: table{
			table:    studiosAliasesJoinTable,
			idColumn: studiosAliasesJoinTable.Col(studioIDColumn),
		},
		stringColumn: studiosAliasesJoinTable.Col(studioAliasColumn),
	}

	studiosURLsTableMgr = &orderedValueTable[string]{
		table: table{
			table:    studiosURLsJoinTable,
			idColumn: studiosURLsJoinTable.Col(studioIDColumn),
		},
		valueColumn: studiosURLsJoinTable.Col(studioURLColumn),
	}

	studiosTagsTableMgr = &joinTable{
		table: table{
			table:    studiosTagsJoinTable,
			idColumn: studiosTagsJoinTable.Col(studioIDColumn),
		},
		fkColumn:     studiosTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
	}

	studiosStashIDsTableMgr = &stashIDTable{
		table: table{
			table:    studiosStashIDsJoinTable,
			idColumn: studiosStashIDsJoinTable.Col(studioIDColumn),
		},
	}
)

var (
	tagTableMgr = &table{
		table:    goqu.T(tagTable),
		idColumn: goqu.T(tagTable).Col(idColumn),
	}

	// formerly: goqu.COALESCE(tagTableMgr.table.Col("sort_name"), tagTableMgr.table.Col("name")).Asc()
	tagTableSort    = goqu.L("COALESCE(tags.sort_name, tags.name) COLLATE NATURAL_CI").Asc()
	tagTableSortSQL = "COALESCE(tags.sort_name, tags.name) COLLATE NATURAL_CI ASC"

	tagsAliasesTableMgr = &stringTable{
		table: table{
			table:    tagsAliasesJoinTable,
			idColumn: tagsAliasesJoinTable.Col(tagIDColumn),
		},
		stringColumn: tagsAliasesJoinTable.Col(tagAliasColumn),
	}

	tagsParentTagsTableMgr = &joinTable{
		table: table{
			table:    tagRelationsJoinTable,
			idColumn: tagRelationsJoinTable.Col(tagChildIDColumn),
		},
		fkColumn:     tagRelationsJoinTable.Col(tagParentIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
	}

	tagsChildTagsTableMgr = *tagsParentTagsTableMgr.invert()

	tagsStashIDsTableMgr = &stashIDTable{
		table: table{
			table:    tagsStashIDsJoinTable,
			idColumn: tagsStashIDsJoinTable.Col(tagIDColumn),
		},
	}
)

var (
	groupTableMgr = &table{
		table:    goqu.T(groupTable),
		idColumn: goqu.T(groupTable).Col(idColumn),
	}

	groupsURLsTableMgr = &orderedValueTable[string]{
		table: table{
			table:    groupsURLsJoinTable,
			idColumn: groupsURLsJoinTable.Col(groupIDColumn),
		},
		valueColumn: groupsURLsJoinTable.Col(groupURLColumn),
	}

	groupsTagsTableMgr = &joinTable{
		table: table{
			table:    groupsTagsJoinTable,
			idColumn: groupsTagsJoinTable.Col(groupIDColumn),
		},
		fkColumn:     groupsTagsJoinTable.Col(tagIDColumn),
		foreignTable: tagTableMgr,
		orderBy:      tagTableSort,
	}

	groupRelationshipTableMgr = &table{
		table: groupRelationsJoinTable,
	}
)

var (
	blobTableMgr = &table{
		table:    goqu.T(blobTable),
		idColumn: goqu.T(blobTable).Col(blobChecksumColumn),
	}
)

var (
	savedFilterTableMgr = &table{
		table:    goqu.T(savedFilterTable),
		idColumn: goqu.T(savedFilterTable).Col(idColumn),
	}
)

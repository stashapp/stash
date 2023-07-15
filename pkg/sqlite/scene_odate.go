package sqlite

const sceneODateTable = "scenes_odates"

type SceneODateStore struct {
	repository
	tableMgr *table
	oDateManager
}

func NewSceneODateStore() *SceneODateStore {
	return &SceneODateStore{
		repository: repository{
			tableName: sceneODateTable,
			idColumn:  idColumn,
		},
		tableMgr:     sceneODateTableMgr,
		oDateManager: oDateManager{sceneODateTableMgr},
	}
}

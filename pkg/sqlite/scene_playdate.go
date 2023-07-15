package sqlite

const scenePlayDateTable = "scenes_playdates"

type ScenePlayDateStore struct {
	repository
	tableMgr *table
	playDateManager
}

func NewScenePlayDateStore() *ScenePlayDateStore {
	return &ScenePlayDateStore{
		repository: repository{
			tableName: scenePlayDateTable,
			idColumn:  idColumn,
		},
		tableMgr:        scenePlayDateTableMgr,
		playDateManager: playDateManager{scenePlayDateTableMgr},
	}
}

import { Injectable } from '@angular/core';

import { StudioListState, SceneListState } from '../shared/models/list-state.model';

@Injectable()
export class StudiosService {
  listState: StudioListState = new StudioListState();
  detailsSceneListState: SceneListState = new SceneListState();

  constructor() {}
}

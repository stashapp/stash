import { Injectable } from '@angular/core';

import { SceneListState, SceneMarkerListState } from '../shared/models/list-state.model';

@Injectable()
export class ScenesService {
  sceneListState: SceneListState = new SceneListState();
  sceneMarkerListState: SceneMarkerListState = new SceneMarkerListState();

  constructor() {}
}

import { Injectable } from '@angular/core';

import { PerformerListState, SceneListState } from '../shared/models/list-state.model';

@Injectable()
export class PerformersService {
  performerListState: PerformerListState = new PerformerListState();
  detailsSceneListState: SceneListState = new SceneListState();

  constructor() {}
}

import { Component, OnInit } from '@angular/core';

import { ScenesService } from '../scenes.service';

@Component({
  selector: 'app-scene-list',
  template: '<app-list [state]="state"></app-list>'
})
export class SceneListComponent implements OnInit {
  state = this.scenesService.sceneListState;

  constructor(private scenesService: ScenesService) {}

  ngOnInit() {}

}

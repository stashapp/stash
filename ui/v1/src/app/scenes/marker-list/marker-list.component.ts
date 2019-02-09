import { Component, OnInit } from '@angular/core';

import { ScenesService } from '../scenes.service';

@Component({
  selector: 'app-marker-list',
  template: '<app-list [state]="state"></app-list>'
})
export class MarkerListComponent implements OnInit {
  state = this.scenesService.sceneMarkerListState;

  constructor(private scenesService: ScenesService) {}

  ngOnInit() {}

}

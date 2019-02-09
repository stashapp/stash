import { Component, OnInit, Input } from '@angular/core';

import { SceneMarkerData } from '../../core/graphql-generated';
import { BaseWallItemComponent } from '../base-wall-item/base-wall-item.component';

@Component({
  selector: 'app-scene-marker-wall-item',
  templateUrl: './scene-marker-wall-item.component.html'
})
export class SceneMarkerWallItemComponent extends BaseWallItemComponent implements OnInit {
  @Input() sceneMarker: SceneMarkerData.Fragment;

  constructor() { super(); }

  ngOnInit() {
    if (!!this.sceneMarker) {
      this.title = this.sceneMarker.title;
      this.imagePath = this.sceneMarker.preview;
      this.videoPath = this.sceneMarker.stream;
    } else {
      this.title = '';
      this.imagePath = '';
    }
  }
}

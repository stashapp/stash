import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';

import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

import { StashService } from '../../core/stash.service';

export enum WallMode {
  Scenes,
  Markers
}

@Component({
  selector: 'app-scene-wall',
  templateUrl: './scene-wall.component.html',
  styleUrls: ['./scene-wall.component.css']
})
export class SceneWallComponent implements OnInit {
  WallMode = WallMode;
  items: any[]; // scenes or scene markers
  markerOptions: any[];
  showingMarkerList = false;
  searchTerm = '';
  searchFormControl = new FormControl();
  mode: WallMode = WallMode.Markers;

  constructor(
    private stashService: StashService
  ) {}

  ngOnInit() {
    this.searchFormControl.valueChanges.pipe(
                            debounceTime(1000),
                            distinctUntilChanged()
                          ).subscribe(term => {
                            this.getScenes(term);
                          });
    this.stashService.markerStrings().valueChanges.subscribe(result => {
      this.markerOptions = result.data.markerStrings;
    });
    this.searchFormControl.setValue(this.searchTerm);
  }

  async getScenes(q: string) {
    this.items = null;
    this.searchTerm = q;
    if (this.mode === WallMode.Scenes) {
      const response = await this.stashService.sceneWall(q).result();
      this.items = response.data.sceneWall;
    } else {
      const response = await this.stashService.markerWall(q).result();
      this.items = response.data.markerWall;
    }
  }

  toggleMode() {
    if (this.mode === WallMode.Scenes) {
      this.mode = WallMode.Markers;
    } else {
      this.mode = WallMode.Scenes;
    }
    this.getScenes(this.searchTerm);
  }

  toggleMarkerList() {
    this.showingMarkerList = !this.showingMarkerList;
  }

  refresh() {
    this.getScenes(this.searchTerm);
  }

  onClickMarker(marker) {
    this.searchTerm = `${marker.title}`;
    this.searchFormControl.setValue(this.searchTerm);
    this.showingMarkerList = false;
  }

  async sortMarkers(by) {
    const result = await this.stashService.markerStrings(null, by).result();
    this.markerOptions = result.data.markerStrings;
  }

}

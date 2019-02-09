import { Component, OnInit, OnChanges, SimpleChanges, Input } from '@angular/core';
import { FormControl } from '@angular/forms';

import { StashService } from '../../core/stash.service';

import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

import { MarkerStrings, SceneMarkerData, SceneData, AllTagsForFilter } from '../../core/graphql-generated';


@Component({
  selector: 'app-scene-detail-marker-manager',
  templateUrl: './scene-detail-marker-manager.component.html',
  styleUrls: ['./scene-detail-marker-manager.component.css']
})
export class SceneDetailMarkerManagerComponent implements OnInit, OnChanges {
  @Input() scene: SceneData.Fragment;
  @Input() player: any;

  showingMarkerModal = false;
  markerOptions: MarkerStrings.Query['markerStrings'];
  filteredMarkerOptions: string[] = [];
  hasFocus = false;
  editingMarker: SceneMarkerData.Fragment;
  deleteClickCount = 0;

  searchFormControl = new FormControl();

  // Form input
  title: string;
  seconds: number;
  primary_tag_id: string;
  tag_ids: string[] = [];

  // From the network
  tags: AllTagsForFilter.AllTags[];

  constructor(private stashService: StashService) {}

  ngOnInit() {
    this.stashService.allTagsForFilter().valueChanges.subscribe(result => {
      this.tags = result.data.allTags;
    });

    this.stashService.markerStrings().valueChanges.subscribe(result => {
      this.markerOptions = result.data.markerStrings;
    });

    this.searchFormControl.valueChanges.pipe(
      debounceTime(400),
      distinctUntilChanged()
    ).subscribe(term => {
      this.filteredMarkerOptions = this.markerOptions.filter(value => {
        return value.title.toLowerCase().includes(term.toLowerCase());
      }).map(value => {
        return value.title;
      }).slice(0, 15);
    });
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['scene']) {
    }
  }

  onSubmit() {
    this.title = this.searchFormControl.value;
    const input = {
      id: null,
      title: this.title,
      seconds: this.seconds,
      scene_id: this.scene.id,
      primary_tag_id: this.primary_tag_id,
      tag_ids: this.tag_ids
    };

    if (this.editingMarker == null) {
      this.stashService.markerCreate(input).subscribe(response => {
        console.log(response);
        this.hideModal();
      }, error => {
        console.log(error);
      });
    } else {
      input.id = this.editingMarker.id;
      this.stashService.markerUpdate(input).subscribe(response => {
        console.log(response);
        this.hideModal();
      }, error => {
        console.log(error);
      });
    }
  }

  onCancel() {
    this.hideModal();
  }

  onClickDelete() {
    this.deleteClickCount += 1;
    if (this.deleteClickCount > 2) {
      this.stashService.markerDestroy(this.editingMarker.id, this.scene.id).subscribe(response => {
        console.log('Delete successfull:', response);
        this.hideModal();
      });
    }
  }

  onClickAddMarker() {
    this.player.pause();
    this.showModal();
  }

  onClickMarker(marker: SceneMarkerData.Fragment) {
    this.player.seek(marker.seconds);
  }

  onClickEditMarker(marker: SceneMarkerData.Fragment) {
    this.showModal(marker);
  }

  onClickMarkerTitle(title: string) {
    this.setTitle(title);
  }

  setHasFocus(hasFocus: boolean) {
    if (hasFocus === false) {
      setTimeout(() => { this.hasFocus = false; }, 400);
    } else {
      this.hasFocus = hasFocus;
    }
  }

  private hideModal() {
    this.showingMarkerModal = false;
    this.editingMarker = null;
  }

  private showModal(marker: SceneMarkerData.Fragment = null) {
    this.deleteClickCount = 0;
    this.showingMarkerModal = true;

    this.setTitle('');
    this.primary_tag_id = null;
    this.tag_ids = [];
    this.seconds = Math.round(this.player.getPosition());

    if (marker == null) { return; }

    this.editingMarker = marker;

    this.setTitle(marker.title);
    this.seconds = marker.seconds;
    this.primary_tag_id = marker.primary_tag.id;
    this.tag_ids = marker.tags.map(value => value.id);
  }

  private setTitle(title: string) {
    this.title = title;
    this.searchFormControl.setValue(title);
  }
}

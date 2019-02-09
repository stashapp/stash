import { Component, OnInit, OnDestroy, Input, AfterViewInit } from '@angular/core';

import { StashService } from '../../core/stash.service';

import {
  DisplayMode,
  FilterMode,
  ListFilter,
  ListState,
  SceneListState,
  GalleryListState,
  PerformerListState,
  StudioListState,
  SceneMarkerListState
} from '../../shared/models/list-state.model';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ListComponent implements OnInit, OnDestroy, AfterViewInit {
  DisplayMode = DisplayMode;
  FilterMode = FilterMode;

  @Input() state: ListState<any>;

  loading = true;

  constructor(private stashService: StashService,
              private router: Router,
              private activatedRoute: ActivatedRoute) {}

  ngOnInit() {}

  ngOnDestroy() {
    this.loading = true;
    this.state.reset();
    this.state.scrollY = window.scrollY;
  }

  ngAfterViewInit() {
    if (!!this.state.scrollY) {
      setTimeout(() => {
        window.scroll(0, this.state.scrollY);
      }, 1);
    } else {
      window.scroll(0, 0);
    }
  }

  async getData() {
    this.loading = true;

    if (this.state instanceof SceneListState) {
      const result = await this.stashService.findScenes(this.state.filter.currentPage, this.state.filter).result();
      this.state.data = result.data.findScenes.scenes;
      this.state.totalCount = result.data.findScenes.count;
    } else if (this.state instanceof GalleryListState) {
      const result = await this.stashService.findGalleries(this.state.filter.currentPage, this.state.filter).result();
      this.state.data = result.data.findGalleries.galleries;
      this.state.totalCount = result.data.findGalleries.count;
    } else if (this.state instanceof PerformerListState) {
      const result = await this.stashService.findPerformers(this.state.filter.currentPage, this.state.filter).result();
      this.state.data = result.data.findPerformers.performers;
      this.state.totalCount = result.data.findPerformers.count;
    } else if (this.state instanceof StudioListState) {
      const result = await this.stashService.findStudios(this.state.filter.currentPage, this.state.filter).result();
      this.state.data = result.data.findStudios.studios;
      this.state.totalCount = result.data.findStudios.count;
    } else if (this.state instanceof SceneMarkerListState) {
      const result = await this.stashService.findSceneMarkers(this.state.filter.currentPage, this.state.filter).result();
      this.state.data = result.data.findSceneMarkers.scene_markers;
      this.state.totalCount = result.data.findSceneMarkers.count;
    }

    this.loading = false;
  }

  onFilterUpdate(filter: ListFilter) {
    console.log('filter update', filter);
    const options = Object.assign({relativeTo: this.activatedRoute, replaceUrl: true}, filter.makeQueryParameters());
    this.router.navigate([], options);
    this.state.filter = filter;
    this.getData();
  }

  getPage(page: number) {
    this.state.filter.currentPage = page;
    this.onFilterUpdate(this.state.filter);
    window.scroll(0, 0);
  }

}

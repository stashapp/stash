import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';
import { StudiosService } from '../studios.service';

import { StudioData } from '../../core/graphql-generated';

import { SceneListState, CustomCriteria } from '../../shared/models/list-state.model';

@Component({
  selector: 'app-studio-detail',
  templateUrl: './studio-detail.component.html',
  styleUrls: ['./studio-detail.component.scss']
})
export class StudioDetailComponent implements OnInit, OnDestroy {
  studio: StudioData.Fragment;
  sceneListState: SceneListState;

  constructor(
    private route: ActivatedRoute,
    private stashService: StashService,
    private studioService: StudiosService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    this.sceneListState = this.studioService.detailsSceneListState;
    this.sceneListState.filter.customCriteria = [];
    this.sceneListState.filter.customCriteria.push(new CustomCriteria('studio_id', id.toString()));

    this.getStudio();
  }

  ngOnDestroy() {}

  async getStudio() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    const result = await this.stashService.findStudio(id).result();
    this.studio = result.data.findStudio;
  }

  onClickEdit() {
    this.router.navigate(['edit'], { relativeTo: this.route });
  }

}

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';
import { PerformersService } from '../performers.service';

import { PerformerData } from '../../core/graphql-generated';

import { SceneListState, CustomCriteria } from '../../shared/models/list-state.model';

@Component({
  selector: 'app-performer-detail',
  templateUrl: './performer-detail.component.html',
  styleUrls: ['./performer-detail.component.css']
})
export class PerformerDetailComponent implements OnInit {
  performer: PerformerData.Fragment;
  sceneListState: SceneListState;

  constructor(
    private route: ActivatedRoute,
    private stashService: StashService,
    private performerService: PerformersService,
    private router: Router
  ) {}

  ngOnInit() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    this.sceneListState = this.performerService.detailsSceneListState;
    this.sceneListState.filter.customCriteria = [];
    this.sceneListState.filter.customCriteria.push(new CustomCriteria('performer_id', id.toString()));

    this.getPerformer();
    window.scrollTo(0, 0);
  }

  getPerformer() {
    const id = parseInt(this.route.snapshot.params['id'], 10);

    this.stashService.findPerformer(id).valueChanges.subscribe(performer => {
      this.performer = performer.data.findPerformer;
    });
  }

  onClickEdit() {
    this.router.navigate(['edit'], { relativeTo: this.route });
  }

  twitterLink(): string {
    return 'http://www.twitter.com/' + this.performer.twitter;
  }

  instagramLink(): string {
    return 'http://www.instagram.com/' + this.performer.instagram;
  }
}

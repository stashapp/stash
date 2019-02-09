import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';

import { TagData } from '../../core/graphql-generated';

@Component({
  selector: 'app-tag-detail',
  template: `
  <div class="ui text menu">
    <div class="right menu">
      <button (click)="onClickEdit()" class="ui button">Edit</button>
    </div>
  </div>
  <!-- TODO: New tag detail screen... -->
  {{tag?.name}}
  `
})
export class TagDetailComponent implements OnInit {
  tag: TagData.Fragment;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private stashService: StashService
  ) {}

  ngOnInit() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    this.getTag(id);
  }

  async getTag(id: number) {
    if (!!id === false) {
      console.error('no id?');
      return;
    }

    const result = await this.stashService.findTag(id).result();
    this.tag = result.data.findTag;
  }

  onClickEdit() {
    this.router.navigate(['edit'], { relativeTo: this.route });
  }
}

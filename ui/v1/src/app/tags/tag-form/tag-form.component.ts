import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';

@Component({
  selector: 'app-tag-form',
  templateUrl: './tag-form.component.html',
  styleUrls: ['./tag-form.component.scss']
})
export class TagFormComponent implements OnInit {
  name: string;

  loading: Boolean = true;

  constructor(private route: ActivatedRoute, private stashService: StashService, private router: Router) {}

  ngOnInit() {
    this.getTag();
  }

  async getTag() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    if (!!id === false) {
      console.log('new tag');
      this.loading = false;
      return;
    }

    const result = await this.stashService.findTag(id).result();
    this.loading = result.loading;

    if (!result.data.findTag) { this.router.navigate(['/tags']); }

    this.name = result.data.findTag.name;
  }

  hasId() {
    return !!this.route.snapshot.params['id'];
  }

  onSubmit() {
    const id = this.route.snapshot.params['id'];

    if (!!id) {
      this.stashService.tagUpdate({
        id: id,
        name: this.name
      }).subscribe(result => {
        this.router.navigate(['/tags', id]);
      });
    } else {
      this.stashService.tagCreate({
        name: this.name
      }).subscribe(result => {
        this.router.navigate(['/tags', result.data.tagCreate.id]);
      });
    }
  }

  onDestroy() {
    const id = this.route.snapshot.params['id'];

    this.stashService.tagDestroy({
      id: id
    }).subscribe(result => {
      this.router.navigate(['/tags']);
    });
  }

}

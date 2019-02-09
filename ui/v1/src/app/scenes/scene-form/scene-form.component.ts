import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';

import { FindSceneForEditing } from '../../core/graphql-generated';

@Component({
  selector: 'app-scene-form',
  templateUrl: './scene-form.component.html',
  styleUrls: ['./scene-form.component.css']
})
export class SceneFormComponent implements OnInit {
  loading = true;

  title: string;
  details: string;
  url: string;
  date: string;
  rating: number;
  gallery_id: string;
  studio_id: string;
  performer_ids: string[] = [];
  tag_ids: string[] = [];

  performers: FindSceneForEditing.Query['allPerformers'];
  tags: FindSceneForEditing.Query['allTags'];
  studios: FindSceneForEditing.Query['allStudios'];
  galleries: FindSceneForEditing.Query['validGalleriesForScene'];

  constructor(
    private route: ActivatedRoute,
    private stashService: StashService,
    private router: Router
  ) {}

  ngOnInit() {
    this.getScene();
  }

  getScene() {
    const id = parseInt(this.route.snapshot.params['id'], 10);

    if (!!id === false) {
      console.log('new scene');
      return;
    }

    this.stashService.findSceneForEditing(id).valueChanges.subscribe(result => {
      this.title = result.data.findScene.title;
      this.details = result.data.findScene.details;
      this.url   = result.data.findScene.url;
      this.date  = result.data.findScene.date;
      this.rating = result.data.findScene.rating;
      this.gallery_id = !!result.data.findScene.gallery ? result.data.findScene.gallery.id : null;
      this.studio_id = !!result.data.findScene.studio ? result.data.findScene.studio.id : null;
      this.performer_ids = result.data.findScene.performers.map(performer => performer.id);
      this.tag_ids = result.data.findScene.tags.map(tag => tag.id);

      this.performers = result.data.allPerformers;
      this.tags = result.data.allTags;
      this.studios = result.data.allStudios;
      this.galleries = result.data.validGalleriesForScene;

      this.loading = result.loading;
    });
  }

  onSubmit() {
    const id = this.route.snapshot.params['id'];
    this.stashService.sceneUpdate({
      id: id,
      title: this.title,
      details: this.details,
      url: this.url,
      date: this.date,
      rating: this.rating,
      studio_id: this.studio_id,
      gallery_id: this.gallery_id,
      performer_ids: this.performer_ids,
      tag_ids: this.tag_ids
    }).subscribe(result => {
      this.router.navigate(['/scenes', id]);
    });
  }

}

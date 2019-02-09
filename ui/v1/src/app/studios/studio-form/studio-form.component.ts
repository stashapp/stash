import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StashService } from '../../core/stash.service';

@Component({
  selector: 'app-studio-form',
  templateUrl: './studio-form.component.html'
})
export class StudioFormComponent implements OnInit, OnDestroy {
  name: string;
  url: string;
  image: string;

  image_path: string;

  loading: Boolean = true;
  imagePreview: string;

  constructor(private route: ActivatedRoute, private stashService: StashService, private router: Router) { }

  ngOnInit() {
    this.getStudio();
  }

  ngOnDestroy() {
  }

  async getStudio() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    if (!!id === false) {
      console.log('new studio');
      this.loading = false;
      return;
    }

    const result = await this.stashService.findStudio(id).result();
    this.loading = result.loading;
    this.name = result.data.findStudio.name;
    this.url = result.data.findStudio.url;

    this.image_path = result.data.findStudio.image_path;
    this.imagePreview = this.image_path;
  }

  onImageChange(event) {
    const file: File = event.target.files[0];
    const reader: FileReader = new FileReader();

    reader.onloadend = (e) => {
      this.image = reader.result as string;
      this.imagePreview = this.image;
    };
    reader.readAsDataURL(file);
  }

  onResetImage(imageInput) {
    imageInput.value = '';
    this.imagePreview = this.image_path;
    this.image = null;
  }

  // onFavoriteChange() {
  //   this.studio.favorite = !this.studio.favorite;
  // }

  onSubmit() {
    const id = this.route.snapshot.params['id'];

    if (!!id) {
      this.stashService.studioUpdate({
        id: id,
        name: this.name,
        url: this.url,
        image: this.image
      }).subscribe(result => {
        this.router.navigate(['/studios', id]);
      });
    } else {
      this.stashService.studioCreate({
        name: this.name,
        url: this.url,
        image: this.image
      }).subscribe(result => {
        this.router.navigate(['/studios', result.data.studioCreate.id]);
      });
    }
  }
}

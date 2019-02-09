import { Component, OnInit, HostListener } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { StashService } from '../../core/stash.service';

import { GalleryImage } from '../../shared/models/gallery.model';
import { GalleryData } from '../../core/graphql-generated';

@Component({
  selector: 'app-gallery-detail',
  templateUrl: './gallery-detail.component.html',
  styleUrls: ['./gallery-detail.component.css']
})
export class GalleryDetailComponent implements OnInit {
  gallery: GalleryData.Fragment;
  displayedImage: GalleryImage = null;

  constructor(
    private route: ActivatedRoute,
    private stashService: StashService
  ) {}

  ngOnInit() {
    this.getGallery();
    window.scrollTo(0, 0);
  }

  async getGallery() {
    const id = parseInt(this.route.snapshot.params['id'], 10);
    const result = await this.stashService.findGallery(id).result();
    this.gallery = result.data.findGallery;
  }

  @HostListener('mousewheel', ['$event'])
  onMousewheel(event) {
    this.displayedImage = null;
  }

  @HostListener('window:mouseup', ['$event'])
  onMouseup(event: MouseEvent) {
    if (event.button !== 0 || !(event.target instanceof HTMLDivElement)) { return; }
    const target: HTMLDivElement = event.target;
    if (target.id !== 'gallery-image') {
      this.displayedImage = null;
    } else {
      window.open(this.displayedImage.path, '_blank');
    }
  }

  onClickEdit() {
    // TODO
    console.log('edit');
  }

  onClickCard(galleryImage: GalleryImage) {
    console.log(galleryImage);
    this.displayedImage = galleryImage;
  }

  onKey(event) {
    const i = this.displayedImage.index;
    console.log(event);

    switch (event.key) {
      case 'ArrowLeft': {
        this.displayedImage = this.gallery.files[i - 1];
        break;
      }

      case 'ArrowRight': {
        this.displayedImage = this.gallery.files[i + 1];
        break;
      }

      case 'ArrowUp': {
        window.open(this.displayedImage.path, '_blank');
        break;
      }

      case 'ArrowDown': {
        this.displayedImage = null;
        break;
      }

      default:
        break;
    }

    event.preventDefault();
  }

}

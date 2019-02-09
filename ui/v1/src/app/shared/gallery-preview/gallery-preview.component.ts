import { Component, OnInit, OnChanges, Input, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';

import { StashService } from '../../core/stash.service';

import { GalleryImage } from '../../shared/models/gallery.model';
import { GalleryData } from '../../core/graphql-generated';

@Component({
  selector: 'app-gallery-preview',
  templateUrl: './gallery-preview.component.html',
  styleUrls: ['./gallery-preview.component.css']
})
export class GalleryPreviewComponent implements OnInit, OnChanges {
  @Input() gallery: GalleryData.Fragment;
  @Input() galleryId: number;
  @Input() type = 'random';
  @Input() numberOfRandomImages = 12;
  @Input() showTitles = true;
  @Input() numberPerRow = 4;
  @Output() clicked: EventEmitter<GalleryImage> = new EventEmitter();
  files: GalleryImage[];

  constructor(
    private router: Router,
    private stashService: StashService
  ) {}

  ngOnInit() {
    if (!!this.galleryId) {
      this.fetchGallery();
    }
  }

  async fetchGallery() {
    const result = await this.stashService.findGallery(this.galleryId).result();
    this.gallery = result.data.findGallery;
    this.setupFiles();
  }

  imagePath(image) {
    return `${image.path}?thumb=true`;
  }

  shuffle(a) {
    for (let i = a.length; i; i--) {
      const j = Math.floor(Math.random() * i);
      [a[i - 1], a[j]] = [a[j], a[i - 1]];
    }
  }

  onClickGallery() {
    if (this.type === 'random') {
      this.router.navigate(['galleries', this.gallery.id]);
    }
  }

  onClickImage(image) {
    if (this.type === 'full') {
      this.clicked.emit(image);
    }
  }

  suiWidthForNumberPerRow(): string {
    switch (this.numberPerRow) {
      case 1: {
        return 'sixteen';
      }
      case 2: {
        return 'eight';
      }
      case 4: {
        return 'four';
      }
      default:
        return 'four';
    }
  }

  setupFiles() {
    if (!this.gallery) { return; }

    this.files = [...this.gallery.files];
    if (this.type === 'random') {
      this.shuffle(this.files);
      this.files = this.files.slice(0, this.numberOfRandomImages);
    } else if (this.type === 'gallery') {

    }
  }

  ngOnChanges(changes: any) {
    if (!!changes.gallery) {
      this.setupFiles();
    }
  }

}

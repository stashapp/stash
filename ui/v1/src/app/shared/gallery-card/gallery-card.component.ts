import { Component, OnInit, Input, HostBinding } from '@angular/core';
import { Router } from '@angular/router';

import { GalleryData } from '../../core/graphql-generated';

@Component({
  selector: 'app-gallery-card',
  templateUrl: './gallery-card.component.html',
  styleUrls: ['./gallery-card.component.css']
})
export class GalleryCardComponent implements OnInit {
  @Input() gallery: GalleryData.Fragment;

  // The host class needs to be card
  @HostBinding('class') class = 'card';

  constructor(
    private router: Router
  ) {}

  ngOnInit() {
  }

  onSelect(): void {
    this.router.navigate(['/galleries', this.gallery.id]);
  }

}

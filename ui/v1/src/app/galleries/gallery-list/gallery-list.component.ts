import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { GalleriesService } from '../galleries.service';

@Component({
  selector: 'app-gallery-list',
  templateUrl: './gallery-list.component.html'
})
export class GalleryListComponent implements OnInit {
  state = this.galleriesService.listState;

  constructor(private galleriesService: GalleriesService,
              private route: ActivatedRoute,
              private router: Router) {}

  ngOnInit() {}

  onClickNew() {
    this.router.navigate(['new'], { relativeTo: this.route });
  }

}

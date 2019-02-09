import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';

import { GalleriesRoutingModule } from './galleries-routing.module';
import { GalleriesService } from './galleries.service';

import { GalleriesComponent } from './galleries/galleries.component';
import { GalleryDetailComponent } from './gallery-detail/gallery-detail.component';
import { GalleryListComponent } from './gallery-list/gallery-list.component';

@NgModule({
  imports: [
    SharedModule,
    GalleriesRoutingModule
  ],
  declarations: [
    GalleriesComponent,
    GalleryDetailComponent,
    GalleryListComponent
  ],
  providers: [
    GalleriesService
  ]
})
export class GalleriesModule { }

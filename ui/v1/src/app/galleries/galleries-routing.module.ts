import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { GalleriesComponent } from './galleries/galleries.component';
import { GalleryListComponent } from './gallery-list/gallery-list.component';
import { GalleryDetailComponent } from './gallery-detail/gallery-detail.component';

const routes: Routes = [
  { path: '',
    component: GalleriesComponent,
    children: [
      { path: '', component: GalleryListComponent },
      { path: ':id', component: GalleryDetailComponent },
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class GalleriesRoutingModule { }

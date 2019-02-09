import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { TagsComponent } from './tags/tags.component';
import { TagListComponent } from './tag-list/tag-list.component';
import { TagDetailComponent } from './tag-detail/tag-detail.component';
import { TagFormComponent } from './tag-form/tag-form.component';

const routes: Routes = [
  { path: '',
    component: TagsComponent,
    children: [
      { path: '', component: TagListComponent },
      { path: 'new', component: TagFormComponent },
      { path: ':id', component: TagDetailComponent },
      { path: ':id/edit', component: TagFormComponent }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class TagsRoutingModule {}

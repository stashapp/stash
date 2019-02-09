import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { StudiosComponent } from './studios/studios.component';
import { StudioListComponent } from './studio-list/studio-list.component';
import { StudioDetailComponent } from './studio-detail/studio-detail.component';
import { StudioFormComponent } from './studio-form/studio-form.component';

const routes: Routes = [
  { path: '',
    component: StudiosComponent,
    children: [
      { path: '', component: StudioListComponent },
      { path: 'new', component: StudioFormComponent },
      { path: ':id', component: StudioDetailComponent },
      { path: ':id/edit', component: StudioFormComponent }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class StudiosRoutingModule {}

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PerformersComponent } from './performers/performers.component';
import { PerformerListComponent } from './performer-list/performer-list.component';
import { PerformerDetailComponent } from './performer-detail/performer-detail.component';
import { PerformerFormComponent } from './performer-form/performer-form.component';

const performersRoutes: Routes = [
  { path: '',
    component: PerformersComponent,
    children: [
      { path: '', component: PerformerListComponent },
      { path: 'new', component: PerformerFormComponent },
      { path: ':id', component: PerformerDetailComponent },
      { path: ':id/edit', component: PerformerFormComponent }
    ]
  }
];

@NgModule({
  imports: [
    RouterModule.forChild(performersRoutes)
  ],
  exports: [
    RouterModule
  ]
})
export class PerformersRoutingModule {}

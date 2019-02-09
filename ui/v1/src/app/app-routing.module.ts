import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PageNotFoundComponent } from './core/page-not-found/page-not-found.component';
import { DashboardComponent } from './core/dashboard/dashboard.component';

const appRoutes: Routes = [
  { path: '', component: DashboardComponent },
  { path: 'scenes', loadChildren: './scenes/scenes.module#ScenesModule' },
  { path: 'galleries', loadChildren: './galleries/galleries.module#GalleriesModule' },
  { path: 'performers', loadChildren: './performers/performers.module#PerformersModule' },
  { path: 'studios', loadChildren: './studios/studios.module#StudiosModule' },
  { path: 'tags', loadChildren: './tags/tags.module#TagsModule' },
  { path: 'settings', loadChildren: './settings/settings.module#SettingsModule' },
  { path: '**', component: PageNotFoundComponent }
];

@NgModule({
  imports: [
    RouterModule.forRoot(appRoutes)
  ],
  exports: [RouterModule],
  providers: []
})
export class AppRoutingModule {}

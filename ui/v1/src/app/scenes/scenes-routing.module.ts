import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ScenesComponent } from './scenes/scenes.component';
import { SceneListComponent } from './scene-list/scene-list.component';
import { SceneDetailComponent } from './scene-detail/scene-detail.component';
import { SceneFormComponent } from './scene-form/scene-form.component';
import { SceneWallComponent } from './scene-wall/scene-wall.component';
import { MarkerListComponent } from './marker-list/marker-list.component';

const scenesRoutes: Routes = [
  { path: 'wall', component: SceneWallComponent },
  { path: 'markers', component: MarkerListComponent },
  { path: '',
    component: ScenesComponent,
    children: [
      { path: '', component: SceneListComponent },
      { path: ':id', component: SceneDetailComponent },
      { path: ':id/edit', component: SceneFormComponent }
    ]
  }
];

@NgModule({
  imports: [
    RouterModule.forChild(scenesRoutes)
  ],
  exports: [
    RouterModule
  ]
})
export class ScenesRoutingModule {}

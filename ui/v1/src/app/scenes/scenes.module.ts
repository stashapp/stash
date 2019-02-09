import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';

import { ScenesRoutingModule } from './scenes-routing.module';
import { ScenesService } from './scenes.service';

import { ScenesComponent } from './scenes/scenes.component';
import { SceneListComponent } from './scene-list/scene-list.component';
import { SceneDetailComponent } from './scene-detail/scene-detail.component';
import { SceneFormComponent } from './scene-form/scene-form.component';
import { SceneWallComponent } from './scene-wall/scene-wall.component';
import { SceneDetailScrubberComponent } from './scene-detail-scrubber/scene-detail-scrubber.component';
import { SceneDetailMarkerManagerComponent } from './scene-detail-marker-manager/scene-detail-marker-manager.component';
import { MarkerListComponent } from './marker-list/marker-list.component';

@NgModule({
  imports: [
    SharedModule,
    ScenesRoutingModule
  ],
  declarations: [
    ScenesComponent,
    SceneListComponent,
    SceneDetailComponent,
    SceneFormComponent,
    SceneWallComponent,
    SceneDetailScrubberComponent,
    SceneDetailMarkerManagerComponent,
    MarkerListComponent
  ],
  providers: [
    ScenesService
  ]
})
export class ScenesModule {}

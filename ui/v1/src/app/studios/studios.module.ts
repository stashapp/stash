import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';

import { StudiosRoutingModule } from './studios-routing.module';
import { StudiosService } from './studios.service';

import { StudiosComponent } from './studios/studios.component';
import { StudioListComponent } from './studio-list/studio-list.component';
import { StudioDetailComponent } from './studio-detail/studio-detail.component';
import { StudioFormComponent } from './studio-form/studio-form.component';

@NgModule({
  imports: [
    SharedModule,
    StudiosRoutingModule
  ],
  declarations: [
    StudiosComponent,
    StudioListComponent,
    StudioDetailComponent,
    StudioFormComponent
  ],
  providers: [
    StudiosService
  ]
})
export class StudiosModule {}

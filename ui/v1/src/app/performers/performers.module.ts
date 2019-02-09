import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { ReactiveFormsModule } from '@angular/forms';

import { PerformersRoutingModule } from './performers-routing.module';
import { PerformersService } from './performers.service';

import { PerformersComponent } from './performers/performers.component';
import { PerformerListComponent } from './performer-list/performer-list.component';
import { PerformerDetailComponent } from './performer-detail/performer-detail.component';
import { PerformerFormComponent } from './performer-form/performer-form.component';

@NgModule({
  imports: [
    ReactiveFormsModule,
    SharedModule,
    PerformersRoutingModule
  ],
  declarations: [
    PerformersComponent,
    PerformerListComponent,
    PerformerDetailComponent,
    PerformerFormComponent
  ],
  providers: [
    PerformersService
  ]
})
export class PerformersModule {}

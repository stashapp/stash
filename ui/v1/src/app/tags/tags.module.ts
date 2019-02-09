import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';

import { TagsRoutingModule } from './tags-routing.module';
import { TagsComponent } from './tags/tags.component';
import { TagListComponent } from './tag-list/tag-list.component';
import { TagFormComponent } from './tag-form/tag-form.component';
import { TagDetailComponent } from './tag-detail/tag-detail.component';

@NgModule({
  imports: [
    SharedModule,
    TagsRoutingModule
  ],
  declarations: [
    TagsComponent,
    TagListComponent,
    TagFormComponent,
    TagDetailComponent
  ]
})
export class TagsModule {}

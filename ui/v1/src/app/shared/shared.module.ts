import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { SuiModule } from 'ng2-semantic-ui';
import { NgxPaginationModule } from 'ngx-pagination';
import { ClipboardModule } from 'ngx-clipboard';
import { LazyLoadImageModule } from 'ng-lazyload-image';

import { SuiPaginationComponent } from './sui-pagination/sui-pagination.component';
import { JwplayerComponent } from './jwplayer/jwplayer.component';
import { SceneCardComponent } from './scene-card/scene-card.component';
import { PerformerCardComponent } from './performer-card/performer-card.component';
import { ListFilterComponent } from './list-filter/list-filter.component';
import { TruncatePipe } from './truncate.pipe';
import { CapitalizePipe } from './capitalize.pipe';
import { AgePipe } from './age.pipe';
import { StudioCardComponent } from './studio-card/studio-card.component';
import { GalleryPreviewComponent } from './gallery-preview/gallery-preview.component';
import { GalleryCardComponent } from './gallery-card/gallery-card.component';
import { ListComponent } from './list/list.component';
import { SceneListItemComponent } from './scene-list-item/scene-list-item.component';
import { SecondsPipe } from './seconds.pipe';
import { VisibleDirective } from './visible.directive';
import { FileSizePipe } from './file-size.pipe';
import { SceneMarkerWallItemComponent } from './scene-marker-wall-item/scene-marker-wall-item.component';
import { SceneWallItemComponent } from './scene-wall-item/scene-wall-item.component';
import { BaseWallItemComponent } from './base-wall-item/base-wall-item.component';
import { ShufflePipe } from './shuffle.pipe';
import { PerformerListItemComponent } from './performer-list-item/performer-list-item.component';
import { FileNamePipe } from './file-name.pipe';

// Import blah.  Include in dec and exports (https://angular.io/guide/ngmodule#shared-modules)

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    SuiModule,
    NgxPaginationModule,
    ClipboardModule,
    LazyLoadImageModule
  ],
  declarations: [
    TruncatePipe,
    SuiPaginationComponent,
    JwplayerComponent,
    SceneCardComponent,
    PerformerCardComponent,
    ListFilterComponent,
    CapitalizePipe,
    AgePipe,
    StudioCardComponent,
    GalleryPreviewComponent,
    GalleryCardComponent,
    ListComponent,
    SceneListItemComponent,
    SecondsPipe,
    VisibleDirective,
    FileSizePipe,
    SceneMarkerWallItemComponent,
    SceneWallItemComponent,
    BaseWallItemComponent,
    ShufflePipe,
    PerformerListItemComponent,
    FileNamePipe,
  ],
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    SuiModule,
    NgxPaginationModule,
    ClipboardModule,
    LazyLoadImageModule,
    SuiPaginationComponent,
    JwplayerComponent,
    SceneCardComponent,
    PerformerCardComponent,
    ListFilterComponent,
    TruncatePipe,
    CapitalizePipe,
    AgePipe,
    StudioCardComponent,
    GalleryPreviewComponent,
    GalleryCardComponent,
    ListComponent,
    SceneListItemComponent,
    SecondsPipe,
    VisibleDirective,
    FileSizePipe,
    SceneMarkerWallItemComponent,
    SceneWallItemComponent,
    BaseWallItemComponent,
    ShufflePipe,
    FileNamePipe
  ],
  providers: [
    FileNamePipe
  ]
})
export class SharedModule { }

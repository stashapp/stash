import { Injectable } from '@angular/core';

import { GalleryListState } from '../shared/models/list-state.model';

@Injectable()
export class GalleriesService {
  listState: GalleryListState = new GalleryListState();

  constructor() { }

}

import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';

import { StashService } from '../../core/stash.service';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html'
})
export class SettingsComponent implements OnInit, OnDestroy {
  progress: number;
  message: string;
  logs: string[];
  statusObservable: Subscription;
  importClickCount = 0;

  constructor(private stashService: StashService) {}

  ngOnInit() {
    this.statusObservable = this.stashService.metadataUpdate().subscribe(response => {
      const result = JSON.parse(response.data.metadataUpdate);

      this.progress = result.progress;
      this.message = result.message;
      this.logs = result.logs;
    });
  }

  ngOnDestroy() {
    if (!this.statusObservable) { return; }
    this.statusObservable.unsubscribe();
  }

  onClickImport() {
    this.importClickCount += 1;
    if (this.importClickCount > 2) {
      this.stashService.metadataImport().refetch();
      this.importClickCount = 0;
    }
  }

  onClickExport() {
    this.stashService.metadataExport().refetch();
  }

  onClickScan() {
    this.stashService.metadataScan().refetch();
  }

  onClickGenerate() {
    this.stashService.metadataGenerate().refetch();
  }

  onClickClean() {
    this.stashService.metadataClean().refetch();
  }

}

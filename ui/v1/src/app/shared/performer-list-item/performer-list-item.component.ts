import { Component, OnInit, Input, HostBinding, OnChanges, SimpleChanges, SimpleChange } from '@angular/core';
import { PerformerData } from '../../core/graphql-generated';
import { StashService } from '../../core/stash.service';
import { ListFilter, FilterMode, Criteria, CriteriaType, CriteriaValueType, CustomCriteria } from '../models/list-state.model';

@Component({
  selector: 'app-performer-list-item',
  templateUrl: './performer-list-item.component.html',
  styleUrls: ['./performer-list-item.component.scss']
})
export class PerformerListItemComponent implements OnInit, OnChanges {
  @Input() performer: PerformerData.Fragment;
  showingScenes = false;
  scenes: any[];

  // The host class needs to be card
  @HostBinding('class') class = 'dark item';

  constructor(private stashService: StashService) { }

  ngOnInit() {}

  async toggleScenes() {
    this.showingScenes = !this.showingScenes;
    await this.getScenes();
  }

  async ngOnChanges(changes: SimpleChanges) {
    // const newPerformerChange: SimpleChange = changes.performer;
    // this.performer = newPerformerChange.currentValue;
    // await this.getScenes();
  }

  async getScenes() {
    const filter = new ListFilter();
    filter.itemsPerPage = 4;
    filter.sortBy = 'random';
    filter.customCriteria = [new CustomCriteria('performer_id', this.performer.id)];

    const result = await this.stashService.findScenes(1, filter).result();
    this.scenes = result.data.findScenes.scenes;
  }

}

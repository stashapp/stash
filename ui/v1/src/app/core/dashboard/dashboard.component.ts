import { Component, OnInit } from '@angular/core';
import { StashService } from '../stash.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  stats: any;

  constructor(private stashService: StashService) {}

  ngOnInit() {
    this.fetchStats();
  }

  async fetchStats() {
    const result = await this.stashService.stats().result();
    this.stats = result.data.stats;
  }

}

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { PerformersService } from '../performers.service';

@Component({
  selector: 'app-performer-list',
  templateUrl: './performer-list.component.html'
})
export class PerformerListComponent implements OnInit {
  state = this.performersService.performerListState;

  constructor(private performersService: PerformersService,
              private route: ActivatedRoute,
              private router: Router) {}

  ngOnInit() {}

  onClickNew() {
    this.router.navigate(['new'], { relativeTo: this.route });
  }

}

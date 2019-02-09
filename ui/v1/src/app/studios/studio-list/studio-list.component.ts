import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { StudiosService } from '../studios.service';

@Component({
  selector: 'app-studio-list',
  templateUrl: './studio-list.component.html'
})
export class StudioListComponent implements OnInit {
  state = this.studiosService.listState;

  constructor(private studiosService: StudiosService,
              private route: ActivatedRoute,
              private router: Router) {}

  ngOnInit() {}

  onClickNew() {
    this.router.navigate(['new'], { relativeTo: this.route });
  }
}

import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <app-navigation-bar></app-navigation-bar>
    <div class="ui main container">
      <router-outlet></router-outlet>
    </div>
  `
})
export class AppComponent {}

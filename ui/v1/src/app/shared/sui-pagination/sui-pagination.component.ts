import { Component, Input, Output, EventEmitter, ChangeDetectionStrategy, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-sui-pagination',
  templateUrl: './sui-pagination.component.html',
  styleUrls: ['./sui-pagination.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None
})
export class SuiPaginationComponent {
  @Input() id: string;
  @Input() maxSize = 7;
  @Input()
  get directionLinks(): boolean {
      return this._directionLinks;
  }
  set directionLinks(value: boolean) {
      this._directionLinks = !!value && <any>value !== 'false';
  }
  @Input()
  get autoHide(): boolean {
      return this._autoHide;
  }
  set autoHide(value: boolean) {
      this._autoHide = !!value && <any>value !== 'false';
  }
  @Input() previousLabel = 'Previous';
  @Input() nextLabel = 'Next';
  @Input() screenReaderPaginationLabel = 'Pagination';
  @Input() screenReaderPageLabel = 'page';
  @Input() screenReaderCurrentLabel = `You're on page`;
  @Output() pageChange: EventEmitter<number> = new EventEmitter<number>();

  private _directionLinks = true;
  private _autoHide = false;

  constructor() {}

}

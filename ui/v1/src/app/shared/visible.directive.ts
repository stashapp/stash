import { Directive, ElementRef, Output, AfterViewInit, OnDestroy, EventEmitter, Inject } from '@angular/core';
import { DOCUMENT } from '@angular/common';

import { Subscription, fromEvent } from 'rxjs';
import { debounceTime, startWith } from 'rxjs/operators';

@Directive({
  // tslint:disable-next-line:directive-selector
  selector: '[visible]'
})
export class VisibleDirective implements AfterViewInit, OnDestroy {

  @Output()
  visibleEvent: EventEmitter<boolean> = new EventEmitter<boolean>();

  scrollSubscription: Subscription;
  resizeSubscription: Subscription;

  constructor(private element: ElementRef, @Inject(DOCUMENT) private document: any) {}

  ngAfterViewInit() {
    this.subscribe();
  }
  ngOnDestroy() {
    this.unsubscribe();
  }

  subscribe() {
    this.scrollSubscription = fromEvent(window, 'scroll').pipe(
      startWith(null),
      debounceTime(2000)
    ).subscribe(() => {
      this.visibleEvent.emit(this.isInViewport());
    });
    this.resizeSubscription = fromEvent(window, 'resize').pipe(
      startWith(null),
      debounceTime(2000)
    ).subscribe(() => {
      this.visibleEvent.emit(this.isInViewport());
    });
  }
  unsubscribe() {
    if (this.scrollSubscription) { this.scrollSubscription.unsubscribe(); }
    if (this.resizeSubscription) { this.resizeSubscription.unsubscribe(); }
  }

  isInViewport(): boolean {
    const rect = this.element.nativeElement.getBoundingClientRect();
    const html = this.document.documentElement;
    const bufferSpace = 400;
    return rect.top >= -(bufferSpace) &&
           rect.left >= 0 &&
           rect.bottom <= (window.innerHeight + bufferSpace || html.clientHeight + bufferSpace) &&
           rect.right <= (window.innerWidth || html.clientWidth);
  }

}

import { Component, OnInit, ViewChild, ElementRef, HostListener } from '@angular/core';

@Component({
  selector: 'app-base-wall-item',
  template: ''
})
export class BaseWallItemComponent implements OnInit {
  private video: any;
  private hoverTimeout: any = null;
  isHovering = false;

  title = '';
  imagePath = '';
  videoPath = '';

  @ViewChild('videoTag')
  set videoTag(videoTag: ElementRef) {
    if (videoTag === undefined) { return; }
    this.video = videoTag.nativeElement;
    this.video.volume = 0.05;
    this.video.loop = true;
    this.video.oncanplay = () => {
      this.video.play();
    };
  }

  constructor() {}

  ngOnInit() {}

  @HostListener('mouseenter', ['$event'])
  onMouseEnter(e) {
    if (!!this.hoverTimeout) { return; }

    const that = this;
    this.hoverTimeout = setTimeout(function() {
      that.configureTimeout(e);
    }, 1000);
  }

  @HostListener('mouseleave')
  onMouseLeave() {
    if (!!this.hoverTimeout) {
      clearTimeout(this.hoverTimeout);
      this.hoverTimeout = null;
    }
    if (this.video !== undefined) {
      this.video.pause();
      this.video.src = '';
    }
    this.isHovering = false;
  }

  @HostListener('mousemove', ['$event'])
  onMouseMove(event: MouseEvent) {
    if (!!this.hoverTimeout) {
      clearTimeout(this.hoverTimeout);
      this.hoverTimeout = null;
    }
    this.configureTimeout(event);
  }

  transitionEnd(event) {
    if (event.target.classList.contains('double-scale')) {
      event.target.style.zIndex = 2;
    } else {
      event.target.style.zIndex = null;
    }
  }

  private configureTimeout(event: MouseEvent) {
    const that = this;
    this.hoverTimeout = setTimeout(function() {
      if (event.target instanceof HTMLElement) {
        const target: HTMLElement = event.target;
        if (target.className === 'scene-wall-item-text-container' ||
            target.offsetParent.className === 'scene-wall-item-text-container') {
          that.configureTimeout(event);
          return;
        }
      }
      that.isHovering = true;
    }, 1000);
  }
}

import {
  Component,
  OnInit,
  OnChanges,
  SimpleChanges,
  Input,
  Output,
  HostListener,
  ViewChild,
  EventEmitter
} from '@angular/core';

import { HttpClient } from '@angular/common/http';

import { SceneData } from '../../core/graphql-generated';

class SceneSpriteItem {
  start: number;
  end: number;
  x: number;
  y: number;
  w: number;
  h: number;
}

@Component({
  selector: 'app-scene-detail-scrubber',
  templateUrl: './scene-detail-scrubber.component.html',
  styleUrls: ['./scene-detail-scrubber.component.css']
})
export class SceneDetailScrubberComponent implements OnInit, OnChanges {
  @Input() scene: SceneData.Fragment;
  @Output() seek: EventEmitter<number> = new EventEmitter();
  @Output() scrolled: EventEmitter<any> = new EventEmitter();

  slider: HTMLElement;
  @ViewChild('scrubberSlider') sliderTag: any;

  indicator: HTMLElement;
  @ViewChild('positionIndicator') indicatorTag: any;

  spriteItems: SceneSpriteItem[] = [];

  private mouseDown = false;
  private last: MouseEvent;
  private start: MouseEvent;
  private velocity = 0;

  private _position = 0;
  getPostion(): number { return this._position; }
  setPosition(newPostion: number, shouldEmit: boolean = true) {
    if (shouldEmit) { this.scrolled.emit(); }

    const midpointOffset = this.slider.clientWidth / 2;

    const bounds = this.getBounds() * -1;
    if (newPostion > midpointOffset) {
      this._position = midpointOffset;
    } else if (newPostion < bounds - midpointOffset) {
      this._position = bounds - midpointOffset;
    } else {
      this._position = newPostion;
    }

    this.slider.style.transform = `translateX(${this._position}px)`;

    const indicatorPosition = ((newPostion - midpointOffset) / (bounds - (midpointOffset * 2)) * this.slider.clientWidth);
    this.indicator.style.transform = `translateX(${indicatorPosition}px)`;
  }

  @HostListener('window:mouseup', ['$event'])
  onMouseup(event: MouseEvent) {
    if (!this.start) { return; }
    this.mouseDown = false;
    const delta = Math.abs(event.clientX - this.start.clientX);
    if (delta < 1 && event.target instanceof HTMLDivElement) {
      const target: HTMLDivElement = event.target;
      let seekSeconds: number = null;

      const spriteIdString = target.getAttribute('data-sprite-item-id');
      if (spriteIdString != null) {
        const spritePercentage = event.offsetX / target.clientWidth;
        const offset = target.offsetLeft + (target.clientWidth * spritePercentage);
        const percentage = offset / this.slider.scrollWidth;
        seekSeconds = percentage * this.scene.file.duration;
      }

      const markerIdString = target.getAttribute('data-marker-id');
      if (markerIdString != null) {
        const marker = this.scene.scene_markers[Number(markerIdString)];
        seekSeconds = marker.seconds;
      }

      if (!!seekSeconds) { this.seek.emit(seekSeconds); }
    } else if (Math.abs(this.velocity) > 25) {
      const newPosition = this.getPostion() + (this.velocity * 10);
      this.setPosition(newPosition);
      this.velocity = 0;
    }
  }

  @HostListener('mousedown', ['$event'])
  onMousedown(event) {
    event.preventDefault();
    this.mouseDown = true;
    this.last = event;
    this.start = event;
    this.velocity = 0;
  }

  @HostListener('mousemove', ['$event'])
  onMousemove(event: MouseEvent) {
    if (!this.mouseDown) { return; }

    // negative dragging right (past), positive left (future)
    const delta = event.clientX - this.last.clientX;

    const movement = event.movementX;
    this.velocity = movement;

    const newPostion = this.getPostion() + delta;
    this.setPosition(newPostion);
    this.last = event;
  }

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.slider = this.sliderTag.nativeElement;
    this.indicator = this.indicatorTag.nativeElement;

    this.slider.style.transform = `translateX(${this.slider.clientWidth / 2}px)`;
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['scene']) {
      this.fetchSpriteInfo();
    }
  }

  fetchSpriteInfo() {
    if (!this.scene) { return; }

    this.http.get(this.scene.paths.vtt, {responseType: 'text'}).subscribe(res => {
      // TODO: This is gnarly
      const lines = res.split('\n');
      if (lines.shift() !== 'WEBVTT') { return; }
      if (lines.shift() !== '') { return; }
      let item = new SceneSpriteItem();
      this.spriteItems = [];
      while (lines.length) {
        const line = lines.shift();

        if (line.includes('#') && line.includes('=') && line.includes(',')) {
          const size = line.split('#')[1].split('=')[1].split(',');
          item.x = Number(size[0]);
          item.y = Number(size[1]);
          item.w = Number(size[2]);
          item.h = Number(size[3]);

          this.spriteItems.push(item);
          item = new SceneSpriteItem();
        } else if (line.includes(' --> ')) {
          const times = line.split(' --> ');

          const start = times[0].split(':');
          item.start = (+start[0]) * 60 * 60 + (+start[1]) * 60 + (+start[2]);

          const end = times[1].split(':');
          item.end = (+end[0]) * 60 * 60 + (+end[1]) * 60 + (+end[2]);
        }
      }
    }, error => {
      console.log(error);
    });
  }

  getBounds(): number {
    return this.slider.scrollWidth - this.slider.clientWidth;
  }

  getStyleForSprite(i) {
    const sprite = this.spriteItems[i];
    const left = sprite.w * i;
    const path = this.scene.paths.vtt.replace('_thumbs.vtt', '_sprite.jpg'); // TODO: Gnarly
    return {
      'width.px': sprite.w,
      'height.px': sprite.h,
      'margin': '0px auto',
      'background-position': -sprite.x + 'px ' + -sprite.y + 'px',
      'background-image': `url(${path})`,
      'left.px': left
    };
  }

  getTagStyle(tag: HTMLDivElement, i: number) {
    if (!this.slider || this.spriteItems.length === 0 || this.getBounds() === 0) { return {}; }

    const marker = this.scene.scene_markers[i];
    const duration = Number(this.scene.file.duration);
    const percentage = marker.seconds / duration;

    // TODO: this doesn't seem necessary anymore.  Double check.
    // Need to offset from the left margin or the tags are slightly off.
    // const offset = Number(window.getComputedStyle(this.slider.offsetParent).marginLeft.replace('px', ''));
    const offset = 0;

    const left = (this.slider.scrollWidth * percentage) - (tag.clientWidth / 2) + offset;
    return {
      'left.px': left,
      'height.px': 20
    };
  }

  goBack() {
    const newPosition = this.getPostion() + this.slider.clientWidth;
    this.setPosition(newPosition);
  }

  goForward() {
    const newPosition = this.getPostion() - this.slider.clientWidth;
    this.setPosition(newPosition);
  }

  public scrollTo(seconds: number) {
    const duration = Number(this.scene.file.duration);
    const percentage = seconds / duration;
    const position = ((this.slider.scrollWidth * percentage) - (this.slider.clientWidth / 2)) * -1;
    this.setPosition(position, false);
  }

}

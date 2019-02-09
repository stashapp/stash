import { Component, OnInit, Input, HostBinding, HostListener, ViewChild } from '@angular/core';

import { SceneData } from '../../core/graphql-generated';

@Component({
  selector: 'app-scene-card',
  templateUrl: './scene-card.component.html',
  styleUrls: ['./scene-card.component.css']
})
export class SceneCardComponent implements OnInit {
  private isPlaying = false;
  private isHovering = false;
  private video: any;
  previewPath: string = null;
  @Input() scene: SceneData.Fragment;

  // The host class needs to be card
  @HostBinding('class') class = 'dark card';
  @ViewChild('videoTag') videoTag: any;

  constructor() {}

  ngOnInit() {
    this.video = this.videoTag.nativeElement;
    this.video.volume = 0.05;
    this.video.onplaying = () => {
      if (this.isHovering === true) {
        this.isPlaying = true;
      } else {
        this.video.pause();
      }
    };
    this.video.onpause = () => this.isPlaying = false;
  }

  @HostListener('mouseenter')
  onMouseEnter() {
    this.isHovering = true;
    if (!this.previewPath) {
      this.previewPath = this.scene.paths.preview;
    }
    if (this.video.paused && !this.isPlaying) {
      this.video.play();
    }
  }

  @HostListener('mouseleave')
  onMouseLeave() {
    this.isHovering = false;
    if (!this.video.paused && this.isPlaying) {
      this.video.pause();
    }
  }

  hasFavoritePerformer(): boolean {
    return this.scene.performers.filter(performer => performer.favorite === true).length > 0;
  }
}

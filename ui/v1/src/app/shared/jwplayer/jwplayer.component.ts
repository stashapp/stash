import { Component, EventEmitter, Input, Output, ElementRef } from '@angular/core';

declare var jwplayer: any;

@Component({
  selector: 'app-jwplayer',
  templateUrl: './jwplayer.component.html',
  styleUrls: ['./jwplayer.component.css']
})
export class JwplayerComponent {
  @Input() public title: string;
  @Input() public file: string;
  @Input() public image: string;
  @Input() public height: string;
  @Input() public width: string;
  @Output() public bufferChange: EventEmitter<any> = new EventEmitter();
  @Output() public complete: EventEmitter<any> = new EventEmitter();
  @Output() public buffer: EventEmitter<any> = new EventEmitter();
  @Output() public error: EventEmitter<any> = new EventEmitter<any>();
  @Output() public play: EventEmitter<any> = new EventEmitter<any>();
  @Output() public start: EventEmitter<any> = new EventEmitter<any>();
  @Output() public fullscreen: EventEmitter<any> = new EventEmitter<any>();
  @Output() public seeked: EventEmitter<any> = new EventEmitter();
  @Output() public time: EventEmitter<any> = new EventEmitter();

  private _player: any = null;

  constructor(private _elementRef: ElementRef) { }

  public get player(): any {
    this._player = this._player || jwplayer(this._elementRef.nativeElement);
    return this._player;
  }

  public setupPlayer(file: string, image?: string, vtt?: string, chaptersVtt?: string) {
    this.player.remove();
    this.player.setup({
      file: file,
      image: image,
      tracks: [
        {
          file: vtt,
          kind: 'thumbnails'
        },
        {
          file: chaptersVtt,
          kind: 'chapters'
        }
      ],
      primary: 'html5',
      autostart: false,
      playbackRateControls: [0.75, 1, 1.5, 2, 3, 4]
    });
    this.handleEventsFor(this.player);
  }

  public handleEventsFor = (player: any) => {
    player.on('bufferChange', this.onBufferChange);
    player.on('buffer', this.onBuffer);
    player.on('complete', this.onComplete);
    player.on('error', this.onError);
    player.on('fullscreen', this.onFullScreen);
    player.on('play', this.onPlay);
    player.on('start', this.onStart);
    player.on('seeked', this.onSeeked);
    player.on('time', this.onTime);
  }

  public onComplete = (options: {}) => this.complete.emit(options);

  public onError = () => this.error.emit();

  public onBufferChange = (options: {
    duration: number,
    bufferPercent: number,
    position: number,
    metadata?: number
  }) => this.bufferChange.emit(options)

  public onBuffer = (options: {
    oldState: string,
    newState: string,
    reason: string
  }) => this.buffer.emit()

  public onStart = (options: {
    oldState: string,
    newState: string,
    reason: string
  }) => this.buffer.emit()

  public onFullScreen = (options: {
    oldState: string,
    newState: string,
    reason: string
  }) => this.buffer.emit()

  public onPlay = (options: {
  }) => this.play.emit()

  public onSeeked = (options: {
  }) => this.seeked.emit()

  public onTime = (options: {
    duration: number,
    position: number,
    viewable: boolean
  }) => this.time.emit(options)

  onKey(event) {
    const currentPlaybackRate = this._player.getPlaybackRate();
    switch (event.key) {

      case '0': {
        this._player.setPlaybackRate(1);
        console.log(`Playback rate: 1`);
        event.preventDefault();
        break;
      }

      case '2': {
        this._player.setPlaybackRate(currentPlaybackRate + 0.5);
        console.log(`Playback rate: ${currentPlaybackRate + 0.5}`);
        event.preventDefault();
        break;
      }

      case '1': {
        this._player.setPlaybackRate(currentPlaybackRate - 0.5);
        console.log(`Playback rate: ${currentPlaybackRate - 0.5}`);
        event.preventDefault();
        break;
      }

      default:
        break;
    }
  }
}

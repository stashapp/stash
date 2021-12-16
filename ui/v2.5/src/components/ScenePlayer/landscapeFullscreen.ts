import videojs from 'video.js';
import window from 'global/window';

// Default options for the plugin.
const defaults = {
  fullscreen: {
    enterOnRotate: true,
    exitOnRotate: true,
    alwaysInLandscapeMode: true,
    iOS: true
  }
};

const { screen } = window;

/* eslint-disable no-console */
screen.lockOrientationUniversal = (mode) => screen.orientation && screen.orientation.lock(mode).then(() => {}, err => console.log(err)) || screen.mozLockOrientation && screen.mozLockOrientation(mode) || screen.msLockOrientation && screen.msLockOrientation(mode);

const angle = () => {
  // iOS
  if (typeof window.orientation === 'number') {
    return window.orientation;
  }
  // Android
  if (screen && screen.orientation && screen.orientation.angle) {
    return window.orientation;
  }
  videojs.log('angle unknown');
  return 0;
};

// Cross-compatibility for Video.js 5 and 6.
const registerPlugin = videojs.registerPlugin || videojs.plugin;
// const dom = videojs.dom || videojs;

/**
 * Function to invoke when the player is ready.
 *
 * This is a great place for your plugin to initialize itself. When this
 * function is called, the player will have its DOM and child components
 * in place.
 *
 * @function onPlayerReady
 * @param    {Player} player
 *           A Video.js player object.
 *
 * @param    {Object} [options={}]
 *           A plain object containing options for the plugin.
 */
const onPlayerReady = (player, options) => {
  player.addClass('vjs-landscape-fullscreen');

  if (options.fullscreen.iOS &&
    videojs.browser.IS_IOS && videojs.browser.IOS_VERSION > 9 &&
    !player.el_.ownerDocument.querySelector('.bc-iframe')) {
    player.tech_.el_.setAttribute('playsinline', 'playsinline');
    player.tech_.supportsFullScreen = function() {
      return false;
    };
  }

  const rotationHandler = () => {
    const currentAngle = angle();

    console.log(currentAngle);
    if (currentAngle === 90 || currentAngle === 270 || currentAngle === -90) {
      if (options.fullscreen.enterOnRotate && player.paused() === false) {
        player.requestFullscreen();
        // screen.lockOrientationUniversal('landscape');
      }
    }
    if (currentAngle === 0 || currentAngle === 180) {
      if (options.fullscreen.exitOnRotate && player.isFullscreen()) {
        player.exitFullscreen();
      }
    }
  };

  if (videojs.browser.IS_IOS) {
    window.addEventListener('orientationchange', rotationHandler);
  } else if (screen && screen.orientation) {
    // addEventListener('orientationchange') is not a user interaction on Android
    screen.orientation.onchange = rotationHandler;
  }

  player.on('fullscreenchange', e => {
    if (videojs.browser.IS_ANDROID || videojs.browser.IS_IOS) {

      if (!angle() && player.isFullscreen() && options.fullscreen.alwaysInLandscapeMode) {
        screen.lockOrientationUniversal('landscape');
      }
    }
  });
};

/**
 * A video.js plugin.
 *
 * In the plugin function, the value of `this` is a video.js `Player`
 * instance. You cannot rely on the player being in a "ready" state here,
 * depending on how the plugin is invoked. This may or may not be important
 * to you; if not, remove the wait for "ready"!
 *
 * @function landscapeFullscreen
 * @param    {Object} [options={}]
 *           An object of options left to the plugin author to define.
 */
const landscapeFullscreen = function(options) {
  if (videojs.browser.IS_ANDROID || videojs.browser.IS_IOS) {
    this.ready(() => {
      onPlayerReady(this, videojs.mergeOptions(defaults, options));
    });
  }
};

// Register the plugin with video.js.
registerPlugin('landscapeFullscreen', landscapeFullscreen);

export default landscapeFullscreen;

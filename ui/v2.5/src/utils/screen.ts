const isMobile = () =>
  window.matchMedia("only screen and (max-width: 576px)").matches;

const isTouch = () => window.matchMedia("(pointer: coarse)").matches;

const isSmallScreen = () => window.matchMedia("(max-width: 1200px)").matches;

const ScreenUtils = {
  isMobile,
  isTouch,
  isSmallScreen,
};

export default ScreenUtils;

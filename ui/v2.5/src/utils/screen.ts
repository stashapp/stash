const isMobile = () =>
  window.matchMedia("only screen and (max-width: 576px)").matches;

const isTouch = () => window.matchMedia("(pointer: coarse)").matches;

const ScreenUtils = {
  isMobile,
  isTouch,
};

export default ScreenUtils;

const isMobile = () =>
  window.matchMedia("only screen and (max-width: 767px)").matches;

const ScreenUtils = {
  isMobile,
};

export default ScreenUtils;

const isMobile = () =>
  window.matchMedia("only screen and (max-width: 576px)").matches;

const ScreenUtils = {
  isMobile,
};

export default ScreenUtils;

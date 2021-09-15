const isMobile = () =>
  window.matchMedia("only screen and (max-width: 767px)").matches;

export default {
  isMobile,
};

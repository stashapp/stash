function determineSlidesToScroll(
  cardCount: number,
  prefered: number,
  isTouch: boolean
) {
  if (isTouch) {
    return 1;
  } else if (cardCount! > prefered) {
    return prefered;
  } else {
    return cardCount;
  }
}

export function getSlickSliderSettings(cardCount: number, isTouch: boolean) {
  return {
    dots: !isTouch,
    arrows: !isTouch,
    infinite: !isTouch && cardCount > 5,
    speed: 300,
    variableWidth: true,
    swipeToSlide: true,
    slidesToShow: cardCount! > 5 ? 5 : cardCount,
    slidesToScroll: determineSlidesToScroll(cardCount!, 5, isTouch),
    responsive: [
      {
        breakpoint: 1909,
        settings: {
          infinite: !isTouch && cardCount > 4,
          slidesToShow: cardCount! > 4 ? 4 : cardCount,
          slidesToScroll: determineSlidesToScroll(cardCount!, 4, isTouch),
        },
      },
      {
        breakpoint: 1542,
        settings: {
          infinite: !isTouch && cardCount > 3,
          slidesToShow: cardCount! > 3 ? 3 : cardCount,
          slidesToScroll: determineSlidesToScroll(cardCount!, 3, isTouch),
        },
      },
      {
        breakpoint: 1170,
        settings: {
          infinite: !isTouch && cardCount > 2,
          slidesToShow: cardCount! > 2 ? 2 : cardCount,
          slidesToScroll: determineSlidesToScroll(cardCount!, 2, isTouch),
        },
      },
      {
        breakpoint: 801,
        settings: {
          infinite: !isTouch && cardCount > 1,
          slidesToShow: 1,
          slidesToScroll: 1,
          dots: cardCount < 6,
        },
      },
    ],
  };
}

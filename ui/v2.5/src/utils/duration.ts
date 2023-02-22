import TextUtils from "./text";

const secondsToString = (seconds: number) => {
  let ret = TextUtils.secondsToTimestamp(seconds);

  if (ret.startsWith("00:")) {
    ret = ret.substr(3);

    if (ret.startsWith("0")) {
      ret = ret.substr(1);
    }
  }

  return ret;
};

const stringToSeconds = (v?: string) => {
  if (!v) {
    return 0;
  }

  const splits = v.split(":");

  if (splits.length > 3) {
    return 0;
  }

  let seconds = 0;
  let factor = 1;
  while (splits.length > 0) {
    const thisSplit = splits.pop();
    if (thisSplit === undefined) {
      return 0;
    }

    const thisInt = parseInt(thisSplit, 10);
    if (Number.isNaN(thisInt)) {
      return 0;
    }

    seconds += factor * thisInt;
    factor *= 60;
  }

  return seconds;
};

const DurationUtils = {
  secondsToString,
  stringToSeconds,
};

export default DurationUtils;

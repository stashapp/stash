import { TextUtils } from "./text";

export class DurationUtils {
  public static secondsToString(seconds : number) {
    let ret = TextUtils.secondsToTimestamp(seconds);

    if (ret.startsWith("00:")) {
      ret = ret.substr(3);

      if (ret.startsWith("0")) {
        ret = ret.substr(1);
      }
    }

    return ret;
  }

  public static stringToSeconds(v : string) {
    if (!v) {
      return 0;
    }
    
    let splits = v.split(":");

    if (splits.length > 3) {
      return 0;
    }

    let seconds = 0;
    let factor = 1;
    while(splits.length > 0) {
      let thisSplit = splits.pop();
      if (thisSplit === undefined) {
        return 0;
      }

      let thisInt = parseInt(thisSplit, 10);
      if (isNaN(thisInt)) {
        return 0;
      }

      seconds += factor * thisInt;
      factor *= 60;
    }

    return seconds;
  }
}
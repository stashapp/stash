import { IntlShape } from "react-intl";

// Typescript currently does not implement the intl Unit interface
type Unit =
  | "byte"
  | "kilobyte"
  | "megabyte"
  | "gigabyte"
  | "terabyte"
  | "petabyte";
const Units: Unit[] = [
  "byte",
  "kilobyte",
  "megabyte",
  "gigabyte",
  "terabyte",
  "petabyte",
];

const truncate = (
  value?: string,
  limit: number = 100,
  tail: string = "..."
) => {
  if (!value) return "";
  return value.length > limit ? value.substring(0, limit) + tail : value;
};

const fileSize = (bytes: number = 0) => {
  if (Number.isNaN(parseFloat(String(bytes))) || !Number.isFinite(bytes))
    return { size: 0, unit: Units[0] };

  let unit = 0;
  let count = bytes;
  while (count >= 1024) {
    count /= 1024;
    unit++;
  }

  return {
    size: count,
    unit: Units[unit],
  };
};

const secondsToTimestamp = (seconds: number) => {
  let ret = new Date(seconds * 1000).toISOString().substr(11, 8);

  if (ret.startsWith("00")) {
    // strip hours if under one hour
    ret = ret.substr(3);
  }
  if (ret.startsWith("0")) {
    // for duration under a minute, leave one leading zero
    ret = ret.substr(1);
  }
  return ret;
};

const fileNameFromPath = (path: string) => {
  if (!!path === false) return "No File Name";
  return path.replace(/^.*[\\/]/, "");
};

const getAge = (dateString?: string | null, fromDateString?: string) => {
  if (!dateString) return 0;

  const birthdate = new Date(dateString);
  const fromDate = fromDateString ? new Date(fromDateString) : new Date();

  let age = fromDate.getFullYear() - birthdate.getFullYear();
  if (
    birthdate.getMonth() > fromDate.getMonth() ||
    (birthdate.getMonth() >= fromDate.getMonth() &&
      birthdate.getDay() > fromDate.getDay())
  ) {
    age -= 1;
  }

  return age;
};

const bitRate = (bitrate: number) => {
  const megabits = bitrate / 1000000;
  return `${megabits.toFixed(2)} megabits per second`;
};

const resolution = (height: number) => {
  if (height >= 240 && height < 480) {
    return "240p";
  }
  if (height >= 480 && height < 720) {
    return "480p";
  }
  if (height >= 720 && height < 1080) {
    return "720p";
  }
  if (height >= 1080 && height < 2160) {
    return "1080p";
  }
  if (height >= 2160) {
    return "4K";
  }
};

const twitterURL = new URL("https://www.twitter.com");
const instagramURL = new URL("https://www.instagram.com");

const sanitiseURL = (url?: string, siteURL?: URL) => {
  if (!url) {
    return url;
  }

  if (url.startsWith("http://") || url.startsWith("https://")) {
    // just return the entire URL
    return url;
  }

  if (siteURL) {
    // if url starts with the site host, then prepend the protocol
    if (url.startsWith(siteURL.host)) {
      return `${siteURL.protocol}//${url}`;
    }

    // otherwise, construct the url from the protocol, host and passed url
    return `${siteURL.protocol}//${siteURL.host}/${url}`;
  }

  // just prepend the protocol - assume https
  return `https://${url}`;
};

const formatDate = (intl: IntlShape, date?: string) => {
  if (!date) {
    return "";
  }

  return intl.formatDate(date, { format: "long" });
};

const TextUtils = {
  truncate,
  fileSize,
  secondsToTimestamp,
  fileNameFromPath,
  age: getAge,
  bitRate,
  resolution,
  sanitiseURL,
  twitterURL,
  instagramURL,
  formatDate,
};

export default TextUtils;

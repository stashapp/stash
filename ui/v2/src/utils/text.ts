export class TextUtils {

  public static truncate(value?: string, limit: number = 100, tail: string = "..."): string {
    if (!value) { return ""; }
    return value.length > limit ? value.substring(0, limit) + tail : value;
  }

  public static fileSize(bytes: number = 0, precision: number = 2): string {
    if (isNaN(parseFloat(String(bytes))) || !isFinite(bytes)) { return "?"; }

    let unit = 0;
    while ( bytes >= 1024 ) {
      bytes /= 1024;
      unit++;
    }

    return bytes.toFixed(+precision) + " " + this.units[unit];
  }

  public static secondsToTimestamp(seconds: number): string {
    return new Date(seconds * 1000).toISOString().substr(11, 8);
  }

  public static fileNameFromPath(path: string): string {
    if (!!path === false) { return "No File Name"; }
    return path.replace(/^.*[\\\/]/, "");
  }

  public static age(dateString?: string, fromDateString?: string): number {
    if (!dateString) { return 0; }

    const birthdate = new Date(dateString);
    const fromDate = !!fromDateString ? new Date(fromDateString) : new Date();

    let age = fromDate.getFullYear() - birthdate.getFullYear();
    if (birthdate.getMonth() > fromDate.getMonth() ||
        (birthdate.getMonth() >= fromDate.getMonth() && birthdate.getDay() > fromDate.getDay())) {
      age -= 1;
    }

    return age;
  }

  public static bitRate(bitrate: number) {
    const megabits = bitrate / 1000000;
    return `${megabits.toFixed(2)} megabits per second`;
  }

  private static units = [
    "bytes",
    "kB",
    "MB",
    "GB",
    "TB",
    "PB",
  ];
}

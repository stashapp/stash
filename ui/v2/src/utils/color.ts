export class ColorUtils {
  public static classForRating(rating: number): string {
    switch (rating) {
      case 5: return "rating-5";
      case 4: return "rating-4";
      case 3: return "rating-3";
      case 2: return "rating-2";
      case 1: return "rating-1";
      default: return "";
    }
  }
}

export class ParserField {
  public field: string;
  public helperText?: string;

  constructor(field: string, helperText?: string) {
    this.field = field;
    this.helperText = helperText;
  }

  public getFieldPattern() {
    return `{${this.field}}`;
  }

  static Title = new ParserField("title");
  static Ext = new ParserField("ext", "File extension");
  static Rating = new ParserField("rating100");

  static I = new ParserField("i", "Matches any ignored word");
  static D = new ParserField("d", "Matches any delimiter (.-_)");

  static Performer = new ParserField("performer");
  static Studio = new ParserField("studio");
  static Tag = new ParserField("tag");

  // date fields
  static Date = new ParserField("date", "YYYY-MM-DD");
  static YYYY = new ParserField("yyyy", "Year");
  static YY = new ParserField("yy", "Year (20YY)");
  static MM = new ParserField("mm", "Two digit month");
  static MMM = new ParserField("mmm", "Three letter month (eg Jan)");
  static DD = new ParserField("dd", "Two digit date");
  static YYYYMMDD = new ParserField("yyyymmdd");
  static YYMMDD = new ParserField("yymmdd");
  static DDMMYYYY = new ParserField("ddmmyyyy");
  static DDMMYY = new ParserField("ddmmyy");
  static MMDDYYYY = new ParserField("mmddyyyy");
  static MMDDYY = new ParserField("mmddyy");

  static validFields = [
    ParserField.Title,
    ParserField.Ext,
    ParserField.D,
    ParserField.I,
    ParserField.Rating,
    ParserField.Performer,
    ParserField.Studio,
    ParserField.Tag,
    ParserField.Date,
    ParserField.YYYY,
    ParserField.YY,
    ParserField.MM,
    ParserField.MMM,
    ParserField.DD,
    ParserField.YYYYMMDD,
    ParserField.YYMMDD,
    ParserField.DDMMYYYY,
    ParserField.DDMMYY,
    ParserField.MMDDYYYY,
    ParserField.MMDDYY,
  ];

  static fullDateFields = [
    ParserField.YYYYMMDD,
    ParserField.YYMMDD,
    ParserField.DDMMYYYY,
    ParserField.DDMMYY,
    ParserField.MMDDYYYY,
    ParserField.MMDDYY,
  ];
}

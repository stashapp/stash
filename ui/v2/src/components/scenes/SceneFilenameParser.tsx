import {
  Card,
  FormGroup,
  InputGroup,
  Button,
  H4,
  Spinner,
  HTMLTable,
  Checkbox,
  H5,
  MenuItem,
  HTMLSelect,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import { IBaseProps } from "../../models";
import { StashService } from "../../core/StashService";
import * as GQL from "../../core/generated-graphql";
import { SlimSceneDataFragment, Maybe } from "../../core/generated-graphql";
import { TextUtils } from "../../utils/text";
import _ from "lodash";
import { ToastUtils } from "../../utils/toasts";
import { ErrorUtils } from "../../utils/errors";
import { Pagination } from "../list/Pagination";
import { Select, ItemRenderer, ItemPredicate } from "@blueprintjs/select";
import { FilterMultiSelect } from "../select/FilterMultiSelect";
  
interface IProps extends IBaseProps {}

class ParserResult<T> {
  public value: Maybe<T>;
  public originalValue: Maybe<T>;
  public set: boolean = false;

  public setOriginalValue(v : Maybe<T>) {
    this.originalValue = v;
    this.value = v;
  }
}

class ParserField {
  public field : string;
  public fieldRegex: RegExp;
  public regex : string;
  public helperText? : string;

  constructor(field: string, regex?: string, helperText?: string, captured?: boolean) {
    if (regex === undefined) {
      regex = ".*";
    }

    if (captured === undefined) {
      captured = true;
    }

    this.field = field;
    this.helperText = helperText;

    this.fieldRegex = new RegExp("\\{" + this.field + "\\}", "g");

    var regexStr = regex;
    if (captured) {
      regexStr = "(" + regexStr + ")";
    }
    this.regex = regexStr;
  }

  public replaceInPattern(pattern : string) {
    return pattern.replace(this.fieldRegex, this.regex);
  }

  public getFieldPattern() {
    return "{" + this.field + "}";
  }

  static Title = new ParserField("title");
  static Ext = new ParserField("ext", ".*$", "File extension", false);

  static I = new ParserField("i", undefined, "Matches any ignored word", false);
  static D = new ParserField("d", "(?:\\.|-|_)", "Matches any delimiter (.-_)", false);

  static Performer = new ParserField("performer");
  static Studio = new ParserField("studio");
  static Tag = new ParserField("tag");

  // date fields
  static Date = new ParserField("date", "\\d{4}-\\d{2}-\\d{2}", "YYYY-MM-DD");
  static YYYY = new ParserField("yyyy", "\\d{4}", "Year");
  static YY = new ParserField("yy", "\\d{2}", "Year (20YY)");
  static MM = new ParserField("mm", "\\d{2}", "Two digit month");
  static DD = new ParserField("dd", "\\d{2}", "Two digit date");
  static YYYYMMDD = new ParserField("yyyymmdd", "\\d{8}");
  static YYMMDD = new ParserField("yymmdd", "\\d{6}");
  static DDMMYYYY = new ParserField("ddmmyyyy", "\\d{8}");
  static DDMMYY = new ParserField("ddmmyy", "\\d{6}");
  static MMDDYYYY = new ParserField("mmddyyyy", "\\d{8}");
  static MMDDYY = new ParserField("mmddyy", "\\d{6}");

  static validFields = [
    ParserField.Title,
    ParserField.Ext,
    ParserField.D,
    ParserField.I,
    ParserField.Performer,
    ParserField.Studio,
    ParserField.Tag,
    ParserField.Date,
    ParserField.YYYY,
    ParserField.YY,
    ParserField.MM,
    ParserField.DD,
    ParserField.YYYYMMDD,
    ParserField.YYMMDD,
    ParserField.DDMMYYYY,
    ParserField.DDMMYY,
    ParserField.MMDDYYYY,
    ParserField.MMDDYY
  ]

  static fullDateFields = [
    ParserField.YYYYMMDD,
    ParserField.YYMMDD,
    ParserField.DDMMYYYY,
    ParserField.DDMMYY,
    ParserField.MMDDYYYY,
    ParserField.MMDDYY
  ];

  public static getParserField(field: string) {
    return ParserField.validFields.find((f) => {
      return f.field === field;
    });
  }

  public static isValidField(field : string) {
    return !!ParserField.getParserField(field);
  }

  public static isFullDateField(field : ParserField) {
    return ParserField.fullDateFields.includes(field);
  }

  public static replacePatternWithRegex(pattern: string) {
    ParserField.validFields.forEach((field) => {
      pattern = field.replaceInPattern(pattern);
    });
    return pattern;
  }
}

interface IPerformerQueryMap {
  query: string,
  results: GQL.SlimPerformerDataFragment[]
}

class SceneParserResult {
  public id: string;
  public filename: string;
  public title: ParserResult<string> = new ParserResult();
  public date: ParserResult<string> = new ParserResult();

  public yyyy : ParserResult<string> = new ParserResult();
  public mm : ParserResult<string> = new ParserResult();
  public dd : ParserResult<string> = new ParserResult();

  public studioId: ParserResult<string> = new ParserResult();
  public tagIds: ParserResult<string[]> = new ParserResult();
  public performerIds: ParserResult<string[]> = new ParserResult();

  public studio : string = "";
  public performers : string[] = [];
  public tags : string[] = [];

  public scene : SlimSceneDataFragment;

  constructor(scene : SlimSceneDataFragment) {
    this.id = scene.id;
    this.filename = TextUtils.fileNameFromPath(scene.path);
    this.title.setOriginalValue(scene.title);
    this.date.setOriginalValue(scene.date);

    this.scene = scene;
  }

  public static validateDate(dateStr: string) {
    var splits = dateStr.split("-");
    if (splits.length != 3) {
      return false;
    }
    
    var year = parseInt(splits[0]);
    var month = parseInt(splits[1]);
    var d = parseInt(splits[2]);

    var date = new Date();
    date.setMonth(month - 1);
    date.setDate(d);

    // assume year must be between 1900 and 2100
    if (year < 1900 || year > 2100) {
      return false;
    }

    if (month < 1 || month > 12) {
      return false;
    }

    // not checking individual months to ensure date is in the correct range
    if (d < 1 || d > 31) {
      return false;
    }

    return true;
  }

  private setDate(field: ParserField, value: string) {
    var yearIndex = 0;
    var yearLength = field.field.split("y").length - 1;
    var dateIndex = 0;
    var monthIndex = 0;

    switch (field) {
      case ParserField.YYYYMMDD:
      case ParserField.YYMMDD:
        monthIndex = yearLength;
        dateIndex = monthIndex + 2;
        break;
      case ParserField.DDMMYYYY:
      case ParserField.DDMMYY:
        monthIndex = 2;
        yearIndex = monthIndex + 2;
        break;
      case ParserField.MMDDYYYY:
      case ParserField.MMDDYY:
        dateIndex = monthIndex + 2;
        yearIndex = dateIndex + 2;
        break;
    }

    var yearValue = value.substring(yearIndex, yearIndex + yearLength);
    var monthValue = value.substring(monthIndex, monthIndex + 2);
    var dateValue = value.substring(dateIndex, dateIndex + 2);

    var fullDate = yearValue + "-" + monthValue + "-" + dateValue;

    // ensure the date is valid
    // only set if new value is different from the old
    if (SceneParserResult.validateDate(fullDate) && this.date.originalValue !== fullDate) {
      this.date.set = true;
      this.date.value = fullDate
    }
  }

  public setField(field: ParserField, value: any) {
    var parserResult : ParserResult<any> | undefined = undefined;

    if (ParserField.isFullDateField(field)) {
      this.setDate(field, value);
      return;
    }

    switch (field) {
      case ParserField.Title:
        parserResult = this.title;
        break;
      case ParserField.Date:
        parserResult = this.date;
        break;
      case ParserField.YYYY:
        parserResult = this.yyyy;
        break;
      case ParserField.YY:
        parserResult = this.yyyy;
        value = "20" + value;
        break;
      case ParserField.MM:
        parserResult = this.mm;
        break;
      case ParserField.DD:
        parserResult = this.dd;
        break;
    }
    // TODO - other fields

    // only set if different from original value
    if (!!parserResult && parserResult.originalValue !== value) {
      parserResult.set = true;
      parserResult.value = value;
    }
  }

  private static setInput(object: any, key: string, parserResult : ParserResult<any>) {
    if (parserResult.set) {
      object[key] = parserResult.value;
    }
  }

  // returns true if any of its fields have set == true
  public isChanged() {
    return this.title.set || this.date.set;
  }

  public toSceneUpdateInput() {
    var ret = {
      id: this.id,
      title: this.scene.title,
      details: this.scene.details,
      url: this.scene.url,
      date: this.scene.date,
      rating: this.scene.rating,
      gallery_id: this.scene.gallery ? this.scene.gallery.id : undefined,
      studio_id: this.scene.studio ? this.scene.studio.id : undefined,
      performer_ids: this.scene.performers.map((performer) => performer.id),
      tag_ids: this.scene.tags.map((tag) => tag.id)
    };

    SceneParserResult.setInput(ret, "title", this.title);
    SceneParserResult.setInput(ret, "date", this.date);
    // TODO - other fields as added

    return ret;
  }
};

class ParseMapper {
  public fields : string[] = [];
  public regex : string = "";
  public matched : boolean = true;

  constructor(pattern : string, ignoreFields : string[]) {
    // escape control characters
    this.regex = pattern.replace(/([\-\.\(\)\[\]])/g, "\\$1");

    // replace {} with wildcard
    this.regex = this.regex.replace(/\{\}/g, ".*");

    // set ignore fields
    ignoreFields = ignoreFields.map((s) => s.replace(/([\-\.\(\)\[\]])/g, "\\$1").trim());
    var ignoreClause = ignoreFields.map((s) => "(?:" + s + ")").join("|");
    ignoreClause = "(?:" + ignoreClause + ")";

    ParserField.I.regex = ignoreClause;

    // replace all known fields with applicable regexes
    this.regex = ParserField.replacePatternWithRegex(this.regex);
    
    var ignoreField = new ParserField("i", ignoreClause, undefined, false);
    this.regex = ignoreField.replaceInPattern(this.regex);

    // find invalid fields
    var foundInvalid = this.regex.match(/\{[A-Za-z]+\}/g);
    if (foundInvalid) {
      throw new Error("Invalid fields: " + foundInvalid.join(", "));
    }

    var fieldExtractor = new RegExp(/\{([A-Za-z]+)\}/);
    var result = pattern.match(fieldExtractor);

    while(!!result && result.index !== undefined) {
      var field = result[1];

      this.fields.push(field);
      pattern = pattern.substring(result.index + result[0].length);
      result = pattern.match(fieldExtractor);
    } 
  }

  private postParse(scene: SceneParserResult) {
    // set the date if the components are set
    if (scene.yyyy.set && scene.mm.set && scene.dd.set) {
      var fullDate = scene.yyyy.value + "-" + scene.mm.value + "-" + scene.dd.value;
      if (SceneParserResult.validateDate(fullDate)) {
        scene.setField(ParserField.Date, scene.yyyy.value + "-" + scene.mm.value + "-" + scene.dd.value);
      }
    }
  }

  public parse(scene : SceneParserResult) {
    var regex = new RegExp(this.regex, "i");

    var result = scene.filename.match(regex);

    if(!result) {
      return false;
    }

    var mapper = this;

    result.forEach((match, index) => {
      if (index === 0) {
        // skip entire match
        return;
      }

      var field = mapper.fields[index - 1];
      var parserField = ParserField.getParserField(field);
      if (!!parserField) {
        scene.setField(parserField, match);
      }
    });

    this.postParse(scene);

    return true;
  }
}

interface IParserInput {
  pattern: string,
  ignoreWords: string[],
  whitespaceCharacters: string,
  capitalizeTitle: boolean
}

interface IParserRecipe extends IParserInput {
  description: string
}

const builtInRecipes = [
  {
    pattern: "{title}",
    ignoreWords: [],
    whitespaceCharacters: "",
    capitalizeTitle: false,
    description: "Filename"
  },
  {
    pattern: "{title}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: "",
    capitalizeTitle: false,
    description: "Without extension"
  },
  {
    pattern: "{}.{yy}.{mm}.{dd}.{title}.XXX.{}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: ""
  },
  {
    pattern: "{}.{yy}.{mm}.{dd}.{title}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: ""
  },
  {
    pattern: "{title}.XXX.{}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: ""
  },
  {
    pattern: "{}.{yy}.{mm}.{dd}.{title}.{i}.{ext}",
    ignoreWords: ["cz", "fr"],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: "Foreign language"
  }
];

// TODO:
// Add mappings for tags, performers, studio

export const SceneFilenameParser: FunctionComponent<IProps> = (props: IProps) => {
  const [parser, setParser] = useState<ParseMapper | undefined>();
  const [parserResult, setParserResult] = useState<SceneParserResult[]>([]);
  const [parserInput, setParserInput] = useState<IParserInput>(initialParserInput());

  const [allTitleSet, setAllTitleSet] = useState<boolean>(false);
  const [allDateSet, setAllDateSet] = useState<boolean>(false);
  
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);
  const [totalItems, setTotalItems] = useState<number>(0);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const updateScenes = StashService.useScenesUpdate(getScenesUpdateData());

  function initialParserInput() {
    return {
      pattern: "{title}.{ext}",
      ignoreWords: [],
      whitespaceCharacters: "._",
      capitalizeTitle: true
    };
  }

  function getQueryFilter(regex : string, page: number, perPage: number) : GQL.FindFilterType {
    return {
      q: regex,
      page: page,
      per_page: perPage
    };
  }

  async function onFind() {
    setParserResult([]);

    if (!parser) {
      return;
    }
    
    setIsLoading(true);
    
    try {
      const response = await StashService.querySceneByPathRegex(getQueryFilter(parser.regex, page, pageSize));

      let result = response.data.findScenesByPathRegex;
      if (!!result) {
        parseResults(result.scenes);
        setTotalItems(result.count);
      }
    } catch (err) {
      ErrorUtils.handle(err);
    }

    setIsLoading(false);
  }


  useEffect(() => {
    onFind();
  }, [page, parser, parserInput]);

  useEffect(() => {
    setPage(1);
    onFind();
  }, [pageSize])

  function onFindClicked(input : IParserInput) {
    var parser;
    try {
      parser = new ParseMapper(input.pattern, input.ignoreWords);
    } catch(err) {
      ErrorUtils.handle(err);
      return;
    }

    setParser(parser);
    setParserInput(input);
    setPage(1);
    setTotalItems(0);
  }

  function getScenesUpdateData() {
    return parserResult.filter((result) => result.isChanged()).map((result) => result.toSceneUpdateInput());
  }

  async function onApply() {
    setIsLoading(true);

    try {
      await updateScenes();
      ToastUtils.success("Updated scenes");
    } catch (e) {
      ErrorUtils.handle(e);
    }

    setIsLoading(false);
  }

  function parseResults(scenes : GQL.SlimSceneDataFragment[]) {
    if (scenes && parser) {
      var result = scenes.map((scene) => {
        var parserResult = new SceneParserResult(scene);
        if(!parser.parse(parserResult)) {
          return undefined;
        }

        // post-process
        if (parserResult.title && !!parserResult.title.value) {
          if (parserInput.whitespaceCharacters) {
            var wsRegExp = parserInput.whitespaceCharacters.replace(/([\-\.\(\)\[\]])/g, "\\$1");
            wsRegExp = "[" + wsRegExp + "]";
            parserResult.title.value = parserResult.title.value.replace(new RegExp(wsRegExp, "g"), " ");
          }

          if (parserInput.capitalizeTitle) {
            parserResult.title.value = parserResult.title.value.replace(/(?:^| )\w/g, function (chr) {
              return chr.toUpperCase();
            });
          }
        }
        
        return parserResult;
      }).filter((r) => !!r) as SceneParserResult[];

      setParserResult(result);
    }
  }

  useEffect(() => {
    var newAllTitleSet = !parserResult.some((r) => {
      return !r.title.set;
    });
    var newAllDateSet = !parserResult.some((r) => {
      return !r.date.set;
    });

    if (newAllTitleSet != allTitleSet) {
      setAllTitleSet(newAllTitleSet);
    }
    if (newAllDateSet != allDateSet) {
      setAllDateSet(newAllDateSet);
    }
  }, [parserResult]);

  function onSelectAllTitleSet(selected : boolean) {
    var newResult = [...parserResult];

    newResult.forEach((r) => {
      r.title.set = selected;
    });

    setParserResult(newResult);
    setAllTitleSet(selected);
  }

  function onSelectAllDateSet(selected : boolean) {
    var newResult = [...parserResult];

    newResult.forEach((r) => {
      r.date.set = selected;
    });

    setParserResult(newResult);
    setAllDateSet(selected);
  }

  interface IParserInputProps {
    input: IParserInput,
    onFind: (input : IParserInput) => void
  }

  function ParserInput(props : IParserInputProps) {
    const [pattern, setPattern] = useState<string>(props.input.pattern);
    const [ignoreWords, setIgnoreWords] = useState<string>(props.input.ignoreWords.join(" "));
    const [whitespaceCharacters, setWhitespaceCharacters] = useState<string>(props.input.whitespaceCharacters);
    const [capitalizeTitle, setCapitalizeTitle] = useState<boolean>(props.input.capitalizeTitle);

    function onFind() {
      props.onFind({
        pattern: pattern,
        ignoreWords: ignoreWords.split(" "),
        whitespaceCharacters: whitespaceCharacters,
        capitalizeTitle: capitalizeTitle
      });
    }

    const ParserRecipeSelect = Select.ofType<IParserRecipe>();

    const renderParserRecipe: ItemRenderer<IParserRecipe> = (input, { handleClick, modifiers }) => {
      if (!modifiers.matchesPredicate) {
        return null;
      }
      return (
        <MenuItem
            key={input.pattern}
            onClick={handleClick}
            text={input.pattern || "{}"}
            label={input.description}
        />
      );
    };

    const parserRecipePredicate: ItemPredicate<IParserRecipe> = (query, item) => {
      return item.pattern.includes(query);
    };

    function setParserRecipe(recipe: IParserInput) {
      setPattern(recipe.pattern);
      setIgnoreWords(recipe.ignoreWords.join(" "));
      setWhitespaceCharacters(recipe.whitespaceCharacters);
      setCapitalizeTitle(recipe.capitalizeTitle);
    }
  
    const ParserFieldSelect = Select.ofType<ParserField>();

    const renderParserField: ItemRenderer<ParserField> = (field, { handleClick, modifiers }) => {
        if (!modifiers.matchesPredicate) {
          return null;
        }
        return (
          <MenuItem
              key={field.field}
              onClick={handleClick}
              text={field.field || "{}"}
              label={field.helperText}
          />
        );
    };

    const parserFieldPredicate: ItemPredicate<ParserField> = (query, item) => {
      return item.field.includes(query);
    };

    const validFields = [new ParserField("", undefined, "Wildcard")].concat(ParserField.validFields);
    
    function addParserField(field: ParserField) {
      setPattern(pattern + field.getFieldPattern());
    }

    const parserFieldSelect = (
      <ParserFieldSelect
        items={validFields}
        onItemSelect={(item) => addParserField(item)}
        itemRenderer={renderParserField}
        itemPredicate={parserFieldPredicate}
      >
        <Button 
          text="Add field" 
          rightIcon="caret-down" 
        />
      </ParserFieldSelect>
    );

    const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120"];

    return (
      <>
        <FormGroup className="inputs">
          <FormGroup 
            label="Filename pattern:" 
            inline={true}
            helperText="Use '\\' to escape literal {} characters"
          >
            <InputGroup
              onChange={(newValue: any) => setPattern(newValue.target.value)}
              value={pattern}
              rightElement={parserFieldSelect}
            />
          </FormGroup>

          <FormGroup>
            <FormGroup label="Ignored words:" inline={true} helperText="Matches with {i}">
              <InputGroup
                onChange={(newValue: any) => setIgnoreWords(newValue.target.value)}
                value={ignoreWords}
              />
            </FormGroup>
          </FormGroup>
          
          <FormGroup>
            <H5>Title</H5>
            <FormGroup label="Whitespace characters:" 
            inline={true}
            helperText="These characters will be replaced with whitespace in the title">
              <InputGroup
                onChange={(newValue: any) => setWhitespaceCharacters(newValue.target.value)}
                value={whitespaceCharacters}
              />
            </FormGroup>
            <Checkbox
              label="Capitalize title"
              checked={capitalizeTitle}
              onChange={() => setCapitalizeTitle(!capitalizeTitle)}
              inline={true}
            />
          </FormGroup>
          
          {/* TODO - mapping stuff will go here */}

          <FormGroup>
            <ParserRecipeSelect
              items={builtInRecipes}
              onItemSelect={(item) => setParserRecipe(item)}
              itemRenderer={renderParserRecipe}
              itemPredicate={parserRecipePredicate}
            >
              <Button 
                text="Select Parser Recipe" 
                rightIcon="caret-down" 
              />
            </ParserRecipeSelect>
          </FormGroup>

          <FormGroup>
              <Button text="Find" onClick={() => onFind()} />
              <HTMLSelect
                style={{flexBasis: "min-content"}}
                options={PAGE_SIZE_OPTIONS}
                onChange={(event) => setPageSize(parseInt(event.target.value))}
                value={pageSize}
                className="filter-item"
              />
          </FormGroup>
        </FormGroup>
      </>
    );
  }

  interface ISceneParserFieldProps {
    parserResult : ParserResult<any>
    className? : string
    onSetChanged : (set : boolean) => void
    onValueChanged : (value : any) => void
    renderOriginalInputField: (props : ISceneParserFieldProps) => JSX.Element
    renderNewInputField: (props : ISceneParserFieldProps, onChange : (event : any) => void) => JSX.Element
  }

  function SceneParserField(props : ISceneParserFieldProps) {

    function maybeValueChanged(value : any) {
      if (value !== props.parserResult.value) {
        props.onValueChanged(value);
      }
    }

    return (
      <>
        <td>
          <Checkbox
            checked={props.parserResult.set}
            inline={true}
            onChange={() => {props.onSetChanged(!props.parserResult.set)}}
          />
        </td>
        <td>
          <FormGroup>
            {props.renderOriginalInputField(props)}
            {props.renderNewInputField(props, (value) => maybeValueChanged(value))}
          </FormGroup>
        </td>
      </>
    );
  }

  function renderOriginalInputGroup(props : ISceneParserFieldProps) {
    return (
      <InputGroup
        key="originalValue"
        className={props.className}
        small={true}
        disabled={true}
        value={props.parserResult.originalValue || ""}
      />
    );
  }

  interface IInputGroupWrapperProps {
    parserResult: ParserResult<any>
    onChange : (event : any) => void
    className? : string
  }

  function InputGroupWrapper(props : IInputGroupWrapperProps) {
    const [value, setValue] = useState<string>(props.parserResult.value);

    useEffect(() => {
      setValue(props.parserResult.value);
    }, [props.parserResult.value]);

    return (
      <InputGroup
        key="newValue"
        className={props.className}
        small={true}
        onChange={(event : any) => {setValue(event.target.value)}}
        onBlur={() => props.onChange(value)}
        disabled={!props.parserResult.set}
        value={value || ""}
        autoComplete={"new-password" /* required to prevent Chrome autofilling */}
      />
    );
  }
  
  function renderNewInputGroup(props : ISceneParserFieldProps, onChange : (value : any) => void) {
    return (
      <InputGroupWrapper
        className={props.className}
        onChange={(value : any) => {onChange(value)}}
        parserResult={props.parserResult}
      />
    );
  }

  function renderOriginalPerformerSelect(props : ISceneParserFieldProps) {
    return (
      <FilterMultiSelect
        type="performers"
        onUpdate={() => {}}
        disabled={true}
        initialIds={props.parserResult.originalValue}
      />
    );
  }

  function renderNewPerformerSelect(props : ISceneParserFieldProps, onChange : (value : any) => void) {
    return (
      <FilterMultiSelect
        type="performers"
        onUpdate={(items) => {
          const ids = items.map((i) => i.id);
          onChange(ids);
        }}
        initialIds={props.parserResult.value}
      />
    );
  }

  interface ISceneParserRowProps {
    scene : SceneParserResult,
    onChange: (changedScene : SceneParserResult) => void
  }

  function SceneParserRow(props : ISceneParserRowProps) {

    function changeParser(result : ParserResult<any>, set : boolean, value : any) {
      var newParser = _.clone(result);
      newParser.set = set;
      newParser.value = value;
      return newParser;
    }

    function onTitleChanged(set : boolean, value: string | undefined) {
      var newResult = _.clone(props.scene);
      newResult.title = changeParser(newResult.title, set, value);
      props.onChange(newResult);
    }

    function onDateChanged(set : boolean, value: string | undefined) {
      var newResult = _.clone(props.scene);
      newResult.date = changeParser(newResult.date, set, value);
      props.onChange(newResult);
    }

    function onPerformerIdsChanged(set : boolean, value: string[] | undefined) {
      var newResult = _.clone(props.scene);
      newResult.performerIds = changeParser(newResult.date, set, value);
      props.onChange(newResult);
    }

    return (
      <>
      <tr className="scene-parser-row">
        <td style={{textAlign: "left"}}>
          {props.scene.filename}
        </td>
        <SceneParserField 
          key="title"
          className="title" 
          parserResult={props.scene.title}
          onSetChanged={(set) => onTitleChanged(set, props.scene.title.value)}
          onValueChanged={(value) => onTitleChanged(props.scene.title.set, value)}
          renderOriginalInputField={renderOriginalInputGroup}
          renderNewInputField={renderNewInputGroup}
        />
        <SceneParserField 
          key="date"
          parserResult={props.scene.date}
          onSetChanged={(set) => onDateChanged(set, props.scene.date.value)}
          onValueChanged={(value) => onDateChanged(props.scene.date.set, value)}
          renderOriginalInputField={renderOriginalInputGroup}
          renderNewInputField={renderNewInputGroup}
        />
        <SceneParserField 
          key="performers"
          parserResult={props.scene.performerIds}
          onSetChanged={(set) => onPerformerIdsChanged(set, props.scene.performerIds.value)}
          onValueChanged={(value) => onPerformerIdsChanged(props.scene.performerIds.set, value)}
          renderOriginalInputField={renderOriginalPerformerSelect}
          renderNewInputField={renderNewPerformerSelect}
        />
        {/*
        <td>
        </td>
        <td>
        </td>*/}
      </tr>
      </>
    )
  }

  function onChange(scene : SceneParserResult, changedScene : SceneParserResult) {
    var newResult = [...parserResult];

    var index = newResult.indexOf(scene);
    newResult[index] = changedScene;

    setParserResult(newResult);
  }

  function renderTable() {
    if (parserResult.length == 0) { return undefined; }

    return (
      <>
      <form autoComplete="off">
        <div className="grid">
          <HTMLTable condensed={true}>
            <thead>
              <tr className="scene-parser-row">
                <th>Filename</th>
                <td>
                  <Checkbox
                    checked={allTitleSet}
                    inline={true}
                    onChange={() => {onSelectAllTitleSet(!allTitleSet)}}
                  />
                </td>
                <th>Title</th>
                <td>
                <Checkbox
                    checked={allDateSet}
                    inline={true}
                    onChange={() => {onSelectAllDateSet(!allDateSet)}}
                  />
                </td>
                <th>Date</th>
                <th>Performers</th>
                {/* TODO <th>Tags</th>
                
                <th>Studio</th>*/}
              </tr>
            </thead>
            <tbody>
              {parserResult.map((scene) => 
                <SceneParserRow 
                  scene={scene} 
                  key={scene.id}
                  onChange={(changedScene) => onChange(scene, changedScene)}/>
              )}
            </tbody>
          </HTMLTable>
        </div>
        <Pagination
          currentPage={page}
          itemsPerPage={pageSize}
          totalItems={totalItems}
          onChangePage={(page) => setPage(page)}
        />
        <Button intent="primary" text="Apply" onClick={() => onApply()}></Button>
      </form>
    </>
    )
  }

  return (
    <Card id="parser-container">
      <H4>Scene Filename Parser</H4>
      <ParserInput
        input={parserInput}
        onFind={(input) => onFindClicked(input)}
      />

      {isLoading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      {renderTable()}
    </Card>
  );
};
  
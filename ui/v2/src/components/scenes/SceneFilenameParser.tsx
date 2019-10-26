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

  class SceneParserResult {
    public id: string;
    public filename: string;
    public title: ParserResult<string> = new ParserResult();
    public date: ParserResult<string> = new ParserResult();

    public yyyy : ParserResult<string> = new ParserResult();
    public mm : ParserResult<string> = new ParserResult();
    public dd : ParserResult<string> = new ParserResult();

    public studioId: ParserResult<string> = new ParserResult();
    public tags: ParserResult<string[]> = new ParserResult();
    public performerIds: ParserResult<string[]> = new ParserResult();

    public scene : SlimSceneDataFragment;

    constructor(scene : SlimSceneDataFragment) {
      this.id = scene.id;
      this.filename = TextUtils.fileNameFromPath(scene.path);
      this.title.setOriginalValue(scene.title);
      this.date.setOriginalValue(scene.date);

      this.scene = scene;
    }

    private setDate(format: string, value: string) {
      var yearIndex = 0;
      var yearLength = format.split("y").length - 1;
      var dateIndex = 0;
      var monthIndex = 0;

      switch (format) {
        case "yyyymmdd":
        case "yymmdd":
          monthIndex = yearLength;
          dateIndex = monthIndex + 2;
          break;
        case "ddmmyyyy":
        case "ddmmyy":
          monthIndex = 2;
          yearIndex = monthIndex + 2;
          break;
        case "mmddyyyy":
        case "mmddyy":
          dateIndex = monthIndex + 2;
          yearIndex = dateIndex + 2;
          break;
      }

      var yearValue = value.substring(yearIndex, yearIndex + yearLength);
      var monthValue = value.substring(monthIndex, monthIndex + 2);
      var dateValue = value.substring(dateIndex, dateIndex + 2);
      this.date.value = yearValue + "-" + monthValue + "-" + dateValue;
      this.date.set = true;
    }

    public setField(field: string, value: any) {
      var parserResult : ParserResult<any> | undefined = undefined;
      switch (field) {
        case "title":
          parserResult = this.title;
          break;
        case "date":
          parserResult = this.date;
          break;
        case "yyyy":
          parserResult = this.yyyy;
          break;
        case "yy":
          parserResult = this.yyyy;
          value = "20" + value;
          break;
        case "mm":
          parserResult = this.mm;
          break;
        case "dd":
          parserResult = this.dd;
          break;
        case "yyyymmdd":
        case "yymmdd":
        case "ddmmyyyy":
        case "ddmmyy":
        case "mmddyyyy":
        case "mmddyy":
          this.setDate(field, value);
          return;
      }
      // TODO - other fields

      if (!!parserResult) {
        parserResult.value = value;
        parserResult.set = true;
      }
    }
  
    private static setInput(object: any, key: string, parserResult : ParserResult<any>) {
      if (parserResult.set) {
        object[key] = parserResult.value;
      }
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
      
      // TODO - ensure there are no unsupported {fields}

      // replace date fields with applicable regexes
      this.regex = this.regex.replace(/\{yyyy\}/g, "(\\d{4})");
      this.regex = this.regex.replace(/\{(?:dd|mm|yy)\}/g, "(\\d{2})");
      this.regex = this.regex.replace(/\{date\}/g, "(\\d{4}-\\d{2}-\\d{2})");
      this.regex = this.regex.replace(/\{(?:yyyymmdd|ddmmyyyy|mmddyyyy)\}/g, "(\\d{8})");
      this.regex = this.regex.replace(/\{(?:yymmdd|ddmmyy|mmddyy)\}/g, "(\\d{6})");

      // replace {i} with ignored fields
      ignoreFields = ignoreFields.map((s) => s.replace(/([\-\.\(\)\[\]])/g, "\\$1").trim());
      var ignoreClause = ignoreFields.map((s) => "(?:" + s + ")").join("|");
      ignoreClause = "(?:" + ignoreClause + ")";
      this.regex = this.regex.replace(/\{i\}/g, ignoreClause);
      
      // replace remaining fields
      this.regex = this.regex.replace(/\{[A-Za-z]+\}/g, "(.*)");

      var fieldExtractor = new RegExp(/\{([A-Za-z]+)\}/);
      var result = pattern.match(fieldExtractor);
      
      while(!!result && result.index !== undefined) {
        this.fields.push(result[1]);
        pattern = pattern.substring(result.index + result[0].length);
        result = pattern.match(fieldExtractor);
      } 
    }

    private postParse(scene: SceneParserResult) {
      // set the date if the components are set
      if (scene.yyyy.set && scene.mm.set && scene.dd.set) {
        scene.setField("date", scene.yyyy.value + "-" + scene.mm.value + "-" + scene.dd.value);
      }
    }

    public parse(scene : SceneParserResult) {
      var regex = new RegExp(this.regex);

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
        scene.setField(field, match);
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

  // TODO:
  // Add {d} for delimiter characters (._-)
  // Add mappings for tags, performers, studio
  // Add drop-down button to add {fields}

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
      
      const response = await StashService.querySceneByPathRegex(getQueryFilter(parser.regex, page, pageSize));
      
      let result = response.data.findScenesByPathRegex;
      if (!!result) {
        parseResults(result.scenes);
        setTotalItems(result.count);
      }

      setIsLoading(false);
    }

    useEffect(() => {
      onFind();
    }, [page, parser, parserInput]);

    function onFindClicked(input : IParserInput) {
      var parser = new ParseMapper(input.pattern, input.ignoreWords);
      setParser(parser);
      setParserInput(input);
      setPage(1);
      setTotalItems(0);
    }

    function getScenesUpdateData() {
      return parserResult.map((result) => result.toSceneUpdateInput());
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

      return (
        <>
          <FormGroup className="inputs">
            <FormGroup label="Filename pattern:" inline={true}>
              <InputGroup
                onChange={(newValue: any) => setPattern(newValue.target.value)}
                value={pattern}
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
                <Button text="Find" onClick={() => onFind()} />
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
    }
  
    function SceneParserField(props : ISceneParserFieldProps) {

      const [value, setValue] = useState<string>(props.parserResult.value);

      function maybeValueChanged() {
        if (value !== props.parserResult.value) {
          props.onValueChanged(value);
        }
      }

      useEffect(() => {
        setValue(props.parserResult.value);
      }, [props.parserResult.value]);

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
              <InputGroup
                key="originalValue"
                className={props.className}
                small={true}
                disabled={true}
                value={props.parserResult.originalValue || ""}
              />
              <InputGroup
                key="newValue"
                className={props.className}
                small={true}
                onChange={(event : any) => {setValue(event.target.value)}}
                onBlur={() => maybeValueChanged()}
                disabled={!props.parserResult.set}
                value={value || ""}
                autoComplete={"new-password" /* required to prevent Chrome autofilling */}
              />
            </FormGroup>
          </td>
        </>
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
          />
          <SceneParserField 
            key="date"
            parserResult={props.scene.date}
            onSetChanged={(set) => onDateChanged(set, props.scene.date.value)}
            onValueChanged={(value) => onDateChanged(props.scene.date.set, value)}
            />
          {/*<td>
          </td>
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
                  {/* TODO <th>Tags</th>
                  <th>Performers</th>
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
          <Button text="Apply" onClick={() => onApply()}></Button>
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
  
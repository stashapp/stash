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
import React, { FunctionComponent, useEffect, useState, useRef } from "react";
import { IBaseProps } from "../../models";
import { StashService } from "../../core/StashService";
import { SlimSceneDataFragment, Maybe } from "../../core/generated-graphql";
import { TextUtils } from "../../utils/text";
import _ from "lodash";
  
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

    constructor(scene : SlimSceneDataFragment) {
      this.id = scene.id;
      this.filename = TextUtils.fileNameFromPath(scene.path);
      this.title.setOriginalValue(scene.title);
      this.date.setOriginalValue(scene.date);
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
          // TODO - special case - set date
          break;
        case "yymmdd":
          // TODO - special case - set date
          break;
        case "ddmmyyyy":
          // TODO - special case - set date
          break;
        case "ddmmyy":
          // TODO - special case - set date
          break;
        case "mmddyyyy":
          // TODO - special case - set date
          break;
        case "mmddyy":
          // TODO - special case - set date
          break;
      }
      // TODO - other fields

      if (!!parserResult) {
        parserResult.value = value;
        parserResult.set = true;
      }
    }
  };

  class ParseMapper {
    public fields : string[] = [];
    public regex : string = "";
    public matched : boolean = true;

    constructor(pattern : string, ignoreFields : string[]) {
      // escape control characters
      this.regex = pattern.replace(/([\-\.\(\)\[\]])/g, "\\$1");

      // replace date fields with applicable regexes
      this.regex = this.regex.replace(/\{yyyy\}/g, "(\\d{4})");

      // replace {i} with ignored fields
      ignoreFields = ignoreFields.map((s) => s.replace(/([\-\.\(\)\[\]])/g, "\\$1").trim());
      var ignoreClause = ignoreFields.map((s) => "(?:" + s + ")").join("|");
      ignoreClause = "(?:" + ignoreClause + ")";
      this.regex = this.regex.replace(/\{i\}/g, ignoreClause);
      
      // replace remaining fields
      this.regex = this.regex.replace(/\{\w+\}/g, "(.*)");

      var regex = new RegExp(/\{(\w+)\}/);
      var result = pattern.match(regex);
      
      while(!!result && result.index !== undefined) {
        this.fields.push(result[1]);
        pattern = pattern.substring(result.index + result[0].length);
        result = pattern.match(regex);
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

  // Outstanding issues:
  // - cannot modify filenames without losing focus due to re-creation of elements
  // - need to add select all/none for fields

  // TODO:
  // Add {d} for delimiter characters (._-)
  // Add mappings for tags, performers, studio
  // Add implementation to apply stuff
  // Add drop-down button to add {fields}

  export const SceneFilenameParser: FunctionComponent<IProps> = (props: IProps) => {
    const [pattern, setPattern] = useState<string>("{title}.{ext}");
    const [ignoreWords, setIgnoreWords] = useState<string>("");
    const [whitespaceCharacters, setWhitespaceCharacters] = useState<string>("._");
    const [capitaliseTitle, setCapitaliseTitle] = useState<boolean>(true);
    const [sceneResults, setSceneResults] = useState<SlimSceneDataFragment[]>([]);
    const [parserResult, setParserResult] = useState<SceneParserResult[]>([]);

    const [ignoreWordsStage, setIgnoreWordsStage] = useState<string>("");

    // Network state
    const [isLoading, setIsLoading] = useState(false);

    function getQueryPattern() {
      // {title}
      var queryPattern = pattern;

      // replace {..} with wildcard
      queryPattern = queryPattern.replace(/\{.*?\}/g, "*");
      
      return queryPattern;
    }

    async function onFind() {
      setIgnoreWords(ignoreWordsStage);
      setIsLoading(true);
      const result = await StashService.querySceneByPath(getQueryPattern());
      
      if (!!result.data.findScenesByFilename) {
        setSceneResults(result.data.findScenesByFilename.scenes);
      }

      setIsLoading(false);
    }

    useEffect(() => {
      if (sceneResults) {

        var parser = new ParseMapper(pattern, ignoreWords.split(" "));

        var result = sceneResults.map((scene) => {
          var parserResult = new SceneParserResult(scene);
          if(!parser.parse(parserResult)) {
            return undefined;
          }

          // post-process
          if (parserResult.title && !!parserResult.title.value) {
            if (whitespaceCharacters) {
              var wsRegExp = whitespaceCharacters.replace(/([\-\.\(\)\[\]])/g, "\\$1");
              wsRegExp = "[" + wsRegExp + "]";
              parserResult.title.value = parserResult.title.value.replace(new RegExp(wsRegExp, "g"), " ");
            }

            if (capitaliseTitle) {
              parserResult.title.value = parserResult.title.value.replace(/(?:^| )\w/g, function (chr) {
                return chr.toUpperCase();
              });
            }
          }
          
          return parserResult;
        }).filter((r) => !!r);

        setParserResult(result as SceneParserResult[]);
      }
    }, [sceneResults, whitespaceCharacters, ignoreWords]);

    interface ISceneParserFieldProps {
      parserResult : ParserResult<any>
      className? : string
      onSetChanged : (set : boolean) => void
      onValueChanged : (value : any) => void
    }
  
    function SceneParserField(props : ISceneParserFieldProps) {
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
              onChange={(event : any) => {props.onValueChanged(event.target.value)}}
              disabled={true /* TODO - make editable !props.parserResult.set*/}
              value={props.parserResult.value || ""}
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
            onSetChanged={(set) => onDateChanged(set, props.scene.title.value)}
            onValueChanged={(value) => onDateChanged(props.scene.title.set, value)}
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
        <div className="grid">
          <HTMLTable condensed={true}>
            <thead>
              <tr className="scene-parser-row">
                <th>Filename</th>
                <td>
                  <Checkbox
                    checked={true /*allTitleSet*/}
                    inline={true}
                    onChange={undefined /*() => {props.onSetChanged(!props.parserResult.set)}*/}
                  />
                </td>
                <th>Title</th>
                <td>
                <Checkbox
                    checked={true /*allDateSet*/}
                    inline={true}
                    onChange={undefined /*() => {props.onSetChanged(!props.parserResult.set)}*/}
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
        <Button text="Apply"></Button>
      </>
      )
    }

    return (
      <Card id="parser-container">
        <H4>Scene Filename Parser</H4>

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
                onChange={(newValue: any) => setIgnoreWordsStage(newValue.target.value)}
                value={ignoreWordsStage}
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
              checked={capitaliseTitle}
              onChange={() => setCapitaliseTitle(!capitaliseTitle)}
              inline={true}
            />
          </FormGroup>
          
          {/* TODO - mapping stuff will go here */}
          <FormGroup>
              <Button text="Find" onClick={() => onFind()} />
          </FormGroup>
        </FormGroup>

        {isLoading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
        {renderTable()}
      </Card>
    );
  };
  
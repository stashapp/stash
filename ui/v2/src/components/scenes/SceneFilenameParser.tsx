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
  TagInput,
  Tree,
  ITreeNode,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState, useRef } from "react";
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
import { FilterSelect } from "../select/FilterSelect";
  
interface IProps extends IBaseProps {}

class ParserResult<T> {
  public value: Maybe<T>;
  public originalValue: Maybe<T>;
  public set: boolean = false;

  public setOriginalValue(v : Maybe<T>) {
    this.originalValue = v;
    this.value = v;
  }

  public setValue(v : Maybe<T>) {
    if (!!v) {
      this.value = v;
      this.set = !_.isEqual(this.value, this.originalValue);
    }
  }
}

class ParserField {
  public field : string;
  public helperText? : string;

  constructor(field: string, helperText?: string) {
    this.field = field;
    this.helperText = helperText;
  }

  public getFieldPattern() {
    return "{" + this.field + "}";
  }

  static Title = new ParserField("title");
  static Ext = new ParserField("ext", "File extension");

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
}
class SceneParserResult {
  public id: string;
  public filename: string;
  public title: ParserResult<string> = new ParserResult();
  public date: ParserResult<string> = new ParserResult();

  public studio: ParserResult<GQL.SlimSceneDataStudio> = new ParserResult();
  public studioId: ParserResult<string> = new ParserResult();
  public tags: ParserResult<GQL.SlimSceneDataTags[]> = new ParserResult();
  public tagIds: ParserResult<string[]> = new ParserResult();
  public performers: ParserResult<GQL.SlimSceneDataPerformers[]> = new ParserResult();
  public performerIds: ParserResult<string[]> = new ParserResult();

  public scene : SlimSceneDataFragment;

  constructor(result : GQL.ParseSceneFilenamesResults) {
    this.scene = result.scene;

    this.id = this.scene.id;
    this.filename = TextUtils.fileNameFromPath(this.scene.path);
    this.title.setOriginalValue(this.scene.title);
    this.date.setOriginalValue(this.scene.date);
    this.performerIds.setOriginalValue(this.scene.performers.map((p) => p.id));
    this.performers.setOriginalValue(this.scene.performers);
    this.tagIds.setOriginalValue(this.scene.tags.map((t) => t.id));
    this.tags.setOriginalValue(this.scene.tags);
    this.studioId.setOriginalValue(this.scene.studio ? this.scene.studio.id : undefined);
    this.studio.setOriginalValue(this.scene.studio);

    this.title.setValue(result.title);
    this.date.setValue(result.date);
    this.performerIds.setValue(result.performer_ids);
    this.tagIds.setValue(result.tag_ids);
    this.studioId.setValue(result.studio_id);

    if (result.performer_ids) {
      this.performers.setValue(result.performer_ids.map((p) => {
        return {
          id: p,
          name: "",
          favorite: false,
          image_path: ""
        };
      }));
    }

    if (result.tag_ids) {
      this.tags.setValue(result.tag_ids.map((t) => {
        return {
          id: t,
          name: "",
        };
      }));
    }

    if (result.studio_id) {
      this.studio.setValue({
        id: result.studio_id,
        name: "",
        image_path: ""
      });
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

interface IParserInput {
  pattern: string,
  ignoreWords: string[],
  whitespaceCharacters: string,
  capitalizeTitle: boolean,
  page: number,
  pageSize: number,
  findClicked: boolean
}

interface IParserRecipe {
  pattern: string,
  ignoreWords: string[],
  whitespaceCharacters: string,
  capitalizeTitle: boolean,
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
  const [parserResult, setParserResult] = useState<SceneParserResult[]>([]);
  const [parserInput, setParserInput] = useState<IParserInput>(initialParserInput());

  const [allTitleSet, setAllTitleSet] = useState<boolean>(false);
  const [allDateSet, setAllDateSet] = useState<boolean>(false);
  const [allPerformerSet, setAllPerformerSet] = useState<boolean>(false);
  const [allTagSet, setAllTagSet] = useState<boolean>(false);
  const [allStudioSet, setAllStudioSet] = useState<boolean>(false);

  const [showFields, setShowFields] = useState<Map<string, boolean>>(initialShowFieldsState());
  
  const [totalItems, setTotalItems] = useState<number>(0);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const updateScenes = StashService.useScenesUpdate(getScenesUpdateData());

  function initialParserInput() {
    return {
      pattern: "{title}.{ext}",
      ignoreWords: [],
      whitespaceCharacters: "._",
      capitalizeTitle: true,
      page: 1,
      pageSize: 20,
      findClicked: false
    };
  }

  function initialShowFieldsState() {
    return new Map<string, boolean>([
      ["Title", true],
      ["Date", true],
      ["Performers", true],
      ["Tags", true],
      ["Studio", true]
    ]);
  }

  function getParserFilter() {
    return {
      q: parserInput.pattern,
      page: parserInput.page,
      per_page: parserInput.pageSize,
      sort: "path",
      direction: GQL.SortDirectionEnum.Asc,
    };
  }

  function getParserInput() {
    return {
      ignoreWords: parserInput.ignoreWords,
      whitespaceCharacters: parserInput.whitespaceCharacters,
      capitalizeTitle: parserInput.capitalizeTitle
    };
  }

  async function onFind() {
    setParserResult([]);

    setIsLoading(true);
    
    try {
      const response = await StashService.queryParseSceneFilenames(getParserFilter(), getParserInput());

      let result = response.data.parseSceneFilenames;
      if (!!result) {
        parseResults(result.results);
        setTotalItems(result.count);
      }
    } catch (err) {
      ErrorUtils.handle(err);
    }

    setIsLoading(false);
  }

  useEffect(() => {
    if(parserInput.findClicked) {
      onFind();
    }
  }, [parserInput]);

  function onPageSizeChanged(newSize : number) {
    var newInput = _.clone(parserInput);
    newInput.page = 1;
    newInput.pageSize = newSize;
    setParserInput(newInput);
  }

  function onPageChanged(newPage : number) {
    if (newPage !== parserInput.page) {
      var newInput = _.clone(parserInput);
      newInput.page = newPage;
      setParserInput(newInput);
    }
  }

  function onFindClicked(input : IParserInput) {
    input.page = 1;
    input.findClicked = true;
    setParserInput(input);
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

  function parseResults(results : GQL.ParseSceneFilenamesResults[]) {
    if (results) {
      var result = results.map((r) => {
        return new SceneParserResult(r);
      }).filter((r) => !!r) as SceneParserResult[];

      setParserResult(result);
      determineFieldsToHide();
    }
  }

  function determineFieldsToHide() {
    var pattern = parserInput.pattern;
    var titleSet = pattern.includes("{title}");
    var dateSet = pattern.includes("{date}") || 
      pattern.includes("{dd}") || // don't worry about other partial date fields since this should be implied
      ParserField.fullDateFields.some((f) => {
        return pattern.includes("{" + f.field + "}");
      });
    var performerSet = pattern.includes("{performer}");
    var tagSet = pattern.includes("{tag}");
    var studioSet = pattern.includes("{studio}");

    var showFieldsCopy = _.clone(showFields);
    showFieldsCopy.set("Title", titleSet);
    showFieldsCopy.set("Date", dateSet);
    showFieldsCopy.set("Performers", performerSet);
    showFieldsCopy.set("Tags", tagSet);
    showFieldsCopy.set("Studio", studioSet);
    setShowFields(showFieldsCopy);
  }

  useEffect(() => {
    var newAllTitleSet = !parserResult.some((r) => {
      return !r.title.set;
    });
    var newAllDateSet = !parserResult.some((r) => {
      return !r.date.set;
    });
    var newAllPerformerSet = !parserResult.some((r) => {
      return !r.performerIds.set;
    });
    var newAllTagSet = !parserResult.some((r) => {
      return !r.tagIds.set;
    });
    var newAllStudioSet = !parserResult.some((r) => {
      return !r.studioId.set;
    });

    if (newAllTitleSet != allTitleSet) {
      setAllTitleSet(newAllTitleSet);
    }
    if (newAllDateSet != allDateSet) {
      setAllDateSet(newAllDateSet);
    }
    if (newAllPerformerSet != allPerformerSet) {
      setAllTagSet(newAllPerformerSet);
    }
    if (newAllTagSet != allTagSet) {
      setAllTagSet(newAllTagSet);
    }
    if (newAllStudioSet != allStudioSet) {
      setAllStudioSet(newAllStudioSet);
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

  function onSelectAllPerformerSet(selected : boolean) {
    var newResult = [...parserResult];

    newResult.forEach((r) => {
      r.performerIds.set = selected;
    });

    setParserResult(newResult);
    setAllPerformerSet(selected);
  }

  function onSelectAllTagSet(selected : boolean) {
    var newResult = [...parserResult];

    newResult.forEach((r) => {
      r.tagIds.set = selected;
    });

    setParserResult(newResult);
    setAllTagSet(selected);
  }

  function onSelectAllStudioSet(selected : boolean) {
    var newResult = [...parserResult];

    newResult.forEach((r) => {
      r.studioId.set = selected;
    });

    setParserResult(newResult);
    setAllStudioSet(selected);
  }

  interface IShowFieldsTreeProps {
    showFields: Map<string, boolean>
    onShowFieldsChanged: (fields : Map<string, boolean>) => void
  }

  function ShowFieldsTree(props : IShowFieldsTreeProps) {
    const [displayFieldsExpanded, setDisplayFieldsExpanded] = useState<boolean>();

    const treeState: ITreeNode[] = [
      {
        id: 0,
        hasCaret: true,
        label: "Display fields",
        childNodes: [
          {
            id: 1,
            label: "Title",
          },
          {
            id: 2,
            label: "Date",
          },
          {
            id: 3,
            label: "Performers",
          },
          {
            id: 4,
            label: "Tags",
          },
          {
            id: 5,
            label: "Studio",
          }
        ]
      }
    ];

    function setNodeState() {
      if (!!treeState[0].childNodes) {
        treeState[0].childNodes.forEach((n) => {
          n.icon = props.showFields.get(n.label as string) ? "tick" : "cross";
        });
      }

      treeState[0].isExpanded = displayFieldsExpanded;
    }

    setNodeState();

    function expandNode() {
      setDisplayFieldsExpanded(true);
    }

    function collapseNode() {
      setDisplayFieldsExpanded(false);
    }

    function handleClick(nodeData: ITreeNode) {
      var field = nodeData.label as string;
      var fieldsCopy = _.clone(props.showFields);
      fieldsCopy.set(field, !fieldsCopy.get(field));
      props.onShowFieldsChanged(fieldsCopy);
    }

    return (
      <Tree
        contents={treeState}
        onNodeClick={handleClick}
        onNodeCollapse={collapseNode}
        onNodeExpand={expandNode}
      />
    );
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
        capitalizeTitle: capitalizeTitle,
        page: 1,
        pageSize: props.input.pageSize,
        findClicked: props.input.findClicked
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

    function setParserRecipe(recipe: IParserRecipe) {
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

    const validFields = [new ParserField("", "Wildcard")].concat(ParserField.validFields);
    
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
            <ShowFieldsTree
              key="showFields"
              showFields={showFields}
              onShowFieldsChanged={(fields) => setShowFields(fields)}
            />
          </FormGroup>

          <FormGroup>
              <Button text="Find" onClick={() => onFind()} />
              <HTMLSelect
                style={{flexBasis: "min-content"}}
                options={PAGE_SIZE_OPTIONS}
                onChange={(event) => onPageSizeChanged(parseInt(event.target.value))}
                value={props.input.pageSize}
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
    fieldName : string
    onSetChanged : (set : boolean) => void
    onValueChanged : (value : any) => void
    originalParserResult? : ParserResult<any>
    renderOriginalInputField: (props : ISceneParserFieldProps) => JSX.Element
    renderNewInputField: (props : ISceneParserFieldProps, onChange : (event : any) => void) => JSX.Element
  }

  function SceneParserField(props : ISceneParserFieldProps) {

    function maybeValueChanged(value : any) {
      if (value !== props.parserResult.value) {
        props.onValueChanged(value);
      }
    }

    if (!showFields.get(props.fieldName)) {
      return null;
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
    var parserResult = props.parserResult;

    if (!!props.originalParserResult) {
      parserResult = props.originalParserResult;
    }

    return (
      <InputGroup
        key="originalValue"
        className={props.className}
        small={true}
        disabled={true}
        value={parserResult.originalValue || ""}
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

  interface HasName {
    name: string
  }

  function renderOriginalSelect(props : ISceneParserFieldProps) {
    var parserResult = props.parserResult;

    if (!!props.originalParserResult) {
      parserResult = props.originalParserResult;
    }

    var elements = [];
    
    if (parserResult.originalValue) {
      if (parserResult.originalValue.map) {
        elements = parserResult.originalValue.map((element : HasName) => (
          element.name
        ));
      } else {
        elements = [parserResult.originalValue.name];
      }
    }

    return (
      <>
      <TagInput
        className={props.className}
        values={elements}
        disabled={true}
      />
      </>
    )
  }

  function renderNewMultiSelect(type: "performers" | "tags", props : ISceneParserFieldProps, onChange : (value : any) => void) {
    return (
      <FilterMultiSelect
        className={props.className}
        type={type}
        onUpdate={(items) => {
          const ids = items.map((i) => i.id);
          onChange(ids);
        }}
        initialIds={props.parserResult.value}
      />
    );
  }

  function renderNewPerformerSelect(props : ISceneParserFieldProps, onChange : (value : any) => void) {
    return renderNewMultiSelect("performers", props, onChange);
  }

  function renderNewTagSelect(props : ISceneParserFieldProps, onChange : (value : any) => void) {
    return renderNewMultiSelect("tags", props, onChange);
  }

  function renderNewStudioSelect(props : ISceneParserFieldProps, onChange : (value : any) => void) {
    return (
      <FilterSelect
        type="studios"
        noSelectionString=""
        className={props.className}
        onSelectItem={(item) => onChange(item ? item.id : undefined)}
        initialId={props.parserResult.value}
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
      newResult.performerIds = changeParser(newResult.performerIds, set, value);
      props.onChange(newResult);
    }

    function onTagIdsChanged(set : boolean, value: string[] | undefined) {
      var newResult = _.clone(props.scene);
      newResult.tagIds = changeParser(newResult.tagIds, set, value);
      props.onChange(newResult);
    }

    function onStudioIdChanged(set : boolean, value: string | undefined) {
      var newResult = _.clone(props.scene);
      newResult.studioId = changeParser(newResult.studioId, set, value);
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
          fieldName="Title"
          className="parser-field-title" 
          parserResult={props.scene.title}
          onSetChanged={(set) => onTitleChanged(set, props.scene.title.value)}
          onValueChanged={(value) => onTitleChanged(props.scene.title.set, value)}
          renderOriginalInputField={renderOriginalInputGroup}
          renderNewInputField={renderNewInputGroup}
        />
        <SceneParserField 
          key="date"
          fieldName="Date"
          className="parser-field-date"
          parserResult={props.scene.date}
          onSetChanged={(set) => onDateChanged(set, props.scene.date.value)}
          onValueChanged={(value) => onDateChanged(props.scene.date.set, value)}
          renderOriginalInputField={renderOriginalInputGroup}
          renderNewInputField={renderNewInputGroup}
        />
        <SceneParserField 
          key="performers"
          fieldName="Performers"
          className="parser-field-performers"
          parserResult={props.scene.performerIds}
          originalParserResult={props.scene.performers}
          onSetChanged={(set) => onPerformerIdsChanged(set, props.scene.performerIds.value)}
          onValueChanged={(value) => onPerformerIdsChanged(props.scene.performerIds.set, value)}
          renderOriginalInputField={renderOriginalSelect}
          renderNewInputField={renderNewPerformerSelect}
        />
        <SceneParserField 
          key="tags"
          fieldName="Tags"
          className="parser-field-tags"
          parserResult={props.scene.tagIds}
          originalParserResult={props.scene.tags}
          onSetChanged={(set) => onTagIdsChanged(set, props.scene.tagIds.value)}
          onValueChanged={(value) => onTagIdsChanged(props.scene.tagIds.set, value)}
          renderOriginalInputField={renderOriginalSelect}
          renderNewInputField={renderNewTagSelect}
        />
        <SceneParserField 
          key="studio"
          fieldName="Studio"
          className="parser-field-studio"
          parserResult={props.scene.studioId}
          originalParserResult={props.scene.studio}
          onSetChanged={(set) => onStudioIdChanged(set, props.scene.studioId.value)}
          onValueChanged={(value) => onStudioIdChanged(props.scene.studioId.set, value)}
          renderOriginalInputField={renderOriginalSelect}
          renderNewInputField={renderNewStudioSelect}
        />
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

  function renderHeader(fieldName: string, allSet: boolean, onAllSet: (set: boolean) => void) {
    if (!showFields.get(fieldName)) {
      return null;
    }

    return (
      <>
      <td>
        <Checkbox
          checked={allSet}
          inline={true}
          onChange={() => {onAllSet(!allSet)}}
        />
      </td>
      <th>{fieldName}</th>
      </>
    )
  }

  function renderTable() {
    if (parserResult.length == 0) { return undefined; }

    return (
      <>
      <div>
        <div className="scene-parser-results">
          <HTMLTable condensed={true}>
            <thead>
              <tr className="scene-parser-row">
                <th>Filename</th>
                {renderHeader("Title", allTitleSet, onSelectAllTitleSet)}
                {renderHeader("Date", allDateSet, onSelectAllDateSet)}
                {renderHeader("Performers", allPerformerSet, onSelectAllPerformerSet)}
                {renderHeader("Tags", allTagSet, onSelectAllTagSet)}
                {renderHeader("Studio", allStudioSet, onSelectAllStudioSet)}
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
          currentPage={parserInput.page}
          itemsPerPage={parserInput.pageSize}
          totalItems={totalItems}
          onChangePage={(page) => onPageChanged(page)}
        />
        <Button intent="primary" text="Apply" onClick={() => onApply()}></Button>
      </div>
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
  
import React, { useState } from "react";
import {
  Button,
  Dropdown,
  DropdownButton,
  Form,
  InputGroup,
} from "react-bootstrap";
import { useIntl } from "react-intl";
import { ParserField } from "./ParserField";
import { ShowFields } from "./ShowFields";

const builtInRecipes = [
  {
    pattern: "{title}",
    ignoreWords: [],
    whitespaceCharacters: "",
    capitalizeTitle: false,
    description: "Filename",
  },
  {
    pattern: "{title}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: "",
    capitalizeTitle: false,
    description: "Without extension",
  },
  {
    pattern: "{}.{yy}.{mm}.{dd}.{title}.XXX.{}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: "",
  },
  {
    pattern: "{}.{yy}.{mm}.{dd}.{title}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: "",
  },
  {
    pattern: "{title}.XXX.{}.{ext}",
    ignoreWords: [],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: "",
  },
  {
    pattern: "{}.{yy}.{mm}.{dd}.{title}.{i}.{ext}",
    ignoreWords: ["cz", "fr"],
    whitespaceCharacters: ".",
    capitalizeTitle: true,
    description: "Foreign language",
  },
];

export interface IParserInput {
  pattern: string;
  ignoreWords: string[];
  whitespaceCharacters: string;
  capitalizeTitle: boolean;
  page: number;
  pageSize: number;
  findClicked: boolean;
  ignoreOrganized: boolean;
}

interface IParserRecipe {
  pattern: string;
  ignoreWords: string[];
  whitespaceCharacters: string;
  capitalizeTitle: boolean;
  description: string;
}

interface IParserInputProps {
  input: IParserInput;
  onFind: (input: IParserInput) => void;
  onPageSizeChanged: (newSize: number) => void;
  showFields: Map<string, boolean>;
  setShowFields: (fields: Map<string, boolean>) => void;
}

export const ParserInput: React.FC<IParserInputProps> = (
  props: IParserInputProps
) => {
  const intl = useIntl();
  const [pattern, setPattern] = useState<string>(props.input.pattern);
  const [ignoreWords, setIgnoreWords] = useState<string>(
    props.input.ignoreWords.join(" ")
  );
  const [whitespaceCharacters, setWhitespaceCharacters] = useState<string>(
    props.input.whitespaceCharacters
  );
  const [capitalizeTitle, setCapitalizeTitle] = useState<boolean>(
    props.input.capitalizeTitle
  );
  const [ignoreOrganized, setIgnoreOrganized] = useState<boolean>(
    props.input.ignoreOrganized
  );

  function onFind() {
    props.onFind({
      pattern,
      ignoreWords: ignoreWords.split(" "),
      whitespaceCharacters,
      capitalizeTitle,
      page: 1,
      pageSize: props.input.pageSize,
      findClicked: props.input.findClicked,
      ignoreOrganized,
    });
  }

  function setParserRecipe(recipe: IParserRecipe) {
    setPattern(recipe.pattern);
    setIgnoreWords(recipe.ignoreWords.join(" "));
    setWhitespaceCharacters(recipe.whitespaceCharacters);
    setCapitalizeTitle(recipe.capitalizeTitle);
  }

  const validFields = [new ParserField("", "Wildcard")].concat(
    ParserField.validFields
  );

  function addParserField(field: ParserField) {
    setPattern(pattern + field.getFieldPattern());
  }

  const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120", "250", "500", "1000"];

  return (
    <Form.Group>
      <Form.Group className="row">
        <Form.Label htmlFor="filename-pattern" className="col-2">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.filename_pattern",
          })}
        </Form.Label>
        <InputGroup className="col-8">
          <Form.Control
            className="text-input"
            id="filename-pattern"
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPattern(e.currentTarget.value)
            }
            value={pattern}
          />
          <InputGroup.Append>
            <DropdownButton
              id="parser-field-select"
              title={intl.formatMessage({
                id: "config.tools.scene_filename_parser.add_field",
              })}
            >
              {validFields.map((item) => (
                <Dropdown.Item
                  key={item.field}
                  onSelect={() => addParserField(item)}
                >
                  <span className="mr-2">{item.field || "{}"}</span>
                  <span className="ml-auto text-muted">{item.helperText}</span>
                </Dropdown.Item>
              ))}
            </DropdownButton>
          </InputGroup.Append>
        </InputGroup>
        <Form.Text className="text-muted row col-10 offset-2">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.escape_chars",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group className="row" controlId="ignored-words">
        <Form.Label className="col-2">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.ignored_words",
          })}
        </Form.Label>
        <InputGroup className="col-8">
          <Form.Control
            className="text-input"
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setIgnoreWords(e.currentTarget.value)
            }
            value={ignoreWords}
          />
        </InputGroup>
        <Form.Text className="text-muted col-10 offset-2">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.matches_with",
          })}
        </Form.Text>
      </Form.Group>

      <h5>{intl.formatMessage({ id: "title" })}</h5>
      <Form.Group className="row">
        <Form.Label htmlFor="whitespace-characters" className="col-2">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.whitespace_chars",
          })}
        </Form.Label>
        <InputGroup className="col-8">
          <Form.Control
            className="text-input"
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setWhitespaceCharacters(e.currentTarget.value)
            }
            value={whitespaceCharacters}
          />
        </InputGroup>
        <Form.Text className="text-muted col-10 offset-2">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.whitespace_chars_desc",
          })}
        </Form.Text>
      </Form.Group>
      <Form.Group>
        <Form.Check
          inline
          className="m-0"
          id="capitalize-title"
          checked={capitalizeTitle}
          onChange={() => setCapitalizeTitle(!capitalizeTitle)}
        />
        <Form.Label htmlFor="capitalize-title">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.capitalize_title",
          })}
        </Form.Label>
      </Form.Group>
      <Form.Group>
        <Form.Check
          inline
          className="m-0"
          id="ignore-organized"
          checked={ignoreOrganized}
          onChange={() => setIgnoreOrganized(!ignoreOrganized)}
        />
        <Form.Label htmlFor="ignore-organized">
          {intl.formatMessage({
            id: "config.tools.scene_filename_parser.ignore_organized",
          })}
        </Form.Label>
      </Form.Group>

      {/* TODO - mapping stuff will go here */}

      <Form.Group>
        <DropdownButton
          variant="secondary"
          id="recipe-select"
          title={intl.formatMessage({
            id: "config.tools.scene_filename_parser.select_parser_recipe",
          })}
          drop="up"
        >
          {builtInRecipes.map((item) => (
            <Dropdown.Item
              key={item.pattern}
              onSelect={() => setParserRecipe(item)}
            >
              <span>{item.pattern}</span>
              <span className="ml-auto text-muted">{item.description}</span>
            </Dropdown.Item>
          ))}
        </DropdownButton>
      </Form.Group>

      <Form.Group>
        <ShowFields
          fields={props.showFields}
          onShowFieldsChanged={(fields) => props.setShowFields(fields)}
        />
      </Form.Group>

      <Form.Group className="row">
        <Button variant="secondary" className="ml-3 col-1" onClick={onFind}>
          {intl.formatMessage({ id: "actions.find" })}
        </Button>
        <Form.Control
          as="select"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            props.onPageSizeChanged(parseInt(e.currentTarget.value, 10))
          }
          defaultValue={props.input.pageSize}
          className="col-1 input-control filter-item"
        >
          {PAGE_SIZE_OPTIONS.map((val) => (
            <option key={val} value={val}>
              {val}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    </Form.Group>
  );
};

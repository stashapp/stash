import React, { useState } from "react";
import {
  Button,
  Dropdown,
  DropdownButton,
  Form,
  InputGroup
} from "react-bootstrap";
import { ParserField } from "./ParserField";
import { ShowFields } from "./ShowFields";

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

export interface IParserInput {
  pattern: string;
  ignoreWords: string[];
  whitespaceCharacters: string;
  capitalizeTitle: boolean;
  page: number;
  pageSize: number;
  findClicked: boolean;
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

  function onFind() {
    props.onFind({
      pattern,
      ignoreWords: ignoreWords.split(" "),
      whitespaceCharacters,
      capitalizeTitle,
      page: 1,
      pageSize: props.input.pageSize,
      findClicked: props.input.findClicked
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

  const PAGE_SIZE_OPTIONS = ["20", "40", "60", "120"];

  return (
    <Form.Group>
      <Form.Group className="row">
        <Form.Label htmlFor="filename-pattern" className="col-2">
          Filename Pattern
        </Form.Label>
        <InputGroup className="col-8">
          <Form.Control
            id="filename-pattern"
            onChange={(e: React.FormEvent<HTMLInputElement>) =>
              setPattern(e.currentTarget.value)
            }
            value={pattern}
          />
          <InputGroup.Append>
            <DropdownButton id="parser-field-select" title="Add Field">
              {validFields.map(item => (
                <Dropdown.Item
                  key={item.field}
                  onSelect={() => addParserField(item)}
                >
                  <span>{item.field}</span>
                  <span className="ml-auto">{item.helperText}</span>
                </Dropdown.Item>
              ))}
            </DropdownButton>
          </InputGroup.Append>
        </InputGroup>
        <Form.Text className="text-muted row col-10 offset-2">
          Use &apos;\&apos; to escape literal {} characters
        </Form.Text>
      </Form.Group>

      <Form.Group className="row" controlId="ignored-words">
        <Form.Label className="col-2">Ignored words</Form.Label>
        <InputGroup className="col-8">
          <Form.Control
            onChange={(e: React.FormEvent<HTMLInputElement>) =>
              setIgnoreWords(e.currentTarget.value)
            }
            value={ignoreWords}
          />
        </InputGroup>
        <Form.Text className="text-muted col-10 offset-2">
          Matches with {"{i}"}
        </Form.Text>
      </Form.Group>

      <h5>Title</h5>
      <Form.Group className="row">
        <Form.Label htmlFor="whitespace-characters" className="col-2">
          Whitespace characters:
        </Form.Label>
        <InputGroup className="col-8">
          <Form.Control
            onChange={(e: React.FormEvent<HTMLInputElement>) =>
              setWhitespaceCharacters(e.currentTarget.value)
            }
            value={whitespaceCharacters}
          />
        </InputGroup>
        <Form.Text className="text-muted col-10 offset-2">
          These characters will be replaced with whitespace in the title
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
        <Form.Label htmlFor="capitalize-title">Capitalize title</Form.Label>
      </Form.Group>

      {/* TODO - mapping stuff will go here */}

      <Form.Group>
        <DropdownButton
          variant="secondary"
          id="recipe-select"
          title="Select Parser Recipe"
        >
          {builtInRecipes.map(item => (
            <Dropdown.Item
              key={item.pattern}
              onSelect={() => setParserRecipe(item)}
            >
              <span>{item.pattern}</span>
              <span className="mr-auto">{item.description}</span>
            </Dropdown.Item>
          ))}
        </DropdownButton>
      </Form.Group>

      <Form.Group>
        <ShowFields
          fields={props.showFields}
          onShowFieldsChanged={fields => props.setShowFields(fields)}
        />
      </Form.Group>

      <Form.Group className="row">
        <Button variant="secondary" className="ml-3 col-1" onClick={onFind}>
          Find
        </Button>
        <Form.Control
          as="select"
          options={PAGE_SIZE_OPTIONS}
          onChange={(e: React.FormEvent<HTMLInputElement>) =>
            props.onPageSizeChanged(parseInt(e.currentTarget.value, 10))
          }
          defaultValue={props.input.pageSize}
          className="col-1 filter-item"
        >
          {PAGE_SIZE_OPTIONS.map(val => (
            <option key={val} value={val}>
              {val}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    </Form.Group>
  );
};

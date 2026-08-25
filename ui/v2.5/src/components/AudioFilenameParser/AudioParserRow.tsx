import React from "react";
import isEqual from "lodash-es/isEqual";
import clone from "lodash-es/clone";
import { Form } from "react-bootstrap";
import {
  ParseAudioFilenamesQuery,
  SlimAudioDataFragment,
} from "src/core/generated-graphql";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
} from "src/components/Shared/Select";
import cx from "classnames";
import { objectTitle } from "src/core/files";

class ParserResult<T> {
  public value?: T;
  public originalValue?: T;
  public isSet: boolean = false;

  public setOriginalValue(value?: T) {
    this.originalValue = value;
    this.value = value;
  }

  public setValue(value?: T) {
    if (value) {
      this.value = value;
      this.isSet = !isEqual(this.value, this.originalValue);
    }
  }
}

export class AudioParserResult {
  public id: string;
  public filename: string;
  public title: ParserResult<string> = new ParserResult<string>();
  public date: ParserResult<string> = new ParserResult<string>();
  public rating: ParserResult<number> = new ParserResult<number>();

  public studio: ParserResult<string> = new ParserResult<string>();
  public tags: ParserResult<string[]> = new ParserResult<string[]>();
  public performers: ParserResult<string[]> = new ParserResult<string[]>();

  public audio: SlimAudioDataFragment;

  constructor(
    result: ParseAudioFilenamesQuery["parseAudioFilenames"]["results"][0]
  ) {
    this.audio = result.audio;

    this.id = this.audio.id;
    this.filename = objectTitle(this.audio);
    this.title.setOriginalValue(this.audio.title ?? undefined);
    this.date.setOriginalValue(this.audio.date ?? undefined);
    this.rating.setOriginalValue(this.audio.rating100 ?? undefined);
    this.performers.setOriginalValue(this.audio.performers.map((p) => p.id));
    this.tags.setOriginalValue(this.audio.tags.map((t) => t.id));
    this.studio.setOriginalValue(this.audio.studio?.id);

    this.title.setValue(result.title ?? undefined);
    this.date.setValue(result.date ?? undefined);
    this.rating.setValue(result.rating100 ?? undefined);

    this.performers.setValue(result.performer_ids ?? undefined);
    this.tags.setValue(result.tag_ids ?? undefined);
    this.studio.setValue(result.studio_id ?? undefined);
  }

  // returns true if any of its fields have set == true
  public isChanged() {
    return (
      this.title.isSet ||
      this.date.isSet ||
      this.rating.isSet ||
      this.performers.isSet ||
      this.studio.isSet ||
      this.tags.isSet
    );
  }

  public toAudioUpdateInput() {
    return {
      id: this.id,
      rating100: this.rating.isSet ? this.rating.value : undefined,
      title: this.title.isSet ? this.title.value : undefined,
      date: this.date.isSet ? this.date.value : undefined,
      studio_id: this.studio.isSet ? this.studio.value : undefined,
      performer_ids: this.performers.isSet ? this.performers.value : undefined,
      tag_ids: this.tags.isSet ? this.tags.value : undefined,
    };
  }
}

interface IAudioParserFieldProps<T> {
  parserResult: ParserResult<T>;
  className?: string;
  onSetChanged: (isSet: boolean) => void;
  onValueChanged: (value: T) => void;
  originalParserResult?: ParserResult<T>;
}

function AudioParserStringField(props: IAudioParserFieldProps<string>) {
  function maybeValueChanged(value: string) {
    if (value !== props.parserResult.value) {
      props.onValueChanged(value);
    }
  }

  const result = props.originalParserResult || props.parserResult;

  return (
    <>
      <td>
        <Form.Check
          checked={props.parserResult.isSet}
          onChange={() => {
            props.onSetChanged(!props.parserResult.isSet);
          }}
        />
      </td>
      <td>
        <Form.Group>
          <Form.Control
            disabled
            className={props.className}
            defaultValue={result.originalValue || ""}
          />
          <Form.Control
            readOnly={!props.parserResult.isSet}
            className={props.className}
            value={props.parserResult.value || ""}
            onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
              maybeValueChanged(event.currentTarget.value)
            }
          />
        </Form.Group>
      </td>
    </>
  );
}

// audio only has the 1-100 rating, so this is a plain number input rather than
// the 1-5 dropdown the scene parser uses
function AudioParserRatingField(
  props: IAudioParserFieldProps<number | undefined>
) {
  function maybeValueChanged(value?: number) {
    if (value !== props.parserResult.value) {
      props.onValueChanged(value);
    }
  }

  const result = props.originalParserResult || props.parserResult;

  return (
    <>
      <td>
        <Form.Check
          checked={props.parserResult.isSet}
          onChange={() => {
            props.onSetChanged(!props.parserResult.isSet);
          }}
        />
      </td>
      <td>
        <Form.Group>
          <Form.Control
            disabled
            className={cx("input-control text-input", props.className)}
            defaultValue={result.originalValue ?? ""}
          />
          <Form.Control
            type="number"
            min={1}
            max={100}
            className={cx("input-control text-input", props.className)}
            readOnly={!props.parserResult.isSet}
            value={props.parserResult.value ?? ""}
            onChange={(event: React.ChangeEvent<HTMLInputElement>) =>
              maybeValueChanged(
                event.currentTarget.value === ""
                  ? undefined
                  : Number.parseInt(event.currentTarget.value, 10)
              )
            }
          />
        </Form.Group>
      </td>
    </>
  );
}

function AudioParserPerformerField(props: IAudioParserFieldProps<string[]>) {
  function maybeValueChanged(value: string[]) {
    if (value !== props.parserResult.value) {
      props.onValueChanged(value);
    }
  }

  const originalPerformers = (props.originalParserResult?.originalValue ??
    []) as string[];
  const newPerformers = props.parserResult.value ?? [];

  return (
    <>
      <td>
        <Form.Check
          checked={props.parserResult.isSet}
          onChange={() => {
            props.onSetChanged(!props.parserResult.isSet);
          }}
        />
      </td>
      <td>
        <Form.Group className={props.className}>
          <PerformerSelect
            isDisabled
            isMulti
            ids={originalPerformers}
            className="parser-field-performers-select"
          />
          <PerformerSelect
            className="parser-field-performers-select"
            isMulti
            isDisabled={!props.parserResult.isSet}
            onSelect={(items) => {
              maybeValueChanged(items.map((i) => i.id));
            }}
            ids={newPerformers}
          />
        </Form.Group>
      </td>
    </>
  );
}

function AudioParserTagField(props: IAudioParserFieldProps<string[]>) {
  function maybeValueChanged(value: string[]) {
    if (value !== props.parserResult.value) {
      props.onValueChanged(value);
    }
  }

  const originalTags = props.originalParserResult?.originalValue ?? [];
  const newTags = props.parserResult.value ?? [];

  return (
    <>
      <td>
        <Form.Check
          checked={props.parserResult.isSet}
          onChange={() => {
            props.onSetChanged(!props.parserResult.isSet);
          }}
        />
      </td>
      <td>
        <Form.Group className={props.className}>
          <TagSelect
            isDisabled
            isMulti
            ids={originalTags}
            className="parser-field-tags-select"
          />
          <TagSelect
            className="parser-field-tags-select"
            isMulti
            isDisabled={!props.parserResult.isSet}
            onSelect={(items) => {
              maybeValueChanged(items.map((i) => i.id));
            }}
            ids={newTags}
          />
        </Form.Group>
      </td>
    </>
  );
}

function AudioParserStudioField(props: IAudioParserFieldProps<string>) {
  function maybeValueChanged(value: string) {
    if (value !== props.parserResult.value) {
      props.onValueChanged(value);
    }
  }

  const originalStudio = props.originalParserResult?.originalValue
    ? [props.originalParserResult?.originalValue]
    : [];
  const newStudio = props.parserResult.value ? [props.parserResult.value] : [];

  return (
    <>
      <td>
        <Form.Check
          checked={props.parserResult.isSet}
          onChange={() => {
            props.onSetChanged(!props.parserResult.isSet);
          }}
        />
      </td>
      <td>
        <Form.Group className={props.className}>
          <StudioSelect
            isDisabled
            ids={originalStudio}
            className="parser-field-studio-select"
          />
          <StudioSelect
            className="parser-field-studio-select"
            isDisabled={!props.parserResult.isSet}
            onSelect={(items) => {
              maybeValueChanged(items[0].id);
            }}
            ids={newStudio}
          />
        </Form.Group>
      </td>
    </>
  );
}

interface IAudioParserRowProps {
  audio: AudioParserResult;
  onChange: (changedAudio: AudioParserResult) => void;
  showFields: Map<string, boolean>;
}

export const AudioParserRow = (props: IAudioParserRowProps) => {
  function changeParser<T>(result: ParserResult<T>, isSet: boolean, value?: T) {
    const newParser = clone(result);
    newParser.isSet = isSet;
    newParser.value = value;
    return newParser;
  }

  function onTitleChanged(set: boolean, value: string) {
    const newResult = clone(props.audio);
    newResult.title = changeParser(newResult.title, set, value);
    props.onChange(newResult);
  }

  function onDateChanged(set: boolean, value: string) {
    const newResult = clone(props.audio);
    newResult.date = changeParser(newResult.date, set, value);
    props.onChange(newResult);
  }

  function onRatingChanged(set: boolean, value?: number) {
    const newResult = clone(props.audio);
    newResult.rating = changeParser(newResult.rating, set, value);
    props.onChange(newResult);
  }

  function onPerformerIdsChanged(set: boolean, value: string[]) {
    const newResult = clone(props.audio);
    newResult.performers = changeParser(newResult.performers, set, value);
    props.onChange(newResult);
  }

  function onTagIdsChanged(set: boolean, value: string[]) {
    const newResult = clone(props.audio);
    newResult.tags = changeParser(newResult.tags, set, value);
    props.onChange(newResult);
  }

  function onStudioIdChanged(set: boolean, value: string) {
    const newResult = clone(props.audio);
    newResult.studio = changeParser(newResult.studio, set, value);
    props.onChange(newResult);
  }

  return (
    <tr className="scene-parser-row">
      <td className="text-left parser-field-filename">
        {props.audio.filename}
      </td>
      {props.showFields.get("Title") && (
        <AudioParserStringField
          key="title"
          className="parser-field-title input-control text-input"
          parserResult={props.audio.title}
          onSetChanged={(isSet) =>
            onTitleChanged(isSet, props.audio.title.value ?? "")
          }
          onValueChanged={(value) =>
            onTitleChanged(props.audio.title.isSet, value)
          }
        />
      )}
      {props.showFields.get("Date") && (
        <AudioParserStringField
          key="date"
          className="parser-field-date input-control text-input"
          parserResult={props.audio.date}
          onSetChanged={(isSet) =>
            onDateChanged(isSet, props.audio.date.value ?? "")
          }
          onValueChanged={(value) =>
            onDateChanged(props.audio.date.isSet, value)
          }
        />
      )}
      {props.showFields.get("Rating") && (
        <AudioParserRatingField
          key="rating"
          className="parser-field-rating"
          parserResult={props.audio.rating}
          onSetChanged={(isSet) =>
            onRatingChanged(isSet, props.audio.rating.value ?? undefined)
          }
          onValueChanged={(value) =>
            onRatingChanged(props.audio.rating.isSet, value)
          }
        />
      )}
      {props.showFields.get("Performers") && (
        <AudioParserPerformerField
          key="performers"
          className="parser-field-performers"
          parserResult={props.audio.performers}
          originalParserResult={props.audio.performers}
          onSetChanged={(set) =>
            onPerformerIdsChanged(set, props.audio.performers.value ?? [])
          }
          onValueChanged={(value) =>
            onPerformerIdsChanged(props.audio.performers.isSet, value)
          }
        />
      )}
      {props.showFields.get("Tags") && (
        <AudioParserTagField
          key="tags"
          className="parser-field-tags"
          parserResult={props.audio.tags}
          originalParserResult={props.audio.tags}
          onSetChanged={(isSet) =>
            onTagIdsChanged(isSet, props.audio.tags.value ?? [])
          }
          onValueChanged={(value) =>
            onTagIdsChanged(props.audio.tags.isSet, value)
          }
        />
      )}
      {props.showFields.get("Studio") && (
        <AudioParserStudioField
          key="studio"
          className="parser-field-studio"
          parserResult={props.audio.studio}
          originalParserResult={props.audio.studio}
          onSetChanged={(set) =>
            onStudioIdChanged(set, props.audio.studio.value ?? "")
          }
          onValueChanged={(value) =>
            onStudioIdChanged(props.audio.studio.isSet, value)
          }
        />
      )}
    </tr>
  );
};

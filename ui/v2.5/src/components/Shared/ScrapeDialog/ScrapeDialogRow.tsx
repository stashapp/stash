import React, { useContext, useState } from "react";
import {
  Form,
  Col,
  Row,
  InputGroup,
  Button,
  FormControl,
} from "react-bootstrap";
import { Icon } from "../Icon";
import clone from "lodash-es/clone";
import { faCheck, faTimes } from "@fortawesome/free-solid-svg-icons";
import { getCountryByISO } from "src/utils/country";
import { CountrySelect } from "../CountrySelect";
import { StringListInput } from "../StringListInput";
import { ImageSelector } from "../ImageSelector";
import { ScrapeResult } from "./scrapeResult";
import { ScrapeDialogContext } from "./ScrapeDialog";

function renderButtonIcon(selected: boolean) {
  const className = selected ? "text-success" : "text-muted";

  return (
    <Icon
      className={`fa-fw ${className}`}
      icon={selected ? faCheck : faTimes}
    />
  );
}

interface IScrapedFieldProps<T> {
  result: ScrapeResult<T>;
}

interface IScrapedRowProps<T> extends IScrapedFieldProps<T> {
  className?: string;
  field: string;
  title: string;
  originalField: React.ReactNode;
  newField: React.ReactNode;
  onChange: (value: ScrapeResult<T>) => void;
  newValues?: React.ReactNode;
}

export const ScrapeDialogRow = <T,>(props: IScrapedRowProps<T>) => {
  const { existingLabel, scrapedLabel } = useContext(ScrapeDialogContext);

  function handleSelectClick(isNew: boolean) {
    const ret = clone(props.result);
    ret.useNewValue = isNew;
    props.onChange(ret);
  }

  if (!props.result.scraped && !props.newValues) {
    return <></>;
  }

  return (
    <Row
      className={`px-3 pt-3 ${props.className ?? ""}`}
      data-field={props.field}
    >
      <Form.Label column lg="3">
        {props.title}
      </Form.Label>

      <Col lg="9">
        <Row>
          <Form.Label column className="d-lg-none column-label">
            {existingLabel}
          </Form.Label>
          <Col lg="6">
            <InputGroup>
              <InputGroup.Prepend className="bg-secondary text-white border-secondary">
                <Button
                  variant="secondary"
                  onClick={() => handleSelectClick(false)}
                >
                  {renderButtonIcon(!props.result.useNewValue)}
                </Button>
              </InputGroup.Prepend>
              {props.originalField}
            </InputGroup>
          </Col>

          <Form.Label column className="d-lg-none column-label">
            {scrapedLabel}
          </Form.Label>
          <Col lg="6">
            <InputGroup>
              <InputGroup.Prepend>
                <Button
                  variant="secondary"
                  onClick={() => handleSelectClick(true)}
                >
                  {renderButtonIcon(props.result.useNewValue)}
                </Button>
              </InputGroup.Prepend>
              {props.newField}
            </InputGroup>
            {props.newValues}
          </Col>
        </Row>
      </Col>
    </Row>
  );
};

interface IScrapedInputGroupProps {
  isNew?: boolean;
  placeholder?: string;
  locked?: boolean;
  result: ScrapeResult<string>;
  onChange?: (value: string) => void;
}

const ScrapedInputGroup: React.FC<IScrapedInputGroupProps> = (props) => {
  return (
    <FormControl
      placeholder={props.placeholder}
      value={props.isNew ? props.result.newValue : props.result.originalValue}
      readOnly={!props.isNew || props.locked}
      onChange={(e) => {
        if (props.isNew && props.onChange) {
          props.onChange(e.target.value);
        }
      }}
      className="bg-secondary text-white border-secondary"
    />
  );
};

interface IScrapedInputGroupRowProps {
  title: string;
  field: string;
  className?: string;
  placeholder?: string;
  result: ScrapeResult<string>;
  locked?: boolean;
  onChange: (value: ScrapeResult<string>) => void;
}

export const ScrapedInputGroupRow: React.FC<IScrapedInputGroupRowProps> = (
  props
) => {
  return (
    <ScrapeDialogRow
      title={props.title}
      field={props.field}
      className={props.className}
      result={props.result}
      originalField={
        <ScrapedInputGroup
          placeholder={props.placeholder || props.title}
          result={props.result}
        />
      }
      newField={
        <ScrapedInputGroup
          placeholder={props.placeholder || props.title}
          result={props.result}
          isNew
          locked={props.locked}
          onChange={(value) =>
            props.onChange(props.result.cloneWithValue(value))
          }
        />
      }
      onChange={props.onChange}
    />
  );
};

interface IScrapedStringListProps {
  isNew?: boolean;
  placeholder?: string;
  locked?: boolean;
  result: ScrapeResult<string[]>;
  onChange?: (value: string[]) => void;
}

const ScrapedStringList: React.FC<IScrapedStringListProps> = (props) => {
  const value = props.isNew
    ? props.result.newValue
    : props.result.originalValue;

  return (
    <StringListInput
      value={value ?? []}
      setValue={(v) => {
        if (props.isNew && props.onChange) {
          props.onChange(v);
        }
      }}
      placeholder={props.placeholder}
      readOnly={!props.isNew || props.locked}
    />
  );
};

interface IScrapedStringListRowProps {
  title: string;
  field: string;
  placeholder?: string;
  result: ScrapeResult<string[]>;
  locked?: boolean;
  onChange: (value: ScrapeResult<string[]>) => void;
}

export const ScrapedStringListRow: React.FC<IScrapedStringListRowProps> = (
  props
) => {
  return (
    <ScrapeDialogRow
      className="string-list-row"
      title={props.title}
      field={props.field}
      result={props.result}
      originalField={
        <ScrapedStringList
          placeholder={props.placeholder || props.title}
          result={props.result}
        />
      }
      newField={
        <ScrapedStringList
          placeholder={props.placeholder || props.title}
          result={props.result}
          isNew
          locked={props.locked}
          onChange={(value) =>
            props.onChange(props.result.cloneWithValue(value))
          }
        />
      }
      onChange={props.onChange}
    />
  );
};

const ScrapedTextArea: React.FC<IScrapedInputGroupProps> = (props) => {
  return (
    <FormControl
      as="textarea"
      placeholder={props.placeholder}
      value={props.isNew ? props.result.newValue : props.result.originalValue}
      readOnly={!props.isNew}
      onChange={(e) => {
        if (props.isNew && props.onChange) {
          props.onChange(e.target.value);
        }
      }}
      className="bg-secondary text-white border-secondary scene-description"
    />
  );
};

export const ScrapedTextAreaRow: React.FC<IScrapedInputGroupRowProps> = (
  props
) => {
  return (
    <ScrapeDialogRow
      title={props.title}
      field={props.field}
      result={props.result}
      originalField={
        <ScrapedTextArea
          placeholder={props.placeholder || props.title}
          result={props.result}
        />
      }
      newField={
        <ScrapedTextArea
          placeholder={props.placeholder || props.title}
          result={props.result}
          isNew
          onChange={(value) =>
            props.onChange(props.result.cloneWithValue(value))
          }
        />
      }
      onChange={props.onChange}
    />
  );
};

interface IScrapedImageProps {
  isNew?: boolean;
  className?: string;
  placeholder?: string;
  result: ScrapeResult<string>;
}

const ScrapedImage: React.FC<IScrapedImageProps> = (props) => {
  const value = props.isNew
    ? props.result.newValue
    : props.result.originalValue;

  if (!value) {
    return <></>;
  }

  return (
    <img className={props.className} src={value} alt={props.placeholder} />
  );
};

interface IScrapedImageRowProps {
  title: string;
  field: string;
  className?: string;
  result: ScrapeResult<string>;
  onChange: (value: ScrapeResult<string>) => void;
}

export const ScrapedImageRow: React.FC<IScrapedImageRowProps> = (props) => {
  return (
    <ScrapeDialogRow
      title={props.title}
      field={props.field}
      result={props.result}
      originalField={
        <ScrapedImage
          result={props.result}
          className={props.className}
          placeholder={props.title}
        />
      }
      newField={
        <ScrapedImage
          result={props.result}
          className={props.className}
          placeholder={props.title}
          isNew
        />
      }
      onChange={props.onChange}
    />
  );
};

interface IScrapedImagesRowProps {
  title: string;
  field: string;
  className?: string;
  result: ScrapeResult<string>;
  images: string[];
  onChange: (value: ScrapeResult<string>) => void;
}

export const ScrapedImagesRow: React.FC<IScrapedImagesRowProps> = (props) => {
  const [imageIndex, setImageIndex] = useState(0);

  function onSetImageIndex(newIdx: number) {
    const ret = props.result.cloneWithValue(props.images[newIdx]);
    props.onChange(ret);
    setImageIndex(newIdx);
  }

  return (
    <ScrapeDialogRow
      title={props.title}
      field={props.field}
      result={props.result}
      originalField={
        <ScrapedImage
          result={props.result}
          className={props.className}
          placeholder={props.title}
        />
      }
      newField={
        <div className="image-selection-parent">
          <ImageSelector
            imageClassName={props.className}
            images={props.images}
            imageIndex={imageIndex}
            setImageIndex={onSetImageIndex}
          />
        </div>
      }
      onChange={props.onChange}
    />
  );
};

interface IScrapedCountryRowProps {
  title: string;
  field: string;
  result: ScrapeResult<string>;
  onChange: (value: ScrapeResult<string>) => void;
  locked?: boolean;
  locale?: string;
}

export const ScrapedCountryRow: React.FC<IScrapedCountryRowProps> = ({
  title,
  field,
  result,
  onChange,
  locked,
  locale,
}) => (
  <ScrapeDialogRow
    title={title}
    field={field}
    result={result}
    originalField={
      <FormControl
        value={
          getCountryByISO(result.originalValue, locale) ?? result.originalValue
        }
        readOnly
        className="bg-secondary text-white border-secondary"
      />
    }
    newField={
      <CountrySelect
        value={result.newValue}
        disabled={locked}
        onChange={(value) => {
          if (onChange) {
            onChange(result.cloneWithValue(value));
          }
        }}
        showFlag={false}
        isClearable={false}
        className="flex-grow-1"
      />
    }
    onChange={onChange}
  />
);

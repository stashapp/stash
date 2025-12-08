import React, { useState } from "react";
import {
  Form,
  Col,
  Row,
  InputGroup,
  Button,
  FormControl,
  Badge,
} from "react-bootstrap";
import { CollapseButton } from "../CollapseButton";
import { Icon } from "../Icon";
import clone from "lodash-es/clone";
import { faCheck, faPlus, faTimes } from "@fortawesome/free-solid-svg-icons";
import { getCountryByISO } from "src/utils/country";
import { CountrySelect } from "../CountrySelect";
import { StringListInput } from "../StringListInput";
import { ImageSelector } from "../ImageSelector";
import { ScrapeResult } from "./scrapeResult";

interface IScrapedFieldProps<T> {
  result: ScrapeResult<T>;
}

interface IScrapedRowProps<T, V> extends IScrapedFieldProps<T> {
  className?: string;
  field: string;
  title: string;
  renderOriginalField: (result: ScrapeResult<T>) => JSX.Element | undefined;
  renderNewField: (result: ScrapeResult<T>) => JSX.Element | undefined;
  onChange: (value: ScrapeResult<T>) => void;
  newValues?: V[];
  onCreateNew?: (index: number) => void;
  getName?: (value: V) => string;
}

function renderButtonIcon(selected: boolean) {
  const className = selected ? "text-success" : "text-muted";

  return (
    <Icon
      className={`fa-fw ${className}`}
      icon={selected ? faCheck : faTimes}
    />
  );
}

export const ScrapeDialogRow = <T, V>(props: IScrapedRowProps<T, V>) => {
  const { getName = () => "" } = props;

  function handleSelectClick(isNew: boolean) {
    const ret = clone(props.result);
    ret.useNewValue = isNew;
    props.onChange(ret);
  }

  function hasNewValues() {
    return props.newValues && props.newValues.length > 0 && props.onCreateNew;
  }

  if (!props.result.scraped && !hasNewValues()) {
    return <></>;
  }

  function renderNewValues() {
    if (!hasNewValues()) {
      return;
    }

    const ret = (
      <>
        {props.newValues!.map((t, i) => (
          <Badge
            className="tag-item"
            variant="secondary"
            key={getName(t)}
            onClick={() => props.onCreateNew!(i)}
          >
            {getName(t)}
            <Button className="minimal ml-2">
              <Icon className="fa-fw" icon={faPlus} />
            </Button>
          </Badge>
        ))}
      </>
    );

    const minCollapseLength = 10;

    if (props.newValues!.length >= minCollapseLength) {
      return (
        <CollapseButton text={`Missing (${props.newValues!.length})`}>
          {ret}
        </CollapseButton>
      );
    }

    return ret;
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
          <Col xs="6">
            <InputGroup>
              <InputGroup.Prepend className="bg-secondary text-white border-secondary">
                <Button
                  variant="secondary"
                  onClick={() => handleSelectClick(false)}
                >
                  {renderButtonIcon(!props.result.useNewValue)}
                </Button>
              </InputGroup.Prepend>
              {props.renderOriginalField(props.result)}
            </InputGroup>
          </Col>
          <Col xs="6">
            <InputGroup>
              <InputGroup.Prepend>
                <Button
                  variant="secondary"
                  onClick={() => handleSelectClick(true)}
                >
                  {renderButtonIcon(props.result.useNewValue)}
                </Button>
              </InputGroup.Prepend>
              {props.renderNewField(props.result)}
            </InputGroup>
            {renderNewValues()}
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

function getNameString(value: string) {
  return value;
}

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
      renderOriginalField={() => (
        <ScrapedInputGroup
          placeholder={props.placeholder || props.title}
          result={props.result}
        />
      )}
      renderNewField={() => (
        <ScrapedInputGroup
          placeholder={props.placeholder || props.title}
          result={props.result}
          isNew
          locked={props.locked}
          onChange={(value) =>
            props.onChange(props.result.cloneWithValue(value))
          }
        />
      )}
      onChange={props.onChange}
      getName={getNameString}
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
      renderOriginalField={() => (
        <ScrapedStringList
          placeholder={props.placeholder || props.title}
          result={props.result}
        />
      )}
      renderNewField={() => (
        <ScrapedStringList
          placeholder={props.placeholder || props.title}
          result={props.result}
          isNew
          locked={props.locked}
          onChange={(value) =>
            props.onChange(props.result.cloneWithValue(value))
          }
        />
      )}
      onChange={props.onChange}
      getName={getNameString}
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
      renderOriginalField={() => (
        <ScrapedTextArea
          placeholder={props.placeholder || props.title}
          result={props.result}
        />
      )}
      renderNewField={() => (
        <ScrapedTextArea
          placeholder={props.placeholder || props.title}
          result={props.result}
          isNew
          onChange={(value) =>
            props.onChange(props.result.cloneWithValue(value))
          }
        />
      )}
      onChange={props.onChange}
      getName={getNameString}
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
      renderOriginalField={() => (
        <ScrapedImage
          result={props.result}
          className={props.className}
          placeholder={props.title}
        />
      )}
      renderNewField={() => (
        <ScrapedImage
          result={props.result}
          className={props.className}
          placeholder={props.title}
          isNew
        />
      )}
      onChange={props.onChange}
      getName={getNameString}
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
      renderOriginalField={() => (
        <ScrapedImage
          result={props.result}
          className={props.className}
          placeholder={props.title}
        />
      )}
      renderNewField={() => (
        <div className="image-selection-parent">
          <ImageSelector
            imageClassName={props.className}
            images={props.images}
            imageIndex={imageIndex}
            setImageIndex={onSetImageIndex}
          />
        </div>
      )}
      onChange={props.onChange}
      getName={getNameString}
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
    renderOriginalField={() => (
      <FormControl
        value={
          getCountryByISO(result.originalValue, locale) ?? result.originalValue
        }
        readOnly
        className="bg-secondary text-white border-secondary"
      />
    )}
    renderNewField={() => (
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
    )}
    onChange={onChange}
    getName={getNameString}
  />
);

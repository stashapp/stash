import React, { useState } from "react";
import { Form, Col, Row, InputGroup, Button, FormControl } from "react-bootstrap";
import { Modal, Icon, StudioSelect, PerformerSelect } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import _ from "lodash";
import { MovieSelect, TagSelect } from "src/components/Shared/Select";

class ScrapeResult<T> {
  public newValue?: T;
  public originalValue?: T;
  public scraped: boolean = false;
  public useNewValue: boolean = false;

  public constructor(originalValue?: T | null, newValue?: T | null) {
    this.originalValue = originalValue ?? undefined;
    this.newValue = newValue ?? undefined;
    this.useNewValue = !!this.newValue && this.newValue !== this.originalValue;
    this.scraped = this.useNewValue;
  }

  public setOriginalValue(value?: T) {
    this.originalValue = value;
    this.newValue = value;
  }

  public cloneWithValue(value?: T) {
    const ret = _.clone(this);

    ret.newValue = value;
    ret.useNewValue = !_.isEqual(ret.newValue, ret.originalValue);

    return ret;
  }

  public getNewValue() {
    if (this.useNewValue) {
      return this.newValue;
    }
  }
}

interface IScrapedFieldProps<T> {
  result: ScrapeResult<T>
}

interface IScrapedRowProps<T> extends IScrapedFieldProps<T> {
  title: string;
  renderOriginalField: (result: ScrapeResult<T>) => JSX.Element | undefined;
  renderNewField: (result: ScrapeResult<T>) => JSX.Element | undefined;
  onChange: (value: ScrapeResult<T>) => void;
}

function renderButtonIcon(selected: boolean) {
  let className = selected ? "text-success" : "text-muted";

  return (
    <Icon className={`fa-fw ${className}`} icon={selected ? "check" : "times"} />
  );
}

function renderScrapedRow<T>(props: IScrapedRowProps<T>) {
  function handleSelectClick(isNew: boolean) {
    const ret = _.clone(props.result);
    ret.useNewValue = isNew;
    props.onChange(ret);
  }

  if (!props.result.scraped) {
    return;
  }

  return (
    <Row className="px-3 pt-3">
      <Form.Label column lg="3">
        {props.title}
      </Form.Label>
      
      <Col lg="9">
        <Row>
          <Col xs="6">
            <InputGroup>
              <InputGroup.Prepend className="bg-secondary text-white border-secondary">
                <Button variant="secondary" onClick={() => handleSelectClick(false)}>
                  {renderButtonIcon(!props.result.useNewValue)}
                </Button>
              </InputGroup.Prepend>
              {props.renderOriginalField(props.result)}
            </InputGroup>
            
          </Col>
          <Col xs="6">
            <InputGroup>
              <InputGroup.Prepend>
                <Button variant="secondary" onClick={() => handleSelectClick(true)}>
                  {renderButtonIcon(props.result.useNewValue)}
                </Button>
              </InputGroup.Prepend>
              {props.renderNewField(props.result)}
            </InputGroup>
          </Col>
        </Row>
      </Col>
    </Row>
  );
}

interface IScrapedInputGroupProps {
  title: string;
  placeholder?: string
  result: ScrapeResult<string>;
}

function renderScrapedInputGroup(props: IScrapedInputGroupProps, isNew?: boolean, onChange?: (value : string) => void) {
  return (
    <FormControl
      placeholder={props.placeholder ?? props.title}
      value={isNew ? props.result.newValue : props.result.originalValue}
      readOnly={!isNew}
      onChange={isNew && onChange ? (e) => onChange(e.target.value) : () => {}}
      className="bg-secondary text-white border-secondary"
    />
  );
}

function renderScrapedInputGroupRow(props: IScrapedInputGroupProps, onChange: (value : ScrapeResult<string>) => void) {
  return renderScrapedRow({
    title: props.title,
    result: props.result,
    renderOriginalField: () => renderScrapedInputGroup(props),
    renderNewField: () => renderScrapedInputGroup(props, true, (value) => onChange(props.result.cloneWithValue(value))),
    onChange,
  });
}

function renderScrapedStudio(result: ScrapeResult<string>, isNew?: boolean, onChange?: (value : string) => void) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ? [resultValue] : [];

  return (
    <StudioSelect
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items[0]?.id);
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedStudioRow(result: ScrapeResult<string>, onChange: (value : ScrapeResult<string>) => void) {
  return renderScrapedRow({
    title: "Studio",
    result,
    renderOriginalField: () => renderScrapedStudio(result),
    renderNewField: () => renderScrapedStudio(result, true, (value) => onChange(result.cloneWithValue(value))),
    onChange,
  });
}

function renderScrapedPerformers(result: ScrapeResult<string[]>, isNew?: boolean, onChange?: (value : string[]) => void) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <PerformerSelect
      isMulti={true}
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedPerformersRow(result: ScrapeResult<string[]>, onChange: (value : ScrapeResult<string[]>) => void) {
  return renderScrapedRow({
    title: "Performers",
    result,
    renderOriginalField: () => renderScrapedPerformers(result),
    renderNewField: () => renderScrapedPerformers(result, true, (value) => onChange(result.cloneWithValue(value))),
    onChange,
  });
}

function renderScrapedMovies(result: ScrapeResult<string[]>, isNew?: boolean, onChange?: (value : string[]) => void) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <MovieSelect
      isMulti={true}
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedMoviesRow(result: ScrapeResult<string[]>, onChange: (value : ScrapeResult<string[]>) => void) {
  return renderScrapedRow({
    title: "Movies",
    result,
    renderOriginalField: () => renderScrapedMovies(result),
    renderNewField: () => renderScrapedMovies(result, true, (value) => onChange(result.cloneWithValue(value))),
    onChange,
  });
}

function renderScrapedTags(result: ScrapeResult<string[]>, isNew?: boolean, onChange?: (value : string[]) => void) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <TagSelect
      isMulti={true}
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedTagsRow(result: ScrapeResult<string[]>, onChange: (value : ScrapeResult<string[]>) => void) {
  return renderScrapedRow({
    title: "Tags",
    result,
    renderOriginalField: () => renderScrapedTags(result),
    renderNewField: () => renderScrapedTags(result, true, (value) => onChange(result.cloneWithValue(value))),
    onChange,
  });
}

function renderScrapedTextArea(props: IScrapedInputGroupProps, isNew?: boolean, onChange?: (value : string) => void) {
  return (
    <FormControl as="textarea"
      placeholder={props.placeholder ?? props.title}
      value={isNew ? props.result.newValue : props.result.originalValue}
      readOnly={!isNew}
      onChange={isNew && onChange ? (e) => onChange(e.target.value) : () => {}}
      className="bg-secondary text-white border-secondary scene-description"
    />
  );
}

function renderScrapedTextAreaRow(props: IScrapedInputGroupProps, onChange: (value : ScrapeResult<string>) => void) {
  return renderScrapedRow({
    title: props.title,
    result: props.result,
    renderOriginalField: () => renderScrapedTextArea(props),
    renderNewField: () => renderScrapedTextArea(props, true, (value) => onChange(props.result.cloneWithValue(value))),
    onChange,
  });
}

function renderScrapedImage(result: ScrapeResult<string>, isNew?: boolean) {
  const value = isNew ? result.newValue : result.originalValue;

  if (!value) {
    return;
  }
  
  return (
    <img
      className="scene-cover"
      src={value}
      alt="Scene cover"
    />
  );
}

function renderScrapedImageRow(result: ScrapeResult<string>, onChange: (value : ScrapeResult<string>) => void) {
  return renderScrapedRow({
    title: "Cover Image",
    result: result,
    renderOriginalField: () => renderScrapedImage(result),
    renderNewField: () => renderScrapedImage(result, true),
    onChange,
  });
}

interface ISceneScrapeDialogProps {
  scene: Partial<GQL.SceneDataFragment>
  scraped: GQL.ScrapedScene;

  onClose: (scrapedScene?: GQL.ScrapedScene) => void;
}

interface HasID {
  id?: string | null;
}

export const SceneScrapeDialog: React.FC<ISceneScrapeDialogProps> = (
  props: ISceneScrapeDialogProps
) => {
  const [title, setTitle] = useState<ScrapeResult<string>>(new ScrapeResult<string>(props.scene.title, props.scraped.title));
  const [url, setURL] = useState<ScrapeResult<string>>(new ScrapeResult<string>(props.scene.url, props.scraped.url));
  const [date, setDate] = useState<ScrapeResult<string>>(new ScrapeResult<string>(props.scene.date, props.scraped.date));
  const [studio, setStudio] = useState<ScrapeResult<string>>(new ScrapeResult<string>(props.scene.studio?.id, props.scraped.studio?.id));
  
  function mapIdObjects(scrapedObjects?: HasID[]): string[] | undefined {
    if (!scrapedObjects) {
      return undefined;
    }
    const ret = scrapedObjects.map(p => p.id).filter(p => {
      return p !== undefined && p !== null;
    }) as string[];

    if (ret.length === 0) {
      return undefined;
    }
  }

  const [performers, setPerformers] = useState<ScrapeResult<string[]>>(new ScrapeResult<string[]>(props.scene.performers?.map(p => p.id), mapIdObjects(props.scraped.performers ?? undefined)));
  const [movies, setMovies] = useState<ScrapeResult<string[]>>(new ScrapeResult<string[]>(props.scene.movies?.map(p => p.movie.id), mapIdObjects(props.scraped.movies ?? undefined)));
  const [tags, setTags] = useState<ScrapeResult<string[]>>(new ScrapeResult<string[]>(props.scene.tags?.map(p => p.id), mapIdObjects(props.scraped.tags ?? undefined)));
  const [details, setDetails] = useState<ScrapeResult<string>>(new ScrapeResult<string>(props.scene.details, props.scraped.details));
  const [image, setImage] = useState<ScrapeResult<string>>(new ScrapeResult<string>(props.scene.paths?.screenshot, props.scraped.image));

  // don't show the dialog if nothing was scraped
  if ([title, url, date, studio, performers, movies, tags, details, image].every(r => !r.scraped)) {
    props.onClose();
    return <></>;
  }

  function makeNewScrapedItem() {
    const newStudio = studio.getNewValue();
    
    return {
      title: title.getNewValue(),
      url: url.getNewValue(),
      date: date.getNewValue(),
      studio: newStudio ? {
        id: newStudio,
        name: "",
      } : undefined,
      performers: performers.getNewValue()?.map(p => {
        return {
          id: p,
          name: "",
        };
      }),
      movies: movies.getNewValue()?.map(m => {
        return {
          id: m,
          name: "",
        };
      }),
      tags: tags.getNewValue()?.map(m => {
        return {
          id: m,
          name: "",
        };
      }),
      details: details.getNewValue(),
      image: image.getNewValue(),
    };
  }

  return (
    <Modal
      show
      icon="pencil-alt"
      header="Scene Scrape Results"
      accept={{ onClick: () => { props.onClose(makeNewScrapedItem())}, text: "Apply" }}
      cancel={{
        onClick: () => props.onClose(),
        text: "Cancel",
        variant: "secondary",
      }}
      modalProps={{size: "lg", dialogClassName: "scrape-dialog"}}
    >
      <div className="dialog-container">
        <Form>
          <Row className="px-3 pt-3">
            <Col lg={{span: 9, offset: 3}}>
              <Row>
                <Form.Label column xs="6">
                  Existing
                </Form.Label>
                <Form.Label column xs="6">
                  Scraped
                </Form.Label>
              </Row>
            </Col>
          </Row>

          {renderScrapedInputGroupRow({title: "Title", result: title}, (value) => setTitle(value))}
          {renderScrapedInputGroupRow({title: "URL", result: url}, (value) => setURL(value))}
          {renderScrapedInputGroupRow({title: "Date", result: date, placeholder: "YYYY-MM-DD"}, (value) => setDate(value))}
          {renderScrapedStudioRow(studio, (value) => setStudio(value))}
          {renderScrapedPerformersRow(performers, (value) => setPerformers(value))}
          {renderScrapedMoviesRow(movies, (value) => setMovies(value))}
          {renderScrapedTagsRow(tags, (value) => setTags(value))}
          {renderScrapedTextAreaRow({title: "Details", result: details}, (value) => setDetails(value))}
          {renderScrapedImageRow(image, (value) => setImage(value))}
        </Form>
      </div>
    </Modal>
  );
};

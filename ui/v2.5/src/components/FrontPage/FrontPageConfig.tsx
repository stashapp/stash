import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, IntlShape, useIntl } from "react-intl";
import { useFindSavedFilters } from "src/core/StashService";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { Button, Form, Modal } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import {
  ISavedFilterRow,
  ICustomFilter,
  FrontPageContent,
  generatePremadeFrontPageContent,
  getFrontPageContent,
} from "src/core/config";

interface IAddSavedFilterModalProps {
  onClose: (content?: FrontPageContent) => void;
  existingSavedFilterIDs: string[];
  candidates: GQL.FindSavedFiltersQuery;
}

const FilterModeToMessageID = {
  [GQL.FilterMode.Galleries]: "galleries",
  [GQL.FilterMode.Images]: "images",
  [GQL.FilterMode.Movies]: "groups",
  [GQL.FilterMode.Groups]: "groups",
  [GQL.FilterMode.Performers]: "performers",
  [GQL.FilterMode.SceneMarkers]: "markers",
  [GQL.FilterMode.Scenes]: "scenes",
  [GQL.FilterMode.Studios]: "studios",
  [GQL.FilterMode.Tags]: "tags",
};

type SavedFilter = Pick<GQL.SavedFilter, "id" | "mode" | "name">;

function filterTitle(intl: IntlShape, f: SavedFilter) {
  const typeMessage = intl.formatMessage({ id: FilterModeToMessageID[f.mode] });
  return `${typeMessage}: ${f.name}`;
}

const AddContentModal: React.FC<IAddSavedFilterModalProps> = ({
  onClose,
  existingSavedFilterIDs,
  candidates,
}) => {
  const intl = useIntl();

  const premadeFilterOptions = useMemo(
    () => generatePremadeFrontPageContent(intl),
    [intl]
  );

  const [contentType, setContentType] = useState(
    "front_page.types.premade_filter"
  );
  const [premadeFilterIndex, setPremadeFilterIndex] = useState<
    number | undefined
  >(0);
  const [savedFilter, setSavedFilter] = useState<string | undefined>();

  function onTypeSelected(t: string) {
    setContentType(t);

    switch (t) {
      case "front_page.types.premade_filter":
        setPremadeFilterIndex(0);
        setSavedFilter(undefined);
        break;
      case "front_page.types.saved_filter":
        setPremadeFilterIndex(undefined);
        setSavedFilter(undefined);
        break;
    }
  }

  function isValid() {
    switch (contentType) {
      case "front_page.types.premade_filter":
        return premadeFilterIndex !== undefined;
      case "front_page.types.saved_filter":
        return savedFilter !== undefined;
    }

    return false;
  }

  const savedFilterOptions = useMemo(() => {
    const ret = [
      {
        value: "",
        text: "",
      },
    ].concat(
      candidates.findSavedFilters
        .filter((f) => {
          return !existingSavedFilterIDs.includes(f.id);
        })
        .map((f) => {
          return {
            value: f.id,
            text: filterTitle(intl, f),
          };
        })
    );

    ret.sort((a, b) => {
      return a.text.localeCompare(b.text);
    });

    return ret;
  }, [candidates, existingSavedFilterIDs, intl]);

  function renderTypeSelect() {
    const options = [
      "front_page.types.premade_filter",
      "front_page.types.saved_filter",
    ];
    return (
      <Form.Group controlId="filter">
        <Form.Label>
          <FormattedMessage id="type" />
        </Form.Label>
        <Form.Control
          as="select"
          value={contentType}
          onChange={(e) => onTypeSelected(e.target.value)}
          className="btn-secondary"
        >
          {options.map((c) => (
            <option key={c} value={c}>
              {intl.formatMessage({ id: c })}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function maybeRenderPremadeFiltersSelect() {
    if (contentType !== "front_page.types.premade_filter") return;

    return (
      <Form.Group controlId="premade-filter">
        <Form.Label>
          <FormattedMessage id="front_page.types.premade_filter" />
        </Form.Label>
        <Form.Control
          as="select"
          value={premadeFilterIndex}
          onChange={(e) => setPremadeFilterIndex(parseInt(e.target.value))}
          className="btn-secondary"
        >
          {premadeFilterOptions.map((c, i) => (
            <option key={i} value={i}>
              {intl.formatMessage({ id: c.message!.id }, c.message!.values)}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function maybeRenderSavedFiltersSelect() {
    if (contentType !== "front_page.types.saved_filter") return;
    return (
      <Form.Group controlId="filter">
        <Form.Label>
          <FormattedMessage id="search_filter.name" />
        </Form.Label>
        <Form.Control
          as="select"
          value={savedFilter}
          onChange={(e) => setSavedFilter(e.target.value)}
          className="btn-secondary"
        >
          {savedFilterOptions.map((c) => (
            <option key={c.value} value={c.value}>
              {c.text}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function doAdd() {
    switch (contentType) {
      case "front_page.types.premade_filter":
        onClose(premadeFilterOptions[premadeFilterIndex!]);
        return;
      case "front_page.types.saved_filter":
        onClose({
          __typename: "SavedFilter",
          savedFilterId: parseInt(savedFilter!),
        });
        return;
    }

    onClose();
  }

  return (
    <Modal show onHide={() => onClose()}>
      <Modal.Header>
        <FormattedMessage id="actions.add" />
      </Modal.Header>
      <Modal.Body>
        <div className="dialog-content">
          {renderTypeSelect()}
          {maybeRenderSavedFiltersSelect()}
          {maybeRenderPremadeFiltersSelect()}
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={() => onClose()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
        <Button onClick={() => doAdd()} disabled={!isValid()}>
          <FormattedMessage id="actions.add" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

interface IFilterRowProps {
  content: FrontPageContent;
  allSavedFilters: SavedFilter[];
  onDelete: () => void;
}

const ContentRow: React.FC<IFilterRowProps> = (props: IFilterRowProps) => {
  const intl = useIntl();

  function title() {
    switch (props.content.__typename) {
      case "SavedFilter":
        const savedFilterId = String(props.content.savedFilterId);
        const savedFilter = props.allSavedFilters.find(
          (f) => f.id === savedFilterId
        );
        if (!savedFilter) return "";
        return filterTitle(intl, savedFilter);
      case "CustomFilter":
        const asCustomFilter = props.content as ICustomFilter;
        if (asCustomFilter.message)
          return intl.formatMessage(
            { id: asCustomFilter.message.id },
            asCustomFilter.message.values
          );
        return asCustomFilter.title ?? "";
    }
  }

  return (
    <div className="recommendation-row">
      <div className="recommendation-row-head">
        <div>
          <h2>{title()}</h2>
        </div>
        <Button
          variant="danger"
          title={intl.formatMessage({ id: "actions.delete" })}
          onClick={() => props.onDelete()}
        >
          <FormattedMessage id="actions.delete" />
        </Button>
      </div>
    </div>
  );
};

interface IFrontPageConfigProps {
  onClose: (content?: FrontPageContent[]) => void;
}

export const FrontPageConfig: React.FC<IFrontPageConfigProps> = ({
  onClose,
}) => {
  const { configuration } = useConfigurationContext();

  const ui = configuration?.ui;

  const { data: allFilters, loading } = useFindSavedFilters();

  const [isAdd, setIsAdd] = useState(false);
  const [currentContent, setCurrentContent] = useState<FrontPageContent[]>([]);
  const [dragIndex, setDragIndex] = useState<number | undefined>();

  useEffect(() => {
    if (!allFilters?.findSavedFilters) {
      return;
    }

    const frontPageContent = getFrontPageContent(ui);
    if (frontPageContent) {
      setCurrentContent(
        // filter out rows where the saved filter no longer exists
        frontPageContent.filter((r) => {
          if (r.__typename === "SavedFilter") {
            const savedFilterId = String(r.savedFilterId);
            return allFilters.findSavedFilters.some(
              (f) => f.id === savedFilterId
            );
          }
          return true;
        })
      );
    }
  }, [allFilters, ui]);

  function onDragStart(event: React.DragEvent<HTMLElement>, index: number) {
    event.dataTransfer.effectAllowed = "move";
    setDragIndex(index);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>, index?: number) {
    if (dragIndex !== undefined && index !== undefined && index !== dragIndex) {
      const newFilters = [...currentContent];
      const moved = newFilters.splice(dragIndex, 1);
      newFilters.splice(index, 0, moved[0]);
      setCurrentContent(newFilters);
      setDragIndex(index);
    }

    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDragOverDefault(event: React.DragEvent<HTMLDivElement>) {
    event.dataTransfer.dropEffect = "move";
    event.preventDefault();
  }

  function onDrop() {
    // assume we've already set the temp filter list
    // feed it up
    setDragIndex(undefined);
  }

  if (loading) {
    return <LoadingIndicator />;
  }

  const existingSavedFilterIDs = currentContent
    .filter(
      (f) =>
        f.__typename === "SavedFilter" && (f as ISavedFilterRow).savedFilterId
    )
    .map((f) => (f as ISavedFilterRow).savedFilterId.toString());

  function addSavedFilter(content?: FrontPageContent) {
    setIsAdd(false);

    if (!content) {
      return;
    }

    setCurrentContent([...currentContent, content]);
  }

  function deleteSavedFilter(index: number) {
    setCurrentContent(currentContent.filter((f, i) => i !== index));
  }

  return (
    <>
      {isAdd && allFilters && (
        <AddContentModal
          candidates={allFilters}
          existingSavedFilterIDs={existingSavedFilterIDs}
          onClose={addSavedFilter}
        />
      )}
      <div className="recommendations-container recommendations-container-edit">
        <div onDragOver={onDragOverDefault}>
          {currentContent.map((content, index) => (
            <div
              key={index}
              draggable
              onDragStart={(e) => onDragStart(e, index)}
              onDragEnter={(e) => onDragOver(e, index)}
              onDrop={() => onDrop()}
            >
              <ContentRow
                key={index}
                allSavedFilters={allFilters!.findSavedFilters}
                content={content}
                onDelete={() => deleteSavedFilter(index)}
              />
            </div>
          ))}
          <div className="recommendation-row recommendation-row-add">
            <div className="recommendation-row-head">
              <Button
                className="recommendations-add"
                variant="primary"
                onClick={() => setIsAdd(true)}
              >
                <FormattedMessage id="actions.add" />
              </Button>
            </div>
          </div>
        </div>
        <div className="recommendations-footer">
          <Button onClick={() => onClose()} variant="secondary">
            <FormattedMessage id={"actions.cancel"} />
          </Button>
          <Button onClick={() => onClose(currentContent)}>
            <FormattedMessage id={"actions.save"} />
          </Button>
        </div>
      </div>
    </>
  );
};

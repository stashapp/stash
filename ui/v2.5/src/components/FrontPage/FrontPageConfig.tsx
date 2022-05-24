import React, { useEffect, useState } from "react";
import { FormattedMessage, IntlShape, useIntl } from "react-intl";
import {
  useFindFrontPageFiltersQuery,
  useFindSavedFilters,
} from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared";
import { Button, Form, Modal } from "react-bootstrap";
import {
  FilterMode,
  FindSavedFiltersQuery,
  SavedFilter,
} from "src/core/generated-graphql";

interface IAddSavedFilterModalProps {
  onClose: (id?: string) => void;
  existingSavedFilterIDs: string[];
  candidates: FindSavedFiltersQuery;
}

const FilterModeToMessageID = {
  [FilterMode.Galleries]: "galleries",
  [FilterMode.Images]: "images",
  [FilterMode.Movies]: "movies",
  [FilterMode.Performers]: "performers",
  [FilterMode.SceneMarkers]: "markers",
  [FilterMode.Scenes]: "scenes",
  [FilterMode.Studios]: "studios",
  [FilterMode.Tags]: "tags",
};

function filterTitle(intl: IntlShape, f: Pick<SavedFilter, "mode" | "name">) {
  return `${intl.formatMessage({ id: FilterModeToMessageID[f.mode] })}: ${
    f.name
  }`;
}

const AddSavedFilterModal: React.FC<IAddSavedFilterModalProps> = ({
  onClose,
  existingSavedFilterIDs,
  candidates,
}) => {
  const intl = useIntl();

  const [value, setValue] = useState("");

  function renderSelect() {
    const options = [
      {
        value: "",
        text: "",
      },
    ].concat(
      candidates.findSavedFilters
        .filter((f) => {
          // markers not currently supported
          return (
            f.mode !== FilterMode.SceneMarkers &&
            !existingSavedFilterIDs.includes(f.id)
          );
        })
        .map((f) => {
          return {
            value: f.id,
            text: filterTitle(intl, f),
          };
        })
    );

    return (
      <Form.Group controlId="filter">
        <Form.Label>
          <FormattedMessage id="search_filter.name" />
        </Form.Label>
        <Form.Control
          as="select"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          className="btn-secondary"
        >
          {options.map((c) => (
            <option key={c.value} value={c.value}>
              {c.text}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  return (
    <Modal show onHide={() => onClose()}>
      <Modal.Header>
        <FormattedMessage id="actions.add" />
      </Modal.Header>
      <Modal.Body>
        <div className="dialog-content">{renderSelect()}</div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={() => onClose()}>
          <FormattedMessage id="actions.cancel" />
        </Button>
        <Button onClick={() => onClose(value)} disabled={value === ""}>
          <FormattedMessage id="actions.add" />
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

interface IFilterRowProps {
  filterID: string;
  header: String;
  index: number;
  onDelete: () => void;
}

const FilterRow: React.FC<IFilterRowProps> = (props: IFilterRowProps) => {
  const intl = useIntl();

  return (
    <div className="recommendation-row">
      <div className="recommendation-row-head">
        <div>
          <h2>{props.header}</h2>
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
  onClose: (filterIDs?: string[]) => void;
}

export const FrontPageConfig: React.FC<IFrontPageConfigProps> = ({
  onClose,
}) => {
  const intl = useIntl();

  const { loading, data } = useFindFrontPageFiltersQuery();
  const { data: foo, loading: loading2 } = useFindSavedFilters();

  const [isAdd, setIsAdd] = useState(false);
  const [currentFilters, setCurrentFilters] = useState<
    Pick<SavedFilter, "id" | "mode" | "name">[]
  >([]);
  const [dragIndex, setDragIndex] = useState<number | undefined>();

  useEffect(() => {
    if (data && data.findFrontPageFilters) {
      setCurrentFilters(data?.findFrontPageFilters);
    }
  }, [data]);

  function onDragStart(event: React.DragEvent<HTMLElement>, index: number) {
    event.dataTransfer.effectAllowed = "move";
    setDragIndex(index);
  }

  function onDragOver(event: React.DragEvent<HTMLElement>, index?: number) {
    if (dragIndex !== undefined && index !== undefined && index !== dragIndex) {
      const newFilters = [...currentFilters];
      const moved = newFilters.splice(dragIndex, 1);
      newFilters.splice(index, 0, moved[0]);
      setCurrentFilters(newFilters);
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

  if (loading || loading2 || !data) {
    return <LoadingIndicator />;
  }

  const existingSavedFilterIDs = currentFilters.map((f) => f.id);

  function addSavedFilter(id?: string) {
    setIsAdd(false);

    if (!id) {
      return;
    }

    const newFilter = foo?.findSavedFilters.find((f) => f.id === id);
    if (newFilter) {
      setCurrentFilters([...currentFilters, newFilter]);
    }
  }

  function deleteSavedFilter(index: number) {
    setCurrentFilters(currentFilters.filter((f, i) => i !== index));
  }

  return (
    <>
      {isAdd && foo && (
        <AddSavedFilterModal
          candidates={foo}
          existingSavedFilterIDs={existingSavedFilterIDs}
          onClose={addSavedFilter}
        />
      )}
      <div className="recommendations-container recommendations-container-edit">
        <div onDragOver={onDragOverDefault}>
          {currentFilters.map((filter, index) => (
            <div
              key={index}
              draggable
              onDragStart={(e) => onDragStart(e, index)}
              onDragEnter={(e) => onDragOver(e, index)}
              onDrop={() => onDrop()}
            >
              <FilterRow
                key={index}
                filterID={filter.id}
                header={filterTitle(intl, filter)}
                index={index}
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
          <Button onClick={() => onClose(currentFilters.map((f) => f.id))}>
            <FormattedMessage id={"actions.save"} />
          </Button>
        </div>
      </div>
    </>
  );
};

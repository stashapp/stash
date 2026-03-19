import { Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { StashBox } from "src/core/generated-graphql";

interface IStashBoxSelectorProps {
  stashBoxes: StashBox[];
  selectedEndpoint: string;
  onEndpointChange: (endpoint: string) => void;
}

export const StashBoxSelector: React.FC<IStashBoxSelectorProps> = ({
  stashBoxes,
  selectedEndpoint,
  onEndpointChange,
}) => {
  return (
    <Form.Control
      as="select"
      value={selectedEndpoint}
      className="input-control"
      disabled={stashBoxes.length < 2}
      onChange={(e) => onEndpointChange(e.target.value)}
    >
      {!stashBoxes.length && (
        <option>
          <FormattedMessage id="tagger.config.no_instances_found" />
        </option>
      )}
      {stashBoxes.map((i) => (
        <option value={i.endpoint} key={i.endpoint}>
          {i.endpoint}
        </option>
      ))}
    </Form.Control>
  );
};

export const StashBoxSelectorField: React.FC<IStashBoxSelectorProps> = (
  props
) => {
  return (
    <Form.Group controlId="scraper">
      <Form.Label>
        <FormattedMessage id="component_tagger.config.source" />
      </Form.Label>
      <div>
        <StashBoxSelector {...props} />
      </div>
    </Form.Group>
  );
};

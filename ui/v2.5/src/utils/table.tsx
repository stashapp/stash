import EditableTextUtils from "./editabletext";

const renderEditableTextTableRow = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>{EditableTextUtils.renderEditableText(options)}</td>
  </tr>
);

const renderTextArea = (options: {
  title: string;
  value: string | undefined;
  isEditing: boolean;
  onChange: (value: string) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>{EditableTextUtils.renderTextArea(options)}</td>
  </tr>
);

const renderInputGroup = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  url?: string;
  onChange: (value: string) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>{EditableTextUtils.renderInputGroup(options)}</td>
  </tr>
);

const renderDurationInput = (options: {
  title: string;
  placeholder?: string;
  value: string | undefined;
  isEditing: boolean;
  asString?: boolean;
  onChange: (value: string | undefined) => void;
}) => {
  return (
    <tr>
      <td>{options.title}</td>
      <td>{EditableTextUtils.renderDurationInput(options)}</td>
    </tr>
  );
};

const renderHtmlSelect = (options: {
  title: string;
  value?: string | number;
  isEditing: boolean;
  onChange: (value: string) => void;
  selectOptions: Array<string | number>;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>{EditableTextUtils.renderHtmlSelect(options)}</td>
  </tr>
);

// TODO: isediting
const renderFilterSelect = (options: {
  title: string;
  type: "performers" | "studios" | "tags";
  initialId: string | undefined;
  onChange: (id: string | undefined) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>{EditableTextUtils.renderFilterSelect(options)}</td>
  </tr>
);

// TODO: isediting
const renderMultiSelect = (options: {
  title: string;
  type: "performers" | "studios" | "tags";
  initialIds: string[] | undefined;
  onChange: (ids: string[]) => void;
}) => (
  <tr>
    <td>{options.title}</td>
    <td>{EditableTextUtils.renderMultiSelect(options)}</td>
  </tr>
);

const TableUtils = {
  renderEditableTextTableRow,
  renderTextArea,
  renderInputGroup,
  renderDurationInput,
  renderHtmlSelect,
  renderFilterSelect,
  renderMultiSelect,
};

export default TableUtils;

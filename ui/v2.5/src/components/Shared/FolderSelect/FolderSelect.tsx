import React, { useState, useEffect } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, InputGroup, Form, Collapse } from "react-bootstrap";
import { Icon } from "../Icon";
import { LoadingIndicator } from "../LoadingIndicator";
import { faEllipsis, faTimes } from "@fortawesome/free-solid-svg-icons";
import { useDebounce } from "src/hooks/debounce";
import TextUtils from "src/utils/text";
import { useDirectoryPaths } from "./useDirectoryPaths";
import { PatchComponent } from "src/patch";

interface IProps {
  currentDirectory: string;
  onChangeDirectory: (value: string) => void;
  defaultDirectories?: string[];
  appendButton?: JSX.Element;
  collapsible?: boolean;
  quotePath?: boolean;
  hideError?: boolean;
}

const _FolderSelect: React.FC<IProps> = ({
  currentDirectory,
  onChangeDirectory,
  defaultDirectories = [],
  appendButton,
  collapsible = false,
  quotePath = false,
  hideError = false,
}) => {
  const intl = useIntl();
  const [showBrowser, setShowBrowser] = useState(false);
  const [path, setPath] = useState(currentDirectory);
  const [displayPath, setDisplayPath] = useState(currentDirectory);

  useEffect(() => {
    setDisplayPath(currentDirectory);
  }, [currentDirectory]);

  const normalizedPath = TextUtils.stripQuotes(path);
  const { directories, parent, error, loading } = useDirectoryPaths(
    normalizedPath,
    hideError
  );

  const selectableDirectories =
    (currentDirectory ? directories : defaultDirectories) ?? defaultDirectories;

  const debouncedSetDirectory = useDebounce(setPath, 250);

  function setInstant(value: string) {
    const stripped = TextUtils.stripQuotes(value);
    const normalizedValue =
      quotePath && stripped.includes(" ")
        ? TextUtils.addQuotes(stripped)
        : stripped;
    setDisplayPath(normalizedValue);
    onChangeDirectory(normalizedValue);
    setPath(normalizedValue);
  }

  function setDebounced(value: string) {
    const stripped = TextUtils.stripQuotes(value);
    setDisplayPath(stripped);
    onChangeDirectory(stripped);
    debouncedSetDirectory(stripped);
  }

  function goUp() {
    if (defaultDirectories?.includes(currentDirectory)) {
      setInstant("");
    } else if (parent) {
      setInstant(parent);
    }
  }

  const topDirectory = currentDirectory && parent && (
    <li className="folder-list-parent folder-list-item">
      <Button variant="link" onClick={() => goUp()} disabled={loading}>
        <span>
          <FormattedMessage id="setup.folder.up_dir" />
        </span>
      </Button>
    </li>
  );

  return (
    <>
      <InputGroup>
        <Form.Control
          className="btn-secondary"
          placeholder={intl.formatMessage({ id: "setup.folder.file_path" })}
          onChange={(e) => {
            setDebounced(e.currentTarget.value);
          }}
          value={displayPath}
          spellCheck={false}
        />

        {appendButton && <InputGroup.Append>{appendButton}</InputGroup.Append>}

        {collapsible && (
          <InputGroup.Append>
            <Button
              variant="secondary"
              onClick={() => setShowBrowser(!showBrowser)}
            >
              <Icon icon={faEllipsis} />
            </Button>
          </InputGroup.Append>
        )}

        {(loading || error) && (
          <InputGroup.Append className="align-self-center">
            {loading ? (
              <LoadingIndicator inline small message="" />
            ) : (
              !hideError && <Icon icon={faTimes} color="red" className="ml-3" />
            )}
          </InputGroup.Append>
        )}
      </InputGroup>

      {!hideError && error !== undefined && (
        <h5 className="mt-4 text-break">Error: {error.message}</h5>
      )}

      <Collapse in={!collapsible || showBrowser}>
        <ul className="folder-list">
          {topDirectory}
          {selectableDirectories.map((dir) => (
            <li key={dir} className="folder-list-item">
              <Button
                variant="link"
                onClick={() => setInstant(dir)}
                disabled={loading}
              >
                <span>{dir}</span>
              </Button>
            </li>
          ))}
        </ul>
      </Collapse>
    </>
  );
};

export const FolderSelect = PatchComponent("FolderSelect", _FolderSelect);

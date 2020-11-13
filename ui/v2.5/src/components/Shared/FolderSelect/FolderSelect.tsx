import React, { useEffect } from "react";
import { FormattedMessage } from "react-intl";
import { Button, InputGroup, Form } from "react-bootstrap";
import { LoadingIndicator } from "src/components/Shared";
import { useDirectory } from "src/core/StashService";

interface IProps {
  currentDirectory: string;
  setCurrentDirectory: (value: string) => void;
  defaultDirectories?: string[];
  appendButton?: JSX.Element;
}

export const FolderSelect: React.FC<IProps> = ({
  currentDirectory,
  setCurrentDirectory,
  defaultDirectories,
  appendButton
}) => {
  const { data, error, loading } = useDirectory(currentDirectory);

  const selectableDirectories: string[] = currentDirectory ? data?.directory.directories ?? defaultDirectories ?? [] : defaultDirectories ?? [];

  useEffect(() => {
    if (currentDirectory === "" && !defaultDirectories && data?.directory.path)
      setCurrentDirectory(data.directory.path);
  }, [currentDirectory, setCurrentDirectory, data, defaultDirectories]);

  function goUp() {
    if (defaultDirectories?.includes(currentDirectory)) {
      setCurrentDirectory("");
    } else if (data?.directory.parent) {
      setCurrentDirectory(data.directory.parent);
    }
  }

  const topDirectory = currentDirectory && data?.directory?.parent ? (
    <li className="folder-list-parent folder-list-item">
      <Button
        variant="link"
        onClick={() => goUp()}
      >
        <FormattedMessage defaultMessage="Up a directory" id="up-dir" />
      </Button>
    </li>
  ) : null;

  return (
    <>
    {error ? <h1>{error.message}</h1> : ""}
    <InputGroup>
      <Form.Control
        placeholder="File path"
        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
          setCurrentDirectory(e.currentTarget.value)
        }
        value={currentDirectory}
        spellCheck={false}
      />
      {appendButton ? (
        <InputGroup.Append>
          {appendButton}
        </InputGroup.Append>
      ) : undefined}
      {!data || !data.directory || loading ? (
        <InputGroup.Append>
          <LoadingIndicator 
            inline 
            small
            message={""}
          />
        </InputGroup.Append>
      ) : undefined}
    </InputGroup>
    <ul className="folder-list">
      {topDirectory}
      {selectableDirectories.map((path) => {
        return (
          <li key={path} className="folder-list-item">
            <Button
              variant="link"
              onClick={() => setCurrentDirectory(path)}
            >
              {path}
            </Button>
          </li>
        );
      })}
    </ul>
    </>
  );
};

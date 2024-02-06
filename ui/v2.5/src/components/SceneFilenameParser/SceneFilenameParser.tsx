/* eslint-disable no-param-reassign, jsx-a11y/click-events-have-key-events, jsx-a11y/no-static-element-interactions */

import React, { useEffect, useState, useCallback, useRef } from "react";
import { Button, Card, Form, Table } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import clone from "lodash-es/clone";
import {
  queryParseSceneFilenames,
  useScenesUpdate,
} from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { Pagination } from "src/components/List/Pagination";
import { IParserInput, ParserInput } from "./ParserInput";
import { ParserField } from "./ParserField";
import { SceneParserResult, SceneParserRow } from "./SceneParserRow";

const initialParserInput = {
  pattern: "{title}.{ext}",
  ignoreWords: [],
  whitespaceCharacters: "._",
  capitalizeTitle: true,
  page: 1,
  pageSize: 20,
  findClicked: false,
  ignoreOrganized: true,
};

const initialShowFieldsState = new Map<string, boolean>([
  ["Title", true],
  ["Date", true],
  ["Rating", true],
  ["Performers", true],
  ["Tags", true],
  ["Studio", true],
]);

export const SceneFilenameParser: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const [parserResult, setParserResult] = useState<SceneParserResult[]>([]);
  const [parserInput, setParserInput] =
    useState<IParserInput>(initialParserInput);
  const prevParserInputRef = useRef<IParserInput>();
  const prevParserInput = prevParserInputRef.current;

  const [allTitleSet, setAllTitleSet] = useState<boolean>(false);
  const [allDateSet, setAllDateSet] = useState<boolean>(false);
  const [allRatingSet, setAllRatingSet] = useState<boolean>(false);
  const [allPerformerSet, setAllPerformerSet] = useState<boolean>(false);
  const [allTagSet, setAllTagSet] = useState<boolean>(false);
  const [allStudioSet, setAllStudioSet] = useState<boolean>(false);

  const [showFields, setShowFields] = useState<Map<string, boolean>>(
    initialShowFieldsState
  );

  const [totalItems, setTotalItems] = useState<number>(0);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [updateScenes] = useScenesUpdate(getScenesUpdateData());

  useEffect(() => {
    prevParserInputRef.current = parserInput;
  }, [parserInput]);

  const determineFieldsToHide = useCallback(() => {
    const { pattern } = parserInput;
    const titleSet = pattern.includes("{title}");
    const dateSet =
      pattern.includes("{date}") ||
      pattern.includes("{dd}") || // don't worry about other partial date fields since this should be implied
      ParserField.fullDateFields.some((f) => {
        return pattern.includes(`{${f.field}}`);
      });
    const ratingSet = pattern.includes("{rating100}");
    const performerSet = pattern.includes("{performer}");
    const tagSet = pattern.includes("{tag}");
    const studioSet = pattern.includes("{studio}");

    const newShowFields = new Map<string, boolean>([
      ["Title", titleSet],
      ["Date", dateSet],
      ["Rating", ratingSet],
      ["Performers", performerSet],
      ["Tags", tagSet],
      ["Studio", studioSet],
    ]);

    setShowFields(newShowFields);
  }, [parserInput]);

  const parseResults = useCallback(
    (
      results: GQL.ParseSceneFilenamesQuery["parseSceneFilenames"]["results"]
    ) => {
      if (results) {
        const result = results
          .map((r) => {
            return new SceneParserResult(r);
          })
          .filter((r) => !!r) as SceneParserResult[];

        setParserResult(result);
        determineFieldsToHide();
      }
    },
    [determineFieldsToHide]
  );

  const parseSceneFilenames = useCallback(() => {
    setParserResult([]);
    setIsLoading(true);

    const parserFilter = {
      q: parserInput.pattern,
      page: parserInput.page,
      per_page: parserInput.pageSize,
      sort: "path",
      direction: GQL.SortDirectionEnum.Asc,
    };

    const parserInputData = {
      ignoreWords: parserInput.ignoreWords,
      whitespaceCharacters: parserInput.whitespaceCharacters,
      capitalizeTitle: parserInput.capitalizeTitle,
      ignoreOrganized: parserInput.ignoreOrganized,
    };

    queryParseSceneFilenames(parserFilter, parserInputData)
      .then((response) => {
        const result = response?.data?.parseSceneFilenames;
        if (result) {
          parseResults(result.results);
          setTotalItems(result.count);
        }
      })
      .catch((err) => Toast.error(err))
      .finally(() => setIsLoading(false));
  }, [parserInput, parseResults, Toast]);

  useEffect(() => {
    // only refresh if parserInput actually changed
    if (prevParserInput === parserInput) {
      return;
    }

    if (parserInput.findClicked) {
      parseSceneFilenames();
    }
  }, [parserInput, parseSceneFilenames, prevParserInput]);

  function onPageSizeChanged(newSize: number) {
    const newInput = clone(parserInput);
    newInput.page = 1;
    newInput.pageSize = newSize;
    setParserInput(newInput);
  }

  function onPageChanged(newPage: number) {
    if (newPage !== parserInput.page) {
      const newInput = clone(parserInput);
      newInput.page = newPage;
      setParserInput(newInput);
    }
  }

  function onFindClicked(input: IParserInput) {
    const newInput = clone(input);
    newInput.page = 1;
    newInput.findClicked = true;
    setParserInput(newInput);
    setTotalItems(0);
  }

  function getScenesUpdateData() {
    return parserResult
      .filter((result) => result.isChanged())
      .map((result) => result.toSceneUpdateInput());
  }

  async function onApply() {
    setIsLoading(true);

    try {
      await updateScenes();
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "scenes" }).toLocaleLowerCase() }
        )
      );
    } catch (e) {
      Toast.error(e);
    }

    setIsLoading(false);

    // trigger a refresh of the results
    onFindClicked(parserInput);
  }

  useEffect(() => {
    const newAllTitleSet = !parserResult.some((r) => {
      return !r.title.isSet;
    });
    const newAllDateSet = !parserResult.some((r) => {
      return !r.date.isSet;
    });
    const newAllRatingSet = !parserResult.some((r) => {
      return !r.rating.isSet;
    });
    const newAllPerformerSet = !parserResult.some((r) => {
      return !r.performers.isSet;
    });
    const newAllTagSet = !parserResult.some((r) => {
      return !r.tags.isSet;
    });
    const newAllStudioSet = !parserResult.some((r) => {
      return !r.studio.isSet;
    });

    setAllTitleSet(newAllTitleSet);
    setAllDateSet(newAllDateSet);
    setAllRatingSet(newAllRatingSet);
    setAllTagSet(newAllPerformerSet);
    setAllTagSet(newAllTagSet);
    setAllStudioSet(newAllStudioSet);
  }, [parserResult]);

  function onSelectAllTitleSet(selected: boolean) {
    const newResult = [...parserResult];

    newResult.forEach((r) => {
      r.title.isSet = selected;
    });

    setParserResult(newResult);
    setAllTitleSet(selected);
  }

  function onSelectAllDateSet(selected: boolean) {
    const newResult = [...parserResult];

    newResult.forEach((r) => {
      r.date.isSet = selected;
    });

    setParserResult(newResult);
    setAllDateSet(selected);
  }

  function onSelectAllRatingSet(selected: boolean) {
    const newResult = [...parserResult];

    newResult.forEach((r) => {
      r.rating.isSet = selected;
    });

    setParserResult(newResult);
    setAllRatingSet(selected);
  }

  function onSelectAllPerformerSet(selected: boolean) {
    const newResult = [...parserResult];

    newResult.forEach((r) => {
      r.performers.isSet = selected;
    });

    setParserResult(newResult);
    setAllPerformerSet(selected);
  }

  function onSelectAllTagSet(selected: boolean) {
    const newResult = [...parserResult];

    newResult.forEach((r) => {
      r.tags.isSet = selected;
    });

    setParserResult(newResult);
    setAllTagSet(selected);
  }

  function onSelectAllStudioSet(selected: boolean) {
    const newResult = [...parserResult];

    newResult.forEach((r) => {
      r.studio.isSet = selected;
    });

    setParserResult(newResult);
    setAllStudioSet(selected);
  }

  function onChange(scene: SceneParserResult, changedScene: SceneParserResult) {
    const newResult = [...parserResult];

    const index = newResult.indexOf(scene);
    newResult[index] = changedScene;

    setParserResult(newResult);
  }

  function renderHeader(
    fieldName: string,
    allSet: boolean,
    onAllSet: (set: boolean) => void
  ) {
    if (!showFields.get(fieldName)) {
      return null;
    }

    return (
      <>
        <th className="w-15">
          <Form.Check
            checked={allSet}
            onChange={() => {
              onAllSet(!allSet);
            }}
          />
        </th>
        <th>{fieldName}</th>
      </>
    );
  }

  function renderTable() {
    if (parserResult.length === 0) {
      return undefined;
    }

    return (
      <>
        <div className="scene-parser-results">
          <Table>
            <thead>
              <tr className="scene-parser-row">
                <th className="parser-field-filename">
                  {intl.formatMessage({
                    id: "config.tools.scene_filename_parser.filename",
                  })}
                </th>
                {renderHeader(
                  intl.formatMessage({ id: "title" }),
                  allTitleSet,
                  onSelectAllTitleSet
                )}
                {renderHeader(
                  intl.formatMessage({ id: "date" }),
                  allDateSet,
                  onSelectAllDateSet
                )}
                {renderHeader(
                  intl.formatMessage({ id: "rating" }),
                  allRatingSet,
                  onSelectAllRatingSet
                )}
                {renderHeader(
                  intl.formatMessage({ id: "performers" }),
                  allPerformerSet,
                  onSelectAllPerformerSet
                )}
                {renderHeader(
                  intl.formatMessage({ id: "tags" }),
                  allTagSet,
                  onSelectAllTagSet
                )}
                {renderHeader(
                  intl.formatMessage({ id: "studio" }),
                  allStudioSet,
                  onSelectAllStudioSet
                )}
              </tr>
            </thead>
            <tbody>
              {parserResult.map((scene) => (
                <SceneParserRow
                  scene={scene}
                  key={scene.id}
                  onChange={(changedScene) => onChange(scene, changedScene)}
                  showFields={showFields}
                />
              ))}
            </tbody>
          </Table>
        </div>
        <Pagination
          currentPage={parserInput.page}
          itemsPerPage={parserInput.pageSize}
          totalItems={totalItems}
          metadataByline={[]}
          onChangePage={(page) => onPageChanged(page)}
        />
        <Button variant="primary" onClick={onApply}>
          <FormattedMessage id="actions.apply" />
        </Button>
      </>
    );
  }

  return (
    <Card id="parser-container" className="col col-sm-9 mx-auto">
      <h4>
        {intl.formatMessage({ id: "config.tools.scene_filename_parser.title" })}
      </h4>
      <ParserInput
        input={parserInput}
        onFind={(input) => onFindClicked(input)}
        onPageSizeChanged={onPageSizeChanged}
        showFields={showFields}
        setShowFields={setShowFields}
      />

      {isLoading && <LoadingIndicator />}
      {renderTable()}
    </Card>
  );
};

export default SceneFilenameParser;

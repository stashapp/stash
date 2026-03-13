import React, { useMemo, useState } from "react";
import {
  Button,
  Form,
  Spinner,
  Table,
  Row,
  Col,
  Card,
  ButtonGroup,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useFindDuplicateImagesQuery } from "src/core/generated-graphql";
import * as GQL from "src/core/generated-graphql";
import { PatchContainerComponent } from "src/patch";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ErrorMessage } from "../Shared/ErrorMessage";
import { FileSize } from "../Shared/FileSize";
import { Pagination } from "src/components/List/Pagination";
import { useHistory } from "react-router-dom";
import { DeleteImagesDialog } from "../Images/DeleteImagesDialog";
import { EditImagesDialog } from "../Images/EditImagesDialog";
import { Icon } from "../Shared/Icon";
import { faPencilAlt, faTrash } from "@fortawesome/free-solid-svg-icons";

const ImageDuplicateCheckerSection = PatchContainerComponent(
  "ImageDuplicateCheckerSection"
);

const ImageDuplicateChecker: React.FC = () => {
  const intl = useIntl();
  const history = useHistory();
  const query = new URLSearchParams(history.location.search);
  const currentPage = Number.parseInt(query.get("page") ?? "1", 10);
  const pageSize = Number.parseInt(query.get("size") ?? "20", 10);
  const hashDistance = Number.parseInt(query.get("distance") ?? "0", 10);

  const [isSearching, setIsSearching] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [checkedImages, setCheckedImages] = useState<Record<string, boolean>>({});
  const [selectedImages, setSelectedImages] = useState<GQL.ImageDataFragment[]>();
  const [deletingImages, setDeletingImages] = useState(false);
  const [editingImages, setEditingImages] = useState(false);

  const { data, loading, error, refetch } = useFindDuplicateImagesQuery({
    variables: { distance: hashDistance },
    skip: !hasSearched,
    fetchPolicy: "network-only",
  });

  const handleSearch = () => {
    setIsSearching(true);
    setHasSearched(true);
    setCheckedImages({});
    refetch({ distance: hashDistance }).finally(() => setIsSearching(false));
  };

  const getGroupTotalSize = (group: GQL.ImageDataFragment[]) => {
    return group.reduce((groupTotal, img) => {
      const imgTotal = img.visual_files.reduce(
        (fileTotal, file) => fileTotal + (file.size ?? 0),
        0
      );
      return groupTotal + imgTotal;
    }, 0);
  };

  const allGroups = useMemo(() => {
    const groups = data?.findDuplicateImages ?? [];
    return [...groups].sort((a, b) => {
      return getGroupTotalSize(b) - getGroupTotalSize(a);
    });
  }, [data?.findDuplicateImages]);

  const pagedGroups = useMemo(() => {
    const start = (currentPage - 1) * pageSize;
    return allGroups.slice(start, start + pageSize);
  }, [allGroups, currentPage, pageSize]);

  const checkCount = Object.keys(checkedImages).filter((id) => checkedImages[id]).length;

  const handleCheck = (checked: boolean, imageID: string) => {
    setCheckedImages({ ...checkedImages, [imageID]: checked });
  };

  const handleDeleteChecked = () => {
    setSelectedImages(allGroups.flat().filter((i) => checkedImages[i.id]));
    setDeletingImages(true);
  };

  const onEdit = () => {
    setSelectedImages(allGroups.flat().filter((i) => checkedImages[i.id]));
    setEditingImages(true);
    setCheckedImages({});
  };

  const onDeleteDialogClosed = (confirmed: boolean) => {
    setDeletingImages(false);
    setSelectedImages(undefined);
    if (confirmed) {
      setCheckedImages({});
      refetch();
    }
  };

  const onEditDialogClosed = (applied: boolean) => {
    setEditingImages(false);
    setSelectedImages(undefined);
    if (applied) {
      refetch();
    }
  };

  if (error) return <ErrorMessage error={error.message} />;

  const renderGroup = (group: GQL.ImageDataFragment[], index: number) => {
    const groupIndex = (currentPage - 1) * pageSize + index + 1;
    return (
      <Card key={groupIndex} className="mb-4">
        <Card.Header className="d-flex justify-content-between align-items-center">
          <h5>Group {groupIndex}</h5>
          <span className="text-muted">
            Total Size: <FileSize size={getGroupTotalSize(group)} />
          </span>
        </Card.Header>
        <Card.Body>
          <Table striped bordered hover responsive size="sm">
            <thead>
              <tr>
                <th style={{ width: "40px" }}></th>
                <th style={{ width: "150px" }}>Image</th>
                <th>Details</th>
                <th style={{ width: "120px" }}>Size</th>
                <th style={{ width: "150px" }}>Dimensions</th>
              </tr>
            </thead>
            <tbody>
              {group.map((img) => {
                const file = img.visual_files[0];
                return (
                  <tr key={img.id}>
                    <td className="text-center align-middle">
                      <Form.Check
                        checked={checkedImages[img.id] || false}
                        onChange={(e) => handleCheck(e.currentTarget.checked, img.id)}
                      />
                    </td>
                    <td>
                      <img
                        src={img.paths.thumbnail || ""}
                        alt={img.title || img.id}
                        style={{
                          maxWidth: "120px",
                          maxHeight: "120px",
                          objectFit: "contain",
                        }}
                      />
                    </td>
                    <td>
                      <div className="fw-bold">{img.title || "(No Title)"}</div>
                      <div className="text-muted small text-truncate" style={{ maxWidth: "400px" }}>
                        {img.visual_files[0]?.path}
                      </div>
                      <div className="mt-1 small">ID: {img.id}</div>
                    </td>
                    <td>
                      <FileSize size={file?.size ?? 0} />
                    </td>
                    <td>
                      {file?.__typename === "ImageFile" || file?.__typename === "VideoFile" ? (
                        <>
                          {file.width} x {file.height}
                        </>
                      ) : (
                        "N/A"
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </Table>
        </Card.Body>
      </Card>
    );
  };

  return (
    <div className="container-fluid py-4">
      <ImageDuplicateCheckerSection>
        {deletingImages && selectedImages && (
          <DeleteImagesDialog
            selected={selectedImages}
            onClose={onDeleteDialogClosed}
          />
        )}
        {editingImages && selectedImages && (
          <EditImagesDialog
            selected={selectedImages}
            onClose={onEditDialogClosed}
          />
        )}
        <Row className="mb-4">
          <Col>
            <h3>
              <FormattedMessage id="config.tools.image_duplicate_checker" />
            </h3>
          </Col>
        </Row>

        <Form className="bg-light p-3 rounded mb-4 shadow-sm">
          <Row className="align-items-end">
            <Col md={3}>
              <Form.Group controlId="distanceInput">
                <Form.Label>PHash Distance</Form.Label>
                <Form.Control
                  type="number"
                  value={hashDistance}
                  min={0}
                  max={10}
                  onChange={(e) => {
                    const val = parseInt(e.target.value) || 0;
                    query.set("distance", val.toString());
                    history.push({ search: query.toString() });
                  }}
                />
                <Form.Text className="text-muted small">
                  0 = exact matches.
                </Form.Text>
              </Form.Group>
            </Col>
            <Col md={2}>
              <Button
                variant="primary"
                className="w-100"
                onClick={handleSearch}
                disabled={isSearching || loading}
              >
                {isSearching || loading ? (
                  <Spinner animation="border" size="sm" />
                ) : (
                  "Search"
                )}
              </Button>
            </Col>
          </Row>
        </Form>

        {loading && <LoadingIndicator />}

        {hasSearched && !loading && !error && allGroups.length === 0 && (
          <div className="text-center py-5 border rounded bg-light">
            <p className="mb-0">No duplicates found with the current distance.</p>
          </div>
        )}

        {hasSearched && !loading && !error && allGroups.length > 0 && (
          <div className="d-flex mb-3 align-items-center">
            <h6 className="me-auto mb-0">
              Found {allGroups.length} duplicate groups
            </h6>
            {checkCount > 0 && (
              <ButtonGroup>
                <OverlayTrigger
                  overlay={
                    <Tooltip id="edit">
                      {intl.formatMessage({ id: "actions.edit" })}
                    </Tooltip>
                  }
                >
                  <Button variant="secondary" onClick={onEdit}>
                    <Icon icon={faPencilAlt} />
                  </Button>
                </OverlayTrigger>
                <OverlayTrigger
                  overlay={
                    <Tooltip id="delete">
                      {intl.formatMessage({ id: "actions.delete" })}
                    </Tooltip>
                  }
                >
                  <Button variant="danger" onClick={handleDeleteChecked}>
                    <Icon icon={faTrash} />
                  </Button>
                </OverlayTrigger>
              </ButtonGroup>
            )}
          </div>
        )}

        {pagedGroups.map((group, index) => renderGroup(group, index))}

        {allGroups.length > pageSize && (
          <div className="d-flex justify-content-center mt-4">
            <Pagination
              currentPage={currentPage}
              totalItems={allGroups.length}
              pageSize={pageSize}
              onChangePage={(page) => {
                query.set("page", page.toString());
                history.push({ search: query.toString() });
              }}
            />
          </div>
        )}
      </ImageDuplicateCheckerSection>
    </div>
  );
};

export default ImageDuplicateChecker;

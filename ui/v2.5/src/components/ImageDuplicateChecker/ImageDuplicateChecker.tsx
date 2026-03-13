import React, { useState } from "react";
import { Button, Form, Spinner } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { useFindDuplicateImagesQuery } from "src/core/generated-graphql";
import { PatchContainerComponent } from "src/patch";

const ImageDuplicateCheckerSection = PatchContainerComponent(
  "ImageDuplicateCheckerSection"
);

const ImageDuplicateChecker: React.FC = () => {
  const [distance, setDistance] = useState(0);
  const [isSearching, setIsSearching] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);

  // We lazily fetch the query only when "Search" is clicked
  const { data, loading, error, refetch } = useFindDuplicateImagesQuery({
    variables: { distance },
    skip: !hasSearched,
    fetchPolicy: "network-only",
  });

  const handleSearch = () => {
    setIsSearching(true);
    setHasSearched(true);
    refetch({ distance }).finally(() => setIsSearching(false));
  };

  const results = data?.findDuplicateImages ?? [];

  return (
    <div className="row image-duplicate-checker">
      <div className="col-md-12">
        <ImageDuplicateCheckerSection>
          <h3>
            <FormattedMessage id="config.tools.image_duplicate_checker" />
          </h3>
          <Form className="d-flex align-items-end mb-4">
            <Form.Group controlId="distanceInput" className="mb-0 me-3">
              <Form.Label>PHash Distance</Form.Label>
              <Form.Control
                type="number"
                value={distance}
                min={0}
                max={10}
                onChange={(e) => setDistance(parseInt(e.target.value) || 0)}
              />
              <Form.Text className="text-muted">
                Distance 0 means exact matches.
              </Form.Text>
            </Form.Group>

            <Button
              variant="primary"
              onClick={handleSearch}
              disabled={isSearching || loading}
            >
              {isSearching || loading ? (
                <Spinner animation="border" size="sm" />
              ) : (
                "Search"
              )}
            </Button>
          </Form>

          {error && (
            <div className="text-danger mb-4">Error: {error.message}</div>
          )}

          {hasSearched && !loading && !error && results.length === 0 && (
            <p>No duplicates found.</p>
          )}

          {results.map((group, index) => {
            if (!group || group.length < 2) return null;
            return (
              <div
                key={index}
                className="duplicate-group mb-4 pb-4 border-bottom"
              >
                <h5>Group {index + 1}</h5>
                {/* ImageList requires an array of items with proper types. We map it nicely. */}
                <div className="d-flex flex-wrap gap-3">
                  {group.map((img) => (
                    <div key={img.id} className="border p-2 rounded">
                      <img
                        src={img.paths.thumbnail || ""}
                        alt={img.title || img.id}
                        style={{
                          maxWidth: "200px",
                          maxHeight: "200px",
                          objectFit: "contain",
                        }}
                      />
                      <div
                        className="mt-2 text-center text-truncate"
                        style={{ maxWidth: "200px" }}
                        title={img.title || img.id}
                      >
                        {img.title || img.id}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </ImageDuplicateCheckerSection>
      </div>
    </div>
  );
};

export default ImageDuplicateChecker;

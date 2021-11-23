import React from "react";
import { Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { useIntl } from "react-intl";

interface IScanOptions {
  options: GQL.ScanMetadataInput;
  setOptions: (s: GQL.ScanMetadataInput) => void;
}

export const ScanOptions: React.FC<IScanOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const intl = useIntl();

  const {
    useFileMetadata,
    stripFileExtension,
    scanGeneratePreviews,
    scanGenerateImagePreviews,
    scanGenerateSprites,
    scanGeneratePhashes,
    scanGenerateThumbnails,
  } = options;

  function setOptions(input: Partial<GQL.ScanMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <Form.Group>
      <Form.Check
        id="use-file-metadata"
        checked={useFileMetadata ?? false}
        label={intl.formatMessage({
          id: "config.tasks.set_name_date_details_from_metadata_if_present",
        })}
        onChange={() => setOptions({ useFileMetadata: !useFileMetadata })}
      />
      <Form.Check
        id="strip-file-extension"
        checked={stripFileExtension ?? false}
        label={intl.formatMessage({
          id: "config.tasks.dont_include_file_extension_as_part_of_the_title",
        })}
        onChange={() => setOptions({ stripFileExtension: !stripFileExtension })}
      />
      <Form.Check
        id="scan-generate-previews"
        checked={scanGeneratePreviews ?? false}
        label={intl.formatMessage({
          id: "config.tasks.generate_video_previews_during_scan",
        })}
        onChange={() =>
          setOptions({ scanGeneratePreviews: !scanGeneratePreviews })
        }
      />
      <div className="d-flex flex-row">
        <div>↳</div>
        <Form.Check
          id="scan-generate-image-previews"
          checked={scanGenerateImagePreviews ?? false}
          disabled={!scanGeneratePreviews}
          label={intl.formatMessage({
            id: "config.tasks.generate_previews_during_scan",
          })}
          onChange={() =>
            setOptions({
              scanGenerateImagePreviews: !scanGenerateImagePreviews,
            })
          }
          className="ml-2 flex-grow"
        />
      </div>
      <Form.Check
        id="scan-generate-sprites"
        checked={scanGenerateSprites ?? false}
        label={intl.formatMessage({
          id: "config.tasks.generate_sprites_during_scan",
        })}
        onChange={() =>
          setOptions({ scanGenerateSprites: !scanGenerateSprites })
        }
      />
      <Form.Check
        id="scan-generate-phashes"
        checked={scanGeneratePhashes ?? false}
        label={intl.formatMessage({
          id: "config.tasks.generate_phashes_during_scan",
        })}
        onChange={() =>
          setOptions({ scanGeneratePhashes: !scanGeneratePhashes })
        }
      />
      <Form.Check
        id="scan-generate-thumbnails"
        checked={scanGenerateThumbnails ?? false}
        label={intl.formatMessage({
          id: "config.tasks.generate_thumbnails_during_scan",
        })}
        onChange={() =>
          setOptions({ scanGenerateThumbnails: !scanGenerateThumbnails })
        }
      />
    </Form.Group>
  );
};

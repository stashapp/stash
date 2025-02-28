import React from "react";
import { useIntl } from "react-intl";
import { Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { NumberField } from "src/utils/form";

export type VideoPreviewSettingsInput = Pick<
  GQL.ConfigGeneralInput,
  | "previewSegments"
  | "previewSegmentDuration"
  | "previewExcludeStart"
  | "previewExcludeEnd"
>;

interface IVideoPreviewInput {
  value: VideoPreviewSettingsInput;
  setValue: (v: VideoPreviewSettingsInput) => void;
}

export const VideoPreviewInput: React.FC<IVideoPreviewInput> = ({
  value,
  setValue,
}) => {
  const intl = useIntl();

  function set(v: Partial<VideoPreviewSettingsInput>) {
    setValue({
      ...value,
      ...v,
    });
  }

  const {
    previewSegments,
    previewSegmentDuration,
    previewExcludeStart,
    previewExcludeEnd,
  } = value;

  return (
    <div>
      <Form.Group id="preview-segments">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_count_head",
          })}
        </h6>
        <NumberField
          className="text-input"
          value={previewSegments?.toString() ?? 1}
          min={1}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({
              previewSegments: Number.parseInt(
                e.currentTarget.value || "1",
                10
              ),
            })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_count_desc",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group id="preview-segment-duration">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_duration_head",
          })}
        </h6>
        <NumberField
          className="text-input"
          value={previewSegmentDuration?.toString() ?? 0}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({
              previewSegmentDuration: Number.parseFloat(
                e.currentTarget.value || "0"
              ),
            })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_duration_desc",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group id="preview-exclude-start">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_start_time_head",
          })}
        </h6>
        <Form.Control
          className="text-input"
          value={previewExcludeStart ?? ""}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({ previewExcludeStart: e.currentTarget.value })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_start_time_desc",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group id="preview-exclude-start">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_end_time_head",
          })}
        </h6>
        <Form.Control
          className="text-input"
          value={previewExcludeEnd ?? ""}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({ previewExcludeEnd: e.currentTarget.value })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_end_time_desc",
          })}
        </Form.Text>
      </Form.Group>
    </div>
  );
};

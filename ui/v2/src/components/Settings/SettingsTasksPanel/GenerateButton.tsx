import {
  Button,
  Checkbox,
  FormGroup,
} from "@blueprintjs/core";
import React, { FunctionComponent, useState } from "react";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { ToastUtils } from "../../../utils/toasts";

interface IProps {}

export const GenerateButton: FunctionComponent<IProps> = () => {
  const [sprites, setSprites] = useState<boolean>(true);
  const [previews, setPreviews] = useState<boolean>(true);
  const [markers, setMarkers] = useState<boolean>(true);
  const [transcodes, setTranscodes] = useState<boolean>(true);

  async function onGenerate() {
    try {
      await StashService.mutateMetadataGenerate({sprites, previews, markers, transcodes});
      ToastUtils.success("Started generating");
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  return (
    <FormGroup
      helperText="Generate supporting image, sprite, video, vtt and other files."
      labelFor="generate"
      inline={true}
    >
      <Checkbox checked={sprites} label="Sprites (for the scene scrubber)" onChange={() => setSprites(!sprites)} />
      <Checkbox
        checked={previews}
        label="Previews (video previews which play when hovering over a scene)"
        onChange={() => setPreviews(!previews)}
      />
      <Checkbox
        checked={markers}
        label="Markers (20 second videos which begin at the given timecode)"
        onChange={() => setMarkers(!markers)}
      />
      <Checkbox
        checked={transcodes}
        label="Transcodes (MP4 conversions of unsupported video formats)"
        onChange={() => setTranscodes(!transcodes)}
      />
      <Button id="generate" text="Generate" onClick={() => onGenerate()} />
    </FormGroup>
  );
};

import React from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Link } from "react-router-dom";
import { Setting } from "./Inputs";
import { SettingSection } from "./SettingSection";
import { PatchContainerComponent } from "src/patch";

const SettingsToolsSection = PatchContainerComponent("SettingsToolsSection");

export const SettingsToolsPanel: React.FC = () => {
  return (
    <SettingSection headingID="config.tools.scene_tools">
      <SettingsToolsSection>
        <Setting
          heading={
            <Link to="/sceneFilenameParser">
              <Button>
                <FormattedMessage id="config.tools.scene_filename_parser.title" />
              </Button>
            </Link>
          }
        />

        <Setting
          heading={
            <Link to="/sceneDuplicateChecker">
              <Button>
                <FormattedMessage id="config.tools.scene_duplicate_checker" />
              </Button>
            </Link>
          }
        />
      </SettingsToolsSection>
    </SettingSection>
  );
};

import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useConfigureUI } from "src/core/StashService";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { Button } from "react-bootstrap";
import { FrontPageConfig } from "./FrontPageConfig";
import { useToast } from "src/hooks/Toast";
import { Control } from "./Control";
import { ConfigurationContext } from "src/hooks/Config";
import {
  FrontPageContent,
  generateDefaultFrontPageContent,
  getFrontPageContent,
} from "src/core/config";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";

const FrontPage: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();

  const [isEditing, setIsEditing] = useState(false);
  const [saving, setSaving] = useState(false);

  const [saveUI] = useConfigureUI();

  const { configuration, loading } = React.useContext(ConfigurationContext);

  useScrollToTopOnMount();

  async function onUpdateConfig(content?: FrontPageContent[]) {
    setIsEditing(false);

    if (!content) {
      return;
    }

    setSaving(true);
    try {
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            frontPageContent: content,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
    setSaving(false);
  }

  if (loading || saving) {
    return <LoadingIndicator />;
  }

  if (isEditing) {
    return <FrontPageConfig onClose={(content) => onUpdateConfig(content)} />;
  }

  const ui = configuration?.ui ?? {};

  if (!ui.frontPageContent) {
    const defaultContent = generateDefaultFrontPageContent(intl);
    onUpdateConfig(defaultContent);
  }

  const frontPageContent = getFrontPageContent(ui);

  return (
    <div className="recommendations-container">
      <div>
        {frontPageContent?.map((content, i) => (
          <Control key={i} content={content} />
        ))}
      </div>
      <div className="recommendations-footer">
        <Button onClick={() => setIsEditing(true)}>
          <FormattedMessage id={"actions.customise"} />
        </Button>
      </div>
    </div>
  );
};

export default FrontPage;

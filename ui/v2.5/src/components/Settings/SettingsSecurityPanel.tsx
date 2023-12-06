import React from "react";
import { ModalSetting, NumberSetting } from "./Inputs";
import { SettingSection } from "./SettingSection";
import * as GQL from "src/core/generated-graphql";
import { Button, Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useSettings } from "./context";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { useGenerateAPIKey } from "src/core/StashService";

type AuthenticationSettingsInput = Pick<
  GQL.ConfigGeneralInput,
  "username" | "password"
>;

interface IAuthenticationInput {
  value: AuthenticationSettingsInput;
  setValue: (v: AuthenticationSettingsInput) => void;
}

const AuthenticationInput: React.FC<IAuthenticationInput> = ({
  value,
  setValue,
}) => {
  const intl = useIntl();

  function set(v: Partial<AuthenticationSettingsInput>) {
    setValue({
      ...value,
      ...v,
    });
  }

  const { username, password } = value;

  return (
    <div>
      <Form.Group id="username">
        <h6>{intl.formatMessage({ id: "config.general.auth.username" })}</h6>
        <Form.Control
          className="text-input"
          value={username ?? ""}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({ username: e.currentTarget.value })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.general.auth.username_desc" })}
        </Form.Text>
      </Form.Group>
      <Form.Group id="password">
        <h6>{intl.formatMessage({ id: "config.general.auth.password" })}</h6>
        <Form.Control
          className="text-input"
          type="password"
          value={password ?? ""}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({ password: e.currentTarget.value })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.general.auth.password_desc" })}
        </Form.Text>
      </Form.Group>
    </div>
  );
};

export const SettingsSecurityPanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();

  const { general, apiKey, loading, error, saveGeneral, refetch } =
    useSettings();

  const [generateAPIKey] = useGenerateAPIKey();

  async function onGenerateAPIKey() {
    try {
      await generateAPIKey({
        variables: {
          input: {},
        },
      });
      refetch();
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onClearAPIKey() {
    try {
      await generateAPIKey({
        variables: {
          input: {
            clear: true,
          },
        },
      });
      refetch();
    } catch (e) {
      Toast.error(e);
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <SettingSection headingID="config.general.auth.authentication">
        <ModalSetting<AuthenticationSettingsInput>
          id="authentication-settings"
          headingID="config.general.auth.credentials.heading"
          subHeadingID="config.general.auth.credentials.description"
          value={{
            username: general.username,
            password: general.password,
          }}
          onChange={(v) => saveGeneral(v)}
          renderField={(value, setValue) => (
            <AuthenticationInput value={value ?? {}} setValue={setValue} />
          )}
          renderValue={(v) => {
            if (v?.username && v?.password)
              return <span>{v?.username ?? ""}</span>;
            return <></>;
          }}
        />

        <div className="setting" id="apikey">
          <div>
            <h3>{intl.formatMessage({ id: "config.general.auth.api_key" })}</h3>

            <div className="value text-break">{apiKey}</div>

            <div className="sub-heading">
              {intl.formatMessage({ id: "config.general.auth.api_key_desc" })}
            </div>
          </div>
          <div>
            <Button
              disabled={!general.username || !general.password}
              onClick={() => onGenerateAPIKey()}
            >
              {intl.formatMessage({
                id: "config.general.auth.generate_api_key",
              })}
            </Button>
            <Button variant="danger" onClick={() => onClearAPIKey()}>
              {intl.formatMessage({
                id: "config.general.auth.clear_api_key",
              })}
            </Button>
          </div>
        </div>

        <NumberSetting
          id="maxSessionAge"
          headingID="config.general.auth.maximum_session_age"
          subHeadingID="config.general.auth.maximum_session_age_desc"
          value={general.maxSessionAge ?? undefined}
          onChange={(v) => saveGeneral({ maxSessionAge: v })}
        />
      </SettingSection>
    </>
  );
};

import { Tab, Nav, Dropdown } from "react-bootstrap";
import React, {
  useCallback,
  useEffect,
  useState,
  useMemo,
  useLayoutEffect,
} from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory, RouteComponentProps } from "react-router-dom";
import { Helmet } from "react-helmet";
import * as GQL from "src/core/generated-graphql";
import {
  mutateMetadataScan,
  useFindAudio,
  useAudioIncrementO,
  useAudioUpdate,
  useAudioIncrementPlayCount,
} from "src/core/StashService";

import { AudioEditPanel } from "./AudioEditPanel";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { Icon } from "src/components/Shared/Icon";
import { Counter } from "src/components/Shared/Counter";
import { useToast } from "src/hooks/Toast";
import Mousetrap from "mousetrap";
import { OrganizedButton } from "src/components/Scenes/SceneDetails/OrganizedButton";
import { useConfigurationContext } from "src/hooks/Config";
import { faEllipsisV } from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectTitle } from "src/core/files";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import TextUtils from "src/utils/text";
import {
  OCounterButton,
  ViewCountButton,
} from "src/components/Shared/CountButton";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { lazyComponent } from "src/utils/lazyComponent";
import cx from "classnames";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PatchComponent, PatchContainerComponent } from "src/patch";
import { AudioMergeModal } from "../AudioMergeDialog";
import { goBackOrReplace } from "src/utils/history";
import { FormattedDate } from "src/components/Shared/Date";
import { StudioLogo } from "src/components/Shared/StudioLogo";

const AudioPlayer = lazyComponent(
  () => import("src/components/AudioPlayer/AudioPlayer")
);

const GalleryViewer = lazyComponent(
  () => import("src/components/Galleries/GalleryViewer")
);

const AudioFileInfoPanel = lazyComponent(() => import("./AudioFileInfoPanel"));
const AudioDetailPanel = lazyComponent(() => import("./AudioDetailPanel"));
const AudioHistoryPanel = lazyComponent(() => import("./AudioHistoryPanel"));
const AudioGroupPanel = lazyComponent(() => import("./AudioGroupPanel"));
const AudioGalleriesPanel = lazyComponent(
  () => import("./AudioGalleriesPanel")
);
const DeleteAudiosDialog = lazyComponent(() => import("../DeleteAudiosDialog"));

const AudioSpecs: React.FC<{
  bitRate?: number | null;
  sampleRate?: number | null;
  audioCodec?: string | null;
}> = ({ bitRate, sampleRate, audioCodec }) => {
  const intl = useIntl();

  const parts: React.ReactNode[] = [];

  if (audioCodec) {
    parts.push(
      <span className="audio-codec" key="codec" data-value={audioCodec}>
        {audioCodec}
      </span>
    );
  }

  if (bitRate) {
    parts.push(
      <span className="bit-rate" key="bitrate" data-value={bitRate}>
        <FormattedMessage
          id="kilobits_per_second"
          values={{
            value: intl.formatNumber(bitRate / 1000, {
              maximumFractionDigits: 0,
            }),
          }}
        />
      </span>
    );
  }

  if (sampleRate) {
    parts.push(
      <span className="sample-rate" key="samplerate" data-value={sampleRate}>
        <FormattedMessage
          id="hertz"
          values={{ value: intl.formatNumber(sampleRate) }}
        />
      </span>
    );
  }

  return (
    <span>
      {parts.map((p, i) => (
        <React.Fragment key={i}>
          {i > 0 && <span className="divider"> | </span>}
          {p}
        </React.Fragment>
      ))}
    </span>
  );
};

interface IProps {
  audio: GQL.AudioDataFragment;
  onDelete: () => void;
}

interface IAudioParams {
  id: string;
}

const AudioPageTabs = PatchContainerComponent<IProps>("AudioPage.Tabs");
const AudioPageTabContent = PatchContainerComponent<IProps>(
  "AudioPage.TabContent"
);

const AudioPage: React.FC<IProps> = PatchComponent("AudioPage", (props) => {
  const { audio, onDelete } = props;

  const Toast = useToast();
  const intl = useIntl();
  const history = useHistory();
  const [updateAudio] = useAudioUpdate();
  const { configuration } = useConfigurationContext();
  const { showStudioText } = configuration?.ui ?? {};

  const [incrementO] = useAudioIncrementO(audio.id);
  const [incrementPlay] = useAudioIncrementPlayCount();

  function incrementPlayCount() {
    incrementPlay({
      variables: {
        id: audio.id,
      },
    });
  }

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const [activeTabKey, setActiveTabKey] = useState("audio-details-panel");

  const [isMerging, setIsMerging] = useState(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  const onIncrementOClick = useCallback(async () => {
    try {
      await incrementO();
    } catch (e) {
      Toast.error(e);
    }
  }, [incrementO, Toast]);

  function setRating(v: number | null) {
    updateAudio({
      variables: {
        input: {
          id: audio.id,
          rating100: v,
        },
      },
    });
  }

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("audio-details-panel"));
    Mousetrap.bind("e", () => setActiveTabKey("audio-edit-panel"));
    Mousetrap.bind("i", () => setActiveTabKey("audio-file-info-panel"));
    Mousetrap.bind("h", () => setActiveTabKey("audio-history-panel"));
    Mousetrap.bind("o", () => {
      onIncrementOClick();
    });
    Mousetrap.bind("d d", () => setIsDeleteAlertOpen(true));

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("i");
      Mousetrap.unbind("h");
      Mousetrap.unbind("o");
      Mousetrap.unbind("d d");
    };
  });

  async function onSave(input: GQL.AudioCreateInput) {
    await updateAudio({
      variables: {
        input: {
          id: audio.id,
          ...input,
        },
      },
    });
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "audio" }).toLocaleLowerCase() }
      )
    );
  }

  const onOrganizedClick = async () => {
    try {
      setOrganizedLoading(true);
      await updateAudio({
        variables: {
          input: {
            id: audio.id,
            organized: !audio.organized,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setOrganizedLoading(false);
    }
  };

  async function onRescan() {
    await mutateMetadataScan({
      paths: [objectPath(audio)],
      rescan: true,
    });

    Toast.success(
      intl.formatMessage(
        { id: "toast.rescanning_entity" },
        {
          count: 1,
          singularEntity: intl
            .formatMessage({ id: "audio" })
            .toLocaleLowerCase(),
        }
      )
    );
  }

  function onDeleteDialogClosed(deleted: boolean) {
    setIsDeleteAlertOpen(false);
    if (deleted) {
      onDelete();
    }
  }

  function maybeRenderMergeDialog() {
    if (!audio.id) return;
    return (
      <AudioMergeModal
        show={isMerging}
        onClose={(mergedId) => {
          setIsMerging(false);
          if (mergedId !== undefined && mergedId !== audio.id) {
            history.replace(`/audios/${mergedId}`);
          }
        }}
        audios={[{ id: audio.id, title: objectTitle(audio) }]}
      />
    );
  }

  function maybeRenderDeleteDialog() {
    if (isDeleteAlertOpen) {
      return (
        <DeleteAudiosDialog selected={[audio]} onClose={onDeleteDialogClosed} />
      );
    }
  }

  const renderOperations = () => (
    <Dropdown>
      <Dropdown.Toggle
        variant="secondary"
        id="operation-menu"
        className="minimal"
        title={intl.formatMessage({ id: "operations" })}
      >
        <Icon icon={faEllipsisV} />
      </Dropdown.Toggle>
      <Dropdown.Menu className="bg-secondary text-white">
        {!!audio.files.length && (
          <Dropdown.Item
            key="rescan"
            className="bg-secondary text-white"
            onClick={() => onRescan()}
          >
            <FormattedMessage id="actions.rescan" />
          </Dropdown.Item>
        )}
        <Dropdown.Item
          key="merge-audio"
          className="bg-secondary text-white"
          onClick={() => setIsMerging(true)}
        >
          <FormattedMessage id="actions.merge" />
          ...
        </Dropdown.Item>
        <Dropdown.Item
          key="delete-audio"
          className="bg-secondary text-white"
          onClick={() => setIsDeleteAlertOpen(true)}
        >
          <FormattedMessage
            id="actions.delete"
            values={{ entityType: intl.formatMessage({ id: "audio" }) }}
          />
        </Dropdown.Item>
      </Dropdown.Menu>
    </Dropdown>
  );

  const renderTabs = () => (
    <Tab.Container
      activeKey={activeTabKey}
      onSelect={(k) => k && setActiveTabKey(k)}
    >
      <div>
        <Nav variant="tabs" className="mr-auto">
          <AudioPageTabs {...props}>
            <Nav.Item>
              <Nav.Link eventKey="audio-details-panel">
                <FormattedMessage id="details" />
              </Nav.Link>
            </Nav.Item>
            {audio.groups.length > 0 ? (
              <Nav.Item>
                <Nav.Link eventKey="audio-group-panel">
                  <FormattedMessage
                    id="countables.groups"
                    values={{ count: audio.groups.length }}
                  />
                </Nav.Link>
              </Nav.Item>
            ) : (
              ""
            )}
            {audio.galleries.length >= 1 ? (
              <Nav.Item>
                <Nav.Link eventKey="audio-galleries-panel">
                  <FormattedMessage
                    id="countables.galleries"
                    values={{ count: audio.galleries.length }}
                  />
                </Nav.Link>
              </Nav.Item>
            ) : undefined}
            <Nav.Item>
              <Nav.Link eventKey="audio-file-info-panel">
                <FormattedMessage id="file_info" />
                <Counter count={audio.files.length} hideZero hideOne />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="audio-history-panel">
                <FormattedMessage id="history" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="audio-edit-panel">
                <FormattedMessage id="actions.edit" />
              </Nav.Link>
            </Nav.Item>
          </AudioPageTabs>
        </Nav>
      </div>

      <Tab.Content>
        <AudioPageTabContent {...props}>
          <Tab.Pane eventKey="audio-details-panel">
            <AudioDetailPanel audio={audio} />
          </Tab.Pane>
          <Tab.Pane eventKey="audio-group-panel">
            <AudioGroupPanel audio={audio} />
          </Tab.Pane>
          {audio.galleries.length >= 1 && (
            <Tab.Pane eventKey="audio-galleries-panel">
              <AudioGalleriesPanel galleries={audio.galleries} />
              {audio.galleries.length === 1 && (
                <GalleryViewer galleryId={audio.galleries[0].id} />
              )}
            </Tab.Pane>
          )}
          <Tab.Pane
            className="file-info-panel"
            eventKey="audio-file-info-panel"
          >
            <AudioFileInfoPanel audio={audio} />
          </Tab.Pane>
          <Tab.Pane eventKey="audio-edit-panel" mountOnEnter>
            <AudioEditPanel
              isVisible={activeTabKey === "audio-edit-panel"}
              audio={audio}
              onSubmit={onSave}
              onDelete={() => setIsDeleteAlertOpen(true)}
            />
          </Tab.Pane>
          <Tab.Pane eventKey="audio-history-panel">
            <AudioHistoryPanel audio={audio} />
          </Tab.Pane>
        </AudioPageTabContent>
      </Tab.Content>
    </Tab.Container>
  );

  const title = objectTitle(audio);

  const file = useMemo(
    () => (audio.files.length > 0 ? audio.files[0] : undefined),
    [audio]
  );

  return (
    <>
      <Helmet>
        <title>{title}</title>
      </Helmet>
      {maybeRenderMergeDialog()}
      {maybeRenderDeleteDialog()}
      <div className="audio-tabs order-xl-first order-last">
        <div>
          <div className="audio-header-container">
            <StudioLogo studio={audio.studio} showText={showStudioText} />
            <h3 className={cx("audio-header", { "no-studio": !audio.studio })}>
              <TruncatedText lineCount={2} text={title} />
            </h3>
          </div>

          <div className="audio-subheader">
            <span className="date" data-value={audio.date}>
              {!!audio.date && <FormattedDate value={audio.date} />}
            </span>
            <span className="duration">
              {!!file?.duration && TextUtils.secondsToTimestamp(file.duration)}
            </span>
            <AudioSpecs
              bitRate={file?.bit_rate}
              sampleRate={file?.sample_rate}
              audioCodec={file?.audio_codec}
            />
          </div>

          <div className="audio-toolbar">
            <span className="audio-toolbar-group">
              <RatingSystem
                value={audio.rating100}
                onSetRating={setRating}
                clickToRate
                withoutContext
              />
            </span>
            <span className="audio-toolbar-group">
              <span>
                <ViewCountButton
                  value={audio.play_count ?? 0}
                  onIncrement={() => incrementPlayCount()}
                />
              </span>
              <span>
                <OCounterButton
                  value={audio.o_counter ?? 0}
                  onIncrement={() => onIncrementOClick()}
                />
              </span>
              <span>
                <OrganizedButton
                  loading={organizedLoading}
                  organized={audio.organized}
                  onClick={onOrganizedClick}
                />
              </span>
              <span>{renderOperations()}</span>
            </span>
          </div>
        </div>
        {renderTabs()}
      </div>
    </>
  );
});

const AudioLoader: React.FC<RouteComponentProps<IAudioParams>> = ({
  location,
  history,
  match,
}) => {
  const { id } = match.params;
  const { data, loading, error } = useFindAudio(id);

  const [audio, setAudio] = useState<GQL.AudioDataFragment>();

  // useLayoutEffect to update before paint
  useLayoutEffect(() => {
    // only update audio when loading is done
    if (!loading) {
      setAudio(data?.findAudio ?? undefined);
    }
  }, [data, loading]);

  const queryParams = useMemo(
    () => new URLSearchParams(location.search),
    [location.search]
  );

  const initialTimestamp = useMemo(() => {
    const t = queryParams.get("t");
    if (!t) return 0;

    const n = Number(t);
    if (Number.isNaN(n)) return 0;
    return n;
  }, [queryParams]);

  const autoplay = queryParams.get("autoplay") === "true";

  function onDelete() {
    goBackOrReplace(history, "/audios");
  }

  if (!audio) {
    if (loading) return <LoadingIndicator />;
    if (error) return <ErrorMessage error={error.message} />;
    return <ErrorMessage error={`No audio found with id ${id}.`} />;
  }

  return (
    <div className="row">
      <AudioPage audio={audio} onDelete={onDelete} />
      <div className="audio-player-container">
        <AudioPlayer
          key="AudioPlayer"
          audio={audio}
          autoplay={autoplay}
          initialTimestamp={initialTimestamp}
        />
      </div>
    </div>
  );
};

export default AudioLoader;

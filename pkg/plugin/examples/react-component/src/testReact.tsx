interface IPluginApi {
  React: typeof React;
  ReactRouterDOM: {
    Link: React.FC<any>;
    Route: React.FC<any>;
  }
  components: Record<string, React.FC<any>>;
  utils: {
    NavUtils: any;
  },
  patch: {
    before: (target: string, fn: Function) => void;
    instead: (target: string, fn: Function) => void;
    after: (target: string, fn: Function) => void;
  },
  register: {
    route: (path: string, component: React.FC<any>) => void;
  }
}

(function () {
  const PluginApi = (window as any).PluginApi as IPluginApi;
  const React = PluginApi.React;

  const {
    HoverPopover,
    TagLink,
  } = PluginApi.components;

  const {
    Link,
  } = PluginApi.ReactRouterDOM;

  const {
    NavUtils
  } = PluginApi.utils;

  const ScenePerformer: React.FC<{
    performer: any;
  }> = ({ performer }) => {
    const popoverContent = React.useMemo(
      () => (
        <div className="scene-performer-popover">
          <Link to={`/performers/${performer.id}`}>
            <img
              className="image-thumbnail"
              alt={performer.name ?? ""}
              src={performer.image_path ?? ""}
            />
          </Link>
        </div>
      ),
      [performer]
    );
  
    return (
      <HoverPopover
        className="scene-card__performer"
        placement="top"
        content={popoverContent}
        leaveDelay={100}
      >
        <a href={NavUtils.makePerformerScenesUrl(performer)}>{performer.name}</a>
      </HoverPopover>
    );
  };

  function SceneDetails(props: any) {
    function maybeRenderPerformers() {
      if (props.scene.performers.length <= 0) return;
  
      return (
        <div className="scene-card__performers">
          {props.scene.performers.map((performer: any) => (
            <ScenePerformer performer={performer} key={performer.id} />
          ))}
        </div>
      );
    }
  
    function maybeRenderTags() {
      if (props.scene.tags.length <= 0) return;
  
      return (
        <div className="scene-card__tags">
          {props.scene.tags.map((tag: any) => (
            <TagLink key={tag.id} tag={tag} />
          ))}
        </div>
      );
    }

    return (
      <div className="scene-card__details">
        <span className="scene-card__date">{props.scene.date}</span>
        {maybeRenderPerformers()}
        {maybeRenderTags()}
      </div>
    );
  }

  PluginApi.patch.instead("SceneCard.Details", function (props: any, _: any, original: any) {
    return <SceneDetails {...props} />;
  });

  const TestPage: React.FC = () => {
    return (
      <div>This is a test page.</div>
    );
  };

  PluginApi.register.route("/plugin/test-react", TestPage);
  
})();
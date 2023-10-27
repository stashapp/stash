interface IPluginApi {
  React: typeof React;
  libraries: {
    ReactRouterDOM: {
      Link: React.FC<any>;
      Route: React.FC<any>;
      NavLink: React.FC<any>;
    },
    Bootstrap: {
      Button: React.FC<any>;
      Nav: React.FC<any> & {
        Link: React.FC<any>;
      };
    },
    FontAwesomeSolid: {
      faEthernet: any;
    },
    Intl: {
      FormattedMessage: React.FC<any>;
    }
  },
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

  const { Button } = PluginApi.libraries.Bootstrap;
  const { faEthernet } = PluginApi.libraries.FontAwesomeSolid;
  const {
    Link,
    NavLink,
  } = PluginApi.libraries.ReactRouterDOM;

  const {
    NavUtils
  } = PluginApi.utils;

  const ScenePerformer: React.FC<{
    performer: any;
  }> = ({ performer }) => {
    // PluginApi.components may not be registered when the outside function is run
    // need to initialise these inside the function component
    const {
      HoverPopover,
    } = PluginApi.components;

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
    const {
      TagLink,
    } = PluginApi.components;

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

  PluginApi.patch.before("SettingsToolsSection", function (props: any) {
    const {
      Setting,
    } = PluginApi.components;

    return [
      {
        children: (
          <>
            {props.children}
            <Setting
              heading={
                <Link to="/plugin/test-react">
                  <Button>
                    Test page
                  </Button>
                </Link>
              }
            />
          </>
        ),
      },
    ];
  });

  PluginApi.patch.before("MainNavBar.UtilityItems", function (props: any) {
    const {
      Icon,
    } = PluginApi.components;

    return [
      {
        children: (
          <>
            {props.children}
            <NavLink
              className="nav-utility"
              exact
              to="/plugin/test-react"
            >
              <Button
                className="minimal d-flex align-items-center h-100"
                title="Test page"
              >
                <Icon icon={faEthernet} />
              </Button>
            </NavLink>
          </>
        )
      }
    ]
  })
})();
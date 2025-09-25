interface IPluginApi {
  React: typeof React;
  GQL: any;
  Event: {
    addEventListener: (event: string, callback: (e: CustomEvent) => void) => void;
  };
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
        Item: React.FC<any>;
      };
      Tab: React.FC<any> & {
        Pane: React.FC<any>;
      }
    },
    FontAwesomeSolid: {
      faEthernet: any;
    },
    Intl: {
      FormattedMessage: React.FC<any>;
    }
  },
  loadableComponents: any;
  components: Record<string, React.FC<any>>;
  utils: {
    NavUtils: any;
    loadComponents: any;
  },
  hooks: any;
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
  const GQL = PluginApi.GQL;

  const { Button, Nav, Tab } = PluginApi.libraries.Bootstrap;
  const { faEthernet } = PluginApi.libraries.FontAwesomeSolid;
  const {
    Link,
    NavLink,
  } = PluginApi.libraries.ReactRouterDOM;

  const {
    NavUtils
  } = PluginApi.utils;

  PluginApi.Event.addEventListener("stash:location", (e) => console.log("Page Changed", e.detail.data.location.pathname, e.detail.data.location.search))

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

  function Overlays() {
    return <span className="example-react-component-custom-overlay">Custom overlay</span>;
  }

  PluginApi.patch.instead("SceneCard.Details", function (props: any, _: any, original: any) {
    return <SceneDetails {...props} />;
  });

  PluginApi.patch.instead("SceneCard.Overlays", function (props: any, _: any, original: (props: any) => any) {  
    return <><Overlays />{original({...props})}</>;
  });

  PluginApi.patch.instead("FrontPage", function (props: any, _: any, original: (props: any) => any) {  
    return <><p>Hello from Test React!</p>{original({...props})}</>;
  });

  const TestPage: React.FC = () => {
    const componentsToLoad = [
      PluginApi.loadableComponents.SceneCard,
      PluginApi.loadableComponents.PerformerSelect,
    ];
    const componentsLoading = PluginApi.hooks.useLoadComponents(componentsToLoad);
    
    const {
      SceneCard,
      LoadingIndicator,
      PerformerSelect,
    } = PluginApi.components;

    // read a random scene and show a scene card for it
    const { data } = GQL.useFindScenesQuery({
      variables: {
        filter: {
          per_page: 1,
          sort: "random",
        },
      },
    });

    const scene = data?.findScenes.scenes[0];

    if (componentsLoading) return (
      <LoadingIndicator />
    );
    
    return (
      <div>
        <div>This is a test page.</div>
        {!!scene && <SceneCard scene={data.findScenes.scenes[0]} />}
        <div>
          <PerformerSelect isMulti onSelect={() => {}} values={[]} />
        </div>
      </div>
    );
  };

  PluginApi.register.route("/plugins/test-react", TestPage);

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
                <Link to="/plugins/test-react">
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
              to="/plugins/test-react"
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
  });

  PluginApi.patch.before("ScenePage.Tabs", function (props: any) {
    return [
      {
        children: (
          <>
            {props.children}
            <Nav.Item>
              <Nav.Link eventKey="test-react-tab">
                Test React tab
              </Nav.Link>
            </Nav.Item>
          </>
        ),
      },
    ];
  });

  PluginApi.patch.before("ScenePage.TabContent", function (props: any) {
    return [
      {
        children: (
          <>
            {props.children}
            <Tab.Pane eventKey="test-react-tab">
              Test React tab content {props.scene.id}
            </Tab.Pane>
          </>
        ),
      },
    ];
  });
})();
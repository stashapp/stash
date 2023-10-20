(function () {
  const PluginApi = window.PluginApi;
  const React = PluginApi.React;
  const e = React.createElement;

  var state = false;
  function MainComponent() {
      return e(
          "button",
          { onClick: function onClick() {
                  state[1](true);
              } },
          state[0] ? "true" : "false"
      );
  }

  function FooPage() {
      return e(
          "button",
          { onClick: function onClick() {
                  state[1](true);
              } },
          state[0] ? "foo true" : "foo false"
      );
  }

  function SceneDetails(props) {
      const performers = React.useMemo(() => {
        return props.scene.performers.map((performer) => 
          e("div", null, performer.name)
        );
      }, [props.scene.performers]);

      return e("div", null, performers);
  }

  PluginApi.register.component("main", MainComponent);
  PluginApi.register.page("foo", FooPage);
  PluginApi.register.cardComponentHook("SceneCard", { Details: SceneDetails });
})();
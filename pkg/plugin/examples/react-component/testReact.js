(function () {
    var state = false;

    function TestReact(props) {
        var state = React.useState(false);

        return React.createElement(
            "button",
            { onClick: function onClick() {
                    state[1](true);
                } },
            state[0] ? "true" : "false"
        );
    }

    addPluginComponent("main", TestReact);
    registerPluginPage("foo", TestReact);
})();
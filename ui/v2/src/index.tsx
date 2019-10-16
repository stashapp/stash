import React from "react";
import { ApolloProvider } from "react-apollo-hooks";
import ReactDOM from "react-dom";
import { BrowserRouter } from "react-router-dom";
import { App } from "./App";
import { StashService } from "./core/StashService";
import "./index.scss";
import * as serviceWorker from "./serviceWorker";

ReactDOM.render((
  <>
  <link rel="stylesheet" type="text/css" href="/css"/>
  <BrowserRouter>
    <ApolloProvider client={StashService.initialize()!}>
      <App />
    </ApolloProvider>
  </BrowserRouter>
  </>
), document.getElementById("root"));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();

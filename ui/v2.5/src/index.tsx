import React from "react";
import { ApolloProvider } from "@apollo/client";
import ReactDOM from "react-dom";
import { BrowserRouter } from "react-router-dom";
import { App } from "./App";
import { getClient } from "./core/StashService";
import { getPlatformURL, getBaseURL } from "./core/createClient";
import "./index.scss";
import * as serviceWorker from "./serviceWorker";

ReactDOM.render(
  <>
    <BrowserRouter basename={getBaseURL()}>
      <ApolloProvider client={getClient()}>
        <App />
      </ApolloProvider>
    </BrowserRouter>
  </>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();

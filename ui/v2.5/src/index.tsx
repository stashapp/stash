import React from "react";
import { ApolloProvider } from "@apollo/client";
import ReactDOM from "react-dom";
import { BrowserRouter } from "react-router-dom";
import { App } from "./App";
import { getClient } from "./core/StashService";
import { getPlatformURL, getBaseURL } from "./core/createClient";
import "./index.scss";

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

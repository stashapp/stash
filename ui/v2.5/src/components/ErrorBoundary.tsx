import React from "react";
import { is_lazy_component_error } from "src/utils/lazy_component";

interface IErrorBoundaryProps {
  children?: React.ReactNode;
}

type ErrorInfo = {
  componentStack: string;
};

interface IErrorBoundaryState {
  error?: Error;
  errorHelp?: string;
  errorInfo?: ErrorInfo;
}

export class ErrorBoundary extends React.Component<
  IErrorBoundaryProps,
  IErrorBoundaryState
> {
  constructor(props: IErrorBoundaryProps) {
    super(props);
    this.state = {};
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    let errorHelp: string | undefined;
    if (is_lazy_component_error(error)) {
      errorHelp =
        "If you recently upgraded Stash, please reload the page or clear your browser cache.";
    }
    this.setState({
      error,
      errorHelp,
      errorInfo,
    });
  }

  public render() {
    const { error, errorHelp, errorInfo } = this.state;
    if (errorInfo) {
      // Error path
      return (
        <div>
          <h2>Something went wrong.</h2>
          {errorHelp && <h5>{errorHelp}</h5>}
          <details className="error-message">
            {error?.toString()}
            <br />
            {errorInfo.componentStack.trim().replaceAll(/^\s*/gm, "    ")}
          </details>
        </div>
      );
    }

    // Normally, just render children
    return this.props.children;
  }
}

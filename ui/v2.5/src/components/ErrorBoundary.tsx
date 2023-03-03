import React from "react";
import { FormattedMessage } from "react-intl";
import { isLazyComponentError } from "src/utils/lazyComponent";

interface IErrorBoundaryProps {
  children?: React.ReactNode;
}

type ErrorInfo = {
  componentStack: string;
};

interface IErrorBoundaryState {
  error?: Error;
  errorHelpId?: string;
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
    let errorHelpId: string | undefined;
    if (isLazyComponentError(error)) {
      errorHelpId = "errors.lazy_component_error_help";
    }
    this.setState({
      error,
      errorHelpId,
      errorInfo,
    });
  }

  public render() {
    const { error, errorHelpId, errorInfo } = this.state;
    if (errorInfo) {
      // Error path
      return (
        <div>
          <h2>
            <FormattedMessage id="errors.something_went_wrong" />
          </h2>
          {errorHelpId && (
            <h5>
              <FormattedMessage id={errorHelpId} />
            </h5>
          )}
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

import { Component, type ReactNode } from "react";
import { AlertTriangle, RefreshCw } from "lucide-react";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false, error: null };

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  render() {
    if (!this.state.hasError) return this.props.children;

    return (
      <div className="flex flex-col items-center justify-center min-h-[50vh] gap-4 p-8 text-center">
        <AlertTriangle className="w-12 h-12 text-amber-400" />
        <h2 className="text-xl font-semibold text-plex-text">Something went wrong</h2>
        <p className="text-sm text-plex-text-muted max-w-md">
          {this.state.error?.message || "An unexpected error occurred."}
        </p>
        <button
          onClick={() => { this.setState({ hasError: false, error: null }); }}
          className="inline-flex items-center gap-2 rounded-xl bg-plex-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-plex-accent-hover"
        >
          <RefreshCw className="w-4 h-4" /> Try Again
        </button>
      </div>
    );
  }
}

import React, { useCallback, useEffect, useRef, useState } from "react";
import { getPlatformURL } from "src/core/createClient";
import "./AssistantWidget.scss";

// ── types ────────────────────────────────────────────────────────────────────

interface IStatus {
  enabled: boolean;
  configured: boolean;
  model: string;
  write_policy: string;
  tools: string[];
}

interface IConfirm {
  id: string;
  name: string;
  input: unknown;
}

type Entry =
  | { type: "user"; text: string }
  | { type: "assistant"; text: string }
  | { type: "tool"; id: string; name: string; status: "running" | "done" | "error"; summary?: string }
  | { type: "confirm"; id: string; name: string; input: unknown; resolved?: string }
  | { type: "error"; text: string };

// ── SSE helpers ──────────────────────────────────────────────────────────────

function llmURL(path: string): string {
  return getPlatformURL(`llm/${path}`).toString();
}

// parseSSE invokes cb(event, data) for each complete `event:/data:` frame in the
// streamed response body.
async function streamSSE(
  resp: Response,
  cb: (event: string, data: any) => void
): Promise<void> {
  const reader = resp.body!.getReader();
  const decoder = new TextDecoder();
  let buf = "";
  for (;;) {
    const { done, value } = await reader.read();
    if (done) break;
    buf += decoder.decode(value, { stream: true });
    let idx: number;
    while ((idx = buf.indexOf("\n\n")) >= 0) {
      const frame = buf.slice(0, idx);
      buf = buf.slice(idx + 2);
      let event = "message";
      let data = "";
      for (const line of frame.split("\n")) {
        if (line.startsWith("event:")) event = line.slice(6).trim();
        else if (line.startsWith("data:")) data += line.slice(5).trim();
      }
      cb(event, data ? JSON.parse(data) : {});
    }
  }
}

// ── component ────────────────────────────────────────────────────────────────

export const AssistantWidget: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [status, setStatus] = useState<IStatus | null>(null);
  const [entries, setEntries] = useState<Entry[]>([]);
  const [input, setInput] = useState("");
  const [busy, setBusy] = useState(false);
  const convId = useRef<string>("");
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open || status) return;
    fetch(llmURL("status"))
      .then((r) => r.json())
      .then(setStatus)
      .catch(() => setStatus({ enabled: false, configured: false, model: "", write_policy: "", tools: [] }));
  }, [open, status]);

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight });
  }, [entries]);

  const handleEvent = useCallback((event: string, data: any) => {
    setEntries((prev) => {
      const next = prev.slice();
      switch (event) {
        case "text": {
          const last = next[next.length - 1];
          if (last && last.type === "assistant") last.text += data.text ?? "";
          else next.push({ type: "assistant", text: data.text ?? "" });
          break;
        }
        case "tool_call":
          next.push({ type: "tool", id: data.id, name: data.name, status: "running" });
          break;
        case "tool_result": {
          const t = next.find((e) => e.type === "tool" && e.id === data.id) as
            | Extract<Entry, { type: "tool" }>
            | undefined;
          if (t) {
            t.status = data.is_error ? "error" : "done";
            t.summary = data.summary;
          }
          break;
        }
        case "confirm_required":
          next.push({ type: "confirm", id: data.id, name: data.name, input: data.input });
          break;
        case "error":
          next.push({ type: "error", text: data.error ?? "unknown error" });
          break;
        case "done":
          if (data.conversationId) convId.current = data.conversationId;
          break;
      }
      return next;
    });
  }, []);

  const sendMessage = useCallback(
    async (message: string) => {
      setBusy(true);
      try {
        const resp = await fetch(llmURL("chat"), {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ conversationId: convId.current, message }),
        });
        if (!resp.ok || !resp.body) {
          const txt = await resp.text();
          handleEvent("error", { error: txt || `HTTP ${resp.status}` });
          return;
        }
        await streamSSE(resp, handleEvent);
      } catch (e) {
        handleEvent("error", { error: String(e) });
      } finally {
        setBusy(false);
      }
    },
    [handleEvent]
  );

  const onSubmit = useCallback(() => {
    const text = input.trim();
    if (!text || busy) return;
    setEntries((prev) => [...prev, { type: "user", text }]);
    setInput("");
    void sendMessage(text);
  }, [input, busy, sendMessage]);

  const onConfirm = useCallback(
    async (entry: Extract<Entry, { type: "confirm" }>, approve: boolean) => {
      if (!approve) {
        setEntries((prev) =>
          prev.map((e) => (e.type === "confirm" && e.id === entry.id ? { ...e, resolved: "declined" } : e))
        );
        void sendMessage("I declined that action. Do not perform it.");
        return;
      }
      try {
        const resp = await fetch(llmURL("confirm"), {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name: entry.name, input: entry.input }),
        });
        const result = await resp.json();
        const summary = result.ok ? JSON.stringify(result.result) : `error: ${result.error}`;
        setEntries((prev) =>
          prev.map((e) => (e.type === "confirm" && e.id === entry.id ? { ...e, resolved: summary } : e))
        );
        void sendMessage(`I approved "${entry.name}". Result: ${summary}. Please continue.`);
      } catch (e) {
        handleEvent("error", { error: String(e) });
      }
    },
    [sendMessage, handleEvent]
  );

  if (!open) {
    return (
      <button className="assistant-fab" title="Library assistant" onClick={() => setOpen(true)}>
        ✦
      </button>
    );
  }

  const unavailable = status && (!status.enabled || !status.configured);

  return (
    <div className="assistant-panel">
      <div className="assistant-header">
        <span className="assistant-title">Library Assistant</span>
        {status?.model && <span className="assistant-model">{status.model}</span>}
        <button className="assistant-close" onClick={() => setOpen(false)}>
          ×
        </button>
      </div>

      <div className="assistant-body" ref={scrollRef}>
        {unavailable && (
          <div className="assistant-notice">
            {!status?.enabled
              ? "The assistant is disabled in settings."
              : "No model gateway configured. Set STASH_ASSISTANT_BASE_URL (and STASH_ASSISTANT_MODEL) to your OpenAI-compatible gateway and restart."}
          </div>
        )}
        {!unavailable && entries.length === 0 && (
          <div className="assistant-notice">
            Ask about your library — e.g. “how many scenes do I have?”, “find unorganized scenes from
            studio X”, or “tag these as review”.
          </div>
        )}
        {entries.map((e, i) => {
          switch (e.type) {
            case "user":
              return (
                <div key={i} className="assistant-msg user">
                  {e.text}
                </div>
              );
            case "assistant":
              return (
                <div key={i} className="assistant-msg assistant">
                  {e.text}
                </div>
              );
            case "tool":
              return (
                <div key={i} className={`assistant-tool ${e.status}`}>
                  <code>{e.name}</code>
                  {e.status === "running" ? " …" : ""}
                  {e.summary ? <span className="assistant-tool-summary">{e.summary}</span> : null}
                </div>
              );
            case "confirm":
              return (
                <div key={i} className="assistant-confirm">
                  <div>
                    Confirm <code>{e.name}</code>: <code>{JSON.stringify(e.input)}</code>
                  </div>
                  {e.resolved ? (
                    <div className="assistant-confirm-resolved">{e.resolved}</div>
                  ) : (
                    <div className="assistant-confirm-actions">
                      <button onClick={() => onConfirm(e, true)}>Approve</button>
                      <button onClick={() => onConfirm(e, false)}>Decline</button>
                    </div>
                  )}
                </div>
              );
            case "error":
              return (
                <div key={i} className="assistant-msg error">
                  {e.text}
                </div>
              );
          }
        })}
        {busy && <div className="assistant-msg assistant thinking">…</div>}
      </div>

      <div className="assistant-input">
        <textarea
          value={input}
          placeholder="Ask the assistant…"
          disabled={!!unavailable}
          onChange={(ev) => setInput(ev.target.value)}
          onKeyDown={(ev) => {
            if (ev.key === "Enter" && !ev.shiftKey) {
              ev.preventDefault();
              onSubmit();
            }
          }}
        />
        <button disabled={busy || !!unavailable} onClick={onSubmit}>
          Send
        </button>
      </div>
    </div>
  );
};

export default AssistantWidget;

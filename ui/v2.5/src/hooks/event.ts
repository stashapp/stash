class StashEvent extends EventTarget {
  dispatch(event: string, id?: string, data?: object) {
    this.dispatchEvent(
      new CustomEvent(`stash:${event}${id ? `:${id}` : ""}`, {
        detail: {
          event: event,
          ...(id ? { id } : {}),
          ...(data ? { data } : {}),
        },
      })
    );
  }
}

const Event = new StashEvent();

export default Event;

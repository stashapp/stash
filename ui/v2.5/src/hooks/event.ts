class StashEvent extends EventTarget {
  dispatch(event: string, id?: string, data?: object) {
    event = `stash:${event}${id ? `:${id}` : ""}`;

    this.dispatchEvent(
      new CustomEvent(event, {
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

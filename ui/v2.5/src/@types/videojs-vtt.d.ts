/* eslint-disable @typescript-eslint/naming-convention */

declare module "videojs-vtt.js" {
  /**
   * A custom JS error object that is reported through the parser's `onparsingerror` callback.
   * It has a name, code, and message property, along with all the regular properties that come with a JavaScript error object.
   *
   * There are two error codes that can be reported back currently:
   * * 0 BadSignature
   * * 1 BadTimeStamp
   *
   * Note: Exceptions other then ParsingError will be thrown and not reported.
   */
  class ParsingError extends Error {
    readonly name: string;
    readonly code: number;
    readonly message: string;
  }

  export namespace WebVTT {
    /**
     * A parser for the WebVTT spec in JavaScript.
     */
    class Parser {
      /**
       * The Parser constructor is passed a window object with which it will create new `VTTCues` and `VTTRegions`
       * as well as an optional `StringDecoder` object which it will use to decode the data that the `parse()` function receives.
       * For ease of use, a `StringDecoder` is provided via `WebVTT.StringDecoder()`.
       * If a custom `StringDecoder` object is passed in it must support the API specified by the #whatwg string encoding spec.
       *
       * @param window the window object to use
       * @param vttjs the vtt.js module
       * @param decoder the decoder to decode `parse()` data with
       */
      constructor(window: Window);
      constructor(window: Window, decoder: TextDecoder);
      constructor(
        window: Window,
        vttjs: typeof import("videojs-vtt.js"),
        decoder: TextDecoder
      );

      /**
       * Callback that is invoked for every region that is correctly parsed. Is passed a `VTTRegion` object.
       */
      onregion?: (cue: VTTRegion) => void;

      /**
       * Callback that is invoked for every cue that is fully parsed. In case of streaming parsing,
       * `oncue` is delayed until the cue has been completely received. Is passed a `VTTCue` object.
       */
      oncue?: (cue: VTTCue) => void;

      /**
       * Is invoked in response to `flush()` and after the content was parsed completely.
       */
      onflush?: () => void;

      /**
       * Is invoked when a parsing error has occurred. This means that some part of the WebVTT file markup is badly formed.
       * Is passed a `ParsingError` object.
       */
      onparsingerror?: (e: ParsingError) => void;

      /**
       * Hands data in some format to the parser for parsing. The passed data format is expected to be decodable by the
       * StringDecoder object that it has. The parser decodes the data and reassembles partial data (streaming), even across line breaks.
       *
       * @param data data to be parsed
       */
      parse(data: string): this;

      /**
       * Indicates that no more data is expected and will force the parser to parse any unparsed data that it may have.
       * Will also trigger `onflush`.
       */
      flush(): this;
    }

    /**
     * Helper to allow strings to be decoded instead of the default binary utf8 data.
     */
    function StringDecoder(): TextDecoder;

    /**
     * Parses the cue text handed to it into a tree of DOM nodes that mirrors the internal WebVTT node structure of the cue text.
     * It uses the window object handed to it to construct new HTMLElements and returns a tree of DOM nodes attached to a top level div.
     *
     * @param window window object to use
     * @param cuetext cue text to parse
     */
    function convertCueToDOMTree(
      window: Window,
      cuetext: string
    ): HTMLDivElement | null;

    /**
     * Converts the cuetext of the cues passed to it to DOM trees - by calling convertCueToDOMTree - and then runs the
     * processing model steps of the WebVTT specification on the divs. The processing model applies the necessary CSS styles
     * to the cue divs to prepare them for display on the web page. During this process the cue divs get added to a block level element (overlay).
     * The overlay should be a part of the live DOM as the algorithm will use the computed styles (only of the divs) to do overlap avoidance.
     *
     * @param overlay A block level element (usually a div) that the computed cues and regions will be placed into.
     */
    function processCues(
      window: Window,
      cues: VTTCue[],
      overlay: Element
    ): void;
  }
}

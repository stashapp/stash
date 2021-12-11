////////////////////////////////////////////////////////////////////////////////
//
// @small-tech/vite-plugin-sri
//
// Subresource integrity (SRI) plugin for Vite (https://vitejs.dev/)
//
// Adds subresource integrity hashes to script and stylesheet
// imports from your index.html file at build time.
//
// If you’re looking for a generic Rollup plugin that does the same thing,
// see rollup-plugin-sri by Jonas Kruckenberg that this one was inspired by:
// https://github.com/JonasKruckenberg/rollup-plugin-sri
//
// Like this? Fund us!
// https://small-tech.org/fund-us
//
// Copyright ⓒ 2021-present Aral Balkan, Small Technology Foundation
// License: ISC.
//
////////////////////////////////////////////////////////////////////////////////

// Stash 2021-12-10: Took changes from pending PR to support variable basepath, disabled remote resource loading

import { createHash } from 'crypto'
import cheerio from 'cheerio'
fs = require('fs');

export default function sri() {
  let config

  return {
    name: 'vite-plugin-sri',
    enforce: 'post',
    apply: 'build',

    configResolved(resolvedConfig) {
      config = resolvedConfig
    },

    async transformIndexHtml(html, context) {
      const bundle = context.bundle

      const getResource = async(path) => {
        const { base } = config
        const pathWithoutBase = base.length && path.startsWith(base) ? path.slice(base.length) : path

        if (pathWithoutBase in bundle) {
          // Load local source from bundle.
          const bundleItem = bundle[pathWithoutBase]
          return bundleItem.code ?? bundleItem.source
        } else {
          throw new Error("A remote resource is in the build!")
        }
      }

      const calculateIntegrityHashes = async (element) => {
        const attributeName = element.attribs.src ? 'src' : 'href'
        const resourcePath = element.attribs[attributeName]
        const source = await getResource(resourcePath)
        if (source === null) {
          console.warn(`Could not resolve resource '${resourcePath}'`)
          return;
        }
        
        element.attribs.integrity = `sha384-${createHash('sha384').update(source).digest().toString('base64')}`
      }

      const $ = cheerio.load(html)
      $.prototype.asyncForEach = async function (callback) {
        for (let index = 0; index < this.length; index++) {
          await callback(this[index], index, this)
        }
      }

      // Implement SRI for scripts and stylesheets.
      const scripts = $('script').filter('[src]')
      const modules = $('link[rel=modulepreload]').filter('[href]')
      const stylesheets = $('link[rel=stylesheet]').filter('[href]')

      await scripts.asyncForEach(calculateIntegrityHashes)
      await modules.asyncForEach(calculateIntegrityHashes)
      await stylesheets.asyncForEach(calculateIntegrityHashes)

      return $.html()
    }
  }
}

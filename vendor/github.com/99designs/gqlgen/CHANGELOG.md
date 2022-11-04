# CHANGELOG
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased](https://github.com/99designs/gqlgen/compare/v0.17.19...HEAD)

<!-- end of if -->
<!-- end of CommitGroups -->
<a name="v0.17.19"></a>
## [v0.17.19](https://github.com/99designs/gqlgen/compare/v0.17.18...v0.17.19) - 2022-09-15
- <a href="https://github.com/99designs/gqlgen/commit/588c6ac137b8ed7aea1bc7c009ea23cb9dec5caa"><tt>588c6ac1</tt></a> release v0.17.19

- <a href="https://github.com/99designs/gqlgen/commit/c671317056298db8073498c8db02120b6f737032"><tt>c6713170</tt></a> v0.17.18 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.18"></a>
## [v0.17.18](https://github.com/99designs/gqlgen/compare/v0.17.17...v0.17.18) - 2022-09-15
- <a href="https://github.com/99designs/gqlgen/commit/1d41c808a93446fca8ff867e957ef552e56f6ae3"><tt>1d41c808</tt></a> release v0.17.18

- <a href="https://github.com/99designs/gqlgen/commit/4dbe2e475f15ce77a498c841ea6c9149ef5ceaba"><tt>4dbe2e47</tt></a> update graphiql to 2.0.7 (<a href="https://github.com/99designs/gqlgen/pull/2375">#2375</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b7cc094a49e3d348cfc457aa76f1640c86cdcae9"><tt>b7cc094a</tt></a> testfix: make apollo federated tracer test more consistent (<a href="https://github.com/99designs/gqlgen/pull/2374">#2374</a>)</summary>

* Update tracing_test.go

* add missing imports

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d096fb9b08531b0dc389a786b6f44add045ea75e"><tt>d096fb9b</tt></a> Update directives (<a href="https://github.com/99designs/gqlgen/pull/2371">#2371</a>)

- <a href="https://github.com/99designs/gqlgen/commit/1acfea2fbdf3564df16f8023f4e736e90a05b909"><tt>1acfea2f</tt></a> Add v0.17.17 changelog

- <a href="https://github.com/99designs/gqlgen/commit/c273adc8ad45e15940bbb6fe211603670d9f3220"><tt>c273adc8</tt></a> v0.17.17 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.17"></a>
## [v0.17.17](https://github.com/99designs/gqlgen/compare/v0.17.16...v0.17.17) - 2022-09-13
- <a href="https://github.com/99designs/gqlgen/commit/d50bc5aca10c5a5dd6a1680b2288c35a61327ade"><tt>d50bc5ac</tt></a> release v0.17.17

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/462025b400e9b792a5afbe320cde4cc952f6b547"><tt>462025b4</tt></a> nil check error before type assertion follow-up from <a href="https://github.com/99designs/gqlgen/pull/2341">#2341</a> (<a href="https://github.com/99designs/gqlgen/pull/2368">#2368</a>)</summary>

* Improve errcode.Set safety

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/59493aff86020d170e58900654d334f5ebc2ceee"><tt>59493aff</tt></a> fix: apollo federation tracer was race prone (<a href="https://github.com/99designs/gqlgen/pull/2366">#2366</a>)</summary>

The tracer was using a global state across different goroutines
Added req headers to operation context to allow it to be fetched in InterceptOperation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fc0185567f2dfc37b38f11283efb9cc1db69e96d"><tt>fc018556</tt></a> Update gqlparser to v2.5.1 (<a href="https://github.com/99designs/gqlgen/pull/2363">#2363</a>)

- <a href="https://github.com/99designs/gqlgen/commit/56574a146bd16a13c9055128ec3c80e96a7c4b29"><tt>56574a14</tt></a> feat: make Playground HTML content compatible with UTF-8 charset (<a href="https://github.com/99designs/gqlgen/pull/2355">#2355</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/182b039d34cb730f432c486ebe763f246937dea4"><tt>182b039d</tt></a> Add `subscriptions.md` recipe to docs (<a href="https://github.com/99designs/gqlgen/pull/2346">#2346</a>)</summary>

* Add `subscriptions.md` recipe to docs

* Fix wrong request type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b66fff16de0b16edc317398a5574fcff2cb39e66"><tt>b66fff16</tt></a> Add omit_getters config option (<a href="https://github.com/99designs/gqlgen/pull/2348">#2348</a>)

- <a href="https://github.com/99designs/gqlgen/commit/2ba8040f20e32d06dc6d5bfacaadc5619a6e66ee"><tt>2ba8040f</tt></a> Update changelog for v0.17.16

- <a href="https://github.com/99designs/gqlgen/commit/8bef8c8061222071e6c814e45bbc33fcabcb3980"><tt>8bef8c80</tt></a> v0.17.16 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.16"></a>
## [v0.17.16](https://github.com/99designs/gqlgen/compare/v0.17.15...v0.17.16) - 2022-08-26
- <a href="https://github.com/99designs/gqlgen/commit/9593ceadd6e07c6fd0f0b0e0c55b9f1bf8ade762"><tt>9593cead</tt></a> release v0.17.16

- <a href="https://github.com/99designs/gqlgen/commit/2390af2db920dc632fe47bc778a24c30495b9efd"><tt>2390af2d</tt></a> Update gqlparser to v2.5.0 (<a href="https://github.com/99designs/gqlgen/pull/2341">#2341</a>)

- <a href="https://github.com/99designs/gqlgen/commit/2a87fe0645fd271e4e71d2b7bde34ecf31bf844c"><tt>2a87fe06</tt></a> feat: update Graphiql to version 2 (<a href="https://github.com/99designs/gqlgen/pull/2340">#2340</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/32e2ccd30e82fc566ca022a65dcc4a67c4b6125a"><tt>32e2ccd3</tt></a> Update yaml to v3 (<a href="https://github.com/99designs/gqlgen/pull/2339">#2339</a>)</summary>

* update yaml to v3

* add missing go entry for yaml on _example

* add missing sum file

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7949117a524be7f8882a61e2d4ade1bedf105107"><tt>7949117a</tt></a> v0.17.15 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.15"></a>
## [v0.17.15](https://github.com/99designs/gqlgen/compare/v0.17.14...v0.17.15) - 2022-08-23
- <a href="https://github.com/99designs/gqlgen/commit/23cc749256b4e2edc4b11ce9e84c643a7bb3194f"><tt>23cc7492</tt></a> release v0.17.15

- <a href="https://github.com/99designs/gqlgen/commit/577a570cdb6b1b9185f24940690a14cdced37a36"><tt>577a570c</tt></a> Markdown formatting fixes (<a href="https://github.com/99designs/gqlgen/pull/2335">#2335</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2b584011fc64a55cbda67f46637a280bf94d9cc1"><tt>2b584011</tt></a> Fix Interface Slice Getter Generation (<a href="https://github.com/99designs/gqlgen/pull/2332">#2332</a>)</summary>

* Make modelgen test fail if generated doesn't build
Added returning list of interface to modelgen test schema

* Implement slice copying when returning interface slices

* Re-generate to satisfy the linter

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/aee57b4c521e527ebc0538b8edfbe610973abf21"><tt>aee57b4c</tt></a> Correct boolean logic (<a href="https://github.com/99designs/gqlgen/pull/2330">#2330</a>)</summary>

Correcting boolean logic issue

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/da0610e11accf3afd34903f03bfc0abd045d07ed"><tt>da0610e1</tt></a> Update changelog for v0.17.14

- <a href="https://github.com/99designs/gqlgen/commit/ddcb524e3321d849505f6937307ef3dcbd3acace"><tt>ddcb524e</tt></a> v0.17.14 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.14"></a>
## [v0.17.14](https://github.com/99designs/gqlgen/compare/v0.17.13...v0.17.14) - 2022-08-18
- <a href="https://github.com/99designs/gqlgen/commit/581bf6eb063a0d6a3cec3b6bc7a16ca10e310a97"><tt>581bf6eb</tt></a> release v0.17.14

- <a href="https://github.com/99designs/gqlgen/commit/d3384377aefb4b7d34ba52f8def6c0a6a3dec27f"><tt>d3384377</tt></a> Update gqlparser

- <a href="https://github.com/99designs/gqlgen/commit/c2d02d352f8d531fa0bd9b246fc152eeb6dbf10a"><tt>c2d02d35</tt></a> More descriptive `not implemented` stubs (<a href="https://github.com/99designs/gqlgen/pull/2328">#2328</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/2327"> #2327</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9f919d2cee464acdaf4a490aeb42d63369dbd572"><tt>9f919d2c</tt></a> Avoid GraphQL to Go Naming Collision with "ToGoModelName" func (<a href="https://github.com/99designs/gqlgen/pull/2322">#2322</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/2321"> #2321</a>)</summary>

* using ReplaceAllStringLiteral

* fixing wordInfo template test

* bumping linter timeout to 5m

* comment cleanup

* some cleanup, adding "ToGoPrivateModelName" func

* adding "ToGoPrivateModelName" func

* refactoring word walker impl and tests

* hopefully making linter happy

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2304c104fc8d26487f50e80e9c5eaee113005a30"><tt>2304c104</tt></a> Include docstrings on interface getters (<a href="https://github.com/99designs/gqlgen/pull/2317">#2317</a>)

- <a href="https://github.com/99designs/gqlgen/commit/f5d603269502b50e19d0ed966e2dfe3ecd74049f"><tt>f5d60326</tt></a> Leverage (*Imports).LookupType when generating interface field getters (<a href="https://github.com/99designs/gqlgen/pull/2315">#2315</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/242c3ba217ee740e37445ce4b14e0808554263f5"><tt>242c3ba2</tt></a> Generate getters for interface fields (<a href="https://github.com/99designs/gqlgen/pull/2314">#2314</a>)</summary>

* Generate getters for interface fields

* Changes to make models_test.go pass

* Use text/template, not html/template

* Re-run go generate ./...

* gofmt a few files that were failing lint checks

* Another gofmt straggler

* Try making the "generated" match the exact whitespace github is disliking

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0d91c893e285cc14330c80643b663cd2bebeb911"><tt>0d91c893</tt></a> Add hackernews graphql api tutorial to other resources (<a href="https://github.com/99designs/gqlgen/pull/2305">#2305</a>)

- <a href="https://github.com/99designs/gqlgen/commit/c2526ba50ff3a69b5eca88a62a571c47f3c245ed"><tt>c2526ba5</tt></a> Update gqlparser to v2.4.7 (<a href="https://github.com/99designs/gqlgen/pull/2300">#2300</a>)

- <a href="https://github.com/99designs/gqlgen/commit/f283124d1cea309e054afb197d16012364b88097"><tt>f283124d</tt></a> <a href="https://github.com/99designs/gqlgen/pull/2298">#2298</a>: fix gqlgen extracting module name from comment line (<a href="https://github.com/99designs/gqlgen/pull/2299">#2299</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/779d7cdd4991e3cf4bf1ecbdea1f02664a56ac8d"><tt>779d7cdd</tt></a> Add support for KeepAlive message in websocket client (<a href="https://github.com/99designs/gqlgen/pull/2293">#2293</a>)</summary>

* Add support for KeepAlive message in websocket client

* rewrite if-else to switch statement

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5a37d1dc079f5212b6e043b0f6889cae7b08dea9"><tt>5a37d1dc</tt></a> v0.17.13 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.13"></a>
## [v0.17.13](https://github.com/99designs/gqlgen/compare/v0.17.12...v0.17.13) - 2022-07-15
- <a href="https://github.com/99designs/gqlgen/commit/e82b6bf1cf311d6af2e280127f47b15ae35ca6ac"><tt>e82b6bf1</tt></a> release v0.17.13

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f0e9047df5f86efbfbceea9c04593bb1f52e06de"><tt>f0e9047d</tt></a> Hide dependencies in `tools.go` from importers (<a href="https://github.com/99designs/gqlgen/pull/2287">#2287</a>)</summary>

Projects that use `go mod vendor` will vendor `github.com/matryer/moq`
despite it not being required at runtime.

Moving `tools.go` to `internal` hides this import from downstream
users and avoids `github.com/matryer/moq` being vendored.

`go generate` of the mocks still works as expected.

The assumption behind the import test broke, so I've pointed it at a
different path that has no Go code. This seems to match the intent
behind the original test for the `internal/code/..` path.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6310e6a736ccbf3bb8caea981553ee7549aea748"><tt>6310e6a7</tt></a> support named interface to Field.CallArgs (<a href="https://github.com/99designs/gqlgen/pull/2289">#2289</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/30493696aacf79090bb5e144a304a5a7df488c67"><tt>30493696</tt></a> fix: return the original error (<a href="https://github.com/99designs/gqlgen/pull/2288">#2288</a>)</summary>

* fix: return the original error

close https://github.com/99designs/gqlgen/issues/2286

* Update error.go

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fb13091df76b47b936224336fe19b15fe310b41d"><tt>fb13091d</tt></a> updated WebSocker InitFunc recipe (<a href="https://github.com/99designs/gqlgen/pull/2275">#2275</a>)

- <a href="https://github.com/99designs/gqlgen/commit/770c09fb9db0485943590b9986afe36818c2a70e"><tt>770c09fb</tt></a> Update changelog for v0.17.12

- <a href="https://github.com/99designs/gqlgen/commit/b4c186a7142c0a151b6a21b40914fe317e13819d"><tt>b4c186a7</tt></a> v0.17.12 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.12"></a>
## [v0.17.12](https://github.com/99designs/gqlgen/compare/v0.17.11...v0.17.12) - 2022-07-04
- <a href="https://github.com/99designs/gqlgen/commit/94c02b0de6d483d87453fc18a7f7625ae4adaa6c"><tt>94c02b0d</tt></a> release v0.17.12

- <a href="https://github.com/99designs/gqlgen/commit/7eb8ba93daacef77ca7266fdfb9e5abc8a720eb7"><tt>7eb8ba93</tt></a> Fix CreateTodo (<a href="https://github.com/99designs/gqlgen/pull/2256">#2256</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0b0e5ce4afc5b503217304f89914b2e903c05fa5"><tt>0b0e5ce4</tt></a> Replace use of strings.Title with cases.Title (<a href="https://github.com/99designs/gqlgen/pull/2268">#2268</a>)</summary>

* github: Test more go versions

* github: Fix ci tests

* github: Increase verbosity, sleep

* github: Drop bash

* github: Test go 1.18 and newer node verisons

* github: Pull out node 16 for now

* github: Only lint 1.16 for now

* cases: Use cases.Title over strings.Title which is deprecated

* gqlgen: Remove use of deprecated strings.Title

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0c11e5fdd8ec4fd7612b857c4c554e1ef463d194"><tt>0c11e5fd</tt></a> parse at beginning of do function (<a href="https://github.com/99designs/gqlgen/pull/2269">#2269</a>)

- <a href="https://github.com/99designs/gqlgen/commit/edb1c585c1c49102dc962e0ac3bd271688e51ecf"><tt>edb1c585</tt></a> Update Changelog for v0.17.11

- <a href="https://github.com/99designs/gqlgen/commit/5e6b52fddab835611513e3572f23716666ebae58"><tt>5e6b52fd</tt></a> v0.17.11 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.11"></a>
## [v0.17.11](https://github.com/99designs/gqlgen/compare/v0.17.10...v0.17.11) - 2022-07-03
- <a href="https://github.com/99designs/gqlgen/commit/ea294c4ea344186c3b41b82d5f1c60138f6ce05e"><tt>ea294c4e</tt></a> release v0.17.11

- <a href="https://github.com/99designs/gqlgen/commit/8ebf75c19d775ddbd12b3d94461b605ef4c5f711"><tt>8ebf75c1</tt></a> Update gqlparser (<a href="https://github.com/99designs/gqlgen/pull/2270">#2270</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b8497f52fde0803277d981405cd2e42ce0455a70"><tt>b8497f52</tt></a> github: Fix CI pipelines (<a href="https://github.com/99designs/gqlgen/pull/2266">#2266</a>)</summary>

* github: Test more go versions

* github: Fix ci tests

* github: Increase verbosity, sleep

* github: Drop bash

* github: Test go 1.18 and newer node verisons

* github: Pull out node 16 for now

* github: Only lint 1.16 for now

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c287a7b0b40cdd8c94077ed54fff257fe383e796"><tt>c287a7b0</tt></a> codegen: fix resolvers execution order (<a href="https://github.com/99designs/gqlgen/pull/2267">#2267</a>)</summary>

* codegen: fix run order of resolver


* fix: update code generate

* fix: update stub, root to generate resolver for input

* fix: added unit-test for input field order

* fix: added test for singlefile

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8481457f2fd6ae711f688fc4726d724df5992b8c"><tt>8481457f</tt></a> gqlgen: Add resolver comment generation and preservation (<a href="https://github.com/99designs/gqlgen/pull/2263">#2263</a>)</summary>

* gqlgen: Add resolver comment generation and preservation

* gqlgen: Regenerate

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/532d46af5b2b97f5b69ebd1ce261a191e2690fa3"><tt>532d46af</tt></a> Make uploads content seekable (<a href="https://github.com/99designs/gqlgen/pull/2247">#2247</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/34bbc450c502919cd46c5eefcc66341ef697c0e8"><tt>34bbc450</tt></a> Use the go:embed API to lookup templates (<a href="https://github.com/99designs/gqlgen/pull/2262">#2262</a>)</summary>

* Switch the templates package internally to read from TemplateFS

Users are expected to pass in the FS by using the embed API.

* Update all usages of templates.Render to use the TemplateFS option

* Fix unit tests

* Fix linter error

* Commit generated changes

Doesn't look like anything has changed though. Maybe just a different
whitespace character.

* Fix test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/53ca207a4c53c78e4dec1e4d34d3b9251dd00b0b"><tt>53ca207a</tt></a> Fix PR links in CHANGELOG.md (<a href="https://github.com/99designs/gqlgen/pull/2257">#2257</a>)</summary>

* fix "PR" regex in CHANGELOG-full-history.tpl.md

* regenerate CHANGELOG.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/53ada82edb7e8bb91059cbf3f344270a668934c5"><tt>53ada82e</tt></a> Replace deprecated ioutil pkg with os & io (<a href="https://github.com/99designs/gqlgen/pull/2254">#2254</a>)</summary>

As of Go 1.16, the same functionality is now provided by package io or
package os, and those implementations should be preferred in new code.

So replacing all usage of ioutil pkg with io & os.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a8f112e0c5b3466de2f550e68cbe872233f052ac"><tt>a8f112e0</tt></a> update changelog

- <a href="https://github.com/99designs/gqlgen/commit/82fbbe4163459cbb8d862c99931fcd015ed756e6"><tt>82fbbe41</tt></a> v0.17.10 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.10"></a>
## [v0.17.10](https://github.com/99designs/gqlgen/compare/v0.17.9...v0.17.10) - 2022-06-13
- <a href="https://github.com/99designs/gqlgen/commit/4ff9ea92b0d90f7fdc7c22bec592fbec1aca60a6"><tt>4ff9ea92</tt></a> release v0.17.10

- <a href="https://github.com/99designs/gqlgen/commit/cac4f40486edd280654412485979ff619238a877"><tt>cac4f404</tt></a> update gqlparser (<a href="https://github.com/99designs/gqlgen/pull/2239">#2239</a>)

- <a href="https://github.com/99designs/gqlgen/commit/d07ec12d69d3db0e5b502be7528884fdb5fb7593"><tt>d07ec12d</tt></a> Use exact capitalization from field names overridden in config (<a href="https://github.com/99designs/gqlgen/pull/2237">#2237</a>)

- <a href="https://github.com/99designs/gqlgen/commit/3a64078299f0417fca48c652620015937cb19c5a"><tt>3a640782</tt></a> fix: <a href="https://github.com/99designs/gqlgen/pull/2234">#2234</a> (<a href="https://github.com/99designs/gqlgen/pull/2235">#2235</a>) Response.Errors in DispatchError function is not PresentedError

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c355df9efc910053e888922edc14170271968671"><tt>c355df9e</tt></a> fix <a href="https://github.com/99designs/gqlgen/pull/1876">#1876</a>: Optional Any type should allow nil values (<a href="https://github.com/99designs/gqlgen/pull/2231">#2231</a>)</summary>

* Anonymous func that checks value of arg type interface for nil

* Added unit test for `CallArgs()`

* Fixed type of argument in unit test

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/65e68108d926faf635285144d7b6670f7f6d9ce4"><tt>65e68108</tt></a> Add config boolean for whether resolvers return pointers (<a href="https://github.com/99designs/gqlgen/pull/2175">#2175</a>)

- <a href="https://github.com/99designs/gqlgen/commit/ddd825ef62f1fa7cbc0824c1696f72a3c67d78e0"><tt>ddd825ef</tt></a> Only make cyclical struct fields pointers (<a href="https://github.com/99designs/gqlgen/pull/2174">#2174</a>)

- <a href="https://github.com/99designs/gqlgen/commit/5a87fe29353e3fee987a39431df0322d12b575f9"><tt>5a87fe29</tt></a> Update websocket.go (<a href="https://github.com/99designs/gqlgen/pull/2223">#2223</a>)

- <a href="https://github.com/99designs/gqlgen/commit/e2edda5d5d02a1496dcf9eb48ab95ecd8f07f018"><tt>e2edda5d</tt></a> Update dataloaders.MD (<a href="https://github.com/99designs/gqlgen/pull/2221">#2221</a>)

- <a href="https://github.com/99designs/gqlgen/commit/3de7d2cf730cc060a27f6d1c815742d1a9f479cd"><tt>3de7d2cf</tt></a> fix: chat example frontend race condition (<a href="https://github.com/99designs/gqlgen/pull/2219">#2219</a>)

- <a href="https://github.com/99designs/gqlgen/commit/11f405724f19fa4b93120746fdf74c1d97f4575b"><tt>11f40572</tt></a> Update Changelog

- <a href="https://github.com/99designs/gqlgen/commit/caca01fb6c64ad078ff195cddf52dd6966e7995e"><tt>caca01fb</tt></a> v0.17.9 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.9"></a>
## [v0.17.9](https://github.com/99designs/gqlgen/compare/v0.17.8...v0.17.9) - 2022-05-26
- <a href="https://github.com/99designs/gqlgen/commit/7f0611b2d19833a740afcfaf5708febff942da2d"><tt>7f0611b2</tt></a> release v0.17.9

- <a href="https://github.com/99designs/gqlgen/commit/738209b26337bc1116be7b0afacc83eae6bb93b0"><tt>738209b2</tt></a> Update gqlparser (<a href="https://github.com/99designs/gqlgen/pull/2216">#2216</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6855b7290cab62a1fc6a26a2b633e0b5bbf248da"><tt>6855b729</tt></a> fix: prevent goroutine leak and CPU spinning at websocket transport (<a href="https://github.com/99designs/gqlgen/pull/2209">#2209</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/2168"> #2168</a>)</summary>

* Added goroutine leak test for chat example

* Improved chat example with proper concurrency


This reverts commit eef7bfaad1b524f9e2fc0c1150fdb321c276069e.

* Improved subscription channel usage

* Regenerated examples and codegen

* Add support for subscription keepalives in websocket client

* Update chat example test

* if else chain to switch


* Revert "Add support for subscription keepalives in websocket client"

This reverts commits 64b882c3c9901f25edc0684ce2a1f9b63443416b and 670cf22272b490005d46dc2bee1634de1cd06d68.

* Fixed chat example race condition

* Fixed chatroom#Messages type

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5f5bfcb97fdb01026cf35a5dc46f1246a30f9b26"><tt>5f5bfcb9</tt></a> fix <a href="https://github.com/99designs/gqlgen/pull/2204">#2204</a> - don't try to embed builtin sources (<a href="https://github.com/99designs/gqlgen/pull/2214">#2214</a>)</summary>

* dont't try to embed builtins

* add test

* generated code

* fix error message string

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8d9d3f125f13dcd19f59072d3c38366dc520758b"><tt>8d9d3f12</tt></a> Check only direct dependencies (<a href="https://github.com/99designs/gqlgen/pull/2205">#2205</a>)

- <a href="https://github.com/99designs/gqlgen/commit/b262e40a485f67d2659e239a156418938d0fe2e9"><tt>b262e40a</tt></a> v0.17.8 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.8"></a>
## [v0.17.8](https://github.com/99designs/gqlgen/compare/v0.17.7...v0.17.8) - 2022-05-25
- <a href="https://github.com/99designs/gqlgen/commit/25367e0a24998aea40f09218f60d1d0e6d1cce4a"><tt>25367e0a</tt></a> release v0.17.8

- <a href="https://github.com/99designs/gqlgen/commit/5a56b69d89c7414e21b2f01e0e5042a26b69c5cb"><tt>5a56b69d</tt></a> Add security workflow with nancy (<a href="https://github.com/99designs/gqlgen/pull/2202">#2202</a>)

- <a href="https://github.com/99designs/gqlgen/commit/482f4ce08e65458cec2dbfaf7d184f1c8fccb129"><tt>482f4ce0</tt></a> Run CI tests on windows (<a href="https://github.com/99designs/gqlgen/pull/2199">#2199</a>)

- <a href="https://github.com/99designs/gqlgen/commit/656045d3fa643b898932c3f5332544b0baed1af4"><tt>656045d3</tt></a> This works on Windows too! (<a href="https://github.com/99designs/gqlgen/pull/2197">#2197</a>)

- <a href="https://github.com/99designs/gqlgen/commit/f6aeed60a508dae102b2b821d3a947e24e5e0826"><tt>f6aeed60</tt></a> Merge branch 'master' of github.com:99designs/gqlgen

- <a href="https://github.com/99designs/gqlgen/commit/d91080be396af96266941499d369d0f8279761b0"><tt>d91080be</tt></a> Update changelog

- <a href="https://github.com/99designs/gqlgen/commit/752d2d7e9fff08c82a6d3ffc1c8c7ffe2a2e9fe2"><tt>752d2d7e</tt></a> v0.17.7 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.7"></a>
## [v0.17.7](https://github.com/99designs/gqlgen/compare/v0.17.6...v0.17.7) - 2022-05-24
- <a href="https://github.com/99designs/gqlgen/commit/2b1dff1b71f89c95e946bbe5948b7061f9c47aa8"><tt>2b1dff1b</tt></a> release v0.17.7

- <a href="https://github.com/99designs/gqlgen/commit/b2087f944d9b9af6e776a9d97662c9e8b86a8c3b"><tt>b2087f94</tt></a> Update module dependencies (<a href="https://github.com/99designs/gqlgen/pull/2192">#2192</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8825ac460b047e22724ed7728c7d7ffbf1b523a9"><tt>8825ac46</tt></a> Fix misprint (<a href="https://github.com/99designs/gqlgen/pull/2187">#2187</a>)</summary>

* Fix misprint

* Fix misprint

* Re-generate

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/41daa5d8dc1e35bdbfe68e95b37c10599b224456"><tt>41daa5d8</tt></a> fix <a href="https://github.com/99designs/gqlgen/pull/2190">#2190</a> - don't use backslash for "embed" paths on windows (<a href="https://github.com/99designs/gqlgen/pull/2191">#2191</a>)

- <a href="https://github.com/99designs/gqlgen/commit/0cce5544f06fd84831ef2ca0f60e16f7f554814d"><tt>0cce5544</tt></a> Update Changelog

- <a href="https://github.com/99designs/gqlgen/commit/26644541aafbcdb46f10a6ff0f5894637227c331"><tt>26644541</tt></a> v0.17.6 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.6"></a>
## [v0.17.6](https://github.com/99designs/gqlgen/compare/v0.17.5...v0.17.6) - 2022-05-23
- <a href="https://github.com/99designs/gqlgen/commit/358d45dcfc2b022fdda9476a37f44c0622607ae9"><tt>358d45dc</tt></a> release v0.17.6

- <a href="https://github.com/99designs/gqlgen/commit/7c95938c5f1278fa14a13a92eb88d117102e0330"><tt>7c95938c</tt></a> Improve operation error handling (<a href="https://github.com/99designs/gqlgen/pull/2184">#2184</a>)

- <a href="https://github.com/99designs/gqlgen/commit/2526f6871166377b4f444ad8d22577a632b0abf4"><tt>2526f687</tt></a> Correct identation (<a href="https://github.com/99designs/gqlgen/pull/2182">#2182</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f7bf453c79d82b01ed4baed894043aaff645bf2f"><tt>f7bf453c</tt></a> Bump dset from 3.1.1 to 3.1.2 in /integration (<a href="https://github.com/99designs/gqlgen/pull/2176">#2176</a>)</summary>

Bumps [dset](https://github.com/lukeed/dset) from 3.1.1 to 3.1.2.
- [Release notes](https://github.com/lukeed/dset/releases)
- [Commits](https://github.com/lukeed/dset/compare/v3.1.1...v3.1.2)

---
updated-dependencies:
- dependency-name: dset
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4cdf70261a9dc8af589399f09b56d0f90606a9fa"><tt>4cdf7026</tt></a> Update getting-started.md (<a href="https://github.com/99designs/gqlgen/pull/2157">#2157</a>)</summary>

Fix getting-started missing fields resolver config

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/eef7bfaad1b524f9e2fc0c1150fdb321c276069e"><tt>eef7bfaa</tt></a> fix: prevents goroutine leak at websocket transport (<a href="https://github.com/99designs/gqlgen/pull/2168">#2168</a>)

- <a href="https://github.com/99designs/gqlgen/commit/b8ec51d8629a24288353b4ee4be70fff3645b03e"><tt>b8ec51d8</tt></a> go: update gqlparser to latest (<a href="https://github.com/99designs/gqlgen/pull/2149">#2149</a>)

- <a href="https://github.com/99designs/gqlgen/commit/ec3e597e7b45e17464cd8c7faaa51e75755ce3cf"><tt>ec3e597e</tt></a> Fix docs bug in field collection (<a href="https://github.com/99designs/gqlgen/pull/2141">#2141</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f6b352316fae4b4fdc6317e24ea94ba48ac29e85"><tt>f6b35231</tt></a> Add argument to WebsocketErrorFunc (<a href="https://github.com/99designs/gqlgen/pull/2124">#2124</a>)</summary>

* Add argument to WebsocketErrorFunc

to determine whether the error ocured on read or write to the websocket.

* Wrap websocket error

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0f016df3ae7ee4898358dc67a491689164297df6"><tt>0f016df3</tt></a> Fix invalid query parameter for playground subscription endpoint (<a href="https://github.com/99designs/gqlgen/pull/2148">#2148</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fb5751ab478603a864977f9fbe70655776d7fb55"><tt>fb5751ab</tt></a> use "embed" in generated code (<a href="https://github.com/99designs/gqlgen/pull/2119">#2119</a>)</summary>

* use "embed" in generated code

* don't use embed for builtins

* working poc

* handle no embeddable sources

* fix dir

* comment

* add test for embedding

* improve error handling

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d38911f1a9d7f0ec39a74a95994d95291f1922c3"><tt>d38911f1</tt></a> Allow absolute https://github.com/99designs/gqlgens to the GraphQL playground (<a href="https://github.com/99designs/gqlgen/pull/2142">#2142</a>)</summary>

* Allow absolute URLs to the GraphQL playground

* Add test for playground URLs

* Close res.Body in playground test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3228f36fec50930483801b27b92658592fab5e87"><tt>3228f36f</tt></a> Update getting-started.md (<a href="https://github.com/99designs/gqlgen/pull/2140">#2140</a>)</summary>

* Update getting-started.md

function rand.Int requires two parameters and returns two value in golang version 1.18.1.

* Highlight the package used so people don't pick crypto/rand

* Revert to original

* Remove extra space

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/33fe0b9b824ec86699059f410505c02659fc6c81"><tt>33fe0b9b</tt></a> Update package.json (<a href="https://github.com/99designs/gqlgen/pull/2138">#2138</a>)</summary>

I added `graphql-ws` because there is no graphql-ws in package.json

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f8e837b824ef4903a60f3cb974ef72fb4718a858"><tt>f8e837b8</tt></a> Use MultipartReader to parse file uploads (<a href="https://github.com/99designs/gqlgen/pull/2135">#2135</a>)</summary>

Use a streaming MultipartReader to parse requests with file
uploads. The GraphQL multipart request specification guarantees
that the operations and map form fields will come first.

There are two reasons motivating this change:

- This allows for file uploads without specifying a specific
  filename.
- This avoids unnecessary copies for requests with more than one
  file. Go's ParseForm already copies the request's body into
  memory or on disk. We were also doing this manually as a second
  step.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/05bfc1fb12f73648833e1055e775e074a6df7eed"><tt>05bfc1fb</tt></a> Upddate Changelog

- <a href="https://github.com/99designs/gqlgen/commit/62f694f0a8cf24f52dffca5823bb44fa1c32f97b"><tt>62f694f0</tt></a> v0.17.5 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.5"></a>
## [v0.17.5](https://github.com/99designs/gqlgen/compare/v0.17.4...v0.17.5) - 2022-04-29
- <a href="https://github.com/99designs/gqlgen/commit/fd97e74eafc898278fd4b74477cb053393672232"><tt>fd97e74e</tt></a> release v0.17.5

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9250f9ac1f90b27da0bd8583ef8dcf0894d70686"><tt>9250f9ac</tt></a> Feature: Add FTV1 Support via Handler (<a href="https://github.com/99designs/gqlgen/pull/2132">#2132</a>)</summary>

* initial support for ftv1 traces via handler

* remove testing json extension

* remove binary from commit and add to .gitignore

* updating go.mod

* updating examples go.sum

* rerunning generate within the examples folder

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fce3a11a9f570ffed3e9035d32deddfb3076c2cf"><tt>fce3a11a</tt></a> feat: added graphql.UnmarshalInputFromContext (<a href="https://github.com/99designs/gqlgen/pull/2131">#2131</a>)</summary>

* feat: added graphql.UnmarshalInputFromContext

* chore: run go generate for _examples

* fix: apply suggestions from code review


* fix: update error cases

* fix: fixed unit-test by update root_.gotpl

* fix: apply suggestions from code review

* fix: update graphql/input.go

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6a24e88147fb2523af0163d7fa84d296b5e32e4d"><tt>6a24e881</tt></a> update instructions to specify package of Role (<a href="https://github.com/99designs/gqlgen/pull/2130">#2130</a>)</summary>

Can't compile with the example unless I also include `model.` for Role.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ccfa245b1eb2657e588bf73f4df0e99f96869cbd"><tt>ccfa245b</tt></a> Ignore protobuf files in coverage (<a href="https://github.com/99designs/gqlgen/pull/2133">#2133</a>)

- <a href="https://github.com/99designs/gqlgen/commit/0465dcb1e8e4177945c2670f15316ac96e4b992a"><tt>0465dcb1</tt></a> Update federation.md (<a href="https://github.com/99designs/gqlgen/pull/2129">#2129</a>)

- <a href="https://github.com/99designs/gqlgen/commit/8f0631dcd3ca6fcfcd3dc6e69f4a92fec54e6dc7"><tt>8f0631dc</tt></a> Update Changelog

- <a href="https://github.com/99designs/gqlgen/commit/41611560d45f1226b860e795bcb35b5ecf09c5b3"><tt>41611560</tt></a> v0.17.4 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.4"></a>
## [v0.17.4](https://github.com/99designs/gqlgen/compare/v0.17.3...v0.17.4) - 2022-04-25
- <a href="https://github.com/99designs/gqlgen/commit/d6de831a28a0f1d8834c5dba4216dcd763814d3f"><tt>d6de831a</tt></a> release v0.17.4

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2a2a3dcb67c7d713e41476eac47e20ab0e21fba7"><tt>2a2a3dcb</tt></a> Feature: Adds Federation 2 Support (<a href="https://github.com/99designs/gqlgen/pull/2115">#2115</a>)</summary>

* fed2 rough support

* autodetection of fed2

* adding basic tests for changes

* fixing docs

* Update plugin/federation/federation.go

* removing custom scalar since it was causing issues

* fixing lint test

* should fix for real this time

* fixing test failures

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/77260e88c853a047e4e61a5357ceda4a5ea26405"><tt>77260e88</tt></a> shorten some generated code (<a href="https://github.com/99designs/gqlgen/pull/2120">#2120</a>)</summary>

* shorten some generated code

* generate examples

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4da17e1c7a59149eb6c2f5d60fcf11a2374b2488"><tt>4da17e1c</tt></a> update modules except mapstructure (<a href="https://github.com/99designs/gqlgen/pull/2118">#2118</a>)</summary>

* Update modules


* Update modules except for mapstructure


* Try to update to v1.3.1

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cddbf02d494e3aeaac3f60d1708b25facc5b767d"><tt>cddbf02d</tt></a> Update Changelog file

- <a href="https://github.com/99designs/gqlgen/commit/8f80f4efe8947b55919ce37291f4e908f57fd8dc"><tt>8f80f4ef</tt></a> v0.17.3 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.3"></a>
## [v0.17.3](https://github.com/99designs/gqlgen/compare/v0.17.2...v0.17.3) - 2022-04-20
- <a href="https://github.com/99designs/gqlgen/commit/0bb262d1a0143f60640f60ebbb516e0f4cd79042"><tt>0bb262d1</tt></a> release v0.17.3

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8d0bd22aff1cdb6ad2e36190e11871b169f8da0a"><tt>8d0bd22a</tt></a> Update gqlparser (<a href="https://github.com/99designs/gqlgen/pull/2109">#2109</a>)</summary>

* Update gqlparser


* Update tests to be NoError

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ec0dea883a2c967d533e5f1530791ad72a08198b"><tt>ec0dea88</tt></a> Fix the ability of websockets to get errors (<a href="https://github.com/99designs/gqlgen/pull/2097">#2097</a>)</summary>

Because DispatchOperation creates tempResponseContext,
which is passed into Exec, which is then used in _Subscription to
generate the next function. Inside the various subscription functions
when generating next the context was captured there.

Which means later when the returned function from DispatchOperation is
called. The responseContext which accumulates the errors is the
tempResponseContext which we no longer have access to to read the errors
out of it.

Instead add a context to next() so that it can be passed through and
accumulated the errors as expected.

Added a unit test for this as well.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e3f04b42f1fc5d4b13dc0579b2ec713f770a4fd0"><tt>e3f04b42</tt></a> Change the error message to be consumer targeted (<a href="https://github.com/99designs/gqlgen/pull/2096">#2096</a>)</summary>

* Change the error message to be slightly more clear

* Rebase on updated origin/master.

Fix the test to not be sensitive to array ordering.
Re-generate on master as there was a schema change.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5a49764956ffb674df2c9bee19455bb1fd3407db"><tt>5a497649</tt></a> Fix websocket subscriptions to not double close. (<a href="https://github.com/99designs/gqlgen/pull/2095">#2095</a>)</summary>

We were closing at the end of the loop and also in the defer.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a15a9bfdbad30b2f5ce7a966ec1190c108c4df3e"><tt>a15a9bfd</tt></a> Update test.yml to be valid

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a1538928a569a09834579db941863ccce28113e3"><tt>a1538928</tt></a> Use Github API to update the docs (<a href="https://github.com/99designs/gqlgen/pull/2101">#2101</a>)</summary>

* Use Github API to update the docs

Instead of a hard-coded version of the docs we want to realease, this
uses the Github API to get the last 20 versions and publish those. This
will allow any script invoking this to make sure to always have the
latest version of the docs

* Reinstate set -e

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3bf437c232f8be30a473cf94495a1014c0583af2"><tt>3bf437c2</tt></a> Update golangci-lint (<a href="https://github.com/99designs/gqlgen/pull/2103">#2103</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/12c6d0bf15431f666d08c4c82581957e1b727898"><tt>12c6d0bf</tt></a> Fix misprint (<a href="https://github.com/99designs/gqlgen/pull/2102">#2102</a>)</summary>

* Fix misprint

* Update websocket_graphql_transport_ws.go

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9f5fad13fa6275139e051788cc5fe8c2b2630428"><tt>9f5fad13</tt></a> Bump minimist from 1.2.5 to 1.2.6 in /integration (<a href="https://github.com/99designs/gqlgen/pull/2085">#2085</a>)</summary>

Bumps [minimist](https://github.com/substack/minimist) from 1.2.5 to 1.2.6.
- [Release notes](https://github.com/substack/minimist/releases)
- [Commits](https://github.com/substack/minimist/compare/1.2.5...1.2.6)

---
updated-dependencies:
- dependency-name: minimist
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/035e1d6eeb81179ddec3d36d8776212d8fe35cd6"><tt>035e1d6e</tt></a> Add AllowedMethods field to transport.Options (<a href="https://github.com/99designs/gqlgen/pull/2080">#2080</a>)</summary>

* Add AllowedMethods field to transport.Options

to enable users to specify allowed HTTP methods.

* Update graphql/handler/transport/options.go

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f0fdb116f45350aabf698c20bf6410283f96bb11"><tt>f0fdb116</tt></a> Add instructions for enabling autobinding (<a href="https://github.com/99designs/gqlgen/pull/2079">#2079</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/12b0b38583e2c7b2174585bf1243a98cbbc2eba6"><tt>12b0b385</tt></a> Bump Playground version (<a href="https://github.com/99designs/gqlgen/pull/2078">#2078</a>)</summary>

* update playground

* enables tabs

* update shas

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1324c3ffb9ff0afef6e9cc41d99b5b4b9bc928b6"><tt>1324c3ff</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/2062">#2062</a> from a8m/childfield</summary>

graphql: add FieldContext.Child field function and enable it in codegen

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bf9caeaee091e32178fe2906894a7c7e72fdd66d"><tt>bf9caeae</tt></a> graphql: add FieldContext.ChildArgs field and enable it in codegen

- <a href="https://github.com/99designs/gqlgen/commit/36fb3dc6733601f96162bc80fccda42e34b3b7ff"><tt>36fb3dc6</tt></a> codegen: allow binding methods with optional variadic arguments (<a href="https://github.com/99designs/gqlgen/pull/2066">#2066</a>)

- <a href="https://github.com/99designs/gqlgen/commit/fba5edd4fa1176ef0f2840f3bb90fe10b9f4b695"><tt>fba5edd4</tt></a> Update Changelog

- <a href="https://github.com/99designs/gqlgen/commit/48b2b7e1521c50d03cabf524bbb78805e3fb023f"><tt>48b2b7e1</tt></a> v0.17.2 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.2"></a>
## [v0.17.2](https://github.com/99designs/gqlgen/compare/v0.17.1...v0.17.2) - 2022-03-21
- <a href="https://github.com/99designs/gqlgen/commit/1f04d38a4441c5de6171400218b9dd25cebb3639"><tt>1f04d38a</tt></a> release v0.17.2

- <a href="https://github.com/99designs/gqlgen/commit/87fc5f22e8fbfa28a180cbf0e7008af9f830273e"><tt>87fc5f22</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1961">#1961</a> for Go 1.18 (<a href="https://github.com/99designs/gqlgen/pull/2052">#2052</a>)

- <a href="https://github.com/99designs/gqlgen/commit/f85d59d30ae055fd89b79aa3d7e3ca1c7fcaedfa"><tt>f85d59d3</tt></a> fixed modelgen test schema (<a href="https://github.com/99designs/gqlgen/pull/2032">#2032</a>)

- <a href="https://github.com/99designs/gqlgen/commit/d873ff8bb9927b302752bd48d7836f2597db558e"><tt>d873ff8b</tt></a> v0.17.1 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.1"></a>
## [v0.17.1](https://github.com/99designs/gqlgen/compare/v0.17.0...v0.17.1) - 2022-03-02
- <a href="https://github.com/99designs/gqlgen/commit/5ea50aee16088ed414be73ca9a59a90f622c9483"><tt>5ea50aee</tt></a> release v0.17.1

- <a href="https://github.com/99designs/gqlgen/commit/a493a4239673c5922281628fc8b94c727398283e"><tt>a493a423</tt></a> Prepare for new release

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9f520a2897cf42750e7290cbd83de6fdf13f2e75"><tt>9f520a28</tt></a> Update golangci-lint and fix resource leak (<a href="https://github.com/99designs/gqlgen/pull/2024">#2024</a>)</summary>

* Fix golangci-lint in CI

* Fix resource leak

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/74baaa14c924871d100b56e8103ec27d678c33d0"><tt>74baaa14</tt></a> fixed model gen for multiple implemented type (<a href="https://github.com/99designs/gqlgen/pull/2021">#2021</a>)

- <a href="https://github.com/99designs/gqlgen/commit/d31cf6bed5712e4015c286498ede649894e48d01"><tt>d31cf6be</tt></a> v0.17.0 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.17.0"></a>
## [v0.17.0](https://github.com/99designs/gqlgen/compare/v0.16.0...v0.17.0) - 2022-03-01
- <a href="https://github.com/99designs/gqlgen/commit/e4be56513300286729b1276de2741ce6a93f3afa"><tt>e4be5651</tt></a> release v0.17.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/082bbff65eaf9931d4637001635d72014033523f"><tt>082bbff6</tt></a> Revert "Update quickstart (<a href="https://github.com/99designs/gqlgen/pull/1850">#1850</a>)" (<a href="https://github.com/99designs/gqlgen/pull/2014">#2014</a>)</summary>

This reverts commit 0ab636144bfc875f86e4d9fd7a2686bc57d5050c.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a58411b804f848aa7e2e4547a1ba768f5dfdc8d3"><tt>a58411b8</tt></a> Embed templates instead of inlining them (<a href="https://github.com/99designs/gqlgen/pull/2019">#2019</a>)

- <a href="https://github.com/99designs/gqlgen/commit/839b50df1c068e6b17adb27a68a26984ea363bcc"><tt>839b50df</tt></a> Test gqlgen generate in CI (<a href="https://github.com/99designs/gqlgen/pull/2017">#2017</a>)

- <a href="https://github.com/99designs/gqlgen/commit/00dc14ad817806840ca4df4a04d7a658f6f38105"><tt>00dc14ad</tt></a> Remove ambient imports (<a href="https://github.com/99designs/gqlgen/pull/2016">#2016</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/45e192ea9fa2af6ed3b16e1a8b5c67276f13d34f"><tt>45e192ea</tt></a> Clean up docs to clarify how to use a particular version (<a href="https://github.com/99designs/gqlgen/pull/2015">#2015</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/1851"> #1851</a>)</summary>

This reverts commit 57a148f6d12572fe585ecfcafafbb7441dbf9cab.

* Update getting-started.md

* Update getting-started.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3a9413f718b217866c89cb88e268e8f2c461fb95"><tt>3a9413f7</tt></a> Fix issue template

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5236fb096802dc66ba9d11d096b2c6fa1ad24b14"><tt>5236fb09</tt></a> fix introspection for description to be nullable (<a href="https://github.com/99designs/gqlgen/pull/2008">#2008</a>)</summary>

* fixed introspection for description to be nullable

* regenerated for integration

* regenerated

* fixed introspection package

* regenerated

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/82fefdb51046ca80506292f25dcb2d636301f865"><tt>82fefdb5</tt></a> support to generate model for intermediate interface (<a href="https://github.com/99designs/gqlgen/pull/1982">#1982</a>)</summary>

* support to generate model for intermediate interface

* go generate ./... in example

* fixed filepath generation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3ec8363554ff17e3ffb3e86c58f1ee2d5689e798"><tt>3ec83635</tt></a> Bump ajv from 6.10.2 to 6.12.6 in /integration (<a href="https://github.com/99designs/gqlgen/pull/2007">#2007</a>)</summary>

Bumps [ajv](https://github.com/ajv-validator/ajv) from 6.10.2 to 6.12.6.
- [Release notes](https://github.com/ajv-validator/ajv/releases)
- [Commits](https://github.com/ajv-validator/ajv/compare/v6.10.2...v6.12.6)

---
updated-dependencies:
- dependency-name: ajv
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9546de2c825c92318230d21d525791dfa3f0f184"><tt>9546de2c</tt></a> Web Socket initialization message timeout (<a href="https://github.com/99designs/gqlgen/pull/2006">#2006</a>)</summary>

* Added an optional timeout to the web socket initialization message read operation.

* Added a fail message to a web socket init read timeout test.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f6ea623003fe8b8f40beb1b545a6dc91a2af0f12"><tt>f6ea6230</tt></a> fixed introspection for schema description and specifiedByhttps://github.com/99designs/gqlgen (<a href="https://github.com/99designs/gqlgen/pull/1986">#1986</a>)</summary>

* fixed introspection for schema description and specifiedByURL

* updated to the master latest

* fixed Description resolver

* updated integration go file

* fixed codegen tests for the latest gqlparser

* updated go mod in example

* go generate

* skip specifiedBy

* regenerate

* fixed schema-expected.graphql for the latest

* fixed integration test to use latest tools

* fixed integration workflow

* use v2.4.0

* fixed sum

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f17ca15e1fab4837ebcca99958e1651034e48852"><tt>f17ca15e</tt></a> Fix broken links in docs (<a href="https://github.com/99designs/gqlgen/pull/1983">#1983</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/1734"> #1734</a>)

- <a href="https://github.com/99designs/gqlgen/commit/a0c856b72e1e001633388644310a388342d7d0ff"><tt>a0c856b7</tt></a> Added a callback error handling function to the websocket and added tests for it. (<a href="https://github.com/99designs/gqlgen/pull/1975">#1975</a>)

- <a href="https://github.com/99designs/gqlgen/commit/cfea9f07627143fd184e8f36448cb501006bc63a"><tt>cfea9f07</tt></a> generate resolvers for input types (<a href="https://github.com/99designs/gqlgen/pull/1950">#1950</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ffa857ef346e87215bd985b5e84330b0f31afe96"><tt>ffa857ef</tt></a> Websocket i/o timeout fix (<a href="https://github.com/99designs/gqlgen/pull/1973">#1973</a>)</summary>

* Renamed "pingMesageType" to "pingMessageType" and refactored websocket_graphqlws.go to look more like websocket_graphql_transport_ws.go for the sake of consistency.

* Made the keep-alive messages graphql-ws only, and the ping-pong messages graphql-transport-ws only (and added tests for it).

* gofmt

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d7da5b0d3b3cffa0cdb30fa7dcf16c87d8434e7e"><tt>d7da5b0d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1958">#1958</a> from 99designs/cleanup-main</summary>

Cleanup main

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/42f32432b068b381d4666a56da50ebf73520831f"><tt>42f32432</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1957">#1957</a> from 99designs/move-init-ci</summary>

Upate init CI step

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/be1647480a2c89bac80c094e0fa7f5ffda7fe602"><tt>be164748</tt></a> Cleanup main

- <a href="https://github.com/99designs/gqlgen/commit/8ea290c0854579f7cb6746fb7e51485fda01a6b6"><tt>8ea290c0</tt></a> Upate init CI step

- <a href="https://github.com/99designs/gqlgen/commit/56bfb1880603486189025f043d8155beaa9f53d2"><tt>56bfb188</tt></a> Fix 1955: only print message on [@key](https://github.com/key) found on interfaces (<a href="https://github.com/99designs/gqlgen/pull/1956">#1956</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/213a085b954945eeaa1fc87f0fedef2c07fe14c1"><tt>213a085b</tt></a> rename "example" dir to "_examples" (<a href="https://github.com/99designs/gqlgen/pull/1734">#1734</a>)</summary>

* rename "example" dir to "_examples"

* fix lint

* Adjust permissions

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9262b35865e4a8c62749a0bdf4cce075877a8b4f"><tt>9262b358</tt></a> fix: typo in dataloader code sample (<a href="https://github.com/99designs/gqlgen/pull/1954">#1954</a>)</summary>

* fix: typo in dataloader code sample

* rename k to key for sample to compile

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a05437332754fab8bebce18ae8371dcdbe05460b"><tt>a0543733</tt></a> remove autobind example (<a href="https://github.com/99designs/gqlgen/pull/1949">#1949</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/06bbca37286edb32db2480db6d200b709ae071a9"><tt>06bbca37</tt></a> docs: migrate dataloaders sample to graph-gophers/dataloader (<a href="https://github.com/99designs/gqlgen/pull/1871">#1871</a>)</summary>

* docs: add dataloader sample

* finish example

* add example

* simplify method

* replace old example

* styling

* Update docs/content/reference/dataloaders.md

* Update docs/content/reference/dataloaders.md

* Update docs/content/reference/dataloaders.md

* Update docs/content/reference/dataloaders.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f9fcfa16a13c64ecfb298ad6cf97b3548b7ee0ff"><tt>f9fcfa16</tt></a> Comment out autobind in the sample config file (<a href="https://github.com/99designs/gqlgen/pull/1872">#1872</a>)</summary>

The reason is that many people using it for the first time copy exactly that configuration example and then open the issues to say it doesn't work.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a30b68de58cbd5fa19e88bd5a198e7ca67147f3b"><tt>a30b68de</tt></a> fix: whitelist VERSION and CURRENT_VERSION env vars (<a href="https://github.com/99designs/gqlgen/pull/1870">#1870</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/76a533b8161ad3ee56b721404650697cce808221"><tt>76a533b8</tt></a> Bump gopkg.in/yaml.v2 from 2.2.4 to 2.2.8 (<a href="https://github.com/99designs/gqlgen/pull/1858">#1858</a>)</summary>

* Bump gopkg.in/yaml.v2 from 2.2.4 to 2.2.8

Bumps [gopkg.in/yaml.v2](https://github.com/go-yaml/yaml) from 2.2.4 to 2.2.8.
- [Release notes](https://github.com/go-yaml/yaml/releases)
- [Commits](https://github.com/go-yaml/yaml/compare/v2.2.4...v2.2.8)

---
updated-dependencies:
- dependency-name: gopkg.in/yaml.v2
  dependency-type: direct:production
...


* Update go sum for example

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/eed4301c7123e329e067600fe25aa1c876c99b8d"><tt>eed4301c</tt></a> Bump node-fetch from 2.6.1 to 2.6.7 in /integration (<a href="https://github.com/99designs/gqlgen/pull/1859">#1859</a>)</summary>

Bumps [node-fetch](https://github.com/node-fetch/node-fetch) from 2.6.1 to 2.6.7.
- [Release notes](https://github.com/node-fetch/node-fetch/releases)
- [Commits](https://github.com/node-fetch/node-fetch/compare/v2.6.1...v2.6.7)

---
updated-dependencies:
- dependency-name: node-fetch
  dependency-type: direct:development
...

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7f5dce6d9eebe829edf20ba67e5c5019921fa1ee"><tt>7f5dce6d</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1776">#1776</a> : Edit and persist headers in GraphiQL (<a href="https://github.com/99designs/gqlgen/pull/1856">#1856</a>)

- <a href="https://github.com/99designs/gqlgen/commit/e0b42f9981825814e9ec7c216d9c393b2486831c"><tt>e0b42f99</tt></a> fix requires directive with nested field when entityResolver directive is used (<a href="https://github.com/99designs/gqlgen/pull/1863">#1863</a>)

- <a href="https://github.com/99designs/gqlgen/commit/25c2cdcb12d574e597e8902ab9e2d94b1e5ef974"><tt>25c2cdcb</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1636">#1636</a> by updating gqlparser (<a href="https://github.com/99designs/gqlgen/pull/1857">#1857</a>)

- <a href="https://github.com/99designs/gqlgen/commit/c161ab382547948feb9a54e7bf9ff17458b8c3c9"><tt>c161ab38</tt></a> fix <a href="https://github.com/99designs/gqlgen/pull/1770">#1770</a> minor error in getting-started.md (<a href="https://github.com/99designs/gqlgen/pull/1771">#1771</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/57a148f6d12572fe585ecfcafafbb7441dbf9cab"><tt>57a148f6</tt></a> Remove outdated version reference so example is always for latest (<a href="https://github.com/99designs/gqlgen/pull/1851">#1851</a>)</summary>

* Also update version reference to next

* Update getting-started.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0ab636144bfc875f86e4d9fd7a2686bc57d5050c"><tt>0ab63614</tt></a> Update quickstart (<a href="https://github.com/99designs/gqlgen/pull/1850">#1850</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a8eba26dec863c5d34905ab6408f970b4d2abdb5"><tt>a8eba26d</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1777">#1777</a> by updating version constant and adding release checklist (<a href="https://github.com/99designs/gqlgen/pull/1848">#1848</a>)</summary>

* Revise to use script 🤦

</details></dd></dl>

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.16.0"></a>
## [v0.16.0](https://github.com/99designs/gqlgen/compare/v0.15.1...v0.16.0) - 2022-01-24
- <a href="https://github.com/99designs/gqlgen/commit/b90f9750f40583823a3e875d6bbe1538ce50f527"><tt>b90f9750</tt></a> Merge branch 'master' of github.com:99designs/gqlgen

- <a href="https://github.com/99designs/gqlgen/commit/99523e44ae67633ecfa714794a209191d3519017"><tt>99523e44</tt></a> Prepare for v0.16.0 release (<a href="https://github.com/99designs/gqlgen/pull/1842">#1842</a>)

- <a href="https://github.com/99designs/gqlgen/commit/0563146c6bd7188b2ae187040c5a7f3d17cc9f89"><tt>0563146c</tt></a> Prepare for v0.16.0 release

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7cefef26f7b714eb9d5be117ba2159d4e40168f3"><tt>7cefef26</tt></a> add PrependPlugin (<a href="https://github.com/99designs/gqlgen/pull/1839">#1839</a>)</summary>

* add PrependPlugin

related: https://github.com/99designs/gqlgen/pull/1838

* added test for PrependPlugin

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/972878a04fc8e9065df212b7cdbe934f141d069b"><tt>972878a0</tt></a> Revert "Fix plugin addition (<a href="https://github.com/99designs/gqlgen/pull/1717">#1717</a>)" (<a href="https://github.com/99designs/gqlgen/pull/1838">#1838</a>)</summary>

This reverts commit f591c8f797e35635fb5eb0e4465c77b6a073896b.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1ed7e050a1f736b7395dd3e61f771d3ddcf80a8d"><tt>1ed7e050</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1832">#1832</a> [@requires](https://github.com/requires) directive when [@entityResolver](https://github.com/entityResolver) is used (<a href="https://github.com/99designs/gqlgen/pull/1833">#1833</a>)</summary>

* fix requires directive for multipleEntity directive


* fix lint

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fcee4c404f52c2debcb8b8acaa31421804f625ea"><tt>fcee4c40</tt></a> Update README.md (<a href="https://github.com/99designs/gqlgen/pull/1836">#1836</a>)</summary>

Corrected a simple grammar typo.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3fb5fd9916896d2c084c646093b3cbd634f61121"><tt>3fb5fd99</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1834">#1834</a>: Implement federation correctly (<a href="https://github.com/99designs/gqlgen/pull/1835">#1835</a>)</summary>

* Fix federation implementation which does not conform to Apollo Federation subgraph specification

* Optimize generated line breaks

* Run go generate

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/986650718267770018139713df874a70e1e79e16"><tt>98665071</tt></a> Imporve gqlgen test cases (<a href="https://github.com/99designs/gqlgen/pull/1773">#1773</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/1765"> #1765</a>)</summary>

* Imporve test cases for init and generate

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5d904d8782833d85ca6fedb47420237b3258eb66"><tt>5d904d87</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1778">#1778</a> from ipfans/gh-pages-patch</summary>

Bump gqlgen.com version list

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/196ee13bc364a4dea800edc0a9e60c6e7e2bd03b"><tt>196ee13b</tt></a> Bump gqlgen.com version

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.15.1"></a>
## [v0.15.1](https://github.com/99designs/gqlgen/compare/v0.15.0...v0.15.1) - 2022-01-16
- <a href="https://github.com/99designs/gqlgen/commit/7102a36bbde485fbbb671499fdde8697232c0725"><tt>7102a36b</tt></a> Prepare for 0.15.1 release

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2b8f50b3b38b129f56e06c3dadccd0cb8667a741"><tt>2b8f50b3</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1765">#1765</a>: Sometimes module info not exists or not loaded. (<a href="https://github.com/99designs/gqlgen/pull/1767">#1767</a>)</summary>

* Remove failing test

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/46502e5e5af713208143f089fa2db5b266aa5fc7"><tt>46502e5e</tt></a> fixed broken link (<a href="https://github.com/99designs/gqlgen/pull/1768">#1768</a>)

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.15.0"></a>
## [v0.15.0](https://github.com/99designs/gqlgen/compare/v0.14.0...v0.15.0) - 2022-01-14
- <a href="https://github.com/99designs/gqlgen/commit/99be19512eb2f7c5f3db3d699eecc8cdd2020d25"><tt>99be1951</tt></a> Prepare for release

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/931271a2f3af4cf50de702855d396cde02a3d99f"><tt>931271a2</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1762">#1762</a>: Reload packages before merging type systems (<a href="https://github.com/99designs/gqlgen/pull/1763">#1763</a>)</summary>

* run gofmt on file

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e5b5e832e5c5ab388af3ebbae8a899f2b362eab1"><tt>e5b5e832</tt></a> Improve performance of MarshalBoolean (<a href="https://github.com/99designs/gqlgen/pull/1757">#1757</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/57664bf0369a843e56f74f4f1951a1808656fe99"><tt>57664bf0</tt></a> Migrate playgrounds to GraphiQL (<a href="https://github.com/99designs/gqlgen/pull/1751">#1751</a>)</summary>

* migrate to GraphiQL playground

* fix lint

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b2a832d51d5d67463d35fc1397b7d6133e1d6b43"><tt>b2a832d5</tt></a> Avoid problems with `val` being undefined in the federation template. (<a href="https://github.com/99designs/gqlgen/pull/1760">#1760</a>)</summary>

* Avoid problems with `val` being undefined in the federation template.

When running gqlgen over our schema, we were seeing errors like:
```
assignments/generated/graphql/service.go:300:4: val declared but not used
```

The generated code looks like this:
```
func entityResolverNameForMobileNavigation(ctx context.Context, rep map[string]interface{}) (string, error) {
        for {
                var (
                        m   map[string]interface{}
                        val interface{}
                        ok  bool
                )
                m = rep
                if _, ok = m["kaid"]; !ok {
                        break
                }
                m = rep
                if _, ok = m["language"]; !ok {
                        break
                }
                return "findMobileNavigationByKaidAndLanguage", nil
        }
        return "", fmt.Errorf("%w for MobileNavigation", ErrTypeNotFound)
}
```

Looking at the code, it's pretty clear that this happens when there
are multiple key-fields, but each of them has only one keyField.Field
entry.  This is because the old code looked at `len(keyFields)` to
decide whether to declare the `val` variable, but looks at
`len(keyField.Field)` for each keyField to decide whether to use the
`val` variable.

The easiest solution, and the one I do in this PR, is to just declare
`val` all the time, and use a null-assignment to quiet the compiler
when it's not used.

* run go generate to update generated files

* run go generate to update moar generated files

* Adding a test for verify that this fixes the issue.

From `plugins/federation`, run the following command and verify that no errors are produced

```
go run github.com/99designs/gqlgen --config testdata/entityresolver/gqlgen.yml
```

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/47015f12e3aa26af251fec67eab50d3388c17efe"><tt>47015f12</tt></a> Added pointer to a solution for `no Go files` err (<a href="https://github.com/99designs/gqlgen/pull/1747">#1747</a>)</summary>

While following the instructions in this getting started guide I run into this error `package github.com/99designs/gqlgen: no Go files` which was pretty annoying to fix. Its a golang issue but for people who are unfamiliar with how the `go generate` command works in vendored projects its a blocker trying to follow the rest of this guide. It will be really nice to at least have a pointer in the guide for people to find a possible solution to the issue while going through the guide. I'm sure many folks have run into this issue given vendoring is now very popular with the latest go releases.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/27a2b210d9137e2f1d341a103542b47dcc783182"><tt>27a2b210</tt></a> Downgrade to Go 1.16 (<a href="https://github.com/99designs/gqlgen/pull/1743">#1743</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/14cfee7002fb07f3f247b13582fe45959f7133de"><tt>14cfee70</tt></a> Support for multiple [@key](https://github.com/key) directives in federation (reworked) (<a href="https://github.com/99designs/gqlgen/pull/1723">#1723</a>)</summary>

* address review comments

- reworked code generation for federation.go
- better checking for missing/incorrect parameters to entity resolver functions
- better tests for generated entity resolvers

Still missing: 
- suggested test for autobind vs non-autobind generation
- could probably clean up generated code spacing, etc

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2747bd5f3c69db7d55db5f10592ecd0accf3499f"><tt>2747bd5f</tt></a> Add CSV and PDF to common initialisms (<a href="https://github.com/99designs/gqlgen/pull/1741">#1741</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/44beadc1d037a0d46c799c85d88dc71a57f3938b"><tt>44beadc1</tt></a> Fix list coercion when using graphql variables (<a href="https://github.com/99designs/gqlgen/pull/1740">#1740</a>)</summary>

* fix(codegen): support coercion of lists in graphql variables

This was broken by an upstream dependency `gqlparser` coercing variables during validation. this has broken the existing coercion process withing `gqlgen`

* test: add list coercion integration tests

* chore: regenerate generated code

* test: update expected schema for integration tests

* chore: run goimports

* chore: regenerate examples

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bd8938d853e0aade3eab106e35f28d586e355124"><tt>bd8938d8</tt></a> fix: automatically register built-in directive goTag (<a href="https://github.com/99designs/gqlgen/pull/1737">#1737</a>)</summary>

* fix: automatically register built-in tag goTag

* doc: add directive config documentation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/497227faf4e43266bf1a8f9ef756b98ef85cfee7"><tt>497227fa</tt></a> Close Websocket Connection on Context close/cancel (<a href="https://github.com/99designs/gqlgen/pull/1728">#1728</a>)</summary>

* Added code to the web socket so it closes when the context is cancelled (with an optional close reason).

* Added a test.

* go fmt


* Fix linter issues about the cancel function being thrown away.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4581fccdb67e05937ba862bf601a18393731bccc"><tt>4581fccd</tt></a> Don't loose field arguments when none match (<a href="https://github.com/99designs/gqlgen/pull/1725">#1725</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/213ecd93c2593a535efc291b065966573d798a5d"><tt>213ecd93</tt></a> Add support for graphql-transport-ws with duplex ping-pong (<a href="https://github.com/99designs/gqlgen/pull/1578">#1578</a>)</summary>

* Add support for graphql-transport-ws with duplex ping-pong

* Add tests for the duplex ping-pong

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ae92c83d7f7d14ab3a29016eb71d577e1c66721e"><tt>ae92c83d</tt></a> add federation tests (<a href="https://github.com/99designs/gqlgen/pull/1719">#1719</a>)

- <a href="https://github.com/99designs/gqlgen/commit/f591c8f797e35635fb5eb0e4465c77b6a073896b"><tt>f591c8f7</tt></a> Fix plugin addition (<a href="https://github.com/99designs/gqlgen/pull/1717">#1717</a>)

- <a href="https://github.com/99designs/gqlgen/commit/8fa6470f9e0eb8d85e90f211a90b0ea38a228f35"><tt>8fa6470f</tt></a> Fix <a href="https://github.com/99designs/gqlgen/pull/1704">#1704</a>: handle [@required](https://github.com/required) nested fields as in [@key](https://github.com/key) (<a href="https://github.com/99designs/gqlgen/pull/1706">#1706</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/af33b7cd2486d52363ba1ed90e1963ce12c22250"><tt>af33b7cd</tt></a> Cleaning up extra return in federation generated code (<a href="https://github.com/99designs/gqlgen/pull/1713">#1713</a>)</summary>

In PR 1709, I introduced GetMany semantics for resolving federated entities.  But I left a couple of extra return statements in the generated code that are not necessary. So Im just cleaning those up here.

Also added `go:generate` in federation entity resolver tests to make it simpler to test.

To test:
```
go generate ./... && cd example/ && go generate ./... && cd ..
go test -race ./... && cd example && go test -race ./... && cd ..
```

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/402a22593f6fe71ff33d805a057de0f331eded5a"><tt>402a2259</tt></a> Optimize performance for binder, imports and packages (Rebased from sbalabanov/master) (<a href="https://github.com/99designs/gqlgen/pull/1711">#1711</a>)</summary>

* Cache go.mod resolution for module name search

* Optimize binder.FindObject() for performance by eliminating repeatitive constructs

* Optimize allocations in packages.Load() function

* Optimize binder.FindObject() by indexing object definitions for each loaded package

* goimports to fix linting

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/237a7e6a7bb1abdd71e957adeb119691bbae5671"><tt>237a7e6a</tt></a> Separate golangci-lint from other jobs (<a href="https://github.com/99designs/gqlgen/pull/1712">#1712</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/50292e99d5cd0021f9fbec6118e406e4da86505b"><tt>50292e99</tt></a> Resolve multiple federated entities in a single entityResolve call (<a href="https://github.com/99designs/gqlgen/pull/1709">#1709</a>)</summary>

* Resolve multiple federated entities in a single entityResolve call

Entity resolver functions can only process one entity at a time. But often we want to resolve all the entities at once so that we can optimize things like database calls. And to do that you need to add you'd need to add batching with abstractions like dataloadgen or batchloader. The drawback here is that the resolver code (the domain logic) gets more complex to implement, test, and debug.

An alternative is to have entity resolvers that can process all the representations in a single call so that domain logic can have access to all the representations up front, which is what Im adding in this PR.

There are a few moving pieces here:
3. When that's configured, the federation plugin will create an entity resolver that will take a list of representations.

Please note that this is very specific to federation and entity resolvers. This does not add support for resolving fields in an entity.

Some of the implementation details worth noting. In order to efficiently process batches of entities, I group them by type so that we can process groups of entities at the same time. The resolution of groups of entities run concurrently in Go routines.  If there is _only_ one type, then that's just processed without concurrency. Entities that don't have multiget enabled will still continue to resolve concurrently with Go routines, and entities that have multiget enabled just get the entire list of representations.

The list of representations that are passed to entity resolvers are strongly types, and the type is generated for you.

There are lots of new tests to ensure that there are no regressions and that the new functionality still functions as expected. To test:
1. Go to `plugin/federation`
2. Generate files with `go run github.com/99designs/gqlgen --config testdata/entityresolver/gqlgen.yml`
3. And run `go test ./...`. Verify they all pass.

You can look at the federated code in `plugin/federation/testdata/entityresolver/gederated/federation.go`

* Added `InputType` in entity to centralize logic for generating types for multiget resolvers.

* reformat and regenerate

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/80713b84cf937b203cbd745603809e0d72f1dd84"><tt>80713b84</tt></a> Adding entity resolver tests for errors, entities with different type… (<a href="https://github.com/99designs/gqlgen/pull/1708">#1708</a>)</summary>

* Adding entity resolver tests for errors, entities with different types, and requires

The tests in this PR are for ensuring we get the expected errors from entity resolvers, that we also handle resolving entities where the representations are for different types, and that requires directive works correctly.

To run tests:
1. Go to `plugin/federation`
2. Generate files with `go run github.com/99designs/gqlgen --config testdata/entityresolver/gqlgen.yml`
3. And run `go test ./...`.  Verify they all pass.

* Fixed test for errors

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ed2d699804c09875c678900937923fcbc4f6e00a"><tt>ed2d6998</tt></a> Replace ! with _ in root.generated file to avoid build conflicts (<a href="https://github.com/99designs/gqlgen/pull/1701">#1701</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/828820afa95e4e66857839cbcbf6584665cef10c"><tt>828820af</tt></a> transport: implement `graphql-transport-ws` ws sub-protocol   (<a href="https://github.com/99designs/gqlgen/pull/1507">#1507</a>)</summary>

* websocket: create `messageExchanger` to handle subprotocol messages

* remove unused type

* typo in comments

* change `graphqlwsMessageType` type to string

* add support for `graphql-transport-ws` subprotocol

* fix chat app example

* update example chat app dependencies

* improve chat app exmaple to use the recommended ws library

* add tests

* removed unused const in tests

* Update example/chat/readme.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/28caa6ce504b4ed206407dcaea8f8418feb91888"><tt>28caa6ce</tt></a> Ignore generated files from test coverage (<a href="https://github.com/99designs/gqlgen/pull/1699">#1699</a>)

- <a href="https://github.com/99designs/gqlgen/commit/7ac988dee1187a27bb1290fe25f7b179f6102e42"><tt>7ac988de</tt></a> Fix linting issue

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/01d3c4f8c0b732ee85b919415322fdb190127fac"><tt>01d3c4f8</tt></a> Entity resolver tests (<a href="https://github.com/99designs/gqlgen/pull/1697">#1697</a>)</summary>

* Moving federation tests to their own folders

Reorganizing the tests in the federation plugin a little bit so make it simpler to add more safely without testdata colliding. This is in anticipation for a follow up PR for adding entity resolver tests.

Run the tests with `go test ./plugin/federation/...` and verify they all pass. Also verify that the testdata/allthething directory has a `generated` directory specific to that test.

NOTE: There is a catch all type of test that I moved to the directory `allthething`.  Open to suggestions for a better name! One potential thing to considere here is to split up the tests that use that testdata and break them down into more specific tests. E.g. Add a multikey test in the testdata/entity.  For now, Im leaving that as a TODO.

* Adding entity resolver tests in the federation plugin

The tests work by sending `_entities` queries with `representation` variables directly to the mocked server, which will allow us to test generated federation code end to end.  For context, the format of the entity query is something like:

```
query($representations:[_Any!]!){_entities(representations:$representations){ ...on Hello{secondary} }}
```

And `representations` are the list of federated keys for the entities being resovled, and they look like

```
representations: [{
   "__typename": "Hello",
   "name":       "federated key value 1",
}, {
   "__typename": "Hello",
   "name":       "federated key value 2",
}]
```

The entity resolver tests are in `plugin/federation/federation_entityresolver_test.go` and they rely on `plugin/federation/testdata/entityresolver`.

To run the tests:
1. Build the entityresolver testdata
  - From plugin/federation, run `go run github.com/99designs/gqlgen --config testdata/entityresolver/gqlgen.yml`
2. Run the tests with `go test ./...` or similar

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b7db36d368c260d90fb5fa6084c295d92c1a001d"><tt>b7db36d3</tt></a> Revert "Support for multiple [@key](https://github.com/key) directives in federation (<a href="https://github.com/99designs/gqlgen/pull/1684">#1684</a>)" (<a href="https://github.com/99designs/gqlgen/pull/1698">#1698</a>)</summary>

This reverts commit 47de912f56cd4bd6da9b74929cd67b8881617026.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4a4b5601a661259cc33d22c8a254803025e8a1f7"><tt>4a4b5601</tt></a> DOC: Fixed indention in example code. (<a href="https://github.com/99designs/gqlgen/pull/1693">#1693</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/47de912f56cd4bd6da9b74929cd67b8881617026"><tt>47de912f</tt></a> Support for multiple [@key](https://github.com/key) directives in federation (<a href="https://github.com/99designs/gqlgen/pull/1684">#1684</a>)</summary>

* add more unit test coverage to plugin/federation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/59a30919a8d8a9b67972cc7e4dd1e425901f15c2"><tt>59a30919</tt></a> Reimplement goTag using FieldMutateHook (<a href="https://github.com/99designs/gqlgen/pull/1682">#1682</a>)</summary>

* Reimplement goTag using a FieldMutateHook

This change does not change the logic of goTag, merely reimplements it using a FieldMutateHook and sets it as the default FieldMutateHook for the modelgen plugin.

* Add repeated tag test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/37a4e7eefa241dd7e5eeab7bbca56b1677d37daf"><tt>37a4e7ee</tt></a> Rename `[@extraTag](https://github.com/extraTag)` directive to `[@goTag](https://github.com/goTag)` and make repeatable (<a href="https://github.com/99designs/gqlgen/pull/1680">#1680</a>)</summary>

* Allow Repeatable `goTag` Directive

* Default to field name if none provided

* Update Docs

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/87f9e436e922a977b6002bb3070f087641001782"><tt>87f9e436</tt></a> Fix nil pointer dereference when an invalid import is bound to a model (<a href="https://github.com/99designs/gqlgen/pull/1676">#1676</a>)</summary>

* Fixes remaining Name field in singlefile test

* Fixes nill pointer dereference when an invalid import is bound to a model

* Only return error if we failed to find type

* Revert "Fixes remaining Name field in singlefile test"

This reverts commit e43ebf7aa80f884afdb3feca90867b1eff593f01.

* Undo change of log.Println -> fmt.Println

Totally accidental, sorry!

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6c65e8f15389e8aad7e20e7ba4a9b3ff4be565d9"><tt>6c65e8f1</tt></a> Update getting-started.md (<a href="https://github.com/99designs/gqlgen/pull/1674">#1674</a>)</summary>

missing an 's' on quoted filename default

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3bbc2a342fc7a0839b25e26ebd0aa25d7f498dac"><tt>3bbc2a34</tt></a> feat: generate resolvers for inputs if fields are missing (<a href="https://github.com/99designs/gqlgen/pull/1404">#1404</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7db941a56e742dca10cbba1e32d4d09458cdc8ef"><tt>7db941a5</tt></a> Fix 1138: nested fieldset support (<a href="https://github.com/99designs/gqlgen/pull/1669">#1669</a>)</summary>

* formatting

* update federation schema to latest Apollo spec


also:
handle extra spaces in FieldSet
upgrade deps in federation integration tests

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/488a31fc12979b825b166cf6f317f4b27ee456a0"><tt>488a31fc</tt></a> ContextMarshaler (<a href="https://github.com/99designs/gqlgen/pull/1652">#1652</a>)</summary>

* Add interface and detection for ContextMarshaler

* Test error on float marshalling

* Revert prettier changes

* Rename context test

* Only use the erroring float printer

* Test that context is passed to marshal functions

* Update scalar docs to include the context

* Generate the examples

* Move ContextMarshaller test code to new followschema

* Resolve conflict a little more


* Replicate sclar test for singlefile

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a626d9b47e54e4439976d3ad333dfefe4c2710d6"><tt>a626d9b4</tt></a> Add ICMP to common initialisms (<a href="https://github.com/99designs/gqlgen/pull/1666">#1666</a>)

- <a href="https://github.com/99designs/gqlgen/commit/db4b5eb71b4959c3f8b68086b753ec3c01c5b4c9"><tt>db4b5eb7</tt></a> Merge Inline Fragment Nested Interface Fields (<a href="https://github.com/99designs/gqlgen/pull/1663">#1663</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8b9737179d3ba08405ab6422fb8c4883dd8e720c"><tt>8b973717</tt></a> Update directives doc page (<a href="https://github.com/99designs/gqlgen/pull/1660">#1660</a>)</summary>

* Update directives doc page

* Add back one beloved piece of jargon

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1f500016aedcb9bc35eb9964b55730efe966ef5e"><tt>1f500016</tt></a> Add follow-schema layout for exec (<a href="https://github.com/99designs/gqlgen/pull/1309">#1309</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/1265"> #1265</a>)</summary>

* Define ExecConfig separate from PackageConfig

When support for writing generated code to a directory instead of
a single file is added, ExecConfig will need additional fields
that will not be relevant to other users of PackageConfig.

* Add single-file, follow-schema layouts

When `ExecLayout` is set to `follow-schema`, output generated code to a
directory instead of a single file. Each file in the output directory
will correspond to a single *.graphql schema file (plus a
root!.generated.go file containing top-level definitions that are not
specific to a single schema file).

`ExecLayout` defaults to `single-file`, which is the current behavior, so
this new functionality is opt-in.

These layouts expose similar functionality to the `ResolverLayout`s with
the same name, just applied to `exec` instead of `resolver`.


* Rebase, regenerate

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/129783590c9cb8c709a2b0e9cbb7b69f022ee78a"><tt>12978359</tt></a> Update GQLgen test client to work with multipart form data (take 2) (<a href="https://github.com/99designs/gqlgen/pull/1661">#1661</a>)</summary>

* Update GQLgen test client to work with multipart form data

Update the GQLgen to support multipart form data, like those present
within the fileupload examples.

- Add missing space between "unsupported encoding " and failing
  content-type header error

(cherry picked from commit 101842f73fb79b10c1299bb40506080e08543ec6)

* Add WithFiles client option for fileupload GQLgen client tests

Add a `WithFiles` GQLgen client option to support the fileupload input
within tests, using the core Golang `os` package and File type, which
converts `os.File`s to their appropriate multipart form data within a
request.

- If there are no files this should just simply convert a
  `application/json` Content-Type to supported `multipart/form-data`

(cherry picked from commit 08ef942416c98a2cadf61223308a3ff3c879d1c9)

* Update fileupload test to use GQLgen test client

Update the fileupload test to use the GQLgen test client and `WithFiles`
option to remove the need for `createUploadRequest` helper with raw http
posts

- Fix setting the Content Type by using the appropriate `http` package
  function to dectect it

  + https://godoc.org/net/http#DetectContentType

(cherry picked from commit 5e573d51440eba9d457adb4186772577b28ef085)

* Update WithFiles option test with multipart Reader

(cherry picked from commit 6dfa3cbe0647138e80a59a0c1d55dd9c900f96f2)

* Update file upload tests `WithFiles` option

Update the file upload tests to use the GQL test client and its
`WithFiles` option to remove the need for a custom raw HTTP post request
builder `createUploadRequest`.

- Also update `WithFiles` option to group & map identical files; e.g.

  ```
    { "0": ["variables.req.0.file", "variables.req.1.file"] }
  ```

(cherry picked from commit 486d9f1b2b200701f9ce6b386736a633547c1441)

* Make sure `WithFiles` does not add duplicates to multipart form data

(cherry picked from commit 0c2364d8495553051d97ab805618b006fcd9eddb)

* Fix use of byte vs string in `WithFiles` tests

(cherry picked from commit ba10b5b1c52a74e63e825ee57c235254e8821e0d)

* Fix strict withFiles option test for race conditions

Fix a problem with how strict the test's expected response was for tests
with files in their request, since it always expected a strict order of
files input that is somewhat random or dependent on what OS it is
running the test on and/or race condition

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7435403cf94ce8147fdd9d473a5469d63e7e5b38"><tt>7435403c</tt></a> Adds RootFieldInterceptor to extension interfaces (<a href="https://github.com/99designs/gqlgen/pull/1647">#1647</a>)</summary>

* Adds RootFieldInterceptor to extension interfaces


* Regenerates example folder


* Re-generate after changes

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8b25c9e005c44ae3730eca83445fa7f7223481d1"><tt>8b25c9e0</tt></a> Add a config option to skip running "go mod tidy" on code generation (<a href="https://github.com/99designs/gqlgen/pull/1644">#1644</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/658195b79d8c9072a90419ce3ddccda4e430ebf0"><tt>658195b7</tt></a> Revert "Update GQLgen test client to work with multipart form data (<a href="https://github.com/99designs/gqlgen/pull/1418">#1418</a>)" (<a href="https://github.com/99designs/gqlgen/pull/1659">#1659</a>)</summary>

This reverts commit 1318f12792e86c76a2cdff9132ebac5b3e30e148.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/41c867658a9eacf94b3c682121b03727e18940d5"><tt>41c86765</tt></a> Revert 1595 (<a href="https://github.com/99designs/gqlgen/pull/1658">#1658</a>)

- <a href="https://github.com/99designs/gqlgen/commit/8359f9749e6fd54be20325ff6aafb05503124238"><tt>8359f974</tt></a> Allow custom websocket upgrader (<a href="https://github.com/99designs/gqlgen/pull/1595">#1595</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1318f12792e86c76a2cdff9132ebac5b3e30e148"><tt>1318f127</tt></a> Update GQLgen test client to work with multipart form data (<a href="https://github.com/99designs/gqlgen/pull/1418">#1418</a>)</summary>

* Update GQLgen test client to work with multipart form data

Update the GQLgen to support multipart form data, like those present
within the fileupload examples.

- Add missing space between "unsupported encoding " and failing
  content-type header error

* Add WithFiles client option for fileupload GQLgen client tests

Add a `WithFiles` GQLgen client option to support the fileupload input
within tests, using the core Golang `os` package and File type, which
converts `os.File`s to their appropriate multipart form data within a
request.

- If there are no files this should just simply convert a
  `application/json` Content-Type to supported `multipart/form-data`

* Update fileupload test to use GQLgen test client

Update the fileupload test to use the GQLgen test client and `WithFiles`
option to remove the need for `createUploadRequest` helper with raw http
posts

- Fix setting the Content Type by using the appropriate `http` package
  function to dectect it

  + https://godoc.org/net/http#DetectContentType

* Update WithFiles option test with multipart Reader

* Update file upload tests `WithFiles` option

Update the file upload tests to use the GQL test client and its
`WithFiles` option to remove the need for a custom raw HTTP post request
builder `createUploadRequest`.

- Also update `WithFiles` option to group & map identical files; e.g.

  ```
    { "0": ["variables.req.0.file", "variables.req.1.file"] }
  ```

* Make sure `WithFiles` does not add duplicates to multipart form data

* Fix use of byte vs string in `WithFiles` tests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6758654c4e28dc0589147c9e962c9d4c1fd44705"><tt>6758654c</tt></a> raise panic when nested [@requires](https://github.com/requires) are used on federation (<a href="https://github.com/99designs/gqlgen/pull/1655">#1655</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f6c35be2128d8d0ec6c0c0d63bc0f135292ab5fe"><tt>f6c35be2</tt></a> Add ReplacePlugin option to replace a specific plugin (<a href="https://github.com/99designs/gqlgen/pull/1657">#1657</a>)</summary>

* Add Helper Option for replacing plugins

* Update recipe to use ReplacePlugin instead of NoPlugin and AddPlugin

* fix linting issue on comment

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f8c46600aa005be4d62e52ca6f4a0467480c58c2"><tt>f8c46600</tt></a> fix double indirect bug (<a href="https://github.com/99designs/gqlgen/pull/1604">#1604</a>) (closes <a href="https://github.com/99designs/gqlgen/issues/1587"> #1587</a>)</summary>

* invalid code generated

* update code generation for pointer-to-pointer updating

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/629c91a2dff9982a5c469f25e8076ab7737e167a"><tt>629c91a2</tt></a> remove extra WithOperationContext call (<a href="https://github.com/99designs/gqlgen/pull/1641">#1641</a>)

- <a href="https://github.com/99designs/gqlgen/commit/35199c49ab02648b518d5653ee25eac3e3627602"><tt>35199c49</tt></a> codegen: ensure Elem present before using (<a href="https://github.com/99designs/gqlgen/pull/1317">#1317</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bfea93cdf3594edf0924e5ecd251eea09a1d35cb"><tt>bfea93cd</tt></a> Reload config packages after generating models (<a href="https://github.com/99designs/gqlgen/pull/1491">#1491</a>)</summary>

If models are generated in a package that has already been loaded, and
that package refers to another package that has already been loaded, we
can find ourselves in a position where it appears that a GQL `union` is
not satisfied.

For example, if we have:

```
union Subject = User
```

with this gqlgen.yml in github.com/wendorf/gqlgen-error/gql:

```
schema:
- schema.graphql
exec:
  filename: generated.go
model:

  filename: models_gen.go
models:
  User:
    model: github.com/wendorf/gqlgen-error/gql.User
  Subject:
    model: github.com/wendorf/gqlgen-error/models.Subject
```

Note that our User model is in the github.com/wendorf/gqlgen-error.gql
package, and our models_gen.go will be generated in that same package.

When we try to run gqlgen, we get this error:

```
merging type systems failed: unable to bind to interface: github.com/wendorf/gqlgen-error/gql.User does not satisfy the interface github.com/wendorf/gqlgen-error/models.Subject
```

Digging deeper, it's because we use types.Implements in
codegen/interface.go, which does a shallow object comparison. Because
the type has been reloaded, it refers to a _different_ interface type
object than the one we're comparing against, and get a false negative.

By clearing the package cache and repopulating it, the whole package
cache is generated at the same time, and comparisons across packages
work.

To see a demo of this, check out
https://github.com/wendorf/gqlgen-error and try the following:

1. Checkout the works-with-v0.10.2 branch and `go generate ./...` to see
   that it works
2. Checkout the breaks-with-v0.13.0 branch (or run go get
   to see errors
3. Checkout the works-with-pull-request branch and `go generate ./...`
   to see that it works again. This branch adds a go.mod replace
   directive to use the gqlgen code in this PR.

The demo starts at v0.10.2 since it is the last release without this
problem. https://github.com/99designs/gqlgen/pull/1020 introduces the
code that fails in this scenario.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9e0817cdc7428ea9f3a1542faaca72ce9c5f317c"><tt>9e0817cd</tt></a> Add graphql schema aware field level hook to modelgen (<a href="https://github.com/99designs/gqlgen/pull/1650">#1650</a>)</summary>

* Add ast aware field level hook to modelgen

Currently, the only mechanism for extending the model generation is to use a BuildMutateHook at the end of the model generation process. This can be quite limiting as the hook only has scope of the model build and not the graphql schema which has been parsed.

This change adds a hook at the end of the field creation process which provides access to the parsed graphql type definition and field definition. This allows for more flexibility for example adding additional tags to the model based off custom directives

* Add recipe for using the modelgen FieldMutateHook

* fix goimport linting issue in models_test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/af2ac061db4e08616cecff2ed74649465ee5fc20"><tt>af2ac061</tt></a> handling unconventional naming used in type names (<a href="https://github.com/99designs/gqlgen/pull/1549">#1549</a>)</summary>

* handling unconventional naming used in type names

* Fix merge resolution mistake

* Fix merge resolution mistake

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/393f755421ae42d207655984dbe6b8b990440384"><tt>393f7554</tt></a> add extraTag directive (<a href="https://github.com/99designs/gqlgen/pull/1173">#1173</a>)

- <a href="https://github.com/99designs/gqlgen/commit/fd1bd7c9b3b3804ce1b90b786cd3fb9281918882"><tt>fd1bd7c9</tt></a> adding support for sending extension with gqlgen client (<a href="https://github.com/99designs/gqlgen/pull/1633">#1633</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/589a774290cfaf8f39d6099650e930c6f10cd670"><tt>589a7742</tt></a> Enable lowercase type names in GraphQL schema to properly render (<a href="https://github.com/99designs/gqlgen/pull/1359">#1359</a>)</summary>

The difficulty with lowercased type names is that in go code any lowercased name is not exported.
This change makes the names title case for go code while preserving the proper case when interacting with the GraphQL schema.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/50f6a2aa603842fcdc158ab135fa117d1716d7e2"><tt>50f6a2aa</tt></a> Fixes <a href="https://github.com/99designs/gqlgen/pull/1653">#1653</a>: update docs and wrap error if not *gqlerror.Error (<a href="https://github.com/99designs/gqlgen/pull/1654">#1654</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7081dedb0efc6ed650118c7fce65ca3bdb33b8de"><tt>7081dedb</tt></a> Bump tmpl from 1.0.4 to 1.0.5 in /integration (<a href="https://github.com/99designs/gqlgen/pull/1627">#1627</a>)</summary>

Bumps [tmpl](https://github.com/daaku/nodejs-tmpl) from 1.0.4 to 1.0.5.
- [Release notes](https://github.com/daaku/nodejs-tmpl/releases)
- [Commits](https://github.com/daaku/nodejs-tmpl/commits/v1.0.5)

---
updated-dependencies:
- dependency-name: tmpl
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5287e4e5f30548d233f58111d128787c673c7f01"><tt>5287e4e5</tt></a> Add QR and KVK to common initialisms (<a href="https://github.com/99designs/gqlgen/pull/1419">#1419</a>)</summary>

* Add QR and KVK to common initialisms

* Update templates.go

* Sort commonInitialisms

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f9df1a46601a23df87d364b33cc9c5564d77edd8"><tt>f9df1a46</tt></a> Update time format for `Time` scalar (<a href="https://github.com/99designs/gqlgen/pull/1648">#1648</a>)</summary>

* Use more precise time format

* update test

* update docs

* Apply suggestions from code review

* Update scalars.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/77c757f0cf9f18de02f3b9e6a235a51dd2c75259"><tt>77c757f0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1640">#1640</a> from minus7/master</summary>

Fix example run instructions

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e60dc7af373aeca15f48df234c1cb90d4909db5d"><tt>e60dc7af</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1619">#1619</a> from Khan/benkraft.mod-tidy-stdout</summary>

Forward `go mod tidy` stdout/stderr

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0c63f1d10f508f528037c94c5cdb9f29af890098"><tt>0c63f1d1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1515">#1515</a> from OpenSourceProjects/time</summary>

Marshaling & Unmarshaling time return initial value

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a3d9e8ce9689533ab8c3ab4b3b4cd22df3cbfa03"><tt>a3d9e8ce</tt></a> Remove redundant favicon (<a href="https://github.com/99designs/gqlgen/pull/1638">#1638</a>)

- <a href="https://github.com/99designs/gqlgen/commit/210c1aa6edf8b34d6e2f35e8c66b3d404cf7e8eb"><tt>210c1aa6</tt></a> Appropriately Handle Falsy Default Field Values (<a href="https://github.com/99designs/gqlgen/pull/1623">#1623</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/47ce074a3c30a981bbeef8f3465fca4330aba783"><tt>47ce074a</tt></a> Fix example run instructions (closes <a href="https://github.com/99designs/gqlgen/issues/1607"> #1607</a>)</summary>

Making ./example a separate Go module [1] broke the `go run` invocations
listed in a few example readmes [2]. Using relative paths from the
respective example directory should be clear enough.

[2]:
example/todo/server/server.go:10:2: no required module provides package github.com/99designs/gqlgen/example/todo; to add it:
	go get github.com/99designs/gqlgen/example/todo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1a0b19feff6f02d2af6631c9d847bc243f8ede39"><tt>1a0b19fe</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d999828375978c666728167f75a208fd727b4b15"><tt>d9998283</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1628">#1628</a> from robertmarsal/patch-1</summary>

Fix typo in the getting-started docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f93f73ac395209cd158862b86b2a0c136bc3e11b"><tt>f93f73ac</tt></a> Fix typo in the getting-started docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2f6919ff74dd165d58c0c2039e3fb1fc1f72b598"><tt>2f6919ff</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1624">#1624</a> from FlymeDllVa/master</summary>

Update disabling Introspection

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c53bc0e53cacb9f9930ec3125388b85820241a27"><tt>c53bc0e5</tt></a> Update disabling Introspection

- <a href="https://github.com/99designs/gqlgen/commit/880cd73dbe0d602168639a7d3f59638473a4a91c"><tt>880cd73d</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/eec81df05e18dd9ac977b5a661964df212bf4627"><tt>eec81df0</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/43b56cbaf3f1de1d1ad379055ab1de157592cf38"><tt>43b56cba</tt></a> Forward `go mod tidy` stdout/stderr</summary>

This is a command that can fail (in my case I think for stupid reasons
in a hell of my own construction, but nonetheless).  Right now we just
get
```
$ go run github.com/Khan/webapp/dev/cmd/gqlgen
tidy failed: go mod tidy failed: exit status 1
exit status 3
```
which is not the most informative.  Now, instead, we'll forward its
output to our own stdout/stderr rather than devnull.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ce7a8ee469108e1c2dd62511249199312b32729a"><tt>ce7a8ee4</tt></a> Fix link in docs

- <a href="https://github.com/99designs/gqlgen/commit/488cf7e8180e653c7a085c137d14734c5897393b"><tt>488cf7e8</tt></a> Update docs/content/getting-started.md

- <a href="https://github.com/99designs/gqlgen/commit/73809f6912ed01c8ea18dbb5d7ea98e803ddb9c7"><tt>73809f69</tt></a> Update getting started

- <a href="https://github.com/99designs/gqlgen/commit/b938e55811966ef5ff96cc911597327af92c493c"><tt>b938e558</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/cacd49a6d0421093bdb67d17551ea0323ceb438a"><tt>cacd49a6</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7d549d6476853a33aaacab76c3d954bed9d6f0cd"><tt>7d549d64</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1617">#1617</a> from 99designs/update-docs-for-go1.17</summary>

Update docs for getting started

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5c52f27c49d65d97a05bba14fc12ce8873d2b959"><tt>5c52f27c</tt></a> Update docs for getting started

- <a href="https://github.com/99designs/gqlgen/commit/41d6926f38922cd0a01d93f30ac8702d9f9bcf48"><tt>41d6926f</tt></a> Replace gitter with discord in contributing.md

- <a href="https://github.com/99designs/gqlgen/commit/24d4edcf128f7fa327d97fd85d087fef2e230943"><tt>24d4edcf</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/2272e05bc00cc9f38e6d3179a8981651812acbba"><tt>2272e05b</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ef4d4a38a229a485bec4f111a1112b47c638382e"><tt>ef4d4a38</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1614">#1614</a> from 99designs/go-1.16</summary>

Also test against 1.16

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/00ed6fb1a74b1c94d195b38025b1721f9c77db90"><tt>00ed6fb1</tt></a> Also test against 1.16

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/473f0671b5d9c3a0bb9afe6c2de2b4f10d9eeef6"><tt>473f0671</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1613">#1613</a> from 99designs/bump-non-module-deps</summary>

Clean up non-module deps

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6960c0c2adbe9e2359e0d7cb332d1eaceccb6b6f"><tt>6960c0c2</tt></a> Bump non-module deps

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bf9b34aae2eb079913aacbb75fa8b3dff7c45c32"><tt>bf9b34aa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1612">#1612</a> from 99designs/update-linter</summary>

Update golangci linter

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/85e7a4a0aae4c80c8522350d17400184fac48882"><tt>85e7a4a0</tt></a> Linting fixes

- <a href="https://github.com/99designs/gqlgen/commit/777dabde381c1c4b1b6bb0316658f65cce22c654"><tt>777dabde</tt></a> Update the linter

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/85dd47bb8ca9547ebc3530f441aab8a99e16b5a7"><tt>85dd47bb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1607">#1607</a> from 99designs/example-module</summary>

[POC/RFC] Split examples into separate go module

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f93fb2489285eef0542c4aa7b11341a8479b606a"><tt>f93fb248</tt></a> Split examples into separate go module

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/890f5f66fb2059ccd760bb7dd708b04166ed274d"><tt>890f5f66</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1610">#1610</a> from 99designs/go-1.17</summary>

Update to go 1.17

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9162c53fc32471c96d874764b47bf71f0c493fce"><tt>9162c53f</tt></a> Fix newlines in error messages

- <a href="https://github.com/99designs/gqlgen/commit/f67a5b2611ca89dfb0ffd3dd4c72cebb8f2532ef"><tt>f67a5b26</tt></a> Update github.com/urfave/cli/v2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1116ea6cdc1f3c66af5e214365af7ed67d464a52"><tt>1116ea6c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1608">#1608</a> from jjmengze/patch-1</summary>

fix Options response header

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/71e5784352a3c6d7b91bf9f2307fcebc959ab8b0"><tt>71e57843</tt></a> Simplify init

- <a href="https://github.com/99designs/gqlgen/commit/a8903ca2aca700217fe29c2b8c262b6ee45959fb"><tt>a8903ca2</tt></a> Wrap errors

- <a href="https://github.com/99designs/gqlgen/commit/a644175b8f80ba1fbc220799b4cd9ce257ccc7ac"><tt>a644175b</tt></a> Update error checks for go 1.17

- <a href="https://github.com/99designs/gqlgen/commit/c6b9f2926b14c0367ffb82db4e2c250c4fc7aab2"><tt>c6b9f292</tt></a> go mod tidy

- <a href="https://github.com/99designs/gqlgen/commit/1c63cfff8f2943d2ab3316f850b89bba83548c05"><tt>1c63cfff</tt></a> Add missing model package file

- <a href="https://github.com/99designs/gqlgen/commit/59da23feb5e135b3c2c0ce976b0214643147e85c"><tt>59da23fe</tt></a> Create a temporary file on init so go recognises the directory as a package

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/682a7d662bd2491213c45b8556461533cee56997"><tt>682a7d66</tt></a> fix Options response header</summary>

operatee the header of ResponseWriter should before WriteHeader called

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ed8054b054c4e7dae1c2ad12d86a8d4d26c194e9"><tt>ed8054b0</tt></a> Update to a post-release version

- <a href="https://github.com/99designs/gqlgen/commit/5216db5849eb5298722ee233cc98fa1e85db3237"><tt>5216db58</tt></a> Fix TestAutobinding test failure by checking the module

- <a href="https://github.com/99designs/gqlgen/commit/90c5eb59b0fde89eb0c42e5b9c5ad276a8335f8f"><tt>90c5eb59</tt></a> go generate

- <a href="https://github.com/99designs/gqlgen/commit/402f44950b4694a3a5fd65d3225e06ccd76fbf9d"><tt>402f4495</tt></a> go fmt

- <a href="https://github.com/99designs/gqlgen/commit/10bb1ef262af49031c3d5224eb094586e6fd8083"><tt>10bb1ef2</tt></a> Go mod tidy

- <a href="https://github.com/99designs/gqlgen/commit/ed210385722431111b8842ec7da7b42051b77d91"><tt>ed210385</tt></a> Update to go 1.17

- <a href="https://github.com/99designs/gqlgen/commit/5c7acc1bc8198ed57a38867ae0c452987da91911"><tt>5c7acc1b</tt></a> Fix imports

- <a href="https://github.com/99designs/gqlgen/commit/d747387036fb300fc384fa6f17a76f6547dc6a1b"><tt>d7473870</tt></a> Update plugin/servergen/server.go

- <a href="https://github.com/99designs/gqlgen/commit/a6c6de6b73741caa1ac02cf8adea79afcd4ed78b"><tt>a6c6de6b</tt></a> Update plugin/resolvergen/resolver.go

- <a href="https://github.com/99designs/gqlgen/commit/de7d19c812bed562d2841c6c08fb080d706195d1"><tt>de7d19c8</tt></a> Update codegen/config/config_test.go

- <a href="https://github.com/99designs/gqlgen/commit/60d80d4aee61337793d2cade8a7ab35c2892613a"><tt>60d80d4a</tt></a> Update cmd/gen.go

- <a href="https://github.com/99designs/gqlgen/commit/a991e3e73ec4d624f6b23124d83198ce51af8ae3"><tt>a991e3e7</tt></a> Update errors to use go1.13 semantics

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8f179be920401bf993630693b0586cfc51dbdd04"><tt>8f179be9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1581">#1581</a> from tsh96/master</summary>

Bypass complexity limit on __Schema queries.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5048f9927e54828b07ea47cd79d1b30b3858d320"><tt>5048f992</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1525">#1525</a> from Code-Hex/fix/support-input-object</summary>

support input object directive

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1e2b303a8cc8a596fc24cb79b560c33bec2c9ad6"><tt>1e2b303a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1526">#1526</a> from epulze/fix/allow-more-types</summary>

allow more than 10 different import sources with types

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e7df3e5c7d33dd485ec928c4be07f421423c722b"><tt>e7df3e5c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1405">#1405</a> from alexsn/subsciption-complete-on-panic</summary>

subscriptions: send complete message on resolver panic

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/06e4fe8810d14649a5b6e2f9c1482eff28c86ddb"><tt>06e4fe88</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1529">#1529</a> from mathieupost/master</summary>

Return type loading errors in config.Binder.FindObject

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a557c90cd741f808f4065e9846b0c59ba1d29f9b"><tt>a557c90c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1340">#1340</a> from bickyeric/master</summary>

serialize ID just like String

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/522cab59d1ecb9813610013a60596da2edf91f33"><tt>522cab59</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1285">#1285</a> from Khan/benkraft.federation</summary>

Resolve requests for federation entities in parallel

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5adb73bbba5375f07cd21d1fe498c6a252b6f933"><tt>5adb73bb</tt></a> add bypass __schema field test case

- <a href="https://github.com/99designs/gqlgen/commit/54cef3ddcdd3c5e106f8347ccd8afd9bbb8bdb44"><tt>54cef3dd</tt></a> Bypass complexity limit on __Schema queries.

- <a href="https://github.com/99designs/gqlgen/commit/f0ccab79549f184b1225dbc946d703400673aa5b"><tt>f0ccab79</tt></a> Return type loading errors in config.Binder.FindObject

- <a href="https://github.com/99designs/gqlgen/commit/91b54787166e1dd26e2f95493e02a43130950105"><tt>91b54787</tt></a> generated go code

- <a href="https://github.com/99designs/gqlgen/commit/1efc152e4f9399a908620d1d4c1198b9662cb181"><tt>1efc152e</tt></a> supported INPUT_OBJECT directive

- <a href="https://github.com/99designs/gqlgen/commit/e82b401ddd97fe807a5066d3702b391ec317a358"><tt>e82b401d</tt></a> allow more than 10 different import sources with types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/481a4e44cbc114382dee087feb9cf18c15907d4a"><tt>481a4e44</tt></a> Marshaling & Unmarshaling time return initial value</summary>

There was a lack of symmetry that would prevent times for being
symmetrical. That is because time.Parse actually parses an RFC3339Nano
implicitly, thereby allowing nanosecond resolution on unmarshaling a
time. Therefore we now marshal into nanoseconds, getting more
information into GraphQL times when querying for a time, and restoring
the symmetry

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/956531936e9a7a0d9012c246153531a9bcf3c7bd"><tt>95653193</tt></a> Resolve requests for federation entities in parallel (closes <a href="https://github.com/99designs/gqlgen/issues/1278"> #1278</a>)</summary>

In apollo federation, we may be asked for data about a list of entities.
These can typically be resolved in parallel, just as with sibling fields
in ordinary GraphQL queries.  Now we do!

I also changed the behavior such that if one lookup fails, we don't
cancel the others.  This is more consistent with the behavior of other
resolvers, and is more natural now that they execute in parallel.  This,
plus panic handling, required a little refactoring.

The examples probably give the clearest picture of the changes. (And the
clearest test; the changed functionality is already exercised by
`integration-test.js` as watching the test server logs will attest.)

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f00e2c3f3d70a58ae85ab05d5f2b8adf543a306d"><tt>f00e2c3f</tt></a> subscriptions: send complete message on resolver panic

- <a href="https://github.com/99designs/gqlgen/commit/fa371b9bb76e478dabc649564e1f465c45057f72"><tt>fa371b9b</tt></a> serialize ID just like String

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.14.0"></a>
## [v0.14.0](https://github.com/99designs/gqlgen/compare/v0.13.0...v0.14.0) - 2021-09-08
- <a href="https://github.com/99designs/gqlgen/commit/56451d92d626be6d15317b44e448c857297ddb68"><tt>56451d92</tt></a> release v0.14.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8e97969b06e6d63160f1252a92bb36972c55a0b6"><tt>8e97969b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1358">#1358</a> from mtsmfm/patch-1</summary>

Create package declaration to run dataloaden

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b978593ca852ca1968a93501eaebec7bdc7fd359"><tt>b978593c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1387">#1387</a> from Khan/benkraft.config</summary>

codegen/config: Add a new API to finish an already-validated config

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/71507dfc1be232bed3d4d08f63d61f8c8a9cdd77"><tt>71507dfc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1408">#1408</a> from max107/patch-1</summary>

int64 support graphql/string.go

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/23577b696e36bf3342c685b20ec48b8484158518"><tt>23577b69</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1460">#1460</a> from snxk/edit-docs-recipe-gin</summary>

Edited the Gin-Gonic Recipe Docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/db6154b9eac60bf3ec82b959c0e6bfe1e79f0bb8"><tt>db6154b9</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cecda16053ffe9f2b689c6e7da987eabe5c2515f"><tt>cecda160</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1464">#1464</a> from frederikhors/patch-1</summary>

Add goreportcard badge

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cc957171fc77df6bfb0749ddd66e1d8d9ca24afe"><tt>cc957171</tt></a> Merge branch 'master' into patch-1

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/023f66df41a1761ab36eed6c68733e1dd85e5608"><tt>023f66df</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1465">#1465</a> from frederikhors/patch-2</summary>

Add coveralls badge

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/50c2028a9574c89f59d7720b84bf720e07a6a974"><tt>50c2028a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1497">#1497</a> from polytomic/stable-introspection</summary>

Return introspection document in stable order

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a0232dd21c2fe51d8d52e6541d7a01d50bbaab4d"><tt>a0232dd2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1603">#1603</a> from 99designs/dependabot/npm_and_yarn/integration/normalize-url-4.5.1</summary>

Bump normalize-url from 4.5.0 to 4.5.1 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4e059eba3d28e0b27445ba71ffe670852c84c096"><tt>4e059eba</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1602">#1602</a> from 99designs/dependabot/npm_and_yarn/integration/ini-1.3.8</summary>

Bump ini from 1.3.5 to 1.3.8 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/43705d459a1340ef0e7c4ee46af97cec592c976c"><tt>43705d45</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1601">#1601</a> from 99designs/dependabot/npm_and_yarn/integration/y18n-3.2.2</summary>

Bump y18n from 3.2.1 to 3.2.2 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1f2465c6d29a731ad1b9dd5b08d57da37ce14043"><tt>1f2465c6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1600">#1600</a> from 99designs/dependabot/npm_and_yarn/integration/browserslist-4.17.0</summary>

Bump browserslist from 4.14.0 to 4.17.0 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bbdebd4c55194a2e18428a33546adbecea99b05d"><tt>bbdebd4c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1599">#1599</a> from 99designs/dependabot/npm_and_yarn/integration/hosted-git-info-2.8.9</summary>

Bump hosted-git-info from 2.8.5 to 2.8.9 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/900a37af53ee048bd10a755d227bebf11b0bc53f"><tt>900a37af</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1598">#1598</a> from 99designs/dependabot/npm_and_yarn/integration/node-fetch-2.6.1</summary>

Bump node-fetch from 2.6.0 to 2.6.1 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9d334cdd222c4616165e0ef168086c4d178d4313"><tt>9d334cdd</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1597">#1597</a> from 99designs/dependabot/npm_and_yarn/integration/ws-7.4.6</summary>

Bump ws from 7.3.1 to 7.4.6 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/56181e8abe857e229a3e63e8d634582647480681"><tt>56181e8a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1365">#1365</a> from frederikhors/add-uint,-uint64,-uint32-types-in-graphql</summary>

add uint, uint64, uint32 types in graphql pkg

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fd133c0b7a2d552e73da63180e3af2e4bf4aa434"><tt>fd133c0b</tt></a> Bump normalize-url from 4.5.0 to 4.5.1 in /integration</summary>

Bumps [normalize-url](https://github.com/sindresorhus/normalize-url) from 4.5.0 to 4.5.1.
- [Release notes](https://github.com/sindresorhus/normalize-url/releases)
- [Commits](https://github.com/sindresorhus/normalize-url/commits)

---
updated-dependencies:
- dependency-name: normalize-url
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/24d8c703c122cd007fe2eb81457a2eae89b49be8"><tt>24d8c703</tt></a> Bump ini from 1.3.5 to 1.3.8 in /integration</summary>

Bumps [ini](https://github.com/isaacs/ini) from 1.3.5 to 1.3.8.
- [Release notes](https://github.com/isaacs/ini/releases)
- [Commits](https://github.com/isaacs/ini/compare/v1.3.5...v1.3.8)

---
updated-dependencies:
- dependency-name: ini
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/de89d3a6f28cffb44704f6f9b6876c8c81d9f7aa"><tt>de89d3a6</tt></a> Bump y18n from 3.2.1 to 3.2.2 in /integration</summary>

Bumps [y18n](https://github.com/yargs/y18n) from 3.2.1 to 3.2.2.
- [Release notes](https://github.com/yargs/y18n/releases)
- [Changelog](https://github.com/yargs/y18n/blob/master/CHANGELOG.md)
- [Commits](https://github.com/yargs/y18n/commits)

---
updated-dependencies:
- dependency-name: y18n
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/13db61111eae250a02ead0cd9faa456d98dc007b"><tt>13db6111</tt></a> Bump browserslist from 4.14.0 to 4.17.0 in /integration</summary>

Bumps [browserslist](https://github.com/browserslist/browserslist) from 4.14.0 to 4.17.0.
- [Release notes](https://github.com/browserslist/browserslist/releases)
- [Changelog](https://github.com/browserslist/browserslist/blob/main/CHANGELOG.md)
- [Commits](https://github.com/browserslist/browserslist/compare/4.14.0...4.17.0)

---
updated-dependencies:
- dependency-name: browserslist
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/94e9406e93a5afee6f193fc936f200287c4f5847"><tt>94e9406e</tt></a> Bump hosted-git-info from 2.8.5 to 2.8.9 in /integration</summary>

Bumps [hosted-git-info](https://github.com/npm/hosted-git-info) from 2.8.5 to 2.8.9.
- [Release notes](https://github.com/npm/hosted-git-info/releases)
- [Changelog](https://github.com/npm/hosted-git-info/blob/v2.8.9/CHANGELOG.md)
- [Commits](https://github.com/npm/hosted-git-info/compare/v2.8.5...v2.8.9)

---
updated-dependencies:
- dependency-name: hosted-git-info
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/36be94fff25cfe702f6b93a7fd595972ceb41a4b"><tt>36be94ff</tt></a> Bump node-fetch from 2.6.0 to 2.6.1 in /integration</summary>

Bumps [node-fetch](https://github.com/node-fetch/node-fetch) from 2.6.0 to 2.6.1.
- [Release notes](https://github.com/node-fetch/node-fetch/releases)
- [Changelog](https://github.com/node-fetch/node-fetch/blob/main/docs/CHANGELOG.md)
- [Commits](https://github.com/node-fetch/node-fetch/compare/v2.6.0...v2.6.1)

---
updated-dependencies:
- dependency-name: node-fetch
  dependency-type: direct:development
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/721158f3cdb00744c33380c05ab020f31b894325"><tt>721158f3</tt></a> Bump ws from 7.3.1 to 7.4.6 in /integration</summary>

Bumps [ws](https://github.com/websockets/ws) from 7.3.1 to 7.4.6.
- [Release notes](https://github.com/websockets/ws/releases)
- [Commits](https://github.com/websockets/ws/compare/7.3.1...7.4.6)

---
updated-dependencies:
- dependency-name: ws
  dependency-type: direct:development
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2b3b721258bb22d0da26790a9383047cf1ef444c"><tt>2b3b7212</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1594">#1594</a> from 99designs/dependabot/npm_and_yarn/integration/tar-6.1.11</summary>

Bump tar from 6.0.5 to 6.1.11 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5b43833db94d42332553ba79103f3d054c461e62"><tt>5b43833d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1582">#1582</a> from 99designs/dependabot/npm_and_yarn/integration/path-parse-1.0.7</summary>

Bump path-parse from 1.0.6 to 1.0.7 in /integration

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/55b028cadc5e16421f0cb06ec9ffa94febee72a4"><tt>55b028ca</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1584">#1584</a> from nullism/patch-1</summary>

Fix spaces -> tabs typo in authentication.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/edf630a3da614949b9f0246299b08a61d25f6635"><tt>edf630a3</tt></a> Bump tar from 6.0.5 to 6.1.11 in /integration</summary>

Bumps [tar](https://github.com/npm/node-tar) from 6.0.5 to 6.1.11.
- [Release notes](https://github.com/npm/node-tar/releases)
- [Changelog](https://github.com/npm/node-tar/blob/main/CHANGELOG.md)
- [Commits](https://github.com/npm/node-tar/compare/v6.0.5...v6.1.11)

---
updated-dependencies:
- dependency-name: tar
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/29133c1152b4d89ab897bf75e0114fffad32351e"><tt>29133c11</tt></a> Fix spaces -> tabs typo in authentication.md</summary>

The indentation here was supposed to be a tab rather than spaces so the readme was off.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/01b25c5534453574296c52065c04d1923b3607c7"><tt>01b25c55</tt></a> Bump path-parse from 1.0.6 to 1.0.7 in /integration</summary>

Bumps [path-parse](https://github.com/jbgutierrez/path-parse) from 1.0.6 to 1.0.7.
- [Release notes](https://github.com/jbgutierrez/path-parse/releases)
- [Commits](https://github.com/jbgutierrez/path-parse/commits/v1.0.7)

---
updated-dependencies:
- dependency-name: path-parse
  dependency-type: indirect
...

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9a214e80158b78443cfb53bb10059df9c36d352e"><tt>9a214e80</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1451">#1451</a> from sanjeevchopra/patch-1</summary>

doc only change: updated sample code for disabling introspection

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/011974377ac4437c060d963c97f9c04f1fd1bfae"><tt>01197437</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1417">#1417</a> from RicCu/patch-1</summary>

Use mutation instead of query in 'Changesets' doc example

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e3293b53d07145d4932a7b05fea534763cc8af12"><tt>e3293b53</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1444">#1444</a> from lisowskibraeden/patch-1</summary>

Update cors.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a4d6785526f2423e5cd7d2a35bf2f6e68ab66bf7"><tt>a4d67855</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1517">#1517</a> from ShivangGoswami/patch-1</summary>

Update apq.md function definition mismatch

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/eb36f04ffde9a706467839b0aadb5671e4ed16a9"><tt>eb36f04f</tt></a> Return introspection document in stable order</summary>

This avoids spurious changes when generating client code using
something like graphql-codegen.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7e38dd46943cc103a82c6bca0c2510e5d1291edc"><tt>7e38dd46</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1568">#1568</a> from DanyHenriquez/patch-1</summary>

Update apq.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/88f2b8a77680f49b07238cdafc03d429e5fb75b7"><tt>88f2b8a7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1572">#1572</a> from talhaguy/dataloaders-doc-casing</summary>

Correct minor casing issue

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/be9a0791a9217a05e105914dbae11e29a651ea54"><tt>be9a0791</tt></a> Update apq.md

- <a href="https://github.com/99designs/gqlgen/commit/3e45ddc151232c5d05eb719f4722a9306f06afa1"><tt>3e45ddc1</tt></a> Correct minor casing issue

- <a href="https://github.com/99designs/gqlgen/commit/145101e439f5460cbe7e85f8618e4de74104b676"><tt>145101e4</tt></a> Update apq.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/843edd9ea507bcf50b02cb43d0543f3fdf0ae875"><tt>843edd9e</tt></a> Update apq.md function definition mismatch</summary>

line 67:  cache, err := NewCache(cfg.RedisAddress, 24*time.Hour)
line 41: func NewCache(redisAddress string, password string,ttl time.Duration) (*Cache, error)

either password should be removed from 41 or added in line 67
Proposed the first one for now.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5ad012e3d7be1127706b9c8a3da0378df3a98ec1"><tt>5ad012e3</tt></a> Revert "Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1511">#1511</a> from a8m/a8m/restore-cwd"</summary>

This reverts commit f4bf1f591b6a3884041876deb64ce0dd70c3c883, reversing
changes made to 3f68ea27a1a9fea2064caf877f7e24d00aa439e6.

Reverting this because it will break existing setups, moving where
generated files get put.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bb59cc43aa5bae2595ec823f8d7e67369e990082"><tt>bb59cc43</tt></a> Add a CHANGELOG.md (<a href="https://github.com/99designs/gqlgen/pull/1512">#1512</a>)

- <a href="https://github.com/99designs/gqlgen/commit/058a365a3608a0d8e9704ee8715eb6c70e7cc902"><tt>058a365a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1456">#1456</a> from skaji/issue-1455

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bf2fdf4401b3c77a4d032572f641787eb99e8b71"><tt>bf2fdf44</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1514">#1514</a> from 99designs/bump-gqlparser</summary>

Bump gqlparser to v2.2.0

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4e881981de33f2e2bf1fef11c2bf833995f60719"><tt>4e881981</tt></a> Bump to gqlparser v2.2.0

- <a href="https://github.com/99designs/gqlgen/commit/1d768a29c960df5d54f6e675e0338619d5b04bfd"><tt>1d768a29</tt></a> Add test covering single element -> slice coercion

- <a href="https://github.com/99designs/gqlgen/commit/f57d1a0285eebce853a6a008da3c9c7b4eb77c57"><tt>f57d1a02</tt></a> Bump gqlparser to master & support repeated directives

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f4bf1f591b6a3884041876deb64ce0dd70c3c883"><tt>f4bf1f59</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1511">#1511</a> from a8m/a8m/restore-cwd</summary>

codegen/config: restore current working directory after changing it

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3f68ea27a1a9fea2064caf877f7e24d00aa439e6"><tt>3f68ea27</tt></a> Special handling for pointers to slices (<a href="https://github.com/99designs/gqlgen/pull/1363">#1363</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c920bdebb02f1fc4c406b9f36a63114556303657"><tt>c920bdeb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1449">#1449</a> from steebchen/feat-prisma-compat</summary>

feat(codegen): handle (v, ok) methods

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3cfc5b14d89ff8250656e40f967d1b0fae9de374"><tt>3cfc5b14</tt></a> codegen/config: restore current working directory after changing it</summary>

Before this commit, a call to config.LoadConfigFromDefaultLocations changed
the working directory to the directory that contains the gqlgen config
file.

This commit changes the implementation to restore the working directory
after loading the config.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/35b80a72f6cdae48cf98c68128d96d9d70e5f756"><tt>35b80a72</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1495">#1495</a> from Niennienzz/improve-apq-doc</summary>

Update apq.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/463debae6b4eb068aff0b882e6ea292bfac0fae2"><tt>463debae</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1503">#1503</a> from nana4gonta/resolve-vulnerability</summary>

Resolve indirect dependency vulnerability in example

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/29e7bccbf7fcb7d8d7f8c47cabd7abdc542cdcc6"><tt>29e7bccb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1501">#1501</a> from 99designs/fix-init-1.16</summary>

Run go mod tidy after code generation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9a4c80abc704d77c9471f5b7ee47d64cbced0348"><tt>9a4c80ab</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1502">#1502</a> from 99designs/rm-chi</summary>

Remove chi from dataloader example

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5f21f9d9ecdedca84810d7fda605c6eddd1f2335"><tt>5f21f9d9</tt></a> Remove chi from dataloader example

- <a href="https://github.com/99designs/gqlgen/commit/e02db808a857ce1ec31861b3b0b54fa4d5cb85e6"><tt>e02db808</tt></a> Run go mod tidy after code generation

- <a href="https://github.com/99designs/gqlgen/commit/8c3e64e1965081ff07bbe9353dd9246817232887"><tt>8c3e64e1</tt></a> Improve APQ documentation

- <a href="https://github.com/99designs/gqlgen/commit/03b57f3e01f34504261550aacb0b08f64843b6ad"><tt>03b57f3e</tt></a> Run go mod tidy

- <a href="https://github.com/99designs/gqlgen/commit/54e387c45e97e7b7922f06baf1c6e57dd5a7ff2e"><tt>54e387c4</tt></a> Resolve indirect dependency vulnerability in example

- <a href="https://github.com/99designs/gqlgen/commit/7985db44855b160c1f2552bedbfed5bc150fc840"><tt>7985db44</tt></a> Mention math.rand for the todo ID (<a href="https://github.com/99designs/gqlgen/pull/1489">#1489</a>)

- <a href="https://github.com/99designs/gqlgen/commit/b995f7f1fa2e18b4016d167739213ec5de95a053"><tt>b995f7f1</tt></a> Make spacing consistent (<a href="https://github.com/99designs/gqlgen/pull/1488">#1488</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/52ded95125beecffe5ad61d37ef942fbac2d726f"><tt>52ded951</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1459">#1459</a> from aaronArinder/getting-started-server-section</summary>

getting started: make running server own section

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/82a8e1bf39aec5225e05deb8e026083d859d50ef"><tt>82a8e1bf</tt></a> Make it clearer what happened on init. (<a href="https://github.com/99designs/gqlgen/pull/1487">#1487</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7258af5f837802cd1673afd0778ee7a76b8c2471"><tt>7258af5f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1458">#1458</a> from aaronArinder/getting-started-wording</summary>

getting started: making the resolver fn section clearer

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4fead4895bc44aff95fd06a4a5a3aa4b184cc2ff"><tt>4fead489</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1452">#1452</a> from fmyd/fix/formatted-query-indent</summary>

prettified some indentation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/58e3225ed3100371d286869e9a2a4b19ec9810e6"><tt>58e3225e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1480">#1480</a> from wilhelmeek/double-bubble</summary>

Bubble Null from List Element to Nearest Nullable Ancestor

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1fac78e9b4d68a76d3ae2fc0b980e7569cc3eb3c"><tt>1fac78e9</tt></a> Add test case for nullable field

- <a href="https://github.com/99designs/gqlgen/commit/469e31bddf27395e5d124694e683cecc44fa00d3"><tt>469e31bd</tt></a> Fix bad test case

- <a href="https://github.com/99designs/gqlgen/commit/635b1aef316c528b6f17416057d41a0948b42a2d"><tt>635b1aef</tt></a> Add Test Case

- <a href="https://github.com/99designs/gqlgen/commit/0b5da15cd87f315af0bc851a50f21362e4312002"><tt>0b5da15c</tt></a> Check in generated code

- <a href="https://github.com/99designs/gqlgen/commit/55b774ba48146540bdef95d8cc027998eca7fd13"><tt>55b774ba</tt></a> Fix type ref

- <a href="https://github.com/99designs/gqlgen/commit/45903a6597846c5ae71d82f156fc3ee3b743ec75"><tt>45903a65</tt></a> Handle nillable list elements

- <a href="https://github.com/99designs/gqlgen/commit/c4bf36c5bd94b64e8d13060a77a7b8ac8050b794"><tt>c4bf36c5</tt></a> Add coveralls badge

- <a href="https://github.com/99designs/gqlgen/commit/269a58ad547f61b4ceaaa6586ed7898cafc041c1"><tt>269a58ad</tt></a> Add goreportcard badge

- <a href="https://github.com/99designs/gqlgen/commit/971da82c8e3d1cf7cca31bb9cfff91cab2a460d3"><tt>971da82c</tt></a> Updated gin.md

- <a href="https://github.com/99designs/gqlgen/commit/41ad51ceefc190b70c4b5fa77c77641ede1d7281"><tt>41ad51ce</tt></a> Edited the Gin-Gonic Recipe Docs

- <a href="https://github.com/99designs/gqlgen/commit/67e652ad974418a9f6fa9cfa4c97eebc3db910bf"><tt>67e652ad</tt></a> getting started: separate example mutation/query

- <a href="https://github.com/99designs/gqlgen/commit/31d339ab390a2c4119fec42c54edabebd96ef730"><tt>31d339ab</tt></a> getting started: make running server own section

- <a href="https://github.com/99designs/gqlgen/commit/aa531ed87f327f0e03d243744a5f5e810d5c1230"><tt>aa531ed8</tt></a> getting started: more wording updates

- <a href="https://github.com/99designs/gqlgen/commit/5b2531aee84fa4f971ea83c083f7113b2d6b7c6e"><tt>5b2531ae</tt></a> getting started: wording update

- <a href="https://github.com/99designs/gqlgen/commit/ada1b928096db2d4cff8d476e62d8a84c41da47e"><tt>ada1b928</tt></a> getting started: updating wording around implementing unimpl fns

- <a href="https://github.com/99designs/gqlgen/commit/23eec79139fd4735d50d397a081edfefad09fd27"><tt>23eec791</tt></a> go generate ./...

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/18678b15ecbcf6075356623fbc0902606e440513"><tt>18678b15</tt></a> Fix data race</summary>

The argument of unmarshalInput may be the same for concurrent use if it pass as graphql "variables".
So we have to copy it before setting default values

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/02b140038d1f192af2ff2cc05a08a691b300ec94"><tt>02b14003</tt></a> fomatted query indent

- <a href="https://github.com/99designs/gqlgen/commit/0e9d9c3a9d072c6e2262bf742705d723be8d2508"><tt>0e9d9c3a</tt></a> updated sample code for disabling introspection

- <a href="https://github.com/99designs/gqlgen/commit/478c3f08b20fadd31c538ddd90bb7e88a4e2c1a9"><tt>478c3f08</tt></a> feat(codegen): handle (v, ok) methods

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5ef5d14f864eb355ffb96d2727a19b61c5d2b362"><tt>5ef5d14f</tt></a> Update cors.md</summary>

I had problems reading this page and applying it to my project. With these changes it worked on my end

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/997da421b0b80884fcb43c8c6a22d747564b301c"><tt>997da421</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1436">#1436</a> from ddouglas/patch-1</summary>

Upgrade graphql-playground to 1.7.26

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/be4514c60a3cef3673b003a545818eb67e620506"><tt>be4514c6</tt></a> Upgrade graphql-playground to 1.7.26

- <a href="https://github.com/99designs/gqlgen/commit/918801eac861c0ceb5cf45969745f674b823ef7c"><tt>918801ea</tt></a> Change 'Changeset' doc example to mutation

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/862762c77bae8b5f119c401d037358cfaf33fa52"><tt>862762c7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1409">#1409</a> from zikaeroh/chi-mod</summary>

Upgrade go-chi to v1.5.1 with module support

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c30ff3ddec18b7a9465e0b81ab0af301d3899141"><tt>c30ff3dd</tt></a> Upgrade go-chi to v1.5.1 with module support

- <a href="https://github.com/99designs/gqlgen/commit/a9c8fabff6d56c9c523ca68764dcd9f9e6cd4f45"><tt>a9c8fabf</tt></a> int64 support

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b484fc27b153639c6d96a3f1df7e952d587749be"><tt>b484fc27</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1401">#1401</a> from oseifrimpong/patch-1</summary>

fix typo

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4cc031afba041dd41b9f533cba3ac1399e8b66cd"><tt>4cc031af</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1394">#1394</a> from j2gg0s/fix-default-recover-func</summary>

bugfix: Default Recover func should return gqlerror.Error

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2af51336b421f45bd5076f5ee144d4d44c15ec54"><tt>2af51336</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1400">#1400</a> from 99designs/sanstale</summary>

Remove stale bot

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/34a442c7980f5ba363f58aaee67b9ddaa77d7520"><tt>34a442c7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1399">#1399</a> from 99designs/prevent-possible-error-deadlock</summary>

Dont hold error lock when calling into error presenters

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1123ba0da6c0cd0f5ab71d7e4c8aae7a5e8f40b4"><tt>1123ba0d</tt></a> Update gin.md</summary>

Changed this:
`In your router file, define the handlers for the GraphQL and Playground endpoints in two different methods and tie then together in the Gin router:
`
to: 
`In your router file, define the handlers for the GraphQL and Playground endpoints in two different methods and tie them together in the Gin router:
`

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/89a9f743240e589c476df6a01e122564314ec215"><tt>89a9f743</tt></a> Remove stale bot</summary>

We tried it, but it's just causing more work both for maintainers and reporters of errors.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4628ef8422c2118ee482dc15a65eedff5144d34c"><tt>4628ef84</tt></a> Dont hold error lock when calling into error presenters</summary>

This can result in a deadlock if error handling code calls GetErrors.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d0d5f7db3d8f09f087abedfe3958423c9f5e4fb9"><tt>d0d5f7db</tt></a> bugfix: Default Recover func should return gqlerror.Error

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/18b5df19bba282e0217dc269f6e9e211b52fc707"><tt>18b5df19</tt></a> codegen/config: Add a new API to finish an already-validated config</summary>

LoadConfig parses the config from yaml, but it does a bunch of other
things too.  We want to parse the config ourselves, so that we can have
extra fields which will be passed to our plugins.  Right now, that means
we either have to duplicate all of LoadConfig, or write the config back
to disk only to ask gqlgen re-parse it.

In this commit, I expose a new function that does all the parts of
LoadConfig other than the actual YAML-reading: that way, a caller who
wants to parse the YAML themselves (or otherwise programmatically
compute the config) can do so without having to write it back to disk.

An alternative would be to move all this logic to Config.Init(), but
that could break existing clients.  Either way would work for us.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0e12bfbfde3b8fc8e54241fcd107c6301c98c6fa"><tt>0e12bfbf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1269">#1269</a> from dqn/new-line-at-the-end-of-file</summary>

Add a new line to end of the file schema.graphqls

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/22c5d1f56eb081104b586c6a73f9324ded90a8b5"><tt>22c5d1f5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1303">#1303</a> from kunalpowar/inline-directives-doc</summary>

Update README.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/88cffee4fc5a29b2899d70d93b4ef64c145e4722"><tt>88cffee4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1356">#1356</a> from maapteh/chore/chat-example-update</summary>

Chore: update Chat example

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1e8c34e540c04b5bb203a385788e0e02111f0afb"><tt>1e8c34e5</tt></a> Dont export  Input

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/de8af66c892a5b2ec8b2ce7ed274003d7706d904"><tt>de8af66c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1360">#1360</a> from Captain-K-101/master</summary>

Update introspection.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0975691550f23be6b68c29b57840e2af4698eac4"><tt>09756915</tt></a> Update introspection docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/651eda40fe4318ec7acd48e7e1f1eb933331c22d"><tt>651eda40</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1374">#1374</a> from rudylee/docs-file-upload-small-typo</summary>

Fix small typo in file upload docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/94252e047b6ab16532c003003511762e6bf4f655"><tt>94252e04</tt></a> singleUpload consistency

- <a href="https://github.com/99designs/gqlgen/commit/c9d346f549ca75ab19097ea8bf90f0d3666cc117"><tt>c9d346f5</tt></a> Fix small typo in file upload docs

- <a href="https://github.com/99designs/gqlgen/commit/9f85161930220becde384662e4da9f4b457ce19f"><tt>9f851619</tt></a> add uint, uint64, uint32 types in graphql

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0625525f663e7fc8d246f5435c1545beef137173"><tt>0625525f</tt></a> Update introspection.md</summary>

updated disabling interospect

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c6a93aa77298d91dac38d7a6bf99c6d14a097b82"><tt>c6a93aa7</tt></a> split layout components to their own part, makes sample more readable

- <a href="https://github.com/99designs/gqlgen/commit/7904ef6fc5f481a4df5938d3cda34919c2160db7"><tt>7904ef6f</tt></a> channel is switchable too

- <a href="https://github.com/99designs/gqlgen/commit/13752055b7c17d5f5d9c0119d54fb65ff623b648"><tt>13752055</tt></a> add some layout for demo :)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/82ca6e24739f63266788b34e55248b68aed0202f"><tt>82ca6e24</tt></a> Create package declaration to run dataloaden</summary>

ref: https://github.com/vektah/dataloaden/issues/35

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bf5491364e91b40372ef9132bc68b9ec303c0a38"><tt>bf549136</tt></a> use Apollo docs styling for the gql var uppercase

- <a href="https://github.com/99designs/gqlgen/commit/36045a3758836ece8954d1ac969e6055d3632bcb"><tt>36045a37</tt></a> do not autofocus

- <a href="https://github.com/99designs/gqlgen/commit/0502228a61387cab237d786a530d41c662216661"><tt>0502228a</tt></a> chore: update example to React hooks and latest Apollo client

- <a href="https://github.com/99designs/gqlgen/commit/e6e64224a32ca35bf543b1cb18e7ccfe65ba824f"><tt>e6e64224</tt></a> update deps

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3a31a752df764738b1f6e99408df3b169d514784"><tt>3a31a752</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1345">#1345</a> from abeltay/fix-alignment</summary>

Fix tab spacing in cors.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0c68337cee2cb15ef0487ae3b9bf902d2d2a96d1"><tt>0c68337c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1346">#1346</a> from abeltay/fix-typo</summary>

Fix typo in migration guide

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/436a88adf607aede4ced31c54342c23fddaf1f95"><tt>436a88ad</tt></a> Fix typo in migration guide

- <a href="https://github.com/99designs/gqlgen/commit/3791f71df5a39d0e57e59ce6c9460627d45a8ab0"><tt>3791f71d</tt></a> Fix tab spacing in cors.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/819e751c2416245370ec00a33ec3b8708aee51c4"><tt>819e751c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1341">#1341</a> from dgraph-io/rajas/fix-gqlgen-1299</summary>

Rajas/fix gqlgen 1299

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/789d02f5c632a018b897e3e5c0a82fe3777e9b54"><tt>789d02f5</tt></a> Requested changes

- <a href="https://github.com/99designs/gqlgen/commit/130ed3f7d9e586aec1960557f9819159986cb53b"><tt>130ed3f7</tt></a> Fix different alias with same name in inline fragment

- <a href="https://github.com/99designs/gqlgen/commit/f4669ba9b54cdda6340c6e0c16d05eff2ee4fa21"><tt>f4669ba9</tt></a> v0.13.0 postrelease bump

- <a href="https://github.com/99designs/gqlgen/commit/07c065946504daa7c9fad0a1d0915713a45c9818"><tt>07c06594</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/1c9f24b2e7f75bb1134318a6d805e3ab3f109ba5"><tt>1c9f24b2</tt></a> remove triming space for schemaDefault

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.13.0"></a>
## [v0.13.0](https://github.com/99designs/gqlgen/compare/v0.12.2...v0.13.0) - 2020-09-21
- <a href="https://github.com/99designs/gqlgen/commit/07c1f93b3d05a07dd7403c1f99793b1976228a48"><tt>07c1f93b</tt></a> release v0.13.0

- <a href="https://github.com/99designs/gqlgen/commit/259f27119bf24ef4806e86334200c216429fbf5c"><tt>259f2711</tt></a> Bump to gqlparser to v2.1.0 Error unwrapping release

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/669a16680de3e31be5650ccdbf6ca4e8f011dcda"><tt>669a1668</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1312">#1312</a> from 99designs/error-wrapping</summary>

Always wrap user errors

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9b948a5f7816eb2764ef85b88f9af582f53a1d78"><tt>9b948a5f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1316">#1316</a> from skaji/is-resolver</summary>

Add IsResolver to FieldContext

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/77aeb47790c737d337c6794e6e57d38e9326a4e9"><tt>77aeb477</tt></a> Point latest docs to v0.12.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e821b97bfbb922589c9eea649f0415ec3454e446"><tt>e821b97b</tt></a> Always wrap user errors (closes <a href="https://github.com/99designs/gqlgen/issues/1305"> #1305</a>)</summary>

Requires use of go 1.13 error unwrapping.

On measure I think I prefer this approach, even though it's a bigger BC break:
- There's less mutex juggling
- It has never felt right to me that we make the user deal with path when overriding the error presenter
- The default error presenter is now incredibly simple

Questions:
- Are we comfortable with supporting 1.13 and up?
- Should we change the signature of `ErrorPresenterFunc` to `func(ctx context.Context, err *gqlerror.Error) *gqlerror.Error`?
    - It always is now, and breaking BC will force users to address the requirement for `errors.As`

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/51b580de1408b934a03614dcdea94f6aa6f25f97"><tt>51b580de</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1324">#1324</a> from bemasher/patch-1</summary>

Fix typos in README.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8b2a023cdb586fdda6227a065cb776318f5fa33e"><tt>8b2a023c</tt></a> Fix typos in README.md

- <a href="https://github.com/99designs/gqlgen/commit/3e5dd956afecb08404a4ff2120b53f5979b44054"><tt>3e5dd956</tt></a> add test for FieldContext.IsResolver

- <a href="https://github.com/99designs/gqlgen/commit/1524989b7d252219502f045e8554c20cd4f34dce"><tt>1524989b</tt></a> go generate

- <a href="https://github.com/99designs/gqlgen/commit/55951163bacbda7399b23824903e9d4a318ebd51"><tt>55951163</tt></a> add IsResolver to FieldContext

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/622316e764b5c296455c558d6fae4b314ca52733"><tt>622316e7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1295">#1295</a> from a-oz/a-oz-patch-1</summary>

Update getting-started.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4c11d9fa30e660180ed7d6be068e69888ad77b6d"><tt>4c11d9fa</tt></a> Update getting-started.md</summary>

fix typo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b4375b04fcb0f7c5ca12d547c88a69c4a8d9f2b5"><tt>b4375b04</tt></a> v0.12.2 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.12.2"></a>
## [v0.12.2](https://github.com/99designs/gqlgen/compare/v0.12.1...v0.12.2) - 2020-08-18
- <a href="https://github.com/99designs/gqlgen/commit/03cebf201ec911411c2c1463ff9b05dfe574bd40"><tt>03cebf20</tt></a> release v0.12.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e3ce560de7c2ceb00076b439495883548c462c78"><tt>e3ce560d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1288">#1288</a> from alexsn/nopath-field-noerror</summary>

avoid computing field path when getting field errors

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/108975c3add6ca4fd1e0d629b82813073d5f49b6"><tt>108975c3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1284">#1284</a> from dgraph-io/jatin/sameFieldSameTypeGettingIgnored</summary>

fix same field name in two different fragments

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/eb424a22657c04b3f2a08d592fe132ca8ff6309f"><tt>eb424a22</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1294">#1294</a> from 99designs/fix-init</summary>

Allow rewriter to work on empty but potentially importable packages

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a87c54adf263fc0ab9f7349363c678b278619bb4"><tt>a87c54ad</tt></a> Allow rewriter to work on empty but potentially importable ckages

- <a href="https://github.com/99designs/gqlgen/commit/8a7f3e64194f56138266694ca131072ab1d68e43"><tt>8a7f3e64</tt></a> clean code

- <a href="https://github.com/99designs/gqlgen/commit/fd0f97ceff73bd9cbca11080c2f36434b0e910a0"><tt>fd0f97ce</tt></a> avoid computing field path when getting field errors

- <a href="https://github.com/99designs/gqlgen/commit/2d59b684a3c41a7573cee1228a7f5bcb4e09392b"><tt>2d59b684</tt></a> ran fmt on test

- <a href="https://github.com/99designs/gqlgen/commit/3a1530755476fc690ffbe769b8f8a582f6f20e0d"><tt>3a153075</tt></a> ran fmt

- <a href="https://github.com/99designs/gqlgen/commit/defd71199ad4dbb0c6a3fa3db6eb1e7ebd0ff97a"><tt>defd7119</tt></a> added test

- <a href="https://github.com/99designs/gqlgen/commit/9fcdbcd1fafad2763d58e34eac0adcbcc5d0d8b4"><tt>9fcdbcd1</tt></a> fix panic test

- <a href="https://github.com/99designs/gqlgen/commit/473d63c02710986ea49749ec6cba44019a82bb9d"><tt>473d63c0</tt></a> change name to alias

- <a href="https://github.com/99designs/gqlgen/commit/849e3eace8b82cb3481db3674068723862842cb8"><tt>849e3eac</tt></a> added check for object defination name

- <a href="https://github.com/99designs/gqlgen/commit/08eee0fc5dbf08af483d7d7a5a9337ef97d8f8dd"><tt>08eee0fc</tt></a> v0.12.1 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.12.1"></a>
## [v0.12.1](https://github.com/99designs/gqlgen/compare/v0.12.0...v0.12.1) - 2020-08-14
- <a href="https://github.com/99designs/gqlgen/commit/0d5f462b25d920a7767ee571e2438dbc835cfbc7"><tt>0d5f462b</tt></a> release v0.12.1

- <a href="https://github.com/99designs/gqlgen/commit/e076b1b03002816516ca6c4a757415cba57d5b13"><tt>e076b1b0</tt></a> Regenerate test server

- <a href="https://github.com/99designs/gqlgen/commit/c952e0de6ac8df5798b04f1d6a53c73f1f691143"><tt>c952e0de</tt></a> v0.12.0 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.12.0"></a>
## [v0.12.0](https://github.com/99designs/gqlgen/compare/v0.11.3...v0.12.0) - 2020-08-14
- <a href="https://github.com/99designs/gqlgen/commit/7030212379f41dea8a1cac2f76f9e56e3054cf24"><tt>70302123</tt></a> Version 0.12.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3b633dfa11874ae5fc8d03da6c963acea6c12a07"><tt>3b633dfa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1267">#1267</a> from ImKcat/master</summary>

Fixed transport not support issue

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c9a27ae3bee99c22ea0af06744257f0184f78e70"><tt>c9a27ae3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1255">#1255</a> from s-ichikawa/fix-object-directive-bug</summary>

Fix bug about OBJECT directive

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e9863af16691f94a68c23f073854930bd754781f"><tt>e9863af1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1276">#1276</a> from Ghvstcode/master</summary>

Documentation Fixes

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/04f6a691577d627caaa0952c13cedecda5455e28"><tt>04f6a691</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1277">#1277</a> from 99designs/direct-pointer-binding</summary>

Support pointers in un/marshal functions

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bef9c8bf3a2d7d63531e4bf9540d550977289a5c"><tt>bef9c8bf</tt></a> Add comments and docs for pointer scalars

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/997efd0395b8d51312cc73f709a7545c928709a0"><tt>997efd03</tt></a> Reintroduce special cast case for string enums</summary>

This reverts commit 89960664d05f0e93ed629a22753b9e30ced2698f.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8561c056cb74834233f33a0ae4d42eafcc3e7c29"><tt>8561c056</tt></a> Replace awkward loop in buildTypes with recursion

- <a href="https://github.com/99designs/gqlgen/commit/d65b04f9ff3fa3435ef103da7ad5b45a59d82868"><tt>d65b04f9</tt></a> Clean up generated code

- <a href="https://github.com/99designs/gqlgen/commit/e1c463a4c873306d634279e21c7bd631d9c54fca"><tt>e1c463a4</tt></a> Linting

- <a href="https://github.com/99designs/gqlgen/commit/89960664d05f0e93ed629a22753b9e30ced2698f"><tt>89960664</tt></a> Remove unused special cast case for string enums

- <a href="https://github.com/99designs/gqlgen/commit/196954bc64771795d179a91bbd1422569f9fdced"><tt>196954bc</tt></a> Bind directly to pointer types when possible, instead of always binding to value types

- <a href="https://github.com/99designs/gqlgen/commit/5b3d08db47c9c3b49687411f6ecbcbc84b66a495"><tt>5b3d08db</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/efd33dab01483cfa54d634690aea67ab914dd72a"><tt>efd33dab</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/f35b162f214ca0ae1461c25fde29d41b55293f16"><tt>f35b162f</tt></a> Fixed transport not support issue

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/39a12e0f1b6d9833f516f0271db0dbfa45c5ec45"><tt>39a12e0f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1134">#1134</a> from seriousben/fix-default-config-no-ast-sources</summary>

Add LoadDefaultConfig to load the schema by default

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1b23cf15b134cd695a11fe899d59c5457778a8be"><tt>1b23cf15</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1264">#1264</a> from 99designs/go-1.14</summary>

Target multiple go versions for CI

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/dbbda22ef42a921190cf52b3f23fa53b54726828"><tt>dbbda22e</tt></a> go 1.14

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ce964c1f46bedff709f8d53356d83a3e983295f4"><tt>ce964c1f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1115">#1115</a> from bowd/add-input-path-for-unmarshaling</summary>

Add input path in unmarshaling errors

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bde4291cfa7669a889db6e5e518218a855ffd433"><tt>bde4291c</tt></a> shadow context to ensure scoped context use

- <a href="https://github.com/99designs/gqlgen/commit/c43990a00ccfb11d950e7b33716d214f2fefec5c"><tt>c43990a0</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

- <a href="https://github.com/99designs/gqlgen/commit/6be2e9df78c81b3fa45ac00717ef7fa505ed6a4f"><tt>6be2e9df</tt></a> fix fileupload example

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ad675f0092bc8973dba5114ffa168379647bde45"><tt>ad675f00</tt></a> Allow custom resolver filenames using `filename_template` option (closes <a href="https://github.com/99designs/gqlgen/issues/1085"> #1085</a>)</summary>

resolve merge conflicts.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fbfdd41c12147ec7fbc307163e4667dd28065626"><tt>fbfdd41c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1262">#1262</a> from sateeshpnv/gqlparser-alias (closes <a href="https://github.com/99designs/gqlgen/issues/1258"> #1258</a>)

- <a href="https://github.com/99designs/gqlgen/commit/99fafc9f19c838f26afc71785da6adbf0a1dbe76"><tt>99fafc9f</tt></a> [issue <a href="https://github.com/99designs/gqlgen/pull/1258">#1258</a>] explicitly add gqlparser alias to vektah/gqlparser/v2 import

- <a href="https://github.com/99designs/gqlgen/commit/49291f234e99878b925946efeb13c5bf1b2c348e"><tt>49291f23</tt></a> fix bug in OBJECT directive

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0fbf293f29f3c2ca01822723f1f43e01cab358d4"><tt>0fbf293f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1248">#1248</a> from sotoslammer/master</summary>

close the connection when run returns

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d7eabafb4085e7b802ab553c13ac62fa6e3331f8"><tt>d7eabafb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1246">#1246</a> from arkhvoid/master</summary>

Fix typo cause memory problem on upload

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/21b223b8d37208bdfb45d49dc7c35a557255e548"><tt>21b223b8</tt></a> Fix typo cause memory problem on upload

- <a href="https://github.com/99designs/gqlgen/commit/cc9c520f1ecf11e5786f2aca8d9cf24ef4af2f2e"><tt>cc9c520f</tt></a> close the connection when run returns

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8494028eac6a22ef26f71ff30a9eb5738a86adff"><tt>8494028e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1243">#1243</a> from 99designs/nilable-nullable-unnmarshal</summary>

Remove a bunch of unneeded nil checks from non-nullable graphql type unmarshalling

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b81138dac0e4ba71f04b7e686e18eb3744e23919"><tt>b81138da</tt></a> Add test for nillable input slice

- <a href="https://github.com/99designs/gqlgen/commit/14d1a4dc0a9242154d3a22787d92ab70239f079a"><tt>14d1a4dc</tt></a> Only return nil for nilable types when the graphql spec would allow it

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3e59a10d4268671b92428ef1b860aebbb73da60b"><tt>3e59a10d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1215">#1215</a> from ddouglas/master</summary>

Adding Missing Header to response

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1650c499548c074d5afb152a90235589dc98d107"><tt>1650c499</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1242">#1242</a> from 99designs/named_map_references</summary>

Do not use pointers on named map types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d11f60218ccb4e5b17d702e736585f02978b69a4"><tt>d11f6021</tt></a> Do not use pointers on named map types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/acaee3615ba86465e710635d7676cf05e017eb9b"><tt>acaee361</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1121">#1121</a> from Khan/extern-only</summary>

Do not require a resolver for "empty" extended types.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/555db6d20b0157674b5ad05f6ce8856c6502db2e"><tt>555db6d2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1224">#1224</a> from frederikhors/patch-1</summary>

Indentation misprint

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/77b37bb290c55008d2ca1653aca587d3f1ea17e5"><tt>77b37bb2</tt></a> Indentation misprint

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a3c38c6574b71389895e5162296ed12e13347349"><tt>a3c38c65</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1221">#1221</a> from longngn/patch-1</summary>

Update dataloaders.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/71182de820e20edc5fd5a2e7363f089ff75bdb9a"><tt>71182de8</tt></a> Update dataloaders.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d81baeed9f212c2728e6d9901bdec99929787ac2"><tt>d81baeed</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1218">#1218</a> from StevenACoffman/patch-1</summary>

Update feature comparison for federation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2c1f234503d94c6397833b018078ee9099ecf643"><tt>2c1f2345</tt></a> Update feature comparison for federation (closes <a href="https://github.com/99designs/gqlgen/issues/5"> #5</a>)

- <a href="https://github.com/99designs/gqlgen/commit/e19d43bcb46df36bf6d8d8699788bfd929a596f5"><tt>e19d43bc</tt></a> Adding test

- <a href="https://github.com/99designs/gqlgen/commit/4a62f0121af3a730fa57bf9f9124beaf6a0809d7"><tt>4a62f012</tt></a> Adding ContentType header to GET request responses

- <a href="https://github.com/99designs/gqlgen/commit/f5de4731aa552bff75d6ddb06f7d7338388c5a34"><tt>f5de4731</tt></a> Add timeout to integration test

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a21a6633b841779b9b720f95e7297db888935993"><tt>a21a6633</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1189">#1189</a> from RichardLindhout/patch-1</summary>

Upgrade to OperationContext and remove duplicate fields to fix https:…

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/543317a28754c46f1679af33dde0c311d73f7ddd"><tt>543317a2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1170">#1170</a> from alexsn/apollotracing/nopanic</summary>

apollotracing: skip field interceptor when on no tracing extension

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d347d97278ec166866b309d458001a17ed5779e0"><tt>d347d972</tt></a> Update stale.yml

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/032854bb8a0877796d8d84d7d619d503beae5d52"><tt>032854bb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1154">#1154</a> from gsgalloway/master</summary>

Add operation context when dispatching

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ccc4eb1db613027376e0bd69da02bdde8914e911"><tt>ccc4eb1d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1188">#1188</a> from k-yomo/update-errors-doc</summary>

Update outdated examples in errors doc

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/628b83c19657e042a0c4e8d694bcabb6ac182b1f"><tt>628b83c1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1198">#1198</a> from ddevault/pgp</summary>

codegen: add PGP to common initialisms

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d881559bb1d9e6a4d61deb296cb6a18a6d8e1476"><tt>d881559b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1202">#1202</a> from whereswaldon/patch-1</summary>

doc: fix typo in embedded struct example

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b6ce42a7a218c33e327e68d7ebe5e49017dbe223"><tt>b6ce42a7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1207">#1207</a> from k-yomo/update-gorilla-websocket</summary>

Update gorilla/websocket to v1.4.2 to resolve vulnerability

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c5bfe9d3d9341a23d1254a38200a6d957c52e21c"><tt>c5bfe9d3</tt></a> Update gorilla/websocket to v1.4.2 to resolve vulnerability

- <a href="https://github.com/99designs/gqlgen/commit/55c16e93ed6d5e63eb34d6d99f3a2db9830ad822"><tt>55c16e93</tt></a> doc: fix typo in embedded struct example

- <a href="https://github.com/99designs/gqlgen/commit/89eb19937a8903b41fd85398a89443e30f63db01"><tt>89eb1993</tt></a> codegen: add PGP to common initialisms

- <a href="https://github.com/99designs/gqlgen/commit/9ab7294d79825e96b00063abecccfabf4286ba9b"><tt>9ab7294d</tt></a> apollotracing: skip field interceptor when on no tracing extension

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/40570d1b4d70070c84915f7a468e406705b3f3ef"><tt>40570d1b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1163">#1163</a> from fwojciec/master</summary>

fix redundant type warning

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3f7f60bf180a405ea24600c82a6e5b24c605ca9f"><tt>3f7f60bf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1181">#1181</a> from tmc/patch-1</summary>

Update getting-started.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6518d8391acf8618b1c39539b9bb7306814898d3"><tt>6518d839</tt></a> Upgrade to OperationContext and remove duplicate fields to fix https://github.com/99designs/gqlgen/pull/1161

- <a href="https://github.com/99designs/gqlgen/commit/632904adf183e25559de42a06d2fde2a0af0ca53"><tt>632904ad</tt></a> Update outdated examples in errors doc

- <a href="https://github.com/99designs/gqlgen/commit/0921915d02741d3021a69ac7834bd76e9bbc38ab"><tt>0921915d</tt></a> Update getting-started.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0a40481343ef344121b5a2a5c0910ff7391aad1f"><tt>0a404813</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1117">#1117</a> from s-ichikawa/object-directive</summary>

Add support for OBJECT directive

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/90ee8dedb8dbea7d447801380639de7d475e62e8"><tt>90ee8ded</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1137">#1137</a> from ddevault/master</summary>

Replace ~ with א in package names

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e4c699dcd4b8b1ccc62c31c158670904092fb374"><tt>e4c699dc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1147">#1147</a> from ddevault/docs</summary>

Add links to godoc to the README and docsite

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/73746621696f25c55bd64d2d81b09fe314667c21"><tt>73746621</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1131">#1131</a> from muraoka/fix-typo</summary>

Fix typo in authentication docs

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ace558b411ec2633cfa4cc8b76d7a035d3216cef"><tt>ace558b4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1124">#1124</a> from OpenSourceProjects/update-apq-documentation</summary>

Update APQ example to reflect newer API

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3c126f9edba3ba3f4c252b003f052743d5fd1c72"><tt>3c126f9e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1119">#1119</a> from skaji/patch-1</summary>

type Person -> type Person struct

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1610039e7302c8bc543421af54e4a5005b8f9d13"><tt>1610039e</tt></a> updated generated code

- <a href="https://github.com/99designs/gqlgen/commit/905e1aadfbd449a045f65b1368a0936daf591143"><tt>905e1aad</tt></a> fix redundant type warning

- <a href="https://github.com/99designs/gqlgen/commit/39ded924092564cd953040c9a04be9ce8800aaa1"><tt>39ded924</tt></a> fix ctx

- <a href="https://github.com/99designs/gqlgen/commit/e7798ff26198eceb79f10d5743558aba02629ec8"><tt>e7798ff2</tt></a> insert operation context

- <a href="https://github.com/99designs/gqlgen/commit/6f78c6ac5c071700b80507c4ad72283a7911acaf"><tt>6f78c6ac</tt></a> Add links to godoc to the README and docsite

- <a href="https://github.com/99designs/gqlgen/commit/9b823a348713911f32e73826db69d61df25163c0"><tt>9b823a34</tt></a> Replace ~ with א in package names (closes <a href="https://github.com/99designs/gqlgen/issues/1136"> #1136</a>)

- <a href="https://github.com/99designs/gqlgen/commit/35a904829a89e9d29d1404ffed882cc5eda64ccf"><tt>35a90482</tt></a> Add LoadDefaultConfig to load the schema by default

- <a href="https://github.com/99designs/gqlgen/commit/07a5494b34560fdfccfdd43f1430cbf064b860a6"><tt>07a5494b</tt></a> Fix typo in docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/04b120c9044a20ac8c14f8329c529862ccdc6b6f"><tt>04b120c9</tt></a> Update APQ example to reflect newer API</summary>

The example in APQ relates to the old handlers. This brings it up to
show how extensions can be used - and uses the new API for registering
plugins that come in the graph.

The cache example now implements the graphql.Cache interface

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/55e0f0db84dfc709c837074d57cd26f9a579e9c5"><tt>55e0f0db</tt></a> Check in a place where `Entity` might be nil now.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1ecd0749dd3ba3ba6ba387d0391e69d410b313de"><tt>1ecd0749</tt></a> Handle the case that all entities are "empty extend".</summary>

In that case, there are no resolvers to write, so we shouldn't emit
any.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0e2666fb2bfca1134661dc835b80b1061550225c"><tt>0e2666fb</tt></a> Run `go fmt`

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/36b5ed834d3444affb6df88b65f4f3ec8a42d488"><tt>36b5ed83</tt></a> Actually, we need to check all-external, not all-key.</summary>

We might well be defining our own type that has only key-fields, but
if they're not external then we're the primary provider of the type

Test plan:
go test ./plugin/federation/

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7e3f5844adb79c776a30a6e01fcb2b373049fb8e"><tt>7e3f5844</tt></a> Do not require a resolver for "empty" extended types.</summary>

Summary:
If our schema has a field with a type defined in another service, then
we need to define an "empty extend" of that type in this service, so
this service knows what the type is like.  But the graphql-server will
never ask us to actually resolve this "empty extend", so we don't
require a resolver function for it.  Example:
```
   type MyType {
      myvar: TypeDefinedInOtherService
   }

   // Federation needs this type, but it doesn't need a resolver for
   // it!  graphql-server will never ask *us* to resolve a
   // TypeDefinedInOtherService; it will ask the other service.
   extend TypeDefinedInOtherService @key(fields: "id") {
      id: ID @extends
   }
```

Test Plan:
I manually tested this on a service (`assignments`) that we have that
fell afoul of this problem.  But I had a hard time adding tests inside
gqlgen because the error happens at validation-time, and the
federation tests are not set up to go that far down the processing
path.

Reviewers: benkraft, lizfaubell, dhruv

Subscribers: #graphql

Differential Revision: https://phabricator.khanacademy.org/D61883

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9c80bb5ba7e315c8735599db945105e5acfe86d8"><tt>9c80bb5b</tt></a> type Person -> type Person struct

- <a href="https://github.com/99designs/gqlgen/commit/ea210929aeb7daf4ea1346a655c056df4935c63b"><tt>ea210929</tt></a> add test for object directive

- <a href="https://github.com/99designs/gqlgen/commit/5c3812cbe30cc452490e240f41919508d58c3009"><tt>5c3812cb</tt></a> merge object directives to field directives

- <a href="https://github.com/99designs/gqlgen/commit/8ea5ba2b3f0851c7de1312a24c14e234d051b0bb"><tt>8ea5ba2b</tt></a> Fix additional missed tests

- <a href="https://github.com/99designs/gqlgen/commit/65be2a6e80792993384534da9794c7a47da78244"><tt>65be2a6e</tt></a> Run generate

- <a href="https://github.com/99designs/gqlgen/commit/fd615cf6d4f829e38f0d52cf2c3de886d65d463b"><tt>fd615cf6</tt></a> Fix linting

- <a href="https://github.com/99designs/gqlgen/commit/61fa9903fa26c97a4fabbf3200e22c3d25ffcbc6"><tt>61fa9903</tt></a> Add documentation for scalad error handling

- <a href="https://github.com/99designs/gqlgen/commit/1aa20f25f5c8b3897e8ba47e4adb181a005ee60c"><tt>1aa20f25</tt></a> Add test to highlight usecase

- <a href="https://github.com/99designs/gqlgen/commit/d98ff1b04ca102e74037ac914dddb595aa9c6808"><tt>d98ff1b0</tt></a> Modify templates to include deeper context nesting

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a1a02615f705292de46bd1c14f1710eedc95cb86"><tt>a1a02615</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1104">#1104</a> from oshalygin/docs/update-query-complexity-initialization</summary>

Update Query Complexity Documentation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c68df3c61f389e065a4607c67e973f46f492cd9f"><tt>c68df3c6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1112">#1112</a> from s-ichikawa/delete-unused-code</summary>

delete unused code

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/dfb6558a6fb4af8a60693b046485aeaaffffc507"><tt>dfb6558a</tt></a> run CI on PRs</summary>

PRs from outside the org arent running CI, hopefully this fixes it.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5149231ca6cd10b1928510fb5f1eef78d29e88ce"><tt>5149231c</tt></a> delete unused code

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6f81ff9273850813bec2a4ff705caab66d64eb88"><tt>6f81ff92</tt></a> Update Query Complexity Documentation</summary>

- This pass at the documentation updates the
  appropriate section regarding query complexity,
  specifically in the way that the http.Handler
  is created.
- The deprecated handler.GraphQL calls were replaced
  with NewDefaultServer.
- Instead of passing along the fixed query complexity
  as a second argument to the now deprecated handler.GraphQL
  func, extension.FixedComplexityLimit is used instead.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f0cd7a703261c5ce6274829686ed9611a2b6deb7"><tt>f0cd7a70</tt></a> update doc site to point to latest version

- <a href="https://github.com/99designs/gqlgen/commit/224ff3454cfd3ff505c6ceca5e78978b94073aa1"><tt>224ff345</tt></a> v0.11.3 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.11.3"></a>
## [v0.11.3](https://github.com/99designs/gqlgen/compare/v0.11.2...v0.11.3) - 2020-03-13
- <a href="https://github.com/99designs/gqlgen/commit/4d73535648395cdb8bfbf5f66345a8eaf1b4e0c7"><tt>4d735356</tt></a> release v0.11.3

- <a href="https://github.com/99designs/gqlgen/commit/4b949f2e69026b51ddd26b71d5efb7b5dc8c6aca"><tt>4b949f2e</tt></a> remove copyright notice at bottom of doc pages

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c5039196612e7afeee14fc6373d8dcfe70eb5ab9"><tt>c5039196</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1094">#1094</a> from 99designs/update-upload-docs</summary>

Update file upload docs with Apollo client usage

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5e3cef245a1e87b64c8e57630fa78dae706d069b"><tt>5e3cef24</tt></a> revert <a href="https://github.com/99designs/gqlgen/pull/1079">#1079</a>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/793b0672cee7c096e4b32f4b132469c455407c15"><tt>793b0672</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1100">#1100</a> from sonatard/fast</summary>

Gnerate to fast by exec codegen.GenerateCode before plugin GenerateCode

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6ac2d1cdcb2b2ff33048a768a0e3cfefe8f29d75"><tt>6ac2d1cd</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1097">#1097</a> from 86/86/update-federation-doc</summary>

Add Enable federation section in federation doc

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/97896eeb76a2d942217cc049af7357ea40a1b9e7"><tt>97896eeb</tt></a> exec codegen.GenerateCode before plugin GenerateCode to fast

- <a href="https://github.com/99designs/gqlgen/commit/44f8ba9ff598df631a1149a296e6584766303665"><tt>44f8ba9f</tt></a> Update licence

- <a href="https://github.com/99designs/gqlgen/commit/94701fb78f7c3f2fff84ee22c5df7389f6f18869"><tt>94701fb7</tt></a> add Enable federation section in federation doc

- <a href="https://github.com/99designs/gqlgen/commit/64190309002d888375931eea2d8597cfe5a8af5c"><tt>64190309</tt></a> Update upload docs with Apollo usage

- <a href="https://github.com/99designs/gqlgen/commit/a538119155cf1ce66094a5a1ba1a201905d3e832"><tt>a5381191</tt></a> v0.11.2 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.11.2"></a>
## [v0.11.2](https://github.com/99designs/gqlgen/compare/v0.11.1...v0.11.2) - 2020-03-05
- <a href="https://github.com/99designs/gqlgen/commit/2ccc0aa65998154a57ddb2fcb37046cbffbd6518"><tt>2ccc0aa6</tt></a> release v0.11.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/78f3da2296a5d69967a7ab09b66338ed9bd94033"><tt>78f3da22</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1050">#1050</a> from technoweenie/executor</summary>

Executor

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b82ee517f528dbdbd1236aca61b7aacdfc633978"><tt>b82ee517</tt></a> Fix CI badge

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/42eff5a9a5011606b5a89556dd211e6e08224b19"><tt>42eff5a9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1057">#1057</a> from RichardLindhout/master</summary>

Upgrade to github.com/urfave/cli/v2

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bb5cb8a3dd51980a9b5c27dc4d5da87e9ab9a1c9"><tt>bb5cb8a3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1086">#1086</a> from 99designs/github-actions</summary>

Use GitHub Actions

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cd2b53f210373729a2d483183db08fd6cce38ebd"><tt>cd2b53f2</tt></a> remove os.Exits

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/587bc81c1a52f259da85421500d971c035d8a0cc"><tt>587bc81c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1074">#1074</a> from yudppp/feature/add_contenttype_for_upload</summary>

Add ContentType to graphql.Upload

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a84d657791b8f11a31b5af768d429b662faca312"><tt>a84d6577</tt></a> graphql/handler: revive the existing around func types

- <a href="https://github.com/99designs/gqlgen/commit/f9bb017b440f74dbb3db744638410dc357bdea18"><tt>f9bb017b</tt></a> graphql/executor_test: ensure operation trace is started before every query

- <a href="https://github.com/99designs/gqlgen/commit/57dd8d9c75b1286b318b9236346e4ab14aa99c86"><tt>57dd8d9c</tt></a> graphql/gqlgen: remove unnecessary convenience method

- <a href="https://github.com/99designs/gqlgen/commit/fb86f7b9a8d768cd69c34e8e6654c06f1af9a7b8"><tt>fb86f7b9</tt></a> graphql/executor: remove the naked return

- <a href="https://github.com/99designs/gqlgen/commit/9ae6bc0b357e1ab5325e60cc4a86cb2072a75e01"><tt>9ae6bc0b</tt></a> graphql/executor: reinit all extension values on every Use() call

- <a href="https://github.com/99designs/gqlgen/commit/f3909a8aa8b5fc00375cfa7905e89cd873ff514f"><tt>f3909a8a</tt></a> graphql/executor: make ext funcs private

- <a href="https://github.com/99designs/gqlgen/commit/df9e7ce3617d082f584edb998d9428fb6ce633ff"><tt>df9e7ce3</tt></a> Run CI on push only

- <a href="https://github.com/99designs/gqlgen/commit/ed76bc923b36dda8437a70dc0f43cbaf21295864"><tt>ed76bc92</tt></a> Update badge

- <a href="https://github.com/99designs/gqlgen/commit/5a1a54463628a2ce3fe9dc67e70e120a42515d56"><tt>5a1a5446</tt></a> Coveralls fixes

- <a href="https://github.com/99designs/gqlgen/commit/41acc753cf76ced6a545d9229579db57c047710e"><tt>41acc753</tt></a> Fix windows line endings

- <a href="https://github.com/99designs/gqlgen/commit/390cea4fea9261b1fe12b689d7185381110c4051"><tt>390cea4f</tt></a> Replace Appveyor with Github Actions

- <a href="https://github.com/99designs/gqlgen/commit/85be072f4c28fd367ab3ce2b5016928acefa67bd"><tt>85be072f</tt></a> Replace CircleCI with Github Actions

- <a href="https://github.com/99designs/gqlgen/commit/8d540db3f8787395c876c3e612966d43bd98a6a1"><tt>8d540db3</tt></a> fix: Add Upload.ContentType test

- <a href="https://github.com/99designs/gqlgen/commit/f21832af91c7b9f7523f4f149c775700860286c1"><tt>f21832af</tt></a> fix: Fixed Upload type document

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b165568cce7d5027ba1ee2da68dad3012b12d189"><tt>b165568c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1071">#1071</a> from kandros/fix-server-path</summary>

fix server path

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9d7648aa95d7479b11685350d6bd71d26e9aecac"><tt>9d7648aa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1072">#1072</a> from wtask/patch-1</summary>

Fix a typo in sql example

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/24400c9b44d7a2dfc177969d13f9ca7cc1f158e0"><tt>24400c9b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1079">#1079</a> from sonatard/remove-unused</summary>

Remove unused code

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a7c79891abdd4323ce74b71521508d064f057f6d"><tt>a7c79891</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1081">#1081</a> from sonatard/fix-plugin-test</summary>

Fix unlink file path in resolvergen test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e7bf75481aa110b533e1cc710a8dbeb88930a6ac"><tt>e7bf7548</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1080">#1080</a> from sonatard/fix-testdata</summary>

Fix test data

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3a61dc00a20288601fdbc6f5192333f2244dc85a"><tt>3a61dc00</tt></a> Fix unlink file path in resolvergen test

- <a href="https://github.com/99designs/gqlgen/commit/df5ac929eb4bbd3090821b5350ba33493a78d58c"><tt>df5ac929</tt></a> Fix test data

- <a href="https://github.com/99designs/gqlgen/commit/b2843f67e04c2eff0579f5f05f3cfe21e9a05ba3"><tt>b2843f67</tt></a> Remove unused code

- <a href="https://github.com/99designs/gqlgen/commit/cff73f71fe116cd7549cd431aa13ea21c377e9b5"><tt>cff73f71</tt></a> Add ContentType to Upload

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f0ebc0dfbe4d260129161ce6b6ce1d7ca44a325c"><tt>f0ebc0df</tt></a> Fix a typo in sql example</summary>

I think todo is referenced to user by user_id field, not by todo.id

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/22a43d776126936526ae37070ccc765983823f23"><tt>22a43d77</tt></a> fix server path

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b788cce5682351522a12ae63e8f1d11be60276a1"><tt>b788cce5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1054">#1054</a> from 99designs/golint-free-resolvers</summary>

suppress golint messages

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c515d403319c4c85799959501f56d927d221e1ce"><tt>c515d403</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1053">#1053</a> from RichardLindhout/patch-3</summary>

Add practical example of getting all the requested fields

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e57cd44598f2446c5f00888c1c545dde4b63edf3"><tt>e57cd445</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1061">#1061</a> from halvdan/patch-1</summary>

Fix mismatching documentation of Todo struct

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1388fa9465d5ba99b8b704472c281193a2926f3f"><tt>1388fa94</tt></a> Fix mismatching documentation of Todo struct</summary>

Mismatch between the code and the getting started documentation.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/294884adaa279d56cf6008098b2c6d1d88d8f961"><tt>294884ad</tt></a> Rollback go.sum and go.mod as per feedback of [@vektah](https://github.com/vektah)

- <a href="https://github.com/99designs/gqlgen/commit/d8acf1655d750fcd069989bafd86179077face7f"><tt>d8acf165</tt></a> Upgrade to github.com/urfave/cli/v2

- <a href="https://github.com/99designs/gqlgen/commit/81bcbe75812169a1d6525df8cfbe0b42ad828151"><tt>81bcbe75</tt></a> suppress golint messages

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/248130792a7a4fb665b29bbab94f2bc43bebe0c4"><tt>24813079</tt></a> Add practical example of getting all the requested fields</summary>

Based on this https://github.com/99designs/gqlgen/issues/954 was tagged as 'need documentation'

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a53ce377c9f7601f09e00da71a958c6f45deca4c"><tt>a53ce377</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1051">#1051</a> from 99designs/has-operation-context</summary>

Add function to check presense of operation context

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/95e453bfe6db62ee00faa2bd66cd64b537721f04"><tt>95e453bf</tt></a> Add function to check presense of operation context

- <a href="https://github.com/99designs/gqlgen/commit/36365c4103361d395d9025dc811d816ad1edb626"><tt>36365c41</tt></a> graphql/executor: move setExtensions()

- <a href="https://github.com/99designs/gqlgen/commit/3acc942111a9051e53d93138363405bce4f2b33b"><tt>3acc9421</tt></a> graphql/executor: ensure Executor implements graphql.GraphExecutor.

- <a href="https://github.com/99designs/gqlgen/commit/f89b973bbec0c14cd69992e8791f8d23ee31c928"><tt>f89b973b</tt></a> graphql/executor: merge ExtensionList into Executor

- <a href="https://github.com/99designs/gqlgen/commit/c16a77c3195fe0a1a297429208c8816ed6962490"><tt>c16a77c3</tt></a> graphql/handler: replace internal executor type

- <a href="https://github.com/99designs/gqlgen/commit/8fa26cec4065fd90bbc7a01cb3cfa933e2b7b461"><tt>8fa26cec</tt></a> graphql/executor: extract an Executor type from graphql/handler

- <a href="https://github.com/99designs/gqlgen/commit/d5d780c5b59a1fc6fd3fe55274fd09c64505917b"><tt>d5d780c5</tt></a> Point latest docs to 0.11.1

- <a href="https://github.com/99designs/gqlgen/commit/abaa0a041172268ffd8605345e6e2ee44847cf4c"><tt>abaa0a04</tt></a> v0.11.1 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.11.1"></a>
## [v0.11.1](https://github.com/99designs/gqlgen/compare/v0.11.0...v0.11.1) - 2020-02-19
- <a href="https://github.com/99designs/gqlgen/commit/11af15a14ba1f3217f1e81a0aeaf053f3f17d56d"><tt>11af15a1</tt></a> release v0.11.1

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bc07188cd6eaa4b5509843c7cc2d20a03623c759"><tt>bc07188c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1038">#1038</a> from 99designs/feat-check-len</summary>

check slice length

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2c3853c8dc0f500339e7349f9256c3a9ac1ef129"><tt>2c3853c8</tt></a> fix whitespace in comparison

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/07a13861cf2a3bcaa96433a3a5e6380f688489f4"><tt>07a13861</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1043">#1043</a> from 99designs/ensure-panic-handlers-get-applied</summary>

Ensure panic handlers get applied

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/156d306d69255cea25ee3cbf4f0b9f57c0b9f09f"><tt>156d306d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1046">#1046</a> from appleboy/patch</summary>

docs(gin): missing import playground

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/26ee1aa1c9a08d15df5cb189c05195a635be6a64"><tt>26ee1aa1</tt></a> docs(gin): missing import playground

- <a href="https://github.com/99designs/gqlgen/commit/3abe5b32e965a6a05164a66110ac7cafda2e89d3"><tt>3abe5b32</tt></a> add test

- <a href="https://github.com/99designs/gqlgen/commit/6ecdb88de5ba7751f0403d9d09c9b573d1b6c635"><tt>6ecdb88d</tt></a> Merge branch 'master' into feat-check-len

- <a href="https://github.com/99designs/gqlgen/commit/2340f7a7ae54c33e5d065d99a1f30906f64ba229"><tt>2340f7a7</tt></a> Ensure panic handlers get applied

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/25d167613df1e8e50dd895394343238b0e9f92ad"><tt>25d16761</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1039">#1039</a> from VitaliiLakusta/patch-1</summary>

Fix link to examples directory in Federation docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4c47ad16341cd1dfa9ec5e35aa0a99a7d102fb66"><tt>4c47ad16</tt></a> Fix link to examples directory in Federation docs

- <a href="https://github.com/99designs/gqlgen/commit/2506dce04f7be87b44c2ec8d99b9fcadc7c7db75"><tt>2506dce0</tt></a> check slice len

- <a href="https://github.com/99designs/gqlgen/commit/1a68df34c397d7fddb161a5ad3ab0b293f683eb9"><tt>1a68df34</tt></a> fix origin/master reference in switcher

- <a href="https://github.com/99designs/gqlgen/commit/199cfedf8be47ed8141530ca9f42a04e22634d4d"><tt>199cfedf</tt></a> remove old docs that no longer run with new layout

- <a href="https://github.com/99designs/gqlgen/commit/556c84843e1f6452bb0983a199772e6f17ff6855"><tt>556c8484</tt></a> fix paths

- <a href="https://github.com/99designs/gqlgen/commit/282100c8205af2f964021c67550e0a118b4e751c"><tt>282100c8</tt></a> use current layout to build old doc content

- <a href="https://github.com/99designs/gqlgen/commit/4c38b8b4ae49aed623d5a0fd5447462109c9dba9"><tt>4c38b8b4</tt></a> v0.11.0 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.11.0"></a>
## [v0.11.0](https://github.com/99designs/gqlgen/compare/v0.10.2...v0.11.0) - 2020-02-17
- <a href="https://github.com/99designs/gqlgen/commit/368597aa18d82bc778e45d8e1f7a817a70ca62a7"><tt>368597aa</tt></a> release v0.11.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e65d62285c5af647fa58ae947fde286f3e0ccc9c"><tt>e65d6228</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1036">#1036</a> from 99designs/update-v011-docs</summary>

Update 0.11 migration docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/11f979365864fd6739fcf9c6d8d6deddfed834ba"><tt>11f97936</tt></a> Update 0.11 migration docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2b3eed303e63e433088efbfa026b7a42e66ff0de"><tt>2b3eed30</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1034">#1034</a> from 99designs/strip-underscores-from-entity-interfaces</summary>

Trim underscores from around go identifiers

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b2d9bfcbfac96d920b3d1dd2ed7eb20f098a1d61"><tt>b2d9bfcb</tt></a> Update stale.yml

- <a href="https://github.com/99designs/gqlgen/commit/1ac8b5aeca013d19ce23c33db0e7c360fdc06095"><tt>1ac8b5ae</tt></a> Update stale.yml

- <a href="https://github.com/99designs/gqlgen/commit/4b9dfa61085d478814efbdd8ac237e3f81d4189b"><tt>4b9dfa61</tt></a> trim underscores from around go identifiers

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7cac3610246fb0b8a9028c88f57e63c41e5cced7"><tt>7cac3610</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1027">#1027</a> from sonatard/response-errors</summary>

propagate resolver errors to response error in ResponseMiddleware

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/14dccc57885df5d5ca4ef347c1b80f5f3648719a"><tt>14dccc57</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1022">#1022</a> from 99designs/feat-gqlparser-117</summary>

example about apply https://github.com/vektah/gqlparser/pull/117

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cf6f76830d3ebbfbce26a96f709edbb11c828551"><tt>cf6f7683</tt></a> bump to gqlparser v2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4ece3857c9abb4dc0c6122af6a72d4d5ce134feb"><tt>4ece3857</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1028">#1028</a> from abhimanyusinghgaur/master</summary>

Respect includeDeprecated for EnumValues

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9638ce0f33d314c9ca5a2154f0b901d26ee7c719"><tt>9638ce0f</tt></a> Fix format

- <a href="https://github.com/99designs/gqlgen/commit/51b921fab81fc6c9f8304f4a0013654693bcc77d"><tt>51b921fa</tt></a> Fix format

- <a href="https://github.com/99designs/gqlgen/commit/07ffcc821c7fbeb9fea3c6b0916367a48f7d5b82"><tt>07ffcc82</tt></a> Respect includeDeprecated for EnuValues

- <a href="https://github.com/99designs/gqlgen/commit/d58434c9f4827a1a5448009b94341498dae87105"><tt>d58434c9</tt></a> propagate resolver errors to response error in ResponseMiddleware

- <a href="https://github.com/99designs/gqlgen/commit/598559252c076967ae3120ee203885954c254a2a"><tt>59855925</tt></a> go mod tidy

- <a href="https://github.com/99designs/gqlgen/commit/e4530da6a4636a6d6430ae93356fa423d7a5ead7"><tt>e4530da6</tt></a> apply https://github.com/vektah/gqlparser/pull/117

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/30e23757e9b50612f3300edd84310efe7eac9d4d"><tt>30e23757</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1020">#1020</a> from 99designs/handle-interfaces-implementing-interfaces</summary>

Handle interfaces that implement interfaces

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b7a58a1c0e4b30a75d97ef69a8593e1ce3914bf2"><tt>b7a58a1c</tt></a> Handle interfaces that implement interfaces

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ab8d62b67dd0dd9a27ad5320b3cb57b0bd76df51"><tt>ab8d62b6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1019">#1019</a> from 99designs/remove-source-reprinting</summary>

Remove source reprinting

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2f0fa0ef4d91a39e850e92e26d8450e5ea985bac"><tt>2f0fa0ef</tt></a> handle schema loading error better

- <a href="https://github.com/99designs/gqlgen/commit/aacc9b1fd6ff8fa91d4b4985b4c115796851b416"><tt>aacc9b1f</tt></a> Remove source reprinting

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e289aaa0b2ee378279710b018f0fb1c9c0da7997"><tt>e289aaa0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1018">#1018</a> from 99designs/federation-docs</summary>

Federation docs and examples

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3045b2cfadc80b4645dbb7228bf7e3cceb74d3c7"><tt>3045b2cf</tt></a> Federation docs and examples

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/656a07d1877081199058c3bc54ed521350a15e72"><tt>656a07d1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1016">#1016</a> from 99designs/federation-entity-type</summary>

Create a non generated federation _Entity type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8850a527a89f9b880531dc7fdcfec1016882b777"><tt>8850a527</tt></a> Create a non generated federation _Entity type

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1d41c2ebf22664de9660365fb4584c4ca1ac776c"><tt>1d41c2eb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1012">#1012</a> from 99designs/federation-config</summary>

Allow configuring the federation output file location

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/afa9a1504edb7da3e8055ed5923ce74647152396"><tt>afa9a150</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1013">#1013</a> from 99designs/feat-error-dispatch</summary>

propagate errors to response context in DispatchError

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/652aa2fb2917fc4c2362abc608e7643ee589daa7"><tt>652aa2fb</tt></a> propagate errors to response context in DispatchError

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0fe1af8c8c55cfcef938e59ee63fac3df1a319df"><tt>0fe1af8c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1011">#1011</a> from Khan/compound-keys</summary>

Compound key support in federation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ad3c1c818f86e8e610733e4be0cf6f755e60af25"><tt>ad3c1c81</tt></a> Allow configuring the federation output file location

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b4a00e6cfd6a7c0da51caa002bc72b7210610db9"><tt>b4a00e6c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1010">#1010</a> from Khan/query-exists</summary>

Make sure there's a Query node before trying to add a field to it.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/65401637e72c047c56ee2901d6341b9d52e60e14"><tt>65401637</tt></a> Adding type with multiple keys to federation test</summary>

Summary: The current federation test schema only has types with single keys (or no keys). Adding a type with multiple keys, including one non-String key, to test compound key federation code gen.

Test Plan: - go test

Reviewers: csilvers, miguel

Differential Revision: https://phabricator.khanacademy.org/D60715

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3f714a46146136c6e4adc4a77a4429488aa2f768"><tt>3f714a46</tt></a> Extending federation to support compound keys per Apollo spec</summary>

Summary:
Compound keys are not yet supported for federation in gqlgen. This diff adds support by modifying the federation plugin to handle a list of key fields on an entity rather than a single top-level key field. It will now look for "find<EntityName>By<KeyField1><KeyField2>..." in the resolver, rather than the original "Find<EntityName>By<KeyField>". The federation plugin does not yet support more complicated FieldSets in the key, such as nested selections.

References:
- Apollo federation spec: https://www.apollographql.com/docs/apollo-server/federation/federation-spec/
- Selection sets: https://graphql.github.io/graphql-spec/draft/#sec-Selection-Sets

Will update https://phabricator.khanacademy.org/D59469 with multiple key changes.

Test Plan:
- Tested Go GQL services using both single- and multiple-key federated types (assignments and content-library in webapp/services)
- Ran gqlgen on non-federated services in webapp to ensure regular generation still works (donations service)
- WIP: creating unit tests; will submit as separate diff

Reviewers: briangenisio, dhruv, csilvers, O4 go-vernors

Reviewed By: dhruv, csilvers, O4 go-vernors

Differential Revision: https://phabricator.khanacademy.org/D59569

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9f2a624bb4b57505d751d8d3e1d3248d11291c31"><tt>9f2a624b</tt></a> Make sure there's a Query node before trying to add a field to it.</summary>

Federation adds some queries to the schema.  There already existed
code to insert a Query node if none existed previously.  But that code
was only put on addEntityToSchema(), and not the other place we update
the query, addServiceToSchema().

Almost always the old code was good enough, since we call
addEntityToSchema() before addServiceToSchema().  But there's on
addServiceToSchema(), so we need to do the query-existence check there
too.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b941b970f0b67a4e102ee1635156a6a2b5a2863b"><tt>b941b970</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1007">#1007</a> from 99designs/handle-invalid-autoload-path</summary>

Give an appropriate error message when autoload isnt a valid package

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/95b1080954c5d447a6afd059a185256ae2e2ed1e"><tt>95b10809</tt></a> bump appveyor go version for consistent behavour

- <a href="https://github.com/99designs/gqlgen/commit/91a9ff97633e52f078687b79c7717f9f5129b9df"><tt>91a9ff97</tt></a> fix bad copy from template

- <a href="https://github.com/99designs/gqlgen/commit/d5d6f830475fcaa5790bb1d6390f3876cd9073a2"><tt>d5d6f830</tt></a> Give an appropriate error message when autoload isnt a valid package

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f7667e127039af2d7fb4252c0bb3a36554634f80"><tt>f7667e12</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1009">#1009</a> from 99designs/interface-regression</summary>

Interface regression

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ffc419f3053b18f6aead8eaba5d346c1ad9e31c5"><tt>ffc419f3</tt></a> Fix interfaces used as normal object types

- <a href="https://github.com/99designs/gqlgen/commit/44cfb92639db2e7f18a1b0a4091b29f120a0dbb4"><tt>44cfb926</tt></a> Test example for interface regression

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0ddb3ef308d1f801928e919cc478d5b8d6653458"><tt>0ddb3ef3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1006">#1006</a> from ravisastryk/entity-directives-lookup</summary>

skip searching directives when entity is found

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/395e1d731969d6a33bdbebc4fb2fa14be5fe2fe4"><tt>395e1d73</tt></a> skip searching directives when entity is found

- <a href="https://github.com/99designs/gqlgen/commit/e1f2282e1331fd38c01fe449f116d47d16583cd6"><tt>e1f2282e</tt></a> bump to go 1.13 in ci

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/34c92eba0b29a49cc592060890263d1677490a3d"><tt>34c92eba</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1003">#1003</a> from 99designs/fix-chat-example</summary>

fix chat example

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6bf88417d867bce8bc3eca3b52a62f56eee994b9"><tt>6bf88417</tt></a> fix chat example

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8ed2ec599b8faed3751177fd4335b1b3c3a79922"><tt>8ed2ec59</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/988">#988</a> from 99designs/package-cache</summary>

Cache all packages.Load calls in a central object

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9ccd7ed7199405ecbcc917ce601b83b72419009b"><tt>9ccd7ed7</tt></a> Cache all packages.Load calls in a central object

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/565619a80674052de24a4df5c9a2dba87ceb72df"><tt>565619a8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/993">#993</a> from 99designs/resolver-generator-v2</summary>

Resolver regenerator

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cf4a3eb455da1c482c20418deceeb44698a7f5bc"><tt>cf4a3eb4</tt></a> keep imports when scattering resolvers between files

- <a href="https://github.com/99designs/gqlgen/commit/da7c1e45788d1b9b2dca9dc080d70c96729a3c89"><tt>da7c1e45</tt></a> Update getting started docs

- <a href="https://github.com/99designs/gqlgen/commit/c233876e67aea3946d73890aeb9b96cfc4106edd"><tt>c233876e</tt></a> fix windows test paths

- <a href="https://github.com/99designs/gqlgen/commit/93713a291eb31279789332d316da7abd98c6276c"><tt>93713a29</tt></a> Add tests for code persistence

- <a href="https://github.com/99designs/gqlgen/commit/3e507e0dec0d2d29eedb2f067d138c509f91a991"><tt>3e507e0d</tt></a> separate resolver stubs by 1 empty line

- <a href="https://github.com/99designs/gqlgen/commit/8a208af51f6f808bbc8ef0b67aafb087e1112ca6"><tt>8a208af5</tt></a> add tests covering ResolverConfig

- <a href="https://github.com/99designs/gqlgen/commit/f8e6196164ee45c005be61e866e2d149d9a2fe9e"><tt>f8e61961</tt></a> set init to use new resolvers by default

- <a href="https://github.com/99designs/gqlgen/commit/dbaf355dbf2d6d2255e1934094fad2fbd69b0441"><tt>dbaf355d</tt></a> copy through any unknown data

- <a href="https://github.com/99designs/gqlgen/commit/e7255580193837f91f6f015b24700787a56017eb"><tt>e7255580</tt></a> copy old imports through before gofmt prunes

- <a href="https://github.com/99designs/gqlgen/commit/6ec365046295574d8e903ee93cf24b7d49b4b180"><tt>6ec36504</tt></a> Copy existing resolver bodies when regenerating new resolvers

- <a href="https://github.com/99designs/gqlgen/commit/9e3b399d4d9f18b1de1ec51acba09406ae9e56ad"><tt>9e3b399d</tt></a> add resolver layout = follow-schema

- <a href="https://github.com/99designs/gqlgen/commit/8a18895e1ec49e383ce2cda79d315c78f5a701ca"><tt>8a18895e</tt></a> Update to latest golangci-lint

- <a href="https://github.com/99designs/gqlgen/commit/f7a67722a6baf2612fa429bd21ceb9c6b9cbed1c"><tt>f7a67722</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/985">#985</a> from Khan/no-key-needed

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fa884991b5faec01f1ded737a64350d782e3a418"><tt>fa884991</tt></a> Correctly generate a federated schema when no entity has a `[@key](https://github.com/key)`.</summary>

Normally, when a service is taking part in graphql federation, it will
services can link to (that is, have an edge pointing to) the type that
this service provides.  The previous federation code assumed that was
the case.

types.  It might seem that would mean the service is unreachable,
since there is no possibility of edges into the service, but there are
and top level Mutation edges.  That is, if a service only provides a
top-level query or top-level mutation, it might not need to define a

This commit updates the federation code to support that use case.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/36aae4aa277847bc5dc2d4fcec5ae0c1d7a1d686"><tt>36aae4aa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/994">#994</a> from 99designs/feat-cache-ctx</summary>

Add context.Context to graphql.Cache interface's methods

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/61e060bdfe4559138a08a00683815a315f85a154"><tt>61e060bd</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/995">#995</a> from alexsn/directiveroot_empty_lines</summary>

Remove empty lines on DirectiveRoot generation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/30c295c4014e2010ef8ed3356890c961247b379b"><tt>30c295c4</tt></a> Remove empty lines on DirectiveRoot generation

- <a href="https://github.com/99designs/gqlgen/commit/85cfa8a3afffad11a99bf4205310f46987f3329d"><tt>85cfa8a3</tt></a> Add context.Context to graphql.Cache interface's methods

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a6c7aafb928f648d0a8106c0a42554abdce53952"><tt>a6c7aafb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/931">#931</a> from fridolin-koch/master</summary>

Fix for Panic if only interfaces shall be generated

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ec4f6b151d4c14d704f27ae7fe341f7ad5ad4883"><tt>ec4f6b15</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/989">#989</a> from 99designs/fix-intermittent-test-ka-failure</summary>

Fix intermittent websocket ka test failure

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/76035df5e63c580004440762edbf6779fe9243db"><tt>76035df5</tt></a> Fix intermittent websocket ka test failure

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/aa407b1f3553ac2aee1939fbe28c85ed5cbfcdf9"><tt>aa407b1f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/979">#979</a> from 99designs/capture-read-times</summary>

Capture read times

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4dd1008659429e99e94e5da0f3401e358a16b69e"><tt>4dd10086</tt></a> fix test race by only stubbing now where we need to

- <a href="https://github.com/99designs/gqlgen/commit/8dbce3cf161f19c132d3cf29aa95851732c7f922"><tt>8dbce3cf</tt></a> Capture the time spent reading requests from the client

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c6b3e2a1ef220cd122ad3c2a6e25bc74c89a7a4c"><tt>c6b3e2a1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/983">#983</a> from vikstrous/name-for-package-global</summary>

single packages.Load for NameForPackage

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ae79e75bc2d8296551e8b88b7b3f8596f038ca94"><tt>ae79e75b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/978">#978</a> from 99designs/pluggable-error-code</summary>

Allow customizing http and websocket status codes for errors

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7f6f1667bd06e4a5f18128c592ef96e01bca97b6"><tt>7f6f1667</tt></a> bump x/tools for consistent import formatting

- <a href="https://github.com/99designs/gqlgen/commit/842fcc11b1481bdb04d2bd711a1c091354b7a96e"><tt>842fcc11</tt></a> review feedback

- <a href="https://github.com/99designs/gqlgen/commit/f0bea5ffcbdfbf231a6d2848b77f0e9c20288702"><tt>f0bea5ff</tt></a> Allow customizing http and websocket status codes for errors

- <a href="https://github.com/99designs/gqlgen/commit/bd50bbcbb3d96bc168c1b5186147be14487e0cc6"><tt>bd50bbcb</tt></a> single packages.Load for NameForPackage

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/28c032d1f3ba55761dbac0cf846c9c66b7abb5e8"><tt>28c032d1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/982">#982</a> from DavidJFelix/patch-1</summary>

fix: explicitly exclude trailing comma from link

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ac67050a7156be5e109d6b67b3fc7073e05fdeb1"><tt>ac67050a</tt></a> fix: explicitly exclude trailing comma from link</summary>

- this looks dumb, but when the page is rendered, the link resolves with the comma, despite the comma being excluded in github rendering.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4e95b363e8799ddf92af4b06cca058342164bfd1"><tt>4e95b363</tt></a> fix some version switcher paths

- <a href="https://github.com/99designs/gqlgen/commit/08369dfe093d010b8a3c81245265232c42f909db"><tt>08369dfe</tt></a> add missing trailing slash on paths

- <a href="https://github.com/99designs/gqlgen/commit/ea347ca7c0a9f9ee5756d8e44b0dbc48e0bd6ed6"><tt>ea347ca7</tt></a> fetch all tags

- <a href="https://github.com/99designs/gqlgen/commit/8c1a8f5777a3b63ce0bde47538fd04bd4c9aa3b3"><tt>8c1a8f57</tt></a> fix branch switching

- <a href="https://github.com/99designs/gqlgen/commit/324efc5cb537cee8df072e097d83e0bf0f57abd8"><tt>324efc5c</tt></a> add origin if missing

- <a href="https://github.com/99designs/gqlgen/commit/cfa2907a017d517ce90a62c4ef978f7bce66e9b7"><tt>cfa2907a</tt></a> Generate docs for all tags

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8218c734fb126b49882b44801228b36f033909d2"><tt>8218c734</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/851">#851</a> from marwan-at-work/federation</summary>

Apollo Federation MVP

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/48dc29c19314cc9f7cbf6c58a9514245fa07a1b6"><tt>48dc29c1</tt></a> go 1.12 generate, 1.14 failed

- <a href="https://github.com/99designs/gqlgen/commit/b2e81787dbd76b43b2969e556dfda10076bdeaf8"><tt>b2e81787</tt></a> update gqlparse to v1.2.1

- <a href="https://github.com/99designs/gqlgen/commit/d2a13d33cdcb27e8e141c54e1ab1fa0aebad2d2b"><tt>d2a13d33</tt></a> update go.mod

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0eef2fe2cf5123990844a6c2bbb2418b044df1e6"><tt>0eef2fe2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/970">#970</a> from spiffyjr/master</summary>

Fix extra trimspace on nillable Unmarshals

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/56b8eef2b7ca0852fff3f3cda80a08df3569868e"><tt>56b8eef2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/974">#974</a> from oshalygin/docs/gqlgen-pg-example-repo</summary>

Add Link to Sample Project with GQLGen and Postgres

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f49936eb9a89336beb3677490ef687fb7b6a408e"><tt>f49936eb</tt></a> Add Link to Sample Project with GQLGen and Postgres</summary>

This is a very straightforward project with numerous details in the README and the official
documentation, but questions continue to pop up around how to use this project, organize the files
and ultimately make data calls to some persistent layer.

The `https://github.com/oshalygin/gqlgen-pg-todo-example` was built in order to show newcomers the
following:
- How to organize their graphql schema, resolvers, dataloaders and models
- How to create a new dataloader
- How to resolve with a dataloader and how to avoid some of the pitfalls(inconsistent db query to keys array order)
- How to map models from a gql schema to structs

While the examples in this project are helpful, they could benefit from more elaborate explanations in the
code as well as the README to help newcomers get started.  This PR is not intended to portray any of the examples
negatively and should not be interpreted as such.  There are many findings/lessons learned from the work that folks
put together in those examples.

README which covers a ton of the details on how to use this project:
- [README](https://github.com/oshalygin/gqlgen-pg-todo-example)

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/db499561277294e5fc368a45fe49a678ba217fe9"><tt>db499561</tt></a> force rebuild

- <a href="https://github.com/99designs/gqlgen/commit/0985a78e18d4f2c629445d4bef9958fed85c5e09"><tt>0985a78e</tt></a> remove debug comments

- <a href="https://github.com/99designs/gqlgen/commit/7f64842577435c7246e5448acc1c1238221140d6"><tt>7f648425</tt></a> add preliminary test_data

- <a href="https://github.com/99designs/gqlgen/commit/c9d6d94b7465d1026d4a75f204cf0c645aa66972"><tt>c9d6d94b</tt></a> add preliminary tests

- <a href="https://github.com/99designs/gqlgen/commit/2345936ea1c0066052bc823c851bad8b8885c888"><tt>2345936e</tt></a> fix integration

- <a href="https://github.com/99designs/gqlgen/commit/aae7486d5fbfff6bc7ca696df2949b48e6c4b80f"><tt>aae7486d</tt></a> go generate

- <a href="https://github.com/99designs/gqlgen/commit/555a95462cec99c7b0549383102238d5b373ed71"><tt>555a9546</tt></a> go generate + remove directives nil check

- <a href="https://github.com/99designs/gqlgen/commit/368d546dc2f6a240ab0cbb165d5dd6bf202be747"><tt>368d546d</tt></a> Apollo Federation MVP

- <a href="https://github.com/99designs/gqlgen/commit/21e0e6762eef0dbdea1b0242b887cd5adad35a4d"><tt>21e0e676</tt></a> Fix extra trimspace on nillable Unmarshals

- <a href="https://github.com/99designs/gqlgen/commit/f869f5a85385745d5854daaa25eab5571b04b245"><tt>f869f5a8</tt></a> remove deprected handler call

- <a href="https://github.com/99designs/gqlgen/commit/f0b83cb16c618ddcad4c26a779d97868dbb9c8a8"><tt>f0b83cb1</tt></a> fix merge conflict

- <a href="https://github.com/99designs/gqlgen/commit/cdf967214d9e800c48ba55ac41e060b1107b0a53"><tt>cdf96721</tt></a> update generated code

- <a href="https://github.com/99designs/gqlgen/commit/21356ce35cc55896ee855c9b3238aa00684ac242"><tt>21356ce3</tt></a> markdown cleanup

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/412a72fe26b093b08d27e90adf0390ad0ea0a7ea"><tt>412a72fe</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/885">#885</a> from 99designs/handler-refactor</summary>

Refactor handler package

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bac79c54bb58d0c7450a6e0f98371c6792ae3a3a"><tt>bac79c54</tt></a> force clean git checkout

- <a href="https://github.com/99designs/gqlgen/commit/dca9e4a5b04f34f1bba32d472c0075ee9d0ea476"><tt>dca9e4a5</tt></a> Add migration docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5106480b4c6332c9488b88b4f9a66b29a666948b"><tt>5106480b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/947">#947</a> from 99designs/handler-oc-handling</summary>

always return OperationContext for postpone process

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/922db1e3182f119f7c9608f71be4503ad69fad56"><tt>922db1e3</tt></a> always return OperationContext for postpone process

- <a href="https://github.com/99designs/gqlgen/commit/8794f03e783d6319b8a9066b4bfe93d651a4a594"><tt>8794f03e</tt></a> v0.10.2 postrelease bump

- <a href="https://github.com/99designs/gqlgen/commit/14dbf1aae6d512f6be5affd0e29686f360eb5579"><tt>14dbf1aa</tt></a> use new handler package in new test

- <a href="https://github.com/99designs/gqlgen/commit/a339a0423ee7cf9562c1767664aad6f15e1b5677"><tt>a339a042</tt></a> panic if operation context is missing when requested

- <a href="https://github.com/99designs/gqlgen/commit/a13a0f5f83eb2c71dbcd658fe2abe263aca0fb2d"><tt>a13a0f5f</tt></a> add docs on extension name conventions

- <a href="https://github.com/99designs/gqlgen/commit/458fa0deda7b9b45fe7b45d098bb82b15aa207d7"><tt>458fa0de</tt></a> Add more interface assertions

- <a href="https://github.com/99designs/gqlgen/commit/d0836b72d5fe632808afff9254e415948ff11680"><tt>d0836b72</tt></a> Expose APQ stats

- <a href="https://github.com/99designs/gqlgen/commit/cf14cf103cff2d6d52146c7f30f46bd6a329aa59"><tt>cf14cf10</tt></a> fix: Fix no code generation for only interfaces

- <a href="https://github.com/99designs/gqlgen/commit/dc76d029e2aea0aa71efb032923d660b809068b4"><tt>dc76d029</tt></a> Merge remote-tracking branch 'origin/master' into handler-refactor

- <a href="https://github.com/99designs/gqlgen/commit/572fb419fc66c56ac103abf65bb54282d46aec29"><tt>572fb419</tt></a> remove all references to deprecated handler package

- <a href="https://github.com/99designs/gqlgen/commit/dc6223462ed4dc85b66e836feb5c5ee58bc363bd"><tt>dc622346</tt></a> Tune allocs for benchmarks

- <a href="https://github.com/99designs/gqlgen/commit/a6f9462634b3408cdd0a8e0c0c14ff86d8f45317"><tt>a6f94626</tt></a> Merge remote-tracking branch 'origin/master' into handler-refactor

- <a href="https://github.com/99designs/gqlgen/commit/c3f938108d172b7097828d1ae2a7ce940b611ae6"><tt>c3f93810</tt></a> fix benchmark

- <a href="https://github.com/99designs/gqlgen/commit/631b48a56ac159cd03dbc67e15f4a8dfef7dc266"><tt>631b48a5</tt></a> remove automatic field stat collection to reduce time calls

- <a href="https://github.com/99designs/gqlgen/commit/a77d9fc29019a78bd0f4459251cf24ce3083723b"><tt>a77d9fc2</tt></a> Add generated stanzas back in

- <a href="https://github.com/99designs/gqlgen/commit/0ee185b811db5d81286e02d90e2bf5cfbd424b9a"><tt>0ee185b8</tt></a> fix duplicate header sends

- <a href="https://github.com/99designs/gqlgen/commit/7cbd75db593854ac3be66fa66a3248a4b0acf6a9"><tt>7cbd75db</tt></a> fix APQ signature

- <a href="https://github.com/99designs/gqlgen/commit/67fa21049567aff73f5ae019683a14b7b03a496d"><tt>67fa2104</tt></a> allow extensions to declare their own stats

- <a href="https://github.com/99designs/gqlgen/commit/e9502ae042f901e85731b39316d5d4687e3709f9"><tt>e9502ae0</tt></a> Make extensions validatable

- <a href="https://github.com/99designs/gqlgen/commit/fc727c9cd7663a874d3dcdccdbb096f385934dfd"><tt>fc727c9c</tt></a> Add a signpost method to handler extension interface

- <a href="https://github.com/99designs/gqlgen/commit/0a39ae206916d607fcf8c7fb07a167d0b48c8933"><tt>0a39ae20</tt></a> add fixed complexity limit

- <a href="https://github.com/99designs/gqlgen/commit/f2ef5ec3d660c7c226ad5392d9d43e38abcc6827"><tt>f2ef5ec3</tt></a> more deprecations and more compat

- <a href="https://github.com/99designs/gqlgen/commit/2898a622b48dd484e2c2365b9655f3b068ba524d"><tt>2898a622</tt></a> rename ResolverContext to FieldContext

- <a href="https://github.com/99designs/gqlgen/commit/092ed95fd7cc9587b88865175a055d8b35a9fb42"><tt>092ed95f</tt></a> collect field timing in generated code

- <a href="https://github.com/99designs/gqlgen/commit/848c627c375d9cc991131d94bbb13f317a723ecf"><tt>848c627c</tt></a> remove DirectiveMiddleware

- <a href="https://github.com/99designs/gqlgen/commit/40f088681b169ef1249258f164b70f37ae7b66b6"><tt>40f08868</tt></a> add NewDefaultServer

- <a href="https://github.com/99designs/gqlgen/commit/1b57bc3eda296ffe15788626c628c6b0a092ced9"><tt>1b57bc3e</tt></a> Rename RequestContext to OperationContext

- <a href="https://github.com/99designs/gqlgen/commit/3476ac44bf70da28b606bd92f7d7e8dd68b371e7"><tt>3476ac44</tt></a> fix linting issues

- <a href="https://github.com/99designs/gqlgen/commit/479abbef5e8f6a00b23db8c5f68e04a6fcd3e9a8"><tt>479abbef</tt></a> update generated code

- <a href="https://github.com/99designs/gqlgen/commit/bc98156929b06eeca1d95f475501470bd9034c2d"><tt>bc981569</tt></a> Combine root handlers in ExecutableSchema into a single Exec method

- <a href="https://github.com/99designs/gqlgen/commit/473a0d256af2a23ed5b59facd0850ccf137d6fa8"><tt>473a0d25</tt></a> Implement bc shim for old handler package

- <a href="https://github.com/99designs/gqlgen/commit/631142cfacc0d6b4c7773fea7750bd9d92cea4c4"><tt>631142cf</tt></a> move writer all the way back to the transport

- <a href="https://github.com/99designs/gqlgen/commit/c7bb03a8fd9b43c47f52b1beafb872d7aee65280"><tt>c7bb03a8</tt></a> merge executable schema entrypoints

- <a href="https://github.com/99designs/gqlgen/commit/e7e913d901fa72237cb1ee3dc22b7530238c8532"><tt>e7e913d9</tt></a> Remove remains of old handler package

- <a href="https://github.com/99designs/gqlgen/commit/8c5340c1ab61c43dc5d1e3fd6feb30d24d95cdb9"><tt>8c5340c1</tt></a> Add complexity limit plugin

- <a href="https://github.com/99designs/gqlgen/commit/0965420a4246492bbac6922da742b157ea968c29"><tt>0965420a</tt></a> Add query document caching

- <a href="https://github.com/99designs/gqlgen/commit/aede7d1cf15f054b1762f9801337bd3e8764b54d"><tt>aede7d1c</tt></a> Add multipart from transport

- <a href="https://github.com/99designs/gqlgen/commit/64cfc9add38004e8741fbe3bdd5a247c61718d80"><tt>64cfc9ad</tt></a> extract shared handler test server stubs

- <a href="https://github.com/99designs/gqlgen/commit/a70e93bcae24130ef3746d89afa48b23e96f4787"><tt>a70e93bc</tt></a> consistently name transports

- <a href="https://github.com/99designs/gqlgen/commit/9d1d77e67df3fd2c75646af7b5de361b9cbe8482"><tt>9d1d77e6</tt></a> split context.go into 3 files

- <a href="https://github.com/99designs/gqlgen/commit/72c47c985f2727ab7d9dff7f6ebaf3614c44a507"><tt>72c47c98</tt></a> rename result handler to response handler

- <a href="https://github.com/99designs/gqlgen/commit/4a69bcd034ade82bacf6b71b4945f4917d2fdfc1"><tt>4a69bcd0</tt></a> Bring operation middleware inline with other handler interfaces

- <a href="https://github.com/99designs/gqlgen/commit/ab5665add4f1b6effe21cb1ef77f7346cad1d59c"><tt>ab5665ad</tt></a> Add result context

- <a href="https://github.com/99designs/gqlgen/commit/c3dbcf83eaa8bc865b7e482fd17f26fa5139485b"><tt>c3dbcf83</tt></a> Add apollo tracing

- <a href="https://github.com/99designs/gqlgen/commit/f00e5fa0791be8e8909923711f7ebde8d2e74c15"><tt>f00e5fa0</tt></a> use plugins instead of middleware so multiple hooks can be configured

- <a href="https://github.com/99designs/gqlgen/commit/a7c5e6600729012283270ad2c653de57772eba6b"><tt>a7c5e660</tt></a> build middleware graph once at startup

- <a href="https://github.com/99designs/gqlgen/commit/2e0c9cab65d4c6a0cd8237a1b15f55a075afb44f"><tt>2e0c9cab</tt></a> mark validation and parse errors separately to execution errors

- <a href="https://github.com/99designs/gqlgen/commit/cb99b42ed0e4974aeb1fc2d9cee43d061c1152cf"><tt>cb99b42e</tt></a> Add websocket transport

- <a href="https://github.com/99designs/gqlgen/commit/eed1515c7abeb08a00291e9ab241f88af860d8aa"><tt>eed1515c</tt></a> Split middlware out of handler package

- <a href="https://github.com/99designs/gqlgen/commit/b5089cac400ddf2ffb00d75c849656731d2cb29e"><tt>b5089cac</tt></a> Split transports into subpackage

- <a href="https://github.com/99designs/gqlgen/commit/d0f683034fbf877457990060a8c2423b1ccfce0d"><tt>d0f68303</tt></a> port json post

- <a href="https://github.com/99designs/gqlgen/commit/afe241b56cd44394a2b32447f7d817a8361f909d"><tt>afe241b5</tt></a> port over tracing

- <a href="https://github.com/99designs/gqlgen/commit/311887d6a9336c1c5f6f9a59752a94afa6be5b52"><tt>311887d6</tt></a> convert APQ to middleware

- <a href="https://github.com/99designs/gqlgen/commit/da986181d7e6ca9da2999fb62d8fbc7c33eda21f"><tt>da986181</tt></a> port over the setter request context middleware

- <a href="https://github.com/99designs/gqlgen/commit/249b602d487fd189787bcd3605ff4c3a459771e9"><tt>249b602d</tt></a> Start drafting new handler interfaces

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.10.2"></a>
## [v0.10.2](https://github.com/99designs/gqlgen/compare/v0.10.1...v0.10.2) - 2019-11-28
- <a href="https://github.com/99designs/gqlgen/commit/f276a4e6773992c572119b22821d375ad008c53d"><tt>f276a4e6</tt></a> release v0.10.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9e989d946989e941985a62c1497d3d2d0abd856c"><tt>9e989d94</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/929">#929</a> from nmaquet/check-nil-interface-ptrs</summary>

Don't crash when interface resolver returns a typed nil

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6f20101c40adf13a6ceef483ce0158b83273afed"><tt>6f20101c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/940">#940</a> from vikstrous/optional-modelgen</summary>

make model generation optional

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9b9dd5620e65c4d1632c71c52acf1c3c12e7ca3d"><tt>9b9dd562</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/942">#942</a> from vikstrous/disable-validation</summary>

add skip_validation flag

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f9f2063a5f77a5cb21d30db1f038d17242f2dbd9"><tt>f9f2063a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/941">#941</a> from vikstrous/qualify-package-path-faster</summary>

shortcut QualifyPackagePath in go module mode

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4db0e6eccc8745ed765f3863d221ea13c57f0bd1"><tt>4db0e6ec</tt></a> keep function private

- <a href="https://github.com/99designs/gqlgen/commit/c06f05b319fc9287110ac0dce2f7e4aafbd34873"><tt>c06f05b3</tt></a> add doc

- <a href="https://github.com/99designs/gqlgen/commit/bd353b3e227f9dadd923ed49bc1a0ddb0c043865"><tt>bd353b3e</tt></a> add skip_validation flag

- <a href="https://github.com/99designs/gqlgen/commit/b829628d3186975544ee375cc23bf5cb778965a5"><tt>b829628d</tt></a> shortcut QualifyPackagePath in go module mode

- <a href="https://github.com/99designs/gqlgen/commit/3a05d2dd985ee4f1e2d3390a65d4a24447a5ecb4"><tt>3a05d2dd</tt></a> add mention in the docs

- <a href="https://github.com/99designs/gqlgen/commit/c2c2d7de0cf8dfb232e33c619d72e85e70e656b8"><tt>c2c2d7de</tt></a> make model generation optional

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d3f6384425e61f39d58819e7fc893b55cfd00d21"><tt>d3f63844</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/939">#939</a> from mjarkk/patch-1</summary>

(docs) graph-gophers now supports Struct Field resolving

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ba3d018929670a1831e58086660c6e062704d815"><tt>ba3d0189</tt></a> graph-gophers now supports Struct Field resolvers

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e747d923d3d0c587a59b5586e1a6dddb2f0f3a7f"><tt>e747d923</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/938">#938</a> from lulucas/master</summary>

modelgen hook docs fixed

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/63be1d5e2a0365cc1eaf23c57826ec47f77eb730"><tt>63be1d5e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1">#1</a> from lulucas/modelgen-hook-patch-1</summary>

modelgen hook docs use plugin poitner

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/33fc16b1a75e745e28378452847f45a943fdd238"><tt>33fc16b1</tt></a> modelgen hook docs use plugin poitner</summary>

and add modelgen package to ModelBuild type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fcfe595e65660e5779a113bc0f9e68e0af750821"><tt>fcfe595e</tt></a> Add a comment

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/599460871b52f129d653b7a9befd4773fe3b5acd"><tt>59946087</tt></a> Add unit test for the interface resolver / typed nil interaction</summary>

This added test shows that the `_Dog_species` automatically generated
resolver will crash unless the extra nil check is added in
`interface.gotpl`.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/201768f0b3590dd4e8e7c52fb506c1bed90abc40"><tt>201768f0</tt></a> Regenerate examples

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/85ca9efe5cfdbf111fe1b4d57fbb38eae151dbfb"><tt>85ca9efe</tt></a> Return graphql.Null in interface resolver when passed a typed nil</summary>

Go's dreaded _typed nil_ strikes again. Nil pointers of struct types
aren't equal to nil interface pointers.

See https://golang.org/doc/faq#nil_error

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/15b30588a1451bbe280660a1d6cf629f50121d86"><tt>15b30588</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/894">#894</a> from 99designs/enum-var-value-coercion</summary>

Improve enum value (with vars) validation timing

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/568433a23bd2c123460a389fcb2e2e03dfe61ef2"><tt>568433a2</tt></a> fix ci failed

- <a href="https://github.com/99designs/gqlgen/commit/0ccfc7e0ebadffe8d59300a2b05dc8cfaa78d5a8"><tt>0ccfc7e0</tt></a> Merge branch 'master' into enum-var-value-coercion

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9cfd817e013b951206bc969ba517c98ff208a11c"><tt>9cfd817e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/897">#897</a> from mskrip/modelgen-hook</summary>

Add possibility to hook into modelgen plugin

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c1e6414834344c20728f9a31b74dacf312713516"><tt>c1e64148</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/900">#900</a> from zannen/master (closes <a href="https://github.com/99designs/gqlgen/issues/896"> #896</a>)

- <a href="https://github.com/99designs/gqlgen/commit/8a8f0a0f8de1c10e1e8749d12108038ed5eac452"><tt>8a8f0a0f</tt></a> Add autogenerated files (<a href="https://github.com/99designs/gqlgen/pull/896">#896</a>)

- <a href="https://github.com/99designs/gqlgen/commit/531729df1303cde002c0bede5a8a0cb11ac4abda"><tt>531729df</tt></a> Move test schema file from example dir into codegen/testserver (<a href="https://github.com/99designs/gqlgen/pull/896">#896</a>)

- <a href="https://github.com/99designs/gqlgen/commit/5144775f6c57ab7c0ca2b8eaa2441ed398042e40"><tt>5144775f</tt></a> Add example to check for regression of <a href="https://github.com/99designs/gqlgen/pull/896">#896</a>

- <a href="https://github.com/99designs/gqlgen/commit/3b5df4ceec3629694cf8ba3f0c62eac8dd66e82e"><tt>3b5df4ce</tt></a> Add check for obviously different TypeReferences (<a href="https://github.com/99designs/gqlgen/pull/896">#896</a>)

- <a href="https://github.com/99designs/gqlgen/commit/fb96756a2095523acb9b59e219eb5861ca41e588"><tt>fb96756a</tt></a> Update generated content (<a href="https://github.com/99designs/gqlgen/pull/896">#896</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fd201a8c8b0f6d95fa4593c8f7cdf3b629f062aa"><tt>fd201a8c</tt></a> Update UniquenessKey for when Element is/isn't nullable (<a href="https://github.com/99designs/gqlgen/pull/896">#896</a>)</summary>

With a schema:
type Query {
  things1: [Thing] # Note the lack of "!"
}

type Subscription {
  things2: [Thing!] # Note the "!"
}

the UniquenessKey for the two lists is the same, which causes non-deterministic output.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2a269dd3e4303fa748be4e2ec3d264f8a29bd6fd"><tt>2a269dd3</tt></a> Add modelgen hook recipe

- <a href="https://github.com/99designs/gqlgen/commit/6ceb76b632ec117b8b22adaddae68ad7f56e36df"><tt>6ceb76b6</tt></a> Test tag generation only by looking up extected tag strings

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1f272d1bd14a9a5b0457238ba28f8070d7352307"><tt>1f272d1b</tt></a> Add possibility to hook into modelgen plugin (closes <a href="https://github.com/99designs/gqlgen/issues/876"> #876</a>)</summary>

This change introduces option to implement custom hook for model
generation plugin without the need to completly copy the whole `modelgen` plugin.

that can be:

```golang
func mutateHook(b *ModelBuild) *ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			field.Tag += ` orm_binding:"` + model.Name + `.`  +  field.Name + `"`
		}
	}

	return b
}

...

func main() {
    p := modelgen.Plugin {
        MutateHook: mutateHook,
    }

    ...
}

```

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/99a55da2cdb986686f72dd6d4c6841dc1a79c688"><tt>99a55da2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/927">#927</a> from matiasanaya/feature/bind-to-embedded-interface</summary>

Bind to embedded interface

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/70e860cce0a3c943f34b796351cc956fa48ab900"><tt>70e860cc</tt></a> Bind to embedded interface method

- <a href="https://github.com/99designs/gqlgen/commit/a745dc7807357e9064292a7978d80fe85c6794cd"><tt>a745dc78</tt></a> Fixes <a href="https://github.com/99designs/gqlgen/pull/843">#843</a>: Bind to embedded struct method or field

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f80cab0662d30abd847bfcb012a3d54e7fe4d8bb"><tt>f80cab06</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/923">#923</a> from 99designs/gqlparser-1.2.0</summary>

Update to gqlparser-1.2.0

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7508f4e560194d17862ae269df40e7cf1968698e"><tt>7508f4e5</tt></a> Update to gqlparser-1.2.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7653a681a9696ba5a2562c976d298e38a408ba1b"><tt>7653a681</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/916">#916</a> from karthikraobr/patch-1</summary>

3->4 scalars

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8faa0e3aad002970214f2e04de0fb3f3186c13ec"><tt>8faa0e3a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/917">#917</a> from colelawrence/patch-1</summary>

docs: Fix typo in title of "Resolvers"

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f7d888f9e95d076fde232362e371912bda070ddd"><tt>f7d888f9</tt></a> Merge branch 'master' into patch-1

- <a href="https://github.com/99designs/gqlgen/commit/d722ac66368529b8f57e0a5c1feea635a2b1bbbe"><tt>d722ac66</tt></a> Update scalars.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1172128c3c7f231e8ca1654a0c978b1f3447736e"><tt>1172128c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/904">#904</a> from cfilby/fix-config-docs</summary>

Minor Documentation Tweaks

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/935f11eda72a3d0761ade303cbc324ab9f4098da"><tt>935f11ed</tt></a> Fix typo in title

- <a href="https://github.com/99designs/gqlgen/commit/026d029cfa1ad86a22363efab43338fd231a2420"><tt>026d029c</tt></a> 3->4 scalars

- <a href="https://github.com/99designs/gqlgen/commit/5eb6bef6f515b31ca5c539ede824241b1befb75f"><tt>5eb6bef6</tt></a> Fix weird indending

- <a href="https://github.com/99designs/gqlgen/commit/756dcf6bb3be6680d6574b01eb31fd378c225bdf"><tt>756dcf6b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/907">#907</a> from lian-yue/patch-1 (closes <a href="https://github.com/99designs/gqlgen/issues/860"> #860</a>)

- <a href="https://github.com/99designs/gqlgen/commit/2a943eed912ab9557d124ea1ab7abfc9dd9fa8e8"><tt>2a943eed</tt></a> Update directive.go (closes <a href="https://github.com/99designs/gqlgen/issues/860"> #860</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/adbceeea04b3b8c27b92760da4bd1a8beae0a913"><tt>adbceeea</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/902">#902</a> from cfilby/fix-int64-marshalling</summary>

Add support for int64 IDs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/13c3d9224b184f8648ff78b6fe860e4dc51aa913"><tt>13c3d922</tt></a> Update id function

- <a href="https://github.com/99designs/gqlgen/commit/37191779306d628648fc888b8cc3dd83e4eb7f3c"><tt>37191779</tt></a> Add more tests

- <a href="https://github.com/99designs/gqlgen/commit/0968e0cbfb660f909f96e3833204f7e17b3c2268"><tt>0968e0cb</tt></a> Fix VSCode Weirdness, validate formatting

- <a href="https://github.com/99designs/gqlgen/commit/a20c96d51a4ee0a2715af0bf3a9e6ba7f8fb5327"><tt>a20c96d5</tt></a> More edits

- <a href="https://github.com/99designs/gqlgen/commit/e9e88b41e0ee2fcb6e3ebd61cea41e5811fcb544"><tt>e9e88b41</tt></a> Stop double indending

- <a href="https://github.com/99designs/gqlgen/commit/9f4df68edc84a087ecab4beaf0cbf4f15a7ffd0b"><tt>9f4df68e</tt></a> More minor doc fixes

- <a href="https://github.com/99designs/gqlgen/commit/7abf0ac3d3bf03a2c66bf2729da5a0b07ba94a11"><tt>7abf0ac3</tt></a> Fix documentation bug

- <a href="https://github.com/99designs/gqlgen/commit/e9730ab90b17ca10ab85b866481f934ea5e957e8"><tt>e9730ab9</tt></a> gofmt

- <a href="https://github.com/99designs/gqlgen/commit/c3930f57e3b4164e1627a8d117872a5bc54599b6"><tt>c3930f57</tt></a> Remove redundant paren, add test

- <a href="https://github.com/99designs/gqlgen/commit/395fc85e02c2be3a6d1b67169919e2555f2b74de"><tt>395fc85e</tt></a> Add support for int64 ids

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/dbc88428d2d13c2de3554ffb361c09c90ac21474"><tt>dbc88428</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/889">#889</a> from thnt/fix-init-with-schema-arg</summary>

fix init not use custom schema filename

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fc4e513fd91773f4bf4707328a0a7c749c7fab64"><tt>fc4e513f</tt></a> add test for https://github.com/vektah/gqlparser/pull/109

- <a href="https://github.com/99designs/gqlgen/commit/dd98bb13d9a3ae85f9afa525091b8c0c1c2fa7c8"><tt>dd98bb13</tt></a> fix init not use custom schema

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4c35356cbe7bf886fd8c59c9f754d2d98f6987b8"><tt>4c35356c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/883">#883</a> from 99designs/handle-invalid-types</summary>

Gracefully handle invalid types from invalid go packages

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/25b7027118f99c097255d0d11e7384898d65b471"><tt>25b70271</tt></a> Gracefully handle invalid types from invalid go packages

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/046054dbda38cc50f8f1e2e3c6073cbcc315c2b1"><tt>046054db</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/882">#882</a> from 99designs/testserver-autobind</summary>

Use autobinding in testserver

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/12c963a4f3b88e3545f5b878216e54f1c2d6b32d"><tt>12c963a4</tt></a> Use autobinding in testserver

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/305116a0a0dc48fff6486a78642e77058365a41c"><tt>305116a0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/879">#879</a> from coderste/patch-1</summary>

Fixed broken GitHub link within the APQ page

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b4867b3f6934446c46845957981d62d2708d8343"><tt>b4867b3f</tt></a> Fixed broken GitHub link within the APQ page</summary>

Small documentation change to fix a broken GitHub link.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9f6b0ee4f5822a9af90b753f36bc01d4a2cfe0a4"><tt>9f6b0ee4</tt></a> v0.10.1 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.10.1"></a>
## [v0.10.1](https://github.com/99designs/gqlgen/compare/v0.10.0...v0.10.1) - 2019-09-25
- <a href="https://github.com/99designs/gqlgen/commit/efb6efe06c6e4fc706440acebf6f81fff85f295c"><tt>efb6efe0</tt></a> release v0.10.1

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/955f3499b245507e817e74417f8179a18b18eb81"><tt>955f3499</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/877">#877</a> from 99designs/fix-websocket-client</summary>

Fix websocket connections on test client

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ef24a1cc1e144f73c9fc71eb514ea67478ac504c"><tt>ef24a1cc</tt></a> Fix websocket connections on test client

- <a href="https://github.com/99designs/gqlgen/commit/c997ec0c922b724b6752a87e0759e9b387ca052e"><tt>c997ec0c</tt></a> v0.10.0 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.10.0"></a>
## [v0.10.0](https://github.com/99designs/gqlgen/compare/v0.9.3...v0.10.0) - 2019-09-24
- <a href="https://github.com/99designs/gqlgen/commit/75a837522ff029e1d0c5349922182c14023649ef"><tt>75a83752</tt></a> release v0.10.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0bc3cc86fae5aef301b93a6e206cb275a053b2a1"><tt>0bc3cc86</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/875">#875</a> from 99designs/fix-clientwide-opts</summary>

Fix client global options

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b43edf5d613c79cbc3d4e26a9446a81f80437a07"><tt>b43edf5d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/874">#874</a> from 99designs/configurable-slice-element-pointers</summary>

Add config option to omit pointers to slice elements

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/921aa9cf9055575d19d967e976eabe5e6aee2872"><tt>921aa9cf</tt></a> Fix client global options

- <a href="https://github.com/99designs/gqlgen/commit/d0098e60acc03cc7314bfc743579be75c46625a8"><tt>d0098e60</tt></a> Add config option to omit pointers to slice elements

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0189328068eb49438a131715a6d9354dc30731db"><tt>01893280</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/819">#819</a> from 99designs/fix-directive-interface-nils</summary>

Fix directives returning nils from optional interfaces

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/34d109754e83d85262d53a5b7df098e07908007c"><tt>34d10975</tt></a> Fix directives returning nils from optional interfaces

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/eea38e55661d6de749d1851e6d4447331a063df7"><tt>eea38e55</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/862">#862</a> from qhenkart/fixes-shareable-link-setting</summary>

fixes shareable link button in playground

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b5e78342937549496dbc3a362ed7a2a8738279c6"><tt>b5e78342</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/870">#870</a> from 99designs/ws-init-ctx</summary>

Allow changing context in websocket init func

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/034aa627cfa5b943497b54595e87b53637a6d1f5"><tt>034aa627</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/871">#871</a> from 99designs/subscription-middleware</summary>

Call middleware and directives for subscriptions

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7b41ca3c13858ca6e69659ade6bc0fc7a54d81da"><tt>7b41ca3c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/872">#872</a> from 99designs/autobind-prefix</summary>

Allow prefixes when using autobind

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/de8e559f5a90fb1120bec15465f10cc5adea74cc"><tt>de8e559f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/854">#854</a> from wabain/nested-map-interface</summary>

Fix for nested fields backed by map or interface

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cc64f331d1024d485ee583cdd0be61e8cf03a506"><tt>cc64f331</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/828">#828</a> from 99designs/feat-rc</summary>

introduce RequestContext#Validate and use it instead of NewRequestContext function

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ed2a853637e3bb944fe7264c22496b5c6e4c5a2e"><tt>ed2a8536</tt></a> Allow prefixes when using autobind

- <a href="https://github.com/99designs/gqlgen/commit/819cc71b92808353fe73c88904aa7057db6abcb3"><tt>819cc71b</tt></a> Call middleware and directives for subscriptions

- <a href="https://github.com/99designs/gqlgen/commit/5a7c5903f64efb240d575ef947b0ed1d59b1a3d0"><tt>5a7c5903</tt></a> Allow changing context in websocket init func

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/17f32d28c29ec45dce407f2f6afac16bdd8d64ca"><tt>17f32d28</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/861">#861</a> from 99designs/refactor-test-client</summary>

Refactor test client

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ed14cf045779a9a485b2f31a523b120c60f463a3"><tt>ed14cf04</tt></a> Update playground.go</summary>

fix formatting

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ee8d7a173d3fb8066dff40e6e9c0b4d1e9260b71"><tt>ee8d7a17</tt></a> Update playground.go</summary>

fix formatting

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/27389951d110511732010349e2c519cfa71319c5"><tt>27389951</tt></a> fixes shareable link button in playground

- <a href="https://github.com/99designs/gqlgen/commit/4162d11e2badb02e7a7e1974d02838a1127fcb47"><tt>4162d11e</tt></a> Refactor test client

- <a href="https://github.com/99designs/gqlgen/commit/8ed6ffc7a183732696564f7cdd800c2bb29e4ea6"><tt>8ed6ffc7</tt></a> Fix for nested fields backed by map or interface

- <a href="https://github.com/99designs/gqlgen/commit/55b2144289debeb5ca104a4d01d36f96be5ed84c"><tt>55b21442</tt></a> Update stale.yml

- <a href="https://github.com/99designs/gqlgen/commit/feebee7d305e02b6ba96eb9307922436a61f99a4"><tt>feebee7d</tt></a> stalebot

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7e643fdc5cc5097f786ef6a386e28feee293fd7a"><tt>7e643fdc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/838">#838</a> from 99designs/fix-directive-nil</summary>

fix directives return nil handling

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f33e09e8109cfcc1ef2dd4096d27d0a4b6eee9c8"><tt>f33e09e8</tt></a> Merge branch 'master' into fix-directive-nil

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8590edef5b8bbf094240f33c0ed696e034ca80e0"><tt>8590edef</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/839">#839</a> from 99designs/fix-nil-directive</summary>

refactor unimplemented directive handling

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1f7ed0d52ac853830c164d5d3ebbed8a77c8d710"><tt>1f7ed0d5</tt></a> refactor unimplemented directive handling

- <a href="https://github.com/99designs/gqlgen/commit/94ad3f2e92e0fd7ffc99566f3c4e1abe7a616acc"><tt>94ad3f2e</tt></a> fix directives return nil handling

- <a href="https://github.com/99designs/gqlgen/commit/5c644a6fbef1a9bc1c50ef6975686711ec31ff28"><tt>5c644a6f</tt></a> v0.9.3 postrelease bump

- <a href="https://github.com/99designs/gqlgen/commit/82758be87d570691febde1f1072a435af4c3920c"><tt>82758be8</tt></a> fix error

- <a href="https://github.com/99designs/gqlgen/commit/edde2d03aa14cb1b6c42748bd40e9c87f6670d12"><tt>edde2d03</tt></a> add OperationName field to RequestContext

- <a href="https://github.com/99designs/gqlgen/commit/830e466ec066f58c93d1f38ada51c6c874b74e19"><tt>830e466e</tt></a> introduce RequestContext#Validate and use it instead of NewRequestContext function

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.9.3"></a>
## [v0.9.3](https://github.com/99designs/gqlgen/compare/v0.9.2...v0.9.3) - 2019-08-16
- <a href="https://github.com/99designs/gqlgen/commit/a7bc468ca1b184a5ce1b07ea331e0121fc56ae82"><tt>a7bc468c</tt></a> release v0.9.3

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fc02cfe83a8f78f36b5e37a63b0a87bf511e94b2"><tt>fc02cfe8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/829">#829</a> from 99designs/fix-2directives</summary>

fix go syntax issue when field has 2 directives

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/924f620c4110ac02b76d302d686f3c58e77948ed"><tt>924f620c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/831">#831</a> from yudppp/patch-1</summary>

Fixed scalar reference documentation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ca4cc732569d70c84aea94dd1f22c5748e23881e"><tt>ca4cc732</tt></a> Fixed scalar documents

- <a href="https://github.com/99designs/gqlgen/commit/cc9fe1450c86effb85642bca00735ef17f0415f8"><tt>cc9fe145</tt></a> fix go syntax issue when field has 2 directives

- <a href="https://github.com/99designs/gqlgen/commit/6b70be0316bce27b048f64a44f82f1e450716d2c"><tt>6b70be03</tt></a> v0.9.2 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.9.2"></a>
## [v0.9.2](https://github.com/99designs/gqlgen/compare/v0.9.1...v0.9.2) - 2019-08-08
- <a href="https://github.com/99designs/gqlgen/commit/4eeacc6e4cb7bedc7c5312b6a3947697ad5cfb55"><tt>4eeacc6e</tt></a> release v0.9.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5628169dd38517e856e3d12b50696b4e9a79d60f"><tt>5628169d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/822">#822</a> from 99designs/windows-import-path-loop</summary>

fix for windows infinite loop

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a861aa524c4d96da8a4b25f75fde5006a29398c0"><tt>a861aa52</tt></a> lint fix

- <a href="https://github.com/99designs/gqlgen/commit/6348a5632123286aae227c456f5178275ddb737a"><tt>6348a563</tt></a> fix for windows infinite loop

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/12893fa4db4bcd5b341aed584dd584d0c6f2b226"><tt>12893fa4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/821">#821</a> from 99designs/fix-init</summary>

Fix config loading during gqlgen init

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5fafe79c050ae34448ee495c703caced6b8f3126"><tt>5fafe79c</tt></a> Fix config loading during gqlgen init

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2599f5607863415a015feb97f1abb9d46359e1ee"><tt>2599f560</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/820">#820</a> from 99designs/keepalive-on-init</summary>

send keepalive on init

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/139e4e8d444b5c6fb27406ad610ca4fce897710a"><tt>139e4e8d</tt></a> More directive docs

- <a href="https://github.com/99designs/gqlgen/commit/f93df34059a0867c7c33bc992aff0e3a5ddb0f14"><tt>f93df340</tt></a> send keepalive on init

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8f0d9b482a9b735810a99c80543c37323592e4e9"><tt>8f0d9b48</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/816">#816</a> from nii236/patch-1</summary>

Update cors.md to allow CORS for websockets

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/297e09c4a99356b8c7ef7ec7922ba86a94d4c435"><tt>297e09c4</tt></a> change origin check

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/410d83225aec8f4cb20e50b2a187a16d6dceded6"><tt>410d8322</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/805">#805</a> from andrey1s/golangci</summary>

enable-all linters on golangci-lint

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/504a96bc30e1a672191c5baedc60916cd44e57d3"><tt>504a96bc</tt></a> set enabled linters

- <a href="https://github.com/99designs/gqlgen/commit/91966ef485331c61177b384bfb029e51eb0a3cb1"><tt>91966ef4</tt></a> add example to lint

- <a href="https://github.com/99designs/gqlgen/commit/bcddd7aad8eb4147f13a845d20d391790f6512c9"><tt>bcddd7aa</tt></a> fix typo in readme

- <a href="https://github.com/99designs/gqlgen/commit/cce06f1d060bbb5093e6c5b992ec907499725403"><tt>cce06f1d</tt></a> update lint in circleci

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/da1c208e00658e40800af7029fd4b522c5f9655c"><tt>da1c208e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/795">#795</a> from oshalygin/feature/issue-794-resolve-dead-readme-link</summary>

Update GraphQL Reference Link

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8343c32c2b0059bd8d04ac56fa3c82bbbb6b908e"><tt>8343c32c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/784">#784</a> from y15e/add-missing-header</summary>

Add a missing "Upload" header

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8302463fb74745c198436381049f7adbd750421e"><tt>8302463f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/797">#797</a> from muesli/format-fixes</summary>

Format import order using goimports

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f2825e09c2331e21d0e3bdbeba4ffbd23cd7a1b0"><tt>f2825e09</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/801">#801</a> from Schparky/patch-1</summary>

Documentation: getting-started edits

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3db5627f2f3987b7cec7e8732b24cc2fdf27fc24"><tt>3db5627f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/807">#807</a> from flrossetto/patch-1</summary>

Fix doc

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ab228f1bd4477dcb2ceea091acef17d35322cbb2"><tt>ab228f1b</tt></a> Update cors.md to allow CORS for websockets

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c4ac93473b017e423823716f2548a633f2527517"><tt>c4ac9347</tt></a> Fix doc</summary>

map[string]{interface} -> map[string]interface{}

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fbbed5b82a01e037c4c08871c35c129f7f008c36"><tt>fbbed5b8</tt></a> use alias when invalid pkg name

- <a href="https://github.com/99designs/gqlgen/commit/2591ea36be96a057f27c571668f303ce95647f8b"><tt>2591ea36</tt></a> fix lint prealloc

- <a href="https://github.com/99designs/gqlgen/commit/3b0e44fecf327e380162c367073ce574849d2201"><tt>3b0e44fe</tt></a> fix lint misspell

- <a href="https://github.com/99designs/gqlgen/commit/6ff62b61668930c1565a59d8900d0b96353d0210"><tt>6ff62b61</tt></a> fix lint gocritic

- <a href="https://github.com/99designs/gqlgen/commit/cb7f482b9167917c71eddb6ea661faf72b5a90ec"><tt>cb7f482b</tt></a> fix lint unparam

- <a href="https://github.com/99designs/gqlgen/commit/620552be324097b65576240b0f886a687bb57f88"><tt>620552be</tt></a> fix lint goimports

- <a href="https://github.com/99designs/gqlgen/commit/477e804eb60c07391b609aaf35e83670484d2059"><tt>477e804e</tt></a> update config golangci

- <a href="https://github.com/99designs/gqlgen/commit/5b203bcca8e36841e7249aa1adb409cb261a695d"><tt>5b203bcc</tt></a> clarify where the go:generate line should be added

- <a href="https://github.com/99designs/gqlgen/commit/2a3df24e66417c6071f281fa4cdba328709c7dca"><tt>2a3df24e</tt></a> Replace the -v flag as described below.

- <a href="https://github.com/99designs/gqlgen/commit/f3eeb6392dd06b43f6819c3ee946dbd949727e7b"><tt>f3eeb639</tt></a> Clarify that the schema file will be generated

- <a href="https://github.com/99designs/gqlgen/commit/3ac17960bb4b8a32c47368567c8b8beef9b90b4d"><tt>3ac17960</tt></a> Missing '*' in Todos resolver example

- <a href="https://github.com/99designs/gqlgen/commit/bd598c2ce3daa89958573f56f40bf01c948f1cc9"><tt>bd598c2c</tt></a> Format import order using goimports

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/419f966d3c30abe85b53594f2fb5c51f2da07cf1"><tt>419f966d</tt></a> Update GraphQL Reference Link (closes <a href="https://github.com/99designs/gqlgen/issues/794"> #794</a>)</summary>

- The link in the readme has been updated to reference a post by
  Iván Corrales Solera, "Dive into GraphQL".  The previous link
  does not resolve, likely because the personal site is no longer
  hosted.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/373359de83dd626d4a19ec57b7b599c58c88ca2c"><tt>373359de</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/781">#781</a> from 99designs/fix-default-directives-init</summary>

Set default directives after parsing config

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ca8b21e31f05a633ac79f79328fb3993cbfe7b8d"><tt>ca8b21e3</tt></a> Add a missing header

- <a href="https://github.com/99designs/gqlgen/commit/8cab5fba1f1c7ebda1a3cffc789a7ab7a2ac2736"><tt>8cab5fba</tt></a> Set default directives after parsing config

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d2c5bf2ae2d8d7da647aded0b3287be4ad2547a9"><tt>d2c5bf2a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/780">#780</a> from zdebra/master</summary>

fixed generating a description to golang comments for enum type

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bf2cc90ec9bda2be035d6fc93f58344577f12172"><tt>bf2cc90e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/768">#768</a> from 99designs/fix-ptr-from-directive</summary>

Fix pointer returns from directive

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/446c3df37f727fd638442669f2e86b12fa688c3a"><tt>446c3df3</tt></a> fixed generating a description to golang comments for enum type

- <a href="https://github.com/99designs/gqlgen/commit/414a4d3414b2b5856851816f0152980098f7b3ab"><tt>414a4d34</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/771">#771</a> from sunfmin/master

- <a href="https://github.com/99designs/gqlgen/commit/4d1484b012d9f0d35d82f1ffa24157e66faac444"><tt>4d1484b0</tt></a> Fix doc for how to use [@goField](https://github.com/goField) directives forceResolver option

- <a href="https://github.com/99designs/gqlgen/commit/6f3d73103dae0fa8f409b8b869421ba1c57f3e90"><tt>6f3d7310</tt></a> Fix pointer returns from directive

- <a href="https://github.com/99designs/gqlgen/commit/21b65112e5952aee4a5cf40da97551ccbf246552"><tt>21b65112</tt></a> v0.9.1 postrelease bump

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.9.1"></a>
## [v0.9.1](https://github.com/99designs/gqlgen/compare/v0.9.0...v0.9.1) - 2019-06-27
- <a href="https://github.com/99designs/gqlgen/commit/b128a29122e8ca8ada5f34cc18338fa7c10fc5b4"><tt>b128a291</tt></a> release v0.9.1

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1bbc0cd60235877c2aa14242f5d4ec8c4bab5083"><tt>1bbc0cd6</tt></a> Update release process to keep tags on master</summary>

this was affecting the version shown in go modules when using commits

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5ffc29754dd71a847f6860f0fb37d75dea367ee7"><tt>5ffc2975</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/764">#764</a> from 99designs/fix-field-directives-on-roots</summary>

fix field schema directives applied to roots

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ef3830b5e951d30dd49c43183dc029c79c338645"><tt>ef3830b5</tt></a> fix field schema directives applied to roots

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/17ee40ba40898db7b8293e537ee6bb0aa953c0b3"><tt>17ee40ba</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/761">#761</a> from 99designs/autobinding</summary>

Autobind models

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b716bfac517ae641461df88c525643b5fcdf184e"><tt>b716bfac</tt></a> Autobind models

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fc3755f1c2f0eb6383c59646388161052aa5e676"><tt>fc3755f1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/732">#732</a> from 99designs/schemaconfig-plugin</summary>

Add a plugin for configuring gqlgen via directives

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c14f8650d72591d6ae7a75e904118c05cf1e291f"><tt>c14f8650</tt></a> Add docs

- <a href="https://github.com/99designs/gqlgen/commit/64aca616f334797818f7272c0c11eccc86d2d93b"><tt>64aca616</tt></a> Merge remote-tracking branch 'origin/master' into schemaconfig-plugin

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5e7e94c80034a986f45305eb6c5ed559259fbd16"><tt>5e7e94c8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/756">#756</a> from andrey1s/field</summary>

generate field defenition and execute field directive

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ad2ca304b532470721a0a6d2a9a78b85eef633cf"><tt>ad2ca304</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/759">#759</a> from 99designs/circle-workflows</summary>

CircleCI workflows

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0fc822ca68f02fe7c510519ca91c7e0a131fbb99"><tt>0fc822ca</tt></a> CircleCI workflows

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2dc8423b7937c9d012ee8886d2854011cf61dee7"><tt>2dc8423b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/758">#758</a> from franxois/patch-1</summary>

Update dataloaders.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d0db28ab9b1327bb539b57cd68a82835b48afc37"><tt>d0db28ab</tt></a> Update dataloaders.md</summary>

Make SQL request use requested IDs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a58ecfe9be6b8aa9201a3221e2838ea7cf5b2f9f"><tt>a58ecfe9</tt></a> add example and test field directive

- <a href="https://github.com/99designs/gqlgen/commit/526beecb981b91e3e880050af655b4caa43f4fb4"><tt>526beecb</tt></a> update generate field

- <a href="https://github.com/99designs/gqlgen/commit/6e9d7dab9458106d2deed0e83910202d6625245b"><tt>6e9d7dab</tt></a> generate types directive by location

- <a href="https://github.com/99designs/gqlgen/commit/dfec7b687fb5b61780941f51a53b195dadca8621"><tt>dfec7b68</tt></a> define fieldDefinition template

- <a href="https://github.com/99designs/gqlgen/commit/be890ab9a1d887a1993aa3403c4e36c294003187"><tt>be890ab9</tt></a> use UnmarshalFunc in args directives implement

- <a href="https://github.com/99designs/gqlgen/commit/dd162f04051c034bfae7c8fb3732975dec449586"><tt>dd162f04</tt></a> define implDirectives template

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/56f3f92b8ee315ef8a3c32484e6b63dd13ae574a"><tt>56f3f92b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/755">#755</a> from 99designs/fix-globbing-windows</summary>

fix globbing on windows

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a4480fb078794d8761a77d5053abf7fd0cb759fc"><tt>a4480fb0</tt></a> fix globbing on windows

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ba176e2efbc717c6da5ca9b604524104c3daddec"><tt>ba176e2e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/754">#754</a> from 99designs/coveralls</summary>

Add coveralls

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f28ed264310d3c62c0a19a5de742d006019a7675"><tt>f28ed264</tt></a> Add coveralls

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f4a69ab5fa6235e840c1059fbf0a670c1ab69177"><tt>f4a69ab5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/744">#744</a> from andrey1s/directive</summary>

add Execute QUERY/MUTATION/SUBSCRIPTION Directives

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/dbd2cc6e5d5eedfad3c857c4d8cb1052a25ff19f"><tt>dbd2cc6e</tt></a> simplify resolver test

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7fed71b6ce42790a2894fb550ef33f186acacd65"><tt>7fed71b6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/728">#728</a> from fgallina/make-generated-resolver-dependent-types-follow-configured-type</summary>

resolvergen: use the resolver type as base name for dependent types

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cb284c568490926d6d25999daab9ee3ad3bc6a06"><tt>cb284c56</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/734">#734</a> from DBL-Lee/master</summary>

Automatic Persisted Queries

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/726a94f4895cb9d7ca031e97893e60eac3ac0e5d"><tt>726a94f4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/750">#750</a> from 99designs/ws-connection-param-check</summary>

[websocket] Add a config to reject initial connection

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/69d7e28241b9847a073f1a335a5cc12e2efddf05"><tt>69d7e282</tt></a> move directive to directives.gotpl

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/090f0bd949536c41d2d0610c0454263bd80f8243"><tt>090f0bd9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/722">#722</a> from marwan-at-work/deps</summary>

resolve all pkg dependencies

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c397be0c409c7566214fa673a504a4f5821d4bf5"><tt>c397be0c</tt></a> Update websocketInitFunc to return error instead of boolean

- <a href="https://github.com/99designs/gqlgen/commit/be18ae1feaedfdba4787a745b8b62d981955912c"><tt>be18ae1f</tt></a> Add a test

- <a href="https://github.com/99designs/gqlgen/commit/a6508b6d4fbec8fcba56091293fc4e5ade3a13aa"><tt>a6508b6d</tt></a> Update typing, function name and small code refactor

- <a href="https://github.com/99designs/gqlgen/commit/e6d791a9b83827bb5054b023baf100c5866d54bd"><tt>e6d791a9</tt></a> Add websocketOnConnectFunc as a config that can be used to validate websocket init requests

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c5acbead96b3c2c307809b8dd37af9dda610a84b"><tt>c5acbead</tt></a> resolvergen: use the resolver type as base name for dependent types</summary>

The template was outputing invalid code since the resolver type was
not used in places like the embedding at {query,mutation}Resolver.

This change also ensures that objects like {query,mutation}Resolver
also use the user provided type name as suffix.

Here's the resulting diff on the code generation with `type:
GeneratedResolver` in the resolver config:

```
diff -u resolver.go resolvernew.go
--- resolver.go 2019-05-26 20:04:15.361969755 -0300
+++ resolvernew.go      2019-05-26 20:04:54.170737786 -0300
@@ -7,20 +7,20 @@
 type GeneratedResolver struct{}

 func (r *GeneratedResolver) Mutation() MutationResolver {
-       return &mutationResolver{r}
+       return &mutationGeneratedResolver{r}
 }
 func (r *GeneratedResolver) Query() QueryResolver {
-       return &queryResolver{r}
+       return &queryGeneratedResolver{r}
 }

-type mutationResolver struct{ *Resolver }
+type mutationGeneratedResolver struct{ *GeneratedResolver }

-func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*Todo, error) {
+func (r *mutationGeneratedResolver) CreateTodo(ctx context.Context, input NewTodo) (*Todo, error) {
        panic("not implemented")
 }

-type queryResolver struct{ *Resolver }
+type queryGeneratedResolver struct{ *GeneratedResolver }

-func (r *queryResolver) Todos(ctx context.Context) ([]*Todo, error) {
+func (r *queryGeneratedResolver) Todos(ctx context.Context) ([]*Todo, error) {
        panic("not implemented")
 }
```

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cfdbc39ac69a08450ebf70dfaa40f4b86225d0a2"><tt>cfdbc39a</tt></a> update QueryDirectives

- <a href="https://github.com/99designs/gqlgen/commit/f32571ee002479dc07c3b81548d15fdc1169bfa6"><tt>f32571ee</tt></a> add SUBSCRIPTION Directive

- <a href="https://github.com/99designs/gqlgen/commit/32462d0f1c2cd6b9f94639ee83da288d25219eae"><tt>32462d0f</tt></a> update example todo add directive with location QUERY and MUTATION

- <a href="https://github.com/99designs/gqlgen/commit/3eec887a69508b9e431d384a9adae8ee53d63b97"><tt>3eec887a</tt></a> add Execute QUERY/MUTATION/SUBSCRIPTION Directives

- <a href="https://github.com/99designs/gqlgen/commit/8fcc186817974f99060f0f815dd3935876607bf0"><tt>8fcc1868</tt></a> format

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e0e1e318bb5bbae97348f3d773fd16d8f7fc8317"><tt>e0e1e318</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/1">#1</a> from radev/master</summary>

Support for external APQ cache

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9873d998b54009721c81b3c13b23f975463aaf02"><tt>9873d998</tt></a> Add APQ documentation with example

- <a href="https://github.com/99designs/gqlgen/commit/48292c1020344b7686efa6667ea16ed7706ece14"><tt>48292c10</tt></a> Support pluggable APQ cache implementations.

- <a href="https://github.com/99designs/gqlgen/commit/694f90aa089c34a80bf3007a6d298a96ba7f132c"><tt>694f90aa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/717">#717</a> from cbelsole/schema_file_globbing (closes <a href="https://github.com/99designs/gqlgen/issues/631"> #631</a>)

- <a href="https://github.com/99designs/gqlgen/commit/9be5aad0cf295796e01ea0955ff53f944a8c5cb9"><tt>9be5aad0</tt></a> Don't inject builtins during schema config

- <a href="https://github.com/99designs/gqlgen/commit/8dc17b470d339f5fcd6e1cf9a3d14dec2df4067e"><tt>8dc17b47</tt></a> support GET for apq

- <a href="https://github.com/99designs/gqlgen/commit/d36932c55ee97e41567eb9c42d49a884388095d0"><tt>d36932c5</tt></a> support automatic persisted query

- <a href="https://github.com/99designs/gqlgen/commit/de75743c1169cc0e51de0049e78c7f7c0bd92cef"><tt>de75743c</tt></a> Add plugin for providing config via schema directives

- <a href="https://github.com/99designs/gqlgen/commit/17a82c37e86df494354c92de4b306a15c11747ee"><tt>17a82c37</tt></a> Provide config to skip generating runtime for a directive

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ba7092c595a7e62f3a2350504f35bd2cc11b0c1b"><tt>ba7092c5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/724">#724</a> from saint1991/patch-1</summary>

added a missing close bracket

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9c1f8f2a10f3664749b51dc8805448881d73c625"><tt>9c1f8f2a</tt></a> added a missing close bracket

- <a href="https://github.com/99designs/gqlgen/commit/3dd8baf528b79dc3945f6a156c7cb7316ef87479"><tt>3dd8baf5</tt></a> resolve all pkg dependencies

- <a href="https://github.com/99designs/gqlgen/commit/1617ff28daba04a67413ba9696c7650e718aa080"><tt>1617ff28</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/718">#718</a> from hh/fix-docs (closes <a href="https://github.com/99designs/gqlgen/issues/714"> #714</a>)

- <a href="https://github.com/99designs/gqlgen/commit/9d332a7d77d025b2f4e9bec74a8c7c06c6129098"><tt>9d332a7d</tt></a> Fixing getting-started documentation

- <a href="https://github.com/99designs/gqlgen/commit/39db147719f48556710f1f5afedf6e79bc8affc4"><tt>39db1477</tt></a> updated docs

- <a href="https://github.com/99designs/gqlgen/commit/e32c82be0f3d562b149c01b08f459b6515b75aca"><tt>e32c82be</tt></a> cleanup

- <a href="https://github.com/99designs/gqlgen/commit/e9389ef8f8eee80eec108983575ed303e99000e9"><tt>e9389ef8</tt></a> added schema file globbing fixes <a href="https://github.com/99designs/gqlgen/pull/631">#631</a>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4f163cbc8466634ddaeca0a4071813a76ec55ea5"><tt>4f163cbc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/713">#713</a> from 99designs/faq</summary>

Add faq section

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3a21b36965067663a85745c1d0df06894be7af67"><tt>3a21b369</tt></a> Add faq section

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.9.0"></a>
## [v0.9.0](https://github.com/99designs/gqlgen/compare/v0.8.3...v0.9.0) - 2019-05-15
- <a href="https://github.com/99designs/gqlgen/commit/ea4652d223c441dc77b31882781ce08488763d67"><tt>ea4652d2</tt></a> release v0.9.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f3c8406d909bd3f7bce89897637648870f7b1295"><tt>f3c8406d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/710">#710</a> from 99designs/slice-pointers</summary>

Use pointers to structs inside slices

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e669d476e6a04bf2fb43ed8e53bb28e91b424a3a"><tt>e669d476</tt></a> fix imports for vendor based projects

- <a href="https://github.com/99designs/gqlgen/commit/315141d9bd2ab14169a88a50f19326c4483d17e3"><tt>315141d9</tt></a> Use pointers to structs inside slices

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9a6a10abe7d927f2d3271acc9bff9058cd070bf9"><tt>9a6a10ab</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/706">#706</a> from 99designs/mapping-primitive</summary>

Fix mapping object types onto go primitives

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a5120054d3bafebd6bbb9557e6b4e773c40d6693"><tt>a5120054</tt></a> fix binding to primitive non leaf types

- <a href="https://github.com/99designs/gqlgen/commit/b0cd95a19710fd92ab9e76422591ac9f0c5a6f31"><tt>b0cd95a1</tt></a> Test mapping object types onto go string

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/eaa61bb56c75a6fdb98ee7b53b81b70b519b35ba"><tt>eaa61bb5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/707">#707</a> from 99designs/gomodules-performance</summary>

make gqlgen generate 10x faster in some projects

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ab961ce000f72314346521009172e6391d107b69"><tt>ab961ce0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/705">#705</a> from 99designs/fix-error-race</summary>

Fix a data race when handling concurrent resolver errors

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/71cc8554135ea424df46baf1f41dd9b43de0e673"><tt>71cc8554</tt></a> make gqlgen generate 10x faster in projects with cgo

- <a href="https://github.com/99designs/gqlgen/commit/cab4babec38b23dc30f0ecfc53a92a1c2dce41fd"><tt>cab4babe</tt></a> Test mapping object types onto go primitives

- <a href="https://github.com/99designs/gqlgen/commit/962470dec347d728e236395c59844ae0acf22333"><tt>962470de</tt></a> Fix a data race when handling concurrent resolver errors

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9ca43ba938f6c67a33986cddf650658fc406fe95"><tt>9ca43ba9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/701">#701</a> from 99designs/modelgen-pointers</summary>

Use pointers when embedding structs in generated structs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4f5e9cf03a047ce69523df374be8dd3991735f65"><tt>4f5e9cf0</tt></a> always use pointers when refering to structs in generated models

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e2ac84807945c048572c49bb73b48aa70845ca04"><tt>e2ac8480</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/704">#704</a> from tul/doc-typo</summary>

Fix typo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/80ebe644b0286a466beea67f6a76c910b083b8fd"><tt>80ebe644</tt></a> Fix typo

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0bd9080964817324cbb046a3833bdfc4674d22fb"><tt>0bd90809</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/700">#700</a> from 99designs/fix-interface-caseing</summary>

Fix interface casing

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5586ee2c0e610ad2f257d2bda577d97437bbb17e"><tt>5586ee2c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/702">#702</a> from 99designs/drop-automatic-zeroisnull</summary>

Drop automatic conversion of IsZero to null

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/75aa99ad2bfc9506c76bb86ff810461b699ec31d"><tt>75aa99ad</tt></a> Drop automatic conversion of IsZero to null

- <a href="https://github.com/99designs/gqlgen/commit/46c40b748d55bcecd236c1ca827cef00808a5363"><tt>46c40b74</tt></a> Fix interface casing (closes <a href="https://github.com/99designs/gqlgen/issues/694"> #694</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e49d44f7464a05347691bfa9ed9eab44bf4ffe5f"><tt>e49d44f7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/689">#689</a> from tgwizard/enforce-request-content-type</summary>

Enforce content type for POST requests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/78f277e9178316fe8ae3e646a5d5c1ba4c9e5d15"><tt>78f277e9</tt></a> run go generate

- <a href="https://github.com/99designs/gqlgen/commit/d4b3de3aff5eeb6ad9a7093e6b3905edeef10f1e"><tt>d4b3de3a</tt></a> Merge remote-tracking branch 'origin/master' into enforce-request-content-type

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f8ef6d2e2be1b481f4fdbb065da5b50236600143"><tt>f8ef6d2e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/668">#668</a> from mbranch/complexity</summary>

Fix: complexity case selection

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c4805049b75965b1305583fda6e60487c727d8a0"><tt>c4805049</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/655">#655</a> from hantonelli/file-upload</summary>

File upload

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5d1dea0a104ff07f657aa18e63097bf1626174a7"><tt>5d1dea0a</tt></a> run go generate

- <a href="https://github.com/99designs/gqlgen/commit/8a0c34a485d06a2eea5e8fb7cbc0aed0a025d80d"><tt>8a0c34a4</tt></a> Merge branch 'master' into file-upload

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4e359aa26d29f8a182010e907a834ee948abf45d"><tt>4e359aa2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/686">#686</a> from qhenkart/master</summary>

Adds default custom scalar of interface{}

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/aeccbce0fe78665cd0d7d55c42b4df142623fe6b"><tt>aeccbce0</tt></a> Update test include an example that uses io.Read interface directly

- <a href="https://github.com/99designs/gqlgen/commit/d9dca642e348a9e8ea8749cd48235bf4da633b63"><tt>d9dca642</tt></a> Improve documentation

- <a href="https://github.com/99designs/gqlgen/commit/f30f1c312f5e49142578705cef05ef753bcc4277"><tt>f30f1c31</tt></a> Fix fmt

- <a href="https://github.com/99designs/gqlgen/commit/54226cdbb2ab5e2f687241357ec9c062667e7b8a"><tt>54226cdb</tt></a> Add bytesReader to reuse read byte array

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/02e9dd8e842028aa3e38da361f04b1d5d406f024"><tt>02e9dd8e</tt></a> Fix complexity case selection</summary>

Use the GraphQL field name rather than the Go field name in the generated
`Complexity` func.

Before this patch, overloading complexity funcs was ineffective because they
were never executed.

It also ensures that overlapping fields are now generated; mapping all possible
field names to the associated complexity func.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bf2d07a47fa57b27434a02021ca8ed868019ba9c"><tt>bf2d07a4</tt></a> moves naming convention to a non-go standard

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d1e8acda1dd71b8f944ab08019e70ae6092ab61a"><tt>d1e8acda</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/687">#687</a> from stereosteve/fix-includeDeprecated</summary>

Fix: omit deprecated fields when includeDeprecated=false

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f7d0b9c824c025856d066a76c320a7752f45b519"><tt>f7d0b9c8</tt></a> Enforce content type for POST requests

- <a href="https://github.com/99designs/gqlgen/commit/7d0b8eecc94c16d46a653758de9c8506b953b8c0"><tt>7d0b8eec</tt></a> Fix: omit deprecated fields when includeDeprecated=false

- <a href="https://github.com/99designs/gqlgen/commit/89c873459212f48d95e16d87d5732460d229255e"><tt>89c87345</tt></a> fix grammar in docs

- <a href="https://github.com/99designs/gqlgen/commit/85643f5dad6ff410ebb78f26539e845f38df37e1"><tt>85643f5d</tt></a> fix import

- <a href="https://github.com/99designs/gqlgen/commit/ca96a1550d4581d9aef9532b8eacdbcf6e0122a0"><tt>ca96a155</tt></a> update docs

- <a href="https://github.com/99designs/gqlgen/commit/1de25d0c65d549e10684d55e674a1ec1748aea86"><tt>1de25d0c</tt></a> adds interface scalar type

- <a href="https://github.com/99designs/gqlgen/commit/43fc53f9a90429c6693bfe88afaa83b1de8591e3"><tt>43fc53f9</tt></a> Improve variable name

- <a href="https://github.com/99designs/gqlgen/commit/b961d34e06ee088251ffbe2929c86db4381dec24"><tt>b961d34e</tt></a> Remove wrapper that is now not required

- <a href="https://github.com/99designs/gqlgen/commit/bb0234760f8d30ce1d142fe117d5b857eff7cb8e"><tt>bb023476</tt></a> Lint code

- <a href="https://github.com/99designs/gqlgen/commit/f8484159adde6204e727c3edc09bb7ca4e264029"><tt>f8484159</tt></a> Modify graphql.Upload to use io.ReadCloser. Change the way upload files are managed.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0306783ed6c9eeb371f5775a33cfc505a9961351"><tt>0306783e</tt></a> Revert "Change graphql.Upload File field to FileData."</summary>

This reverts commit 7ade7c2

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/afe33f73875beca92e917742c1f49c1f6145018b"><tt>afe33f73</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/680">#680</a> from asp24/collect-fields-performance</summary>

Better CollectFields performance

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7ba1b3b298695a7382abda437d3417d0aeae9c5b"><tt>7ba1b3b2</tt></a> graphql.CollectFields now accept *RequestContext as first arg It was done because RequestContext is a part of executionContext and can be passed directly without extraction from ctx. This is increasing performance when model depth is high

- <a href="https://github.com/99designs/gqlgen/commit/5dfa2285234fc7fe2839bd567a09351e6f6f2136"><tt>5dfa2285</tt></a> Pre-allocate mem for collectFields() method result slice

- <a href="https://github.com/99designs/gqlgen/commit/88cdbdf10d5e2da58b06249202fbb2769e21d0bb"><tt>88cdbdf1</tt></a> Rename getOrCreateField to getOrCreateAndAppendField to describe behaviour

- <a href="https://github.com/99designs/gqlgen/commit/a74abc4738315f36fee98ad53ac02a0b2093e36a"><tt>a74abc47</tt></a> Early return in shouldIncludeNode if directives empty

- <a href="https://github.com/99designs/gqlgen/commit/7ade7c21d83d37397bd0fa42dc7bfadcdd492fa9"><tt>7ade7c21</tt></a> Change graphql.Upload File field to FileData.

- <a href="https://github.com/99designs/gqlgen/commit/da52e810cf98ffe9a57cb0a8ecf68df97f6ac6f7"><tt>da52e810</tt></a> Extend test and don't close form file.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1c95d42a4cd8e20fd41dba4170e02e03ad7dacbe"><tt>1c95d42a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/678">#678</a> from jonatasbaldin/gin-context-recipe</summary>

Fix unset key and comment block at Gin recipe docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0b39c44526fa300bb9341edd1e11fa71e7cdcd74"><tt>0b39c445</tt></a> Fix unset key and comment block

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5aa6a20b72ec0d90dc676abcd796ad33b2e39a83"><tt>5aa6a20b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/673">#673</a> from marwan-at-work/tpl</summary>

codegen/templates: allow templates to be passed in options instead of…

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/37fd067e8099e2c11706db7c2d00f27bc8c76dd0"><tt>37fd067e</tt></a> fix typo

- <a href="https://github.com/99designs/gqlgen/commit/e69b739955acd2e2bd744ec84abc61661ca6a626"><tt>e69b7399</tt></a> add docs to the templates package

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8cae895b99cad57406851d11b2fd92c628f18ff6"><tt>8cae895b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/676">#676</a> from jonatasbaldin/gin-context-recipe</summary>

Add recipe to use gin.Context

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/40c7b9524d94b342a6169653483dd8f87061f1ec"><tt>40c7b952</tt></a> update test name

- <a href="https://github.com/99designs/gqlgen/commit/5418a290f381f85c89562cd7366d3f6ee98f3777"><tt>5418a290</tt></a> Add recipe to use gin.Context

- <a href="https://github.com/99designs/gqlgen/commit/16f392eeaffb9c3e9e973288c5cbe23c118ac0ad"><tt>16f392ee</tt></a> add unit test

- <a href="https://github.com/99designs/gqlgen/commit/a0ee7172289ff0bd0c00d45f50b80766a7f93f48"><tt>a0ee7172</tt></a> codegen/templates: allow templates to be passed in options instead of os files

- <a href="https://github.com/99designs/gqlgen/commit/2cf7f452dd9ca5b4baa1399c501926a51debf59d"><tt>2cf7f452</tt></a> Fix comments (add request size limit, remove useless comments, improve decoding and function signature, improve documentation)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5ff60925ce417b4689cd745361bb800608192dfe"><tt>5ff60925</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/665">#665</a> from ezeql/patch-1</summary>

update README.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b42e1ba633f1da7845159bd2d6b0903ada51e843"><tt>b42e1ba6</tt></a> update README.md</summary>

fix link

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d3770395a23e6ec92e06886036f1832f4ddf9af2"><tt>d3770395</tt></a> Fix tests.

- <a href="https://github.com/99designs/gqlgen/commit/2c1f8573321521e6ca70c73f6a78192a9ca38066"><tt>2c1f8573</tt></a> Fix lint errors.

- <a href="https://github.com/99designs/gqlgen/commit/73b3a5366cd145534014215d9c7597619d64b506"><tt>73b3a536</tt></a> Fmt graphql.go

- <a href="https://github.com/99designs/gqlgen/commit/83cde4b69e95758e8a819cc5ab02ab9f1a3301b3"><tt>83cde4b6</tt></a> Fix tests. Improve code format.

- <a href="https://github.com/99designs/gqlgen/commit/425849a697d475c70ba496f079d9cc5f46f457b2"><tt>425849a6</tt></a> Improve fileupload example readme. Update scalars.md. Add file-upload.md

- <a href="https://github.com/99designs/gqlgen/commit/849d4b1eaa5bf2141ef30e9a0d850f851770dcbe"><tt>849d4b1e</tt></a> Make uploadMaxMemory configurable

- <a href="https://github.com/99designs/gqlgen/commit/fc318364521228b1cf208fb6157e77a4d2af0ff0"><tt>fc318364</tt></a> Improve format, inline const.

- <a href="https://github.com/99designs/gqlgen/commit/662dc3372caa4aca09474eb67ff0c35b838af3f3"><tt>662dc337</tt></a> Move Upload to injected if defined in the schema as scalars

- <a href="https://github.com/99designs/gqlgen/commit/f244442e86a96c99333610972134cf90a081fddc"><tt>f244442e</tt></a> Fix merge. Remove regexp check.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bf79bc9289961295db8067f662a8b9a5619d5c30"><tt>bf79bc92</tt></a> Merge branch 'master' into next</summary>

# Conflicts:
#	codegen/config/config.go
#	handler/graphql.go
#	handler/graphql_test.go

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bd4aeaa6596443a044d7f0c0df4b1a8ed894f461"><tt>bd4aeaa6</tt></a> Merge remote-tracking branch 'upstream/master'

- <a href="https://github.com/99designs/gqlgen/commit/3a6f2fb7e7e83d7de00fc6198795257f81c40d59"><tt>3a6f2fb7</tt></a> Improve test code

- <a href="https://github.com/99designs/gqlgen/commit/239bc46f6f54f7d7f28f6f182b93a91dbc75d50d"><tt>239bc46f</tt></a> Add comments

- <a href="https://github.com/99designs/gqlgen/commit/be8d6d1249bca523133950a631686267467c0d96"><tt>be8d6d12</tt></a> Improve test

- <a href="https://github.com/99designs/gqlgen/commit/4d92696bf6e87404f85aaea0a5bfd8c15b7744b7"><tt>4d92696b</tt></a> Clean up code and add tests

- <a href="https://github.com/99designs/gqlgen/commit/2c414edcbb3f672afea9cb7d1cfa822129674b84"><tt>2c414edc</tt></a> Improve and add tests

- <a href="https://github.com/99designs/gqlgen/commit/68446e17d6157036cd59d897a35724dac43647bc"><tt>68446e17</tt></a> Revert change to websocket_test

- <a href="https://github.com/99designs/gqlgen/commit/61c1cb9cbb5b1cb154901526beec94041c8b35b2"><tt>61c1cb9c</tt></a> Improve examples

- <a href="https://github.com/99designs/gqlgen/commit/493d9375b767f85333e4f7d70d1bab9c0e47e475"><tt>493d9375</tt></a> Improve examples

- <a href="https://github.com/99designs/gqlgen/commit/3c5f8bb9f7e47ee3e39cc2cd254191d274e1c59f"><tt>3c5f8bb9</tt></a> Improve some examples

- <a href="https://github.com/99designs/gqlgen/commit/db7a03b117ba2fd74346f669df9b926e46b891c7"><tt>db7a03b1</tt></a> Improve tests and names

- <a href="https://github.com/99designs/gqlgen/commit/c493d1b96c33ab73ab27332cd9ba41a9f073170b"><tt>c493d1b9</tt></a> Revert changing to websocket_test

- <a href="https://github.com/99designs/gqlgen/commit/998f7674cae96d0cf7eca9eaffffc4cab3e55b0a"><tt>998f7674</tt></a> Revert changing the stub file

- <a href="https://github.com/99designs/gqlgen/commit/a7e95c597673673cb71a9f3e8f64a33f3dc1a197"><tt>a7e95c59</tt></a> Fix tests. Improve file generation

- <a href="https://github.com/99designs/gqlgen/commit/10beedb30fe05e94891cbffdc314498b4cc005c7"><tt>10beedb3</tt></a> Remove not required file

- <a href="https://github.com/99designs/gqlgen/commit/5afb6b4059b8d13206d19ad5ad4277698c118370"><tt>5afb6b40</tt></a> Add file upload to default schema

- <a href="https://github.com/99designs/gqlgen/commit/9c17ce33f74900fcb2db4285f30b1d57143759e6"><tt>9c17ce33</tt></a> Add file upload

- <a href="https://github.com/99designs/gqlgen/commit/b454621d436be44ef082bb82c841fe404a91639f"><tt>b454621d</tt></a> Add support to upload files.

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.8.3"></a>
## [v0.8.3](https://github.com/99designs/gqlgen/compare/v0.8.2...v0.8.3) - 2019-04-03
- <a href="https://github.com/99designs/gqlgen/commit/010a79b66f08732cb70d133dcab297a8ee895572"><tt>010a79b6</tt></a> release v0.8.3

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3623f7fcd78ac6e8ddaff3ecfb5b0e006c8e862a"><tt>3623f7fc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/650">#650</a> from andcan/plugin-funcmap</summary>

Allow plugins to provide additional template funcs

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a2e5936250f2807a1d1faebd1a2bf0d918af533d"><tt>a2e59362</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/652">#652</a> from andrey1s/extraBuiltins</summary>

add extra builtins types when no type exists

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c93d92ba10391cd0a1e44e493ed38ff85bb2acb0"><tt>c93d92ba</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/654">#654</a> from sharkyze/fix-introscpetion-doc</summary>

doc: fix mistake on introspection doc page

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/93e72b589b09e7582f072d9e911decb2d897c971"><tt>93e72b58</tt></a> doc: fix error on introspection doc page

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ef2e51ba3e4b6408c48f4e44e0b048947196719b"><tt>ef2e51ba</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/637">#637</a> from 99designs/fix-is-slice</summary>

Fix Mapping Custom Scalar to Slice

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e5ff6bc2c3bedb73df62f81a35056957c061f81c"><tt>e5ff6bc2</tt></a> add extra builtins types when no type exists

- <a href="https://github.com/99designs/gqlgen/commit/8225f63a89a8ee272d454399fa295785bb8416fc"><tt>8225f63a</tt></a> Allow plugins to provide additional template funcs

- <a href="https://github.com/99designs/gqlgen/commit/7b533df1a00be66e2ebe7a4268add80a5aa123e1"><tt>7b533df1</tt></a> Update ISSUE_TEMPLATE.md

- <a href="https://github.com/99designs/gqlgen/commit/055157f979f06a490979726db7ec20d8d9a634d4"><tt>055157f9</tt></a> Update ISSUE_TEMPLATE.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a148229cff7c5cb02b61caba81d7142ec9ef9948"><tt>a148229c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/644">#644</a> from Sauraus/master</summary>

Fix Gin installation instruction

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/52624e5372f260077e90645732dd24b8d1fad317"><tt>52624e53</tt></a> Fix Gin installation instruction</summary>

Current `go get gin` instruction results in an error from Go: `package gin: unrecognized import path "gin" (import path does not begin with hostname)`

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/515f22547466f30865a313e46d9ae499468d0cf1"><tt>515f2254</tt></a> Add test case for custom scalar to slice

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2284a3eb7c2d41458c60f2ca603220d4e7ffa80e"><tt>2284a3eb</tt></a> Improve IsSlice logic to check GQL def</summary>

Currently TypeReference.IsSlice only looks at the Go type to decide.
This should also take into account the GraphQL type as well, to cover
cases such as a scalar mapping to []byte

</details></dd></dl>

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.8.2"></a>
## [v0.8.2](https://github.com/99designs/gqlgen/compare/v0.8.1...v0.8.2) - 2019-03-18
- <a href="https://github.com/99designs/gqlgen/commit/ee06517c25deb254fa6708609ee5fd3fb3fbdbf2"><tt>ee06517c</tt></a> release v0.8.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8ac8a1f8142aaa0a4ce12d30c9680236efbffb03"><tt>8ac8a1f8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/635">#635</a> from 99designs/fix-inject-builtin-scalars</summary>

Only Inject Builtin Scalars if Defined in Schema

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d10e048e5be9ec34451a8c180d19cbcdcedc753c"><tt>d10e048e</tt></a> Add docs for built-in scalar implementations

- <a href="https://github.com/99designs/gqlgen/commit/d27e6eb65e8fb62ac9893c2992157a5050923c57"><tt>d27e6eb6</tt></a> Add example case for object type overriding builtin scalar

- <a href="https://github.com/99designs/gqlgen/commit/d567d5c8f737fc9870e4772e9e04e5e0dbe04e7a"><tt>d567d5c8</tt></a> Inject non-spec builtin values only if defined

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3e39b57a98fac8d4db1a85e0b998bd83e29fb9c5"><tt>3e39b57a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/634">#634</a> from 99designs/fallback-to-string</summary>

Use graphql.String for types wrapping a basic string

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a2cce0d14984402fccd9e8abb9f286a80e7296fb"><tt>a2cce0d1</tt></a> Use graphql.String for types wrapping a basic string

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fc05501b87c73faf116b1096490978254b4f2436"><tt>fc05501b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/633">#633</a> from 99designs/fix-union-pointers</summary>

Fix Having Pointers to Union Types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f02dabb7bad39e049323e05b171feb8f387d257a"><tt>f02dabb7</tt></a> Add test case for union pointer

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8257d423e8c2d936bf03de1aa4f6215dfdf7229c"><tt>8257d423</tt></a> Check Go type rather than GQL type for ptr</summary>

This is probably a more correct way to check whether we should wrap the
type in a pointer or not, rather than looking at the GrapQL definition.
There may be use-cases where a GraphQL interface/union might be mapped
to a Go stuct.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5df0938f0dfa6418457a23b31c5bd1b3bc7e879d"><tt>5df0938f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/628">#628</a> from 99designs/fix-ambient-imports</summary>

Move ambient imports into cmd package

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8e1590d784e27a0dc3fbb4779453ad3a732f1e53"><tt>8e1590d7</tt></a> Move ambient imports into cmd package</summary>

The getting started docs for dep suggest creating a local gqlgen script,
however these ambient import are in the root, so dep misses them.

This was changed in 0.8 but the ambient imports weren't moved.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/58744de96123b48aea4728f844b01694ff2b8ae9"><tt>58744de9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/622">#622</a> from 99designs/handle-complexity-root-collisions</summary>

Handle colliding fields in complexity root gracefully

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c889b3148f631f9e6a9467ee60cf3d21c0333ff8"><tt>c889b314</tt></a> Handle colliding fields in complexity root gracefully

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/26c395b02c59e39b264ab3f9750f845e660ad99c"><tt>26c395b0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/620">#620</a> from codyleyhan/cl/error</summary>

Allow user to supply path to gqlerror

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/12cf01aa60fef7be26b6c78ac3139b407ed7a455"><tt>12cf01aa</tt></a> Allow user to supply path to gqlerror

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/932322b6da351176889e41102f032202d1b0778c"><tt>932322b6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/619">#619</a> from 99designs/nil-slices</summary>

Support returning nulls from slices

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a48c55b24e22490ef3b6baffec247cf41b2f07f3"><tt>a48c55b2</tt></a> Support returning nulls from slices

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2b270e4d5f469bf764fff2e7c76ee3abc6d2aaa8"><tt>2b270e4d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/618">#618</a> from codyleyhan/cl/method</summary>

Adds way to determine if a resolver is a function call or value

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/af6dc16d8b2637c961fc5b8ad76418670a3df1e4"><tt>af6dc16d</tt></a> Add test for IsMethod in resolver

- <a href="https://github.com/99designs/gqlgen/commit/27e97535903e956eaf48acb08951b49539c7f80f"><tt>27e97535</tt></a> Expose IsMethod to resolver context

- <a href="https://github.com/99designs/gqlgen/commit/f52726dec03743d4e96aa0f56e2b6569d55beaba"><tt>f52726de</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ac2422e3ccf40361d7fc597c7e8d55df314c87c3"><tt>ac2422e3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/614">#614</a> from wesovilabs/master</summary>

Adding entry for workshop

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/db4f7255b11fc4468e7642894f007fde9e2b1102"><tt>db4f7255</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/613">#613</a> from icco/patch-2</summary>

Upgrade graphql-playground to 1.7.20

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/163bfc76c47e6663a8514020db26f67f0786c8b9"><tt>163bfc76</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/612">#612</a> from 99designs/maps-changesets</summary>

Maps as changesets

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6aa9dfc65af4ebd6c5373b854b40fe0a97f6dfbc"><tt>6aa9dfc6</tt></a> Adding entry for workshop

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/08f936e1ec91f1af2bca22de8511cc13b47c6e09"><tt>08f936e1</tt></a> Upgrade graphql-playground to 1.7.20</summary>

CSS didn't change but js did.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8fb1fafdcf641f4520b726d209cf2eb685dd69f3"><tt>8fb1fafd</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/611">#611</a> from 99designs/gqlparser-1.1.2</summary>

Bump gqlparser to 1.1.2

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/37983a5f1c799ee856fe59677578301fac3405d7"><tt>37983a5f</tt></a> remove some invalid test schema

- <a href="https://github.com/99designs/gqlgen/commit/765ff73865ac2d61eb8b6c9fe4534e5a4fbc072e"><tt>765ff738</tt></a> Add some docs on maps

- <a href="https://github.com/99designs/gqlgen/commit/0a92ca465691ea96186481471ec1ba01d6ecfaf8"><tt>0a92ca46</tt></a> Support map[string]interface{} in return types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ac56112b0cb83260f8f9b18ee5f03bb0af6f6905"><tt>ac56112b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/610">#610</a> from tgwizard/dynamic-complexity</summary>

Allow configuring the complexity limit dynamically per request

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a89050aa1f1ffb93950eccb48bf1156da9587850"><tt>a89050aa</tt></a> Bump gqlparser to 1.1.2

- <a href="https://github.com/99designs/gqlgen/commit/dd2881455f20f37ec601791bf09e32db9e928790"><tt>dd288145</tt></a> Allow configuring the complexity limit dynamically per request

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/485ddf3051a289874577c8384a22d7e58b199e72"><tt>485ddf30</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/605">#605</a> from 99designs/fix-default-scalars</summary>

Fix default scalars

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3ca2599adf049dba59dbf904de1ab0aeb14e2da3"><tt>3ca2599a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/606">#606</a> from jonatasbaldin/add-gin-recipe</summary>

Add Gin recipe

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/386eede91a4d5b3b2b35944d08716c4ed0e4886a"><tt>386eede9</tt></a> Add Gin recipe

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/22be59d159135252154e46e21b0e0cbdf0fb23b9"><tt>22be59d1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/604">#604</a> from cevou/arg-scalar</summary>

Fix directives on args with custom type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d02736dcdd9860685c2ea0c34a9a7e7626316d0f"><tt>d02736dc</tt></a> Added test for fix directives on args with custom type

- <a href="https://github.com/99designs/gqlgen/commit/30d235bc781c1976c08def2ca5283998f4786d76"><tt>30d235bc</tt></a> Fix default scalars

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d7b5dc283de948282815a8addaa4ef36d8253358"><tt>d7b5dc28</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/591">#591</a> from 99designs/fix-577</summary>

Fix mixed case name handling in ToGo, ToGoPrivate

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bef6c0a960bb9646f617af8c337b17b6c3e63da1"><tt>bef6c0a9</tt></a> Fix directives on args with custom type

- <a href="https://github.com/99designs/gqlgen/commit/bc386d79c8cee441951f71814be1a61bbd4b9a5b"><tt>bc386d79</tt></a> Fix mixed case name handling in ToGo, ToGoPrivate

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.8.1"></a>
## [v0.8.1](https://github.com/99designs/gqlgen/compare/v0.8.0...v0.8.1) - 2019-03-07
- <a href="https://github.com/99designs/gqlgen/commit/229185e45e9f411de393ee22f0daf0c30ad83812"><tt>229185e4</tt></a> release v0.8.1

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d872af63addc93b6b8e8df37a82bd33e312f5b59"><tt>d872af63</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/582">#582</a> from demdxx/master</summary>

Load the playground sources from HTTPS by default

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8e66832f1f8fad88af08b05e0f5f0a36a7b4e0a4"><tt>8e66832f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/589">#589</a> from 99designs/fix-autocasing-modelgen-bugs</summary>

Fix autocasing modelgen bugs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/de3b7cb825c9676c720826472c7c54b0c0301ce0"><tt>de3b7cb8</tt></a> Fix autocasing modelgen bugs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8e00703ebe967db83bf5544e169a8b0cc5895866"><tt>8e00703e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/588">#588</a> from 99designs/fix-default-scalar-implementation-regression</summary>

Fix default scalar implementation regression

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b27139ed5c290214d979e374a23689fc36eed78a"><tt>b27139ed</tt></a> Fix default scalar implementation regression

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/737a59a31d6663987127ff604f671d326e509337"><tt>737a59a3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/579">#579</a> from 99designs/fix-camelcase</summary>

Take care about commonInitialisms in ToGo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/52838cca91034b6e26733eaa896e3c426be50eed"><tt>52838cca</tt></a> fix ci

- <a href="https://github.com/99designs/gqlgen/commit/2c3783f18f3b679ece4fffed18fa839977e8359d"><tt>2c3783f1</tt></a> some refactoring

- <a href="https://github.com/99designs/gqlgen/commit/eb4536743c4dc507df32da1fe7a581052f7c438c"><tt>eb453674</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/dcd208d91603475f4cb3505e89a71aeb53d0c52f"><tt>dcd208d9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/584">#584</a> from 99designs/fix-deprecated-directive

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5ba8c8ead3ca891601215263b5969e4942c52a6b"><tt>5ba8c8ea</tt></a> Add builtin flag for build directives</summary>

These have an internal implementation and should be excluded from the
DirectiveRoot. In the future this may be a func that plugins could use
to add custom implementations.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b8526698d46f6a3e7726bfdd6c33163e06cb2ea1"><tt>b8526698</tt></a> Load the playground sources from HTTPS by default</summary>

For some browsers on non-secure domains resources from CDN doesn't loads, so I made all cdn.jsdelivr.net resources of the playground by HTTPS by default

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6ea48ff69aa370cc8fe7f498149f55859c3c11dd"><tt>6ea48ff6</tt></a> Take care about commonInitialisms in ToCamel

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1968a7bcfff5a640b4fb80aecc485d404bf596e3"><tt>1968a7bc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/576">#576</a> from jflam/patch-1</summary>

Update README.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/44becbbe3afc314ba19778a18b0e7733645eaf24"><tt>44becbbe</tt></a> Update README.md</summary>

Fixed typo in MD link ttps -> https

</details></dd></dl>

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.8.0"></a>
## [v0.8.0](https://github.com/99designs/gqlgen/compare/v0.7.2...v0.8.0) - 2019-03-04
- <a href="https://github.com/99designs/gqlgen/commit/f24e79d00425f1bbf13fbc79f0230ff4b2037955"><tt>f24e79d0</tt></a> release v0.8.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/55df9b8d926d238ede66a29cd7b38513ab2bb2f7"><tt>55df9b8d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/574">#574</a> from 99designs/next</summary>

v0.8.0

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/aedcc68ada4d9c299e6a1c96c56059787a193403"><tt>aedcc68a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/573">#573</a> from 99designs/plugin-docs</summary>

Very rough first pass at plugin docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8f91cf56cf39de4f3845395f3564224452d9b95b"><tt>8f91cf56</tt></a> Very rough first pass at plugin docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3d9ad75ebd191d578a5b8085f150be6c9e42b3d2"><tt>3d9ad75e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/572">#572</a> from 99designs/handle-nonexistant-directories-when-genreating-packagenames</summary>

Handle non-existant directories when generating default package names

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/08923334725c963ffcff948efe030881d658c38a"><tt>08923334</tt></a> Handle non-existant directories when generating default package names

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2ef4b443a873c8be7ccabdbba39b941d71cd64d1"><tt>2ef4b443</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/571">#571</a> from 99designs/automatically-bind-to-int32-int64</summary>

Automatically bind to int32 and int64

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2888e96c01bbcb4614f636bbc5dd9ecf4af45284"><tt>2888e96c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/570">#570</a> from 99designs/vendor-packages-workaround</summary>

Workaround for using packages with vendored code

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fb87dc3942e4179aa1eb8007f8f3d384fd4b9fca"><tt>fb87dc39</tt></a> Automatically bind to int32 and int64

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f2d9c3f74f11947f4ae1e6fdaa12d2095ab30518"><tt>f2d9c3f7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/569">#569</a> from 99designs/improve-introduction</summary>

Introduction Improvements

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1e7aab63fc134090a47d1f3634bab4d8f3d35291"><tt>1e7aab63</tt></a> Workaround for using packages with vendored code

- <a href="https://github.com/99designs/gqlgen/commit/5c692e294e71aa6846192ce71e90c983d760d114"><tt>5c692e29</tt></a> User README as canonical introduction

- <a href="https://github.com/99designs/gqlgen/commit/25bdf3d6bff9790ffb0eb8ef7b4422e4008386f7"><tt>25bdf3d6</tt></a> Consolidate Introduction documents

- <a href="https://github.com/99designs/gqlgen/commit/d81670d8bf568c4d5ed5b8db1274c3cf496db416"><tt>d81670d8</tt></a> Add initial contributing guidelines

- <a href="https://github.com/99designs/gqlgen/commit/d9a9a532178aff2def9992c39db37fcd79092fa2"><tt>d9a9a532</tt></a> playground: secure CDN resources with Subresource Integrity

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cb38b4be9c6baa228a67e897d1bdfa205142ea7b"><tt>cb38b4be</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/568">#568</a> from MichaelMure/secured-playground</summary>

playground: secure CDN resources with Subresource Integrity

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0258e1a29685aa8fd0adda5f9a7d57fb60f7ce22"><tt>0258e1a2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/565">#565</a> from steebchen/next</summary>

Fix cli config getters

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6ad1d97e52a5c003f8014dc0ed00663caf4a1118"><tt>6ad1d97e</tt></a> Move feature comparison

- <a href="https://github.com/99designs/gqlgen/commit/37cbbd6de2268ca9b0930abc3e39abc72e685e58"><tt>37cbbd6d</tt></a> playground: secure CDN resources with Subresource Integrity

- <a href="https://github.com/99designs/gqlgen/commit/da12fd11020c4d2449f5abb2545ce28d4dde75dd"><tt>da12fd11</tt></a> Fix cli config getters

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/51266b8f7ab04827838098a328ef9fcd70b545a2"><tt>51266b8f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/554">#554</a> from 99designs/fix-missing-recover</summary>

Recover from panics in unlikly places

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/67795c95b21a60fda303b35e509a447653bc8351"><tt>67795c95</tt></a> Recover from panics in unlikly places

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/56163b4584483d70ba5dd2cb5b968e63447c36e6"><tt>56163b45</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/553">#553</a> from 99designs/getting-started-0.8</summary>

Update Getting Started for 0.8 and Go Modules

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0bd120b5c27fa81fc3f60f662496bc3328f3ad86"><tt>0bd120b5</tt></a> Update dep code as well

- <a href="https://github.com/99designs/gqlgen/commit/6c5760320cc3834d46cd190fa906bdf7c033a5da"><tt>6c576032</tt></a> Update getting started with 0.8 generated code

- <a href="https://github.com/99designs/gqlgen/commit/ba761dcf38224e9a3da08bc678fe2e2415d0df28"><tt>ba761dcf</tt></a> Reintroduce main package in root

- <a href="https://github.com/99designs/gqlgen/commit/cdc575a23d5a8191e51176214a4aad0cbea5feb2"><tt>cdc575a2</tt></a> Update getting started with Go Modules support

- <a href="https://github.com/99designs/gqlgen/commit/378510e5cd93d65003b605a7451021e0d7d3b533"><tt>378510e5</tt></a> Move Getting Started above Configuration

- <a href="https://github.com/99designs/gqlgen/commit/d261b3fbb107b328ef6ffcfdbbc0e0903b1c5767"><tt>d261b3fb</tt></a> Fix navigation font weights

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/327a1a34f92860957148a011ce8e043224cf8cc5"><tt>327a1a34</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/551">#551</a> from 99designs/improved-collect-fields-api</summary>

Improved Collect Fields API and Documentation

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6439f197d0d84d22b1abf3068bd7982bb3c98c22"><tt>6439f197</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/552">#552</a> from 99designs/always-return-struct-pointers</summary>

Always return *Thing from resolvers for structs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/318639bbbb71985778ca2ee796bfe22fe231808e"><tt>318639bb</tt></a> Always return *Thing from resolvers for structs

- <a href="https://github.com/99designs/gqlgen/commit/e61b3e0be1fc52d94f0bbc80f7c486bb9eeb4f40"><tt>e61b3e0b</tt></a> Add Field Collection docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ef0223cfdf8a17e8dc5f83ef6a14e4f5a12ddd49"><tt>ef0223cf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/541">#541</a> from 99designs/fix-underscore-only-fields</summary>

Allow underscore only fields and naming collisions to be aliased explicitly

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/58b2c74f6d7dd5ecf4503d1adab5bcd3838aef73"><tt>58b2c74f</tt></a> drive by config fix

- <a href="https://github.com/99designs/gqlgen/commit/f6c52666a789164d9efcc22473a5911c56c34b00"><tt>f6c52666</tt></a> Add a test for aliasing different cases (closes <a href="https://github.com/99designs/gqlgen/issues/376"> #376</a>)

- <a href="https://github.com/99designs/gqlgen/commit/8c2d15ee698737bc1c524ecc11a96deb1c7253fc"><tt>8c2d15ee</tt></a> Fix underscore only fields (closes <a href="https://github.com/99designs/gqlgen/issues/473"> #473</a>)

- <a href="https://github.com/99designs/gqlgen/commit/0eb8b5c158eb2cfcad62444869ef077be1c7e1e0"><tt>0eb8b5c1</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/015d02ebca21137a2bf5aff0e379956da14c2628"><tt>015d02eb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/542">#542</a> from Elgarni/add-more-validation-checks-on-yml-config-file</summary>

Add more validation checks on .yml config file

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/647c62a555775a278f1df528a3e71d14ad39320f"><tt>647c62a5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/550">#550</a> from 99designs/fix-unstable-marshaler-func</summary>

Fix unstable external marshaler funcs with same name as type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3a8bf33f05dbdaacf8a2bba0114977eeab1615b9"><tt>3a8bf33f</tt></a> Add CollectAllFields test cases

- <a href="https://github.com/99designs/gqlgen/commit/9ebe77175f676c267fb9091042ec3da209479659"><tt>9ebe7717</tt></a> Fix unstable external marshaler funcs with same name as type

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a1195843c84e17593a5b22c811484bb0f0b3d63b"><tt>a1195843</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/544">#544</a> from enjoylife/fix-directive</summary>

Fix directives on fields with custom scalars

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/dc925c462705af9e61357d038b6a3f869f1f3157"><tt>dc925c46</tt></a> Added a test for config checking

- <a href="https://github.com/99designs/gqlgen/commit/b56cb659d36b090b0c26277c061d2b9a4ed2c4a9"><tt>b56cb659</tt></a> Refactored config check so that it runs after being normalized

- <a href="https://github.com/99designs/gqlgen/commit/dc6a7a36272c22a045dad6e8a7dd47d5f1a41a0c"><tt>dc6a7a36</tt></a> Add CollectAllFields helper method

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a2e61b3d628bf5e0f5b12e961f82a68fba066448"><tt>a2e61b3d</tt></a> Added a model and used directive on an input field within the integration schema</summary>

Added to the integration schema such that the build will catch the directive bug in question.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0b0e4a91a0827c66fd3d6f6b5a58177be96ce860"><tt>0b0e4a91</tt></a> Fix directives on fields with custom scalars

- <a href="https://github.com/99designs/gqlgen/commit/8ac0f6e4e0ae01485bf240286533c2cc1da42d81"><tt>8ac0f6e4</tt></a> Removed redundant semicolons

- <a href="https://github.com/99designs/gqlgen/commit/3645cd3ecf6c11cde04e65493afb7d3db34e04dd"><tt>3645cd3e</tt></a> Add more validation checks on .yml config file

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1b8b1ea16cb53e976af44f7d126446b5eba52d97"><tt>1b8b1ea1</tt></a> Fix typo in README</summary>

Fix typo in README in selection example directory to point to the selection example, not the todo example.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/66120d8fbff69609195a6483e53de00ccb5b54dd"><tt>66120d8f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/535">#535</a> from awiede/master</summary>

Fix typo in README

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fcacf200a365c1ad4201d17c8243a7d598765922"><tt>fcacf200</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b9819b21789c0543161f77dd811a8b208ec17f0f"><tt>b9819b21</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/540">#540</a> from 99designs/check-is-zero</summary>

Automatically convert IsZero to null

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/03a655dc9e9d7a1db9271a385b3c3b0854fceec7"><tt>03a655dc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/526">#526</a> from 99designs/union-fragment-bug</summary>

Union Fragment Bug Fix

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/99e9f41fd92a958226c960260de82cc0c0fd85f2"><tt>99e9f41f</tt></a> Use Implements for type Implementors in codegen

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ccca823f7257048049df495b88e0c33818db5a4b"><tt>ccca823f</tt></a> Separate out conditionals in collect fields</summary>

These conditions are not really related, and I missed the second
conditional when reading through the first time.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/efe8b026d562b95b6cd95f3df3b64e4624ed2b01"><tt>efe8b026</tt></a> Add reproducable test cases

- <a href="https://github.com/99designs/gqlgen/commit/306da15f43db5b796bf7a77e7897e5f8772a7fa9"><tt>306da15f</tt></a> Automatically convert IsZero to null

- <a href="https://github.com/99designs/gqlgen/commit/f81c61d3f8fa25b6cc1200f148dfefcd266810f7"><tt>f81c61d3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/539">#539</a> from 99designs/test-nullable-interface-pointers (closes <a href="https://github.com/99designs/gqlgen/issues/484"> #484</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f5200c80a42e5442c34dd97d4d680c240d5b4c46"><tt>f5200c80</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/498">#498</a> from vilterp/playground-content-type</summary>

add `content-type: text/html` header to playground handler

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/de148d133667b3cf0574b2054b8d5579f5c5db1a"><tt>de148d13</tt></a> Test for <a href="https://github.com/99designs/gqlgen/pull/484">#484</a>

- <a href="https://github.com/99designs/gqlgen/commit/9a48a007bd7ee36752a0879c075bb3470d243ebd"><tt>9a48a007</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/538">#538</a> from 99designs/test-input-marshalling (closes <a href="https://github.com/99designs/gqlgen/issues/487"> #487</a>)

- <a href="https://github.com/99designs/gqlgen/commit/7a82ab43a060b334d19c51a353a8924857dd39bb"><tt>7a82ab43</tt></a> Test for <a href="https://github.com/99designs/gqlgen/pull/487">#487</a>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/48a7e07f4fd28b8bca8c270ee824c5f99436cfaa"><tt>48a7e07f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/537">#537</a> from 99designs/stub-generation</summary>

Stub generation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/787b38d8589a5b4169464a46ab5bd1e87573221c"><tt>787b38d8</tt></a> Break testserver tests down into smaller files using stubs

- <a href="https://github.com/99designs/gqlgen/commit/c5e3dd44959d59bef6830cbe2652498e49b53089"><tt>c5e3dd44</tt></a> add stub generation plugin

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/43db679a9d39ed4939b2c06aeed5b6ee95b961cd"><tt>43db679a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/534">#534</a> from 99designs/multiple-bind-types</summary>

Multiple bind types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b26b915ea6f5a8529cc63c1f5285118426868510"><tt>b26b915e</tt></a> Move input validation into gqlparser see https://github.com/vektah/gqlparser/pull/96

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7d394222022613afcbae9b7dbc3083eb7ce0c1fe"><tt>7d394222</tt></a> Fix typo in README</summary>

Fix typo in README in selection example directory to point to the selection example, not the todo example.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/42131868df3225cc1960e205dc56f769c03bb285"><tt>42131868</tt></a> Linting fixes

- <a href="https://github.com/99designs/gqlgen/commit/956d03063de0e021ec7d23e4a6695fdafe0a01f6"><tt>956d0306</tt></a> Arg type binding

- <a href="https://github.com/99designs/gqlgen/commit/6af3d85da38a4768598f2a0a07be2da6c3f1b5f3"><tt>6af3d85d</tt></a> Allow multiple field bind types

- <a href="https://github.com/99designs/gqlgen/commit/3015624baf4f280a1b6812b39ddb1c2506ccd4a7"><tt>3015624b</tt></a> Regen dataloader with correct version

- <a href="https://github.com/99designs/gqlgen/commit/50f7d9c846f48a0824946aa4f3ca0b7d9fc572cd"><tt>50f7d9c8</tt></a> Add input field directives back in

- <a href="https://github.com/99designs/gqlgen/commit/8047b82ad25c4f71741a53360f814a576a06c4b3"><tt>8047b82a</tt></a> Fix nullability checks in new marshalling

- <a href="https://github.com/99designs/gqlgen/commit/b3f139c9374854d87d61a0628031136d1f215653"><tt>b3f139c9</tt></a> Cleanup field/method bind code

- <a href="https://github.com/99designs/gqlgen/commit/cf94d3ba6a4a7cd513454a8cb4b250f1db68ee53"><tt>cf94d3ba</tt></a> Removed named types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/82ded32137fcddf9b295bb283d118183bba983cb"><tt>82ded321</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/532">#532</a> from 99designs/fix-missing-json-content-type</summary>

Fix set header to JSON earlier in GraphQL response

Update the GraphQL handler to set the Response Header to JSON earlier for
error messages to be returned as JSON and not text/html.

Fixes https://github.com/99designs/gqlgen/issues/519

## Notes:
- Add checks for JSON Content-Type checks in decode bad queries tests

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b4c5a074b8250cb6470efbe18e5f2b028437fa37"><tt>b4c5a074</tt></a> Fix set header to JSON earlier in GraphQL response</summary>

Update the GraphQL handler to set the Response Header to JSON earlier for
error messages to be returned as JSON and not text/html.

Fixes https://github.com/99designs/gqlgen/issues/519

== Notes:
- Add checks for JSON Content-Type checks in decode bad queries tests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/533b08b698ae8254d97091ec9e6e0087f7ccaa30"><tt>533b08b6</tt></a> remove wonky input directives

- <a href="https://github.com/99designs/gqlgen/commit/60473555e7668de29c188448558123d9dc8edb3b"><tt>60473555</tt></a> Shared arg unmarshaling logic

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a7c8abe6d89935c129e982fc91e0afb6db07dc9f"><tt>a7c8abe6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/529">#529</a> from 99designs/websocket-keepalive</summary>

Add websocket keepalive support

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/555d7468922e2a27411e688f783dcda5c450554c"><tt>555d7468</tt></a> Remove TypeDefinition from interface building

- <a href="https://github.com/99designs/gqlgen/commit/cfa012de44b74987653560296ac9571a385c31dd"><tt>cfa012de</tt></a> Enable websocket connection keepalive by default

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c5b9b5a812e5b87b4504a64ac201f3f8f0d27d37"><tt>c5b9b5a8</tt></a> Use constant tick rate for websocket keepalive</summary>

Some clients (e.g. apollographql/subscriptions-transport-ws) expect a
constant tick rate for the keepalive, not just a keepalive after x
duration of inactivity.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/693753fcc69329fba282ad5f1d69c02979cbce08"><tt>693753fc</tt></a> Add websocket keepalive support

- <a href="https://github.com/99designs/gqlgen/commit/162afad73b653d9456d21ec38d00c3476ab2dde4"><tt>162afad7</tt></a> enums dont exist in runtime

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d0b6485b23f259588e9ecd34e23924c741d9c6f2"><tt>d0b6485b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/525">#525</a> from 99designs/stop-grc-panic</summary>

Stop GetResolverContext from panicking when missing

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/78cfff48aa5895a9eefe35615aae3b7b33ed11b1"><tt>78cfff48</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/528">#528</a> from 99designs/fix-todo-directive</summary>

Fix Todo Example Directive

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5e1bcfaf71d81d8e119758e3ec47e1cdfb0bf092"><tt>5e1bcfaf</tt></a> Remove parent check in directive</summary>

This should always be true, and currently has a bug when comparing
pointers to structs. Can just be removed.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c1b50cecda176fadfb98f8b84205058f8d1c7a5b"><tt>c1b50cec</tt></a> Stop GetResolverContext from panicking when missing

- <a href="https://github.com/99designs/gqlgen/commit/44aabbd3e92bc7582f92e74d1d657a626b0a1962"><tt>44aabbd3</tt></a> Move all build steps back into file containing defs

- <a href="https://github.com/99designs/gqlgen/commit/4e49d489cf589001179b68efe5ee2f5f1f823070"><tt>4e49d489</tt></a> Merge object build and bind

- <a href="https://github.com/99designs/gqlgen/commit/97764aec135905c4be05ad769b4f680a982b70ba"><tt>97764aec</tt></a> move generated gotpl to top

- <a href="https://github.com/99designs/gqlgen/commit/d380eccfd6721769835f1c3ed486ec5a1e5f9b2e"><tt>d380eccf</tt></a> promote args partial to full template

- <a href="https://github.com/99designs/gqlgen/commit/1bc51010eb895d675573cae4243207ea3aa6477e"><tt>1bc51010</tt></a> Everything is a plugin

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/055fb4bc9a6aae2ca92c50deb011b089d6fea5d0"><tt>055fb4bc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/514">#514</a> from 99designs/gomod</summary>

Add support for go modules

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/48eb6c521259c1b6f80c3849318d7878f1223bf6"><tt>48eb6c52</tt></a> Update appveyor

- <a href="https://github.com/99designs/gqlgen/commit/9e02a977daaba56f7c189a54bcc13f2d53e0d2f2"><tt>9e02a977</tt></a> fix integration test

- <a href="https://github.com/99designs/gqlgen/commit/251e8514d637b99bb0e52030dc25483f5e74085d"><tt>251e8514</tt></a> Add support for go modules

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/62175eab825058eadd79eddf40978db8a194a40b"><tt>62175eab</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/502">#502</a> from 99designs/model-plugin</summary>

Model plugin

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0f8844932c526fa14dc81ed14b036537cd5eca35"><tt>0f884493</tt></a> linting fixes

- <a href="https://github.com/99designs/gqlgen/commit/c6eb1a854225c2618ee6f281401a90a58f535c59"><tt>c6eb1a85</tt></a> Extract model generation into a plugin

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d3f1195ce28ba0dab5b7df9a4044777fbfc999db"><tt>d3f1195c</tt></a> add `content-type: text/html` header to playground handler</summary>

This ensures that the browser doesn't think it should download the page
instead of rendering it, if the handler goes through a gzipping
middleware.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f94b4b7883e0e8e0db78146a4b4efa8a50f917ff"><tt>f94b4b78</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/497">#497</a> from azavorotnii/small_fixes</summary>

Small fixes

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/21769d937f1d0f96d31e95d1adcd0f17da13b009"><tt>21769d93</tt></a> Ensure no side affect from preceding tests in wont_leak_goroutines test

- <a href="https://github.com/99designs/gqlgen/commit/10f4ccde2a5fd79d8cd0fd66d79ceed19dbcd2de"><tt>10f4ccde</tt></a> newRequestContext: remove redundant else part

- <a href="https://github.com/99designs/gqlgen/commit/a76e022803d49a5e4cdc95504766e2060fb7124a"><tt>a76e0228</tt></a> Add cache usage for websocket connection

- <a href="https://github.com/99designs/gqlgen/commit/940db1f962c24864fea7e41362f383e163bedd1a"><tt>940db1f9</tt></a> Fix cacheSize usage in handler

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fba9a37816b6ff55f80ded3b941a7e85390decb3"><tt>fba9a378</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/492">#492</a> from 99designs/unified-merge-pass</summary>

Unified merge pass

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a7f719a3982c24bd36c1f60545a146f59d494cbe"><tt>a7f719a3</tt></a> update appveyour to not rely on main

- <a href="https://github.com/99designs/gqlgen/commit/f46b7c8e38a8309243a46002cd9bb47b2f750dcb"><tt>f46b7c8e</tt></a> Reclaim main package for public interface to code generator

- <a href="https://github.com/99designs/gqlgen/commit/6b829037855aa73c0dc3f41553cb40e7604217d0"><tt>6b829037</tt></a> Extract builder object

- <a href="https://github.com/99designs/gqlgen/commit/87b37b0c30b98f311d18db56c404679bf40b68e5"><tt>87b37b0c</tt></a> Replace string based type comparisons with recursive types.Type check

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/82b1917d512e24a1160e0f65f8334257a8abc2eb"><tt>82b1917d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/490">#490</a> from 99designs/bind-directly-to-types</summary>

Bind directly to AST types, instead of copying out random bits

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1d86f9883da24d175611d2268b1ecf4485030a85"><tt>1d86f988</tt></a> extract argument construction

- <a href="https://github.com/99designs/gqlgen/commit/4b85d1b0363824814bab903fe8d86d661102aa53"><tt>4b85d1b0</tt></a> Merge buildInput into buildObject

- <a href="https://github.com/99designs/gqlgen/commit/db33d7b71c20737e2237139e46c9fc910c6a9d90"><tt>db33d7b7</tt></a> Extract graphql go merge into its own package

- <a href="https://github.com/99designs/gqlgen/commit/afc773b1133ae173c285e68c5b4d2a8f0e24209d"><tt>afc773b1</tt></a> Use ast definition directly, instead of copying

- <a href="https://github.com/99designs/gqlgen/commit/8298acb0903f76ae31f92c66ae251f53aa4c22fd"><tt>8298acb0</tt></a> bind to types.Types in field / arg references too

- <a href="https://github.com/99designs/gqlgen/commit/38add2c22a6bc13fecc9bc9180b5c48b208b4cc0"><tt>38add2c2</tt></a> Remove definition embedding, use normal field instead

- <a href="https://github.com/99designs/gqlgen/commit/950ff42c2668399e50d749b8d56e2208e9155142"><tt>950ff42c</tt></a> Bind to types.Type directly to remove TypeImplementation

- <a href="https://github.com/99designs/gqlgen/commit/70c852eb59e98136ae6595b0e756ae94000d5a85"><tt>70c852eb</tt></a> Add lookup by go type to import collection

- <a href="https://github.com/99designs/gqlgen/commit/eb1011617b052446b44239ec2fb9b6f6ee9cfdde"><tt>eb101161</tt></a> Remove aliased types, to be replaced by allowing multiple backing types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e79252b0ea6efabfd9131fab703089ed23ff2d38"><tt>e79252b0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/488">#488</a> from 99designs/refactor-config</summary>

Refactor config

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4138a3728d3dec69e509b09ff7f3fe7257f01cf2"><tt>4138a372</tt></a> rename generator receiver

- <a href="https://github.com/99designs/gqlgen/commit/bec38c7e304c2bdc2e4ff32661e1d1972cdb1484"><tt>bec38c7e</tt></a> Extract config into its own package

- <a href="https://github.com/99designs/gqlgen/commit/34b878713c95839e085cef8e8fbe377cf7dae725"><tt>34b87871</tt></a> Rename core types to have clearer meanings

- <a href="https://github.com/99designs/gqlgen/commit/f10fc649fe4db6a10d048d2bcf2520121a4bbcac"><tt>f10fc649</tt></a> Merge remote-tracking branch 'origin/next' into HEAD

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/dd972081ca802d30d94a9a6e42879491c572edca"><tt>dd972081</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/486">#486</a> from nicovogelaar/feature/list-of-enums</summary>

add list of enums

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1140dd85782572b148e8ec2c1d9526f87b763a1e"><tt>1140dd85</tt></a> add unit test for list of enums

- <a href="https://github.com/99designs/gqlgen/commit/1e3e5e9b36587311b8b632668f83b53278d3f8eb"><tt>1e3e5e9b</tt></a> add list of enums

- <a href="https://github.com/99designs/gqlgen/commit/f87ea6e85f867d169760c1bb2a73842f82194363"><tt>f87ea6e8</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/473f4f0c566042e93d7ae9ae3cee49697386a766"><tt>473f4f0c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/465">#465</a> from 99designs/performance-improvments</summary>

Performance improvments

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f9ee6ce0a09fd587c5786e18d2f089df83a10b55"><tt>f9ee6ce0</tt></a> return arg in middleware

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5c8b1e24546c78ddc00a05e2605f894d66f2911f"><tt>5c8b1e24</tt></a> Avoid unnessicary goroutines</summary>

goos: linux
goarch: amd64
pkg: github.com/99designs/gqlgen/example/starwars
BenchmarkSimpleQueryNoArgs-8      300000             25093 ns/op            6453 B/op        114 allocs/op
PASS
ok      github.com/99designs/gqlgen/example/starwars    10.807s

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b0ffa22a8a06fb912b5ef36260b2025bf66e1356"><tt>b0ffa22a</tt></a> Remove strconv.Quote call in hot path to avoid some allocs</summary>

go test -benchtime=5s -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/99designs/gqlgen/example/starwars
BenchmarkSimpleQueryNoArgs-8      200000             32125 ns/op            6277 B/op        118 allocs/op
PASS
ok      github.com/99designs/gqlgen/example/starwars    9.768s

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2cf5a5b8049198e130ccdcd8068cfe321b82e17e"><tt>2cf5a5b8</tt></a> Add a benchmark</summary>

go test -benchtime=5s -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/99designs/gqlgen/example/starwars
BenchmarkSimpleQueryNoArgs-8      200000             32680 ns/op            6357 B/op        126 allocs/op
PASS
ok      github.com/99designs/gqlgen/example/starwars    9.901s

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5e0456febd0b1c61fa09fe16f1dbd1d8a56c5b84"><tt>5e0456fe</tt></a> fix fmt anf metalint generated code

- <a href="https://github.com/99designs/gqlgen/commit/b32ebe1480704a34bc941553c689e0c4d746bcd5"><tt>b32ebe14</tt></a> check nullable value for go1.10

- <a href="https://github.com/99designs/gqlgen/commit/d586bb61f1b4da0a9a916118809cb6dc0c996bbd"><tt>d586bb61</tt></a> use arg value for the ResolveArgs

- <a href="https://github.com/99designs/gqlgen/commit/e201bcb5540ecae0c9c9d96868cfce1ba2f79227"><tt>e201bcb5</tt></a> set default nil arg to ResolverContext

- <a href="https://github.com/99designs/gqlgen/commit/6fa6364004f9140e62b2b873f0fb5bf344dd68f7"><tt>6fa63640</tt></a> remove empty line in generated files

- <a href="https://github.com/99designs/gqlgen/commit/139ed9fb048a34598cc153f3e8ef3b6b6fa8b220"><tt>139ed9fb</tt></a> fix go10 assign exist variable by eq

- <a href="https://github.com/99designs/gqlgen/commit/428c6300f3f6ddbde7182e56a096307447686109"><tt>428c6300</tt></a> add nullable argument to directives

- <a href="https://github.com/99designs/gqlgen/commit/740960331927b0d7af1e80fca0f3bc4825fa74de"><tt>74096033</tt></a> move chainFieldMiddleware to generate code for BC

- <a href="https://github.com/99designs/gqlgen/commit/be51904c52bfec461626a7765b0c8ef2f063ffe6"><tt>be51904c</tt></a> check nullable arguments

- <a href="https://github.com/99designs/gqlgen/commit/6b0050940400bc02da4ccb6827e6ce1237f7b221"><tt>6b005094</tt></a> add test directives generate

- <a href="https://github.com/99designs/gqlgen/commit/047f2ebcaebd043aed4a38cba282b2f5b5602521"><tt>047f2ebc</tt></a> update inline template

- <a href="https://github.com/99designs/gqlgen/commit/a13b31e9b40aa9098fcf90fb947b9d854b53dbe1"><tt>a13b31e9</tt></a> metalinter

- <a href="https://github.com/99designs/gqlgen/commit/526bef0bafb42b15952e724570dd16c7366da34e"><tt>526bef0b</tt></a> generate servers and add path to error

- <a href="https://github.com/99designs/gqlgen/commit/29770d6485a69abf86585da6ecb68fb4fe75b4fa"><tt>29770d64</tt></a> resolve = in template

- <a href="https://github.com/99designs/gqlgen/commit/3a729cc3c60f1aafdee405513d43cdf43693269d"><tt>3a729cc3</tt></a> update recursive middleware

- <a href="https://github.com/99designs/gqlgen/commit/8b3e634e5cc1ba0bf4317629f05872752d0f2ba7"><tt>8b3e634e</tt></a> update tempate and set Dump public

- <a href="https://github.com/99designs/gqlgen/commit/e268bb75797fd9cb8a753cc0312fff3e564a20d5"><tt>e268bb75</tt></a> Merge remote-tracking branch 'upstream/master' into directives

- <a href="https://github.com/99designs/gqlgen/commit/e8f0578de98539ea16fea15f0033b83615192f96"><tt>e8f0578d</tt></a> add execute ARGUMENT_DEFINITION and INPUT_FIELD_DEFINITION directive

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.7.2"></a>
## [v0.7.2](https://github.com/99designs/gqlgen/compare/v0.7.1...v0.7.2) - 2019-02-05
- <a href="https://github.com/99designs/gqlgen/commit/da1e07f5876c0fb79cbad19006f7135be08590d6"><tt>da1e07f5</tt></a> release v0.7.2

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8c0562c17743ea26cc316e9ff4cd509054b35287"><tt>8c0562c1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/530">#530</a> from 99designs/websocket-keepalive-master</summary>

Add websocket keepalive support

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/43fdb7da02b53094b44ce6a268a3e845bedd0967"><tt>43fdb7da</tt></a> Suppress staticcheck lint check on circleci

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9c4b877afd9d6396d110957c9947eb62ad9409b7"><tt>9c4b877a</tt></a> Use constant tick rate for websocket keepalive</summary>

Some clients (e.g. apollographql/subscriptions-transport-ws) expect a
constant tick rate for the keepalive, not just a keepalive after x
duration of inactivity.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d36d3dc567476909853fd3bb0a0b3ba28b24ed9f"><tt>d36d3dc5</tt></a> Add websocket keepalive support

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/39216361225c6dc824331ad1a218c9c931cc0985"><tt>39216361</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/476">#476</a> from svanburen/patch-1</summary>

Update config.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9f6f2bb8a58dfc8d8da6714297a460d7d85e55cd"><tt>9f6f2bb8</tt></a> Update config.md</summary>

Add a missed word and add an apostrophe

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c033f5fcfdcab6be940356b3d99e3924e7919314"><tt>c033f5fc</tt></a> Fix edit link positioning

- <a href="https://github.com/99designs/gqlgen/commit/b3f163d828f65797884a15197f31e49a17f55408"><tt>b3f163d8</tt></a> Add not about relative generate path

- <a href="https://github.com/99designs/gqlgen/commit/675ba773946c772b3bc405f7f9fb7cfe4c8b9a47"><tt>675ba773</tt></a> Update errors.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5c870a489da3eecd57b4152f0501b65cd682f4af"><tt>5c870a48</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/461">#461</a> from ryota548/patch-1</summary>

Update getting-started.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9bcd27c1a179a8ab2ac546287f1d699b11245af7"><tt>9bcd27c1</tt></a> Update getting-started.md</summary>

modify `graph/graph.go` to `resolver.go`

</details></dd></dl>

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.7.1"></a>
## [v0.7.1](https://github.com/99designs/gqlgen/compare/v0.7.0...v0.7.1) - 2018-11-29
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3a7f37c7e22a8fedce430c4d340ad5c1351198f4"><tt>3a7f37c7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/455">#455</a> from 99designs/fix-deprecated-fields</summary>

Fix deprecated fields

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b365333ba4015017a05664d382773954be7d71db"><tt>b365333b</tt></a> Fix graphiql deprecating all fields

- <a href="https://github.com/99designs/gqlgen/commit/99610be997bb8f906c2f27cfe010e69693ad2e9e"><tt>99610be9</tt></a> Get chat example up to date

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.7.0"></a>
## [v0.7.0](https://github.com/99designs/gqlgen/compare/v0.6.0...v0.7.0) - 2018-11-28
- <a href="https://github.com/99designs/gqlgen/commit/a81fe5037b2492cdd312a7d8c875677da4b1f6c9"><tt>a81fe503</tt></a> release v0.7.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4bfc82d782409044f07e002c926f867fdb14ac8d"><tt>4bfc82d7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/453">#453</a> from 99designs/deprecate-binary</summary>

Add Deprecation Warning to Binary

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8dd29b8548320a3a01a8b7645bc79d5b216edd62"><tt>8dd29b85</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/454">#454</a> from 99designs/update-gqlparser</summary>

Update gqlparser to latest

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/747c3f9c41903a77781f50d919335d4b29a215b8"><tt>747c3f9c</tt></a> Update gqlparser to latest

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d6d9885fac16079f92aeda0f6571cc0e7697b0ac"><tt>d6d9885f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/416">#416</a> from 99designs/improved-getting-started</summary>

Improve Getting Started Documentation — No Binary Approach

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d22f03c62d582f36872c85507b75d095c0ec4fc9"><tt>d22f03c6</tt></a> Add deprecation warning

- <a href="https://github.com/99designs/gqlgen/commit/878f3945f1ec84aa26400f67a714fd8cc7db40e3"><tt>878f3945</tt></a> Minor fixes to getting started code examples

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6a02657c5a257173034efe93fa8cecd49da7d990"><tt>6a02657c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/447">#447</a> from 99designs/disable-introspection</summary>

Add config option to disable introspection

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b9fbb642f1f90bf8b173cfabdfc88d24ec36344d"><tt>b9fbb642</tt></a> Mention recursive-ness of generate ./...

- <a href="https://github.com/99designs/gqlgen/commit/e236d8f36d8b09e034928761523a21dead7014d5"><tt>e236d8f3</tt></a> Remove generate command from resolver.go

- <a href="https://github.com/99designs/gqlgen/commit/04a72430f5550aa0a052aad20878ccf401eb3b23"><tt>04a72430</tt></a> Re-add final touches section to getting started

- <a href="https://github.com/99designs/gqlgen/commit/3a7a506259cce1a7da4ea9ebe74e1358d3c7bed0"><tt>3a7a5062</tt></a> Add handler import to root cmd

- <a href="https://github.com/99designs/gqlgen/commit/9dba96d51a60114b2a33a2c296dde02f2306fe41"><tt>9dba96d5</tt></a> Fix GraphQL capitalisation

- <a href="https://github.com/99designs/gqlgen/commit/1dfaf637be0b4f7ec2d4063f383e8288c5994c25"><tt>1dfaf637</tt></a> Minor updates to getting started from feedback

- <a href="https://github.com/99designs/gqlgen/commit/94b95d976b490a2342f6f9a18f7b7b36f0e996e7"><tt>94b95d97</tt></a> Some CSS fixes

- <a href="https://github.com/99designs/gqlgen/commit/a36fffd213e5860a4d57604829811bc4fc4c0ec7"><tt>a36fffd2</tt></a> Updated getting started with new no-binary approach

- <a href="https://github.com/99designs/gqlgen/commit/601354b364911267763ac8536d9a738edf15a147"><tt>601354b3</tt></a> Add blockquote breakout style

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6bea1d88358ced8438311deb7ac488674ebdecf4"><tt>6bea1d88</tt></a> Merge remote-tracking branch 'origin/master' into disable-introspection</summary>

Regenerate

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e4bad0e6bd11386d8436e1c268a02ed875b52ef5"><tt>e4bad0e6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/449">#449</a> from 99designs/increase-float-precision</summary>

Increase float precision

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c5589792b0f69c28a536a13c12dbf2a8f5e11b01"><tt>c5589792</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/450">#450</a> from 99designs/import-refactor</summary>

Refactor import handling

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/62f0d085b5a019254b7699e2f02ee0e61c96f52d"><tt>62f0d085</tt></a> Edit copy for introspection docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/63fc2753eb74995bcb65f64914bbd913114cf4da"><tt>63fc2753</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/452">#452</a> from cemremengu/patch-1</summary>

Fix typo in directives.md

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/da31e8eda2125363b99cc84890ba9da9e0e6cf3f"><tt>da31e8ed</tt></a> Update directives.md</summary>

Fix small typo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/83e33c135ca0a065c625922cdc8d46808cb73107"><tt>83e33c13</tt></a> Remove a debug print

- <a href="https://github.com/99designs/gqlgen/commit/6c57591462cae14ddb9dfff8279e4e4493c44b33"><tt>6c575914</tt></a> fix doc indentation

- <a href="https://github.com/99designs/gqlgen/commit/f03b32d3539a231e1c572a92a52090b73a8a3762"><tt>f03b32d3</tt></a> Use new import handling code

- <a href="https://github.com/99designs/gqlgen/commit/c45546e57a1f73a57d763cde673fe48462314191"><tt>c45546e5</tt></a> Increase float precision

- <a href="https://github.com/99designs/gqlgen/commit/77f2e2847d8c171563e824cfadc402e8001882ac"><tt>77f2e284</tt></a> Start moving import management to templates

- <a href="https://github.com/99designs/gqlgen/commit/c114346d88dbe979b7e5e7c09eacef2bba0f7500"><tt>c114346d</tt></a> Decouple loader creation from schema

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9d636e780e5b1aa3a1290b41e590d207764848a5"><tt>9d636e78</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/448">#448</a> from 99designs/update-gqlparser</summary>

Update to latest gqlparser

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d6ce42df131b868afdec66a00966f7a334233c16"><tt>d6ce42df</tt></a> Update to latest gqlparser

- <a href="https://github.com/99designs/gqlgen/commit/b0acd078ac03f165a4d48a18e8e1743086e80270"><tt>b0acd078</tt></a> Add config option to disable introspection

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f9c880b6ee3a4423320468c1304237af7ee4e8b6"><tt>f9c880b6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/446">#446</a> from 99designs/fix-flakey-test</summary>

Fix flakey goroutine test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5461e967f847df9d8c51d76f193c69b31882e89f"><tt>5461e967</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/445">#445</a> from 99designs/remove-graphqlgen</summary>

Remove graphqlgen link

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8a5039d8fe5157003196c995505ae8d1b1a59676"><tt>8a5039d8</tt></a> Fix flakey goroutine test

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4b082518034eeae70f74c32ed77b6d590853df25"><tt>4b082518</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/439">#439</a> from snormore/pointer-slice</summary>

Fix type binding validation for slices of pointers like []*foo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/293b9eaf5c88e3d524da38cd7421ed0ba038d1b2"><tt>293b9eaf</tt></a> Remove graphqlgen link

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/77b27884d0d4d5238d9e72ba881ae1b142e6abca"><tt>77b27884</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/443">#443</a> from mgutz/patch-1</summary>

fix generate stubs sentence

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ae1c77327b5b91900ac3531521ee0cdf2d56501a"><tt>ae1c7732</tt></a> fix generate stubs sentence

- <a href="https://github.com/99designs/gqlgen/commit/827dac5e0991bb6368f6d47642e9b1e7a232cf4d"><tt>827dac5e</tt></a> Fix type binding validation for slices of pointers like []*foo

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f7932b40ee0f75b35428af917cd05193b7f3414f"><tt>f7932b40</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/435">#435</a> from matiasanaya/update-readme</summary>

Update README.md comparison with graph-gophers

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a816208b938269ed569bf6d9b8e9884c28547dfc"><tt>a816208b</tt></a> Update README.md comparison with graph-gophers

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d25e3b4b2ba627f306dcd703cb98506f3268837b"><tt>d25e3b4b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/422">#422</a> from gracenoah/model-method-context</summary>

accept an optional ctx parameter on model methods

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0ac6fa5758104b70852eb0c90317fa21f1cd4ecf"><tt>0ac6fa57</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/434">#434</a> from urakozz/patch-1</summary>

Tracer: fixed nil pointer issue

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d4f7c954a52a1d704fce96971e6ff0c179bb3437"><tt>d4f7c954</tt></a> Update context.go

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4c4ccf471f769a402b2231519fd2fcb15d6e5e4a"><tt>4c4ccf47</tt></a> Update context.go</summary>

Right now code generated with latest master fails since there are usages of Trace but there is no any single write to this variable

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5faf3a2bdf592c8445655b915c56e2762e305dc8"><tt>5faf3a2b</tt></a> re-generate

- <a href="https://github.com/99designs/gqlgen/commit/6fed89478253c95b400cb55dfb8d6b7fe47fb776"><tt>6fed8947</tt></a> rebase fixes

- <a href="https://github.com/99designs/gqlgen/commit/4c10ba55bb493839605673a63701a4cb84b364e2"><tt>4c10ba55</tt></a> fix generated code

- <a href="https://github.com/99designs/gqlgen/commit/8066edb719066f3b306b6e6165ff9e959f298deb"><tt>8066edb7</tt></a> add tests

- <a href="https://github.com/99designs/gqlgen/commit/9862c30f17df6dbcc1c3712dc3cff523950ee918"><tt>9862c30f</tt></a> mention contexts on model methods in docs

- <a href="https://github.com/99designs/gqlgen/commit/602a83d6f1ac630b0abf6732375643b48f1b38db"><tt>602a83d6</tt></a> make ctx method resolvers concurrent

- <a href="https://github.com/99designs/gqlgen/commit/497551202238b0befd67d69b1da547bfff660948"><tt>49755120</tt></a> accept an optional ctx parameter on model methods

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/02a1935255fadced7f122dff7f2a9f54546c9d61"><tt>02a19352</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/429">#429</a> from 99designs/refactor-gofmt</summary>

apply go fmt ./...

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6a77af136c090502fdffbc0ca19f731080c26021"><tt>6a77af13</tt></a> apply gofmt on ./.circleci/test.sh

- <a href="https://github.com/99designs/gqlgen/commit/c656dc3127d29b102bdcd286ac94dc00b2a9600a"><tt>c656dc31</tt></a> apply go fmt ./...

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3f598bdc8c72c9f20cb47046aac6837090005452"><tt>3f598bdc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/427">#427</a> from anurag/patch-1</summary>

Fix docs typo

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cac61bb2035545137ac96c9d0ac102b71f2f3169"><tt>cac61bb2</tt></a> Fix docs typo

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9f4afe3a6eb09e9ba12fbf5591ce3c0f06ad48be"><tt>9f4afe3a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/425">#425</a> from 99designs/render</summary>

Switch to hosting docs on render.com

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9875e74bdc2d071773429111d632eec38f978320"><tt>9875e74b</tt></a> Switch to hosting docs on render.com</summary>

Render.com has offered to host our static site for free, and have
a pretty simple setup for rebuilding on merge to master. I've
switched the DNS records and updated the docs.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/981fd10a677716bd3e5b803465e2084e4fac3723"><tt>981fd10a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/419">#419</a> from 99designs/fix-capture-ctx</summary>

fix unexpected ctx variable capture on Tracing

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/027803d23914a1b23082fe7391972a890565d24b"><tt>027803d2</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/2b090de9ebe65dc3c23642d5c6b57d5e7d40d0de"><tt>2b090de9</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/d3238d54013d07d5b543aa01b453fdafa6ac7b3d"><tt>d3238d54</tt></a> chore

- <a href="https://github.com/99designs/gqlgen/commit/a2c33f13a501e0b2058f6730cf1ba104a72edfda"><tt>a2c33f13</tt></a> write ctx behavior test & refactoring tracer test

- <a href="https://github.com/99designs/gqlgen/commit/5c28d0116ee17a921ab891f06e555fda7cf7ca61"><tt>5c28d011</tt></a> fix unexpected ctx variable capture on Tracing

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4bda3bc1291bdc0bc44e3057b5229d987eeecde2"><tt>4bda3bc1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/411">#411</a> from 99designs/feat-geterrors</summary>

add GetErrors to RequestContext

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a4eaa400c2cdfd0aa6fd4471a5899a93deb460f4"><tt>a4eaa400</tt></a> add tests for RequestContext#GetErrors

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/53f33f7722464e3063ecd15f053dfe3d79928dff"><tt>53f33f77</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/410">#410</a> from 99designs/move-tracing-to-contrib</summary>

Move tracing to contrib

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/19403832ced12b9c610fb773a184eacf4ec8e3f6"><tt>19403832</tt></a> add GetErrors to RequestContext

- <a href="https://github.com/99designs/gqlgen/commit/f0dbce5a30b444ece286db2b7bac21dca01de174"><tt>f0dbce5a</tt></a> Move tracing to contrib

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a3a92775eee4365784bcf35a5cb550027c680fd7"><tt>a3a92775</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/409">#409</a> from 99designs/graphql-playground-1.7.8</summary>

Bump to the latest version of graphql-playground

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d2648580b441915b0a57b37b70b0e136bc20ea61"><tt>d2648580</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/402">#402</a> from 99designs/feat-opencensus</summary>

add Tracer for OpenCensus

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7286e2445ea6c30040bd12c87f92d36cd21faed5"><tt>7286e244</tt></a> fix shadowing

- <a href="https://github.com/99designs/gqlgen/commit/af38cc5a74ff978e1e266a1e0bb5d0a50be4dd4f"><tt>af38cc5a</tt></a> Bump to the latest version of graphql-playground

- <a href="https://github.com/99designs/gqlgen/commit/8bbb5eb79bfcc448851883510aa01d99de3f721c"><tt>8bbb5eb7</tt></a> fix some tests

- <a href="https://github.com/99designs/gqlgen/commit/256e741f8ae41174693757544aea18568f2f8226"><tt>256e741f</tt></a> add complexityLimit and operationComplexity to StartOperationExecution

- <a href="https://github.com/99designs/gqlgen/commit/4e7e6a1c7167e938e611dcaffa9c25b34e4ecc02"><tt>4e7e6a1c</tt></a> Merge branch 'master' into feat-opencensus

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/926ad17adfe84c26db5f4882c289941fe654af31"><tt>926ad17a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/403">#403</a> from 99designs/feat-complexity</summary>

copy complexity to RequestContext

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2d3026cbaf390d61791e9d892cfc629365e3b137"><tt>2d3026cb</tt></a> Merge branch 'master' into feat-complexity

- <a href="https://github.com/99designs/gqlgen/commit/59ef91ad7bb0790d361f720b89cac9accea00106"><tt>59ef91ad</tt></a> merge master

- <a href="https://github.com/99designs/gqlgen/commit/c9368904b5cab04655a1ce2666346c0070f37ed0"><tt>c9368904</tt></a> Merge branch 'master' into feat-opencensus

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b26ee6b4ea235d699f34eee09cd569d814a4bafc"><tt>b26ee6b4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/404">#404</a> from 99designs/feat-apollo-tracing</summary>

add apollo-tracing support

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fd4f55877352659feabd73640fb0816caec3dee9"><tt>fd4f5587</tt></a> fix timing issue

- <a href="https://github.com/99designs/gqlgen/commit/91e3e88d0f212865d9d2f0bd907330eb61974edf"><tt>91e3e88d</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/a905efa88e85bed2f9df597c5bf8581da4a7fab1"><tt>a905efa8</tt></a> fix lint warning

- <a href="https://github.com/99designs/gqlgen/commit/b2ba5f86b7e668dcd0c13a490a55015a1ec3fb88"><tt>b2ba5f86</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/561be1c070c2d3237b38235bc81fd1c4f18d153b"><tt>561be1c0</tt></a> add Apollo Tracing sample implementation

- <a href="https://github.com/99designs/gqlgen/commit/83c7b2cba6c21075913f46515b5ff1483ca619e0"><tt>83c7b2cb</tt></a> add Start/EndOperationParsing & Start/EndOperationValidation methods to Tracer

- <a href="https://github.com/99designs/gqlgen/commit/b5305d75c79e9749ea90c7487a54de88ca61be28"><tt>b5305d75</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/784dc01fdb4f0759e59e32bb48814b94760ca00b"><tt>784dc01f</tt></a> oops...

- <a href="https://github.com/99designs/gqlgen/commit/a027ac21c773ed1bf71ec6017e5cafbd140305a2"><tt>a027ac21</tt></a> copy complexity to RequestContext

- <a href="https://github.com/99designs/gqlgen/commit/ececa23c60cafd25454d0c2d45f89f6e0549b8f4"><tt>ececa23c</tt></a> add Tracer for OpenCensus

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0d5c65b6dc1c80d7809b862fd0b9ad3247926b0f"><tt>0d5c65b6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/400">#400</a> from 99designs/fix-ci</summary>

fix Circle CI test

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/00d11794f4779c4d0755336ecfaf6547a84306da"><tt>00d11794</tt></a> add mutex to logger

- <a href="https://github.com/99designs/gqlgen/commit/884d35c6d7519dd3910e95f2522db6d79380e991"><tt>884d35c6</tt></a> fix race condition

- <a href="https://github.com/99designs/gqlgen/commit/f70cedc2bac7e4582e9b50d8b34cb49d41dcc8d2"><tt>f70cedc2</tt></a> fix Circle CI test

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1b17b5a2489094911e701a299aedee3d1d1a2319"><tt>1b17b5a2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/392">#392</a> from 99designs/feat-tracer</summary>

Add Tracer layer

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/184e48cbd3c1391a2dd1434b17dc0bcceaf41661"><tt>184e48cb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/396">#396</a> from 99designs/remove-ci-exclusion</summary>

Run generate ./... and test ./... in circle

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fd5d9ecae750465e74475112f6af0496404db87b"><tt>fd5d9eca</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/395">#395</a> from 99designs/feat-extension-example</summary>

add Type System Extension syntax example

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/686c71a47162855e93a9b8c69d3102b4c1ce63ee"><tt>686c71a4</tt></a> Run generate ./... and test ./... in circle

- <a href="https://github.com/99designs/gqlgen/commit/304d3495819f8cbe163d664ea03736b1dd0107e4"><tt>304d3495</tt></a> fix https://github.com/99designs/gqlgen

- <a href="https://github.com/99designs/gqlgen/commit/85322586f166296310dcd6f0855cbda9c65a8362"><tt>85322586</tt></a> address comment

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/195f952b91dbeec23f39ba2d91c9bf6c96f42423"><tt>195f952b</tt></a> fix CI failed</summary>

AppVeyor handle this test, But Circle CI is not

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b5b767c42cae2a9eaebdde1afe9dde1f3accc412"><tt>b5b767c4</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/d723844b9bc3b4115b837399e6d4705ddbf4c0cb"><tt>d723844b</tt></a> add Type System Extension syntax example

- <a href="https://github.com/99designs/gqlgen/commit/df685ef7f771453efed1c4bbfe47bf740202a7f3"><tt>df685ef7</tt></a> change timing of EndFieldExecution calling

- <a href="https://github.com/99designs/gqlgen/commit/94b7ab02b4ca683d4385717433ea5dae2b6138d6"><tt>94b7ab02</tt></a> refactor Tracer interface signature that fit to apollo-tracing specs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8eb2675a439e98293a8d799dbd78290f6137d3c2"><tt>8eb2675a</tt></a> Revert "change field marshaler return process that make it easy to insert other processing"</summary>

This reverts commit 583f98047f5d1b6604d87e7b8d6f8fd38082d459.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c8af48cdaec63ab37226bcc2f2f4b1e071d8e709"><tt>c8af48cd</tt></a> rename Tracer method name

- <a href="https://github.com/99designs/gqlgen/commit/a3060e80a1926cffb080932725262f540a6e55f3"><tt>a3060e80</tt></a> refactor Tracer signature

- <a href="https://github.com/99designs/gqlgen/commit/d319afe6623c04e99f1186685c0384e43103a790"><tt>d319afe6</tt></a> add support request level tracer

- <a href="https://github.com/99designs/gqlgen/commit/1c5aedde3d72e4dd306c0120e6031d66d359d592"><tt>1c5aedde</tt></a> add support field level tracer

- <a href="https://github.com/99designs/gqlgen/commit/583f98047f5d1b6604d87e7b8d6f8fd38082d459"><tt>583f9804</tt></a> change field marshaler return process that make it easy to insert other processing

- <a href="https://github.com/99designs/gqlgen/commit/ab4752c28debbe543601e16dc7861b7973407b5c"><tt>ab4752c2</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3447dd2d8b361c532d38688474d04315788edec1"><tt>3447dd2d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/389">#389</a> from 99designs/multiple-schemas</summary>

Support multiple schemas

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a230eb049370be09434b451ed8e913e68a134ad1"><tt>a230eb04</tt></a> Support multiple schemas

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/20a5b6c7f2f6a9feabe96dcde68ce0b6d23f4982"><tt>20a5b6c7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/369">#369</a> from vetcher/master</summary>

reverse errors and data order in response

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f1f043b9d4e97d3dc8ad1559110e8ff8688b4afe"><tt>f1f043b9</tt></a> reverse 'data' and 'error' fields order in failure tests

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3eab22a33bf8a9ff99482c899dd184c589445c7b"><tt>3eab22a3</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/370">#370</a> from rodrigo-brito/fix-underscore</summary>

Underscore on field name finder

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0ad3d3ce71b852dbb62a72babf3f5352eae47100"><tt>0ad3d3ce</tt></a> fix on struct name finder

- <a href="https://github.com/99designs/gqlgen/commit/42e110453498b2690ce4412b384881a3bf55d0c5"><tt>42e11045</tt></a> reverse errors and data order in response

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.6.0"></a>
## [v0.6.0](https://github.com/99designs/gqlgen/compare/v0.5.1...v0.6.0) - 2018-10-03
- <a href="https://github.com/99designs/gqlgen/commit/6f486bde038887adf67c3e3766624ef111ea95cf"><tt>6f486bde</tt></a> release v0.6.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7833d0cbf3fd7fb82ddadf8c19a9284554f48250"><tt>7833d0cb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/365">#365</a> from 99designs/dont-guess-imports</summary>

Don't let goimports guess import paths

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/732be3959b402bbd3b864c5f40f475640f1334c5"><tt>732be395</tt></a> Don't let goimports guess import paths

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/926eb9d814747bf3726313d397a31cd7dbddddd1"><tt>926eb9d8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/364">#364</a> from 99designs/query-cache-test</summary>

Add a stress test on query cache

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bab70df5bf66675b16bf1945ce902658a4fdaed2"><tt>bab70df5</tt></a> Add a stress test on query cache

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8448176179aa4492d6cb2962b6155bdeaae2774a"><tt>84481761</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/362">#362</a> from 99designs/fix-error-docs</summary>

fix error docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/23b58f6d5c03e69e05eb0d862be9286912a70151"><tt>23b58f6d</tt></a> fix error docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8f0ef777fbb502ed453fa64861a8a0ca59fcacef"><tt>8f0ef777</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/361">#361</a> from 99designs/revert-360-revert-335-typed-interfaces</summary>

Revert "Revert "Generate typed interfaces for gql interfaces & unions""

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/77257d1e593f35e6b639b3120b8ee6ba4dc7d4a5"><tt>77257d1e</tt></a> Revert "Revert "Generate typed interfaces for gql interfaces & unions""

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1cae19bb114a3bee95ebac0e2e0e47ecdb59ec46"><tt>1cae19bb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/359">#359</a> from 99designs/fix-null-arg-error</summary>

Fix Issue With Argument Pointer Type

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ee862717e87b39a46d01efc031f272ac26fd7b0b"><tt>ee862717</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/360">#360</a> from 99designs/revert-335-typed-interfaces</summary>

Revert "Generate typed interfaces for gql interfaces & unions"

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/02658647f2de7ad601c8eee43417d322b4060ccc"><tt>02658647</tt></a> Revert "Generate typed interfaces for gql interfaces & unions"

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bc35d730cf4a22df34faafe066c77f891b750b9d"><tt>bc35d730</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/335">#335</a> from 99designs/typed-interfaces</summary>

Generate typed interfaces for gql interfaces & unions

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/48724dea899491c1aa75b825047d9c1ef66029e8"><tt>48724dea</tt></a> Removed redundant file

- <a href="https://github.com/99designs/gqlgen/commit/2432ab3cfcc043e31537ac21b550a6b3faf5bfcc"><tt>2432ab3c</tt></a> Fix other tests with pointer change

- <a href="https://github.com/99designs/gqlgen/commit/20add1267ff738b7a6f976de81afab53b22e50da"><tt>20add126</tt></a> Fix test case

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f5c034019889ffb19c28a0c10919c209411def54"><tt>f5c03401</tt></a> Do not strip ptr for args with defaults</summary>

This fails if a client still sends a null value.  If an arg is nullable
but has a default, then null is still a valid value to send through.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0c399270a5c29e617d0cfc147d1809325ad8b8cc"><tt>0c399270</tt></a> Add test case

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b836a976a12a3ea70b0e5b1767b6aceefb8a9fa6"><tt>b836a976</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/358">#358</a> from 99designs/fix-embedded-pointer</summary>

Fix Embedded Pointer

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d3e27553742d7559996fb9cfb310e87c5774fea4"><tt>d3e27553</tt></a> Bump gqlparser to latest master

- <a href="https://github.com/99designs/gqlgen/commit/b8af0c811747c48190126e2d2b4006e718362756"><tt>b8af0c81</tt></a> Use types.Implements to check if an interface implementor accepts value recievers

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2ab05daf9c7864073fcbf8eeb03328223dc66df2"><tt>2ab05daf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/353">#353</a> from 99designs/resolver-ctx-parenting</summary>

Parent middleware generated contexts

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/faf0416bf719ea00f40a96a061567807cef16827"><tt>faf0416b</tt></a> Parent resolver generated contexts

- <a href="https://github.com/99designs/gqlgen/commit/caa474c6ac53f3408ebd24ac3d285247bf6c6f8f"><tt>caa474c6</tt></a> Check for embedded pointer when finding field on struct

- <a href="https://github.com/99designs/gqlgen/commit/f302b4082be7f143304f3e2c39419b2297efef03"><tt>f302b408</tt></a> Added reproduce test case

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/14cf46bc15514a35ac23f8f8b980203cb7bb31da"><tt>14cf46bc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/348">#348</a> from gissleh/feat-websocket-initpayload</summary>

Added parsing of the websocket init message payload

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3147d914a6c2bc0f297a2c8bbe4eaf2c64be0552"><tt>3147d914</tt></a> Updated example in docs to use handler.GetInitPayload instead of graphql.GetInitPayload

- <a href="https://github.com/99designs/gqlgen/commit/32f0b843d8d0e834993caa8b31279862793e137f"><tt>32f0b843</tt></a> Moved InitPayload from graphql to handler package, updated test to import it from there.

- <a href="https://github.com/99designs/gqlgen/commit/01923de635c3c33cbbfd571150973c095f8806a8"><tt>01923de6</tt></a> Moved initPayload to wsConnection member, changed wsConnection.init to return false on invalid payload

- <a href="https://github.com/99designs/gqlgen/commit/25268ef991d11155af3c2abc7294bdc53698ec53"><tt>25268ef9</tt></a> Added information about it under recipes/authentication doc

- <a href="https://github.com/99designs/gqlgen/commit/575f28e0305e991ca7b2c8e60fa793efc1260f0a"><tt>575f28e0</tt></a> Fixed graphql.GetInitPayload panic if payload is nil.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/380828fa768a570b4209762a9cb5ff7236a29bfb"><tt>380828fa</tt></a> Added parsing of the websocket init message payload, and making it available via the context passed to resolvers.</summary>

* Added GetInitPayload(ctx) function to graphql
* Added WithInitPayload(ctx) function to graphql
* Added WebsocketWithPayload method to client.Client (Websocket calls it with a nil payload for backwards compability)
* Added tests for these changes in codegen/testserver/generated_test

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2bd1cc2e669c685d41cacf82c8b730c04c44fef5"><tt>2bd1cc2e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/334">#334</a> from 99designs/support-response-extensions</summary>

Support Extensions in Response

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8fdf4fbbcd24cd78942a3c3f7f3533c59ee5273f"><tt>8fdf4fbb</tt></a> Add test case for extension response

- <a href="https://github.com/99designs/gqlgen/commit/60196b87614965d3926c0a1172974ed8fedbdf4e"><tt>60196b87</tt></a> Add extensions to response struct

- <a href="https://github.com/99designs/gqlgen/commit/cbde0ea97359831026ecc22d2e69adc2c3cd22ad"><tt>cbde0ea9</tt></a> Generate typed interfaces for gql interfaces & unions

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.5.1"></a>
## [v0.5.1](https://github.com/99designs/gqlgen/compare/v0.5.0...v0.5.1) - 2018-09-13
- <a href="https://github.com/99designs/gqlgen/commit/636435b68700211441303f1a5ed92f3768ba5774"><tt>636435b6</tt></a> release v0.5.1

- <a href="https://github.com/99designs/gqlgen/commit/bfb48f2f833c6ab7f2981035b61efdf773dcddba"><tt>bfb48f2f</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/869215a7e69f227f4869f901d319267d1061289d"><tt>869215a7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/339">#339</a> from 99designs/fix-subscription-goroutine-leak</summary>

Fix gouroutine leak when using subscriptions

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/535dd24bf4986186af9fbac5f5965e853fcbdb4f"><tt>535dd24b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/338">#338</a> from codyleyhan/cl/docs</summary>

Adds docs for how resolvers are bound

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/baa99fc58711afed393a3952296a4a9bc754494c"><tt>baa99fc5</tt></a> cleaned up resolver doc

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/647fbbc95b4b3926bedf1ad84e10bc38f050bc68"><tt>647fbbc9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/340">#340</a> from chris-ramon/patch-1</summary>

README.md: Updates `graphql-go/graphql` features.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/729e09c8add3eb78480bc1e55140a5ca3bf3e426"><tt>729e09c8</tt></a> README.md: Updates `graphql-go/graphql` features.</summary>

- Subscription support: https://github.com/graphql-go/graphql/issues/49#issuecomment-404909227
- Concurrency support: https://github.com/graphql-go/graphql/issues/389
- Dataloading support: https://github.com/graphql-go/graphql/pull/388

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/229a81be8b3ad6bcbd3974f2cc078d366ebade7c"><tt>229a81be</tt></a> Fix gouroutine leak when using subscriptions

- <a href="https://github.com/99designs/gqlgen/commit/c15a70ffbb19d8875504f8fde90bb3ff4c5ddd7c"><tt>c15a70ff</tt></a> Adds docs for how resolvers are bound

- <a href="https://github.com/99designs/gqlgen/commit/35c15c940d3b83909551818ed5dc2dd5cd926c6a"><tt>35c15c94</tt></a> Add link to talk by Christopher Biscardi

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/72edf98a67f4efbe44cf03e4922dcbfd0a1bf91a"><tt>72edf98a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/331">#331</a> from edsrzf/arg-refactor</summary>

Refactor arg codegen

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/31505ff44b7fac9941292bb0152c53921e0fcf4a"><tt>31505ff4</tt></a> Use arg function for generated Complexity method

- <a href="https://github.com/99designs/gqlgen/commit/ebdbeba01cb5c62bb3583bcf1171e7f64d5dea1e"><tt>ebdbeba0</tt></a> Just realized "if not" is allow in templates

- <a href="https://github.com/99designs/gqlgen/commit/861a805c7ae8128d2aa61ec1b9ba72bb28cec024"><tt>861a805c</tt></a> Regenerate code

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/639727b644f35584a4bf00a7e3f2331bffdf08df"><tt>639727b6</tt></a> Refactor arg codegen</summary>

Now a function is generated for each field and directive that has
arguments. This function can be used by both field methods as well as
the `Complexity` method.

The `args.gotpl` template now generates the code for this function, so
its purpose is a little different than it used to be.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8026e63b05391e9313a6450207e49fb03c8418f9"><tt>8026e63b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/330">#330</a> from edsrzf/string-compare</summary>

Use built-in less than operator instead of strings.Compare

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c770b4e75b1206130441c8d2a06e05ea44b4715a"><tt>c770b4e7</tt></a> Use built-in less than operator instead of strings.Compare

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.5.0"></a>
## [v0.5.0](https://github.com/99designs/gqlgen/compare/v0.4.4...v0.5.0) - 2018-08-31
- <a href="https://github.com/99designs/gqlgen/commit/5bc4665fab378aa7fe6b81bef968ed608aad1477"><tt>5bc4665f</tt></a> release v0.5.0

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b48c6b92dcdfc7a797c34250478a2b1d1dc486c8"><tt>b48c6b92</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/326">#326</a> from 99designs/version</summary>

Add version const

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/14587a5f44f051fdd733fe80102194edd368d84f"><tt>14587a5f</tt></a> Add version const

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7d44dd6bfe71faf85b6df5e651157210d573b6cd"><tt>7d44dd6b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/315">#315</a> from edsrzf/query-complexity</summary>

Query complexity calculation and limits

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2ab857eec3baef22867cef3a898eb24b0eb65d14"><tt>2ab857ee</tt></a> Merge branch 'master' into query-complexity

- <a href="https://github.com/99designs/gqlgen/commit/6e408d5d64d45e3797fd5fd5b1bf5acc7c50a094"><tt>6e408d5d</tt></a> Interfaces take max complexity of implementors

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d08b9c4a36a8b7c3a22ec1f328a6bef9c95ec1e8"><tt>d08b9c4a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/325">#325</a> from edsrzf/no-get-mutations</summary>

Only allow query operations on GET requests

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/82a28b5735f59524aa3f513f9c893a4c48b6d104"><tt>82a28b57</tt></a> Only allow query operations on GET requests (closes <a href="https://github.com/99designs/gqlgen/issues/317"> #317</a>)</summary>

This mitigates the risk of CSRF attacks.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/239b1d2277632767defaa6d72c82a765d7e87ff1"><tt>239b1d22</tt></a> Don't emit complexity fields for reserved objects

- <a href="https://github.com/99designs/gqlgen/commit/8da5d61b045ee4ea230ecac6706c8857d8f9081d"><tt>8da5d61b</tt></a> Generate complexity for all fields. Fix bugs. Re-generate examples.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/40943c6d92d2158db473aac699277dc6a95b95bb"><tt>40943c6d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/322">#322</a> from 99designs/drop-old-flags</summary>

Drop old cli flags

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8c17eea9ed7fb449e84f89c3f295b119a636c80c"><tt>8c17eea9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/320">#320</a> from andrioid/master</summary>

Description added to generated Model code

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/988b367a542f44fac6e01148503b2e2e7e13fd5d"><tt>988b367a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/316">#316</a> from 99designs/feat-concurrent-each-element</summary>

use goroutine about processing each array elements

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e5265ac2842f5ee39d266e8814d1209e6d4c9625"><tt>e5265ac2</tt></a> Fix complexity template bug

- <a href="https://github.com/99designs/gqlgen/commit/7c0400454230b3bb01d49ba310b7047b064cdc65"><tt>7c040045</tt></a> now with field values

- <a href="https://github.com/99designs/gqlgen/commit/08ab33bedf493581c932313d6cc14e7b5722faf0"><tt>08ab33be</tt></a> starting to look better

- <a href="https://github.com/99designs/gqlgen/commit/e834f6b90453aaf97d9442c238ad2ff1676463ba"><tt>e834f6b9</tt></a> Query complexity docs

- <a href="https://github.com/99designs/gqlgen/commit/a0158a4edd009fbfd6f67a1ec63d1b69c56b719b"><tt>a0158a4e</tt></a> Drop old cli flags

- <a href="https://github.com/99designs/gqlgen/commit/bb78d2faee1270626c670d21699bbde7e682cd93"><tt>bb78d2fa</tt></a> go generate ./..

- <a href="https://github.com/99designs/gqlgen/commit/2488e1b3c0b30b36cc0a0289a4813d21ff672824"><tt>2488e1b3</tt></a> Merge branch 'master' of https://github.com/99designs/gqlgen

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f6a733aea7c01e71ce43efecd1622b61eb8b537c"><tt>f6a733ae</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/308">#308</a> from codyleyhan/tags</summary>

Finds fields by configurable struct tag

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f7aeb88adbe30d19b2c37f8ed41e0322d27f5ef4"><tt>f7aeb88a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/321">#321</a> from 99designs/remove-typemap</summary>

Remove support for the old json typemap

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d63449b91ae1b95c1faa0dbd1100928e2d1b8641"><tt>d63449b9</tt></a> Remove support for the old json typemap

- <a href="https://github.com/99designs/gqlgen/commit/fce4c722a818665a1b0277693a5123b6f166f4a8"><tt>fce4c722</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/8c3aed7d199521533d59e302a9dff4d2cd643aff"><tt>8c3aed7d</tt></a> Merge branch 'master' into query-complexity

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cecd84c698b8ce171e0cd9604215405de248e765"><tt>cecd84c6</tt></a> Add complexity package tests</summary>

Also some small behavior fixes to complexity calculations.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/002ea4761fff26f723d0bbf0119b61b2f5c4f816"><tt>002ea476</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/318">#318</a> from edsrzf/query-cache</summary>

Add query cache

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fcd700b6f9613c01381f5fa78d03700c53b05343"><tt>fcd700b6</tt></a> Panic on lru cache creation error

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/78c570790aa4010613756e9d66f632dd500de091"><tt>78c57079</tt></a> Add query cache</summary>

This commit adds a query cache with a configurable maximum size.
Past this size, queries are evicted from the cache on an LRU basis.

The default cache size is 1000, chosen fairly arbitrarily. If the size
is configured with a non-positive value, then the cache is disabled.

Also ran `dep ensure` to add the new dependency to `Gopkg.lock`.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/076f9eac313b2d2460f4f735a8e59283a2862950"><tt>076f9eac</tt></a> removed dirt

- <a href="https://github.com/99designs/gqlgen/commit/6ae82383f52b2a6c479f982d04215ddfda28f806"><tt>6ae82383</tt></a> trying to get description with generated models

- <a href="https://github.com/99designs/gqlgen/commit/7d6f8ed4b4d5e64eb98dca853aa58fef9aee8784"><tt>7d6f8ed4</tt></a> fixes case where embeded structs would cause no field to be found

- <a href="https://github.com/99designs/gqlgen/commit/02873495e1f85ae83e38cf79e52a2122a845986f"><tt>02873495</tt></a> use goroutine about processing each array elements

- <a href="https://github.com/99designs/gqlgen/commit/40f904a6a3e07145532db7ef08b09d8ec221cbd9"><tt>40f904a6</tt></a> Merge branch 'master' of github.com:99designs/gqlgen into tags

- <a href="https://github.com/99designs/gqlgen/commit/56768d6ba53088c47390944813b0c13cb97e4ae4"><tt>56768d6b</tt></a> adds tests for findField

- <a href="https://github.com/99designs/gqlgen/commit/556b93ac9f76f72cedb4be5189a1c014dea4da04"><tt>556b93ac</tt></a> Run go generate ./...

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2dcb2dd8c3aea0d993b19c93660b4579e404c53f"><tt>2dcb2dd8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/314">#314</a> from 99designs/directive-obj</summary>

Add obj to Directives

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0e2aaa9ec6e6be2cafe04b878f42444275056a41"><tt>0e2aaa9e</tt></a> Merge branch 'master' of github.com:99designs/gqlgen into tags

- <a href="https://github.com/99designs/gqlgen/commit/7cfd9772cb10bb4cb75cef188c0c4635aecb8663"><tt>7cfd9772</tt></a> fixes field selection priority

- <a href="https://github.com/99designs/gqlgen/commit/238a7e2fe255b6e9bd2106f3d0663ef7909fc62e"><tt>238a7e2f</tt></a> Add complexity support to codegen, handler

- <a href="https://github.com/99designs/gqlgen/commit/95ed529b11a57ad717b92efed2767031e7c55d9f"><tt>95ed529b</tt></a> New complexity package

- <a href="https://github.com/99designs/gqlgen/commit/1fda3edefc6e9a028903fafeca7f59bedb796e6f"><tt>1fda3ede</tt></a> Add obj to Directives

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9b24710218507d6420c9710f92ab33d59594584e"><tt>9b247102</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/301">#301</a> from 99designs/feat-directive-parent</summary>

add Result field to ResolverContext

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9ec385d1553ad4af3db23109d708a17848132db4"><tt>9ec385d1</tt></a> Merge branch 'tags' of github.com:codyleyhan/gqlgen into tags

- <a href="https://github.com/99designs/gqlgen/commit/c5849929105209bcdaba2f86ab65cb8b21d6190d"><tt>c5849929</tt></a> adds binding by passed tag

- <a href="https://github.com/99designs/gqlgen/commit/6ef2035b06a98435f9f9fd5b7d8a67b86c7da51d"><tt>6ef2035b</tt></a> refactor set Result timing

- <a href="https://github.com/99designs/gqlgen/commit/568a72e9edde1d564f64e1e22267b76670d12853"><tt>568a72e9</tt></a> add some refactor

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/50588a8af8eedc8c85f388cf17a83ec5077bb39e"><tt>50588a8a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/299">#299</a> from 99designs/test-init-on-windows</summary>

Test gqlgen init on windows

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9148adfc5d3d88e341408acb40b4ee910a7d7a03"><tt>9148adfc</tt></a> Test gqlgen init on windows

- <a href="https://github.com/99designs/gqlgen/commit/c7fd841666d8fbd3496a7da63abb9c9ced3f1c61"><tt>c7fd8416</tt></a> Merge branch 'master' into feat-directive-parent

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3f8a601b5b1d129bdccf79aa72787897701a0027"><tt>3f8a601b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/312">#312</a> from 99designs/validate-gopath</summary>

Validate gopath when running gqlgen

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/77e6955279ca8c844bd1bd5541f4fd7f793164cd"><tt>77e69552</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/310">#310</a> from 99designs/sitemap-404s</summary>

Remove 404s from sitemap

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0b6cedfbab61fd9fbe93c5ac36ce25fff48dd4d1"><tt>0b6cedfb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/311">#311</a> from jekaspekas/fix-mapstructure-err</summary>

fix mapstructure unit test error

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b07736ef8cfd03ba2a649c70cf8bfa3667102ecc"><tt>b07736ef</tt></a> Validate gopath when running gqlgen

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b082227dc3f43337ea17c0a1944e0e7397c31e0f"><tt>b082227d</tt></a> fix mapstructure unit test error</summary>

fix unit test error "mapstructure: result must be a pointer". It appears instead of resolver returned error.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/25b12cb600ae04866a089831c9fb4c01d2e53ab4"><tt>25b12cb6</tt></a> Remove 404s from sitemap

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4a6f505d968843836998a154d06fe8f46e7b598c"><tt>4a6f505d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/309">#309</a> from 99designs/pr-template</summary>

Add a PR template

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/64f3518edb7448bceaa4357a70f2e2a65fbc4c58"><tt>64f3518e</tt></a> run generate

- <a href="https://github.com/99designs/gqlgen/commit/a81147dfed21e1725e5489ddfebfab6b4ea1bd7e"><tt>a81147df</tt></a> Add a PR template

- <a href="https://github.com/99designs/gqlgen/commit/15d8d4ad219ee7121206c7edb332ed56db018795"><tt>15d8d4ad</tt></a> Merge branch 'introspection-directive-args' into HEAD

- <a href="https://github.com/99designs/gqlgen/commit/12efa2d5d8c01a2359721b04a27a58c353567e34"><tt>12efa2d5</tt></a> add tests

- <a href="https://github.com/99designs/gqlgen/commit/95b6f323a880d9034cacf048f0db0c783ca772e9"><tt>95b6f323</tt></a> finds fields by json struct tag

- <a href="https://github.com/99designs/gqlgen/commit/07ee49f3162553b41d45cee11ed0b96ecfe5d745"><tt>07ee49f3</tt></a> Added args to introspection scheme directives.

- <a href="https://github.com/99designs/gqlgen/commit/e57464fef03faff2664fa249de8f9c0e821ed910"><tt>e57464fe</tt></a> refactor ResolverContext#indicies and suppress lint error

- <a href="https://github.com/99designs/gqlgen/commit/09e4bf8c481a8fbdc2ed34bad80051ded1a2023e"><tt>09e4bf8c</tt></a> add Result field instead of ParentObject field

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b8695fb5223ba49fce18aad984af831b710c8b60"><tt>b8695fb5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/304">#304</a> from 99designs/newline-for-init-response</summary>

Put newline at end of `gqlgen init` output

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fabc6f8ff011cf969e318032d3eabdc934632292"><tt>fabc6f8f</tt></a> Merge branch 'master' into feat-directive-parent

- <a href="https://github.com/99designs/gqlgen/commit/e53d224e2ad3137cf93deb5baae638197a3e73f6"><tt>e53d224e</tt></a> Merge branch 'master' into feat-directive-parent

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/de750645ac12568d8909b52f539d625c6a4ee62c"><tt>de750645</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/298">#298</a> from 99designs/handle-response-nulls</summary>

Nulls in required fields should cause errors and bubble

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c855272921a9d21a9eb321276b36a9ac65c4593e"><tt>c8552729</tt></a> Put newline at end of gqlgen init output

- <a href="https://github.com/99designs/gqlgen/commit/072363c777fa35f726d5c6f0626720dab296bc3b"><tt>072363c7</tt></a> add ParentObject field to ResolverContext

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e15d78906a5ab4d528c4ca4b13aaeae2dd944a14"><tt>e15d7890</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/300">#300</a> from 99designs/fix-starwars-connection-example</summary>

fix connection example

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d6acec162547ee261961de22da28cd219df4a24d"><tt>d6acec16</tt></a> fix connection example

- <a href="https://github.com/99designs/gqlgen/commit/7d1cdacabfb6dcea8367468d97ec29be02164a2c"><tt>7d1cdaca</tt></a> Nulls in required fields should cause errors and bubble

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2c4e6cbf05e6920d119790b6f9d6262b87a56e3b"><tt>2c4e6cbf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/294">#294</a> from 99designs/simplfy-concurrent-resolvers</summary>

Simplfy concurrent resolver logic

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7926c688bc50c55b8a2e5c308c3c60857b8fef52"><tt>7926c688</tt></a> Simplfy concurrent resolver logic

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="v0.4.4"></a>
## [v0.4.4](https://github.com/99designs/gqlgen/compare/0.4.3...v0.4.4) - 2018-08-21
- <a href="https://github.com/99designs/gqlgen/commit/6f6622c6b78098660f03d38fb8f0d459d428bdbe"><tt>6f6622c6</tt></a> Bump gqlparser to latest version

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/72659af418c34428b706f6dfc100a678540c8acd"><tt>72659af4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/297">#297</a> from 99designs/fix-dep-pruning</summary>

Explicitly import ambient imports so dep doesn't prune them

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cac3c729ad5fbb5d0e4de0b156c7a0f6f5453b24"><tt>cac3c729</tt></a> Explicitly import ambient imports so dep doesn't prune them

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e6af26e097a8046ad2463c3c580b4777ad54f848"><tt>e6af26e0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/296">#296</a> from heww/master</summary>

sort directives by name when gen

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fd09cd9931347475cfc83f67685ff6ef7e815f6b"><tt>fd09cd99</tt></a> sort directives by name when gen

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/719172670d1c5f80e819bf1376edb8e1f9ed59f3"><tt>71917267</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/292">#292</a> from m4ppi/fix-doc</summary>

Fix broken links in docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/05c73d9f5eac7c84319e894febf9b74b8ae76336"><tt>05c73d9f</tt></a> Fix broken links in docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5a0b56aa613bf7604688f2c5f8f10ec586aec835"><tt>5a0b56aa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/285">#285</a> from 99designs/fix-force-type</summary>

Stop force resolver from picking up types from matching fields

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/31478cf4f74564682f6b0160867661c25b0bbe78"><tt>31478cf4</tt></a> Stop force resolver from picking up types from matching fields

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ebdcf7401de7853e33940a3be357cd7b19b543be"><tt>ebdcf740</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/283">#283</a> from 99designs/speed-up-tests</summary>

Speed up tests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/36e84073852d13ee782733030c457bb68aab8a03"><tt>36e84073</tt></a> Speed up tests

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.4.3"></a>
## [0.4.3](https://github.com/99designs/gqlgen/compare/0.4.2...0.4.3) - 2018-08-10
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3575c289486fce174c941b63749b9bbb88c3ca90"><tt>3575c289</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/281">#281</a> from 99designs/introspection-default-args</summary>

Fix missing default args on types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b808253f02667601d6162acc2a941a504d5a95c2"><tt>b808253f</tt></a> Fix missing default args on types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bf235296103837d477e5d05b062c20f399a51553"><tt>bf235296</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/282">#282</a> from 99designs/flakey-tests</summary>

Remove sleeps in tests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e9c68f08502011d73b20a68dd74b9ed103f9ebe7"><tt>e9c68f08</tt></a> make appveyor less flakey

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.4.2"></a>
## [0.4.2](https://github.com/99designs/gqlgen/compare/0.4.1...0.4.2) - 2018-08-10
- <a href="https://github.com/99designs/gqlgen/commit/06b00d459e44d7a7e29094992d56697fcf8b0f2b"><tt>06b00d45</tt></a> Update README.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5c379a338ee2a7d46da0f0cbab6427d00aa93fc3"><tt>5c379a33</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/279">#279</a> from 99designs/integration-tests</summary>

Integration tests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7f20bdef2615f55f3ea90cd429cd2664ee6e4208"><tt>7f20bdef</tt></a> disable tty for jest

- <a href="https://github.com/99designs/gqlgen/commit/bb0a89a0fd94d1b1a9cf456c6094ed73e701c61f"><tt>bb0a89a0</tt></a> exclude generated code from tests

- <a href="https://github.com/99designs/gqlgen/commit/c2bcff795b4c2d0e0730d5bc84c7a36addd26571"><tt>c2bcff79</tt></a> regenerate

- <a href="https://github.com/99designs/gqlgen/commit/45e22cb1f117b1b277c23e2483cb12f041528e91"><tt>45e22cb1</tt></a> Add introspection schema check

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/53109cd081a7e3f6a0304a0f205eac7fa9cd6b03"><tt>53109cd0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/270">#270</a> from 99designs/feat-handlers</summary>

stop pickup "github.com/vektah/gqlgen/handler" from GOPATH

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ae82b94af59cf825f2c45cc4d7453b0cd136f867"><tt>ae82b94a</tt></a> convert existing tests to jest

- <a href="https://github.com/99designs/gqlgen/commit/f04820b1a64ebff617ce6a0adc98900f7ba521e7"><tt>f04820b1</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/88730e2cce1358a3a420ca9e18660a36f417c8b7"><tt>88730e2c</tt></a> Convert test directory into integration test server

- <a href="https://github.com/99designs/gqlgen/commit/f372b1c920835b873005e33e77e79733267ec93f"><tt>f372b1c9</tt></a> Use docker in docker for the existing testsuite

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0eb08ab9545252332e683fbd912c39dbd9dbc821"><tt>0eb08ab9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/274">#274</a> from 99designs/fix-variable-validation-data</summary>

Prevent executing queries on variable validation failures

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/47a7ac35e34664992754d55d468b2ac09628475a"><tt>47a7ac35</tt></a> Prevent executing queries on variable validation failures

- <a href="https://github.com/99designs/gqlgen/commit/e6e323d02785d0a25f0ccbb88926bc03f6df8a47"><tt>e6e323d0</tt></a> stop pickup "github.com/vektah/gqlgen/handler" from GOPATH

- <a href="https://github.com/99designs/gqlgen/commit/e6005f6b9205ea489b453614858488d46eb48672"><tt>e6005f6b</tt></a> fix mobile nav

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5cdbc9751f61597d92519e4b406674b6a53f6650"><tt>5cdbc975</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/267">#267</a> from 99designs/authentication-docs</summary>

Authentication docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1871c4ce2c38a4d4b27191a08094dfbda626f17c"><tt>1871c4ce</tt></a> Add bold variant of Roboto to docs

- <a href="https://github.com/99designs/gqlgen/commit/fc9fba099d58d1782fabd6a9da7ae213301b1824"><tt>fc9fba09</tt></a> Some minor edits to authentication docs

- <a href="https://github.com/99designs/gqlgen/commit/d151ec8d9cfbb279343104ffda312dec939a402e"><tt>d151ec8d</tt></a> Add docs on user authentication

- <a href="https://github.com/99designs/gqlgen/commit/8db3c143559c65ff3317e17b3b6307afa90c02cb"><tt>8db3c143</tt></a> Add structure to menu

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c57619e0879f25bc749940b8dba3e75910b4e5eb"><tt>c57619e0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/260">#260</a> from 99designs/init-improvements</summary>

Init Config Improvement

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/336b62ec685d04a065c2163cd728877165b3e2ea"><tt>336b62ec</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/266">#266</a> from 99designs/lint-friendly-decollision</summary>

Make keyword decollision more lint friendly

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2acbc245889a300d00c4943983afc9963b3db912"><tt>2acbc245</tt></a> Make keyword decollision more lint friendly

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f12f08a78c27a8d1736cfe19cec983db0c91cdd1"><tt>f12f08a7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/264">#264</a> from 99designs/docs</summary>

CORS docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a2a7c0e7863242de8b77411f8ff1c4bcbaff41ee"><tt>a2a7c0e7</tt></a> Eliminate font resize popin

- <a href="https://github.com/99designs/gqlgen/commit/8a7ed618ff56c297d6590683e8ed149d5f40b734"><tt>8a7ed618</tt></a> Fix errors docs

- <a href="https://github.com/99designs/gqlgen/commit/96e6aab249b4caf1248cd290fa214474caf3d406"><tt>96e6aab2</tt></a> Add CORS docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0ab1c685eb45f580dc52bfe143f8185f3dc363ef"><tt>0ab1c685</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/263">#263</a> from 99designs/add-logo</summary>

Add logo to doc site

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/6d39f868b2129f90f62c3e7d4ca104bfca7eb6a8"><tt>6d39f868</tt></a> Add logo to doc site

- <a href="https://github.com/99designs/gqlgen/commit/d7241728f83cbb8e524ce2e3d765c022a466f5c5"><tt>d7241728</tt></a> Better error on init if file exists

- <a href="https://github.com/99designs/gqlgen/commit/fb03bad9f1ebf86526bce85091272ed50ba46a68"><tt>fb03bad9</tt></a> Run init even if config is found

- <a href="https://github.com/99designs/gqlgen/commit/52b78793bd1bb1343032aeb59724f15dbe628f41"><tt>52b78793</tt></a> Fix hard-coded server filename in init

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.4.1"></a>
## [0.4.1](https://github.com/99designs/gqlgen/compare/0.4.0...0.4.1) - 2018-08-04
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/42f10ec9122abaac7b9cf03444f35b6c5cb5f53d"><tt>42f10ec9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/255">#255</a> from 99designs/introspection-fixes</summary>

Fix introspection api

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7400221c3a5e8cd8917726e9e92522679c2acfbe"><tt>7400221c</tt></a> Fix introspection api

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b35804bac53cf90fa4179dcff6cf6b3b47126c5e"><tt>b35804ba</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/254">#254</a> from oskanberg/patch-1</summary>

Fix typo in introduction docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/84552437d36e5a4124a1af886f01184e16661b57"><tt>84552437</tt></a> Fix typo in introduction docs

- <a href="https://github.com/99designs/gqlgen/commit/b5a48e3e76630d733860c314a6119bed1f224b67"><tt>b5a48e3e</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/c20bb134fa69173a3708160f01ecb79c4276b096"><tt>c20bb134</tt></a> update badges

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.4.0"></a>
## [0.4.0](https://github.com/99designs/gqlgen/compare/0.3.0...0.4.0) - 2018-08-03
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7b5a3d7473f375bb81bd8efe1a08e69a932e6706"><tt>7b5a3d74</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/247">#247</a> from 99designs/next</summary>

0.4.0 Release

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c0be9c9943982ce21a0ff47655c9f4f99034d489"><tt>c0be9c99</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/251">#251</a> from 99designs/rewrite-imports</summary>

Rewrite import paths

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4361401a903bda5d84220b8cb41d8cef3c11f720"><tt>4361401a</tt></a> Rewrite import paths

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f042328a2c75ea771472390e9f1bc33d7cad75f0"><tt>f042328a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/252">#252</a> from 99designs/move-doc-site</summary>

Move doc site

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/658a24d9dcda158b451f5f21535ce2363eb188f8"><tt>658a24d9</tt></a> Move doc site

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/07b7e6ca88acceb1882789fa180109d2a54331dd"><tt>07b7e6ca</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/248">#248</a> from 99designs/json-usenumber</summary>

use json.Decoder.UseNumber() when unmarshalling vars

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/95fe07fef6e87653242067346d7f3e99c0589e5c"><tt>95fe07fe</tt></a> use json.Decoder.UseNumber() when unmarshalling vars

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c555f54cead11d8885d24eb6f7e11260ac930450"><tt>c555f54c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/245">#245</a> from vektah/new-feature-docs</summary>

New feature docs

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/825840aacdf3f160372add5e714dc6e7e42566db"><tt>825840aa</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/244">#244</a> from vektah/array-coercion</summary>

Add implicit value to array coercion

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/90b4076951ef8b5962ec2276b38434d957ae6c94"><tt>90b40769</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/246">#246</a> from vektah/fix-introspection</summary>

Fix introspection

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ef208c76db36eabcabf94c0576e4f18f194e54c1"><tt>ef208c76</tt></a> add docs for resolver generation

- <a href="https://github.com/99designs/gqlgen/commit/e44d798d6e01c3d317199c0345e1a8db5b1bf865"><tt>e44d798d</tt></a> Add directives docs

- <a href="https://github.com/99designs/gqlgen/commit/62d4c8aa60a8187de695b21fdd40858af2b85b87"><tt>62d4c8aa</tt></a> Ignore __ fields in instrospection

- <a href="https://github.com/99designs/gqlgen/commit/bc204c64a892622db810a4b729603f696dda639e"><tt>bc204c64</tt></a> Update getting started guide

- <a href="https://github.com/99designs/gqlgen/commit/b38c580ab3b9828bb5d91fad941274e03c6a0d15"><tt>b38c580a</tt></a> Return the correct mutation & subscription type

- <a href="https://github.com/99designs/gqlgen/commit/9397920c4abf1cc940eee17e022666bf742a62f5"><tt>9397920c</tt></a> Add field name config docs

- <a href="https://github.com/99designs/gqlgen/commit/d2265f3d0a40525359328a52a3d30467f330baa5"><tt>d2265f3d</tt></a> Add implicit value to array coercion

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/191c8ba020f10a1f000cfe5925b972f00807ab6c"><tt>191c8ba0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/239">#239</a> from vektah/directive-args</summary>

Directive args

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3bef596d022a4c58c30e7f1c73adc3b7dec918d3"><tt>3bef596d</tt></a> regenerate

- <a href="https://github.com/99designs/gqlgen/commit/4f37d17028f85eb6e12009d26bea8604e332f766"><tt>4f37d170</tt></a> Add directive args

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f78a6046f87d4ba6d7dc81421b189fa2e772741a"><tt>f78a6046</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/241">#241</a> from vektah/feat-lintfree</summary>

Make more golint free generated code

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/19b5817589c3eaeaf1cbace84e1318c8af33c14b"><tt>19b58175</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

- <a href="https://github.com/99designs/gqlgen/commit/c3fa1a55981538ea9d5b6c9bef10a1c19880588a"><tt>c3fa1a55</tt></a> Merge branch 'next' into feat-lintfree

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/17bfa2cbd4bfccffe704eefbf4d11007ad193e92"><tt>17bfa2cb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/240">#240</a> from vektah/doc-fonts</summary>

Use fonts from golang styleguide

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/64ef0571cc62cae0c71ddb3e1f1dfe6369e6d6e3"><tt>64ef0571</tt></a> Use fonts from golang styleguide

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6b532383c4176fb9b0b683eddc267a0e15ab7481"><tt>6b532383</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/237">#237</a> from vektah/feat-fieldmapping</summary>

Add model field mapping

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4fb721aeec445aa05c595d4b89edf71a9500ab7d"><tt>4fb721ae</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/bf43ab3ddfddfba0309093ffa561b9c7f590eeb2"><tt>bf43ab3d</tt></a> Merge branch 'next' into feat-fieldmapping

- <a href="https://github.com/99designs/gqlgen/commit/353319caf4905a2c5917660db1f10794c37729fd"><tt>353319ca</tt></a> Refactor GoVarName and GoMethodName to GoFieldName etc...

- <a href="https://github.com/99designs/gqlgen/commit/d7e24664af0c9d143b002598ae6cae686eebb59e"><tt>d7e24664</tt></a> Add method support

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/17bcb322d07c82ac794884b572f3acbe59b8bbc0"><tt>17bcb322</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/236">#236</a> from vektah/generate-handler-on-init</summary>

Generate server on running init

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/600f4675970cbca0b376f204fe75e820906db863"><tt>600f4675</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/238">#238</a> from vektah/variable-validation</summary>

Add missing variable validation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d6a76254d197c14902762ff944b0af32126d7b6f"><tt>d6a76254</tt></a> Add missing variable validation

- <a href="https://github.com/99designs/gqlgen/commit/121e8db49d3e3441d04cc74d59799d497408cf44"><tt>121e8db4</tt></a> Generate server on running init

- <a href="https://github.com/99designs/gqlgen/commit/108bb6b4f73b8d0c627602e748373ab64cfb0826"><tt>108bb6b4</tt></a> Rename govarname to modelField

- <a href="https://github.com/99designs/gqlgen/commit/f7f6f9166ab71b67713276597347d429b4691398"><tt>f7f6f916</tt></a> Make more lint friendly

- <a href="https://github.com/99designs/gqlgen/commit/69eab93811af49085104fb1aca7822a3c62392b4"><tt>69eab938</tt></a> Add model field mapping

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ffee020c33783b5e9e2d7f91eff01800b09d6b29"><tt>ffee020c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/235">#235</a> from vektah/generate-resolver-on-init</summary>

Generate resolver on init

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/df95f0032b47a110c0fb2f3f6e8fb879420b05f4"><tt>df95f003</tt></a> Generate code after init

- <a href="https://github.com/99designs/gqlgen/commit/58831ac11446af5960153aed1a7ae84b88ec1506"><tt>58831ac1</tt></a> Generate resolver if configured

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7031264d468ad298a4a562c26cce0f746e2ea5e2"><tt>7031264d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/229">#229</a> from vektah/fix-init-command</summary>

Fixing init command

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/078bc9853f2cc46ec0cec9bc8f57b8a3a7758724"><tt>078bc985</tt></a> Fixing init command</summary>

The init command always return file already exists if there are no
configFilename specified

This is caused by codegen.LoadDefaultConfig() hiding the loading details
and always return the default config with no error while the init
command code expects it to tell us if config exists in default
locations.

To avoid confusion I have splitted the loading config from default
locations out into its own method so we can handle different cases
better.

Additionally I also moved default config into a method so we always
generating new a config instead of passing it around and potentially
mutating the default config.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/803711e941fa9465ebd4ffa7565cba84412a26f9"><tt>803711e9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/221">#221</a> from vektah/middleware-stack</summary>

Implement FieldMiddleware Stack

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/0ec918bf88145a813c99eafa5502d1de7627e54d"><tt>0ec918bf</tt></a> Switch GoName to Name|ucFirst

- <a href="https://github.com/99designs/gqlgen/commit/5dc104ebe5fe54e6eb32c58807a9335e09f248d2"><tt>5dc104eb</tt></a> Add middleware example for Todo

- <a href="https://github.com/99designs/gqlgen/commit/73a8e3a3386fa46b2f2678693faa3c52b26e09b2"><tt>73a8e3a3</tt></a> Fix some issues with directive middlewares

- <a href="https://github.com/99designs/gqlgen/commit/8416324766817e69d984d6e39e51e68758795b8f"><tt>84163247</tt></a> Regenerate

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0e16f1fcdadac39c8cab520956b4987a1938cb7e"><tt>0e16f1fc</tt></a> Generate FieldMiddleware</summary>

Moves it off of RequestContext and into generated land.  This change
has a basic implementation of how directive middlewares might work.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2748a19b2cd34ba6db06aee47a446b242d8824ff"><tt>2748a19b</tt></a> Require Config object into NewExecutableSchema

- <a href="https://github.com/99designs/gqlgen/commit/09242061c7ce1d6e304f06344914c7bc6788c8b7"><tt>09242061</tt></a> Add Directives to Build

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/69e790c27430454c03cc1f97b5fa425f3df6442d"><tt>69e790c2</tt></a> Add *Field to CollectedField</summary>

We need the Field Definition so that we can run directive middlewares
for this field.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d6813f6d47ce90123739c4c8ce6ab9333623e2e0"><tt>d6813f6d</tt></a> Generarte

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/764c6fda0cd7199a51c7e9d806565bf8e6679799"><tt>764c6fda</tt></a> Refactor ResolverMiddleware to FieldMiddleware</summary>

This will allow us to include DirectiveMiddleware in the same middleware
setup, that will run after Resolver middlewares.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7226e573a85b9ac69e054e7c1a73122a1c4afc7d"><tt>7226e573</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/225">#225</a> from rongfengliang/patch-1</summary>

Update getting-started.md

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/66593ffe5339e707c11870a2afc21e2da29042a4"><tt>66593ffe</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

- <a href="https://github.com/99designs/gqlgen/commit/8714f7fbb97089cc7bd56ee1d00e9d234c226220"><tt>8714f7fb</tt></a> hush metalinter

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0dfb92a79af66024f70a323faef7ddc72bd2b83b"><tt>0dfb92a7</tt></a> Update getting-started.md</summary>

CreateTodo  UserID   input should be UserId not User

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0fa7977f9a778a0ba8b56e7272702d508bd390f9"><tt>0fa7977f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/217">#217</a> from vektah/resolver-middleware-all</summary>

Run Resolver Middleware For All Fields

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7292be78338605c49ef96c55595d53aba243682a"><tt>7292be78</tt></a> Rename CastType to AliasedType</summary>

This field stores a Ref if a type is a builtin that has been aliased. In
most cases if this is set, we want to use this as the type signature
instead of the named type resolved from the schema.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ec928cadf952e43d95092194a3aed42f2174e207"><tt>ec928cad</tt></a> Regenerate examples

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/97f131842addf705f39d56c21cd094cebeca435f"><tt>97f13184</tt></a> Remove comment about ResolverMiddleware</summary>

Not true anymore!

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b512176ccd04fca974685829c93345357a0c9cf1"><tt>b512176c</tt></a> Run resolver middleware for all fields

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f67f8390f31b2a498900aa4e3db817b8da3e704f"><tt>f67f8390</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/218">#218</a> from vektah/remove-old-resolvers</summary>

Remove old resolvers

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1a3e4e9940602bb047924e623e641f6e7c40cff0"><tt>1a3e4e99</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/220">#220</a> from vektah/feat-race</summary>

turn back -race option

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/40989b193c9885f94f6ceee0e66865992865af26"><tt>40989b19</tt></a> turn back -race option

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1ba61fcb262123cefb506c899160292fc1dab542"><tt>1ba61fcb</tt></a> Update test & examples to use new resolver pattern</summary>

* chat
* dataloader
* scalar
* selection
* starwars
* todo

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3870896111fe12e5320880a0f1baf65a1baa3776"><tt>38708961</tt></a> Stop generating two types of resolvers</summary>

In recent refactor we introduced a new pattern of resolvers which is
better structured and more readable. To keep Gqlgen backward compatible
we started generate two styles of resolvers side by side.

It is now time to sunset the old resolver. This commit removes the old
resolver and update the generation code to use the new resolver
directly.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ffe42658fed0fd990c80c1cb58684c79e9c33642"><tt>ffe42658</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/208">#208</a> from vektah/directives-skip-include

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a69071e3a59ae80bef327cbc9f823933d5ec4794"><tt>a69071e3</tt></a> Pass context to CollectFields instead of RequestContext</summary>

Internally it can still get to RequestContext as required.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d02d17ae1d05a2efd3915bff84153feb6497f458"><tt>d02d17ae</tt></a> Add method for generating method name from field

- <a href="https://github.com/99designs/gqlgen/commit/c7ff32086724ca8a1fd8bf437cc1f7417c47a619"><tt>c7ff3208</tt></a> Update gqlparser version to include default resolution

- <a href="https://github.com/99designs/gqlgen/commit/ce17cd9034ff42a6881429ac198afb572e6950f2"><tt>ce17cd90</tt></a> Add default value test case

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cbfae3d31bd8ee298103c65d83d9eef7272d44b6"><tt>cbfae3d3</tt></a> Add skip/include test cases</summary>

Adds a set of test cases for skip and include directives to the todo
example. Also now conforms to spec if both are included.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ea0f821cc6c8e661f4b60f26aeec12ac90bcf3e0"><tt>ea0f821c</tt></a> Add skip/include directive implementation</summary>

This is a snowflake implementation for skip/include directives based on
the graphql-js implementation.  Skip takes precedence here.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ebfde103e0ca294e8ce6ba131419344a1be67048"><tt>ebfde103</tt></a> Pass request context through to CollectFields

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bab7abb21b7328380e9ff8985db81fc523160158"><tt>bab7abb2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/210">#210</a> from vektah/feat-init</summary>

introduce gen & init subcommand

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6ba508f96ab7dd962ce4ac3a43f3bd2572e1bf4c"><tt>6ba508f9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/214">#214</a> from vektah/gqlparser-schema-validation</summary>

Bump gqlparser to get schema validation

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/138b4ceafd405c1d0e288e996d2f5ca74f3e179a"><tt>138b4cea</tt></a> Bump gqlparser to get schema validation

- <a href="https://github.com/99designs/gqlgen/commit/08d7f7d0cdd6bb85ce42f70cbb57e4e16b41824b"><tt>08d7f7d0</tt></a> Merge branch 'next' into feat-init

- <a href="https://github.com/99designs/gqlgen/commit/39f9dbf6d245f5b1b1c942c7afcda31519ca1113"><tt>39f9dbf6</tt></a> fix error from breaking change

- <a href="https://github.com/99designs/gqlgen/commit/41147f6f907d2144b0499ca141ed48fc8b97e9c6"><tt>41147f6f</tt></a> update Gopkg.lock

- <a href="https://github.com/99designs/gqlgen/commit/87d8fbeaa5f2e5363e717e808ba42e0de60af5f7"><tt>87d8fbea</tt></a> remove unused flag

- <a href="https://github.com/99designs/gqlgen/commit/eff49d048962d2804181a93d29654335de7b74f3"><tt>eff49d04</tt></a> support init subcommand

- <a href="https://github.com/99designs/gqlgen/commit/c5810170eb49ef72c55adf3de2789cf1f50348e9"><tt>c5810170</tt></a> introduce cobra library

- <a href="https://github.com/99designs/gqlgen/commit/c3c20f8f609b0e58e3da1dc26d481becd90d8b8a"><tt>c3c20f8f</tt></a> Merge remote-tracking branch 'origin/master' into HEAD

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/90df37f6a80713a4c6faf0e9887036189c29453a"><tt>90df37f6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/205">#205</a> from vektah/forward-credential-to-graphql-endpoint</summary>

Use original credential for query request in playground

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/52343745e6e491a072e4ea51bc624bf2a911c159"><tt>52343745</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/206">#206</a> from vektah/validation-locations</summary>

Update gqlparser for validation locations

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f4d31aa429fe5e5ea20e5e88e114c2d29a551a71"><tt>f4d31aa4</tt></a> Update gqlparser for validation locations

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9d473f8b6585a1919cca3b7224ef0cca9767dae9"><tt>9d473f8b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/203">#203</a> from vektah/99designs-announcement</summary>

Announcement: 99designs is now sponsoring gqlgen

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c2f1570d1a9b0cbe2fdc97f568b4fa4005eab608"><tt>c2f1570d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/204">#204</a> from vektah/gqlparser-prelude</summary>

Use shared prelude

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/004ec6a96a06b8a08815311c8e8170661c91a836"><tt>004ec6a9</tt></a> Add 99designs sponsorship news

- <a href="https://github.com/99designs/gqlgen/commit/548aed142366e4fc6e37527db629d08e7f2903c2"><tt>548aed14</tt></a> Use shared prelude

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/edb3ea4e725657409ff4c96de59b49821b6225e5"><tt>edb3ea4e</tt></a> Use original credential for query request in playg</summary>

Currently the playground doesn't forward any credentials when making
query calls. This can cause problems if your playground requires
credential logins.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f855a89c8a873d1d8523bb3c1d4b47330778da65"><tt>f855a89c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/201">#201</a> from cocorambo/remove-trailing-println</summary>

Remove trailing Println

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c41a6c36cf9adad011de060817bae8cf66f79933"><tt>c41a6c36</tt></a> Remove trailing Println

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2692d3e0aa5a6c5936301a822c5c71320ad96dbc"><tt>2692d3e0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/197">#197</a> from vektah/new-parser</summary>

Integrate gqlparser

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5796d47d3c62f32f128e85c037bf35f31b458a11"><tt>5796d47d</tt></a> Integrate gqlparser

- <a href="https://github.com/99designs/gqlgen/commit/55179a61aeefeb044f023d7fa0342518b90af76e"><tt>55179a61</tt></a> Update badges

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/01a4c67737c161aebd8e76c3a9406140ecdd895f"><tt>01a4c677</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/195">#195</a> from jonstaryuk/master</summary>

Update playground version

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c52f24aff789512ede45cb2cf16a670f769f6f2d"><tt>c52f24af</tt></a> Update playground version to 1.6.2

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.3.0"></a>
## [0.3.0](https://github.com/99designs/gqlgen/compare/0.2.5...0.3.0) - 2018-07-14
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/381b34691fd93829e50ba8821412dc3467ec4821"><tt>381b3469</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/194">#194</a> from vektah/multiline-comments</summary>

Fix multiline comments

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/112d68a680fc6a6d053eabaadff3ba391f2bd1b6"><tt>112d68a6</tt></a> only build master branch

- <a href="https://github.com/99designs/gqlgen/commit/4b3778e32cffe09c745f12b130aae3b08f281902"><tt>4b3778e3</tt></a> Fix multiline comments

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/eb44925cb17032879939a5352a2fbdc930f79320"><tt>eb44925c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/193">#193</a> from vektah/validate-method-returns</summary>

validate method return types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/164acaed8876c96f0b9f726fd4fdc5e59f79aad9"><tt>164acaed</tt></a> validate method return types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f478f816529a9045ad39a90df36629250445a317"><tt>f478f816</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/192">#192</a> from vektah/strict-config</summary>

Strict config

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/a1c02e7b771e7898aff1590326b324c3e373a702"><tt>a1c02e7b</tt></a> Strict config

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/533dcba7f4062ee4c090ea5f3bdceafe29e3bce0"><tt>533dcba7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/191">#191</a> from vektah/nullable-list-elements</summary>

Support nullable list elements

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e0bf6afd14fbba795a0823248a05547c0d4fc520"><tt>e0bf6afd</tt></a> Support nullable list elements

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0780bf2ecce2409042b359ddfb86f34a59417ef4"><tt>0780bf2e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/190">#190</a> from vektah/generated-forced-resolvers</summary>

Allow forcing resolvers on generated types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bf1823cdf0e531aeda296090a5fd249efa8131ec"><tt>bf1823cd</tt></a> Allow forcing resolvers on generated types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/febd0358f4f7373e619597f1199c8e4f71270329"><tt>febd0358</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/186">#186</a> from vektah/error-redux</summary>

Error redux

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/b884239a5170fcff95f9ceb1dbf7027671cedd95"><tt>b884239a</tt></a> clarify error response ordering

- <a href="https://github.com/99designs/gqlgen/commit/58e32bbf601cba00e6ecc40e49067ffacb2a9bf3"><tt>58e32bbf</tt></a> Drop custom graphql error methods

- <a href="https://github.com/99designs/gqlgen/commit/d390f9c649e4cad098e4c20f2019cb1459ca642d"><tt>d390f9c6</tt></a> Errors redux

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.2.5"></a>
## [0.2.5](https://github.com/99designs/gqlgen/compare/0.2.4...0.2.5) - 2018-07-13
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0a9709db44324b53cb37ee23eff512213678362d"><tt>0a9709db</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/188">#188</a> from vektah/fix-windows-gopath</summary>

Fix windows gopath issue

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ea4f26c69e8b86e57ed46db71fec96e4dc6742be"><tt>ea4f26c6</tt></a> more fixes

- <a href="https://github.com/99designs/gqlgen/commit/1066953dfc6f39dfc4324190064deaf8793eaec3"><tt>1066953d</tt></a> Appveyor config

- <a href="https://github.com/99designs/gqlgen/commit/f08d8b61c5c68e1e9d9b5cd411ba465e986c63ae"><tt>f08d8b61</tt></a> Fix windows gopath issue

- <a href="https://github.com/99designs/gqlgen/commit/9ade6b7a62b2ad75b96ace7171c6c500fdccb137"><tt>9ade6b7a</tt></a> Update gettingstarted to use new resolvers

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.2.4"></a>
## [0.2.4](https://github.com/99designs/gqlgen/compare/0.2.3...0.2.4) - 2018-07-10
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ac9e5a66f8790951063927b9c971f99cdaae7a2f"><tt>ac9e5a66</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/180">#180</a> from vektah/import-alias-before-finalize</summary>

Fix a bug custom scalar marshallers in external packages

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/160ebab52f7ed9d95944dadaa07050a08e81ac36"><tt>160ebab5</tt></a> Fix a bug custom scalar marshallers in external packages

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/43212c04c9ffa65fa3c44d3179f8b123addae767"><tt>43212c04</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/179">#179</a> from vektah/models-config-error</summary>

Improve Output Filename and Package Handling

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/936bc76e4fc8dd23110eefde1c6e054feaaf50c4"><tt>936bc76e</tt></a> Better handling of generated package name

- <a href="https://github.com/99designs/gqlgen/commit/5d3c8ed2a58189dbe8c09b8854446e503e463f37"><tt>5d3c8ed2</tt></a> Inline ImportPath strings

- <a href="https://github.com/99designs/gqlgen/commit/fc43a92ad6b1370cbb319fa37ca8cb12f9d59226"><tt>fc43a92a</tt></a> Check that exec and model filenames end in *.go

- <a href="https://github.com/99designs/gqlgen/commit/6d38f77d0841210a48b64a62573a6460fff93a62"><tt>6d38f77d</tt></a> Handle package name mismatch with dirname

- <a href="https://github.com/99designs/gqlgen/commit/ebf1b2a5688787fa6e3e3474e77080990e96a875"><tt>ebf1b2a5</tt></a> Add error message when specifying path in package name

- <a href="https://github.com/99designs/gqlgen/commit/c8355f48e90da4d8acba0ab5aa124e87db6fbb2d"><tt>c8355f48</tt></a> Check models config for package-only specs

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.2.3"></a>
## [0.2.3](https://github.com/99designs/gqlgen/compare/0.2.2...0.2.3) - 2018-07-08
- <a href="https://github.com/99designs/gqlgen/commit/6391596d0b6d3fb06b412bffbb7e18e3bc1e3044"><tt>6391596d</tt></a> Add some basic docs on the new config file

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a9c3af86ee8cda101e1e6044407db9d447da9f86"><tt>a9c3af86</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/176">#176</a> from vektah/config-search-paths</summary>

Search for config

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/25cfbf082d8a07c4f9f247acb5452631056160e1"><tt>25cfbf08</tt></a> Search for config

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/bff3356bec9b9b3ea70ff157fee5a9fa9421ab2a"><tt>bff3356b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/175">#175</a> from vektah/lint-all-packages</summary>

gometalinter should cover all packages

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/61f3717348e532ea4c072268a40844e3f758a1b1"><tt>61f37173</tt></a> gometalinter should cover all packages

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ce6570448d36a9805c89bf4d071f1943872cc02e"><tt>ce657044</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/173">#173</a> from vvakame/feat-resolver-hint</summary>

add resolver option support to field

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/57b8279e2bb11344280410baa8e9a4c11721955d"><tt>57b8279e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/172">#172</a> from vvakame/feat-newconfig</summary>

switch to .gqlgen.yml

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fcfceefbfc9e91a5d25da702de673036198354f8"><tt>fcfceefb</tt></a> add resolver option support to field

- <a href="https://github.com/99designs/gqlgen/commit/c7ce1cbbbf9082d023d9cfc2e4279c7e077dfc86"><tt>c7ce1cbb</tt></a> update docs

- <a href="https://github.com/99designs/gqlgen/commit/42948153981d2fe84c715b20b396966dd74d5c09"><tt>42948153</tt></a> move to .gqlgen.yml

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/325c45a40b41dec948abd1138cc8f84ae815b285"><tt>325c45a4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/171">#171</a> from vvakame/add-gitignore</summary>

add .idea/ to .gitignore

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/aa4cec9b05374bfa196d29fdd31fc67865373b3e"><tt>aa4cec9b</tt></a> add .idea/ to .gitignore

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.2.2"></a>
## [0.2.2](https://github.com/99designs/gqlgen/compare/0.2.1...0.2.2) - 2018-07-05
- <a href="https://github.com/99designs/gqlgen/commit/f79b6a52ef73871a2f0d2d57b15a77078439c3b1"><tt>f79b6a52</tt></a> cleanup new config

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f0a08617588220df36b0938dd5e0dbe0f2a06538"><tt>f0a08617</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/163">#163</a> from vvakame/feat-types-json</summary>

support .gqlgen.yml

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/faf095fc412ac24799165e3e6217a46a34396cb8"><tt>faf095fc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/166">#166</a> from vektah/validate-at-end</summary>

Validate at end

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/fca1e08e1796afd152cbf8ad94c4de16fa1faebc"><tt>fca1e08e</tt></a> shh errcheck

- <a href="https://github.com/99designs/gqlgen/commit/cc78971ee7f5e83d49d5bb7e0c4bea30fc12e4a8"><tt>cc78971e</tt></a> Dont show compilation errors until after codegen

- <a href="https://github.com/99designs/gqlgen/commit/9f6ff0cf7a9e14380c09a03029959199315a4455"><tt>9f6ff0cf</tt></a> Convert todo example to new resolver syntax

- <a href="https://github.com/99designs/gqlgen/commit/8577ceab9d2715b84d76aa38a6dc2bb20fd95889"><tt>8577ceab</tt></a> address comment

- <a href="https://github.com/99designs/gqlgen/commit/86dcce730bad2f06b48faf7f8ea1f27af668de50"><tt>86dcce73</tt></a> Add format check to -typemap argument

- <a href="https://github.com/99designs/gqlgen/commit/5debbc6acfc7aa0b5ace5c96d775401aef4ad85f"><tt>5debbc6a</tt></a> Implement types.yaml parsing

- <a href="https://github.com/99designs/gqlgen/commit/ecf56003d5b366805930e62f787ec338f57d9543"><tt>ecf56003</tt></a> Refactor types.json parsing

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b16e84295a1e72d27c6b96784c0266058d8716bb"><tt>b16e8429</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/159">#159</a> from vektah/enum-only-generation</summary>

Dont skip model generation if there are enums defined

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3f751a407d6c33f5cc61998983ec02d1c644fc26"><tt>3f751a40</tt></a> Dont skip model generation if there are enums defined

- <a href="https://github.com/99designs/gqlgen/commit/588aeacb5fb32b6b3e4ee818fec784eae2277956"><tt>588aeacb</tt></a> more tutorial fixes

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/dc472965103858af8dd0af0f951f453f88ea3f3e"><tt>dc472965</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/157">#157</a> from johncurley/fix-docs-argument</summary>

Updated mutation to take correct argument

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/88a84f83320c19fb01b6935af8d4fd34652344fc"><tt>88a84f83</tt></a> Updated mutation to take correct argument

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/404f0b0d4844035f971135864bb4b20e98761b22"><tt>404f0b0d</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/151">#151</a> from qdentity/fix-longer-gopath</summary>

Fix bug with multiple GOPATH full package name resolving

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f66e2b3b8c7be8915dbf8cdc2a7d57a907ff31f0"><tt>f66e2b3b</tt></a> Fix bug with multiple GOPATH full package name resolving</summary>

This commit fixes the bug where GOPATH values that are longer than the input package name cause 'slice bounds out of range'  errors.

</details></dd></dl>

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.2.1"></a>
## [0.2.1](https://github.com/99designs/gqlgen/compare/0.2.0...0.2.1) - 2018-06-26
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cb87a2cb66f5a64749f6464900b1c12bca47ed67"><tt>cb87a2cb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/147">#147</a> from vektah/import-overhaul</summary>

Improve import handling

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9fa3f0fbdc5691042e6ca21e9574d87715838318"><tt>9fa3f0fb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/134">#134</a> from mastercactapus/small-interfaces</summary>

add lint-friendly small interfaces option for resolvers

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e8c30acdce72d1f39b9308f1b548a89f7a11316c"><tt>e8c30acd</tt></a> fix template error on generated defaults (<a href="https://github.com/99designs/gqlgen/pull/146">#146</a>)</summary>

* fix template error on generated defaults

* go fmt

* add test for default fix

* .

* add key sort for default values

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/769a97e2e6903960df001e5953e5369e86d6432e"><tt>769a97e2</tt></a> fix race in chat example test - t.Parallel() doesn't guarantee parallel execution - moved goroutine so the test can execute independently

- <a href="https://github.com/99designs/gqlgen/commit/5b77e4c22dd929db47cdf23535e14bab5b793a93"><tt>5b77e4c2</tt></a> remove deprecation warning for now

- <a href="https://github.com/99designs/gqlgen/commit/59a5d7520fc936e9432b9ec8484ad9175b3ddc5b"><tt>59a5d752</tt></a> remove trailing S

- <a href="https://github.com/99designs/gqlgen/commit/b04846f6da0cfdb341a38e76ca8069a6cd792ef7"><tt>b04846f6</tt></a> fix time race in scalar test

- <a href="https://github.com/99designs/gqlgen/commit/a80b720fc230172c8c3a7414b2ad06175d94cba4"><tt>a80b720f</tt></a> name updates, deprecation, some code comments

- <a href="https://github.com/99designs/gqlgen/commit/2bbbe0546d12ba343db5eea61aeba896c127a6c4"><tt>2bbbe054</tt></a> Merge branch 'master' into small-interfaces

- <a href="https://github.com/99designs/gqlgen/commit/4ffa2b24d4b4c53b934d8a168d54c50e08dde9b6"><tt>4ffa2b24</tt></a> case insensitive compare to determine self package

- <a href="https://github.com/99designs/gqlgen/commit/c0158f5418b679d1af358ce6f7f9a9d3ecf4fcf0"><tt>c0158f54</tt></a> make sure colliding imports are stable

- <a href="https://github.com/99designs/gqlgen/commit/abf85a104ab1b06c6c181c7dae74d84b3d88628c"><tt>abf85a10</tt></a> get package name from package source

- <a href="https://github.com/99designs/gqlgen/commit/a39c63a5ef9dadec023241f049d511380dbce189"><tt>a39c63a5</tt></a> remove a random json tag from tutorial

- <a href="https://github.com/99designs/gqlgen/commit/f48cbf03b9df1d3bcedf44f8b23fa7f18b6c909a"><tt>f48cbf03</tt></a> tutorial fixes

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0a85d4f2df870106b309ab11c467070de5214200"><tt>0a85d4f2</tt></a> Update generated headers to match convention. (<a href="https://github.com/99designs/gqlgen/pull/139">#139</a>)</summary>

* Update generated.gotpl

* Update models.gotpl

* Update data.go

* update go generate

* revert code changes

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4a6827bdcc106187fb8ee4b70deb19931f5514ee"><tt>4a6827bd</tt></a> Update getting started guide

- <a href="https://github.com/99designs/gqlgen/commit/a21f32731d8f354f29aa11f1f62b5bb2da977ad8"><tt>a21f3273</tt></a> Use recognized `Code generated` header

- <a href="https://github.com/99designs/gqlgen/commit/038c6fd2d499c89dea430255ffaf583510fd6016"><tt>038c6fd2</tt></a> change from `ShortResolver` to `ShortResolvers` - prevents possible collision with an object type named `Short`

- <a href="https://github.com/99designs/gqlgen/commit/0bc592cd070d8d54fd6ceb23032f994918bc6db8"><tt>0bc592cd</tt></a> run go generate

- <a href="https://github.com/99designs/gqlgen/commit/db2cec072a8f617844393c248248f94500f7749a"><tt>db2cec07</tt></a> fix template formatting

- <a href="https://github.com/99designs/gqlgen/commit/59ee1b5cf3789aa98907436ab446373249e88527"><tt>59ee1b5c</tt></a> from probably makes more sense

- <a href="https://github.com/99designs/gqlgen/commit/620f7fb42c4be6c5bc46bf56b41d0f0903adb9f0"><tt>620f7fb4</tt></a> add "short" resolver interface

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<a name="0.2.0"></a>
## [0.2.0](https://github.com/99designs/gqlgen/releases/tag/0.2.0) - 2018-06-21
<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d26ef2a2622e005e6047c924ef83fcbec83ea46c"><tt>d26ef2a2</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/136">#136</a> from tianhai82/master</summary>

fix GOPATH case mismatch issue

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a34b4de4cdf37401d82871970f7696b857cd63ce"><tt>a34b4de4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/137">#137</a> from appleboy/patch-1</summary>

fix example links

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c1cde36c18c84633c4adcc8de9d6f97e48b7ec31"><tt>c1cde36c</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/133">#133</a> from mastercactapus/skip-type-mismatch</summary>

skip struct fields with incompatible types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c1b4574cc5f9c09925f391ebab839d30c63f7f2b"><tt>c1b4574c</tt></a> fix example links

- <a href="https://github.com/99designs/gqlgen/commit/63976d5fd90bf374c9b0553ac34c3e68bec88310"><tt>63976d5f</tt></a> fix GOPATH case mismatch issue

- <a href="https://github.com/99designs/gqlgen/commit/8771065fa6a95ac7cf8c0548f6f1da3e0d23818f"><tt>8771065f</tt></a> skip fields with incompatible types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/40d9a11be10846d2d69c8731c503ddf434b93146"><tt>40d9a11b</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/127">#127</a> from jon-walton/windows-path-slash</summary>

convert windows input path separators to slash

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7db9d122bb7b2145dac61b69a1947c5d211c7623"><tt>7db9d122</tt></a> convert windows input path separators to slash

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a5f7260161f98d29b79c3484f80279ce42990dfc"><tt>a5f72601</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/122">#122</a> from vektah/json-encoding-fixes</summary>

Fix json string encoding

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f207c62c817cfe47e58407639fbf47491a6da3fd"><tt>f207c62c</tt></a> review feedback

- <a href="https://github.com/99designs/gqlgen/commit/578d8415192341986a44ee0e1acf5f623534f5ac"><tt>578d8415</tt></a> Fix json string encoding

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e9b406669f567b4804958606dfd68f641dbce6a3"><tt>e9b40666</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/123">#123</a> from vektah/drop-fk-generation</summary>

BC Break: Stop generating foreign keys in models

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a8419e20b5181f1397761cf17cd4da773c2873c9"><tt>a8419e20</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/124">#124</a> from vektah/fix-backtick-escaping</summary>

Fix backtick escaping

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/47eaff4de24fbdda9c0d78df76182748e47035ab"><tt>47eaff4d</tt></a> Fix backtick escaping

- <a href="https://github.com/99designs/gqlgen/commit/a5c02e6c1d9700b620d6735ed38e4f023f26bbd9"><tt>a5c02e6c</tt></a> BC Break: Stop generating foreign keys in models

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/94d5c89eabde2bef4efd024c23587474d38ddf94"><tt>94d5c89e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/120">#120</a> from andrewmunro/bugfix/fix-panic-on-invalid-array-type</summary>

Fixing panic when non array value is passed to array type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5680ee49b7ff8972d6c8998e3c301675bc6e30d0"><tt>5680ee49</tt></a> Adding dataloader test to confirm no panic on malformed array query

- <a href="https://github.com/99designs/gqlgen/commit/55cc161f6fb9f44d15fe520b1f6bcff12b3a1db6"><tt>55cc161f</tt></a> Fixing panic when non array value is passed to array type

- <a href="https://github.com/99designs/gqlgen/commit/6b3b338d5f9c8b5a80ad4ea1e2e37aa58677ea9d"><tt>6b3b338d</tt></a> Add gitter link to readme

- <a href="https://github.com/99designs/gqlgen/commit/6c823beb069be4fed6f0218ec5ff6a5211968b56"><tt>6c823beb</tt></a> add doc publish script

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a25232d8cce19899b610fcaffbda5cee3d1f4bab"><tt>a25232d8</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/113">#113</a> from mikeifomin/patch-1</summary>

Fix typo in url dataloaden

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3a129c77a73b340bd9e04ecfdbebe65e9779f47f"><tt>3a129c77</tt></a> Fix typo in url dataloaden

- <a href="https://github.com/99designs/gqlgen/commit/e1fd79fed15f60c47471d901c8250ab56aff1c55"><tt>e1fd79fe</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/111">#111</a> from imiskolee/master (closes <a href="https://github.com/99designs/gqlgen/issues/110"> #110</a>)

- <a href="https://github.com/99designs/gqlgen/commit/e38cb497d72a1452c04ed2b82195f6b7cb142038"><tt>e38cb497</tt></a> 1. fix bug: <a href="https://github.com/99designs/gqlgen/pull/110">#110</a>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/3990eacf7d8d99143a69b249ef164787ed00e2ee"><tt>3990eacf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/108">#108</a> from imiskolee/master</summary>

generate json tag to model field  by gql name.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/abb7502af6f8df4726b5d96233a32871668b3787"><tt>abb7502a</tt></a> 1. run go generate

- <a href="https://github.com/99designs/gqlgen/commit/e1f90946d1d81737ef40e2bcf8cecdb770d34f5f"><tt>e1f90946</tt></a> 1. add json tag in models_gen.go 2. use gqlname to model filed json tag.

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/35e0971773350fc64949ed5526bf74d7ea2cd574"><tt>35e09717</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/107">#107</a> from vektah/fix-vendor-normalization</summary>

Fix vendor normalization

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/63ee41996e5466e27ea2711923c546eef41c6183"><tt>63ee4199</tt></a> Fix vendor normalization</summary>

When refering to vendored types in fields a type assertion would fail. This
PR makes sure that both paths are normalized to not include the vendor
directory.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2a437c23b379cdcad7b7f4ca2e14a5c6075123a9"><tt>2a437c23</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/105">#105</a> from vektah/keyword-input-args</summary>

Automatically add a _ suffix to reserved words

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/26ac13ffec364baf5d6db0a3d6bb613c5fba25ea"><tt>26ac13ff</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/104">#104</a> from vektah/new-request-context</summary>

Add a NewRequestContext method

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/309e5c6db1c7af0da9b75b14aeff16885264122d"><tt>309e5c6d</tt></a> Automatically add a _ suffix to reserved words

- <a href="https://github.com/99designs/gqlgen/commit/a2fb14213d99f82edc5ef1a0ae44999c0fa1a707"><tt>a2fb1421</tt></a> Add a NewRequestContext method

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ab6e65bd7d3b9dec26340074435cffec27cfe8d8"><tt>ab6e65bd</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/97">#97</a> from vektah/add-input-defaults</summary>

Default values for input unmarshalers

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1cd80c4a529688a8a713a9f5755c678e14db0e8c"><tt>1cd80c4a</tt></a> Default values for input unmarshalers

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/79c69d15a44a35ee31aa4048a67a3240b3264636"><tt>79c69d15</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/96">#96</a> from vektah/refactor-tests</summary>

Refactor tests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7b1c819850ea5571e0c19530151bd1e0d3e02b32"><tt>7b1c8198</tt></a> Refactor tests

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0c7bdfc682bc065ff44752ac4fbefd455c3fcfeb"><tt>0c7bdfc6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/95">#95</a> from vektah/custom-error-types</summary>

Custom error types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/4bdc1e1f2a9f6442e65c6e0115d576dbc6566147"><tt>4bdc1e1f</tt></a> regenerate

- <a href="https://github.com/99designs/gqlgen/commit/20250f189458d46e133a2cd583f0c503b40d51b2"><tt>20250f18</tt></a> Add fully customizable resolver errors

- <a href="https://github.com/99designs/gqlgen/commit/a0f66c8801e1d76bac4e3c16fbcf39a6403acc66"><tt>a0f66c88</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/8f62d505c2c1d38252f6dda848d072b0aae88456"><tt>8f62d505</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/a1043da696875b205ae3447ca75f58741fb4780e"><tt>a1043da6</tt></a> Add feature comparison table to readme

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/22128e0ed65686e6482cb1b0a27363e67db69733"><tt>22128e0e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/93">#93</a> from vektah/input-type-error-handling</summary>

Input type error handling

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e7539f110a3d90ff9a0bce861e7c024fd91e2a02"><tt>e7539f11</tt></a> Add an error message when using types inside inputs

- <a href="https://github.com/99designs/gqlgen/commit/a780ce694891bd12fafe93eecc0e30cf534480e0"><tt>a780ce69</tt></a> Add a better error message when passing a type into an input

- <a href="https://github.com/99designs/gqlgen/commit/0424f0434ff34db8d4e8fc769e69f7205db05d49"><tt>0424f043</tt></a> Refactor main so tests can execute the generator

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ab3803a6683be2941e5817ccb9f34460f2b39ac9"><tt>ab3803a6</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/89">#89</a> from vektah/opentracing-parent-span</summary>

Add parent opentracing span around root query/mutation/resolvers

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d157ac353535af70d5c235d493281fefc8111d73"><tt>d157ac35</tt></a> Add context to recover func</summary>

This makes the request and resolver contexts available during panic
so that you can log the incoming query, user info etc with your bug
tracker

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3ceaa18941bfa5cfd26e1d1eaae499050f72c91d"><tt>3ceaa189</tt></a> add request middleware

- <a href="https://github.com/99designs/gqlgen/commit/877f75a07cca3abc836437608cf2e7d499ad3ff8"><tt>877f75a0</tt></a> remove debugging trace (closes <a href="https://github.com/99designs/gqlgen/issues/81"> #81</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/091d25ab6e9250c5eca7e5144bf05455ed1a8754"><tt>091d25ab</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/87">#87</a> from jon-walton/windows-paths</summary>

fix package paths on windows

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/53a6e8141d870f28176433233d98573fc83a890a"><tt>53a6e814</tt></a> fix package paths on windows

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/546b7b7607c9b6137253e528ab5cef9508ba7410"><tt>546b7b76</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/84">#84</a> from yamitzky/master</summary>

Fix collectFields to handle aliased fields properly

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ba2ecb166b6754a4684dcedff7cab532c025496d"><tt>ba2ecb16</tt></a> Add test case for aliased field

- <a href="https://github.com/99designs/gqlgen/commit/78f3a56cb3670a203f9484af57c05f54c94345ea"><tt>78f3a56c</tt></a> Fix collectFields to handle aliased fields

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4d2eece0b6a79c7c34b5485dc0678fe1fd1690e3"><tt>4d2eece0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/77">#77</a> from vektah/opentracing</summary>

Add resolver middleware Add opentracing support

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f0def668d3a7edb75e88830013ff957b3cf7e363"><tt>f0def668</tt></a> better opentracing tags

- <a href="https://github.com/99designs/gqlgen/commit/600bff7ae34e28ceab663317f07c3b583cdd59c0"><tt>600bff7a</tt></a> bump metalinter deadline

- <a href="https://github.com/99designs/gqlgen/commit/2e32c12162aa3952734d1ef47678b94fa279b275"><tt>2e32c121</tt></a> regenerate code

- <a href="https://github.com/99designs/gqlgen/commit/5b9085072d77a72e0d138a09ace6946f15fee827"><tt>5b908507</tt></a> opentracing middleware

- <a href="https://github.com/99designs/gqlgen/commit/57adb244df6b509c2ecf4bfb15ea7d4ecd00e17c"><tt>57adb244</tt></a> Add resolver middleware

- <a href="https://github.com/99designs/gqlgen/commit/28d0c81f077b8b01dc715de9cd96820fc197c36e"><tt>28d0c81f</tt></a> capture args in map

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e266fab9b98129526e78178a927e28f3242dc002"><tt>e266fab9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/75">#75</a> from mathewbyrne/fix-import-dash</summary>

Replace Invalid Characters in Package Name with an Underscore

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b0d79115d07506920fea5f3f93086c9b4849f9e3"><tt>b0d79115</tt></a> Replace invalid package characters with an underscore</summary>

This will sanatise local import names to a valid go identifier by
replacing any non-word characters with an underscore.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/66a915034875ec0ff4be588406272159d35c342e"><tt>66a91503</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/72">#72</a> from vektah/custom-enum</summary>

Add support for custom enums

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/61a34a7428aa897f1c14d5d3e5e47a655c2c236c"><tt>61a34a74</tt></a> Add support for custom enums

- <a href="https://github.com/99designs/gqlgen/commit/74ac827a9a417cd0bb4bfb0276cd61278ba0ae6e"><tt>74ac827a</tt></a> move docs to new domain

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ebcc94d153a1cf336ec5f866ac6c9514fba3f57e"><tt>ebcc94d1</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/70">#70</a> from vektah/models-in-separate-package</summary>

Allow generated models to go into their own package

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9a5321317b96de2404c56d0a4dcae12f4beff78f"><tt>9a532131</tt></a> Allow generated models to go into their own package

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6129fd266ac2d662404531df0070c47d402b4b35"><tt>6129fd26</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/69">#69</a> from vektah/support-options</summary>

Support OPTIONS requests

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/af38cf0571bbf4e43231f764508330d818fdcc5b"><tt>af38cf05</tt></a> Support OPTIONS requests

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/893ead12b72ade2fe965dab33bbf0c7a4179bcad"><tt>893ead12</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/67">#67</a> from vektah/raw-schema-string</summary>

Use a raw string for schema

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/af6178a7a782ca1e03200929d7c67371f6855cf8"><tt>af6178a7</tt></a> Use a raw string for schema

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c0753bed4b8687082b559540abbbdfdbf1a65a1b"><tt>c0753bed</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/66">#66</a> from vektah/generate-enums</summary>

Generate enums

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/85a51268892a5c33fdebd1bf51eb181d8b1e1b2a"><tt>85a51268</tt></a> Generate enums

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/71c4e2655c7ed81068f278598aa6cc2d9eed5b32"><tt>71c4e265</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/65">#65</a> from vektah/context</summary>

Make field selections available in context

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c60336bff6e6367f71d48279bceaeb7cbd331299"><tt>c60336bf</tt></a> regenerate

- <a href="https://github.com/99designs/gqlgen/commit/c5ccfe4e720b2ece5599ebe3fdf7a0473012acc2"><tt>c5ccfe4e</tt></a> Add an example for getting the selection sets from ctx

- <a href="https://github.com/99designs/gqlgen/commit/e7007746dea1694c97a6741fee7ec1351b5ef350"><tt>e7007746</tt></a> add fields to resolver context

- <a href="https://github.com/99designs/gqlgen/commit/40918d52d4f222f7d89bd09bdee6c0d2963a6720"><tt>40918d52</tt></a> move request scoped data into context

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4e13262e5b1b30b8f85fb580f00d1bc6f213ca8e"><tt>4e13262e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/64">#64</a> from vektah/vendor-gen-path</summary>

Fix vendored import paths in generated models

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/2ff9f32f6e6195f9ee8b760e872fc2f8b052de07"><tt>2ff9f32f</tt></a> Fix vendored import paths

- <a href="https://github.com/99designs/gqlgen/commit/630a3cfc60cb0dff9d32431b11c966bdcc596f58"><tt>630a3cfc</tt></a> failing test

- <a href="https://github.com/99designs/gqlgen/commit/99dec54c55dd7ec38d9e338ea643087955974363"><tt>99dec54c</tt></a> fix missing deps

- <a href="https://github.com/99designs/gqlgen/commit/652c567e145310921ee1159dc97d13965c46f932"><tt>652c567e</tt></a> Remove missing field warning and add test for scalar resolvers (closes <a href="https://github.com/99designs/gqlgen/issues/63"> #63</a>)

- <a href="https://github.com/99designs/gqlgen/commit/3dc87e1b08d7ba4eeae95b738dff14a631691603"><tt>3dc87e1b</tt></a> gtm

- <a href="https://github.com/99designs/gqlgen/commit/c76c34342e71c59a53e3307664d65c09eecc86a3"><tt>c76c3434</tt></a> Add dataloader tutorial

- <a href="https://github.com/99designs/gqlgen/commit/449fe8f823760340cd8569d7caed00fb97613eb4"><tt>449fe8f8</tt></a> Optimize frontmatter

- <a href="https://github.com/99designs/gqlgen/commit/b90ae60e80d3df104e94f2d35b7f8f687bbf738c"><tt>b90ae60e</tt></a> flatten menus

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a508ecb07fc8c5d714bba99f5174247d99e7eaca"><tt>a508ecb0</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/45">#45</a> from dvic/fix-resolver-public-errors</summary>

Retain orignal resolver error and support overriding error message

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ab4e7010a6ed1b75a5ab334b543806a2b1889628"><tt>ab4e7010</tt></a> Retain orignal resolver error and support overriding error message (closes <a href="https://github.com/99designs/gqlgen/issues/38"> #38</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a05a18d5a4bdc85962d17945659df4a7360dbf91"><tt>a05a18d5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/61">#61</a> from vektah/import-resolver-collisions</summary>

Deal with import collisions better

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d81ea2c23f9ddab0ca9fa8d9a34183f08a146a1d"><tt>d81ea2c2</tt></a> Deal with import collisions better

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/fb131a94816a31ade5dbeedd9c7fb15c87aaa1a1"><tt>fb131a94</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/59">#59</a> from vektah/map-support</summary>

Add map[string]interface{} escape hatch

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/49d921647337e2a54fca53b71d5ae3907e681f12"><tt>49d92164</tt></a> Add map[string]interface{} escape hatch

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/5abdba16ebfe032ef6c1b965d81b53099a3efc69"><tt>5abdba16</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/57">#57</a> from vektah/null-input-fields</summary>

Null input fields

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f8add9d2c79ce57f5b7776fd4265f7a72bebc4a3"><tt>f8add9d2</tt></a> remove more unneeded whitespace

- <a href="https://github.com/99designs/gqlgen/commit/84b066170081d0c2d16d06ec89e39323062e5d0e"><tt>84b06617</tt></a> Allow nulls in input fields

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/54fbe16a2788b870cd966cb8ad6581b4f7b016e3"><tt>54fbe16a</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/56">#56</a> from vektah/getting-started-fixes</summary>

Getting started fixes

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/17fd17a4318c9cb02d69a2bdd9dd3176857497d6"><tt>17fd17a4</tt></a> Update the tutorial

- <a href="https://github.com/99designs/gqlgen/commit/e65d2a5ac92eb9528d4f80c391e284571ff66d18"><tt>e65d2a5a</tt></a> detect correct FK type

- <a href="https://github.com/99designs/gqlgen/commit/b66cfa03ed53f8639b3cf0cf20a4ff8ebd604194"><tt>b66cfa03</tt></a> small fixes to entry point

- <a href="https://github.com/99designs/gqlgen/commit/0b62315a89d6c865b6d554789db9f817503d7da8"><tt>0b62315a</tt></a> Create ISSUE_TEMPLATE.md

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f3a70dacafc04b6895c39eb0f8789022bd476254"><tt>f3a70dac</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/55">#55</a> from vektah/fix-input-ptr-unpacking</summary>

Fix ptr unpacking in input fields

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/10541f1922656c7ed48153421484edb3da973a79"><tt>10541f19</tt></a> Fix ptr unpacking in input fields

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/15b3af2d5a6ed7240d09fb4c21fd905c8c316aa7"><tt>15b3af2d</tt></a> Fix value receivers for unions too</summary>

fixes 42

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/46103bdc1bc88c8cb6182f1d339c3b562c4e12b6"><tt>46103bdc</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/53">#53</a> from vektah/docs</summary>

Docs

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d0211a0a6e2096fb609279d93c75d432921854f0"><tt>d0211a0a</tt></a> Custom scalar docs

- <a href="https://github.com/99designs/gqlgen/commit/e6ed4de5ab7f3346b2133c756008f6deb998d79a"><tt>e6ed4de5</tt></a> Update readme link

- <a href="https://github.com/99designs/gqlgen/commit/51f08a9ec5a1622a057502ee67390c9329f8c19e"><tt>51f08a9e</tt></a> start of docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ac832dea46f2f1b80660f781c5bc9e40e6482101"><tt>ac832dea</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/51">#51</a> from vektah/support-embedding</summary>

Support embedding in models

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9d710712becfe5ee08c7e9d0940ddd14145f77f7"><tt>9d710712</tt></a> add embedding support

- <a href="https://github.com/99designs/gqlgen/commit/0980df0e999d0a7a7cf9fb98e280782d29a5e862"><tt>0980df0e</tt></a> Embedding example

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cb34e6db2695cc16b72a661a3bb8c5b57e6bb2a0"><tt>cb34e6db</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/50">#50</a> from vektah/valuer-receiver</summary>

Don't generate value receivers for types that cant fit the interface

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ec5f5e66a1303bc6dcca06b717e8e1cdd45e6b96"><tt>ec5f5e66</tt></a> check for valuer receivers before generating type switch

- <a href="https://github.com/99designs/gqlgen/commit/dc898409dca42e0ae392f4b768321545d30434c3"><tt>dc898409</tt></a> add test case

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/8ef253cbc3ae69e6da29b38b04a86287b2cc944a"><tt>8ef253cb</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/49">#49</a> from vektah/default-entrypoints</summary>

default to Query / Mutation / Subscription if no entry points are specified

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/302058a7705f008c6877373e2e701332b6e50469"><tt>302058a7</tt></a> Use default entry points for Query/Mutation/Subscription

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/13949fdf6977b05c3e4809e6678ac6a5c13c86e5"><tt>13949fdf</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/37">#37</a> from vektah/generate-interfaces</summary>

generate interfaces

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/acc45bf0bcda7e85be14d2489b193c1ed0228d4f"><tt>acc45bf0</tt></a> generate interfaces

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e47d5038e43963754073d2ea745cd8206e3a5756"><tt>e47d5038</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/34">#34</a> from vektah/root-types-only</summary>

Only bind to types in the root package scope

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/ffe972a878d485bae6d82130947f95833e9987f5"><tt>ffe972a8</tt></a> Only bind to types in the root package scope

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/0e78c0aeb9b73a57634e6934fa2e13234ec31d9a"><tt>0e78c0ae</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/33">#33</a> from vektah/unset-arguments</summary>

Allow unset arguments

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/bc9e0e54f69b88ff709f8d79d55f3491b30ff45e"><tt>bc9e0e54</tt></a> Allow unset arguments

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/947183514df552507c78193d95d62665f03820a9"><tt>94718351</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/31">#31</a> from vektah/recover-handler</summary>

Customizable recover func

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/e4e249ea5103c190226131fad48adfdfa6a3f551"><tt>e4e249ea</tt></a> Customizable recover func

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/69277045fd045dfd3f47508f2491809c1e91826f"><tt>69277045</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/30">#30</a> from vektah/complex-input-types</summary>

Fix complex input types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/9b64dd22cee0454c5cc48d2921e12213d531e6cc"><tt>9b64dd22</tt></a> Fix complex input types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1d074b89baedb19b7011976a213d7169a4f24793"><tt>1d074b89</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/29">#29</a> from vektah/multi-stage-model-build</summary>

Split model generation into its own stage

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cf580c24c96f9064706accdaaa80513e3c8d350e"><tt>cf580c24</tt></a> Split model generation into its own stage

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/926384db2ff70cf2bdbf7608feeb69666ddd9919"><tt>926384db</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/28">#28</a> from vektah/default-args</summary>

add default args

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/68c54a14debcc6a7b325d7e838c6a6b516d0dfda"><tt>68c54a14</tt></a> add default args

- <a href="https://github.com/99designs/gqlgen/commit/d63128f6ac1b0749b44b2386632174fe76b3be21"><tt>d63128f6</tt></a> appease the linting gods

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/7b6b124ed53950c6e038fa641e391e8bf1369333"><tt>7b6b124e</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/20">#20</a> from vektah/codegen-cleanup</summary>

Codegen cleanup

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/78c34cb3ca028519a8e4a83211a7e9a4d6f8ec3a"><tt>78c34cb3</tt></a> regenerate

- <a href="https://github.com/99designs/gqlgen/commit/5ebd157c798c7742656f802f9ecad5ada6c29183"><tt>5ebd157c</tt></a> Only use one gofunc per subscription

- <a href="https://github.com/99designs/gqlgen/commit/79a70376e65dd097929deb048cd21f43d6b4fa48"><tt>79a70376</tt></a> Move generated field resolvers into separate methods

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/e676abe4cdfc1ea97f9747cf0eb1d6a80c059b3e"><tt>e676abe4</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/19">#19</a> from vektah/generate-input-types</summary>

Generate input models too

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f094e79c3a24e820485cdf631eaeb666f8571447"><tt>f094e79c</tt></a> Generate input models too

- <a href="https://github.com/99designs/gqlgen/commit/1634f0882f3d469d0542cd5cd2943658df49a145"><tt>1634f088</tt></a> Add a missed error check

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4feb1689aa3b397c4ff6a1ff97c33f9a5a61fff2"><tt>4feb1689</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/18">#18</a> from vektah/array-input-args</summary>

Fix input array processing

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/98176297239463a57c35d05607051ecdfb4c59e9"><tt>98176297</tt></a> Fix input array processing

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/4880497fc0431e2a83e5525c1f2668f7c15328cd"><tt>4880497f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/16">#16</a> from vektah/better-templates</summary>

Better templates

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/278df9de416c542904b94c6a68b4b2f8f67966d3"><tt>278df9de</tt></a> Better templates

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/f3731c73087b7d4233fb76f5f8518e13446cc6f9"><tt>f3731c73</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/14">#14</a> from vektah/autogenerate-models</summary>

Autogenerate models

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cfe902a0b5c8ca9e8c12bccbb679530df2dff22a"><tt>cfe902a0</tt></a> Autogenerate models

- <a href="https://github.com/99designs/gqlgen/commit/287bf7f43a14a4e689f9a5ff7812c6a9f36076a0"><tt>287bf7f4</tt></a> more docs

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/9d896f4018f85b6c3db31b2221b83a49aeae8260"><tt>9d896f40</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/13">#13</a> from vektah/autocast</summary>

Automatically add type conversions around wrapped types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/85fa63b9570357b8ce954b88ccfd2ba2dd437d15"><tt>85fa63b9</tt></a> Automatically add type conversions around wrapped types

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/c8c2e40fe967a455c9a68dff69d1b5b100828ff6"><tt>c8c2e40f</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/11">#11</a> from vektah/subscriptions</summary>

Add support for subscriptions

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/d514b82940a4aeb3312b3703e926f9c209f910f9"><tt>d514b829</tt></a> Add some go tests to the chat app

- <a href="https://github.com/99designs/gqlgen/commit/ec2916d91c7cf3ebd0803b82b064b9b830f8e9e0"><tt>ec2916d9</tt></a> chat example for subscriptions using CRA+apollo

- <a href="https://github.com/99designs/gqlgen/commit/8f93bf8d8856091e5d19103911c74dd0e1d6fe1f"><tt>8f93bf8d</tt></a> get arg errors working in both contexts

- <a href="https://github.com/99designs/gqlgen/commit/62a18ff1ebd8ff1d4b370f217b03f63ac7641d6f"><tt>62a18ff1</tt></a> Update generator to build a new ExecutableSchema interface

- <a href="https://github.com/99designs/gqlgen/commit/c082c3a443d3547a83e0a594f06c69e3b58f8dea"><tt>c082c3a4</tt></a> prevent concurrent writes in subscriptions

- <a href="https://github.com/99designs/gqlgen/commit/f555aec6a3d3965288f67ebafe2b491f749be16d"><tt>f555aec6</tt></a> switch to graphql playground for better subscription support

- <a href="https://github.com/99designs/gqlgen/commit/182195413ab1573286a3e0ea05ca5ad8550e360b"><tt>18219541</tt></a> add websocket support to the handler

- <a href="https://github.com/99designs/gqlgen/commit/d4c7f3b988c4d74eeb67e55be40a45549e4498f7"><tt>d4c7f3b9</tt></a> update resolver definition to use channels for subscriptions

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d0244d24425231e7c021ef972993350482b69ff7"><tt>d0244d24</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/10">#10</a> from vektah/newtypes</summary>

User defined custom types

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5d86eeb60b9a82015e28deaf085349498b5ebb97"><tt>5d86eeb6</tt></a> fix jsonw test

- <a href="https://github.com/99designs/gqlgen/commit/944ee0884b565d82c2f567dc64f01bff1f6b2806"><tt>944ee088</tt></a> regenerate

- <a href="https://github.com/99designs/gqlgen/commit/4722a855959112565548f5b2fc48891afd8d641a"><tt>4722a855</tt></a> add scalar example

- <a href="https://github.com/99designs/gqlgen/commit/83b001aeb00773681279bc489197c08fff6794c2"><tt>83b001ae</tt></a> rename marshaler methods

- <a href="https://github.com/99designs/gqlgen/commit/e0b7c25f10b0b089109ede0edb5ac55c56e91272"><tt>e0b7c25f</tt></a> move collectFields out of generated code

- <a href="https://github.com/99designs/gqlgen/commit/146c65380cec9e5f5bf59641d84faafecb31c07b"><tt>146c6538</tt></a> generate input object unpackers

- <a href="https://github.com/99designs/gqlgen/commit/d94cfb1f8654a42c8ab21e80b54ee2a50dc744e3"><tt>d94cfb1f</tt></a> allow primitive scalars to be redefined

- <a href="https://github.com/99designs/gqlgen/commit/402e073076977744699cdf3ab961bd9967125561"><tt>402e0730</tt></a> rename jsonw to graphql

- <a href="https://github.com/99designs/gqlgen/commit/3e7d80dfe2afc488ea1a76f7ec727b5634bfca11"><tt>3e7d80df</tt></a> Update README.md

- <a href="https://github.com/99designs/gqlgen/commit/9c77e7a05659baef60d470502dd4ae903ac641a0"><tt>9c77e7a0</tt></a> Update dataloaden dep

- <a href="https://github.com/99designs/gqlgen/commit/530f7895b79a60c7b0fb820126e727dffa2c3dfa"><tt>530f7895</tt></a> Make gql client work with older versions of mapstructure

- <a href="https://github.com/99designs/gqlgen/commit/5c04d1adaddd2b67b5052caa33488a4a2c011df0"><tt>5c04d1ad</tt></a> __typename support

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/51292db98260ba3dd75630a87104d69b057a20b4"><tt>51292db9</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/4">#4</a> from vektah/cleanup-type-binding</summary>

Cleanup schema binding code

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/c89a8774650d41a4a99bf5b90d0c69c4a7a166a3"><tt>c89a8774</tt></a> Cleanup schema binding code

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/030954a5d49ccc5129c10edd921d169a33f79be4"><tt>030954a5</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/2">#2</a> from ulrikstrid/patch-1</summary>

Fix typo in README

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/cb507bd07bb4cdd65fd2f64a34d2dfa1e7319008"><tt>cb507bd0</tt></a> Fix typo in README

- <a href="https://github.com/99designs/gqlgen/commit/e3167785fce58eef451015cefd60b0a015a86b0b"><tt>e3167785</tt></a> Fix template loading from inside vendor

- <a href="https://github.com/99designs/gqlgen/commit/261b52ce9fb650099a8d4eb3c2105c5fa04b67a1"><tt>261b52ce</tt></a> fix an error handling bug

- <a href="https://github.com/99designs/gqlgen/commit/1da57f59f06576884f4d9ce793914f6ae8b572d0"><tt>1da57f59</tt></a> Split starwars models out from resolvers

- <a href="https://github.com/99designs/gqlgen/commit/743b2cf9741096ebddbbcad4d3e61c5c98fe5800"><tt>743b2cf9</tt></a> fix indenting

- <a href="https://github.com/99designs/gqlgen/commit/fb2d5817ebb2c447880b5b3525b56954ca9c742a"><tt>fb2d5817</tt></a> use gorunpkg to vendor go generate binaries

- <a href="https://github.com/99designs/gqlgen/commit/7f4d0405e0aad3c01d00533af603f029512d2f71"><tt>7f4d0405</tt></a> encourage dep use

- <a href="https://github.com/99designs/gqlgen/commit/3276c7824ad7f9de4ae0b83f00517e6eeb7ec623"><tt>3276c782</tt></a> Do not bind to unexported vars or methods

- <a href="https://github.com/99designs/gqlgen/commit/5fabffaf1765f9404fa6a1a4286df8888e6984ab"><tt>5fabffaf</tt></a> heading tweaks

- <a href="https://github.com/99designs/gqlgen/commit/e032c1d5e9a24c8813a282e895747946a0bb7bdc"><tt>e032c1d5</tt></a> Prior art

- <a href="https://github.com/99designs/gqlgen/commit/45b79a1e050299e9268f5c95588f90b3c1cb13e7"><tt>45b79a1e</tt></a> Add a test for multidimensional arrays

- <a href="https://github.com/99designs/gqlgen/commit/ec73a50a45e576bb8cc3ba886244e96d2c03eb32"><tt>ec73a50a</tt></a> fix race

- <a href="https://github.com/99designs/gqlgen/commit/75a3a05c2b34442c334f74baef551424a57be4b2"><tt>75a3a05c</tt></a> Dont execute mutations concurrently

- <a href="https://github.com/99designs/gqlgen/commit/3900a41db7a9afaade738ddcfe20b7463013ae6a"><tt>3900a41d</tt></a> tidy up json writing

- <a href="https://github.com/99designs/gqlgen/commit/0dcf7f6b0a792d98bac361e9a353ec6310d9a097"><tt>0dcf7f6b</tt></a> add circle ci badge

- <a href="https://github.com/99designs/gqlgen/commit/2c9bf21cb301fd1c3b42e6444002d4e1c4b16e2f"><tt>2c9bf21c</tt></a> get dataloaden

- <a href="https://github.com/99designs/gqlgen/commit/4fff3241e1a237bd15dd7300adec4a2954159bf5"><tt>4fff3241</tt></a> install dataloaden in ci

- <a href="https://github.com/99designs/gqlgen/commit/951f41b2c13bd2057bd271ca70e1c7a59d791920"><tt>951f41b2</tt></a> circle ci

- <a href="https://github.com/99designs/gqlgen/commit/8fa5f628c2c03f4ca8e23554cb7bce7a0f4e6038"><tt>8fa5f628</tt></a> less whitespace

- <a href="https://github.com/99designs/gqlgen/commit/c76f3b9883b107811f983902438920bfc7104dc6"><tt>c76f3b98</tt></a> clean up template layout

- <a href="https://github.com/99designs/gqlgen/commit/4a6cea5e50be1e1721d7ac2eb0089c325f68e068"><tt>4a6cea5e</tt></a> readme fixes

- <a href="https://github.com/99designs/gqlgen/commit/b814ad52598a681550b0421e09b4ea6e701cecae"><tt>b814ad52</tt></a> rename repo

- <a href="https://github.com/99designs/gqlgen/commit/9c79a37adfbfa205e05b78dd50081c7108406e8d"><tt>9c79a37a</tt></a> Cleanup and add tests

- <a href="https://github.com/99designs/gqlgen/commit/5afb5caa9f7ddd94a7f1f817be1589735e80c90d"><tt>5afb5caa</tt></a> update dataloaden

- <a href="https://github.com/99designs/gqlgen/commit/d00fae08a66e6f8a1eb4dbdaf72b7adb423e94b5"><tt>d00fae08</tt></a> Add dataloader example

- <a href="https://github.com/99designs/gqlgen/commit/86cdf3a0a7c57af5e9f9649c7d0c176c52f200bd"><tt>86cdf3a0</tt></a> Fix package resolution

- <a href="https://github.com/99designs/gqlgen/commit/41306cbab515668d3f6ba42d41a21f9779bb8f50"><tt>41306cba</tt></a> Better readme

- <a href="https://github.com/99designs/gqlgen/commit/ce5e38ed04f98aeb9eac7d2304b81004e7998bc0"><tt>ce5e38ed</tt></a> Add GET query param support to handler

- <a href="https://github.com/99designs/gqlgen/commit/dd9a8e4d1179532784ac9344cbf69071dc6240bb"><tt>dd9a8e4d</tt></a> parallel execution

- <a href="https://github.com/99designs/gqlgen/commit/4468127eeed9f747fdf590620906725d7312e9a9"><tt>4468127e</tt></a> pointer juggling

- <a href="https://github.com/99designs/gqlgen/commit/9e99c14929c18b624f2899551c79b1667df2afac"><tt>9e99c149</tt></a> Use go templates to generate code

- <a href="https://github.com/99designs/gqlgen/commit/41f74970749b289303bc1d1ded8b0d9a0bb99adb"><tt>41f74970</tt></a> Support go versions earlier than 1.9

- <a href="https://github.com/99designs/gqlgen/commit/c20ef3d0d5c448cbf377f407d4b277566eb5a1aa"><tt>c20ef3d0</tt></a> add missing nulls

- <a href="https://github.com/99designs/gqlgen/commit/bb753776138e12b08c0ceaa8b3059e0ba3cbdb5b"><tt>bb753776</tt></a> Use goimports instead of gofmt on generated code

- <a href="https://github.com/99designs/gqlgen/commit/c2cf38354c2db094761f7021cd9979d77fe280b8"><tt>c2cf3835</tt></a> coerce types between similar types

- <a href="https://github.com/99designs/gqlgen/commit/5297dd4090db7564fadc389845e5d743d087dfd2"><tt>5297dd40</tt></a> Add support for RFC3339 formatted Time as time.Time

- <a href="https://github.com/99designs/gqlgen/commit/61291ce9c1215156ceab87ec6b0fcd4c821c304f"><tt>61291ce9</tt></a> support vendor

- <a href="https://github.com/99designs/gqlgen/commit/6d437d7ea42017741f9209a95db23491dfc5c1d0"><tt>6d437d7e</tt></a> allow map[string]interface{} arg types

- <a href="https://github.com/99designs/gqlgen/commit/39a8090a4800c5788572d0de1a4a4fe223bf6847"><tt>39a8090a</tt></a> cleanup

- <a href="https://github.com/99designs/gqlgen/commit/a9352e3239cf58ff86e630a2d324b2292f6c6654"><tt>a9352e32</tt></a> gometalinter pass

- <a href="https://github.com/99designs/gqlgen/commit/9ab81d671b8a9417baa7d9e622440b3dbe819bda"><tt>9ab81d67</tt></a> Finish fleshing out the connection example

- <a href="https://github.com/99designs/gqlgen/commit/e04b1e50e8437fc44e6bc8d861dc1a0a71bfd38d"><tt>e04b1e50</tt></a> inline supporting runtime funcs

- <a href="https://github.com/99designs/gqlgen/commit/9cedf0122d4c1164475e69d038ecabfbd8e77267"><tt>9cedf012</tt></a> complex arg handling

- <a href="https://github.com/99designs/gqlgen/commit/0c9c009f2d0f47d5491a37ac7d561000feb60ca7"><tt>0c9c009f</tt></a> Clean up json writer

- <a href="https://github.com/99designs/gqlgen/commit/e7e18c401d874362132fcac3b8349130c8e265ea"><tt>e7e18c40</tt></a> much cleaner generated code

- <a href="https://github.com/99designs/gqlgen/commit/6a76bbf6d8f1fdc2f084f4de107381707dad2565"><tt>6a76bbf6</tt></a> Interfaces and starwars example

- <a href="https://github.com/99designs/gqlgen/commit/29110e76b56bbb4265f5a8327021a26d3d0e5cbb"><tt>29110e76</tt></a> Generate ESS to remove it interface{} casts completly

- <a href="https://github.com/99designs/gqlgen/commit/2f358e7daacfcfd861372d2b891e19ac909c47e2"><tt>2f358e7d</tt></a> graphiql autocomplete working

- <a href="https://github.com/99designs/gqlgen/commit/2e2c3135e6940657647d8d091962cda880638a5f"><tt>2e2c3135</tt></a> create separate type objects in prep for fragment support

- <a href="https://github.com/99designs/gqlgen/commit/22c0ad0a230ce02fe058aaaa2f971aee897983d5"><tt>22c0ad0a</tt></a> Add basic introspection support

- <a href="https://github.com/99designs/gqlgen/commit/c1c2cb6440efd45c92f4125029756923a16773d4"><tt>c1c2cb64</tt></a> Code generation

- <a href="https://github.com/99designs/gqlgen/commit/4be5ac84f45789d327fd18c5ec2bbcdf5a936659"><tt>4be5ac84</tt></a> args

- <a href="https://github.com/99designs/gqlgen/commit/bde800e19db009e91eb0386d240857b73ed122f8"><tt>bde800e1</tt></a> imports

- <a href="https://github.com/99designs/gqlgen/commit/596554da3c31f7b4de5fd8c81c3b4f4f36ebdb72"><tt>596554da</tt></a> start of code generator

- <a href="https://github.com/99designs/gqlgen/commit/62fa8184f7f9ec767f2b9a74482224e98b77df1f"><tt>62fa8184</tt></a> split generated vs exec code

- <a href="https://github.com/99designs/gqlgen/commit/0ea104cd9cec6b427261d7357eebff4e9a75a9bd"><tt>0ea104cd</tt></a> remove internal package

- <a href="https://github.com/99designs/gqlgen/commit/f81371e8431e8136985d9742fbea157e25bc34fd"><tt>f81371e8</tt></a> Args

- <a href="https://github.com/99designs/gqlgen/commit/01896b3bbd66b3e6131ab1242c34cc5b1761d78a"><tt>01896b3b</tt></a> Hand written codegen example

- <a href="https://github.com/99designs/gqlgen/commit/5a756bda0de8154469bc50db5a481ea86a90ea09"><tt>5a756bda</tt></a> Rewrite paths and add readme

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/b46637030579abd312c5eea21d36845b6e9e7ca4"><tt>b4663703</tt></a> trace: Log graphql.variables rather than tag</summary>

According to the OT documentation tag values can be numeric types, strings, or
bools. The behavior of other tag value types is undefined at the OpenTracing
level. For example `github.com/lightstep/lightstep-tracer-go` generates error
events.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/5d3b13f2e2215f7d53c88b20fd2c37f0a37b5ffd"><tt>5d3b13f2</tt></a> Support context injection in testing

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/beff08417b04a91f0a7618b55786d63eec833b9f"><tt>beff0841</tt></a> Separate literal arg parsing cases (int, float)</summary>

This change allows the ID scalar implementation to more semantically
handle the case for unmarshalling integer IDs.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/ab1dd4b5bf46b60897a048fab59628abe736812c"><tt>ab1dd4b5</tt></a> Add tests for ID scalar input</summary>

This commit adds two tests cases for ID scalar input:
- a string literal
- an integer literal

Both of these literal types are covered by the GraphQL specification as
valid input for the ID scalar.

Reference the ID section of the spec for more information:
http://facebook.github.io/graphql/October2016/#sec-ID

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/d8c57437fcbad670b16d06538151dfddc72eda88"><tt>d8c57437</tt></a> Extract ID scalar implmentation</summary>

This change moves the ID scalar implementation out of `graphql.go` and
into its own file `id.go` for consistency with the Time scalar
implementation.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/10eb949b8f4439212ef19d5d924c47c10c110cc7"><tt>10eb949b</tt></a> cleaned up example to use MustParseSchema

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/52080e1f0951c75dd4addba57db27fc4de1d0019"><tt>52080e1f</tt></a> Rename friendsConenctionArgs to friendsConnectionArgs</summary>

Fix spelling error in friendsConnectionArgs type

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3965041f0afca9aa0b87f45421eaa1244aada88f"><tt>3965041f</tt></a> Update GraphiQL interface (add history)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/6b9bc3e2dcb85a6f0d08b133d170ec3e60fe55d0"><tt>6b9bc3e2</tt></a> Add `(*Schema).Validate` (<a href="https://github.com/99designs/gqlgen/pull/99">#99</a>)</summary>

* Add `(*Schema).Validate`

This adds a `Validate` method to the schema, which allows you to find out if a query is valid without actually running it. This is valuable when you have a client with static queries and want to statically determine whether they are valid.

* Fix Validate doc string

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/7f3f7120c8849d4f3ed4c3c9421b9c05a4625d84"><tt>7f3f7120</tt></a> Set content-type header to `application/json`

- <a href="https://github.com/99designs/gqlgen/commit/c76ff4d892ae945f1c1a7b31cd19262d8031d597"><tt>c76ff4d8</tt></a> moved packer into separate package

- <a href="https://github.com/99designs/gqlgen/commit/073edccda088f21eb03050184190bbe1e66289eb"><tt>073edccd</tt></a> updated tests from graphql-js

- <a href="https://github.com/99designs/gqlgen/commit/3a9ac3683963f80ac4bb7446ff1870c760c50b72"><tt>3a9ac368</tt></a> validation: improved overlap check

- <a href="https://github.com/99designs/gqlgen/commit/f86c8b01d949de962be8254b486e7b3e99c52489"><tt>f86c8b01</tt></a> allow multiple schemas in tests

- <a href="https://github.com/99designs/gqlgen/commit/77750960d3c3002eb2b17268834795782448e654"><tt>77750960</tt></a> validation: OverlappingFieldsCanBeMerged

- <a href="https://github.com/99designs/gqlgen/commit/e7ca4fde4e2fe3aa5c71ae6ff8af605728ec6dd9"><tt>e7ca4fde</tt></a> refactor: remove SelectionSet type

- <a href="https://github.com/99designs/gqlgen/commit/7aad6ba78fc19e4b89b0eae12263e79f56eeb97b"><tt>7aad6ba7</tt></a> refactor: use schema.NamedType

- <a href="https://github.com/99designs/gqlgen/commit/fddcbcb791d050972e8eb4177c7e48b14db38824"><tt>fddcbcb7</tt></a> resolves <a href="https://github.com/99designs/gqlgen/pull/92">#92</a>: fix processing of negative scalars during parse literals

- <a href="https://github.com/99designs/gqlgen/commit/48c1a0fb9c1adcec045be7453f74eed6b77f6419"><tt>48c1a0fb</tt></a> Small fix based on feedback.

- <a href="https://github.com/99designs/gqlgen/commit/e90d10895bd124a6519d08e4c18fa49a4f592b5f"><tt>e90d1089</tt></a> allow custom types as input arguments

- <a href="https://github.com/99designs/gqlgen/commit/dd3d39e28aa2c58ad04f5a608d8ff5c15f9db518"><tt>dd3d39e2</tt></a> fix panic when variable name not declared

- <a href="https://github.com/99designs/gqlgen/commit/c2bc105ff947eb09cdc5df7f452813c67472e0a0"><tt>c2bc105f</tt></a> validation: NoUnusedVariables

- <a href="https://github.com/99designs/gqlgen/commit/4aff2976b4cc7aea5311d994920aeb0023b09c47"><tt>4aff2976</tt></a> refactor

- <a href="https://github.com/99designs/gqlgen/commit/0933d24133db6ffbe62da026e62e6e5e4c7711d2"><tt>0933d241</tt></a> validation: VariablesInAllowedPosition

- <a href="https://github.com/99designs/gqlgen/commit/83e2f31aa8b8fd94d6fb0d47dbf1676907e07631"><tt>83e2f31a</tt></a> validation: NoUndefinedVariables

- <a href="https://github.com/99designs/gqlgen/commit/c39ffecaa7a19ae0d5703923a7e7c7d018d29e23"><tt>c39ffeca</tt></a> validation: PossibleFragmentSpreads

- <a href="https://github.com/99designs/gqlgen/commit/47c5cde7110bca8e5c78771db80d8c216dc1cd18"><tt>47c5cde7</tt></a> validation: UniqueInputFieldNames

- <a href="https://github.com/99designs/gqlgen/commit/94cb291812ee54a780bc70d065532e2be952be95"><tt>94cb2918</tt></a> big refactoring around literals

- <a href="https://github.com/99designs/gqlgen/commit/3d63ae8037964a34f0660baa752973f549de8ee3"><tt>3d63ae80</tt></a> some refactoring

- <a href="https://github.com/99designs/gqlgen/commit/969dab9d2fc07c21a64c0566a6374a2ef7854950"><tt>969dab9d</tt></a> merged lexer into package "common"

- <a href="https://github.com/99designs/gqlgen/commit/a9de61717bc8ac77dcf5c4e537b12d0788c34c4f"><tt>a9de6171</tt></a> renamed lexer.Literal to lexer.BasicLit

- <a href="https://github.com/99designs/gqlgen/commit/88c492bbb06d76cb7f07883e7194e73b4fcb93a9"><tt>88c492bb</tt></a> validation: NoFragmentCycles (closes <a href="https://github.com/99designs/gqlgen/issues/38"> #38</a>)

- <a href="https://github.com/99designs/gqlgen/commit/d39712c819b716ef2e87097bffb37988e13af9e4"><tt>d39712c8</tt></a> refactor addErrMultiLoc

- <a href="https://github.com/99designs/gqlgen/commit/ee5e1c3baa2b0eb4854942596e82f7dea686a83f"><tt>ee5e1c3b</tt></a> validation: updated tests

- <a href="https://github.com/99designs/gqlgen/commit/490ad6b2b9bc4e3fa0f89b04dd12ee1d9dd0f1bc"><tt>490ad6b2</tt></a> validation: NoUnusedFragments

- <a href="https://github.com/99designs/gqlgen/commit/da85f09dd939b282cf7df78870d0bd71ca6d6681"><tt>da85f09d</tt></a> add path to errors on resolver error or panic (closes <a href="https://github.com/99designs/gqlgen/issues/86"> #86</a>)

- <a href="https://github.com/99designs/gqlgen/commit/04cb2550483c3cf827b6668b69383e1f363e8dd4"><tt>04cb2550</tt></a> allow structs without pointers (closes <a href="https://github.com/99designs/gqlgen/issues/78"> #78</a>)

- <a href="https://github.com/99designs/gqlgen/commit/4c40b305eb4c8a2abad853dd53c6df51f050199f"><tt>4c40b305</tt></a> show all locations in error string

- <a href="https://github.com/99designs/gqlgen/commit/5c26f320e3a296c653c6b3763280883cf4e4b416"><tt>5c26f320</tt></a> fix limiter

- <a href="https://github.com/99designs/gqlgen/commit/dbc3f0a094e9cf25aeeb6106d78587ca878613e1"><tt>dbc3f0a0</tt></a> fix composing of fragments (closes <a href="https://github.com/99designs/gqlgen/issues/75"> #75</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/213a5d013a4e0f7e540931cc1f3581d183412232"><tt>213a5d01</tt></a> Warn if an interface's resolver has a ToTYPE implementation that does not return two values.</summary>

Currently this instead crashes fairly inscrutably at runtime here: https://github.com/neelance/graphql-go/blob/master/internal/exec/exec.go#L117

An alternate fix would be to check len(out) there and perhaps rely on out[0] being nil to continue if there's only one return value.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/00c4c5743b2fff2500fbd0643d3a43e297139221"><tt>00c4c574</tt></a> Fix panic when resolver is not a pointer

- <a href="https://github.com/99designs/gqlgen/commit/d0df6d8a2d50fdc903c88b60f1fe925020df2790"><tt>d0df6d8a</tt></a> small cleanup

- <a href="https://github.com/99designs/gqlgen/commit/036945e2bef3692493bfdbd42d1a8231ae5d45a3"><tt>036945e2</tt></a> fix hang on panic (fixes <a href="https://github.com/99designs/gqlgen/pull/82">#82</a>)

- <a href="https://github.com/99designs/gqlgen/commit/01ab5128e53e25d42cc15dd04168cacb310fc0a9"><tt>01ab5128</tt></a> Add supports for snake case (<a href="https://github.com/99designs/gqlgen/pull/77">#77</a>)

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/67e6f91d3f998c19d0a23d292852af2589a03d15"><tt>67e6f91d</tt></a> use encoding/json to encode scalars</summary>

There are some edge cases that are better handled by the proven encoder of encoding/json, for example special characters in strings.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/3f1cb6f83af94ebfc621b7dd2f83bbaa70046252"><tt>3f1cb6f8</tt></a> implement user defined logger (with sensible defaults)

- <a href="https://github.com/99designs/gqlgen/commit/b357f4641697a00d450101cd21349a0adc03d130"><tt>b357f464</tt></a> built-in json encoding

- <a href="https://github.com/99designs/gqlgen/commit/2d828770c3ce02a9d98ec9fbcb46621f37004298"><tt>2d828770</tt></a> refactor: collect fields to resolve

- <a href="https://github.com/99designs/gqlgen/commit/32f8b6ba2bd8cdb042f47c48e98ede15df1b4531"><tt>32f8b6ba</tt></a> refactor: replaced MetaField

- <a href="https://github.com/99designs/gqlgen/commit/b95c566e05321524e962333dbc202806adaf3586"><tt>b95c566e</tt></a> simplify schema introspection

- <a href="https://github.com/99designs/gqlgen/commit/4200a584f2986af42ab07cba6a794e89dc8f6b54"><tt>4200a584</tt></a> split internal/exec into multiple packages

- <a href="https://github.com/99designs/gqlgen/commit/c11687a72110091da7964ca370fcbc9c5bc912ee"><tt>c11687a7</tt></a> refactored internal/exec

- <a href="https://github.com/99designs/gqlgen/commit/bd742d84b6026f484c7d05e6f5dec90a113ec2ca"><tt>bd742d84</tt></a> WIP

- <a href="https://github.com/99designs/gqlgen/commit/d09dd543c7f5fbcb68f7cf40df9a2d005f04c283"><tt>d09dd543</tt></a> added SchemaOpt

- <a href="https://github.com/99designs/gqlgen/commit/1dcc5753f3f91f1269617e2350124b52189719ed"><tt>1dcc5753</tt></a> fix Schema.ToJSON

- <a href="https://github.com/99designs/gqlgen/commit/4f07e397b1ea3e26621bffdd8b261cb77a47ee9a"><tt>4f07e397</tt></a> pass variable types to tracer

- <a href="https://github.com/99designs/gqlgen/commit/36e6c97e53891214b0d6286a3a2720f2d1dbc790"><tt>36e6c97e</tt></a> readme: remove outdated section about opentracing

- <a href="https://github.com/99designs/gqlgen/commit/0b143cca346bd478e684aaa66fba557890a8f8ce"><tt>0b143cca</tt></a> refactor: apply before exec

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/a992060271937699b725b5a72b8bfbcdc06ee7ca"><tt>a9920602</tt></a> pluggable tracer</summary>

Improved performance while keeping flexibility.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/58d3d5b8274933506d071426c63408c7820d6bcc"><tt>58d3d5b8</tt></a> refactored exec.Request

- <a href="https://github.com/99designs/gqlgen/commit/9dd714ec00dbce57bfc1df09b8a59bd557211873"><tt>9dd714ec</tt></a> refactored execField some more

- <a href="https://github.com/99designs/gqlgen/commit/a43ef2411e0c59d63c0f4575d13cf73a4aeb51e2"><tt>a43ef241</tt></a> refactor: meta fields

- <a href="https://github.com/99designs/gqlgen/commit/48931d17313dd9f706a5f4b36188e84ab60cb847"><tt>48931d17</tt></a> refactor fieldExec

- <a href="https://github.com/99designs/gqlgen/commit/ee95710db59e7abf49eac7accd9f81cfc7a60bfb"><tt>ee95710d</tt></a> small cleanup

- <a href="https://github.com/99designs/gqlgen/commit/84baade55e292ffcc1be7dc0d0ba0c5805d3aa11"><tt>84baade5</tt></a> perf: create span label only once

- <a href="https://github.com/99designs/gqlgen/commit/a16ed60054ea3a6ccb8c118367f6af785338657c"><tt>a16ed600</tt></a> improved concurrency architecture

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/aef3d9cf7adefb68aabbb3d20c99db9beb087f98"><tt>aef3d9cf</tt></a> Add testing.go into its own package (gqltesting)</summary>

This is done so that "testing" (and especially its registered cli flags)
aren't included in any production builds.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/f78108a335e635dffd75aa83a5f2114c1f76ca18"><tt>f78108a3</tt></a> validation: meta fields

- <a href="https://github.com/99designs/gqlgen/commit/c6ab2374926b50199eca22c8f63e53f21b37e714"><tt>c6ab2374</tt></a> added empty file to make CI happy

- <a href="https://github.com/99designs/gqlgen/commit/d59c1709c5e8daf56c18275df9fd3793a6430a44"><tt>d59c1709</tt></a> fix introspection of default value

- <a href="https://github.com/99designs/gqlgen/commit/42608a035615d89dfd4ea7e8543d1f847d72bbb3"><tt>42608a03</tt></a> clean up now unnecessary check

- <a href="https://github.com/99designs/gqlgen/commit/e45f26dd458a2171a12ed77ed867529ff51eb01d"><tt>e45f26dd</tt></a> validation: UniqueDirectivesPerLocation

- <a href="https://github.com/99designs/gqlgen/commit/dcf7e59f141277e4da7d8e5f671aeb4bcbec1dfe"><tt>dcf7e59f</tt></a> validation: UniqueFragmentNames

- <a href="https://github.com/99designs/gqlgen/commit/eeaa510b3ec2d6babfe823350df93bea71896c90"><tt>eeaa510b</tt></a> validation: UniqueOperationNames

- <a href="https://github.com/99designs/gqlgen/commit/a5a11604b187c8104f8688babe9fade58dc3797f"><tt>a5a11604</tt></a> refactor: Loc on Field

- <a href="https://github.com/99designs/gqlgen/commit/b5919db43e5c1a7fbeab1e91e022e2ce03aaa6c4"><tt>b5919db4</tt></a> validation: UniqueVariableNames

- <a href="https://github.com/99designs/gqlgen/commit/8632753a6219518df212d6519f43d9d4666f8d8c"><tt>8632753a</tt></a> validation: ScalarLeafs

- <a href="https://github.com/99designs/gqlgen/commit/4584498444b623e12480bbc03cb8447a9156b63a"><tt>45844984</tt></a> validation: ProvidedNonNullArguments

- <a href="https://github.com/99designs/gqlgen/commit/c741ea84f9bb6f39c9f41672a350831a2526071a"><tt>c741ea84</tt></a> validation: VariablesAreInputTypes

- <a href="https://github.com/99designs/gqlgen/commit/0875d74f6aebbdcc1226def90090bb7772b44b9c"><tt>0875d74f</tt></a> validation: UniqueArgumentNames

- <a href="https://github.com/99designs/gqlgen/commit/1fdab07f1df987d1b6beb19bc628bd01b39a5128"><tt>1fdab07f</tt></a> validation: LoneAnonymousOperation

- <a href="https://github.com/99designs/gqlgen/commit/090df527284031385ebc1073bf431ffa0df62cb0"><tt>090df527</tt></a> validation: KnownTypeNames

- <a href="https://github.com/99designs/gqlgen/commit/f99ca95eea69cffbfa899e9408a99367b3d96550"><tt>f99ca95e</tt></a> refactor: validation context

- <a href="https://github.com/99designs/gqlgen/commit/8aac28174b16d5d519f158c955abf5986914129d"><tt>8aac2817</tt></a> validation: KnownFragmentNames

- <a href="https://github.com/99designs/gqlgen/commit/eae3efc9479546a60af4b8ddb68c36e074d00325"><tt>eae3efc9</tt></a> validation: KnownDirectives

- <a href="https://github.com/99designs/gqlgen/commit/70581168ba2fea03a35f1f271ffe2f09d3c5c635"><tt>70581168</tt></a> refactor: separate InlineFragment and FragmentSpread

- <a href="https://github.com/99designs/gqlgen/commit/d6aec0d65bdb976f49f32fa71f7d7648f0f9a110"><tt>d6aec0d6</tt></a> renamed schema.Directive to DirectiveDecl

- <a href="https://github.com/99designs/gqlgen/commit/b616eeca1e84157410bf8531897bf9cfd390809c"><tt>b616eeca</tt></a> validation: KnownArgumentNames

- <a href="https://github.com/99designs/gqlgen/commit/885af6079bd4e417283418e0f3a7a1bb716ba5cf"><tt>885af607</tt></a> refactor: Location without pointer

- <a href="https://github.com/99designs/gqlgen/commit/5a40251c951f5ae584f77b6d97b04255a3da06b6"><tt>5a40251c</tt></a> tests: filter errors to currently tested rule

- <a href="https://github.com/99designs/gqlgen/commit/9c054f5304bd0032b0c654f7351b35a842fdf88c"><tt>9c054f53</tt></a> refactor: lexer.Ident

- <a href="https://github.com/99designs/gqlgen/commit/254afa8a3c1fbee15faf0ce61b7982fe6dc54c4a"><tt>254afa8a</tt></a> validation: fragment type

- <a href="https://github.com/99designs/gqlgen/commit/b6ef81af178ee49eeb4b1f24e96aecf98a86251d"><tt>b6ef81af</tt></a> added test export script

- <a href="https://github.com/99designs/gqlgen/commit/95a4ecd841312e27b542fb147295923832fdc7bc"><tt>95a4ecd8</tt></a> validation: fields

- <a href="https://github.com/99designs/gqlgen/commit/c387449f4deea4a7531c87e7f10b66ccef089b6f"><tt>c387449f</tt></a> validation: default values

- <a href="https://github.com/99designs/gqlgen/commit/44c6e634bae779e1a115b3443df5d73eb18696d7"><tt>44c6e634</tt></a> validation: arguments

- <a href="https://github.com/99designs/gqlgen/commit/30dcc339f36601cf099e4853093d93b1c63387b1"><tt>30dcc339</tt></a> directive arguments as slice

- <a href="https://github.com/99designs/gqlgen/commit/d331ac27e70a409e7602ee2ee7964ab12f7946b0"><tt>d331ac27</tt></a> input values as slice

- <a href="https://github.com/99designs/gqlgen/commit/615afd61aacfdc6f408c518e3a8e4c0928d97209"><tt>615afd61</tt></a> fields as slice

- <a href="https://github.com/99designs/gqlgen/commit/607599043f6a3b3d1875405f22efb3706c759933"><tt>60759904</tt></a> arguments as slice

- <a href="https://github.com/99designs/gqlgen/commit/f7d9ff4e09a02f4b09800eaf3604a36df90514e6"><tt>f7d9ff4e</tt></a> refactor literals

- <a href="https://github.com/99designs/gqlgen/commit/2e1fef012d13bf299d1a5b9949571fffc9ef7abb"><tt>2e1fef01</tt></a> keep track of location of arguments

- <a href="https://github.com/99designs/gqlgen/commit/29e0b375539193232c7b0298d5283773a2dfdc47"><tt>29e0b375</tt></a> added EnumValue type

- <a href="https://github.com/99designs/gqlgen/commit/aa868e8d461160e65058b1976484dcb74e264b0b"><tt>aa868e8d</tt></a> resolve fragments early

- <a href="https://github.com/99designs/gqlgen/commit/adeb53d684ee2f38b387c9063e11dd598ba9adc2"><tt>adeb53d6</tt></a> remove resolver from query package

- <a href="https://github.com/99designs/gqlgen/commit/2e23573fa55611dde28cbb8b0e18ce24960ee786"><tt>2e23573f</tt></a> parse directive decl without arguments

- <a href="https://github.com/99designs/gqlgen/commit/36f8ba8ba769fe7d76bb6f5954cad2b4d721b3e9"><tt>36f8ba8b</tt></a> fix introspection of default value (closes <a href="https://github.com/99designs/gqlgen/issues/65"> #65</a>)

- <a href="https://github.com/99designs/gqlgen/commit/e06f58558c5a1b319f1bea12bc157cf2d25f0aa9"><tt>e06f5855</tt></a> support for "deprecated" directive on enum values

- <a href="https://github.com/99designs/gqlgen/commit/498fe3961c3058c45f5ad90b1212a48adfbba266"><tt>498fe396</tt></a> support for [@deprecated](https://github.com/deprecated) directive on fields (fixes <a href="https://github.com/99designs/gqlgen/pull/64">#64</a>)

- <a href="https://github.com/99designs/gqlgen/commit/93ddece9c068b9e3ce6f0c85f8517e60047fb9f5"><tt>93ddece9</tt></a> refactor: DirectiveArgs

- <a href="https://github.com/99designs/gqlgen/commit/8f5605a1414369c7da0ee6b5e424d554a8b2e718"><tt>8f5605a1</tt></a> refactor directives

- <a href="https://github.com/99designs/gqlgen/commit/faf5384a347efef90900c63ff0c1ee800d005cb8"><tt>faf5384a</tt></a> simplify parseArguments

- <a href="https://github.com/99designs/gqlgen/commit/b2c2e906436885fac70ee32d4a05af4cc82d56e8"><tt>b2c2e906</tt></a> some more docs

- <a href="https://github.com/99designs/gqlgen/commit/f45165236be72481237a834a86313ff65338e8d9"><tt>f4516523</tt></a> added some method documentations

- <a href="https://github.com/99designs/gqlgen/commit/91bd7f887b6c9e45856c74db19e2081de48d1e60"><tt>91bd7f88</tt></a> improved meta schema

- <a href="https://github.com/99designs/gqlgen/commit/10dc8ee62965c02c6bca323fd6516814fdadb303"><tt>10dc8ee6</tt></a> added support for directive declarations in schema

- <a href="https://github.com/99designs/gqlgen/commit/28028f6677bd02cc1da876f7aa49d229004090cd"><tt>28028f66</tt></a> readme: more info on current project status

- <a href="https://github.com/99designs/gqlgen/commit/e9afca38415b36b594b585bf044a49e723e79713"><tt>e9afca38</tt></a> hint in error if method only exists on pointer type (fixes <a href="https://github.com/99designs/gqlgen/pull/60">#60</a>)

- <a href="https://github.com/99designs/gqlgen/commit/356ebd93204134d83349e001302811d9a47850bc"><tt>356ebd93</tt></a> nicer error messages (fixes <a href="https://github.com/99designs/gqlgen/pull/56">#56</a>)

- <a href="https://github.com/99designs/gqlgen/commit/e413f4edabe636dc05c0a5fabce95b60ca0fbf70"><tt>e413f4ed</tt></a> make gocyclo happy

- <a href="https://github.com/99designs/gqlgen/commit/6e92795e8e04c3d961a30ea313cdb617c1195d9c"><tt>6e92795e</tt></a> fix spelling

- <a href="https://github.com/99designs/gqlgen/commit/306e27ef29eb170e61557d02551effacdd802627"><tt>306e27ef</tt></a> gofmt -s

- <a href="https://github.com/99designs/gqlgen/commit/612317b28bf7e9a793d7b11d2cac81109ad1f988"><tt>612317b2</tt></a> fix ToJSON

- <a href="https://github.com/99designs/gqlgen/commit/728e57a9c3c084ec90938d7b4346bd004f3c305e"><tt>728e57a9</tt></a> improved doc for MaxParallelism

- <a href="https://github.com/99designs/gqlgen/commit/e8590a10d5acadaa06f4655d877544677d6983bb"><tt>e8590a10</tt></a> don't execute any further resolvers after context got cancelled

- <a href="https://github.com/99designs/gqlgen/commit/644435cc1084b8393d3f2d780a0547a535500787"><tt>644435cc</tt></a> added MaxParallelism

- <a href="https://github.com/99designs/gqlgen/commit/21802a339d523fae0fe0b8cee0e6ad5005c85bec"><tt>21802a33</tt></a> readme: add Sourcegraph badge

- <a href="https://github.com/99designs/gqlgen/commit/5b2978fcb1baf5d104a51d2638e2204abac0f5fd"><tt>5b2978fc</tt></a> added support for recursive input values

- <a href="https://github.com/99designs/gqlgen/commit/8c84afb1a622fdf7a44ce471a4f9070137044056"><tt>8c84afb1</tt></a> improved structure of "make exec" code

- <a href="https://github.com/99designs/gqlgen/commit/d5a6ca4953dbade97f08c7e05d1d92df612167c8"><tt>d5a6ca49</tt></a> make sure internal types don't get exposed

- <a href="https://github.com/99designs/gqlgen/commit/c9d4d865c19532267da57d7b681a4a8f39b40eec"><tt>c9d4d865</tt></a> fixed some null handling

- <a href="https://github.com/99designs/gqlgen/commit/a336dd4be28093e214c2404a2af88f32ac978b05"><tt>a336dd4b</tt></a> added request.resolveVar

- <a href="https://github.com/99designs/gqlgen/commit/943f80f48b8255bdef89a4edae33106f5f6e2dc4"><tt>943f80f4</tt></a> added unmarshalerPacker type

- <a href="https://github.com/99designs/gqlgen/commit/f77f73392b42621b62fa7a0af690f85abd647334"><tt>f77f7339</tt></a> refactored non-null handling in packer

- <a href="https://github.com/99designs/gqlgen/commit/ae0f1689b8e8ae5ab847fa09d3ea6c7d14d2419e"><tt>ae0f1689</tt></a> remove hasDefault flag from makePacker

- <a href="https://github.com/99designs/gqlgen/commit/9cbad485080affc6d94126f55fdfdabf996262e7"><tt>9cbad485</tt></a> allow Unmarshaler for all types, not just scalars

- <a href="https://github.com/99designs/gqlgen/commit/f565a119801e2faa7fecb6ca7386e7aa933e61c6"><tt>f565a119</tt></a> refactored "make exec" code

- <a href="https://github.com/99designs/gqlgen/commit/07a09e5d93da69b5e571076540e15bc35414e2e8"><tt>07a09e5d</tt></a> properly check scalar types of result values

- <a href="https://github.com/99designs/gqlgen/commit/ecceddec6e3bbc4df96536a211a9a71f0919e47d"><tt>ecceddec</tt></a> Add ResolverError field to QueryError for post processing

- <a href="https://github.com/99designs/gqlgen/commit/b7c59ab9f042d11a73cf1c5fedc75538d8bca6e6"><tt>b7c59ab9</tt></a> renamed type

- <a href="https://github.com/99designs/gqlgen/commit/5817d30019edf5984356016bceea085af3f26bc8"><tt>5817d300</tt></a> moved some introspection code into new package, added Schema.Introspect

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/cdef8563513e825d0c77b7ef13856169780e05af"><tt>cdef8563</tt></a> removed SchemaBuilder</summary>

It is not necessary any more. Simpler API wins.

</details></dd></dl>

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/518a5fe7542d07e61d163a1c199cce9ac0b6a078"><tt>518a5fe7</tt></a> Merge pull request <a href="https://github.com/99designs/gqlgen/pull/45">#45</a> from nicksrandall/master</summary>

fix wrong import statement

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/8112e7191fa2341244e176b3b00025229f9afc53"><tt>8112e719</tt></a> fix wrong import statement

- <a href="https://github.com/99designs/gqlgen/commit/7fafcc6ec1a2a6c69f8775fb30f933ca352ab9db"><tt>7fafcc6e</tt></a> allow single value as implicit list (fixes <a href="https://github.com/99designs/gqlgen/pull/41">#41</a>)

- <a href="https://github.com/99designs/gqlgen/commit/2b513d7e6c46b01de0cbcaeb1b68424dd3c043be"><tt>2b513d7e</tt></a> improved custom types

- <a href="https://github.com/99designs/gqlgen/commit/191422c4ccf1fef693582e7f18bc263e6f2e09d2"><tt>191422c4</tt></a> merged code for value coercion and packing

- <a href="https://github.com/99designs/gqlgen/commit/232356b393b32a9db53bc49d35c6a08ede9dad26"><tt>232356b3</tt></a> introspection for "skip" and "include" directives (fixes <a href="https://github.com/99designs/gqlgen/pull/30">#30</a>)

- <a href="https://github.com/99designs/gqlgen/commit/2e10f7b852261d9afa33637e6978fecd511320d1"><tt>2e10f7b8</tt></a> readme: spec version and link (fixes <a href="https://github.com/99designs/gqlgen/pull/35">#35</a>)

- <a href="https://github.com/99designs/gqlgen/commit/61eca4c7621fe86c75ad96ddaffb6107086e65ea"><tt>61eca4c7</tt></a> pretty print SchemaBuilder.ToJSON

- <a href="https://github.com/99designs/gqlgen/commit/5e09ced15c3be3cd26ace9c2dc8edc9554edaf74"><tt>5e09ced1</tt></a> fix "null" value for empty descriptions of types

- <a href="https://github.com/99designs/gqlgen/commit/33cd194fc2cf0d2fa2c1768d0b0e60d2d0606237"><tt>33cd194f</tt></a> SchemaBuilder.ToJSON instead of SchemaToJSON (fixes <a href="https://github.com/99designs/gqlgen/pull/29">#29</a>)

- <a href="https://github.com/99designs/gqlgen/commit/fff173bcbe1bdc521f919c91e897078163d6c552"><tt>fff173bc</tt></a> proper error message when using non-input type as input (<a href="https://github.com/99designs/gqlgen/pull/19">#19</a>)

- <a href="https://github.com/99designs/gqlgen/commit/b94f2afe22033bcb4b1a68d64adab00986abeea4"><tt>b94f2afe</tt></a> improved support for null

- <a href="https://github.com/99designs/gqlgen/commit/4130d540fbff248ea1990a36ec7b1198befc28db"><tt>4130d540</tt></a> added support for input object literals

- <a href="https://github.com/99designs/gqlgen/commit/663e466fda86388ee5239a6d32162f39f3eafbac"><tt>663e466f</tt></a> moved some code into separate file

- <a href="https://github.com/99designs/gqlgen/commit/728e071e474353acb30d2530c50044acf7b48918"><tt>728e071e</tt></a> added support for lists as input values (fixes <a href="https://github.com/99designs/gqlgen/pull/19">#19</a>)

- <a href="https://github.com/99designs/gqlgen/commit/86f0f14544112d824efa3d586f47ccc1b04a1d26"><tt>86f0f145</tt></a> fix Float literals

- <a href="https://github.com/99designs/gqlgen/commit/b07f277bb052b43f27580766945d9d4a80a79e8e"><tt>b07f277b</tt></a> raise error on unexported input field (fixes <a href="https://github.com/99designs/gqlgen/pull/24">#24</a>)

- <a href="https://github.com/99designs/gqlgen/commit/4838c6f3bf4d198be8ebc81f5f39cf08bfd7c27b"><tt>4838c6f3</tt></a> fix optional input fields (fixes <a href="https://github.com/99designs/gqlgen/pull/25">#25</a>)

- <a href="https://github.com/99designs/gqlgen/commit/a15deed4f8fe354dc77feb60a45b29eb5671de72"><tt>a15deed4</tt></a> better way to implement GraphQL interfaces (<a href="https://github.com/99designs/gqlgen/pull/23">#23</a>)

- <a href="https://github.com/99designs/gqlgen/commit/7a66d0e02d9fdb4150e045360c1008dafe92d63b"><tt>7a66d0e0</tt></a> add support for description comments (fixes <a href="https://github.com/99designs/gqlgen/pull/20">#20</a>)

- <a href="https://github.com/99designs/gqlgen/commit/0b3be40c0717ac6955d2ddd0021bb9e16d4faac7"><tt>0b3be40c</tt></a> improved tracing

- <a href="https://github.com/99designs/gqlgen/commit/da879f4f78a67c3031a9ff9537cc49437ef38113"><tt>da879f4f</tt></a> small improvements to readme

- <a href="https://github.com/99designs/gqlgen/commit/f3f24cf6f1946333562bd7f90ef4f7afe08f00c6"><tt>f3f24cf6</tt></a> added some documentation

- <a href="https://github.com/99designs/gqlgen/commit/38598d83ded236fac7b1f7a8383d72cc42bbf26c"><tt>38598d83</tt></a> added CI badge to readme

- <a href="https://github.com/99designs/gqlgen/commit/bab81332446865f44ea10d773f024da9e0d9414a"><tt>bab81332</tt></a> starwars example: fix pagination panic (<a href="https://github.com/99designs/gqlgen/pull/12">#12</a>)

- <a href="https://github.com/99designs/gqlgen/commit/5ce3ca69fa6e3c7db69914f4ccb26ac9bf71cbd3"><tt>5ce3ca69</tt></a> testing: proper error on invalid ExpectedResult

- <a href="https://github.com/99designs/gqlgen/commit/8f7d2b1efd9e96e423c0dc0ecf7dbc905e1b5ff0"><tt>8f7d2b1e</tt></a> added relay.Handler

- <a href="https://github.com/99designs/gqlgen/commit/fce75a50a4f393bbf1bcef85835d406b03228b87"><tt>fce75a50</tt></a> properly coerce Int input values (<a href="https://github.com/99designs/gqlgen/pull/8">#8</a>)

- <a href="https://github.com/99designs/gqlgen/commit/0dd38747e3d5907dfc56d098bda073e2f5c34a4b"><tt>0dd38747</tt></a> star wars example: pass operation name and variables (<a href="https://github.com/99designs/gqlgen/pull/8">#8</a>)

- <a href="https://github.com/99designs/gqlgen/commit/3b7efd5cb6e73337e26a67ca05750b6e1b02320e"><tt>3b7efd5c</tt></a> fix __typename for concrete object types (fixes <a href="https://github.com/99designs/gqlgen/pull/9">#9</a>)

- <a href="https://github.com/99designs/gqlgen/commit/35667edabfbc519b7d9fc5e1ca0d3c4943448df5"><tt>35667eda</tt></a> testing tools

- <a href="https://github.com/99designs/gqlgen/commit/84571820f69feb18afd28609d997feac4e1bfab3"><tt>84571820</tt></a> only create schema once for tests

- <a href="https://github.com/99designs/gqlgen/commit/de113f969220b95e48cb64439069a13370dc26d6"><tt>de113f96</tt></a> added MustParseSchema

- <a href="https://github.com/99designs/gqlgen/commit/d5e5f6096fb8fbf48725ecca9014c47c31de58f4"><tt>d5e5f609</tt></a> improved structure for tests

- <a href="https://github.com/99designs/gqlgen/commit/947a1a3a8a25fd0821569a402dd433c751c58787"><tt>947a1a3a</tt></a> added package with tools for Relay

- <a href="https://github.com/99designs/gqlgen/commit/65f3e2b186c2bd9cb439505dd92be586ee6205de"><tt>65f3e2b1</tt></a> fix SchemaToJSON

- <a href="https://github.com/99designs/gqlgen/commit/cec7cea1c3b771e99b0c1fb45dd3fbf530c20d71"><tt>cec7cea1</tt></a> better error handling

- <a href="https://github.com/99designs/gqlgen/commit/e3386b067b33fb73ca885fc5b25c5a1249134671"><tt>e3386b06</tt></a> improved type coercion and explicit ID type

- <a href="https://github.com/99designs/gqlgen/commit/2ab9d765d642ae067a0055a4b2bde56eb6fcd461"><tt>2ab9d765</tt></a> support for custom scalars (fixes <a href="https://github.com/99designs/gqlgen/pull/3">#3</a>)

- <a href="https://github.com/99designs/gqlgen/commit/bdfd5ce306598d599f7d26bc9c4f9a70fe0c0c66"><tt>bdfd5ce3</tt></a> use custom error type less

- <a href="https://github.com/99designs/gqlgen/commit/0a7a37d1a7a6a8c8e84ac4b3d0f427ea3c55892c"><tt>0a7a37d1</tt></a> more flexible API for creating a schema

- <a href="https://github.com/99designs/gqlgen/commit/bd20a165e0aa4428f40424d68846b68d6e330a65"><tt>bd20a165</tt></a> improved type handling

- <a href="https://github.com/99designs/gqlgen/commit/ffa9fea4d939e28562475e1beac0826176fd64fc"><tt>ffa9fea4</tt></a> renamed GraphQLError to QueryError

- <a href="https://github.com/99designs/gqlgen/commit/fcfa135a03366e19b4900b6381ab5e06b3e58b8c"><tt>fcfa135a</tt></a> refactor

- <a href="https://github.com/99designs/gqlgen/commit/c28891d831baca59c10027b31a67deda775b3fff"><tt>c28891d8</tt></a> added support for OpenTracing

- <a href="https://github.com/99designs/gqlgen/commit/2cf7fcc8b709f160d5fcb3d1746fbe9e18967782"><tt>2cf7fcc8</tt></a> added SchemaToJSON

- <a href="https://github.com/99designs/gqlgen/commit/f6b498ac52dbb4cc87b22872f3b9844dadd1ffb5"><tt>f6b498ac</tt></a> stricter type mapping for input values

- <a href="https://github.com/99designs/gqlgen/commit/3c15e177dc0be85c5aa1d71a18a9c7966c5db34f"><tt>3c15e177</tt></a> execute mutations serially

- <a href="https://github.com/99designs/gqlgen/commit/1faf666161862d92bf8396bf89048d297cd4f850"><tt>1faf6661</tt></a> fix missing error

- <a href="https://github.com/99designs/gqlgen/commit/de9b7fed219a29adb0973819c9d7affb57d78d51"><tt>de9b7fed</tt></a> add support for mutations

- <a href="https://github.com/99designs/gqlgen/commit/094061d8ce65dbc0177a90a4bb2c6c4ab548b230"><tt>094061d8</tt></a> improved error handling a bit

- <a href="https://github.com/99designs/gqlgen/commit/cdb088d6e0df8357c1caba39dc93208753f611a4"><tt>cdb088d6</tt></a> refactor: args as input object

- <a href="https://github.com/99designs/gqlgen/commit/b06d39411d92439009084d342c1f3cdd98cb6339"><tt>b06d3941</tt></a> refactor: values

- <a href="https://github.com/99designs/gqlgen/commit/4fd33958e6645b471b742f07148c29e1ab19e155"><tt>4fd33958</tt></a> refactor: improved type system

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/1d03e667370a66100eead7f198eb457cd29a58f3"><tt>1d03e667</tt></a> refactor: new package "common"</summary>

package "query" does not depend on "schema" any more

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/1a959516cacda256f52b3ba2f5cf2e7fd03c85da"><tt>1a959516</tt></a> refactor

- <a href="https://github.com/99designs/gqlgen/commit/f8cb11c10ab15ad4a547a1695efe538a7b5209a8"><tt>f8cb11c1</tt></a> example/starwars: use interface base type to make type assertions nicer

- <a href="https://github.com/99designs/gqlgen/commit/746da4b8e1dff2ec72167060487e1441506a97c3"><tt>746da4b8</tt></a> star wars example: friendsConnection

- <a href="https://github.com/99designs/gqlgen/commit/bec45364913260dd699d295837336a4df824b7a6"><tt>bec45364</tt></a> parse type in query

- <a href="https://github.com/99designs/gqlgen/commit/be87a8fa7cfe42a2ea4748a525d0490a73e155c0"><tt>be87a8fa</tt></a> remove unused code

- <a href="https://github.com/99designs/gqlgen/commit/042e306a2478445bb9eecb7fb43b04892c8ee85f"><tt>042e306a</tt></a> simpler way to resolve type refs in schema

- <a href="https://github.com/99designs/gqlgen/commit/7cbf85fb2507b5b46af44e84ea195639fd95b202"><tt>7cbf85fb</tt></a> improved type checking of arguments

- <a href="https://github.com/99designs/gqlgen/commit/2b6460ae9c13be41f1aafbc7e881446f47c9c558"><tt>2b6460ae</tt></a> type check for scalars

- <a href="https://github.com/99designs/gqlgen/commit/17034fe713fd4396fc8483e9c2dc124833377fcb"><tt>17034fe7</tt></a> improved null handling

- <a href="https://github.com/99designs/gqlgen/commit/e6b6fbcab4297dd19e706d3981f7b6a98a581ee8"><tt>e6b6fbca</tt></a> small cleanup

- <a href="https://github.com/99designs/gqlgen/commit/7b8cd1bce34ce543781376f42062bc7234e8d8d1"><tt>7b8cd1bc</tt></a> meta schema from graphql-js

- <a href="https://github.com/99designs/gqlgen/commit/9333c0b396c6c19c306ba2d5e0b0b9fe024abd41"><tt>9333c0b3</tt></a> introspection: inputFields

- <a href="https://github.com/99designs/gqlgen/commit/c4faac56a02f7743de55bc20b261fb749578afa8"><tt>c4faac56</tt></a> introspection: ofType

- <a href="https://github.com/99designs/gqlgen/commit/86da849221ab53ae653eef12c12afc7604f603fd"><tt>86da8492</tt></a> introspection: interfaces and possibleTypes

- <a href="https://github.com/99designs/gqlgen/commit/20dbb84517b342ecec87e6ca8fba0951e5aa5d71"><tt>20dbb845</tt></a> proper nil support for lists

- <a href="https://github.com/99designs/gqlgen/commit/2e3369aebc265d07b5b1b025fcb6aa768c70ccaf"><tt>2e3369ae</tt></a> resolve types in schema package

- <a href="https://github.com/99designs/gqlgen/commit/7da95f4a932a8ce4dae85f7b8924294d041f55d9"><tt>7da95f4a</tt></a> introspection: enum values

- <a href="https://github.com/99designs/gqlgen/commit/cb423e6e1d36044a00a6d8ee62336844294e54fb"><tt>cb423e6e</tt></a> improved handling of scalar types

- <a href="https://github.com/99designs/gqlgen/commit/5b07780f940edfc6f7ee4c5d0bdd3f598a158998"><tt>5b07780f</tt></a> introspection: original order for fields and args

- <a href="https://github.com/99designs/gqlgen/commit/1e2d180c24bcdf8cfa945b3491bcceb8bbde9565"><tt>1e2d180c</tt></a> introspection: arguments

- <a href="https://github.com/99designs/gqlgen/commit/f21131bbc269ce4804ad7065f508abd14f131c10"><tt>f21131bb</tt></a> refactored schema to be more in line with introspection

- <a href="https://github.com/99designs/gqlgen/commit/0152d4f21a0b204df5a0f6df5c398a29734f3e68"><tt>0152d4f2</tt></a> introspection: currently no descriptions and deprecations

- <a href="https://github.com/99designs/gqlgen/commit/ad5689bbed8bd2fc7ac5e7005cc30ab2c16472e5"><tt>ad5689bb</tt></a> field introspection

- <a href="https://github.com/99designs/gqlgen/commit/2749d81451ae5033ab78095634dfd058d4e77c60"><tt>2749d814</tt></a> removed query.TypeReference

<dl><dd><details><summary><a href="https://github.com/99designs/gqlgen/commit/2eb105ec0bfe26ca1e6272c363a922394a9d95ef"><tt>2eb105ec</tt></a> Revert "resolve scalar types in exec"</summary>

This reverts commit fb3a6fc969b0c8c286c7d024a108f5696627639c.

</details></dd></dl>

- <a href="https://github.com/99designs/gqlgen/commit/40682d680d3613b866629e768cd23b39c1415346"><tt>40682d68</tt></a> removed exec.typeRefExec

- <a href="https://github.com/99designs/gqlgen/commit/64ea90fec366a3dcc09c428e177d96496a4d374a"><tt>64ea90fe</tt></a> makeWithType

- <a href="https://github.com/99designs/gqlgen/commit/2966f213e10bf0d1ca5b338da57558f42a7e05c6"><tt>2966f213</tt></a> added nonNilExec

- <a href="https://github.com/99designs/gqlgen/commit/c12a8ad39cdc3d4f484bc70e2f636f1eca0400d4"><tt>c12a8ad3</tt></a> added support for ints and floats in query

- <a href="https://github.com/99designs/gqlgen/commit/0f85412bbdf9bf4b1d984bd949022867b02204cd"><tt>0f85412b</tt></a> improved example

- <a href="https://github.com/99designs/gqlgen/commit/22ce46d1adeed7cf7022a223dfa7fc4b72e500da"><tt>22ce46d1</tt></a> support for optional error result

- <a href="https://github.com/99designs/gqlgen/commit/0fe56128d58946add9f4b03b9e3047a0dd2eb697"><tt>0fe56128</tt></a> optional context parameter

- <a href="https://github.com/99designs/gqlgen/commit/f1bc9b21f69b8e66144c2fe78faf52ea4e677a7b"><tt>f1bc9b21</tt></a> syntax errors with proper line and column

- <a href="https://github.com/99designs/gqlgen/commit/ae299efc1456ae143d6ffdf10c0ad9e71ef840b9"><tt>ae299efc</tt></a> proper response format

- <a href="https://github.com/99designs/gqlgen/commit/9619721b0ce1d21c82c35cc1bbb101b48da42e9e"><tt>9619721b</tt></a> added support for contexts

- <a href="https://github.com/99designs/gqlgen/commit/267fc6316b9ac4dc46dba9cdda70bd4057b09d4d"><tt>267fc631</tt></a> refactor

- <a href="https://github.com/99designs/gqlgen/commit/2e56e7ea619c95e419b298658b4042958d0934b3"><tt>2e56e7ea</tt></a> renamed NewSchema to ParseSchema

- <a href="https://github.com/99designs/gqlgen/commit/356b6e6bb14143f021d28e13070dad6482d2d1b9"><tt>356b6e6b</tt></a> added godoc badge

- <a href="https://github.com/99designs/gqlgen/commit/03f2e72dd506695e0ab30e73f8b8d18cd92df601"><tt>03f2e72d</tt></a> added README.md

- <a href="https://github.com/99designs/gqlgen/commit/1134562aaca2c8d08974dd33427093af79eeb5d4"><tt>1134562a</tt></a> added non-null type

- <a href="https://github.com/99designs/gqlgen/commit/8fa415513ce4b83bad465887c3d61c5c665ada80"><tt>8fa41551</tt></a> renamed Input to InputObject

- <a href="https://github.com/99designs/gqlgen/commit/6f2399aa0eec93f27bb784a7436e4da08cd92f30"><tt>6f2399aa</tt></a> introspection: type kind

- <a href="https://github.com/99designs/gqlgen/commit/e2c58f2f7163522224bd7b39fced1a44cd3967f4"><tt>e2c58f2f</tt></a> refactor: schema types for interface and input

- <a href="https://github.com/99designs/gqlgen/commit/0c8c9436ae2809df20f918956002aec688f2a1f7"><tt>0c8c9436</tt></a> introspection: __type

- <a href="https://github.com/99designs/gqlgen/commit/99a37521cb7af233924ae365edaf8497ba4bf3b6"><tt>99a37521</tt></a> refactoring: calculate "implemented by" in schema package

- <a href="https://github.com/99designs/gqlgen/commit/1cac7e5657f680593a1b2c3bfdfa4e3cee952cde"><tt>1cac7e56</tt></a> introspection: queryType

- <a href="https://github.com/99designs/gqlgen/commit/cc348faf3c98e7107d2ed484c90063094c22e2bc"><tt>cc348faf</tt></a> first bit of introspection

- <a href="https://github.com/99designs/gqlgen/commit/fb3a6fc969b0c8c286c7d024a108f5696627639c"><tt>fb3a6fc9</tt></a> resolve scalar types in exec

- <a href="https://github.com/99designs/gqlgen/commit/4cb8dcc015ba05002165d496a38a6e9ecb05fdf1"><tt>4cb8dcc0</tt></a> panic handlers

- <a href="https://github.com/99designs/gqlgen/commit/c7a528d4df37c4211cf3303ff75a68bc02a2e99d"><tt>c7a528d4</tt></a> proper error handling when creating schema

- <a href="https://github.com/99designs/gqlgen/commit/ae37381cb14fa6d851ef3a709cb5ba69ab0196af"><tt>ae37381c</tt></a> add support for __typename

- <a href="https://github.com/99designs/gqlgen/commit/4057080f8dcb75e6207fb9f1e357329738dcda50"><tt>4057080f</tt></a> add support for union types

- <a href="https://github.com/99designs/gqlgen/commit/d304a418586a1c33ccfcd8df017f70122fcd6d62"><tt>d304a418</tt></a> attribute source of star wars schema

- <a href="https://github.com/99designs/gqlgen/commit/0fcab871feeb20d445ba02aa411fa042ba6f47f9"><tt>0fcab871</tt></a> added LICENSE

- <a href="https://github.com/99designs/gqlgen/commit/0dc0116d69be9bc31688961cd653c047e602e4d2"><tt>0dc0116d</tt></a> support for inline fragments

- <a href="https://github.com/99designs/gqlgen/commit/f5e7d0709417463bd70730e244e2e8e515d14009"><tt>f5e7d070</tt></a> support for type assertions

- <a href="https://github.com/99designs/gqlgen/commit/fcb853c628f31d439af7c76de2f86e12317f561d"><tt>fcb853c6</tt></a> refactoring: addResultFn

- <a href="https://github.com/99designs/gqlgen/commit/741343f809a4bb85dc93cd4405738da44d886ce6"><tt>741343f8</tt></a> explicit fragment spread exec

- <a href="https://github.com/99designs/gqlgen/commit/73759258e589f4ded5977bd6a63a1048ecb7954e"><tt>73759258</tt></a> all missing stubs for star wars example

- <a href="https://github.com/99designs/gqlgen/commit/edc78e2bb7f2ad1aab86c4e6876d4b01e56fbcb6"><tt>edc78e2b</tt></a> parallelism

- <a href="https://github.com/99designs/gqlgen/commit/fb63371482c981cc1ef73be35245f0b75b58e8de"><tt>fb633714</tt></a> collect fields

- <a href="https://github.com/99designs/gqlgen/commit/08f02a2b149efb2c9d3fa1e5a090eecb07b8235e"><tt>08f02a2b</tt></a> execs

- <a href="https://github.com/99designs/gqlgen/commit/d70d16c4ebc4a9a8513ccab945b4e0ad77de6cfe"><tt>d70d16c4</tt></a> added server example

- <a href="https://github.com/99designs/gqlgen/commit/6f9a89db8ecf538d72ece47158d2ef7f487d8a06"><tt>6f9a89db</tt></a> separate example/starwars package

- <a href="https://github.com/99designs/gqlgen/commit/e4060db594dbd7ca6a98e7bb9cf967b6bc51e063"><tt>e4060db5</tt></a> added support for directives

- <a href="https://github.com/99designs/gqlgen/commit/89b066523b38cdbed48f1b6df03ed87bac737caf"><tt>89b06652</tt></a> added support for variables

- <a href="https://github.com/99designs/gqlgen/commit/78065ecbd8480e39506e7443141bad031caaf68b"><tt>78065ecb</tt></a> added support for enums

- <a href="https://github.com/99designs/gqlgen/commit/18645e60bbc9a9124fba220912a482a1d3fc0238"><tt>18645e60</tt></a> added support for query fragments

- <a href="https://github.com/99designs/gqlgen/commit/84f532b9b25363a224ddff3e415773d88de1b152"><tt>84f532b9</tt></a> added support for aliases

- <a href="https://github.com/99designs/gqlgen/commit/59d2a619ad146e6ec34b96601d4c10e9689c0a77"><tt>59d2a619</tt></a> improved support for arguments

- <a href="https://github.com/99designs/gqlgen/commit/edce4ec8712f3bf56252773292b1a22379e18e46"><tt>edce4ec8</tt></a> proper star wars data

- <a href="https://github.com/99designs/gqlgen/commit/d6ffc01de6704abc68346c6a07abe8a25f504341"><tt>d6ffc01d</tt></a> syntax support for full star wars schema

- <a href="https://github.com/99designs/gqlgen/commit/b582410448f091f9e2e8ad6975b2a8cf8cd09c01"><tt>b5824104</tt></a> support for comments

- <a href="https://github.com/99designs/gqlgen/commit/2f9ce9b48c3a85f538a66a1a556f5a62d3c3a20a"><tt>2f9ce9b4</tt></a> support for entry points

- <a href="https://github.com/99designs/gqlgen/commit/0b3d103849df9e267894eaf61d597887818d718d"><tt>0b3d1038</tt></a> support for arguments

- <a href="https://github.com/99designs/gqlgen/commit/cff8b3020fc132fc493a7bffdd77ff46ef7253e5"><tt>cff8b302</tt></a> support for arrays

- <a href="https://github.com/99designs/gqlgen/commit/565e59f53ab2e59f41c0eee010ab9525349c35a3"><tt>565e59f5</tt></a> schema package

- <a href="https://github.com/99designs/gqlgen/commit/1ae71ba2d5990082bdff59622fb636ed6105ed8e"><tt>1ae71ba2</tt></a> query package

- <a href="https://github.com/99designs/gqlgen/commit/42c13e7a09f05d98f6f7541199a66f140651ad1b"><tt>42c13e7a</tt></a> named types, complex objects

- <a href="https://github.com/99designs/gqlgen/commit/bf64e5dad2916d25a04a95d007cff263307b964a"><tt>bf64e5da</tt></a> initial commit

 <!-- end of Commits -->
<!-- end of Else -->

<!-- end of If NoteGroups -->
<!-- end of Versions -->
<!-- end of If Versions -->

# Start here: Gallery People draft

Your work is saved on GitHub in the [draft branch](https://github.com/ucbbrla9011-jpg/stash232323/tree/codex/gallery-people-draft-20261027). This is a checkpoint you can return to; it is **not a finished plugin or Windows installer**.

## What we are building

Keep Stash's normal galleries and image viewer. On your Windows machine, select images and choose **Analyze and tag selected images**. The planned plugin will process those images locally, suggest groups of the same person, and add hair-colour and visible skin-tone tags. You will be able to label groups, correct matches, and filter images by person or appearance.

Only selected images will be analyzed. New imports will not trigger analysis automatically. The face groups and appearance tags are separate so, for example, a person's hair colour can change between photographs.

## What is done

- Reviewed the original [Stash project](https://github.com/stashapp/stash), its plugin/API contracts, and existing face-recognition/community plugins.
- Confirmed this fork's base matches the reviewed upstream `develop` commit. The full Stash source is already present.
- Saved a detailed [compatibility plan](docs/GALLERY_PEOPLE_PLAN.md) for Windows and libraries with thousands of images.
- Saved an early Python prototype and its dependency/model manifests. It is clearly labeled unfinished.
- Prepared cloud Go and pnpm tooling; installed locked UI dependencies and verified Go module checksums.
- Successfully generated the UI's GraphQL client, passed its TypeScript check, and built the existing Stash frontend. The build emitted upstream deprecation/bundle-size warnings.
- Cloned the additional repository, `ucbbrla9011-jpg/Stashnewprivate`, on its selected branch `6221`. It is a separate, clean checkout; no work was removed or moved into it.
- Created a separate GitHub draft branch and a verified local recovery archive. `develop` was not changed by this checkpoint.

## What is not finished

- Backend generation was stopped at the requested handoff; its caches are retained. A complete backend build, running server, gallery smoke test, and backend test suite are still outstanding.
- The prototype's UI does not exist yet, and the prototype still contains a scan-all task that must be replaced with selected-image analysis.
- Hair-colour and skin-tone analysis are planned, not implemented.
- Face model download received a cloud proxy `403 Forbidden`. The face models have not been loaded or evaluated.
- Windows installation, cancellation/restart, recognition quality, correction/undo, and large-library performance have not been tested.

Do not install `plugins/gallery-people` in your real Stash library yet. Its configuration references frontend files that have not been written.

## Your next steps — no Git experience needed

1. Open the [saved draft branch](https://github.com/ucbbrla9011-jpg/stash232323/tree/codex/gallery-people-draft-20261027). Read this file and the compatibility plan there. The branch is the GitHub backup of the proposed work.
2. To download a copy, choose **Code → Download ZIP** while viewing that branch. The ZIP contains source and plans, not a ready-to-run Windows app.
3. If you want a draft pull request now, use [GitHub's pull-request page for this branch](https://github.com/ucbbrla9011-jpg/stash232323/pull/new/codex/gallery-people-draft-20261027). Choose `develop` as the base and `codex/gallery-people-draft-20261027` as the compare branch, then choose **Create draft pull request**. Do not merge the unfinished prototype.
4. In the cloud environment settings, review the saved network additions: `storage.googleapis.com` for Go artifacts, `media.githubusercontent.com` for model downloads, and `api.github.com` for pull-request creation. Preserve your existing domains. Saving the draft alone does not enable these connections; review/save the environment settings and publish when you want to activate a prepared environment version.
5. To continue later, tell Codex: **“Resume the Gallery People draft. Read README-GALLERY-PEOPLE.md and docs/GALLERY_PEOPLE_STATUS.md first. Preserve both selected repositories. Finish baseline validation before implementing the selected-image workflow.”**

GitHub accepted the branch push. Automatic draft-PR creation failed because the cloud proxy denied access to `api.github.com`; no pull request was created by the agent.

## Next implementation steps

1. Finish the ordinary Stash backend build, startup, and relevant image/gallery tests using isolated sample data.
2. Add the native selected-image action and review page through the plugin API, preserving existing gallery controls and other plugins.
3. Evaluate local face and appearance models, their licenses, and Windows support. Decide auto-tagging thresholds from measured results; preserve uncertain predictions for review.
4. Implement naming, merge/split/reject, additive tags and performer assignments, durable corrections, resumable jobs, and provenance-aware undo.
5. Test a fresh Windows installation and package a downloadable application with simple setup/start launchers. The final user should not need to build Go or React code.

## Recovery details

- Cloud work checkout: `/workspace/stash232323`.
- Additional checkout: `/workspace/Stashnewprivate`.
- Local archive: `/workspace/checkpoints/gallery-people-draft-20261027.tar.gz`.
- Archive manifest/checksum: `/workspace/checkpoints/gallery-people-draft-20261027.json`.
- Cloud tooling/cache/logs: `/workspace/stash-cloud/`.
- Full handoff: [GALLERY_PEOPLE_STATUS.md](docs/GALLERY_PEOPLE_STATUS.md).

The GitHub branch preserves source and documentation. It does not contain installed cloud dependencies, images, credentials, face data, or a published environment snapshot.

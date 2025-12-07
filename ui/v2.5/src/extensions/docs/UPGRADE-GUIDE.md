# Upstream Upgrade Guide

This guide explains how to upgrade from upstream stashapp/stash and reapply fork customizations.

## Upstream Relationship

- **Upstream repo:** `stashapp/stash`
- **Upstream branch:** `develop` (main development branch)
- **Current baseline:** v0.29.3 (the release this fork was based on)

When upgrading, you typically merge from `upstream/develop`, not release tags.

## Quick Reference

After every upstream upgrade, you need to:

1. **Re-apply patches** - The files documented in `extensions/patches/` will be overwritten
2. **Regenerate GraphQL** - Run `yarn generate` if graphql files changed
3. **Verify build** - Run `yarn build && yarn test`

---

## Step-by-Step Upgrade Process

### Step 1: Create a Backup Branch

```bash
git checkout -b backup-before-upgrade
git push origin backup-before-upgrade
git checkout main
```

### Step 2: Fetch and Merge Upstream

```bash
git fetch upstream
git merge upstream/develop  # merge from upstream's develop branch
# or: git rebase upstream/develop
```

### Step 3: Resolve Conflicts

Extensions directory should have **no conflicts** (self-contained).

For conflicts in other files, refer to the appropriate patch doc.

### Step 4: Verify Extensions are Intact

```bash
# Extensions should be untouched after merge
# (upstream doesn't have an extensions folder)
ls src/extensions/
# Should show all our folders: components, docs, filters, hooks, etc.
```

### Step 5: Re-apply Patches

Run this command to see which patched files changed in upstream:

```bash
# Compare against HEAD before the merge (or your last sync point)
# Replace LAST_SYNC with the commit/tag before you merged (e.g., HEAD~1 after merge)
LAST_SYNC="HEAD~1"  # or a specific commit hash

# List files that have patches AND changed in upstream
git diff --name-only $LAST_SYNC -- \
  src/components/List/ListFilter.tsx \
  src/components/List/ListTable.tsx \
  src/components/List/Pagination.tsx \
  src/components/Galleries/GalleryCard.tsx \
  src/components/Groups/GroupCard.tsx \
  src/components/MainNavbar.tsx \
  src/components/Scenes/SceneListTable.tsx \
  src/components/Settings/SettingsInterfacePanel/SettingsInterfacePanel.tsx \
  src/components/ScenePlayer/ScenePlayer.tsx \
  src/components/Scenes/SceneDetails/ \
  src/components/Galleries/GalleryDetails/ \
  src/components/Images/ImageDetails/ \
  src/components/Shared/ \
  src/components/FrontPage/ \
  src/models/ \
  src/utils/ \
  graphql/
```

If any files are listed, consult the corresponding patch document.

### Step 6: Regenerate GraphQL Types

```bash
yarn generate
```

### Step 7: Build and Test

```bash
yarn build
yarn test --run
```

**The test suite includes upgrade verification tests** that check:
- Critical utility functions exist and work
- Type structures are valid
- Import paths resolve correctly
- Regression prevention (labels, empty maps, loading states)

If tests fail after an upgrade, check:
- `upgrade-verification.test.ts` failures indicate broken extension imports
- Other test failures may indicate API changes in upstream

### Step 8: Manual Verification

Test these features in the browser:
- [ ] Scene list with filters
- [ ] Sidebar filter counts
- [ ] Scene player (settings menu)
- [ ] AI recommendations (front page)
- [ ] Detail pages (Gallery, Image, Scene)

---

## Patch Application Checklist

Use this checklist after every upgrade:

### Always Apply (Import Changes)

These files have import paths pointing to extensions:

| File | Import From |
|------|-------------|
| `Scenes/Scenes.tsx` | `src/extensions/components/Scene/Scene` |
| `Shared/TagLink.tsx` | `src/extensions/components/GalleryPopover` |
| `List/EditFilterDialog.tsx` | `src/extensions/ui` |
| `List/ItemList.tsx` | `src/extensions/ui` |
| `List/ListToolbar.tsx` | `src/extensions/ui` |
| `Galleries/GalleryDetails/Gallery.tsx` | `src/extensions/components` |
| `Images/ImageDetails/Image.tsx` | `src/extensions/components` |
| `Groups/GroupDetails/GroupPerformersPanel.tsx` | `src/extensions/facets/enhanced` |
| `Groups/GroupDetails/GroupScenesPanel.tsx` | `src/extensions/facets/enhanced` |
| `FrontPage/Control.tsx` | `src/extensions/components` |

**Check:** After merge, verify these imports still point to extensions.

### Conditional Apply (If Upstream Changed)

Only re-apply if the upstream file was modified. Use `$LAST_SYNC` (your previous merge point):

```bash
# Set this to your last sync point (commit before merge, or previous baseline)
LAST_SYNC="HEAD~1"  # adjust as needed
```

| Patch File | Check Command |
|------------|---------------|
| `list-components.md` | `git diff $LAST_SYNC -- src/components/List/ListFilter.tsx src/components/List/ListTable.tsx src/components/List/Pagination.tsx` |
| `card-components.md` | `git diff $LAST_SYNC -- src/components/Galleries/GalleryCard.tsx src/components/Groups/GroupCard.tsx` |
| `navigation-settings.md` | `git diff $LAST_SYNC -- src/components/MainNavbar.tsx src/components/Scenes/SceneListTable.tsx src/components/Settings/` |
| `scene-player.md` | `git diff $LAST_SYNC -- src/components/ScenePlayer/ScenePlayer.tsx` |
| `scene-details.md` | `git diff $LAST_SYNC -- src/components/Scenes/SceneDetails/` |
| `detail-panels.md` | `git diff $LAST_SYNC -- src/components/Galleries/GalleryDetails/ src/components/Images/ImageDetails/` |
| `shared-components.md` | `git diff $LAST_SYNC -- src/components/Shared/` |
| `frontpage.md` | `git diff $LAST_SYNC -- src/components/FrontPage/` |
| `models-types.md` | `git diff $LAST_SYNC -- src/models/` |
| `utils.md` | `git diff $LAST_SYNC -- src/utils/` |
| `graphql-frontend.md` | `git diff $LAST_SYNC -- graphql/` |
| `miscellaneous.md` | `git diff $LAST_SYNC -- src/locales/` |

---

## Automated Check Script

Add this script to check patch status after merge:

```bash
#!/bin/bash
# save as: scripts/check-patches.sh

# Usage: ./check-patches.sh [COMPARE_REF]
# COMPARE_REF: The git ref to compare against (default: HEAD~1)
#   - After merge: HEAD~1 (compares to state before merge)
#   - Specific commit: abc123
#   - Previous sync tag: last-upstream-sync

COMPARE_REF=${1:-"HEAD~1"}

echo "Checking patched files against $COMPARE_REF..."
echo "(This compares your current state to the reference point)"
echo ""

# Files with patches
PATCHED_FILES=(
  "src/components/List/ListFilter.tsx"
  "src/components/List/ListTable.tsx"
  "src/components/List/Pagination.tsx"
  "src/components/Galleries/GalleryCard.tsx"
  "src/components/Groups/GroupCard.tsx"
  "src/components/MainNavbar.tsx"
  "src/components/Scenes/SceneListTable.tsx"
  "src/components/ScenePlayer/ScenePlayer.tsx"
  "src/components/Shared/ClearableInput.tsx"
  "src/components/Shared/CollapseButton.tsx"
  "src/components/Shared/DetailItem.tsx"
  "src/components/Shared/GridCard/GridCard.tsx"
  "src/components/Shared/Sidebar.tsx"
  "src/components/Shared/TagLink.tsx"
  "src/components/FrontPage/Control.tsx"
  "src/components/FrontPage/FrontPageConfig.tsx"
  "src/models/list-filter/types.ts"
  "src/models/sceneQueue.ts"
  "src/utils/caption.ts"
  "src/utils/screen.ts"
)

echo "Files changed that need patch review:"
echo "======================================"

for file in "${PATCHED_FILES[@]}"; do
  if git diff --quiet $COMPARE_REF -- "ui/v2.5/$file" 2>/dev/null; then
    : # No changes
  else
    echo "  ⚠️  $file"
  fi
done

echo ""
echo "Run 'git diff $COMPARE_REF -- ui/v2.5/<file>' to see changes"
echo ""
echo "Tip: Create a tag before each merge to make future comparisons easier:"
echo "  git tag last-upstream-sync"
```

---

## When Things Go Wrong

### Build Fails

1. Check for TypeScript errors - imports may have changed
2. Run `yarn generate` to update GraphQL types
3. Check patch files for updated diffs

### Tests Fail

1. Review test output for specific failures
2. Check if upstream changed component APIs
3. Update extension code if needed

### Runtime Errors

1. Check browser console for errors
2. Verify imports are resolving correctly
3. Check if GraphQL schema changed (regenerate types)

---

## Sync History

Track your upstream sync points here:

| Date | Upstream Commit/Ref | Notes |
|------|---------------------|-------|
| Dec 2024 | v0.29.3 | Initial fork baseline |

**Tip:** After each successful merge from upstream, add an entry and optionally tag it:

```bash
git tag upstream-sync-YYYY-MM-DD
```

This makes future `git diff` comparisons easier.


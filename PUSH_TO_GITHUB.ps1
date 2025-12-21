# Stashapp Branch Creation and Push Script
# Run this in PowerShell

cd C:\Users\Admin\Desktop\myrepo\stashapp

Write-Host "=== STEP 1: Current Git Status ===" -ForegroundColor Cyan
git status

Write-Host "`n=== STEP 2: Current Branch ===" -ForegroundColor Cyan
git branch

Write-Host "`n=== STEP 3: Remote Configuration ===" -ForegroundColor Cyan
git remote -v

Write-Host "`n=== STEP 4: Creating feature/scene-segments branch ===" -ForegroundColor Cyan
git checkout -b feature/scene-segments 2>&1

Write-Host "`n=== STEP 5: Adding scene segment files ===" -ForegroundColor Cyan
git add pkg/models/model_scene_segment.go
git add pkg/models/repository_scene_segment.go
git add pkg/sqlite/scene_segment.go
git add pkg/sqlite/migrations/76_scene_segments.up.sql
git add pkg/sqlite/migrations/76_scene_segments.down.sql
git add graphql/schema/types/scene-segment.graphql
git add internal/api/resolver_model_scene_segment.go
git add internal/api/resolver_mutation_scene_segment.go
git add internal/api/resolver_query_scene_segment.go
git add ui/v2.5/src/components/Scenes/SceneDetails/SceneSegmentForm.tsx
git add ui/v2.5/src/components/Scenes/SceneDetails/SceneSegmentsPanel.tsx

Write-Host "`n=== STEP 6: Checking what will be committed ===" -ForegroundColor Cyan
git status

Write-Host "`n=== STEP 7: Committing changes ===" -ForegroundColor Cyan
git commit -m "feat: implement scene segments feature for issue #3530"

Write-Host "`n=== STEP 8: Pushing to GitHub ===" -ForegroundColor Cyan
git push -u origin feature/scene-segments

Write-Host "`n=== DONE! ===" -ForegroundColor Green
Write-Host "Check https://github.com/SBALAVIGNESH123/stash/tree/feature/scene-segments" -ForegroundColor Yellow

# FAST FINISH - Copy & Run in NEW PowerShell

cd C:\Users\Admin\Desktop\myrepo\stashapp

# 1. Setup (30 sec)
if (!(Test-Path "ui\v2.5\build")) { mkdir ui\v2.5\build }
echo $null > ui\v2.5\build\index.html

# 2. Install gqlgen (10 sec)
go install github.com/99designs/gqlgen@latest

# 3. Generate GraphQL (20 sec)
go generate ./cmd/stash

# 4. Build backend (1 min)
go build ./cmd/stash

# 5. Install npm packages (3 min)
cd ui\v2.5
npm install

# 6. Build frontend (5 min)
npm run build

# 7. Done! 
cd ..\..
Write-Host "✅ BUILD COMPLETE! Run: .\stash.exe" -ForegroundColor Green
Write-Host "Then test & submit PR for $450!" -ForegroundColor Yellow

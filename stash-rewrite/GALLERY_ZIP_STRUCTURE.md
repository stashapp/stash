# Gallery Zip Infrastructure - Final Structure

## File Organization

```
Stash.Core/
  Entities/
    Gallery.cs              ← Gallery entity (existing)
    Image.cs                ← Image entity (existing)
    File.cs                 ← File entities including GalleryFile (existing)

  Galleries/                ← NEW: Gallery domain logic
    Zip/
      IZipFileReader.cs     ← Interface for zip operations
      ZipFileReader.cs      ← Implementation using System.IO.Compression
      ZipGalleryReader.cs   ← High-level gallery-specific operations
    GalleryServiceExtensions.cs  ← DI registration
```

## Why This Structure?

### Domain-Driven Design
- **Gallery entity** lives in `Stash.Core/Entities`
- **Gallery logic** lives in `Stash.Core/Galleries`
- Keeps related code together

### Clear Ownership
- All gallery-specific infrastructure in one namespace
- Easy to find: "Need gallery zip logic? → Stash.Core.Galleries.Zip"
- Room for expansion: parsing, serving, processing

### Matches Entity Pattern
- Images → `Stash.Core/Entities/Image.cs`
- Galleries → `Stash.Core/Entities/Gallery.cs` + `Stash.Core/Galleries/`
- Future: Videos could get `Stash.Core/Videos/` for transcoding logic

## Namespace Structure

```csharp
// Core entity
namespace Stash.Core.Entities;
public class Gallery { }

// Gallery infrastructure
namespace Stash.Core.Galleries.Zip;
public interface IZipFileReader { }
public class ZipFileReader : IZipFileReader { }
public class ZipGalleryReader { }

// Service registration
namespace Stash.Core.Galleries;
public static class GalleryServiceExtensions
{
    public static IServiceCollection AddGalleryServices(this IServiceCollection services)
    {
        services.AddSingleton<IZipFileReader, ZipFileReader>();
        services.AddSingleton<ZipGalleryReader>();
        return services;
    }
}
```

## Usage in Application

### Program.cs
```csharp
using Stash.Core.Galleries;

// Register gallery services
builder.Services.AddGalleryServices();
```

### ScanService.cs
```csharp
using Stash.Core.Galleries.Zip;

public class ScanService(
    ZipGalleryReader zipGalleryReader, // Injected
    // ... other dependencies
)
{
    private async Task ProcessGalleryFileAsync(StashContext db, string path, CancellationToken ct)
    {
        // Use zipGalleryReader to extract images
        var imageEntries = await zipGalleryReader.GetImageEntriesAsync(path, ct);

        // Create Image entities for each image in the zip
        foreach (var entry in imageEntries)
        {
            // Create ImageFile with ZipFileId linking to parent zip
            // Create Image entity
            // Link to Gallery via ImageGallery junction table
        }
    }
}
```

## Benefits of This Structure

### 1. **Cohesion**
- Gallery entity + gallery logic in same project (Stash.Core)
- No artificial separation between data and behavior

### 2. **Discoverability**
- Developer looking for gallery features checks `Stash.Core/Galleries`
- All gallery infrastructure in one place

### 3. **Extensibility**
```
Stash.Core/Galleries/
  Zip/              ← Archive handling
  Parsing/          ← Future: Gallery metadata extraction
  Serving/          ← Future: Image serving optimization
  Thumbnails/       ← Future: Gallery preview generation
```

### 4. **No Premature Abstraction**
- Didn't create separate `Stash.Media` or `Stash.Galleries` project
- Kept it simple: entity + supporting logic in Core
- Can extract to separate project later if needed

## Comparison to Alternatives

### ❌ Stash.Media/Zip (original)
- Problem: "Media" is too generic
- Problem: Zips are gallery-specific, not general media
- Problem: Separates gallery logic from Gallery entity

### ❌ Stash.Galleries (separate project)
- Problem: Overkill for current needs
- Problem: Extra project adds complexity
- Maybe later: If galleries grow significantly

### ✅ Stash.Core/Galleries (chosen)
- Benefit: Gallery logic with Gallery entity
- Benefit: Simple and discoverable
- Benefit: Room to grow without restructuring

## Integration Points

### Current Integration (Implemented)
1. **ScanService** uses `ZipGalleryReader` to extract images during scan
2. **Program.cs** registers services via `AddGalleryServices()`
3. **Gallery entity** has `ImageGalleries` collection populated by scan

### Next Steps (TODO)
1. **GalleryController** - Add endpoint to serve images by index
2. **Image dimension detection** - Extract width/height for gallery images
3. **Cover image selection** - Use `FindCoverImageIndexAsync()`
4. **Caching layer** - Cache extracted images for performance

## Migration from Old Structure

### What Changed
- Moved: `Stash.Media/Zip/*` → `Stash.Core/Galleries/Zip/*`
- Renamed: `AddMediaServices()` → `AddGalleryServices()`
- Namespace: `Stash.Media.Zip` → `Stash.Core.Galleries.Zip`
- Deleted: Entire `Stash.Media` project

### What Stayed the Same
- All class names unchanged
- All method signatures unchanged
- All documentation unchanged
- Only namespaces and file locations changed

## Testing the Structure

### Build Verification
```bash
cd /Users/mellownen/Documents/VSCode/Rewrite/stash/stash-rewrite
dotnet build
# Should succeed with 0 errors
```

### Runtime Verification
```bash
# After restarting backend:
# 1. Run a scan
# 2. Check that zip files create galleries
# 3. Verify images are extracted and linked
# 4. Check database for ImageGallery records
```

## Future Enhancements

### Short Term
- [ ] Image serving endpoint: `GET /api/galleries/{id}/images/{index}`
- [ ] Dimension detection for gallery images
- [ ] Gallery cover image selection

### Medium Term
- [ ] RAR/7-Zip support (using SharpCompress)
- [ ] Character encoding detection for legacy archives
- [ ] Gallery image caching layer

### Long Term
- [ ] `Stash.Core/Videos/` for video-specific logic (transcoding, etc.)
- [ ] `Stash.Core/Images/` for image-specific logic (thumbnails, etc.)
- [ ] Extract to `Stash.Galleries` project if domain grows significantly

using Microsoft.EntityFrameworkCore;
using Stash.Core.Entities;

namespace Stash.Data;

public class StashContext : DbContext
{
    public StashContext(DbContextOptions<StashContext> options) : base(options) { }

    // Core entities
    public DbSet<Scene> Scenes => Set<Scene>();
    public DbSet<Performer> Performers => Set<Performer>();
    public DbSet<Tag> Tags => Set<Tag>();
    public DbSet<Studio> Studios => Set<Studio>();
    public DbSet<Gallery> Galleries => Set<Gallery>();
    public DbSet<Image> Images => Set<Image>();
    public DbSet<Group> Groups => Set<Group>();
    public DbSet<SceneMarker> SceneMarkers => Set<SceneMarker>();
    public DbSet<SavedFilter> SavedFilters => Set<SavedFilter>();
    public DbSet<GalleryChapter> GalleryChapters => Set<GalleryChapter>();

    // Users
    public DbSet<User> Users => Set<User>();
    public DbSet<UserRole> UserRoles => Set<UserRole>();

    // Extensions
    public DbSet<ExtensionData> ExtensionData => Set<ExtensionData>();

    // Files & Folders
    public DbSet<Folder> Folders => Set<Folder>();
    public DbSet<VideoFile> VideoFiles => Set<VideoFile>();
    public DbSet<ImageFile> ImageFiles => Set<ImageFile>();
    public DbSet<GalleryFile> GalleryFiles => Set<GalleryFile>();
    public DbSet<FileFingerprint> FileFingerprints => Set<FileFingerprint>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Apply all configurations from this assembly
        modelBuilder.ApplyConfigurationsFromAssembly(typeof(StashContext).Assembly);

        // TPH for file hierarchy
        modelBuilder.Entity<BaseFileEntity>()
            .HasDiscriminator<string>("FileType")
            .HasValue<VideoFile>("Video")
            .HasValue<ImageFile>("Image")
            .HasValue<GalleryFile>("Gallery");

        modelBuilder.Entity<BaseFileEntity>()
            .HasMany(f => f.Fingerprints)
            .WithOne(fp => fp.File)
            .HasForeignKey(fp => fp.FileId)
            .OnDelete(DeleteBehavior.Cascade);
    }

    public override int SaveChanges()
    {
        UpdateTimestamps();
        return base.SaveChanges();
    }

    public override Task<int> SaveChangesAsync(CancellationToken cancellationToken = default)
    {
        UpdateTimestamps();
        return base.SaveChangesAsync(cancellationToken);
    }

    private void UpdateTimestamps()
    {
        var entries = ChangeTracker.Entries()
            .Where(e => e.State == EntityState.Added || e.State == EntityState.Modified);

        foreach (var entry in entries)
        {
            if (entry.Entity is BaseEntity entity)
            {
                entity.UpdatedAt = DateTime.UtcNow;
                if (entry.State == EntityState.Added)
                    entity.CreatedAt = DateTime.UtcNow;
            }
            else if (entry.Entity is BaseFileEntity file)
            {
                file.UpdatedAt = DateTime.UtcNow;
                if (entry.State == EntityState.Added)
                    file.CreatedAt = DateTime.UtcNow;
            }
            else if (entry.Entity is Folder folder)
            {
                folder.UpdatedAt = DateTime.UtcNow;
                if (entry.State == EntityState.Added)
                    folder.CreatedAt = DateTime.UtcNow;
            }
        }
    }
}

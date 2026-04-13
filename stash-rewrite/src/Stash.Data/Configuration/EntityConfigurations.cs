using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Stash.Core.Entities;

namespace Stash.Data.Configuration;

public class SceneConfiguration : IEntityTypeConfiguration<Scene>
{
    public void Configure(EntityTypeBuilder<Scene> builder)
    {
        builder.ToTable("scenes");
        builder.HasKey(s => s.Id);
        builder.Property(s => s.Rating).HasAnnotation("Range", new[] { 1, 100 });
        builder.Property(s => s.CustomFields).HasColumnType("jsonb");

        builder.HasOne(s => s.Studio)
            .WithMany(st => st.Scenes)
            .HasForeignKey(s => s.StudioId)
            .OnDelete(DeleteBehavior.SetNull);

        builder.HasMany(s => s.Urls).WithOne(u => u.Scene).HasForeignKey(u => u.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(s => s.Files).WithOne(f => f.Scene).HasForeignKey(f => f.SceneId).OnDelete(DeleteBehavior.SetNull);
        builder.HasMany(s => s.SceneMarkers).WithOne(m => m.Scene).HasForeignKey(m => m.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(s => s.StashIds).WithOne(si => si.Scene).HasForeignKey(si => si.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(s => s.PlayHistory).WithOne(h => h.Scene).HasForeignKey(h => h.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(s => s.OHistory).WithOne(h => h.Scene).HasForeignKey(h => h.SceneId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(s => s.Title);
        builder.HasIndex(s => s.StudioId);
        builder.HasIndex(s => s.Date);
        builder.HasIndex(s => s.Rating);
    }
}

public class SceneTagConfiguration : IEntityTypeConfiguration<SceneTag>
{
    public void Configure(EntityTypeBuilder<SceneTag> builder)
    {
        builder.ToTable("scene_tags");
        builder.HasKey(st => new { st.SceneId, st.TagId });
        builder.HasOne(st => st.Scene).WithMany(s => s.SceneTags).HasForeignKey(st => st.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(st => st.Tag).WithMany(t => t.SceneTags).HasForeignKey(st => st.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class ScenePerformerConfiguration : IEntityTypeConfiguration<ScenePerformer>
{
    public void Configure(EntityTypeBuilder<ScenePerformer> builder)
    {
        builder.ToTable("scene_performers");
        builder.HasKey(sp => new { sp.SceneId, sp.PerformerId });
        builder.HasOne(sp => sp.Scene).WithMany(s => s.ScenePerformers).HasForeignKey(sp => sp.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(sp => sp.Performer).WithMany(p => p.ScenePerformers).HasForeignKey(sp => sp.PerformerId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class SceneGalleryConfiguration : IEntityTypeConfiguration<SceneGallery>
{
    public void Configure(EntityTypeBuilder<SceneGallery> builder)
    {
        builder.ToTable("scene_galleries");
        builder.HasKey(sg => new { sg.SceneId, sg.GalleryId });
        builder.HasOne(sg => sg.Scene).WithMany(s => s.SceneGalleries).HasForeignKey(sg => sg.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(sg => sg.Gallery).WithMany(g => g.SceneGalleries).HasForeignKey(sg => sg.GalleryId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class SceneGroupConfiguration : IEntityTypeConfiguration<SceneGroup>
{
    public void Configure(EntityTypeBuilder<SceneGroup> builder)
    {
        builder.ToTable("scene_groups");
        builder.HasKey(sg => new { sg.SceneId, sg.GroupId });
        builder.HasOne(sg => sg.Scene).WithMany(s => s.SceneGroups).HasForeignKey(sg => sg.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(sg => sg.Group).WithMany(g => g.SceneGroups).HasForeignKey(sg => sg.GroupId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class PerformerConfiguration : IEntityTypeConfiguration<Performer>
{
    public void Configure(EntityTypeBuilder<Performer> builder)
    {
        builder.ToTable("performers");
        builder.HasKey(p => p.Id);
        builder.Property(p => p.Name).IsRequired().HasMaxLength(500);
        builder.Property(p => p.CustomFields).HasColumnType("jsonb");

        builder.HasMany(p => p.Urls).WithOne(u => u.Performer).HasForeignKey(u => u.PerformerId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(p => p.Aliases).WithOne(a => a.Performer).HasForeignKey(a => a.PerformerId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(p => p.StashIds).WithOne(si => si.Performer).HasForeignKey(si => si.PerformerId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(p => p.Name);
        builder.HasIndex(p => p.Favorite);
        builder.HasIndex(p => p.Rating);
    }
}

public class PerformerTagConfiguration : IEntityTypeConfiguration<PerformerTag>
{
    public void Configure(EntityTypeBuilder<PerformerTag> builder)
    {
        builder.ToTable("performer_tags");
        builder.HasKey(pt => new { pt.PerformerId, pt.TagId });
        builder.HasOne(pt => pt.Performer).WithMany(p => p.PerformerTags).HasForeignKey(pt => pt.PerformerId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(pt => pt.Tag).WithMany(t => t.PerformerTags).HasForeignKey(pt => pt.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class TagConfiguration : IEntityTypeConfiguration<Tag>
{
    public void Configure(EntityTypeBuilder<Tag> builder)
    {
        builder.ToTable("tags");
        builder.HasKey(t => t.Id);
        builder.Property(t => t.Name).IsRequired().HasMaxLength(500);
        builder.Property(t => t.CustomFields).HasColumnType("jsonb");

        builder.HasMany(t => t.Aliases).WithOne(a => a.Tag).HasForeignKey(a => a.TagId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(t => t.StashIds).WithOne(si => si.Tag).HasForeignKey(si => si.TagId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(t => t.Name).IsUnique();
        builder.HasIndex(t => t.Favorite);
    }
}

public class TagParentConfiguration : IEntityTypeConfiguration<TagParent>
{
    public void Configure(EntityTypeBuilder<TagParent> builder)
    {
        builder.ToTable("tag_parents");
        builder.HasKey(tp => new { tp.ParentId, tp.ChildId });
        builder.HasOne(tp => tp.Parent).WithMany(t => t.ChildRelations).HasForeignKey(tp => tp.ParentId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(tp => tp.Child).WithMany(t => t.ParentRelations).HasForeignKey(tp => tp.ChildId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class StudioConfiguration : IEntityTypeConfiguration<Studio>
{
    public void Configure(EntityTypeBuilder<Studio> builder)
    {
        builder.ToTable("studios");
        builder.HasKey(s => s.Id);
        builder.Property(s => s.Name).IsRequired().HasMaxLength(500);
        builder.Property(s => s.CustomFields).HasColumnType("jsonb");

        builder.HasOne(s => s.Parent).WithMany(s => s.Children).HasForeignKey(s => s.ParentId).OnDelete(DeleteBehavior.SetNull);
        builder.HasMany(s => s.Urls).WithOne(u => u.Studio).HasForeignKey(u => u.StudioId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(s => s.Aliases).WithOne(a => a.Studio).HasForeignKey(a => a.StudioId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(s => s.StashIds).WithOne(si => si.Studio).HasForeignKey(si => si.StudioId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(s => s.Name);
        builder.HasIndex(s => s.ParentId);
    }
}

public class StudioTagConfiguration : IEntityTypeConfiguration<StudioTag>
{
    public void Configure(EntityTypeBuilder<StudioTag> builder)
    {
        builder.ToTable("studio_tags");
        builder.HasKey(st => new { st.StudioId, st.TagId });
        builder.HasOne(st => st.Studio).WithMany(s => s.StudioTags).HasForeignKey(st => st.StudioId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(st => st.Tag).WithMany(t => t.StudioTags).HasForeignKey(st => st.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class GalleryConfiguration : IEntityTypeConfiguration<Gallery>
{
    public void Configure(EntityTypeBuilder<Gallery> builder)
    {
        builder.ToTable("galleries");
        builder.HasKey(g => g.Id);
        builder.Property(g => g.CustomFields).HasColumnType("jsonb");

        builder.HasOne(g => g.Studio).WithMany(s => s.Galleries).HasForeignKey(g => g.StudioId).OnDelete(DeleteBehavior.SetNull);
        builder.HasOne(g => g.Folder).WithMany().HasForeignKey(g => g.FolderId).OnDelete(DeleteBehavior.SetNull);
        builder.HasMany(g => g.Urls).WithOne(u => u.Gallery).HasForeignKey(u => u.GalleryId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(g => g.Files).WithOne(f => f.Gallery).HasForeignKey(f => f.GalleryId).OnDelete(DeleteBehavior.SetNull);
        builder.HasMany(g => g.Chapters).WithOne(c => c.Gallery).HasForeignKey(c => c.GalleryId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(g => g.Title);
        builder.HasIndex(g => g.StudioId);
    }
}

public class GalleryTagConfiguration : IEntityTypeConfiguration<GalleryTag>
{
    public void Configure(EntityTypeBuilder<GalleryTag> builder)
    {
        builder.ToTable("gallery_tags");
        builder.HasKey(gt => new { gt.GalleryId, gt.TagId });
        builder.HasOne(gt => gt.Gallery).WithMany(g => g.GalleryTags).HasForeignKey(gt => gt.GalleryId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(gt => gt.Tag).WithMany(t => t.GalleryTags).HasForeignKey(gt => gt.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class GalleryPerformerConfiguration : IEntityTypeConfiguration<GalleryPerformer>
{
    public void Configure(EntityTypeBuilder<GalleryPerformer> builder)
    {
        builder.ToTable("gallery_performers");
        builder.HasKey(gp => new { gp.GalleryId, gp.PerformerId });
        builder.HasOne(gp => gp.Gallery).WithMany(g => g.GalleryPerformers).HasForeignKey(gp => gp.GalleryId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(gp => gp.Performer).WithMany(p => p.GalleryPerformers).HasForeignKey(gp => gp.PerformerId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class ImageConfiguration : IEntityTypeConfiguration<Image>
{
    public void Configure(EntityTypeBuilder<Image> builder)
    {
        builder.ToTable("images");
        builder.HasKey(i => i.Id);
        builder.Property(i => i.CustomFields).HasColumnType("jsonb");

        builder.HasOne(i => i.Studio).WithMany(s => s.Images).HasForeignKey(i => i.StudioId).OnDelete(DeleteBehavior.SetNull);
        builder.HasMany(i => i.Urls).WithOne(u => u.Image).HasForeignKey(u => u.ImageId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(i => i.Files).WithOne(f => f.Image).HasForeignKey(f => f.ImageId).OnDelete(DeleteBehavior.SetNull);

        builder.HasIndex(i => i.Title);
        builder.HasIndex(i => i.StudioId);
    }
}

public class ImageTagConfiguration : IEntityTypeConfiguration<ImageTag>
{
    public void Configure(EntityTypeBuilder<ImageTag> builder)
    {
        builder.ToTable("image_tags");
        builder.HasKey(it => new { it.ImageId, it.TagId });
        builder.HasOne(it => it.Image).WithMany(i => i.ImageTags).HasForeignKey(it => it.ImageId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(it => it.Tag).WithMany(t => t.ImageTags).HasForeignKey(it => it.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class ImagePerformerConfiguration : IEntityTypeConfiguration<ImagePerformer>
{
    public void Configure(EntityTypeBuilder<ImagePerformer> builder)
    {
        builder.ToTable("image_performers");
        builder.HasKey(ip => new { ip.ImageId, ip.PerformerId });
        builder.HasOne(ip => ip.Image).WithMany(i => i.ImagePerformers).HasForeignKey(ip => ip.ImageId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(ip => ip.Performer).WithMany(p => p.ImagePerformers).HasForeignKey(ip => ip.PerformerId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class ImageGalleryConfiguration : IEntityTypeConfiguration<ImageGallery>
{
    public void Configure(EntityTypeBuilder<ImageGallery> builder)
    {
        builder.ToTable("image_galleries");
        builder.HasKey(ig => new { ig.ImageId, ig.GalleryId });
        builder.HasOne(ig => ig.Image).WithMany(i => i.ImageGalleries).HasForeignKey(ig => ig.ImageId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(ig => ig.Gallery).WithMany(g => g.ImageGalleries).HasForeignKey(ig => ig.GalleryId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class GroupConfiguration : IEntityTypeConfiguration<Group>
{
    public void Configure(EntityTypeBuilder<Group> builder)
    {
        builder.ToTable("groups");
        builder.HasKey(g => g.Id);
        builder.Property(g => g.Name).IsRequired().HasMaxLength(500);
        builder.Property(g => g.CustomFields).HasColumnType("jsonb");

        builder.HasOne(g => g.Studio).WithMany(s => s.Groups).HasForeignKey(g => g.StudioId).OnDelete(DeleteBehavior.SetNull);
        builder.HasMany(g => g.Urls).WithOne(u => u.Group).HasForeignKey(u => u.GroupId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(g => g.Name);
        builder.HasIndex(g => g.StudioId);
    }
}

public class GroupTagConfiguration : IEntityTypeConfiguration<GroupTag>
{
    public void Configure(EntityTypeBuilder<GroupTag> builder)
    {
        builder.ToTable("group_tags");
        builder.HasKey(gt => new { gt.GroupId, gt.TagId });
        builder.HasOne(gt => gt.Group).WithMany(g => g.GroupTags).HasForeignKey(gt => gt.GroupId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(gt => gt.Tag).WithMany(t => t.GroupTags).HasForeignKey(gt => gt.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class GroupRelationConfiguration : IEntityTypeConfiguration<GroupRelation>
{
    public void Configure(EntityTypeBuilder<GroupRelation> builder)
    {
        builder.ToTable("group_relations");
        builder.HasKey(gr => new { gr.ContainingGroupId, gr.SubGroupId });
        builder.HasOne(gr => gr.ContainingGroup).WithMany(g => g.SubGroupRelations).HasForeignKey(gr => gr.ContainingGroupId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(gr => gr.SubGroup).WithMany(g => g.ContainingGroupRelations).HasForeignKey(gr => gr.SubGroupId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class SceneMarkerConfiguration : IEntityTypeConfiguration<SceneMarker>
{
    public void Configure(EntityTypeBuilder<SceneMarker> builder)
    {
        builder.ToTable("scene_markers");
        builder.HasKey(m => m.Id);
        builder.HasOne(m => m.PrimaryTag).WithMany().HasForeignKey(m => m.PrimaryTagId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(m => m.Scene).WithMany(s => s.SceneMarkers).HasForeignKey(m => m.SceneId).OnDelete(DeleteBehavior.Cascade);
        builder.HasIndex(m => m.SceneId);
    }
}

public class SceneMarkerTagConfiguration : IEntityTypeConfiguration<SceneMarkerTag>
{
    public void Configure(EntityTypeBuilder<SceneMarkerTag> builder)
    {
        builder.ToTable("scene_marker_tags");
        builder.HasKey(smt => new { smt.SceneMarkerId, smt.TagId });
        builder.HasOne(smt => smt.SceneMarker).WithMany(m => m.SceneMarkerTags).HasForeignKey(smt => smt.SceneMarkerId).OnDelete(DeleteBehavior.Cascade);
        builder.HasOne(smt => smt.Tag).WithMany().HasForeignKey(smt => smt.TagId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class SavedFilterConfiguration : IEntityTypeConfiguration<SavedFilter>
{
    public void Configure(EntityTypeBuilder<SavedFilter> builder)
    {
        builder.ToTable("saved_filters");
        builder.HasKey(f => f.Id);
        builder.Property(f => f.FindFilter).HasColumnType("jsonb");
        builder.Property(f => f.ObjectFilter).HasColumnType("jsonb");
        builder.Property(f => f.UIOptions).HasColumnType("jsonb");
    }
}

public class GalleryChapterConfiguration : IEntityTypeConfiguration<GalleryChapter>
{
    public void Configure(EntityTypeBuilder<GalleryChapter> builder)
    {
        builder.ToTable("gallery_chapters");
        builder.HasKey(c => c.Id);
    }
}

public class UserConfiguration : IEntityTypeConfiguration<User>
{
    public void Configure(EntityTypeBuilder<User> builder)
    {
        builder.ToTable("users");
        builder.HasKey(u => u.Id);
        builder.Property(u => u.Username).IsRequired().HasMaxLength(200);
        builder.Property(u => u.PasswordHash).IsRequired().HasMaxLength(500);
        builder.Property(u => u.ApiKey).HasMaxLength(200);

        builder.HasIndex(u => u.Username).IsUnique();
        builder.HasIndex(u => u.ApiKey).IsUnique().HasFilter("\"ApiKey\" IS NOT NULL");

        builder.HasMany(u => u.Roles).WithOne(r => r.User).HasForeignKey(r => r.UserId).OnDelete(DeleteBehavior.Cascade);
    }
}

public class UserRoleConfiguration : IEntityTypeConfiguration<UserRole>
{
    public void Configure(EntityTypeBuilder<UserRole> builder)
    {
        builder.ToTable("user_roles");
        builder.HasKey(r => r.Id);
        builder.Property(r => r.Role).IsRequired().HasMaxLength(100);
        builder.HasIndex(r => new { r.UserId, r.Role }).IsUnique();
    }
}

public class FolderConfiguration : IEntityTypeConfiguration<Folder>
{
    public void Configure(EntityTypeBuilder<Folder> builder)
    {
        builder.ToTable("folders");
        builder.HasKey(f => f.Id);
        builder.Property(f => f.Path).IsRequired();
        builder.HasOne(f => f.ParentFolder).WithMany(f => f.SubFolders).HasForeignKey(f => f.ParentFolderId).OnDelete(DeleteBehavior.Cascade);
        builder.HasMany(f => f.Files).WithOne(file => file.ParentFolder).HasForeignKey(file => file.ParentFolderId).OnDelete(DeleteBehavior.Cascade);

        builder.HasIndex(f => f.Path).IsUnique();
    }
}

public class BaseFileEntityConfiguration : IEntityTypeConfiguration<BaseFileEntity>
{
    public void Configure(EntityTypeBuilder<BaseFileEntity> builder)
    {
        builder.ToTable("files");
        builder.HasKey(f => f.Id);
        builder.Property(f => f.Basename).IsRequired();
        builder.HasIndex(f => new { f.ParentFolderId, f.Basename }).IsUnique();
    }
}

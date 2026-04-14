using Microsoft.EntityFrameworkCore;
using Stash.Core.Entities;
using Stash.Core.Interfaces;

namespace Stash.Data.Repositories;

public class PerformerRepository : IPerformerRepository
{
    private readonly StashContext _db;
    public PerformerRepository(StashContext db) => _db = db;

    public async Task<Performer?> GetByIdAsync(int id, CancellationToken ct = default)
        => await _db.Performers.FindAsync([id], ct);

    public async Task<Performer?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Performers
            .Include(p => p.Urls)
            .Include(p => p.Aliases)
            .Include(p => p.PerformerTags).ThenInclude(pt => pt.Tag)
            .Include(p => p.StashIds)
            .AsSplitQuery()
            .FirstOrDefaultAsync(p => p.Id == id, ct);

    public async Task<IReadOnlyList<Performer>> GetAllAsync(CancellationToken ct = default)
        => await _db.Performers.AsNoTracking().ToListAsync(ct);

    public async Task<Performer> AddAsync(Performer entity, CancellationToken ct = default)
    {
        _db.Performers.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Performer entity, CancellationToken ct = default)
    {
        _db.Performers.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Performers.FindAsync([id], ct);
        if (entity != null)
        {
            _db.Performers.Remove(entity);
            await _db.SaveChangesAsync(ct);
        }
    }

    public async Task<int> CountAsync(CancellationToken ct = default)
        => await _db.Performers.CountAsync(ct);

    public async Task<(IReadOnlyList<Performer> Items, int TotalCount)> FindAsync(PerformerFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Performers
            .Include(p => p.PerformerTags).ThenInclude(pt => pt.Tag)
            .AsSplitQuery()
            .AsQueryable();

        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Name))
                query = query.Where(p => EF.Functions.ILike(p.Name, $"%{filter.Name}%"));
            if (filter.Favorite.HasValue)
                query = query.Where(p => p.Favorite == filter.Favorite.Value);
            if (filter.Rating.HasValue)
                query = query.Where(p => p.Rating >= filter.Rating.Value);
            if (filter.TagIds?.Count > 0)
                query = query.Where(p => p.PerformerTags.Any(pt => filter.TagIds.Contains(pt.TagId)));
        }

        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
        {
            var q = findFilter.Q;
            query = query.Where(p =>
                EF.Functions.ILike(p.Name, $"%{q}%") ||
                (p.Disambiguation != null && EF.Functions.ILike(p.Disambiguation, $"%{q}%")) ||
                p.Aliases.Any(a => EF.Functions.ILike(a.Alias, $"%{q}%")));
        }

        var totalCount = await query.CountAsync(ct);

        var sort = findFilter?.Sort ?? "name";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "name" => desc ? query.OrderByDescending(p => p.Name) : query.OrderBy(p => p.Name),
            "rating" => desc ? query.OrderByDescending(p => p.Rating) : query.OrderBy(p => p.Rating),
            "created_at" => desc ? query.OrderByDescending(p => p.CreatedAt) : query.OrderBy(p => p.CreatedAt),
            "birthdate" => desc ? query.OrderByDescending(p => p.Birthdate) : query.OrderBy(p => p.Birthdate),
            "scene_count" => desc
                ? query.OrderByDescending(p => p.ScenePerformers.Count)
                : query.OrderBy(p => p.ScenePerformers.Count),
            "image_count" => desc
                ? query.OrderByDescending(p => p.ImagePerformers.Count)
                : query.OrderBy(p => p.ImagePerformers.Count),
            "gallery_count" => desc
                ? query.OrderByDescending(p => p.GalleryPerformers.Count)
                : query.OrderBy(p => p.GalleryPerformers.Count),
            "height" => desc ? query.OrderByDescending(p => p.HeightCm) : query.OrderBy(p => p.HeightCm),
            "weight" => desc ? query.OrderByDescending(p => p.Weight) : query.OrderBy(p => p.Weight),
            "tag_count" => desc
                ? query.OrderByDescending(p => p.PerformerTags.Count)
                : query.OrderBy(p => p.PerformerTags.Count),
            "random" => query.OrderBy(_ => EF.Functions.Random()),
            _ => desc ? query.OrderByDescending(p => p.UpdatedAt) : query.OrderBy(p => p.UpdatedAt),
        };

        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);

        return (items, totalCount);
    }
}

public class TagRepository : ITagRepository
{
    private readonly StashContext _db;
    public TagRepository(StashContext db) => _db = db;

    public async Task<Tag?> GetByIdAsync(int id, CancellationToken ct = default)
        => await _db.Tags.FindAsync([id], ct);

    public async Task<Tag?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Tags
            .Include(t => t.Aliases)
            .Include(t => t.ParentRelations).ThenInclude(tp => tp.Parent)
            .Include(t => t.ChildRelations).ThenInclude(tp => tp.Child)
            .Include(t => t.StashIds)
            .AsSplitQuery()
            .FirstOrDefaultAsync(t => t.Id == id, ct);

    public async Task<Tag?> GetByNameAsync(string name, CancellationToken ct = default)
        => await _db.Tags.FirstOrDefaultAsync(t => t.Name == name, ct);

    public async Task<IReadOnlyList<Tag>> GetAllAsync(CancellationToken ct = default)
        => await _db.Tags.AsNoTracking().OrderBy(t => t.Name).ToListAsync(ct);

    public async Task<Tag> AddAsync(Tag entity, CancellationToken ct = default)
    {
        _db.Tags.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Tag entity, CancellationToken ct = default)
    {
        _db.Tags.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Tags.FindAsync([id], ct);
        if (entity != null)
        {
            _db.Tags.Remove(entity);
            await _db.SaveChangesAsync(ct);
        }
    }

    public async Task<int> CountAsync(CancellationToken ct = default)
        => await _db.Tags.CountAsync(ct);

    public async Task<(IReadOnlyList<Tag> Items, int TotalCount)> FindAsync(TagFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Tags
            .Include(t => t.Aliases)
            .AsQueryable();

        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Name))
                query = query.Where(t => EF.Functions.ILike(t.Name, $"%{filter.Name}%"));
            if (filter.Favorite.HasValue)
                query = query.Where(t => t.Favorite == filter.Favorite.Value);
        }

        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
        {
            var q = findFilter.Q;
            query = query.Where(t =>
                EF.Functions.ILike(t.Name, $"%{q}%") ||
                (t.Description != null && EF.Functions.ILike(t.Description, $"%{q}%")) ||
                t.Aliases.Any(a => EF.Functions.ILike(a.Alias, $"%{q}%")));
        }

        var totalCount = await query.CountAsync(ct);

        var sort = findFilter?.Sort ?? "name";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "name" => desc ? query.OrderByDescending(t => t.Name) : query.OrderBy(t => t.Name),
            _ => desc ? query.OrderByDescending(t => t.UpdatedAt) : query.OrderBy(t => t.UpdatedAt),
        };

        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);

        return (items, totalCount);
    }
}

public class StudioRepository : IStudioRepository
{
    private readonly StashContext _db;
    public StudioRepository(StashContext db) => _db = db;

    public async Task<Studio?> GetByIdAsync(int id, CancellationToken ct = default) => await _db.Studios.FindAsync([id], ct);

    public async Task<Studio?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Studios
            .Include(s => s.Parent).Include(s => s.Children)
            .Include(s => s.Urls).Include(s => s.Aliases)
            .Include(s => s.StudioTags).ThenInclude(st => st.Tag)
            .Include(s => s.StashIds)
            .AsSplitQuery()
            .FirstOrDefaultAsync(s => s.Id == id, ct);

    public async Task<IReadOnlyList<Studio>> GetAllAsync(CancellationToken ct = default)
        => await _db.Studios.AsNoTracking().OrderBy(s => s.Name).ToListAsync(ct);

    public async Task<Studio> AddAsync(Studio entity, CancellationToken ct = default)
    {
        _db.Studios.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Studio entity, CancellationToken ct = default)
    {
        _db.Studios.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Studios.FindAsync([id], ct);
        if (entity != null) { _db.Studios.Remove(entity); await _db.SaveChangesAsync(ct); }
    }

    public async Task<int> CountAsync(CancellationToken ct = default) => await _db.Studios.CountAsync(ct);

    public async Task<(IReadOnlyList<Studio> Items, int TotalCount)> FindAsync(StudioFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Studios.Include(s => s.StudioTags).ThenInclude(st => st.Tag).Include(s => s.StashIds).AsSplitQuery().AsQueryable();
        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Name)) query = query.Where(s => EF.Functions.ILike(s.Name, $"%{filter.Name}%"));
            if (filter.Favorite.HasValue) query = query.Where(s => s.Favorite == filter.Favorite.Value);
            if (filter.ParentId.HasValue) query = query.Where(s => s.ParentId == filter.ParentId.Value);
            if (filter.TagIds?.Count > 0) query = query.Where(s => s.StudioTags.Any(st => filter.TagIds.Contains(st.TagId)));
        }
        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
            query = query.Where(s => EF.Functions.ILike(s.Name, $"%{findFilter.Q}%"));

        var totalCount = await query.CountAsync(ct);
        var sort = findFilter?.Sort ?? "name";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "name" => desc ? query.OrderByDescending(s => s.Name) : query.OrderBy(s => s.Name),
            _ => desc ? query.OrderByDescending(s => s.UpdatedAt) : query.OrderBy(s => s.UpdatedAt),
        };
        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);
        return (items, totalCount);
    }
}

public class GalleryRepository : IGalleryRepository
{
    private readonly StashContext _db;
    public GalleryRepository(StashContext db) => _db = db;

    public async Task<Gallery?> GetByIdAsync(int id, CancellationToken ct = default) => await _db.Galleries.FindAsync([id], ct);

    public async Task<Gallery?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Galleries
            .Include(g => g.Studio).Include(g => g.Urls)
            .Include(g => g.GalleryTags).ThenInclude(gt => gt.Tag)
            .Include(g => g.GalleryPerformers).ThenInclude(gp => gp.Performer)
            .Include(g => g.Chapters)
            .Include(g => g.Files).ThenInclude(f => f.ParentFolder)
            .Include(g => g.Files).ThenInclude(f => f.Fingerprints)
            .Include(g => g.Folder)
            .Include(g => g.SceneGalleries)
            .AsSplitQuery()
            .FirstOrDefaultAsync(g => g.Id == id, ct);

    public async Task<IReadOnlyList<Gallery>> GetAllAsync(CancellationToken ct = default)
        => await _db.Galleries.AsNoTracking().ToListAsync(ct);

    public async Task<Gallery> AddAsync(Gallery entity, CancellationToken ct = default)
    {
        _db.Galleries.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Gallery entity, CancellationToken ct = default)
    {
        _db.Galleries.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Galleries.FindAsync([id], ct);
        if (entity != null) { _db.Galleries.Remove(entity); await _db.SaveChangesAsync(ct); }
    }

    public async Task<int> CountAsync(CancellationToken ct = default) => await _db.Galleries.CountAsync(ct);

    public async Task<(IReadOnlyList<Gallery> Items, int TotalCount)> FindAsync(GalleryFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Galleries.Include(g => g.GalleryTags).ThenInclude(gt => gt.Tag).AsSplitQuery().AsQueryable();
        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Title)) query = query.Where(g => g.Title != null && EF.Functions.ILike(g.Title, $"%{filter.Title}%"));
            if (filter.Organized.HasValue) query = query.Where(g => g.Organized == filter.Organized.Value);
            if (filter.StudioId.HasValue) query = query.Where(g => g.StudioId == filter.StudioId.Value);
            if (filter.TagIds?.Count > 0) query = query.Where(g => g.GalleryTags.Any(gt => filter.TagIds.Contains(gt.TagId)));
            if (filter.PerformerIds?.Count > 0) query = query.Where(g => g.GalleryPerformers.Any(gp => filter.PerformerIds.Contains(gp.PerformerId)));
        }
        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
            query = query.Where(g => (g.Title != null && EF.Functions.ILike(g.Title, $"%{findFilter.Q}%")));

        var totalCount = await query.CountAsync(ct);
        var sort = findFilter?.Sort ?? "updated_at";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "title" => desc ? query.OrderByDescending(g => g.Title) : query.OrderBy(g => g.Title),
            _ => desc ? query.OrderByDescending(g => g.UpdatedAt) : query.OrderBy(g => g.UpdatedAt),
        };
        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);
        return (items, totalCount);
    }
}

public class ImageRepository : IImageRepository
{
    private readonly StashContext _db;
    public ImageRepository(StashContext db) => _db = db;

    public async Task<Image?> GetByIdAsync(int id, CancellationToken ct = default) => await _db.Images.FindAsync([id], ct);

    public async Task<Image?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Images
            .Include(i => i.Studio).Include(i => i.Urls)
            .Include(i => i.ImageTags).ThenInclude(it => it.Tag)
            .Include(i => i.ImagePerformers).ThenInclude(ip => ip.Performer)
            .Include(i => i.ImageGalleries)
            .Include(i => i.Files)
            .AsSplitQuery()
            .FirstOrDefaultAsync(i => i.Id == id, ct);

    public async Task<IReadOnlyList<Image>> GetAllAsync(CancellationToken ct = default)
        => await _db.Images.AsNoTracking().ToListAsync(ct);

    public async Task<Image> AddAsync(Image entity, CancellationToken ct = default)
    {
        _db.Images.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Image entity, CancellationToken ct = default)
    {
        _db.Images.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Images.FindAsync([id], ct);
        if (entity != null) { _db.Images.Remove(entity); await _db.SaveChangesAsync(ct); }
    }

    public async Task<int> CountAsync(CancellationToken ct = default) => await _db.Images.CountAsync(ct);

    public async Task<(IReadOnlyList<Image> Items, int TotalCount)> FindAsync(ImageFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Images.Include(i => i.ImageTags).ThenInclude(it => it.Tag).AsSplitQuery().AsQueryable();
        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Title)) query = query.Where(i => i.Title != null && EF.Functions.ILike(i.Title, $"%{filter.Title}%"));
            if (filter.Organized.HasValue) query = query.Where(i => i.Organized == filter.Organized.Value);
            if (filter.StudioId.HasValue) query = query.Where(i => i.StudioId == filter.StudioId.Value);
            if (filter.GalleryId.HasValue) query = query.Where(i => i.ImageGalleries.Any(ig => ig.GalleryId == filter.GalleryId.Value));
            if (filter.TagIds?.Count > 0) query = query.Where(i => i.ImageTags.Any(it => filter.TagIds.Contains(it.TagId)));
            if (filter.PerformerIds?.Count > 0) query = query.Where(i => i.ImagePerformers.Any(ip => filter.PerformerIds.Contains(ip.PerformerId)));
        }
        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
            query = query.Where(i => (i.Title != null && EF.Functions.ILike(i.Title, $"%{findFilter.Q}%")));

        var totalCount = await query.CountAsync(ct);
        var sort = findFilter?.Sort ?? "updated_at";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "title" => desc ? query.OrderByDescending(i => i.Title) : query.OrderBy(i => i.Title),
            _ => desc ? query.OrderByDescending(i => i.UpdatedAt) : query.OrderBy(i => i.UpdatedAt),
        };
        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);
        return (items, totalCount);
    }
}

public class GroupRepository : IGroupRepository
{
    private readonly StashContext _db;
    public GroupRepository(StashContext db) => _db = db;

    public async Task<Group?> GetByIdAsync(int id, CancellationToken ct = default) => await _db.Groups.FindAsync([id], ct);

    public async Task<Group?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Groups
            .Include(g => g.Studio).Include(g => g.Urls)
            .Include(g => g.GroupTags).ThenInclude(gt => gt.Tag)
            .Include(g => g.SceneGroups).ThenInclude(sg => sg.Scene)
            .Include(g => g.SubGroupRelations)
            .Include(g => g.ContainingGroupRelations)
            .AsSplitQuery()
            .FirstOrDefaultAsync(g => g.Id == id, ct);

    public async Task<IReadOnlyList<Group>> GetAllAsync(CancellationToken ct = default)
        => await _db.Groups.AsNoTracking().OrderBy(g => g.Name).ToListAsync(ct);

    public async Task<Group> AddAsync(Group entity, CancellationToken ct = default)
    {
        _db.Groups.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Group entity, CancellationToken ct = default)
    {
        _db.Groups.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Groups.FindAsync([id], ct);
        if (entity != null) { _db.Groups.Remove(entity); await _db.SaveChangesAsync(ct); }
    }

    public async Task<int> CountAsync(CancellationToken ct = default) => await _db.Groups.CountAsync(ct);

    public async Task<(IReadOnlyList<Group> Items, int TotalCount)> FindAsync(GroupFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Groups.Include(g => g.GroupTags).ThenInclude(gt => gt.Tag).AsSplitQuery().AsQueryable();
        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Name)) query = query.Where(g => EF.Functions.ILike(g.Name, $"%{filter.Name}%"));
            if (filter.StudioId.HasValue) query = query.Where(g => g.StudioId == filter.StudioId.Value);
        }
        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
            query = query.Where(g => EF.Functions.ILike(g.Name, $"%{findFilter.Q}%"));

        var totalCount = await query.CountAsync(ct);
        var sort = findFilter?.Sort ?? "name";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "name" => desc ? query.OrderByDescending(g => g.Name) : query.OrderBy(g => g.Name),
            _ => desc ? query.OrderByDescending(g => g.UpdatedAt) : query.OrderBy(g => g.UpdatedAt),
        };
        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);
        return (items, totalCount);
    }
}

public class SavedFilterRepository : ISavedFilterRepository
{
    private readonly StashContext _db;
    public SavedFilterRepository(StashContext db) => _db = db;

    public async Task<SavedFilter?> GetByIdAsync(int id, CancellationToken ct = default) => await _db.SavedFilters.FindAsync([id], ct);
    public async Task<IReadOnlyList<SavedFilter>> GetAllAsync(CancellationToken ct = default) => await _db.SavedFilters.AsNoTracking().ToListAsync(ct);
    public async Task<IReadOnlyList<SavedFilter>> GetByModeAsync(Core.Enums.FilterMode mode, CancellationToken ct = default)
        => await _db.SavedFilters.Where(f => f.Mode == mode).AsNoTracking().ToListAsync(ct);

    public async Task<SavedFilter> AddAsync(SavedFilter entity, CancellationToken ct = default)
    {
        _db.SavedFilters.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(SavedFilter entity, CancellationToken ct = default) { _db.SavedFilters.Update(entity); await _db.SaveChangesAsync(ct); }
    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.SavedFilters.FindAsync([id], ct);
        if (entity != null) { _db.SavedFilters.Remove(entity); await _db.SaveChangesAsync(ct); }
    }
    public async Task<int> CountAsync(CancellationToken ct = default) => await _db.SavedFilters.CountAsync(ct);
}

public class SceneMarkerRepository : ISceneMarkerRepository
{
    private readonly StashContext _db;
    public SceneMarkerRepository(StashContext db) => _db = db;

    public async Task<SceneMarker?> GetByIdAsync(int id, CancellationToken ct = default) => await _db.SceneMarkers.FindAsync([id], ct);
    public async Task<IReadOnlyList<SceneMarker>> GetAllAsync(CancellationToken ct = default) => await _db.SceneMarkers.AsNoTracking().ToListAsync(ct);
    public async Task<IReadOnlyList<SceneMarker>> GetBySceneIdAsync(int sceneId, CancellationToken ct = default)
        => await _db.SceneMarkers.Include(m => m.PrimaryTag).Where(m => m.SceneId == sceneId).AsNoTracking().ToListAsync(ct);

    public async Task<SceneMarker> AddAsync(SceneMarker entity, CancellationToken ct = default)
    {
        _db.SceneMarkers.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(SceneMarker entity, CancellationToken ct = default) { _db.SceneMarkers.Update(entity); await _db.SaveChangesAsync(ct); }
    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.SceneMarkers.FindAsync([id], ct);
        if (entity != null) { _db.SceneMarkers.Remove(entity); await _db.SaveChangesAsync(ct); }
    }
    public async Task<int> CountAsync(CancellationToken ct = default) => await _db.SceneMarkers.CountAsync(ct);
}

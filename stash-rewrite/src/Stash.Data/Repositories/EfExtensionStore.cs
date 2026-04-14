using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Stash.Core.Entities;
using Stash.Plugins;

namespace Stash.Data.Repositories;

/// <summary>
/// EF Core implementation of IExtensionStore, scoped to a single extension ID.
/// </summary>
public class EfExtensionStore : IExtensionStore
{
    private readonly StashContext _db;
    private readonly string _extensionId;

    public EfExtensionStore(StashContext db, string extensionId)
    {
        _db = db;
        _extensionId = extensionId;
    }

    public async Task<string?> GetAsync(string key, CancellationToken ct = default)
    {
        var entry = await _db.ExtensionData
            .FirstOrDefaultAsync(e => e.ExtensionId == _extensionId && e.Key == key, ct);
        return entry?.Value;
    }

    public async Task SetAsync(string key, string value, CancellationToken ct = default)
    {
        var entry = await _db.ExtensionData
            .FirstOrDefaultAsync(e => e.ExtensionId == _extensionId && e.Key == key, ct);

        if (entry is not null)
        {
            entry.Value = value;
            entry.UpdatedAt = DateTime.UtcNow;
        }
        else
        {
            _db.ExtensionData.Add(new ExtensionData
            {
                ExtensionId = _extensionId,
                Key = key,
                Value = value
            });
        }
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(string key, CancellationToken ct = default)
    {
        var entry = await _db.ExtensionData
            .FirstOrDefaultAsync(e => e.ExtensionId == _extensionId && e.Key == key, ct);
        if (entry is not null)
        {
            _db.ExtensionData.Remove(entry);
            await _db.SaveChangesAsync(ct);
        }
    }

    public async Task<Dictionary<string, string>> GetAllAsync(CancellationToken ct = default)
    {
        return await _db.ExtensionData
            .Where(e => e.ExtensionId == _extensionId)
            .ToDictionaryAsync(e => e.Key, e => e.Value, ct);
    }
}

/// <summary>
/// Factory that creates scoped IExtensionStore instances for each extension.
/// </summary>
public class EfExtensionStoreFactory : IExtensionStoreFactory
{
    private readonly IServiceProvider _services;

    public EfExtensionStoreFactory(IServiceProvider services)
    {
        _services = services;
    }

    public IExtensionStore CreateStore(string extensionId)
    {
        var scope = _services.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();
        return new EfExtensionStore(db, extensionId);
    }
}

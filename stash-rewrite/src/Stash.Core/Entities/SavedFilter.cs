using Stash.Core.Enums;

namespace Stash.Core.Entities;

public class SavedFilter : BaseEntity
{
    public FilterMode Mode { get; set; }
    public string Name { get; set; } = string.Empty;
    public string? FindFilter { get; set; } // JSON
    public string? ObjectFilter { get; set; } // JSON
    public string? UIOptions { get; set; } // JSON
}

namespace Stash.Core.Entities;

/// <summary>
/// Key-value storage for extension state, scoped per extension.
/// </summary>
public class ExtensionData
{
    public required string ExtensionId { get; set; }
    public required string Key { get; set; }
    public required string Value { get; set; }
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
}

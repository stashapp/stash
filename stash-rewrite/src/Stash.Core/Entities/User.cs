namespace Stash.Core.Entities;

public class User : BaseEntity
{
    public string Username { get; set; } = string.Empty;
    public string PasswordHash { get; set; } = string.Empty;
    public string? ApiKey { get; set; }
    public bool IsAdmin { get; set; }
    public bool IsEnabled { get; set; } = true;
    public DateTime? LastLoginAt { get; set; }
    public ICollection<UserRole> Roles { get; set; } = [];
}

/// <summary>Named roles assigned to users to control access.</summary>
public class UserRole
{
    public int Id { get; set; }
    public int UserId { get; set; }
    public string Role { get; set; } = string.Empty;
    public User? User { get; set; }
}

/// <summary>Well-known permission strings checked by middleware and UI.</summary>
public static class Permissions
{
    public const string ViewScenes = "scenes:read";
    public const string EditScenes = "scenes:write";
    public const string DeleteScenes = "scenes:delete";

    public const string ViewImages = "images:read";
    public const string EditImages = "images:write";
    public const string DeleteImages = "images:delete";

    public const string ViewPerformers = "performers:read";
    public const string EditPerformers = "performers:write";

    public const string ViewStudios = "studios:read";
    public const string EditStudios = "studios:write";

    public const string ViewTags = "tags:read";
    public const string EditTags = "tags:write";

    public const string RunJobs = "jobs:run";
    public const string ViewLogs = "logs:read";
    public const string ManageConfig = "config:write";
    public const string ManageUsers = "users:manage";
}

/// <summary>Predefined role definitions mapping role names to permissions.</summary>
public static class Roles
{
    public const string Admin = "admin";
    public const string Editor = "editor";
    public const string Viewer = "viewer";

    public static readonly IReadOnlyDictionary<string, string[]> PermissionMap =
        new Dictionary<string, string[]>
        {
            [Admin] = [
                Permissions.ViewScenes, Permissions.EditScenes, Permissions.DeleteScenes,
                Permissions.ViewImages, Permissions.EditImages, Permissions.DeleteImages,
                Permissions.ViewPerformers, Permissions.EditPerformers,
                Permissions.ViewStudios, Permissions.EditStudios,
                Permissions.ViewTags, Permissions.EditTags,
                Permissions.RunJobs, Permissions.ViewLogs,
                Permissions.ManageConfig, Permissions.ManageUsers,
            ],
            [Editor] = [
                Permissions.ViewScenes, Permissions.EditScenes,
                Permissions.ViewImages, Permissions.EditImages,
                Permissions.ViewPerformers, Permissions.EditPerformers,
                Permissions.ViewStudios, Permissions.EditStudios,
                Permissions.ViewTags, Permissions.EditTags,
                Permissions.RunJobs, Permissions.ViewLogs,
            ],
            [Viewer] = [
                Permissions.ViewScenes, Permissions.ViewImages,
                Permissions.ViewPerformers, Permissions.ViewStudios,
                Permissions.ViewTags,
            ],
        };

    public static IEnumerable<string> GetPermissions(string role) =>
        PermissionMap.TryGetValue(role, out var perms) ? perms : [];
}

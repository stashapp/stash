using Microsoft.Extensions.DependencyInjection;
using Stash.Core.Entities.Galleries.Zip;

namespace Stash.Core.Entities.Galleries;

/// <summary>
/// Extension methods for registering gallery-related services in the DI container.
/// </summary>
public static class GalleryServiceExtensions
{
    /// <summary>
    /// Registers gallery infrastructure services (zip reading, parsing, etc.).
    /// Call this from Program.cs to add gallery services to the application.
    /// </summary>
    /// <param name="services">Service collection</param>
    /// <returns>Service collection for chaining</returns>
    public static IServiceCollection AddGalleryServices(this IServiceCollection services)
    {
        // Register zip file reading services
        // Singleton is appropriate because these services are stateless
        // and can be safely shared across all requests
        services.AddSingleton<IZipFileReader, ZipFileReader>();
        services.AddSingleton<ZipGalleryReader>();

        // TODO: Add other gallery services here as they're implemented:
        // - Cover image detection
        // - Gallery metadata extraction
        // - Image dimension detection
        // - Thumbnail generation for gallery previews

        return services;
    }
}

using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Stash.Core.Interfaces;
using Stash.Data.Repositories;

namespace Stash.Data;

public static class DataServiceExtensions
{
    public static IServiceCollection AddStashData(this IServiceCollection services, string connectionString)
    {
        services.AddDbContext<StashContext>(options =>
        {
            options.UseNpgsql(connectionString, npgsqlOptions =>
            {
                npgsqlOptions.UseVector();
                npgsqlOptions.MigrationsAssembly(typeof(StashContext).Assembly.FullName);
            });
        });

        services.AddScoped<ISceneRepository, SceneRepository>();
        services.AddScoped<IPerformerRepository, PerformerRepository>();
        services.AddScoped<ITagRepository, TagRepository>();
        services.AddScoped<IStudioRepository, StudioRepository>();
        services.AddScoped<IGalleryRepository, GalleryRepository>();
        services.AddScoped<IImageRepository, ImageRepository>();
        services.AddScoped<IGroupRepository, GroupRepository>();
        services.AddScoped<ISavedFilterRepository, SavedFilterRepository>();
        services.AddScoped<ISceneMarkerRepository, SceneMarkerRepository>();

        return services;
    }
}

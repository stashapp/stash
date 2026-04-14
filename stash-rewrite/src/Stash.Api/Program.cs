using System.IO.Compression;
using System.Text;
using System.Text.Json;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.ResponseCompression;
using Microsoft.AspNetCore.SignalR;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Serilog;
using Stash.Api.Hubs;
using Stash.Api.Services;
using Stash.Core.Events;
using Stash.Core.Interfaces;
using Stash.Data;
using Stash.Plugins;

// Ensure enough threads for async I/O under concurrent load
ThreadPool.SetMinThreads(Environment.ProcessorCount * 4, Environment.ProcessorCount * 4);

Log.Logger = new LoggerConfiguration()
    .WriteTo.Console()
    .CreateBootstrapLogger();

try
{
    var builder = WebApplication.CreateBuilder(args);

    // Serilog
    builder.Host.UseSerilog((context, services, configuration) => configuration
        .ReadFrom.Configuration(context.Configuration)
        .ReadFrom.Services(services)
        .Enrich.FromLogContext()
        .MinimumLevel.Override("Microsoft", Serilog.Events.LogEventLevel.Warning)
        .MinimumLevel.Override("System", Serilog.Events.LogEventLevel.Warning)
        .WriteTo.Console()
        .WriteTo.Sink(new SignalRLogSink()));

    // Bind configuration
    var stashConfig = builder.Configuration.GetSection("Stash");
    builder.Services.Configure<StashConfiguration>(stashConfig);
    builder.Services.Configure<AuthConfig>(stashConfig.GetSection("Auth"));
    builder.Services.Configure<PostgresConfig>(stashConfig.GetSection("Postgres"));

    // Register a singleton StashConfiguration instance so all consumers share the same mutable object
    var stashCfgInstance = stashConfig.Get<StashConfiguration>() ?? new StashConfiguration();
    builder.Services.AddSingleton(stashCfgInstance);

    // Database - EF Core + PostgreSQL
    var pgSection = stashConfig.GetSection("Postgres");
    var connectionString = pgSection.GetValue<string>("ConnectionString");
    if (string.IsNullOrEmpty(connectionString))
    {
        // Build from individual settings (managed or external)
        var pgPort = pgSection.GetValue<int?>("Port") ?? 5433;
        var pgDb = pgSection.GetValue<string>("Database") ?? "stash";
        connectionString = $"Host=127.0.0.1;Port={pgPort};Database={pgDb};Username=postgres;Trust Server Certificate=true;Minimum Pool Size=10;Maximum Pool Size=200;Timeout=15;Command Timeout=30";
    }
    builder.Services.AddStashData(connectionString);

    // Event bus (singleton for cross-service communication)
    builder.Services.AddSingleton<IEventBus, EventBus>();

    // Job service (background task processing)
    builder.Services.AddSingleton<JobService>();
    builder.Services.AddSingleton<IJobService>(sp => sp.GetRequiredService<JobService>());
    builder.Services.AddHostedService(sp => sp.GetRequiredService<JobService>());

    // Application services
    builder.Services.AddSingleton<IThumbnailService, ThumbnailService>();
    builder.Services.AddSingleton<IFingerprintService, FingerprintService>();
    builder.Services.AddScoped<IScanService, ScanService>();
    builder.Services.AddScoped<IStreamService, StreamService>();
    builder.Services.AddScoped<IAutoTagService, AutoTagService>();
    builder.Services.AddScoped<ICleanService, CleanService>();
    builder.Services.AddScoped<IBackupService, BackupService>();
    builder.Services.AddSingleton<IBlobService, BlobService>();
    builder.Services.AddSingleton<ConfigService>();
    builder.Services.AddSingleton<ScraperService>();
    builder.Services.AddHttpClient<StashBoxService>();

    // Extension system
    var extensionsDataDir = Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
        "stash", "extensions");
    Directory.CreateDirectory(extensionsDataDir);
    var extensionContext = new ExtensionContext
    {
        Configuration = builder.Configuration,
        DataDirectory = extensionsDataDir
    };
    var extensionManager = new ExtensionManager(extensionContext);
    // Discover .NET plugin DLLs from extensions directory
    extensionManager.DiscoverExtensions(extensionsDataDir);
    // Register built-in extensions (POC demonstrations)
    extensionManager.Register(new Stash.Api.Extensions.ThemeCollectionExtension());
    extensionManager.Register(new Stash.Api.Extensions.SceneAnalyticsExtension());
    extensionManager.Register(new Stash.Api.Extensions.CustomHomeExtension());
    extensionManager.Register(new Stash.Api.Extensions.SystemToolsExtension());
    extensionManager.Register(new Stash.Api.Extensions.NotificationSettingsExtension());
    extensionManager.Register(new Stash.Api.Extensions.EnhancedDeleteDialogExtension());
    extensionManager.Register(new Stash.Api.Extensions.AuditLogExtension());
    builder.Services.AddSingleton(extensionManager);
    builder.Services.AddSingleton<IExtensionStoreFactory>(sp => new Stash.Data.Repositories.EfExtensionStoreFactory(sp));
    extensionManager.ConfigureServices(builder.Services);

    // Managed PostgreSQL — auto-downloads and runs a local PG instance
    var pgManaged = pgSection.GetValue<bool?>("Managed") ?? true;
    if (pgManaged)
        builder.Services.AddHostedService<PostgresManagerService>();

    // SignalR
    builder.Services.AddSignalR();

    // Auth
    var authConfig = stashConfig.GetSection("Auth");
    var jwtSecret = authConfig.GetValue<string>("JwtSecret") ?? Guid.NewGuid().ToString();
    var authEnabled = authConfig.GetValue<bool>("Enabled");

    builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
        .AddJwtBearer(options =>
        {
            options.TokenValidationParameters = new TokenValidationParameters
            {
                ValidateIssuer = true,
                ValidateAudience = true,
                ValidateLifetime = true,
                ValidateIssuerSigningKey = true,
                ValidIssuer = "Stash",
                ValidAudience = "Stash",
                IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtSecret))
            };
            // Allow SignalR to authenticate via query string
            options.Events = new JwtBearerEvents
            {
                OnMessageReceived = context =>
                {
                    var accessToken = context.Request.Query["access_token"];
                    var path = context.HttpContext.Request.Path;
                    if (!string.IsNullOrEmpty(accessToken) && path.StartsWithSegments("/hubs"))
                        context.Token = accessToken;
                    return Task.CompletedTask;
                }
            };
        });
    builder.Services.AddAuthorization();

    // MVC + OpenAPI
    builder.Services.AddControllers()
        .AddJsonOptions(options =>
        {
            options.JsonSerializerOptions.Converters.Add(new System.Text.Json.Serialization.JsonStringEnumConverter(JsonNamingPolicy.CamelCase));
        });
    builder.Services.AddOpenApi();
    builder.Services.AddEndpointsApiExplorer();
    builder.Services.AddSwaggerGen();

    // Response compression — reduces 22KB scene lists to ~2KB
    builder.Services.AddResponseCompression(options =>
    {
        options.EnableForHttps = true;
        options.Providers.Add<BrotliCompressionProvider>();
        options.Providers.Add<GzipCompressionProvider>();
    });
    builder.Services.Configure<BrotliCompressionProviderOptions>(options => options.Level = CompressionLevel.Fastest);
    builder.Services.Configure<GzipCompressionProviderOptions>(options => options.Level = CompressionLevel.Fastest);

    // Output caching for read-heavy API endpoints
    builder.Services.AddOutputCache(options =>
    {
        options.AddBasePolicy(b => b.NoCache());
        options.AddPolicy("ShortCache", b => b.Expire(TimeSpan.FromSeconds(1)).SetVaryByQuery("*").SetLocking(false));
    });

    // In-memory cache for POST query results
    builder.Services.AddMemoryCache();

    // CORS - allow frontend dev server
    builder.Services.AddCors(options =>
    {
        options.AddDefaultPolicy(policy =>
        {
            policy.WithOrigins("http://localhost:5173", "http://localhost:3000")
                .AllowAnyHeader()
                .AllowAnyMethod()
                .AllowCredentials();
        });
    });

    var app = builder.Build();

    // Middleware pipeline
    // UseSerilogRequestLogging removed — adds 3-5ms per request overhead

    if (app.Environment.IsDevelopment())
    {
        app.MapOpenApi();
        app.UseSwagger();
        app.UseSwaggerUI();
    }

    app.UseResponseCompression();
    app.UseCors();
    app.UseOutputCache();

    if (authEnabled)
    {
        app.UseAuthentication();
        app.UseAuthorization();
    }

    app.MapControllers();
    app.MapHub<JobHub>("/hubs/jobs");
    app.MapHub<LogHub>("/hubs/logs");
    extensionManager.MapEndpoints(app);

    // Serve SPA static files (production)
    app.UseDefaultFiles();
    app.UseStaticFiles();
    app.MapFallbackToFile("index.html");

    var port = stashConfig.GetValue<int?>("Port") ?? 9999;
    app.Urls.Add($"http://0.0.0.0:{port}");

    // Initialize SignalR log sink with hub context
    SignalRLogSink.SetHubContext(app.Services.GetRequiredService<IHubContext<LogHub>>());

    // Start hosted services (including managed PostgreSQL) before database creation.
    await app.StartAsync();

    // Load saved user config (stash-config.json) and apply on top of appsettings.json
    var configSvc = app.Services.GetRequiredService<ConfigService>();
    var savedConfig = await configSvc.LoadSavedConfigAsync();
    if (savedConfig != null)
    {
        await configSvc.SaveConfigAsync(savedConfig); // applies to live IOptions
        Log.Information("Loaded user configuration from {Path}", configSvc.ConfigPath);
    }

    // Auto-migrate database + pre-warm EF Core and connection pool
    using (var scope = app.Services.CreateScope())
    {
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();
        await db.Database.EnsureCreatedAsync();

        // Pre-warm: compile EF Core query cache, prime connection pool, JIT hot paths
        _ = await db.Scenes.CountAsync();
        _ = await db.Scenes.AsNoTracking()
            .Include(s => s.Files).ThenInclude(f => f.Fingerprints)
            .Include(s => s.SceneTags).ThenInclude(st => st.Tag)
            .Include(s => s.ScenePerformers).ThenInclude(sp => sp.Performer)
            .Take(1).AsSingleQuery().ToListAsync();
        Log.Information("EF Core and connection pool pre-warmed");
    }

    // Initialize extensions after database is ready
    await extensionManager.InitializeAllAsync(app.Services);

    Log.Information("Stash starting on port {Port}", port);
    await app.WaitForShutdownAsync();

    // Graceful shutdown for extensions
    await extensionManager.ShutdownAllAsync();
}
catch (Exception ex)
{
    Log.Fatal(ex, "Application terminated unexpectedly");
}
finally
{
    await Log.CloseAndFlushAsync();
}

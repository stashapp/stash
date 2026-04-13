using System.Text;
using System.Text.Json;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.SignalR;
using Microsoft.IdentityModel.Tokens;
using Serilog;
using Stash.Api.Hubs;
using Stash.Api.Services;
using Stash.Core.Events;
using Stash.Core.Interfaces;
using Stash.Data;
using Stash.Plugins;

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
        connectionString = $"Host=127.0.0.1;Port={pgPort};Database={pgDb};Username=postgres;Trust Server Certificate=true";
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
    builder.Services.AddSingleton(extensionManager);
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
    app.UseSerilogRequestLogging();

    if (app.Environment.IsDevelopment())
    {
        app.MapOpenApi();
        app.UseSwagger();
        app.UseSwaggerUI();
    }

    app.UseCors();

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

    // Auto-migrate database
    using (var scope = app.Services.CreateScope())
    {
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();
        await db.Database.EnsureCreatedAsync();
    }

    // Initialize extensions after database is ready
    await extensionManager.InitializeAllAsync(app.Services);

    Log.Information("Stash starting on port {Port}", port);
    await app.WaitForShutdownAsync();
}
catch (Exception ex)
{
    Log.Fatal(ex, "Application terminated unexpectedly");
}
finally
{
    await Log.CloseAndFlushAsync();
}

using Microsoft.AspNetCore.Mvc;
using Stash.Api.Services;
using Stash.Core.DTOs;
using Stash.Core.Interfaces;
using Stash.Plugins;
using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class PluginsController(
    ExtensionManager extensionManager,
    IJobService jobService,
    StashConfiguration config,
    ConfigService configService,
    ILogger<PluginsController> logger) : ControllerBase
{
    [HttpGet]
    public ActionResult<List<PluginDto>> ListPlugins()
    {
        var manifest = extensionManager.GetAggregatedManifest();
        var plugins = extensionManager.GetAllExtensions()
            .Select(ext => new PluginDto(
                ext.Id,
                ext.Name,
                ext.Description,
                ext.Version,
                !config.DisabledPlugins.Contains(ext.Id),
                GetPluginTasks(ext)
            ))
            .ToList();

        // Also scan for Python plugins
        foreach (var pluginDir in GetPluginDirectories())
        {
            var pyPlugins = DiscoverPythonPlugins(pluginDir);
            plugins.AddRange(pyPlugins);
        }

        return Ok(plugins);
    }

    [HttpGet("tasks")]
    public ActionResult<List<PluginTaskDto>> ListTasks()
    {
        var tasks = new List<PluginTaskDto>();
        foreach (var ext in extensionManager.GetAllExtensions())
            tasks.AddRange(GetPluginTasks(ext));

        // Python plugin tasks
        foreach (var pluginDir in GetPluginDirectories())
        {
            foreach (var plugin in DiscoverPythonPlugins(pluginDir))
                tasks.AddRange(plugin.Tasks);
        }

        return Ok(tasks);
    }

    [HttpPost("run-task")]
    public ActionResult<object> RunTask([FromBody] RunPluginTaskDto dto)
    {
        // Check if it's a Python plugin
        foreach (var pluginDir in GetPluginDirectories())
        {
            var configPath = Path.Combine(pluginDir, dto.PluginId, "plugin.yml");
            if (!System.IO.File.Exists(configPath))
                configPath = Path.Combine(pluginDir, dto.PluginId, "plugin.yaml");
            if (!System.IO.File.Exists(configPath)) continue;

            var jobId = jobService.Enqueue($"plugin:{dto.PluginId}", $"Running {dto.PluginId}/{dto.TaskName}", async (progress, ct) =>
            {
                progress.Report(0, $"Starting plugin task {dto.TaskName}...");
                var scriptDir = Path.Combine(pluginDir, dto.PluginId);
                var entryPoint = FindPythonEntryPoint(scriptDir);

                if (entryPoint != null)
                {
                    await RunPythonPluginAsync(entryPoint, dto.TaskName, dto.Args, ct);
                }
                else
                {
                    logger.LogWarning("No entry point found for plugin {PluginId}", dto.PluginId);
                }

                progress.Report(1, "Done");
            });

            return Ok(new { jobId });
        }

        // Check .NET extensions
        var ext = extensionManager.GetAllExtensions().FirstOrDefault(e => e.Id == dto.PluginId);
        if (ext == null) return NotFound($"Plugin '{dto.PluginId}' not found");

        var netJobId = jobService.Enqueue($"plugin:{dto.PluginId}", $"Running {dto.PluginId}/{dto.TaskName}", async (progress, ct) =>
        {
            progress.Report(0, "Running plugin task...");
            // .NET plugin tasks would be invoked via reflection or IPluginTask interface
            await Task.CompletedTask;
            progress.Report(1, "Done");
        });

        return Ok(new { jobId = netJobId });
    }

    [HttpPost("settings")]
    public async Task<IActionResult> UpdateSettings([FromBody] PluginSettingsDto dto)
    {
        foreach (var (pluginId, enabled) in dto.EnabledMap)
        {
            if (enabled)
            {
                extensionManager.EnableExtension(pluginId);
                config.DisabledPlugins.Remove(pluginId);
            }
            else
            {
                extensionManager.DisableExtension(pluginId);
                config.DisabledPlugins.Add(pluginId);
            }
        }

        await configService.SaveCurrentConfigAsync();
        return Ok();
    }

    [HttpGet("{pluginId}/config")]
    public ActionResult<Dictionary<string, object?>> GetPluginConfig(string pluginId)
    {
        if (config.PluginConfigurations.TryGetValue(pluginId, out var cfg))
            return Ok(cfg);
        return Ok(new Dictionary<string, object?>());
    }

    [HttpPost("{pluginId}/config")]
    public async Task<IActionResult> SetPluginConfig(string pluginId, [FromBody] Dictionary<string, object?> values)
    {
        config.PluginConfigurations[pluginId] = values;
        await configService.SaveCurrentConfigAsync();
        return Ok();
    }

    [HttpPost("reload")]
    public async Task<IActionResult> ReloadPlugins()
    {
        await extensionManager.ReloadAllAsync();
        return Ok(new { message = "Plugins reloaded" });
    }

    // ===== Package Management =====

    [HttpGet("packages/installed")]
    public ActionResult<List<PackageDto>> GetInstalledPackages([FromQuery] string? type)
    {
        var packages = new List<PackageDto>();
        foreach (var pluginDir in GetPluginDirectories())
        {
            if (!Directory.Exists(pluginDir)) continue;
            foreach (var dir in Directory.GetDirectories(pluginDir))
            {
                var name = Path.GetFileName(dir);
                packages.Add(new PackageDto(name, "", "local", "", type ?? "plugin", true, "local"));
            }
        }
        return Ok(packages);
    }

    [HttpGet("packages/available")]
    public ActionResult<List<PackageDto>> GetAvailablePackages([FromQuery] string? type, [FromQuery] string? source)
    {
        // In a full implementation, this would query package sources
        return Ok(new List<PackageDto>());
    }

    [HttpPost("packages/install")]
    public ActionResult<object> InstallPackages([FromBody] InstallPackagesDto dto)
    {
        var jobId = jobService.Enqueue("install-packages", "Installing packages", async (progress, ct) =>
        {
            for (var i = 0; i < dto.Packages.Count; i++)
            {
                var pkg = dto.Packages[i];
                progress.Report((double)(i + 1) / dto.Packages.Count, $"Installing {pkg.Id}...");
                logger.LogInformation("Installing package {Id} from {Source}", pkg.Id, pkg.SourceUrl);
                // In a full implementation, would download and extract the package
                await Task.Delay(100, ct); // Placeholder
            }
        });

        return Ok(new { jobId });
    }

    [HttpPost("packages/uninstall")]
    public IActionResult UninstallPackages([FromBody] List<string> packageIds)
    {
        foreach (var id in packageIds)
        {
            foreach (var pluginDir in GetPluginDirectories())
            {
                var dir = Path.Combine(pluginDir, id);
                if (Directory.Exists(dir))
                {
                    Directory.Delete(dir, true);
                    logger.LogInformation("Uninstalled package {Id}", id);
                }
            }
        }
        return Ok(new { uninstalled = packageIds.Count });
    }

    // ===== Helpers =====

    private List<string> GetPluginDirectories()
    {
        var dirs = new List<string>();
        var defaultDir = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
            "stash", "plugins");
        dirs.Add(defaultDir);

        if (config.ExtensionPaths?.Count > 0)
            dirs.AddRange(config.ExtensionPaths);

        return dirs;
    }

    private static List<PluginTaskDto> GetPluginTasks(IExtension ext)
    {
        // Default task for any extension
        return [new PluginTaskDto("run", $"Run {ext.Name}")];
    }

    private List<PluginDto> DiscoverPythonPlugins(string pluginDir)
    {
        var plugins = new List<PluginDto>();
        if (!Directory.Exists(pluginDir)) return plugins;

        var deserializer = new DeserializerBuilder()
            .WithNamingConvention(CamelCaseNamingConvention.Instance)
            .IgnoreUnmatchedProperties()
            .Build();

        foreach (var dir in Directory.GetDirectories(pluginDir))
        {
            var configPath = Path.Combine(dir, "plugin.yml");
            if (!System.IO.File.Exists(configPath))
                configPath = Path.Combine(dir, "plugin.yaml");
            if (!System.IO.File.Exists(configPath)) continue;

            var id = Path.GetFileName(dir);
            var entryPoint = FindPythonEntryPoint(dir);
            if (entryPoint == null) continue;

            var name = id;
            var description = $"Python plugin: {id}";
            var version = "1.0.0";
            string? url = null;
            var tasks = new List<PluginTaskDto> { new("run", $"Run {id}") };
            var settingsSchema = new List<PluginSettingSchemaDto>();

            try
            {
                var yamlText = System.IO.File.ReadAllText(configPath);
                var parsed = deserializer.Deserialize<Dictionary<string, object?>>(yamlText);
                if (parsed != null)
                {
                    if (parsed.TryGetValue("name", out var n) && n is string ns) name = ns;
                    if (parsed.TryGetValue("description", out var d) && d is string ds) description = ds;
                    if (parsed.TryGetValue("version", out var v) && v != null) version = v.ToString()!;
                    if (parsed.TryGetValue("url", out var u) && u is string us) url = us;

                    // Parse tasks
                    if (parsed.TryGetValue("exec", out var exec) && exec is Dictionary<object, object> execDict)
                    {
                        tasks.Clear();
                        if (execDict.TryGetValue("tasks", out var tasksObj) && tasksObj is List<object> tasksList)
                        {
                            foreach (var t in tasksList)
                            {
                                if (t is Dictionary<object, object> td)
                                {
                                    var taskName = td.TryGetValue("name", out var tn) ? tn?.ToString() ?? "run" : "run";
                                    var taskDesc = td.TryGetValue("description", out var tdesc) ? tdesc?.ToString() ?? "" : "";
                                    tasks.Add(new PluginTaskDto(taskName, taskDesc));
                                }
                            }
                        }
                    }

                    // Parse settings schema
                    if (parsed.TryGetValue("settings", out var settings) && settings is Dictionary<object, object> settingsDict)
                    {
                        foreach (var (key, val) in settingsDict)
                        {
                            var settingName = key.ToString()!;
                            var settingType = "STRING";
                            string? displayName = null;
                            string? settingDescription = null;

                            if (val is Dictionary<object, object> sd)
                            {
                                if (sd.TryGetValue("type", out var st)) settingType = st?.ToString()?.ToUpperInvariant() ?? "STRING";
                                if (sd.TryGetValue("displayName", out var dn)) displayName = dn?.ToString();
                                if (sd.TryGetValue("description", out var sdesc)) settingDescription = sdesc?.ToString();
                            }

                            settingsSchema.Add(new PluginSettingSchemaDto(settingName, settingType, displayName, settingDescription));
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                logger.LogWarning(ex, "Failed to parse plugin.yml for {PluginId}", id);
            }

            plugins.Add(new PluginDto(
                id,
                name,
                description,
                version,
                !config.DisabledPlugins.Contains(id),
                tasks,
                settingsSchema.Count > 0 ? settingsSchema : null,
                url
            ));
        }

        return plugins;
    }

    private static string? FindPythonEntryPoint(string dir)
    {
        var mainPy = Path.Combine(dir, "main.py");
        if (System.IO.File.Exists(mainPy)) return mainPy;

        var initPy = Path.Combine(dir, "__init__.py");
        if (System.IO.File.Exists(initPy)) return initPy;

        return Directory.GetFiles(dir, "*.py").FirstOrDefault();
    }

    private async Task RunPythonPluginAsync(string scriptPath, string taskName, Dictionary<string, string>? args, CancellationToken ct)
    {
        var pythonPath = FindPython();
        if (pythonPath == null)
        {
            logger.LogError("Python not found in PATH");
            return;
        }

        using var process = new System.Diagnostics.Process();
        process.StartInfo = new System.Diagnostics.ProcessStartInfo
        {
            FileName = pythonPath,
            Arguments = $"\"{scriptPath}\"",
            WorkingDirectory = Path.GetDirectoryName(scriptPath),
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
            CreateNoWindow = true,
        };

        // Pass task name and args as environment variables
        process.StartInfo.EnvironmentVariables["STASH_TASK"] = taskName;
        if (args != null)
        {
            foreach (var (key, value) in args)
                process.StartInfo.EnvironmentVariables[$"STASH_ARG_{key.ToUpperInvariant()}"] = value;
        }

        process.Start();
        var output = await process.StandardOutput.ReadToEndAsync(ct);
        var error = await process.StandardError.ReadToEndAsync(ct);
        await process.WaitForExitAsync(ct);

        if (process.ExitCode != 0)
            logger.LogWarning("Python plugin {Script} exited with code {Code}: {Error}", scriptPath, process.ExitCode, error);
        else if (!string.IsNullOrEmpty(output))
            logger.LogInformation("Python plugin output: {Output}", output.TrimEnd());
    }

    private static string? FindPython()
    {
        var names = new[] { "python3", "python" };
        foreach (var name in names)
        {
            var pathEnv = Environment.GetEnvironmentVariable("PATH") ?? "";
            foreach (var dir in pathEnv.Split(Path.PathSeparator))
            {
                var exts = OperatingSystem.IsWindows() ? new[] { ".exe", ".cmd", ".bat" } : new[] { "" };
                foreach (var ext in exts)
                {
                    var fullPath = Path.Combine(dir, name + ext);
                    if (System.IO.File.Exists(fullPath)) return fullPath;
                }
            }
        }
        return null;
    }
}

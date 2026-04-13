using Stash.Core.DTOs;
using Stash.Core.Interfaces;
using YamlDotNet.Serialization;

namespace Stash.Api.Services;

public class ScraperService
{
    private static readonly string[] SupportedExtensions = [".yml", ".yaml"];

    private readonly StashConfiguration _config;
    private readonly ILogger<ScraperService> _logger;
    private readonly IDeserializer _deserializer;
    private readonly Lock _sync = new();
    private IReadOnlyList<ScraperSummaryDto> _cached = [];

    public ScraperService(StashConfiguration config, ILogger<ScraperService> logger)
    {
        _config = config;
        _logger = logger;
        _deserializer = new DeserializerBuilder()
            .IgnoreUnmatchedProperties()
            .Build();
    }

    public IReadOnlyList<ScraperSummaryDto> GetScrapers()
    {
        lock (_sync)
        {
            if (_cached.Count == 0)
                _cached = LoadScrapers();

            return _cached;
        }
    }

    public IReadOnlyList<ScraperSummaryDto> ReloadScrapers()
    {
        lock (_sync)
        {
            _cached = LoadScrapers();
            return _cached;
        }
    }

    private IReadOnlyList<ScraperSummaryDto> LoadScrapers()
    {
        var summaries = new List<ScraperSummaryDto>();
        var seenFiles = new HashSet<string>(StringComparer.OrdinalIgnoreCase);

        foreach (var directory in _config.Scraping.ScraperDirectories.Where(path => !string.IsNullOrWhiteSpace(path)))
        {
            if (!Directory.Exists(directory))
                continue;

            IEnumerable<string> files;
            try
            {
                files = Directory.EnumerateFiles(directory, "*.*", SearchOption.AllDirectories)
                    .Where(file => SupportedExtensions.Contains(Path.GetExtension(file), StringComparer.OrdinalIgnoreCase));
            }
            catch (Exception ex)
            {
                _logger.LogWarning(ex, "Failed to enumerate scraper directory {Directory}", directory);
                continue;
            }

            foreach (var file in files)
            {
                if (!seenFiles.Add(file))
                    continue;

                try
                {
                    summaries.AddRange(ParseScraperFile(file));
                }
                catch (Exception ex)
                {
                    _logger.LogWarning(ex, "Failed to load scraper definition from {File}", file);
                }
            }
        }

        return summaries
            .OrderBy(summary => summary.Name, StringComparer.OrdinalIgnoreCase)
            .ThenBy(summary => summary.EntityType, StringComparer.OrdinalIgnoreCase)
            .ToList();
    }

    private IReadOnlyList<ScraperSummaryDto> ParseScraperFile(string file)
    {
        using var stream = File.OpenRead(file);
        using var reader = new StreamReader(stream);
        var definition = _deserializer.Deserialize<ScraperManifest>(reader);

        var scraperId = Path.GetFileNameWithoutExtension(file);
        var scraperName = string.IsNullOrWhiteSpace(definition.Name)
            ? scraperId
            : definition.Name.Trim();
        var summaries = new List<ScraperSummaryDto>();

        AddSummary(
            summaries,
            scraperId,
            scraperName,
            "scene",
            file,
            byName: definition.SceneByName,
            byFragments: [definition.SceneByFragment, definition.SceneByQueryFragment],
            byUrls: definition.SceneByUrl
        );
        AddSummary(
            summaries,
            scraperId,
            scraperName,
            "performer",
            file,
            byName: definition.PerformerByName,
            byFragments: [definition.PerformerByFragment],
            byUrls: definition.PerformerByUrl
        );
        AddSummary(
            summaries,
            scraperId,
            scraperName,
            "gallery",
            file,
            byFragments: [definition.GalleryByFragment],
            byUrls: definition.GalleryByUrl
        );
        AddSummary(
            summaries,
            scraperId,
            scraperName,
            "image",
            file,
            byFragments: [definition.ImageByFragment],
            byUrls: definition.ImageByUrl
        );
        AddSummary(
            summaries,
            scraperId,
            scraperName,
            "group",
            file,
            byUrls: [.. definition.GroupByUrl, .. definition.MovieByUrl]
        );

        return summaries;
    }

    private static void AddSummary(
        ICollection<ScraperSummaryDto> summaries,
        string scraperId,
        string scraperName,
        string entityType,
        string file,
        ByNameDefinition? byName = null,
        IEnumerable<ByFragmentDefinition?>? byFragments = null,
        IEnumerable<ByUrlDefinition>? byUrls = null)
    {
        var supportedScrapes = new List<string>();
        var urls = new HashSet<string>(StringComparer.OrdinalIgnoreCase);

        if (byName != null)
            supportedScrapes.Add("Name");

        if (byFragments?.Any(definition => definition != null) == true)
            supportedScrapes.Add("Fragment");

        if (byUrls?.Any() == true)
        {
            supportedScrapes.Add("URL");
            foreach (var url in byUrls.SelectMany(definition => definition.Url ?? []))
            {
                if (!string.IsNullOrWhiteSpace(url))
                    urls.Add(url.Trim());
            }
        }

        if (supportedScrapes.Count == 0)
            return;

        summaries.Add(new ScraperSummaryDto(
            Id: $"{scraperId}:{entityType}",
            Name: scraperName,
            EntityType: entityType,
            SupportedScrapes: supportedScrapes,
            Urls: urls.OrderBy(url => url, StringComparer.OrdinalIgnoreCase).ToList(),
            SourcePath: file
        ));
    }

    private sealed class ScraperManifest
    {
        [YamlMember(Alias = "name")]
        public string? Name { get; init; }

        [YamlMember(Alias = "performerByName")]
        public ByNameDefinition? PerformerByName { get; init; }

        [YamlMember(Alias = "performerByFragment")]
        public ByFragmentDefinition? PerformerByFragment { get; init; }

        [YamlMember(Alias = "performerByURL")]
        public List<ByUrlDefinition> PerformerByUrl { get; init; } = [];

        [YamlMember(Alias = "sceneByName")]
        public ByNameDefinition? SceneByName { get; init; }

        [YamlMember(Alias = "sceneByFragment")]
        public ByFragmentDefinition? SceneByFragment { get; init; }

        [YamlMember(Alias = "sceneByQueryFragment")]
        public ByFragmentDefinition? SceneByQueryFragment { get; init; }

        [YamlMember(Alias = "sceneByURL")]
        public List<ByUrlDefinition> SceneByUrl { get; init; } = [];

        [YamlMember(Alias = "galleryByFragment")]
        public ByFragmentDefinition? GalleryByFragment { get; init; }

        [YamlMember(Alias = "galleryByURL")]
        public List<ByUrlDefinition> GalleryByUrl { get; init; } = [];

        [YamlMember(Alias = "imageByFragment")]
        public ByFragmentDefinition? ImageByFragment { get; init; }

        [YamlMember(Alias = "imageByURL")]
        public List<ByUrlDefinition> ImageByUrl { get; init; } = [];

        [YamlMember(Alias = "groupByURL")]
        public List<ByUrlDefinition> GroupByUrl { get; init; } = [];

        [YamlMember(Alias = "movieByURL")]
        public List<ByUrlDefinition> MovieByUrl { get; init; } = [];
    }

    private sealed class ByNameDefinition
    {
        [YamlMember(Alias = "queryURL")]
        public string? QueryUrl { get; init; }
    }

    private sealed class ByFragmentDefinition
    {
        [YamlMember(Alias = "queryURL")]
        public string? QueryUrl { get; init; }
    }

    private sealed class ByUrlDefinition
    {
        [YamlMember(Alias = "url")]
        public List<string> Url { get; init; } = [];

        [YamlMember(Alias = "queryURL")]
        public string? QueryUrl { get; init; }
    }
}
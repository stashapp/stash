using System.Globalization;
using System.Net.Http.Json;
using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.EntityFrameworkCore;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Enums;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public class StashBoxService
{
    private const string PerformerFragment = """
fragment PerformerFields on Performer {
  id
  name
  disambiguation
  aliases
  gender
  deleted
  merged_into_id
  urls {
    url
  }
  images {
    url
  }
  birth_date
  death_date
  ethnicity
  country
  eye_color
  hair_color
  height
  measurements {
    band_size
    cup_size
    waist
    hip
  }
  breast_type
  career_start_year
  career_end_year
  tattoos {
    location
    description
  }
  piercings {
    location
    description
  }
}
""";

    private const string SearchPerformerQuery = """
query SearchPerformer($term: String!) {
  searchPerformer(term: $term) {
    ... PerformerFields
  }
}
""" + PerformerFragment;

    private const string FindPerformerByIdQuery = """
query FindPerformerByID($id: ID!) {
  findPerformer(id: $id) {
    ... PerformerFields
  }
}
""" + PerformerFragment;

        private const string SearchStudioQuery = """
query SearchStudio($term: String!) {
  searchStudio(term: $term) {
    ... StudioFields
  }
}
""" + StudioFragment;

        private const string FindStudioByIdQuery = """
query FindStudioByID($id: ID!) {
  findStudio(id: $id) {
    ... StudioFields
  }
}
""" + StudioFragment;

        private const string StudioFragment = """
fragment StudioFields on Studio {
    id
    name
    aliases
    urls {
        url
    }
    images {
        url
    }
    parent {
        id
        name
    }
}
""";

        private const string TagFragment = """
fragment TagFields on Tag {
    id
    name
    description
    aliases
}
""";

        private const string FingerprintFragment = """
fragment FingerprintFields on Fingerprint {
    algorithm
    hash
    duration
}
""";

        private const string SceneFragment = """
fragment SceneFields on Scene {
    id
    title
    code
    details
    director
    duration
    date
    urls {
        url
    }
    images {
        url
    }
    studio {
        ... StudioFields
    }
    tags {
        ... TagFields
    }
    performers {
        performer {
            ... PerformerFields
        }
    }
    fingerprints {
        ... FingerprintFields
    }
}
""" + StudioFragment + TagFragment + FingerprintFragment + PerformerFragment;

        private const string SearchSceneQuery = """
query SearchScene($term: String!) {
    searchScene(term: $term) {
        ... SceneFields
    }
}
""" + SceneFragment;

        private const string FindSceneByIdQuery = """
query FindSceneByID($id: ID!) {
    findScene(id: $id) {
        ... SceneFields
    }
}
""" + SceneFragment;

        private const string FindScenesByFingerprintsQuery = """
query FindScenesBySceneFingerprints($fingerprints: [[FingerprintQueryInput!]!]!) {
    findScenesBySceneFingerprints(fingerprints: $fingerprints) {
        ... SceneFields
    }
}
""" + SceneFragment;

    private const string MeQuery = """
query Me {
  me {
    name
  }
}
""";

    private readonly HttpClient _httpClient;
    private readonly StashConfiguration _config;
    private readonly StashContext _db;
    private readonly IBlobService _blobService;
    private readonly ILogger<StashBoxService> _logger;
    private readonly JsonSerializerOptions _jsonOptions = new()
    {
        PropertyNameCaseInsensitive = true,
        DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
    };

    public StashBoxService(HttpClient httpClient, StashConfiguration config, StashContext db, IBlobService blobService, ILogger<StashBoxService> logger)
    {
        _httpClient = httpClient;
        _config = config;
        _db = db;
        _blobService = blobService;
        _logger = logger;
    }

    public async Task<StashBoxValidationResultDto> ValidateAsync(StashBoxDto input, CancellationToken ct)
    {
        var box = ToConfigBox(input);

        try
        {
            var response = await SendQueryAsync<StashBoxMeQueryResponse>(box, MeQuery, null, ct);
            var username = response.Me?.Name?.Trim();
            if (!string.IsNullOrWhiteSpace(username))
            {
                return new StashBoxValidationResultDto(true, $"Successfully authenticated as {username}", username);
            }

            return new StashBoxValidationResultDto(false, "Invalid or expired API key.", null);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to validate stash-box endpoint {Endpoint}", box.Endpoint);
            return new StashBoxValidationResultDto(false, MapValidationError(ex), null);
        }
    }

    public async Task<IReadOnlyList<StashBoxPerformerMatchDto>> SearchPerformersAsync(string term, string? endpoint, CancellationToken ct)
    {
        if (string.IsNullOrWhiteSpace(term))
            return [];

        var boxes = ResolveBoxes(endpoint);
        var results = new List<StashBoxPerformerMatchDto>();
        var strictEndpoint = !string.IsNullOrWhiteSpace(endpoint);

        foreach (var box in boxes)
        {
            try
            {
                var response = await SendQueryAsync<StashBoxSearchPerformerResponse>(box, SearchPerformerQuery, new { term }, ct);
                results.AddRange(response.SearchPerformer.Select(remote => ToMatchDto(box, remote)));
            }
            catch (Exception ex) when (!strictEndpoint)
            {
                _logger.LogWarning(ex, "Skipping stash-box performer search for {Endpoint}", box.Endpoint);
            }
        }

        return results
            .OrderByDescending(match => string.Equals(match.Name, term, StringComparison.OrdinalIgnoreCase))
            .ThenBy(match => match.Deleted)
            .ThenBy(match => match.Name, StringComparer.OrdinalIgnoreCase)
            .ThenBy(match => match.StashBoxName, StringComparer.OrdinalIgnoreCase)
            .ToList();
    }

    public async Task<StashBoxPerformerMatchDto?> GetPerformerMatchAsync(string endpoint, string performerId, CancellationToken ct)
    {
        var box = ResolveBox(endpoint);
        var performer = await GetRemotePerformerAsync(box, performerId, ct);
        if (performer == null)
            return null;

        if (!string.IsNullOrWhiteSpace(performer.MergedIntoId))
        {
            var merged = await GetRemotePerformerAsync(box, performer.MergedIntoId, ct);
            if (merged != null)
                performer = merged;
        }

        return ToMatchDto(box, performer);
    }

    public async Task<bool> MergePerformerAsync(Performer performer, string endpoint, string performerId, CancellationToken ct)
    {
        var box = ResolveBox(endpoint);
        var remote = await GetRemotePerformerAsync(box, performerId, ct);
        if (remote == null)
            return false;

        if (!string.IsNullOrWhiteSpace(remote.MergedIntoId))
        {
            var merged = await GetRemotePerformerAsync(box, remote.MergedIntoId, ct);
            if (merged != null)
                remote = merged;
        }

        ApplyRemotePerformer(performer, box.Endpoint, remote);
        await DownloadPerformerImageAsync(performer, remote, ct);
        return true;
    }

    // ===== Studio Stash-Box Methods =====

    public async Task<IReadOnlyList<StashBoxStudioMatchDto>> SearchStudiosAsync(string term, string? endpoint, CancellationToken ct)
    {
        if (string.IsNullOrWhiteSpace(term))
            return [];

        var boxes = ResolveBoxes(endpoint);
        var results = new List<StashBoxStudioMatchDto>();
        var strictEndpoint = !string.IsNullOrWhiteSpace(endpoint);

        foreach (var box in boxes)
        {
            try
            {
                var response = await SendQueryAsync<StashBoxSearchStudioResponse>(box, SearchStudioQuery, new { term }, ct);
                results.AddRange(response.SearchStudio.Select(remote => ToStudioMatchDto(box, remote)));
            }
            catch (Exception ex) when (!strictEndpoint)
            {
                _logger.LogWarning(ex, "Skipping stash-box studio search for {Endpoint}", box.Endpoint);
            }
        }

        return results
            .OrderByDescending(m => string.Equals(m.Name, term, StringComparison.OrdinalIgnoreCase))
            .ThenBy(m => m.Name, StringComparer.OrdinalIgnoreCase)
            .ThenBy(m => m.StashBoxName, StringComparer.OrdinalIgnoreCase)
            .ToList();
    }

    public async Task<bool> MergeStudioAsync(Studio studio, string endpoint, string studioId, CancellationToken ct)
    {
        var box = ResolveBox(endpoint);
        var remote = await GetRemoteStudioAsync(box, studioId, ct);
        if (remote == null)
            return false;

        studio.Name = remote.Name.Trim();
        MergeAliases(studio, remote.Aliases);
        MergeUrls(studio, remote.Urls.Select(u => u.Url));
        UpsertStashId(studio.StashIds, box.Endpoint, remote.Id, id => id.Endpoint, id => id.StashId, (id, value) => id.StashId = value, value => new StudioStashId { Endpoint = box.Endpoint, StashId = value });
        await DownloadStudioImageAsync(studio, remote, ct);

        // Resolve parent studio
        if (remote.Parent != null && studio.ParentId == null)
        {
            var parent = await _db.Studios
                .Include(s => s.StashIds)
                .FirstOrDefaultAsync(s => s.StashIds.Any(id => id.Endpoint == box.Endpoint && id.StashId == remote.Parent.Id), ct)
                ?? await _db.Studios
                    .Include(s => s.StashIds)
                    .FirstOrDefaultAsync(s => s.Name == remote.Parent.Name, ct);

            if (parent == null)
            {
                parent = new Studio { Name = remote.Parent.Name };
                parent.StashIds.Add(new StudioStashId { Endpoint = box.Endpoint, StashId = remote.Parent.Id });
                _db.Studios.Add(parent);
            }
            studio.Parent = parent;
        }

        return true;
    }

    private async Task<StashBoxRemoteStudio?> GetRemoteStudioAsync(StashBoxInstance box, string studioId, CancellationToken ct)
    {
        try
        {
            var response = await SendQueryAsync<StashBoxFindStudioResponse>(box, FindStudioByIdQuery, new { id = studioId }, ct);
            return response.FindStudio;
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to fetch studio {StudioId} from {Endpoint}", studioId, box.Endpoint);
            return null;
        }
    }

    private static StashBoxStudioMatchDto ToStudioMatchDto(StashBoxInstance box, StashBoxRemoteStudio studio)
    {
        return new StashBoxStudioMatchDto(
            Endpoint: box.Endpoint,
            StashBoxName: string.IsNullOrWhiteSpace(box.Name) ? box.Endpoint : box.Name,
            Id: studio.Id,
            Name: studio.Name,
            ImageUrl: studio.Images.FirstOrDefault()?.Url,
            Aliases: studio.Aliases
                .Where(a => !string.IsNullOrWhiteSpace(a))
                .Distinct(StringComparer.OrdinalIgnoreCase)
                .ToList(),
            Urls: studio.Urls
                .Select(u => u.Url)
                .Where(u => !string.IsNullOrWhiteSpace(u))
                .Distinct(StringComparer.OrdinalIgnoreCase)
                .ToList(),
            ParentName: studio.Parent?.Name
        );
    }

    public async Task<IReadOnlyList<StashBoxSceneMatchDto>> SearchScenesAsync(Scene scene, string? term, string? endpoint, CancellationToken ct)
    {
        var boxes = ResolveBoxes(endpoint);
        var strictEndpoint = !string.IsNullOrWhiteSpace(endpoint);
        var results = new List<StashBoxSceneMatchDto>();
        var sceneTitle = term ?? scene.Title;
        var sceneDuration = GetSceneDurationSeconds(scene);

        foreach (var box in boxes)
        {
            try
            {
                if (string.IsNullOrWhiteSpace(term))
                {
                    var existingStashId = scene.StashIds.FirstOrDefault(stashId => string.Equals(stashId.Endpoint, box.Endpoint, StringComparison.OrdinalIgnoreCase));
                    if (existingStashId != null)
                    {
                        var existing = await GetSceneMatchAsync(box.Endpoint, existingStashId.StashId, ct);
                        if (existing != null)
                        {
                            results.Add(existing);
                            continue;
                        }
                    }

                    var fingerprintQuery = BuildFingerprintQuery(scene);
                    if (fingerprintQuery.Count > 0)
                    {
                        var fingerprintResponse = await SendQueryAsync<StashBoxFindScenesByFingerprintsResponse>(
                            box,
                            FindScenesByFingerprintsQuery,
                            new { fingerprints = new[] { fingerprintQuery } },
                            ct);

                        foreach (var remote in fingerprintResponse.FindScenesBySceneFingerprints.SelectMany(batch => batch))
                        {
                            results.Add(await ToSceneMatchDtoAsync(box, remote, ct));
                        }
                        if (fingerprintResponse.FindScenesBySceneFingerprints.Any(batch => batch.Count > 0))
                            continue;
                    }
                }

                var effectiveTerm = string.IsNullOrWhiteSpace(term) ? scene.Title : term;
                if (string.IsNullOrWhiteSpace(effectiveTerm))
                    continue;

                var searchResponse = await SendQueryAsync<StashBoxSearchSceneResponse>(box, SearchSceneQuery, new { term = effectiveTerm }, ct);
                foreach (var remote in searchResponse.SearchScene)
                {
                    results.Add(await ToSceneMatchDtoAsync(box, remote, ct));
                }
            }
            catch (Exception ex) when (!strictEndpoint)
            {
                _logger.LogWarning(ex, "Skipping stash-box scene search for {Endpoint}", box.Endpoint);
            }
        }

        return results
            .GroupBy(match => $"{match.Endpoint}::{match.Id}", StringComparer.OrdinalIgnoreCase)
            .Select(group => group.First())
            .OrderByDescending(match => string.Equals(match.Title, sceneTitle, StringComparison.OrdinalIgnoreCase))
            .ThenBy(match => GetDurationDifference(sceneDuration, match.Duration))
            .ThenBy(match => match.Title ?? string.Empty, StringComparer.OrdinalIgnoreCase)
            .ThenBy(match => match.StashBoxName, StringComparer.OrdinalIgnoreCase)
            .ToList();
    }

    public async Task<StashBoxSceneMatchDto?> GetSceneMatchAsync(string endpoint, string sceneId, CancellationToken ct)
    {
        var box = ResolveBox(endpoint);
        var scene = await GetRemoteSceneAsync(box, sceneId, ct);
        return scene == null ? null : await ToSceneMatchDtoAsync(box, scene, ct);
    }

    public async Task<bool> MergeSceneAsync(Scene scene, string endpoint, string sceneId, StashBoxSceneImportRequestDto? importConfig, CancellationToken ct)
    {
        var box = ResolveBox(endpoint);
        var remote = await GetRemoteSceneAsync(box, sceneId, ct);
        if (remote == null)
            return false;

        await ApplyRemoteSceneAsync(scene, box.Endpoint, remote, importConfig, ct);
        return true;
    }

    private async Task<StashBoxRemotePerformer?> GetRemotePerformerAsync(StashBoxInstance box, string performerId, CancellationToken ct)
    {
        var response = await SendQueryAsync<StashBoxFindPerformerResponse>(box, FindPerformerByIdQuery, new { id = performerId }, ct);
        return response.FindPerformer;
    }

    private async Task<StashBoxRemoteScene?> GetRemoteSceneAsync(StashBoxInstance box, string sceneId, CancellationToken ct)
    {
        var response = await SendQueryAsync<StashBoxFindSceneResponse>(box, FindSceneByIdQuery, new { id = sceneId }, ct);
        return response.FindScene;
    }

    private async Task ApplyRemoteSceneAsync(Scene scene, string endpoint, StashBoxRemoteScene remote, StashBoxSceneImportRequestDto? importConfig, CancellationToken ct)
    {
        var setCoverImage = importConfig?.SetCoverImage ?? true;
        var setTags = importConfig?.SetTags ?? true;
        var setPerformers = importConfig?.SetPerformers ?? true;
        var setStudio = importConfig?.SetStudio ?? true;
        var onlyExistingTags = importConfig?.OnlyExistingTags ?? false;
        var onlyExistingPerformers = importConfig?.OnlyExistingPerformers ?? false;
        var onlyExistingStudio = importConfig?.OnlyExistingStudio ?? false;
        var markOrganized = importConfig?.MarkOrganized ?? false;
        var excludedTagNames = importConfig?.ExcludedTagNames?.ToHashSet(StringComparer.OrdinalIgnoreCase);
        var excludedPerformerNames = importConfig?.ExcludedPerformerNames?.ToHashSet(StringComparer.OrdinalIgnoreCase);
        var studioOverride = MatchSceneEntityOverride(importConfig?.StudioOverride, remote.Studio?.Id, remote.Studio?.Name);
        var performerOverrides = importConfig?.PerformerOverrides;
        var tagOverrides = importConfig?.TagOverrides;

        scene.Title = Coalesce(scene.Title, remote.Title) ?? scene.Title;
        scene.Code = Coalesce(scene.Code, remote.Code) ?? scene.Code;
        scene.Details = Coalesce(scene.Details, remote.Details) ?? scene.Details;
        scene.Director = Coalesce(scene.Director, remote.Director) ?? scene.Director;
        scene.Date = ParseDate(remote.Date) ?? scene.Date;
        if (markOrganized) scene.Organized = true;

        MergeSceneUrls(scene, remote.Urls.Select(url => url.Url));

        if (setStudio && remote.Studio != null)
        {
            var studio = await ResolveSceneStudioAsync(remote.Studio, endpoint, studioOverride, ct, allowCreate: !onlyExistingStudio);
            if (studio != null)
            {
                scene.Studio = studio;
                scene.StudioId = studio.Id;
            }
        }

        if (setTags)
        {
            foreach (var remoteTag in remote.Tags)
            {
                var tagOverride = MatchSceneEntityOverride(tagOverrides, remoteTag.Id, remoteTag.Name);
                if (GetSceneEntityOverrideAction(tagOverride) == SceneEntityOverrideAction.Skip)
                    continue;
                if (tagOverride == null && excludedTagNames != null && excludedTagNames.Contains(remoteTag.Name))
                    continue;
                var tag = await ResolveSceneTagAsync(remoteTag, endpoint, tagOverride, ct, allowCreate: !onlyExistingTags);
                if (tag == null)
                    continue;
                if (!scene.SceneTags.Any(link => link.TagId == tag.Id))
                {
                    scene.SceneTags.Add(new SceneTag { SceneId = scene.Id, TagId = tag.Id, Tag = tag });
                }
            }
        }

        if (setPerformers)
        {
            foreach (var remotePerformer in remote.Performers.Select(appearance => appearance.Performer).OfType<StashBoxRemotePerformer>())
            {
                var performerOverride = MatchSceneEntityOverride(performerOverrides, remotePerformer.Id, remotePerformer.Name);
                if (GetSceneEntityOverrideAction(performerOverride) == SceneEntityOverrideAction.Skip)
                    continue;
                if (performerOverride == null && excludedPerformerNames != null && remotePerformer.Name != null && excludedPerformerNames.Contains(remotePerformer.Name))
                    continue;
                var performer = await ResolveScenePerformerAsync(remotePerformer, endpoint, performerOverride, ct, allowCreate: !onlyExistingPerformers);
                if (performer == null)
                    continue;
                if (!scene.ScenePerformers.Any(link => link.PerformerId == performer.Id))
                {
                    scene.ScenePerformers.Add(new ScenePerformer { SceneId = scene.Id, PerformerId = performer.Id, Performer = performer });
                }
            }
        }

        // Download scene cover image
        if (setCoverImage && remote.Images.Count > 0)
        {
            await DownloadSceneCoverAsync(scene.Id, remote.Images[0].Url, ct);
        }

        var stashId = scene.StashIds.FirstOrDefault(id => string.Equals(id.Endpoint, endpoint, StringComparison.OrdinalIgnoreCase));
        if (stashId == null)
        {
            scene.StashIds.Add(new SceneStashId { Endpoint = endpoint, StashId = remote.Id, SceneId = scene.Id });
        }
        else
        {
            stashId.StashId = remote.Id;
        }
    }

    private async Task<Performer?> ResolveScenePerformerAsync(
        StashBoxRemotePerformer remote,
        string endpoint,
        StashBoxSceneEntityOverrideDto? entityOverride,
        CancellationToken ct,
        bool allowCreate)
    {
        return GetSceneEntityOverrideAction(entityOverride) switch
        {
            SceneEntityOverrideAction.Skip => null,
            SceneEntityOverrideAction.Existing when entityOverride?.LocalId is int localId => await _db.Performers.FirstOrDefaultAsync(performer => performer.Id == localId, ct),
            SceneEntityOverrideAction.Create => await FindOrCreatePerformerAsync(remote, endpoint, ct, allowCreate: true),
            _ => await FindOrCreatePerformerAsync(remote, endpoint, ct, allowCreate: allowCreate),
        };
    }

    private async Task<Studio?> ResolveSceneStudioAsync(
        StashBoxRemoteStudio remote,
        string endpoint,
        StashBoxSceneEntityOverrideDto? entityOverride,
        CancellationToken ct,
        bool allowCreate)
    {
        return GetSceneEntityOverrideAction(entityOverride) switch
        {
            SceneEntityOverrideAction.Skip => null,
            SceneEntityOverrideAction.Existing when entityOverride?.LocalId is int localId => await _db.Studios.FirstOrDefaultAsync(studio => studio.Id == localId, ct),
            SceneEntityOverrideAction.Create => await FindOrCreateStudioAsync(remote, endpoint, ct, allowCreate: true),
            _ => await FindOrCreateStudioAsync(remote, endpoint, ct, allowCreate: allowCreate),
        };
    }

    private async Task<Tag?> ResolveSceneTagAsync(
        StashBoxRemoteTag remote,
        string endpoint,
        StashBoxSceneEntityOverrideDto? entityOverride,
        CancellationToken ct,
        bool allowCreate)
    {
        return GetSceneEntityOverrideAction(entityOverride) switch
        {
            SceneEntityOverrideAction.Skip => null,
            SceneEntityOverrideAction.Existing when entityOverride?.LocalId is int localId => await _db.Tags.FirstOrDefaultAsync(tag => tag.Id == localId, ct),
            SceneEntityOverrideAction.Create => await FindOrCreateTagAsync(remote, endpoint, ct, allowCreate: true),
            _ => await FindOrCreateTagAsync(remote, endpoint, ct, allowCreate: allowCreate),
        };
    }

    private static StashBoxSceneEntityOverrideDto? MatchSceneEntityOverride(
        IEnumerable<StashBoxSceneEntityOverrideDto>? overrides,
        string? remoteId,
        string? name)
    {
        if (overrides == null)
            return null;

        return overrides.FirstOrDefault(entityOverride =>
            (!string.IsNullOrWhiteSpace(remoteId) && string.Equals(entityOverride.RemoteId, remoteId, StringComparison.OrdinalIgnoreCase)) ||
            (!string.IsNullOrWhiteSpace(name) && string.Equals(entityOverride.Name, name, StringComparison.OrdinalIgnoreCase)));
    }

    private static StashBoxSceneEntityOverrideDto? MatchSceneEntityOverride(
        StashBoxSceneEntityOverrideDto? entityOverride,
        string? remoteId,
        string? name)
    {
        if (entityOverride == null)
            return null;

        return MatchSceneEntityOverride(new[] { entityOverride }, remoteId, name);
    }

    private static SceneEntityOverrideAction GetSceneEntityOverrideAction(StashBoxSceneEntityOverrideDto? entityOverride)
    {
        return entityOverride?.Action.Trim().ToLowerInvariant() switch
        {
            "skip" => SceneEntityOverrideAction.Skip,
            "create" => SceneEntityOverrideAction.Create,
            "existing" => SceneEntityOverrideAction.Existing,
            _ => SceneEntityOverrideAction.Auto,
        };
    }

    private enum SceneEntityOverrideAction
    {
        Auto,
        Skip,
        Create,
        Existing,
    }

    private void ApplyRemotePerformer(Performer performer, string endpoint, StashBoxRemotePerformer remote)
    {
        performer.Name = remote.Name.Trim();
        performer.Disambiguation = string.IsNullOrWhiteSpace(remote.Disambiguation) ? performer.Disambiguation : remote.Disambiguation.Trim();
        performer.Gender = MapGender(remote.Gender) ?? performer.Gender;
        performer.Birthdate = ParseDate(remote.BirthDate) ?? performer.Birthdate;
        performer.DeathDate = ParseDate(remote.DeathDate) ?? performer.DeathDate;
        performer.Country = Coalesce(performer.Country, remote.Country);
        performer.Ethnicity = Coalesce(performer.Ethnicity, HumanizeGraphQlEnum(remote.Ethnicity));
        performer.EyeColor = Coalesce(performer.EyeColor, HumanizeGraphQlEnum(remote.EyeColor));
        performer.HairColor = Coalesce(performer.HairColor, HumanizeGraphQlEnum(remote.HairColor));
        performer.HeightCm = remote.Height > 0 ? remote.Height.Value : performer.HeightCm;
        performer.Measurements = Coalesce(performer.Measurements, FormatMeasurements(remote.Measurements));
        performer.FakeTits = Coalesce(performer.FakeTits, HumanizeGraphQlEnum(remote.BreastType));
        performer.CareerStart = remote.CareerStartYear > 0 ? new DateOnly(remote.CareerStartYear.Value, 1, 1) : performer.CareerStart;
        performer.CareerEnd = remote.CareerEndYear > 0 ? new DateOnly(remote.CareerEndYear.Value, 1, 1) : performer.CareerEnd;
        performer.Tattoos = Coalesce(performer.Tattoos, FormatBodyModifications(remote.Tattoos));
        performer.Piercings = Coalesce(performer.Piercings, FormatBodyModifications(remote.Piercings));

        var aliases = remote.Aliases
            .Where(alias => !string.IsNullOrWhiteSpace(alias))
            .Select(alias => alias.Trim())
            .Where(alias => !string.Equals(alias, remote.Name, StringComparison.OrdinalIgnoreCase));
        MergeAliases(performer, aliases);
        MergeUrls(performer, remote.Urls.Select(url => url.Url));

        var stashId = performer.StashIds.FirstOrDefault(id => string.Equals(id.Endpoint, endpoint, StringComparison.OrdinalIgnoreCase));
        if (stashId == null)
        {
            performer.StashIds.Add(new PerformerStashId
            {
                Endpoint = endpoint,
                StashId = remote.Id,
            });
        }
        else
        {
            stashId.StashId = remote.Id;
        }
    }

    private async Task DownloadPerformerImageAsync(Performer performer, StashBoxRemotePerformer remote, CancellationToken ct)
    {
        if (performer.ImageBlobId != null || remote.Images.Count == 0)
            return;

        try
        {
            var imageUrl = remote.Images[0].Url;
            using var response = await _httpClient.GetAsync(imageUrl, HttpCompletionOption.ResponseHeadersRead, ct);
            if (!response.IsSuccessStatusCode) return;

            var contentType = response.Content.Headers.ContentType?.MediaType ?? "image/jpeg";
            await using var stream = await response.Content.ReadAsStreamAsync(ct);
            performer.ImageBlobId = await _blobService.StoreBlobAsync(stream, contentType, ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to download performer image for {Name}", performer.Name);
        }
    }

    private async Task DownloadStudioImageAsync(Studio studio, StashBoxRemoteStudio remote, CancellationToken ct)
    {
        if (studio.ImageBlobId != null || remote.Images.Count == 0)
            return;

        try
        {
            var imageUrl = remote.Images[0].Url;
            using var response = await _httpClient.GetAsync(imageUrl, HttpCompletionOption.ResponseHeadersRead, ct);
            if (!response.IsSuccessStatusCode) return;

            var contentType = response.Content.Headers.ContentType?.MediaType ?? "image/png";
            await using var stream = await response.Content.ReadAsStreamAsync(ct);
            studio.ImageBlobId = await _blobService.StoreBlobAsync(stream, contentType, ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to download studio image for {Name}", studio.Name);
        }
    }

    private async Task DownloadSceneCoverAsync(int sceneId, string imageUrl, CancellationToken ct)
    {
        try
        {
            var generatedPath = _config.GeneratedPath;
            if (string.IsNullOrEmpty(generatedPath)) return;

            var hash = Convert.ToHexStringLower(System.Security.Cryptography.SHA256.HashData(BitConverter.GetBytes(sceneId)));
            var thumbPath = Path.Combine(generatedPath, "screenshots", hash[..2], $"{sceneId}.jpg");
            if (File.Exists(thumbPath)) return;

            using var response = await _httpClient.GetAsync(imageUrl, HttpCompletionOption.ResponseHeadersRead, ct);
            if (!response.IsSuccessStatusCode) return;

            var dir = Path.GetDirectoryName(thumbPath)!;
            Directory.CreateDirectory(dir);
            await using var stream = await response.Content.ReadAsStreamAsync(ct);
            await using var fileStream = new FileStream(thumbPath, FileMode.Create, FileAccess.Write, FileShare.None);
            await stream.CopyToAsync(fileStream, ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to download scene cover for scene {SceneId}", sceneId);
        }
    }

    private static void MergeAliases(Performer performer, IEnumerable<string> aliases)
    {
        var existing = performer.Aliases
            .Select(alias => alias.Alias)
            .ToHashSet(StringComparer.OrdinalIgnoreCase);

        foreach (var alias in aliases)
        {
            if (existing.Add(alias))
            {
                performer.Aliases.Add(new PerformerAlias { Alias = alias, PerformerId = performer.Id });
            }
        }
    }

    private static void MergeUrls(Performer performer, IEnumerable<string> urls)
    {
        var existing = performer.Urls
            .Select(url => url.Url)
            .ToHashSet(StringComparer.OrdinalIgnoreCase);

        foreach (var url in urls.Where(url => !string.IsNullOrWhiteSpace(url)).Select(url => url.Trim()))
        {
            if (existing.Add(url))
            {
                performer.Urls.Add(new PerformerUrl { Url = url, PerformerId = performer.Id });
            }
        }
    }

    private static void MergeSceneUrls(Scene scene, IEnumerable<string> urls)
    {
        var existing = scene.Urls
            .Select(url => url.Url)
            .ToHashSet(StringComparer.OrdinalIgnoreCase);

        foreach (var url in urls.Where(url => !string.IsNullOrWhiteSpace(url)).Select(url => url.Trim()))
        {
            if (existing.Add(url))
            {
                scene.Urls.Add(new SceneUrl { Url = url, SceneId = scene.Id });
            }
        }
    }

    private async Task<Performer?> FindOrCreatePerformerAsync(StashBoxRemotePerformer remote, string endpoint, CancellationToken ct, bool allowCreate = true)
    {
        var performer = await _db.Performers
            .Include(entity => entity.StashIds)
            .Include(entity => entity.Aliases)
            .Include(entity => entity.Urls)
            .FirstOrDefaultAsync(entity => entity.StashIds.Any(stashId => stashId.Endpoint == endpoint && stashId.StashId == remote.Id), ct)
            ?? await _db.Performers
                .Include(entity => entity.StashIds)
                .Include(entity => entity.Aliases)
                .Include(entity => entity.Urls)
                .FirstOrDefaultAsync(entity => entity.Name == remote.Name, ct);

        if (performer == null && !allowCreate)
        {
            return null;
        }

        if (performer == null)
        {
            performer = new Performer { Name = remote.Name };
            _db.Performers.Add(performer);
        }

        ApplyRemotePerformer(performer, endpoint, remote);
        await DownloadPerformerImageAsync(performer, remote, ct);
        return performer;
    }

    private async Task<Studio?> FindOrCreateStudioAsync(StashBoxRemoteStudio remote, string endpoint, CancellationToken ct, bool allowCreate = true)
    {
        var studio = await _db.Studios
            .Include(entity => entity.StashIds)
            .Include(entity => entity.Aliases)
            .Include(entity => entity.Urls)
            .FirstOrDefaultAsync(entity => entity.StashIds.Any(stashId => stashId.Endpoint == endpoint && stashId.StashId == remote.Id), ct)
            ?? await _db.Studios
                .Include(entity => entity.StashIds)
                .Include(entity => entity.Aliases)
                .Include(entity => entity.Urls)
                .FirstOrDefaultAsync(entity => entity.Name == remote.Name, ct);

        if (studio == null && !allowCreate)
        {
            return null;
        }

        if (studio == null)
        {
            studio = new Studio { Name = remote.Name };
            _db.Studios.Add(studio);
        }

        studio.Name = remote.Name.Trim();
        MergeAliases(studio, remote.Aliases);
        MergeUrls(studio, remote.Urls.Select(url => url.Url));
        UpsertStashId(studio.StashIds, endpoint, remote.Id, id => id.Endpoint, id => id.StashId, (id, value) => id.StashId = value, value => new StudioStashId { Endpoint = endpoint, StashId = value });

        // Download studio image
        await DownloadStudioImageAsync(studio, remote, ct);

        // Resolve parent studio
        if (remote.Parent != null && studio.ParentId == null)
        {
            var parent = await _db.Studios
                .Include(s => s.StashIds)
                .FirstOrDefaultAsync(s => s.StashIds.Any(id => id.Endpoint == endpoint && id.StashId == remote.Parent.Id), ct)
                ?? await _db.Studios
                    .Include(s => s.StashIds)
                    .FirstOrDefaultAsync(s => s.Name == remote.Parent.Name, ct);

            if (parent == null)
            {
                parent = new Studio { Name = remote.Parent.Name };
                parent.StashIds.Add(new StudioStashId { Endpoint = endpoint, StashId = remote.Parent.Id });
                _db.Studios.Add(parent);
            }
            studio.Parent = parent;
        }

        return studio;
    }

    private async Task<Tag?> FindOrCreateTagAsync(StashBoxRemoteTag remote, string endpoint, CancellationToken ct, bool allowCreate = true)
    {
        var tag = await _db.Tags
            .Include(entity => entity.StashIds)
            .Include(entity => entity.Aliases)
            .FirstOrDefaultAsync(entity => entity.StashIds.Any(stashId => stashId.Endpoint == endpoint && stashId.StashId == remote.Id), ct)
            ?? await _db.Tags
                .Include(entity => entity.StashIds)
                .Include(entity => entity.Aliases)
                .FirstOrDefaultAsync(entity => entity.Name == remote.Name, ct);

        if (tag == null && !allowCreate)
        {
            return null;
        }

        if (tag == null)
        {
            tag = new Tag { Name = remote.Name };
            _db.Tags.Add(tag);
        }

        tag.Name = remote.Name.Trim();
        tag.Description = Coalesce(tag.Description, remote.Description) ?? tag.Description;
        MergeAliases(tag, remote.Aliases);
        UpsertStashId(tag.StashIds, endpoint, remote.Id, id => id.Endpoint, id => id.StashId, (id, value) => id.StashId = value, value => new TagStashId { Endpoint = endpoint, StashId = value });
        return tag;
    }

    private static void MergeAliases(Studio studio, IEnumerable<string> aliases)
    {
        var existing = studio.Aliases.Select(alias => alias.Alias).ToHashSet(StringComparer.OrdinalIgnoreCase);
        foreach (var alias in aliases.Where(alias => !string.IsNullOrWhiteSpace(alias)).Select(alias => alias.Trim()).Where(alias => !string.Equals(alias, studio.Name, StringComparison.OrdinalIgnoreCase)))
        {
            if (existing.Add(alias))
                studio.Aliases.Add(new StudioAlias { Alias = alias, StudioId = studio.Id });
        }
    }

    private static void MergeUrls(Studio studio, IEnumerable<string> urls)
    {
        var existing = studio.Urls.Select(url => url.Url).ToHashSet(StringComparer.OrdinalIgnoreCase);
        foreach (var url in urls.Where(url => !string.IsNullOrWhiteSpace(url)).Select(url => url.Trim()))
        {
            if (existing.Add(url))
                studio.Urls.Add(new StudioUrl { Url = url, StudioId = studio.Id });
        }
    }

    private static void MergeAliases(Tag tag, IEnumerable<string> aliases)
    {
        var existing = tag.Aliases.Select(alias => alias.Alias).ToHashSet(StringComparer.OrdinalIgnoreCase);
        foreach (var alias in aliases.Where(alias => !string.IsNullOrWhiteSpace(alias)).Select(alias => alias.Trim()).Where(alias => !string.Equals(alias, tag.Name, StringComparison.OrdinalIgnoreCase)))
        {
            if (existing.Add(alias))
                tag.Aliases.Add(new TagAlias { Alias = alias, TagId = tag.Id });
        }
    }

    private static void UpsertStashId<TStashId>(ICollection<TStashId> stashIds, string endpoint, string stashId, Func<TStashId, string> getEndpoint, Func<TStashId, string> getStashId, Action<TStashId, string> setStashId, Func<string, TStashId> create)
    {
        var existing = stashIds.FirstOrDefault(item => string.Equals(getEndpoint(item), endpoint, StringComparison.OrdinalIgnoreCase));
        if (existing == null)
        {
            stashIds.Add(create(stashId));
            return;
        }

        if (!string.Equals(getStashId(existing), stashId, StringComparison.OrdinalIgnoreCase))
            setStashId(existing, stashId);
    }

    private IReadOnlyList<StashBoxInstance> ResolveBoxes(string? endpoint)
    {
        if (string.IsNullOrWhiteSpace(endpoint))
            return _config.Scraping.StashBoxes;

        return [ResolveBox(endpoint)];
    }

    private StashBoxInstance ResolveBox(string endpoint)
    {
        return _config.Scraping.StashBoxes.FirstOrDefault(box => string.Equals(box.Endpoint, endpoint, StringComparison.OrdinalIgnoreCase))
            ?? throw new InvalidOperationException($"Configured stash-box endpoint not found: {endpoint}");
    }

    private async Task<T> SendQueryAsync<T>(StashBoxInstance box, string query, object? variables, CancellationToken ct)
    {
        using var request = new HttpRequestMessage(HttpMethod.Post, box.Endpoint);
        if (!string.IsNullOrWhiteSpace(box.ApiKey))
            request.Headers.TryAddWithoutValidation("ApiKey", box.ApiKey);

        request.Content = JsonContent.Create(new StashBoxGraphQlRequest(query, variables), options: _jsonOptions);

        using var response = await _httpClient.SendAsync(request, HttpCompletionOption.ResponseHeadersRead, ct);
        var payload = await response.Content.ReadAsStringAsync(ct);

        if (payload.Contains("<!doctype", StringComparison.OrdinalIgnoreCase) || payload.Contains("<html", StringComparison.OrdinalIgnoreCase))
            throw new InvalidOperationException("Invalid endpoint");

        if (!response.IsSuccessStatusCode)
            throw new InvalidOperationException(string.IsNullOrWhiteSpace(payload) ? response.ReasonPhrase ?? "Request failed" : payload);

        var graphQl = JsonSerializer.Deserialize<StashBoxGraphQlResponse<T>>(payload, _jsonOptions)
            ?? throw new InvalidOperationException("Empty response from server");

        if (graphQl.Errors.Count > 0)
            throw new InvalidOperationException(string.Join("; ", graphQl.Errors.Select(error => error.Message)));

        if (graphQl.Data == null)
            throw new InvalidOperationException("No response from server");

        return graphQl.Data;
    }

    private static StashBoxInstance ToConfigBox(StashBoxDto dto) => new()
    {
        Endpoint = dto.Endpoint.Trim(),
        ApiKey = dto.ApiKey?.Trim() ?? string.Empty,
        Name = dto.Name?.Trim() ?? string.Empty,
        MaxRequestsPerMinute = dto.MaxRequestsPerMinute > 0 ? dto.MaxRequestsPerMinute : 240,
    };

    private static StashBoxPerformerMatchDto ToMatchDto(StashBoxInstance box, StashBoxRemotePerformer performer)
    {
        return new StashBoxPerformerMatchDto(
            Endpoint: box.Endpoint,
            StashBoxName: string.IsNullOrWhiteSpace(box.Name) ? box.Endpoint : box.Name,
            Id: performer.Id,
            Name: performer.Name,
            Disambiguation: performer.Disambiguation,
            Gender: HumanizeGraphQlEnum(performer.Gender),
            BirthDate: performer.BirthDate,
            Country: performer.Country,
            ImageUrl: performer.Images.FirstOrDefault()?.Url,
            Deleted: performer.Deleted,
            MergedIntoId: performer.MergedIntoId,
            Aliases: performer.Aliases
                .Where(alias => !string.IsNullOrWhiteSpace(alias))
                .Distinct(StringComparer.OrdinalIgnoreCase)
                .ToList(),
            Urls: performer.Urls
                .Select(url => url.Url)
                .Where(url => !string.IsNullOrWhiteSpace(url))
                .Distinct(StringComparer.OrdinalIgnoreCase)
                .ToList()
        );
    }

    private async Task<StashBoxSceneMatchDto> ToSceneMatchDtoAsync(StashBoxInstance box, StashBoxRemoteScene scene, CancellationToken ct)
    {
        var studioCandidate = await BuildStudioCandidateAsync(box.Endpoint, scene.Studio, ct);
        var performerCandidates = await BuildPerformerCandidatesAsync(box.Endpoint, scene, ct);
        var tagCandidates = await BuildTagCandidatesAsync(box.Endpoint, scene, ct);

        return new StashBoxSceneMatchDto(
            Endpoint: box.Endpoint,
            StashBoxName: string.IsNullOrWhiteSpace(box.Name) ? box.Endpoint : box.Name,
            Id: scene.Id,
            Title: scene.Title,
            Code: scene.Code,
            Date: scene.Date,
            Director: scene.Director,
            Details: scene.Details,
            StudioName: scene.Studio?.Name,
            ImageUrl: scene.Images.FirstOrDefault()?.Url,
            Duration: scene.Duration,
            PerformerNames: performerCandidates.Select(candidate => candidate.Name).ToList(),
            TagNames: tagCandidates.Select(candidate => candidate.Name).ToList(),
            Urls: scene.Urls.Select(url => url.Url).Where(url => !string.IsNullOrWhiteSpace(url)).Distinct(StringComparer.OrdinalIgnoreCase).ToList(),
            FingerprintAlgorithms: scene.Fingerprints.Select(fingerprint => fingerprint.Algorithm).Distinct(StringComparer.OrdinalIgnoreCase).ToList(),
            Fingerprints: scene.Fingerprints.Select(fp => new StashBoxFingerprintDto(fp.Algorithm, fp.Hash, fp.Duration)).ToList(),
            StudioCandidate: studioCandidate,
            PerformerCandidates: performerCandidates,
            TagCandidates: tagCandidates
        );
    }

    private async Task<StashBoxEntityCandidateDto?> BuildStudioCandidateAsync(string endpoint, StashBoxRemoteStudio? remoteStudio, CancellationToken ct)
    {
        if (remoteStudio == null || string.IsNullOrWhiteSpace(remoteStudio.Name))
            return null;

        var localId = await _db.Studios
            .Where(studio => studio.Name == remoteStudio.Name || studio.StashIds.Any(stashId => stashId.Endpoint == endpoint && stashId.StashId == remoteStudio.Id))
            .Select(studio => (int?)studio.Id)
            .FirstOrDefaultAsync(ct);

        return new StashBoxEntityCandidateDto(remoteStudio.Id, remoteStudio.Name.Trim(), localId.HasValue, localId);
    }

    private async Task<List<StashBoxEntityCandidateDto>> BuildPerformerCandidatesAsync(string endpoint, StashBoxRemoteScene scene, CancellationToken ct)
    {
        var remotePerformers = scene.Performers
            .Select(appearance => appearance.Performer)
            .OfType<StashBoxRemotePerformer>()
            .Where(performer => !string.IsNullOrWhiteSpace(performer.Name))
            .GroupBy(performer => performer.Id, StringComparer.OrdinalIgnoreCase)
            .Select(group => group.First())
            .ToList();

        if (remotePerformers.Count == 0)
            return [];

        var remoteIds = remotePerformers.Select(performer => performer.Id).Distinct(StringComparer.OrdinalIgnoreCase).ToList();
        var remoteNames = remotePerformers.Select(performer => performer.Name.Trim()).Distinct(StringComparer.OrdinalIgnoreCase).ToList();

        var matchedByStashId = remoteIds.Count == 0
            ? []
            : await _db.Performers
                .SelectMany(performer => performer.StashIds
                    .Where(stashId => stashId.Endpoint == endpoint && remoteIds.Contains(stashId.StashId))
                    .Select(stashId => new { stashId.StashId, PerformerId = performer.Id }))
                .ToListAsync(ct);

        var matchedByName = remoteNames.Count == 0
            ? []
            : await _db.Performers
                .Where(performer => remoteNames.Contains(performer.Name))
                .Select(performer => new { performer.Name, performer.Id })
                .ToListAsync(ct);

        var idsByStashId = new Dictionary<string, int>(StringComparer.OrdinalIgnoreCase);
        foreach (var match in matchedByStashId)
        {
            idsByStashId.TryAdd(match.StashId, match.PerformerId);
        }

        var idsByName = new Dictionary<string, int>(StringComparer.OrdinalIgnoreCase);
        foreach (var match in matchedByName)
        {
            idsByName.TryAdd(match.Name, match.Id);
        }

        return remotePerformers.Select(remotePerformer =>
        {
            var name = remotePerformer.Name.Trim();
            var exists = idsByStashId.TryGetValue(remotePerformer.Id, out var localId) || idsByName.TryGetValue(name, out localId);
            return new StashBoxEntityCandidateDto(remotePerformer.Id, name, exists, exists ? localId : null);
        }).ToList();
    }

    private async Task<List<StashBoxEntityCandidateDto>> BuildTagCandidatesAsync(string endpoint, StashBoxRemoteScene scene, CancellationToken ct)
    {
        var remoteTags = scene.Tags
            .Where(tag => !string.IsNullOrWhiteSpace(tag.Name))
            .GroupBy(tag => tag.Id, StringComparer.OrdinalIgnoreCase)
            .Select(group => group.First())
            .ToList();

        if (remoteTags.Count == 0)
            return [];

        var remoteIds = remoteTags.Select(tag => tag.Id).Distinct(StringComparer.OrdinalIgnoreCase).ToList();
        var remoteNames = remoteTags.Select(tag => tag.Name.Trim()).Distinct(StringComparer.OrdinalIgnoreCase).ToList();

        var matchedByStashId = remoteIds.Count == 0
            ? []
            : await _db.Tags
                .SelectMany(tag => tag.StashIds
                    .Where(stashId => stashId.Endpoint == endpoint && remoteIds.Contains(stashId.StashId))
                    .Select(stashId => new { stashId.StashId, TagId = tag.Id }))
                .ToListAsync(ct);

        var matchedByName = remoteNames.Count == 0
            ? []
            : await _db.Tags
                .Where(tag => remoteNames.Contains(tag.Name))
                .Select(tag => new { tag.Name, tag.Id })
                .ToListAsync(ct);

        var idsByStashId = new Dictionary<string, int>(StringComparer.OrdinalIgnoreCase);
        foreach (var match in matchedByStashId)
        {
            idsByStashId.TryAdd(match.StashId, match.TagId);
        }

        var idsByName = new Dictionary<string, int>(StringComparer.OrdinalIgnoreCase);
        foreach (var match in matchedByName)
        {
            idsByName.TryAdd(match.Name, match.Id);
        }

        return remoteTags.Select(remoteTag =>
        {
            var name = remoteTag.Name.Trim();
            var exists = idsByStashId.TryGetValue(remoteTag.Id, out var localId) || idsByName.TryGetValue(name, out localId);
            return new StashBoxEntityCandidateDto(remoteTag.Id, name, exists, exists ? localId : null);
        }).ToList();
    }

    private static List<object> BuildFingerprintQuery(Scene scene)
    {
        return scene.Files
            .SelectMany(file => file.Fingerprints)
            .Select(fingerprint =>
            {
                var algorithm = fingerprint.Type.ToLowerInvariant() switch
                {
                    "md5" => "MD5",
                    "oshash" => "OSHASH",
                    "phash" => "PHASH",
                    _ => null,
                };

                return algorithm == null || string.IsNullOrWhiteSpace(fingerprint.Value)
                    ? null
                    : new { algorithm, hash = fingerprint.Value } as object;
            })
            .Where(item => item != null)
            .Distinct()
            .Cast<object>()
            .ToList();
    }

    private static int? GetSceneDurationSeconds(Scene scene)
    {
        var maxDuration = scene.Files.Select(file => file.Duration).DefaultIfEmpty().Max();
        return maxDuration > 0 ? (int?)Math.Round(maxDuration) : null;
    }

    private static int GetDurationDifference(int? localDuration, int? remoteDuration)
    {
        if (!localDuration.HasValue && !remoteDuration.HasValue) return 0;
        if (!localDuration.HasValue || !remoteDuration.HasValue) return int.MaxValue;
        return Math.Abs(localDuration.Value - remoteDuration.Value);
    }

    private static string? Coalesce(string? currentValue, string? nextValue)
    {
        return string.IsNullOrWhiteSpace(nextValue) ? currentValue : nextValue.Trim();
    }

    private static DateOnly? ParseDate(string? value)
    {
        return DateOnly.TryParse(value, CultureInfo.InvariantCulture, DateTimeStyles.None, out var date)
            ? date
            : null;
    }

    private static string? FormatMeasurements(StashBoxRemoteMeasurements? measurements)
    {
        if (measurements == null || measurements.BandSize <= 0 || string.IsNullOrWhiteSpace(measurements.CupSize) || measurements.Waist <= 0 || measurements.Hip <= 0)
            return null;

        return $"{measurements.BandSize}{measurements.CupSize}-{measurements.Waist}-{measurements.Hip}";
    }

    private static string? FormatBodyModifications(List<StashBoxBodyModification>? items)
    {
        if (items == null || items.Count == 0)
            return null;

        return string.Join("; ", items.Select(item => string.IsNullOrWhiteSpace(item.Description) ? item.Location : $"{item.Location}, {item.Description}"));
    }

    private static GenderEnum? MapGender(string? value)
    {
        return value?.ToUpperInvariant() switch
        {
            "MALE" => GenderEnum.Male,
            "FEMALE" => GenderEnum.Female,
            "TRANSGENDER_MALE" => GenderEnum.TransgenderMale,
            "TRANSGENDER_FEMALE" => GenderEnum.TransgenderFemale,
            "INTERSEX" => GenderEnum.Intersex,
            "NON_BINARY" => GenderEnum.NonBinary,
            _ => null,
        };
    }

    private static string? HumanizeGraphQlEnum(string? value)
    {
        if (string.IsNullOrWhiteSpace(value))
            return null;

        var parts = value.Split('_', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
        return string.Join(' ', parts.Select(part => CultureInfo.InvariantCulture.TextInfo.ToTitleCase(part.ToLowerInvariant())));
    }

    private static string MapValidationError(Exception ex)
    {
        var message = ex.Message.ToLowerInvariant();
        return message switch
        {
            _ when message.Contains("doctype") || message.Contains("<html") => "Invalid endpoint",
            _ when message.Contains("connection refused") || message.Contains("no such host") || message.Contains("name or service not known") => "No response from server",
            _ when message.Contains("signature is invalid") || message.Contains("unauthorized") || message.Contains("forbidden") => "Invalid or expired API key.",
            _ when message.Contains("illegal base64 data") || message.Contains("token contains an invalid number of segments") || message.Contains("malformed") => "Malformed API key.",
            _ => $"Unknown error: {ex.Message}",
        };
    }

    private sealed record StashBoxGraphQlRequest(string Query, object? Variables);

    private sealed record StashBoxGraphQlResponse<T>
    {
        public T? Data { get; init; }
        public List<StashBoxGraphQlError> Errors { get; init; } = [];
    }

    private sealed record StashBoxGraphQlError(string Message);

    private sealed record StashBoxMeQueryResponse(StashBoxMeUser? Me);

    private sealed record StashBoxMeUser(string Name);

    private sealed record StashBoxSearchPerformerResponse(List<StashBoxRemotePerformer> SearchPerformer);

    private sealed record StashBoxFindPerformerResponse(StashBoxRemotePerformer? FindPerformer);

    private sealed record StashBoxSearchSceneResponse(List<StashBoxRemoteScene> SearchScene);

    private sealed record StashBoxFindSceneResponse(StashBoxRemoteScene? FindScene);

    private sealed record StashBoxSearchStudioResponse(List<StashBoxRemoteStudio> SearchStudio);

    private sealed record StashBoxFindStudioResponse(StashBoxRemoteStudio? FindStudio);

    private sealed record StashBoxFindScenesByFingerprintsResponse(List<List<StashBoxRemoteScene>> FindScenesBySceneFingerprints);

    private sealed record StashBoxRemotePerformer(
        string Id,
        string Name,
        string? Disambiguation,
        List<string> Aliases,
        string? Gender,
        bool Deleted,
        [property: JsonPropertyName("merged_into_id")] string? MergedIntoId,
        List<StashBoxRemoteUrl> Urls,
        List<StashBoxRemoteImage> Images,
        [property: JsonPropertyName("birth_date")] string? BirthDate,
        [property: JsonPropertyName("death_date")] string? DeathDate,
        string? Ethnicity,
        string? Country,
        [property: JsonPropertyName("eye_color")] string? EyeColor,
        [property: JsonPropertyName("hair_color")] string? HairColor,
        int? Height,
        StashBoxRemoteMeasurements? Measurements,
        [property: JsonPropertyName("breast_type")] string? BreastType,
        [property: JsonPropertyName("career_start_year")] int? CareerStartYear,
        [property: JsonPropertyName("career_end_year")] int? CareerEndYear,
        List<StashBoxBodyModification>? Tattoos,
        List<StashBoxBodyModification>? Piercings
    );

    private sealed record StashBoxRemoteUrl(string Url);

    private sealed record StashBoxRemoteImage(string Url);

    private sealed record StashBoxRemoteScene(
        string Id,
        string? Title,
        string? Code,
        string? Details,
        string? Director,
        int? Duration,
        string? Date,
        List<StashBoxRemoteUrl> Urls,
        List<StashBoxRemoteImage> Images,
        StashBoxRemoteStudio? Studio,
        List<StashBoxRemoteTag> Tags,
        List<StashBoxRemotePerformerAppearance> Performers,
        List<StashBoxRemoteFingerprint> Fingerprints
    );

    private sealed record StashBoxRemotePerformerAppearance(StashBoxRemotePerformer? Performer);

    private sealed record StashBoxRemoteStudio(string Id, string Name, List<string> Aliases, List<StashBoxRemoteUrl> Urls, List<StashBoxRemoteImage> Images, StashBoxRemoteStudioParent? Parent);
    private sealed record StashBoxRemoteStudioParent(string Id, string Name);

    private sealed record StashBoxRemoteTag(string Id, string Name, string? Description, List<string> Aliases);

    private sealed record StashBoxRemoteFingerprint(string Algorithm, string Hash, int? Duration);

    private sealed record StashBoxRemoteMeasurements(
        [property: JsonPropertyName("band_size")] int BandSize,
        [property: JsonPropertyName("cup_size")] string? CupSize,
        int Waist,
        int Hip
    );

    private sealed record StashBoxBodyModification(string Location, string? Description);
}
using Microsoft.AspNetCore.Mvc;
using Stash.Api.Services;
using Stash.Core.DTOs;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class SystemController(
    ISceneRepository sceneRepo, IImageRepository imageRepo,
    IGalleryRepository galleryRepo, IPerformerRepository performerRepo,
    IStudioRepository studioRepo, ITagRepository tagRepo,
    IGroupRepository groupRepo, ConfigService configService,
    ScraperService scraperService, StashBoxService stashBoxService) : ControllerBase
{
    [HttpGet("status")]
    public ActionResult<SystemStatusDto> GetStatus()
    {
        return Ok(new SystemStatusDto(
            Version: GetType().Assembly.GetName().Version?.ToString() ?? "0.1.0",
            AppDir: AppContext.BaseDirectory,
            ConfigFile: configService.ConfigPath,
            DatabasePath: "PostgreSQL"
        ));
    }

    [HttpGet("stats")]
    public async Task<ActionResult<StatsDto>> GetStats(CancellationToken ct)
    {
        var sceneCt = await sceneRepo.CountAsync(ct);
        var imageCt = await imageRepo.CountAsync(ct);
        var galleryCt = await galleryRepo.CountAsync(ct);
        var performerCt = await performerRepo.CountAsync(ct);
        var studioCt = await studioRepo.CountAsync(ct);
        var tagCt = await tagRepo.CountAsync(ct);
        var groupCt = await groupRepo.CountAsync(ct);

        return Ok(new StatsDto(sceneCt, imageCt, galleryCt, performerCt, studioCt, tagCt, groupCt, 0, 0));
    }

    [HttpGet("config")]
    public ActionResult<StashConfigDto> GetConfig()
    {
        return Ok(configService.GetConfig());
    }

    [HttpPut("config")]
    public async Task<ActionResult<StashConfigDto>> SaveConfig([FromBody] StashConfigDto config)
    {
        await configService.SaveConfigAsync(config);
        return Ok(configService.GetConfig());
    }

    [HttpGet("scrapers")]
    public ActionResult<IReadOnlyList<ScraperSummaryDto>> GetScrapers()
    {
        return Ok(scraperService.GetScrapers());
    }

    [HttpPost("scrapers/reload")]
    public ActionResult<IReadOnlyList<ScraperSummaryDto>> ReloadScrapers()
    {
        return Ok(scraperService.ReloadScrapers());
    }

    [HttpPost("stash-boxes/validate")]
    public async Task<ActionResult<StashBoxValidationResultDto>> ValidateStashBox([FromBody] StashBoxDto stashBox, CancellationToken ct)
    {
        return Ok(await stashBoxService.ValidateAsync(stashBox, ct));
    }
}

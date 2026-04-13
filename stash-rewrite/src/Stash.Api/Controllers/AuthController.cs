using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.IdentityModel.Tokens;
using Stash.Core.DTOs;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class AuthController(StashConfiguration configuration) : ControllerBase
{
    [HttpPost("login")]
    [AllowAnonymous]
    public ActionResult<LoginResponse> Login([FromBody] LoginRequest request)
    {
        var config = configuration.Auth;
        if (!config.Enabled)
            return Ok(new LoginResponse(GenerateToken(config, "anonymous"), "anonymous"));

        if (request.Username != config.Username ||
            !BCrypt.Net.BCrypt.Verify(request.Password, config.HashedPassword))
            return Unauthorized(new { message = "Invalid credentials" });

        return Ok(new LoginResponse(GenerateToken(config, request.Username), request.Username));
    }

    [HttpPost("logout")]
    public IActionResult Logout()
    {
        // JWT is stateless; client should discard token
        return Ok(new { message = "Logged out" });
    }

    [HttpGet("me")]
    [Authorize]
    public ActionResult<object> GetCurrentUser()
    {
        var username = User.FindFirstValue(ClaimTypes.Name) ?? "anonymous";
        return Ok(new { username });
    }

    private static string GenerateToken(AuthConfig config, string username)
    {
        var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(config.JwtSecret));
        var credentials = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);

        var claims = new[]
        {
            new Claim(ClaimTypes.Name, username),
            new Claim(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString())
        };

        var token = new JwtSecurityToken(
            issuer: "Stash",
            audience: "Stash",
            claims: claims,
            expires: DateTime.UtcNow.AddMinutes(Math.Max(config.MaxSessionAgeMinutes, 1)),
            signingCredentials: credentials
        );

        return new JwtSecurityTokenHandler().WriteToken(token);
    }
}

using Microsoft.AspNetCore.SignalR;

namespace Stash.Api.Hubs;

public class JobHub : Hub
{
    public override async Task OnConnectedAsync()
    {
        await Clients.Caller.SendAsync("ConnectionEstablished", Context.ConnectionId);
        await base.OnConnectedAsync();
    }
}

public class LogHub : Hub
{
    public override async Task OnConnectedAsync()
    {
        await base.OnConnectedAsync();
    }
}

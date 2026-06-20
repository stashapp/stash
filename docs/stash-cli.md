# Stash CLI 使用说明

`stash-cli` 是本地终端版 Stash 入口。它直接打开已有的 Stash SQLite 数据库，不启动 Stash server，也不会创建新的 SQLite 数据库。写入数据库前请先停止正在使用同一数据库的 Stash server。

## 构建与启动

```bash
make stash-cli
go run ./cmd/stash-cli --print-config
go run ./cmd/stash-cli --config ~/.config/stash-cli/config.toml --check
go run ./cmd/stash-cli --config ~/.config/stash-cli/config.toml
```

`--check` 只验证配置和数据库，并在启用 `scan_on_startup` 时执行启动扫描；不进入 TUI。

## 配置示例

```toml
database_path = "/home/tan/git/stash/.local/stash-go.sqlite"
media_dirs = ["/mnt/media/videos", "/run/media/tan/remote/videos"]
scan_on_startup = true
display_fields = ["name", "duration", "date"]
graphics_mode = "auto"
cache_dir = "~/.cache/stash-cli"
log_file = "~/.local/state/stash-cli/stash-cli.log"
log_level = "info"
log_stdout = false
ffmpeg_path = "/usr/bin/ffmpeg"
ffprobe_path = "/usr/bin/ffprobe"

[cover_fallback]
generate_with_ffmpeg = false

[blobs]
# Match the Stash server config. Use filesystem when blobs_storage: FILESYSTEM.
storage = "filesystem"
path = "/home/tan/git/stash/.local/blobs"

[stash_box]
endpoint = "https://stashdb.org/graphql"
api_key = ""
max_requests_per_minute = 240
```

`graphics_mode` 支持 `auto`、`kitty` 和 `list-only`。`kitty` 会强制使用 Kitty graphics，适合 Kitty/Ghostty 等已确认支持 Kitty 协议的终端；`auto` 会依赖终端探测，失败时可能回退到 glyph。远程视频请先通过 `sshfs`、NFS、SMB 等方式挂载成本机路径，再写入 `media_dirs`。

`log_stdout = false` 会把日志只写入 `log_file`，避免破坏 TUI 画面。调试时可临时设为 `true` 或用 `tail -f ~/.local/state/stash-cli/stash-cli.log` 观察日志。

`database_path` 必须指向 Stash server 已经初始化过的 sqlite 文件。`stash-cli` 不负责创建或迁移数据库；如果文件不存在，请先运行 Stash server 完成初始化。

`[blobs]` 必须和同一个 Stash server 的 blob 配置匹配。主 Stash 的 `blobs_storage: FILESYSTEM` 对应 `storage = "filesystem"` 和相同的 `blobs_path`；如果主 Stash 使用数据库存储 blob，则使用 `storage = "database"`。

`[stash_box]` 只用于显式抓取官方封面。扫描阶段只创建/更新已有数据库中的记录，不会默认通过 ffmpeg 生成封面，也不会自动联网抓取官方封面。`[cover_fallback].generate_with_ffmpeg = true` 时，官方封面抓取失败后才会用 ffmpeg 生成本地封面。

## Slash Commands

TUI 底部输入栏支持：

- `/search <query>`：搜索名称，也支持 `tag:<name>`、`performer:<name>`、`actress:<name>`。
- `/clear`：清除搜索条件。
- `/scan`：按配置的 `media_dirs` 扫描本地或已挂载的远程目录。
- `/cover fetch`：对当前选中的 scene 使用 stash-box fingerprints 抓取官方封面，并写入数据库。
- `/cover fetch-all`：对当前搜索条件匹配的所有 scene 批量抓取官方封面，不限当前 40 条列表页；每秒刷新进度，完成后刷新列表。官方抓取先使用 fingerprint；fingerprint 未命中时会从文件名提取番号并做精确 code/title 匹配。默认不会生成本地封面；仅当 `[cover_fallback].generate_with_ffmpeg = true` 时，官方抓取失败后才会用 ffmpeg 生成封面并计入 generated。每条抓取结果和加载封面的源图尺寸会以 info 级别写入 `log_file`。
- `/view grid|list`：切换显示模式。
- `/edit key=value ...`：编辑选中 scene，例如 `/edit title=Demo rating=80 organized=true watched=true`。
- `/help`：显示命令列表。
- `/quit`：退出。

第一版不会执行 tag/performer/studio 管理、scraper、文件删除/移动或终端内视频播放。`favorite` 不是 scene 字段，`/edit favorite=true` 会返回不支持。

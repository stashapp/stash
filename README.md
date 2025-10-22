# Stash

[![构建](https://github.com/stashapp/stash/actions/workflows/build.yml/badge.svg?branch=develop&event=push)](https://github.com/stashapp/stash/actions/workflows/build.yml)
[![Docker-Hub-Downloads](https://img.shields.io/docker/pulls/stashapp/stash.svg)](https://hub.docker.com/r/stashapp/stash 'DockerHub')
[![GitHub-Sponsoren](https://img.shields.io/github/sponsors/stashapp?logo=github)](https://github.com/sponsors/stashapp)
[![Open-Collective-Unterstützer](https://img.shields.io/opencollective/backers/stashapp?logo=opencollective)](https://opencollective.com/stashapp)
[![Go-Report-Card](https://goreportcard.com/badge/github.com/stashapp/stash)](https://goreportcard.com/report/github.com/stashapp/stash)
[![Discord](https://img.shields.io/discord/559159668438728723.svg?logo=discord)](https://discord.gg/2TsNFKt)
[![GitHub-Veröffentlichung (neueste nach Datum)](https://img.shields.io/github/v/release/stashapp/stash?logo=github)](https://github.com/stashapp/stash/releases/latest)
[![GitHub-Issues nach Label](https://img.shields.io/github/issues-raw/stashapp/stash/bounty)](https://github.com/stashapp/stash/labels/bounty)

### **Stash 是一个用 Go 编写的自托管 Web 应用程序，用于组织和提供您的成人影片。**

> **本仓库提供了rkmpp的硬件编解码支持**

![演示图片](docs/readme_assets/demo_image.png)

*   Stash 从互联网上收集关于您收藏的视频的信息，并且可以通过社区构建的插件进行扩展，以支持大量的内容制作者和网站。
*   Stash 支持多种视频和图像格式。
*   您可以标记视频并稍后找到它们。
*   Stash 提供有关表演者、标签、工作室等的统计信息。

您可以[观看一个 SFW（工作安全）的演示视频](https://vimeo.com/545323354)来了解它的实际操作。

更多信息，您可以查阅[文档](https://docs.stashapp.cc)或[阅读应用内手册](ui/v2.5/src/docs/en)。

# 安装 Stash

```yaml
services:
  stash:
    image: a15355447898a/stash-rk:latest
    container_name: stash-rk
    restart: always
    privileged: true
    ports:
      - "9999:9999"
    logging:
      driver: "json-file"
      options:
        max-file: "10"
        max-size: "2m"
    environment:
      - STASH_STASH=/video/
      - STASH_GENERATED=/generated/
      - STASH_METADATA=/metadata/
      - STASH_CACHE=/cache/
      - STASH_PORT=9999
    devices:
      - /dev/dri:/dev/dri           # GPU/显示接口
      - /dev/dma_heap:/dev/dma_heap # DMA 内存堆
      - /dev/mpp_service:/dev/mpp_service # MPP VPU 服务核心
      - /dev/rga:/dev/rga           # RGA 2D 加速器
      - /dev/mali0:/dev/mali0       # Mali GPU 设备，用于 OpenCL (HDR 色调映射)
    volumes:
      - /etc/localtime:/etc/localtime:ro
      ## 此目录存储配置文件、抓取器和插件
      - /docker-file/stash/config:/root/.stash
      ## 指向你的数据集合，比如你存放视频的地方
      - /docker-file/stash/video:/video
      ## Stash 的元数据存储位置
      - /docker-file/stash/metadata:/metadata
      ## 其他缓存内容存储位置
      - /docker-file/stash/cache:/cache
      ## 存储二进制大对象（比如场景封面、图片等）
      - /docker-file/stash/blobs:/blobs
      ## 存储生成的内容（截图、预览图、转码文件、精灵图等）
      - /docker-file/stash/generated:/generated
```

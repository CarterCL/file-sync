# file-sync

根据配置文件从远程下载文件到指定目录。支持预处理、格式转换、多任务并发。

## 使用方法

```bash
./file-sync -c config.yaml
```

参数：
- `-c`：配置文件路径（默认 `./config.yaml`）

## 配置文件示例

```yaml
default-ua: Mozilla/5.0

sync-tasks:
  - tag: example
    file-pairs:
      - url: https://example.com/data.txt
        path: ./output/data.txt
        convert: true
        extensions:
          - .bak
          - .old

      - url: https://example.com/other.txt
        path: ./output/other.txt
        ua: CustomUA/1.0
```

### 字段说明

| 字段 | 说明 |
|------|------|
| `default-ua` | 默认 User-Agent |
| `sync-tasks` | 同步任务列表 |
| `tag` | 任务标识 |
| `url` | 下载地址 |
| `path` | 保存路径 |
| `convert` | 是否去掉行首的 `-` 和 `'` |
| `ua` | 自定义 User-Agent（覆盖默认） |
| `extensions` | 预处理时追加到文件末尾的内容 |

---

# file-sync

Download files from remote URLs to local paths based on a YAML config. Supports preprocessing, format conversion, and concurrent tasks.

## Usage

```bash
./file-sync -c config.yaml
```

Flags:
- `-c`：config file path (default `./config.yaml`)

## Config Example

```yaml
default-ua: Mozilla/5.0

sync-tasks:
  - tag: example
    file-pairs:
      - url: https://example.com/data.txt
        path: ./output/data.txt
        convert: true
        extensions:
          - .bak
          - .old

      - url: https://example.com/other.txt
        path: ./output/other.txt
        ua: CustomUA/1.0
```

### Fields

| Field | Description |
|-------|-------------|
| `default-ua` | Default User-Agent header |
| `sync-tasks` | List of sync tasks |
| `tag` | Task label |
| `url` | Source URL to download |
| `path` | Destination file path |
| `convert` | Strip leading `-` and `'` from each line |
| `ua` | Per-file User-Agent (overrides default) |
| `extensions` | Lines appended to the file during preprocessing |

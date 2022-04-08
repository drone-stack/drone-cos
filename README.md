# qcloud cos

## 参数

- PLUGIN_DEBUG 是否开启调试模式
- PLUGIN_PAUSE 是否开启sleep模式
- PLUGIN_PROXY 设置代理模式

- PLUGIN_BUCKET bucket name
- PLUGIN_ACCESSKEY/PLUGIN_SECRETKEY accesskey/secretkey
- PLUGIN_REGION/PLUGIN_ENDPOINT region或者endpoint 

```yaml
region: ap-shanghai
endpoint: cos.ap-shanghai.myqcloud.com
```

- PLUGIN_SOURCE 源文件路径
- PLUGIN_TARGET 目标文件路径
- PLUGIN_STRIP_PREFIX 是否去除前缀

> 路径不包含*, 类似`aaa`

- PLUGIN_AUTOTIME 自动添加时间目录 如 `/aaa/` 则会自动添加 `/aaa/0407/`
- PLUGIN_TIMEFORMAT Go时间格式 如 `0102`

## 使用示例

> for k8s

```yaml
  - name: build-stable
    image: ysicing/drone-plugin-cos
    privileged: true
    pull: always
    settings:
      region: ap-shanghai
      bucket:
        from_secret: s3-bucket
      accesskey:
        from_secret: s3-access-key
      secretkey:
        from_secret: s3-secret-key
      source: dist/*
      target:
        from_secret: s3-stable-path
    when:
      event:
      - tag

  - name: build-edge
    image: ysicing/drone-plugin-cos
    privileged: true
    pull: always
    settings:
      region: ap-shanghai
      autotime: true
      bucket:
        from_secret: s3-bucket
      accesskey:
        from_secret: s3-access-key
      secretkey:
        from_secret: s3-secret-key
      source: dist/*
      target:
        from_secret: s3-edge-path
    when:
      branch:
      - master
```
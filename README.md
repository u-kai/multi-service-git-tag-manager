# Multi Service Git Tag Manager

- 一つのリポジトリで複数のサービス用の git tag を管理するためのツール
- 以下のような状況で利用することを想定しています
  - 1 つのリポジトリで複数のサービスを管理している
  - それぞれのサービスに対して独自のバージョンタグを管理したい
  - 一つのコミットに対して複数のサービス用のバージョンタグを付与し、コンテナイメージのタグ付けと連動したい

## Usage

```bash
$ msgtm tag v1.0.0 -s service1 service2 service3 -i HEAD
create service1-v1.0.0 and service2-v1.0.0 and service3-v1.0.0 to HEAD

$ git tag
```

- service の config file を作成したい
- 最初は service を羅列する
- tag を付けていくごとに、または config file の update 毎に config file を更新する

```yaml
services:
  - service1
    - latest
        - description: latest version
        - tag: service1-v1.0.0
        - commit: 0123456789abcdef
    - prev
        - description: previous version
        - tag: service1-v0.9.0
        - commit: 0123456789abcdef
  - service2
    - latest
        - description: latest version
        - tag: service2-v1.0.0
        - commit: 0123456789abcdef
  ...
```

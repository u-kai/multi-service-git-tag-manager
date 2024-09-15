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

- どうやって全てのサービスに対してタグをつける？
- HEAD とかってどうやって解釈する？そのまま Git に渡す？
- auto increment をできるようにしたい
- メッセージどうする?
  - service 毎にできたらいいけどな
    - vim を呼ぶか
  - 全てのサービスに対しても同じメッセージをつけることができるといいかも
  - AI とか使って、当該サービスのディレクトリの変更に対していコメントとかつけてもらったらめちゃいいかも

# api-sample-with-echo-ddd

## 使用技術

- **言語**: Go 1.23.1
- **フレームワーク**: Echo
- **アーキテクチャ**: DDD（ドメイン駆動設計）
- **データベース**: MyeSQL
- **ORM ライブラリ**: GORM

## ディレクトリ構成

```
/cmd            # エントリーポイント
/src       # アプリケーション内部コード
  /domain       # ドメイン層
  /usecase      # ユースケース層
  /interface    # インターフェース層
  /infrastructure # インフラ層
/pkg            # 外部パッケージ
```

## References

- https://github.com/gs1068/golang-ddd-sample

# Goのセットアップ

Go言語を始める上で学んだ開発環境のセットアップ方法について



## 前提条件

- goがインストールされていること



## Modulesを使ったプロジェクトの作成

[Go Modules](https://github.com/golang/go/wiki/Modules#quick-start-example)は現在主流になっているパッケージ管理ツール。

- Node.jsのnpmのような機能

```bash
# Go モジュールを初期化
go mod init <module_name>
```

Webに公開する場合は、`github.com/hibikihsh/api-sample-with-echo-ddd` のような形にする

ローカルで開発する場合は, `api-module` のような形でOK

[Go のモジュール管理【バージョン 1.17 改訂版】](https://zenn.dev/spiegel/articles/20210223-go-module-aware-mode?__hstc=53729015.64e060de20f355b50d595249802a3708.1738557069101.1748524600358.1748650237327.72&__hssc=53729015.2.1748650237327&__hsfp=688677955)



## パッケージと構成

ソースコードに含まれる関数や変数は必ず１つのパッケージに属する。 また同じフォルダ内に属するファイルは全て同一のパッケージである必要がある。

```
api-sample-with-echo-ddd
  ┣ main.go
  ┗ go.mod
```



## コード

```bash
# main.go ファイルを作成
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

- main 関数がエントリポイントとなる



### コードの実行

```bash
# プログラムを実行
go run main.go

# または、バイナリをビルドしてから実行
go build main.go
./main
```



## ビルド方法

下記コマンドで実行ファイルを作成することができる。

```
% go build
```



## 参考資料

- [Go言語の開発環境セットアップとサンプルプロジェクト作成] https://dev.classmethod.jp/articles/go-setup-and-sample/
- [Go 公式ドキュメント](https://golang.org/doc/)
- [Go Tour](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

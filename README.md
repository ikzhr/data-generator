# data-generator
リレーショナルデータベースのテーブルサンプルデータを生成させるためのツール。

## 使い方
下記の実行例では300000行のsampleA.tsvと100行のsampleB.tsvが生成される。
```sh
> generate -c ./sample/gtypes.yml ./sample/schema.xml
Table: sampleA  Rows: 300000
Table: sampleB  Rows: 100
```

オプションはヘルプ(-hオプション)参照。(READMEに書くと二重管理になるのでいずれ消す)
```sh
> generate -h
Usage of ./generate:
  -b int
        Buffer size of writer (default 1048576)
  -c string
        Path to column types def file
  -e string
        Extension of file to output (default "tsv")
  -n int
        The number of values generated (default 10)
  -o string
        Path to dir to output (default "./")
  -s string
        Separator of output data (default "\t")
  -t string
        List of table names
  -w int
        Woker size to generate table data (default 1)
```

## 開発者向け

### セットアップ
[go](https://golang.org/doc/install) & [dep](https://golang.github.io/dep/docs/installation.html)
をインストールした後、レポジトリルートで以下を実行して依存パッケージをインストール。

```sh
> dep ensure
```

### build
```sh
> go build generate.go
```

##  TODO

### コード整理
- 変数名, 関数名の整理・修正

### 機能追加
- range, 連続値のカラムに対して分布関数の種類とパラメータを指定可能にする?
- 異なるデータセットを生成するために乱数生成のseedを外部から指定可能にする
- カラム同士に関連性があるケースの考慮
	- カラムAとカラムBに大小関係がある
		- 例: day_from, day_to で day_to >= day_from という関係を満たす必要がある
	- カラムAはカラムBの値を包含している
		- 例: (商品コードカラム) = (年度カラム) + (商品固有のコード値)
- DateTimeのrangeで最大値を出力対象とする
	- 最大値は開区間として今のままの仕様とするのもありか？
# きふわらべ　第１７回ＵＥＣ杯コンピュータ囲碁大会　Go言語版

コンピュータ囲碁の思考エンジン。  
[きふわらべ　第１４回ＵＥＣ杯コンピュータ囲碁大会](https://github.com/muzudho/kifuwarabe-uec14) をベースにしています。  


# 環境設定

（第１４回版のことは）ブログにも書いてある：  
📖 [ブログ](https://qiita.com/muzudho1/items/cea62be01f7418bbf150)  


## Go 言語の環境設定

(1) PC に Go言語をインストールする：  
📖 [GO](https://go.dev/)  

Go 言語のインストール先は `C:\Program Files\Go\` にした。  
インストール後、PCの再起動をしておいた方が確実かと思う。  


(2) Visual Studio Code （以下、VSCode と呼称）の Extensions から、 `Go` （Go Team at Google）をインストールしておく。  


## ソースコードのデプロイ先

Input:  

```shell
echo %HOMEPATH%
```

Output:  

```
\Users\muzud
```

ユーザー・ディレクトリーが `\Users\muzud` と分かった。  
Go ではユーザー・ディレクトリーの下に `go/src` というフォルダーを作るものなので、そのようにフォルダーを作る。  

さらに `go/src` フォルダーの下に `github.com\muzudho\kifuwarabe-uec17` というフォルダーを作った。  

つまり、📁 `C:\Users\muzud\go\src\github.com\muzudho\kifuwarabe-uec17` というフォルダーを作った。  

この中にソースコードを置くことにする。  


## 動作確認

VSCode のターミナルに以下のコマンドを打鍵：  

PowerShell ではなく、Command Prompt を使う。  

Input:  

```shell
go version
```

Output:  

```
go version go1.25.3 windows/amd64
```


## Go のワークスペース作成

Go 言語の開発環境は、ワークスペースという名前で呼ぶそう。（C#言語ではソリューションと呼ぶもの）  

📄 `go.work` ファイルが既存でなければ、work コマンドを使って、ワークスペースを作る。  

```shell
go work init
```

`go.work` というファイルが作られる。  

`.gitignore` ファイルを作成する。  

📄 `.gitignore`:  

```
go.work
go.work.sum
```


## Go のワークスペースの使用

どのワークスペースを使用するか：  

```shell
go work use .
```

これで `go.work` ファイルに `use .` が追加される。  


もし、他にも使っているワークスペースがあれば、そのフォルダー名を指定する：  

```shell
go work use kernel
go work use debugger
```


## Go のモジュール作成

📄 `go.mod` ファイルが既存でなければ、以下のコマンドを打鍵：  

```shell
go mod init github.com/muzudho/kifuwarabe-uec17
```


## 必要なパッケージのダウンロード

以下のコマンドを打鍵：  

```shell
go mod tidy
```

📄 `go.sum` というファイルが更新される。パッケージのバージョンなども記載されている。  


他のサブ・フォルダーが既存で、そこに `go.mod` ファイルがあれば、そのディレクトリーで `go mod tity` コマンドを叩いた方がいいのだろうか？分からない。  


## プログラムの設定ファイル

📄 `engine.toml` が設定ファイルです。  
実行ファイルを実行するには、この設定ファイルを同じディレクトリーに置いておく必要があります。  


## プログラムのエントリー・ポイント

Go言語ではフォルダーを、ファイルを小分けにするただの入れ物として使うには向いていない。  
大量のファイルをトップレベルに置いておくことも気にしないことにする。  

とりあえず、トップレベルのフォルダーに 📄 `main.go` ファイルを置いておく。  
その中には `main` 関数がある。  
これを Go言語のプログラムのエントリー・ポイントとする。  


## プログラムの実行

Input:  

```shell
go run .
```


## 実行ファイルの作成

Input:

```shell
go build
```

これで以下の１つの実行ファイルが作成されます：  

📄 `kifuwarabe-uec17.exe`  


## 便利な .bat ファイル

必要ではないが便利にするものとして、 Windows 用に、 📁 `bin` フォルダー下に 📄 `*.bat` ファイルを置いています。  
使うときは、使いやすいようにしてください。  

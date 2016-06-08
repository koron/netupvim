# netupvim

[Jump to English](#english)

netupvim は Windows 用の Vim (香り屋版) をネットワーク経由で更新、修復、もしく
はインストールするためのプログラムです。

## 使い方

### セットアップ

ダウンロードしたアーカイブを展開し、以下の3つのファイルを Vim をインストールし
たフォルダ (もしくはこれからインストールするフォルダ) にコピーしてください。

*   netupvim.exe
*   UPDATE.bat
*   RESTORE.bat

#### プロキシーを使う

会社内で利用するなど、HTTP/HTTPS のアクセスにプロキシを使う必要がある場合は、環
境変数 `HTTPS_PROXY` 及び `HTTP_PROXY` を設定してください。これにより netupvim
のネットワークアクセスはすべてプロキシ経由になります。

環境変数名    |設定値の例
--------------|----------------------------
`HTTPS_PROXY` |`https://my.proxy.url:8443`
`HTTP_PROXY`  |`http://my.proxy.url:8080`

Windows10 での環境変数の設定方法は、以下のページを参照してください。

* [Windows 10 で環境変数を設定する](http://waman.hatenablog.com/entry/2015/12/09/085415)

### 更新(通常のアップデート)

UPDATE.bat をダブルクリックして実行してください。しばらく待つと Vim の差分
アップデートが完了します。Vim を起動中でも更新できます。更新完了後に Vim を再起
動すると、アップデートされた Vim が起動します。

Vim をインストールしていない状態で更新を実行すると、インストールになります。

### 修復(リストア)

インストールされた Vim のファイルが壊れてしまった場合には、リストアを実行すると
修復できます。

リストアするには RESTORE.bat を実行してください。しばらく待つと Vim の修復が完
了します。Vim は起動中でも修復は実行できますが、修復完了後に Vim を再起動してく
ださい。

### 問題が起こったら

更新、修復の実行時に問題が発生した場合には、ログファイルを以下に報告してくださ
い。

<https://github.com/koron/netupvim/issues/new>

ログファイルは `netupvim\log` というフォルダの下に、実行時刻をファイル名として
保存されています。例: `20160502T021805+0900.log`

## エキスパート向け情報

### 設定ファイル

netupvim は、実行時のカレントディレクトリにある設定ファイル netupvim.ini を起動
時に読み込みます。

### 開発版の利用

開発版を利用したい場合には、設定ファイルの `source` プロパティを以下のように設
定してください。

```ini
# 開発版
source = "develop"
```

また人柱版を利用したい場合には、以下のように設定してください。

```ini
# 人柱版
source = "canary"
```

同様に、<https://github.com/vim/vim-win32-installer> で配布されている本家の最新
のVimを利用できます。設定ファイルに以下の記述をしてください。

```ini
# vim/vim-win32-installer 版
source = "vim.org"
```

これらの版はあくまでも開発・実験用であり、予告なく不安定な動作の Vim が配信され
る可能性があることに留意してください。

また、一度 netupvim を実行した後で `source` プロパティを変更した場合の動作は未
定義です。直近でサポートする予定はありません。

### 実行回数制限

netupvim は GitHub API の回数制限の影響を受けます。そのため短時間に何度も実行す
ると(1時間に50回)、更新チェックに失敗するようになります。IPアドレス単位での制限
となるため、ルーターを通して複数のコンピューターが接続している場合には、一括で
制限を受けることに注意してください。

この制限を緩和するには GitHub のユーザーとトークンを準備し、netupvim へ設定して
ください。ユーザーとトークンを用いることで、制限回数は1時間あたり5000回に拡張さ
れます。設定には、設定ファイルの `github_user` と `github_token`、もしくは環境
変数の `NETUPVIM_GITHUB_USER` と `NETUPVIM_GITHUB_TOKEN` を使ってください。設定
ファイルと環境変数の両方を設定した場合には、設定ファイルのものが優先されます。
以下は netupvim.ini の設定例です。

```ini
github_user = "koron"
github_token = "0123456789abcdef0123456789abcdef'
```

FIXME: トークンの作成・取得方法を解説する

---

# English

netupvim is a program to update, restore or install Vim (+kaoriya version) for
Windows.

## How to use

### Setup

Extract a downloaded archive file, then copy below three files into the folder
where you have installed Vim (or you are going to install).

*   netupvim.exe
*   UPDATE.bat
*   RESTORE.bat

#### HTTPS/HTTP Proxy

You should setup two environment variables `HTTPS_PROXY` and `HTTP_PROXY`, when
you use netupvim at the network which need proxy to access HTTPS/HTTP.
Netupvim will use the proxy to access HTTPS/HTTP when these environment
variables are available.

Env Name      |Value Example
--------------|----------------------------
`HTTPS_PROXY` |`https://my.proxy.url:8443`
`HTTP_PROXY`  |`http://my.proxy.url:8080`

Please refer below links to set environment variables on Windows10.

* http://superuser.com/a/949573
* https://youtu.be/C-U9SGaNbwY

### Update

Double click UPDATE.bat and execute it.  After a while, netupvim has finished
to update your Vim.  You can update Vim, if it is executing.  You should
restart Vim after updated.

If you execute UPDATE.bat when you have not installed Vim, it will install.

### Restore

You can restore Vim files, if there are some broken.

Double click RESTORE.bat to restore those.  It will download recent Vim's
archive and extract/install it forcibly after a while.  You can restore Vim, it
is executing, but you should restart Vim, after finished to restore.

### When met trouble

When you met some troubles, plesae send log file to the issue tracker.

<https://github.com/koron/netupvim/issues/new>

Netupvim's log files are saved into `netupvim\log` folder with name which
determined from the time to execute.  Ex: `20160502T021805+0900.log`

## For Expert

TODO: translate me.

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

*   [Windows 10 で環境変数を設定する][2]

ユーザー認証(ユーザー名とパスワードの入力)が必要なプロキシを利用している場合に
は、環境変数で指定するURLに以下のようなフォーマットで追加してください。

    https://{username}:{password}@{host}:{port}
    http://{username}:{password}@{host}:{port}

`username` が `foo@example.com` で `password` が `bar123456` の場合には、設定す
る値は以下のようになるでしょう。

    https://foo%40example.com:bar123456@my.proxy.url:8443
    http://foo%40example.com:bar123456@my.proxy.url:8080

`username` と `password` はそれぞれ URL エンコードする必要があります。そのため
上記の例では `@` を `%40` に変換しています。

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

このファイルで設定できる項目は以下のとおりです。

項目                    |説明
------------------------|-----------------------------------------------------
`source`                |取得・更新するリリースの種類。release, develop, canary, vim.org のいずれか。デフォルトは release
`target_dir`            |更新対象のディレクトリ。デフォルトはカレントディレクトリで、通常は指定する必要はない
`cpu`                   |CPUの種類: x86, amd64 のどちらかで、デフォルトは自動判定
`github_token`          |更新確認を頻繁に行えるようにするためのトークン。取得方法は別セクションを参照。環境変数 `NETUPVIM_GITHUB_TOKEN` でも設定できる
`github_verbose`        |GitHub との通信をデバッグするためのオプション
`download_timeout`      |ダウンロードのタイムアウト。デフォルトは "5m"
`log_rotate_count`      |ログローテーションの世代数
`exe_rotate_count`      |実行ファイルローテーションの世代数
`disable_self_update`   |netupvim 自身の更新を抑制する
                               
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
ると(1時間に60回)、更新チェックに失敗するようになります。IPアドレス単位での制限
となるため、ルーターを通して複数のコンピューターが接続している場合には、一括で
制限を受けることに注意してください。

この制限を緩和するには GitHub の Personal access token (以下トークン) を作成
し、netupvim へ設定してください。トークンを設定することで、制限回数は1時間あた
り5000回に拡張されます。設定には、設定ファイルの `github_token`、もしくは環境変
数の `NETUPVIM_GITHUB_TOKEN` を使ってください。設定ファイルと環境変数の両方を設
定した場合には、設定ファイルのものが優先されます。以下は netupvim.ini の設定例
です。

```ini
github_token = "0123456789abcdef0123456789abcdef'
```

トークンを作成する方法は [Creating an access token for command-line use][1] を
参照してください。netupvim で利用するトークンはいかなるスコープ(権限)も必要とし
ていません。そのため参照先の手順5. "Select the scopes you with to grant to this
token" では、1つもスコープを選択しないで構いません。

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

*   <http://superuser.com/a/949573>
*   <https://youtu.be/C-U9SGaNbwY>

If you use the proxy which requires authorization (input username and
password), you should add those username and password into the URL in above
environment variables in format below.

    https://{username}:{password}@{host}:{port}
    http://{username}:{password}@{host}:{port}

Example: `username` is `foo@example.com` and `password` is `bar123456`, it
would be like this:

    https://foo%40example.com:bar123456@my.proxy.url:8443
    http://foo%40example.com:bar123456@my.proxy.url:8080

Both `username` and `password` should be encoded as URL.  Therefore `@` is
converted to `%40` in above example.

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

When you met some troubles, please send log file to the issue tracker.

<https://github.com/koron/netupvim/issues/new>

Netupvim's log files are saved into `netupvim\log` folder with name which
determined from the time to execute.  Ex: `20160502T021805+0900.log`

## For Expert

### Configuration file

Netupvim reads a configuration file "netupvim.ini" in current directory when
started.

Configurable items by the file are listed in below:

Item                    |Description
------------------------|-----------------------------------------------------
`source`                | Channel to update: one of "release", "develop", "canary" or "vim.org". Default is "release".
`target_dir`            | Direct to update. Default is current directory, and shouldn't be set usually.
`cpu`                   | CPU architecture: one of "x86" or "amd64". Default will be detected automatically.
`github_token`          | The GitHub's token to check update more frequently. See other section for more details. It can be set by `NETUPVIM_GITHUB_TOKEN` env.
`github_verbose`        | Enable debug log for communication with GitHub.
`download_timeout`      | Timeout for download operations. Default is "5m".
`log_rotate_count`      | Number of generations for log file rotation.
`exe_rotate_count`      | Number of generations for ".exe" file rotation.
`disable_self_update`   | Disable netupvim's self update.

### TODO: translate other sections.

[1]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/
[2]: http://waman.hatenablog.com/entry/2015/12/09/085415

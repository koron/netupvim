# netupvim

[Jump to English](#english)

netupvim は Windows に Vim (香り屋版) をネットワーク経由で更新、修復、もしくは
インストールするためのプログラムです。

## 使い方

### セットアップ

ダウンロードしたアーカイブを展開し、以下の3つのファイルを Vim をインストールし
たフォルダ (もしくはこれからインストールするフォルダ) にコピーしてください。

*   netupvim.exe
*   UPDATE.bat
*   RESTORE.bat

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

### 実行回数制限

netupvim は GitHub API の回数制限の影響を受けます。そのため短時間に何度も実行す
ると(1時間に50回程度以上の頻度で)、更新チェックに失敗するようになります。制限時
にはIPアドレス単位での制限となるため、ルーターを通して複数のコンピューターが接
続されている場合には、一括で制限を受けるため注意してください。

### 開発版の利用

開発版を利用したい場合には、netupvim.exe と同じ場所に netupvim.ini という名前の
ファイルを置き、以下の内容を記述してください。

```ini
# 開発版
source = "develop"
```

また人柱版を利用したい場合には、netupvim.ini の内容は以下のようにしてください。

```ini
# 人柱版
source = "canary"
```

これらの版はあくまでも開発・実験用であり、予告なく不安定な動作の Vim が配信され
る可能性があることに留意してください。

---

# English

netupvim is a program to update, restore or install Vim (+kaoriya version) for
Windows.

## How to use

### Setup

Unarchive a downloaded archive file, then copy below three files into your folder where you install Vim (or you are going to install).

*   netupvim.exe
*   UPDATE.bat
*   RESTORE.bat

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

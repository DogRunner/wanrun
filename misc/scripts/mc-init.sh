#!/bin/sh

set -e  # エラーが発生したらスクリプトを終了
set -x  # 実行されるコマンドを出力（デバッグ用）

# MinIOサーバーへの接続設定
mc alias set wanrun http://minio:9000 admin adminpass

# バケットの作成
mc mb wanrun/cms|| true

# ファイルのアップロード
mc cp /s3_local_init_data/* wanrun/cms|| true

# 確認のため、アップロードしたファイルをリスト表示
mc ls wanrun/cms

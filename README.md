# GraylogPractice
實作Graylog收集器

--

# 環境準備

使用系統 -> Mac

容器開發測試 -> Docker

需安裝準備

- Elasticsearch 
- MongoDB
- Graylog
- `go get -u github.com/lestrrat-go/file-rotatelogs`
- `go get -u github.com/rifflock/lfshook`
- `go get -u github.com/sirupsen/logrus`

Graylog Docker-Compose.yml 如附件

docker-compose -f docker-compose.yml up -d

環境如成功建立後，Graylog 將會於localhost:9000啟用

account/pwd 皆為 admin

# 設置Graylog Input

於上方工具列 System->Inputs

Notice : 選單可拖拉

選擇GELF TCP 設定port :12201

設置完成後 即可於 Terminal 測試是否能接收log

測試語法如下

 `echo '{"version": "3.1415926","host":"terminal666","short_message":"A short message that helps you identify wt is going on","full_message":"Backtrace here\n\nmore stuff","level":2,"_user_id":7443,"_some_info":"Aja","_some_env_var":"Attk"}' | nc -w 1 127.0.0.1 12201`
  
實作用則是設定 GELF UDP port:12201

go run 專案且至對應路徑觸發其log機制即可

## config 格式
```
set:
  graylogHost: 127.0.0.1:12201
  # 是否於根目錄建立log檔案
  createFile: true
  # 檔案保存時間
  LogTimeLimit: 168
  # 檔案切割時間
  LogSliceTime: 24
```

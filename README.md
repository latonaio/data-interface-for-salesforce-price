# data-interface-for-salesforce-price

## 概要
data-interface-for-salesforce-price は、salesforce の価格オブジェクト取得に必要なデータの整形、および作成時に salesforce から返ってきた response の MySQL への格納を行うマイクロサービスです。

## 動作環境
data-interface-for-salesforce-price は、aion-coreのプラットフォーム上での動作を前提としています。  
使用する際は、事前に下記の通りAIONの動作環境を用意してください。     
  
* OS: Linux OS     
* CPU: ARM/AMD/Intel     
* Kubernetes     
* [AION](https://github.com/latonaio/aion-core)のリソース      

## セットアップ
1. 以下のコマンドを実行して、docker imageを作成してください。
```
$ cd /path/to/data-interface-for-salesforce-price
$ make docker-build
```

2. 本マイクロサービスは DB に MySQL を使用します。MySQL に関する設定を、 `data-interface-for-salesforce-price.yaml` の環境変数に記述してください。

| env_name | description |
| --- | --- |
| MYSQL_HOST | ホスト名 |
| MYSQL_PORT | ポート番号 |
| MYSQL_USER | ユーザー名 |
| MYSQL_PASSWORD | パスワード |
| MYSQL_DBNAME | データベース名 |
| MAX_OPEN_CONNECTION | 最大コネクション数 |
| MAX_IDLE_CONNECTION | アイドル状態の最大コネクション数 |
| KANBANADDR: | kanban のアドレス |
| TZ | タイムゾーン |

## 起動方法
以下のコマンドを実行して、podを立ち上げてください。
```
$ cd /path/to/data-interface-for-salesforce-price
$ kubectl apply -f data-interface-for-salesforce-price.yaml
```

## kanban との通信
### kanban(ui-backend) から受信するデータ
kanban から受信する metadata に下記の情報を含む必要があります。

| key | value |
| --- | --- |
| method | get |
| object | Price |
| connection_type | request |
| districtId | 価格が紐づくID |

具体例: 
```example
# metadata (map[string]interface{}) の中身

"method": "get"
"object": "Price"
"districtId":      "a080l000007iaF4AAI",
"connection_type": "request"
```

### kanban に送信するデータ

#### PriceRecord
kanban に送信する metadata は下記の情報を含みます。

| key | type | description |
| --- | --- | --- |
| method | string | get |
| object | string | PriceRecord |
| connection_key | string | price_get |
| query_params | map[string]string | クエリパラメータ |

具体例: 
```example
# metadata (map[string]interface{}) の中身

"connection_key": "price_get",
"method":         "get",
"object":         "PriceRecord",
"query_params": map[string]string{
    "platId": "a080l000007iaF4AAI",
},
```

#### PriceRecordSeriesNumber
PriceRecord のデータ取得後に実行されます。
kanban に送信する metadata には下記の情報を含みます。

| key | type | description |
| --- | --- | --- |
| method | string | get |
| object | string | PriceRecordSeriesNumber |
| connection_key | string | price_get |
| query_params | map[string]string | クエリパラメータ |

具体例: 
```example
# metadata (map[string]interface{}) の中身

"connection_key": "price_get",
"method":         "get",
"object":         "PriceRecord",
"query_params": map[string]string{
    "priceRecordId": "a080l000007iaK5IDS",
},
```

## kanban(salesforce-api-kube) から受信するデータ
kanban からの受信可能データは下記の形式です

### PriceRecord
kanban から受信する metadata に下記の情報を含む必要があります。

| key | value |
| --- | --- |
| key | 文字列 "PriceRecord" を指定 |
| content | PriceRecord の詳細情報を含む JSON 配列 |
| connection_type | response |

具体例:
```example
# metadata (map[string]interface{}) の中身

"key": "PriceRecord"
"content": "[{xxxxxxxxxxx}]"
"connection_type": "response"
```

### PriceRecordSeriesNumber
kanban から受信する metadata に下記の情報を含む必要があります。

| key | value |
| --- | --- |
| key | 文字列 "PriceRecordSeriesNumber" を指定 |
| content | PriceRecordSeriesNumber の詳細情報を含む JSON 配列 | 
| connection_type | response |

具体例:
```example
# metadata (map[string]interface{}) の中身

"key": "PriceRecordSeriesNumber"
"content": "[{xxxxxxxxxxx}]"
"connection_type": "response"
```


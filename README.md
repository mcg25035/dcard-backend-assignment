# dcard-backend-intern-assignment
本專案採用 Sqlite 作為資料庫系統，後端框架選用 Golang 的 Fiber，並使用 Redis 作為快取。整體架構基於 Restful API 設計原則構建。<br>
伺服器每3分鐘會將 Sqlite 的資料和 Redis 的資料同步一遍，在建立廣告的時候會即時更新Redis，每次只會在Redis上留下接下來五分鐘有效的廣告。

## 專案結構
- Data Handlers:
> - dbHandler: 負責與 Sqlite 資料庫進行互動。
> - cacheHandler: 管理 Redis 快取的讀寫操作。
這兩者形成了 dataHandler，用於處理所有資料相關的邏輯。

- HTTP Handlers:
> - 處理 HTTP 請求的 httpHandler 包含兩部分：
> - POST /api/v1/ad: 用於創建廣告。
> - GET /api/v1/ad: 用於獲取廣告資訊。
這些API的具體規範跟題目中的說明一樣。

## 前提條件
在運行本專案之前，請確保已安裝 Redis 並設置為無密碼模式。在使用過程中，如果遇到性能瓶頸或意外崩潰，可能需要調整 Redis 的記憶體限制至少為 1GB。


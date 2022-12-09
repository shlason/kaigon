# kaigon

一個可以聊天和互相追蹤的論壇系統。

API 文件：[點我](https://kaigon.sidesideeffect.io/swagger/index.html)

## 功能規劃
- 使用者相關
    - 建立使用者
    - 更新使用者相關顯示資訊 (暱稱、大頭貼、自我介紹、第三方登入綁定)
    - 更新即時通知設定的啟用 (聊天室有新訊息、自己的發文有人回應或留言、追蹤的人有新動態等等)
    - 修改密碼和 Email
- 驗證和授權相關
    - 產生 Captcha 圖片，當帳號建立或登入時
    - Google OAuth
    - 使用者 Email 驗證
    - auth token (JWT), refresh token
- 聊天室
    - 個人和群組聊天室
    - 客製化個人聊天室設定和聊天室成員統一可見的設定
    - 已讀功能
- 看板 (開發中)
    - 特定話題的分類
- 文章 (開發中)
    - 文章會建立在特定看板底下
    - 文章可以留言或表情回應
    - 文章發表者可以即時收到相關的回應通知
    - 文章可以收藏追蹤
    - 有 hashtag
- 搜尋 (開發中)
    - 搜尋使用者
    - 搜尋文章
    - 搜尋 hashtag 相關文章
    - 搜尋特定看板
- 上傳圖片
    - 上傳至 S3

## 目前功能的實現方式
使用者自身相關：
- Email 辦帳號：透過發送驗證信來驗證 Email，發送驗證信目前是簡單的先使用 Go 原生的 package `net/smtp` 來進行發送，寄件者則是使用自己的 domain `sidesideeffect.io` 來申請 Google Workspace 的帳號，用以寄送該服務的所有郵件。
- 登入、登出：帳號登入成功後會在 response 上加上 cookie 用以儲存長效期的 refresh_token，前端再藉由 Call 一隻專門拿短效期的 auth_token API，該 API 會檢查 cookie 的 refresh_token 是否有效，並回傳相對應使用者的 auth_token，auth_token 的效期為 15 分鐘，refresh_token 的效期為 20 天。
若在不同裝置登入，也會拿到相同的 refresh_token 這樣可以相對方便的實現全站登出 (每個裝置登出)。
- 第三方登入：目前支援 Google 登入，若是一開始使用 Email 建立帳號．後續也可以在後台設定使用 Google 來快速登入，會有一張資料表專門關聯使用者帳號和第三方登入之間的連接關係。
- 忘記密碼：會發送一封設定密碼的信給使用者，該信中有重設密碼的連結，該連結帶有儲存在 Redis 的短效期驗證 token，會在之後使用者設定好新密碼時，一起帶給 API 來檢查這個重設密碼的請求是否合法。
- 驗證 Email：跟忘記密碼一樣原理，會帶有驗證 token 來檢查請求是否合法。

驗證相關：
- Captcha 圖片：建立 Captcha 時會產生一組 UUID，該 UUID 會關聯特定隨機產出的 Captcha code，而圖片就是藉由這 Captcha code 再做一些圖片處理來產生的，Captcha 本身有時效因此同一份 UUID 可以一直更新來獲得不同的 captcha code 繼續用下去。
- 聊天室 WebSocket token：因為 WebSocket 沒辦法使用 Header 因此只能待驗證 token 在 URL 上，但直接用登入的 auth token 覺得會有安全疑慮，因此產生一個非常短效期的 token 來作為連接 WebSocket 用。

聊天室相關：
- 聊天訊息：目前使用 MongoDB 來作為儲存手段，是因為當初在比較時發現 MongoDB 相較於 MySQL 來說對於寫入、多筆資料寫入的效能都好很多，然後覺得聊天訊息會是一個非常需要頻繁寫入的情境，因此選用 MongoDB。
- 連接管理：目前是一台機器去連接多個 Client 的 WebSocket，在實現上，我把多個 WebSocket 的連接都開一個 goroutine 出來，並且有個 Client Manager 使用 channel 作為單位，來管理這些 Client 的連接並透過 channel 向各個 goroutine 做溝通，在連上和斷線或是需要主動發送聊天訊息時都是由 Client Manager 來進行匹配和管理。

上傳圖片相關：
- S3: 使用 AWS 的 IAM 建立一個給 server 用的 image uploader role 並且透過 Cloudflare CDN 來提升圖片的載入速度。

## 技術棧
- Go
- Gin (管理路由和 middleware)
- `gopkg.in/ini.v1` (用來讀取設定檔 *.ini)
- JWT (作為登入用 token 格式，並擁有短效期和長效期的 token 流程，auth token, refresh token)
- MySQL (存放使用者相關個人資料)
- GORM
- MongoDB (存放聊天室相關訊息和文章相關資料)
- Redis (作為一些有效期的臨時 token 的存放地)
- swaggo/swag (產 API 文件)
- autocert (產 SSL 憑證)
- github Actions (用以 build 專案和部署到 EC2 上)
- Cloudflare DNS
- AWS EC2
- AWS S3
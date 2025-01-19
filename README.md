# Go Cache Toolkit

一個簡單、彈性且高效能的 Go 快取工具套件，專為需要快速管理記憶體快取、設定存活時間以及自動清理過期資料的應用場景而設計，並支援多協程環境的安全操作。

## 功能特色

- **泛型支援**：輕鬆存取任意類型資料，免除繁瑣的類型轉換。
- **併發安全**：基於 `sync.Map` 的設計，確保多協程操作的一致性。
- **自動清理**：內建過期檢查機制，自動清除無效資料。
- **多鍵組合**：支援多參數組合產生唯一鍵值，滿足複雜快取場景。
- **彈性設定**：支援每筆快取資料自訂存活時間，並提供預設存活時間設定。

## 使用方式

### 安裝模組

```bash
go get github.com/yttsai1511/go-cache-toolkit
```

### 導入模組

```go
import "github.com/yttsai1511/go-cache-toolkit"
```

## 使用範例

```go
// 設定快取資料，存活時間為 3 秒
cache.SetWithTTL("Hello, Go!", 3 * time.Second, "greeting")

value, err := cache.Get[string]("greeting")
if err != nil {
	fmt.Println("Error:", err)
	return
}

// 輸出快取內容
fmt.Println("Cache value:", *value)

// 等待超過存活時間
time.Sleep(5 * time.Second)

// 資料已過期
value, err = cache.Get[string]("greeting") 
if err != nil {
	fmt.Println("Error:", err)
	return
}

// Error:
// Item for key greeting has expired
```

## 文件說明
| 方法名稱 | 功能描述 |
| :----- | :----- |
| `Set` | 快速設定快取資料，使用全域預設的存活時間 |
| `SetWithTTL` | 設定快取資料，並自訂每筆資料的存活時間 |
| `Get` | 取得指定鍵的快取資料，支援泛型返回類型 |
| `Delete` | 刪除指定鍵的快取資料 |
| `CleanExpired` | 清理所有過期的快取資料 |
| `SetDefaultTTL` | 設定全域預設存活時間 |
| `GenerateKey` | 根據多參數產生唯一鍵值 |

## 授權

此專案基於 GPLv3 授權條款，詳情請參閱 LICENSE 文件。
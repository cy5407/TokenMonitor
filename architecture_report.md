# TokenMonitor 專案架構分析報告

## 1. 專案概覽

本專案 **TokenMonitor** 是一個多語言、功能全面的 AI Token 使用量監控與分析工具。從檔案結構來看，它似乎旨在提供從即時監控、成本計算、活動分析到報告生成的一整套解決方案。

- **主要技術棧**:
  - **Go (Golang)**: 用於開發核心後端邏輯，包括計算、分析與監控等高性能任務。
  - **JavaScript (Node.js)**: 用於開發客戶端工具、CLI 介面、整合腳本以及可能的 Web 介面。
  - **PowerShell**: 用於 Windows 環境下的安裝、部署與管理腳本。
  - **Shell Script**: 用於類 Unix 環境下的安裝腳本。

## 2. 頂層目錄結構

專案根目錄下的結構清晰，分離了原始碼、腳本、文件、測試和設定檔。

- **`.git/`**: Git 版本控制目錄。
- **`.github/`**: 包含 GitHub Actions 的 CI/CD 工作流程設定 (例如 `release.yml`)。
- **`.kiro/`**: 似乎是 Kiro 開發工具的特定目錄，用於任務管理、規格、掛鉤 (hooks) 和客製化指令。
- **`build/`**: 推測為建置輸出目錄 (目前為空)。
- **`data/`**: 推測用於存放應用程式產生的資料或日誌 (目前為空)。
- **`docs/`**: 存放專案的用戶和開發者文件。
- **`node_modules/`**: Node.js 專案依賴庫。
- **`scripts/`**: 存放各種輔助腳本，用於開發、部署和特定任務。
- **`src/`**: 核心應用程式原始碼。
- **`templates/`**: 包含專案範本，例如用於建立新的 npm 套件。
- **`tests/`**: 包含各種類型的測試檔案，從單元測試資料到整合測試。
- **`config.yaml`**: 主要的應用程式設定檔。
- **`package.json`**: Node.js 專案定義與依賴管理。
- **`README.md`**: 專案的主要說明文件。

## 3. 核心原始碼分析 (`src`)

原始碼目錄 `src` 分為 `go` 和 `js` 兩部分，顯示了雙核心的架構。

### 3.1 Go 應用程式 (`src/go`)

Go 應用程式遵循標準的專案佈局，是整個系統的核心。

- **`main.go`**: 應用程式進入點。
- **`go.mod`, `go.sum`**: Go 模組依賴管理。
- **`cmd/`**: 存放 Cobra CLI 框架的指令。每個 `.go` 檔案 (如 `analyze.go`, `cost.go`, `report.go`) 代表一個子指令，建構出一個功能豐富的 CLI 工具。
- **`internal/`**: 存放專案內部使用的套件，不對外開放。
  - **`analyzer/`**: 活動分析器，用於分析 Token 使用行為。
  - **`calculator/`**: Token 計算器，是 Token 計算的核心邏輯。
  - **`config/`**: 設定檔管理器。
  - **`cost/`**: 成本計算器，根據定價模型計算費用。
  - **`errors/`**: 自定義錯誤處理。
  - **`reporter/`**: 報告產生器，可以輸出 CSV、HTML 等格式。
  - **`monitor/`**: 即時監控邏輯。
  - **其他**: `interfaces`, `storage`, `testutils`, `types`, `utils` 等提供了輔助功能和定義。

### 3.2 JavaScript 應用程式 (`src/js`)

JavaScript 部分似乎提供了更多面向用戶的工具或與 Go 核心互動的介面。

- **`Professional-Token-Cli.js`**: 一個專業級的 Token CLI 工具。
- **`Enhanced-Token-Reporter.js`**: 增強版的 Token 報告工具。
- **`Test-Token-Monitoring.js`**: 用於測試監控功能的腳本。
- **`Token-Monitor-Integration.js`**: 可能用於整合不同部分的監控功能。

## 4. 腳本與自動化

### `scripts/`

此目錄包含大量實用腳本，多數為 PowerShell (`.ps1`)，顯示專案對 Windows 環境的良好支援。

- **`Deploy-Token-Monitor.ps1`**: 部署工具。
- **`Universal-Token-Monitor.js`**: 一個通用的監控腳本。
- **`*hunter.ps1`**: 推測是用於在檔案系統中搜尋和分析 Token 的工具。

### `.github/workflows/`

- **`release.yml`**: 定義了在 GitHub 上建立 release 的自動化流程。

## 5. 測試結構 (`tests`)

測試目錄 `tests` 組織良好，涵蓋了不同層面的測試需求。

- **`ai_content_generation/`**: 存放用於測試的大型文本檔案。
- **`data/`**: 存放測試用的資料，如 JSON 報告和觸發器檔案。
- **`kiro_integration_test/`**: 針對 Kiro 工具的整合測試。
- **`monitor_diagnosis/`**: 用於監控系統的診斷腳本。
- **`reports/`**: 存放大型的測試報告檔案。

## 6. 文件 (`docs`)

`docs` 目錄提供了非常完整的文件，涵蓋了架構、部署、使用指南等。

- **`Architecture.md`**: 系統架構圖。
- **`Deployment-Guide.md`**: 部署指南。
- **`Quick-Start.md`**: 快速入門指南。
- **`Usage-Guide.md`**: 使用手冊。

## 7. Kiro 工具整合 (`.kiro`)

`.kiro` 目錄的存在表明專案深度整合了名為 Kiro 的開發輔助工具。

- **`hooks/`**: 包含在特定事件（如檔案儲存、git commit）觸發的腳本。
- **`specs/`**: 包含專案規格、需求和任務分解文件。
- **`steering/`**: 包含專案的開發規範，如檔案命名約定。
- **`statusbar/`**: 狀態列相關的客製化。

## 8. 架構總結與重構建議

### 架構總結

- **混合式架構 (Hybrid Architecture)**: 專案結合了 Go 的高性能後端與 JavaScript/PowerShell 的靈活腳本/前端，形成了一個強大的混合式系統。Go 負責核心運算，而 JS/PS 負責使用者互動與自動化。
- **CLI 驅動**: Go 應用程式透過 Cobra 框架提供了一個強大的命令列介面，這是與系統核心互動的主要方式。
- **高度可設定與自動化**: 專案擁有 `config.yaml` 設定檔，並透過大量腳本和 GitHub Actions 實現了高度自動化。
- **開發工具整合**: `.kiro` 目錄顯示專案依賴特定的開發工具鏈，這可能提高了開發效率，但也可能增加新成員的學習曲線。

### 給 Claude 的重構建議方向

1.  **統一介面與程式碼整合**:
    -   **問題**: `src/go` 和 `src/js` 中的功能似乎有重疊（例如都有報告和 CLI 功能）。`scripts` 目錄中的 JS 和 PS 腳本也提供了類似的功能。這可能導致功能分散和維護困難。
    -   **建議**: 評估是否可以將所有核心邏輯集中到 Go 應用程式中，並讓 Go 應用程式以 API (例如 RESTful 或 gRPC) 的形式提供服務。JavaScript 和 PowerShell 腳本可以作為該 API 的客戶端，而不是各自實現業務邏輯。

2.  **前後端分離**:
    -   **問題**: 目前的 JS 檔案散落在 `src/js` 和 `scripts` 中，職責不清。
    -   **建議**: 如果有 Web UI 的需求，可以考慮建立一個標準的前端專案 (例如使用 React, Vue)，並將其與 Go 後端完全分離。如果只是 CLI 工具，可以考慮將 `src/js` 中的工具重構為 Go CLI 的子功能或插件。

3.  **設定檔管理**:
    -   **問題**: 專案中存在 `config.yaml`, `package.json`, `.kiro/settings/token-monitor.json` 等多個設定來源。
    -   **建議**: 考慮將設定檔策略統一。例如，讓 Go 應用程式成為唯一的設定讀取者，其他工具透過環境變數或 CLI 參數從 Go 主程式獲取設定。

4.  **腳本語言收斂**:
    -   **問題**: 同時使用 PowerShell 和 JavaScript 進行腳本編寫，增加了維護的複雜性。
    -   **建議**: 評估是否可以將大部分腳本統一為一種語言。例如，使用 Node.js 來編寫跨平台的安裝和部署腳本，以取代 PowerShell 和 Shell Script 的組合。

5.  **依賴關係簡化**:
    -   **問題**: `.kiro` 工具的深度整合可能使專案難以在沒有該工具的環境中進行開發和建置。
    -   **建議**: 評估 `.kiro` 的核心功能，並考慮是否能用更通用的工具（例如 Makefiles, just, 或僅使用 npm scripts）來替代，以降低專案的特殊依賴。

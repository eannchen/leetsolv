[English](./README.md) | [繁體中文](./README.zh-TW.md) | [简体中文](./README.zh-CN.md)

# LeetSolv
[![Release](https://img.shields.io/github/release/eannchen/leetsolv.svg)](https://github.com/eannchen/leetsolv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/eannchen/leetsolv)](https://goreportcard.com/report/github.com/eannchen/leetsolv)
[![CI/CD](https://github.com/eannchen/leetsolv/actions/workflows/ci.yml/badge.svg)](https://github.com/eannchen/leetsolv/actions/workflows/ci.yml)

**LeetSolv** 是一個命令列工具，專為 **資料結構與演算法 (DSA)** 問題複習而設計，带有 **间隔重复 (spaced repetition)** 功能。它由一個客製化的 [SuperMemo 2](https://en.wikipedia.org/wiki/SuperMemo) 演算法驅動，該演算法結合了 **熟悉度**、**重要性** 和 **推理** 等變數。這個方法避免了死記硬背，透過 **刻意練習** 幫助您掌握複雜的演算法。

*想知道這與 Anki 之類的工具有何不同嗎？請參閱 [常見問題](#常見問題)。*

> ***0️⃣ 零依賴哲學**: 完全用純 Go 語言實現，沒有第三方庫、API 或外部工具。甚至避免使用一些標準套件，以便完全控制底層實現——突顯了該專案對基礎知識的關注。有關更多詳細資訊，請參閱 [MOTIVATION.md](document/MOTIVATION.md)*。

**LeetSolv 的工作流程：**

使用 LeetSolv 的建議每日工作流程。

```mermaid
graph LR
A{你正在做什麼？} --> B[一個新問題];
A --> C[LeetSolv 上的一個到期問題];

D[練習這個問題] --> E[新增/更新到 LeetSolv];

B --> D
C --> D

E --> F((SM-2 演算法排定複習));
```

![Demo](document/image/DEMO_header.gif)

## 目錄
- [LeetSolv](#leetsolv)
  - [目錄](#目錄)
  - [快速安裝](#快速安裝)
    - [自動安裝 (Linux/macOS)](#自動安裝-linuxmacos)
    - [手動下載 (所有平台)](#手動下載-所有平台)
    - [驗證安裝](#驗證安裝)
  - [複習排程系統](#複習排程系統)
    - [自適應 SM-2 演算法](#自適應-sm-2-演算法)
    - [到期優先級評分](#到期優先級評分)
    - [間隔增長曲線](#間隔增長曲線)
  - [問題管理](#問題管理)
    - [CRUD + 復原](#crud--復原)
    - [資料隱私 \& 安全](#資料隱私--安全)
    - [CLI 介面](#cli-介面)
  - [使用方法](#使用方法)
  - [配置](#配置)
  - [支援](#支援)
    - [常見問題](#常見問題)
      - [問：為什麼要使用 LeetSolv 而不是 Anki 卡片？](#問為什麼要使用-leetsolv-而不是-anki-卡片)
      - [問：我應該加入所有之前解決過的問題嗎？](#問我應該加入所有之前解決過的問題嗎)
      - [問：使用一段時間後，我累積了太多到期問題。](#問使用一段時間後我累積了太多到期問題)
    - [文件](#文件)
  - [發展規劃](#發展規劃)
    - [自適應 SM-2 演算法](#自適應-sm-2-演算法-1)
    - [安裝 \& 執行](#安裝--執行)
    - [功能](#功能)
  - [許可證](#許可證)

## 快速安裝

### 自動安裝 (Linux/macOS)

下載並執行安裝腳本

```bash
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash
```

或者先下載，然後執行

```bash
wget https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh
chmod +x install.sh
./install.sh
```

> *對於 Windows，請使用下面的 [手動下載](#手動下載-所有平台) 方法。未來計劃支援 Windows。*

### 手動下載 (所有平台)
1. 前往 [Releases](https://github.com/eannchen/leetsolv/releases)
2. 下載適用於您平台的二進位檔案：
   - **Linux**: `leetsolv-linux-amd64` 或 `leetsolv-linux-arm64`
   - **macOS**: `leetsolv-darwin-amd64` 或 `leetsolv-darwin-arm64`
   - **Windows**: `leetsolv-windows-amd64.exe` 或 `leetsolv-windows-arm64.exe`

### 驗證安裝
```bash
leetsolv version
leetsolv help
```

> *有關詳細的安裝和配置說明，請參閱 [INSTALL.md](document/INSTALL.md)*

## 複習排程系統

### 自適應 SM-2 演算法

新增問題時，LeetSolv 會應用具有自定義因素（**熟悉度**、**重要性**和**推理**）的 SM-2 演算法，以使用容易度因子計算下一次複習日期。

- **容易度因子**：用於計算下一次複習間隔的經典 SM-2 演算法決定因素。
- **熟悉度量表**：5 級熟悉度評估（`VeryHard`、`Hard`、`Medium`、`Easy`、`VeryEasy`）用於個人評分。
- **重要性量表**：4 層優先級系統（`Low`、`Medium`、`High`、`Critical`）用於問題優先級排序。
- **推理量表**：3 級記憶評估（`Reasoned`、`Partial`、`Full recall`），對薄弱的推理進行懲罰。此設計強化了通過推理而不是死記硬背來掌握資料結構和演算法的目標。
- **到期懲罰（可選）**：對過期複習問題的自動懲罰系統。
- **隨機化（可選）**：隨機排程變化以防止日期過度擬合。

```mermaid
graph TD
    A[將 DSA 問題新增到 LeetSolv] --> B[演算法應用調整]
    B --> C[熟悉度量表]
    B --> D[重要性量表]
    B --> E[推理量表]
    B --> F[到期懲罰，可選]
    B --> G[隨機化，可選]


    H[演算法使用 SM-2 容易度因子計算]
    C --> H
    D --> H
    E --> H
    F --> H
    G --> H

    H --> I[確定下一次複習]
```

> *在此處查看以了解有關該演算法的更多資訊：[間隔增長曲線](#間隔增長曲線)*

### 到期優先級評分
使用 SM-2 時，由於用戶有不同的排程和學習習慣，因此到期複習很容易累積。為了應對這一挑戰，LeetSolv 引入了一種到期優先級評分功能，**允許用戶根據優先級分數對到期問題進行優先級排序**。

- **多因素評分**：重要性、熟悉度、過期天數、複習次數和容易度因子的組合決定了到期優先級。

```mermaid
graph LR
    A[問題] --> B[優先級分數計算]
    B --> C[重要性權重]
    B --> D[過期權重]
    B --> E[熟悉度權重]
    B --> F[複習次數]
    B --> G[容易度因子]

    C --> H[最終優先級分數]
    D --> H
    E --> H
    F --> H
    G --> H

    H --> I[按分數排序]
    I --> J[前 K 個到期問題]
    I --> K[前 K 個即將到來的問題]
```

> *預設情況下，優先級分數使用以下公式計算：(1.5×重要性)+(0.5×過期天數)+(3.0×熟悉度)+(-1.5×複習次數)+(-1.0×容易度因子)*

**到期優先級列表演示：**
![Demo](document/image/DEMO_due_scoring.gif)

### 間隔增長曲線

LeetSolv 根據重要性、熟悉度和推理調整 SM-2 間隔。**重要問題會更頻繁地複習，而較容易的問題的間隔會逐漸變長。**

以下圖表展示了不同重要性級別的複習間隔如何隨時間增長，顯示了預設的增長模式：

![SM2 Critical](document/image/SM2_CRITICAL.png)
**重要問題**：最短的間隔，頻繁的複習，以確保掌握最重要的概念。

![SM2 High](document/image/SM2_HIGH.png)
**高重要性**：適度的間隔，平衡了頻率和效率，適用於重要問題。

![SM2 Medium](document/image/SM2_MEDIUM.png)
**中等重要性**：標準間隔，遵循經典的 SM-2 漸進式練習。

![SM2 Low](document/image/SM2_LOW.png)
**低重要性**：較長的間隔，適用於不需要頻繁複習的問題。

> *未來的迭代可能會允許直接在配置文件中修改演算法參數。*

## 問題管理

### CRUD + 復原
- **摘要視圖**：總問題、到期問題和即將到來的問題的概述，帶有[到期優先級評分](#到期優先級評分)。
- **列表視圖**：列出所有問題，帶有分頁。
- **基於 Trie 的搜尋 \& 篩選**：通過關鍵字、熟悉度、重要性、複習次數和到期狀態進行快速搜尋和篩選。
- **獲取問題詳細資訊**：檢索有關特定問題的詳細資訊。
- **新增/更新問題**：建立或修改帶有 URL 和註釋的問題。
- **移除問題**：按 ID 和 URL 刪除問題。
- **歷史記錄 \& 復原**：追蹤更改並還原最近的操作。

```mermaid
graph TD
    A[用戶命令] --> B{命令類型}
    B -->|搜尋/列表| C[應用篩選器 \& 搜尋] --> D[基於 Trie 的搜尋] --> E[分頁結果]
    B -->|新增/更新| F[通過 SM-2 演算法排定複習] --> I[更新資料 \& 索引]
    B -->|移除| G[移除資料 \& 索引]
    B -->|復原| H[從歷史記錄還原資料 \& 索引]

    J[新增到歷史記錄]
    I --> J
    G --> J
```

**搜尋、歷史記錄、刪除、復原功能演示：**
![Demo](document/image/DEMO_mgmt.gif)

### 資料隱私 \& 安全

- **無資料收集**：LeetSolv 不會將用戶資料上傳到網際網路。
- **原子寫入**：所有更新都使用帶有原子替換的臨時檔案，以保證一致性並防止資料丟失。

```mermaid
graph LR
    A[寫入新資料] -->|到| B[臨時檔案]
    B --> C[重新命名臨時檔案]
    C -->|替換| D[原始檔案]
```

### CLI 介面
- **互動模式**：用戶和程式之間的輪流對話。
- **批次模式**：直接從命令列執行命令。
- **別名支援**：方便的直觀快捷方式（例如，`ls`、`rm`、`cfg`）。
- **分頁**：用於大型問題集的有效導航。
- **清晰輸出**：結構化的、顏色編碼的 CLI 輸出。

**分頁演示：**

![Demo](document/image/DEMO_pagination.gif)

## 使用方法

LeetSolv 可以互動方式執行，也可以直接從終端機傳遞命令來執行。

```bash
# 啟動互動模式
leetsolv

# 或直接執行命令
leetsolv add https://leetcode.com/problems/two-sum
leetsolv status

# 取得說明
leetsolv help
```

[檢視完整使用指南 (USAGE.md)](document/USAGE.md)

## 配置

LeetSolv 可以使用環境變數或 JSON 設定檔進行自訂。這允許您更改檔案路徑、演算法參數和評分權重。

如需所有可用選項、預設值和範例的完整清單，請參閱詳細的設定指南：

[檢視完整設定指南 (CONFIGURATION.md)](document/CONFIGURATION.md)

## 支援

### 常見問題

#### 問：為什麼要使用 LeetSolv 而不是 Anki 卡片？

答：Anki 非常適合**記憶幾秒鐘就能完成的事實**，**但對於 DSA 來說，它可能會適得其反**。LeetSolv 的自訂 SM-2 演算法 **延長** 間隔複習時間，並使用您對**推理**、**熟悉度**和問題**重要性**的輸入來建立一個加深您解決問題能力的排程，而不是僅僅檢查您是否記住了答案。

👉 提示：Anki 適合用來記憶事實，而 LeetSolv 則是用來安排深度、刻意的解題練習。

#### 問：我應該加入所有之前解決過的問題嗎？

答：**不應該。** LeetSolv 不是一個已解決問題的資料庫，而是一個間隔重複排程器。
只加入你真正想要複習的問題。排程演算法依賴你加入問題的日期來計算複習時間，所以大量加入會創建一個不切實際的排程和一個壓倒性的到期問題堆。

👉 提示：如果你想複習一個幾個月前解決的舊問題，先重新解決它，然後在那天將它添加到 LeetSolv。這樣「上次看到」的日期才是準確的。

#### 問：使用一段時間後，我累積了太多到期問題。

答：這是SM-2 演算法的**本質** —— 如果你跳過幾天或一次加入許多問題，到期列表可能會迅速增長。
為了使其易於管理，LeetSolv 引入了 **[到期優先級評分](#到期優先級評分)**，它按重要性、熟悉度、逾期天數、複習次數和容易程度對到期問題進行排名。
不要一次清除所有內容，只需專注於 **最高優先級的問題** 即可。其餘的可以安全地等到以後。

👉 提示：一旦你掌握了一個問題並重新解決了幾次，你就可以安全地將其從 LeetSolv 中刪除。目標不是永遠追蹤所有內容，而是專注於仍然需要間隔練習的問題。


### 文件
- **[INSTALL.md](document/INSTALL.md)**：完整的安裝指南，包含疑難排解
- **[USAGE.md](document/USAGE.md)**：命令行使用指南
- **[CONFIGURATION.md](document/CONFIGURATION.md)**：配置選項和環境變數
- **[DEVELOPMENT_GUIDE.md](document/DEVELOPMENT_GUIDE.md)**：開發工作流程、CI/CD 和貢獻指南
- **[MOTIVATION.md](document/MOTIVATION.md)**：專案動機和設計說明
- **本 README**：專案概述和快速入門

## 發展規劃

我們的發展規劃是以使 LeetSolv 成為刻意練習的强大工具為目標。有建議嗎？隨時 [開啟 issue](https://github.com/eannchen/leetsolv/issues)！

### 自適應 SM-2 演算法

- 分離 `熟悉度` 為 `推理熟悉度` 和 `編碼熟悉度`
  - 閱讀 [ALGORITHM_ROADMAP.md](document/ALGORITHM_ROADMAP.md) 了解更多細節
- 使 SM-2 演算法可由用戶自定義

### 安裝 \& 執行

- 支援通過套件管理器（Homebrew、Chocolatey、apt）安裝
- 支援作為 Docker 容器執行
- 支援自動更新程式功能
- 新增 Windows 安裝腳本

### 功能

- 在問題條目中提供「提示」欄位
- 提供每日新增限制功能
- 提供標記功能
- 提供匯出功能
- 新增對來自其他平台的 DSA 問題的支援
- 實作模糊搜尋功能

## 許可證

本專案已根據 [LICENSE](LICENSE) 檔案中指定的條款獲得許可。

---

**LeetSolv** - 一個基於自訂 SM-2 演算法的 DSA 間隔重複 CLI 工具，幫助你透過刻意練習掌握演算法。

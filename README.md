# 塔防三国（Three Kingdoms: Siege）

- **游戏中文名：**塔防三国
- **游戏英文名：**Three Kingdoms: Siege
- **开发者：**重庆沛雨霖科技有限公司（PYLinTech）

## 项目概述

《塔防三国》是一款以三国历史为背景的三人实时联机 PvPvP 策略塔防游戏，融合城池防御、武将统兵、军团作战及领土争夺玩法。游戏采用真实三维场景，通过正交摄像机呈现 2.5D 等距视角，兼顾立体表现与策略操作体验。

## 最终技术选型

项目采用 **React 19 + TypeScript + React Three Fiber v9 + Three.js + Zustand + Vite**。桌面端使用 **Tauri 2**，移动端使用 **Capacitor**，后端采用 **Go + WebSocket**。

| 技术 | 主要职责 |
|---|---|
| React 19 | 游戏 UI、页面及交互组件 |
| TypeScript | 前端统一开发语言 |
| Three.js | 三维场景、模型、动画与特效渲染 |
| React Three Fiber v9 | 以 React 组件组织 3D 场景 |
| Drei | 摄像机控制、模型加载等辅助能力 |
| Tailwind CSS | UI 样式与跨屏幕适配 |
| Zustand | 客户端状态管理 |
| Vite | 开发环境与项目构建 |
| Tauri 2 | Windows、macOS 桌面客户端 |
| Capacitor | Android、iOS 移动客户端 |
| WebSocket | 实时联机通信 |
| Go | 权威游戏服务器及战斗逻辑 |

## 技术架构

### 三维战场

使用 Three.js 实现真实三维地图、城池、防御塔、武将、士兵与战斗特效，通过 React Three Fiber 组织场景组件。采用正交摄像机呈现 2.5D 等距视角，支持地图缩放和平移、单位选择、建筑部署及战斗动画。

针对大量单位同屏作战，采用 InstancedMesh、模型复用、分级细节（LOD）和动画插值等方式优化渲染性能。

### 游戏 UI

使用 React、HTML、CSS 构建登录界面、游戏大厅、武将选择、建造菜单、技能栏及战斗结算界面。UI 与 3D 渲染层相互独立，通过 Zustand 状态和统一事件机制通信；使用自定义新中式美术资源实现复杂布局和交互。

### 多人联机

客户端通过 WebSocket 与 Go 游戏服务器保持实时通信。服务器采用权威架构，负责房间管理、战斗判定、资源结算、单位状态及胜负判定；客户端负责提交操作输入、场景渲染和状态插值，以保证对战公平及状态一致。

### 跨平台发布

共享核心游戏逻辑、渲染代码和 UI 组件，并按平台适配操作方式、屏幕尺寸与硬件性能：

- **浏览器：** Vite 构建，支持主流桌面及移动浏览器。
- **桌面：** Tauri 2 封装 Windows、macOS 客户端。
- **移动：** Capacitor 封装 Android、iOS 客户端。

## 优势与注意事项

### 方案优势

1. React、TypeScript、Three.js 生态成熟，文档和社区资源丰富。
2. TSX、TypeScript、CSS、JSON 易于 AI 辅助生成、修改和重构。
3. 真实 3D 模型结合正交摄像机，可实现适合策略塔防的 2.5D 视觉效果。
4. UI 不依赖游戏引擎专用控件，设计自由度高，便于维护和调整。
5. 多平台复用核心代码，降低开发与长期维护成本。

### 技术注意事项

React Three Fiber 是三维渲染框架，并非完整游戏引擎；地图编辑、战斗系统、寻路及单位管理等需自行设计。Tauri 和 Capacitor 主要通过系统 WebView 运行，移动端图形性能可能受设备与 WebView 环境影响，需针对 GPU、内存、渲染分辨率及同屏单位数专项优化。

## 项目结构

- `frontend/`：前端应用（目标为 React + TypeScript + React Three Fiber + Vite）。
- `backend/`：Go 后端服务及权威战斗逻辑。

## 许可证与素材授权

- **代码：**采用 GNU Affero General Public License v3.0 only（AGPL-3.0-only），详见 [`LICENSE`](LICENSE)。
- **Logo、品牌与原创美术：**除非有单独授权说明，否则保留权利；代码许可证不授予这些素材的使用权。详见 [`ASSET-LICENSE.md`](ASSET-LICENSE.md)。
- **第三方素材：**继续遵循各自原始许可证及署名要求，本项目不对其重新授权。

## 音频资源规范

项目音频统一使用 **M4A 容器与 AAC-LC 编码**。

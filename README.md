# 塔防三国

英文名：Three Kingdoms: Siege  
开发者：重庆沛雨霖科技有限公司（PYLinTech）

## 项目简介

《塔防三国》（Three Kingdoms: Siege）是一款以三国历史为背景的三人实时联机 PvPvP 策略塔防游戏，融合城池防御、武将统兵、军团作战及领土争夺玩法。游戏采用三维场景与正交摄像机，呈现 2.5D 等距视角，兼顾立体表现与策略操作体验。

## 技术选型

项目规划采用 React 19、TypeScript、React Three Fiber v9、Three.js、Zustand 和 Vite 构建前端；桌面端计划使用 Tauri 2，移动端计划使用 Capacitor，后端采用 Go 和 WebSocket。

| 技术 | 用途 |
| --- | --- |
| React 19 | 游戏界面与交互组件 |
| TypeScript | 前端开发语言 |
| Three.js | 三维场景、模型、动画与特效渲染 |
| React Three Fiber v9 | 以 React 组件组织三维场景 |
| Drei | 摄像机控制、模型加载等辅助功能 |
| Tailwind CSS | 界面样式与跨屏幕适配 |
| Zustand | 客户端状态管理 |
| Vite | 前端开发与构建 |
| Tauri 2 | Windows、macOS 桌面客户端 |
| Capacitor | Android、iOS 移动客户端 |
| WebSocket | 实时联机通信 |
| Go | 游戏服务器与战斗逻辑 |

以上为项目技术规划，具体实现以代码仓库当前内容为准。

## 技术架构

### 三维战场

使用 Three.js 实现地图、城池、防御塔、武将、士兵与战斗特效，并通过 React Three Fiber 组织场景组件。采用正交摄像机呈现 2.5D 等距视角，支持地图缩放和平移、单位选择、建筑部署及战斗动画。

针对大量单位同屏作战，可采用 InstancedMesh、模型复用、分级细节（LOD）和动画插值等方式优化渲染性能。

### 游戏界面

使用 React、HTML 和 CSS 构建登录界面、游戏大厅、武将选择、建造菜单、技能栏及战斗结算界面。界面与三维渲染层相互独立，通过 Zustand 状态和统一事件机制通信。

### 多人联机

客户端通过 WebSocket 与 Go 游戏服务器通信。服务器负责房间管理、战斗判定、资源结算、单位状态及胜负判定；客户端提交操作输入并负责场景渲染和状态插值。

### 跨平台发布

项目计划共享核心游戏逻辑、渲染代码和界面组件，并按平台适配操作方式、屏幕尺寸与硬件性能：

- 浏览器：使用 Vite 构建，支持主流桌面及移动浏览器。
- 桌面：计划使用 Tauri 2 封装 Windows 和 macOS 客户端。
- 移动：计划使用 Capacitor 封装 Android 和 iOS 客户端。

## 技术考量

React Three Fiber 是三维渲染框架，并非完整游戏引擎；地图编辑、战斗系统、寻路及单位管理等需要自行设计。Tauri 和 Capacitor 主要通过系统 WebView 运行，移动端图形性能可能受设备与 WebView 环境影响，需针对 GPU、内存、渲染分辨率及同屏单位数进行优化。

## 项目结构

```text
.
├── frontend/   # 前端应用
└── backend/    # Go 后端服务
```

## 许可证与素材授权

- 代码采用 GNU Affero 通用公共许可证第 3 版（仅此版本，AGPL-3.0-only），详见 [LICENSE](LICENSE)。
- Logo、品牌与原创美术除另有授权外均保留权利；代码许可证不授予这些素材的使用权，详见 [ASSET-LICENSE.md](ASSET-LICENSE.md)。
- 第三方素材遵循其各自的许可证及署名要求，本项目不对其重新授权。

## 音频资源规范

项目音频统一使用 M4A 容器与 AAC-LC 编码。

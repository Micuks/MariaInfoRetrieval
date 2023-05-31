# MariaInfoRetrieval

## 安装和配置
### 后端
1. 确保已经安装Go. 如果没有, 请安装
2. 克隆该项目: `git clone https://github.com/Micuks/MariaInfoRetrieval`
3. 进入后端文件夹 `cd backend`
4. 运行服务, 并根据报错安装对应的包, 主要有`gin`和`logrus`: `go run main.go`

后端现在应该运行在`http://localhost:9011`

### 前端
1. 确保已经安装Node.js和npm/yarn/pnpm
2. 进入前端文件夹: `cd frontend`
3. 安装需要的Node.js包: `npm install`或者`yarn install`或者`pnpm install`
4. 运行前端服务器: `npm start`或者`yarn start`或者`pnpm start`

前端现在应该运行在`http://localhost:9001`

### Python微服务
1. 确保已经安装Python 3
2. 进入微服务文件夹: `cd image_micro_service`
3. 创建并进入虚拟环境: `python -m venv .env && source .env/bin/activate`
4. 安装需要的包: `pip install -r requirements.txt`
5. 运行微服务: `python server.py`

按照报错提示安装各种依赖后, 微服务现在应该运行在`http://localhost:9021`

## 使用方法
一旦前端和后端成功运行，您就可以开始使用这个系统了。

### 文档搜索
1. 打开浏览器并导航至 `http://localhost:9001`

2. 在搜索栏输入您的查询内容，然后点击"搜索"。

3. 系统将返回与您的查询相关的文档列表。将显示文档的相关性、标题、匹配内容、URL和日期。

4. 您可以点击文档标题查看文档的全部内容。

5. 您也可以对搜索结果提供反馈，这将用于改进文档的排名。

### 图像搜索
1. 点击"图像搜索"选项卡。

2. 上传一张图片，系统将返回与该图片相关的关键词列表。

### 反馈
您的反馈对我们改善系统非常宝贵。请随时通过点击"反馈"按钮提供反馈。

### 贡献
欢迎贡献.

### 许可
该项目在MIT许可下授权 - 详情请见LICENSE.md文件。
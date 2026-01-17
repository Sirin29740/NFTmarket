import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
// import axios from 'axios' // 如果你还没有创建 API 客户端，可以暂时注释

const app = createApp(App)

// 注册路由
app.use(router)

// 注册全局属性 (如果需要 Axios)
// app.config.globalProperties.$axios = axios

// 挂载应用到 HTML 模板中的 <div id="app"></div> 元素
app.mount('#app')
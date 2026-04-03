import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 1. 导入 Wagmi 相关
import { WagmiPlugin } from '@wagmi/vue'
import { config } from './wagmi'

// 2. 导入 Vue Query 相关 (修复点：使用 VueQueryPlugin)
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'

const app = createApp(App)

// 3. 初始化 QueryClient
const queryClient = new QueryClient()

app.use(router)

// 4. 注册 Wagmi 插件
app.use(WagmiPlugin, { config })

// 5. 注册 Vue Query 插件 (修复点：代替之前的 QueryClientProvider)
app.use(VueQueryPlugin, { queryClient })

app.mount('#app')
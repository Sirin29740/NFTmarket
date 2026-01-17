// frontend/src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'

// 导入页面组件
import LoginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'
import ProfilePage from '../views/ProfilePage.vue'

// 假设认证状态检查函数
// function isAuthenticated() {
//   return localStorage.getItem('token') !== null
// }

const router = createRouter({
    // 使用 History 模式，对应 Go 后端的 NoRoute 配置
    history: createWebHistory(import.meta.env.BASE_URL),

    routes: [
        {
            path: '/',
            redirect: '/login' // 默认重定向到登录页
        },
        {
            path: '/login',
            name: 'Login',
            component: LoginPage
        },
        {
            path: '/register',
            name: 'Register',
            component: RegisterPage
        },
        {
            path: '/profile',
            name: 'Profile',
            component: ProfilePage,
            meta: { requiresAuth: true }
        }
    ]
})
router.beforeEach((to, from, next) => {
    // 检查目标路由是否需要认证
    const requiresAuth = to.meta.requiresAuth;
    // 使用 'jwt_token' 键名，与登录/拦截器保持一致
    const token = localStorage.getItem('jwt_token');

    if (requiresAuth && !token) {
        // 如果需要认证但没有 Token
        console.warn("Auth Guard: Missing token, redirecting to login.");
        // 重定向到登录页，并带上目标路径作为查询参数，方便登录后返回
        next({ name: 'Login', query: { redirect: to.fullPath } });
    } else {
        // 允许通过 (有 Token 或不需要认证)
        next();
    }
});

export default router
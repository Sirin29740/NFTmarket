

import axios from 'axios';
// 如果你想使用 Vue Router 的 router.push，你需要确保在这个模块中可以访问到 router 实例。
// 简单示例使用 window.location.href 进行跳转。

const api = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json'
    }
});

// -----------------------------------------------------
// 修正请求拦截器：确保 Token trim() 消除格式问题
// -----------------------------------------------------
api.interceptors.request.use(config => {
    const token = localStorage.getItem('jwt_token');

    if (token) {
        // 增加 trim() 确保 Token 字符串不含首尾空格，解决 'Token 格式错误'
        const trimmedToken = token.trim();
        config.headers.Authorization = `Bearer ${trimmedToken}`;
    }
    return config;
});


// -----------------------------------------------------
// 修正响应拦截器：统一处理 401 并跳转
// -----------------------------------------------------
api.interceptors.response.use(
    (response) => {
        return response; // 200 OK，不执行任何操作
    },
    (error) => {
        if (error.response && error.response.status === 401) {

            // 1. 打印日志确认 Token 被移除
            console.log("401 Unauthorized received. Clearing token and redirecting.");

            // 2. 移除 Token：这是导致刷新丢失的原因，但也是处理过期Token的必要操作
            localStorage.removeItem('jwt_token');
            localStorage.removeItem('user_info');

            // 3. 执行跳转（使用原生的方式，或者如果你能引入 router 实例，使用 router.push）
            // 阻止当前页面继续加载，并跳转到登录页
            if (window.location.pathname !== '/login') {
                // 避免无限重定向循环
                window.location.href = '/login';
            }

            // 阻止 Promise 链继续执行，防止 Profile.vue 中的 catch 块被触发
            return Promise.reject(error);
        }

        // 其他错误 (500, 404 等)
        return Promise.reject(error);
    }
);

export default api;
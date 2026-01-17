<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2 class="title">✨ 用户登录</h2>
      <p class="subtitle">欢迎回来，请输入您的账号信息</p>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名 / 邮箱</label>
          <input
              type="text"
              id="username"
              v-model="form.username"
              placeholder="请输入用户名或邮箱"
              required
              autocomplete="username"
          />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input
              type="password"
              id="password"
              v-model="form.password"
              placeholder="请输入密码"
              required
              autocomplete="current-password"
          />
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>

        <button type="submit" :disabled="isLoading">
          {{ isLoading ? '验证并登录...' : '立即登录' }}
        </button>
      </form>

      <p class="link-footer">
        还没有账号？<router-link to="/register" class="auth-link">免费注册</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
// 假设 'api' 是您的 axios/fetch 封装
import api from '@/api/client';

const router = useRouter();
const form = ref({ username: '', password: '' });
const error = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  error.value = '';
  isLoading.value = true;
  try {
    const response = await api.post('/login', form.value);
    const { token, user } = response.data;

    // 存储 Token 和用户信息
    localStorage.setItem('jwt_token', token.trim());
    localStorage.setItem('user_info', JSON.stringify(user));

    // 登录成功，跳转到个人资料页
    router.push('/profile');
  } catch (err) {
    // 错误处理逻辑
    if (err.response && err.response.data && err.response.data.message) {
      error.value = err.response.data.message;
    } else {
      console.error(err);
      error.value = '登录请求失败，请检查网络连接或稍后重试。';
    }
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
/* 定义 CSS 变量 */
:root {
  --primary-color: #4a90e2;
  --primary-hover: #3a76c4;
  --text-dark: #333333;
  --card-bg: #fff;
  --error-color: #e54c3c;
  --bg-light: #f5f7fa;
  --shadow-dark: 0 6px 18px rgba(0,0,0,0.08);
}

/* 整体布局 */
.auth-wrapper {
  min-height: 30vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--bg-light);
  padding: 20px;
}

/* 登录卡片 */
.auth-card {
  width: 100%;
  max-width: 380px;
  padding: 30px 35px;
  background: var(--card-bg);
  border-radius: 10px;
  box-shadow: var(--shadow-dark);
  text-align: center;
}

.title {
  color: var(--primary-color);
  margin-bottom: 5px;
  font-size: 26px;
  font-weight: 700;
}

.subtitle {
  color: #888;
  margin-bottom: 25px;
  font-size: 14px;
}

/* 表单组间距优化 */
.form-group {
  margin-bottom: 20px; /* 增加间距 */
  text-align: left;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-size: 13px;
  color: var(--text-dark);
  font-weight: 500;
}

/* 输入框样式 */
input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 11px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 15px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.2);
  outline: none;
}

/* 登录按钮样式 (重点优化部分) */
button[type="submit"] {
  width: 100%;
  padding: 14px; /* 增加内边距，让按钮更高 */
  margin-top: 25px; /* 增加与上方输入框的间距，使其更突出 */
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 17px; /* 字体稍微增大 */
  font-weight: 700; /* 加粗字体 */
  letter-spacing: 0.5px; /* 增加字母间距，提升视觉强度 */
  transition: background-color 0.3s ease, transform 0.1s ease;
  cursor: pointer;
}

button[type="submit"]:hover:not(:disabled) {
  background-color: #4a90e2;
  transform: translateY(-1px); /* 悬停时轻微上浮效果 */
  box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3); /* 添加阴影 */
}

button[type="submit"]:disabled {
  background-color: #a0a0a0;
  cursor: not-allowed;
  opacity: 0.8;
}

/* 错误信息样式 */
.error-message {
  color: var(--error-color);
  margin: 15px 0 0 0; /* 调整间距 */
  padding: 10px;
  border: 1px solid var(--error-color);
  background-color: #fef0f0;
  border-radius: 6px;
  font-size: 13px;
  text-align: left;
}

/* 底部链接 */
.link-footer {
  margin-top: 30px; /* 增加与按钮的间距 */
  font-size: 14px;
  color: #666;
}

.auth-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 600;
}

.auth-link:hover {
  text-decoration: underline;
}
</style>
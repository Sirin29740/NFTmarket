<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2 class="title">🚀 新用户注册</h2>
      <p class="subtitle">创建您的账号，进入 CryptoMarket</p>

      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
              type="text"
              id="username"
              v-model="form.username"
              placeholder="必填，用于登录和展示"
              required
              autocomplete="new-username"
          />
        </div>

        <div class="form-group">
          <label for="phone">电话 (选填)</label>
          <input
              type="tel"
              id="phone"
              v-model="form.phone"
              placeholder="可选，用于找回密码"
              autocomplete="tel"
          />
        </div>

        <div class="form-group">
          <label for="email">邮箱 (选填)</label>
          <input
              type="email"
              id="email"
              v-model="form.email"
              placeholder="可选，接收通知和验证"
              autocomplete="email"
          />
        </div>

        <div class="form-group">
          <label for="password">密码 (最少8位)</label>
          <input
              type="password"
              id="password"
              v-model="form.password"
              placeholder="请设置您的登录密码"
              required
              minlength="8"
              autocomplete="new-password"
          />
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>
        <p v-if="success" class="success-message">{{ success }}</p>

        <button type="submit" :disabled="isLoading" class="primary-action-btn">
          {{ isLoading ? '正在创建账号...' : '注册并登录' }}
        </button>
      </form>

      <p class="link-footer">
        已有账号？<router-link to="/login" class="auth-link">立即登录</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '@/api/client';

const router = useRouter();
const form = ref({ username: '', phone: '', email: '', password: '' });
const error = ref('');
const success = ref('');
const isLoading = ref(false);

const handleRegister = async () => {
  error.value = '';
  success.value = '';
  isLoading.value = true;

  try {
    const payload = {
      username: form.value.username,
      password: form.value.password,
      email: form.value.email || undefined,
      phone: form.value.phone || undefined
    };

    const response = await api.post('/register', payload);

    const { token, user } = response.data;

    if (token) localStorage.setItem('jwt_token', token.trim());
    localStorage.setItem('user_info', JSON.stringify(user));

    success.value = '注册成功！正在跳转...';
    setTimeout(() => {
      router.push('/profile');
    }, 1000);

  } catch (err) {
    if (err.response && err.response.data && err.response.data.message) {
      error.value = err.response.data.message;
    } else {
      console.error(err);
      error.value = '注册失败，请检查数据格式或网络。';
    }
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
:root {
  --primary-color: #4a90e2;
  --primary-hover: #3a76c4;
  --text-dark: #333333;
  --card-bg: #fff;
  --error-color: #e54c3c;
  --success-color: #28a745;
  --bg-light: #f5f7fa;
  --shadow-dark: 0 6px 18px rgba(0,0,0,0.08);
}

.auth-wrapper {
  min-height: 30vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--bg-light);
  padding: 20px;
}

.auth-card {
  width: 100%;
  max-width: 380px; /* 缩小卡片宽度，更紧凑 */
  padding: 30px 35px; /* 缩小内边距 */
  background: var(--card-bg);
  border-radius: 16px;
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
  margin-bottom: 20px; /* 缩小底部间距 */
  font-size: 14px;
}

.form-group {
  margin-bottom: 15px; /* 缩小表单间距 */
  text-align: left;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-size: 13px;
  color: var(--text-dark);
  font-weight: 500;
}

input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 15px;
  box-sizing: border-box;
}

input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.2);
  outline: none;
}

.primary-action-btn {
  width: 100%;
  padding: 12px;
  margin-top: 15px;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.primary-action-btn:hover:not(:disabled) {
  background-color: #4a90e2;
  transform: translateY(-1px);
}

.primary-action-btn:disabled {
  background-color: #a0a0a0;
  cursor: not-allowed;
  opacity: 0.8;
}

.error-message {
  color: var(--error-color);
  text-align: center;
  margin: 12px 0;
  padding: 8px;
  border: 1px solid var(--error-color);
  background-color: #fef0f0;
  border-radius: 4px;
  font-size: 13px;
}

.success-message {
  color: var(--success-color);
  text-align: center;
  margin: 12px 0;
  padding: 8px;
  border: 1px solid var(--success-color);
  background-color: #f0fff4;
  border-radius: 4px;
  font-size: 13px;
}

.link-footer {
  margin-top: 20px;
  font-size: 13px;
  color: #666;
}

.auth-link {
  color: var(--primary-color);
  text-decoration: none;
}

.auth-link:hover {
  text-decoration: underline;
}
</style>

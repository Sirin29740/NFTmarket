<template>
  <div id="app-layout">
    <!-- 头部导航 -->
    <header class="app-header">
      <h1 class="logo-title">
        <span class="highlight">Crypto</span>Market
      </h1>

      <nav class="main-nav">
        <router-link to="/login" class="nav-link">登录</router-link>
        <router-link to="/register" class="nav-link primary-btn">注册</router-link>
        <router-link to="/profile" class="nav-link profile-link">个人资料</router-link>
      </nav>
    </header>

    <!-- 页面内容 -->
    <main class="app-main">
      <div class="page-container">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
// Layout.vue 不需要逻辑
</script>

<style scoped>
/* 主题色 */
:root {
  --primary-color: #4a90e2;
  --primary-color-dark: #3a76c4;
  --text-light: #ffffff;
  --text-dark: #333333;
  --bg-light: #f5f7fa;
  --shadow-dark: 0 6px 18px rgba(0,0,0,0.08);
}

/* 布局基础 */
#app-layout {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
}

/* --- 头部 --- */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 40px;
  background-color: var(--text-light);
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.logo-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-dark);
}

.highlight {
  color: var(--primary-color);
}

/* --- 导航 --- */
.main-nav {
  display: flex;
  gap: 15px;
  align-items: center;
}

.nav-link {
  text-decoration: none;
  color: var(--text-dark);
  font-weight: 500;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.3s ease;
}
.nav-link.router-link-active {
  color: var(--primary-color);
}
.nav-link:hover {
  background-color: rgba(0,0,0,0.05);
}

.primary-btn {
  background-color: var(--primary-color);
  color: var(--text-light);
  font-weight: 600;
}
.primary-btn:hover {
  background-color: var(--primary-color-dark);
}

.profile-link {
  border: 1px solid rgba(0,0,0,0.2);
  padding: 6px 12px;
}

/* --- 页面容器 --- */
.app-main {
  min-height: calc(100vh - 70px);
  background-color: var(--bg-light);
  display: flex;
  justify-content: center;
  padding-top: 40px;
}

.page-container {
  width: 100%;
  max-width: 900px;
  padding: 20px;
}
</style>

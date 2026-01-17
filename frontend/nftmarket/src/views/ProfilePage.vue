<template>
  <div class="profile-card">
    <div class="profile-title">
      <span class="icon">👤</span>
      <h2>个人资料</h2>
      <p>您的账户信息总览</p>
    </div>

    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>正在加载用户资料...</p>
    </div>

    <div v-else-if="error || !userInfo" class="error-state">
      <p>{{ error || '获取资料失败，请检查登录状态。' }}</p>
      <button @click="handleLogout" class="logout-btn">
        {{ error ? '重新登录' : '去登录' }}
      </button>
    </div>

    <div v-else class="profile-info">
      <div class="info-item">
        <label>用户ID</label>
        <span>{{ userInfo.user_id }}</span>
      </div>
      <div class="info-item">
        <label>用户名</label>
        <span>{{ userInfo.username }}</span>
      </div>
      <div class="info-item">
        <label>邮箱</label>
        <span>{{ userInfo.email || '未设置' }}</span>
      </div>
      <div class="info-item">
        <label>电话</label>
        <span>{{ userInfo.phone || '未设置' }}</span>
      </div>

      <div class="upload-section">
        <h3>上传头像 / 文件 (IPFS)</h3>
        <input type="file" ref="fileInput" @change="handleFileChange" accept="image/*" />

        <button
            @click="uploadFile"
            :disabled="isUploading || !selectedFile"
            class="upload-btn"
        >
          {{ isUploading ? '上传中...' : '上传到 IPFS' }}
        </button>

        <p v-if="uploadStatus" :class="uploadStatus.type === 'error' ? 'upload-error' : 'upload-success'">
          {{ uploadStatus.message }}
          <a v-if="uploadStatus.url" :href="uploadStatus.url" target="_blank" class="ipfs-link">查看文件</a>
        </p>
      </div>

      <button @click="handleLogout" class="logout-btn">🚪 退出登录</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '@/api/client'; // 假设这是您的 axios 实例

const router = useRouter();
const userInfo = ref(null);
const isLoading = ref(true);
const error = ref('');

// --- 文件上传状态管理 ---
const fileInput = ref(null);
const selectedFile = ref(null);
const isUploading = ref(false);
const uploadStatus = ref(null); // { type: 'success'|'error', message: '', url: '' }


// --- 核心业务逻辑：获取资料 ---
const fetchProfile = async () => {
  const token = localStorage.getItem('jwt_token');
  if (!token) {
    error.value = '未检测到登录凭证。';
    isLoading.value = false;
    // 自动跳转到登录页
    router.push('/login');
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    // 假设 /api/profile 需要 Authorization Header
    const response = await api.get('/api/profile');
    const data = response.data.data || response.data;

    if (data && (data.user_id || data.username)) {
      userInfo.value = data;
      error.value = '';
    } else {
      error.value = '用户资料为空或不完整';
    }
  } catch (err) {
    if (err.response?.status === 401) {
      error.value = '登录凭证已过期或无效，请重新登录。';
    } else {
      error.value = err.response?.data?.message || '获取用户资料失败，请检查网络。';
    }
  } finally {
    isLoading.value = false;
  }
};


// --- 核心业务逻辑：退出登录 ---
const handleLogout = () => {
  localStorage.removeItem('jwt_token');
  localStorage.removeItem('user_info');
  router.push('/login');
};


// --- 文件上传逻辑：处理文件选择 ---
const handleFileChange = (event) => {
  selectedFile.value = event.target.files ? event.target.files[0] : null;
  uploadStatus.value = null; // 清除之前的状态
};

// --- 文件上传逻辑：执行上传 ---
const uploadFile = async () => {
  if (!selectedFile.value) {
    uploadStatus.value = { type: 'error', message: '请先选择一个文件。' };
    return;
  }

  isUploading.value = true;
  uploadStatus.value = null;

  const formData = new FormData();
  // 关键：字段名 'image' 必须与 Go Gin 后端 c.FormFile("image") 匹配
  formData.append('image', selectedFile.value);

  try {
    // 调用 Gin 后端 API /api/upload-image
    const response = await api.post('/api/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        // 假设您的 API 客户端已自动添加了 token
      },
    });

    const imageUrl = response.data.image_url;

    uploadStatus.value = {
      type: 'success',
      message: '文件上传成功并已发布到 IPFS！',
      url: imageUrl
    };
    // 清空文件选择框和文件状态
    selectedFile.value = null;
    fileInput.value.value = '';

  } catch (err) {
    const errorMessage = err.response?.data?.error || err.response?.data?.details || '文件上传失败，请检查后端服务。';
    uploadStatus.value = { type: 'error', message: errorMessage };
  } finally {
    isUploading.value = false;
  }
};

onMounted(fetchProfile);
</script>

<style scoped>
/* 容器和标题基础样式 */
.profile-card {
  background: #fff;
  padding: 40px 50px;
  border-radius: 16px;
  box-shadow: 0 6px 18px rgba(0,0,0,0.08);
  width: 100%;
}
.profile-title {
  text-align: center;
  margin-bottom: 30px;
}
.profile-title .icon {
  font-size: 36px;
}
.profile-title h2 {
  margin: 8px 0 4px;
  font-size: 28px;
  font-weight: 700;
}
.profile-title p {
  color: #888;
  font-size: 14px;
}
.profile-info {
  margin-top: 20px;
}
.info-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  font-size: 16px;
}
.info-item:last-child {
  border-bottom: none;
}

/* 退出登录按钮样式 */
.logout-btn {
  width: 100%;
  margin-top: 30px;
  padding: 12px 0;
  background-color: #4a90e2;
  color: white;
  font-size: 16px;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}
.logout-btn:hover {
  background-color: #3a76c4;
}

/* --- 文件上传区域样式 --- */
.upload-section {
  margin-top: 40px;
  padding: 20px;
  border: 1px dashed #d1d1d1; /* 边框颜色更柔和 */
  background-color: #f9f9f9;
  border-radius: 8px;
  text-align: left;
}
.upload-section h3 {
  font-size: 18px;
  margin-bottom: 15px;
  color: #333;
}
.upload-btn {
  width: 100%;
  margin-top: 15px;
  padding: 10px 0;
  background-color: #2ecc71; /* 绿色系按钮，代表文件操作 */
  color: white;
  font-size: 15px;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}
.upload-btn:hover:not(:disabled) {
  background-color: #27ae60;
}
.upload-btn:disabled {
  background-color: #a0a0a0;
  cursor: not-allowed;
}

/* 上传反馈 */
.upload-success {
  margin-top: 15px;
  color: #27ae60;
  padding: 8px;
  background-color: #e6ffe6;
  border-radius: 4px;
  font-size: 14px;
  word-break: break-all;
}
.upload-error {
  margin-top: 15px;
  color: #e54c3c;
  padding: 8px;
  background-color: #fef0f0;
  border-radius: 4px;
  font-size: 14px;
}
.ipfs-link {
  color: #4a90e2;
  margin-left: 10px;
  text-decoration: underline;
}

/* 加载动画 */
.loading-state, .error-state {
  text-align: center;
  padding: 30px 0;
}
.spinner {
  border: 4px solid rgba(0,0,0,0.1);
  border-top: 4px solid #4a90e2;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  margin: 0 auto 15px;
  animation: spin 1s linear infinite;
}
@keyframes spin { 0%{transform:rotate(0deg);} 100%{transform:rotate(360deg);} }

</style>
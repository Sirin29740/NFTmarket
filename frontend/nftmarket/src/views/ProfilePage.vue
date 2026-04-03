<template>
  <div class="profile-container">
    <div class="user-header">
      <div v-if="isLoadingUser" class="loading-text">加载用户信息中...</div>

      <div v-else class="user-info">
        <div class="avatar">👤</div>
        <div class="details">
          <h2>{{ userInfo?.username || '未登录' }}</h2>
          <p>用户ID: {{ userInfo?.user_id || '---' }}</p>

          <p class="wallet-tag" v-if="isConnected">
            🟢 已连接: {{ address?.slice(0,6) }}...{{ address?.slice(-4) }}
          </p>
          <p class="wallet-tag warn" v-else>🔴 钱包未连接</p>
        </div>

        <div class="header-actions">
          <button v-if="!isConnected" @click="connectWallet" class="btn-connect">连接钱包</button>
          <button v-else @click="disconnect()" class="btn-disconnect">断开</button>
          <button @click="handleLogout" class="btn-logout">退出登录</button>
        </div>
      </div>
    </div>

    <hr class="divider" />

    <div class="mint-section">
      <h3>🎨 铸造 NFT 到区块链</h3>
      <div class="form-group">
        <label>接收者钱包地址 (To) *</label>
        <input
                v-model="nftForm.toAddress"
                placeholder="连接钱包后自动填充"
                :disabled="isMinting"
        />
        <small class="hint">建议连接钱包以确保地址准确</small>
      </div>
      <div class="form-group">
        <label>NFT 名称 *</label>
        <input v-model="nftForm.name" placeholder="作品名称" :disabled="isMinting" />
      </div>

      <div class="form-group">
        <label>描述 *</label>
        <textarea v-model="nftForm.description" placeholder="描述该作品..." rows="2" :disabled="isMinting"></textarea>
      </div>

      <div class="upload-area">
        <label>上传作品文件 *</label>
        <div v-if="nftPreviewUrl" class="image-preview">
          <img :src="nftPreviewUrl" />
          <button @click="removeImage" class="remove-btn" v-if="!isMinting">×</button>
        </div>
        <div v-else class="upload-box" @click="triggerFileSelect" @dragover.prevent="dragover = true" @dragleave.prevent="dragover = false" @drop.prevent="handleDrop" :class="{ 'drag-over': dragover }">
          <p>点击或拖拽图片</p>
          <input type="file" ref="fileInput" @change="handleFileSelect" accept="image/*" style="display: none" />
        </div>
      </div>

      <div v-if="isMinting" class="status-indicator">
        <div class="spinner"></div>
        <p>{{ mintStatusText }}</p>
        <div class="progress-bar-bg">
          <div class="progress-bar-fill" :style="{ width: uploadProgress + '%' }"></div>
        </div>
      </div>

      <button
              @click="startFullMintProcess"
              :disabled="isMinting || !nftForm.name || !nftForm.toAddress || !nftForm.file"
              class="btn-mint"
      >
        {{ isMinting ? '处理中...' : '🎯 立即铸造' }}
      </button>

      <div v-if="resultMsg" :class="['result-box', resultMsg.type]">
        <p>{{ resultMsg.text }}</p>
        <div class="links" v-if="resultMsg.type === 'success'">
          <a :href="resultMsg.metadataUrl" target="_blank">📄 元数据 (IPFS)</a>
          <a :href="resultMsg.txUrl" target="_blank" v-if="resultMsg.txHash">⛓️ 交易哈希</a>
        </div>
      </div>
    </div>

    <hr class="divider" />

    <div class="nft-list-section">
      <div class="list-header">
        <h3>🖼️ 我的收藏仓库</h3>
        <button @click="fetchNFTs" class="btn-refresh" :disabled="isLoadingNFTs">
          {{ isLoadingNFTs ? '同步中...' : '🔄 刷新列表' }}
        </button>
      </div>

      <div v-if="isLoadingNFTs" class="nft-grid">
        <div v-for="i in 3" :key="i" class="skeleton-card"></div>
      </div>

      <div v-else-if="nftList && nftList.length > 0" class="nft-grid">
        <div v-for="nft in nftList" :key="nft.tokenId" class="nft-card">
          <div class="nft-image-wrapper">
            <img :src="parseIpfs(nft.image)" :alt="nft.name" @error="(e) => e.target.src = 'https://via.placeholder.com/200?text=Error'" />
          </div>
          <div class="nft-info">
            <h4 :title="nft.name">{{ nft.name || '未命名作品' }}</h4>
            <p class="nft-desc">{{ nft.description || '暂无描述' }}</p>
            <p class="token-id">Token ID: #{{ nft.tokenId }}</p>
          </div>
        </div>
      </div> <div v-else class="empty-state">
      <p>暂无 NFT 资产</p>
    </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import api from '@/api/client';
  // 引入 wagmi hooks
  import { useAccount, useConnect, useDisconnect } from '@wagmi/vue'
  import { injected } from '@wagmi/vue/connectors'

  const router = useRouter();

  // --- Wagmi 逻辑 ---
  const { address, isConnected } = useAccount();
  const { connect } = useConnect();
  const { disconnect } = useDisconnect();

  const connectWallet = () => {
    connect({ connector: injected() });
  };

  // 监听地址变化：只要钱包切换，自动更新表单地址并刷新列表
  watch(address, (newAddr) => {
    if (newAddr) {
      nftForm.value.toAddress = newAddr;
      fetchNFTs();
    }
  });

  // --- 原有状态保持不变 ---
  const userInfo = ref(null);
  const isLoadingUser = ref(true);
  const nftForm = ref({ toAddress: '', name: '', description: '', file: null });
  const fileInput = ref(null);
  const nftPreviewUrl = ref('');
  const isMinting = ref(false);
  const uploadProgress = ref(0);
  const mintStatusText = ref('');
  const resultMsg = ref(null);
  const dragover = ref(false);
  const nftList = ref([]);
  const isLoadingNFTs = ref(false);

  onMounted(async () => {
    try {
      const res = await api.get('/api/profile');
      userInfo.value = res.data.data || res.data;

      // 逻辑优先度：如果钱包已连接，用钱包地址；否则用 profile 里的地址
      if (isConnected.value && address.value) {
        nftForm.value.toAddress = address.value;
      } else if (userInfo.value.wallet_address) {
        nftForm.value.toAddress = userInfo.value.wallet_address;
      }

      fetchNFTs();
    } catch (err) {
      if (err.response?.status === 401) router.push('/login');
    } finally {
      isLoadingUser.value = false;
    }
  });

  // --- 获取 NFT 列表 (修复了之前的字符串模板错误) ---
  const fetchNFTs = async () => {
    // 优先使用当前连接的钱包地址查询，没有则使用 userInfo 的
    const targetAddr = address.value || userInfo.value?.wallet_address;
    if (!targetAddr) return;

    isLoadingNFTs.value = true;
    try {
      // 修正点：使用了反引号 `` 并在 URL 中正确引用变量
      const res = await api.get(`/api/getnftlist?address=${targetAddr}`);
      nftList.value = res.data || [];
    } catch (err) {
      console.error("获取 NFT 失败:", err);
    } finally {
      isLoadingNFTs.value = false;
    }
  };

  // --- 其他原有方法保持完全一致 ---
  // 1. 定义一个解析函数
  const parseIpfs = (url) => {
    if (!url) return 'https://via.placeholder.com/200?text=No+Data';

    // // 打印一下，确保你能看到真实的值（在 F12 控制台看）
    // console.log("原始 URL:", url);
    //
    // // 只要 URL 里面包含 ipfs.io/ipfs/，管它前面有没有 https
    // if (url.includes('ipfs.io/ipfs/')) {
    //   // 这种方法最稳：直接用 /ipfs/ 分割，取后面那段哈希
    //   const hash = url.split('/ipfs/')[1];
    //   return `https://4everland.io/ipfs/${hash}`;
    // }
    //
    // // 兼容直接以 ipfs:// 开头的情况
    // if (url.startsWith('ipfs://')) {
    //   const hash = url.replace('ipfs://', '');
    //   return `https://4everland.io/ipfs/${hash}`;
    // }

    // 如果已经是别的正常链接，直接返回
    return url;
  };
  // // 2. 这里的 handleImgError 是最后的兜底
  // const handleImgError = (e) => {
  //   // 如果 Gateway 也挂了，显示这个
  //   e.target.src = 'https://via.placeholder.com/200?text=IPFS+Load+Error';
  // };
  const triggerFileSelect = () => fileInput.value.click();
  const handleFileSelect = (e) => prepareFile(e.target.files[0]);
  const handleDrop = (e) => { dragover.value = false; prepareFile(e.dataTransfer.files[0]); };
  const prepareFile = (file) => { if (!file) return; nftForm.value.file = file; nftPreviewUrl.value = URL.createObjectURL(file); };
  const removeImage = () => { nftForm.value.file = null; nftPreviewUrl.value = ''; };

  const startFullMintProcess = async () => {
    if (!nftForm.value.toAddress.startsWith('0x')) { alert("无效地址"); return; }
    isMinting.value = true;
    resultMsg.value = null;
    try {
      mintStatusText.value = "上传中...";
      const formData = new FormData();
      formData.append('name', nftForm.value.name);
      formData.append('description', nftForm.value.description);
      formData.append('image', nftForm.value.file);

      const uploadRes = await api.post('/api/upload', formData, {
        onUploadProgress: (p) => { uploadProgress.value = Math.round((p.loaded * 50) / p.total); },
        headers: { 'Content-Type': 'multipart/form-data' }
      });

      const { token_uri, metadata_url } = uploadRes.data;
      mintStatusText.value = "铸造中...";
      const mintRes = await api.post('/api/mint', { to: nftForm.value.toAddress, token_uri });

      resultMsg.value = { type: 'success', text: '成功', metadataUrl: metadata_url, txHash: mintRes.data.tx_hash, txUrl: `https://sepolia.etherscan.io/tx/${mintRes.data.tx_hash}` };
      setTimeout(fetchNFTs, 2000);
    } catch (err) {
      resultMsg.value = { type: 'error', text: err.message };
    } finally { isMinting.value = false; }
  };

  const handleLogout = () => { localStorage.removeItem('jwt_token'); router.push('/login'); };
</script>

<style scoped>
  /* 在原有样式基础上添加几个钱包按钮样式 */
  .header-actions { display: flex; gap: 10px; align-items: center; margin-left: auto; }
  .btn-connect { background: #f6851b; color: white; border: none; padding: 6px 12px; border-radius: 6px; cursor: pointer; font-weight: bold; }
  .btn-disconnect { background: #eee; color: #666; border: none; padding: 6px 12px; border-radius: 6px; cursor: pointer; }
  .wallet-tag.warn { color: #ff4d4f; background: #fff2f0; }

  /* 以下是你原有的所有样式 ... */
  .profile-container { max-width: 800px; margin: 2rem auto; padding: 2rem; background: #fff; border-radius: 16px; box-shadow: 0 10px 30px rgba(0,0,0,0.08); }
  .user-info { display: flex; align-items: center; gap: 15px; text-align: left; }
  .avatar { font-size: 2rem; background: #f3f0ff; padding: 12px; border-radius: 50%; }
  .wallet-tag { font-family: monospace; font-size: 12px; color: #8a2be2; background: #f3f0ff; padding: 2px 8px; border-radius: 4px; }
  .btn-logout { color: #ff4d4f; border: 1px solid #ff4d4f; background: none; padding: 5px 12px; border-radius: 6px; cursor: pointer; font-size: 13px; }
  .divider { border: 0; border-top: 1px solid #f0f0f0; margin: 2rem 0; }
  .form-group { margin-bottom: 1.2rem; text-align: left; }
  .form-group label { display: block; margin-bottom: 6px; font-weight: 600; color: #333; }
  .form-group input, .form-group textarea { width: 100%; padding: 12px; border: 1px solid #e0e0e0; border-radius: 8px; box-sizing: border-box; }
  .upload-box { border: 2px dashed #d1b9ff; padding: 30px; text-align: center; border-radius: 10px; cursor: pointer; color: #8a2be2; background: #faf8ff; }
  .image-preview img { width: 100%; max-height: 250px; object-fit: contain; }
  .btn-mint { width: 100%; padding: 16px; background: linear-gradient(135deg, #8a2be2, #6a1b9a); color: white; border: none; border-radius: 10px; font-size: 1.1rem; font-weight: bold; cursor: pointer; }
  .nft-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 20px; }
  .nft-card { background: #fff; border: 1px solid #f0f0f0; border-radius: 14px; overflow: hidden; transition: all 0.3s; }
  .nft-image-wrapper { width: 100%; aspect-ratio: 1/1; background: #f8f8f8; display: flex; align-items: center; justify-content: center; }
  .nft-image-wrapper img { width: 100%; height: 100%; object-fit: cover; }
  .nft-info { padding: 12px; }
  .status-indicator { margin: 1.5rem 0; text-align: center; }
  .progress-bar-bg { height: 6px; background: #eee; border-radius: 10px; overflow: hidden; }
  .progress-bar-fill { height: 100%; background: #8a2be2; transition: width 0.3s; }
  .skeleton-card { height: 210px; background: #f6f6f6; border-radius: 14px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.6; } 50% { opacity: 1; } 100% { opacity: 0.6; } }
</style>
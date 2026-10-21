<template>
  <div class="verify-app">
    <!-- 页面头部 -->
    <header class="verify-header">
      <div class="header-container">
        <div class="logo-section">
          <img src="/logo.jpg" alt="logo" class="logo-img" />
          <div class="title-block">
            <h1 class="logo-title">香港國際合唱節</h1>
          <p class="logo-subtitle">HONG KONG INTERNATIONAL<br>CHORAL FESTIVAL</p>
          </div>
        </div>
        <div class="verification-badge">
          <span class="badge-icon">✓</span>
          <span class="badge-text">官方验证</span>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="main-content">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-section">
        <a-spin size="large" tip="正在验证获奖信息...">
          <div class="loading-placeholder"></div>
        </a-spin>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-section">
        <a-result status="error" title="验证失败" :sub-title="error">
          <template #extra>
            <a-button type="primary" @click="retry">重新验证</a-button>
          </template>
        </a-result>
      </div>

      <!-- 验证成功 - 显示获奖信息 -->
      <div v-else-if="awardInfo" class="award-section">
        <!-- 验证成功标识 -->
        <div class="verification-success">
          <div class="success-icon">✓</div>
          <h2 class="success-title">验证成功</h2>
          <p class="success-subtitle">以下获奖信息已通过官方验证</p>
        </div>

        <!-- 获奖信息卡片 -->
        <div class="award-card">
          <div class="award-header">
            <div class="award-type-badge" :class="awardInfo.award_type">
              个人或团队获奖
            </div>
            <div class="award-year">{{ awardInfo.year }}年</div>
          </div>

          <div class="award-content">
            <div class="info-row">
              <span class="label">组别：</span>
              <span class="value">{{ awardInfo.category }}</span>
            </div>
            <div class="award-level">{{ awardInfo.award_level }}</div>

            <!-- 个人获奖信息 -->
            <div v-if="awardInfo.award_type === 'person'" class="person-info">
              <div class="info-row">
                <span class="label">获奖人员或获奖团队名称：</span>
                <span class="value">{{ awardInfo.person_name }}</span>
              </div>
              <!-- <div class="info-row">
                <span class="label">证件号码：</span>
                <span class="value">{{ maskIdNumber(awardInfo.id_number) }}</span>
              </div> -->
            </div>

            <!-- 团体获奖信息 -->
            <div v-else class="group-info">
              <div class="info-row">
                <span class="label">获奖人员或获奖团队名称：</span>
                <span class="value">{{ awardInfo.group_name }}</span>
              </div>
              <div class="info-row">
                <span class="label">指挥：</span>
                <div class="conductor-info">
                  <span class="conductor-name">{{ awardInfo.members?.[0]?.name || '暂无指挥信息' }}</span>
                </div>
              </div>
              <div class="info-row">
                <span class="label">团体成员：</span>
                <div class="members-list">
                  <div v-for="(member, index) in (awardInfo.members || []).slice(1)" :key="index" class="member-item">
                    <span class="member-name">{{ member.name }}</span>
                  </div>
                </div>
              </div>
              
              <!-- 团体证书部分 -->
              <div class="certificate-section">
                <h3 class="certificate-title">
                  <span class="title-icon">🎓</span>
                  获奖证书
                </h3>
                <div class="certificates-container">
                  <!-- 团体证书 -->
                  <div class="certificate-card">
                    <div class="certificate-header">
                      <span class="header-icon">🏆</span>
                      团体证书
                    </div>
                    <div class="certificate-preview" @click="viewCertificate(awardInfo.certificate_image, '团体证书')">
                      <template v-if="awardInfo.certificate_image">
                        <img 
                          :src="processImageUrl(awardInfo.certificate_image)" 
                          alt="团体证书" 
                          class="certificate-thumbnail"
                          @error="handleImageError($event, '团体证书')"
                        />
                        <div class="certificate-overlay">
                          <span class="view-text">
                            <span class="view-icon">🔍</span>
                            点击查看
                          </span>
                        </div>
                      </template>
                      <div v-else class="certificate-placeholder">
                        <span class="placeholder-icon">📄</span>
                        <span class="placeholder-text">暂无证书</span>
                      </div>
                    </div>
                    <div class="certificate-actions">
                      <a-button 
                        type="primary" 
                        :disabled="!awardInfo.certificate_image"
                        :loading="downloadLoading === 'team'"
                        @click="downloadCertificate(awardInfo.certificate_image, '团体证书.jpg', 'team')"
                      >
                        <template v-if="downloadLoading !== 'team'">
                          <download-outlined /> 下载证书
                        </template>
                        <template v-else>
                          下载中...
                        </template>
                      </a-button>
                    </div>
                  </div>

                  <!-- 成员证书 -->
                  <template v-if="awardInfo.members && awardInfo.members.length > 0">
                    <div v-for="(member, index) in awardInfo.members" :key="index" class="certificate-card">
                      <div class="certificate-header">
                        <span class="header-icon">📜</span>
                        {{ member.name }} {{ index === 0 ? '(指挥)' : '' }} 证书
                      </div>
                      <div class="certificate-preview" @click="member.certificate_image && viewCertificate(member.certificate_image, member.name + ' 证书')">
                        <template v-if="member.certificate_image">
                          <img 
                            :src="processImageUrl(member.certificate_image)" 
                            alt="成员证书" 
                            class="certificate-thumbnail"
                            @error="handleImageError($event, member.name + '的证书')"
                          />
                          <div class="certificate-overlay">
                            <span class="view-text">
                              <span class="view-icon">🔍</span>
                              点击查看
                            </span>
                          </div>
                        </template>
                        <div v-else class="certificate-placeholder">
                          <span class="placeholder-icon">📄</span>
                          <span class="placeholder-text">暂无证书</span>
                        </div>
                      </div>
                      <div class="certificate-actions">
                        <a-button 
                          type="primary" 
                          :disabled="!member.certificate_image"
                          :loading="downloadLoading === `member-${index}`"
                          @click="downloadCertificate(member.certificate_image, member.name + '证书.jpg', `member-${index}`)"
                        >
                          <template v-if="downloadLoading !== `member-${index}`">
                            <download-outlined /> 下载证书
                          </template>
                          <template v-else>
                            下载中...
                          </template>
                        </a-button>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 验证信息 -->
        <div class="verification-info">
          <div class="verification-item">
            <span class="verification-label">验证时间：</span>
            <span class="verification-value">{{ formatDate(new Date()) }}</span>
          </div>
          <div class="verification-item">
            <span class="verification-label">验证状态：</span>
            <span class="verification-value verified">已验证 ✓</span>
          </div>
        </div>
      </div>

      <!-- 无数据状态 -->
      <div v-else class="no-data-section">
        <a-result status="404" title="未找到获奖信息" sub-title="请检查验证链接是否正确，或联系相关部门确认。" />
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="verify-footer">
      <div class="footer-container">
        <p class="footer-org">香港國際合唱節組委會</p>
        <p class="footer-info">地址：深圳市南山区后海滨路保利剧院  电话：0755-86371696</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { message } from 'ant-design-vue'
  import { queryPersonAwards, queryGroupAwards } from './api/index.js'
  import { fireConfetti } from './utils/confetti'

  // 响应式数据
  const loading = ref(true)
  const error = ref(null)
  const awardInfo = ref(null)
  const downloadLoading = ref(null) // 用于跟踪证书下载状态

  // 处理图片加载错误
  const handleImageError = (event, imageType) => {
    event.target.style.display = 'none';
    const errorDiv = document.createElement('div');
    errorDiv.className = 'image-error';
    errorDiv.innerHTML = `
      <div class="error-icon">❌</div>
      <div class="error-text">${imageType}加载失败</div>
    `;
    event.target.parentNode.appendChild(errorDiv);
  }

  // 获取URL参数
  const getUrlParams = () => {
    const urlParams = new URLSearchParams(window.location.search)
    const params = {
      type: urlParams.get('type'), // person 或 group
      person_name: urlParams.get('person_name'),
      id_number: urlParams.get('id_number'),
      year: urlParams.get('year'),
      group_name: urlParams.get('group_name'),
      category: urlParams.get('category'),
      award_id: parseInt(urlParams.get('award_id'))
    }

    // 检查必需参数 — 区分个人和团体
    if (!params.type) {
      throw new Error('缺少必要参数: type')
    }

    if (params.type === 'person') {
      const requiredPersonParams = ['person_name', 'id_number', 'year', 'category']
      for (const param of requiredPersonParams) {
        if (!params[param]) {
          throw new Error(`缺少必要参数: ${param}`)
        }
      }
    } else if (params.type === 'group') {
      const requiredGroupParams = ['group_name', 'year', 'category']
      for (const param of requiredGroupParams) {
        if (!params[param]) {
          throw new Error(`缺少必要参数: ${param}`)
        }
      }
    } else {
      throw new Error(`无效的查询类型: ${params.type}`)
    }

    return params
  }

  // 验证获奖信息
  const verifyAward = async () => {
    try {
      loading.value = true
      error.value = null

      const params = getUrlParams()
      // 参数验证已在 getUrlParams 中完成

      let response
      if (params.type === 'person') {
        // 个人获奖查询
        response = await queryPersonAwards({
          person_name: params.person_name,
          id_number: params.id_number,
          year: params.year ? parseInt(params.year) : undefined,
          category: params.category
        })
      } else if (params.type === 'group') {
        // 团体获奖查询
        response = await queryGroupAwards({
          person_name: params.person_name,
          id_number: params.id_number,
          year: params.year ? parseInt(params.year) : undefined,
          group_name: params.group_name,
          category: params.category
        })

      } else {
        throw new Error('无效的查询类型')
      }
      console.log(response);
      if (response && response.awards && response.awards.length > 0) {
        // 优先根据 award_id 精确匹配，否则取第一条
        let targetAward = response.awards[0]
        if (params.award_id) {
          const exactMatch = response.awards.find(award => award.id === params.award_id)
          if (exactMatch) {
            targetAward = exactMatch
          }
        }

        if (!targetAward) {
          awardInfo.value = null
          error.value = '未找到匹配的获奖记录';
        } else {
          // 设置获奖类型
          targetAward.award_type = params.type;
          
          // 处理证书图片路径
          if (targetAward.certificate_image) {
            console.log('团体证书:', targetAward.certificate_image);
          }
          
          // 处理成员证书
          if (targetAward.members) {
            targetAward.members = targetAward.members.map(member => {
              if (member.certificate_image) {
                console.log('成员证书:', member.certificate_image);
              }
              return member;
            });
          }
          
          awardInfo.value = targetAward;
          console.log('获奖信息:', JSON.parse(JSON.stringify(awardInfo.value)));
          
          // 添加礼花效果
          fireConfetti()
        }

      } else {
        awardInfo.value = null
      }
    } catch (err) {
      console.error('验证失败:', err)
      error.value = err.response?.data?.detail || err.message || '验证过程中发生错误'
    } finally {
      loading.value = false
    }
  }

  // 重新验证
  const retry = () => {
    verifyAward()
  }

  // 脱敏身份证号
  const maskIdNumber = (idNumber) => {
    if (!idNumber || idNumber.length < 8) return idNumber
    return idNumber.substring(0, 4) + '****' + idNumber.substring(idNumber.length - 4)
  }

  // 格式化日期
  const formatDate = (date) => {
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  // 处理图片URL
  const processImageUrl = (imageUrl) => {
    if (!imageUrl) return '';
    // 如果是完整URL，直接返回
    if (imageUrl.startsWith('data:image/') || imageUrl.startsWith('http')) {
      return imageUrl;
    }
    // 如果是相对路径，添加基础URL
    if (imageUrl.startsWith('/')) {
      return `${window.location.origin}${imageUrl}`;
    }
    return imageUrl;
  };

  // 查看证书大图
  const viewCertificate = (imageUrl, title) => {
    if (!imageUrl) {
      message.warning('证书图片不存在');
      return;
    }

    const processedUrl = processImageUrl(imageUrl);
    console.log('预览证书:', processedUrl);

    // 创建模态框显示证书大图
    const modal = document.createElement('div')
    modal.className = 'certificate-modal'
    modal.innerHTML = `
    <div class="modal-overlay">
      <div class="modal-content">
        <div class="modal-header">${title || '获奖证书'}</div>
        <img
          src="${processedUrl}"
          alt="${title || '获奖证书'}"
          class="modal-image"
          onerror="this.style.display='none';this.nextElementSibling.style.display='block';"
        />
        <div class="modal-error" style="display:none;text-align:center;padding:20px;color:#ff4d4f;">
          图片加载失败
        </div>
        <button class="modal-close">×</button>
      </div>
    </div>
  `

    document.body.appendChild(modal)

    // 点击关闭
    modal.addEventListener('click', (e) => {
      if (e.target.classList.contains('modal-overlay') || e.target.classList.contains('modal-close')) {
        document.body.removeChild(modal)
      }
    })
  }

  // 下载证书
  const downloadCertificate = async (imageUrl, fileName, loadingKey) => {
    if (downloadLoading.value) {
      message.warning('正在处理其他下载，请稍候');
      return;
    }

    try {
      if (!imageUrl) {
        message.warning('证书图片不存在');
        return;
      }

      downloadLoading.value = loadingKey;
      console.log('开始下载证书:', imageUrl);

      const processedUrl = processImageUrl(imageUrl);
      console.log('处理后的URL:', processedUrl);

      // 用 fetch 获取图片后下载
      const response = await fetch(processedUrl);
      if (!response.ok) {
        throw new Error(`图片获取失败: ${response.status} ${response.statusText}`);
      }

      const contentType = response.headers.get('content-type');
      console.log('图片Content-Type:', contentType);
      
      if (!contentType || !contentType.startsWith('image/')) {
        throw new Error('下载内容不是有效的图片格式');
      }

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      message.success('证书下载成功');

    } catch (error) {
      console.error('证书下载失败:', error);
      message.error(error.message || '证书下载失败，请稍后重试');
    } finally {
      downloadLoading.value = null;
    }
  }

  // 组件挂载时开始验证
  onMounted(() => {
    verifyAward()
  })
</script>

<style scoped>
  .verify-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #FFFAF0 0%, #FFF8DC 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  }

  /* 头部样式 */
  .verify-header {
    background: linear-gradient(135deg, #FFE4B5 0%, #FFF8DC 100%);
    backdrop-filter: blur(10px);
    box-shadow: 0 4px 15px rgba(222, 184, 135, 0.25);
    padding: 20px 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .header-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .logo-section {
    color: white;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }

  .logo-img {
    height: 56px;
    margin: 0 auto 0 auto;
    display: block;
  }

  .title-block {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
  }

  .logo-title {
    font-size: 28px;
    font-weight: 700;
    margin: 0;
    background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    letter-spacing: 1px;
    text-align: center;
  }

  .logo-subtitle {
    font-size: 16px;
    margin: 2px 0 0 0;
    font-weight: 300;
    color: #C14B4B;
    text-align: center;
    letter-spacing: 0.5px;
  }

  .verification-badge {
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.95));
    box-shadow: 0 4px 15px rgba(222, 184, 135, 0.3);
    border-radius: 25px;
    padding: 8px 20px;
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
    border: 2px solid #F9BE50;
    backdrop-filter: blur(5px);
    transition: all 0.3s ease;
  }

  .badge-icon {
    background: linear-gradient(135deg, #52c41a, #95de64);
    color: white;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: bold;
    box-shadow: 0 2px 6px rgba(82, 196, 26, 0.3);
  }

  /* 主要内容区域 */
  .main-content {
    flex: 1;
    padding: 40px 20px;
    max-width: 800px;
    margin: 0 auto;
    width: 100%;
  }

  /* 加载状态 */
  .loading-section {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
  }

  .loading-placeholder {
    width: 200px;
    height: 200px;
  }

  /* 验证成功标识 */
  .verification-success {
    text-align: center;
    margin-bottom: 30px;
    animation: fadeInDown 0.8s ease-out;
  }

  .success-icon {
    width: 80px;
    height: 80px;
    background: #52c41a;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 40px;
    color: white;
    margin: 0 auto 20px;
    font-weight: bold;
    box-shadow: 0 8px 25px rgba(82, 196, 26, 0.3);
    animation: bounceIn 1s ease-out 0.3s both;
  }

  .success-title {
    font-size: 32px;
    font-weight: 700;
    color: white;
    margin: 0 0 10px 0;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .success-subtitle {
    font-size: 16px;
    color: rgba(255, 255, 255, 0.9);
    margin: 0;
  }

  /* 获奖信息卡片 */
  .award-card {
    background: linear-gradient(135deg, #fff 0%, #f8f9fa 100%);
    border-radius: 16px;
    box-shadow: 0 20px 40px rgba(255, 155, 155, 0.15);
    overflow: hidden;
    margin-bottom: 30px;
    animation: fadeInUp 0.8s ease-out 0.2s both;
    border: 1px solid rgba(255, 155, 155, 0.2);
  }

  .award-header {
    background: linear-gradient(135deg, #FF9B9B 0%, #FFD699 100%);
    padding: 20px 30px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .award-type-badge {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    box-shadow: 0 2px 8px rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.3);
    backdrop-filter: blur(5px);
  }

  .award-type-badge.person {
    background: rgba(255, 155, 155, 0.3);
    border: 1px solid rgba(255, 155, 155, 0.5);
  }

  .award-type-badge.group {
    background: rgba(255, 214, 153, 0.3);
    border: 1px solid rgba(255, 214, 153, 0.5);
  }

  .award-year {
    color: white;
    font-size: 18px;
    font-weight: 600;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  }

  .award-content {
    padding: 30px;
  }

  .award-title {
    font-size: 24px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 10px 0;
    line-height: 1.3;
  }

  .award-level {
    font-size: 18px;
    color: #FF6B6B;
    font-weight: 700;
    margin-bottom: 25px;
    padding: 10px 20px;
    background: linear-gradient(135deg, rgba(255, 155, 155, 0.15) 0%, rgba(255, 214, 153, 0.15) 100%);
    border-radius: 20px;
    display: inline-block;
    border: 1px solid rgba(255, 155, 155, 0.3);
    box-shadow: 0 2px 8px rgba(255, 155, 155, 0.1);
  }

  /* 信息行样式 */
  .info-row {
    display: flex;
    margin-bottom: 15px;
    align-items: flex-start;
  }

  .label {
    font-weight: 600;
    color: #FF9B9B;
    min-width: 100px;
    margin-right: 15px;
    font-size: 15px;
  }

  .value {
    color: #333;
    flex: 1;
    font-weight: 500;
    font-size: 15px;
  }

  .description-text {
    line-height: 1.6;
  }

  /* 指挥信息 */
  .conductor-info {
    flex: 1;
    background: linear-gradient(135deg, #FF9B9B 0%, #FFD699 100%);
    border-radius: 12px;
    padding: 16px 20px;
    margin-bottom: 20px;
    box-shadow: 0 4px 12px rgba(255, 155, 155, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  .conductor-name {
  font-size: 18px;
  font-weight: 600;
  color: white;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
  }

  /* 成员列表 */
  .members-list {
    flex: 1;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
    margin-top: 8px;
  }

  .member-item {
    background: linear-gradient(135deg, #FF9B9B 0%, #FFD699 100%);
    border-radius: 8px;
    padding: 12px 16px;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
    transition: all 0.3s ease;
    box-shadow: 0 2px 8px rgba(255, 155, 155, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  .member-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 155, 155, 0.3);
    background: linear-gradient(135deg, #FF8B8B 0%, #FFB347 100%);
  }

  .member-name {
  font-weight: 600;
  color: white;
    font-size: 14px;
    line-height: 1.4;
  }

  .member-id {
    color: #666;
    font-size: 14px;
  }

  /* 响应式布局调整 */
  @media (max-width: 768px) {
    .members-list {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 8px;
    }

    .member-item {
      padding: 10px 12px;
    }

    .conductor-info {
      padding: 12px 16px;
      margin-bottom: 16px;
    }

    .conductor-name {
      font-size: 16px;
    }
  }

  /* 证书样式 */
  .certificate-section {
    margin-top: 30px;
    padding: 25px;
    background: #f8f9fa;
    border-radius: 16px;
    box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.05);
  }

  .certificate-title {
    font-size: 20px;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 2px solid rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .title-icon {
    font-size: 24px;
  }

  .certificate-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    background: #f0f2f5;
    border-radius: 8px;
    padding: 20px;
    gap: 12px;
  }

  .placeholder-icon {
    font-size: 32px;
    opacity: 0.5;
  }

  .placeholder-text {
    color: #999;
    font-size: 14px;
  }

  .certificates-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
    margin-top: 20px;
  }

  .certificate-card {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }

  .certificate-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  .certificate-header {
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.15);
    color: white;
    font-weight: 600;
    font-size: 16px;
    text-align: center;
  }

  .certificate-preview {
    position: relative;
    cursor: pointer;
    aspect-ratio: 4/3;
    overflow: hidden;
  }

  .certificate-thumbnail {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .certificate-preview:hover .certificate-thumbnail {
    transform: scale(1.05);
  }

  .certificate-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .certificate-preview:hover .certificate-overlay {
    opacity: 1;
  }

  .view-text {
    color: white;
    font-size: 16px;
    font-weight: 600;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
    padding: 8px 16px;
    border: 2px solid white;
    border-radius: 20px;
    background: rgba(0, 0, 0, 0.3);
  }

  .certificate-actions {
    padding: 12px;
    text-align: center;
    background: #f8f9fa;
    border-top: 1px solid #eee;
  }

  /* 证书模态框样式 */
  :global(.certificate-modal .modal-header) {
  background: rgba(0, 0, 0, 0.15);
    color: white;
    padding: 16px;
    font-size: 18px;
    font-weight: 600;
    text-align: center;
    border-radius: 8px 8px 0 0;
    margin-bottom: 20px;
  }

  :global(.certificate-modal .modal-content) {
    background: white;
    border-radius: 8px;
    overflow: hidden;
    max-width: 90vw;
    max-height: 90vh;
  }

  :global(.certificate-modal .modal-image) {
    max-width: 100%;
    max-height: calc(90vh - 100px);
    display: block;
    margin: 0 auto;
    padding: 20px;
  }

  /* 图片错误样式 */
  .image-error {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: #f5f5f5;
    padding: 20px;
    text-align: center;
  }

  .error-icon {
    font-size: 24px;
    margin-bottom: 8px;
  }

  .error-text {
    color: #ff4d4f;
    font-size: 14px;
  }

  /* 证书卡片图标 */
  .header-icon {
    margin-right: 8px;
    font-size: 18px;
  }

  .view-icon {
    margin-right: 6px;
    font-size: 16px;
  }

  /* 无证书状态 */
  .no-certificates {
    grid-column: 1 / -1;
    padding: 40px;
    background: white;
    border-radius: 12px;
    text-align: center;
  }

  /* 响应式调整 */
  @media (max-width: 768px) {
    .certificates-container {
      grid-template-columns: 1fr;
      gap: 16px;
    }

    .certificate-section {
      padding: 16px;
      margin-top: 20px;
    }

    .certificate-title {
      font-size: 18px;
      margin-bottom: 16px;
      padding-bottom: 12px;
    }

    .certificate-header {
      padding: 10px;
      font-size: 14px;
    }

    .certificate-actions {
      padding: 10px;
    }

    :global(.certificate-modal .modal-header) {
      padding: 12px;
      font-size: 16px;
    }
  }

  /* 验证信息 */
  .verification-info {
    background: linear-gradient(135deg, #FFE4B5 0%, #FFF8DC 100%);
    border-radius: 12px;
    padding: 20px;
    backdrop-filter: blur(10px);
    box-shadow: 0 4px 15px rgba(222, 184, 135, 0.25);
    animation: fadeIn 0.8s ease-out 0.4s both;
    border: 1px solid rgba(255, 255, 255, 0.3);
  }

  .verification-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }

  .verification-item:last-child {
    margin-bottom: 0;
  }

  .verification-label {
    font-weight: 500;
    display: inline-block;
  }

  .verification-value {
    font-weight: 600;
    display: inline-block;
  }

  .verification-value.verified {
    color: #52c41a;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  }

  /* 页脚 */
  .verify-footer {
    background: linear-gradient(135deg, #FFE4B5 0%, #FFF8DC 100%);
    text-align: center;
    padding: 20px 0 10px 0;
    backdrop-filter: blur(10px);
    box-shadow: 0 -4px 15px rgba(222, 184, 135, 0.25);
    border-top: 1px solid rgba(255, 255, 255, 0.2);
  }

  .footer-container {
    max-width: 800px;
    margin: 0 auto;
  }

  .footer-org {
    font-size: 1.1rem;
    font-weight: 500;
    margin: 0 0 4px 0;
    letter-spacing: 1px;
    background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    display: inline-block;
  }

  .footer-info {
    font-size: 0.95rem;
    margin: 0;
    letter-spacing: 0.5px;
    background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    display: inline-block;
  }

  /* 证书模态框 */
  :global(.certificate-modal) {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 9999;
  }

  :global(.modal-overlay) {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
  }

  :global(.modal-content) {
    position: relative;
    max-width: 90vw;
    max-height: 90vh;
  }

  :global(.modal-image) {
    max-width: 100%;
    max-height: 100%;
    border-radius: 8px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  }

  :global(.modal-close) {
    position: absolute;
    top: -15px;
    right: -15px;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: white;
    border: none;
    font-size: 24px;
    font-weight: bold;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    transition: transform 0.2s ease;
  }

  :global(.modal-close:hover) {
    transform: scale(1.1);
  }

  /* 动画 */
  @keyframes fadeInDown {
    from {
      opacity: 0;
      transform: translateY(-30px);
    }

    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }

    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes bounceIn {
    0% {
      opacity: 0;
      transform: scale(0.3);
    }

    50% {
      opacity: 1;
      transform: scale(1.05);
    }

    70% {
      transform: scale(0.9);
    }

    100% {
      opacity: 1;
      transform: scale(1);
    }
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }

    to {
      opacity: 1;
    }
  }

  /* 响应式设计 */
  @media (max-width: 768px) {
    .header-container {
      flex-direction: column;
      gap: 15px;
      text-align: center;
      align-items: center;
    }

    .main-content {
      padding: 20px 15px;
    }

    .award-content {
      padding: 20px;
    }

    .award-title {
      font-size: 20px;
    }

    .success-title {
      font-size: 24px;
    }

    .info-row {
      flex-direction: column;
      align-items: flex-start;
    }

    .label {
      margin-bottom: 5px;
      margin-right: 0;
    }

    .member-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 5px;
    }
    .logo-img {
      height: 40px;
    }
    .logo-title {
      font-size: 1.2rem;
    }
    .logo-subtitle {
      font-size: 0.9rem;
    }
  }
  @media (max-width: 768px) {
    .footer-org {
      font-size: 1rem;
    }
    .footer-info {
      font-size: 0.85rem;
    }
  }

  .badge-text {
    font-size: 15px;
    background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    letter-spacing: 1px;
  }

  .verification-badge:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(222, 184, 135, 0.4);
    border-color: #E6904F;
  }
</style>

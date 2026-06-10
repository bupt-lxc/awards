<template>
  <div class="group-query-container">
    <!-- 页面头部 -->
    <div class="header-section">
      <div class="header-content">
        <h1 class="main-title">
          团队获奖记录查询
        </h1>
        <p class="subtitle">通过团队名称和参赛编码查询团队获奖信息，指挥姓名可选</p>
      </div>
    </div>

    <!-- 查询表单区域 -->
    <div class="query-section">
      <a-card class="query-card" :bordered="false">
        <template #title>
          <div class="card-title">
            <span class="title-icon">🔍</span>
            查询条件
          </div>
        </template>
        
        <a-form :model="queryForm" layout="vertical" @finish="onQuery">
          <a-row :gutter="[24, 16]">
            <a-col :xs="24" :sm="24" :md="12">
              <a-form-item label="团队名称" name="group_name" :rules="[{ required: true, message: '请选择团队名称' }]">
                <a-auto-complete
                  v-model:value="queryForm.group_name"
                  :options="filteredTeamOptions"
                  placeholder="请输入或选择团队名称"
                  size="large"
                  allow-clear
                  @search="handleTeamSearch"
                  @select="handleTeamSelect"
                  :default-active-first-option="false"
                  style="width: 100%"
                  :dropdown-match-select-width="false"
                  :dropdown-style="{ minWidth: '500px', maxHeight: '300px', overflowY: 'auto' }"
                >
                  <template #option="{ value: team }">
                    <div class="team-option">
                      {{ team }}
                    </div>
                  </template>
                </a-auto-complete>
                <div class="search-tips">
                  <small style="color: #666;">
                    💡 搜索提示：可直接输入城市名（如"广州"、"深圳"）、学校名（如"实验小学"）或关键词（如"童声"、"合唱团"）
                  </small>
                </div>
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="参赛编码" name="id_number" :rules="[{ required: true, message: '请输入参赛编码' }]">
                <a-input 
                  v-model:value="queryForm.id_number" 
                  placeholder="输入参赛编码"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="指挥姓名" name="person_name">
                <a-input 
                  v-model:value="queryForm.person_name" 
                  placeholder="输入指挥姓名（可选）"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="获奖年份">
                <a-select 
                  v-model:value="queryForm.year" 
                  placeholder="请选择年份（可选）"
                  size="large"
                  show-search
                  allow-clear
                >
                  <a-select-option v-for="year in availableYears" :key="year" :value="year">
                    {{ year }}年
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="组别">
                <a-select 
                  v-model:value="queryForm.category" 
                  placeholder="请选择组别（可选）"
                  size="large"
                  show-search
                  allow-clear
                >
                  <a-select-option value="表演合唱">表演合唱</a-select-option>
                  <a-select-option value="儿童合唱 I">儿童合唱 I</a-select-option>
                  <a-select-option value="儿童合唱 II">儿童合唱 II</a-select-option>
                  <a-select-option value="儿童合唱 III">儿童合唱 III</a-select-option>
                  <a-select-option value="流行合唱">流行合唱</a-select-option>
                  <a-select-option value="民谣合唱">民谣合唱</a-select-option>
                  <a-select-option value="青年教师 I">青年教师 I</a-select-option>
                  <a-select-option value="青年教师 II">青年教师 II</a-select-option>
                  <a-select-option value="青年室内 I">青年室内 I</a-select-option>
                  <a-select-option value="青年室内 II">青年室内 II</a-select-option>
                  <a-select-option value="少儿合唱 I">少儿合唱 I</a-select-option>
                  <a-select-option value="少儿合唱 II">少儿合唱 II</a-select-option>
                  <a-select-option value="少儿合唱 III">少儿合唱 III</a-select-option>
                  <a-select-option value="幼儿合唱">幼儿合唱</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12" :md="6">
              <CaptchaInput
                ref="captchaRef"
                v-model:value="queryForm.captcha"
              />
            </a-col>
            <a-col :span="24">
              <div class="action-buttons">
                <a-button 
                  type="primary" 
                  html-type="submit" 
                  :loading="loading"
                  size="large"
                  class="query-btn"
                >
                  <template #icon>
                    <search-outlined />
                  </template>
                  <span class="btn-text">开始查询</span>
                  <div class="btn-ripple"></div>
                </a-button>
              </div>
            </a-col>
            <a-col :span="24">
              <div class="action-buttons">
                <a-button 
                  size="large" 
                  @click="resetQuery"
                  class="reset-btn"
                >
                  <template #icon>
                    <reload-outlined />
                  </template>
                  重置条件
                </a-button>
              </div>
            </a-col>
          </a-row>
        </a-form>
      </a-card>
    </div>

    <!-- 查询结果区域 -->
    <div class="results-section" v-if="queryResults.length > 0 || hasQueried || errorMessage">
      <a-card class="results-card" :bordered="false">
        <template #title>
          <div class="card-title">
            <span class="title-icon">📋</span>
            查询结果
            <a-tag color="blue" v-if="queryResults.length > 0">
              共找到 {{ queryResults.length }} 条记录
            </a-tag>
            <a-tag color="red" v-if="errorMessage">
              查询出错
            </a-tag>
          </div>
        </template>
        
        <!-- 错误信息展示 -->
        <div v-if="errorMessage" class="error-section">
          <a-alert
            message="查询失败"
            :description="errorMessage"
            type="error"
            show-icon
            closable
            @close="clearError"
            class="error-alert"
          >
            <template #icon>
              <exclamation-circle-outlined />
            </template>
            <template #action>
              <a-button size="small" type="primary" @click="retryQuery">
                重试查询
              </a-button>
            </template>
          </a-alert>
        </div>
        
        <div v-else-if="queryResults.length === 0 && hasQueried" class="no-results">
          <a-empty description="未找到符合条件的获奖记录" />
        </div>
        
        <div v-else-if="queryResults.length > 0" class="results-grid">
          <div 
            v-for="(award, index) in queryResults" 
            :key="index" 
            class="award-card"
          >
            <div class="award-header">
              <h3 class="award-title">{{ award.group_name }}</h3>
              <a-tag :color="getAwardLevelColor(award.award_level)" class="award-level">
                {{ award.award_level }}
              </a-tag>
            </div>
            <div class="award-info-list">
              <div class="info-row"><span class="info-label">团体名称：</span><span class="info-value">{{ award.group_name }}</span></div>
              <div class="info-row"><span class="info-label">组别：</span><span class="info-value">{{ award.category || '-' }}</span></div>
              <div class="info-row"><span class="info-label">获奖年份：</span><span class="info-value">{{ award.year }}年</span></div>
              <div class="info-row"><span class="info-label">奖项等级：</span><span class="info-value"><a-tag :color="getAwardLevelColor(award.award_level)">{{ award.award_level }}</a-tag></span></div>
              <div class="info-row" v-if="award.description"><span class="info-label">获奖说明：</span><span class="info-value">{{ award.description }}</span></div>
            </div>
            <div class="award-actions">
              <div class="action-btns-center">
                <a-button type="default" size="large" class="action-btn view-btn" @click="viewCertificate(award)" :disabled="!award.certificate_image">
                  <template #icon><EyeOutlined /></template>
                  查看证书
                </a-button>
                <a-button type="primary" size="large" class="action-btn download-btn" @click="downloadCertificate(award)" :disabled="!award.certificate_image">
                  <template #icon><DownloadOutlined /></template>
                  下载证书
                </a-button>
              </div>
            </div>

            <!-- 团体成员信息 -->
            <div class="members-section" v-if="award.members && award.members.length > 0">
              <div class="members-title">团体成员：</div>
              <div class="members-list">
                <div 
                  v-for="(member, memberIndex) in award.members.filter(m => m.name === queryForm.person_name && m.id_number === queryForm.id_number)" 
                  :key="memberIndex" 
                  class="member-item"
                  :class="{ 'highlighted': isTargetMember(member) }"
                >
                  <div class="member-info">
                    <span class="member-name">{{ member.name }}</span>
                    <span class="member-id">{{ member.id_number }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </a-card>
    </div>

    <!-- 证书查看模态框 -->
    <a-modal
      v-model:visible="certificateModalVisible"
      :footer="null"
      width="min(90vw, 500px)"
      :centered="true"
      title="获奖证书"
      @cancel="certificateModalVisible=false"
      wrapClassName="certificate-modal"
    >
      <template #closeIcon>
        <a-button type="text" shape="circle" size="large" class="close-modal-btn" @click="certificateModalVisible=false">
          <template #icon><CloseOutlined /></template>
        </a-button>
      </template>
      <img v-if="certificateModalImg" :src="certificateModalImg" alt="获奖证书" style="width:100%;border-radius:10px;box-shadow:0 2px 16px #aaa;" />
    </a-modal>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, ReloadOutlined, EyeOutlined, DownloadOutlined, ExclamationCircleOutlined, CloseOutlined } from '@ant-design/icons-vue'
import { queryGroupAwards } from '../api/index.js'
import QRCode from 'qrcode'
import CaptchaInput from './CaptchaInput.vue'
import { fireConfetti } from '../utils/confetti'

export default {
  name: 'GroupQueryView',
  components: {
    SearchOutlined,
    ReloadOutlined,
    EyeOutlined,
    DownloadOutlined,
    CaptchaInput,
    ExclamationCircleOutlined,
    CloseOutlined
  },
  setup() {
    const loading = ref(false)
    const hasQueried = ref(false)
    const queryResults = ref([])
    const captchaRef = ref(null);
    const errorMessage = ref('');
    
    const queryForm = reactive({
      group_name: '',
      id_number: '',
      person_name: '',
      year: null,
      category: '',
      captcha: ''
    })
    
    // 生成可选年份
    const availableYears = computed(() => {
      const currentYear = new Date().getFullYear()
      const years = []
      for (let i = currentYear; i >= currentYear - 20; i--) {
        years.push(i)
      }
      return years
    })
    
    // 团队名称选项列表
    const teamOptions = [
      'Balatime choir №3(Kazakhstan, Astana)',
      '广州市花都区青少年文化宫花之蕾合唱B团',
      '武汉市光谷第一小学布谷童声合唱团',
      '福州糖果少儿合唱团',
      '江苏省东台市实验小学小红花合唱团',
      '杭州寻心少儿合唱团C团',
      '新西兰The Cloud童声合唱团',
      '合肥市屯溪路小学启明星合唱团',
      '珠海大剧院童声合唱团(扬帆少年合唱团)',
      '茂名市茂南第一小学合唱团',
      '广州市荔湾区乐贤坊小学合唱团',
      '北京佳音童声合唱团',
      '苏州市吴韵少年艺术团・合唱团二团',
      '华夏四季合唱团&东周小学合唱团',
      '广州市花都区秀全街九潭初级中学合唱团',
      '杭州市天杭实验学校"天杭之声"合唱团',
      'State chamber choir of The Republic of Tatarstan (Russia)',
      '福鼎市硖门畲乡"金凤凰"童声合唱团',
      '深圳繁星女声合唱团',
      '意大利Rom Ensemble合唱团',
      '随园歌手',
      'Chamber choir "Bel Canto"(Russia, Omsk)',
      '上海市学生艺术团普陀区青少年教育活动中心合唱团',
      '广州市花都区青少年文化宫花之蕾合唱A团',
      '太阳照照童声合唱团',
      '松狮小侬侬童声合唱团',
      '广州市黄广少年合唱团',
      '内蒙古小青马合唱团',
      '济南大野少年合唱团',
      '南山文体中心剧院小橙堡童声合唱团绽放团',
      '芜湖市少年宫合唱团',
      '深圳市福田区外国语学校(景秀)合唱团',
      '南山文体中心剧院小橙堡童声合唱团萌芽团',
      '深圳保利天堂鸟童声合唱乐芽团',
      '镇江星光童声合唱团-杜鹃团',
      '成都城市音乐厅青少年合唱团C2团',
      '成都城市音乐厅青少年合唱团C1团'
    ]
    
    // 过滤后的团队选项
    const filteredTeamOptions = ref([])

    // 团队名称搜索处理函数
    const handleTeamSearch = (value) => {
      if (!value) {
        filteredTeamOptions.value = teamOptions.map(team => ({ value: team }))
        return
      }

      const inputLower = value.toLowerCase()
      const filtered = teamOptions.filter(team => {
        const teamLower = team.toLowerCase()
        
        // 1. 直接包含输入内容
        if (teamLower.includes(inputLower)) return true
        
        // 2. 支持拼音首字母搜索
        const pinyinMap = {
          'gz': '广州', 'wh': '武汉', 'fz': '福州', 'js': '江苏', 'hz': '杭州', 'hf': '合肥',
          'zh': '珠海', 'mm': '茂名', 'bj': '北京', 'sz': '苏州', 'hx': '华夏', 'fd': '福鼎',
          'sh': '上海', 'ty': '太阳', 'ss': '松狮', 'hg': '黄广', 'nm': '内蒙古', 'jn': '济南',
          'ns': '南山', 'wh': '芜湖', 'sf': '深圳', 'zj': '镇江', 'cd': '成都'
        }
        
        for (const [pinyin, chinese] of Object.entries(pinyinMap)) {
          if (inputLower.includes(pinyin) && teamLower.includes(chinese)) return true
        }
        
        // 3. 支持关键词搜索
        const keywords = ['合唱团', '童声', '少年', '青年', '女声', '混声', '童声合唱团', '少年合唱团']
        for (const keyword of keywords) {
          if (inputLower.includes(keyword.toLowerCase()) && teamLower.includes(keyword)) return true
        }
        
        return false
      })

      filteredTeamOptions.value = filtered.map(team => ({
        value: team,
      }))
    }

    // 初始化显示所有选项
    onMounted(() => {
      filteredTeamOptions.value = teamOptions.map(team => ({ value: team }))
    })

    // 选择团队名称
    const handleTeamSelect = (value) => {
      queryForm.group_name = value
    }
    
    // 获取奖项等级颜色
    const getAwardLevelColor = (level) => {
      const colorMap = {
        '金奖': 'gold',
        '银奖': 'silver', 
        '铜奖': '#cd7f32',
        '一等奖': 'red',
        '二等奖': 'orange',
        '三等奖': 'green',
        '优秀奖': 'blue',
        '特等奖': 'purple'
      }
      return colorMap[level] || 'default'
    }
    
    // 判断是否为目标成员（高亮显示）
    const isTargetMember = (member) => {
      return member.name === queryForm.person_name && member.id_number === queryForm.id_number
    }
    
    // 查询函数
    const onQuery = async () => {
      try {
        // 验证验证码
        const captchaText = captchaRef.value.getCaptchaText();
        if (queryForm.captcha !== captchaText) {
          message.error('验证码错误');
          captchaRef.value.refreshCaptcha();
          return;
        }

        loading.value = true
        hasQueried.value = true
        
        const queryData = {
          group_name: queryForm.group_name,
          id_number: queryForm.id_number
        }
        
        // 只有当有值时才添加可选参数
        if (queryForm.person_name) {
          queryData.person_name = queryForm.person_name
        }
        if (queryForm.year) {
          queryData.year = queryForm.year
        }
        if (queryForm.category) {
          queryData.category = queryForm.category
        }
        
        const response = await queryGroupAwards(queryData)
        
        queryResults.value = response.awards || []
        
        if (queryResults.value.length === 0) {
          message.info('未找到符合条件的获奖记录')
        } else {
          message.success(`查询成功，找到 ${queryResults.value.length} 条记录`)
          // 添加礼花效果
          fireConfetti()
        }
      } catch (error) {
        console.error('查询失败:', error)
        
        // 提供更友好的错误信息
        let userFriendlyError = '查询失败，请稍后重试'
        
        if (error.response) {
          const status = error.response.status
          const errorData = error.response.data
          
          switch (status) {
            case 400:
              if (errorData?.error) {
                userFriendlyError = `请求参数错误：${errorData.error}`
              } else {
                userFriendlyError = '请求参数有误，请检查输入信息'
              }
              break
            case 401:
              userFriendlyError = '认证失败，请重新登录'
              break
            case 403:
              userFriendlyError = '权限不足，无法访问此功能'
              break
            case 404:
              userFriendlyError = '查询服务暂时不可用，请稍后重试'
              break
            case 500:
              userFriendlyError = '服务器内部错误，请稍后重试'
              break
            default:
              if (errorData?.error) {
                userFriendlyError = errorData.error
              } else {
                userFriendlyError = `查询失败（错误代码：${status}），请稍后重试`
              }
          }
        } else if (error.request) {
          userFriendlyError = '网络连接失败，请检查网络设置后重试'
        } else if (error.message) {
          if (error.message.includes('timeout')) {
            userFriendlyError = '查询超时，请稍后重试'
          } else {
            userFriendlyError = `查询出错：${error.message}`
          }
        }
        
        errorMessage.value = userFriendlyError
        queryResults.value = []
        
        // 显示错误提示
        message.error(userFriendlyError)
      } finally {
        loading.value = false
        // 刷新验证码
        captchaRef.value.refreshCaptcha();
        queryForm.captcha = '';
      }
    }
    
    // 重置查询
    const resetQuery = () => {
      Object.assign(queryForm, {
        group_name: '',
        id_number: '',
        person_name: '',
        year: null,
        category: '',
        captcha: ''
      })
      queryResults.value = []
      hasQueried.value = false
      errorMessage.value = ''
      // 刷新验证码
      captchaRef.value.refreshCaptcha();
    }

    // 重试查询
    const retryQuery = () => {
      onQuery();
    };

    // 清除错误信息
    const clearError = () => {
      errorMessage.value = '';
    };
    
    // 生成专属验证链接
    const getVerifyUrl = (award) => {
      const baseUrl = window.location.origin;
      const params = new URLSearchParams({
        type: 'group',
        year: award.year,
        group_name: award.group_name,
        category: award.category || '',
        award_id: award.id
      });
      return `${baseUrl}/verify.html?${params.toString()}`;
    };

    // 合成证书图片（加二维码和防伪字样）
    const composeCertificateImage = async (award) => {
      return new Promise(async (resolve, reject) => {
        try {
          // 确保证书图片数据是base64格式
          if (!award.certificate_image) {
            throw new Error('证书图片数据为空');
          }

          const img = new window.Image();
          img.crossOrigin = 'anonymous';
          
          // 如果已经是完整的base64格式，直接使用
          if (award.certificate_image.startsWith('data:image')) {
            img.src = award.certificate_image;
          } else {
            // 否则添加前缀
            img.src = `data:image/jpeg;base64,${award.certificate_image}`;
          }

          img.onload = async () => {
            const canvas = document.createElement('canvas');
            canvas.width = img.width;
            canvas.height = img.height;
            const ctx = canvas.getContext('2d');
            ctx.drawImage(img, 0, 0);

            // 生成防伪二维码
            const qrSize = Math.floor(Math.min(img.width, img.height) * 0.13);
            const verifyUrl = getVerifyUrl(award);
            console.log('生成防伪链接:', verifyUrl); // 调试用
            const qrCodeDataUrl = await QRCode.toDataURL(verifyUrl, { 
              width: qrSize, 
              margin: 0,
              errorCorrectionLevel: 'H' // 提高容错率
            });

            const qrImg = new window.Image();
            qrImg.src = qrCodeDataUrl;
            qrImg.onload = () => {
              const padding = 32;
              const qrX = img.width - qrSize - padding;
              const qrY = padding;
              
              // 添加白色背景提高二维码可见度
              ctx.fillStyle = 'white';
              ctx.fillRect(qrX - 5, qrY - 5, qrSize + 10, qrSize + 10);
              
              // 绘制二维码
              ctx.drawImage(qrImg, qrX, qrY, qrSize, qrSize);
              
              // 添加防伪文字
              ctx.font = `${Math.floor(qrSize/5)}px Arial`;
              ctx.fillStyle = '#333';
              ctx.textAlign = 'center';
              ctx.shadowColor = 'white';
              ctx.shadowBlur = 2;
              ctx.fillText('防伪标识', qrX + qrSize/2, qrY + qrSize + Math.floor(qrSize/3));
              
              resolve(canvas.toDataURL('image/png'));
            };
            qrImg.onerror = (e) => {
              console.error('二维码加载失败:', e);
              reject(new Error('二维码生成失败'));
            };
          };
          img.onerror = (e) => {
            console.error('证书图片加载失败:', e);
            reject(new Error('证书图片加载失败'));
          };
        } catch (e) {
          console.error('证书处理失败:', e);
          reject(e);
        }
      });
    };

    // 查看证书大图（合成后）
    const certificateModalVisible = ref(false);
    const certificateModalImg = ref('');
    const viewCertificate = async (award) => {
      if (!award.certificate_image) {
        message.error('该获奖记录暂无证书');
        return;
      }

      certificateModalImg.value = '';
      certificateModalVisible.value = true;
      try {
        const loadingMessage = message.loading('正在加载证书...', 0);
        console.log('开始处理证书图片:', award.certificate_image.substring(0, 100)); // 调试用
        certificateModalImg.value = await composeCertificateImage(award);
        loadingMessage();
        console.log('证书加载成功');
      } catch (e) {
        console.error('证书加载失败:', e);
        message.error(e.message || '证书加载失败，请稍后重试');
        certificateModalVisible.value = false;
      }
    };

    // 下载证书图片（合成后）
    const downloadCertificate = async (award) => {
      if (!award.certificate_image) {
        message.error('该获奖记录暂无证书');
        return;
      }

      try {
        const loadingMessage = message.loading('正在生成证书...', 0);
        console.log('开始生成下载证书:', award.certificate_image.substring(0, 100)); // 调试用
        const url = await composeCertificateImage(award);
        loadingMessage();

        // 生成文件名
        const fileName = `${award.year}年${award.category || ''}${award.group_name}${award.award_level}证书.png`;
        console.log('准备下载证书:', fileName);
        
        // 创建下载链接
        const link = document.createElement('a');
        link.href = url;
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        
        message.success('证书下载成功');
      } catch (e) {
        console.error('证书生成失败:', e);
        message.error(e.message || '证书生成失败，请稍后重试');
      }
    };
    
    return {
      loading,
      hasQueried,
      queryResults,
      queryForm,
      availableYears,
      teamOptions,
      filteredTeamOptions,
      handleTeamSearch,
      handleTeamSelect,
      getAwardLevelColor,
      isTargetMember,
      onQuery,
      resetQuery,
      downloadCertificate,
      viewCertificate,
      certificateModalVisible,
      certificateModalImg,
      captchaRef,
      errorMessage,
      retryQuery,
      clearError
    }
  }
}
</script>

<style scoped>
.group-query-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #e3f0ff 0%, #f6fbff 100%);
  padding: 20px;
}

.header-section {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px 20px;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
}

.main-title {
  font-size: 2.6rem;
  font-weight: 700;
  color: #1976d2;
  margin-bottom: 16px;
  text-shadow: none;
}

.subtitle {
  display: none;
}

.query-section {
  max-width: 800px;
  margin: 0 auto 40px;
}

.query-card {
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1f2937;
}

.title-icon {
  font-size: 1.2rem;
}

.query-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 8px;
  height: 40px;
  font-weight: 600;
}

.reset-btn {
  border-radius: 8px;
  height: 40px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.results-section {
  max-width: 1200px;
  margin: 0 auto;
}

.results-card {
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.error-section {
  margin-bottom: 20px;
}

.error-alert {
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.error-alert .ant-alert-message {
  font-weight: 600;
  color: #cf1322;
}

.error-alert .ant-alert-description {
  color: #434343;
  line-height: 1.6;
  margin-top: 8px;
}

.error-alert .ant-alert-action {
  margin-top: 12px;
}

.error-alert .ant-alert-action .ant-btn {
  border-radius: 6px;
  font-weight: 500;
}

.no-results {
  padding: 40px;
  text-align: center;
}

.results-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 32px;
  margin-top: 32px;
}

.award-card {
  width: 420px;
  margin: 0 auto;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(102,126,234,0.08);
  padding: 28px 24px 18px 24px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.award-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  margin-bottom: 12px;
}

.award-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: #667eea;
  margin: 0;
}

.award-level {
  margin-left: 8px;
  font-size: 1rem;
}

.award-info-list {
  width: 100%;
  margin-bottom: 18px;
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  font-size: 1.08rem;
}

.info-label {
  font-weight: 600;
  color: #764ba2;
  min-width: 100px;
  text-align: right;
  margin-right: 10px;
}

.info-value {
  flex: 1;
  color: #222;
  word-break: break-all;
}

.members-section {
  margin-top: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 10px 12px;
}

.members-title {
  font-weight: 600;
  color: #667eea;
  margin-bottom: 6px;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.member-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}

.member-info {
  display: flex;
  flex-direction: column;
}

.member-name {
  font-weight: 500;
  color: #333;
}

.member-id {
  font-size: 0.95rem;
  color: #888;
}

/* 证书操作按钮样式 */
.award-actions {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-top: 18px;
}

.action-btns-center {
  display: flex;
  gap: 24px;
  justify-content: center;
  width: 100%;
}

.action-btn {
  border-radius: 8px;
  font-size: 1.08rem;
  font-weight: 500;
  padding: 0 22px;
  height: 44px;
  box-shadow: 0 2px 8px rgba(102,126,234,0.07);
  transition: box-shadow 0.2s, background 0.2s, color 0.2s, transform 0.2s;
}

.view-btn {
  background: #f5f7fa;
  color: #667eea;
  border: 1.5px solid #e0e7ff;
}

.view-btn:hover {
  background: #e0e7ff;
  color: #4c51bf;
  border-color: #c3bffa;
  transform: translateY(-1px) scale(1.03);
}

.download-btn {
  background: #3b82f6;
  color: #fff;
  border: none;
}

.download-btn:hover {
  background: #2563eb;
  color: #fff;
  transform: translateY(-1px) scale(1.03);
}

.close-modal-btn {
  font-size: 28px;
  color: #667eea;
  background: #fff;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(102,126,234,0.13);
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1200;
}

.search-tips {
  margin-top: 8px;
  padding-left: 10px;
}

.team-option {
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 4px;
  white-space: normal;
  word-break: break-word;
  line-height: 1.4;
}

.team-option:hover {
  background-color: #f0f5ff;
}

:deep(.ant-select-dropdown) {
  padding: 8px;
  border-radius: 8px;
  box-shadow: 0 3px 6px -4px rgba(0,0,0,.12), 0 6px 16px 0 rgba(0,0,0,.08), 0 9px 28px 8px rgba(0,0,0,.05);
}

:deep(.ant-auto-complete-option) {
  margin: 2px 0;
}

@media (max-width: 768px) {
  .header-section {
    padding: 24px 8px;
  }
  .header-content {
    padding: 0 4px;
    max-width: 100vw;
  }
  .main-title {
    text-align: center;
    word-break: break-all;
    max-width: 100vw;
    margin: 0 auto;
  }
  .query-section {
    max-width: 100vw;
    margin-left: auto;
    margin-right: auto;
    padding: 0 4px;
  }
  .results-grid {
    grid-template-columns: 1fr;
  }
  
  .group-query-container {
    padding: 10px;
  }
  
  .members-list {
    grid-template-columns: 1fr;
  }
}
@media (max-width: 480px) {
  .main-title {
    font-size: 1.2rem;
  }
  .query-section {
    padding: 0 2vw;
  }
}
@media (max-width: 600px) {
  .back-section {
    position: fixed;
    top: 10px;
    left: 10px;
    z-index: 2000;
    width: auto;
  }
  .back-btn {
    font-size: 1.1rem;
    padding: 10px 18px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(102,126,234,0.13);
    background: #fff;
    color: #667eea;
    border: 1.5px solid #e0e7ff;
  }
  .group-query-container {
    padding: 0 2vw 10vw 2vw;
  }
  .award-card {
    width: 98vw;
    padding: 14px 2vw 10px 2vw;
    min-width: unset;
    max-width: 100vw;
  }
  .award-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  .award-details {
    margin-bottom: 10px;
  }
  .member-actions {
    flex-direction: column;
    gap: 10px;
    align-items: stretch;
    width: 100%;
  }
  .action-btn {
    width: 100%;
    min-width: 0;
    font-size: 1rem;
    height: 40px;
    padding: 0 10px;
  }
  .award-title {
    font-size: 1.1rem;
  }
  .label {
    min-width: 80px;
    font-size: 1rem;
  }
  .detail-item {
    font-size: 1rem;
    margin-bottom: 7px;
  }
  .award-level {
    font-size: 1rem;
    padding: 0 8px;
  }
  .members-section {
    padding: 8px 4px;
  }
}
@media (max-width: 400px) {
  .award-card {
    padding: 8px 1vw 6px 1vw;
  }
  .award-title {
    font-size: 1rem;
  }
  .label {
    min-width: 60px;
    font-size: 0.95rem;
  }
  .detail-item {
    font-size: 0.95rem;
  }
}
</style>
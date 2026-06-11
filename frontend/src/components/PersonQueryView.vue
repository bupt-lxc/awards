<template>
  <div class="person-query-container">
    <!-- 页面头部 -->
    <div class="header-section">
      <div class="header-content">
        <h1 class="main-title">
          个人/团队获奖查询
        </h1>
        <p class="subtitle">查询个人或团队获奖信息</p>
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
            <a-col :xs="24" :sm="24" :md="8">
              <a-form-item label="查询类型" name="query_type">
                <a-radio-group v-model:value="queryForm.query_type" size="large" class="query-type-radio">
                  <a-radio-button value="person">个人获奖查询</a-radio-button>
                  <a-radio-button value="group">团队获奖查询</a-radio-button>
                </a-radio-group>
              </a-form-item>
            </a-col>

            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="获奖年份" name="year" :rules="[{ required: true, message: '请选择获奖年份' }]">
                <a-select 
                  v-model:value="queryForm.year" 
                  placeholder="请选择年份"
                  size="large"
                  show-search
                >
                  <a-select-option v-for="year in availableYears" :key="year" :value="year">
                    {{ year }}年
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="24" :md="12">
              <a-form-item 
                :label="queryForm.query_type === 'group' ? '团队名称' : '个人姓名'" 
                :name="queryForm.query_type === 'group' ? 'group_name' : 'person_name'" 
                :rules="unifiedRules"
              >
                <a-auto-complete
                  v-if="queryForm.query_type === 'group'"
                  v-model:value="queryForm.group_name"
                  placeholder="请输入团队名称关键词"
                  size="large"
                  :options="filteredGroupOptions"
                  @search="handleGroupSearch"
                  class="team-name-select"
                  :dropdownStyle="{ maxHeight: '400px', overflow: 'auto' }"
                >
                  <template #option="{ value, label }">
                    <div style="white-space: normal; word-break: break-all; padding: 8px 12px;">
                      {{ label }}
                    </div>
                  </template>
                </a-auto-complete>
                <a-input 
                  v-else
                  v-model:value="queryForm.person_name" 
                  placeholder="输入个人姓名"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="证书编码" name="id_number" :rules="[{ required: true, message: '请输入证书编码' }]">
                <a-input 
                  v-model:value="queryForm.id_number" 
                  placeholder="输入证书编码"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="组别">
                <a-input 
                  v-model:value="queryForm.category" 
                  placeholder="输入组别（可选）"
                  size="large"
                  allow-clear
                />
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
                  size="large" 
                  :loading="loading" 
                  @click="onQuery"
                  class="query-btn"
                >
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
                  重置条件
                </a-button>
              </div>
            </a-col>
          </a-row>
        </a-form>
      </a-card>
    </div>

    <!-- 查询结果区域 -->
    <div class="results-section" v-if="queryResults.length > 0 || hasQueried">
      <a-card class="results-card" :bordered="false">
        <template #title>
          <div class="card-title">
            <span class="title-icon">📋</span>
            查询结果
            <a-tag color="blue" v-if="queryResults.length > 0">
              共找到 {{ queryResults.length }} 条记录
            </a-tag>
          </div>
        </template>
        
        <div v-if="queryResults.length === 0 && hasQueried" class="no-results">
          <a-empty description="未找到符合条件的获奖记录" />
        </div>
        
        <div v-else class="results-grid">
          <div 
            v-for="(award, index) in queryResults" 
            :key="index" 
            class="award-card"
          >
            <div class="award-header">
              <h3 class="award-title">{{ queryForm.query_type === 'group' ? award.group_name : award.person_name }}</h3>
              <a-tag :color="getAwardLevelColor(award.award_level)" class="award-level-tag">
                {{ award.award_level }}
              </a-tag>
            </div>
            <div class="award-info-list">
              <div class="info-row"><span class="info-label">姓名/团队名称：</span><span class="info-value">{{ queryForm.query_type === 'group' ? award.group_name : award.person_name }}</span></div>
              <div class="info-row"><span class="info-label">证书编码：</span><span class="info-value">{{ queryForm.query_type === 'group' ? '-' : award.id_number }}</span></div>
              <div class="info-row"><span class="info-label">组别：</span><span class="info-value">{{ award.category || '-' }}</span></div>
              <div class="info-row"><span class="info-label">获奖年份：</span><span class="info-value">{{ award.year }}年</span></div>
              <div class="info-row"><span class="info-label">奖项等级：</span><span class="info-value"><a-tag :color="getAwardLevelColor(award.award_level)">{{ award.award_level }}</a-tag></span></div>
              <div class="info-row" v-if="award.description"><span class="info-label">获奖说明：</span><span class="info-value">{{ award.description }}</span></div>
            </div>
            <div class="award-actions" v-if="award.certificate_image">
              <div class="action-btns-center">
                <a-button type="default" size="large" class="action-btn view-btn" @click="viewCertificate(award)">
                  <template #icon><EyeOutlined /></template>
                  查看证书
                </a-button>
                <a-button type="primary" size="large" class="action-btn download-btn" @click="downloadCertificate(award)">
                  <template #icon><DownloadOutlined /></template>
                  下载证书
                </a-button>
              </div>
            </div>
          </div>
        </div>
      </a-card>
    </div>
    <!-- 证书大图弹窗 -->
    <a-modal v-model:open="certificateModalVisible" :footer="null" width="min(90vw, 500px)" title="获奖证书" @cancel="certificateModalVisible=false" :closable="true" wrapClassName="certificate-modal">
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
import { SearchOutlined, ReloadOutlined, EyeOutlined, DownloadOutlined, ArrowLeftOutlined, CloseOutlined } from '@ant-design/icons-vue'
import { queryPersonAwards, queryGroupAwards } from '../api/index.js'
import QRCode from 'qrcode'
import CaptchaInput from './CaptchaInput.vue'
import { fireConfetti } from '../utils/confetti'

export default {
  name: 'PersonQueryView',
  components: {
    SearchOutlined,
    ReloadOutlined,
    EyeOutlined,
    DownloadOutlined,
    ArrowLeftOutlined,
    CloseOutlined,
    CaptchaInput,
    VNodes: {
      functional: true,
      render: (_, ctx) => ctx.attrs.vnodes
    }
  },
  emits: ['go-home'],
  setup(props, { emit }) {
    const loading = ref(false)
    const hasQueried = ref(false)
    const queryResults = ref([])
    const captchaRef = ref(null);
    const queryForm = reactive({
      query_type: 'person', // 默认选择个人查询
      year: null,
      person_name: '',
      id_number: '',
      category: '',
      captcha: '',
      group_name: '',
      custom_group_name: ''
    })

    // 团队名称选项（隐藏的完整列表）
    const groupNameOptions = ref([
      { value: 'Balatime choir №3(Kazakhstan,Astana)', label: 'Balatime choir №3(Kazakhstan,Astana)' },
      { value: '广州市花都区青少年文化宫花之蕾合唱B团', label: '广州市花都区青少年文化宫花之蕾合唱B团' },
      { value: '武汉市光谷第一小学布谷童声合唱团', label: '武汉市光谷第一小学布谷童声合唱团' },
      { value: '福州糖果少儿合唱团', label: '福州糖果少儿合唱团' },
      { value: '江苏省东台市实验小学小红花合唱团', label: '江苏省东台市实验小学小红花合唱团' },
      { value: '广州市花都区花东镇中心小学花儿合唱团', label: '广州市花都区花东镇中心小学花儿合唱团' },
      { value: '杭州寻心少儿合唱团C团', label: '杭州寻心少儿合唱团C团' },
      { value: '新西兰The Cloud童声合唱团', label: '新西兰The Cloud童声合唱团' },
      { value: '合肥市屯溪路小学启明星合唱团', label: '合肥市屯溪路小学启明星合唱团' },
      { value: '珠海大剧院童声合唱团（扬帆少年合唱团）', label: '珠海大剧院童声合唱团（扬帆少年合唱团）' },
      { value: '茂名市茂南第一小学合唱团', label: '茂名市茂南第一小学合唱团' },
      { value: '广州市荔湾区乐贤坊小学合唱团', label: '广州市荔湾区乐贤坊小学合唱团' },
      { value: '杭州寻心少儿合唱团A团', label: '杭州寻心少儿合唱团A团' },
      { value: '武汉市青少年宫少年合唱团', label: '武汉市青少年宫少年合唱团' },
      { value: '北京佳音童声合唱团', label: '北京佳音童声合唱团' },
      { value: '苏州市吴韵少年艺术团·合唱团二团', label: '苏州市吴韵少年艺术团·合唱团二团' },
      { value: '华夏四季合唱团&东周小学合唱团', label: '华夏四季合唱团&东周小学合唱团' },
      { value: '广州市花都区秀全街九潭初级中学合唱团', label: '广州市花都区秀全街九潭初级中学合唱团' },
      { value: '杭州市天杭实验学校“天杭之声”合唱团', label: '杭州市天杭实验学校“天杭之声”合唱团（流行合唱）' },
      { value: '苏州吴韵少年艺术团·合唱团二团', label: '苏州吴韵少年艺术团·合唱团二团' },
      { value: 'State chamber choir of The Republic of Tatarstan(Russia)', label: 'State chamber choir of The Republic of Tatarstan(Russia)' },
      { value: '福鼎市茂门畲乡"金凤凰"童声合唱团', label: '福鼎市茂门畲乡"金凤凰"童声合唱团' },
      { value: '深圳繁星女声合唱团', label: '深圳繁星女声合唱团' },
      { value: '意大利Rom Ensemble合唱团', label: '意大利Rom Ensemble合唱团' },
      { value: '随园歌手', label: '随园歌手' },
      { value: 'Chamber choir "Bel Canto"(Russia,Omsk)', label: 'Chamber choir "Bel Canto"(Russia,Omsk)' },
      { value: '上海市学生艺术团普陀区青少年教育活动中心合唱团', label: '上海市学生艺术团普陀区青少年教育活动中心合唱团' },
      { value: '太阳照照童声合唱团', label: '太阳照照童声合唱团' },
      { value: '松狮小侬侬童声合唱团', label: '松狮小侬侬童声合唱团' },
      { value: '广州市黄广少年合唱团', label: '广州市黄广少年合唱团' },
      { value: '内蒙古小青马合唱团', label: '内蒙古小青马合唱团' },
      { value: '济南大野少年合唱团', label: '济南大野少年合唱团' },
      { value: '南山文体中心剧院小橙堡童声合唱团绽放团', label: '南山文体中心剧院小橙堡童声合唱团绽放团' },
      { value: '芜湖市少年宫合唱团', label: '芜湖市少年宫合唱团' },
      { value: '深圳市福田区外国语学校（景秀）合唱团', label: '深圳市福田区外国语学校（景秀）合唱团' },
      { value: '南山文体中心剧院小橙堡童声合唱团萌芽团', label: '南山文体中心剧院小橙堡童声合唱团萌芽团' },
      { value: '深圳保利天堂鸟童声合唱乐芽团', label: '深圳保利天堂鸟童声合唱乐芽团' },
      { value: '镇江星光童声合唱团-杜鹃团', label: '镇江星光童声合唱团-杜鹃团' },
      { value: '成都城市音乐厅青少年合唱团C2团', label: '成都城市音乐厅青少年合唱团C2团' },
      { value: '广州市花都区青少年文化宫花之蕾合唱A团', label: '广州市花都区青少年文化宫花之蕾合唱A团' },
      { value: '合肥市屯溪路小学启明星合唱团', label: '合肥市屯溪路小学启明星合唱团' },
      { value: '广州市花都区秀全街九潭初级中学合唱团', label: '广州市花都区秀全街九潭初级中学合唱团' },
      { value: '杭州市天杭实验学校“天杭之声”合唱团', label: '杭州市天杭实验学校“天杭之声”合唱团（少儿合唱）' },
      { value: '成都城市音乐厅青少年合唱团C1团', label: '成都城市音乐厅青少年合唱团C1团' }
    ])

    // 过滤后的选项
    const filteredGroupOptions = ref([])

    // 搜索处理函数
    const handleGroupSearch = (searchText) => {
      if (!searchText) {
        filteredGroupOptions.value = []
        return
      }

      // 将搜索文本转换为小写并分割为关键词数组
      const keywords = searchText.toLowerCase().split(/\s+/)

      // 对每个团队名称进行评分
      const scoredOptions = groupNameOptions.value.map(option => {
        const label = option.label.toLowerCase()
        // 计算每个关键词的匹配分数
        const score = keywords.reduce((acc, keyword) => {
          if (label.includes(keyword)) {
            // 如果关键词完全匹配，给更高的分数
            return acc + (label === keyword ? 2 : 1)
          }
          return acc
        }, 0)
        return { ...option, score }
      })

      // 过滤出有匹配分数的选项并按分数排序
      const matchedOptions = scoredOptions
        .filter(option => option.score > 0)
        .sort((a, b) => b.score - a.score)
        .slice(0, 5) // 只显示前5个最匹配的结果
        .map(({ value, label }) => ({ value, label }))

      filteredGroupOptions.value = matchedOptions
    }

    // 初始化时设置默认选中的团队
    onMounted(() => {
      if (queryForm.group_name && groupNameOptions.value.length > 0) {
        const found = groupNameOptions.value.find(opt => opt.value === queryForm.group_name)
        if (found) {
          queryForm.group_name = found.value
        }
      }
    })

    // 添加自定义团队名称
    const addCustomGroup = () => {
      if (queryForm.custom_group_name && queryForm.custom_group_name.trim()) {
        const newGroup = {
          value: queryForm.custom_group_name.trim(),
          label: queryForm.custom_group_name.trim()
        }
        if (!groupNameOptions.value.find(option => option.value === newGroup.value)) {
          groupNameOptions.value.push(newGroup)
        }
        queryForm.group_name = newGroup.value
        queryForm.custom_group_name = ''
      }
    }
    const availableYears = computed(() => {
      const currentYear = new Date().getFullYear()
      const years = []
      for (let i = currentYear; i >= currentYear - 20; i--) {
        years.push(i)
      }
      return years
    })
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
    const onQuery = async () => {
      try {
        // 验证验证码
        const captchaText = captchaRef.value.getCaptchaText();
        if (queryForm.captcha !== captchaText) {
          message.error('验证码错误');
          captchaRef.value.refreshCaptcha();
          return;
        }

        // 验证必填字段
        console.log('当前查询类型:', queryForm.query_type);
        console.log('个人姓名:', queryForm.person_name);
        console.log('团队名称:', queryForm.group_name);
        
        if (queryForm.query_type === 'person' && !queryForm.person_name) {
          message.error('请输入个人姓名');
          return;
        }
        if (queryForm.query_type === 'group' && !queryForm.group_name) {
          message.error('请选择团队名称');
          return;
        }

        loading.value = true
        hasQueried.value = true
        
        // 根据查询类型构建查询参数
        const queryParams = {
          year: queryForm.year,
          id_number: queryForm.id_number,
          category: queryForm.category || undefined
        };
        
        if (queryForm.query_type === 'person') {
          queryParams.person_name = queryForm.person_name;
        } else {
          queryParams.group_name = queryForm.group_name;
        }
        
        console.log('查询参数:', queryParams); // 添加调试日志

        let response;
        if (queryForm.query_type === 'group') {
          response = await queryGroupAwards(queryParams);
        } else {
          response = await queryPersonAwards(queryParams);
        }
        queryResults.value = (response.awards || []).map(a => ({ ...a, award_type: queryForm.query_type }))
        if (queryResults.value.length === 0) {
          message.info('未找到符合条件的获奖记录')
        } else {
          message.success(`查询成功，找到 ${queryResults.value.length} 条记录`)
          // 添加礼花效果
          fireConfetti()
        }
      } catch (error) {
        console.error('查询失败:', error)
        message.error(error.response?.data?.detail || '查询失败，请稍后重试')
        queryResults.value = []
      } finally {
        loading.value = false
        // 刷新验证码
        captchaRef.value.refreshCaptcha();
        queryForm.captcha = '';
      }
    }
    // 动态验证规则
    const personNameRules = computed(() => {
      if (queryForm.query_type === 'person') {
        return [{ required: true, message: '请输入个人姓名' }];
      }
      return [];
    });

    const groupNameRules = computed(() => {
      if (queryForm.query_type === 'group') {
        return [{ required: true, message: '请选择或输入团队名称' }];
      }
      return [];
    });

    // 统一的验证规则
    const unifiedRules = computed(() => {
      if (queryForm.query_type === 'person') {
        return personNameRules.value;
      } else {
        return groupNameRules.value;
      }
    });

    const resetQuery = () => {
      Object.assign(queryForm, {
        query_type: 'person',
        year: null,
        person_name: '',
        id_number: '',
        category: '',
        captcha: '',
        group_name: '',
        custom_group_name: ''
      })
      queryResults.value = []
      hasQueried.value = false
      // 刷新验证码
      captchaRef.value.refreshCaptcha();
    }
    const certificateModalVisible = ref(false);
    const certificateModalImg = ref('');
    const getVerifyUrl = (award) => {
      const baseUrl = window.location.origin;
      const params = new URLSearchParams({
        type: 'person',
        year: award.year,
        person_name: award.person_name,
        id_number: award.id_number,
        category: award.category || '',
        award_id: award.id
      });
      return `${baseUrl}/verify.html?${params.toString()}`;
    };
    const composeCertificateImage = async (award) => {
      return new Promise(async (resolve, reject) => {
        try {
          const img = new window.Image();
          img.crossOrigin = 'anonymous';
          img.src = award.certificate_image;
          img.onload = async () => {
            const canvas = document.createElement('canvas');
            canvas.width = img.width;
            canvas.height = img.height;
            const ctx = canvas.getContext('2d');
            ctx.drawImage(img, 0, 0);
            const qrSize = Math.floor(Math.min(img.width, img.height) * 0.13);
            const qrCodeDataUrl = await QRCode.toDataURL(getVerifyUrl(award), { width: qrSize, margin: 0 });
            const qrImg = new window.Image();
            qrImg.src = qrCodeDataUrl;
            qrImg.onload = () => {
              const padding = 32;
              const qrX = img.width - qrSize - padding;
              const qrY = padding;
              ctx.drawImage(qrImg, qrX, qrY, qrSize, qrSize);
              ctx.font = `${Math.floor(qrSize/5)}px Arial`;
              ctx.fillStyle = '#333';
              ctx.textAlign = 'center';
              ctx.shadowColor = 'white';
              ctx.shadowBlur = 2;
              ctx.fillText('防伪标识', qrX + qrSize/2, qrY + qrSize + Math.floor(qrSize/3));
              resolve(canvas.toDataURL('image/png'));
            };
            qrImg.onerror = reject;
          };
          img.onerror = reject;
        } catch (e) {
          reject(e);
        }
      });
    };
    const viewCertificate = async (award) => {
      certificateModalImg.value = '';
      certificateModalVisible.value = true;
      try {
        certificateModalImg.value = await composeCertificateImage(award);
      } catch (e) {
        message.error('证书图片加载失败');
        certificateModalVisible.value = false;
      }
    };
    const downloadCertificate = async (award) => {
      try {
        const url = await composeCertificateImage(award);
        const link = document.createElement('a');
        link.href = url;
        link.download = `${award.group_name || award.person_name || 'certificate'}.png`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } catch (e) {
        message.error('证书图片生成失败');
      }
    };
    // 判断是否为移动端
    const isMobile = computed(() => window.innerWidth <= 768)
    // 返回主页事件
    const emitGoHome = () => {
      emit('go-home')
    }
    return {
      loading,
      hasQueried,
      queryResults,
      queryForm,
      availableYears,
      getAwardLevelColor,
      onQuery,
      resetQuery,
      downloadCertificate,
      viewCertificate,
      certificateModalVisible,
      certificateModalImg,
      isMobile,
      emitGoHome,
      captchaRef,
      // 团队搜索相关
      filteredGroupOptions,
      handleGroupSearch,
      personNameRules,
      groupNameRules,
      unifiedRules
    }
  }
}
</script>

<style scoped>
.person-query-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #FFF9F0 0%, #FFF6E5 100%);
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
  color: #FF9B9B;
  margin-bottom: 16px;
  text-shadow: none;
}

.subtitle {
  display: none;
}

.query-section {
  margin-bottom: 40px;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
}

.query-card {
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(255, 155, 155, 0.08);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 155, 155, 0.1);
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
  background: linear-gradient(135deg, #FF9B9B 0%, #FFD699 100%);
  border: none;
  border-radius: 8px;
  height: 40px;
  font-weight: 600;
}

.action-buttons {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.reset-btn {
  border-radius: 8px;
  height: 40px;
}

.results-section {
  max-width: 1200px;
  margin: 0 auto;
}

.results-card {
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(255, 155, 155, 0.08);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 155, 155, 0.1);
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
  width: 400px;
  margin: 0 auto;
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 4px 24px rgba(255, 155, 155, 0.08);
  padding: 32px 28px 22px 28px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  transition: box-shadow 0.2s, transform 0.2s;
  border: 1.5px solid #FFE5E5;
}

.award-card:hover {
  box-shadow: 0 8px 32px rgba(255, 155, 155, 0.12);
  transform: translateY(-2px) scale(1.012);
}

.award-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  width: 100%;
  margin-bottom: 18px;
  gap: 12px;
}

.award-title {
  font-size: 1.35rem;
  font-weight: 700;
  color: #FF9B9B;
  margin: 0;
}

.award-level-tag {
  font-size: 1.08rem;
  font-weight: 600;
  border-radius: 8px;
  padding: 0 12px;
  margin-left: 8px;
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
  color: #FFB347;
  min-width: 120px;
  text-align: right;
  margin-right: 10px;
  white-space: nowrap;
}

.info-value {
  flex: 1;
  color: #222;
  word-break: break-all;
}

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
  box-shadow: 0 2px 8px rgba(255, 155, 155, 0.05);
  transition: box-shadow 0.2s, background 0.2s, color 0.2s, transform 0.2s;
}

.view-btn {
  background: #FFF9F0;
  color: #FF9B9B;
  border: 1.5px solid #FFE5E5;
}

.view-btn:hover {
  background: #FFE5E5;
  color: #FF7B7B;
  border-color: #FFD6D6;
  transform: translateY(-1px) scale(1.03);
}

.download-btn {
  background: linear-gradient(135deg, #FF9B9B 0%, #FFD699 100%);
  color: #fff;
  border: none;
}

.download-btn:hover {
  background: linear-gradient(135deg, #FF7B7B 0%, #FFB347 100%);
  color: #fff;
  transform: translateY(-1px) scale(1.03);
}

.mobile-back-btn {
  position: fixed;
  top: 16px;
  left: 16px;
  z-index: 1100;
  font-size: 28px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  border-radius: 50%;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.close-modal-btn {
  font-size: 28px;
  color: #FF9B9B;
  background: #FFF9F0;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(255, 155, 155, 0.1);
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1200;
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
  .person-query-container {
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
  .award-info-list {
    margin-bottom: 10px;
  }
  .award-actions, .action-btns-center {
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
  .info-label {
    min-width: 100px;
    font-size: 1rem;
  }
  .info-row {
    font-size: 1rem;
    margin-bottom: 7px;
  }
  .award-level-tag {
    font-size: 1rem;
    padding: 0 8px;
  }
  
  /* 移动端团队名称选择框优化 */
  .team-name-select :deep(.ant-select-selector) {
    min-height: 44px;
    padding: 6px 12px;
  }
  
  .team-name-select :deep(.ant-select-selection-item) {
    font-size: 13px;
    line-height: 1.4;
  }
  
  .team-name-select :deep(.ant-select-dropdown) {
    max-width: 90vw !important;
    min-width: 280px !important;
  }
  
  .mobile-back-btn {
    top: 10px;
    left: 10px;
    width: 44px;
    height: 44px;
    font-size: 24px;
  }
  .close-modal-btn {
    font-size: 26px;
    width: 40px;
    height: 40px;
  }
}
@media (max-width: 400px) {
  .award-card {
    padding: 8px 1vw 6px 1vw;
  }
  .award-title {
    font-size: 1rem;
  }
  .info-label {
    min-width: 90px;
    font-size: 0.95rem;
  }
  .info-row {
    font-size: 0.95rem;
  }
}

/* 查询类型单选按钮样式 */
.query-type-radio {
  width: 100%;
  display: flex;
  justify-content: flex-start;
  gap: 8px;
}

.query-type-radio :deep(.ant-radio-button-wrapper) {
  flex: 1;
  text-align: center;
  height: 40px;
  line-height: 38px;
  font-size: 14px;
  border-radius: 6px;
  transition: all 0.3s ease;
  border: 1px solid #d9d9d9;
  margin-right: 0;
  border-left-width: 1px !important;
}

.query-type-radio :deep(.ant-radio-button-wrapper:first-child) {
  border-radius: 6px;
  border-left-width: 1px !important;
}

.query-type-radio :deep(.ant-radio-button-wrapper:last-child) {
  border-radius: 6px;
}

.query-type-radio :deep(.ant-radio-button-wrapper::before) {
  display: none !important;
}

.query-type-radio :deep(.ant-radio-button-wrapper-checked) {
  background: linear-gradient(135deg, #FF9B9B 0%, #FFD699 100%);
  border-color: transparent;
  box-shadow: 0 2px 8px rgba(255, 155, 155, 0.2);
  color: white;
}

.query-type-radio :deep(.ant-radio-button-wrapper-checked:hover) {
  background: linear-gradient(135deg, #FF7B7B 0%, #FFB347 100%);
  border-color: transparent;
  color: white;
}

.query-type-radio :deep(.ant-radio-button-wrapper:hover) {
  color: #e53e3e;
}

/* 团队名称选择框样式 */
.team-name-select {
  width: 100%;
}

.team-name-select :deep(.ant-select-selector) {
  min-height: 40px;
  height: auto !important;
  padding: 4px 11px;
  border-radius: 6px;
}

.team-name-select :deep(.ant-select-selection-item) {
  white-space: normal !important;
  word-break: break-all !important;
  height: auto !important;
  line-height: 1.5 !important;
  padding: 2px 0;
  overflow: visible;
  text-overflow: unset;
  font-size: 14px;
  color: #1a1a1a;
}

.team-name-select :deep(.ant-select-selection-search) {
  height: auto !important;
  line-height: 1.5 !important;
}

.team-name-select :deep(.ant-select-selection-placeholder) {
  line-height: 32px;
}

/* 确保下拉选项也能完整显示 */
.team-name-select :deep(.ant-select-dropdown) {
  max-width: 600px !important;
  min-width: 300px !important;
}

.team-name-select :deep(.ant-select-item) {
  padding: 8px 12px;
}

.team-name-select :deep(.ant-select-item-option-content) {
  white-space: normal !important;
  word-break: break-all !important;
  line-height: 1.5 !important;
  font-size: 14px;
}

:deep(.ant-select-dropdown) {
  max-width: 400px !important;
}

:deep(.ant-select-item-option-content) {
  white-space: normal;
  word-break: break-all;
  padding: 8px 0;
  line-height: 1.5;
}
</style>
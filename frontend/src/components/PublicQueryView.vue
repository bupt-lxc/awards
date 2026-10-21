<template>
  <div class="public-query-container">
    <!-- 页面头部 -->
    <div class="header-section">
      <div class="header-content">
        <h1 class="main-title">
          <!-- <span class="title-icon">🏆</span> -->
          获奖记录查询
        </h1>
        <p class="subtitle">查询各类获奖信息，支持多种筛选条件</p>
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
                  <a-radio-button value="person">个人查询</a-radio-button>
                  <a-radio-button value="group">团队成员查询</a-radio-button>
                </a-radio-group>
              </a-form-item>
            </a-col>

            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="获奖年份" name="year">
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
            
            <a-col :xs="24" :sm="12" :md="6" v-if="queryForm.query_type === 'group'">
              <a-form-item label="团队名称" name="group_name">
                <a-input
                  v-model:value="queryForm.group_name"
                  placeholder="输入团队名称（可选）"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6" v-if="queryForm.query_type === 'person'">
              <a-form-item label="个人姓名" name="person_name">
                <a-input
                  v-model:value="queryForm.person_name"
                  placeholder="输入个人姓名（可选）"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="参赛编码" name="id_number">
                <a-input
                  v-model:value="queryForm.id_number"
                  placeholder="输入参赛编码（可选）"
                  size="large"
                  allow-clear
                />
              </a-form-item>
            </a-col>
            
            <a-col :xs="24" :sm="12" :md="6">
              <a-form-item label="组别">
                <a-input 
                  v-model:value="queryForm.category" 
                  placeholder="输入组别"
                  size="large"
                  allow-clear
                />
              </a-form-item>
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
    <div class="results-section" v-if="queryResults.length > 0 || hasQueried">
      <a-card class="results-card" :bordered="false">
        <template #title>
          <div class="card-title">
            <span class="title-icon">📋</span>
            查询结果
            <a-tag color="blue" class="result-count">{{ queryResults.length }} 条记录</a-tag>
          </div>
        </template>
        
        <div v-if="queryResults.length === 0 && hasQueried" class="no-results">
          <a-empty description="未找到符合条件的获奖记录">
            <template #image>
              <div class="empty-icon">🔍</div>
            </template>
          </a-empty>
        </div>
        
        <div v-else class="results-grid">
          <div 
            v-for="record in queryResults" 
            :key="record.id" 
            class="result-card"
            @click="viewDetail(record)"
          >
            <div class="card-header">
              <div class="award-type">
                <a-tag 
                  :color="record.award_type === 'group' ? '#1890ff' : '#52c41a'"
                  class="type-tag"
                >
                  {{ record.award_type === 'group' ? '团体奖项' : '个人奖项' }}
                </a-tag>
              </div>
              <div class="award-year">{{ record.year }}年</div>
            </div>
            
            <div class="card-content">
              <h3 class="award-title">
                {{ record.award_type === 'group' ? record.group_name : record.person_name }}
              </h3>
              <div class="award-details">
                <div class="detail-item">
                  <span class="detail-label">组别：</span>
                  <span class="detail-value">{{ record.category }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">等级：</span>
                  <a-tag :color="getAwardLevelColor(record.award_level)" class="level-tag">
                    {{ record.award_level }}
                  </a-tag>
                </div>
              </div>
              <div class="award-description" v-if="record.description">
                {{ record.description.length > 60 ? record.description.substring(0, 60) + '...' : record.description }}
              </div>
            </div>
            
            <div class="card-footer">
              <span class="view-detail">点击查看详情 →</span>
            </div>
          </div>
        </div>
      </a-card>
    </div>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="获奖记录详情"
      :footer="null"
      width="600px"
      class="detail-modal"
    >
      <div v-if="selectedRecord" class="detail-content">
        <div class="detail-header">
          <h2 class="detail-title">
            {{ selectedRecord.award_type === 'group' ? selectedRecord.group_name : selectedRecord.person_name }}
          </h2>
          <div class="detail-tags">
            <a-tag 
              :color="selectedRecord.award_type === 'group' ? '#1890ff' : '#52c41a'"
              class="type-tag-large"
            >
              {{ selectedRecord.award_type === 'group' ? '团体奖项' : '个人奖项' }}
            </a-tag>
            <a-tag color="orange" class="year-tag">{{ selectedRecord.year }}年</a-tag>
          </div>
        </div>
        
        <a-divider />
        
        <div class="detail-info">
          <div class="info-row">
            <div class="info-item">
              <div class="info-label">组别</div>
              <div class="info-value">{{ selectedRecord.category }}</div>
            </div>
            <div class="info-item">
              <div class="info-label">奖项等级</div>
              <div class="info-value">
                <a-tag :color="getAwardLevelColor(selectedRecord.award_level)" class="level-tag-large">
                  {{ selectedRecord.award_level }}
                </a-tag>
              </div>
            </div>
          </div>
          
          <div class="info-row" v-if="selectedRecord.award_type === 'person' && selectedRecord.id_number">
            <div class="info-item">
              <div class="info-label">证件号码</div>
              <div class="info-value">{{ selectedRecord.id_number }}</div>
            </div>
          </div>
          
          <div class="info-row" v-if="selectedRecord.award_type === 'group' && selectedRecord.members && selectedRecord.members.length > 0">
            <div class="info-item full-width">
              <div class="info-label">团体成员</div>
              <div class="info-value">
                <div class="members-list">
                  <a-tag 
                    v-for="(member, index) in selectedRecord.members" 
                    :key="index" 
                    class="member-tag"
                    color="blue"
                  >
                    {{ member.name }} ({{ member.id_number }})
                  </a-tag>
                </div>
              </div>
            </div>
          </div>
          
          <div class="info-row" v-if="selectedRecord.description">
            <div class="info-item full-width">
              <div class="info-label">获奖描述</div>
              <div class="info-value description">{{ selectedRecord.description }}</div>
            </div>
          </div>
          
          <div class="info-row">
            <div class="info-item">
              <div class="info-label">创建时间</div>
              <div class="info-value">{{ formatDate(selectedRecord.created_at) }}</div>
            </div>
            <div class="info-item" v-if="selectedRecord.updated_at">
              <div class="info-label">更新时间</div>
              <div class="info-value">{{ formatDate(selectedRecord.updated_at) }}</div>
            </div>
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 修改logo区域 -->
    <div class="logo-container">
      <img src="../assets/cca-logo.png" alt="中国合唱协会" class="org-logo" />
      <img src="../assets/cccf-logo.png" alt="中国少年儿童合唱节" class="festival-logo" />
    </div>

    <!-- ICP备案号 -->
    <div class="icp-container">
      <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer" class="icp-link">
        备案/许可证编号为：苏ICP备2026042125号
      </a>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, computed, h } from 'vue';
import { message } from 'ant-design-vue';
import { SearchOutlined, ReloadOutlined, TrophyOutlined, UserOutlined, TeamOutlined, CheckOutlined } from '@ant-design/icons-vue';
import { awardAPI } from '../api/index.js';
import { useRouter } from 'vue-router';
import ccaLogo from '../assets/cca-logo.png';
import cccfLogo from '../assets/cccf-logo.png';
import { fireConfetti } from '../utils/confetti';

const loading = ref(false);
const hasQueried = ref(false);
const queryResults = ref([]);
const detailModalVisible = ref(false);
const selectedRecord = ref(null);

const queryForm = reactive({
  query_type: 'person', // 默认选择个人查询
  year: new Date().getFullYear(),
  group_name: '',
  custom_group_name: '',
  person_name: '',
  id_number: '',
  category: ''
});

// 团队名称选项
const groupNameOptions = ref([
  { value: 'Balatime choir №3(Kazakhstan,Astana)', label: 'Balatime choir №3(Kazakhstan,Astana)' },
  { value: '广州市花都区青少年文化宫花之蕾合唱团', label: '广州市花都区青少年文化宫花之蕾合唱团' },
  { value: '武汉市光谷第一小学布谷童声合唱团', label: '武汉市光谷第一小学布谷童声合唱团' },
  { value: '福州糖果少儿合唱团', label: '福州糖果少儿合唱团' },
  { value: '江苏省东台市实验小学小红花合唱团', label: '江苏省东台市实验小学小红花合唱团' },
  { value: '广州市花都区花东镇中心小学花儿合唱团', label: '广州市花都区花东镇中心小学花儿合唱团' },
  { value: '杭州寻心少儿合唱团C团', label: '杭州寻心少儿合唱团C团' },
  { value: '新西兰The Cloud童声合唱团', label: '新西兰The Cloud童声合唱团' },
  // ... 其他团队名称
]);

// 搜索过滤函数
const filterOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

// 添加自定义团队名称
const addCustomGroup = () => {
  if (queryForm.custom_group_name && queryForm.custom_group_name.trim()) {
    const newGroup = {
      value: queryForm.custom_group_name.trim(),
      label: queryForm.custom_group_name.trim()
    };
    if (!groupNameOptions.value.find(option => option.value === newGroup.value)) {
      groupNameOptions.value.push(newGroup);
    }
    queryForm.group_name = newGroup.value;
    queryForm.custom_group_name = '';
  }
};

const router = useRouter();

// 生成可选年份
const availableYears = computed(() => {
  const currentYear = new Date().getFullYear();
  const years = [];
  for (let year = currentYear; year >= 2000; year--) {
    years.push(year);
  }
  return years;
});

// 动态验证规则
const yearRules = computed(() => {
  return [{ required: true, message: '请选择年份' }];
});

const personNameRules = computed(() => {
  if (queryForm.query_type === 'person') {
    return [{ required: true, message: '请输入个人姓名' }];
  }
  return [];
});

const idNumberRules = computed(() => {
  return [{ required: true, message: '请输入参赛编码' }];
});

const groupNameRules = computed(() => {
  if (queryForm.query_type === 'group') {
    return [{ required: true, message: '请选择或输入团队名称' }];
  }
  return [];
});

// VNodes 渲染函数
const VNodes = {
  functional: true,
  render: (_, ctx) => ctx.attrs.vnodes,
};

// 查询获奖记录
const onQuery = async () => {
  try {
    loading.value = true;
    hasQueried.value = true;
    
    const queryParams = {
      year: queryForm.year || undefined
    };
    
    // 只有在字段有值时才添加到查询参数
    if (queryForm.person_name && queryForm.person_name.trim()) {
      queryParams.person_name = queryForm.person_name.trim();
    }
    if (queryForm.id_number && queryForm.id_number.trim()) {
      queryParams.id_number = queryForm.id_number.trim();
    }
    if (queryForm.group_name && queryForm.group_name.trim()) {
      queryParams.group_name = queryForm.group_name.trim();
    }
    if (queryForm.category && queryForm.category.trim()) {
      queryParams.category = queryForm.category.trim();
    }
    
    // 根据查询类型调用不同的API
    let response;
    console.log('[PublicQueryView] 查询类型:', queryForm.query_type, '查询参数:', JSON.stringify(queryParams));
    if (queryForm.query_type === 'group') {
      response = await awardAPI.queryGroupAwards(queryParams);
    } else {
      response = await awardAPI.queryPersonAwards(queryParams);
    }
    console.log('[PublicQueryView] 原始返回:', JSON.stringify(response));
    // 给每条结果标记 award_type，模板用它区分渲染
    const awardType = queryForm.query_type === 'group' ? 'group' : 'person';
    queryResults.value = (response.awards || []).map(a => ({ ...a, award_type: awardType }));
    
    message.success(`查询完成，找到 ${queryResults.value.length} 条记录`);
    
    // 如果有结果，添加礼花效果
    if (queryResults.value.length > 0) {
      fireConfetti();
    }
    
  } catch (error) {
    console.error('查询失败:', error);
    
    if (error.response?.data?.detail) {
      message.error(`查询失败: ${error.response.data.detail}`);
    } else {
      message.error('查询失败，请检查网络连接');
    }
    
    queryResults.value = [];
  } finally {
    loading.value = false;
  }
};

// 重置查询
const resetQuery = () => {
  Object.assign(queryForm, {
    query_type: 'person', // 重置为个人查询
    year: new Date().getFullYear(),
    group_name: '',
    custom_group_name: '',
    person_name: '',
    id_number: '',
    category: ''
  });
  queryResults.value = [];
  hasQueried.value = false;
};

// 查看详情
const viewDetail = (record) => {
  selectedRecord.value = record;
  detailModalVisible.value = true;
};

// 获取奖项等级颜色
const getAwardLevelColor = (level) => {
  const colorMap = {
    '金奖': '#faad14',
    '银奖': '#bfbfbf',
    '铜奖': '#d4b106',
    '特等奖': '#f5222d',
    '一等奖': '#fa541c',
    '二等奖': '#fa8c16',
    '三等奖': '#fadb14',
    '优秀奖': '#1890ff'
  };
  return colorMap[level] || '#d9d9d9';
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleString('zh-CN');
};

const goToPersonQuery = () => {
  router.push('/person-query');
};

const goToGroupQuery = () => {
  router.push('/group-query');
};

onMounted(() => {
  // 页面加载时可以执行一些初始化操作
});
</script>

<style scoped>
.public-query-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 0;
}

.header-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 0 80px;
  text-align: center;
  color: white;
  position: relative;
  overflow: hidden;
}

.header-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="50" cy="50" r="1" fill="%23ffffff" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>') repeat;
  pointer-events: none;
}

.header-content {
  position: relative;
  z-index: 1;
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}

.main-title {
  font-size: 3.5rem;
  font-weight: 700;
  margin: 0 0 20px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  letter-spacing: 2px;
}

.title-icon {
  font-size: 3rem;
  margin-right: 15px;
  display: inline-block;
  /* animation: bounce 2s infinite; */
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-10px);
  }
  60% {
    transform: translateY(-5px);
  }
}

.subtitle {
  font-size: 1.3rem;
  opacity: 0.9;
  margin: 0;
  font-weight: 300;
}

.query-section {
  margin-top: -40px;
  padding: 0 20px 40px;
  position: relative;
  z-index: 2;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
}

.query-card {
  max-width: 1200px;
  margin: 0 auto;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.query-card :deep(.ant-card-head) {
  background: linear-gradient(135deg, #f6f9fc 0%, #ffffff 100%);
  border-bottom: 1px solid #e8f4fd;
}

.query-card :deep(.ant-card-body) {
  padding: 40px;
}

.card-title {
  font-size: 1.4rem;
  font-weight: 600;
  color: #2c3e50;
  display: flex;
  align-items: center;
}

.card-title .title-icon {
  font-size: 1.5rem;
  margin-right: 10px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.query-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 10px;
  height: 45px;
  padding: 0 30px;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.query-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.btn-text {
  position: relative;
  z-index: 2;
}

.btn-ripple {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
}

.query-btn:active .btn-ripple {
  width: 300px;
  height: 300px;
}

.reset-btn {
  border-radius: 10px;
  height: 45px;
  padding: 0 30px;
  font-weight: 600;
  border: 2px solid #d9d9d9;
  transition: all 0.3s ease;
}

.reset-btn:hover {
  border-color: #667eea;
  color: #667eea;
}

.results-section {
  padding: 0 20px 60px;
}

.results-card {
  max-width: 1200px;
  margin: 0 auto;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.results-card :deep(.ant-card-head) {
  background: linear-gradient(135deg, #f6f9fc 0%, #ffffff 100%);
  border-bottom: 1px solid #e8f4fd;
}

.result-count {
  margin-left: 15px;
  font-weight: 600;
}

.no-results {
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
  padding: 20px 0;
}

.result-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  cursor: pointer;
  border: 1px solid #f0f0f0;
  position: relative;
  overflow: hidden;
}

.result-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.result-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.result-card:hover::before {
  transform: scaleX(1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.type-tag {
  font-weight: 600;
  border-radius: 20px;
  padding: 4px 12px;
}

.award-year {
  font-size: 0.9rem;
  color: #8c8c8c;
  font-weight: 500;
}

.award-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 16px;
  line-height: 1.4;
}

.award-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  align-items: center;
}

.detail-label {
  font-size: 0.9rem;
  color: #8c8c8c;
  margin-right: 8px;
  min-width: 50px;
}

.detail-value {
  font-weight: 500;
  color: #2c3e50;
}

.level-tag {
  font-weight: 600;
  border-radius: 12px;
}

.members-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
}

.member-tag {
  margin: 0;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
}

.award-description {
  font-size: 0.9rem;
  color: #666;
  line-height: 1.5;
  margin-bottom: 16px;
  background: #f8f9fa;
  padding: 12px;
  border-radius: 8px;
}

.card-footer {
  text-align: right;
}

.view-detail {
  font-size: 0.9rem;
  color: #667eea;
  font-weight: 500;
  transition: color 0.3s ease;
}

.result-card:hover .view-detail {
  color: #764ba2;
}

.detail-modal :deep(.ant-modal-header) {
  background: linear-gradient(135deg, #f6f9fc 0%, #ffffff 100%);
  border-bottom: 1px solid #e8f4fd;
}

.detail-content {
  padding: 20px 0;
}

.detail-header {
  text-align: center;
  margin-bottom: 20px;
}

.detail-title {
  font-size: 1.8rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 16px;
}

.detail-tags {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.type-tag-large {
  font-weight: 600;
  border-radius: 20px;
  padding: 6px 16px;
  font-size: 0.9rem;
}

.year-tag {
  font-weight: 600;
  border-radius: 20px;
  padding: 6px 16px;
  font-size: 0.9rem;
}

.detail-info {
  margin-top: 20px;
}

.info-row {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
}

.info-item {
  flex: 1;
}

.info-item.full-width {
  flex: none;
  width: 100%;
}

.info-label {
  font-size: 0.9rem;
  color: #8c8c8c;
  margin-bottom: 8px;
  font-weight: 500;
}

.info-value {
  font-size: 1rem;
  color: #2c3e50;
  font-weight: 600;
}

.info-value.description {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
  line-height: 1.6;
  font-weight: 400;
}

.level-tag-large {
  font-weight: 600;
  border-radius: 12px;
  padding: 4px 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-title {
    font-size: 2.1rem;
    text-align: center;
    word-break: break-all;
    max-width: 100vw;
    margin: 0 auto;
  }
  .subtitle {
    font-size: 1rem;
    text-align: center;
    word-break: break-all;
    max-width: 100vw;
    margin: 0 auto;
  }
  .query-section {
    max-width: 100vw;
    margin-left: auto;
    margin-right: auto;
    padding: 0 4px 24px 4px;
  }
}
@media (max-width: 480px) {
  .main-title {
    font-size: 1.1rem;
  }
  .subtitle {
    font-size: 0.8rem;
  }
  .query-section {
    padding: 0 2vw 12vw 2vw;
  }
}

.logo-container {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 60px;
  margin-top: 40px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 800px;
}

.org-logo {
  height: 80px;
  width: auto;
  object-fit: contain;
}

.festival-logo {
  height: 100px;
  width: auto;
  object-fit: contain;
}

.icp-container {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 16px;
  padding: 8px 0 20px 0;
}

.icp-link {
  color: #666;
  font-size: 13px;
  text-decoration: none;
  transition: color 0.3s;
}

.icp-link:hover {
  color: #1677ff;
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
}

.query-type-radio :deep(.ant-radio-button-wrapper:first-child) {
  border-radius: 6px;
}

.query-type-radio :deep(.ant-radio-button-wrapper:last-child) {
  border-radius: 6px;
}

.query-type-radio :deep(.ant-radio-button-wrapper-checked) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.query-type-radio :deep(.ant-radio-button-wrapper-checked:hover) {
  background: linear-gradient(135deg, #5a71d6 0%, #6a43a1 100%);
  border-color: transparent;
}

.query-type-radio :deep(.ant-radio-button-wrapper:hover) {
  color: #667eea;
}
</style>
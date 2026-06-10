<template>
  <div class="monitor-container">
    <div class="monitor-header">
      <h2>系统监控面板</h2>
      <div class="header-actions">
        <a-button @click="refreshData" :loading="loading" type="primary">
          <reload-outlined />
          刷新数据
        </a-button>
        <a-button @click="resetStats" :loading="resetting" danger>
          <delete-outlined />
          重置统计
        </a-button>
        <a-switch v-model:checked="autoRefresh" @change="toggleAutoRefresh">
          <template #checkedChildren>自动刷新</template>
          <template #unCheckedChildren>手动刷新</template>
        </a-switch>
      </div>
    </div>

    <!-- 系统健康状态 -->
    <a-row :gutter="16" class="monitor-cards">
      <a-col :span="6">
        <a-card class="health-card" :class="healthStatusClass">
          <a-statistic
            title="系统健康度"
            :value="systemHealth.health_score"
            suffix="%"
            :value-style="{ color: healthColor }"
          />
          <div class="health-status">{{ healthStatusText }}</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="在线用户"
            :value="monitorData.online_count"
            prefix="👥"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="QPS"
            :value="monitorData.request_stats?.qps || 0"
            :precision="2"
            prefix="⚡"
          />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="平均响应时间"
            :value="monitorData.request_stats?.avg_response_time || 0"
            suffix="ms"
            :precision="2"
            prefix="⏱️"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- 系统指标 -->
    <a-row :gutter="16" class="monitor-section">
      <a-col :span="12">
        <a-card title="系统资源使用情况" class="metrics-card">
          <div class="metric-item">
            <div class="metric-label">内存使用率</div>
            <a-progress 
              :percent="monitorData.system_metrics?.memory_usage || 0" 
              :status="getProgressStatus(monitorData.system_metrics?.memory_usage)"
              :show-info="true"
            />
            <div class="metric-detail">
              {{ formatBytes(monitorData.system_metrics?.memory_used) }} / {{ formatBytes(monitorData.system_metrics?.memory_total) }}
            </div>
          </div>
          
          <div class="metric-item">
            <div class="metric-label">Goroutine 数量</div>
            <a-progress 
              :percent="Math.min((monitorData.system_metrics?.goroutines || 0) / 10, 100)" 
              :status="getProgressStatus((monitorData.system_metrics?.goroutines || 0) / 10)"
              :show-info="true"
            />
            <div class="metric-detail">{{ monitorData.system_metrics?.goroutines || 0 }} 个</div>
          </div>
          
          <div class="metric-item">
            <div class="metric-label">数据库连接数</div>
            <a-progress 
              :percent="Math.min((monitorData.system_metrics?.db_connections || 0) * 2, 100)" 
              :status="getProgressStatus((monitorData.system_metrics?.db_connections || 0) * 2)"
              :show-info="true"
            />
            <div class="metric-detail">{{ monitorData.system_metrics?.db_connections || 0 }} 个</div>
          </div>
          
          <div class="metric-item">
            <div class="metric-label">系统运行时间</div>
            <div class="metric-detail">{{ formatUptime(monitorData.system_metrics?.uptime) }}</div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :span="12">
        <a-card title="请求统计" class="metrics-card">
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-value">{{ monitorData.request_stats?.total_requests || 0 }}</div>
              <div class="stat-label">总请求数</div>
            </div>
            <div class="stat-item success">
              <div class="stat-value">{{ monitorData.request_stats?.success_requests || 0 }}</div>
              <div class="stat-label">成功请求</div>
            </div>
            <div class="stat-item error">
              <div class="stat-value">{{ monitorData.request_stats?.error_requests || 0 }}</div>
              <div class="stat-label">错误请求</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ getErrorRate() }}%</div>
              <div class="stat-label">错误率</div>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 在线用户列表 -->
    <a-card title="在线用户" class="monitor-section">
      <template #extra>
        <a-button @click="refreshUsers" :loading="loading" size="small">
          <reload-outlined />
          刷新
        </a-button>
      </template>
      
      <a-table
        :dataSource="monitorData.online_users"
        :columns="userColumns"
        :pagination="{ pageSize: 10 }"
        size="small"
        :loading="loading"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'is_admin'">
            <a-tag :color="record.is_admin ? 'red' : 'blue'">
              {{ record.is_admin ? '管理员' : '普通用户' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'last_active'">
            {{ formatTime(record.last_active) }}
          </template>
          <template v-if="column.key === 'login_time'">
            {{ formatTime(record.login_time) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-button 
              size="small" 
              danger 
              @click="kickUser(record.user_id)"
              :disabled="record.is_admin"
            >
              踢出
            </a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 请求历史 -->
    <a-card title="最近请求" class="monitor-section">
      <template #extra>
        <a-button @click="refreshRequests" :loading="loading" size="small">
          <reload-outlined />
          刷新
        </a-button>
      </template>
      
      <a-table
        :dataSource="requestHistory"
        :columns="requestColumns"
        :pagination="{ pageSize: 20 }"
        size="small"
        :loading="loading"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ record.status }}
            </a-tag>
          </template>
          <template v-if="column.key === 'timestamp'">
            {{ formatTime(record.timestamp) }}
          </template>
          <template v-if="column.key === 'response_time'">
            <span :class="{ 'slow-request': record.response_time > 1000 }">
              {{ record.response_time }}ms
            </span>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 数据库统计 -->
    <a-card title="数据库统计" class="monitor-section">
      <template #extra>
        <a-button @click="refreshDatabase" :loading="loading" size="small">
          <reload-outlined />
          刷新
        </a-button>
      </template>
      
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic
            title="个人获奖记录"
            :value="databaseStats.person_awards_count"
            prefix="👤"
          />
        </a-col>
        <a-col :span="6">
          <a-statistic
            title="团体获奖记录"
            :value="databaseStats.group_awards_count"
            prefix="👥"
          />
        </a-col>
        <a-col :span="6">
          <a-statistic
            title="团体成员记录"
            :value="databaseStats.group_members_count"
            prefix="👨‍👩‍👧‍👦"
          />
        </a-col>
        <a-col :span="6">
          <a-statistic
            title="数据库大小"
            :value="databaseStats.database_size_mb"
            suffix="MB"
            prefix="💾"
          />
        </a-col>
      </a-row>
      
      <div v-if="databaseStats.table_sizes && databaseStats.table_sizes.length > 0" class="table-sizes">
        <h4>表大小详情</h4>
        <a-table
          :dataSource="databaseStats.table_sizes"
          :columns="tableSizeColumns"
          :pagination="false"
          size="small"
        />
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { 
  ReloadOutlined, 
  DeleteOutlined 
} from '@ant-design/icons-vue';
import { monitorAPI } from '../api/index.js';

// 响应式数据
const loading = ref(false);
const resetting = ref(false);
const autoRefresh = ref(true);
const refreshInterval = ref(null);

const monitorData = ref({
  online_users: [],
  online_count: 0,
  request_stats: {
    total_requests: 0,
    success_requests: 0,
    error_requests: 0,
    avg_response_time: 0,
    qps: 0
  },
  system_metrics: {
    cpu_usage: 0,
    memory_usage: 0,
    memory_total: 0,
    memory_used: 0,
    goroutines: 0,
    db_connections: 0,
    uptime: 0
  }
});

const systemHealth = ref({
  health_score: 100,
  health_status: 'healthy',
  issues: []
});

const requestHistory = ref([]);
const databaseStats = ref({});

// 表格列定义
const userColumns = [
  { title: '用户名', dataIndex: 'username', key: 'username' },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip' },
  { title: '用户类型', dataIndex: 'is_admin', key: 'is_admin' },
  { title: '最后活动', dataIndex: 'last_active', key: 'last_active' },
  { title: '登录时间', dataIndex: 'login_time', key: 'login_time' },
  { title: '操作', key: 'action', width: 80 }
];

const requestColumns = [
  { title: '时间', dataIndex: 'timestamp', key: 'timestamp', width: 150 },
  { title: '方法', dataIndex: 'method', key: 'method', width: 80 },
  { title: '路径', dataIndex: 'path', key: 'path' },
  { title: '状态码', dataIndex: 'status', key: 'status', width: 80 },
  { title: '响应时间', dataIndex: 'response_time', key: 'response_time', width: 100 },
  { title: 'IP', dataIndex: 'ip', key: 'ip', width: 120 }
];

const tableSizeColumns = [
  { title: '表名', dataIndex: 'table_name', key: 'table_name' },
  { title: '大小(MB)', dataIndex: 'Size in MB', key: 'size' }
];

// 计算属性
const healthStatusClass = computed(() => {
  const score = systemHealth.value.health_score;
  if (score >= 90) return 'health-excellent';
  if (score >= 70) return 'health-good';
  if (score >= 50) return 'health-warning';
  return 'health-danger';
});

const healthColor = computed(() => {
  const score = systemHealth.value.health_score;
  if (score >= 90) return '#52c41a';
  if (score >= 70) return '#1890ff';
  if (score >= 50) return '#faad14';
  return '#ff4d4f';
});

const healthStatusText = computed(() => {
  const status = systemHealth.value.health_status;
  const statusMap = {
    'healthy': '健康',
    'warning': '警告',
    'unhealthy': '异常'
  };
  return statusMap[status] || '未知';
});

// 方法
const refreshData = async () => {
  loading.value = true;
  try {
    await Promise.all([
      loadMonitorData(),
      loadSystemHealth(),
      loadRequestHistory(),
      loadDatabaseStats()
    ]);
  } catch (error) {
    console.error('刷新数据失败:', error);
    message.error('刷新数据失败');
  } finally {
    loading.value = false;
  }
};

const loadMonitorData = async () => {
  try {
    const response = await monitorAPI.getMonitorData();
    monitorData.value = response.data;
  } catch (error) {
    console.error('加载监控数据失败:', error);
  }
};

const loadSystemHealth = async () => {
  try {
    const response = await monitorAPI.getSystemHealth();
    systemHealth.value = response.data;
  } catch (error) {
    console.error('加载系统健康状态失败:', error);
  }
};

const loadRequestHistory = async () => {
  try {
    const response = await monitorAPI.getRequestHistory();
    requestHistory.value = response.data.request_history || [];
  } catch (error) {
    console.error('加载请求历史失败:', error);
  }
};

const loadDatabaseStats = async () => {
  try {
    const response = await monitorAPI.getDatabaseStats();
    databaseStats.value = response.data;
  } catch (error) {
    console.error('加载数据库统计失败:', error);
  }
};

const refreshUsers = async () => {
  await loadMonitorData();
};

const refreshRequests = async () => {
  await loadRequestHistory();
};

const refreshDatabase = async () => {
  await loadDatabaseStats();
};

const resetStats = async () => {
  resetting.value = true;
  try {
    await monitorAPI.resetStats();
    message.success('统计数据已重置');
    await refreshData();
  } catch (error) {
    console.error('重置统计失败:', error);
    message.error('重置统计失败');
  } finally {
    resetting.value = false;
  }
};

const kickUser = async (userId) => {
  Modal.confirm({
    title: '确认踢出用户',
    content: '确定要踢出这个用户吗？',
    onOk: async () => {
      try {
        await monitorAPI.kickUser(userId);
        message.success('用户已踢出');
        await refreshUsers();
      } catch (error) {
        console.error('踢出用户失败:', error);
        message.error('踢出用户失败');
      }
    }
  });
};

const toggleAutoRefresh = (checked) => {
  if (checked) {
    startAutoRefresh();
  } else {
    stopAutoRefresh();
  }
};

const startAutoRefresh = () => {
  refreshInterval.value = setInterval(refreshData, 5000); // 每5秒刷新
};

const stopAutoRefresh = () => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value);
    refreshInterval.value = null;
  }
};

// 工具方法
const formatBytes = (bytes) => {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatUptime = (seconds) => {
  if (!seconds) return '0秒';
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);
  
  let result = '';
  if (days > 0) result += `${days}天`;
  if (hours > 0) result += `${hours}小时`;
  if (minutes > 0) result += `${minutes}分钟`;
  if (secs > 0) result += `${secs}秒`;
  
  return result || '0秒';
};

const formatTime = (timeStr) => {
  if (!timeStr) return '-';
  return new Date(timeStr).toLocaleString();
};

const getProgressStatus = (percent) => {
  if (percent >= 90) return 'exception';
  if (percent >= 70) return 'active';
  return 'normal';
};

const getErrorRate = () => {
  const stats = monitorData.value.request_stats;
  if (!stats || stats.total_requests === 0) return 0;
  return ((stats.error_requests / stats.total_requests) * 100).toFixed(2);
};

const getStatusColor = (status) => {
  if (status >= 200 && status < 300) return 'green';
  if (status >= 300 && status < 400) return 'blue';
  if (status >= 400 && status < 500) return 'orange';
  return 'red';
};

// 生命周期
onMounted(() => {
  refreshData();
  if (autoRefresh.value) {
    startAutoRefresh();
  }
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>

<style scoped>
.monitor-container {
  padding: 24px;
  background: #f5f5f5;
  min-height: 100vh;
}

.monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.monitor-header h2 {
  margin: 0;
  color: #1890ff;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.monitor-cards {
  margin-bottom: 24px;
}

.health-card {
  text-align: center;
  border-radius: 8px;
  transition: all 0.3s;
}

.health-card.health-excellent {
  border-color: #52c41a;
  background: linear-gradient(135deg, #f6ffed 0%, #d9f7be 100%);
}

.health-card.health-good {
  border-color: #1890ff;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
}

.health-card.health-warning {
  border-color: #faad14;
  background: linear-gradient(135deg, #fffbe6 0%, #fff1b8 100%);
}

.health-card.health-danger {
  border-color: #ff4d4f;
  background: linear-gradient(135deg, #fff2f0 0%, #ffccc7 100%);
}

.health-status {
  margin-top: 8px;
  font-weight: 500;
  font-size: 14px;
}

.monitor-section {
  margin-bottom: 24px;
}

.metrics-card {
  height: 100%;
}

.metric-item {
  margin-bottom: 24px;
}

.metric-item:last-child {
  margin-bottom: 0;
}

.metric-label {
  font-weight: 500;
  margin-bottom: 8px;
  color: #262626;
}

.metric-detail {
  margin-top: 4px;
  font-size: 12px;
  color: #8c8c8c;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  border-radius: 8px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
}

.stat-item.success {
  background: #f6ffed;
  border-color: #b7eb8f;
}

.stat-item.error {
  background: #fff2f0;
  border-color: #ffccc7;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1890ff;
  margin-bottom: 4px;
}

.stat-item.success .stat-value {
  color: #52c41a;
}

.stat-item.error .stat-value {
  color: #ff4d4f;
}

.stat-label {
  font-size: 12px;
  color: #8c8c8c;
}

.table-sizes {
  margin-top: 16px;
}

.slow-request {
  color: #ff4d4f;
  font-weight: 500;
}

:deep(.ant-card-head) {
  border-bottom: 1px solid #f0f0f0;
}

:deep(.ant-card-body) {
  padding: 16px;
}

:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
}

:deep(.ant-progress-text) {
  font-weight: 500;
}
</style>

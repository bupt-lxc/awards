<template>
  <div>
    <a-row :gutter="[16, 16]">
      <!-- 按年份统计 -->
      <a-col :span="12">
        <a-card title="按年份统计" :loading="loading">
          <a-table 
            :dataSource="statistics.year_statistics" 
            :columns="yearColumns" 
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'count'">
                <a-tag color="blue">{{ record.count }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      
      <!-- 按奖项类型统计 -->
      <a-col :span="12">
        <a-card title="按奖项类型统计" :loading="loading">
          <a-table 
            :dataSource="statistics.type_statistics" 
            :columns="typeColumns" 
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'award_type'">
                <a-tag :color="record.award_type === '团体' ? 'blue' : 'green'">
                  {{ record.award_type }}
                </a-tag>
              </template>
              <template v-if="column.key === 'count'">
                <a-tag color="blue">{{ record.count }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      
      <!-- 按奖项等级统计 -->
      <a-col :span="12">
        <a-card title="按奖项等级统计" :loading="loading">
          <a-table 
            :dataSource="statistics.level_statistics" 
            :columns="levelColumns" 
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'award_level'">
                <a-tag :color="getAwardLevelColor(record.award_level)">
                  {{ record.award_level }}
                </a-tag>
              </template>
              <template v-if="column.key === 'count'">
                <a-tag color="blue">{{ record.count }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      
      <!-- 按组别统计 -->
      <a-col :span="12">
        <a-card title="按组别统计" :loading="loading">
          <a-table 
            :dataSource="statistics.category_statistics" 
            :columns="categoryColumns" 
            :pagination="{ pageSize: 5 }"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'count'">
                <a-tag color="blue">{{ record.count }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { awardAPI } from '../api/index.js';

const loading = ref(false);
const statistics = ref({
  year_statistics: [],
  type_statistics: [],
  level_statistics: [],
  category_statistics: []
});

// 表格列定义
const yearColumns = [
  {
    title: '年份',
    dataIndex: 'year',
    key: 'year'
  },
  {
    title: '获奖数量',
    dataIndex: 'count',
    key: 'count'
  }
];

const typeColumns = [
  {
    title: '奖项类型',
    dataIndex: 'award_type',
    key: 'award_type'
  },
  {
    title: '获奖数量',
    dataIndex: 'count',
    key: 'count'
  }
];

const levelColumns = [
  {
    title: '奖项等级',
    dataIndex: 'award_level',
    key: 'award_level'
  },
  {
    title: '获奖数量',
    dataIndex: 'count',
    key: 'count'
  }
];

const categoryColumns = [
  {
    title: '组别',
    dataIndex: 'category',
    key: 'category'
  },
  {
    title: '获奖数量',
    dataIndex: 'count',
    key: 'count'
  }
];

// 获取统计数据
const fetchStatistics = async () => {
  try {
    loading.value = true;
    
    // 统计功能暂时不可用
    message.info('统计功能正在维护中，暂时不可用');
    
    // 设置空数据
    statistics.value = {
      year_statistics: [],
      type_statistics: [],
      level_statistics: [],
      category_statistics: []
    };
    
  } catch (error) {
    console.error('获取统计数据失败:', error);
    message.error('统计功能暂时不可用');
  } finally {
    loading.value = false;
  }
};

// 获取总记录数
const getTotalCount = () => {
  return statistics.value.type_statistics.reduce((total, item) => total + item.count, 0);
};

// 获取指定类型的数量
const getTypeCount = (type) => {
  const item = statistics.value.type_statistics.find(item => item.award_type === type);
  return item ? item.count : 0;
};

// 获取指定等级的数量
const getLevelCount = (level) => {
  const item = statistics.value.level_statistics.find(item => item.award_level === level);
  return item ? item.count : 0;
};

// 获取奖项等级颜色
const getAwardLevelColor = (level) => {
  const colorMap = {
    '金奖': 'gold',
    '银奖': '#c0c0c0',
    '铜奖': '#cd7f32',
    '特等奖': 'red',
    '一等奖': 'volcano',
    '二等奖': 'orange',
    '三等奖': 'yellow',
    '优秀奖': 'blue'
  };
  return colorMap[level] || 'default';
};

// 组件挂载时获取数据
onMounted(() => {
  fetchStatistics();
});
</script>

<style scoped>
.ant-statistic {
  text-align: center;
}
</style>
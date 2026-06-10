<template>
  <div>
    <!-- 管理员提示 -->
    <a-alert 
      v-if="!userInfo" 
      message="需要管理员权限" 
      description="记录管理功能需要管理员权限，请先登录。" 
      type="warning" 
      show-icon
      style="margin-bottom: 16px;" 
    />
    
    <div v-else>
      <!-- 搜索和筛选 -->
      <a-card title="记录管理" style="margin-bottom: 16px;">
        <!-- 性能信息 -->
        <div style="margin-bottom: 16px; padding: 8px 12px; background: #f6ffed; border: 1px solid #b7eb8f; border-radius: 4px;">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <span style="color: #52c41a; font-weight: 500;">
              🚀 轻量级模式已启用 - 证书图片按需加载
            </span>
            <div style="font-size: 12px; color: #666;">
              记录数: {{ performanceInfo.recordCount }} | 
              内存使用: ~{{ performanceInfo.memoryUsage }}KB |
              <span style="color: #1890ff; cursor: pointer;" @click="clearMemory" title="点击清理内存">
                🧹 清理内存
              </span>
            </div>
          </div>
        </div>
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input-search
              v-model:value="searchText"
              placeholder="搜索团体名称或个人姓名"
              enter-button="搜索"
              @search="onSearch"
              allow-clear
            />
          </a-col>
          <a-col :span="4">
            <a-select
              v-model:value="filterType"
              placeholder="奖项类型"
              style="width: 100%"
              allow-clear
              @change="onFilterChange"
            >
              <a-select-option value="group">团体</a-select-option>
              <a-select-option value="person">个人</a-select-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select
              v-model:value="filterYear"
              placeholder="年份"
              style="width: 100%"
              allow-clear
              @change="onFilterChange"
            >
              <a-select-option v-for="year in availableYears" :key="year" :value="year">
                {{ year }}
              </a-select-option>
            </a-select>
          </a-col>
          <a-col :span="6">
            <a-input
              v-model:value="filterCategory"
              placeholder="组别筛选"
              allow-clear
              @change="onFilterChange"
            />
          </a-col>
          <a-col :span="4">
            <a-button type="primary" @click="loadAllRecords" :loading="loading">
              刷新数据
            </a-button>
            <a-button 
              v-if="certificateLoading" 
              type="default" 
              :loading="certificateLoading"
              style="margin-left: 8px;"
            >
              加载证书中...
            </a-button>
            <a-button 
              type="default" 
              @click="clearMemory"
              style="margin-left: 8px;"
              title="清理已加载的证书图片，释放内存"
            >
              🧹 清理内存
            </a-button>
            <a-button 
              type="default" 
              @click="advancedClearMemory"
              style="margin-left: 8px;"
              title="高级清理：清理所有缓存和临时数据"
            >
              🚀 高级清理
            </a-button>
          </a-col>
        </a-row>
      </a-card>

      <!-- 记录列表 -->
      <a-card title="获奖记录列表">
        <a-tabs v-model:activeKey="activeTab" @change="onTabChange">
          <!-- 个人记录标签页 -->
          <a-tab-pane key="person" tab="个人记录">
            <template #tab>
              <span>
                <user-outlined />
                个人记录 ({{ personRecords.length }})
              </span>
            </template>
            
            <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
              <span>共 {{ personRecords.length }} 条个人记录</span>
              <a-button 
                type="primary" 
                danger 
                @click="batchDelete" 
                :disabled="selectedRowKeys.length === 0"
                v-if="selectedRowKeys.length > 0"
              >
                批量删除 ({{ selectedRowKeys.length }})
              </a-button>
            </div>
            
            <a-table 
              :dataSource="personRecords" 
              :columns="personColumns" 
              :pagination="{ 
                pageSize: 15, 
                showSizeChanger: true, 
                showQuickJumper: true,
                showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条记录`
              }"
              :loading="loading"
              row-key="id"
              :row-selection="{
                selectedRowKeys: selectedRowKeys,
                onChange: onSelectChange,
                getCheckboxProps: (record) => ({
                  name: record.id,
                }),
              }"
              :scroll="{ x: 1000 }"
              @change="handleTableChange"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'award_level'">
                  <a-tag :color="getAwardLevelColor(record.award_level)">
                    {{ record.award_level }}
                  </a-tag>
                </template>
                
                <template v-if="column.key === 'created_at'">
                  {{ formatDate(record.created_at) }}
                </template>
                
                <template v-if="column.key === 'description'">
                  <a-tooltip :title="record.description" v-if="record.description">
                    <span>{{ record.description.length > 20 ? record.description.substring(0, 20) + '...' : record.description }}</span>
                  </a-tooltip>
                  <span v-else style="color: #999;">-</span>
                </template>
                
                <template v-if="column.key === 'action'">
                  <a-space>
                    <a-button 
                      type="link" 
                      size="small" 
                      @click="viewRecord(record)"
                    >
                      查看
                    </a-button>
                    <a-button 
                      type="link" 
                      size="small" 
                      @click="editRecord(record)"
                    >
                      编辑
                    </a-button>
                    <a-button 
                      type="link" 
                      size="small" 
                      @click="generateVerifyLink(record)"
                    >
                      生成链接
                    </a-button>
                    <a-button 
                      type="link" 
                      danger 
                      size="small" 
                      @click="deleteRecord(record)"
                    >
                      删除
                    </a-button>
                  </a-space>
                </template>
              </template>
            </a-table>
          </a-tab-pane>
          
          <!-- 团队记录标签页 -->
          <a-tab-pane key="group" tab="团队记录">
            <template #tab>
              <span>
                <team-outlined />
                团队记录 ({{ groupRecords.length }})
              </span>
            </template>
            
            <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
              <span>共 {{ groupRecords.length }} 条团队记录</span>
              <a-button 
                type="primary" 
                danger 
                @click="batchDelete" 
                :disabled="selectedRowKeys.length === 0"
                v-if="selectedRowKeys.length > 0"
              >
                批量删除 ({{ selectedRowKeys.length }})
              </a-button>
            </div>
            
            <a-table 
              :dataSource="groupRecords" 
              :columns="groupColumns" 
              :pagination="{ 
                pageSize: 15, 
                showSizeChanger: true, 
                showQuickJumper: true,
                showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条记录`
              }"
              :loading="loading"
              row-key="id"
              :row-selection="{
                selectedRowKeys: selectedRowKeys,
                onChange: onSelectChange,
                getCheckboxProps: (record) => ({
                  name: record.id,
                }),
              }"
              :scroll="{ x: 1200 }"
              @change="handleTableChange"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'award_level'">
                  <a-tag :color="getAwardLevelColor(record.award_level)">
                    {{ record.award_level }}
                  </a-tag>
                </template>
                
                <template v-if="column.key === 'members'">
                  <div v-if="record.members && record.members.length > 0" class="members-display">
                    <a-tag 
                      v-for="(member, index) in record.members.slice(0, 3)" 
                      :key="index" 
                      color="blue"
                      style="margin: 2px;"
                    >
                      {{ member.name }}
                    </a-tag>
                    <a-tag v-if="record.members.length > 3" color="default">
                      +{{ record.members.length - 3 }}人
                    </a-tag>
                  </div>
                  <span v-else style="color: #999;">无成员信息</span>
                </template>
                
                <template v-if="column.key === 'created_at'">
                  {{ formatDate(record.created_at) }}
                </template>
                
                <template v-if="column.key === 'description'">
                  <a-tooltip :title="record.description" v-if="record.description">
                    <span>{{ record.description.length > 20 ? record.description.substring(0, 20) + '...' : record.description }}</span>
                  </a-tooltip>
                  <span v-else style="color: #999;">-</span>
                </template>
                
                <template v-if="column.key === 'action'">
                  <a-space>
                    <a-button 
                      type="link" 
                      size="small" 
                      @click="viewRecord(record)"
                    >
                      查看
                    </a-button>
                    <a-button 
                      type="link" 
                      size="small" 
                      @click="editRecord(record)"
                    >
                      编辑
                    </a-button>
                    <a-button 
                      type="link" 
                      size="small" 
                      @click="generateVerifyLink(record)"
                    >
                      生成链接
                    </a-button>
                    <a-button 
                      type="link" 
                      danger 
                      size="small" 
                      @click="deleteRecord(record)"
                    >
                      删除
                    </a-button>
                  </a-space>
                </template>
              </template>
            </a-table>
          </a-tab-pane>
        </a-tabs>
       </a-card>
    </div>
    
    <!-- 查看详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="获奖记录详情"
      :footer="null"
      width="600px"
    >
      <div v-if="selectedRecord">
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="记录ID">{{ selectedRecord.id }}</a-descriptions-item>
          <a-descriptions-item label="奖项类型">
            <a-tag :color="selectedRecord.award_type === 'group' ? 'blue' : 'green'">
              {{ selectedRecord.award_type === 'group' ? '团体' : '个人' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="年份">{{ selectedRecord.year }}</a-descriptions-item>
          <a-descriptions-item label="团体名称" v-if="selectedRecord.award_type === 'group'">
            {{ selectedRecord.group_name || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="个人姓名" v-if="selectedRecord.award_type === 'person'">
            {{ selectedRecord.person_name || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="证件号码" v-if="selectedRecord.award_type === 'person'">
            {{ selectedRecord.id_number || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="团体成员" v-if="selectedRecord.award_type === 'group' && selectedRecord.members" :span="2">
            <div v-if="selectedRecord.members && selectedRecord.members.length > 0">
              <a-tag v-for="(member, index) in selectedRecord.members" :key="index" style="margin: 2px;">
                {{ member.name }} ({{ member.id_number }})
              </a-tag>
            </div>
            <span v-else>无成员信息</span>
          </a-descriptions-item>
          <a-descriptions-item label="组别">{{ selectedRecord.category }}</a-descriptions-item>
          <a-descriptions-item label="奖项等级">
            <a-tag :color="getAwardLevelColor(selectedRecord.award_level)">
              {{ selectedRecord.award_level }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间" :span="2">
            {{ formatDate(selectedRecord.created_at) }}
          </a-descriptions-item>
          <a-descriptions-item label="更新时间" :span="2" v-if="selectedRecord.updated_at">
            {{ formatDate(selectedRecord.updated_at) }}
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ selectedRecord.description || '无描述' }}
          </a-descriptions-item>
          
          <!-- 个人获奖证书图片 -->
          <a-descriptions-item label="获奖证书" :span="2" v-if="selectedRecord.award_type === 'person'">
            <div style="text-align: center;">
              <template v-if="selectedRecord.certificate_image">
                <div v-if="selectedRecord.certificate_image === 'has_certificate'" style="padding: 20px; background: #f5f5f5; border-radius: 6px;">
                  <a-spin size="large" />
                  <div style="margin-top: 8px; color: #666;">正在加载证书图片...</div>
                </div>
                <img 
                  v-else
                  :src="processImageUrl(selectedRecord.certificate_image)" 
                  alt="获奖证书" 
                  style="max-width: 400px; max-height: 300px; border: 1px solid #d9d9d9; border-radius: 6px; cursor: pointer;"
                  @click="previewImage(selectedRecord.certificate_image, selectedRecord.person_name + ' - 获奖证书')"
                  @error="handleImageError($event, '个人证书')"
                />
              </template>
              <div v-else style="color: #999; font-style: italic;">暂无证书图片</div>
            </div>
          </a-descriptions-item>
          
          <!-- 团体成员证书图片 -->
          <a-descriptions-item label="成员证书" :span="2" v-if="selectedRecord.award_type === 'group' && selectedRecord.members">
            <div v-if="selectedRecord.members.some(m => m.certificate_image)" style="display: flex; flex-wrap: wrap; gap: 10px;">
              <div 
                v-for="(member, index) in selectedRecord.members.filter(m => m.certificate_image)" 
                :key="index"
                style="text-align: center;"
              >
                <div style="margin-bottom: 5px; font-size: 12px; color: #666;">{{ member.name }}</div>
                <div v-if="member.certificate_image === 'has_certificate'" style="width: 120px; height: 90px; background: #f5f5f5; border: 1px solid #d9d9d9; border-radius: 4px; display: flex; align-items: center; justify-content: center;">
                  <a-spin size="small" />
                </div>
                <img 
                  v-else
                  :src="processImageUrl(member.certificate_image)" 
                  :alt="member.name + ' 证书'" 
                  style="width: 120px; height: 90px; object-fit: cover; border: 1px solid #d9d9d9; border-radius: 4px; cursor: pointer;"
                  @click="previewImage(member.certificate_image, member.name + ' - 获奖证书')"
                  @error="handleImageError($event, member.name + '的证书')"
                />
              </div>
            </div>
            <div v-else style="color: #999; font-style: italic;">暂无成员证书图片</div>
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-modal>
    
    <!-- 编辑记录模态框 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑获奖记录"
      :footer="null"
      width="600px"
    >
      <a-form
        v-if="editForm"
        :model="editForm"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
        :rules="editRules"
        ref="editFormRef"
        @finish="saveEdit"
      >
        <a-form-item label="奖项类型" name="award_type">
          <a-radio-group v-model:value="editForm.award_type" button-style="solid">
            <a-radio-button value="group">团体</a-radio-button>
            <a-radio-button value="person">个人</a-radio-button>
          </a-radio-group>
        </a-form-item>
        
        <a-form-item label="年份" name="year">
          <a-input-number v-model:value="editForm.year" :min="2000" :max="2030" style="width: 120px" />
        </a-form-item>
        
        <a-form-item label="团体名称" name="group_name" v-if="editForm.award_type === 'group'">
          <a-input v-model:value="editForm.group_name" placeholder="请输入团体名称" />
        </a-form-item>
        
        <a-form-item label="个人姓名" name="person_name" v-if="editForm.award_type === 'person'">
          <a-input v-model:value="editForm.person_name" placeholder="请输入个人姓名" />
        </a-form-item>
        
        <a-form-item label="证件号码" name="id_number" v-if="editForm.award_type === 'person'">
          <a-input v-model:value="editForm.id_number" placeholder="请输入证件号码" />
        </a-form-item>
        
        <!-- 个人证书图片上传 -->
        <a-form-item label="获奖证书" name="certificate_image" v-if="editForm.award_type === 'person'">
          <a-upload
            v-model:file-list="personCertificateFileList"
            name="certificate"
            list-type="picture-card"
            class="certificate-uploader"
            :show-upload-list="false"
            :before-upload="beforeUpload"
            @change="handlePersonCertificateChange"
          >
            <div v-if="editForm.certificate_image">
              <img :src="editForm.certificate_image" alt="证书" style="width: 100%; height: 100%; object-fit: cover;" />
            </div>
            <div v-else>
              <plus-outlined />
              <div style="margin-top: 8px">上传证书</div>
            </div>
          </a-upload>
        </a-form-item>
        
        <!-- 团体成员编辑 -->
        <a-form-item label="团体成员" v-if="editForm.award_type === 'group'" :wrapper-col="{ span: 18 }">
          <div style="margin-bottom: 8px;">
            <a-input 
              v-model:value="editNewMember.name" 
              placeholder="成员姓名" 
              style="width: 120px; margin-right: 8px;"
            />
            <a-input 
              v-model:value="editNewMember.id_number" 
              placeholder="证件号码" 
              style="width: 150px; margin-right: 8px;"
            />
            <a-button type="primary" size="small" @click="addEditMember">添加成员</a-button>
          </div>
          <div v-if="editForm.members && editForm.members.length > 0" class="members-edit-list">
            <div 
              v-for="(member, index) in editForm.members" 
              :key="index" 
              class="member-edit-card"
            >
              <div class="member-edit-info">
                <div class="member-edit-text">
                  <span class="member-name">{{ member.name }}</span>
                  <span class="member-id">({{ member.id_number }})</span>
                </div>
                <a-button 
                  type="text" 
                  danger 
                  size="small" 
                  @click="removeEditMember(index)"
                  title="删除成员"
                >
                  <delete-outlined />
                </a-button>
              </div>
              
              <!-- 成员证书上传 -->
              <div class="member-certificate-upload">
                <div class="certificate-upload-label">证书图片：</div>
                <a-upload
                  :file-list="[]"
                  name="certificate"
                  list-type="picture-card"
                  class="member-certificate-uploader"
                  :show-upload-list="false"
                  :before-upload="beforeUpload"
                  @change="(info) => handleMemberCertificateChange(info, index)"
                >
                  <div v-if="member.certificate_image">
                    <img :src="member.certificate_image" alt="成员证书" style="width: 100%; height: 100%; object-fit: cover;" />
                  </div>
                  <div v-else class="upload-placeholder">
                    <plus-outlined />
                    <div style="margin-top: 4px; font-size: 12px;">上传证书</div>
                  </div>
                </a-upload>
              </div>
            </div>
          </div>
          <div v-else style="color: #999; font-size: 12px;">暂无成员</div>
        </a-form-item>
        
        <a-form-item label="组别" name="category">
          <a-input v-model:value="editForm.category" placeholder="请输入组别" />
        </a-form-item>
        
        <a-form-item label="奖项等级" name="award_level">
          <a-input v-model:value="editForm.award_level" placeholder="请输入奖项等级" />
        </a-form-item>
        
        <a-form-item label="描述" name="description">
          <a-textarea v-model:value="editForm.description" placeholder="请输入获奖描述" :rows="3" />
        </a-form-item>
        
        <a-form-item :wrapper-col="{ offset: 6, span: 18 }">
          <a-space>
            <a-button type="primary" html-type="submit" :loading="editLoading">
              保存修改
            </a-button>
            <a-button @click="cancelEdit">
              取消
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 新增成员二维码弹窗 -->
    <a-modal v-model:open="memberQrModalVisible" title="选择成员生成专属二维码" :footer="null" width="500px">
      <div v-if="memberQrList.length > 0">
        <a-list :dataSource="memberQrList" bordered>
          <template #renderItem="{ item }">
            <a-list-item>
              <span>{{ item.name }}（{{ item.id_number }}）</span>
              <a-button type="primary" size="small" style="margin-left: 16px;" @click="showMemberQr(item)">生成二维码</a-button>
            </a-list-item>
          </template>
        </a-list>
      </div>
      <div v-else style="text-align: center; color: #999;">暂无成员信息</div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { awardAPI } from '../api/index.js';
import { ExclamationCircleOutlined, PlusOutlined, DeleteOutlined, SearchOutlined, UserOutlined, TeamOutlined } from '@ant-design/icons-vue';
import { createVNode } from 'vue';
import QRCode from 'qrcode';

const loading = ref(false);
const allRecords = ref([]);
const userInfo = ref(null);
const searchText = ref('');
const filterType = ref(undefined);
const filterYear = ref(undefined);
const filterCategory = ref('');
const selectedRowKeys = ref([]);
const detailModalVisible = ref(false);
const selectedRecord = ref(null);
const activeTab = ref('person');
const editModalVisible = ref(false);
const editForm = ref(null);
const editLoading = ref(false);
const editFormRef = ref(null);
const editNewMember = ref({ name: '', id_number: '' });
const personCertificateFileList = ref([]);
const memberQrModalVisible = ref(false);
const memberQrList = ref([]);

// 添加新的响应式变量来存储表格的筛选状态
const tableFilters = ref({
  award_type: null,
  year: null,
  group_name: null,
  category: null,
  award_level: null
});

// 编辑表单验证规则
const editRules = {
  award_type: [{ required: true, message: '请选择奖项类型' }],
  year: [{ required: true, message: '请输入年份' }],
  group_name: [
    {
      validator: (rule, value) => {
        if (editForm.value?.award_type === 'group' && !value) {
          return Promise.reject('请输入团体名称');
        }
        return Promise.resolve();
      }
    }
  ],
  person_name: [
    {
      validator: (rule, value) => {
        if (editForm.value?.award_type === 'person' && !value) {
          return Promise.reject('请输入个人姓名');
        }
        return Promise.resolve();
      }
    }
  ],
  id_number: [
    {
      validator: (rule, value) => {
        if (editForm.value?.award_type === 'person' && !value) {
          return Promise.reject('请输入证件号码');
        }
        return Promise.resolve();
      }
    }
  ],
  category: [{ required: true, message: '请输入组别' }],
  award_level: [{ required: true, message: '请选择奖项等级' }]
};

// 计算属性 - 分离个人和团队记录
const personRecords = computed(() => {
  return filteredRecords.value.filter(record => !record.group_name || record.group_name.trim() === '')
})

const groupRecords = computed(() => {
  return filteredRecords.value.filter(record => record.group_name && record.group_name.trim() !== '')
})

// 个人记录列定义
const personColumns = computed(() => {
  const uniqueCategories = Array.from(new Set(personRecords.value
    .filter(r => r.category)
    .map(r => r.category)));

  const uniqueAwardLevels = Array.from(new Set(personRecords.value
    .filter(r => r.award_level)
    .map(r => r.award_level)));

  return [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80
    },
    {
      title: '个人姓名',
      dataIndex: 'person_name',
      key: 'person_name',
      width: 120,
      ellipsis: true
    },
    {
      title: '证件号码',
      dataIndex: 'id_number',
      key: 'id_number',
      width: 150
    },
    {
      title: '年份',
      dataIndex: 'year',
      key: 'year',
      width: 80
    },
    {
      title: '组别',
      dataIndex: 'category',
      key: 'category',
      width: 120,
      ellipsis: true,
      filters: uniqueCategories.map(cat => ({ text: cat, value: cat }))
    },
    {
      title: '奖项等级',
      dataIndex: 'award_level',
      key: 'award_level',
      width: 100,
      filters: uniqueAwardLevels.map(level => ({ text: level, value: level }))
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 150
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      width: 200,
      ellipsis: true
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right'
    }
  ];
})

// 团队记录列定义
const groupColumns = computed(() => {
  const uniqueGroupNames = Array.from(new Set(groupRecords.value
    .filter(r => r.group_name)
    .map(r => r.group_name)));

  const uniqueCategories = Array.from(new Set(groupRecords.value
    .filter(r => r.category)
    .map(r => r.category)));

  const uniqueAwardLevels = Array.from(new Set(groupRecords.value
    .filter(r => r.award_level)
    .map(r => r.award_level)));

  return [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80
    },
    {
      title: '团体名称',
      dataIndex: 'group_name',
      key: 'group_name',
      width: 200,
      ellipsis: true,
      filters: uniqueGroupNames.map(name => ({ text: name, value: name }))
    },
    {
      title: '团队成员',
      key: 'members',
      width: 250
    },
    {
      title: '年份',
      dataIndex: 'year',
      key: 'year',
      width: 80
    },
    {
      title: '组别',
      dataIndex: 'category',
      key: 'category',
      width: 120,
      ellipsis: true,
      filters: uniqueCategories.map(cat => ({ text: cat, value: cat }))
    },
    {
      title: '奖项等级',
      dataIndex: 'award_level',
      key: 'award_level',
      width: 100,
      filters: uniqueAwardLevels.map(level => ({ text: level, value: level }))
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 150
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      width: 200,
      ellipsis: true
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right'
    }
  ];
})

// 计算可用年份
const availableYears = computed(() => {
  const years = [...new Set(allRecords.value.map(record => record.year))];
  return years.sort((a, b) => b - a);
});

// 过滤后的记录
const filteredRecords = computed(() => {
  if (!allRecords.value) return [];
  
  let filtered = [...allRecords.value];
  
  // 搜索文本过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase();
    filtered = filtered.filter(record => 
      (record.group_name && record.group_name.toLowerCase().includes(search)) ||
      (record.person_name && record.person_name.toLowerCase().includes(search))
    );
  }
  
  // 顶部筛选条件
  if (filterType.value) {
    filtered = filtered.filter(record => record.award_type === filterType.value);
  }
  
  if (filterYear.value) {
    filtered = filtered.filter(record => record.year === filterYear.value);
  }
  
  if (filterCategory.value) {
    filtered = filtered.filter(record => 
      record.category && record.category.toLowerCase().includes(filterCategory.value.toLowerCase())
    );
  }

  // 表格列筛选
  if (tableFilters.value.award_type) {
    filtered = filtered.filter(record => record.award_type === tableFilters.value.award_type);
  }
  
  if (tableFilters.value.year) {
    filtered = filtered.filter(record => record.year === tableFilters.value.year);
  }
  
  if (tableFilters.value.group_name) {
    filtered = filtered.filter(record => record.group_name === tableFilters.value.group_name);
  }
  
  if (tableFilters.value.category) {
    filtered = filtered.filter(record => record.category === tableFilters.value.category);
  }
  
  if (tableFilters.value.award_level) {
    filtered = filtered.filter(record => record.award_level === tableFilters.value.award_level);
  }
  
  return filtered;
});

// 检查用户登录状态
const checkUserInfo = () => {
  const savedUserInfo = localStorage.getItem('user_info');
  if (savedUserInfo) {
    try {
      userInfo.value = JSON.parse(savedUserInfo);
    } catch (error) {
      console.error('解析用户信息失败:', error);
      userInfo.value = null;
    }
  }
};

// 加载所有记录
const loadAllRecords = async () => {
  try {
    loading.value = true;
    // 添加时间戳参数避免缓存
    const timestamp = new Date().getTime();
    
    // 使用轻量级接口，避免加载证书图片数据
    const response = await awardAPI.getAllAwardsLight(timestamp);
    
    // 合并个人和团体获奖记录
    const allAwards = [];
    
    // 添加个人获奖记录，并标记类型
    if (response.person_awards && response.person_awards.length > 0) {
      response.person_awards.forEach(award => {
        if (!award.id) {
          console.warn('警告：发现无效的个人记录', award);
          return;
        }
        allAwards.push({
          ...award,
          award_type: 'person',
          // 轻量级模式下不包含证书图片，只保留标记
          certificate_image: award.has_certificate ? 'has_certificate' : ''
        });
      });
    }
    
    // 添加团体获奖记录，并标记类型
    if (response.group_awards && response.group_awards.length > 0) {
      response.group_awards.forEach(award => {
        if (!award.id) {
          console.warn('警告：发现无效的团体记录', award);
          return;
        }
        
        // 处理成员数据（轻量级版本）
        let members = award.members || [];
        if (Array.isArray(members)) {
          members = members.map(member => ({
            ...member,
            // 轻量级模式下不包含证书图片，只保留标记
            certificate_image: member.has_certificate ? 'has_certificate' : ''
          }));
        }
        
        allAwards.push({
          ...award,
          members: members,
          award_type: 'group',
          // 轻量级模式下不包含证书图片，只保留标记
          certificate_image: award.has_certificate ? 'has_certificate' : ''
        });
      });
    }
    
    // 数据完整性检查
    const validAwards = allAwards.filter(award => {
      if (!award.id || !award.year) {
        console.warn('警告：发现无效的记录', award);
        return false;
      }
      return true;
    });
    
    allRecords.value = validAwards;
    console.log(`加载记录完成：总计 ${validAwards.length} 条有效记录（轻量级模式）`);
    message.success(`成功加载 ${validAwards.length} 条记录（个人：${response.total_person || 0}，团体：${response.total_group || 0}）`);
    
    // 更新性能信息
    updatePerformanceInfo();
  } catch (error) {
    console.error('加载记录失败:', error);
    if (error.response?.data?.error) {
      message.error(`加载记录失败: ${error.response.data.error}`);
    } else {
      message.error('加载记录失败，请检查网络连接');
    }
  } finally {
    loading.value = false;
  }
};

// 搜索
const onSearch = () => {
  // 搜索逻辑已在计算属性中处理
};

// 筛选变化
const onFilterChange = () => {
  // 筛选逻辑已在计算属性中处理
};

// 选择变化
const onSelectChange = (selectedKeys) => {
  selectedRowKeys.value = selectedKeys;
};

// 查看记录详情
const viewRecord = async (record) => {
  try {
    certificateLoading.value = true;
    
    // 如果是轻量级模式且需要证书图片，先加载证书
    if (record.certificate_image === 'has_certificate') {
      await loadCertificateImage(record);
    }
    
    // 如果是团体记录，加载成员证书
    if (record.award_type === 'group' && record.members) {
      await loadGroupMemberCertificates(record);
    }
    
    selectedRecord.value = record;
    detailModalVisible.value = true;
  } catch (error) {
    console.error('查看记录详情失败:', error);
    message.error('查看记录详情失败');
  } finally {
    certificateLoading.value = false;
  }
};

// 按需加载证书图片
const loadCertificateImage = async (record) => {
  try {
    // 如果已经有证书图片数据，直接返回
    if (record.certificate_image && record.certificate_image !== 'has_certificate') {
      return record.certificate_image;
    }
    
    // 如果没有证书标记，返回空
    if (!record.has_certificate) {
      return null;
    }
    
    // 从服务器获取证书图片
    const response = await awardAPI.getAwardDetail(record.id);
    
    // 根据返回的数据结构获取证书图片
    let certificateImage = null;
    if (response.type === 'person' && response.award) {
      certificateImage = response.award.certificate_image;
    } else if (response.type === 'group' && response.award) {
      certificateImage = response.award.certificate_image;
    }
    
    if (certificateImage) {
      // 更新记录中的证书图片
      record.certificate_image = certificateImage;
      return certificateImage;
    }
    
    return null;
  } catch (error) {
    console.error('加载证书图片失败:', error);
    message.error('加载证书图片失败');
    return null;
  }
};

// 按需加载团体成员证书图片
const loadGroupMemberCertificates = async (record) => {
  try {
    if (!record.members || record.award_type !== 'group') {
      return;
    }
    
    // 检查是否有需要加载证书的成员
    const membersNeedingCertificates = record.members.filter(member => 
      member.certificate_image === 'has_certificate'
    );
    
    if (membersNeedingCertificates.length === 0) {
      return;
    }
    
    // 获取团体成员的证书信息
    const response = await awardAPI.getGroupMembersWithCertificates(record.id);
    
    // 更新成员的证书图片
    if (response && Array.isArray(response)) {
      response.forEach(memberData => {
        const member = record.members.find(m => 
          m.name === memberData.member_name && m.id_number === memberData.member_id_number
        );
        if (member && memberData.certificate_image) {
          member.certificate_image = memberData.certificate_image;
        }
      });
    }
  } catch (error) {
    console.error('加载团体成员证书失败:', error);
    message.error('加载团体成员证书失败');
  }
};

// 编辑记录
const editRecord = (record) => {
  editForm.value = {
    id: record.id,
    award_type: record.award_type,
    year: record.year,
    group_name: record.group_name || '',
    person_name: record.person_name || '',
    id_number: record.id_number || '',
    category: record.category,
    award_level: record.award_level,
    description: record.description || '',
    certificate_image: record.certificate_image || '', // 添加这一行
    members: record.members ? [...record.members] : []
  };
  
  // 重置新成员输入
  editNewMember.value = { name: '', id_number: '' };
  
  editModalVisible.value = true;
};

// 验证图片路径（现为静态文件路径如 /static/certificates/uuid.png）
const validateImageData = (imageData) => {
  if (!imageData) return true;  // 允许为空
  return imageData.startsWith('/static/certificates/');
};

// 保存编辑
const saveEdit = async () => {
  try {
    editLoading.value = true;
    
    // 验证表单数据
    await editFormRef.value.validate();
    
    // 构建更新数据
    const updateData = {
      id: editForm.value.id,
      award_type: editForm.value.award_type,
      year: editForm.value.year,
      category: editForm.value.category,
      award_level: editForm.value.award_level,
      description: editForm.value.description
    };
    
    // 根据奖项类型添加相应字段并调用对应API
    if (editForm.value.award_type === 'group') {
      // 验证团体数据
      if (!editForm.value.group_name) {
        message.error('请输入团体名称');
        return;
      }
      
      updateData.group_name = editForm.value.group_name;
      
      // 验证并处理成员数据
      if (!editForm.value.members || editForm.value.members.length === 0) {
        message.warning('当前没有添加任何成员');
        return;
      }
      
      // 转换成员数据格式并验证
      updateData.members = [];
      for (const member of editForm.value.members) {
        if (!member.name || !member.id_number) {
          message.error(`成员信息不完整：${member.name || '未知成员'}`);
          return;
        }
        
        // 验证证书图片格式
        if (member.certificate_image && !validateImageData(member.certificate_image)) {
          message.error(`成员 ${member.name} 的证书图片格式无效`);
          return;
        }
        
        updateData.members.push({
          name: member.name,
          id_number: String(member.id_number),
          certificate_image: member.certificate_image || ''
        });
      }
      
      // 调用更新API
      await awardAPI.updateGroupAward(updateData);
      
    } else {
      // 验证个人数据
      if (!editForm.value.person_name || !editForm.value.id_number) {
        message.error('请填写完整的个人信息');
        return;
      }
      
      // 验证证书图片格式
      if (editForm.value.certificate_image && !validateImageData(editForm.value.certificate_image)) {
        message.error('证书图片格式无效');
        return;
      }
      
      updateData.person_name = editForm.value.person_name;
      updateData.id_number = String(editForm.value.id_number);
      updateData.certificate_image = editForm.value.certificate_image || '';
      
      // 调用更新API
      await awardAPI.updatePersonAward(updateData);
    }
    
    message.success('修改成功');
    
    // 重新加载所有记录以获取最新数据
    await loadAllRecords();
    
    // 关闭编辑模态框
    editModalVisible.value = false;
    editForm.value = null;
    editNewMember.value = { name: '', id_number: '' };
    personCertificateFileList.value = [];
    
  } catch (error) {
    console.error('修改失败:', error);
    if (error.response?.data?.detail) {
      message.error(`修改失败: ${error.response.data.detail}`);
    } else if (error.errors) {
      // 表单验证错误
      const errorMessages = Object.values(error.errors).flat();
      message.error(`请检查表单: ${errorMessages.join(', ')}`);
    } else {
      message.error('修改失败，请检查网络连接或稍后重试');
    }
  } finally {
    editLoading.value = false;
  }
};

// 添加编辑成员
const addEditMember = () => {
  if (editNewMember.value.name && editNewMember.value.id_number) {
    // 检查是否已存在相同的成员
    const exists = editForm.value.members.some(member => 
      member.name === editNewMember.value.name && member.id_number === editNewMember.value.id_number
    );
    
    if (exists) {
      message.warning('该成员已存在');
      return;
    }
    
    editForm.value.members.push({
      name: editNewMember.value.name,
      id_number: String(editNewMember.value.id_number)
    });
    
    // 清空输入框
    editNewMember.value = { name: '', id_number: '' };
    
    message.success('成员添加成功');
  } else {
    message.warning('请输入成员姓名和证件号码');
  }
};

// 移除编辑成员
const removeEditMember = (index) => {
  editForm.value.members.splice(index, 1);
  message.success('成员移除成功');
};

// 取消编辑
const cancelEdit = () => {
  editModalVisible.value = false;
  editForm.value = null;
  editNewMember.value = { name: '', id_number: '' };
  personCertificateFileList.value = [];
  if (editFormRef.value) {
    editFormRef.value.resetFields();
  }
};

// 图片上传前的检查
const beforeUpload = (file) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png' || file.type === 'image/jpg';
  if (!isJpgOrPng) {
    message.error('只能上传 JPG/PNG 格式的图片!');
    return false;
  }
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    message.error('图片大小不能超过 2MB!');
    return false;
  }
  return false; // 阻止自动上传，手动处理
};

// 图片上传到服务器获取静态文件路径
const uploadImageToServer = async (file) => {
  const response = await awardAPI.uploadCertificate(file, 'edit');
  return response.data.file_path;
};

// 个人证书图片上传处理
const handlePersonCertificateChange = async (info) => {
  try {
    if (info.file) {
      loading.value = true;
      const filePath = await uploadImageToServer(info.file.originFileObj || info.file);
      editForm.value.certificate_image = filePath;
      message.success('证书图片上传成功');
    }
  } catch (error) {
    console.error('证书图片处理失败:', error);
    message.error('证书图片处理失败，请重试');
  } finally {
    loading.value = false;
  }
};

// 团体成员证书图片上传处理
const handleMemberCertificateChange = async (info, memberIndex) => {
  try {
    if (info.file && editForm.value.members[memberIndex]) {
      loading.value = true;
      const filePath = await uploadImageToServer(info.file.originFileObj || info.file);
      editForm.value.members[memberIndex].certificate_image = filePath;
      message.success(`成员 ${editForm.value.members[memberIndex].name} 的证书图片上传成功`);
    }
  } catch (error) {
    console.error('成员证书图片处理失败:', error);
    message.error('成员证书图片处理失败，请重试');
  } finally {
    loading.value = false;
  }
};

// 删除单个记录
const deleteRecord = (record) => {
  if (!record || !record.id || !record.award_type) {
    message.error('无效的记录数据，无法删除');
    return;
  }

  Modal.confirm({
    title: '确认删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode('div', {}, [
      createVNode('p', { style: 'margin-bottom: 16px; color: #ff4d4f; font-weight: bold;' }, '确定要删除这条获奖记录吗？此操作不可恢复！'),
      createVNode('div', { style: 'background: #f5f5f5; padding: 12px; border-radius: 4px;' }, [
        createVNode('p', {}, `类型：${record.award_type === 'group' ? '团体' : '个人'}`),
        createVNode('p', {}, `年份：${record.year}`),
        createVNode('p', {}, `${record.award_type === 'group' ? '团体名称' : '个人姓名'}：${record.award_type === 'group' ? record.group_name : record.person_name}`),
        createVNode('p', {}, `组别：${record.category}`),
        createVNode('p', { style: 'margin-bottom: 0;' }, `奖项等级：${record.award_level}`)
      ])
    ]),
    okText: '确定删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        loading.value = true;
        
        // 添加重试逻辑
        let retryCount = 0;
        const maxRetries = 3;
        
        while (retryCount < maxRetries) {
          try {
            await awardAPI.deleteAward(record.id, record.award_type);
            break; // 如果成功，跳出重试循环
          } catch (error) {
            retryCount++;
            if (retryCount === maxRetries) {
              throw error; // 最后一次重试失败，抛出错误
            }
            await new Promise(resolve => setTimeout(resolve, 1000 * retryCount)); // 延迟重试
          }
        }
        
        message.success('删除成功');
        
        // 重新加载数据以确保数据同步
        await loadAllRecords();
        
        // 清除选中状态
        selectedRowKeys.value = selectedRowKeys.value.filter(key => key !== record.id);
        
      } catch (error) {
        console.error('删除失败:', error);
        if (error.response?.data?.detail) {
          message.error(`删除失败: ${error.response.data.detail}`);
        } else {
          message.error('删除失败，请检查网络连接或稍后重试');
        }
      } finally {
        loading.value = false;
      }
    }
  });
};

// 批量删除
const batchDelete = () => {
  // 验证选中的记录
  const selectedRecords = allRecords.value.filter(record => selectedRowKeys.value.includes(record.id));
  const invalidRecords = selectedRecords.filter(record => !record.id || !record.award_type);
  
  if (invalidRecords.length > 0) {
    message.error('选中的记录中包含无效数据，请重新选择');
    return;
  }

  Modal.confirm({
    title: '确认批量删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode('div', {}, [
      createVNode('p', { style: 'margin-bottom: 16px; color: #ff4d4f; font-weight: bold;' }, 
        `确定要删除选中的 ${selectedRowKeys.value.length} 条获奖记录吗？此操作不可恢复！`
      ),
      createVNode('div', { style: 'background: #f5f5f5; padding: 12px; border-radius: 4px; max-height: 200px; overflow-y: auto;' },
        selectedRecords.map(record => 
          createVNode('p', { style: 'margin-bottom: 8px;' }, 
            `${record.award_type === 'group' ? '团体' : '个人'} - ${record.award_type === 'group' ? record.group_name : record.person_name} (${record.year}年 ${record.category})`
          )
        )
      )
    ]),
    okText: '确定删除',
    okType: 'danger',
    cancelText: '取消',
    width: 500,
    async onOk() {
      try {
        loading.value = true;
        
        // 逐个删除选中的记录，添加重试逻辑
        let successCount = 0;
        let failedCount = 0;
        const maxRetries = 3;
        
        for (const record of selectedRecords) {
          let retryCount = 0;
          let success = false;
          
          while (retryCount < maxRetries && !success) {
            try {
              await awardAPI.deleteAward(record.id, record.award_type);
              success = true;
              successCount++;
            } catch (error) {
              retryCount++;
              if (retryCount === maxRetries) {
                console.error(`删除记录失败 (ID: ${record.id}):`, error);
                failedCount++;
              } else {
                await new Promise(resolve => setTimeout(resolve, 1000 * retryCount));
              }
            }
          }
        }
        
        // 重新加载数据以确保数据同步
        await loadAllRecords();
        
        // 清空选中状态
        selectedRowKeys.value = [];
        
        // 显示结果消息
        if (successCount === selectedRecords.length) {
          message.success(`成功删除所有 ${successCount} 条记录`);
        } else if (successCount > 0) {
          message.warning(`部分删除成功：${successCount} 条成功，${failedCount} 条失败`);
        } else {
          message.error(`删除失败：所有 ${failedCount} 条记录均删除失败`);
        }
      } catch (error) {
        console.error('批量删除失败:', error);
        if (error.response?.data?.detail) {
          message.error(`批量删除失败: ${error.response.data.detail}`);
        } else {
          message.error('批量删除失败，请检查网络连接或稍后重试');
        }
      } finally {
        loading.value = false;
      }
    }
  });
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

// 生成验证链接
const generateVerifyLink = async (record) => {
  if (record.award_type === 'group' && record.members && record.members.length > 0) {
    // 团体奖项，弹窗列出所有成员
    memberQrList.value = record.members.map(member => ({
      ...member,
      group_name: record.group_name,
      year: record.year,
      category: record.category,
      award_id: record.id
    }));
    memberQrModalVisible.value = true;
    return;
  }
  // 个人奖项逻辑保持不变
  try {
    const baseUrl = window.location.origin;
    const params = new URLSearchParams();
    params.append('type', record.award_type);
    params.append('year', record.year);
    params.append('category', record.category);
    params.append('award_level', record.award_level);
    params.append('award_id', record.id);
    if (record.award_type === 'person') {
      params.append('person_name', record.person_name || '');
      params.append('id_number', record.id_number || '');
    }
    const verifyUrl = `${baseUrl}/verify.html?${params.toString()}`;
    const qrCodeDataUrl = await QRCode.toDataURL(verifyUrl, {
      width: 200,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    });
    Modal.info({
      title: '验证链接和二维码',
      content: createVNode('div', { style: 'text-align: center;' }, [
        createVNode('div', { style: 'margin-bottom: 16px;' }, [
          createVNode('p', { style: 'margin-bottom: 8px; font-weight: bold;' }, '扫描二维码或复制链接进行验证：'),
          createVNode('img', { 
            src: qrCodeDataUrl, 
            alt: '验证二维码',
            style: 'display: block; margin: 0 auto 16px auto; border: 1px solid #d9d9d9; border-radius: 6px;'
          }),
          createVNode('div', { style: 'background: #f5f5f5; padding: 12px; border-radius: 6px; word-break: break-all; font-family: monospace; font-size: 12px;' }, verifyUrl),
          createVNode('div', { style: 'margin-top: 12px;' }, [
            createVNode('button', {
              style: 'background: #1890ff; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer;',
              onClick: () => {
                navigator.clipboard.writeText(verifyUrl).then(() => {
                  message.success('验证链接已复制到剪贴板');
                }).catch(() => {
                  message.error('复制失败，请手动复制链接');
                });
              }
            }, '复制链接')
          ])
        ])
      ]),
      width: 400,
      okText: '关闭'
    });
  } catch (error) {
    console.error('生成验证链接失败:', error);
    message.error('生成验证链接失败');
  }
};

// 处理图片URL
const processImageUrl = (imageUrl) => {
  if (!imageUrl) return '';
  
  // 处理轻量级模式下的证书标记
  if (imageUrl === 'has_certificate') {
    return null; // 返回null表示需要按需加载
  }
  
  // 如果是完整URL或静态路径，直接返回
  if (imageUrl.startsWith('data:image/') || imageUrl.startsWith('http')) {
    return imageUrl;
  }
  // 如果是相对路径，添加基础URL
  if (imageUrl.startsWith('/')) {
    return `${window.location.origin}${imageUrl}`;
  }
  return imageUrl;
};

// 图片预览功能
const previewImage = (imageUrl, title) => {
  if (!imageUrl) {
    message.warning('没有可预览的图片');
    return;
  }

  const processedUrl = processImageUrl(imageUrl);
  
  Modal.info({
    title: title || '图片预览',
    content: createVNode('div', { style: 'text-align: center;' }, [
      createVNode('img', {
        src: processedUrl,
        alt: '预览图片',
        style: 'max-width: 100%; max-height: 70vh; border: 1px solid #d9d9d9; border-radius: 6px;',
        onError: (e) => {
          e.target.style.display = 'none';
          const errorDiv = document.createElement('div');
          errorDiv.style.color = '#ff4d4f';
          errorDiv.style.padding = '20px';
          errorDiv.innerHTML = '图片加载失败';
          e.target.parentNode.appendChild(errorDiv);
        }
      })
    ]),
    width: 'auto',
    okText: '关闭'
  });
};

// 团体成员二维码弹窗中，生成某成员二维码
const showMemberQr = async (member) => {
  try {
    const baseUrl = window.location.origin;
    const params = new URLSearchParams({
      type: 'group',
      year: member.year,
      group_name: member.group_name,
      category: member.category || '',
      person_name: member.name,
      id_number: member.id_number,
      award_id: member.award_id
    });
    const verifyUrl = `${baseUrl}/verify.html?${params.toString()}`;
    const qrCodeDataUrl = await QRCode.toDataURL(verifyUrl, {
      width: 200,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    });
    Modal.info({
      title: `成员【${member.name}】的验证二维码`,
      content: createVNode('div', { style: 'text-align: center;' }, [
        createVNode('div', { style: 'margin-bottom: 16px;' }, [
          createVNode('p', { style: 'margin-bottom: 8px; font-weight: bold;' }, '扫描二维码或复制链接进行验证：'),
          createVNode('img', { 
            src: qrCodeDataUrl, 
            alt: '验证二维码',
            style: 'display: block; margin: 0 auto 16px auto; border: 1px solid #d9d9d9; border-radius: 6px;'
          }),
          createVNode('div', { style: 'background: #f5f5f5; padding: 12px; border-radius: 6px; word-break: break-all; font-family: monospace; font-size: 12px;' }, verifyUrl),
          createVNode('div', { style: 'margin-top: 12px;' }, [
            createVNode('button', {
              style: 'background: #1890ff; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer;',
              onClick: () => {
                navigator.clipboard.writeText(verifyUrl).then(() => {
                  message.success('验证链接已复制到剪贴板');
                }).catch(() => {
                  message.error('复制失败，请手动复制链接');
                });
              }
            }, '复制链接')
          ])
        ])
      ]),
      width: 400,
      okText: '关闭'
    });
  } catch (error) {
    console.error('生成成员验证二维码失败:', error);
    message.error('生成成员验证二维码失败');
  }
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleString('zh-CN');
};

// 处理图片加载错误
const handleImageError = (event, imageType) => {
  console.error(`图片加载失败: ${imageType}`);
  event.target.style.display = 'none';
  const errorDiv = document.createElement('div');
  errorDiv.style.color = '#ff4d4f';
  errorDiv.style.fontSize = '12px';
  errorDiv.style.marginTop = '4px';
  errorDiv.innerHTML = `${imageType}加载失败`;
  event.target.parentNode.appendChild(errorDiv);
};

// 标签页切换处理
const onTabChange = (key) => {
  activeTab.value = key;
  selectedRowKeys.value = []; // 切换标签页时清空选中状态
};

// 添加表格筛选变化处理函数
const handleTableChange = (pagination, filters, sorter) => {
  // 更新表格筛选状态
  tableFilters.value = {
    award_type: filters.award_type?.[0] || null,
    year: filters.year?.[0] || null,
    group_name: filters.group_name?.[0] || null,
    category: filters.category?.[0] || null,
    award_level: filters.award_level?.[0] || null
  };
};

// 性能监控
const performanceInfo = ref({
  recordCount: 0,
  memoryUsage: 0,
  loadTime: 0
});

// 证书加载状态
const certificateLoading = ref(false);

// 内存清理功能
const clearMemory = () => {
  try {
    // 清理已加载的证书图片数据
    if (allRecords.value) {
      allRecords.value.forEach(record => {
        // 将已加载的证书图片重置为标记状态
        if (record.certificate_image && record.certificate_image !== 'has_certificate' && record.has_certificate) {
          record.certificate_image = 'has_certificate';
        }
        
        // 清理团体成员的证书图片
        if (record.members && Array.isArray(record.members)) {
          record.members.forEach(member => {
            if (member.certificate_image && member.certificate_image !== 'has_certificate' && member.has_certificate) {
              member.certificate_image = 'has_certificate';
            }
          });
        }
      });
    }
    
    // 清理选中的记录
    selectedRecord.value = null;
    
    // 清理编辑表单
    if (editForm.value) {
      editForm.value.certificate_image = '';
      if (editForm.value.members) {
        editForm.value.members.forEach(member => {
          member.certificate_image = '';
        });
      }
    }
    
    // 清理文件列表
    personCertificateFileList.value = [];
    memberCertificateFileList.value = [];
    
    // 强制垃圾回收（如果浏览器支持）
    if (window.gc) {
      window.gc();
    }
    
    // 更新性能信息
    updatePerformanceInfo();
    
    message.success('内存清理完成！已释放证书图片占用的内存');
    console.log('内存清理完成');
  } catch (error) {
    console.error('内存清理失败:', error);
    message.error('内存清理失败');
  }
};

// 高级内存清理功能
const advancedClearMemory = async () => {
  try {
    // 基础清理
    clearMemory();
    
    // 清理图片缓存
    const images = document.querySelectorAll('img');
    images.forEach(img => {
      if (img.src && img.src.startsWith('data:image/')) {
        img.src = '';
      }
    });
    
    // 清理Canvas缓存
    const canvases = document.querySelectorAll('canvas');
    canvases.forEach(canvas => {
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
      }
    });
    
    // 清理localStorage中的临时数据
    const keysToRemove = [];
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && (key.includes('temp_') || key.includes('cache_'))) {
        keysToRemove.push(key);
      }
    }
    keysToRemove.forEach(key => localStorage.removeItem(key));
    
    // 清理sessionStorage
    sessionStorage.clear();
    
    // 强制页面重绘
    document.body.style.display = 'none';
    document.body.offsetHeight; // 触发重排
    document.body.style.display = '';
    
    message.success('高级内存清理完成！已清理所有缓存和临时数据');
    console.log('高级内存清理完成');
  } catch (error) {
    console.error('高级内存清理失败:', error);
    message.error('高级内存清理失败');
  }
};

// 自动内存清理（定期清理）
const autoClearMemory = () => {
  // 每5分钟自动清理一次内存
  setInterval(() => {
    if (allRecords.value && allRecords.value.length > 0) {
      console.log('执行自动内存清理...');
      clearMemory();
    }
  }, 5 * 60 * 1000); // 5分钟
};

// 更新性能信息
const updatePerformanceInfo = () => {
  if (allRecords.value) {
    performanceInfo.value.recordCount = allRecords.value.length;
    
    // 估算内存使用（每条记录大约1KB）
    performanceInfo.value.memoryUsage = Math.round(allRecords.value.length * 1);
    
    // 如果有performance API，获取加载时间
    if (window.performance && window.performance.memory) {
      performanceInfo.value.memoryUsage = Math.round(window.performance.memory.usedJSHeapSize / 1024 / 1024);
    }
  }
};

// 智能内存监控
const memoryMonitor = () => {
  // 每30秒检查一次内存使用情况
  setInterval(() => {
    if (window.performance && window.performance.memory) {
      const usedMemory = window.performance.memory.usedJSHeapSize / 1024 / 1024; // MB
      const totalMemory = window.performance.memory.totalJSHeapSize / 1024 / 1024; // MB
      const memoryUsagePercent = (usedMemory / totalMemory) * 100;
      
      // 如果内存使用超过80%，提醒用户清理
      if (memoryUsagePercent > 80) {
        message.warning(`内存使用率较高 (${memoryUsagePercent.toFixed(1)}%)，建议清理内存以提升性能`);
      }
      
      // 更新性能信息
      performanceInfo.value.memoryUsage = Math.round(usedMemory * 1024); // 转换为KB
    }
  }, 30 * 1000); // 30秒
};

// 页面加载时自动加载记录
onMounted(() => {
  checkUserInfo();
  if (userInfo.value) {
    loadAllRecords();
    // 显示优化提示
    message.info('系统已优化为轻量级模式，证书图片将按需加载，提升页面性能');
    // 启动自动内存清理
    autoClearMemory();
    // 启动智能内存监控
    memoryMonitor();
  }
});

// 导出组件数据和方法
defineExpose({
  userInfo,
  loading,
  allRecords,
  searchText,
  filterType,
  filterYear,
  filterCategory,
  selectedRowKeys,
  detailModalVisible,
  selectedRecord,
  activeTab,
  editModalVisible,
  editForm,
  editLoading,
  editFormRef,
  editNewMember,
  personCertificateFileList,
  memberQrModalVisible,
  memberQrList,
  tableFilters,
  editRules,
  personRecords,
  groupRecords,
  personColumns,
  groupColumns,
  availableYears,
  filteredRecords,
  checkUserInfo,
  loadAllRecords,
  onSearch,
  onFilterChange,
  onSelectChange,
  viewRecord,
  editRecord,
  saveEdit,
  addEditMember,
  removeEditMember,
  cancelEdit,
  beforeUpload,
  handlePersonCertificateChange,
  handleMemberCertificateChange,
  deleteRecord,
  batchDelete,
  getAwardLevelColor,
  generateVerifyLink,
  showMemberQr,
  previewImage,
  formatDate,
  onTabChange,
  handleTableChange,
  performanceInfo,
  updatePerformanceInfo,
  certificateLoading,
  clearMemory,
  advancedClearMemory,
  autoClearMemory,
  memoryMonitor
});
</script>

<style scoped>
.ant-table {
  font-size: 14px;
}

.ant-descriptions-item-label {
  font-weight: bold;
}

/* 证书上传样式 */
.certificate-uploader .ant-upload {
  width: 120px;
  height: 120px;
}

.certificate-uploader .ant-upload-select-picture-card i {
  font-size: 32px;
  color: #999;
}

.certificate-uploader .ant-upload-select-picture-card .ant-upload-text {
  margin-top: 8px;
  color: #666;
}

/* 团体成员编辑样式 */
.members-edit-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.member-edit-card {
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  padding: 12px;
  background: #fafafa;
}

.member-edit-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.member-edit-text {
  display: flex;
  align-items: center;
  gap: 8px;
}

.member-name {
  font-weight: 500;
  color: #1f2937;
}

.member-id {
  color: #6b7280;
  font-size: 0.9rem;
}

.member-certificate-upload {
  display: flex;
  align-items: center;
  gap: 12px;
}

.certificate-upload-label {
  font-size: 0.9rem;
  color: #6b7280;
  min-width: 70px;
}

.member-certificate-uploader .ant-upload {
  width: 80px;
  height: 80px;
}

.member-certificate-uploader .upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
}

.member-certificate-uploader .upload-placeholder .anticon {
  font-size: 20px;
  margin-bottom: 4px;
}

/* 成员显示样式 */
.members-display {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.members-display .ant-tag {
  margin: 0;
  font-size: 12px;
  line-height: 20px;
  border-radius: 10px;
}
</style>
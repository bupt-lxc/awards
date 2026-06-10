<template>
  <a-tabs v-model:activeKey="activeKey">
    <a-tab-pane key="single" tab="单条录入">
      <a-form :model="formState" :label-col="labelCol" :wrapper-col="wrapperCol" :rules="rules" ref="formRef">
        <a-form-item label="类型" name="award_type">
          <a-radio-group v-model:value="formState.award_type" button-style="solid">
            <a-radio-button value="group">团体</a-radio-button>
            <a-radio-button value="person">个人</a-radio-button>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="年份" name="year">
          <a-input-number v-model:value="formState.year" :min="2000" :max="2030" style="width: 120px" />
        </a-form-item>

        <a-form-item label="团名" name="group_name">
          <a-input v-model:value="formState.group_name" placeholder="请输入团体名称" />
        </a-form-item>

        <!-- 团体成员信息 -->
        <a-form-item label="团体成员" v-if="formState.award_type === 'group'">
          <div style="border: 1px solid #d9d9d9; border-radius: 6px; padding: 16px; background-color: #fafafa;">
            <div style="margin-bottom: 16px;">
              <div style="display: flex; align-items: center; gap: 12px; margin-bottom: 12px;">
                <a-input 
                  v-model:value="newMember.name" 
                  placeholder="成员姓名" 
                  style="width: 150px;"
                />
                <a-input 
                  v-model:value="newMember.id_number" 
                  placeholder="证件号码" 
                  style="width: 180px;"
                />
                <a-button type="primary" @click="addMember" :disabled="!newMember.name || !newMember.id_number">
                  添加成员
                </a-button>
              </div>
              
              <!-- 成员证书上传 -->
              <div style="display: flex; align-items: center; gap: 12px;">
                <span style="font-size: 12px; color: #666;">成员证书（可选）：</span>
                <a-upload
                  :before-upload="handleMemberCertificateUpload"
                  :show-upload-list="false"
                  accept="image/*"
                >
                  <a-button size="small" type="dashed">
                    <upload-outlined /> 上传证书
                  </a-button>
                </a-upload>
                <span v-if="newMember.certificate_image" style="font-size: 12px; color: #52c41a;">✓ 已选择证书图片</span>
                <a-button v-if="newMember.certificate_image" size="small" type="link" @click="previewMemberCertificate">预览</a-button>
                <a-button v-if="newMember.certificate_image" size="small" type="link" danger @click="removeMemberCertificate">删除</a-button>
              </div>
            </div>
            
            <div v-if="formState.members && formState.members.length > 0">
              <div style="margin-bottom: 8px; font-weight: 500;">成员列表 ({{ formState.members.length }}人)：</div>
              <div style="max-height: 200px; overflow-y: auto;">
                <div 
                  v-for="(member, index) in formState.members" 
                  :key="index"
                  style="margin-bottom: 8px; padding: 8px; border: 1px solid #e8e8e8; border-radius: 4px; background: white;"
                >
                  <div style="display: flex; justify-content: space-between; align-items: center;">
                    <span>{{ member.name }} ({{ member.id_number }})</span>
                    <div style="display: flex; gap: 8px; align-items: center;">
                      <span v-if="member.certificate_image" style="font-size: 12px; color: #52c41a;">📄 有证书</span>
                      <a-button v-if="member.certificate_image" size="small" type="link" @click="previewMemberCertificateInList(member)">预览证书</a-button>
                      <a-button size="small" danger @click="removeMember(index)">删除</a-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-else style="color: #999; font-style: italic;">
              暂无成员，请添加团体成员信息
            </div>
          </div>
        </a-form-item>

        <a-form-item label="姓名" name="person_name" v-if="formState.award_type === 'person'">
          <a-input v-model:value="formState.person_name" placeholder="请输入个人姓名" />
        </a-form-item>

        <a-form-item label="证件号码" name="id_number" v-if="formState.award_type === 'person'">
          <a-input v-model:value="formState.id_number" placeholder="请输入证件号码" />
        </a-form-item>

        <a-form-item label="组别" name="category">
          <a-input v-model:value="formState.category" style="width: 200px" placeholder="请输入组别">
          </a-input>
        </a-form-item>

        <a-form-item label="奖项等级" name="award_level">
          <a-input v-model:value="formState.award_level" style="width: 120px" placeholder="请输入奖项等级">
          </a-input>
        </a-form-item>

        <a-form-item label="描述" name="description">
          <a-textarea v-model:value="formState.description" placeholder="请输入获奖描述（可选）" :rows="3" />
        </a-form-item>

        <!-- 个人获奖证书图片上传 -->
        <a-form-item label="获奖证书" v-if="formState.award_type === 'person'">
          <div style="border: 1px solid #d9d9d9; border-radius: 6px; padding: 16px; background-color: #fafafa;">
            <a-upload
              :before-upload="handlePersonCertificateUpload"
              :show-upload-list="false"
              accept="image/*"
              style="margin-bottom: 16px;"
            >
              <a-button type="dashed" style="width: 100%; height: 80px;">
                <div>
                  <upload-outlined style="font-size: 20px; margin-bottom: 8px;" />
                  <div>点击上传获奖证书图片</div>
                  <div style="font-size: 12px; color: #999;">支持 JPG、PNG、GIF 格式</div>
                </div>
              </a-button>
            </a-upload>
            
            <div v-if="formState.certificate_image" style="display: flex; align-items: center; justify-content: center; gap: 12px;">
              <span style="font-size: 14px; color: #52c41a;">✓ 已选择证书图片</span>
              <a-button size="small" type="link" @click="previewPersonCertificate">预览</a-button>
              <a-button size="small" type="link" danger @click="removePersonCertificate">删除</a-button>
            </div>
          </div>
        </a-form-item>

        <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
          <a-button type="primary" @click="onSubmit" :loading="loading">提交</a-button>
          <a-button style="margin-left: 10px" @click="onReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-tab-pane>

    <a-tab-pane key="batch" tab="Excel批量上传">
      <div style="padding: 24px;">
        <div style="margin-bottom: 24px; text-align: center;">
          <a-button type="link" @click="downloadTemplate" size="large">
            <download-outlined /> 下载Excel模板
          </a-button>
        </div>

        <div style="margin-bottom: 24px;">
          <div style="font-size: 14px; margin-bottom: 12px;">1. 上传Excel文件：</div>
        <a-upload 
          :before-upload="handleExcelUpload" 
          :show-upload-list="false" 
          accept=".xlsx,.xls"
          class="upload-area"
        >
          <a-button type="dashed" :loading="batchLoading" size="large" style="width: 100%; height: 120px;">
            <div>
              <upload-outlined style="font-size: 24px; margin-bottom: 8px;" />
              <div style="font-size: 16px;">点击选择Excel文件</div>
              <div style="font-size: 12px; color: #999; margin-top: 8px;">
                支持.xlsx和.xls格式<br>
                  请确保Excel文件包含以下列：类型、年份、团名、姓名、证件号码、组别、奖项等级、描述<br>
                  <strong>图片命名规则：</strong>请将证书图片命名为"证件号码.jpg/png"，系统会自动匹配
              </div>
            </div>
          </a-button>
        </a-upload>
        </div>

        <div style="margin-bottom: 24px;">
          <div style="font-size: 14px; margin-bottom: 12px;">2. 上传证书图片（可选）：</div>
          <a-upload
            :before-upload="handleBatchImageUpload"
            :multiple="true"
            :show-upload-list="true"
            accept="image/*"
            class="upload-area"
            v-model:fileList="batchImageFiles"
          >
            <a-button type="dashed" size="large" style="width: 100%;">
              <div>
                <upload-outlined style="font-size: 16px; margin-right: 8px;" />
                <span>选择证书图片</span>
              </div>
            </a-button>
          </a-upload>
          <div style="font-size: 12px; color: #666; margin-top: 8px;">
            提示：请将图片文件命名为对应的证件号码，如：110101199001011234.jpg
          </div>
        </div>

        <!-- 预览解析结果 -->
        <div v-if="excelData.length > 0" style="margin-top: 24px;">
          <h3>预览数据 ({{ excelData.length }}条记录)</h3>
          <div style="margin-bottom: 12px; font-size: 12px; color: #666;">
            已匹配证书图片：{{ getMatchedImagesCount() }} 张
          </div>
          <a-table 
            :dataSource="excelData" 
            :columns="previewColumns" 
            :pagination="{ pageSize: 10 }" 
            size="small"
            :scroll="{ x: 800 }"
            style="margin-top: 16px;"
          >
            <template #bodyCell="{ column, record, index }">
              <template v-if="column.key === 'action'">
                <a-button size="small" danger @click="removeExcelRow(index)">删除</a-button>
              </template>
              <template v-if="column.key === 'certificate_status'">
                <template v-if="record.award_type === 'person'">
                  <div v-if="record.certificate_image">
                    <a-tag color="success" @click="() => previewBatchCertificate(record.certificate_image)" style="cursor: pointer;">
                      ✓ {{ record.certificate_filename }}
                    </a-tag>
                  </div>
                  <a-tag v-else color="default">未匹配</a-tag>
                </template>
                <template v-else-if="record.award_type === 'group'">
                  <div v-if="!record.members || record.members.length === 0" color="default">无成员</div>
                  <div v-else>
                     <div v-for="member in record.members" :key="member.id_number" style="margin-bottom: 4px;">
                       <a-tag v-if="member.certificate_image" color="success" @click="() => previewBatchCertificate(member.certificate_image)" style="cursor: pointer;">
                         ✓ {{ member.certificate_filename }}
                       </a-tag>
                       <a-tag v-else color="default">
                         {{ member.name }}: 未匹配
                       </a-tag>
                     </div>
                  </div>
                </template>
              </template>
            </template>
          </a-table>

          <div style="margin-top: 24px; text-align: center;">
            <a-button type="primary" @click="batchSubmit" :loading="batchLoading" :disabled="excelData.length === 0" size="large">
              批量提交 ({{ excelData.length }}条记录)
            </a-button>
            <a-button style="margin-left: 16px" @click="clearExcelData" size="large">清空数据</a-button>
          </div>
        </div>
      </div>
    </a-tab-pane>
  </a-tabs>
  <a-modal v-if="batchProgress.visible" :visible="batchProgress.visible" :footer="null" :closable="false" :maskClosable="false" centered width="360px">
    <div style="text-align:center; padding:32px 0;">
      <a-progress :percent="batchProgress.percent" status="active" />
      <div style="margin-top:16px; font-size:16px;">{{ batchProgress.text }}</div>
    </div>
  </a-modal>
</template>
<script setup>
import { reactive, toRaw, ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { UploadOutlined, DownloadOutlined } from '@ant-design/icons-vue';
import { awardAPI } from '../api/index.js';
import * as XLSX from 'xlsx';

// 表单引用
const formRef = ref();
const loading = ref(false);
const batchLoading = ref(false);
const excelData = ref([]);
const activeKey = ref('single');
const batchProgress = ref({ visible: false, percent: 0, text: '' });

// 表单数据
const formState = reactive({
  award_type: 'group',
  year: new Date().getFullYear(),
  group_name: '',
  person_name: '',
  id_number: '',
  category: '',
  award_level: '',
    description: '',
    certificate_image: '',
    members: []
  });

// 新成员输入
const newMember = reactive({
  name: '',
  id_number: '',
  certificate_image: ''
});

// 表单验证规则
const rules = {
  award_type: [{ required: true, message: '请选择奖项类型', trigger: 'change' }],
  year: [{ required: true, message: '请输入年份', trigger: 'blur' }],
  group_name: [
    {
      required: true,
      message: '请输入团体名称',
      trigger: 'blur',
      validator: (rule, value) => {
        if (formState.award_type === 'group' && !value) {
          return Promise.reject('请输入团体名称');
        }
        return Promise.resolve();
      }
    }
  ],
  person_name: [
    {
      required: true,
      message: '请输入个人姓名',
      trigger: 'blur',
      validator: (rule, value) => {
        if (formState.award_type === 'person' && !value) {
          return Promise.reject('请输入个人姓名');
        }
        return Promise.resolve();
      }
    }
  ],
  id_number: [
    {
      required: true,
      message: '请输入证件号码',
      trigger: 'blur',
      validator: (rule, value) => {
        if (formState.award_type === 'person' && !value) {
          return Promise.reject('请输入证件号码');
        }
        return Promise.resolve();
      }
    }
  ],
  category: [{ required: true, message: '请选择组别', trigger: 'change' }],
  award_level: [{ required: true, message: '请选择奖项等级', trigger: 'change' }]
};

// 批量图片上传文件列表
const batchImageFiles = ref([]);

// 个人获奖证书上传处理
const handlePersonCertificateUpload = async (file) => {
  try {
    if (!file.type.startsWith('image/')) {
      message.error(`请选择图片文件（${file.name} 不是图片格式）`);
      return false;
    }
    if (file.size > 10 * 1024 * 1024) {
      message.error(`图片大小不能超过10MB（${file.name} 过大）`);
      return false;
    }
    // 上传图片到服务器，获取静态文件路径
    const response = await awardAPI.uploadCertificate(file, 'single');
    formState.certificate_image = response.data.file_path;
    message.success(`证书图片上传成功：${file.name}`);
  } catch (error) {
    message.error(`图片上传失败：${file.name}`);
    console.error('图片上传失败:', error);
  }
  return false;
};

// 新增：预览个人证书
const previewPersonCertificate = () => {
  if (formState.certificate_image) {
    const newWindow = window.open();
    if (newWindow) {
      newWindow.document.write(`<body style="margin:0; display:flex; justify-content:center; align-items:center; background-color:#333;"><img src="${formState.certificate_image}" style="max-width:100%; max-height:100%;" /></body>`);
    } else {
      message.error('预览窗口被浏览器拦截，请允许弹窗。');
    }
  }
};

// 删除个人获奖证书
const removePersonCertificate = () => {
  formState.certificate_image = '';
  message.success('证书图片已删除');
};

// 团体成员证书上传处理
const handleMemberCertificateUpload = async (file) => {
  try {
    if (!file.type.startsWith('image/')) {
      message.error('请选择图片文件');
      return false;
    }
    if (file.size > 10 * 1024 * 1024) {
      message.error('图片大小不能超过10MB');
      return false;
    }
    const response = await awardAPI.uploadCertificate(file, 'member');
    newMember.certificate_image = response.data.file_path;
    message.success('成员证书图片上传成功');
  } catch (error) {
    console.error('图片上传失败:', error);
    message.error('图片上传失败');
  }
  return false;
};

// 提交表单
const onSubmit = async () => {
  try {
    // 表单验证
    await formRef.value.validate();
    // 单条录入校验图片
    if (formState.award_type === 'person' && !formState.certificate_image) {
      message.error('请先上传获奖证书图片');
      return;
    }
    // 团体奖项成员验证
    if (formState.award_type === 'group' && (!formState.members || formState.members.length === 0)) {
      message.error('团体奖项必须添加至少一个成员');
      return;
    }
    loading.value = true;
    const submitData = {
      award_type: formState.award_type,
      year: formState.year,
      category: formState.category,
      award_level: formState.award_level,
      description: formState.description || undefined
    };
    if (formState.award_type === 'group') {
      submitData.group_name = formState.group_name;
      submitData.members = formState.members.map(m => ({
        ...m,
        id_number: String(m.id_number)
      }));
    } else {
      submitData.person_name = formState.person_name;
      submitData.id_number = String(formState.id_number);
      submitData.certificate_image = formState.certificate_image || undefined;
    }
    console.log('[单条录入] 提交数据:', JSON.parse(JSON.stringify(submitData)));
    let response;
    if (submitData.award_type === 'person') {
      response = await awardAPI.uploadPersonAward(submitData);
    } else {
      response = await awardAPI.uploadGroupAward(submitData);
    }
    message.success('获奖记录提交成功！');
    console.log('提交成功:', response);
    onReset();
  } catch (error) {
    handleApiError(error, '提交');
  } finally {
    loading.value = false;
  }
};

// 添加成员
const addMember = () => {
  if (newMember.name && newMember.id_number) {
    // 检查是否已存在相同的成员
    const exists = formState.members.some(member => 
      member.name === newMember.name && member.id_number === newMember.id_number
    );
    
    if (exists) {
      message.warning('该成员已存在');
      return;
    }
    
    formState.members.push({
      name: newMember.name,
      id_number: String(newMember.id_number),
      certificate_image: newMember.certificate_image || null
    });
    
    // 清空输入框
    newMember.name = '';
    newMember.id_number = '';
    newMember.certificate_image = '';
    
    message.success('成员添加成功');
  }
};

// 移除成员
const removeMember = (index) => {
  formState.members.splice(index, 1);
  message.success('成员移除成功');
};

// 重置表单
const onReset = () => {
  formRef.value.resetFields();
  Object.assign(formState, {
    award_type: 'group',
    year: new Date().getFullYear(),
    group_name: '',
    person_name: '',
    id_number: '',
    category: '',
    award_level: '',
    description: '',
    certificate_image: '',
    members: []
  });
  
  // 重置新成员输入
  newMember.name = '';
  newMember.id_number = '';
  newMember.certificate_image = '';
};

// 预览表格列定义
const previewColumns = ref([
  { title: '类型', dataIndex: 'award_type', key: 'award_type', width: 80 },
  { title: '年份', dataIndex: 'year', key: 'year', width: 80 },
  { title: '团名', dataIndex: 'group_name', key: 'group_name', width: 120 },
  { title: '姓名', dataIndex: 'person_name', key: 'person_name', width: 100 },
  { title: '证件号码', dataIndex: 'id_number', key: 'id_number', width: 150 },
  { title: '团体成员', dataIndex: 'members_display', key: 'members_display', width: 200, ellipsis: true },
  { title: '组别', dataIndex: 'category', key: 'category', width: 120 },
  { title: '奖项等级', dataIndex: 'award_level', key: 'award_level', width: 100 },
  { title: '证书状态', key: 'certificate_status', width: 100 },
  { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: '操作', key: 'action', width: 80, fixed: 'right' }
]);

// Excel文件上传处理
const handleExcelUpload = (file) => {
  const reader = new FileReader();

  reader.onload = (e) => {
    try {
      const data = new Uint8Array(e.target.result);
      const workbook = XLSX.read(data, { type: 'array' });
      const sheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[sheetName];
      const jsonData = XLSX.utils.sheet_to_json(worksheet);

      // 校验Excel表头
      const requiredHeaders = ['类型', '年份', '组别', '奖项等级'];
      const actualHeaders = Object.keys(jsonData[0] || {});
      const missingHeaders = requiredHeaders.filter(h => !actualHeaders.includes(h));

      if (missingHeaders.length > 0) {
        message.error(`Excel文件缺少必要的列: ${missingHeaders.join(', ')}。请下载最新模板。`);
        return;
      }
      
      console.log('Excel原始数据:', jsonData);

      // 确保所有证件号码都是字符串格式
      const rawData = jsonData.map((row, index) => {
        const award_type = row['类型'] || row['奖项类型'] || row['type'] || '';
        const year = parseInt(row['年份'] || row['year'] || new Date().getFullYear());
        const group_name = row['团名'] || row['团体名称'] || row['group_name'] || '';
        const person_name = row['姓名'] || row['个人姓名'] || row['person_name'] || '';
        const id_number = String(row['证件号码'] || row['id_number'] || '').trim(); // 确保是字符串并去除空格
        const category = row['组别'] || row['category'] || '';
        const award_level = row['奖项等级'] || row['等级'] || row['award_level'] || '';
        const description = row['描述'] || row['description'] || '';

        // 打印每条记录的证件号码，用于调试
        console.log(`记录 ${index + 1} 的证件号码:`, id_number);

        return {
          key: index,
          award_type: award_type === '团体' ? 'group' : award_type === '个人' ? 'person' : award_type,
          year,
          group_name,
          person_name,
          id_number,
          category,
          award_level,
          description,
          has_certificate: false
        };
      }).filter(item => item.award_type && item.year && item.category && item.award_level);

      if (rawData.length === 0) {
        console.log('Excel解析结果:', rawData);
        message.error('Excel文件中没有找到有效数据，请检查文件格式');
        return;
      }

      // 处理团体成员合并逻辑
      const processedData = processGroupMembers(rawData);
      
      // 打印处理后的数据，用于调试
      console.log('Excel处理后的数据:', processedData);
      
      excelData.value = processedData;

      const personCount = processedData.filter(p => p.award_type === 'person').length;
      const groupCount = processedData.filter(p => p.award_type === 'group').length;
      message.success(`成功解析 ${personCount} 条个人记录和 ${groupCount} 个团体记录。`);

    } catch (error) {
      console.error('Excel解析失败:', error);
      message.error('Excel文件解析失败，请检查文件格式');
    }
  };

  reader.readAsArrayBuffer(file);
  return false; // 阻止自动上传
};

// 处理团体成员合并
const processGroupMembers = (rawData) => {
  const groupMap = new Map();
  const personAwards = [];
  const processedData = [];
  
  rawData.forEach(item => {
    const itemWithStringId = {
      ...item,
      id_number: String(item.id_number).trim(),
      certificate_image: item.certificate_image || ''
    };
    if (item.award_type === 'group' && item.group_name && item.person_name && item.id_number) {
      const groupKey = `${item.group_name}_${item.year}_${item.category}_${item.award_level}`;
      if (!groupMap.has(groupKey)) {
        groupMap.set(groupKey, {
          award_type: 'group',
          year: item.year,
          group_name: item.group_name,
          category: item.category,
          award_level: item.award_level,
          description: item.description,
          members: []
        });
      }
      groupMap.get(groupKey).members.push({
        name: item.person_name,
        id_number: String(item.id_number).trim(),
        certificate_image: item.certificate_image || ''
      });
    } else if (item.award_type === 'person' && item.person_name && item.id_number) {
      personAwards.push(itemWithStringId);
    }
  });
  let keyIndex = 0;
  groupMap.forEach(groupAward => {
    processedData.push({
      ...groupAward,
      key: keyIndex++,
      members_display: groupAward.members.length > 0 ? 
        `${groupAward.members.length}人: ${groupAward.members.map(m => m.name).join(', ')}` : 
        '无成员信息'
    });
  });
  personAwards.forEach(personAward => {
    processedData.push({
      ...personAward,
      key: keyIndex++
    });
  });
  return processedData;
};

// 删除Excel数据行
const removeExcelRow = (index) => {
  excelData.value.splice(index, 1);
  message.success('已删除该行数据');
};

// 清空Excel数据
const clearExcelData = () => {
  excelData.value = [];
  message.success('已清空所有数据');
};

// 批量提交Excel数据
const batchSubmit = async () => {
  if (excelData.value.length === 0) {
    message.warning('没有数据可提交');
    return;
  }

  // 检查未匹配证书的记录
  let missing = [];
  excelData.value.forEach(item => {
    if (item.award_type === 'person') {
      if (!item.certificate_image) {
        missing.push(`个人-${item.person_name}(${item.id_number})`);
      }
    } else if (item.award_type === 'group' && item.members) {
      item.members.forEach(m => {
        if (!m.certificate_image) {
          missing.push(`团体成员-${m.name}(${m.id_number})`);
        }
      });
    }
  });

  if (missing.length > 0) {
    const confirmed = await new Promise(resolve => {
      Modal.confirm({
        title: '确认提交',
        content: createVNode('div', {}, [
          createVNode('p', { style: 'color: #ff4d4f; margin-bottom: 16px;' },
            `发现有 ${missing.length} 条记录未匹配证书图片。`
          ),
          createVNode('p', { style: 'margin-top: 16px;' }, '是否仍要继续提交？')
        ]),
        okText: '继续提交',
        cancelText: '返回检查',
        width: 500,
        onOk: () => resolve(true),
        onCancel: () => resolve(false),
      });
    });
    if (!confirmed) return;
  }

  try {
    batchLoading.value = true;
    batchProgress.value = { visible: true, percent: 0, text: '准备上传...' };

    const BATCH_SIZE = 5;
    const allPersonAwards = [];
    const allGroupAwards = [];

    batchProgress.value.text = '正在准备数据...';
    batchProgress.value.percent = 5;

    for (const item of excelData.value) {
      const data = {
        award_type: item.award_type,
        year: item.year,
        category: item.category,
        award_level: item.award_level,
        description: item.description || undefined
      };

      if (item.award_type === 'group') {
        data.group_name = item.group_name;
        if (item.members && item.members.length > 0) {
          data.members = item.members.map(m => ({
            name: m.name,
            id_number: String(m.id_number),
            certificate_image: m.certificate_image || undefined
          }));
        } else {
          data.members = [];
        }
        if (item.certificate_image) {
          data.certificate_image = item.certificate_image;
        }
        allGroupAwards.push(data);
      } else if (item.award_type === 'person') {
        data.person_name = item.person_name;
        data.id_number = String(item.id_number);
        data.certificate_image = item.certificate_image || undefined;
        allPersonAwards.push(data);
      }
    }

    // 计算总批次
    const totalBatches = Math.ceil(allPersonAwards.length / BATCH_SIZE) + Math.ceil(allGroupAwards.length / BATCH_SIZE);
    let currentBatch = 0;
    let totalUploaded = 0;
    const results = [];

    // 上传个人记录
    let currentBatchSize = BATCH_SIZE;
    for (let i = 0; i < allPersonAwards.length; i += currentBatchSize) {
      currentBatch++;
      batchProgress.value.text = `正在上传个人记录: ${currentBatch}/${totalBatches} (批量大小: ${currentBatchSize})`;
      batchProgress.value.percent = Math.round((currentBatch / totalBatches) * 100);
      
      const batch = allPersonAwards.slice(i, i + currentBatchSize);
      console.log(`[批量上传] 个人记录批次 ${currentBatch}:`, batch);
      
      try {
        const personResponse = await awardAPI.batchUploadPersonAwards(batch);
        totalUploaded += personResponse.uploaded_count || personResponse.successful || 0;
        results.push(`个人获奖记录${personResponse.uploaded_count || personResponse.successful || 0}条`);
      } catch (error) {
        // 如果是413错误且批量大小大于1，则减小批量大小重试
        if (error.response?.status === 413 && currentBatchSize > 1) {
          currentBatchSize = Math.max(1, Math.floor(currentBatchSize / 2));
          message.warning(`检测到数据过大，自动调整批量大小为 ${currentBatchSize}，正在重试...`);
          i -= currentBatchSize; // 重新处理当前批次
          currentBatch--; // 重置批次计数
          continue;
        }
        throw error; // 重新抛出其他错误
      }
    }

    // 上传团体记录
    currentBatchSize = BATCH_SIZE; // 重置批量大小
    for (let i = 0; i < allGroupAwards.length; i += currentBatchSize) {
      currentBatch++;
      batchProgress.value.text = `正在上传团体记录: ${currentBatch}/${totalBatches} (批量大小: ${currentBatchSize})`;
      batchProgress.value.percent = Math.round((currentBatch / totalBatches) * 100);
      
      const batch = allGroupAwards.slice(i, i + currentBatchSize);
      console.log(`[批量上传] 团体记录批次 ${currentBatch}:`, batch);
      
      try {
        const groupResponse = await awardAPI.batchUploadGroupAwards(batch);
        totalUploaded += groupResponse.uploaded_count || groupResponse.successful || 0;
        results.push(`团体获奖记录${groupResponse.uploaded_count || groupResponse.successful || 0}条`);
      } catch (error) {
        // 如果是413错误且批量大小大于1，则减小批量大小重试
        if (error.response?.status === 413 && currentBatchSize > 1) {
          currentBatchSize = Math.max(1, Math.floor(currentBatchSize / 2));
          message.warning(`检测到数据过大，自动调整批量大小为 ${currentBatchSize}，正在重试...`);
          i -= currentBatchSize; // 重新处理当前批次
          currentBatch--; // 重置批次计数
          continue;
        }
        throw error; // 重新抛出其他错误
      }
    }

    batchProgress.value.text = '上传完成';
    batchProgress.value.percent = 100;
    setTimeout(() => { batchProgress.value.visible = false; }, 1200);

    message.success(`批量上传成功！共上传${totalUploaded}条记录（${results.join('，')}）`);
    clearExcelData();
    batchImageFiles.value = [];
  } catch (error) {
    console.error('[批量上传] 错误:', error);
    handleApiError(error, '批量提交');
    batchProgress.value.visible = false;
  } finally {
    batchLoading.value = false;
  }
};

// 统一处理API错误，提供更友好的提示
const handleApiError = (error, operation = '操作') => {
  console.error(`${operation}失败:`, error);
  
  // 处理413错误（请求体过大）
  if (error.response?.status === 413) {
    message.error(`${operation}失败：请求数据过大，请减少批量上传的记录数量或压缩图片大小。建议每批不超过5条记录。`);
    return;
  }
  
  if (error.response?.data?.detail) {
    const detail = error.response.data.detail;
    if (typeof detail === 'string') {
      if (detail.includes('Duplicate entry')) {
        message.error(`${operation}失败：存在重复记录，请检查证件号码是否已存在。`);
        return;
      }
      if (detail.includes('validation failed')) {
        message.error(`${operation}失败：数据校验不通过，请检查填写的内容。`);
        return;
      }
    }
    message.error(`${operation}失败: ${typeof detail === 'string' ? detail : '请查看控制台日志'}`);
  } else if (error.message) {
    message.error(`${operation}失败: ${error.message}`);
  } else {
    message.error(`${operation}失败，请检查网络连接或稍后重试。`);
  }
};

// 辅助函数：标准化ID，用于精确匹配
const normalizeId = (id) => {
  if (id === null || id === undefined) return '';
  // 转换为字符串，移除所有空白，全角转半角，并转为小写
  return String(id)
    .replace(/\s/g, '') // 移除所有空白字符
    .replace(/[\uff01-\uff5e]/g, (ch) => String.fromCharCode(ch.charCodeAt(0) - 0xfee0)) // 全角转半角
    .toLowerCase();
};

// 获取已匹配图片的数量
const getMatchedImagesCount = () => {
  let count = 0;
  if (!excelData.value) return 0;
  excelData.value.forEach(record => {
    if (record.award_type === 'person' && record.certificate_image) {
      count++;
    } else if (record.award_type === 'group' && Array.isArray(record.members)) {
      count += record.members.filter(member => member.certificate_image).length;
    }
  });
  return count;
};

// 预览成员证书（添加成员时）
const previewMemberCertificate = () => {
  if (newMember.certificate_image) {
    const newWindow = window.open();
    newWindow.document.write(`
      <html>
        <head><title>证书预览</title></head>
        <body style="margin: 0; display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #f0f0f0;">
          <img src="${newMember.certificate_image}" style="max-width: 90%; max-height: 90%; border: 1px solid #ddd;" />
        </body>
      </html>
    `);
  }
};

// 删除成员证书（添加成员时）
const removeMemberCertificate = () => {
  newMember.certificate_image = '';
  message.success('成员证书图片已删除');
};

// 预览成员证书（成员列表中）
const previewMemberCertificateInList = (member) => {
  if (member.certificate_image) {
    const newWindow = window.open();
    newWindow.document.write(`
      <html>
        <head><title>${member.name} - 证书预览</title></head>
        <body style="margin: 0; display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #f0f0f0;">
          <div style="text-align: center;">
            <h3 style="margin-bottom: 20px;">${member.name} 的获奖证书</h3>
            <img src="${member.certificate_image}" style="max-width: 90%; max-height: 80%; border: 1px solid #ddd;" />
          </div>
        </body>
      </html>
    `);
  }
};

const downloadTemplate = () => {
  // 创建包含证书图片列的Excel模板
  const templateData = [
    {
      '类型': 'person',
      '年份': 2024,
      '团名': '',
      '姓名': '张三',
      '证件号码': '110101199001011234',
      '组别': '教师组',
      '奖项等级': '一等奖',
      '描述': '优秀教学成果',
      '证书图片': ''
    },
    {
      '类型': 'group',
      '年份': 2024,
      '团名': '创新团队',
      '姓名': '李四',
      '证件号码': '110101199002022345',
      '组别': '科研组',
      '奖项等级': '特等奖',
      '描述': '科技创新项目',
      '证书图片': ''
    },
    {
      '类型': 'group',
      '年份': 2024,
      '团名': '创新团队',
      '姓名': '王五',
      '证件号码': '110101199003033456',
      '组别': '科研组',
      '奖项等级': '特等奖',
      '描述': '科技创新项目',
      '证书图片': ''
    }
  ];

  // 创建工作簿
  const ws = XLSX.utils.json_to_sheet(templateData);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, '获奖记录模板');

  // 下载文件
  XLSX.writeFile(wb, '获奖记录上传模板.xlsx');
  message.success('模板下载成功');
};

// 批量图片上传处理 — 上传到服务器获取静态路径，再匹配记录
const handleBatchImageUpload = async (file) => {
  try {
    if (!file.type.startsWith('image/')) {
      message.error(`请选择图片文件（${file.name} 不是图片格式）`);
      return false;
    }
    if (file.size > 10 * 1024 * 1024) {
      message.error(`图片大小不能超过10MB（${file.name} 过大）`);
      return false;
    }

    // 上传到服务器获取静态文件路径
    const response = await awardAPI.uploadCertificate(file, 'batch');
    const filePath = response.data.file_path;

    const nameWithoutExt = file.name.replace(/\.[^.]+$/, '');
    let matched = false;

    for (const record of excelData.value) {
      if (record.award_type === 'person') {
        if (normalizeId(record.id_number) === normalizeId(nameWithoutExt)) {
          record.certificate_image = filePath;
          record.certificate_filename = file.name;
          matched = true;
          break;
        }
      } else if (record.award_type === 'group') {
        if (nameWithoutExt === record.group_name) {
          record.certificate_image = filePath;
          record.certificate_filename = file.name;
          matched = true;
          break;
        } else if (record.members) {
          for (const member of record.members) {
            if (normalizeId(member.id_number) === normalizeId(nameWithoutExt)) {
              member.certificate_image = filePath;
              member.certificate_filename = file.name;
              matched = true;
              break;
            }
          }
          if (matched) break;
        }
      }
    }

    excelData.value = [...excelData.value];

    if (matched) {
      message.success(`成功匹配证书图片：${file.name}`);
    } else {
      message.warning(`未找到匹配的记录：${file.name}，请检查证件号或团体名称是否一致`);
    }

    return false;
  } catch (error) {
    console.error('图片处理失败:', error);
    message.error(`图片处理失败：${error.message || '未知错误'}`);
    return false;
  }
};

// 图片上传统一使用 awardAPI.uploadCertificate(file, idNumber)
// 返回的 response.data.file_path 即为静态文件路径（如 /static/certificates/uuid.png）

// 新增：预览批量上传的证书
const previewBatchCertificate = (imagePath) => {
  if (!imagePath) {
    message.warning('没有可预览的图片。');
    return;
  }
  const newWindow = window.open();
  newWindow.document.write(`<body style="margin:0; display:flex; justify-content:center; align-items:center; background-color:#333;"><img src="${imagePath}" style="max-width:100%; max-height:100%;" /></body>`);
};

const labelCol = {
  style: {
    width: '150px',
  },
};
const wrapperCol = {
  span: 14,
};
</script>

<style scoped>
.upload-area {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 24px auto;
  max-width: 600px;
}
</style>
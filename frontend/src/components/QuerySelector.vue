<template>
  <div class="query-selector-container">
    <!-- 页面头部 -->
    <div class="header-section">
      <div class="header-content">
        <h1 class="main-title">
          获奖记录查询系统
        </h1>
        <p class="subtitle">欢迎使用获奖记录查询系统</p>
      </div>
    </div>

    <!-- 选择区域 -->
    <div class="selector-section">
      <div class="selector-grid">
        <!-- 个人/团队查询卡片 -->
        <div class="query-card person-card" @click="selectQueryType('person')" style="max-width: 600px; margin: 0 auto;">
          <div class="card-icon">
            👤
          </div>
          <h3 class="card-title">个人/团队奖项查询</h3>
          <p class="card-description">
            查询个人或团队获奖记录<br>
            需要提供：年份、姓名/团队名称、获奖编码
          </p>
          <div class="card-features">
            <div class="feature-item">✓ 精确查询获奖信息</div>
            <div class="feature-item">✓ 支持按年份筛选</div>
            <div class="feature-item">✓ 详细获奖记录展示</div>
          </div>
          <a-button type="primary" size="large" class="select-btn">
            开始查询
          </a-button>
        </div>
      </div>
    </div>

    <!-- 功能说明 -->
    <div class="info-section">
      <a-card class="info-card" :bordered="false">
        <template #title>
          <div class="info-title">
            查询说明
          </div>
        </template>
        
        <div class="info-content">
          <div class="info-item">
            <h4>查询说明</h4>
            <p>您可以通过输入获奖年份、姓名/团队名称和获奖编码来查询具体的获奖记录。系统会根据您提供的信息进行精确匹配，并显示详细的获奖信息。</p>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script>
export default {
  name: 'QuerySelector',
  emits: ['select-query-type'],
  setup(props, { emit }) {
    const selectQueryType = (type) => {
      emit('select-query-type', type)
    }
    
    return {
      selectQueryType
    }
  }
}
</script>

<style scoped>
.query-selector-container {
  min-height: calc(100vh - 80px);
  background: linear-gradient(135deg, #FFFAF0 0%, #FFF8DC 100%);
  padding: 40px 20px;
  position: relative;
  overflow: hidden;
}

.query-selector-container::before,
.query-selector-container::after {
  content: '';
  position: absolute;
  border-radius: 50%;
  z-index: 0;
}

.query-selector-container::before {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(222, 184, 135, 0.08) 0%, rgba(255, 228, 181, 0.05) 50%, transparent 70%);
  top: -200px;
  right: -200px;
}

.query-selector-container::after {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(255, 228, 181, 0.08) 0%, rgba(222, 184, 135, 0.05) 50%, transparent 70%);
  bottom: -100px;
  left: -100px;
}

.header-section {
  background: linear-gradient(135deg, #FFE4B5 0%, #FFF8DC 100%);
  padding: 40px;
  border-radius: 24px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(25, 118, 210, 0.15);
  margin-bottom: 40px;
  text-align: center;
}

.header-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M50,0 C77.6142375,0 100,22.3857625 100,50 C100,77.6142375 77.6142375,100 50,100 C22.3857625,100 0,77.6142375 0,50 C0,22.3857625 22.3857625,0 50,0 Z M50,10 C27.90861,10 10,27.90861 10,50 C10,72.09139 27.90861,90 50,90 C72.09139,90 90,72.09139 90,50 C90,27.90861 72.09139,10 50,10 Z' fill='%23ffffff' fill-opacity='0.05'/%3E%3C/svg%3E") repeat;
  opacity: 0.1;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.main-title {
  font-size: 2.8rem;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  position: relative;
  display: inline-block;
  padding-bottom: 8px;
  text-shadow: none;
}

.main-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 80%;
  height: 4px;
  background: #ffc107;
  border-radius: 2px;
}

.subtitle {
  font-size: 1.2rem;
  color: #C14B4B;
  margin: 0;
  letter-spacing: 1px;
  position: relative;
  display: inline-block;
  padding: 0 30px;
}

.subtitle::before,
.subtitle::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 20px;
  height: 1px;
  background: rgba(255, 255, 255, 0.5);
}

.subtitle::before {
  right: 100%;
}

.subtitle::after {
  left: 100%;
}

.selector-section {
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

.selector-grid {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.query-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(10px);
  border: 2px solid transparent;
}

.query-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(222, 184, 135, 0.03) 0%, rgba(255, 228, 181, 0.03) 100%);
  opacity: 0;
  transition: opacity 0.4s ease;
}

.query-card:hover {
  transform: translateY(-8px);
  border-color: #DEB887;
  box-shadow: 0 20px 40px rgba(222, 184, 135, 0.25);
}

.query-card:hover::before {
  opacity: 1;
}

.card-icon {
  font-size: 3rem;
  margin-bottom: 24px;
  display: inline-block;
  position: relative;
  color: #DEB887;
}

.card-icon::after {
  content: '';
  position: absolute;
  width: 60px;
  height: 60px;
  background: rgba(25, 118, 210, 0.1);
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: -1;
}

.card-title {
  font-size: 1.8rem;
  font-weight: 600;
  background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 16px;
}

.card-description {
  font-size: 1.1rem;
  color: #C14B4B;
  margin-bottom: 24px;
  line-height: 1.6;
  font-weight: 500;
}

.card-features {
  text-align: left;
  margin: 24px 0;
  padding: 0 10px;
}

.feature-item {
  color: #C14B4B;
  margin-bottom: 12px;
  font-size: 1rem;
  display: flex;
  align-items: center;
  position: relative;
  padding-left: 28px;
  font-weight: 500;
}

.feature-item::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: #dd6b20;
  font-weight: bold;
}

.select-btn {
  background: linear-gradient(90deg, #E6904F, #F9BE50);
  border: none;
  border-radius: 24px;
  height: 48px;
  padding: 0 32px;
  font-weight: 600;
  font-size: 1.1rem;
  letter-spacing: 1px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(222, 184, 135, 0.15);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: white;
}

.select-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(249, 190, 80, 0.35);
  background: linear-gradient(90deg, #F9BE50, #E6904F);
}



@media (max-width: 1200px) {
  .selector-grid {
    gap: 30px;
    padding: 10px;
  }
}

@media (max-width: 768px) {
  .header-section {
    padding: 24px 8px;
  }
  .header-content {
    gap: 8px;
    padding: 0 4px;
    max-width: 100vw;
  }
  .main-title, .subtitle {
    text-align: center;
    word-break: break-all;
    max-width: 100vw;
    margin: 0 auto;
  }
  .selector-section {
    max-width: 100vw;
    margin: 0 auto;
    padding: 0 4px;
  }
  .selector-grid {
    grid-template-columns: 1fr;
    gap: 20px;
    padding: 0;
  }
}
@media (max-width: 480px) {
  .main-title {
    font-size: 1.2rem;
  }
  .subtitle {
    font-size: 0.9rem;
  }
  .selector-section {
    padding: 0 2vw;
  }
}

.info-section {
  max-width: 1000px;
  margin: 60px auto 0;
  padding: 0 20px;
}

.info-card {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
  backdrop-filter: blur(10px);
}

.info-title {
  font-size: 1.5rem;
  background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 1.8rem;
}

.info-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  padding: 20px 0;
}

.info-item {
  padding: 20px;
  background: rgba(222, 184, 135, 0.03);
  border-radius: 16px;
  transition: all 0.3s ease;
}

.info-item:hover {
  background: rgba(222, 184, 135, 0.06);
}

.info-item h4 {
  color: #DEB887;
  font-size: 1.2rem;
  margin-bottom: 12px;
  font-weight: 600;
}

.info-item p {
  color: #C14B4B;
  line-height: 1.6;
  margin: 0;
  font-weight: 500;
}
</style>
<template>
  <div class="public-app">
    <!-- 顶部导航栏 -->
    <header class="top-header">
      <div class="header-container new-header-layout">
        <img src="/logo.jpg" alt="logo" class="logo-img-large" />
        <div class="title-section">
          <h1 class="logo-title">香港國際合唱節</h1>
          <p class="logo-subtitle">HONG KONG INTERNATIONAL CHORAL FESTIVAL</p>
        </div>
        <div class="nav-section">
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="main-content">
      <div class="carousel-section">
        <div class="carousel-container">
          <transition-group name="fade" tag="div" class="carousel-wrapper">
            <div v-for="(img, idx) in carouselImages" 
                 :key="img" 
                 v-show="currentCarouselIndex === idx" 
                 class="carousel-item">
              <img :src="img" class="carousel-img" :alt="`轮播图 ${idx + 1}`" />
              <div class="carousel-content">
                <div class="carousel-text">
                  <h2 class="carousel-title">{{ carouselTitles[idx] }}</h2>
                  <p class="carousel-description">{{ carouselDescriptions[idx] }}</p>
                </div>
              </div>
            </div>
          </transition-group>
          
          <!-- 轮播图指示器 -->
          <div class="carousel-indicators">
            <span v-for="(_, idx) in carouselImages" 
                  :key="idx"
                  :class="['indicator-dot', { active: currentCarouselIndex === idx }]"
                  @click="setCarouselIndex(idx)">
            </span>
          </div>
          
          <!-- 轮播图控制按钮 -->
          <button class="carousel-control prev" @click="prevSlide" aria-label="上一张">
            <i class="carousel-arrow">&#10094;</i>
          </button>
          <button class="carousel-control next" @click="nextSlide" aria-label="下一张">
            <i class="carousel-arrow">&#10095;</i>
          </button>
        </div>
      </div>
      
      <!-- 返回按钮 -->
      <div v-if="currentQueryType" class="back-section">
        <a-button 
          type="text" 
          size="large" 
          @click="goBack"
          class="back-btn"
        >
          <template #icon>
            <arrow-left-outlined />
          </template>
          返回选择页面
        </a-button>
      </div>
      
      <!-- 查询选择器 -->
      <QuerySelector 
        v-if="!currentQueryType" 
        @select-query-type="handleSelectQueryType" 
      />
      
      <!-- 个人查询 -->
      <PersonQueryView v-else-if="currentQueryType === 'person'" />
      
      <!-- 团体查询 -->
      <GroupQueryView v-else-if="currentQueryType === 'group'" />
    </main>

    <!-- 页脚 -->
    <footer class="app-footer">
      <div class="footer-container">
        <div class="footer-bottom">
          <p>香港國際合唱節組委會  地址：深圳市南山区后海滨路保利剧院  电话：0755-86371696</p>
          <p class="icp-text"><a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer">备案/许可证编号为：苏ICP备2026042125号</a></p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { ArrowLeftOutlined } from '@ant-design/icons-vue'
import QuerySelector from './components/QuerySelector.vue'
import PersonQueryView from './components/PersonQueryView.vue'
import GroupQueryView from './components/GroupQueryView.vue'

// 当前查询类型
const currentQueryType = ref(null)

// 选择查询类型
const handleSelectQueryType = (type) => {
  currentQueryType.value = type
}

// 返回选择页面
const goBack = () => {
  currentQueryType.value = null
}



const carouselImages = [
  '/logo.jpg',
  '/logo.jpg',
  '/logo.jpg',
  '/logo.jpg',
]

const carouselTitles = [
  '2025南山国际合唱周',
  '共创音乐盛典',
  '青春激扬梦想',
  '艺术交流平台'
]

const carouselDescriptions = [
  '汇聚全球青少年音乐才华，共襄国际文化盛事',
  '打造世界级合唱艺术交流平台',
  '让音乐之声传递青春与梦想',
  '促进国际文化艺术交流与发展'
]

const currentCarouselIndex = ref(0)

const setCarouselIndex = (index) => {
  currentCarouselIndex.value = index
}

const nextSlide = () => {
  currentCarouselIndex.value = (currentCarouselIndex.value + 1) % carouselImages.length
}

const prevSlide = () => {
  currentCarouselIndex.value = currentCarouselIndex.value === 0 
    ? carouselImages.length - 1 
    : currentCarouselIndex.value - 1
}

// 自动轮播
let autoplayInterval
onMounted(() => {
  autoplayInterval = setInterval(nextSlide, 5000)
})

onUnmounted(() => {
  clearInterval(autoplayInterval)
})
</script>

<style scoped>
.public-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #FFF9F0 0%, #FFF6E5 100%);
  position: relative;
  overflow-x: hidden;
}
.public-app::before {
  content: '';
  position: fixed;
  left: -120px;
  top: 80px;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 155, 155, 0.15) 0%, rgba(255, 214, 153, 0.1) 80%, transparent 100%);
  z-index: 0;
  pointer-events: none;
  filter: blur(2px);
}
.public-app::after {
  content: '';
  position: fixed;
  right: -100px;
  bottom: 0;
  width: 320px;
  height: 320px;
  background: radial-gradient(circle, rgba(255, 214, 153, 0.15) 0%, rgba(255, 155, 155, 0.1) 80%, transparent 100%);
  z-index: 0;
  pointer-events: none;
  filter: blur(2px);
}

.carousel-section {
  width: 100%;
  position: relative;
  overflow: hidden;
  padding: 20px;
  background: rgba(26, 26, 26, 0.05);
}

.carousel-container {
  position: relative;
  width: 100%;
  aspect-ratio: 21/9;
  max-height: 320px;
  overflow: hidden;
  max-width: 1000px;
  margin: 0 auto;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.carousel-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
}

.carousel-item {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.carousel-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: all 0.8s ease;
  object-position: center 40%;
  transform: scale(1.02); /* 轻微放大以避免边缘可能出现的白边 */
  filter: brightness(0.95); /* 轻微调暗以增加文字可读性 */
}

.carousel-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 20px 60px 40px;
}

.carousel-text {
  text-align: center;
  color: #fff;
  max-width: 800px;
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.8s ease;
  background: rgba(0, 0, 0, 0.3);
  padding: 20px 30px;
  border-radius: 16px;
  backdrop-filter: blur(8px);
}

.carousel-item.active .carousel-text {
  opacity: 1;
  transform: translateY(0);
}

.carousel-title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 20px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  line-height: 1.2;
}

.carousel-description {
  font-size: 1.3rem;
  line-height: 1.6;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
  margin: 0;
}

/* 轮播图控制按钮 */
.carousel-control {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10;
  backdrop-filter: blur(2px);
}

.carousel-control:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.5);
}

.carousel-control.prev {
  left: 30px;
}

.carousel-control.next {
  right: 30px;
}

.carousel-arrow {
  color: #fff;
  font-size: 24px;
  font-style: normal;
}

/* 轮播图指示器 */
.carousel-indicators {
  position: absolute;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
  z-index: 10;
  background: rgba(0, 0, 0, 0.3);
  padding: 8px 16px;
  border-radius: 20px;
}

.indicator-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator-dot.active {
  background: #fff;
  transform: scale(1.2);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.8s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 返回按钮样式 */
.back-section {
  position: fixed;
  top: 100px;
  left: 20px;
  z-index: 1000;
}

.back-btn {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  padding: 8px 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  color: #667eea;
  font-weight: 600;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 1);
  color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

/* 顶部导航栏样式 */
.top-header {
  background: linear-gradient(135deg, #FFE4B5 0%, #FFF8DC 100%);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 10000;
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  height: 64px;
  position: relative;
  justify-content: space-between;
}

.header-container.new-header-layout {
  gap: 16px;
}

/* Logo部分 */
.logo-img-large {
  height: 48px;
  width: auto;
  flex-shrink: 0;
  justify-self: start;
}

/* 标题部分 */
.title-section {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  text-align: center;
}

.logo-title {
  font-size: 1.8rem;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  letter-spacing: 1px;
  white-space: nowrap;
  text-align: center;
  text-shadow: none;
}

.logo-subtitle {
  font-size: 0.95rem;
  margin: 2px 0 0;
  font-weight: 300;
  letter-spacing: 0.5px;
  color: #C14B4B;
  white-space: nowrap;
  text-align: center;
}

/* 导航部分 */
.nav-section {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 20px;
}

.nav-title {
  color: #fff;
  font-size: 1.1rem;
  font-weight: 500;
  white-space: nowrap;
}



/* 主要内容区域 */
.main-content {
  flex: 1;
  min-height: calc(100vh - 80px - 200px); /* 减去header和footer的高度 */
}

/* 页脚样式 */
.app-footer {
  background: linear-gradient(135deg, #FFE4B5 0%, #FFF8DC 100%);
  color: #fff;
  margin-top: auto;
}

.footer-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 24px 20px;
}

.footer-bottom {
  text-align: center;
}
.footer-bottom p {
  margin: 4px 0;
  font-size: 1rem;
  background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  display: block;
  font-weight: 500;
}

.footer-bottom .icp-text a {
  font-size: 1rem;
  font-weight: 500;
  background: linear-gradient(90deg, #C14B4B, #E6904F, #F9BE50);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  text-decoration: none;
}

/* 响应式设计 */
@media (max-width: 1400px) {
  .logo-img-large {
    margin-left: 0;
  }
}

@media (max-width: 1200px) {
  .header-container {
    grid-template-columns: 160px 1fr 180px;
    gap: 16px;
  }

  .carousel-container {
    max-width: 800px;
    max-height: 280px;
    aspect-ratio: 18/9;
  }
  
  .carousel-content {
    padding: 20px 40px 30px;
  }
  
  .carousel-text {
    padding: 20px;
    max-width: 600px;
  }
  
  .carousel-title {
    font-size: 2rem;
    margin-bottom: 12px;
  }
  
  .carousel-description {
    font-size: 1.1rem;
  }
}

@media (max-width: 900px) {
  .header-container {
    grid-template-columns: 120px 1fr 160px;
    padding: 0 16px;
  }

  .logo-img-large {
    height: 56px;
  }

  .logo-title {
    font-size: 1.8rem;
  }

  .logo-subtitle {
    font-size: 0.9rem;
  }

  .nav-title {
    font-size: 1rem;
  }

  .carousel-section {
    padding: 15px;
  }

  .carousel-container {
    max-height: 240px;
    border-radius: 12px;
    aspect-ratio: 16/9;
  }
  
  .carousel-content {
    padding: 20px 30px 30px;
  }
  
  .carousel-text {
    padding: 16px;
    max-width: 500px;
  }
  
  .carousel-title {
    font-size: 1.8rem;
    margin-bottom: 10px;
  }
  
  .carousel-description {
    font-size: 1rem;
    line-height: 1.4;
  }
  
  .carousel-control {
    width: 36px;
    height: 36px;
  }
  
  .carousel-control.prev {
    left: 10px;
  }
  
  .carousel-control.next {
    right: 10px;
  }

  .carousel-img {
    object-position: center 45%;
  }
}

@media (max-width: 768px) {
  .header-container {
    flex-direction: column;
    padding: 10px;
    height: auto;
    gap: 10px;
  }
  .logo-img-large {
    height: 40px;
    margin: 0;
  }
  .title-section {
    position: relative;
    left: 0;
    transform: none;
    padding: 5px 0;
  }
  .logo-title {
    font-size: 1.1rem;
    margin: 0;
    white-space: normal;
  }
  .logo-subtitle {
    font-size: 0.95rem;
    line-height: 1.3;
    margin: 5px 0 0;
  }
  .nav-section {
    margin: 5px 0;
  }
}

@media (max-width: 480px) {
  .header-container {
    padding: 0;
    gap: 8px;
    width: 100vw;
    max-width: 100vw;
  }

  .logo-img-large {
    height: 40px;
  }

  .logo-title {
    font-size: 1.4rem;
  }

  .logo-subtitle {
    font-size: 0.8rem;
  }

  .admin-btn {
    padding: 0 12px;
  }

  .carousel-container {
    height: 200px;
  }
  
  .carousel-content {
    padding: 15px 20px 30px;
  }
  
  .carousel-text {
    padding: 15px;
  }
  
  .carousel-title {
    font-size: 1.6rem;
    margin-bottom: 8px;
  }
  
  .carousel-description {
    font-size: 1rem;
    line-height: 1.4;
  }
  
  .carousel-control {
    width: 35px;
    height: 35px;
  }
  
  .carousel-indicators {
    bottom: 20px;
    gap: 8px;
    padding: 6px 12px;
  }
  
  .indicator-dot {
    width: 8px;
    height: 8px;
  }
}
</style>
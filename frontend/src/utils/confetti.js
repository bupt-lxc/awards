import confetti from 'canvas-confetti';

// 默认的礼花配置
const defaultConfettiConfig = {
  particleCount: 100,
  spread: 70,
  origin: { y: 0.6 },
  colors: ['#ff0000', '#00ff00', '#0000ff', '#ffff00', '#ff00ff', '#00ffff'],
  startVelocity: 30,
  gravity: 0.8,
  ticks: 200,
  zIndex: 9999,
  disableForReducedMotion: true, // 为了可访问性
  scalar: 1 // 缩放因子，移动端可以调整
};

// 检查是否为移动设备
const isMobile = () => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
};

// 创建礼花效果
export const fireConfetti = (customConfig = {}) => {
  // 如果是移动设备，调整一些参数
  const mobileAdjustments = isMobile() ? {
    particleCount: 50, // 减少粒子数量
    spread: 45, // 减少扩散范围
    startVelocity: 20, // 降低初始速度
    scalar: 0.7 // 缩小整体比例
  } : {};

  // 合并配置
  const config = {
    ...defaultConfettiConfig,
    ...mobileAdjustments,
    ...customConfig
  };

  // 创建左右两侧的礼花
  const leftConfetti = {
    ...config,
    origin: { x: 0.2, y: 0.6 },
    angle: 60
  };

  const rightConfetti = {
    ...config,
    origin: { x: 0.8, y: 0.6 },
    angle: 120
  };

  // 发射礼花
  confetti(leftConfetti);
  confetti(rightConfetti);

  // 中间的礼花
  setTimeout(() => {
    confetti({
      ...config,
      origin: { x: 0.5, y: 0.7 },
      spread: 90,
      angle: 90
    });
  }, 200);
}; 
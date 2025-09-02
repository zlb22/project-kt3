import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
interface addedProps {
  [key: string]: number
}
export const useWorkspaceStore = defineStore('workspace', () => {
  const step = ref(1)

  const showToolbar = ref(false) // 是否显示工具栏
  const selectedProp = ref<any>(null) // 选中的道具
  const showRecord = ref(false) // 是否显示录音
  const remainTime = ref(480) // 剩余时间
  const showTimeModal = ref(false) // 是否显示时间弹窗
  const timeId = ref<number | any>(null) // 定时器id
  const running = ref(false)

  // 已添加的道具列表
  const addedProps = ref<addedProps>({
    Pentagon: 0,
    Rect: 0,
    Triangle: 0,
    Hook: 0,
    Circle: 0,
    Cone: 0,
    TriangularPyramid: 0,
    Cube: 0,
    Cylinder: 0,
    Semicircle: 0,
    Line: 0,
    Curve: 0
  })

  // 道具列表，需要提前设置默认值
  const propsList = ref([
    {
      type: 'Pentagon', // 道具类型
      name: '五边形', // 道具标题
      top: 250,
      left: 450,
      width: 81,
      height: 77,
      src: new URL('@/assets/images/pentagon_icon.png', import.meta.url).href
    },
    {
      type: 'Rect', // 道具类型
      name: '矩形', // 道具标题
      top: 300,
      left: 300,
      width: 72,
      height: 72,
      fill: '#C9AFFF',
      stroke: '#A36EF1',
      strokeWidth: 1,
      src: new URL('@/assets/images/rect_icon.png', import.meta.url).href
    },
    {
      type: 'Triangle', // 道具类型
      name: '三角形', // 道具标题
      top: 300,
      left: 500,
      width: 73,
      height: 67,
      fill: '#C9AFFF',
      stroke: '#A36EF1',
      strokeWidth: 1,
      src: new URL('@/assets/images/triangle_icon.png', import.meta.url).href
    },
    {
      type: 'Hook', // 道具类型
      name: '弯钩', // 道具标题
      top: 200,
      left: 400,
      width: 52,
      height: 74,
      src: new URL('@/assets/images/hook_icon.png', import.meta.url).href
    },
    {
      type: 'Circle', // 道具类型
      name: '圆形', // 道具标题
      top: 300,
      left: 300,
      radius: 36,
      width: 72,
      height: 72,
      fill: '#C9AFFF',
      stroke: '#A36EF1',
      strokeWidth: 1,
      src: new URL('@/assets/images/circle_icon.png', import.meta.url).href
    },
    {
      type: 'Semicircle', // 道具类型
      name: '半圆', // 道具标题
      top: 250,
      left: 450,
      width: 80,
      height: 42,
      src: new URL('@/assets/images/semicircle_icon.png', import.meta.url).href
    },
    {
      type: 'TriangularPyramid', // 道具类型
      name: '三棱锥', // 道具标题
      top: 150,
      left: 490,
      width: 77,
      height: 76,
      src: new URL('@/assets/images/triangular_pyramid_icon.png', import.meta.url).href
    },
    {
      type: 'Cube', // 道具类型
      name: '正方体', // 道具标题
      top: 300,
      left: 520,
      width: 73,
      height: 73,
      src: new URL('@/assets/images/cube_icon.png', import.meta.url).href
    },
    {
      type: 'Cylinder', // 道具类型
      name: '圆柱', // 道具标题
      top: 280,
      left: 530,
      width: 74,
      height: 75,
      src: new URL('@/assets/images/cylinder_icon.png', import.meta.url).href
    },
    {
      type: 'Cone', // 道具类型
      name: '圆锥', // 道具标题
      top: 130,
      left: 450,
      width: 74,
      height: 77,
      src: new URL('@/assets/images/cone_icon.png', import.meta.url).href
    },
    
    {
      type: 'Line',
      name: '直线',
      top: 250,
      left: 550,
      width: 60,
      height: 60,
      stroke: '#C4A5FF',
      strokeWidth: 4,
      src: new URL('@/assets/images/line_icon.png', import.meta.url).href
    },
    {
      type: 'Curve',
      name: '曲线',
      stroke: '#C4A5FF',
      strokeWidth: 4,
      top: 250,
      left: 550,
      width: 42,
      height: 60,
      src: new URL('@/assets/images/curve_icon.png', import.meta.url).href
    }
  ])

  function setStep(value: number) {
    console.log('Setting step to', value)
    step.value = value
  }

  function countProps(type: string) {
    addedProps.value[type] = addedProps.value[type] + 1
  }

  // 是否显示录音
  function isShowRecord(val: boolean) {
    showRecord.value = val
  }
  // 隐藏工具栏
  function hideToolbar() {
    showToolbar.value = false
  }
  // 设置剩余时间
  function setRemainTime(time: number) {
    remainTime.value = time
  }
  // 是否显示时间提示弹窗
  function isShowTimeModal(val: boolean) {
    showTimeModal.value = val
  }

  function setRunning(val: boolean) {
    running.value = val
  }

  function startCountdown() {
    if (remainTime.value > 0 && !running.value) {
      running.value = true
      const timerId = setInterval(() => {
        if (remainTime.value > 0) {
          remainTime.value--
        } else {
          showTimeModal.value = true
          clearInterval(timerId)
          running.value = false
        }
      }, 1000)
      timeId.value = timerId
    }
  }

  return {
    step,
    propsList,
    setStep,
    countProps,
    addedProps,
    showToolbar,
    selectedProp,
    showRecord,
    isShowRecord,
    hideToolbar,
    remainTime,
    setRemainTime,
    showTimeModal,
    isShowTimeModal,
    timeId,
    running,
    setRunning,
    startCountdown
  }
})

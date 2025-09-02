<template>
  <div class="second-wrapper">
    <div class="workspace-container" ref="canvasContainer">
      <PropsBar @add-prop="addProp" v-if="!workspaceStore.showRecord" />
      <ForwardAndBackwardControl v-if="!workspaceStore.showRecord" />
      <ToolBar v-if="workspaceStore.showToolbar" @flipX="flipX" @flipY="flipY" @bringToFront="bringToFront"
        @sendToBack="sendToBack" @bringForward="bringForward" @sendBackwards="sendBackwards" />
      <canvas id="canvas" class="edit-box" ref="canvasRef"></canvas>
      <MyRecord v-if="workspaceStore.showRecord"></MyRecord>
      <div class="opacity-overlay" v-if="workspaceStore.showRecord"></div>
      <ConfirmCommit v-if="fileStore.confirmCommitStatus"></ConfirmCommit>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { fabric } from 'fabric'
import PropsBar from './components/PropsBar.vue'
import ForwardAndBackwardControl from './components/ForwardAndBackwardControl.vue'
import ConfirmCommit from './components/ConfirmCommit.vue'
import propsMap from '../core/propsMap'
import { useOperationHistory } from '../stores/operationHistory'
import ToolBar from './components/ToolBar.vue'
import { useWorkspaceStore } from '@/stores/workspace'
import MyRecord from './components/MyRecord.vue'
import { useFileStore } from '@/stores/files'

const fileStore = useFileStore()
const canvasContainer = ref<any>(null)
let canvas: any = null // 画布
let selectedProp: any = null // 选中的道具
const operationHistory = useOperationHistory()
const workspaceStore = useWorkspaceStore()
const canvasRef = ref<any>(null)
const timer = ref()

onMounted(() => {
  initCanvas()
  handleClickEvent()
  deleteEventHandler()

  timer.value = setInterval(() => {
    operationHistory.updateStayOperaDuration()
  }, 1000)

  watch(
    () => workspaceStore.showRecord,
    (newVal) => {
      if (newVal) {
        canvas.discardActiveObject().renderAll()
      }
    }
  )
})

watch(
  () => fileStore.confirmCommitStatus,
  (newVal) => {
    if (newVal) {
      takeScreenshot()
      console.log(fileStore.imgFile)
    }
  }
)
/**
 * 初始化画布
 */
const initCanvas = () => {
  const width = canvasContainer.value.clientWidth
  const height = canvasContainer.value.clientHeight

  canvas = new fabric.Canvas('canvas', {
    backgroundColor: '#E4EBFD', // 设置画布背景颜色
    selection: false, // 禁止多选
    preserveObjectStacking: true, // 保持对象层次
  })

  canvas.setWidth(width)
  canvas.setHeight(height)

  fabric.Object.prototype.cornerColor = '#FFBE14' // 设置控制点颜色
  fabric.Object.prototype.cornerSize = 7 // 控制点大小
  fabric.Object.prototype.transparentCorners = false // 控制点不透明
  fabric.Object.prototype.borderColor = '#FFBE14' // 控制线颜色

  operationHistory.bindCanvasEvent(canvas)
  fabric.Image.fromURL(
    new URL('@/assets/images/bridge_new.png', import.meta.url).href,
    function (img: any) {
      const scaleRatio = height / img.height
      img.scale(scaleRatio)
      canvas.setBackgroundImage(img, canvas.renderAll.bind(canvas))
    }
  )
}

/**
 * 添加道具
 * @param item
 */
const addProp = (item: any) => {
  const { type } = item
  const propConstructor = propsMap[type]
  const Prop = new propConstructor(item, canvas)
  const { propData } = Prop
  if (propData) {
    canvas.add(propData)
  }
}
/**
 * 点击事件
 */
const handleClickEvent = () => {
  canvas.on('mouse:up', (event: any) => {
    setTimeout(() => {
      const activeObject = canvas.getActiveObject()

      // 页面上没有选中道具
      if (!activeObject) {
        handleOutsideClick()
      } else {
        // console.log('Clicked on an object')
      }
    })
  })
}

/**
 * 处理点击道具外事件
 */
const handleOutsideClick = () => {
  // 隐藏工具栏
  workspaceStore.showToolbar = false
  workspaceStore.selectedProp = null
}

/**
 * 键盘删除事件
 */
const deleteEventHandler = () => {
  document.onkeydown = (e: any) => {
    if (e.keyCode === 8 || e.keyCode === 46) {
      workspaceStore.hideToolbar()
      deleteProp()
    }
  }
}
/**
 * 删除道具
 */
const deleteProp = () => {
  selectedProp = canvas.getActiveObject()
  selectedProp.set({
    selectable: false,
    visible: false
  })

  canvas.discardActiveObject()
  canvas.renderAll()

  canvas.fire('object:deleted', { target: selectedProp })
}
// 水平翻转
const flipX = () => {
  const hookProp = canvas.getActiveObject()
  hookProp.flipX = !hookProp.flipX
  canvas.renderAll()

  canvas.fire('object:flipX', { target: hookProp })
}
// 垂直翻转
const flipY = () => {
  const hookProp = canvas.getActiveObject()
  hookProp.flipY = !hookProp.flipY
  canvas.renderAll()

  canvas.fire('object:flipY', { target: hookProp })
}
// 置于顶层
const bringToFront = () => {
  const hookProp = canvas.getActiveObject()
  canvas.bringToFront(hookProp)
  canvas.renderAll()

  canvas.fire('object:bringToFront', { target: hookProp })
}
// 置于底层
const sendToBack = () => {
  const hookProp = canvas.getActiveObject()
  canvas.sendToBack(hookProp)
  canvas.renderAll()


  canvas.fire('object:sendToBack', { target: hookProp })
}
// 上移一层
const bringForward = () => {
  // 由于有组合元素 需要找到要去的位置  然后通过 .moveTo 移动过去
  const hookProp = canvas.getActiveObject()
  let toIndex = -1

  console.log('bringForward====>', canvas.getObjects(), hookProp)

  const objects = canvas.getObjects().map((item: any) => {
    if (item.name === 'curve-line' && item.idObj) {
      Object.assign(item, item.idObj)
    }
    return item
  })

  const index = objects.findIndex((item: any) => {
    if (hookProp.name === 'curve-line' && hookProp.idObj) {
      Object.assign(hookProp, hookProp.idObj)
    }

    return item.id === hookProp.id
  })

  for (let i = index + 1; i < objects.length; i++) {
    if (objects[i].id && objects[i].get('visible') === true) {
      toIndex = i
      break
    }
  }

  console.log('toIndex====>', toIndex)

  if (toIndex === -1) {
    return
  }

  hookProp.moveTo(toIndex)

  canvas.renderAll()



  canvas.fire('object:bringForward', { target: hookProp })
}
// 下移一层
const sendBackwards = () => {
  // 由于有组合元素 需要找到要去的位置  然后通过 .moveTo 移动过去
  const hookProp = canvas.getActiveObject()
  let toIndex = -1

  console.log('bringForward====>', canvas.getObjects(), hookProp)

  const objects = canvas.getObjects().map((item: any) => {
    if (item.name === 'curve-line' && item.idObj) {
      Object.assign(item, item.idObj)
    }
    return item
  })
  const index = objects.findIndex((item: any) => {
    if (hookProp.name === 'curve-line' && hookProp.idObj) {
      Object.assign(hookProp, hookProp.idObj)
    }

    return item.id === hookProp.id
  })

  for (let i = index - 1; i >= 0; i--) {
    if (objects[i].id && objects[i].get('visible') === true) {
      toIndex = i
      break
    }
  }

  if (toIndex === -1) {
    return
  }

  hookProp.moveTo(toIndex)
  canvas.renderAll()

  canvas.fire('object:sendBackwards', { target: hookProp })
}
//截图
const takeScreenshot = () => {
  // 使用箭头函数来确保`this`的上下文正确
  const callback = (blob: Blob) => {
    const imgFile = new File([blob], 'screenshot.png', { type: 'image/png' });
    fileStore.setImgFile(imgFile);
  };
  // 调用toBlob并将回调函数传递给它
  canvasRef.value.toBlob(callback, 'image/jpeg', 0.95);
};
</script>
<style scoped lang="scss">
.second-wrapper {
  position: absolute;
  width: 100vw;
  height: calc(100vh - 100px);
  top: 100px;
  left: 0;
  z-index: 1;

  .workspace-container {
    width: 100%;
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: row;

    .edit-box {
      width: calc(100% - 340px);
      height: 100%;
    }

    .opacity-overlay {
      position: absolute;
      left: 0;
      width: calc(100% - 485px);
      height: 100%;
      z-index: 10000;
    }
  }
}
</style>

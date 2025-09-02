<template>
  <div class="props-bar-wrapper">
    <img src="../../assets/images/props_title.png" width="84" height="36" class="props-title" />
    <div class="props-container">
      <div
        v-for="(item, index) in workspaceStore.propsList"
        :key="index"
        @click="addProp(item)"
        class="item-box"
      >
        <img :src="item.src" :mytype="item.type" alt="" :width="item.width" :height="item.height" />
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onBeforeUnmount } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'
const workspaceStore = useWorkspaceStore()
const emits = defineEmits(['addProp'])
const addProp = (item: any) => {
  emits('addProp', item)
}
// 鼠标点击元素时距离元素左上角的偏移量
const eleOffset = {
  x: 0,
  y: 0
}
// 拖拽开始
const dragstart = (e: any) => {
  const type = e.target?.attributes?.mytype?.nodeValue
  if (type) {
    eleOffset.x = e.offsetX
    eleOffset.y = e.offsetY
  }
}
// 拖拽结束
const dragend = (e: any) => {
  const type = e.target?.attributes?.mytype?.nodeValue
  if (type) {
    const obj = workspaceStore.propsList.filter((item: any) => item.type === type)[0]
    const leftX = e.pageX - eleOffset.x
    const topY = e.pageY - eleOffset.y
    const menuRect = document.querySelector('.props-bar-wrapper')?.getBoundingClientRect()
    const headRect = document.querySelector('.header-wrapper')?.getBoundingClientRect()
    const boxLeft = menuRect?.width || 340
    const boxTop = headRect?.height || 100
    const canvasLeft = leftX - boxLeft
    const canvasTop = topY - boxTop
    const newObj = JSON.parse(JSON.stringify(obj))
    newObj.left = canvasLeft
    newObj.top = canvasTop
    console.log(type, obj, canvasLeft, canvasTop)
    emits('addProp', newObj)
  }
}
document.addEventListener('dragstart', dragstart)
document.addEventListener('dragend', dragend)

onBeforeUnmount(() => {
  document.addEventListener('dragstart', dragstart)
  document.removeEventListener('dragend', dragend)
})
</script>
<style scoped lang="scss">
.props-bar-wrapper {
  top: 100px;
  width: 340px;
  height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 101;
  background-color: #ffffff;
  flex: 0 0 auto;
  overflow: auto;
  &::-webkit-scrollbar {
    display: none;
  }
  .props-title {
    margin-top: 30px;
    margin-left: 39px;
    align-self: flex-start;
  }
  .props-container {
    width: 260px;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: space-between;
    padding-top: 20px;

    .item-box {
      width: 124px;
      height: 124px;
      background-color: #e8efff;
      border-radius: 10px;
      flex: 0 0 auto;
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      cursor: pointer;
      margin-bottom: 20px;
      user-select: none;
    }
  }
}
</style>

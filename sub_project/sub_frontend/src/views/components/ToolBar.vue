<template>
  <div class="tool-bar-wrapper">
    <div class="level">
      <div class="top-level level-btn" @click="bringToFront">
        <img src="../../assets/images/top_icon.png" width="24" height="24" />
        <span>置于顶层</span>
      </div>
      <div class="bottom-level level-btn" @click="sendToBack">
        <img src="../../assets/images/bottom_icon.png" width="24" height="24" />
        <span>置于底层</span>
      </div>
      <div class="up-level level-btn" @click="bringForward">
        <img src="../../assets/images/down_icon.png" width="24" height="24" />
        <span>上移一层</span>
      </div>
      <div class="down-level level-btn" @click="sendBackwards">
        <img src="../../assets/images/up_icon.png" width="24" height="24" />
        <span>下移一层</span>
      </div>
    </div>
    <div class="flip" v-if="workspaceStore?.selectedProp?.type === 'Hook'">
      <div @click="flipX">
        <img src="../../assets/images/horizontal_flip_icon.png" width="24" height="24" />
        <span>水平翻转</span>
      </div>
      <div @click="flipY">
        <img src="../../assets/images/vertical_flip_icon.png" width="24" height="24" />
        <span>垂直翻转</span>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useWorkspaceStore } from '@/stores/workspace'
const workspaceStore = useWorkspaceStore()
const emits = defineEmits(['flipX', 'flipY', 'bringToFront', 'sendToBack', 'bringForward', 'sendBackwards'])
const flipX = () => {
  emits('flipX')
}
const flipY = () => {
  emits('flipY')
}
const bringToFront = () => {
  emits('bringToFront')
}
const sendToBack = () => {
  emits('sendToBack')
}
const bringForward = () => {
  emits('bringForward')
}
const sendBackwards = () => {
  emits('sendBackwards')
}
</script>
<style scoped lang="scss">
.tool-bar-wrapper {
  position: fixed;
  top: 125px;
  height: 68px;
  min-width: 336px;
  border-radius: 53px;
  background: rgba(255, 255, 255, 0.3);
  box-shadow: 0px -1px 0px #ffffff;
  backdrop-filter: blur(2px);
  display: flex;
  flex-direction: row;
  z-index: 1001;
  justify-content: space-evenly;
  align-items: center;
  transform: translateX(-50%);
  left: calc(50% + 170px);
  .level {
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
    align-items: center;
    width: 350px;
    height: 100%;
    .level-btn {
      height: 100%;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      cursor: pointer;
      & > span {
        display: inline-block;
        font-size: 14px;
        line-height: 21px;
        letter-spacing: 1%;
        color: rgba(0, 0, 0, 0.6);
      }
    }
  }
  .flip {
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
    align-items: center;
    width: 196px;
    height: 37px;
    border-left: 1px solid rgba(0, 0, 0, 0.2);
    div {
      height: 100%;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      cursor: pointer;
      & > span {
        display: inline-block;
        font-size: 14px;
        line-height: 21px;
        letter-spacing: 1%;
        color: rgba(0, 0, 0, 0.6);
      }
    }
  }
}
</style>

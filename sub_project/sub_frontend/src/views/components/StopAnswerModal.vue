<template>
  <div class="stop-answer-wrapper">
    <div class="content">
      <div class="title">结束绘制</div>
      <div class="text">绘制时间已结束，点击确认后将进入下一步</div>
      <div class="confirm-btn" @click="confrim">确认</div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { useWorkspaceStore } from '@/stores/workspace'
import { useFileStore } from '@/stores/files'
const workspaceStore = useWorkspaceStore()
const fileStore = useFileStore()
const confrim = () => {
  fileStore.changeItemStep(2)
  workspaceStore.isShowRecord(true)
  workspaceStore.hideToolbar()
  workspaceStore.isShowTimeModal(false)
}
</script>
<style lang="scss" scoped>
.stop-answer-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: #0000006b;
  z-index: 100001;
  .content {
    position: absolute;
    width: 648px;
    height: 272px;
    transform: translate(-50%, -50%);
    top: 50%;
    left: 50%;
    background: #ffffff;
    border-radius: 20px;
    box-sizing: border-box;
    padding: 35px;
    display: flex;
    flex-direction: column;
    .title {
      font-size: 26px;
      font-weight: 500;
      color: #333333;
    }
    .text {
      font-size: 20px;
      line-height: 30px;
      color: #999999;
      margin-top: 28px;
    }
    .confirm-btn {
      width: 118px;
      height: 52px;
      border-radius: 10px;
      font-size: 18px;
      color: #ffffff;
      background: #7a62f5;
      text-align: center;
      line-height: 52px;
      margin-top: 54px;
      margin-left: auto;
      cursor: pointer;
    }
  }
}
</style>

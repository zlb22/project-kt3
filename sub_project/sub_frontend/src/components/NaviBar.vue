<template>
  <div class="header-wrapper">
    <div class="title">
      <p>创造性问题解决</p>
      <div class="icon" @click="showNotice" ref="tipsRef">
        <img src="https://static0.xesimg.com/xpp-parent-fe/1.28.0/info2.png" alt="" />
        <div class="pop-tips" v-if="isShowTips">
          <div class="triangle-icon"></div>
          <div class="tips-content">
            <div class="tips-title">任务提示</div>
            <p>
              欢迎来到创造力的世界，这个世界属于你。
              接下来你将要看到一条宽阔的河流，河对岸有一箱货物。请利用图形元素设计出一个独特的、切实可行的装置将货运送过河。图形元素可以代表任何材质，并且可以变换大小形状，你可以自由创作。
              你有8分钟的时间操作，完成后，你有1分钟来阐述自己的方案设计，最后提交结束任务。
            </p>
          </div>
        </div>
      </div>
    </div>
    <div class="change-bar">
      <div class="step" v-for="(item, index) in items" :key="index">
        <div :class="['step-number', item.selected ? 'step-number-active' : '']">
          {{ item.number }}
        </div>
        <span class="step-content">{{ item.content }}</span>
      </div>
    </div>
    <div class="tail">
      <Countdown class="count-down-wrapper" v-show="!workspaceStore.showRecord"></Countdown>

      <div
        v-if="workspaceStore.step === 2 && !workspaceStore.showRecord"
        class="finsh-btn btn"
        @click="finish"
      >
        完成绘制
      </div>
      <div class="next-btn" v-if="workspaceStore.showRecord && workspaceStore.step === 2">
        <div
          class="btn"
          @click="goback"
          :class="workspaceStore.remainTime === 0 ? 'back-btn-disabled' : 'back-btn'"
        >
          返回绘制
        </div>
        <div
          class="commit-btn btn"
          @click="commit"
          ref="commitRef"
          :style="
            fileStore.commitStatus
              ? 'background-color : rgba(122, 98, 245, 1)'
              : 'background-color : rgba(0, 0, 0, 0.2)'
          "
        >
          提交
        </div>
      </div>
      <div v-if="workspaceStore.step === 3" class="user-info">
        <img src="../assets/images/user.png" alt="" />
        <p>用户{{ fileStore.uid }}</p>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'
import { useFileStore } from '@/stores/files'
import Countdown from './Countdown.vue'

const workspaceStore = useWorkspaceStore()
const emits = defineEmits(['continueCount'])
const fileStore = useFileStore()
console.log(fileStore.uid)
const items = fileStore.items
const tipsRef = ref<any>(null)

const commitRef = ref<HTMLElement | null>(null)
const isShowTips = ref(false)
const showNotice = (e: any) => {
  isShowTips.value = true
}
onMounted(() => {
  window.addEventListener('click', (e) => {
    if (tipsRef.value && !tipsRef.value.contains(e.target as Node)) {
      isShowTips.value = false
    }
  })
})
/**
 * 完成绘制
 */
const finish = () => {
  if (workspaceStore.running) {
    clearInterval(workspaceStore.timeId)
    workspaceStore.setRunning(false)
  }
  workspaceStore.isShowRecord(true)
  workspaceStore.hideToolbar()
  fileStore.changeItemStep(2)
  console.log('finish')
}
const goback = () => {
  console.log('workspaceStore.remainTime', workspaceStore.remainTime)
  if (workspaceStore.remainTime === 0) return
  // 继续计时
  workspaceStore.startCountdown()
  // 隐藏录制组件
  workspaceStore.isShowRecord(false)

  fileStore.changeItemStep(1)
  console.log('goback')
}
const commit = () => {
  if (fileStore.commitStatus === true) {
    fileStore.setConfirmCommitStatus(true)
  }
}
</script>
<style scoped lang="scss">
.btn {
  width: 155px;
  height: 54px;
  border-radius: 10px;
  text-align: center;
  line-height: 54px;
  cursor: pointer;
}
.header-wrapper {
  box-sizing: border-box;
  position: absolute;
  top: 0;
  left: 0;
  padding: 0 38px;
  width: 100vw;
  height: 100px;
  background-color: #ffffff;
  z-index: 1000;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0px 1px 0px rgba(0, 0, 0, 0.05);
  .title {
    display: flex;
    align-items: center;
    font-size: 24px;
    font-weight: 500;
    line-height: 36px;
    text-align: left;
    color: rgba(51, 51, 51, 1);
    .icon {
      width: 30px;
      height: 30px;
      margin-left: 11px;
      cursor: pointer;
      display: flex;
      position: relative;
      img {
        align-items: center;
        width: 100%;
      }
    }
  }
  .change-bar {
    .step {
      display: inline-block;
      line-height: 48px;
      .step-number {
        display: inline-block;
        width: 48px;
        height: 48px;
        border-radius: 50%;
        text-align: center;
        background-color: rgba(0, 0, 0, 0.2);
        color: rgba(255, 255, 255, 1);
        font-size: 24px;
        font-weight: 500;

        &.step-number-active {
          background-color: rgba(122, 98, 245, 1);
        }
      }
      .step-content {
        margin-left: 9px;
        font-size: 24px;
        font-weight: 500;
        line-height: 36px;
      }
    }
    .step:nth-child(2)::after,
    .step:nth-child(1)::after {
      vertical-align: super;
      margin: 0 20px;
      display: inline-block;
      content: '';
      width: 70px;
      border: 1px solid rgba(204, 204, 204, 1);
    }
  }

  .tail {
    display: flex;
    align-items: center;
    .finsh-btn {
      font-size: 24px;
      font-weight: 400;
      color: rgba(255, 255, 255, 1);
      background-color: #7a62f5;
    }
    .next-btn {
      display: flex;
      justify-content: space-around;
      .back-btn {
        color: rgba(122, 98, 245, 1);
        font-size: 24px;
        font-weight: 400;
        background: rgba(122, 98, 245, 0.2);
      }
      .back-btn-disabled {
        color: #ffffff;
        font-size: 24px;
        font-weight: 400;
        background: rgba(0, 0, 0, 0.2);
      }
      .commit-btn {
        margin-left: 10px;
        font-size: 24px;
        font-weight: 400;
        background-color: #00000033;
        color: #ffffff;
      }
    }
    .user-info {
      display: flex;
      align-items: center;
      img {
        width: 30px;
        height: 30px;
      }
      p {
        margin-left: 6px;
        font-family: Source Han Sans CN;
        font-size: 24px;
        font-weight: 500;
        line-height: 36px;
        color: #00000099;
      }
    }
  }
  .pop-tips {
    position: absolute;
    width: 514px;
    height: 308px;
    top: 30px;
    left: -80px;

    box-sizing: border-box;
    .triangle-icon {
      position: relative;
      width: 20px;
      height: 20px;
      transform: rotate(45deg);
      border-left: 1px solid rgba(0, 0, 0, 0.1);
      border-top: 1px solid rgba(0, 0, 0, 0.1);
      border-top-left-radius: 4px;
      left: 85px;
      top: 2px;
      z-index: 99999;
      background-color: #ffffff;
    }
    .tips-content {
      position: relative;
      z-index: 99998;
      margin-top: -9px;
      background-color: #ffffff;
      border: 1px solid rgba(0, 0, 0, 0.1);
      box-shadow: 0px 4px 11.7px rgba(0, 0, 0, 0.05);
      border-radius: 8px;
      padding: 24px 30px;
      .tips-title {
        font-size: 24px;
        font-weight: 500;
        line-height: 36px;
        color: rgba(51, 51, 51, 1);
      }
      & > p {
        margin-top: 11px;
        font-size: 20px;
        line-height: 30px;
        color: rgba(0, 0, 0, 0.4);
        font-weight: 400;
      }
    }
  }
}
</style>

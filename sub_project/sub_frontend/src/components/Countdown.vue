<template>
  <div class="countdown-wrapper">
    <span class="text-label" :style="{ color: workspaceStore.remainTime <= 30 ? 'red' : '#999999' }"
      >倒计时</span
    >
    <div class="time-box">
      <img
        v-if="workspaceStore.remainTime > 30"
        src="../assets/images/clock_icon.png"
        alt=""
        width="25"
        height="25"
      />
      <img v-else src="../assets/images/clock_red_icon.png" alt="" width="25" height="25" />
      <span class="time" :style="{ color: workspaceStore.remainTime <= 30 ? 'red' : '#999999' }">{{
        formattedTime
      }}</span>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'
const workspaceStore = useWorkspaceStore()

onMounted(() => {
  workspaceStore.startCountdown()
})

const formattedTime = computed(() => {
  const minutes = Math.floor(workspaceStore.remainTime / 60)
  const seconds = workspaceStore.remainTime % 60
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})
</script>
<style lang="scss" scoped>
.countdown-wrapper {
  margin-right: 60px;
  
  width: 191px;
  height: 36px;
  background: #ffffff4d;
  box-shadow: 0px -1px 0px 0px #ffffff;
  border-radius: 53px;
  display: flex;
  flex-direction: row;
  align-items: center;
  box-sizing: border-box;
  justify-content: space-between;
  .text-label {
    display: inline-block;
    font-size: 24px;
    font-weight: 500;
    color: rgba(153, 153, 153, 1);
    // color: #999999;
  }
  .time-box {
    display: flex;
    flex-direction: row;
    width: 97px;
    height: 100%;
    align-items: center;
    .time {
      display: inline-block;
      font-size: 24px;
      letter-spacing: 1%;
      font-weight: 500;
      margin-left: 4px;
      color: rgba(153, 153, 153, 1);
    }
    img {
      width: 25px;
      height: 25px;
    }
  }
}
</style>

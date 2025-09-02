<template>
  <div class="main-wrapper">
    <NaviBar v-if="workspaceStore.step !== 1" />
    <FirstStep v-if="workspaceStore.step === 1" />
    <SecondStep v-if="workspaceStore.step === 2" />
    <ThirdStep v-if="workspaceStore.step === 3" />
    <StopAnswerModal
      v-if="workspaceStore.remainTime === 0 && workspaceStore.showTimeModal"
    ></StopAnswerModal>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import NaviBar from '@/components/NaviBar.vue'
import FirstStep from './FirstStep.vue'
import SecondStep from './SecondStep.vue'
import ThirdStep from './ThirdStep.vue'
import { useWorkspaceStore } from '@/stores/workspace'
import StopAnswerModal from './components/StopAnswerModal.vue'

const workspaceStore = useWorkspaceStore()

const handleBeforeUnload = (event: any) => {
  const message = '你确定要离开吗？';
  event.returnValue = message; // 兼容某些浏览器
  return message;
}

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload);
})
</script>
<style scoped lang="scss">
.main-wrapper {
  position: relative;
  width: 100vw;
  height: 100vh;
  background-color: white;
}
</style>

<template>
  <div class="record-wrapper">
    <p>请用文字描述你的解决方案（建议 50-300 字）</p>
    <div class="text-panel">
      <textarea
        v-model="text"
        class="text-input"
        :maxlength="maxLen"
        placeholder="请输入你的思路阐述..."
        @input="onInput"
      ></textarea>
      <div class="text-meta">
        <span :class="{ warn: remain < 0 }">{{ text.length }}/{{ maxLen }}</span>
      </div>
    </div>
  </div>
  
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useFileStore } from '../../stores/files'

const store = useFileStore()
const maxLen = 1000
const text = ref<string>(store.voiceText || '')
const remain = computed(() => maxLen - text.value.length)

const onInput = () => {
  store.setVoiceText(text.value)
}

watch(
  () => store.voiceText,
  (val) => {
    if (val !== text.value) text.value = val || ''
  }
)
</script>

<style scoped lang="scss">
.record-wrapper {
  position: absolute;
  width: 485px;
  height: calc(100vh - 100px);
  top: 0;
  right: 0;
  background-color: #fff;
  & > p {
    position: absolute;
    top: 24px;
    left: 24px;
    font-size: 18px;
    font-weight: 500;
    line-height: 24px;
  }
  .text-panel {
    position: absolute;
    top: 64px;
    left: 24px;
    right: 24px;
    bottom: 24px;
    display: flex;
    flex-direction: column;
  }
  .text-input {
    flex: 1;
    width: 100%;
    resize: none;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 12px;
    font-size: 14px;
    line-height: 1.6;
    outline: none;
  }
  .text-meta {
    display: flex;
    justify-content: flex-end;
    margin-top: 8px;
    color: #999;
    .warn { color: #e03231; }
  }
}
</style>

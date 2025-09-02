<template>
  <div class="forward-backward-control">
    <div :class="['backward-icon', backwardAble ? '' : 'backward-icon-disable']">
      <img src="../../assets/images/backward-icon.png" alt="backward" @click="backwardHandle" />
    </div>
    <div :class="['forward-icon', forwardAble ? '' : 'forward-icon-disable']">
      <img src="../../assets/images/forward-icon.png" alt="forward" @click="forwardHandle" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useOperationHistory } from '../../stores/operationHistory'

const operationHistory = useOperationHistory()

const forwardAble = computed(() => operationHistory.historyStackTemp.length > 0)

const backwardAble = computed(() => operationHistory.historyStack.length > 0)

const backwardHandle = () => {
  console.log('backwardHandle')
  operationHistory.backwardHanle()
}

const forwardHandle = () => {
  if (forwardAble.value) {
    console.log('forwardHandle')
    operationHistory.forwardHandle()
  }
}
</script>

<style lang="scss" scoped>
.forward-backward-control {
  position: fixed;
  z-index: 2;
  top: 130px;
  right: 40px;
  width: 225px;
  height: 58px;
  border-radius: 29px;
  box-shadow: 0px -1px 0px 0px rgba(255, 255, 255, 1);
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  justify-content: space-between;
  user-select: none;

  .backward-icon,
  .forward-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 50%;
    position: relative;
    img {
      width: 32px;
      height: 32px;
      cursor: pointer;
    }

    &.backward-icon-disable {
      &::after {
        content: '';
        display: block;
        position: absolute;
        z-index: 10;
        border-radius: 29px 0 0 29px;
        width: 100%;
        height: 100%;
        background-color: rgba(251, 251, 251, 0.8);
        pointer-events: none;
      }
    }

    &.forward-icon-disable {
      &::after {
        content: '';
        display: block;
        position: absolute;
        z-index: 10;
        border-radius: 0 29px 29px 0;
        width: 100%;
        height: 100%;
        background-color: rgba(251, 251, 251, 0.8);
        pointer-events: none;
      }
    }
  }

  .backward-icon::after {
    content: '';
    display: block;
    height: 30px;
    width: 1px;
    background: rgba(0, 0, 0, 0.2);
    position: absolute;
    right: 0;
  }
}
</style>

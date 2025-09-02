<template>
  <div class="third-wrapper" ref="thirdPage">
    <div class="logo" :style="logoStyle">
      <img v-if="prograss!==100" src="../assets/images/uploading.png" :width="logoWidth" :height="logoHeight" alt="">
      <img v-else src="../assets/images/uploaded.png" :width="logoWidth" :height="logoHeight" alt="">
    </div>
    <p v-if="prograss!==100" class="desc" :style="descStyle">上传中 {{ prograss }}%，请勿关闭..</p>
    <p v-else class="desc" :style="descStyle">恭喜完成！</p>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import{ config} from '@/utils/update'
const prograss = computed(() => config.authProgress)
const thirdPage= ref<any>(null)
let descWith = ref(0)
let logoWidth = ref(0)
let logoHeight = ref(0)
let descStyle: any = ref(null)
let logoStyle = ref({})
onMounted(() => {
  const height = thirdPage.value.clientHeight
  const width = thirdPage.value.clientWidth
  // 屏幕适配比例
  const ratio = height / 1080
  logoWidth.value = 198 * ratio
  logoHeight.value = 198 * ratio
  descWith.value = 300 * ratio
  logoStyle.value = {
    top: Math.round(420 * ratio) + 'px',
    left: Math.round((width - logoWidth.value) / 2) + 'px'
  }
  descStyle.value = {
    top: Math.round(620 * ratio) + 'px',
    left: Math.round((width - descWith.value) / 2) + 'px',
    fontSize: 27 * ratio + 'px',
    lineHeight: 41 * ratio + 'px',
    width:350 * ratio + 'px'
  }
})
</script>
<style scoped lang="scss">
.third-wrapper {
  position: absolute;
  width: 100vw;
  height: calc(100vh - 100px);
  top: 100px;
  left: 0;
  background-color: rgba(228, 235, 253, 1);
  .logo{
    display: flex;
    flex-direction: column;
    position: absolute;
    img{
      position: absolute;
    }
  }
  .desc{
    position: absolute;
    font-weight: 400;
    text-align: center;
    color: rgba(153, 153, 153, 1);
  }
}
</style>

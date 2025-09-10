<template>
  <div class="start-view">
    <img class="bottom-bg" src="../assets/images/bgBottom.png" width="1920" height="300" />

    <div class="welcome-view" ref="startPage" v-if="!startVedioInfo">
      <div class="logo" :style="logoStyle">
        <img src="../assets/images/logo.png" :width="logoWidth" :height="logoHeight" />
      </div>
      <div class="content" ref="textContainer" :style="textStyle">
        <p class="title" :style="titleStyle">创造性问题解决</p>
        <!-- 更新文案 -->
        <p class="desc" :style="descStyle">
          欢迎来到创造力的世界，这个世界属于你。
          <br />
          接下来你将要看到一条宽阔的河流，河对岸有一箱货物。请利用图形元素设计出一个独特的、切实可行
          <br />
          的装置将货运送过河。图形元素可以代表任何材质，并且可以变换大小形状，你可以自由创作。
          <br />
          你有8分钟的时间操作，完成后，你有1分钟来阐述自己的方案设计，最后提交结束任务。
        </p>
      </div>
      <div class="start-btn" :style="startBtnStyle" @click="start">开始</div>
    </div>

    <div class="vedio-view" v-if="vedioShow">
      <video
        class="vedio"
        :style="vedioStyle"
        src="https://static0.xesimg.com/xpp-parent-fe/topic-three-web/vedio.mp4"
        autoplay
        @ended="vedioEnd"
      ></video>
      <!-- 跳过按钮 -->
      <div
        class="skip-button"
        :style="{ position: 'absolute', background: 'rgba(255, 255, 255, 0.9)', color: '#333', top: '50px', right: '50px', width: '120px', height: '50px', fontSize: '18px', borderRadius: '25px', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', zIndex: 10000, border: '2px solid #7A62F5' }"
        @click="skipVideo"
      >
        跳过视频
      </div>
      <div
        class="vedio-button"
        :style="{ ...startBtnStyle, background: !btnAbled ? 'gray' : 'rgba(122, 98, 245, 1)' }"
        @click="nextStep"
      >
        已知悉
      </div>
    </div>

    <!-- Manual info input disabled due to unified login; keep block hidden -->
    <div class="info-input-view" v-if="false">
      <div class="input-container" :style="inputContainerStyle">
        <div class="input-title" :style="inputTitleStyle">请输入你的信息</div>
        <ul>
          <li :style="inputLiStyle">
            学校：
            <div class="inputTextBox" :style="inputTextBoxStyle">
              <input
                type="text"
                placeholder="请输入学校名称"
                :style="inputTextStyle"
                v-model="userInfo.school"
              />
            </div>
          </li>
          <li :style="inputLiStyle">
            姓名：
            <div class="inputTextBox" :style="inputTextBoxStyle">
              <input
                type="text"
                placeholder="请输入姓名"
                :style="inputTextStyle"
                v-model="userInfo.name"
              />
            </div>
          </li>
          <li :style="inputLiStyle">
            学号：
            <div class="inputTextBox" :style="inputTextBoxStyle">
              <input
                type="text"
                placeholder="请输入学号"
                :style="inputTextStyle"
                v-model="userInfo.num"
              />
            </div>
          </li>
          <li :style="inputLiStyle">
            年级：
            <div class="inputTextBox" :style="inputTextBoxStyle">
              <input
                type="text"
                placeholder="请输入年级"
                :style="inputTextStyle"
                v-model="userInfo.grade"
              />
            </div>
          </li>
          <li :style="inputLiStyle">
            班级：
            <div class="inputTextBox" :style="inputTextBoxStyle">
              <input
                type="text"
                placeholder="请输入班级"
                :style="inputTextStyle"
                v-model="userInfo.class"
              />
            </div>
          </li>
        </ul>
      </div>
      <div
        class="confirm-btn"
        :style="{ ...startBtnStyle, background: comfirmStatus ? 'rgba(122, 98, 245, 1)' : 'grey' }"
        @click="confirm"
      >
        确认
      </div>
    </div>
  </div>
</template>
s
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { apiInstance } from '@/utils/axios'
import { useWorkspaceStore } from '@/stores/workspace'
import { useFileStore } from '@/stores/files'
import { getUrlParam } from '@/utils/url'
import { login } from '@/api/login'
const workspaceStore = useWorkspaceStore()
const fileStore = useFileStore()
interface School {
  id: number
  name: string
}
interface Grade {
  id: number
  name: string
}
interface List {
  school: School[]
  grade: Grade[]
}
const list = computed(() => fileStore.list as List)

const userInfo = ref({
  school: '',
  name: '',
  num: '',
  grade: '',
  class: ''
})
const comfirmStatus = computed(
  () =>
    userInfo.value.school !== '' &&
    userInfo.value.name !== '' &&
    userInfo.value.num !== '' &&
    userInfo.value.grade !== '' &&
    userInfo.value.class !== ''
)
let logoWidth = ref(0)
let logoHeight = ref(0)
let logoStyle: any = ref(null)
let textStyle: any = ref(null)
let titleStyle: any = ref(null)
let descStyle: any = ref(null)
let startBtnStyle: any = ref(null)
let inputTitleStyle: any = ref(null)
let inputContainerStyle: any = ref(null)
let inputLiStyle: any = ref(null)
let inputTextBoxStyle = ref<any>(null)
let inputTextStyle = ref<any>(null)
let selectStyle = ref<any>(null)
let vedioStyle = ref<any>(null)

const startPage = ref<any>(null)
const textContainer = ref<any>(null)

const startVedioInfo = ref(false)
const vedioShow = ref(false) //切换信息输入模式--的状态，用于管理页面
const startInputInfo = ref(false)
const btnAbled = ref(false)
onMounted(() => {
  const height = startPage.value.clientHeight
  const width = startPage.value.clientWidth
  // 屏幕适配比例
  const ratio = height / 1080
  // logo样式
  logoWidth.value = 260 * ratio
  logoHeight.value = 186 * ratio
  logoStyle.value = {
    top: Math.round(213 * ratio) + 'px',
    left: Math.round((width - logoWidth.value) / 2) + 'px'
  }
  textStyle.value = {
    top: Math.round(439 * ratio) + 'px'
  }
  titleStyle.value = {
    fontSize: Math.round(48 * ratio) + 'px',
    lineHeight: Math.round(72 * ratio) + 'px'
  }
  descStyle.value = {
    fontSize: Math.floor(32 * ratio) + 'px',
    lineHeight: Math.floor(48 * ratio) + 'px',
    width: Math.floor(1440 * ratio) + 'px'
  }
  startBtnStyle.value = {
    width: Math.round(556 * ratio) + 'px',
    height: Math.round(88 * ratio) + 'px',
    bottom: Math.round(99 * ratio) + 'px',
    left: Math.round((width - 556 * ratio) / 2) + 'px',
    lineHeight: Math.round(88 * ratio) + 'px',
    fontSize: Math.round(24 * ratio) + 'px',
    borderRadius: Math.round((88 * ratio) / 2) + 'px'
  }
  inputTitleStyle.value = {
    fontSize: Math.round(32 * ratio) + 'px',
    lineHeight: Math.round(36 * ratio) + 'px',
    marginBottom: Math.round(60 * ratio) + 'px'
  }
  inputContainerStyle.value = {
    // width: Math.round(400 * ratio) + 'px',
    // height: Math.round(76 * ratio) + 'px',
    top: Math.round(119 * ratio) + 'px',
    left: Math.round((width - 536 * ratio) / 2) + 'px'
  }
  inputLiStyle.value = {
    width: Math.round(536 * ratio) + 'px',
    fontSize: Math.round(32 * ratio) + 'px',
    lineHeight: Math.round(36 * ratio) + 'px',
    marginBottom: Math.round(31 * ratio) + 'px'
  }
  inputTextBoxStyle.value = {
    width: Math.ceil(400 * ratio) + 'px',
    height: Math.ceil(76 * ratio) + 'px',
    marginLeft: Math.round(24 * ratio) + 'px',
    borderRadius: Math.round(4 * ratio) + 'px'
  }
  inputTextStyle.value = {
    width: Math.round(334 * ratio) + 'px',
    height: Math.round(76 * ratio) + 'px',
    fontSize: Math.round(24 * ratio) + 'px',
    lineHeight: Math.round(36 * ratio) + 'px',
    paddingLeft: Math.round(40 * ratio) + 'px',
    paddingRight: Math.round(26 * ratio) + 'px'
  }
  selectStyle.value = {
    width: Math.round(336 * ratio) + 'px',
    height: Math.round(76 * ratio) + 'px',
    fontSize: Math.round(24 * ratio) + 'px',
    lineHeight: Math.round(36 * ratio) + 'px',
    paddingLeft: Math.round(40 * ratio) + 'px',
    paddingRight: Math.round(26 * ratio) + 'px'
  }
  vedioStyle.value = {
    width: Math.round(1280 * ratio) + 'px',
    height: Math.round(448 * ratio) + 'px',
    top: Math.round(236 * ratio) + 'px'
  }
})
const start = () => {
  startVedioInfo.value = true
  vedioShow.value = true
  //debugger参数可以直接点击跳转
  if (getUrlParam('debugger') === 'true') {
    btnAbled.value = true
  }
}
const vedioEnd = () => {
  btnAbled.value = true
}
const skipVideo = () => {
  btnAbled.value = true
}
const nextStep = () => {
  if (!btnAbled.value) return
  vedioShow.value = false
  // Directly enter step 2 since user info is taken from unified login
  workspaceStore.setStep(2)
}
// On mount, auto-fetch unified user info and login to acquire uid
onMounted(async () => {
  try {
    const me = await apiInstance.get('/api/auth/me')
    const body = me?.data
    if (body && body.username) {
      const loginData = {
        username: body.username,
        school: body.school || '',
        grade: body.grade || ''
      }
      try {
        const result = await login(loginData)
        fileStore.setUid(result.data.data.uid)
      } catch (e) {
        console.error('keti3 student login failed:', e)
      }
    }
  } catch (e) {
    console.error('fetch /api/auth/me failed:', e)
  }
})
</script>
<style scoped lang="scss">
.start-view {
  width: 100vw;
  height: 100vh;
  background-color: rgba(228, 235, 253, 1);

  .bottom-bg {
    position: absolute;
    width: 100%;
    bottom: 0;
    z-index: 10;
  }

  .welcome-view {
    width: 100%;
    height: 100%;
    position: relative;

    .logo {
      position: absolute;
    }

    .content {
      position: absolute;
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      z-index: 20;

      .title {
        font-weight: 700;
      }

      .desc {
        font-family: Source Han Sans CN;
        margin-top: 23px;
        font-weight: 500;
        text-align: center;
        color: rgba(122, 98, 245, 1);
      }
    }

    .start-btn {
      position: absolute;
      background: rgba(122, 98, 245, 1);
      text-align: center;
      color: #ffffff;
      cursor: pointer;
      z-index: 100;
    }
  }

  .info-input-view {
    .input-container {
      position: absolute;
      text-align: center;
      z-index: 1000;

      .input-title {
        font-family: Microsoft YaHei;
        font-size: 32px;
        font-weight: 400;
        line-height: 36px;
        color: rgba(122, 98, 245, 1);
      }

      ul {
        width: 100%;
        list-style: none;
        display: flex;
        flex-direction: column;

        li {
          font-family: Microsoft YaHei;
          font-weight: 400;
          text-align: left;
          color: rgba(34, 34, 34, 1);
          display: flex;
          align-items: center;

          .inputTextBox {
            background-color: rgb(255, 255, 255, 1);
            overflow: hidden;

            input,
            select {
              text-align: left;
              font-family: Microsoft YaHei;
              font-weight: 400;
              color: rgba(153, 153, 153, 1);
              border: none;
            }
          }
        }
      }
    }

    .confirm-btn {
      position: absolute;
      background: rgba(122, 98, 245, 1);
      text-align: center;
      color: #ffffff;
      cursor: pointer;
      z-index: 100;
    }
  }

  .vedio-view {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;

    .vedio {
      position: absolute;
      z-index: 9999;
    }

    .vedio-button {
      position: absolute;
      z-index: 9999;
      text-align: center;
      color: #ffffff;
      cursor: pointer;
    }
  }
}
</style>

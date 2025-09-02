<template>
  <div class="record-wrapper">
    <p>请一分钟描述你的解决方案</p>
    <audio ref="audioRef" :src="audioBlobUrl" @timeupdate="timeUpdate" @ended="onEnded"></audio>
    <div class="record-start" v-if="recordStart">
      <div class="record-start-btn btn" @click="startRecord">
        <img src="../../assets/images/record.png" alt="" />
        <p>点击开始录制</p>
      </div>
    </div>
    <div class="recording" v-if="recording">
      <div class="recording-btn btn" @click="stopRecord">
        <img src="../../assets/images/recording.png" alt="" />
        <p class="timer" :style="timer <= 10 ? { color: 'rgba(224, 50, 49, 1)' } : {}">
          倒计时 {{ formatTime(timer) }}
        </p>
      </div>
      <Music class="music"></Music>
    </div>
    <div class="record-end" v-if="recordEnd">
      <p>录制成功！</p>
      <div class="btns">
        <div class="play-btn btn" @click="controlAudio">
          <img src="../../assets/images/play.png" alt="" />
          <p>开始播放</p>
        </div>
        <div class="re-record-btn btn" @click="reRecord">
          <img src="../../assets/images/record.png" alt="" />
          <p>重新录制</p>
        </div>
      </div>
    </div>
    <div class="play" v-if="play">
      <Music class="music"></Music>
      <p class="current-time">{{ currentTime }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import Music from './Music.vue'
import { useFileStore } from '../../stores/files'
const recordStart = ref(true)
const recording = ref(false)
const recordEnd = ref(false)
const play = ref(false)

const mediaStream = ref<MediaStream | null>(null) // 用于存储用户媒体流
const audioChunks = ref<Blob[]>([]) // 用于存储录制的音频片段 ArrayBuffer
const mediaRecorder = ref<MediaRecorder | null>(null) // 用于存储录制器实例
const audioBlobUrl = ref('') // 用于存储音频的URL
const timer = ref<number>(60) // 用于存储计时器
let countdownInterval: number | null = null // 用于存储计时器

const audioRef = ref<HTMLAudioElement | null>(null)
const currentTime = ref<string>('00:00')
// 获取用户媒体设备
const getUserMedia = (constraints: MediaStreamConstraints): Promise<MediaStream> => {
  if (navigator.mediaDevices === undefined) {
    ;(navigator as any).mediaDevices = {}
  }
  const legacyGetUserMedia =
    (navigator as any).getUserMedia ||
    (navigator as any).webkitGetUserMedia ||
    (navigator as any).mozGetUserMedia
  if (!navigator.mediaDevices.getUserMedia && legacyGetUserMedia) {
    navigator.mediaDevices.getUserMedia = function (
      constraints: MediaStreamConstraints
    ): Promise<MediaStream> {
      return new Promise((resolve, reject) => {
        legacyGetUserMedia.call(navigator, constraints, resolve, reject)
      })
    }
  }
  return navigator.mediaDevices.getUserMedia(constraints)
}
//开始录音
const startRecord = async () => {
  try {
    mediaStream.value = await getUserMedia({ audio: true })
    mediaRecorder.value = new MediaRecorder(mediaStream.value)
    audioChunks.value = []
    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
      playAudio()
    }
    mediaRecorder.value.start()
    startCountdown()
    recording.value = true
    recordStart.value = false
    useFileStore().setCommitStatus(false)
    console.log(useFileStore().commitStatus)
  } catch (error) {
    console.error('获取麦克风权限失败:', error)
  }
}
//停止录音
const stopRecord = () => {
  if (mediaRecorder.value) {
    mediaRecorder.value.stop()
    stopCountdown()
    recordEnd.value = true
    recording.value = false
  } else {
    console.warn('mediaRecorder.value is null, cannot call stop method.')
  }
}
//开始倒计时
const startCountdown = () => {
  countdownInterval = window.setInterval(() => {
    if (timer.value > 0) {
      timer.value--
    } else {
      // 当倒计时结束时自动停止录音
      stopRecord()
      stopCountdown()
    }
  }, 1000)
}
//停止倒计时
const stopCountdown = () => {
  if (countdownInterval !== null) {
    window.clearInterval(countdownInterval)
    countdownInterval = null
  }
}
//重新录制
const reRecord = () => {
  audioChunks.value = []
  audioBlobUrl.value = ''
  timer.value = 60
  startRecord()
  recordStart.value = true
  recordEnd.value = false
  useFileStore().setCommitStatus(false)
}
//生成音频
const playAudio = () => {
  if (audioChunks.value.length > 0) {
    const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
    audioBlobUrl.value = URL.createObjectURL(audioBlob)
    const audioFile = new File([audioBlob], 'audio.webm', { type: 'audio/webm' })
    useFileStore().setAudioFile(audioFile)
    useFileStore().setCommitStatus(true)
  } else {
    console.warn('没有录音数据可播放')
  }
}

const controlAudio = () => {
  audioRef.value?.play()
  play.value = true
  recordEnd.value = false
}

// 播放结束改变状态
const onEnded = () => {
  currentTime.value = '00:00'
  recordEnd.value = true
  play.value = false
}
// 播放时间更新t
const timeUpdate = () => {
  if (audioRef.value) {
    currentTime.value = formatTime(audioRef.value.currentTime)
  }
}
const formatTime = (seconds: number): string => {
  // 假设我们希望将秒数转换为 "分钟:秒" 的格式
  const minutes = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${minutes < 10 ? '0' : ''}${minutes}:${secs < 10 ? '0' : ''}${secs}`
}
</script>
'
<style scoped lang="scss">
//声明一个按钮样式
.btn {
  width: 100px;
  height: 100px;
  z-index: 2;
  img {
    width: 100%;
  }
  p {
    margin-top: 21px;
  }
}
p {
  text-align: center;
}
.record-wrapper {
  position: absolute;
  width: 485px;
  height: calc(100vh - 100px);
  top: 0;
  right: 0;
  background-color: rgba(255, 255, 255, 1);
  & > p {
    position: absolute;
    top: 39px;
    left: calc(50% - 144px);
    font-size: 24px;
    font-weight: 500;
    line-height: 36px;
    text-align: left;
  }
  .record-start {
    position: absolute;
    top: 416px;
    left: calc(50% - 50px);
    height: 148px;
    width: 100px;
    & > p {
      font-size: 16px;
      font-weight: 500;
      line-height: 24px;
      text-align: left;
    }
  }
  .recording {
    position: absolute;
    top: 416px;
    left: calc(50% - 50px);
    height: 148px;
    width: 100px;
    & > .music {
      pointer-events: none;
      position: absolute;
      top: 16px;
      left: calc(50% - 200px);
    }
    .timer {
      font-family: Source Han Sans CN;
      font-size: 14px;
      font-weight: 400;
      line-height: 21px;
      text-align: left;
      color: rgba(122, 98, 245, 1);
      position: relative;
      left: calc(50% - 40.5px);
    }
  }
  .record-end {
    & > p {
      position: absolute;
      top: 292px;
      left: calc(50% - 60px);
      font-family: Source Han Sans CN;
      font-size: 24px;
      font-weight: 500;
      line-height: 36px;
      text-align: left;
      color: rgba(122, 98, 245, 0.6);
    }
    .btns {
      width: 260px;
      position: absolute;
      display: flex;
      justify-content: space-between;
      top: 416px;
      left: calc(50% - 130px);
    }
  }
  .play {
    position: absolute;
    top: 416px;
    left: calc(50% - 148px);
    & > .music {
      width: 295.7px;
    }
    & > .current-time {
      position: absolute;

      left: calc(50% - 17.5px);
      font-family: Source Han Sans CN;
      font-size: 14px;
      font-weight: 400;
      line-height: 21px;
      text-align: left;
      color: rgba(122, 98, 245, 1);
    }
  }
}
</style>

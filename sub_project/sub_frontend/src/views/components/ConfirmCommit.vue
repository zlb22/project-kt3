<template>
    <div class="confirm-commit-wrapper">
      <div class="content">
        <div class="title">提交确认</div>
        <div class="text">提交后将上传本次方案，上传后不可修改，确认要提交吗？</div>
        <div class="btns">
          <div class="cancel-btn" @click="cancel">取消</div>
          <div class="confirm-btn" @click="confirm">确认</div>
        </div>
      </div>
    </div>
</template>

<script lang="ts" setup>
import { useFileStore } from '@/stores/files'
import { upload, setConfig} from '@/utils/update'
import { useOperationHistory } from '@/stores/operationHistory'
import { useWorkspaceStore } from '@/stores/workspace'
import axios from 'axios'
import { saveData } from '@/api'
const operationHistory = useOperationHistory()
const workspaceStore = useWorkspaceStore()
const fileStore = useFileStore()
const confirm = () => {
  try {
    upload({ img: fileStore.imgFile, audio: fileStore.audioFile }, 'test', axios).then(
      (res: any) => {
        // const origin=window.location.origin
        saveData({
        uid: fileStore.uid,
        op_type: 'submit',
        voice_url: res.audioUrl,
        screenshot_url: res.imgUrl,
        op_time: new Date().toISOString(),
        data_after: operationHistory.operationLog
      }).then(result=>{
        if(result.data.errcode===0){
          setConfig(100)
        }
      }).catch(err=>{
        console.log(err)
      })
      }
    )
    workspaceStore.setStep(3)
    fileStore.changeItemStep(3)
  } catch (e) {
    console.log(e)
  }
}
const cancel=()=>{
  fileStore.setConfirmCommitStatus(false)
}
</script>


<style lang="scss" scoped>
.confirm-commit-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: #0000006b;
  z-index: 1000001;
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
    .btns{
      display: flex;
      flex-direction: row;
      margin-top: 54px;
      margin-left: auto;
      .confirm-btn {
      width: 118px;
      height: 52px;
      border-radius: 10px;
      font-size: 18px;
      color: #ffffff;
      background: #7a62f5;
      text-align: center;
      line-height: 52px;
      margin-left: 21px;
      cursor: pointer;
    }
    .cancel-btn {
      width: 118px;
      height: 52px;
      border-radius: 10px;
      font-size: 18px;
      color: rgba(51, 51, 51, 1);
      border: 1px solid rgba(224, 224, 224, 1);
      background: rgba(255, 255, 255, 1);
      text-align: center;
      line-height: 52px;
      cursor: pointer;
    }
    }
  }
}
</style>
  
import type { AxiosResponse } from 'axios'
import axios from 'axios'
import { reactive, ref } from 'vue'
import { useFileStore } from '@/stores/files';
import { uploadToServer } from '@/api'
//地区源
export enum EnumRegion {
  HangZhou = 'oss-cn-xxx',
  ShangHai = 'oss-cn-xxx',
  QingDao = 'oss-cn-xxx',
  BeiJing = 'oss-cn-xxx',
  ShenZhen = 'oss-cn-xxx',
  HongKong = 'oss-cn-xxx'
}
// 上传数据类型
export interface PresignedItem {
  key: string
  url: string
  method: 'PUT'
  content_type: string
  public_url: string
}

export interface UploadData {
  Bucket: string
  HttpsDomain: string
  Presigned: PresignedItem[]
  AccessKeyId: string
  AccessKeySecret: string
  SecurityToken: string
  Expiration: string
}
// 返回数据类型
export interface backData<t> {
  errcode: number
  errmsg: string
  data: t
  trace: string
}
// 上传参数类型
export interface dataType {
  parallel: number
  partSize: number
  authProgress: number
  imgProgress: number
  audioProgress: number
  imgFlag: boolean
}
// axios发送器类型
export interface AxiosSender {
  post: <T = backData<UploadData>, R = AxiosResponse<T>, D = any>(url: string) => Promise<R>
}
// 上传参数
export const config = reactive<dataType>({
  parallel: 4,
  partSize: 1024 * 1024,
  authProgress: 0,
  imgProgress: 0,
  audioProgress: 0,
  imgFlag: false
})
export const setConfig = (data: number) => {
  config.authProgress = data
}
export async function upload (e: { img:File , audio:File}, uploadDir: string, axiosSender: AxiosSender) {
  try {
    const form = new FormData()
    if (e.img) form.append('img', e.img)
    if (e.audio) form.append('audio', e.audio)

    const resp = await uploadToServer(form, (ev) => {
      const p = ev.total ? ev.loaded / ev.total : 0
      // 这里是整体表单的进度，按两段各 49% 比例近似映射
      // 为简单起见，整体进度映射到 authProgress，前端展示即可
      config.authProgress = Math.min(98, Math.floor(p * 98))
    })

    const body = resp.data as any
    const ok = resp.status === 200 && (
      (typeof body?.errcode !== 'undefined' ? body.errcode === 0 : body?.code === 0)
    )
    if (!ok) {
      // 打印一条简要日志便于定位
      console.log('upload response:', resp.status, body)
      const msg = body?.errmsg || body?.message || '上传失败'
      throw new Error(msg)
    }

    const payload = body?.data || body?.result || {}
    const { imgUrl, audioUrl } = payload
    if (!imgUrl && !audioUrl) throw new Error('后端未返回文件 URL')
    return { imgUrl, audioUrl }
  } catch (err) {
    console.log(err)
    throw err
  }
}
// 保留占位符，若未来需要切回阿里云 OSS，可在此处实现 SDK 上传版本

import type { AxiosResponse } from 'axios'
import axios from 'axios'
import { reactive, ref } from 'vue'
import { useFileStore } from '@/stores/files';
import { uploadToServerWithUid, getOssAuth } from '@/api'
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
  // 直传优先：调用 /web/keti3/oss/auth 获取预签名 URL，然后使用 PUT 上传；失败则回退到代理上传
  const tryDirect = async () => {
    // 1) 申请预签名 URL（img/audio 可选）
    const content_types: Record<string, string> = {}
    if (e.img) content_types.img = e.img.type || 'image/*'
    if (e.audio) content_types.audio = e.audio.type || 'audio/*'

    const authResp = await getOssAuth({ content_types })
    const body = authResp.data as any
    const ok = authResp.status === 200 && (typeof body?.errcode !== 'undefined' ? body.errcode === 0 : body?.code === 0)
    if (!ok) {
      const msg = body?.errmsg || body?.message || '获取直传授权失败'
      throw new Error(msg)
    }

    const data = body.data || {}
    const presigned: PresignedItem[] = data.Presigned || []
    // 按 key 匹配 img/audio
    const imgItem = presigned.find(p => p.key?.includes('/img'))
    const audioItem = presigned.find(p => p.key?.includes('/audio'))

    // 2) 逐项 PUT 上传（带进度近似映射）
    const doPut = async (item: PresignedItem | undefined, file: File | undefined, onProgress: (p:number)=>void) => {
      if (!item || !file) return undefined
      const total = file.size || 1
      // fetch 不直接提供进度，这里用分片近似上报，避免引入额外依赖
      // 简化：一次性 PUT，不做真正分片，仅在开始/结束两次映射进度
      onProgress(0.01)
      const resp = await fetch(item.url, {
        method: 'PUT',
        headers: { 'Content-Type': file.type || item.content_type || 'application/octet-stream' },
        body: file
      })
      if (!(resp.status >= 200 && resp.status < 300)) {
        throw new Error(`PUT failed: ${resp.status}`)
      }
      onProgress(1)
      return item.public_url
    }

    let progressed = 0
    const mapProgress = (delta:number) => {
      progressed = Math.min(1, progressed + delta)
      config.authProgress = Math.min(98, Math.floor(progressed * 98))
    }

    const imgUrl = await doPut(imgItem, e.img, (p) => mapProgress((p - 0) * 0.5))
    const audioUrl = await doPut(audioItem, e.audio, (p) => mapProgress((p - 0) * 0.5))

    return { imgUrl, audioUrl }
  }

  try {
    return await tryDirect()
  } catch (err) {
    console.warn('Direct upload failed, fallback to proxy:', err)
    try {
      const form = new FormData()
      if (e.img) form.append('img', e.img)
      if (e.audio) form.append('audio', e.audio)
      const fileStore = useFileStore()
      const resp = await uploadToServerWithUid(form, fileStore.uid, (ev) => {
        const p = ev.total ? ev.loaded / ev.total : 0
        config.authProgress = Math.min(98, Math.floor(p * 98))
      })
      const body = resp.data as any
      const ok = resp.status === 200 && (typeof body?.errcode !== 'undefined' ? body.errcode === 0 : body?.code === 0)
      if (!ok) {
        console.log('upload response:', resp.status, body)
        const msg = body?.errmsg || body?.message || '上传失败'
        throw new Error(msg)
      }
      const payload = body?.data || body?.result || {}
      const { imgUrl, audioUrl } = payload
      if (!imgUrl && !audioUrl) throw new Error('后端未返回文件 URL')
      return { imgUrl, audioUrl }
    } catch (fallbackErr) {
      console.log(fallbackErr)
      throw fallbackErr
    }
  }
}
// 保留占位符，若未来需要切回阿里云 OSS，可在此处实现 SDK 上传版本

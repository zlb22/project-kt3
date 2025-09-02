import OSS from 'ali-oss'
import type { AxiosResponse } from 'axios'
import axios from 'axios'
import { reactive, ref } from 'vue'
import { useFileStore } from '@/stores/files';
import { getOssAuth } from '@/api'
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
export interface UploadData {
  AccessKeyId: string
  AccessKeySecret: string
  Bucket: string
  Domain: string
  Expiration: string
  SecurityToken: string
  HttpsDomain: string
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
export async function upload (e: { img:File , audio:File}, uploadDir: string, axios: AxiosSender) {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    try{
      const origin=window.location.origin
      // const { data } = await axios!.post<backData<UploadData>>(`${origin}/web/keti3/oss/auth`)
      const { data } = await getOssAuth()
      if (data.errcode !== 0) {
        throw new Error(data.errmsg)
      }
      const client = new OSS({
        region: EnumRegion.BeiJing,
        accessKeyId: data.data.AccessKeyId,
        accessKeySecret: data.data.AccessKeySecret,
        stsToken: data.data.SecurityToken,
        bucket: data.data.Bucket
      })
      const date= new Date().getTime();
      const imgKey = `img${useFileStore().uid}${date}`;
      await uploadFile(client, imgKey, e.img, 'image');
      const imgUrl = `${data.data.HttpsDomain}/${imgKey}`;

      // 上传音频
      const audioKey = `audio${useFileStore().uid}${date}`;
      await uploadFile(client, audioKey, e.audio, 'audio');
      const audioUrl = `${data.data.HttpsDomain}/${audioKey}`;

      return { imgUrl, audioUrl };
    }catch(err){
      console.log(err)
    }

}
async function uploadFile(client: OSS, key: string, file: File, type: 'image' | 'audio'): Promise<void> {
  const { parallel, partSize } = config;

  // 初始化分片上传
  const multipartUpload = client.multipartUpload(key, file, {
    parallel,
    partSize,
    progress: (p: number) => {
      if(type === 'image'){
        config.imgProgress = p*49
      }else{
        config.audioProgress = p*49
      }
      config.authProgress = config.imgProgress + config.audioProgress
      console.log('authProgress======',config.authProgress)
    }
  })
  //这里上传进度只到98%，等到这里上传ali-oss成功以后，再将返回的url以及日志等上传到后端以后进度改写为100%
  // 开始上传
  const result = await multipartUpload
}

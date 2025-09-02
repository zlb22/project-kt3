/**
 * 五边形
 */

import { fabric } from 'fabric'
import { useWorkspaceStore } from '@/stores/workspace'
import { fabricControl } from './ControlsConfig'
class PentagonProp {
  propData: any
  constructor(data: any, canvas: any) {
    this.propData = null
    this.createContext(data, canvas)
  }
  createContext(data: any, canvas: any) {
    const { left, top, src, width, height, type, name } = data
    const pentagonProp = (fabric.Image as any).fromURL(src, function (img: any) {
      // 图片原始宽度和高度
      const originalWidth = img.width
      const originalHeight = img.height

      // 计算缩放比例
      const scaleX = data.width / originalWidth
      const scaleY = data.height / originalHeight
      const scale = Math.min(scaleX, scaleY)

      img.scale(scale)
      const workspaceStore = useWorkspaceStore()
      workspaceStore.countProps(type)
      const addedProps = workspaceStore.addedProps
      const count = addedProps[type] || 0
      img.set({
        top,
        left,
        type,
        title: `${name}_${count}`,
        id: `${type}_${count}`
      })
      img.controls = fabricControl
      img.on('selected', function () {
        console.log('选中五边形')
        workspaceStore.selectedProp = img
        workspaceStore.showToolbar = true
      })
      canvas.add(img)
    })
  }
}

export default PentagonProp

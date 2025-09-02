/**
 * 矩形
 */

import { fabric } from 'fabric'
import { useWorkspaceStore } from '@/stores/workspace'
import { fabricControl } from './ControlsConfig'
//拓展IRectOptions属性
interface CustomRectOptions extends fabric.IRectOptions {
  id?: string
  title?: string
}

class RectProp {
  propData: any
  constructor(data: any) {
    this.propData = null
    this.createContext(data)
  }
  createContext(data: any) {
    const { top, left, width, height, fill, stroke, strokeWidth, type, name } = data
    const option: CustomRectOptions = {
      top,
      left,
      width,
      height,
      fill,
      stroke,
      strokeWidth,
      type,
      name,
      controls: fabricControl,
    }
    const workspaceStore = useWorkspaceStore()
    workspaceStore.countProps(type)
    const addedProps = workspaceStore.addedProps
    const count = addedProps[type] || 0
    const newOption = Object.assign({}, option, {
      id: `${type}_${count}`,
      title: `${name}_${count}`
    })
    const rect = new fabric.Rect(newOption)
    rect.on('selected', function () {
      console.log('选中矩形')
      workspaceStore.selectedProp = rect
      workspaceStore.showToolbar = true
    })
    this.propData = rect
  }
}

export default RectProp

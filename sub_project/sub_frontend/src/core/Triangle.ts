/**
 * 三角形
 */

import { fabric } from 'fabric'
import { useWorkspaceStore } from '@/stores/workspace'
import { fabricControl } from './ControlsConfig'
class TriangleProp {
  propData: any
  constructor(data: any) {
    this.propData = null
    this.createContext(data)
  }
  createContext(data: any) {
    const { top, left, width, height, fill, stroke, strokeWidth, type, name } = data
    const option = {
      top,
      left,
      width,
      height,
      fill,
      stroke,
      strokeWidth,
      type,
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
    const triangle = new fabric.Triangle(newOption)
    triangle.on('selected', function () {
      console.log('选中三角形')
      workspaceStore.selectedProp = triangle
      workspaceStore.showToolbar = true
    })
    this.propData = triangle
  }
}

export default TriangleProp

import { fabric } from 'fabric'
import { useWorkspaceStore } from '@/stores/workspace'
import { fabricControl } from './ControlsConfig'
interface CustomOptions extends fabric.ICircleOptions {
  id?: string
  title?: string
}
/**
 * 圆
 */
class CircleProp {
  propData: any
  constructor(data: any) {
    this.propData = null
    this.createContext(data)
  }

  createContext(data: any) {
    const { left, top, radius, fill, stroke, strokeWidth, type, name } = data
    const workspaceStore = useWorkspaceStore()
    workspaceStore.countProps(type)
    const addedProps = workspaceStore.addedProps
    const count = addedProps[type] || 0
    const option: CustomOptions = {
      left,
      top,
      radius,
      fill,
      stroke,
      strokeWidth,
      type,
      name,
      controls: fabricControl,
    }
    const newOption = Object.assign({}, option, {
      id: `${type}_${count}`,
      title: `${name}_${count}`
    })
    const circle = new fabric.Circle(newOption)
    circle.on('selected', function () {
      console.log('选中圆')
      workspaceStore.selectedProp = circle
      workspaceStore.showToolbar = true
    })
    this.propData = circle
  }
}

export default CircleProp

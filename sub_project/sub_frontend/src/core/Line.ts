/**
 * 直线
 */
import { fabric } from 'fabric'
import { useWorkspaceStore } from '@/stores/workspace'

class LineProp {
  propData: any
  constructor(data: any) {
    this.propData = null
    this.createContext(data)
  }

  createContext(data: any) {
    const { left: x1, top: y1, stroke, strokeWidth, type, name } = data

    const x2 = x1, y2 = y1 + 200;
    const workspaceStore = useWorkspaceStore()
    workspaceStore.countProps(type)
    const addedProps = workspaceStore.addedProps
    const count = addedProps[type] || 0
    const newOption = Object.assign({}, {
      id: `${type}_${count}`,
      title: `${name}_${count}`
    })

    const path = new fabric.Path(`M ${x1} ${y1} L ${x2} ${y2}`, {
      stroke,
      strokeWidth: strokeWidth || 2,
      hasBorders: true, // 隐藏边框
      selectable: true, // 确保对象可选中
      strokeUniform: true, // 确保线条宽度不会随着缩放而改变
      ...newOption,
      controls: {
        mtr: new fabric.Control({
          x: 0,
          y: 0,
          offsetX: 50,
          offsetY: 0,
          actionHandler: fabric.controlsUtils.rotationWithSnapping,
          cursorStyleHandler: () =>
            'url(https://static0.xesimg.com/xpp-parent-fe/topic-three-web/refresh-icon.png), auto',
          actionName: 'rotate',
          render: fabric.controlsUtils.renderCircleControl,
          cornerSize: 24
        }),
        br: fabric.Object.prototype.controls.br,
        tl:  fabric.Object.prototype.controls.tl,
      }
    })

    // 获取路径的宽度和高度
    const pathWidth = path.getScaledWidth()
    const pathHeight = path.getScaledHeight()

    console.log('pathWidth===>', pathWidth, pathHeight)

    path.set({
      angle: -45
    })

    path.on('selected', function () {
      console.log('选中直线')
      workspaceStore.selectedProp = path
      workspaceStore.showToolbar = true
    })
    this.propData = path
  }
}

export default LineProp

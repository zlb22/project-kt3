import { fabric } from 'fabric'

export const customRotationControl = new fabric.Control({
  ...fabric.Object.prototype.controls.mtr,
  actionHandler: fabric.controlsUtils.rotationWithSnapping,
  cursorStyleHandler: () =>
    'url(https://static0.xesimg.com/xpp-parent-fe/topic-three-web/refresh-icon.png), auto',
  actionName: 'rotate',
  render: fabric.controlsUtils.renderCircleControl,
  cornerSize: 24
})

export const fabricControl = {
  ...fabric.Object.prototype.controls, // 保持其他控制点的默认功能
  mtr: customRotationControl
}

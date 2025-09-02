/**
 * 曲线
 */
import { fabric } from 'fabric'
import { useWorkspaceStore } from '@/stores/workspace'
interface CustomPathOptions extends fabric.IPathOptions {
  type?: string
  id?: string
  title?: string
}
let isInitCurve = false
class CurveProp {
  propData: any
  isInit: any
  curveArr: string[]
  pcPoint: any
  plPoint: any
  prPoint: any
  canvas: any
  config: any
  nowCurve: any
  id: string
  constructor(data: any, canvas: any) {
    this.config = data
    this.id = ''
    this.propData = null
    this.canvas = canvas
    this.curveArr = ['curve-pl', 'curve-pc', 'curve-pr']
    // 曲线大小位置调整
    this.pcPoint = {
      left: -32 + this.config.left,
      top: 35 + this.config.top
    }
    this.plPoint = {
      left: 41 + this.config.left,
      top: 65 + this.config.top - 5
    }
    this.prPoint = {
      left: 41 + this.config.left,
      top: 5 + this.config.top - 5
    }

    this.initEvent()

    this.createContext()
  }
  // 两端的点
  makeCurveCircle = (left: number, top: number, line1: any, line2: any, line3: any) => {
    const c: any = new fabric.Rect({
      left: left,
      top: top,
      width: 10,
      height: 10,
      fill: '#fff',
      stroke: '#000',
      originX: 'center',
      originY: 'center'
    })

    c.hasBorders = c.hasControls = false
    c.line1 = line1
    c.line2 = line2
    c.line3 = line3

    return c
  }
  // 中间的点
  makeCurvePoint = (left: number, top: number, line1: any, line2: any, line3: any) => {
    const c: any = new fabric.Circle({
      left: left,
      top: top,
      strokeWidth: 0.5,
      originX: 'center',
      originY: 'center',
      radius: 5,
      fill: 'rgb(145,210,89)',
      stroke: 'rgb(119,168,76)'
    })

    c.hasBorders = c.hasControls = false

    c.line1 = line1
    c.line2 = line2
    c.line3 = line3

    return c
  }
  // 绘制线
  makeCurveLine = (myPath: string, pl: any, pc: any, pr: any) => {
    const curveLine: any = new fabric.Path(myPath, {
      fill: '',
      stroke: this.config.stroke,
      strokeWidth: this.config.strokeWidth,
      objectCaching: false
    })
    curveLine.name = 'curve-line'

    curveLine.path[0][1] = pl.left
    curveLine.path[0][2] = pl.top

    curveLine.path[1][1] = pc.left
    curveLine.path[1][2] = pc.top

    curveLine.path[1][3] = pr.left
    curveLine.path[1][4] = pr.top
    // curveLine.hasBorders = curveLine.hasControls = false
    const workspaceStore = useWorkspaceStore()
    curveLine.on('selected', function () {
      console.log('选中直线')
      workspaceStore.selectedProp = curveLine
      workspaceStore.showToolbar = true
    })
    return curveLine
  }
  // 重绘线
  reLine = (line: any, pl: any, pc: any, pr: any) => {
    const p = [pl, pc, pr]
    this.canvas.remove(line)
    const curveLine = this.makeCurveLine(line.path, pl, pc, pr)
    curveLine.curveObj = { pl, pr, curveLine, pc }
    pl.scaleX = pl.scaleY = 1
    pc.scaleX = pc.scaleY = 1
    pr.scaleX = pr.scaleY = 1
    for (let i = 0; i < 3; i++) {
      if (p[i].line1 && p[i].line1.name === 'curve-line') {
        p[i].line1 = curveLine
      }
      if (p[i].line2 && p[i].line2.name === 'curve-line') {
        p[i].line2 = curveLine
      }
      if (p[i].line3 && p[i].line3.name === 'curve-line') {
        p[i].line3 = curveLine
      }
    }
    curveLine.idObj = line.idObj
    this.cacheCurve(curveLine.idObj.id,curveLine)
    this.canvas.add(curveLine)
  }
  // 隐藏所有曲线元素的操作点
  hidenAllPoint = () => {
    this.canvas.getObjects().forEach((item: any) => {
      if (item.name && this.curveArr.includes(item.name)) {
        this.hidePoint(item.curveObj.pc, item.curveObj.pl, item.curveObj.pr, true)
      }
    })
  }
  // 禁用所有曲线元素的操作点
  forbiddenAllPoint = (activeObject: any) => {
    activeObject.setControlVisible('ml', false)
    activeObject.setControlVisible('mb', false)
    activeObject.setControlVisible('mr', false)
    activeObject.setControlVisible('mt', false)
    activeObject.setControlVisible('tl', false)
    activeObject.setControlVisible('bl', false)
    activeObject.setControlVisible('tr', false)
    activeObject.setControlVisible('br', false)
    activeObject.setControlVisible('mtr', false)
  }
  // 鼠标按下
  onObjectMouseDown = (e: any) => {
    const activeObject = e.target
    if (!(activeObject && activeObject.name && this.curveArr.includes(activeObject.name))) {
      this.hidenAllPoint()
    }
    // 如果是曲线则不支持缩放和抻拉
    if (activeObject?.name === 'curve-line') {
      this.forbiddenAllPoint(activeObject)
    }
  }
  // 鼠标抬起
  onObjectMouseUp = (e: any) => {
    const activeObject = e.target
    // 如果点击的是曲线内部元素才触发
    if(activeObject?.name === 'curve-line'&&activeObject?.idObj){
      this.cacheCurve(activeObject.idObj.id,activeObject)
      if(activeObject.moving){
        this.canvas.fire('object:modified', { target: Object.assign(activeObject.idObj,{curveLine:{left:activeObject.left,top:activeObject.top}}) })
      }
      activeObject.moving = false
    }
    if (this.canvas.getActiveObject()?.name === 'curve-line') return
    if (activeObject && activeObject.name && this.curveArr.includes(activeObject.name)) {
      const pl = activeObject.curveObj.pl
      const pc = activeObject.curveObj.pc
      const pr = activeObject.curveObj.pr
      const curveLine = pl.line1
      this.reLine(curveLine, pl, pc, pr)
      this.hidePoint(pc, pl, pr, false)
      this.canvas.fire('object:curveSkew', Object.assign(activeObject.idObj,{curveLine:{left:curveLine.left,top:curveLine.top}}))
    }
  }
  // 移动
  onObjectMoving = (e: any) => {
    // 拖拽曲线移动
    if (e.target.name === 'curve-line') {
      const p = e.target
      const pl = p.curveObj.pl
      const pc = p.curveObj.pc
      const pr = p.curveObj.pr
      pl.left = pl.left + e.e.movementX
      pl.top = pl.top + e.e.movementY

      pc.left = pc.left + e.e.movementX
      pc.top = pc.top + e.e.movementY

      pr.left = pr.left + e.e.movementX
      pr.top = pr.top + e.e.movementY
      pc.setCoords()
      pl.setCoords()
      pr.setCoords()
      e.target.moving = true
    }
    // 拖拽操作点移动
    if (this.curveArr.includes(e.target.name)) {
      const p = e.target
      if (p.line1) {
        p.line1.path[0][1] = p.left
        p.line1.path[0][2] = p.top
      } else if (p.line3) {
        p.line3.path[1][3] = p.left
        p.line3.path[1][4] = p.top
      } else if (p.line2) {
        p.line2.path[1][1] = p.left
        p.line2.path[1][2] = p.top
      }
    }
  }
  // 双击
  onObjectDblclick = (e: any) => {
    this.canvas.discardActiveObject()
    const activeObject = e.target
    if (activeObject && activeObject.name === 'curve-line') {
      const pl = activeObject.curveObj.pl
      const pc = activeObject.curveObj.pc
      const pr = activeObject.curveObj.pr
      const myPoint = [pl, pc, pr]
      myPoint.forEach((item) => {
        if (item.line1) {
          item.line1.path[0][1] = item.left
          item.line1.path[0][2] = item.top
        } else if (item.line3) {
          item.line3.path[1][3] = item.left
          item.line3.path[1][4] = item.top
        } else if (item.line2) {
          item.line2.path[1][1] = item.left
          item.line2.path[1][2] = item.top
        }
      })
      const curveLine = pl.line1
      this.reLine(curveLine, pl, pc, pr)
      this.hidePoint(pc, pl, pr, false)
    }
  }
  // 隐藏展示操作点
  hidePoint = (pc: any, pl: any, pr: any, flag: boolean) => {
    pc.set('visible', !flag)
    pl.set('visible', !flag)
    pr.set('visible', !flag)
    this.canvas.bringToFront(pc)
    this.canvas.bringToFront(pl)
    this.canvas.bringToFront(pr)
    this.canvas.requestRenderAll()
  }
  // 创建id
  createId = () => {
    const workspaceStore = useWorkspaceStore()
    workspaceStore.countProps(this.config.type)
    const addedProps = workspaceStore.addedProps
    const count = addedProps[this.config.type] || 0
    return {
      type: this.config.type,
      title: `${this.config.name}_${count}`,
      id: `${this.config.type}_${count}`
    }
  }
  // 缓存全局元素曲线
  cacheCurve = (id:string,curveLine: any) => {
    console.log(id,curveLine)
    if((window as any).myCurveLine){
      (window as any).myCurveLine[id] = curveLine
    }else{
      (window as any).myCurveLine = {};
      (window as any).myCurveLine[id] = curveLine
    }
  }
  // 获取当前曲线
  getNowCurve = (id:string) => {
    if((window as any).myCurveLine){
      return (window as any).myCurveLine[id]
    }
  }
  // 创建曲线元素
  createCurveLine = () => {
    const idObj = this.createId()
    // 曲线
    const curveLine = this.makeCurveLine(
      'M ' +
        this.plPoint.left +
        ' ' +
        this.plPoint.top +
        ' Q ' +
        this.pcPoint.left +
        ', ' +
        this.pcPoint.top +
        ', ' +
        this.prPoint.left +
        ', ' +
        this.prPoint.top +
        '',
      this.plPoint,
      this.pcPoint,
      this.prPoint
    )
    curveLine.name = 'curve-line'
    curveLine.idObj = idObj
    this.canvas.add(curveLine)
    // 中间点
    const pc = this.makeCurvePoint(this.pcPoint.left, this.pcPoint.top, null, curveLine, null)
    pc.name = 'curve-pc'
    pc.idObj = idObj
    this.canvas.add(pc)
    // 左点
    const pl = this.makeCurveCircle(this.plPoint.left, this.plPoint.top, curveLine, pc, null)
    pl.name = 'curve-pl'
    pl.idObj = idObj
    this.canvas.add(pl)
    // 右点
    const pr = this.makeCurveCircle(this.prPoint.left, this.prPoint.top, null, pc, curveLine)
    pr.name = 'curve-pr'
    pr.idObj = idObj
    this.canvas.add(pr)

    pc.curveObj = { pl, pr, curveLine, pc }
    pl.curveObj = { pl, pr, curveLine, pc }
    pr.curveObj = { pl, pr, curveLine, pc }
    curveLine.curveObj = { pl, pr, curveLine, pc }
    this.hidePoint(pc, pl, pr, true);
    
    this.cacheCurve(idObj.id,curveLine)
    this.canvas.fire('object:added', { target: Object.assign(this, idObj, {curveLine:{left:curveLine.left,top:curveLine.top}}) })
  }
  // 初始化事件
  initEvent = () => {
    if (!isInitCurve) {
      this.canvas.on('mouse:down', (opt: any) => {
        console.log('mouse:down', opt)
        this.onObjectMouseDown(opt)
      })
      this.canvas.on('mouse:up', (opt: any) => {
        console.log('mouse:up', opt)
        this.onObjectMouseUp(opt)
      })
      this.canvas.on('object:moving', (opt: any) => {
        // console.log('object:moving',opt)
        this.onObjectMoving(opt)
      })
      this.canvas.on('mouse:dblclick', (opt: any) => {
        console.log('mouse:dblclick', opt)
        this.onObjectDblclick(opt)
      })
      this.canvas.on('selection:created', (opt: any) => {
        console.log('selection:created', opt)
        // 如果是曲线则不支持缩放和抻拉
        if (opt?.selected[0]?.name === 'curve-line') {
          this.forbiddenAllPoint(opt.selected[0])
        }
      })
      isInitCurve = true
    }
  }
  createContext() {
    this.createCurveLine()
  }
  redrawLine(nowObj:any){
    const nowCurve = this.getNowCurve(this.id)
    console.log('redrawLine',nowObj)
    const leftCha = nowObj.left - nowCurve.left
    const topCha =  nowObj.top - nowCurve.top
    nowCurve.left = nowObj.left
    nowCurve.top = nowObj.top
    nowCurve.setCoords()
    nowCurve.curveObj.curveLine.setCoords()
    nowCurve.curveObj.pl.left = nowCurve.curveObj.pl.left + leftCha
    nowCurve.curveObj.pl.top = nowCurve.curveObj.pl.top + topCha
    nowCurve.curveObj.pc.left = nowCurve.curveObj.pc.left + leftCha
    nowCurve.curveObj.pc.top = nowCurve.curveObj.pc.top + topCha
    nowCurve.curveObj.pr.left = nowCurve.curveObj.pr.left + leftCha
    nowCurve.curveObj.pr.top = nowCurve.curveObj.pr.top + topCha
    nowCurve.curveObj.pl.setCoords()
    nowCurve.curveObj.pc.setCoords()
    nowCurve.curveObj.pr.setCoords()
    this.hidePoint(nowCurve.curveObj.pc, nowCurve.curveObj.pl, nowCurve.curveObj.pr, true)
  }
  visibleRawLine(flag:boolean){
    const nowCurve = this.getNowCurve(this.id)
    nowCurve.set('visible', !flag)
  }
  toObject(){
    return {}
  }
}

export default CurveProp

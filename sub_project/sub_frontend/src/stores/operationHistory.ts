import { ref } from 'vue'
import { defineStore } from 'pinia'

import { useWorkspaceStore } from './workspace'

enum OperationLogType {
  ADD = '新增对象',
  MOVE = '移动对象',
  DELETE = '删除对象',
  SCALE = '调整缩放比例',
  ROTATE = '调整旋转角度',
  ASPECT_RATIO = '调整宽高比例',
  FLIPX = '水平翻转',
  FLIPY = '垂直翻转',
  CURVESKEW = '扭曲对象',
  UNDO = 'Undo', // 撤销
  REDO = 'Redo' // 前进
}

function isImageObject(target: any) {
  // 检查对象是否是 fabric.Image 类型
  if (target.type === 'image') {
    return true
  }

  // 检查对象是否包含图像数据
  if (target.getSrc && typeof target.getSrc === 'function') {
    const src = target.getSrc()
    if (src && src.length > 0) {
      return true
    }
  }

  return false
}

// 判断道具是否为store 列表中的
function inPropsList(type: any) {
  // 获取 workspace store
  const workspaceStore = useWorkspaceStore()

  const list = workspaceStore.propsList
  return list.some((item) => item.type === type)
}

// id 作为key，将target 存放到数组中 记录变换前后位置
const OperationLogObjMap: Record<string, any[]> = {}

export const useOperationHistory = defineStore('operationHistory', () => {
  // 历史记录
  const historyStack: any = ref([] as any[])

  // 点击返回按钮记录临时存放
  const historyStackTemp = ref([] as any[])

  // 操作日志
  const operationLog = ref([] as any[])
  // 操作页面停留时间 单位s  最大 8 * 60 s
  const stayOperaDur = ref(0)
  const maxStayOperaDur = ref(8 * 60)

  const canvasObj: any = ref(null)

  // 控件被操作
  // 点击前进需要将 historyStackTemp 一次次重新存放到 historyStack
  const push = (record: any) => {
    // 获取对象的层级映射
    const getObjectLayerMapping = () => {
      if (!canvasObj.value) return {}

      return canvasObj.value.getObjects().reduce(
        (acc: any, obj: any, index: any) => {
          if (obj.id) {
            acc[obj.id] = index
          }
          if (obj.idObj && obj.idObj.id && obj.name === 'curve-line') {
            acc[obj.idObj.id] = index
          }
          return acc
        },
        {} as Record<string, number>
      )
    }

    const layerMap = getObjectLayerMapping()

    console.log('layerMap===>', layerMap)

    historyStack.value.push({
      ...record,
      target: Object.assign(record.target, { id: record.target.id || record.target.idObj.id }),
      layerMap
    })
  }

  // 点击后退需要将 historyStack 尾部的移出到 historyStackTemp
  const pop = () => {
    if (historyStack.value.length <= 0) return

    const lastRecord = historyStack.value.pop()
    const prevRecord = historyStack.value.findLast(
      (item: any) => item.target === lastRecord.target || item.target.id === lastRecord.target.id
    )

    // console.log('pop===>', lastRecord, prevRecord)

    switch (lastRecord.action) {
      case 'deleted':
        if (lastRecord.target?.idObj && lastRecord.target?.idObj?.type === 'Curve') {
          lastRecord.target.set({
            selectable: true,
            visible: true
          })
        } else {
          prevRecord.target.set({
            selectable: true,
            visible: true
          })
        }
        break
      case 'flipX':
      case 'flipY':
      case 'modified':
        if (prevRecord.target.type === 'Curve') {
          console.log('prevRecord.state===>', prevRecord.state)
          prevRecord.target.redrawLine(prevRecord.state)
        } else {
          prevRecord.target.set(prevRecord.state)
          prevRecord.target.setCoords()
        }

        break

      case 'bringToFront':
      case 'sendToBack':
      case 'bringForward':
      case 'sendBackwards':
        console.log('lastRecord.target===>', lastRecord.target)
        lastRecord.target.moveTo(
          historyStack.value[historyStack.value.length - 1].layerMap[lastRecord.target.id]
        )
        break
      case 'added':
        // 隐藏并移除选择框
        if (lastRecord.target.type === 'Curve') {
          lastRecord.target.visibleRawLine(true)
        } else {
          lastRecord.target.set({
            selectable: false,
            visible: false
          })
          canvasObj.value?.discardActiveObject()
        }

        break
    }

    canvasObj.value?.renderAll() // 重新渲染画布
    historyStackTemp.value.push(lastRecord)
  }

  const bindCanvasEvent = (canvas: any) => {
    canvasObj.value = canvas
    // 监听对象的修改完成事件
    canvas.on('object:modified', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return

      if (target.type === 'Curve') {
        // 根据target.id 找正确可以调用到的 target
        const realTar = historyStack.value.find(
          (item: any) => item.target.id === target.id && item.action === 'added'
        )

        push({
          target: realTar.target,
          action: 'modified',
          timestamp: Date.now(),
          state: target.curveLine
        })

        addOperationLogObjMap({ ...realTar.target, ...target.curveLine }, 'modified')
      } else {
        push({
          target,
          action: 'modified',
          timestamp: Date.now(),
          state: target.toObject()
        })

        addOperationLogObjMap(target, 'modified')
      }
    })

    canvas.on('object:flipX', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return

      push({
        target,
        action: 'flipX',
        timestamp: Date.now(),
        state: target.toObject()
      })

      addOperationLogObjMap(target, 'flipX')
    })

    canvas.on('object:flipY', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return

      push({
        target,
        action: 'flipY',
        timestamp: Date.now(),
        state: target.toObject()
      })

      addOperationLogObjMap(target, 'flipY')
    })

    canvas.on('object:bringToFront', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if (
        (target.id && !inPropsList(target.id.split('_')[0])) ||
        (target.idObj?.id && !inPropsList(target.idObj.id.split('_')[0]))
      )
        return

      push({
        target,
        action: 'bringToFront',
        timestamp: Date.now(),
        state: target.toObject()
      })

      // addOperationLogObjMap(target, 'bringToFront')
    })
    canvas.on('object:sendToBack', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if (
        (target.id && !inPropsList(target.id.split('_')[0])) ||
        (target.idObj?.id && !inPropsList(target.idObj.id.split('_')[0]))
      )
        return

      push({
        target,
        action: 'sendToBack',
        timestamp: Date.now(),
        state: target.toObject()
      })

      // addOperationLogObjMap(target, 'sendToBack')
    })

    canvas.on('object:curveSkew', (e: any) => {
      console.log('11111111', e)
      addOperationLogObjMap({ ...e, ...e.curveLine, action: 'curveSkew' }, 'curveSkew')
    })

    canvas.on('object:bringForward', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if (
        (target.id && !inPropsList(target.id.split('_')[0])) ||
        (target.idObj?.id && !inPropsList(target.idObj.id.split('_')[0]))
      )
        return

      push({
        target,
        action: 'bringForward',
        timestamp: Date.now(),
        state: target.toObject()
      })

      // addOperationLogObjMap(target, 'bringForward')
    })
    canvas.on('object:sendBackwards', (e: any) => {
      const target = e.target
      console.log('target===>', target)

      if (
        (target.id && !inPropsList(target.id.split('_')[0])) ||
        (target.idObj?.id && !inPropsList(target.idObj.id.split('_')[0]))
      )
        return

      push({
        target,
        action: 'sendBackwards',
        timestamp: Date.now(),
        state: target.toObject()
      })

      // addOperationLogObjMap(target, 'sendBackwards')
    })

    // 监听对象的添加事件
    canvas.on('object:added', (e: any) => {
      const target = e.target

      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return

      console.log('object:added===>', e.target)

      if (target.type === 'Curve') {
        push({
          target,
          action: 'added',
          timestamp: Date.now(),
          state: target.curveLine
        })

        addOperationLogObjMap({ ...target, ...target.curveLine }, 'added')
      } else {
        push({
          target,
          action: 'added',
          timestamp: Date.now(),
          state: target.toObject()
        })

        addOperationLogObjMap(target, 'added')
      }
    })

    canvas.on('object:deleted', (e: any) => {
      const target = e.target

      if (
        (target.id && !inPropsList(target.id.split('_')[0])) ||
        (target?.idObj?.id && !inPropsList(target.idObj.id.split('_')[0]))
      )
        return
      console.log('object:deleted====', target)

      if (target.idObj?.type === 'Curve') {
        const realTar = historyStack.value.find((item: any) => item.target.id === target.idObj.id)
        push({
          target,
          action: 'deleted',
          timestamp: Date.now(),
          state: target.idObj.curveLine
        })

        addOperationLogObjMap({ ...realTar.target, ...target.curveLine }, 'deleted')
      } else {
        push({
          target,
          action: 'deleted',
          timestamp: Date.now(),
          state: target.toObject()
        })

        addOperationLogObjMap(target, 'deleted')
      }
    })

    // 监听对象的修改开始事件
    canvas.on('object:scaling', (e: any) => {
      const target = e.target

      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return

      addOperationLogObjMap(e.target, 'scaling')
    })

    canvas.on('object:moving', (e: any) => {
      const target = e.target
      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return

      addOperationLogObjMap(e.target, 'moving')
    })

    canvas.on('object:rotating', (e: any) => {
      const target = e.target
      if ((target.id && !inPropsList(target.id.split('_')[0])) || !target.id) return
      addOperationLogObjMap(e.target, 'rotating')
    })
  }

  const backwardHanle = () => {
    if (historyStack.value.length > 0) {
      canvasObj.value?.discardActiveObject()
      pop()

      pushOperationLog(null, null, OperationLogType.UNDO)
    }
  }
  const forwardHandle = () => {
    if (historyStackTemp.value.length > 0) {
      canvasObj.value?.discardActiveObject()

      const record = historyStackTemp.value.pop()

      if (record.target.type === 'Curve' || record.target.idObj?.type === 'Curve') {
        if (record.action === 'added') {
          record.target.visibleRawLine(false)
        }

        if (record.action === 'deleted') {
          record.target.set({
            selectable: false,
            visible: false
          })
        }

        if (record.action === 'modified') {
          record.target.redrawLine(record.state)
        }

        if (
          ['bringToFront', 'sendToBack', 'bringForward', 'sendBackwards'].indexOf(record.action) >
          -1
        ) {
          record.target.moveTo(record.layerMap[record.target.id])
        }
      } else {
        record.target.set(record.state)
        record.target.setCoords()

        if (record.action === 'added') {
          record.target.set({
            selectable: true,
            visible: true
          })
        } else if (
          ['bringToFront', 'sendToBack', 'bringForward', 'sendBackwards'].indexOf(record.action) >
          -1
        ) {
          record.target.moveTo(record.layerMap[record.target.id])
        }
      }

      push(record)
      canvasObj.value?.renderAll() // 重新渲染画布

      pushOperationLog(null, null, OperationLogType.REDO)
    }
  }

  const updateStayOperaDuration = () => {
    stayOperaDur.value++
  }

  // 映射后端接口字段用来提交
  const pushOperationLog = (target: any, preTarget: any, type: OperationLogType) => {
    console.log('pushOperationLog===>', target, preTarget, type)

    const log = target
      ? {
          op_time: stayOperaDur.value,
          op_type: (type === OperationLogType.FLIPX || type === OperationLogType.FLIPY) ? '翻转对象' : type,
          op_object: target?.id?.split('_')[0],
          object_name: target?.title?.split('_')[0],
          object_no: Number(target?.id?.split('_')[1]),
          data_before: JSON.stringify({}),
          data_after: JSON.stringify({
            flipX: target?.flipX ? '水平翻转' : '',
            flipY: target?.flipY ? '垂直翻转' : '',
            skewX: target?.skewX || '',
            skewY: target?.skewY || '',

            width: target?.width * target.scaleX ? target?.width * target.scaleX + '' : '',
            hegiht: target?.height * target.scaleY ? target?.height * target.scaleY + '' : '',

            scaleX:
              target?.scaleX * (isImageObject(target) ? 4 : 1)
                ? target?.scaleX * (isImageObject(target) ? 4 : 1) + ''
                : '', // 使用图片的道具 缩放比例需要 * 4
            scaleY:
              target?.scaleY * (isImageObject(target) ? 4 : 1)
                ? target?.scaleY * (isImageObject(target) ? 4 : 1) + ''
                : '',

            // 操作的起点坐标
            start_coords: {
              value: preTarget ? `${preTarget?.left},${preTarget?.top}` : ''
            },

            // 操作的终点坐标
            end_coords: {
              value: `${target?.left},${target?.top}`
            },
            // 缩放
            scale: {
              value: target?.scaleX
                ? target?.scaleX * (isImageObject(target) ? 4 : 1) + '' || ''
                : ''
            },
            // 宽高
            aspect: {
              value:
                (target?.width * target.scaleX) / (target?.height * target.scaleY)
                  ? (target?.width * target.scaleX) / (target?.height * target.scaleY) + ''
                  : ''
            },
            // 旋转
            rotate: {
              value: target?.angle ? target?.angle + '' : ''
            },
            // 扭曲
            distortion: {
              value: type === OperationLogType.CURVESKEW ? '扭曲' : ''
            },
            flip: {
              value: type === OperationLogType.FLIPX ? '水平翻转' : type === OperationLogType.FLIPY ? '垂直翻转' : ''
            }
          })
        }
      : {
          op_type: type,
          op_time: stayOperaDur.value
        }

    operationLog.value.push(log)
  }

  // 新增记录map
  const addOperationLogObjMap = (target: any, type: string) => {
    console.log('addOperationLogObjMap===>', target)

    const newLog = {
      action: type,
      target: {
        left: target.left,
        top: target.top,
        height: target.height,
        width: target.width,
        scaleX: target.scaleX,
        scaleY: target.scaleY,
        angle: target.angle,
        flipX: target.flipX,
        flipY: target.flipY,
        visible: target.visible,
        selectable: target.selectable,
        skewX: target.skewX,
        skewY: target.skewY,
        cropX: target.cropX,
        cropY: target.cropY
      }
    }

    if (!OperationLogObjMap[target.id]) OperationLogObjMap[target.id] = []

    // 需要先判断当前的type 若为 modified, 则需要将 OperationLogObjMap[target.id] 净化
    if (type === 'modified' && OperationLogObjMap[target.id].length > 0) {
      // 将 scaling moving rotating 触发的结果 合一
      let lastOperationLog = OperationLogObjMap[target.id][OperationLogObjMap[target.id].length - 1]
      while (
        lastOperationLog.action === 'scaling' ||
        lastOperationLog.action === 'moving' ||
        lastOperationLog.action === 'rotating'
      ) {
        OperationLogObjMap[target.id].pop()
        lastOperationLog = OperationLogObjMap[target.id][OperationLogObjMap[target.id].length - 1]

        if (lastOperationLog.action === 'scaling') {
          newLog.action = 'scaled'
        }

        if (lastOperationLog.action === 'moving') {
          newLog.action = 'moved'
        }

        if (lastOperationLog.action === 'rotating') {
          newLog.action = 'rotated'
        }
      }
    }

    if (type === 'scaled' && OperationLogObjMap[target.id].length > 0) {
      const lastLogTarget =
        OperationLogObjMap[target.id][OperationLogObjMap[target.id].length - 1].target

      if (
        (lastLogTarget.scaleX !== newLog.target.scaleX &&
          lastLogTarget.scaleY === newLog.target.scaleY) ||
        (lastLogTarget.scaleX === newLog.target.scaleX &&
          lastLogTarget.scaleY !== newLog.target.scaleY)
      ) {
        newLog.action = 'aspectRatio'
      } else if (
        lastLogTarget.scaleX === newLog.target.scaleX &&
        lastLogTarget.scaleY === newLog.target.scaleY
      ) {
        console.log('宽高比没变动 ===> ')
      }
    }

    if (type === 'curveSkew') {
      newLog.action = 'curveSkew'
    }

    if (target.type === 'Curve' && type === 'modified') {
      newLog.action = 'moved'
    }

    console.log('addOperationLogObjMap=222==>', newLog)

    const preTarget =
      OperationLogObjMap[target.id].length > 0
        ? OperationLogObjMap[target.id][OperationLogObjMap[target.id].length - 1].target
        : null

    OperationLogObjMap[target.id].push(newLog)

    switch (newLog.action) {
      case 'added':
        pushOperationLog(target, preTarget, OperationLogType.ADD)
        break
      case 'scaled':
        pushOperationLog(target, preTarget, OperationLogType.SCALE)
        break
      case 'aspectRatio':
        pushOperationLog(target, preTarget, OperationLogType.ASPECT_RATIO)
        break
      case 'moved':
        pushOperationLog(target, preTarget, OperationLogType.MOVE)
        break
      case 'rotated':
        pushOperationLog(target, preTarget, OperationLogType.ROTATE)
        break
      case 'deleted':
        pushOperationLog(target, preTarget, OperationLogType.DELETE)
        break
      case 'flipX':
        pushOperationLog(target, preTarget, OperationLogType.FLIPX)
        break
      case 'flipY':
        pushOperationLog(target, preTarget, OperationLogType.FLIPY)
        break
      case 'curveSkew':
        pushOperationLog(target, preTarget, OperationLogType.CURVESKEW)
        break

      default:
        break
    }
  }

  return {
    canvasObj,
    historyStack,
    historyStackTemp,
    stayOperaDur,
    bindCanvasEvent,
    backwardHanle,
    forwardHandle,
    updateStayOperaDuration,
    operationLog
  }
})

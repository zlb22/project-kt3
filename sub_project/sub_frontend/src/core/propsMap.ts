import RectProp from './Rect'
import LineProp from './Line'
import CircleProp from './Circle'
import TriangleProp from './Triangle'
import PentagonProp from './Pentagon'
import SemiCircleProp from './Semicircle'
import HookProp from './Hook'
import ConeProp from './Cone'
import TriangularPyramidProp from './TriangularPyramid'
import CubeProp from './Cube'
import CylinderProp from './Cylinder'
import CurveProp from './Curve'

const propsMap: any = {
  Rect: RectProp,
  Line: LineProp,
  Circle: CircleProp,
  Triangle: TriangleProp,
  Pentagon: PentagonProp,
  Semicircle: SemiCircleProp,
  Hook: HookProp,
  Cone: ConeProp,
  TriangularPyramid: TriangularPyramidProp,
  Cube: CubeProp,
  Cylinder: CylinderProp,
  Curve: CurveProp
}

export default propsMap

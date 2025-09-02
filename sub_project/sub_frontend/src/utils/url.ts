//获取参数用于调试
export const getUrlParam = (key: string) => {
    const search = location.href.split("?")[1] || "";
    const paramsAry = search.split("&");
    for(const item of paramsAry) {
      if (item !== '' && key === item.split("=")[0]) {
        return item.split("=")[1];
      }
    }
    return null
}